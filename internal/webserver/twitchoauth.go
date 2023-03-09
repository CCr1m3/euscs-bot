package webserver

import (
	"context"
	"encoding/gob"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/env"
	"github.com/haashi/omega-strikers-bot/internal/models"
	"github.com/nicklaw5/helix"
	"golang.org/x/oauth2"
)

var twitchoauth2 oauth2.Config

func initTwitchAuth(s *mux.Router) {
	twitchoauth2 = oauth2.Config{
		RedirectURL:  env.Twitch.OAuth2RedirectURL,
		ClientID:     env.Twitch.OAuth2ID,
		ClientSecret: env.Twitch.OAuth2Secret,
		Scopes:       []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://id.twitch.tv/oauth2/authorize",
			TokenURL:  "https://id.twitch.tv/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
	gob.Register(oauth2.Token{})
	s.HandleFunc("/login", twitchAuthHandler)
	s.HandleFunc("/redirect", twitchRedirectHandler)
	s.HandleFunc("/logout", twitchLogoutHandler)
}

func twitchAuthHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	discordId := ctx.Value(models.CallerIDKey).(string)
	player, err := db.GetPlayerById(ctx, discordId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if player.TwitchID != "" {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		session, err := store.Get(r, "session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var sessionID string
		sessionIDraw, ok := session.Values["ID"]
		if ok {
			sessionID = sessionIDraw.(string)
		}
		state := uuid.New().String()
		authorizedStates[state] = sessionID
		http.Redirect(w, r, twitchoauth2.AuthCodeURL(state), http.StatusTemporaryRedirect)
	}
}

func twitchRedirectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var sessionID string
	sessionIDraw, ok := session.Values["ID"]
	if ok {
		sessionID = sessionIDraw.(string)
	}
	state := r.FormValue("state")
	if authorizedStates[state] != sessionID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("State does not match."))
		return
	}
	token, err := twitchoauth2.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	twitchSession, err := helix.NewClient(&helix.Options{ClientID: env.Twitch.OAuth2ID, UserAccessToken: token.AccessToken})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	user, err := twitchSession.GetUsers(nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	discordId := ctx.Value(models.CallerIDKey).(string)
	player, err := db.GetPlayerById(ctx, discordId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	player.TwitchID = user.Data.Users[0].ID
	err = db.UpdatePlayer(ctx, player)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func twitchLogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	discordId := ctx.Value(models.CallerIDKey).(string)
	player, err := db.GetPlayerById(ctx, discordId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	player.TwitchID = ""
	err = db.UpdatePlayer(ctx, player)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
