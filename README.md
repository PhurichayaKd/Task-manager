# Task Manager

## Overview
Task Manager is a web application designed to help users efficiently manage their tasks. It includes a Go-based backend and a frontend built with vanilla JavaScript. The project also uses Docker for containerization.

## Features
- User authentication
- Task creation, retrieval, update, and deletion
- Language and theme switching
- Responsive UI

## Technologies Used
- **Backend**: Go, Gin framework
- **Frontend**: Vanilla JavaScript, HTML, CSS
- **Database**: PostgreSQL
- **Containerization**: Docker, Docker Compose

## Prerequisites
- Docker and Docker Compose installed
- Go installed

## Setup Instructions

### 1. Clone the Repository
```bash
git clone https://github.com/your-username/task-manager.git
cd task-manager
```

### 2. Start Docker Containers
```bash
docker-compose up -d
```

### 3. Run the Backend Server
```bash
cd cmd/server
go run main.go
```

### 4. Open the Frontend
- Navigate to the `frontend/vanilla` folder
- Open `index.html` in your browser

## Environment Variables
Ensure the following environment variables are set:
- `DB_DSN`: Database connection string (e.g., `postgres://app:app@db:5432/taskdb`)

## Folder Structure
```
.
├── cmd
│   └── server
│       └── main.go
├── internal
│   ├── api
│   │   └── task_handler.go
│   ├── config
│   │   └── config.go
│   └── db
│       └── migrate
│           └── 0001_init.sql
├── frontend
│   └── vanilla
│       ├── index.html
│       ├── index.js
│       └── styles.css
├── docker-compose.yml
└── README.md
```

## License
This project is licensed under the MIT License.
