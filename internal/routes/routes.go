package routes

import (
	"gateWay/internal/handler/auth"
	"gateWay/internal/handler/profile"
	"gateWay/internal/handler/projects"
	"gateWay/internal/handler/teams"
	"gateWay/internal/helpers"
	authv1 "github.com/EvgGo/proto/proto/gen/go/sso"
	workspacev1 "github.com/EvgGo/proto/proto/gen/go/teamAndProjects"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	"time"
)

const (
	maxJSONBodyBytes = int64(1 << 20) // 1 MiB - обычно достаточно для auth/profile
)

func NewRoutes(log *slog.Logger, authClient authv1.AuthClient, profileClient authv1.UserProfileClient,
	workspaceTeamsClient workspacev1.TeamsClient, workspaceProjectsClient workspacev1.ProjectsClient, skillClient authv1.SkillsClient) *chi.Mux {

	const requestMaxAge = 300

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Timeout(15 * time.Second))

	// CORS middleware
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"}, // TODO: заменить на список доменов
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Request-Id"},
		ExposedHeaders:   []string{"X-Request-Id"},
		AllowCredentials: true,
		MaxAge:           requestMaxAge,
	}).Handler)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8082/swagger/doc.json"), //The url pointing to API definition
	))

	r.Route("/auth", func(r chi.Router) {
		// Публичные
		r.With(allowJSON, httprate.LimitByIP(5, 1*time.Minute)).
			Post("/register", auth.RegisterHandler(log, authClient))

		// login - ограничиваем жестче (анти-брутфорс)
		r.With(allowJSON, httprate.LimitByIP(10, 1*time.Minute)).
			Post("/login", auth.LoginHandler(log, authClient))

		// refresh - чтобы не спамили, но и не ломать UX
		r.With(allowJSON, httprate.LimitByIP(30, 1*time.Minute)).
			Post("/refresh", auth.RefreshHandler(log, authClient))

		// forgot/reset - тоже ограничиваем (спам почты/токенов)
		r.With(allowJSON, httprate.LimitByIP(5, 10*time.Minute)).
			Post("/forgot-password", auth.ForgotPasswordHandler(log, authClient))

		r.With(allowJSON, httprate.LimitByIP(10, 10*time.Minute)).
			Post("/reset-password", auth.ResetPasswordHandler(log, authClient))

		// Приватные (нужен access token)
		r.Group(func(r chi.Router) {
			r.Use(helpers.RequireBearerAuth)

			r.With(allowJSON, httprate.LimitByIP(20, 1*time.Minute)).
				Post("/logout", auth.LogoutHandler(log, authClient))

			r.With(allowJSON, httprate.LimitByIP(5, 10*time.Minute)).
				Post("/change-password", auth.ChangePasswordHandler(log, authClient))
		})
	})

	r.Route("/users", func(r chi.Router) {
		// Публичные
		r.With(httprate.LimitByIP(60, 1*time.Minute)).Get("/", profile.ListUsersHandler(log, profileClient))
		r.With(httprate.LimitByIP(120, 1*time.Minute)).Get("/{user_id}", profile.GetProfileHandler(log, profileClient))
		r.With(httprate.LimitByIP(60, 1*time.Minute)).
			Get("/by-ids", profile.GetProfilesByIdsHandler(log, profileClient))

		// Приватные
		r.Group(func(r chi.Router) {
			r.Use(helpers.RequireBearerAuth)

			// Очень полезно: если фронт зациклился, это не даст улететь в тысячи запросов
			r.With(httprate.LimitByIP(10, 1*time.Minute)).Get("/me", profile.GetMeHandler(log, profileClient))
			r.With(httprate.LimitByIP(10, 1*time.Minute)).With(allowJSON).Patch("/me", profile.UpdateMeHandler(log, profileClient))
		})
	})

	r.Route("/skills", func(r chi.Router) {
		// Публичные
		r.With(httprate.LimitByIP(60, 1*time.Minute)).Get("/", profile.ListSkillsHandler(log, skillClient))
		r.With(httprate.LimitByIP(120, 1*time.Minute)).Get("/{skill_id}", profile.GetSkillHandler(log, skillClient))
		r.With(httprate.LimitByIP(60, 1*time.Minute)).Post("/by-ids", profile.GetSkillsByIdsHandler(log, skillClient))
		// Приватные
		r.Group(func(r chi.Router) {

			r.Use(helpers.RequireBearerAuth)

			r.With(httprate.LimitByIP(10, 1*time.Minute)).With(allowJSON).Post("/", profile.CreateSkillHandler(log, skillClient))
			r.With(httprate.LimitByIP(10, 1*time.Minute)).Delete("/{skill_id}", profile.DeleteSkillHandler(log, skillClient))
		})
	})

	r.Route("/teams", func(r chi.Router) {
		// только авторизованным
		r.Use(helpers.RequireBearerAuth)

		// CreateTeam / ListTeams
		r.With(allowJSON, httprate.LimitByIP(20, 1*time.Minute)).
			Post("/", teams.CreateTeamHandler(log, workspaceTeamsClient))

		r.With(httprate.LimitByIP(60, 1*time.Minute)).
			Get("/", teams.ListTeamsHandler(log, workspaceTeamsClient))

		// Get / Update / Delete
		r.Route("/{team_id}", func(r chi.Router) {
			r.With(httprate.LimitByIP(120, 1*time.Minute)).
				Get("/", teams.GetTeamHandler(log, workspaceTeamsClient))

			r.With(allowJSON, httprate.LimitByIP(30, 1*time.Minute)).
				Patch("/", teams.UpdateTeamHandler(log, workspaceTeamsClient))

			r.With(httprate.LimitByIP(10, 1*time.Minute)).
				Delete("/", teams.DeleteTeamHandler(log, workspaceTeamsClient))

			// members
			r.Route("/members", func(r chi.Router) {
				r.With(httprate.LimitByIP(120, 1*time.Minute)).
					Get("/", teams.ListTeamMembersHandler(log, workspaceTeamsClient))

				// update duties / remove member
				r.With(allowJSON, httprate.LimitByIP(30, 1*time.Minute)).
					Patch("/{user_id}", teams.UpdateTeamMemberHandler(log, workspaceTeamsClient))

				r.With(httprate.LimitByIP(20, 1*time.Minute)).
					Delete("/{user_id}", teams.RemoveTeamMemberHandler(log, workspaceTeamsClient))
			})
		})
	})

	r.Route("/projects", func(r chi.Router) {
		r.Use(helpers.RequireBearerAuth)

		// CreateProject
		r.With(allowJSON, httprate.LimitByIP(20, 1*time.Minute)).
			Post("/", projects.CreateProjectHandler(log, workspaceProjectsClient))

		// ListProjects (мои / команды / фильтры)
		r.With(httprate.LimitByIP(60, 1*time.Minute)).
			Get("/", projects.ListProjectsHandler(log, workspaceProjectsClient))

		// Public list/search  для авторизованных
		r.With(httprate.LimitByIP(120, 1*time.Minute)).
			Get("/public", projects.ListPublicProjectsHandler(log, workspaceProjectsClient))

		r.With(httprate.LimitByIP(60, 1*time.Minute)).
			Get("/join-requests/manageable", projects.ListManageableProjectJoinRequestBucketsHandler(log, workspaceProjectsClient))

		// Project by id
		r.Route("/{project_id}", func(r chi.Router) {
			// GetProject
			r.With(httprate.LimitByIP(120, 1*time.Minute)).
				Get("/", projects.GetProjectHandler(log, workspaceProjectsClient))

			// UpdateProject
			r.With(allowJSON, httprate.LimitByIP(30, 1*time.Minute)).
				Patch("/", projects.UpdateProjectHandler(log, workspaceProjectsClient))

			// DeleteProject
			r.With(httprate.LimitByIP(10, 1*time.Minute)).
				Delete("/", projects.DeleteProjectHandler(log, workspaceProjectsClient))

			r.With(allowJSON, httprate.LimitByIP(30, 1*time.Minute)).
				Post("/open", projects.SetProjectOpenHandler(log, workspaceProjectsClient))

			// Members
			r.Route("/members", func(r chi.Router) {
				r.With(httprate.LimitByIP(120, 1*time.Minute)).
					Get("/", projects.ListProjectMembersHandler(log, workspaceProjectsClient))

				// add member
				r.With(allowJSON, httprate.LimitByIP(20, 1*time.Minute)).
					Post("/", projects.AddProjectMemberHandler(log, workspaceProjectsClient))

				// remove member
				r.With(httprate.LimitByIP(20, 1*time.Minute)).
					Delete("/{user_id}", projects.RemoveProjectMemberHandler(log, workspaceProjectsClient))

				// rights patch
				r.With(allowJSON, httprate.LimitByIP(20, 1*time.Minute)).
					Patch("/{user_id}/rights", projects.UpdateProjectMemberRightsHandler(log, workspaceProjectsClient))
			})

			// Join requests
			r.Route("/join-requests", func(r chi.Router) {
				// create request (user)
				r.With(allowJSON, httprate.LimitByIP(10, 1*time.Minute)).
					Post("/", projects.RequestJoinProjectHandler(log, workspaceProjectsClient))

				// list requests (managers)
				r.With(httprate.LimitByIP(60, 1*time.Minute)).
					Get("/", projects.ListProjectJoinRequestsHandler(log, workspaceProjectsClient))

				// approve/reject (managers)
				r.With(allowJSON, httprate.LimitByIP(20, 1*time.Minute)).
					Post("/{request_id}/approve", projects.ApproveProjectJoinRequestHandler(log, workspaceProjectsClient))

				r.With(httprate.LimitByIP(60, 1*time.Minute)).
					Get("/details", projects.ListProjectJoinRequestDetailsHandler(log, workspaceProjectsClient))

				r.With(httprate.LimitByIP(60, 1*time.Minute)).
					Get("/me", projects.GetMyProjectJoinRequestHandler(log, workspaceProjectsClient))

				r.With(allowJSON, httprate.LimitByIP(20, 1*time.Minute)).
					Post("/{request_id}/reject", projects.RejectProjectJoinRequestHandler(log, workspaceProjectsClient))

			})
		})
	})

	return r
}

func allowJSON(next http.Handler) http.Handler {
	return middleware.AllowContentType("application/json")(next)
}
