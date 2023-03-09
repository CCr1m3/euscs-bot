package markov

import (
	"context"
	"flag"
	"math/rand"
	"strings"

	"github.com/euscs/euscs-bot/internal/db"
	"github.com/euscs/euscs-bot/internal/discord"
	"github.com/euscs/euscs-bot/internal/models"
	log "github.com/sirupsen/logrus"
)

var reset = flag.Bool("resetmarkov", false, "")

func Init() {
	if *reset {
		ctx := context.Background()
		log.Info("fetching all messages from discord server")
		fetchAllMessages(ctx)
		log.Info("loading all messages into db")
		loadMarkovFromFile(ctx)
		log.Info("done loading messages into db")
	}
	session := discord.GetSession()
	session.AddHandler(messageHandler)
}

func parse(message string) []string {
	words := strings.Fields(message)
	return words
}

func learn(ctx context.Context, message string) error {
	words := parse(message)
	ms := make([]*models.Markov, 0)
	if len(words) > 1 {
		ms = append(ms, &models.Markov{Word1: "__start__", Word2: words[0], Word3: words[1]})
		for i := range words {
			if i == len(words)-2 {
				ms = append(ms, &models.Markov{Word1: words[i], Word2: words[i+1], Word3: "__end__"})
				break
			} else {
				ms = append(ms, &models.Markov{Word1: words[i], Word2: words[i+1], Word3: words[i+2]})
			}
		}
	}
	if len(words) == 1 {
		ms = append(ms, &models.Markov{Word1: "__start__", Word2: words[0], Word3: "__end__"})
	}
	return db.AddMarkovOccurences(ctx, ms)
}

func getNextWord(ctx context.Context, word1 string, word2 string) string {
	occurences, err := db.GetMarkovOccurencesAndTotal(ctx, word1, word2)
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

func getRandomStartingWord(ctx context.Context) string {
	occurences, err := db.GetStartingMarkovOccurences(ctx)
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

func GenerateRandomMessage(ctx context.Context) string {
	words := make([]string, 0)
	lastWord := "__start__"
	word := getRandomStartingWord(ctx)
	words = append(words, word)
	i := 0
	for word != "__end__" {
		nextWord := getNextWord(ctx, lastWord, word)
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
