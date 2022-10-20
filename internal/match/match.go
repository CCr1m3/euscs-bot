package match

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/discord"
	log "github.com/sirupsen/logrus"
)

type Match struct {
	db.Match
}

func GetByThreadId(threadID string) (*Match, error) {
	dbMatch, err := db.GetMatchByThreadID(threadID)
	if err != nil {
		return nil, err
	}
	return &Match{*dbMatch}, nil
}

func New(team1 []*db.Player, team2 []*db.Player) error {
	matchId := rand.Intn(math.MaxInt32)
	channelId := os.Getenv("channelid")
	session := discord.GetSession()
	match := Match{}
	match.ID = fmt.Sprintf("%d", matchId)
	mentionMessage := ""
	for i := range team1 {
		mentionMessage += "<@" + team1[i].DiscordID + ">"
	}
	mentionMessage += " vs "
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
		Content: fmt.Sprintf("Lobby code : %d\nFirst user in both team was assigned goalie, you are free to discuss here to change roles.\nUse this thread to share in game name to make teams, or to chat.\nPlease report match result with /result in this thread.", matchId),
	})
	if err != nil {
		return err
	}
	match.Team1 = team1
	match.Team2 = team2
	return match.Save()
}

func (m *Match) Close(team1Score int, team2Score int) error {
	session := discord.GetSession()
	channelId := os.Getenv("channelid")
	members, _ := session.ThreadMembers(m.ThreadID)
	for _, member := range members {
		_ = session.ThreadMemberRemove(member.ID, member.UserID)
	}
	_, err := session.ChannelDelete(m.ThreadID)
	if err != nil {
		log.Errorf(err.Error())
	}
	message, _ := session.ChannelMessage(channelId, m.MessageID)
	editedMessage := message.Content + fmt.Sprintf(" | Final score : %d - %d", team1Score, team2Score)
	_, err = session.ChannelMessageEdit(message.ChannelID, message.ID, editedMessage)
	if err != nil {
		log.Errorf(err.Error())
	}
	err = m.SetScore(team1Score, team2Score)
	if err != nil {
		log.Errorf(err.Error())
	}
	return err
}

func Init() {
	go func() {
		for {
			time.Sleep(60 * time.Second)
			deleteOldMatches()
		}
	}()
}

func deleteOldMatches() {
	matches, err := db.GetMatchesByTimestamp()
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	for _, m := range matches {
		match, err := GetByThreadId(m.ThreadID)
		if err != nil {
			log.Errorf(err.Error())
			return
		}
		cleanDelay := time.Minute * 15
		if os.Getenv("mode") == "dev" {
			cleanDelay = time.Minute
		}
		if time.Since(time.Unix(int64(m.Timestamp), 0)) > cleanDelay {
			log.Infof("cleaning match %s", m.ID)
			err = match.Close(0, 0)
			if err != nil {
				log.Errorf(err.Error())
				return
			}
		} else {
			break
		}
	}
}
