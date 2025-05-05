# Todo API

This is a simple Todo list API built with Go, Gin, GORM, and MariaDB. It demonstrates basic CRUD operations, authentication, and deployment using Docker that I have learned from [Skooldio](https://www.skooldio.com/courses/developing-robust-api-services-with-go)

## Features

* Create, Read, and Delete Todos
* JWT-based authentication for protected endpoints
* Rate limiting on specific endpoints
* Health checks (`/healthz`)
* Graceful shutdown
* Dockerized setup using Docker Compose

## Getting Started

### Prerequisites

* [Go](https://golang.org/doc/install) (version 1.18 or later recommended)
* [Docker](https://docs.docker.com/get-docker/)
* [Docker Compose](https://docs.docker.com/compose/install/)
* [Make](https://www.gnu.org/software/make/) (optional, for using Makefile commands)

### Installation

1. **Clone the repository:**

    ```bash
    git clone <repository-url>
    cd todo
    ```

2. **Create environment file:**
Copy the example environment variables:

    ```bash
    cp local.env.example local.env
    ```

*Note: You might need to create `local.env.example` based on your `local.env` file.*
Update `local.env` with your desired settings, especially the `SIGN` key for JWT.

### Running the Application

#### Option 1: Using Docker Compose (Recommended)

This method starts both the API and the MariaDB database in containers.

1. **Build and run the services:**

```bash
docker-compose up --build
```

The API will be available at `http://localhost:8081`.

#### Option 2: Running Locally (Requires MariaDB Running Separately)

1. **Start MariaDB:**

    You can use the provided Makefile target to start a MariaDB container:

    ```bash
    make maria
    ```

    *Ensure port 3306 is free or modify the port mapping in the `Makefile`.*

2. **Build the Application:**

    ```bash
    make build
    ```

3. **Run the Application:**

    ```bash
    make run
    ```

    The API will be available at `http://localhost:8081`.

## API Endpoints

* `GET /healthz`: Health check endpoint. Returns `200 OK`.
* `GET /pingz`: Simple ping endpoint. Returns `{"message": "pong"}`.
* `GET /limitz`: Rate-limited endpoint.
* `GET /x`: Returns build information (`{"buildcommit": "...", "buildtime": "..."}`).
* `GET /tokenz`: Generates a JWT token for accessing protected routes.
* **Protected Routes (require `Authorization: Bearer <token>` header):**
* `POST /todos`: Create a new todo. Request body: `{"text": "Your todo text"}`.
* `GET /todos`: Get all todos.
* `DELETE /todos/:id`: Delete a todo by its ID.
