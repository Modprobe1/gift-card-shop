# How to Build the Exchange Project

This is a gift card exchange application with a Go backend, HTML/CSS/JS frontend, and MongoDB database.

## Project Structure
```
.
├── backend/           # Go API server
│   ├── cmd/           # Main application entry point
│   ├── handlers/      # HTTP handlers
│   ├── models/        # Data models
│   ├── Dockerfile     # Backend container config
│   └── go.mod         # Go dependencies
├── frontend/          # Frontend web application
│   ├── dist/          # Built frontend files (HTML/CSS/JS)
│   └── Dockerfile     # Frontend container config
└── docker-compose.yml # Multi-container deployment config
```

## Build Options

### Option 1: Using Docker Compose (Recommended)

This is the easiest way to build and run the entire application:

```bash
# Build and start all services (frontend, backend, MongoDB)
docker-compose up --build

# Or run in background
docker-compose up --build -d

# Stop all services
docker-compose down
```

**Access the application:**
- Frontend: http://localhost (port 80)
- Backend API: http://localhost:8080
- MongoDB: localhost:27017

### Option 2: Build Individual Components

#### Backend (Go)
```bash
cd backend

# Install dependencies and build
go mod tidy
cd cmd
go build -o ../exchange-api .

# Run the backend
cd ..
./exchange-api
```

**Requirements:**
- Go 1.21 or later
- MongoDB running on localhost:27017

#### Frontend
The frontend is already built and located in `frontend/dist/`. To serve it:

```bash
cd frontend/dist

# Using Python's built-in server
python3 -m http.server 8000

# Or using Node.js serve
npx serve .

# Or any other static file server
```

### Option 3: Docker Build Individual Services

#### Build Backend Docker Image
```bash
cd backend
docker build -t exchange-backend .
docker run -p 8080:8080 exchange-backend
```

#### Build Frontend Docker Image
```bash
cd frontend
docker build -t exchange-frontend .
docker run -p 80:80 exchange-frontend
```

## Dependencies

### Backend Dependencies (from go.mod)
- Gin web framework (github.com/gin-gonic/gin v1.9.1)
- MongoDB driver (go.mongodb.org/mongo-driver v1.13.1)
- UUID generator (github.com/google/uuid v1.3.0)

### Database
- MongoDB (automatically included in docker-compose setup)

## Quick Start

1. **Clone/Extract the project**
2. **Make sure Docker and Docker Compose are installed**
3. **Run the application:**
   ```bash
   docker-compose up --build
   ```
4. **Access the application at http://localhost**

## Development

For development, you can run each service separately:

1. **Start MongoDB:**
   ```bash
   docker run -d -p 27017:27017 mongo
   ```

2. **Start Backend:**
   ```bash
   cd backend
   go mod tidy
   cd cmd && go run . &
   ```

3. **Serve Frontend:**
   ```bash
   cd frontend/dist
   python3 -m http.server 8000
   ```

## Troubleshooting

- **Port conflicts:** Make sure ports 80, 8080, and 27017 are available
- **Docker issues:** Try `docker-compose down` and rebuild with `docker-compose up --build`
- **Go module issues:** Run `go mod tidy` in the backend directory
- **Permission issues:** Make sure Docker has proper permissions on your system

## Application Features

Based on the file structure, this appears to be an exchange application with:
- Main exchange interface (`index.html`)
- Admin panel (`admin.html`) 
- Transaction interface (`tx.html`)
- Go backend API with MongoDB integration