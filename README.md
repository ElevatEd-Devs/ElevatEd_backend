# 🎓 ElevatEd Backend

Backend repository for the ElevatEd Education Platform - a comprehensive learning management system built with Go (Golang) and PostgreSQL.

[![Go Version](https://img.shields.io/badge/Go-1.24.5-blue.svg)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Latest-blue.svg)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/Docker-Enabled-blue.svg)](https://www.docker.com/)

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Getting Started](#getting-started)
- [Development](#development)
- [Docker Setup](#docker-setup)
- [API Documentation](#api-documentation)
- [Project Structure](#project-structure)
- [Testing](#testing)
- [Contributing](#contributing)
- [Team & Support](#team--support)

## 🌟 Overview

ElevatEd is a modern education platform designed to streamline learning experiences. The backend provides a robust RESTful API built with Go, offering high performance, scalability, and reliability for educational services including user management, appointments, events, and grading systems.

## ✨ Features

### Core Functionality
- **Authentication & Authorization** - Secure user authentication with JWT tokens
- **User Management** - Comprehensive user profile and account management
- **Appointment System** - Schedule and manage educational appointments
- **Event Management** - Create and organize educational events
- **Grading System** - Track and manage student grades and assessments
- **RESTful API** - Well-structured REST endpoints for all operations

### Technical Features
- **High Performance** - Built with Go for optimal speed and efficiency
- **Database Integration** - PostgreSQL with proper migrations and transactions
- **Live Reloading** - Hot reload during development with Air
- **Containerization** - Docker support for consistent environments
- **Network Security** - Tailscale integration for secure networking
- **Comprehensive Testing** - Built-in API testing with .http files

## 🛠️ Tech Stack

### Backend
- **Go (Golang)** v1.24.5 - High-performance backend language
- **PostgreSQL** - Robust relational database
- **Chi Router** (likely) - HTTP routing and middleware

### DevOps & Infrastructure
- **Docker** - Containerization and deployment
- **Air** - Live reload for Go apps
- **Tailscale** - Secure networking and remote access

### Development Tools
- **Postman / Bruno** - API testing
- **HTTP Files** - In-editor API testing

## 🚀 Getting Started

### Prerequisites

Ensure you have Go installed on your system:

```bash
# Check if Go is installed
go version

# If not installed, visit: https://go.dev/doc/install
```

**Required Go Version:** v1.24.5

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/ElevatEd-Devs/ElevatEd_backend.git
   cd ElevatEd_backend
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   - Create a `.env` file in the root directory
   - Contact the backend team on Discord for required environment variables
   - Typical variables include:
     ```
     DATABASE_URL=postgresql://user:password@localhost:5432/elevated_db
     JWT_SECRET=your_jwt_secret_key
     PORT=3000
     ```

4. **Install Air for live reloading** (optional but recommended)
   ```bash
   go install github.com/air-verse/air@latest
   ```

### Running the Application

#### With Live Reloading (Recommended)

```bash
air
```

This starts the server with hot reload enabled, automatically restarting when you make code changes.

#### Standard Run

```bash
go run main.go
```

The server will start on `http://localhost:3000` (or your configured port).

### Ubuntu Users: Air Setup

If you encounter issues running Air on Ubuntu:

```bash
# Add Go binaries to PATH
export PATH=$PATH:$HOME/go/bin

# Reload shell configuration
source ~/.bashrc

# Run Air
air
```

Consider adding the export command to your `~/.bashrc` or `~/.zshrc` for persistence.

## 🐳 Docker Setup

### Building and Running with Docker

1. **Build the Docker image**
   ```bash
   docker build -t elevated_backend .
   ```

2. **Run the container**
   ```bash
   docker run -p 3000:3000 --rm \
     -v $(pwd):/app \
     -v /app/tmp \
     --name docker-air \
     elevated_backend
   ```

This setup:
- Maps port 3000 from container to host
- Mounts your code for live reloading
- Automatically removes the container on exit
- Names the container `docker-air` for easy management

### Docker Compose (if available)

```bash
docker-compose up
```

## 📚 API Documentation

The API is organized into several modules with comprehensive documentation:

### Available Endpoints

#### Authentication
- **Docs:** [Auth Endpoints](./docs/authEndPoints.md)
- User registration, login, logout
- JWT token management
- Password reset and recovery

#### Appointments
- **Docs:** [Appointment Endpoints](./docs/appointmentEndPoints.md)
- Schedule appointments
- View/update/cancel appointments
- Appointment notifications

#### Events
- **Docs:** [Event Endpoints](./docs/eventEndPoints.md)
- Create and manage educational events
- Event registration and attendance
- Event categories and filtering

#### Grades
- **Docs:** [Grade Endpoints](./docs/gradeEndPoint.md)
- Grade submission and retrieval
- Grade calculations and statistics
- Student performance tracking

### Testing APIs

#### Using HTTP Files

The repository includes an `apitest.http` file for quick API testing:

1. Open `apitest.http` in VS Code with the REST Client extension
2. Click "Send Request" above any endpoint
3. View responses directly in the editor

![HTTP Test Example](https://private-user-images.githubusercontent.com/107510378/464916754-c9a77824-75f4-438a-9366-96309653422a.png)

#### Using API Clients

**Recommended Tools:**
- [Postman](https://www.postman.com/) - Feature-rich API testing platform
- [Bruno](https://www.usebruno.com/) - Open-source alternative to Postman

## 📁 Project Structure

```
ElevatEd_backend/
├── database/           # Database connection and migrations
├── docs/              # API documentation
├── functions/         # Business logic and utilities
├── handler/           # HTTP request handlers
├── router/            # Route definitions and middleware
├── structs/           # Data models and structures
├── .air.toml          # Air configuration for live reload
├── .gitignore         # Git ignore rules
├── Dockerfile         # Docker configuration
├── apitest.http       # HTTP file for API testing
├── go.mod             # Go module definition
├── go.sum             # Go dependency checksums
├── main.go            # Application entry point
└── README.md          # This file
```

### Directory Breakdown

- **`database/`** - PostgreSQL connection, schema, and migration files
- **`handler/`** - HTTP handlers that process requests and return responses
- **`router/`** - Route configuration with middleware (auth, CORS, logging)
- **`functions/`** - Reusable business logic and helper functions
- **`structs/`** - Go structs for request/response models and database entities
- **`docs/`** - Detailed API documentation for each module

### API Testing

Use the included `apitest.http` file or API testing tools to validate endpoints:

```bash
# Example using curl
curl -X GET http://localhost:3000/api/v1/health
```

## 🤝 Contributing

We welcome contributions from the community! Here's how to get started:

### Workflow

1. **Pick an issue** from the [Issues tab](https://github.com/ElevatEd-Devs/ElevatEd_backend/issues)
2. **Create a branch** for your work
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/issue-name
   ```

![Branch Creation](https://private-user-images.githubusercontent.com/107510378/464909681-bb721a49-2d2d-44e2-b26a-d2e10ca01036.png)

3. **Make your changes** following Go best practices
   - Write clean, idiomatic Go code
   - Add tests for new functionality
   - Update documentation as needed

4. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add new feature" # or "fix: resolve issue"
   ```

5. **Push to your branch**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Create a Pull Request**
   - Provide a clear description of changes
   - Reference related issues
   - Alert the team on Discord for review

### Code Standards

- Follow Go's [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use meaningful variable and function names
- Write comments for complex logic
- Keep functions small and focused
- Run `go fmt` before committing
- Ensure all tests pass

### Review Process

- All PRs require peer review before merging
- Address review comments promptly
- Keep the team updated on Discord
- Be open to feedback and suggestions

### Getting Help

1. **Issues** - Report bugs or request features via [GitHub Issues](https://github.com/ElevatEd-Devs/ElevatEd_backend/issues)
2. **Documentation** - Check the `/docs` folder for detailed guides

### Active Contributors

Check our [Contributors page](https://github.com/ElevatEd-Devs/ElevatEd_backend/graphs/contributors) to see everyone who's helped build ElevatEd!

## 🔧 Troubleshooting

### Common Issues

**Port already in use:**
```bash
# Find process using port 3000
lsof -i :3000

# Kill the process
kill -9 <PID>
```

**Air not found:**
```bash
# Ensure Go bin is in PATH
export PATH=$PATH:$(go env GOPATH)/bin

# Reinstall Air
go install github.com/air-verse/air@latest
```

**Database connection errors:**
- Verify PostgreSQL is running
- Check `.env` file configuration
- Ensure database credentials are correct

**Module download issues:**
```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download
```

## 🔗 Related Repositories

- [ElevatEd Frontend](https://github.com/ElevatEd-Devs/ElevatEd_frontend) - React frontend

## 📊 Project Status

- **Active Development** ✅
- **Open Issues:** [View Issues](https://github.com/ElevatEd-Devs/ElevatEd_backend/issues)

---

**Built with by the ElevatEd Development Team**

For more information, visit our [GitHub Organization](https://github.com/ElevatEd-Devs)
