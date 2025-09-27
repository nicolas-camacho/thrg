package user

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type LoginPageData struct {
	Error string
}

var loginTmpl *template.Template

func init() {
	var err error
	loginTmpl, err = template.ParseFiles("web/login.html")
	if err != nil {
		log.Fatalf("Failed to parse login template: %v", err)
	}
}

func ServeLoginPageHandler(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			loginTmpl.Execute(w, LoginPageData{Error: ""})
			return
		}

		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				log.Printf("Failed to parse form: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			username := r.FormValue("username")
			password := r.FormValue("password")

			user, err := repo.Authenticate(r.Context(), username, password)
			if err != nil {
				data := LoginPageData{Error: "Invalid username or password"}
				w.WriteHeader(http.StatusUnauthorized)
				loginTmpl.Execute(w, data)
				return
			}

			if user.Role != RoleAdmin {
				data := LoginPageData{Error: "Access denied, admin only"}
				w.WriteHeader(http.StatusForbidden)
				loginTmpl.Execute(w, data)
				return
			}

			log.Printf("Admin %s logged in successfully", user.Username)
			http.Redirect(w, r, "/admin/dashboard", http.StatusSeeOther)
			return
		}
	}
}

type setupAdminRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SetupAdminHandler(repo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		exists, err := repo.CheckAdminExists(ctx)
		if err != nil {
			http.Error(w, "Failed to check admin existence", http.StatusInternalServerError)
			return
		}
		if exists {
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte(`{"error":"Admin user already exists"}`))
			return
		}

		var req setupAdminRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if req.Username == "" || req.Password == "" {
			http.Error(w, "Username and password are required", http.StatusBadRequest)
			return
		}

		if err := repo.CreateAdmin(ctx, req.Username, req.Password); err != nil {
			log.Printf("Failed to create admin user: %v", err)
			http.Error(w, "Failed to create admin user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message":"Admin user created successfully"}`))
	}
}
