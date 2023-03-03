package webserver

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/haashi/omega-strikers-bot/internal/models"
)

func newAuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		discordIDraw, ok := session.Values["discordID"]
		if ok {
			discordID := discordIDraw.(string)
			ctx := context.WithValue(context.Background(), models.UUIDKey, uuid.New())
			ctx = context.WithValue(ctx, models.CallerIDKey, discordID)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}
