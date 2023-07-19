# Authentication Service in Go

This is a simple authentication service built in Go. It includes functionalities such as user registration, login, and viewing all users.

## Project Structure

\```
.
├── main.go
├── go.mod
├── go.sum
└── pkg
    └── handlers
        ├── register.go
        ├── login.go
        └── home.go
└── templates
    ├── register.html
    ├── login.html
    └── home.html
\```

### Requirements

- Go (v1.16 or later)
- Postgres

### Environment Variables

You need to set the following environment variables for database configuration and security:

- `DB_HOST`: The hostname of your Postgres database.
- `DB_PORT`: The port number where your Postgres is running.
- `DB_USER`: The username of your Postgres database.
- `DB_PASSWORD`: The password of your Postgres database.
- `DB_NAME`: The name of your Postgres database.
- `SECRET_KEY`: The secret key used for session cookies.

You can set these variables in a `.env` file at the root of your project.

### Running the Project

1. Make sure you have created your database locally:

2. Run the Go server:

\```bash
go run cmd/server/main.go
\```

This will start the server on `localhost:8080`.
