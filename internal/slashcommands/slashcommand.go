package slashcommands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/haashi/omega-strikers-bot/internal/discord"
	log "github.com/sirupsen/logrus"
)

type SlashCommand interface {
	Name() string
	Description() string
	Run(s *discordgo.Session, i *discordgo.InteractionCreate)
	Options() []*discordgo.ApplicationCommandOption
}

var registeredCommands []*discordgo.ApplicationCommand
var commands = []SlashCommand{Ping{}, Join{}, Leave{}, Result{}, Rank{}}

func Init() {
	session := discord.GetSession()

	commandHandlers := make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate))
	for _, command := range commands {
		commandHandlers[command.Name()] = command.Run
	}
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
	log.Println("adding commands...")
	registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))
	for i, command := range commands {
		appCommand := &discordgo.ApplicationCommand{
			Name:        command.Name(),
			Description: command.Description(),
			Options:     command.Options(),
		}
		cmd, err := session.ApplicationCommandCreate(session.State.User.ID, discord.GuildID, appCommand)
		if err != nil {
			log.Fatalf("Cannot create '%v' command: %v", command.Name(), err)
		}
		registeredCommands[i] = cmd
	}
}

func Stop() {
	session := discord.GetSession()

	log.Println("removing commands...")
	// We need to fetch the commands, since deleting requires the command ID.
	// We are doing this from the returned commands on line 375, because using
	// this will delete all the commands, which might not be desirable, so we
	// are deleting only the commands that we added.
	registeredCommands, err := session.ApplicationCommands(session.State.User.ID, discord.GuildID)
	if err != nil {
		log.Errorf("Could not fetch registered commands: %v", err)
	}

	for _, v := range registeredCommands {
		err := session.ApplicationCommandDelete(session.State.User.ID, discord.GuildID, v.ID)
		if err != nil {
			log.Errorf("cannot delete '%v' command: %v", v.Name, err)
		}
	}
}
