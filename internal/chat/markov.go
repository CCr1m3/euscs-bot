package chat

import (
	"bufio"
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

func learn(message string) error {
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
	return db.AddMarkovOccurences(ms)
}

func loadMarkovFromFile() {
	readFile, err := os.Open("messages")
	if err != nil {
		log.Fatal("failed to open file: " + err.Error())
		return
	}
	err = db.DeleteAllMarkov()
	if err != nil {
		log.Fatal("failed to drop table markov: " + err.Error())
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
				log.Fatal("failed to save markov occurences: " + err.Error())
			}
			ms = make([]*models.Markov, 0)
		}
	}
	err = db.AddMarkovOccurences(ms)
	if err != nil {
		log.Fatal("failed to save markov occurences: " + err.Error())
	}

	readFile.Close()
	os.Remove("messages")
}
