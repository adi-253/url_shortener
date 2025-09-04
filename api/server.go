package api

import (
	"encoding/json"
	"log/slog"
	"os"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/teris-io/shortid"
)


var logger = slog.New(slog.NewTextHandler(os.Stderr, nil))

// var storeUrl = make(map[string]string)
type ShortenRequest struct {
    URL string `json:"url"`
}

type ShortenResponse struct {
    ShortURL string `json:"short_url"`
}


type Server struct{
	Router *mux.Router
	storeUrl map[string]string
}

func InitServer() *Server{
	s:= &Server{
		Router: mux.NewRouter(),
		storeUrl: make(map[string]string),
	}
	s.Routes()

	logger.Info("Server Created")
	return s
}

func (s *Server) Routes() {
	s.Router.HandleFunc("/post_url", s.shorten_url()).Methods("POST")
	s.Router.HandleFunc("/{url}", s.redirect()).Methods("GET")
}

func (s *Server) shorten_url() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ShortenRequest
		var out ShortenResponse

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		short_id, _ := shortid.Generate()
		s.storeUrl[short_id] = req.URL
		out.ShortURL = "http://localhost:8080/" + short_id
		logger.Info("Created shorturl")
		if err := json.NewEncoder(w).Encode(out); err!=nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) redirect() http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request)  {
		vars := mux.Vars(r)
		shortened_url := vars["url"]
		original_url := s.storeUrl[shortened_url]

		http.Redirect(w, r, original_url, http.StatusSeeOther)
		logger.Info("Redirected to Original Url")

	}
}