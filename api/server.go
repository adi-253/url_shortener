	package api

	import (
		"encoding/json"
		"log/slog"
		"net/http"
		"os"

		"github.com/adi-253/url_shortener/database"
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
			original_url := database.Fetch(shortened_url)
			http.Redirect(w, r, original_url, http.StatusSeeOther)
			logger.Info("Redirected to Original Url")
		}
	}

	// func (s *Server) shortenHandler(w http.ResponseWriter, r *http.Request) {
	// 	if r.Method != http.MethodPost {
	// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 		return
	// 	}

	// 	longUrl := r.FormValue("long_url")
	// 	shortUrl := r.FormValue("short_url")

	// 	if longUrl == "" || shortUrl == "" {
	// 		http.Error(w, "Missing long_url or short_url", http.StatusBadRequest)
	// 		return
	// 	}

	// 	database.UpdateDB(longUrl, shortUrl)
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte("URL added to the database"))
	// }