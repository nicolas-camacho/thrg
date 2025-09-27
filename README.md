# thrg

This is a Go web application that uses the Chi router, GORM for database interaction, and PostgreSQL as the database. The application is containerized using Docker and includes session management and token-based registration.

## Prerequisites

- Docker
- Docker Compose

## Getting Started

1.  **Create a `.env` file:**

    Create a `.env` file in the root of the project with the following content. Make sure to use a strong, randomly generated string for `SESSION_SECRET`. The application loads environment variables from this file using `godotenv`.

    ```
    POSTGRES_USER=user_app
    POSTGRES_PASSWORD=password_segura
    POSTGRES_DB=rol_db
    SESSION_SECRET=<your_strong_session_secret>
    DB_HOST=localhost
    DB_PORT=5432
    PORT=8080
    APP_ENV=development # Set to 'production' for secure cookies
    ```

2.  **Run the application:**

    ```bash
    docker-compose up --build
    ```

    The application will be available at `http://localhost:8080` (or the port specified in `PORT`).

## Technologies Used

-   [Go](https://golang.org/)
-   [Chi](https://github.com/go-chi/chi)
-   [GORM](https://gorm.io/)
-   [gorilla/sessions](https://github.com/gorilla/sessions)
-   [PostgreSQL](https://www.postgresql.org/)
-   [Docker](https://www.docker.com/)
-   [godotenv](https://github.com/joho/godotenv)

## API Endpoints

-   `GET /ping`: Returns a `pong` response. Used for health checks.
-   `POST /api/admin/setup`: Sets up the initial admin user.
-   `GET /admin/login`: Serves the admin login page.
-   `GET /admin/dashboard`: Serves the admin dashboard page, which allows generating and listing registration tokens. This endpoint is protected and requires authentication.
-   `POST /admin/api/tokens`: Generates a new registration token.
-   `GET /admin/api/tokens`: Lists all registration tokens.
-   `GET /admin/logout`: Logs out the current user.

## Token-based Registration

This application includes a feature for token-based user registration. An admin can generate a registration token from the dashboard that can be used by a new user to register.