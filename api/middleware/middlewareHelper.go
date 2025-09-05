package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	types "github.com/Triyaambak/nfs/types"

	"github.com/golang-jwt/jwt/v5"
)

func getAuthToken(r *http.Request) (authToken string, err error) {
	jwt := r.Header.Get("Authorization")
	if jwt == "" {
		return "", errors.New("missing jwt - auth token in header")
	}

	idx := strings.Index(jwt, " ")
	if idx == -1 {
		return "", errors.New("invalid jwt format")
	}

	authTokenPrefix := jwt[:idx]
	if authTokenPrefix != "Bearer" {
		return "", errors.New("invalid jwt format , No prefix of Bearer")
	}

	authToken = jwt[idx+1:]

	return authToken, nil
}

func validateJWT(serverConfig *types.ServerConfig, authToken string) (uid, gid, name, group string, err error) {

	token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return (*serverConfig).Secret, nil
	})

	if err != nil {
		return "", "", "", "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		okUid := false
		okGid := false
		okName := false
		okGroup := false
		uid, okUid = claims["uid"].(string)
		gid, okGid = claims["gid"].(string)
		name, okName = claims["name"].(string)
		group, okGroup = claims["group"].(string)

		if !okUid || !okGid || !okName || !okGroup {
			return "", "", "", "", errors.New("failed to parse uid ,gid, name , group from token claims")
		}
	}

	return uid, gid, name, group, nil
}
