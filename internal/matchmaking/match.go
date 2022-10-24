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
	channelId := os.Getenv("channelid")
	session := discord.GetSession()
	match := &models.Match{}
	match.ID = fmt.Sprintf("%d", matchId)
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

func CloseMatch(match *models.Match, team1Score int, team2Score int) error {
	session := discord.GetSession()
	channelId := os.Getenv("channelid")
	members, _ := session.ThreadMembers(match.ThreadID)
	for _, member := range members {
		err := session.ThreadMemberRemove(member.ID, member.UserID)
		if err != nil {
			log.Errorf("failed to kick players from match thread:" + err.Error())
		}
	}
	_, err := session.ChannelDelete(match.ThreadID)
	if err != nil {
		log.Errorf("failed to deleted match thread:" + err.Error())
	}

	match.Team1Score = team1Score
	match.Team2Score = team2Score
	if team1Score > team2Score {
		match.State = models.MatchStateTeam1Won
	} else if team2Score > team1Score {
		match.State = models.MatchStateTeam2Won
	} else {
		match.State = models.MatchStateCanceled
	}

	message, err := session.ChannelMessage(channelId, match.MessageID)
	if err != nil {
		log.Errorf("failed to get match message:" + err.Error())
	}
	var editedMessage string
	if match.State == models.MatchStateCanceled {
		err = session.ChannelMessageDelete(channelId, match.MessageID)
		if err != nil {
			log.Errorf("failed to deleted match message:" + err.Error())
		}
	} else {
		team1elo := 0
		team2elo := 0
		for _, p := range match.Team1 {
			team1elo += p.Elo
		}
		for _, p := range match.Team2 {
			team2elo += p.Elo
		}
		team1ratingChange, team2ratingChange := eloChanges(team1elo, team2elo, match.State == models.MatchStateTeam1Won)
		for _, p := range match.Team1 {
			p.Elo += team1ratingChange
			err := db.UpdatePlayer(p)
			if err != nil {
				log.Errorf("failed to update rating of player %s : "+err.Error(), p.DiscordID)
			}
		}
		for _, p := range match.Team2 {
			p.Elo += team2ratingChange
			err := db.UpdatePlayer(p)
			if err != nil {
				log.Errorf("failed to update rating of player %s : "+err.Error(), p.DiscordID)
			}
		}
		editedMessage = message.Content + fmt.Sprintf("\nFinal score : %d - %d", team1Score, team2Score)
		editedMessage += fmt.Sprintf("\nElo changes : %d vs %d", team1ratingChange, team2ratingChange)
		_, err = session.ChannelMessageEdit(message.ChannelID, message.ID, editedMessage)
		if err != nil {
			log.Errorf("failed to edit match message:" + err.Error())
		}
	}

	err = db.UpdateMatch(match)
	if err != nil {
		log.Errorf("failed to update match:" + err.Error())
	}
	return err
}

func deleteOldMatches() {
	matches, err := db.GetRunningMatchesOrderedByTimestamp()
	if err != nil {
		log.Errorf("failed to fetch running matches by timestamp:" + err.Error())
		return
	}
	for _, match := range matches {
		if err != nil {
			log.Errorf("failed to fetch running matches by timestamp:" + err.Error())
			return
		}
		cleanDelay := time.Minute * 15
		if os.Getenv("mode") == "dev" {
			cleanDelay = time.Minute
		}
		if time.Since(time.Unix(int64(match.Timestamp), 0)) > cleanDelay {
			log.Infof("cleaning match %s", match.ID)
			if os.Getenv("mode") == "dev" {
				err = CloseMatch(match, rand.Intn(6), rand.Intn(6))
			} else {
				err = CloseMatch(match, 0, 0)
			}
			if err != nil {
				log.Errorf("failed to close match:" + err.Error())
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
