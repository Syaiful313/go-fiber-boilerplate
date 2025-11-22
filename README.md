# Go Fiber API Boilerplate

Template RESTful API menggunakan Go Fiber dengan fitur lengkap untuk pengembangan aplikasi modern. Proyek ini menyediakan struktur yang bersih dan scalable dengan best practices untuk API development.

## ✨ Fitur

- 🚀 **Go Fiber** - Framework web yang cepat dan minimalis
- 🗄️ **GORM** - ORM yang powerful untuk Go
- 🐘 **PostgreSQL** - Database relasional dengan Docker support
- 🔐 **JWT Authentication** - Autentikasi bearer token dengan expired 24 jam
- 🔒 **Password Security** - Hashing menggunakan bcrypt
- 📧 **Email System** - Forgot/reset password via SMTP
- 📁 **File Upload** - Upload gambar ke Cloudinary (dengan pembersihan aset lama)
- 🛡️ **Middleware** - CORS, Authentication, Error Handling, File Upload
- 🐳 **Docker Ready** - Development dengan Docker Compose
- 🔄 **Hot Reload** - Development dengan Air
- 📊 **Clean Architecture** - Struktur yang terorganisir (Controllers, Services, Models)

## 🏗️ Struktur Proyek

```
go-fiber-boilerplate/
├── cmd/                           # Entry point aplikasi
│   └── main.go
├── config/                        # Konfigurasi aplikasi
│   └── config.go
├── database/                      # Database setup
│   └── database.go
├── internal/
│   ├── controllers/              # HTTP handlers
│   │   ├── auth_controller.go
│   │   └── sample_controller.go
│   ├── middlewares/              # Middleware functions
│   │   ├── auth_middleware.go
│   │   ├── cors_middleware.go
│   │   ├── error_middleware.go
│   │   └── uploader_middleware.go
│   ├── models/                   # Data models & DTOs
│   │   ├── user.go
│   │   └── sample.go
│   ├── routes/                   # Route definitions
│   │   ├── routes.go
│   │   ├── auth_router.go
│   │   └── sample_router.go
│   └── services/                 # Business logic
│       ├── auth_service.go
│       ├── sample_service.go
│       └── cloudinary.service.go
├── pkg/                          # Shared packages
│   └── response/
│       └── response.go
├── utils/                        # Utility functions
│   ├── jwt.go
│   ├── password.go
│   ├── email.go
│   └── token.go
├── .env.example                  # Environment template
├── docker-compose.yml            # Docker services
├── Makefile                      # Build automation
└── README.md
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make (optional)

### Installation

1. **Clone repository:**

   ```bash
   git clone <repository-url>
   cd go-fiber-boilerplate
   ```

2. **Setup environment:**

   ```bash
   cp .env.example .env
   # Edit .env dengan konfigurasi Anda
   ```

3. **Start database:**

   ```bash
   make docker-up
   # atau
   docker-compose up -d
   ```

4. **Install dependencies:**

   ```bash
   go mod download
   ```

5. **Run aplikasi:**
   ```bash
   make dev
   # atau
   air
   ```

Server akan berjalan di `http://localhost:8000`

## ⚙️ Environment Variables

```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=admin
DB_NAME=go_fiber_db

# Server
PORT=8000
JWT_SECRET=your_jwt_secret_key_here
RESET_TOKEN_SECRET=your_reset_token_secret_here

# CORS
CORS_ALLOWED_ORIGINS=http://localhost:3000
CORS_ALLOW_CREDENTIALS=false

# Email (SMTP)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your_email@gmail.com
SMTP_PASSWORD=your_app_password
FROM_EMAIL=noreply@yourapp.com

# Frontend URL (untuk reset password)
FRONTEND_URL=http://localhost:3000

# Cloudinary (untuk upload gambar)
CLOUDINARY_CLOUD_NAME=your_cloud_name
CLOUDINARY_API_KEY=your_api_key
CLOUDINARY_API_SECRET=your_api_secret
```

## 📋 API Endpoints

### Health Check

```
GET /api/health
```

### Authentication

```
POST /auth/register          # Register user baru
POST /auth/login             # Login user
POST /auth/forgot-password   # Forgot password
POST /auth/reset-password    # Reset password
```

### Samples

```
GET    /samples              # Publik, daftar sample (pagination, max 50 per halaman, tanpa data user sensitif)
GET    /samples/:id          # Dilindungi, detail sample beserta data user
POST   /samples              # Dilindungi, create sample (JSON atau multipart)
PATCH  /samples/:id          # Dilindungi, update sample (JSON atau multipart)
DELETE /samples/:id          # Dilindungi, delete sample (hapus gambar Cloudinary bila ada)
```

## 📝 Request Examples

### Register User

```bash
curl -X POST http://localhost:8000/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

### Login

```bash
curl -X POST http://localhost:8000/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Create Sample (dengan token)

```bash
curl -X POST http://localhost:8000/samples \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Sample Title",
    "description": "Sample Description"
  }'
```

> Gunakan `multipart/form-data` jika ingin menyertakan gambar:

```bash
curl -X POST http://localhost:8000/samples \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "title=Sample Title" \
  -F "description=Sample Description" \
  -F "image=@/path/to/your-image.jpg"
```

### Forgot Password

```bash
curl -X POST http://localhost:8000/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com"
  }'
```

### Reset Password

```bash
curl -X POST http://localhost:8000/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "token": "reset_token_from_email",
    "new_password": "newpassword123"
  }'
```

## 🛠️ Development Commands

```bash
# Build aplikasi
make build

# Run aplikasi
make run

# Development dengan hot reload
make dev

# Run tests
make test

# Start Docker services
make docker-up

# Stop Docker services
make docker-down

# Clean build artifacts
make clean

# Full setup untuk development
make setup
```

## 🔧 Features Detail

### Authentication System

- JWT-based authentication
- Password hashing dengan bcrypt
- Email verification untuk forgot password
- Reset password dengan secure token

### File Upload System

- Upload gambar ke Cloudinary
- Validasi file type dan size (JPEG/PNG, batas ukuran dari middleware)
- Multiple image variants (thumbnail, small, medium, large)
- Secure file handling dan penghapusan aset lama saat update/delete sample

### Database Features

- Auto migration
- Soft delete support
- Relationship management
- Pagination support (perPage dibatasi, tidak ada mode `all`)

### Middleware

- **CORS**: Configurable cross-origin resource sharing
- **Auth**: JWT token validation
- **Error**: Centralized error handling
- **Upload**: File upload validation dan processing

### Email System

- SMTP support dengan template HTML
- Forgot password email
- Password reset confirmation email

## 🐳 Docker Support

Development environment dengan PostgreSQL dan Adminer:

```yaml
# docker-compose.yml menyediakan:
- PostgreSQL database (port 5432)
- PostgreSQL test database (port 5433)
- Adminer web interface (port 8080)
```

Access database via Adminer: `http://localhost:8080`

## 🧪 Testing

```bash
# Run semua tests
go test -v ./...

# Run tests dengan coverage
go test -v -cover ./...

# Run specific test
go test -v ./internal/services/
```

## 📚 Architecture

Proyek ini menggunakan **Clean Architecture** dengan separation of concerns:

- **Controllers**: Handle HTTP requests/responses
- **Services**: Business logic dan validations
- **Models**: Data structures dan database models
- **Utils**: Reusable utility functions
- **Middleware**: Cross-cutting concerns
- **Config**: Application configuration

## 🔄 Migration dari Express

Template ini merupakan equivalent dari Express.js boilerplate dengan keuntungan:

- **Performance**: Go Fiber lebih cepat dari Express
- **Memory**: Konsumsi memory yang efisien
- **Concurrency**: Built-in goroutines
- **Deployment**: Single binary executable
- **Type Safety**: Static typing

## 🎯 Best Practices

- Environment-based configuration
- Proper error handling dengan custom error types
- Structured logging
- Input validation
- Security headers
- Database connection pooling
- Graceful shutdown handling
- Pagination dan listing publik dibatasi untuk mencegah data dump; data user sensitif tidak dikirim di endpoint publik

## 📄 Response Format

API menggunakan consistent response format:

```json
{
  "message": "Success message",
  "data": {
    // response data
  }
}
```

Error responses:

```json
{
  "error": "Error message"
}
```

## 🚦 Status Codes

- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict
- `500` - Internal Server Error
