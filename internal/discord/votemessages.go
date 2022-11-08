package discord

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func CreateVoteMessage(channelID string, content string, choices []string) (*discordgo.Message, error) {
	s := GetSession()
	discMessage, err := s.ChannelMessageSend(channelID, content)
	if err != nil {
		log.Error("failed to send message: " + err.Error())
		return nil, err
	}
	for _, choice := range choices {
		err = s.MessageReactionAdd(discMessage.ChannelID, discMessage.ID, choice)
		if err != nil {
			log.Errorf("failed add reaction: " + err.Error())
			return nil, err
		}
	}
	return discMessage, nil
}

func FetchVoteResults(message *discordgo.Message, choices []string, allowedVoters []string) ([][]*discordgo.User, error) {
	log.Debugf("fetching vote results for message %s with choices %v", message.ID, choices)
	s := GetSession()
	alreadyVoted := make(map[string]bool)
	voteReactions := make([][]*discordgo.User, 0)
	for _, choice := range choices {
		voteReactionsChoice := make([]*discordgo.User, 0)
		reactions, err := s.MessageReactions(message.ChannelID, message.ID, choice, 20, "", "")
		if err != nil {
			log.Errorf("failed to get reactions: " + err.Error())
			return nil, err
		}
		for _, reaction := range reactions {
			reacterID := reaction.ID
			allowed := true
			if len(allowedVoters) > 0 {
				allowed = false
				for _, id := range allowedVoters {
					if id == reacterID {
						allowed = true
					}
				}
			}
			if (!allowed || alreadyVoted[reacterID]) && reacterID != s.State.User.ID {
				err = s.MessageReactionRemove(message.ChannelID, message.ID, choice, reacterID)
				if err != nil {
					log.Errorf("failed to remove reactions: " + err.Error())
					return nil, err
				}
			} else {
				alreadyVoted[reacterID] = true
				voteReactionsChoice = append(voteReactionsChoice, reaction)
			}
		}
		voteReactions = append(voteReactions, voteReactionsChoice)
	}
	return voteReactions, nil
}
