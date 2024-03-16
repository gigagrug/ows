package auth

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"ows/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	*sql.DB
}

func (db *DB) Register(w http.ResponseWriter, r *http.Request) {
	var user model.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		slog.Error(err.Error())
		http.Error(w, "error getting json data", http.StatusBadRequest)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 15)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error creating user", http.StatusBadRequest)
	}

	if _, err := db.Exec(`INSERT INTO "User" (email, password) VALUES ($1, $2)`, user.Email, hash); err != nil {
		slog.Error(err.Error())
		http.Error(w, "error creating user", http.StatusInternalServerError)
		return
	}

	slog.Info("Register: ok")
}

func (db *DB) Login(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string
		Password string
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		slog.Error(err.Error())
		http.Error(w, "error getting json data", http.StatusBadRequest)
	}

	var user model.User
	if err := db.QueryRow(`SELECT id, email, password FROM "User" WHERE email = $1`, body.Email).Scan(&user.ID, &user.Email, &user.Password); err != nil {
		slog.Error(err.Error())
		http.Error(w, "error getting user", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		slog.Error(err.Error())
		http.Error(w, "wrong email or password", http.StatusBadRequest)
	}

	expiration := time.Now().Add(time.Hour * 24 * 30)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": expiration.Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "jwt issue", http.StatusBadRequest)
	}
	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Expires:  expiration,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, &cookie)
	slog.Info("Login: ok")
}

func (db *DB) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}

	http.SetCookie(w, &cookie)
	slog.Info("Logout: ok")
}
