package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

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

			uidStr, gidStr, name, group, err := validateJWT(serverConfig, authToken)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			uid, err := strconv.Atoi(uidStr)
			if err != nil {
				http.Error(w, fmt.Sprintf("invalid string value for uid : %s , could not convert to int", uidStr), http.StatusBadRequest)
			}
			gid, err := strconv.Atoi(gidStr)
			if err != nil {
				http.Error(w, fmt.Sprintf("invalid string value for uid : %s , could not convert to int", uidStr), http.StatusBadRequest)
			}

			ctxData := types.ContextDataType{
				Uid:   uid,
				Gid:   gid,
				Name:  name,
				Group: group,
			}

			ctx := context.WithValue(r.Context(), serverConfig.ContextKey, &ctxData)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
