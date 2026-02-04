# Forum Web Application

A modern forum web application built with Go, featuring user authentication, post management, comment system, and reaction functionality. This project demonstrates clean architecture principles with a layered approach including handlers, services, repositories, and middleware.

## 📋 Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Database](#database)
- [Authentication & Middleware](#authentication--middleware)
- [Project Structure](#project-structure)
- [Features](#features)
- [Installation & Setup](#installation--setup)
- [Usage](#usage)
- [Endpoints](#api-endpoints)
- [Docker Deployment](#docker-deployment)
- [Development](#development)

## 🌟 Overview

This forum application provides a complete community platform where users can:
- Register and authenticate securely
- Create and manage posts with categories
- Comment on posts and interact with content
- Like/dislike posts and comments
- View personalized profile pages with activity tracking
- Browse content by categories with pagination

The application follows modern web development practices with server-side rendering using Go's `html/template` package and implements a clean separation of concerns through a layered architecture.

## 🏗️ Architecture

The application follows a clean architecture pattern with the following layers:

```
├── cmd/web/           # Application entry point
├── internal/          # Private application code
│   ├── handlers/      # HTTP request handlers
│   ├── services/      # Business logic layer
│   ├── repository/    # Data access layer
│   ├── middleware/    # HTTP middleware
│   ├── models/        # Data models
│   └── utils/         # Utility functions
├── pkg/               # Public packages
├── templates/         # HTML templates
└── static/           # Static assets (CSS, JS, images)
```

### Layer Responsibilities

- **Handlers**: Handle HTTP requests/responses, input validation, and template rendering
- **Services**: Implement business logic and coordinate between handlers and repositories
- **Repository**: Manage database operations and data persistence
- **Middleware**: Handle cross-cutting concerns like authentication and request processing
- **Models**: Define data structures and domain entities

## 🗄️ Database

The application uses SQLite as its database engine with automatic initialization and schema management.

### Database Configuration

- **Database File**: `forum.db` (created automatically)
- **Schema Location**: `internal/database/schema.sql`
- **Driver**: `github.com/mattn/go-sqlite3`
- **Connection Pool**: Max 10 concurrent connections
- **Foreign Keys**: Enabled for referential integrity

### Database Tables

#### User Table
```sql
CREATE TABLE user (
    id TEXT PRIMARY KEY,           -- UUID
    email TEXT NOT NULL UNIQUE,   -- User email
    username TEXT NOT NULL UNIQUE, -- Display name
    password TEXT NOT NULL,        -- Bcrypt hashed password
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

#### Post Table
```sql
CREATE TABLE post (
    id TEXT PRIMARY KEY,           -- UUID
    user_id TEXT NOT NULL,        -- Foreign key to user
    title TEXT NOT NULL,          -- Post title
    content TEXT NOT NULL,        -- Post content
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user(id)
);
```

#### Comment Table
```sql
CREATE TABLE comment (
    id TEXT PRIMARY KEY,           -- UUID
    user_id TEXT NOT NULL,        -- Comment author
    post_id TEXT NOT NULL,        -- Associated post
    content TEXT NOT NULL,        -- Comment content
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user(id),
    FOREIGN KEY (post_id) REFERENCES post(id)
);
```

#### Category Table
```sql
CREATE TABLE category (
    id TEXT PRIMARY KEY,           -- UUID
    name TEXT NOT NULL UNIQUE     -- Category name
);
```

#### Post-Category Junction Table
```sql
CREATE TABLE post_category (
    post_id TEXT NOT NULL,
    category_id TEXT NOT NULL,
    PRIMARY KEY (post_id, category_id),
    FOREIGN KEY (post_id) REFERENCES post(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES category(id) ON DELETE CASCADE
);
```

#### Reaction Table
```sql
CREATE TABLE reaction (
    id TEXT PRIMARY KEY,           -- UUID
    user_id TEXT NOT NULL,        -- User who reacted
    post_id TEXT,                 -- Target post (nullable)
    comment_id TEXT,              -- Target comment (nullable)
    type TEXT NOT NULL CHECK(type IN ('like', 'dislike')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user(id),
    -- Ensures reaction is either on post OR comment, not both
    CHECK ((post_id IS NULL AND comment_id IS NOT NULL) OR 
           (post_id IS NOT NULL AND comment_id IS NULL))
);
```

#### Session Table
```sql
CREATE TABLE session (
    id TEXT PRIMARY KEY,           -- Session token (UUID)
    user_id TEXT NOT NULL,        -- Associated user
    expires_at DATETIME NOT NULL, -- Session expiration
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES user(id)
);
```

### Database Indexes

Performance-optimized indexes are automatically created


### Database Initialization

The database initializes automatically on first run:

1. **Connection Setup**: Opens SQLite connection with foreign key constraints enabled
2. **Schema Creation**: Executes SQL schema from `internal/database/schema.sql`
3. **Data Population**: Seeds initial categories (General Discussion, Tech News, etc.)
4. **Validation**: Verifies database connectivity and structure


## 🔐 Authentication & Middleware

### Authentication System

The application implements a secure session-based authentication system:

#### Password Security
- **Hashing**: Uses bcrypt with default cost factor (10)
- **Validation**: Enforces strong password requirements
- **Storage**: Never stores plaintext passwords

#### Session Management
- **Token Generation**: UUID-based session tokens
- **Storage**: Server-side session storage in database
- **Expiration**: 24-hour session lifetime
- **Security**: HttpOnly, Secure, and SameSite cookie attributes


### Middleware System

#### Authentication Middleware (`middleware.Auth`)

Handles authentication for all requests:

```go
func (m *Middleware) Auth(next http.Handler) http.Handler {
    // 1. Extract session cookie
    // 2. Validate session in database
    // 3. Check session expiration
    // 4. Load user data
    // 5. Set user context for handlers
    // 6. Clean up invalid/expired sessions
}
```

**Features:**
- Automatic session validation
- Context-based user injection
- Graceful handling of invalid sessions
- Automatic cookie cleanup

#### Authorization Middleware (`middleware.RequireAuth`)

Protects routes requiring authentication:

```go
func (m *Middleware) RequireAuth(next http.Handler) http.Handler {
    // 1. Check if user exists in context
    // 2. Redirect to login if not authenticated
    // 3. Continue to protected handler if authenticated
}
```

#### User Context Management

```go
// Retrieve authenticated user in any handler:
user := middleware.GetUser(r)
if user != nil {
    // User is authenticated
    userID := user.ID
    username := user.Username
}
```

### Security Features

- **CSRF Protection**: SameSite cookie attribute prevents CSRF attacks
- **Session Hijacking Prevention**: Server-side session validation
- **Secure Cookies**: HTTPOnly and Secure flags when using HTTPS
- **Input Sanitization**: HTML escaping for user-generated content
- **SQL Injection Protection**: Parameterized queries throughout

## 📁 Project Structure

```
forum-app/
├── cmd/
│   └── web/
│       └── main.go                 # Application entry point
├── internal/
│   ├── database/
│   │   ├── database.go            # Database initialization
│   │   └── schema.sql             # Database schema
│   ├── handlers/
│   │   ├── auth_handlers.go       # Authentication endpoints
│   │   ├── create_comment_handler.go
│   │   ├── create_post_handler.go
│   │   ├── error_handler.go       # Error page rendering
│   │   ├── handlers_helpers.go    # Handler utilities
│   │   ├── home_handler.go        # Homepage logic
│   │   ├── post_by_category_handler.go
│   │   ├── post_handler.go        # Individual post view
│   │   ├── profile_handler.go     # User profile
│   │   └── reaction_handlers.go   # Like/dislike functionality
│   ├── middleware/
│   │   └── middleware.go          # Authentication middleware
│   ├── models/
│   │   └── models.go              # Data structures
│   ├── repository/
│   │   └── repository.go          # Database operations
│   ├── routes/
│   │   └── routes.go              # Route configuration
│   ├── services/
│   │   ├── auth_service.go        # Authentication business logic
│   │   ├── category_service.go    # Category management
│   │   ├── comment_service.go     # Comment operations
│   │   ├── post_service.go        # Post management
│   │   └── reaction_service.go    # Reaction system
│   └── utils/
│       ├── paginate_num.go        # Pagination helpers
│       └── validation/
│           └── validation.go      # Input validation
├── pkg/
│   ├── UUID/
│   │   └── uuid.go                # UUID generation
│   ├── crypt/
│   │   └── password.go            # Password hashing
│   └── session/
│       └── session.go             # Session utilities
├── static/                        # Static assets
│   ├── *.css                     # Stylesheets
│   ├── *.js                      # JavaScript files
│   └── files/                    # Images and icons
├── templates/                     # HTML templates
│   ├── create-post.html
│   ├── errors.html
│   ├── index.html
│   ├── login.html
│   ├── post.html
│   ├── profile.html
│   └── signup.html
├── build.sh                      # Docker build script
├── cleanup.sh                    # Docker cleanup script
├── docker-compose.yml            # Docker Compose configuration
├── Dockerfile                    # Docker image definition
├── Makefile                      # Build automation
├── go.mod                        # Go module definition
└── go.sum                        # Go module checksums
```

## ✨ Features

### Core Functionality
- **User Registration & Login**: Secure account creation and authentication
- **Post Management**: Create, view, and categorize posts
- **Comment System**: Threaded discussions on posts
- **Reaction System**: Like/dislike posts and comments
- **Category Filtering**: Browse posts by topic categories
- **User Profiles**: View personal activity and statistics

### UI/UX Features
- **Responsive Design**: Mobile-friendly interface
- **Pagination**: Efficient content browsing
- **Real-time Updates**: Dynamic content loading
- **Error Handling**: User-friendly error pages
- **Form Validation**: Client and server-side validation

### Technical Features
- **Session Management**: Secure user sessions
- **Database Migrations**: Automatic schema management
- **Input Sanitization**: XSS protection
- **Performance Optimization**: Efficient database queries
- **Docker Support**: Containerized deployment

## 🚀 Installation & Setup

### Prerequisites

- **Go**: Version 1.19 or higher
- **SQLite**: Included with Go SQLite driver
- **Docker** (optional): For containerized deployment

### Local Development Setup

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd forum-app
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Run the application:**
   ```bash
   go run cmd/web/main.go
   ```

4. **Access the application:**
   ```
   http://localhost:8080
   ```

The database will be automatically created and initialized on first run.

## 📖 Usage

### Using the Makefile

The project includes a comprehensive Makefile for common operations:

```bash
# Run the application locally
make run

# Build Docker image
make docker-build

# Start application with Docker Compose or Podman Compose
make docker-run / podman-run 

# Stop Docker containers with Docker Compose or Podman Compose
make docker-stop / podman-stop

# Build and run using build script
make build-docker

# Clean up Docker resources
make cleanup
```

### Direct Commands

```bash
# Development server
go run cmd/web/main.go

# Build binary
go build -o forum cmd/web/main.go

# Run tests
go test ./...
```

## 🐳 Docker Deployment

### Using Docker Compose (Recommended)

1. **Start the application:**
   ```bash
   make docker-run
   # or
   docker-compose up -d
   ```

2. **View logs:**
   ```bash
   docker-compose logs -f
   ```

3. **Stop the application:**
   ```bash
   make docker-stop
   # or
   docker-compose down
   ```

### Using Build Script

The `build.sh` script provides automated Docker deployment:

```bash
# Make script executable and run
chmod +x build.sh
./build.sh

# Or use Makefile
make build-docker
```

**Script features:**
- Automatic cleanup of existing containers
- Image building with error handling
- Container deployment with health checks
- Status reporting and access information

### Manual Docker Commands

```bash
# Build image
docker build -t forum-app .

# Run container
docker run -d --name forum -p 8080:8080 forum-app

# View logs
docker logs forum

# Stop container
docker stop forum
```

### Docker Configuration

The application includes:
- **Multi-stage Dockerfile** for optimized image size
- **Health checks** for container monitoring
- **Volume mounting** for persistent data storage
- **Environment configuration** for deployment flexibility

## 🔧 Development

### Adding New Features

1. **Models**: Define data structures in `internal/models/`
2. **Repository**: Add database operations in `internal/repository/`
3. **Services**: Implement business logic in `internal/services/`
4. **Handlers**: Create HTTP handlers in `internal/handlers/`
5. **Routes**: Register routes in `internal/routes/routes.go`
6. **Templates**: Add HTML templates in `templates/`

### Database Migrations

To modify the database schema:

1. Update `internal/database/schema.sql`
2. Delete existing `forum.db` file
3. Restart application for automatic recreation

### Environment Variables

The application supports the following environment variables:

- `PORT`: Server port (default: 8080)
- `DB_PATH`: Database file path (default: forum.db)

## 📋 Endpoints

### Public Routes
- `GET /` - Homepage with post listings
- `GET /login` - Login page
- `GET /signup` - Registration page
- `POST /login` - Authenticate user
- `POST /signup` - Register new user
- `GET /posts/{id}` - View individual post
- `GET /categories/{id}/posts` - Posts by category

### Protected Routes (Require Authentication)
- `POST /logout` - User logout
- `GET /profile` - User profile page
- `GET /create-post` - Create post form
- `POST /create-post` - Submit new post
- `POST /posts/{id}/comments` - Add comment
- `POST /post/{id}/{reaction}` - React to post
- `POST /comment/{id}/{reaction}` - React to comment

### Static Assets
- `/static/` - CSS, JavaScript, and image files

## 🛠️ Scripts and Automation

### Build Script (`build.sh`)

Automated Docker build and deployment:

```bash
#!/bin/bash
# Features:
- Container cleanup and removal
- Docker image building with error handling
- Container deployment with port mapping
- Status reporting and access information
- Color-coded output for better visibility
```

### Cleanup Script (`cleanup.sh`)

Docker resource cleanup:

```bash
#!/bin/bash
# Removes:
- Forum-specific containers and images
- Unused Docker objects (containers, images, volumes, networks)
- Provides system status after cleanup
```

### Makefile Targets

- `run`: Start development server
- `docker-build`: Build Docker image
- `docker-run`: Start with Docker Compose
- `docker-stop`: Stop Docker containers
- `build-docker`: Execute build script
- `cleanup`: Execute cleanup script

---

## 📄 License

This project is created for educational purposes as part of zone01 assignment.

## 🤝 Contributors

- Dilhan (daslamac)
- Taha (tcavuslu)
- Thanos (tkalompr)
- Ypatios (ychaniot)

---

**Built with ❤️ using Go, SQLite, and Docker**