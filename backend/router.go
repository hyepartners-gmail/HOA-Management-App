package main

import (
	"net/http"

	"github.com/hyepartners-gmail/HOA-Management-App/backend/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func setupRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.CORSMiddleware)

	// Global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Auth routes
	r.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", LoginHandler)
		r.Post("/forgot-password", ForgotPasswordHandler)
		r.Post("/reset-password", ResetPasswordHandler)
	})

	// Cabin routes (admin + board)
	r.Route("/api/cabins", func(r chi.Router) {
		r.Use(JWTMiddleware, RoleMiddleware("admin", "president", "secretary", "treasurer"))
		r.Get("/", GetCabinsHandler)
		r.Post("/", CreateCabinHandler)
		r.Get("/{id}", GetCabinByIDHandler)
		r.Put("/{id}", UpdateCabinHandler)
	})

	// Owner routes (admin only for writes)
	r.Route("/api/owners", func(r chi.Router) {
		r.Use(JWTMiddleware, RoleMiddleware("admin", "president", "secretary", "treasurer"))
		r.Get("/", GetOwnersHandler)
		r.Group(func(r chi.Router) {
			r.Use(RoleMiddleware("admin"))
			r.Post("/", CreateOwnerHandler)
			r.Put("/{id}", UpdateOwnerHandler)
		})
	})

	// Portal routes by role
	r.Route("/api/portal", func(r chi.Router) {
		r.Use(JWTMiddleware)
		r.Route("/admin", func(r chi.Router) {
			r.Use(RoleMiddleware("admin"))
			r.Get("/dashboard", AdminDashboardHandler)
		})
		r.Route("/president", func(r chi.Router) {
			r.Use(RoleMiddleware("president"))
			r.Get("/dashboard", PresidentDashboardHandler)
		})
		r.Route("/secretary", func(r chi.Router) {
			r.Use(RoleMiddleware("secretary"))
			r.Get("/dashboard", SecretaryDashboardHandler)
		})
		r.Route("/treasurer", func(r chi.Router) {
			r.Use(RoleMiddleware("treasurer"))
			r.Get("/dashboard", TreasurerDashboardHandler)
		})
		r.Route("/cabin-owner", func(r chi.Router) {
			r.Use(RoleMiddleware("cabin_owner"))
			r.Get("/dashboard", CabinOwnerDashboardHandler)
		})
	})

	r.Route("/api/invoices", func(r chi.Router) {
		r.Use(JWTMiddleware, RoleMiddleware("admin", "treasurer", "cabin_owner"))
		r.Get("/{id}/pdf", GetInvoicePDFURLHandler)
	})

	r.Route("/api/notifications", func(r chi.Router) {
		r.Use(JWTMiddleware)
		r.Get("/", ListNotificationsHandler)
		r.Post("/seen", UpdateLastSeenNotificationHandler)

		r.Group(func(r chi.Router) {
			r.Use(RoleMiddleware("admin", "president"))
			r.Post("/", CreateNotificationHandler)
		})
	})

	r.Route("/api/message-board", func(r chi.Router) {
		r.Use(JWTMiddleware, RoleMiddleware("admin", "president", "secretary", "treasurer", "cabin_owner"))

		r.Get("/", ListPostsHandler)
		r.Post("/", CreatePostHandler)
		r.Delete("/{postID}", DeletePostHandler)

		r.Get("/{postID}/comments", ListCommentsHandler)
		r.Post("/{postID}/comments", CreateCommentHandler)
	})

	r.Route("/api/newsletters", func(r chi.Router) {
		r.Use(JWTMiddleware)

		r.Get("/", ListNewslettersHandler)

		// Admin + President can create/publish
		r.Group(func(r chi.Router) {
			r.Use(RoleMiddleware("admin", "president"))
			r.Post("/", CreateNewsletterHandler)
			r.Post("/{id}/publish", PublishNewsletterHandler)
		})
	})

	r.Route("/api/rules", func(r chi.Router) {
		r.Use(JWTMiddleware, RoleMiddleware("admin", "president", "secretary"))
		r.Get("/", GetHOARulesHandler)
		r.Put("/", UpdateHOARulesHandler) // admin only
	})

	r.Route("/api/proxies", func(r chi.Router) {
		r.Use(JWTMiddleware)
		r.Post("/", SubmitProxyHandler) // cabin owner
		r.Get("/", RoleMiddleware("admin", "secretary"), GetProxiesHandler)
	})

	r.Route("/api/agenda-requests", func(r chi.Router) {
		r.Use(JWTMiddleware)
		r.Post("/", RoleMiddleware("cabin_owner"), SubmitAgendaRequestHandler)
		r.Get("/", RoleMiddleware("secretary"), GetAgendaRequestsHandler)
	})

	r.Route("/api/minutes", func(r chi.Router) {
		r.Use(JWTMiddleware)
		r.Post("/", RoleMiddleware("secretary"), UploadMeetingMinutesHandler)
		r.Get("/", GetMeetingMinutesHandler) // public/cabin owner access
	})

	r.Route("/api/assessments", func(r chi.Router) {
		r.Use(JWTMiddleware)

		r.Group(func(r chi.Router) {
			r.Use(RoleMiddleware("admin"))
			r.Post("/", TriggerAssessmentHandler)
		})

		r.Group(func(r chi.Router) {
			r.Use(RoleMiddleware("cabin_owner", "treasurer"))
			r.Get("/", GetMyAssessmentsHandler)
		})
	})

	// Service Requests
	r.Route("/api/service-requests", func(r chi.Router) {
		r.Use(JWTMiddleware)

		r.Post("/", SubmitServiceRequestHandler)

		r.Group(func(r chi.Router) {
			r.Use(RoleMiddleware("admin"))
			r.Get("/", GetAllServiceRequestsHandler)
			r.Put("/status", UpdateServiceRequestStatusHandler)
		})
	})

	// Talent Directory
	r.Route("/api/talent", func(r chi.Router) {
		r.Get("/", GetPublicTalentHandler)
		r.Post("/", SubmitTalentHandler)

		r.Group(func(r chi.Router) {
			r.Use(JWTMiddleware, RoleMiddleware("admin"))
			r.Get("/all", GetAllTalentHandler)
			r.Put("/approve", ToggleTalentApprovalHandler)
		})
	})

	// FAQ public and admin
	r.Route("/api/faq", func(r chi.Router) {
		r.Get("/", GetFAQsHandler)

		r.Group(func(r chi.Router) {
			r.Use(JWTMiddleware, RoleMiddleware("admin"))
			r.Post("/", SaveFAQHandler)
		})
	})

	r.Route("/api/me", func(r chi.Router) {
		r.Use(JWTMiddleware)
		r.Get("/", GetProfileHandler)
		r.Put("/", UpdateProfileHandler)
		r.Put("/password", UpdatePasswordHandler)
	})

	r.Route("/api/communications", func(r chi.Router) {
		r.Use(JWTMiddleware)
		r.Get("/", ListCommunicationsHandler)
		r.Get("/{id}", GetCommunicationHandler)
	})

	r.Route("/api/documents", func(r chi.Router) {
		r.Use(JWTMiddleware)
		r.Get("/", ListDocumentsHandler)

		// Upload requires admin
		r.Group(func(r chi.Router) {
			r.Use(RoleMiddleware("admin"))
			r.Post("/", UploadDocumentHandler)
		})
	})

	r.Route("/api/audit-log", func(r chi.Router) {
		r.Use(JWTMiddleware, RoleMiddleware("admin"))
		r.Get("/", ListAuditLogsHandler)
	})

	r.Route("/api/polls", func(r chi.Router) {
		r.Use(JWTMiddleware)
		r.Get("/", ListPollsHandler)
		r.Post("/{id}/vote", SubmitVoteHandler)

		r.Group(func(r chi.Router) {
			r.Use(RoleMiddleware("president", "secretary", "treasurer", "admin"))
			r.Post("/", CreatePollHandler)
		})
	})

	return r
}
