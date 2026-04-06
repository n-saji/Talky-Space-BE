package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	authmw "talky-space-be/middleware"
	"talky-space-be/service"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

type Handler struct {
	service        *service.Service
	requestTimeout time.Duration
}

func NewHandler(service *service.Service, requestTimeout time.Duration) *Handler {
	return &Handler{
		service:        service,
		requestTimeout: requestTimeout,
	}
}

// func (h *Handler) GetRouter() *gin.Engine {

// 	router := gin.Default()
// 	router.Use(cors.New(cors.Config{
// 		AllowOrigins:     []string{"http://localhost:5050", "http://localhost:3000", "http://localhost:5173", "http://192.168.1.54:3000"},
// 		AllowMethods:     []string{"PUT", "PATCH", "GET", "DELETE", "POST"},
// 		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Set-Cookie", "token", "account_id"},
// 		ExposeHeaders:    []string{"Content-Length", "Set-Cookie", "token", "account_id"},
// 		AllowCredentials: true,
// 		MaxAge:           12 * time.Hour,
// 	}))

// 	h.RoutingUser(&router.RouterGroup)
// 	h.AuthenticationChannel(&router.RouterGroup)
// 	h.RoutingChannel(&router.RouterGroup)

// 	protected := router.Group("/")
// 	protected.Use(jwtMiddleware.AuthMiddleware())
// 	{
// 		// Add protected routes here
// 		h.RoutingMessage(protected)
// 		h.RoutingWebSockets(protected)
// 		h.ChatroomChannel(protected)
// 	}

// 	return router
// }

func NewRouter(ctx context.Context, jobService *service.Service, requestTimeout time.Duration) http.Handler {
	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			origin := req.Header.Get("Origin")
			if origin == "http://localhost:3000" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Vary", "Origin")
			}

			if req.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, req)
		})
	})
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	handler := NewHandler(jobService, requestTimeout)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", handler.LoginHandler)
		r.Post("/refresh", handler.RefreshHandler)
		r.Post("/logout", handler.LogoutHandler)
	})

	r.Route("/channel", func(r chi.Router) {
		r.Use(authmw.AuthMiddleware())
		handler.RoutingChannel(r)
	})

	r.Route("/chatrooms", func(r chi.Router) {
		r.Use(authmw.AuthMiddleware())
		handler.ChatroomChannel(r)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/register", handler.CreateUser)
		r.Group(func(r chi.Router) {
			r.Use(authmw.AuthMiddleware())
			r.Get("/me", handler.GetUser)
			r.Put("/update", handler.UpdateUser)
			r.Delete("/delete", handler.DeleteUser)
			r.Get("/look-up", handler.LookUpUser)
		})
	})

	r.Route("/messages", func(r chi.Router) {
		r.Use(authmw.AuthMiddleware())
		handler.RoutingMessage(r)
	})

	r.Route("/ws", func(r chi.Router) {
		r.Use(authmw.AuthMiddleware())
		handler.RoutingWebSockets(r)
	})

	return r
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.WriteHeader(status)
	writeJSON(w, status, map[string]any{"error": map[string]string{"code": code, "message": message}})
}
