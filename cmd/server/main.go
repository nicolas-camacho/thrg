package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"github.com/nicolas-camacho/thrg/internal/token"
	"github.com/nicolas-camacho/thrg/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	sessionName = "admin_session"
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
	connStr := getDBConnectionString()

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL!")

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
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}

	user.Store = store
	log.Println("Session store initialized.")

	userRepo := user.NewRepository(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})

	//r.Post("/api/admin/setup", user.SetupAdminHandler(userRepo))

	r.Handle("/admin/login", user.ServeLoginPageHandler(userRepo))

	r.Group(func(r chi.Router) {
		r.Use(user.AdminAuthMiddleware)

		r.Get("/admin/dashboard", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Welcome to the admin dashboard!"))
		})
	})

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
