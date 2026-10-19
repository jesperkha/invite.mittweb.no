package http

import (
	"context"
	"net/http"
	"time"

	"singel-sondag/config"
	"singel-sondag/domain"
	"singel-sondag/http/routes"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jesperkha/notifier"
)

type Server struct {
	mux     *chi.Mux
	config  *config.Config
	cleanup func()
}

type Dependencies struct {
	Logger       domain.Logger
	Notif        *notifier.Notifier
	Config       *config.Config
	RegisterRepo domain.RegisterRepository
}

func Run(deps Dependencies) {
	mux := chi.NewMux()

	mux.Use(requestLogger(deps.Logger))
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Mount routes
	mux.Get("/", routes.InviteHandler(deps.RegisterRepo))
	mux.Get("/favicon.svg", routes.StaticFile("web/favicon.svg"))
	mux.Post("/register", routes.RegisterHandler(deps.RegisterRepo))
	mux.Get("/admin", routes.AdminPageHandler(deps.Config))
	mux.Post("/admin/login", routes.AdminLoginHandler(deps.Config))
	mux.Post("/admin/logout", routes.AdminLogoutHandler)
	mux.Group(func(admin chi.Router) {
		admin.Use(routes.AdminAuth(deps.Config))
		admin.Get("/admin/api/event", routes.GetEventHandler(deps.RegisterRepo))
		admin.Put("/admin/api/event", routes.UpdateEventHandler(deps.RegisterRepo))
		admin.Get("/admin/api/registrations", routes.GetRegistrationsHandler(deps.RegisterRepo))
		admin.Delete("/admin/api/registrations/{id}", routes.DeleteRegistrationHandler(deps.RegisterRepo))
		admin.Delete("/admin/api/registrations", routes.DeleteAllRegistrationsHandler(deps.RegisterRepo))
	})

	cleanup := func() {
		// Add cleanup here...
		// Runs after server shutdown and before exit
	}

	s := &Server{
		mux:     mux,
		config:  deps.Config,
		cleanup: cleanup,
	}

	s.ListenAndServe(deps.Notif, deps.Logger)
}

func (s *Server) ListenAndServe(notif *notifier.Notifier, logger domain.Logger) {
	done, finish := notif.Register()
	ctx := context.Background()

	server := &http.Server{
		Addr:    s.config.Port,
		Handler: s.mux,
	}

	go func() {
		<-done
		if err := server.Shutdown(ctx); err != nil {
			logger.Error(ctx, "unexpected shutdown", "error", err.Error())
		}

		logger.Info(ctx, "cleaning up")
		s.cleanup()
		finish()
	}()

	logger.Info(ctx, "server is running", "port", s.config.Port)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error(ctx, "server closed with error", "error", err.Error())
	}
}

func requestLogger(logger domain.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := &responseWriter{ResponseWriter: w, status: 200}

			next.ServeHTTP(ww, r)

			duration := time.Since(start)
			logger.Info(r.Context(), "http request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", ww.status,
				"duration", duration.String(),
				"remote_addr", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
