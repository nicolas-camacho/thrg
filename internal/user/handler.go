package user

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nicolas-camacho/thrg/internal/token"
)

type LoginPageData struct {
	Error string
}

var loginTmpl *template.Template
var dashboardTmpl *template.Template

func init() {
	var err error
	loginTmpl, err = template.ParseFiles("web/login.html")
	if err != nil {
		log.Fatalf("Failed to parse login template: %v", err)
	}

	dashboardTmpl, err = template.ParseFiles("web/dashboard.html")
	if err != nil {
		log.Fatalf("Failed to parse dashboard template: %v", err)
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
				
								if err := LoginUser(w, r, user.ID); err != nil {
									log.Printf("Failed to log in user: %v", err)
									http.Error(w, "Internal Server Error", http.StatusInternalServerError)
									return
								}
				
								log.Printf("Admin %s logged in successfully (ID: %s)", user.Username, user.ID)
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

		_, err = repo.CreateUser(ctx, req.Username, req.Password, RoleAdmin)
		if err != nil {
			log.Printf("Failed to create admin user: %v", err)
			http.Error(w, "Failed to create admin user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message":"Admin user created successfully"}`))
	}
}

func DashboardHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/dashboard.html")
	}
}

func LogoutHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := Store.Get(r, SessionName)
		if err != nil {
			http.Redirect(w, r, "/admin/login", http.StatusTemporaryRedirect)
			return
		}

		session.Values[userKey] = nil
		session.Values["role"] = nil
		session.Options.MaxAge = -1
		if err := session.Save(r, w); err != nil {
			log.Printf("Error saving session during logout: %v", err)
		}

		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	}
}

type RegistrationTokenRepository interface {
	ValidateAndUseToken(ctx context.Context, tokenValue string, userID uuid.UUID) (*token.RegistrationToken, error)
}

type RegisterPlayerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

func RegisterPlayerHandler(userRepo *Repository, tokenRepo *token.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterPlayerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if req.Username == "" || req.Password == "" || req.Token == "" {
			http.Error(w, "Username, password, and token are required", http.StatusBadRequest)
			return
		}

		newUser, err := userRepo.CreateUser(r.Context(), req.Username, req.Password, RolePlayer)
		if err != nil {
			log.Printf("Failed to create player user: %v", err)
			if err.Error() == "username already exists" {
				http.Error(w, "Username already exists", http.StatusConflict)
			} else {
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
			}
			return
		}

		_, err = tokenRepo.ValidateAndUseToken(r.Context(), req.Token, newUser.ID)
		if err != nil {
			if deleteErr := userRepo.DeleteUser(r.Context(), newUser.ID); deleteErr != nil {
				log.Printf("Failed to delete user after token validation failure: %v", deleteErr)
			}

			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "User registered successfully",
			"userId":  newUser.ID.String(),
		})
	}
}

func ListPlayersHandler(userRepo *Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		players, err := userRepo.GetAllPlayers(r.Context())
		if err != nil {
			log.Printf("Error al listar jugadores: %v", err)
			http.Error(w, "Error interno al obtener la lista de jugadores.", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// DTO para la respuesta JSON
		type PlayerDTO struct {
			ID        uuid.UUID `json:"ID"`
			Username  string    `json:"Username"`
			Role      string    `json:"Role"`
			CreatedAt time.Time `json:"CreatedAt"`
		}

		dtos := make([]PlayerDTO, len(players))
		for i, p := range players {
			dtos[i] = PlayerDTO{
				ID:        p.ID,
				Username:  p.Username,
				Role:      p.Role,
				CreatedAt: p.CreatedAt,
			}
		}

		if err := json.NewEncoder(w).Encode(dtos); err != nil {
			log.Printf("Error al codificar JSON de jugadores: %v", err)
		}
	}
}
