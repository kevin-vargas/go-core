package middleware

import (
	"net/http"
	"strings"

	"github.com/kevin-vargas/go-core/jwt"
	"github.com/kevin-vargas/go-core/jwt/ctx"
)

type Middleware func(next http.Handler) http.Handler

const (
	bearer = "Bearer "
)

func Auth(j *jwt.JWT) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			token = strings.ReplaceAll(token, bearer, "")
			claim, err := j.Validate(token)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return

			}
			ctx := ctx.WithUsername(r.Context(), claim.Username)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
