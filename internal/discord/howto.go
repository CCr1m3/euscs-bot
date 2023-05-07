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
	return fmt.Sprintf("Hello my name is Ai.Mi and I am this Discord's assistant. My available commands:\n\n``/sync`` to synchronize your Omega Strikers account to your Discord account. I will assign you a rank based on the one you have in-game. Synchronizing an account that does not belong to you is harshly punishable. If you wish to unsynchronize your account for any reason, contact a moderator.\n\n``/update`` for me to update your rank based on the one you have in-game.\n\n``/join`` to join my custom queue and ``/leave`` to leave it. My queue uses your in-game ELO to ensure fair matches of quality. You are free to play Omega Strikers while queueing. Just join up with the rest after you finish your game when I find a match for you.\nWhen I find a match for you I will ping you in %s. I'll give you further instructions there.\n\n``/predict`` in a thread of an ongoing match to predict who will win. If you guess correctly, I'll reward you with credits!\n\n``/credits`` to check your Ai.Mi credit balance. You can get credits by winning and losing in my queue and correctly predicting matches, as well as simply chatting!\n\nYou can spend 10 Ai.Mi credits by tagging me with a message, I will respond!\n\nHave fun!", MatchesChannel.Mention())
}
