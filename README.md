# thrg

`thrg` es una aplicación web desarrollada en Go que implementa un sistema de autenticación dual para administradores y jugadores, con un mecanismo de registro basado en tokens. La aplicación está diseñada para ser desplegada fácilmente mediante Docker.

## Características

-   **Autenticación de Administrador:** Acceso seguro para administradores a un panel de control.
-   **Panel de Control de Administrador:**
    -   Generación de tokens de registro para nuevos jugadores.
    -   Visualización de todos los tokens generados, su estado (disponible o usado) y qué jugador lo utilizó.
    -   Lista de todos los jugadores registrados en el sistema.
-   **Registro de Jugadores por Token:** Los nuevos usuarios solo pueden registrarse utilizando un token válido proporcionado por un administrador.
-   **Autenticación de Jugadores:** Los jugadores pueden iniciar sesión para acceder a una página de juego.
-   **Roles de Usuario:** Diferenciación clara entre roles de `admin` y `player`.
-   **Contenerización:** Totalmente compatible con Docker para un despliegue y desarrollo sencillos.

## Tecnologías Utilizadas

-   **Backend:** [Go](https://golang.org/)
-   **Enrutador HTTP:** [Chi](https://github.com/go-chi/chi)
-   **ORM:** [GORM](https://gorm.io/) para la interacción con la base de datos.
-   **Base de Datos:** [PostgreSQL](https://www.postgresql.org/)
-   **Autenticación:** [gorilla/sessions](https://github.com/gorilla/sessions) para el manejo de sesiones.
-   **Contenerización:** [Docker](https://www.docker.com/) y [Docker Compose](https://docs.docker.com/compose/)
-   **Variables de Entorno:** [godotenv](https://github.com/joho/godotenv)

## Cómo Empezar

### Prerrequisitos

-   [Docker](https://www.docker.com/get-started)
-   [Docker Compose](https://docs.docker.com/compose/install/)

### Instalación y Ejecución

1.  **Clona el repositorio:**

    ```bash
    git clone <URL_DEL_REPOSITORIO>
    cd thrg
    ```

2.  **Crea un archivo `.env`:**

    Crea un archivo `.env` en la raíz del proyecto, basándote en el `.env.example`. Asegúrate de generar un valor seguro y aleatorio para `SESSION_SECRET`.

    ```env
    POSTGRES_USER=user_app
    POSTGRES_PASSWORD=password_segura
    POSTGRES_DB=rol_db
    SESSION_SECRET=<your_strong_session_secret>
    DB_HOST=db
    DB_PORT=5432
    PORT=8080
    APP_ENV=development # Cambia a 'production' para cookies seguras
    ```

3.  **Inicia la aplicación con Docker Compose:**

    Este comando construirá la imagen de la aplicación Go y levantará los contenedores de la aplicación y la base de datos.

    ```bash
    docker-compose up --build
    ```

    La aplicación estará disponible en `http://localhost:8080`.

## Endpoints de la API

### Autenticación y Configuración

-   `POST /api/admin/setup`: Crea el primer usuario administrador. Solo puede ser ejecutado una vez.
-   `POST /api/player/register`: Registra a un nuevo jugador utilizando un token válido.
-   `POST /api/player/login`: Inicia sesión como jugador.

### Rutas de Administrador

-   `GET /admin/login`: Página de inicio de sesión para administradores.
-   `POST /admin/login`: Procesa el formulario de inicio de sesión del administrador.
-   `GET /admin/dashboard`: Panel de control del administrador (ruta protegida).
-   `GET /admin/logout`: Cierra la sesión del administrador.
-   `POST /admin/api/tokens`: (API) Genera un nuevo token de registro.
-   `GET /admin/api/tokens`: (API) Lista todos los tokens de registro.
-.
-   `GET /admin/api/players`: (API) Lista todos los jugadores registrados.

### Rutas de Jugador

-   `GET /player/login`: Página de inicio de sesión para jugadores.
-   `GET /player/game`: Página principal del juego para jugadores autenticados (ruta protegida).
-   `GET /player/logout`: Cierra la sesión del jugador.

## Estructura del Proyecto

```
├── cmd/server/main.go      # Punto de entrada de la aplicación
├── internal/               # Lógica de negocio principal
│   ├── contextutil/        # Utilidades de contexto
│   ├── core/               # Modelos de dominio principales
│   ├── token/              # Lógica para tokens (modelo, repositorio, handler)
│   └── user/               # Lógica para usuarios (modelo, repositorio, handler, auth)
├── web/                    # Archivos HTML del frontend
├── .env.example            # Ejemplo de variables de entorno
├── Dockerfile              # Define la imagen Docker de la aplicación
├── docker-compose.yml      # Orquesta los servicios de la aplicación
└── go.mod                  # Dependencias del proyecto
```
