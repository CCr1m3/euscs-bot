package slashcommands

import (
	"github.com/bwmarrin/discordgo"
)

type Result struct{}

func (p Result) Name() string {
	return "result"
}

func (p Result) Description() string {
	return "Allow you to report a result using scores : team1 vs team2"
}

func (p Result) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "team1-score",
			Description: "Score",
			Required:    true,
		},
		{
			Type:        discordgo.ApplicationCommandOptionInteger,
			Name:        "team2-score",
			Description: "Score",
			Required:    true,
		},
	}
}

func (p Result) Run(s *discordgo.Session, i *discordgo.InteractionCreate) {
	/*match, err := match.GetByThreadId(i.ChannelID)
	if err != nil {
		log.Warningf("failed to find match by threadID %s : "+err.Error(), i.ChannelID)
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This channel is not a match lobby",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Fatal("failed to send message")
		}
		return
	}
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	var message string = fmt.Sprintf("User %s reported score : %d - %d", i.Member.Mention(), optionMap["team1-score"].IntValue(), optionMap["team2-score"].IntValue())
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	})
	if err != nil {
		log.Fatal("failed to send message")
	}

	match.Close(int(optionMap["team1-score"].IntValue()), int(optionMap["team2-score"].IntValue()))*/
}
