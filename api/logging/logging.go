package logging

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"ows/model"
	"time"

	"github.com/go-playground/validator/v10"
)

type DB struct {
	*sql.DB
}

// Logging Service
func (db *DB) GetLogServices(w http.ResponseWriter, r *http.Request) {
	var logService model.LogService
	user := r.Context().Value("key").(string)

	rows, err := db.Query(`SELECT id, name, created_at, updated_at, user_id, project_id FROM "LogService" WHERE user_id = $1 ORDER BY created_at ASC`, user)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error getting projects", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	logServices := []model.LogService{}
	for rows.Next() {
		if err := rows.Scan(&logService.ID, &logService.Name, &logService.CreatedAt, &logService.UpdatedAt, &logService.UserID, &logService.ProjectID); err != nil {
			slog.Error(err.Error())
			http.Error(w, "error getting project", http.StatusInternalServerError)
			return
		}
		logServices = append(logServices, logService)
	}

	jsonData, err := json.Marshal(logServices)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error getting projects", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonData)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error sending projects", http.StatusInternalServerError)
		return
	}

	slog.Info("get log services")
}

func (db *DB) GetLogService(w http.ResponseWriter, r *http.Request) {
	var logService model.LogService
	id := r.PathValue("logServiceId")
	user := r.Context().Value("key").(string)

	row := db.QueryRow(`SELECT id, name, created_at, updated_at, project_id FROM "LogService" WHERE id = $1 AND user_id = $2`, id, user)
	if err := row.Scan(&logService.ID, &logService.Name, &logService.CreatedAt, &logService.UpdatedAt, &logService.ProjectID); err != nil {
		slog.Error(err.Error())
		http.Error(w, "error getting project", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(logService)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error getting project", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonData)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error sending project", http.StatusInternalServerError)
		return
	}

	slog.Info("get log service")
}

func (db *DB) CreateLogService(w http.ResponseWriter, r *http.Request) {
	var logService model.LogService
	user := r.Context().Value("key").(string)

	if err := json.NewDecoder(r.Body).Decode(&logService); err != nil {
		slog.Error(err.Error())
		http.Error(w, "error getting json data", http.StatusBadRequest)
	}

	validate := validator.New()
	if err := validate.Struct(&logService); err != nil {
		slog.Error(err.Error())
		http.Error(w, "name not inputed", http.StatusInternalServerError)
		return
	}

	if _, err := db.Exec(`INSERT INTO "LogService" (name, updated_at, user_id, project_id) VALUES ($1, $2, $3, $4)`, logService.Name, time.Now(), user, logService.ProjectID); err != nil {
		slog.Error(err.Error())
		http.Error(w, "error creating project", http.StatusInternalServerError)
		return
	}

	slog.Info("create log services")
}

func (db *DB) PutLogService(w http.ResponseWriter, r *http.Request) {
	var logService model.LogService
	id := r.PathValue("logServiceId")
	user := r.Context().Value("key").(string)

	if err := json.NewDecoder(r.Body).Decode(&logService); err != nil {
		http.Error(w, "error getting json data", http.StatusBadRequest)
	}

	validate := validator.New()
	if err := validate.Struct(&logService); err != nil {
		slog.Error(err.Error())
		http.Error(w, "name not inputed", http.StatusInternalServerError)
		return
	}

	_, err := db.Exec(`UPDATE "LogService" SET name = $1, updated_at = $2 WHERE id = $3 AND user_id = $4`, logService.Name, time.Now(), id, user)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error updating project", http.StatusInternalServerError)
		return
	}

	slog.Info("put log services")
}

func (db *DB) DeleteLogService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("logServiceId")
	user := r.Context().Value("key").(string)

	_, err := db.Exec(`DELETE FROM "LogService" WHERE id = $1 AND user_id = $2`, id, user)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error deleting project", http.StatusInternalServerError)
		return
	}

	slog.Info("delete log services")
}

// Logging
func (db *DB) GetLogs(w http.ResponseWriter, r *http.Request) {
	var log model.Log
	user := r.Context().Value("key").(string)

	rows, err := db.Query(`SELECT id, timestamp, severity, message, info, log_service_id, user_id FROM "Log" WHERE user_id = $1 ORDER BY timestamp ASC`, user)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error getting projects", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	logs := []model.Log{}
	for rows.Next() {
		if err := rows.Scan(&log.ID, &log.Timestamp, &log.Severity, &log.Message, &log.Info, &log.LogServiceID, &log.UserID); err != nil {
			slog.Error(err.Error())
			http.Error(w, "error getting project", http.StatusInternalServerError)
			return
		}
		logs = append(logs, log)
	}

	jsonData, err := json.Marshal(logs)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error getting projects", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonData)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error sending projects", http.StatusInternalServerError)
		return
	}

	slog.Info("get logs")
}

func (db *DB) CreateLog(w http.ResponseWriter, r *http.Request) {
	var log model.Log
	user := r.Context().Value("key").(string)

	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		slog.Error(err.Error())
		http.Error(w, "error getting json data", http.StatusBadRequest)
	}

	validate := validator.New()
	if err := validate.Struct(&log); err != nil {
		slog.Error(err.Error())
		http.Error(w, "name not inputed", http.StatusInternalServerError)
		return
	}

	if _, err := db.Exec(`INSERT INTO "Log" (timestamp, severity, message, info, log_service_id, user_id) VALUES ($1, $2, $3, $4, $5, $6)`, time.Now(), log.Severity, log.Message, log.Info, log.LogServiceID, user); err != nil {
		slog.Error(err.Error())
		http.Error(w, "error creating project", http.StatusInternalServerError)
		return
	}

	slog.Info("create log")
}

func (db *DB) DeleteLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("logId")
	user := r.Context().Value("key").(string)

	_, err := db.Exec(`DELETE FROM "Log" WHERE id = $1 AND user_id = $2`, id, user)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error deleting project", http.StatusInternalServerError)
		return
	}

	slog.Info("delete log")
}
