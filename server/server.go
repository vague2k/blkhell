package server

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vague2k/blkhell/views/assets"
)

type Server struct {
	Port string

	router *chi.Mux
	server http.Server
}

func NewServer(port string) *Server {
	s := &Server{
		Port:   ":" + port,
		router: chi.NewRouter(),
	}

	s.router.Use(middleware.Logger)
	return s
}

// TODO: should implement graceful shutdown
func (s *Server) Run() error {
	s.server = http.Server{
		Addr:    s.Port,
		Handler: s.router,
	}
	return s.server.ListenAndServe()
}

func (s *Server) SetupAssetsRoutes() {
	isDevelopment := os.Getenv("GO_ENV") != "production"

	assetHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isDevelopment {
			w.Header().Set("Cache-Control", "no-store")
		}

		var fs http.Handler
		if isDevelopment {
			fs = http.FileServer(http.Dir("./views/assets"))
		} else {
			fs = http.FileServer(http.FS(assets.Assets))
		}

		fs.ServeHTTP(w, r)
	})

	s.router.Handle("GET /assets/*", http.StripPrefix("/assets/", assetHandler))
}
