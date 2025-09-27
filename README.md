# thrg

This is a Go web application that uses the Chi router, GORM for database interaction, and PostgreSQL as the database. The application is containerized using Docker.

## Prerequisites

- Docker
- Docker Compose

## Getting Started

1.  **Create a `.env` file:**

    Create a `.env` file in the root of the project with the following content:

    ```
    POSTGRES_USER=user_app
    POSTGRES_PASSWORD=password_segura
    POSTGRES_DB=rol_db
    ```

2.  **Run the application:**

    ```bash
    docker-compose up --build
    ```

    The application will be available at `http://localhost:8080`.

## Technologies Used

-   [Go](https://golang.org/)
-   [Chi](https://github.com/go-chi/chi)
-   [GORM](https://gorm.io/)
-   [PostgreSQL](https://www.postgresql.org/)
-   [Docker](https://www.docker.com/)

## API Endpoints

-   `GET /ping`: Returns a `pong` response. Used for health checks.
-   `GET /admin/login`: Serves the admin login page.
-   `GET /admin/dashboard`: Serves the admin dashboard page.
