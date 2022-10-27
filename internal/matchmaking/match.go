package matchmaking

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/discord"
	"github.com/haashi/omega-strikers-bot/internal/models"
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
	mentionMessage += "\nversus\n"
	for i := range team2 {
		mentionMessage += "<@" + team2[i].DiscordID + ">"
	}
	initialMessage, err := session.ChannelMessageSendComplex(channelId, &discordgo.MessageSend{
		Content: fmt.Sprintf("ID:%d\n%s", matchId, mentionMessage),
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
		Content: fmt.Sprintf("Lobby code : %d\nFirst user in both team was assigned goalie, you are free to discuss here to change roles.\nUse this thread to share in game name to make teams, or to chat.\nPlease report match result with /result in this thread.", matchId),
	})
	if err != nil {
		return err
	}
	match.Team1 = team1
	match.Team2 = team2
	return db.CreateMatch(match)
}

func VoteCancelMatch(match *models.Match) {
	var message string = "A cancel request has been sent for this match.\nPlease react to this message to confirm."
	log.Debugf("getting players confirmation of cancellation of match %s", match.ID)
	s := discord.GetSession()
	discMessage, err := s.ChannelMessageSend(match.ThreadID, message)
	if err != nil {
		log.Error("failed to send message: " + err.Error())
		return
	}
	if err != nil {
		log.Error("failed to get message: " + err.Error())
		return
	}
	err = s.MessageReactionAdd(discMessage.ChannelID, discMessage.ID, "✅")
	if err != nil {
		log.Errorf("failed add reaction: " + err.Error())
		return
	}
	err = s.MessageReactionAdd(discMessage.ChannelID, discMessage.ID, "❌")
	if err != nil {
		log.Errorf("failed add reaction: " + err.Error())
		return
	}

	for {
		reactionsOK, err := s.MessageReactions(discMessage.ChannelID, discMessage.ID, "✅", 20, "", "")
		if err != nil {
			log.Errorf("failed to get reactions: " + err.Error())
			return
		}
		playersOK := 0
		for _, reaction := range reactionsOK {
			playerID := reaction.ID
			inMatch := false
			for _, p := range match.Team1 {
				if p.DiscordID == playerID {
					inMatch = true
				}
			}
			for _, p := range match.Team2 {
				if p.DiscordID == playerID {
					inMatch = true
				}
			}
			if !inMatch && playerID != s.State.User.ID {
				err = s.MessageReactionRemove(discMessage.ChannelID, discMessage.ID, "✅", playerID)
				if err != nil {
					log.Errorf("failed to remove reactions: " + err.Error())
				}
			} else {
				playersOK++
			}
		}
		reactionsNOK, err := s.MessageReactions(discMessage.ChannelID, discMessage.ID, "❌", 20, "", "")
		if err != nil {
			log.Errorf("failed to get reactions: " + err.Error())
			return
		}
		playersNOK := 0
		for _, reaction := range reactionsNOK {
			playerID := reaction.ID
			inMatch := false
			for _, p := range match.Team1 {
				if p.DiscordID == playerID {
					inMatch = true
				}
			}
			for _, p := range match.Team2 {
				if p.DiscordID == playerID {
					inMatch = true
				}
			}
			if !inMatch && playerID != s.State.User.ID {
				err = s.MessageReactionRemove(discMessage.ChannelID, discMessage.ID, "❌", playerID)
				if err != nil {
					log.Errorf("failed to remove reactions: " + err.Error())
				}
			} else {
				playersNOK++
			}
		}
		requiredReactions := 4
		if os.Getenv("mode") == "dev" {
			requiredReactions = 1
		}
		if playersOK > requiredReactions {
			log.Debugf("players confirmed cancellation of match %s", match.ID)
			err = CancelMatch(match)
			if err != nil {
				log.Errorf("failed to cancel match %s: "+err.Error(), match.ID)
				return
			}
			return
		} else if playersNOK > requiredReactions {
			log.Debugf("players refused cancellation of match %s", match.ID)
			err = s.ChannelMessageDelete(discMessage.ChannelID, discMessage.ID)
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
			return
		}
		time.Sleep(time.Second)
	}
}

func VoteResultMatch(match *models.Match, team1Score int, team2Score int) {
	var message string = fmt.Sprintf("Reported score : (%d-%d).\nPlease react to this message to confirm score.", team1Score, team2Score)
	log.Debugf("getting players confirmation of score (%d-%d) of match %s", team1Score, team2Score, match.ID)
	s := discord.GetSession()
	discMessage, err := s.ChannelMessageSend(match.ThreadID, message)
	if err != nil {
		log.Error("failed to send message: " + err.Error())
		return
	}
	match.Team1Score = team1Score
	match.Team2Score = team2Score
	match.State = models.MatchStateVoteInProgress
	err = db.UpdateMatch(match)
	if err != nil {
		log.Errorf("failed to update match: " + err.Error())
		return
	}
	err = s.MessageReactionAdd(discMessage.ChannelID, discMessage.ID, "✅")
	if err != nil {
		log.Errorf("failed add reaction: " + err.Error())
		return
	}
	err = s.MessageReactionAdd(discMessage.ChannelID, discMessage.ID, "❌")
	if err != nil {
		log.Errorf("failed add reaction: " + err.Error())
		return
	}
	for {
		reactionsOK, err := s.MessageReactions(discMessage.ChannelID, discMessage.ID, "✅", 20, "", "")
		if err != nil {
			log.Errorf("failed to get reactions: " + err.Error())
			return
		}
		playersOK := 0
		for _, reaction := range reactionsOK {
			playerID := reaction.ID
			inMatch := false
			for _, p := range match.Team1 {
				if p.DiscordID == playerID {
					inMatch = true
				}
			}
			for _, p := range match.Team2 {
				if p.DiscordID == playerID {
					inMatch = true
				}
			}
			if !inMatch && playerID != s.State.User.ID {
				err = s.MessageReactionRemove(discMessage.ChannelID, discMessage.ID, "✅", playerID)
				if err != nil {
					log.Errorf("failed to remove reactions: " + err.Error())
				}
			} else {
				playersOK++
			}
		}
		reactionsNOK, err := s.MessageReactions(discMessage.ChannelID, discMessage.ID, "❌", 20, "", "")
		if err != nil {
			log.Errorf("failed to get reactions: " + err.Error())
			return
		}
		playersNOK := 0
		for _, reaction := range reactionsNOK {
			playerID := reaction.ID
			inMatch := false
			for _, p := range match.Team1 {
				if p.DiscordID == playerID {
					inMatch = true
				}
			}
			for _, p := range match.Team2 {
				if p.DiscordID == playerID {
					inMatch = true
				}
			}
			if !inMatch && playerID != s.State.User.ID {
				err = s.MessageReactionRemove(discMessage.ChannelID, discMessage.ID, "❌", playerID)
				if err != nil {
					log.Errorf("failed to remove reactions: " + err.Error())
				}
			} else {
				playersNOK++
			}
		}
		requiredReactions := 4
		if os.Getenv("mode") == "dev" {
			requiredReactions = 1
		}
		if playersOK > requiredReactions {
			log.Debugf("players confirmed score (%d-%d) of match %s", match.Team1Score, match.Team2Score, match.ID)
			err = CloseMatch(match)
			if err != nil {
				log.Errorf("failed to close match %s: "+err.Error(), match.ID)
				return
			}
			return
		} else if playersNOK > requiredReactions {
			log.Debugf("players refused score (%d-%d) of match %s", match.Team1Score, match.Team2Score, match.ID)
			err = s.ChannelMessageDelete(discMessage.ChannelID, discMessage.ID)
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
			return
		}
		time.Sleep(time.Second)
	}
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
	}

	message, err := session.ChannelMessage(channelId, match.MessageID)
	if err != nil {
		log.Errorf("failed to get match message: " + err.Error())
		return err
	}
	editedMessage := message.Content + fmt.Sprintf("\nFinal score : %d - %d", match.Team1Score, match.Team2Score)
	_, err = session.ChannelMessageEdit(message.ChannelID, message.ID, editedMessage)
	if err != nil {
		log.Errorf("failed to edit match message: " + err.Error())
		return err
	}
	err = db.UpdateMatch(match)
	if err != nil {
		log.Errorf("failed to update match: " + err.Error())
	}
	return err
}

func CancelMatch(match *models.Match) error {
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
	match.State = models.MatchStateCanceled
	message, err := session.ChannelMessage(channelId, match.MessageID)
	if err != nil {
		log.Errorf("failed to get match message: " + err.Error())
		return err
	}
	editedMessage := message.Content + "\nCanceled"
	_, err = session.ChannelMessageEdit(message.ChannelID, message.ID, editedMessage)
	if err != nil {
		log.Errorf("failed to edit match message: " + err.Error())
		return err
	}
	err = db.UpdateMatch(match)
	if err != nil {
		log.Errorf("failed to update match: " + err.Error())
	}
	return err
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
				err = CancelMatch(match)
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
