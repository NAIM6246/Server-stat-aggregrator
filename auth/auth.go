package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/naim6246/Server-stat-aggregrator/configs"
	"github.com/naim6246/Server-stat-aggregrator/models"
)

const (
	KeyUserID       = "id"
	KeyTokenExpired = "exp"
	KeyUser         = "user"
)

func GenerateToken(user *models.User) (string, error) {
	cfg := configs.GetAppConfig()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims[KeyUserID] = user.AccesID
	claims[KeyTokenExpired] = time.Now().Add(time.Hour * time.Duration(cfg.AccesInTTLHour)).Unix()
	signedToken, err := token.SignedString([]byte(cfg.Secret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := configs.GetAppConfig()
		extractToken := func() string {
			bearerToken := r.Header.Get("Authorization")
			strArr := strings.Split(bearerToken, " ")
			if len(strArr) == 2 {
				return strArr[1]
			}
			return ""
		}
		if r.Header["Authorization"] != nil {
			tokenString := extractToken()
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("invalid authorization token")
				}
				return []byte(cfg.Secret), nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(err)
				return
			}
			if token.Valid {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(errors.New("invalid authorization token"))
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(errors.New("an authorization header is required"))
		}
	})
}
