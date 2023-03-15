package api

import (
	"encoding/json"
	"net/http"

	"github.com/euscs/euscs-bot/internal/db"
	"github.com/euscs/euscs-bot/internal/static"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func Init(s *mux.Router) {
	s.HandleFunc("", apiHandler)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(static.CallerIDKey).(string)
	player, err := db.GetPlayerByID(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Hello from API: ", "userID": userID, "twitchID": player.TwitchID})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		log.Error(err)
	}
}
