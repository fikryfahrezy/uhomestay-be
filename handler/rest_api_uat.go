package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	httpSwagger "github.com/swaggo/http-swagger"

	mw "github.com/PA-D3RPLA/d3if43-htt-uhomestay/middleware"
)

func (p *RestApiConf) RestApiHandlerUat() {
	sentryHandler := sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})
	trxMidd := mw.NewTrxMiddleware(p.PosgrePool)

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	corsMidd := cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	// Enable httprate request limiter of 100 requests per minute.
	//
	// In the code example below, rate-limiting is bound to the request IP address
	// via the LimitByIP middleware handler.
	//
	// To have a single rate-limiter for all requests, use httprate.LimitAll(..).
	//
	// Please see _example/main.go for other more, or read the library code.
	rateLMidd := httprate.LimitByIP(100, 1*time.Minute)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Add pkg or example for go-chi
	// Ref:https://github.com/getsentry/sentry-go/issues/143
	// Important: Chi has a middleware stack and thus it is important to put the
	// Sentry handler on the appropriate place. If using middleware.Recoverer,
	// the Sentry middleware must come afterwards (and configure it with
	// Repanic: true).
	r.Use(sentryHandler.Handle)

	r.Use(rateLMidd)
	r.Use(corsMidd)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("route does not exist"))
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("method is not valid"))
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	})
	r.Get("/info", func(w http.ResponseWriter, r *http.Request) {
		info := "Built: " + p.BuildDate + ", Commit: " + p.CommitHash
		w.Write([]byte(info))
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/docs/swagger.yaml"), // The url pointing to API definition
	))

	r.Get("/registerform", p.DashboardDeps.RegisterFormUat)
	r.Get("/positionform", p.DashboardDeps.PositionFormUat)
	r.Get("/periodform", p.DashboardDeps.PeriodFormUat)
	r.Post("/positions", p.DashboardDeps.PostPositionUat)
	r.With(trxMidd).Post("/periods", p.DashboardDeps.PostPeriodUat)
	r.With(trxMidd).Post("/members", p.DashboardDeps.PostMemberUat)

	r.With(trxMidd).Post("/api/v1/register", p.DashboardDeps.PostRegisterMemberUat)
	r.Post("/api/v1/login/members", p.DashboardDeps.PostLoginMemberUat)
	r.Post("/api/v1/login/admins", p.DashboardDeps.PostLoginAdminUat)

	r.Get("/api/v1/members", p.DashboardDeps.GetMembersUat)
	r.Get("/api/v1/members/{id}", p.DashboardDeps.GetMemberUat)
	r.Get("/api/v1/profile", p.DashboardDeps.GetProfileMemberUat)
	r.With(trxMidd).Post("/api/v1/members", p.DashboardDeps.PostMemberUat)
	r.With(trxMidd).Put("/api/v1/members", p.DashboardDeps.PutMemberProfileUat)
	r.With(trxMidd).Put("/api/v1/members/{id}", p.DashboardDeps.PutMemberUat)
	r.With(trxMidd).Delete("/api/v1/members/{id}", p.DashboardDeps.DeleteMemberUat)
	r.With(trxMidd).Patch("/api/v1/members/{id}", p.DashboardDeps.PatchMemberApprovalUat)

	r.Get("/api/v1/periods", p.DashboardDeps.GetPeriodsUat)
	r.Get("/api/v1/periods/active", p.DashboardDeps.GetActivePeriodUat)
	r.Get("/api/v1/periods/{id}/structures", p.DashboardDeps.GetPeriodStructureUat)
	r.With(trxMidd).Post("/api/v1/periods", p.DashboardDeps.PostPeriodUat)
	r.Post("/api/v1/periods/goals", p.DashboardDeps.PostGoalUat)
	r.With(trxMidd).Put("/api/v1/periods/{id}", p.DashboardDeps.PutPeriodUat)
	r.With(trxMidd).Delete("/api/v1/periods/{id}", p.DashboardDeps.DeletePeriodUat)
	r.With(trxMidd).Patch("/api/v1/periods/{id}/status", p.DashboardDeps.PatchPeriodStatusUat)
	r.Get("/api/v1/periods/{id}/goal", p.DashboardDeps.GetOrgPeriodGoalUat)

	r.Get("/api/v1/positions", p.DashboardDeps.GetPositionsUat)
	r.Get("/api/v1/positions/levels", p.DashboardDeps.GetPositionLevelsUat)
	r.Post("/api/v1/positions", p.DashboardDeps.PostPositionUat)
	r.With(trxMidd).Put("/api/v1/positions/{id}", p.DashboardDeps.PutPositionsUat)
	r.With(trxMidd).Delete("/api/v1/positions/{id}", p.DashboardDeps.DeletePositionUat)

	r.Get("/api/v1/documents", p.DashboardDeps.GetDocumentsUat)
	r.Post("/api/v1/documents/dir", p.DashboardDeps.PostDirDocumentUat)
	r.Post("/api/v1/documents/file", p.DashboardDeps.PostFileDocumentUat)
	r.With(trxMidd).Put("/api/v1/documents/dir/{id}", p.DashboardDeps.PutDirDocumentUat)
	r.With(trxMidd).Put("/api/v1/documents/file/{id}", p.DashboardDeps.PutFileDocumentUat)
	r.Get("/api/v1/documents/{id}", p.DashboardDeps.GetDocumentChildrenUat)
	r.With(trxMidd).Delete("/api/v1/documents/{id}", p.DashboardDeps.DeleteDocumentUat)

	r.Post("/api/v1/histories", p.DashboardDeps.PostHistoryUat)
	r.Get("/api/v1/histories", p.DashboardDeps.GetHistoryUat)

	r.Get("/api/v1/blogs", p.DashboardDeps.GetBlogsUat)
	r.Get("/api/v1/blogs/{id}", p.DashboardDeps.GetBlogUat)
	r.Post("/api/v1/blogs", p.DashboardDeps.PostBlogUat)
	r.Put("/api/v1/blogs/{id}", p.DashboardDeps.PutBlogsUat)
	r.Delete("/api/v1/blogs/{id}", p.DashboardDeps.DeleteBlogUat)
	r.Post("/api/v1/blogs/image", p.DashboardDeps.PostImageUat)

	r.Get("/api/v1/cashflows", p.DashboardDeps.GetCashflowsUat)
	r.Post("/api/v1/cashflows", p.DashboardDeps.PostCashflowUat)
	r.Put("/api/v1/cashflows/{id}", p.DashboardDeps.PutCashflowUat)
	r.Delete("/api/v1/cashflows/{id}", p.DashboardDeps.DeleteCashflowUat)

	r.Put("/api/v1/dues/members/monthly/{id}", p.DashboardDeps.PutMemberDuesUat)
	r.Patch("/api/v1/dues/members/monthly/{id}", p.DashboardDeps.PatchMemberDuesUat)
	r.Post("/api/v1/dues/members/monthly/{id}", p.DashboardDeps.PostMemberDuesUat)
	r.Get("/api/v1/dues/members/{id}", p.DashboardDeps.GetMemberDuesUat)
	r.Get("/api/v1/dues/{id}/members", p.DashboardDeps.GetMembersDuesUat)

	r.Get("/api/v1/dues", p.DashboardDeps.GetDuesUat)
	r.Post("/api/v1/dues", p.DashboardDeps.PostDuesUat)
	r.Get("/api/v1/dues/{id}/check", p.DashboardDeps.GetPaidDuesUat)
	r.Put("/api/v1/dues/{id}", p.DashboardDeps.PutDuesUat)
	r.Delete("/api/v1/dues/{id}", p.DashboardDeps.DeleteDuesUat)

	r.Get("/api/v1/dashboard", p.DashboardDeps.GetPublicDashboardUat)
	r.Get("/api/v1/dashboard/private", p.DashboardDeps.GetPrivateDashboardUat)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "docs"))
	ChiFileServer(r, "/docs", filesDir)

	http.ListenAndServe(fmt.Sprintf(":%s", p.Conf.Port), r)
}
