package discord

import (
	"fmt"
)

func initHowTo() error {
	channelMessages, err := session.ChannelMessages(HowToChannel.ID, 100, "", "", "")
	if err != nil {
		return err
	}
	if len(channelMessages) == 0 {
		_, err := session.ChannelMessageSend(HowToChannel.ID, howtomessage())
		if err != nil {
			return err
		}
	} else {
		_, err := session.ChannelMessageEdit(HowToChannel.ID, channelMessages[0].ID, howtomessage())
		if err != nil {
			return err
		}
	}
	return nil
}

func howtomessage() string {
	return fmt.Sprintf(`Hello my name is Ai.Mi and I am this discord's assitant.
To access my features, you must first use the /link command to link your omega strikers account to your discord account. Linking to someone else account can be punished by moderators.
If you wish to unlink for any reason, contact a mod.
This will allow to assign your rank role and the usage of /update to manually update your role once you rank up (only peak rank is updated).
Once linked, you can access the discord matchmaking queue using /join or leave it using /leave. The matchmaking queue will use your in game rank to ensure balanced matches.
You are encouraged to both queue solo queue and here to find faster matches.
Once a match has been found, you will be pinged into a thread inside the %s channel. Further instructions will be given inside the channel on how to report in game score.
You gain currency from participating, and more by winning. You can use /currency to know how much coins you have.
I also respond to pings, so feel free to @ me to get a funny message (or maybe not, I am bot after all). My response costs you currency.`, MatchesChannel.Mention())
}
