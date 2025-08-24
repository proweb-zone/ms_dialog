package middleware

import (
	"log"
	"ms_dialog/internal/app/repository"
	"ms_dialog/internal/config"
	"ms_dialog/internal/db/postgres"
	"net/http"
	"strings"
)

func CheckAccess(config *config.Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Извлечение токена из заголовка Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header missing", http.StatusUnauthorized)
				return
			}

			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			token := bearerToken[1]

			//db := postgres.Connect(config)
			masterURL := []string{"user=postgres password=yourpassword dbname=master sslmode=disable"}
			slaveURLs := []string{
				"user=postgres password=yourpassword dbname=slave1 sslmode=disable",
				"user=postgres password=yourpassword dbname=slave2 sslmode=disable",
			}

			dataSource, err := postgres.NewReplicationRoutingDataSource(masterURL, slaveURLs, true)
			if err != nil {
				log.Fatal(err)
			}
			defer dataSource.Close()

			userRepository := repository.InitUserRepository(dataSource)
			auth, _ := userRepository.CheckToken(token)

			if auth != nil && len(auth.Token) > 0 {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "Error: check Token Bearer", http.StatusUnauthorized)
				return
			}

		})
	}
}
