package webserver

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/haashi/omega-strikers-bot/internal/api"
	"github.com/haashi/omega-strikers-bot/web"
	log "github.com/sirupsen/logrus"
)

type spaHandler struct {
	staticFS   embed.FS
	staticPath string
	indexPath  string
}

var store = sessions.NewCookieStore([]byte(os.Getenv("session_key")))

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	file := r.URL.Path
	/*if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}*/
	// prepend the path with the path to the static directory
	file = path.Join(h.staticPath, file)

	_, err := h.staticFS.Open(file)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		index, err := h.staticFS.ReadFile(path.Join(h.staticPath, h.indexPath))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusAccepted)
		w.Write(index)
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// get the subdirectory of the static dir
	statics, _ := fs.Sub(h.staticFS, h.staticPath)
	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.FS(statics)).ServeHTTP(w, r)
}

func Init() {
	log.Info("starting web service")
	r := mux.NewRouter()
	sapi := r.PathPrefix("/api").Subrouter()
	sapi.Use(newAuthHandler)
	api.Init(sapi)
	sauth := r.PathPrefix("/auth").Subrouter()
	initAuth(sauth)
	spa := spaHandler{staticFS: web.StaticFiles, staticPath: "dist", indexPath: "index.html"}
	r.PathPrefix("/").Handler(spa)
	err := http.ListenAndServe(":9000", r)
	if err != nil {
		log.Fatal("failed to launch web service")
	}
}
