	package api

	import (
		"encoding/json"
		"log/slog"
		"net/http"
		"os"

		"github.com/adi-253/url_shortener/database"
		"github.com/adi-253/url_shortener/cache"
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

	type Server struct {
		Router   *mux.Router
	}

	func InitServer() *Server {
		s := &Server{
			Router: mux.NewRouter(),
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
			database.UpdateDB(req.URL, short_id)
			out.ShortURL = "http://localhost:8080/" + short_id
			logger.Info("Created shorturl")
			if err := json.NewEncoder(w).Encode(out); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	func (s *Server) redirect() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		shortened_url := vars["url"]

		var original_url string
		var err error

		// Step 1: Try Redis cache
		original_url, err = cache.CacheGet(shortened_url)
		if err != nil {
			// Step 2: Fallback to DB (no error, check empty string)
			original_url = database.Fetch(shortened_url)
			if original_url == "" {
				http.Error(w, "URL not found", http.StatusNotFound)
				logger.Error("Short URL not found in DB or cache")
				return
			}

			// Step 3: Put in cache (optional)
			if err := cache.CachePut(shortened_url, original_url); err != nil {
				logger.Info("Could not cache result", "error", err.Error())
			}
		}

		// Step 4: Redirect to original
		http.Redirect(w, r, original_url, http.StatusSeeOther)
		logger.Info("Redirected to original URL", "url", original_url)
	}
}	
