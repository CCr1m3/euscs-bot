package chat

import (
	"bufio"
	"math/rand"
	"os"
	"strings"

	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/models"
	log "github.com/sirupsen/logrus"
)

func parse(message string) []string {
	words := strings.Fields(message)
	return words
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

func GenerateRandomMessage() string {
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

func LoadMarkovFromFile(filename string) {
	readFile, err := os.Open(filename)
	if err != nil {
		log.Error("failed to open file: " + err.Error())
		return
	}
	err = db.DeleteAllMarkov()
	if err != nil {
		log.Error("failed to drop table markov: " + err.Error())
		return
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	ms := make([]*models.Markov, 0)
	for fileScanner.Scan() {
		message := fileScanner.Text()
		words := parse(message)
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
		if len(ms) > 400 {
			err = db.AddMarkovOccurences(ms)
			if err != nil {
				log.Error("failed to save markov occurences: " + err.Error())
			}
			ms = make([]*models.Markov, 0)
		}

	}

	readFile.Close()
}
