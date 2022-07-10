package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/config"
	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/dashboard"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/jackc/pgx/v4/pgxpool"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/jwt"
	mw "github.com/PA-D3RPLA/d3if43-htt-uhomestay/middleware"
)

type RestApiConf struct {
	BuildDate     string
	CommitHash    string
	Conf          config.Config
	PosgrePool    *pgxpool.Pool
	DashboardDeps *dashboard.DashboardDeps
}

func NewRestApi(
	buildDate string,
	commitHash string,
	conf config.Config,
	posgrePool *pgxpool.Pool,
	dashboardDeps *dashboard.DashboardDeps,
) *RestApiConf {
	return &RestApiConf{
		BuildDate:     buildDate,
		CommitHash:    commitHash,
		Conf:          conf,
		PosgrePool:    posgrePool,
		DashboardDeps: dashboardDeps,
	}
}

func (p *RestApiConf) RestApiHandler() {
	sentryHandler := sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})
	jwtMidd := jwt.NewMiddleware(p.Conf.JwtKey, p.Conf.JwtIssuerUrl, p.Conf.JwtAudiences, &jwt.JwtPrivateClaim{})
	adminJwtMidd := jwt.NewMiddleware(p.Conf.JwtKey, p.Conf.JwtIssuerUrl, p.Conf.JwtAudiences, &jwt.JwtPrivateAdminClaim{})
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

	r.Get("/registerform", p.DashboardDeps.RegisterForm)
	r.Get("/positionform", p.DashboardDeps.PositionForm)
	r.Get("/periodform", p.DashboardDeps.PeriodForm)
	r.Post("/positions", p.DashboardDeps.PostPosition)
	r.With(trxMidd).Post("/periods", p.DashboardDeps.PostPeriod)
	r.With(trxMidd).Post("/members", p.DashboardDeps.PostMember)

	r.With(trxMidd).Post("/api/v1/register", p.DashboardDeps.PostRegisterMember)
	r.Post("/api/v1/login/members", p.DashboardDeps.PostLoginMember)
	r.Post("/api/v1/login/admins", p.DashboardDeps.PostLoginAdmin)

	r.Get("/api/v1/members", p.DashboardDeps.GetMembers)
	r.Get("/api/v1/members/{id}", p.DashboardDeps.GetMember)
	r.With(jwtMidd).Get("/api/v1/profile", p.DashboardDeps.GetProfileMember)
	r.With(adminJwtMidd).With(trxMidd).Post("/api/v1/members", p.DashboardDeps.PostMember)
	r.With(jwtMidd).With(trxMidd).Put("/api/v1/members", p.DashboardDeps.PutMemberProfile)
	r.With(adminJwtMidd).With(trxMidd).Put("/api/v1/members/{id}", p.DashboardDeps.PutMember)
	r.With(adminJwtMidd).With(trxMidd).Delete("/api/v1/members/{id}", p.DashboardDeps.DeleteMember)
	r.With(adminJwtMidd).With(trxMidd).Patch("/api/v1/members/{id}", p.DashboardDeps.PatchMemberApproval)

	r.Get("/api/v1/periods", p.DashboardDeps.GetPeriods)
	r.Get("/api/v1/periods/active", p.DashboardDeps.GetActivePeriod)
	r.Get("/api/v1/periods/{id}/structures", p.DashboardDeps.GetPeriodStructure)
	r.With(adminJwtMidd).With(trxMidd).Post("/api/v1/periods", p.DashboardDeps.PostPeriod)
	r.With(adminJwtMidd).Post("/api/v1/periods/goals", p.DashboardDeps.PostGoal)
	r.With(adminJwtMidd).With(trxMidd).Put("/api/v1/periods/{id}", p.DashboardDeps.PutPeriod)
	r.With(adminJwtMidd).With(trxMidd).Delete("/api/v1/periods/{id}", p.DashboardDeps.DeletePeriod)
	r.With(adminJwtMidd).With(trxMidd).Patch("/api/v1/periods/{id}/status", p.DashboardDeps.PatchPeriodStatus)
	r.Get("/api/v1/periods/{id}/goal", p.DashboardDeps.GetOrgPeriodGoal)

	r.Get("/api/v1/positions", p.DashboardDeps.GetPositions)
	r.Get("/api/v1/positions/levels", p.DashboardDeps.GetPositionLevels)
	r.With(adminJwtMidd).Post("/api/v1/positions", p.DashboardDeps.PostPosition)
	r.With(adminJwtMidd).With(trxMidd).Put("/api/v1/positions/{id}", p.DashboardDeps.PutPositions)
	r.With(adminJwtMidd).With(trxMidd).Delete("/api/v1/positions/{id}", p.DashboardDeps.DeletePosition)

	r.Get("/api/v1/documents", p.DashboardDeps.GetDocuments)
	r.With(adminJwtMidd).Post("/api/v1/documents/dir", p.DashboardDeps.PostDirDocument)
	r.With(adminJwtMidd).Post("/api/v1/documents/file", p.DashboardDeps.PostFileDocument)
	r.With(adminJwtMidd).With(trxMidd).Put("/api/v1/documents/dir/{id}", p.DashboardDeps.PutDirDocument)
	r.With(adminJwtMidd).With(trxMidd).Put("/api/v1/documents/file/{id}", p.DashboardDeps.PutFileDocument)
	r.Get("/api/v1/documents/{id}", p.DashboardDeps.GetDocumentChildren)
	r.With(adminJwtMidd).With(trxMidd).Delete("/api/v1/documents/{id}", p.DashboardDeps.DeleteDocument)

	r.With(adminJwtMidd).Post("/api/v1/histories", p.DashboardDeps.PostHistory)
	r.Get("/api/v1/histories", p.DashboardDeps.GetHistory)

	r.Get("/api/v1/blogs", p.DashboardDeps.GetBlogs)
	r.Get("/api/v1/blogs/{id}", p.DashboardDeps.GetBlog)
	r.With(adminJwtMidd).Post("/api/v1/blogs", p.DashboardDeps.PostBlog)
	r.With(adminJwtMidd).Put("/api/v1/blogs/{id}", p.DashboardDeps.PutBlogs)
	r.With(adminJwtMidd).Delete("/api/v1/blogs/{id}", p.DashboardDeps.DeleteBlog)
	r.With(adminJwtMidd).Post("/api/v1/blogs/image", p.DashboardDeps.PostImage)

	r.Get("/api/v1/cashflows", p.DashboardDeps.GetCashflows)
	r.With(adminJwtMidd).Post("/api/v1/cashflows", p.DashboardDeps.PostCashflow)
	r.With(adminJwtMidd).Put("/api/v1/cashflows/{id}", p.DashboardDeps.PutCashflow)
	r.With(adminJwtMidd).Delete("/api/v1/cashflows/{id}", p.DashboardDeps.DeleteCashflow)

	r.With(adminJwtMidd).Put("/api/v1/dues/members/monthly/{id}", p.DashboardDeps.PutMemberDues)
	r.With(adminJwtMidd).Patch("/api/v1/dues/members/monthly/{id}", p.DashboardDeps.PatchMemberDues)
	r.With(jwtMidd).Post("/api/v1/dues/members/monthly/{id}", p.DashboardDeps.PostMemberDues)
	r.Get("/api/v1/dues/members/{id}", p.DashboardDeps.GetMemberDues)
	r.Get("/api/v1/dues/{id}/members", p.DashboardDeps.GetMembersDues)

	r.Get("/api/v1/dues", p.DashboardDeps.GetDues)
	r.With(adminJwtMidd).Post("/api/v1/dues", p.DashboardDeps.PostDues)
	r.Get("/api/v1/dues/{id}/check", p.DashboardDeps.GetPaidDues)
	r.With(adminJwtMidd).Put("/api/v1/dues/{id}", p.DashboardDeps.PutDues)
	r.With(adminJwtMidd).Delete("/api/v1/dues/{id}", p.DashboardDeps.DeleteDues)

	r.Get("/api/v1/dashboard", p.DashboardDeps.GetPublicDashboard)
	r.With(adminJwtMidd).Get("/api/v1/dashboard/private", p.DashboardDeps.GetPrivateDashboard)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "docs"))
	ChiFileServer(r, "/docs", filesDir)

	http.ListenAndServe(fmt.Sprintf(":%s", p.Conf.Port), r)
}

// Ref:
// https://github.com/go-chi/chi/blob/master/_examples/fileserver/main.go
//
// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func ChiFileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
