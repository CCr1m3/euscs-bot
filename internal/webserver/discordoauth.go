package webserver

import (
	"context"
	"encoding/gob"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/haashi/omega-strikers-bot/internal/db"
	"github.com/haashi/omega-strikers-bot/internal/models"
	"github.com/haashi/omega-strikers-bot/internal/scheduled"
	"golang.org/x/oauth2"
)

var discordoauth2 oauth2.Config

var authorizedStates map[string]string

func initAuth(s *mux.Router) {
	discordoauth2 = oauth2.Config{
		RedirectURL:  os.Getenv("discordoauth2redirectURL"),
		ClientID:     os.Getenv("discordoauth2id"),
		ClientSecret: os.Getenv("discordoauth2secret"),
		Scopes:       []string{"identify"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://discord.com/api/oauth2/authorize",
			TokenURL:  "https://discord.com/api/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
	authorizedStates = make(map[string]string)
	gob.Register(oauth2.Token{})
	s.HandleFunc("/login", authHandler)
	s.HandleFunc("/redirect", redirectHandler)
	s.HandleFunc("/logout", logoutHandler)
	scheduled.TaskManager.Add(scheduled.Task{ID: "clearAuthorizedStates", Run: clearAuthorizedStates, Frequency: time.Hour})
}

func clearAuthorizedStates() {
	authorizedStates = make(map[string]string)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, ok := session.Values["discordID"]
	if ok {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		sessionID := uuid.New().String()
		state := uuid.New().String()
		session.Values["ID"] = sessionID
		authorizedStates[state] = sessionID
		session.Save(r, w)
		http.Redirect(w, r, discordoauth2.AuthCodeURL(state), http.StatusTemporaryRedirect)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Options.MaxAge = -1
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(context.Background(), models.UUIDKey, uuid.New())
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
	token, err := discordoauth2.Exchange(context.Background(), r.FormValue("code"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	discordSession, err := discordgo.New(token.TokenType + " " + token.AccessToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	user, err := discordSession.User("@me")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	session.Values["discordID"] = user.ID
	db.CreatePlayer(ctx, user.ID)
	err = session.Save(r, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
