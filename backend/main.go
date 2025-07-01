package main

import (
	"log"
	"net/http"
	"os"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/handlers"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/middleware"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// API routes
	r.Route("/api", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)

		// Auth routes (no middleware inside these handlers)
		r.Group(func(r chi.Router) {
			r.Post("/login", handlers.LoginHandler)
			r.Post("/reset-password-request", handlers.PasswordResetRequestHandler)
			r.Post("/reset-password", handlers.PasswordResetHandler)
		})

		// User + Cabin
		r.Get("/cabins", handlers.GetCabins)
		r.Get("/cabins/{id}", handlers.GetCabinByID)
		r.Post("/cabins", handlers.CreateCabin)
		r.Put("/cabins/{id}", handlers.UpdateCabin)

		r.Get("/owners", handlers.GetOwners)
		r.Post("/owners", handlers.CreateOwner)
		r.Put("/owners/{id}", handlers.UpdateOwner)
	})

	// Serve frontend static files
	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.Dir("frontend_dist"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	http.ListenAndServe(":"+port, r)
}
