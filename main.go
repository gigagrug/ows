package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"ows/api/auth"
	"ows/api/logging"
	"ows/api/project"
	"ows/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	auth := &auth.DB{DB: db}
	mux.HandleFunc("POST /api/user/register/{$}", auth.Register)
	mux.HandleFunc("POST /api/user/login/{$}", auth.Login)
	mux.HandleFunc("POST /api/user/logout/{$}", auth.Logout)

	project := &project.DB{DB: db}
	mux.HandleFunc("GET /api/project/getProjects/{$}", Chain(project.GetProjects, AuthCheck(db)))
	mux.HandleFunc("GET /api/project/getProject/{projectId}/{$}", Chain(project.GetProject, AuthCheck(db)))
	mux.HandleFunc("POST /api/project/createProject/{$}", Chain(project.CreateProjects, AuthCheck(db)))
	mux.HandleFunc("PUT /api/project/updateProject/{projectId}/{$}", Chain(project.UpdateProjects, AuthCheck(db)))
	mux.HandleFunc("DELETE /api/project/deleteProject/{projectId}/{$}", Chain(project.DeleteProjects, AuthCheck(db)))

	logging := &logging.DB{DB: db}
	// Logging Service
	mux.HandleFunc("GET /api/logService/getLogServices/{$}", Chain(logging.GetLogServices, AuthCheck(db)))
	mux.HandleFunc("GET /api/logService/getLogService/{logServiceId}/{$}", Chain(logging.GetLogService, AuthCheck(db)))
	mux.HandleFunc("POST /api/logService/createLogService/{$}", Chain(logging.CreateLogService, AuthCheck(db)))
	mux.HandleFunc("PUT /api/logService/updateLogService/{logServiceId}/{$}", Chain(logging.PutLogService, AuthCheck(db)))
	mux.HandleFunc("DELETE /api/logService/deleteLogService/{logServiceId}/{$}", Chain(logging.DeleteLogService, AuthCheck(db)))
	// Logs
	mux.HandleFunc("GET /api/log/getLogs/{$}", Chain(logging.GetLogs, AuthCheck(db)))
	mux.HandleFunc("POST /api/log/createLog/{$}", Chain(logging.CreateLog, AuthCheck(db)))
	mux.HandleFunc("DELETE /api/log/deleteLog/{logId}/{$}", Chain(logging.DeleteLog, AuthCheck(db)))

	if err := http.ListenAndServe(os.Getenv("PORT"), middleware(mux)); err != nil {
		log.Fatal(err)
	}
}

func middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		if r.Method == "OPTIONS" {
			return
		}
		fmt.Println(r.Method, r.URL.Path, time.Since(time.Now()))
		h.ServeHTTP(w, r)
	})
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

func AuthCheck(db *sql.DB) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := r.Cookie("Authorization")
			if err != nil {
				slog.Error(err.Error())
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(tokenString.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(os.Getenv("SECRET")), nil
			})
			if err != nil {
				slog.Error(err.Error())
				http.Error(w, "jwt not working", http.StatusUnauthorized)
				return
			}
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				if float64(time.Now().Unix()) > claims["exp"].(float64) {
					slog.Error("Token has expired")
					http.Error(w, "Token has expired", http.StatusUnauthorized)
					return
				}

				var userS model.User
				if err := db.QueryRow(`SELECT id FROM "User" WHERE id = $1`, claims["sub"].(string)).Scan(&userS.ID); err != nil {
					slog.Error(err.Error())
					http.Error(w, "error retrieving user from the database", http.StatusInternalServerError)
					return
				}

				ctx := r.Context()
				ctx = context.WithValue(ctx, "key", userS.ID)
				r = r.WithContext(ctx)
			} else {
				slog.Error("Invalid token")
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			f(w, r)
			slog.Info("Middleware works")
		}
	}
}

func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}
