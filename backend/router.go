package main

import (
	"net/http"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/handlers"
	"github.com/hyepartners-gmail/HOA-Management-App/backend/middleware"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func setupRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.CORSMiddleware)

	// Global middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// Auth routes
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", handlers.LoginHandler)
		r.Post("/forgot-password", handlers.PasswordResetRequestHandler)
		r.Post("/reset-password", handlers.PasswordResetHandler)
	})

	// Cabin routes (admin + board)
	r.Route("/api/cabins", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware, middleware.RoleMiddleware("admin", "president", "secretary", "treasurer"))
		r.Get("/", handlers.GetCabins)
		r.Post("/", handlers.CreateCabin)
		r.Get("/{id}", handlers.GetCabinByID)
		r.Put("/{id}", handlers.UpdateCabin)
	})

	// Owner routes (admin only for writes)
	r.Route("/api/owners", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware, middleware.RoleMiddleware("admin", "president", "secretary", "treasurer"))
		r.Get("/", handlers.GetOwners)
		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleMiddleware("admin"))
			r.Post("/", handlers.CreateOwner)
			r.Put("/{id}", handlers.UpdateOwner)
		})
	})

	// // Portal routes by role
	// r.Route("/api/portal", func(r chi.Router) {
	// 	r.Use(middleware.JWTMiddleware)
	// 	r.Route("/admin", func(r chi.Router) {
	// 		r.Use(middleware.RoleMiddleware("admin"))
	// 		r.Get("/dashboard", handlers.AdminDashboardHandler)
	// 	})
	// 	r.Route("/president", func(r chi.Router) {
	// 		r.Use(middleware.RoleMiddleware("president"))
	// 		r.Get("/dashboard", handlers.PresidentDashboardHandler)
	// 	})
	// 	r.Route("/secretary", func(r chi.Router) {
	// 		r.Use(middleware.RoleMiddleware("secretary"))
	// 		r.Get("/dashboard", handlers.SecretaryDashboardHandler)
	// 	})
	// 	r.Route("/treasurer", func(r chi.Router) {
	// 		r.Use(middleware.RoleMiddleware("treasurer"))
	// 		r.Get("/dashboard", handlers.TreasurerDashboardHandler)
	// 	})
	// 	r.Route("/cabin-owner", func(r chi.Router) {
	// 		r.Use(middleware.RoleMiddleware("cabin_owner"))
	// 		r.Get("/dashboard", handlers.CabinOwnerDashboardHandler)
	// 	})
	// })

	r.Route("/api/invoices", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware, middleware.RoleMiddleware("admin", "treasurer", "cabin_owner"))
		r.Get("/{id}/pdf", handlers.GetInvoicePDFURLHandler)
	})

	r.Route("/api/notifications", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)
		r.Get("/", handlers.ListNotificationsHandler)
		r.Post("/seen", handlers.UpdateLastSeenNotificationHandler)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleMiddleware("admin", "president"))
			r.Post("/", handlers.CreateNotificationHandler)
		})
	})

	r.Route("/api/message-board", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware, middleware.RoleMiddleware("admin", "president", "secretary", "treasurer", "cabin_owner"))

		r.Get("/", handlers.ListPostsHandler)
		r.Post("/", handlers.CreatePostHandler)
		r.Delete("/{postID}", handlers.DeletePostHandler)

		r.Get("/{postID}/comments", handlers.ListCommentsHandler)
		r.Post("/{postID}/comments", handlers.CreateCommentHandler)
	})

	r.Route("/api/newsletters", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)

		r.Get("/", handlers.ListNewslettersHandler)

		// Admin + President can create/publish
		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleMiddleware("admin", "president"))
			r.Post("/", handlers.CreateNewsletterHandler)
			r.Post("/{id}/publish", handlers.PublishNewsletterHandler)
		})
	})

	r.Route("/api/rules", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware, middleware.RoleMiddleware("admin", "president", "secretary"))
		r.Get("/", handlers.GetHOARulesHandler)
		r.Put("/", handlers.UpdateHOARulesHandler) // admin only
	})

	r.Route("/api/proxies", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)
		r.Post("/", handlers.SubmitProxyHandler) // cabin owner

		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleMiddleware("admin", "president", "secretary"))
			r.Get("/", handlers.GetProxiesHandler)
		})
	})

	r.Route("/api/agenda-requests", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleMiddleware("cabin_owner"))
			r.Post("/", handlers.SubmitAgendaRequestHandler)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleMiddleware("secretary"))
			r.Get("/", handlers.GetAgendaRequestsHandler)
		})
	})

	r.Route("/api/minutes", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleMiddleware("secretary"))
			r.Post("/", handlers.UploadMeetingMinutesHandler)
		})

		r.Get("/", handlers.GetMeetingMinutesHandler) // accessible to all authenticated users
	})

	r.Route("/api/assessments", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleMiddleware("admin"))
			r.Post("/", handlers.TriggerAssessmentHandler)
		})

		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleMiddleware("cabin_owner", "treasurer"))
			r.Get("/", handlers.GetMyAssessmentsHandler)
		})
	})

	// Service Requests
	r.Route("/api/service-requests", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)

		r.Post("/", handlers.SubmitServiceRequestHandler)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleMiddleware("admin"))
			r.Get("/", handlers.GetAllServiceRequestsHandler)
			r.Put("/status", handlers.UpdateServiceRequestStatusHandler)
		})
	})

	// Talent Directory
	r.Route("/api/talent", func(r chi.Router) {
		r.Get("/", handlers.GetPublicTalentHandler)
		r.Post("/", handlers.SubmitTalentHandler)

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTMiddleware, middleware.RoleMiddleware("admin"))
			r.Get("/all", handlers.GetAllTalentHandler)
			r.Put("/approve", handlers.ToggleTalentApprovalHandler)
		})
	})

	// FAQ public and admin
	r.Route("/api/faq", func(r chi.Router) {
		r.Get("/", handlers.GetFAQsHandler)

		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTMiddleware, middleware.RoleMiddleware("admin"))
			r.Post("/", handlers.SaveFAQHandler)
		})
	})

	r.Route("/api/me", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)
		r.Get("/", handlers.GetProfileHandler)
		r.Put("/", handlers.UpdateProfileHandler)
		r.Put("/password", handlers.UpdatePasswordHandler)
	})

	r.Route("/api/communications", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)
		r.Get("/", handlers.ListCommunicationsHandler)
		r.Get("/{id}", handlers.GetCommunicationHandler)
	})

	r.Route("/api/documents", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)
		r.Get("/", handlers.ListDocumentsHandler)

		// Upload requires admin
		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleMiddleware("admin"))
			r.Post("/", handlers.UploadDocumentHandler)
		})
	})

	r.Route("/api/audit-log", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware, middleware.RoleMiddleware("admin"))
		r.Get("/", handlers.ListAuditLogsHandler)
	})

	r.Route("/api/polls", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware)
		r.Get("/", handlers.ListPollsHandler)
		r.Post("/{id}/vote", handlers.SubmitVoteHandler)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleMiddleware("president", "secretary", "treasurer", "admin"))
			r.Post("/", handlers.CreatePollHandler)
		})
	})

	r.Route("/api/invoices", func(r chi.Router) {
		r.Use(middleware.JWTMiddleware, middleware.RoleMiddleware("admin", "treasurer", "cabin_owner"))
		r.Get("/", handlers.GetInvoicesHandler)
		r.Get("/{id}", handlers.GetInvoiceByIDHandler)
		r.Get("/{id}/pdf", handlers.GetInvoicePDFURLHandler)

		r.Group(func(r chi.Router) {
			r.Use(middleware.RoleMiddleware("admin"))
			r.Post("/", handlers.CreateInvoiceHandler)
			r.Post("/manual-payment", handlers.RecordManualPaymentHandler)
			r.Post("/{id}/mark-paid", handlers.MarkInvoicePaidHandler)
			r.Post("/generate-quarterly", handlers.GenerateQuarterlyInvoicesHandler)
		})
	})

	return r
}
