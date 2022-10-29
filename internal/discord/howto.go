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
	return fmt.Sprintf("Hello my name is Ai.Mi and I am this discord's assistant. My available commands:\n\n``/link`` to link your Omega Strikers account to your discord account. I will assign you a rank based on the one you have ingame. Linking an account that does not belong to you is harshly punishable. If you wish to unlink your account for any reason contact a moderator.\n\n``/update`` for me to update your rank based on the one you have ingame.\n\n``/join`` to join my custom queue and ``/leave`` to leave it. My queue uses your ingame elo to ensure fair matches of quality. You are free to play Omega Strikers while queueing. Just join up with the rest after you finish your game when I find a match for you.\nWhen I find a match for you I will ping you in #matches. I'll give you further instructions there.\n\n``/predict`` in a thread of an ongoing match to predict who will win. If you guess correctly, I'll reward you with credits!\n\n``/credits`` to check your Ai.Mi credit balance. You can get credits by winning (20) and losing (10) in my queue and correctly predicting matches (5).\n\nYou can spend 20 Ai.Mi credits by tagging me with a message, I will respond!\n\nHave fun!", MatchesChannel.Mention())
}
