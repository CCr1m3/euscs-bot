package chat

import (
	"flag"
	"math/rand"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/discord"
	log "github.com/sirupsen/logrus"
)

var reset = flag.Bool("resetmarkov", false, "")

func Init() {
	if *reset {
		log.Info("fetching all messages from discord server")
		fetchAllMessages()
		log.Info("loading all messages into db")
		loadMarkovFromFile()
		log.Info("done loading messages into db")
	}
	session := discord.GetSession()
	session.AddHandler(messageHandler)
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	r := regexp.MustCompile(s.State.User.Mention())
	if r.MatchString(m.Content) {
		player, err := db.GetPlayerById(m.Author.ID)
		if err != nil {
			log.Error("failed to get player: " + err.Error())
			return
		}
		if player.Currency >= 20 {
			player.Currency -= 20
			_, err := s.ChannelMessageSend(m.ChannelID, generateRandomMessage())
			if err != nil {
				log.Error("failed to send message: " + err.Error())
				return
			}
			err = db.UpdatePlayer(player)
			if err != nil {
				log.Error("failed to update player: " + err.Error())
				return
			}
		}
	} else {
		err := learn(m.Content)
		if err != nil {
			log.Error("failed to learn message: " + err.Error())
		}
	}
}

func getNextWord(word1 string, word2 string) string {
	occurences, err := db.GetMarkovOccurencesAndTotal(word1, word2)
	if err != nil {
		log.Error("error getting occurences: " + err.Error())
	}
	if len(occurences) > 0 {
		r := rand.Intn(occurences[0].Total)
		for _, occurence := range occurences {
			if r < occurence.Count {
				return occurence.Word3
			}
			r -= occurence.Count
		}
	}
	return "__end__"
}

func getRandomStartingWord() string {
	occurences, err := db.GetStartingMarkovOccurences()
	if err != nil {
		log.Error("error getting occurences: " + err.Error())
	}
	if len(occurences) > 0 {
		r := rand.Intn(occurences[0].Total)
		for _, occurence := range occurences {
			if r < occurence.Count {
				return occurence.Word2
			}
			r -= occurence.Count
		}
	}
	return "__end__"
}

func generateRandomMessage() string {
	words := make([]string, 0)
	lastWord := "__start__"
	word := getRandomStartingWord()
	words = append(words, word)
	i := 0
	for word != "__end__" {
		nextWord := getNextWord(lastWord, word)
		if nextWord == "__end__" || i > 12 {
			break
		}
		words = append(words, nextWord)
		lastWord = word
		word = nextWord
		i++
	}
	return strings.Join(words, " ")
}
