package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/haashi/omega-strikers-bot/internal/models"
	"github.com/sirupsen/logrus"
)

func Init(s *mux.Router) {
	s.HandleFunc("", apiHandler)
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(models.CallerIDKey)
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Hello from API: " + userID.(string)})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	//w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		logrus.Error(err)
	}
}