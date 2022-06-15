package handler_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/PA-D3RPLA/d3if43-htt-uhomestay/jwt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"
	mw "github.com/PA-D3RPLA/d3if43-htt-uhomestay/middleware"

	"github.com/go-chi/chi/v5"
)

var db *pgxpool.Pool

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.1",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		dbConfig, err := pgxpool.ParseConfig(databaseUrl)
		if err != nil {
			return err
		}

		db, err = pgxpool.ConnectConfig(context.Background(), dbConfig)
		if err != nil {
			return err
		}

		return db.Ping(context.Background())
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	// Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestPrivateJWTRoute(t *testing.T) {
	jwtKey := []byte("test")
	jwtIssuerUrl := "http://localhost:8080"
	jwtAudiences := []string{"test"}

	jwtMidd := jwt.NewMiddleware(jwtKey, jwtIssuerUrl, jwtAudiences, &jwt.JwtPrivateClaim{})

	testCases := []struct {
		name               string
		setHeader          func(r *http.Request)
		expectedStatusCode int
	}{
		{
			name: "Access Private Route Success",
			setHeader: func(r *http.Request) {
				jwtToken, _ := jwt.Sign(
					"",
					"token",
					jwtIssuerUrl,
					jwtKey,
					jwtAudiences,
					time.Time{},
					time.Now().Add(time.Hour),
					time.Time{},
					jwt.JwtPrivateClaim{
						Uid: "12345678-1234-1234-1234-123456789012",
					})
				r.Header.Set("Authorization", "Bearer "+jwtToken)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Access Private Route Fail, Wrong JWT Key",
			setHeader: func(r *http.Request) {
				jwtToken, _ := jwt.Sign(
					"",
					"token",
					jwtIssuerUrl,
					[]byte("wrong-key"),
					jwtAudiences,
					time.Time{},
					time.Now().Add(time.Hour),
					time.Time{},
					jwt.JwtPrivateClaim{
						Uid: "12345678-1234-1234-1234-123456789012",
					})
				r.Header.Set("Authorization", "Bearer "+jwtToken)
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "Access Private Route Fail, Not Audience",
			setHeader: func(r *http.Request) {
				jwtToken, _ := jwt.Sign(
					"",
					"token",
					jwtIssuerUrl,
					jwtKey,
					[]string{"stranger"},
					time.Time{},
					time.Now().Add(time.Hour),
					time.Time{},
					jwt.JwtPrivateClaim{
						Uid: "12345678-1234-1234-1234-123456789012",
					})
				r.Header.Set("Authorization", "Bearer "+jwtToken)
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "Access Private Route Fail, Token Expired",
			setHeader: func(r *http.Request) {
				jwtToken, _ := jwt.Sign(
					"",
					"token",
					jwtIssuerUrl,
					jwtKey,
					jwtAudiences,
					time.Time{},
					time.Now().Add(-time.Hour),
					time.Time{},
					jwt.JwtPrivateClaim{
						Uid: "12345678-1234-1234-1234-123456789012",
					})
				r.Header.Set("Authorization", "Bearer "+jwtToken)
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "Access Private Route Fail, Not JWT Token",
			setHeader: func(r *http.Request) {
				r.Header.Set("Authorization", "Brear bla-bla-bla")
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Access Private Route Fail, Not Bearer Token",
			setHeader: func(r *http.Request) {
				r.Header.Set("Authorization", "bla-bla-bla")
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Access Private Route Fail, Authorization Header not Provided",
			setHeader: func(r *http.Request) {
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	r := chi.NewRouter()
	r.With(jwtMidd).Get("/private", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/private", nil)
			c.setHeader(req)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			body, err := io.ReadAll(rr.Body)
			if err != nil {
				log.Fatal(err)
			}

			if rr.Code != c.expectedStatusCode {
				t.Logf("%s", body)
				t.Fatalf("Expected response code %d. Got %d\n", c.expectedStatusCode, rr.Code)
			}
		})
	}
}

func TestPrivateJWTAdminRoute(t *testing.T) {
	jwtKey := []byte("test")
	jwtIssuerUrl := "http://localhost:8080"
	jwtAudiences := []string{"test"}

	jwtMidd := jwt.NewMiddleware(jwtKey, jwtIssuerUrl, jwtAudiences, &jwt.JwtPrivateAdminClaim{})

	testCases := []struct {
		name               string
		setHeader          func(r *http.Request)
		expectedStatusCode int
	}{
		{
			name: "Access Private Route Success",
			setHeader: func(r *http.Request) {
				jwtToken, _ := jwt.Sign(
					"",
					"token",
					jwtIssuerUrl,
					jwtKey,
					jwtAudiences,
					time.Time{},
					time.Now().Add(time.Hour),
					time.Time{},
					jwt.JwtPrivateAdminClaim{
						Uid:     "12345678-1234-1234-1234-123456789012",
						IsAdmin: true,
					})
				r.Header.Set("Authorization", "Bearer "+jwtToken)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Access Private Route Fail, JWT is not Admin JWT",
			setHeader: func(r *http.Request) {
				jwtToken, _ := jwt.Sign(
					"",
					"token",
					jwtIssuerUrl,
					[]byte("wrong-key"),
					jwtAudiences,
					time.Time{},
					time.Now().Add(time.Hour),
					time.Time{},
					jwt.JwtPrivateAdminClaim{
						Uid:     "12345678-1234-1234-1234-123456789012",
						IsAdmin: false,
					})
				r.Header.Set("Authorization", "Bearer "+jwtToken)
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "Access Private Route Fail, Wrong JWT Key",
			setHeader: func(r *http.Request) {
				jwtToken, _ := jwt.Sign(
					"",
					"token",
					jwtIssuerUrl,
					[]byte("wrong-key"),
					jwtAudiences,
					time.Time{},
					time.Now().Add(time.Hour),
					time.Time{},
					jwt.JwtPrivateAdminClaim{
						Uid:     "12345678-1234-1234-1234-123456789012",
						IsAdmin: true,
					})
				r.Header.Set("Authorization", "Bearer "+jwtToken)
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "Access Private Route Fail, Not Audience",
			setHeader: func(r *http.Request) {
				jwtToken, _ := jwt.Sign(
					"",
					"token",
					jwtIssuerUrl,
					jwtKey,
					[]string{"stranger"},
					time.Time{},
					time.Now().Add(time.Hour),
					time.Time{},
					jwt.JwtPrivateAdminClaim{
						Uid:     "12345678-1234-1234-1234-123456789012",
						IsAdmin: true,
					})
				r.Header.Set("Authorization", "Bearer "+jwtToken)
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "Access Private Route Fail, Token Expired",
			setHeader: func(r *http.Request) {
				jwtToken, _ := jwt.Sign(
					"",
					"token",
					jwtIssuerUrl,
					jwtKey,
					jwtAudiences,
					time.Time{},
					time.Now().Add(-time.Hour),
					time.Time{},
					jwt.JwtPrivateAdminClaim{
						Uid:     "12345678-1234-1234-1234-123456789012",
						IsAdmin: true,
					})
				r.Header.Set("Authorization", "Bearer "+jwtToken)
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name: "Access Private Route Fail, Not JWT Token",
			setHeader: func(r *http.Request) {
				r.Header.Set("Authorization", "Brear bla-bla-bla")
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Access Private Route Fail, Not Bearer Token",
			setHeader: func(r *http.Request) {
				r.Header.Set("Authorization", "bla-bla-bla")
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Access Private Route Fail, Authorization Header not Provided",
			setHeader: func(r *http.Request) {
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	r := chi.NewRouter()
	r.With(jwtMidd).Get("/private", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/private", nil)
			c.setHeader(req)

			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)

			body, err := io.ReadAll(rr.Body)
			if err != nil {
				log.Fatal(err)
			}

			if rr.Code != c.expectedStatusCode {
				t.Logf("%s", body)
				t.Fatalf("Expected response code %d. Got %d\n", c.expectedStatusCode, rr.Code)
			}
		})
	}
}

func TestSomething(t *testing.T) {
	trxMidd := mw.NewTrxMiddleware(db)

	r := chi.NewRouter()
	r.With(trxMidd).Get("/trx", func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value(arbitary.TrxX{}).(pgx.Tx)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	req, _ := http.NewRequest("GET", "/trx", nil)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatal("Expected response code")
	}
}
