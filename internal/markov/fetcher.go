package markov

import (
	"bufio"
	"context"
	"os"
	"strings"
	"time"

	"github.com/euscs/euscs-bot/internal/db"
	"github.com/euscs/euscs-bot/internal/discord"
	log "github.com/sirupsen/logrus"
)

func loadMarkovFromFile(ctx context.Context) {
	readFile, err := os.Open("messages")
	if err != nil {
		log.Fatal("failed to open file: " + err.Error())
		return
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	ms := make([]*db.Markov, 0)
	for fileScanner.Scan() {
		message := fileScanner.Text()
		words := parse(message)
		if len(words) > 1 {
			ms = append(ms, &db.Markov{Word1: "__start__", Word2: words[0], Word3: words[1]})
			for i := range words {
				if i == len(words)-2 {
					ms = append(ms, &db.Markov{Word1: words[i], Word2: words[i+1], Word3: "__end__"})
					break
				} else {
					ms = append(ms, &db.Markov{Word1: words[i], Word2: words[i+1], Word3: words[i+2]})
				}
			}
		}
		if len(words) == 1 {
			ms = append(ms, &db.Markov{Word1: "__start__", Word2: words[0], Word3: "__end__"})

		}
		if len(ms) > 400 {
			err = db.AddMarkovOccurences(ctx, ms)
			if err != nil {
				log.Fatal("failed to save markov occurences: " + err.Error())
			}
			ms = make([]*db.Markov, 0)
		}
	}
	err = db.AddMarkovOccurences(ctx, ms)
	if err != nil {
		log.Fatal("failed to save markov occurences: " + err.Error())
	}
	readFile.Close()
}

func fetchAllMessages(ctx context.Context) {
	session := discord.GetSession()
	f, _ := os.Create("messages")
	var err error
	channels, err := session.GuildChannels(discord.GuildID)
	if err != nil {
		log.Errorf("failed to get guild channels: " + err.Error())
	}
	countIte, countLines := 1, 0
	for _, channel := range channels {
		//get all messages
		log.Debugf("getting messages from %s", channel.Name)
		messages, err := session.ChannelMessages(channel.ID, 100, "", "", "")
		if err != nil {
			log.Errorf("failed to get messages: " + err.Error())
		}
		for len(messages) != 0 {
			for _, message := range messages {
				if message.Author.ID == session.State.User.ID || len(message.Content) > 100 {
					continue
				}
				_, err = f.WriteString(strings.ToLower(message.Content) + "\n")
				if err != nil {
					log.Errorf("failed to write messages: " + err.Error())
				}
			}
			lenMsgs := len(messages)
			lastMessage := messages[lenMsgs-1]
			countLines += lenMsgs
			messages, err = session.ChannelMessages(channel.ID, 100, lastMessage.ID, "", "")
			if err != nil {
				// definitely need a better way of catching this error
				if strings.Contains(err.Error(), "500 Internal Server Error") {
					log.Info("failed to request messages" + err.Error())
					log.Info("reattempting request...")
					time.Sleep(time.Second / 10)
					messages, err = session.ChannelMessages(channel.ID, 100, lastMessage.ID, "", "")
					if err != nil {
						log.Errorf("failed to reattempt request" + err.Error())
					}
				} else {
					log.Errorf("failed to get messages: " + err.Error())
				}
			}
			countIte++
			if countIte == 10 {
				log.Info("done reading lines: ", countLines)
				loadMarkovFromFile(ctx)
				err = os.Remove("messages")
				if err != nil {
					log.Error("failed to delete markovfile", err.Error)
				}
				f, _ = os.Create("messages")
				countIte = 0
			}

		}
	}
}
