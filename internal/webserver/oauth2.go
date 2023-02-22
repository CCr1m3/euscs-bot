package webserver

import (
	"context"
	"encoding/gob"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

var conf oauth2.Config

var state = "random" //to-do, generate state values at login

func initAuth(s *mux.Router) {
	conf = oauth2.Config{
		RedirectURL: os.Getenv("oauth2redirectURL"),
		// This next 2 lines must be edited before running this.
		ClientID:     os.Getenv("oauth2id"),
		ClientSecret: os.Getenv("oauth2secret"),
		Scopes:       []string{"identify"},
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://discord.com/api/oauth2/authorize",
			TokenURL:  "https://discord.com/api/oauth2/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
	gob.Register(oauth2.Token{})
	s.HandleFunc("", authHandler)
	s.HandleFunc("/redirect", redirectHandler)
	s.HandleFunc("/logout", logoutHandler)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "sessionid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, ok := session.Values["discordID"]
	if ok {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	} else {
		http.Redirect(w, r, conf.AuthCodeURL(state), http.StatusTemporaryRedirect)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "sessionid")
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
	session, err := store.Get(r, "sessionid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.FormValue("state") != state {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("State does not match."))
		return
	}
	token, err := conf.Exchange(context.Background(), r.FormValue("code"))
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
	err = session.Save(r, w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
