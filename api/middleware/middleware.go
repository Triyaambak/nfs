package middleware

import (
	"context"
	"net/http"

	types "github.com/Triyaambak/nfs/types"
)

func AuthMiddle(serverConfig *types.ServerConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authToken, err := getAuthToken(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			uid, gid, err := validateJWT(serverConfig, authToken)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "uid", uid)
			ctx = context.WithValue(ctx, "gid", gid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
