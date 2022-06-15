package middleware

import (
	"context"
	"net/http"

	arbitary "github.com/PA-D3RPLA/d3if43-htt-uhomestay/arbitrary"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Ref:
// https://gist.github.com/Boerworz/b683e46ae0761056a636
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// Ref:
// https://t.me/golangID/146225
// https://t.me/golangID/150320
// HTTP middleware setting a value on the request context
func NewTrxMiddleware(db *pgxpool.Pool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tx, err := db.Begin(context.Background())
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("error in tx middleware"))
				return
			}

			defer tx.Rollback(context.Background())

			// create new context from `r` request context, and assign key `"TrxX{}"`
			// to value of `"pgx.Tx"`
			ctx := context.WithValue(r.Context(), arbitary.TrxX{}, tx)

			// Capture status code from handler
			lrw := NewLoggingResponseWriter(w)

			// call the next handler in the chain, passing the response writer and
			// the updated request object with the new context value.
			//
			// note: context.Context values are nested, so any previously set
			// values will be accessible as well, and the new `"TrxX{}"` key
			// will be accessible from this point forward.
			next.ServeHTTP(lrw, r.WithContext(ctx))

			if lrw.statusCode < http.StatusBadRequest {
				err = tx.Commit(context.Background())
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("error in final commit transaction"))
					return
				}
			}
		})
	}
}
