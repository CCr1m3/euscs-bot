package interactions

import (
	"github.com/bwmarrin/discordgo"
	"github.com/euscs/euscs-bot/internal/discord"
)

type Interaction interface {
	Name() string
	Run(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var interactions = []Interaction{AcceptInvite{}, RefuseInvite{}}

func Init() {
	session := discord.GetSession()
	interactionHandlers := make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
	for _, interaction := range interactions {
		interactionHandlers[interaction.Name()] = interaction.Run
	}
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionMessageComponent {
			if h, ok := interactionHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})
}
