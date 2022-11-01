package matchmaking

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/chat"
	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/discord"
	"github.com/haashi/omega-strikers-bot/internal/models"
	"github.com/haashi/omega-strikers-bot/internal/scheduled"
	log "github.com/sirupsen/logrus"
)

func GetMatchByThreadId(threadID string) (*models.Match, error) {
	match, err := db.GetMatchByThreadID(threadID)
	if err != nil {
		return nil, err
	}
	return match, nil
}

func createNewMatch(team1 []*models.Player, team2 []*models.Player) error {
	matchId := rand.Intn(math.MaxInt32)
	channelId := discord.MatchesChannel.ID
	session := discord.GetSession()
	match := &models.Match{}
	match.ID = fmt.Sprintf("%d", matchId)
	match.Timestamp = int(time.Now().Unix())
	log.Infof("creating new match %s", match.ID)
	mentionMessage := ""
	for i := range team1 {
		mentionMessage += "<@" + team1[i].DiscordID + ">"
	}
	mentionMessage += " VS "
	for i := range team2 {
		mentionMessage += "<@" + team2[i].DiscordID + ">"
	}
	initialMessage, err := session.ChannelMessageSendComplex(channelId, &discordgo.MessageSend{
		Content: fmt.Sprintf("ID:%d | %s", matchId, mentionMessage),
	})
	if err != nil {
		return err
	}
	match.MessageID = initialMessage.ID
	thread, err := session.MessageThreadStartComplex(initialMessage.ChannelID, initialMessage.ID, &discordgo.ThreadStart{
		Name:                fmt.Sprintf("%d", matchId),
		AutoArchiveDuration: 1440,
		Invitable:           false,
	})
	if err != nil {
		return err
	}
	match.ThreadID = thread.ID
	_, err = session.ChannelMessageSendComplex(thread.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Lobby code : %d\nFirst user in both team was assigned goalie, you are free to discuss here to change roles.\nUse this thread to share in game name to make teams, or to chat.\nPlease report match result with ``/result`` in this thread or cancel the match with ``/cancel`` if people are missing.", matchId),
	})
	if err != nil {
		return err
	}
	match.Team1 = team1
	match.Team2 = team2
	return db.CreateMatch(match)
}

func VoteCancelMatch(match *models.Match) {
	var content string = "A cancel request has been sent for this match.\nPlease react to this message to confirm."
	log.Debugf("getting players confirmation of cancellation of match %s", match.ID)
	message, err := chat.CreateVoteMessage(match.ThreadID, content, []string{"✅", "❌"})
	if err != nil {
		log.Errorf("failed to create confirmation message: " + err.Error())
		return
	}
	match.State = models.MatchStateVoteInProgress
	match.Team1Score = 0
	match.Team2Score = 0
	match.VoteMessageID = message.ID
	err = db.UpdateMatch(match)
	if err != nil {
		log.Errorf("failed to update reaction: " + err.Error())
		return
	}
	scheduled.TaskManager.Add(scheduled.Task{ID: "matchvote" + match.ID, Frequency: time.Second, Run: func() { handleMatchVoteResult(match) }})
}

func handleMatchVoteResult(match *models.Match) {
	s := discord.GetSession()
	voteMessage, err := s.ChannelMessage(match.ThreadID, match.VoteMessageID)
	if err != nil {
		log.Errorf("failed to get vote message: " + err.Error())
		return
	}
	allowedVoter := make([]string, 0)
	for _, p := range match.Team1 {
		allowedVoter = append(allowedVoter, p.DiscordID)
	}
	for _, p := range match.Team2 {
		allowedVoter = append(allowedVoter, p.DiscordID)
	}
	reactions, err := chat.FetchVoteResults(voteMessage, []string{"✅", "❌"}, allowedVoter)
	if err != nil {
		log.Errorf("failed to fetch votes: " + err.Error())
		return
	}
	playersOK := len(reactions[0])
	playersNOK := len(reactions[1])
	requiredReactions := 4
	if os.Getenv("mode") == "dev" {
		requiredReactions = 1
	}
	if playersOK > requiredReactions {
		log.Debugf("players confirmed match %s", match.ID)
		err = CloseMatch(match)
		if err != nil {
			log.Errorf("failed to cancel match %s: "+err.Error(), match.ID)
			return
		}
		scheduled.TaskManager.Cancel(scheduled.Task{ID: "matchvote" + match.ID})
		return
	} else if playersNOK > requiredReactions {
		log.Debugf("players refused confirmation of match %s", match.ID)
		s := discord.GetSession()
		err = s.ChannelMessageDelete(voteMessage.ChannelID, voteMessage.ID)
		if err != nil {
			log.Errorf("failed to delete message: " + err.Error())
			return
		}
		match.State = models.MatchStateInProgress
		err = db.UpdateMatch(match)
		if err != nil {
			log.Errorf("failed to update match: " + err.Error())
			return
		}
		scheduled.TaskManager.Cancel(scheduled.Task{ID: "matchvote" + match.ID})
		return
	}
}

func VoteResultMatch(match *models.Match, team1Score int, team2Score int) {
	var content string = fmt.Sprintf("Reported score : (%d-%d).\nPlease react to this message to confirm score.", team1Score, team2Score)
	log.Debugf("getting players confirmation of score (%d-%d) of match %s", team1Score, team2Score, match.ID)
	message, err := chat.CreateVoteMessage(match.ThreadID, content, []string{"✅", "❌"})
	if err != nil {
		log.Errorf("failed to create confirmation message: " + err.Error())
		return
	}
	match.State = models.MatchStateVoteInProgress
	match.Team1Score = team1Score
	match.Team2Score = team2Score
	match.VoteMessageID = message.ID
	err = db.UpdateMatch(match)
	if err != nil {
		log.Errorf("failed to update reaction: " + err.Error())
		return
	}
	scheduled.TaskManager.Add(scheduled.Task{ID: "matchvote" + match.ID, Frequency: time.Millisecond * 300, Run: func() { handleMatchVoteResult(match) }})
}

func CloseMatch(match *models.Match) error {
	session := discord.GetSession()
	channelId := discord.MatchesChannel.ID
	members, _ := session.ThreadMembers(match.ThreadID)
	for _, member := range members {
		err := session.ThreadMemberRemove(member.ID, member.UserID)
		if err != nil {
			log.Errorf("failed to kick players from match thread: " + err.Error())
		}
	}
	archive := true
	lock := true
	_, err := session.ChannelEdit(match.ThreadID, &discordgo.ChannelEdit{Archived: &archive, Locked: &lock})
	if err != nil {
		log.Errorf("failed to lock match thread: " + err.Error())
	}

	if match.Team1Score > match.Team2Score {
		match.State = models.MatchStateTeam1Won
	} else if match.Team2Score > match.Team1Score {
		match.State = models.MatchStateTeam2Won
	} else {
		match.State = models.MatchStateCanceled
	}

	message, err := session.ChannelMessage(channelId, match.MessageID)
	if err != nil {
		log.Errorf("failed to get match message: " + err.Error())
		return err
	}
	var editedMessage string
	if match.State == models.MatchStateCanceled {
		editedMessage = "~~" + message.Content + "~~" + " | Canceled"
	} else {
		editedMessage = message.Content + fmt.Sprintf(" | Final score : %d - %d", match.Team1Score, match.Team2Score)
	}

	_, err = session.ChannelMessageEdit(message.ChannelID, message.ID, editedMessage)
	if err != nil {
		log.Errorf("failed to edit match message: " + err.Error())
		return err
	}
	err = db.UpdateMatch(match)
	if err != nil {
		log.Errorf("failed to update match: " + err.Error())
	}
	if match.State != models.MatchStateCanceled {
		players := append(match.Team1, match.Team2...)
		log.Debugf("paying out players for match %s", match.ID)
		if match.State == models.MatchStateTeam1Won {
			for _, p := range match.Team1 {
				p.Credits += 10
			}
		} else if match.State == models.MatchStateTeam2Won {
			for _, p := range match.Team2 {
				p.Credits += 10
			}
		}
		for _, p := range players {
			p.Credits += 10
			err = db.UpdatePlayer(p)
			if err != nil {
				log.Errorf("failed to update player %s: "+err.Error(), p.DiscordID)
			}
		}
		log.Debugf("paying out predictions for match %s", match.ID)
		predictions, err := db.GetPlayersPredictionOnMatch(match)
		if err != nil {
			log.Errorf("failed to get predictions for match %s: "+err.Error(), match.ID)
		}
		for _, pred := range predictions {
			if match.State == models.MatchState(pred.Team) {
				pred.Player.Credits += 10
				err = db.UpdatePlayer(&pred.Player)
				if err != nil {
					log.Errorf("failed to update player %s: "+err.Error(), pred.DiscordID)
				}
			}
		}
	}

	return err
}

func threadCleanUp() {
	session := discord.GetSession()
	channelID := discord.MatchesChannel.ID
	archivedSince := time.Now().Add(-time.Hour * 4)
	if os.Getenv("mode") == "dev" {
		archivedSince = time.Now().Add(-time.Minute * 10)
	}
	threads, err := session.ThreadsArchived(channelID, &archivedSince, 100)
	if err != nil {
		log.Error("could not get archived threads: " + err.Error())
		return
	}
	for _, thread := range threads.Threads {
		_, err = session.ChannelDelete(thread.ID)
		if err != nil {
			log.Errorf("could not delete thread %s: "+err.Error(), thread.ID)
		}
	}
}

func deleteOldMatches() {
	matches, err := db.GetRunningMatchesOrderedByTimestamp()
	if err != nil {
		log.Errorf("failed to fetch running matches by timestamp: " + err.Error())
		return
	}
	for _, match := range matches {
		if err != nil {
			log.Errorf("failed to fetch running matches by timestamp: " + err.Error())
			return
		}
		cleanDelay := time.Minute * 15
		if os.Getenv("mode") == "dev" {
			cleanDelay = time.Minute
		}
		if time.Since(time.Unix(int64(match.Timestamp), 0)) > cleanDelay {
			log.Infof("cleaning match %s", match.ID)
			if os.Getenv("mode") == "dev" {
				r := rand.Intn(2)
				if r == 0 {
					match.Team1Score = 2
				} else {
					match.Team2Score = 2
				}
				err = CloseMatch(match)
			} else {
				match.State = models.MatchStateCanceled
				err = CloseMatch(match)
			}
			if err != nil {
				log.Errorf("failed to close match: " + err.Error())
				return
			}
		} else {
			break
		}
	}
}

func IsPlayerInMatch(playerID string) (bool, error) {
	p, err := getOrCreatePlayer(playerID)
	if err != nil {
		return false, err
	}
	return db.IsPlayerInMatch(p)
}
