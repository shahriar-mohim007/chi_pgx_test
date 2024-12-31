package httpserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"go_chi_pgx/state"
	"time"
)

func routes(s *state.State) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	corsOptions := cors.Options{
		AllowedOrigins:   []string{"http://localhost"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
	}
	r.Use(cors.New(corsOptions).Handler)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/users", HandleRegisterUser(s))
		r.Post("/users/activate", HandleActivateUser(s))
		r.Post("/token/auth", HandleLogin(s))
		r.Post("/token/refresh", HandleRefreshToken(s))
	})

	r.Route("/api/v1/contacts", func(r chi.Router) {
		r.Use(AuthMiddleware(s))
		r.Get("/", HandlerGetAllContacts(s))
		r.Post("/", HandlerCreateContact(s))
		r.Get("/{id}", HandlerGetContactByID(s))
		r.Patch("/{id}", HandlerPatchContactByID(s))
		r.Delete("/{id}", HandlerDeleteContactByID(s))
	})

	return r
}
