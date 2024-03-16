package project

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

func (db *DB) GetProjects(w http.ResponseWriter, r *http.Request) {
	var project model.Project
	user := r.Context().Value("key").(string)

	rows, err := db.Query(`SELECT id, name,created_at, updated_at, user_id FROM "Project" WHERE user_id = $1 ORDER BY created_at ASC`, user)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error getting projects", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	projects := []model.Project{}
	for rows.Next() {
		if err := rows.Scan(&project.ID, &project.Name, &project.CreatedAt, &project.UpdatedAt, &project.UserID); err != nil {
			slog.Error(err.Error())
			http.Error(w, "error getting project", http.StatusInternalServerError)
			return
		}
		projects = append(projects, project)
	}

	jsonData, err := json.Marshal(projects)
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

	slog.Info("getprojects")
}

func (db *DB) GetProject(w http.ResponseWriter, r *http.Request) {
	var project model.Project
	id := r.PathValue("projectId")
	user := r.Context().Value("key").(string)

	row := db.QueryRow(`SELECT id, name, created_at, updated_at FROM "Project" WHERE id = $1 AND user_id = $2`, id, user)
	if err := row.Scan(&project.ID, &project.Name, &project.CreatedAt, &project.UpdatedAt); err != nil {
		slog.Error(err.Error())
		http.Error(w, "error getting project", http.StatusInternalServerError)
		return
	}

	jsonData, err := json.Marshal(project)
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

	slog.Info("getproject")
}

func (db *DB) CreateProjects(w http.ResponseWriter, r *http.Request) {
	var project model.Project
	user := r.Context().Value("key").(string)

	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		slog.Error(err.Error())
		http.Error(w, "error getting json data", http.StatusBadRequest)
	}

	validate := validator.New()
	if err := validate.Struct(&project); err != nil {
		slog.Error(err.Error())
		http.Error(w, "name not inputed", http.StatusInternalServerError)
		return
	}

	if _, err := db.Exec(`INSERT INTO "Project" (name, updated_at, user_id) VALUES ($1, $2, $3)`, project.Name, time.Now(), user); err != nil {
		slog.Error(err.Error())
		http.Error(w, "error creating project", http.StatusInternalServerError)
		return
	}

	slog.Info("createprojects")
}

func (db *DB) UpdateProjects(w http.ResponseWriter, r *http.Request) {
	var project model.Project
	id := r.PathValue("projectId")
	user := r.Context().Value("key").(string)

	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, "error getting json data", http.StatusBadRequest)
	}

	_, err := db.Exec(`UPDATE "Project" SET name = $1, updated_at = $2 WHERE id = $3 AND user_id = $4`, project.Name, time.Now(), id, user)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error updating project", http.StatusInternalServerError)
		return
	}

	slog.Info("updateprojects")
}

func (db *DB) DeleteProjects(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("projectId")
	user := r.Context().Value("key").(string)

	_, err := db.Exec(`DELETE FROM "Project" WHERE id = $1 AND user_id = $2`, id, user)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, "error deleting project", http.StatusInternalServerError)
		return
	}

	slog.Info("deleteprojects")
}
