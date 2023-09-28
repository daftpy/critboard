package main

import (
	"critboard-backend/api/authAPI"
	"critboard-backend/api/feedbackAPI"
	"critboard-backend/api/submissionsAPI"
	"critboard-backend/api/uploadAPI"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
)

func InitializeRouter(db *pgxpool.Pool) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		Debug:            true,
	}).Handler)

	// Routes
	r.Post("/uploads", uploadAPI.UploadFile(db))
	r.Post("/submissions/link/create", submissionsAPI.CreateLink(db))
	r.Post("/submissions/file/create", submissionsAPI.CreateFile(db))
	r.Get("/submissions/recent/{count}", submissionsAPI.GetRecent(db))
	r.Get("/submissions/{id}", submissionsAPI.Get(db))
	r.Get("/submissions/{id}/feedback", feedbackAPI.Get(db))
	r.Post("/submissions/{id}/feedback", feedbackAPI.Create(db))

	r.Patch("/feedback/{id}", feedbackAPI.Update(db))
	r.Patch("/feedback/{id}/remove", feedbackAPI.Remove(db))
	r.Get("/feedback/{id}/replies", feedbackAPI.Get(db))
	r.Post("/feedback/{id}/replies", feedbackAPI.Create(db))

	r.Get("/auth/twitch", authAPI.TwitchAuthHandler(db))
	r.Get("/oauth/callback", authAPI.TwitchCallbackHandler(db))

	return r
}
