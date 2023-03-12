package markov

import (
	"bufio"
	"context"
	"os"
	"strings"

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
	err = db.DeleteAllMarkov()
	if err != nil {
		log.Fatal("failed to drop table markov: " + err.Error())
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
	os.Remove("messages")
}

func fetchAllMessages(ctx context.Context) {
	session := discord.GetSession()
	f, _ := os.Create("messages")
	channels, err := session.GuildChannels(discord.GuildID)
	if err != nil {
		log.Errorf("failed to get guild channels: " + err.Error())
	}
	for _, channel := range channels {
		//get all messages
		log.Debugf("getting messages from %s", channel.Name)
		messages, err := session.ChannelMessages(channel.ID, 100, "", "", "")
		if err != nil {
			log.Errorf("failed to get messages: " + err.Error())
		}
		for len(messages) != 0 {
			for _, message := range messages {
				if message.Author.ID == session.State.User.ID {
					continue
				}
				_, err = f.WriteString(strings.ToLower(message.Content) + "\n")
				if err != nil {
					log.Errorf("failed to write messages: " + err.Error())
				}
			}
			lastMessage := messages[len(messages)-1]
			messages, err = session.ChannelMessages(channel.ID, 100, lastMessage.ID, "", "")
			if err != nil {
				log.Errorf("failed to get messages: " + err.Error())
			}
		}
	}
}
