# Forum Web Application

A modern forum web application built with Go, featuring multi-authentication system (Email/Password + OAuth), post management, comment system, and reaction functionality. This project demonstrates clean architecture principles with a layered approach including handlers, services, repositories, and middleware.

## 📋 Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Database](#database)
- [Authentication & Middleware](#authentication--middleware)
- [OAuth Integration](#oauth-integration)
- [Project Structure](#project-structure)
- [Features](#features)
- [Installation & Setup](#installation--setup)
- [OAuth Configuration](#oauth-configuration)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Testing Guide](#testing-guide)
- [Docker Deployment](#docker-deployment)
- [Development](#development)

## 🌟 Overview

This forum application provides a complete community platform where users can:

- Register and authenticate securely using multiple methods (Email/Password, Google OAuth, GitHub OAuth)
- Create and manage posts with categories
- Comment on posts and interact with content
- Like/dislike posts and comments
- View personalized profile pages with activity tracking
- Browse content by categories with pagination
- Link multiple authentication methods to a single account

The application follows modern web development practices with server-side rendering using Go's `html/template` package and implements a clean separation of concerns through a layered architecture with comprehensive OAuth integration.

## 🏗️ Architecture

The application follows a clean architecture pattern with the following layers:

```
├── cmd/web/           # Application entry point
├── internal/          # Private application code
│   ├── handlers/      # HTTP request handlers (including OAuth)
│   ├── services/      # Business logic layer (auth + OAuth services)
│   ├── repository/    # Data access layer
│   ├── middleware/    # HTTP middleware
│   ├── models/        # Data models
│   ├── config/        # OAuth configuration
│   └── utils/         # Utility functions
├── pkg/               # Public packages
├── templates/         # HTML templates (with OAuth UI)
└── static/           # Static assets (CSS, JS, images)
```

### Layer Responsibilities

- **Handlers**: Handle HTTP requests/responses, input validation, template rendering, and OAuth flows
- **Services**: Implement business logic including traditional authentication and OAuth processing
- **Repository**: Manage database operations and data persistence
- **Middleware**: Handle cross-cutting concerns like authentication and request processing
- **Models**: Define data structures and domain entities
- **Config**: Manage OAuth application configurations

## 🗄️ Database

The application uses SQLite as its database engine with automatic initialization and schema management, enhanced with OAuth support.

### Database Configuration

- **Database File**: `forum.db` (created automatically)
- **Schema Location**: `internal/database/schema.sql`
- **Driver**: `github.com/mattn/go-sqlite3`
- **Connection Pool**: Max 10 concurrent connections
- **Foreign Keys**: Enabled for referential integrity

### Database Tables

#### Enhanced User Table (OAuth Support)

```sql
CREATE TABLE user (
    id TEXT PRIMARY KEY,              -- UUID
    email TEXT NOT NULL UNIQUE,      -- User email
    username TEXT NOT NULL UNIQUE,   -- Display name
    password TEXT,                   -- Bcrypt hashed password (nullable for OAuth-only users)
    oauth_provider TEXT, --'google', 'github', or NULL for regular users
    oauth_provider_id TEXT, -- OAuth provider's unique user ID
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(oauth_provider, oauth_provider_id) -- Ensure each OAuth account can only be linked once
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

### Enhanced Database Indexes (OAuth Optimized)

Performance-optimized indexes including OAuth lookups:

```sql
CREATE INDEX idx_posts_user_id ON post(user_id);
CREATE INDEX idx_comments_post_id ON comment(post_id);
CREATE INDEX idx_comments_user_id ON comment(user_id);
CREATE INDEX idx_reactions_post_id ON reaction(post_id);
CREATE INDEX idx_reactions_comment_id ON reaction(comment_id);
CREATE INDEX idx_reactions_user_id ON reaction(user_id);
CREATE INDEX idx_sessions_user_id ON session(user_id);
-- OAuth-specific indexes
CREATE INDEX idx_users_google_id ON user(google_id);
CREATE INDEX idx_users_github_id ON user(github_id);
```

### Database Initialization

The database initializes automatically on first run:

1. **Connection Setup**: Opens SQLite connection with foreign key constraints enabled
2. **Schema Creation**: Executes SQL schema from `internal/database/schema.sql`
3. **Data Population**: Seeds initial categories (General Discussion, Tech News, etc.)
4. **OAuth Support**: Creates indexes for OAuth ID lookups
5. **Validation**: Verifies database connectivity and structure

## 🔐 Authentication & Middleware

### Multi-Authentication System

The application implements a comprehensive authentication system supporting multiple methods:

#### Traditional Authentication

- **Password Security**: Uses bcrypt with default cost factor (10)
- **Validation**: Enforces strong password requirements
- **Storage**: Never stores plaintext passwords

#### OAuth Authentication

- **Google OAuth 2.0**: Full integration with Google Sign-In
- **GitHub OAuth**: Complete GitHub authentication flow
- **Account Linking**: Connect OAuth accounts to existing users
- **Username Selection**: Custom username choice for new OAuth users

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
- OAuth session support

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
```

### Security Features

- **CSRF Protection**: SameSite cookie attribute prevents CSRF attacks
- **Session Hijacking Prevention**: Server-side session validation
- **Secure Cookies**: HTTPOnly and Secure flags when using HTTPS
- **Input Sanitization**: HTML escaping for user-generated content
- **SQL Injection Protection**: Parameterized queries throughout
- **OAuth State Parameter**: CSRF protection for OAuth flows
- **Token Validation**: Secure OAuth token exchange and validation

## 🔗 OAuth Integration

### Supported Providers

#### Google OAuth 2.0

- **Scope**: `openid profile email`
- **User Info**: Name, email, profile picture
- **Account Linking**: Automatic linking via email or manual selection

#### GitHub OAuth

- **Scope**: `user:email`
- **User Info**: Username, email, profile information
- **Account Linking**: Automatic linking via email or manual selection

### OAuth Configuration

The application uses a `config.json` file for OAuth settings:

```json
{
  "google_client_id": "GOOGLE CLIENT ID",
  "google_client_secret": "GOOGLE SECRET KEY",
  "github_client_id": "GITHUB CLIENT ID",
  "github_client_secret": "GITHUB SECRET KEY",
  "app_url": "http://localhost:8080"
}
```

### OAuth User Flow

#### New User Registration via OAuth

1. User clicks "Sign in with Google/GitHub" on login/signup page
2. Redirected to OAuth provider for authentication
3. User authorizes the application
4. Callback handler receives user information
5. If email doesn't exist in database:
   - Create new user record with OAuth ID
   - Redirect to username selection page
   - User selects unique username
   - Account creation completed
6. User is logged in and redirected to homepage

#### Existing User Login via OAuth

1. User clicks "Sign in with Google/GitHub"
2. OAuth flow completes
3. If email matches existing user:
   - Link OAuth ID to existing account
   - Log user in immediately
   - Redirect to homepage

#### Account Linking

1. Authenticated user can link additional OAuth accounts
2. OAuth flow completes for logged-in user
3. OAuth ID is added to existing user record
4. User can now login via any linked method

## ✨ Features

### Core Functionality

- **Multi-Authentication System**: Email/password, Google OAuth, and GitHub OAuth
- **OAuth Integration**: Social login with Google and GitHub, with username selection for new users
- **Account Linking**: Connect multiple authentication methods to a single account
- **User Registration & Login**: Secure account creation and authentication with multiple methods
- **Post Management**: Create, view, and categorize posts
- **Comment System**: Threaded discussions on posts
- **Reaction System**: Like/dislike posts and comments
- **Category Filtering**: Browse posts by topic categories
- **User Profiles**: View personal activity and statistics with linked account information

### UI/UX Features

- **Responsive Design**: Mobile-friendly interface
- **OAuth Button Integration**: Seamless social login buttons
- **Username Selection UI**: Guided username setup for OAuth users
- **Account Linking Interface**: Manage connected accounts
- **Pagination**: Efficient content browsing
- **Real-time Updates**: Dynamic content loading
- **Error Handling**: User-friendly error pages
- **Form Validation**: Client and server-side validation

### Technical Features

- **Session Management**: Secure user sessions with UUID-based tokens
- **OAuth State Management**: CSRF protection for OAuth flows
- **Database Migrations**: Automatic schema management with OAuth support
- **Input Sanitization**: XSS protection
- **Performance Optimization**: Efficient database queries with pagination and OAuth indexes
- **Docker Support**: Containerized deployment
- **Multi-Provider OAuth**: Support for multiple OAuth providers
- **Account Consolidation**: Merge OAuth accounts with existing users
- **Flexible Authentication**: Users can login via any linked method

## 🚀 Installation & Setup

### Prerequisites

- **Go**: Version 1.19 or higher
- **SQLite**: Included with Go SQLite driver
- **Docker** (optional): For containerized deployment
- **OAuth Apps**: Google and/or GitHub OAuth application credentials

### Local Development Setup

1. **Clone the repository:**

   ```bash
   git clone https://platform.zone01.gr/git/ychaniot/forum-authentication.git
   cd forum-app
   ```

2. **Install dependencies:**

   ```bash
   go mod download
   ```

3. **Configure OAuth (see OAuth Configuration section)**

4. **Run the application:**

   ```bash
   go run cmd/web/main.go
   ```

5. **Access the application:**
   ```
   http://localhost:8080
   ```

The database will be automatically created and initialized on first run.

## 🔧 OAuth Configuration

### Setting Up OAuth Applications

#### Google OAuth Setup

1. **Go to Google Cloud Console**: https://console.cloud.google.com/
2. **Create a new project** or select existing one
3. **Enable Google+ API** and **Google People API**
4. **Create OAuth 2.0 credentials**:
   - Application type: Web application
   - Authorized redirect URIs: `http://localhost:8080/auth/google/callback`
   - For production: `https://yourdomain.com/auth/google/callback`
5. **Note down Client ID and Client Secret**

#### GitHub OAuth Setup

1. **Go to GitHub Settings**: https://github.com/settings/developers
2. **Click "New OAuth App"**
3. **Fill in application details**:
   - Application name: Your Forum App
   - Homepage URL: `http://localhost:8080`
   - Authorization callback URL: `http://localhost:8080/auth/github/callback`
   - For production: `https://yourdomain.com/auth/github/callback`
4. **Note down Client ID and Client Secret**

## 📖 Usage

### Using the Makefile

The project includes a comprehensive Makefile for common operations:

```bash
# Run the application locally
make run

# Build Docker image
make docker-build

# Start application with Docker Compose
make docker-run

# Stop Docker containers
make docker-stop

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

### Authentication Usage

#### Traditional Authentication

1. **Sign Up**: Create account with email and password
2. **Login**: Authenticate with email and password
3. **Password Requirements**: Minimum 8 characters with complexity rules

#### OAuth Authentication

1. **Sign Up/Login**: Click "Sign in with Google" or "Sign in with GitHub"
2. **Authorization**: Approve application access
3. **Username Selection**: Choose unique username (new users only)
4. **Account Linking**: Connect OAuth to existing account via email match

#### Account Management

1. **Profile Page**: View all linked authentication methods
2. **Link Accounts**: Add additional OAuth providers to existing account
3. **Flexible Login**: Use any linked method to access account

## 📋 API Endpoints

### Public Routes

- `GET /` - Homepage with post listings
- `GET /login` - Login page with OAuth options
- `GET /signup` - Registration page with OAuth options
- `POST /login` - Authenticate user (traditional)
- `POST /signup` - Register new user (traditional)
- `GET /posts/{id}` - View individual post
- `GET /categories/{id}/posts` - Posts by category

### OAuth Routes

- `GET /auth/google` - Initiate Google OAuth flow
- `GET /auth/github` - Initiate GitHub OAuth flow
- `GET /auth/google/callback` - Handle Google OAuth callback
- `GET /auth/github/callback` - Handle GitHub OAuth callback
- `GET /auth/username` - Username selection page for new OAuth users
- `POST /auth/set-username` - Process username selection for OAuth users

### Protected Routes (Require Authentication)

- `POST /logout` - User logout (any authentication method)
- `GET /profile` - User profile page with linked accounts
- `GET /create-post` - Create post form
- `POST /create-post` - Submit new post
- `POST /posts/{id}/comments` - Add comment
- `POST /post/{id}/{reaction}` - React to post
- `POST /comment/{id}/{reaction}` - React to comment

### Static Assets

- `/static/` - CSS, JavaScript, and image files (including OAuth provider icons)

## 🧪 Testing Guide

### Manual Testing Steps

#### Google OAuth Testing

1. **New User Registration**:

   ```
   1. Go to /login
   2. Click "Sign in with Google"
   3. Complete Google authentication
   4. Should redirect to username selection page
   5. Enter unique username
   6. Should create account and login
   7. Verify user appears in database with google_id
   ```

2. **Existing User Login**:

   ```
   1. Create user via traditional signup
   2. Go to /login
   3. Click "Sign in with Google" (use same email)
   4. Should link Google account and login
   5. Verify google_id is added to existing user record
   ```

3. **Account Linking**:
   ```
   1. Login with traditional method
   2. Go to /profile
   3. Click "Link Google Account"
   4. Complete OAuth flow
   5. Should return to profile with Google account linked
   ```

#### GitHub OAuth Testing

1. **New User Registration**:

   ```
   1. Go to /signup
   2. Click "Sign in with GitHub"
   3. Complete GitHub authentication
   4. Should redirect to username selection page
   5. Enter unique username
   6. Should create account and login
   7. Verify user appears in database with github_id
   ```

2. **Existing User Login**:
   ```
   1. Create user via traditional signup
   2. Go to /login
   3. Click "Sign in with GitHub" (use same email)
   4. Should link GitHub account and login
   5. Verify github_id is added to existing user record
   ```

#### Account Linking Verification

1. **Multiple OAuth Providers**:

   ```
   1. Create account via Google OAuth
   2. Link GitHub account via profile page
   3. Logout
   4. Login via GitHub - should access same account
   5. Login via Google - should access same account
   6. Verify both google_id and github_id in database
   ```

2. **Traditional + OAuth**:
   ```
   1. Create account via email/password
   2. Link Google account
   3. Link GitHub account
   4. Verify can login via any of the 3 methods
   5. All should access the same user account
   ```

#### User Rights Verification

1. **OAuth User Capabilities**:

   ```
   1. Create account via OAuth
   2. Create posts - should work
   3. Add comments - should work
   4. React to posts/comments - should work
   5. All functionality should be identical to traditional users
   ```

2. **Profile Information**:
   ```
   1. Login via OAuth
   2. Go to /profile
   3. Should display:
      - Username (chosen during signup)
      - Email (from OAuth provider)
      - Linked accounts section
      - User activity/statistics
   ```

### Expected Behavior

#### OAuth Flow Success Indicators

- ✅ Successful redirect to OAuth provider
- ✅ Proper callback handling with authorization code
- ✅ User information retrieval from OAuth provider
- ✅ Account creation or linking in database
- ✅ Session creation and cookie setting
- ✅ Redirect to appropriate page (username selection or homepage)

#### Error Handling Verification

- ❌ Invalid OAuth configuration → Error message displayed
- ❌ OAuth provider denial → Graceful redirect to login with message
- ❌ Username already taken → Username selection page with error
- ❌ Database errors → Appropriate error page displayed

## 🐳 Docker Deployment

### Using Docker Compose (Recommended)

1. **Ensure OAuth configuration is set up** (config.json or environment variables)

2. **Start the application:**

   ```bash
   make docker-run
   # or
   docker-compose up -d
   ```

3. **View logs:**

   ```bash
   docker-compose logs -f
   ```

4. **Stop the application:**
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

# Run container with OAuth configuration
docker run -d --name forum \
  -p 8080:8080 \
  -v $(pwd)/config.json:/app/config.json \
  forum-app

# View logs
docker logs forum

# Stop container
docker stop forum
```

### Docker Configuration

The application includes:

- **Multi-stage Dockerfile** for optimized image size
- **Health checks** for container monitoring
- **Volume mounting** for persistent data storage and OAuth configuration
- **Environment configuration** for deployment flexibility
- **OAuth configuration mounting** for secure credential management

## 🔧 Development

### Adding New Features

1. **Models**: Define data structures in `internal/models/`
2. **Repository**: Add database operations in `internal/repository/`
3. **Services**: Implement business logic in `internal/services/`
4. **Handlers**: Create HTTP handlers in `internal/handlers/`
5. **Routes**: Register routes in `internal/routes/routes.go`
6. **Templates**: Add HTML templates in `templates/`

### Adding New OAuth Providers

1. **Add provider configuration** to `config.json`
2. **Create provider service** in `internal/services/oauth_service.go`
3. **Add OAuth handlers** in `internal/handlers/oauth_handlers.go`
4. **Update user model** with provider ID field
5. **Add database migration** for provider ID column
6. **Update templates** with provider-specific buttons
7. **Add provider routes** in `internal/routes/routes.go`

### Database Migrations

To modify the database schema:

1. Update `internal/database/schema.sql`
2. Delete existing `forum.db` file
3. Restart application for automatic recreation

### Environment Variables

The application supports the following environment variables:

- `PORT`: Server port (default: 8080)
- `DB_PATH`: Database file path (default: forum.db)
- `GOOGLE_CLIENT_ID`: Google OAuth client ID
- `GOOGLE_CLIENT_SECRET`: Google OAuth client secret
- `GITHUB_CLIENT_ID`: GitHub OAuth client ID
- `GITHUB_CLIENT_SECRET`: GitHub OAuth client secret

## 🛠️ Scripts and Automation

### Build Script (`build.sh`)

Automated Docker build and deployment:

```bash
#!/bin/bash
# Features:
- Container cleanup and removal
- Docker image building with error handling
- Container deployment with port mapping
- OAuth configuration mounting
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

This project is created for educational purposes as part of a school assignment.

## 🤝 Contributors

Dilhan (daslamac)
Taha (tcavuslu)
Thanos (tkalompr)
Ypatios (ychaniot)

---

**Built with ❤️ using Go, SQLite, OAuth 2.0, and Docker**
