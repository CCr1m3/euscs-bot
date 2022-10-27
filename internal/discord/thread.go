package discord

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func threadCleanUp() {
	channelID := MatchesChannel.ID
	archivedSince := time.Now().Add(-time.Hour * 4)
	if os.Getenv("mode") == "dev" {
		archivedSince = time.Now().Add(-time.Minute * 10)
	}
	threads, err := session.ThreadsArchived(channelID, &archivedSince, 100)
	if err != nil {
		log.Error("could not get archived threads: " + err.Error())
		return
	}
	for _, thread := range threads.Threads {
		_, err = session.ChannelDelete(thread.ID)
		if err != nil {
			log.Errorf("could not delete thread %s: "+err.Error(), thread.ID)
		}
	}
}
