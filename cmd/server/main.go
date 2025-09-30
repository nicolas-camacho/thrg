package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/nicolas-camacho/thrg/internal/token"
	"github.com/nicolas-camacho/thrg/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	adminSessionName  = "admin_session"
	playerSessionName = "player_session"
)

var store *sessions.CookieStore

func getDBConnectionString() string {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables")
	}

	gob.Register(uuid.UUID{})
	log.Println("Type uuid.UUID registered with gob.")

	connStr := getDBConnectionString()

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL!")

	log.Println("Verifying extension 'uuid-ossp'...")
	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	if err != nil {
		log.Fatalf("Failed to create extension 'uuid-ossp': %v", err)
	}
	log.Println("Extension 'uuid-ossp' is available.")

	log.Println("Running database migrations...")

	err = db.AutoMigrate(
		&user.User{},
		&token.RegistrationToken{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	log.Println("Database migrated successfully!")

	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   int((time.Hour * 24).Seconds()),
		HttpOnly: true,
		Secure:   os.Getenv("APP_ENV") == "production",
		SameSite: http.SameSiteLaxMode,
	}

	user.Store = store
	log.Println("Session store initialized.")

	userRepo := user.NewRepository(db)
	tokenRepo := token.NewRepository(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	//API ROUTES
	r.Post("/api/player/register", user.RegisterPlayerHandler(userRepo, tokenRepo))
	r.Post("/api/player/login", user.PlayerLoginHandler(userRepo, playerSessionName))
	r.Post("/api/admin/setup", user.SetupAdminHandler(userRepo))

	// Admin routes
	r.Handle("/admin/login", user.ServeLoginPageHandler(userRepo))
	r.Get("/admin/logout", func(w http.ResponseWriter, r *http.Request) {
		if err := user.LogoutUser(w, r, adminSessionName); err != nil {
			log.Printf("Failed to log out user: %v", err)
		}
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
	})
	r.Group(func(r chi.Router) {
		r.Use(user.AdminAuthMiddleware(adminSessionName))

		r.Get("/admin/dashboard", user.DashboardHandler())
		r.Post("/admin/api/tokens", token.GenerateTokenHandler(tokenRepo))
		r.Get("/admin/api/tokens", token.ListTokensHandler(tokenRepo, userRepo))
		r.Get("/admin/api/players", user.ListPlayersHandler(userRepo))
	})

	// Player routes
	r.Get("/player/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/player_login.html")
	})
	r.Get("/player/logout", func(w http.ResponseWriter, r *http.Request) {
		if err := user.LogoutUser(w, r, playerSessionName); err != nil {
			log.Printf("Failed to log out user: %v", err)
		}
		http.Redirect(w, r, "/player/login", http.StatusSeeOther)
	})
	r.Group(func(r chi.Router) {
		r.Use(user.PlayerAuthMiddleware(playerSessionName))

		r.Get("/player/game", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "web/game.html")
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
