# Go Fiber Boilerplate

Boilerplate untuk membangun RESTful API menggunakan Go Fiber, GORM, dan PostgreSQL. Proyek ini merupakan migrasi dari Express Prisma Boilerplate ke teknologi Go.

## Fitur

- 🚀 **Go Fiber** - Framework web yang cepat dan minimalis
- 🗄️ **GORM** - ORM yang powerful untuk Go
- 🐘 **PostgreSQL** - Database relasional yang robust
- 🔐 **JWT Authentication** - Sistem autentikasi yang aman
- 🔒 **Password Hashing** - Menggunakan bcrypt
- 📝 **CRUD Operations** - Operasi Create, Read, Update, Delete
- 🐳 **Docker Support** - Containerization dengan Docker Compose
- 🔄 **Hot Reload** - Development dengan Air
- 📋 **Middleware** - CORS, Authentication, Error Handling
- 🧪 **Testing Ready** - Struktur untuk unit dan integration testing

## Struktur Proyek

```
go-fiber-boilerplate/
├── cmd/                    # Entry point aplikasi
│   └── main.go
├── config/                 # Konfigurasi aplikasi
│   └── config.go
├── database/              # Konfigurasi database
│   └── database.go
├── handlers/              # HTTP handlers
│   ├── auth.go
│   └── sample.go
├── middleware/            # Middleware functions
│   ├── auth.go
│   ├── cors.go
│   └── error.go
├── models/               # Data models
│   ├── user.go
│   └── sample.go
├── routes/               # Route definitions
│   └── routes.go
├── utils/                # Utility functions
│   ├── jwt.go
│   └── password.go
├── .env.example          # Environment variables template
├── docker-compose.yml    # Docker configuration
├── Makefile             # Build automation
└── README.md
```

## Instalasi

### Prasyarat

- Go 1.21 atau lebih baru
- Docker dan Docker Compose
- Make (opsional, untuk menggunakan Makefile)

### Langkah Instalasi

1. **Clone repository:**
   ```bash
   git clone <repository-url>
   cd go-fiber-boilerplate
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Setup environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env sesuai dengan konfigurasi Anda
   ```

4. **Start database dengan Docker:**
   ```bash
   docker-compose up -d
   ```

5. **Run aplikasi:**
   ```bash
   go run cmd/main.go
   ```

   Atau menggunakan Makefile:
   ```bash
   make start
   ```

## Environment Variables

Buat file `.env` berdasarkan `.env.example`:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=admin
DB_NAME=go_fiber_db

# Server Configuration
PORT=8000

# JWT Configuration
JWT_SECRET=your_jwt_secret_key_here
```

## API Endpoints

### Authentication

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register user baru |
| POST | `/api/v1/auth/login` | Login user |
| GET | `/api/v1/profile` | Get user profile (protected) |

### Samples

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/samples` | Get all samples (protected) |
| GET | `/api/v1/samples/:id` | Get sample by ID (protected) |
| POST | `/api/v1/samples` | Create new sample (protected) |
| PUT | `/api/v1/samples/:id` | Update sample (protected) |
| DELETE | `/api/v1/samples/:id` | Delete sample (protected) |

### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | Health check endpoint |

## Contoh Request

### Register User

```bash
curl -X POST http://localhost:8000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "firstName": "John",
    "lastName": "Doe"
  }'
```

### Login

```bash
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Create Sample (dengan token)

```bash
curl -X POST http://localhost:8000/api/v1/samples \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "Sample Title",
    "description": "Sample Description"
  }'
```

## Development

### Hot Reload dengan Air

1. **Install Air:**
   ```bash
   go install github.com/cosmtrek/air@latest
   ```

2. **Run dengan hot reload:**
   ```bash
   air
   ```

   Atau menggunakan Makefile:
   ```bash
   make dev
   ```

### Makefile Commands

```bash
make build          # Build aplikasi
make run            # Run aplikasi
make dev            # Run dengan hot reload
make test           # Run tests
make clean          # Clean build artifacts
make docker-up      # Start Docker containers
make docker-down    # Stop Docker containers
make setup          # Full development setup
```

## Testing

```bash
# Run all tests
go test -v ./...

# Atau menggunakan Makefile
make test
```

## Docker

### Start Services

```bash
docker-compose up -d
```

### Stop Services

```bash
docker-compose down
```

### Database Access

- **Host:** localhost
- **Port:** 5432
- **Database:** go_fiber_db
- **Username:** postgres
- **Password:** admin

## Migrasi dari Express Prisma

Proyek ini merupakan migrasi dari Express Prisma Boilerplate dengan perubahan berikut:

### Teknologi Stack

| Express Prisma | Go Fiber |
|----------------|----------|
| Node.js + Express | Go + Fiber |
| Prisma ORM | GORM |
| TypeScript | Go |
| npm/yarn | Go modules |

### Fitur yang Dipertahankan

- ✅ JWT Authentication
- ✅ Password Hashing
- ✅ CRUD Operations
- ✅ Database Relations
- ✅ Middleware Support
- ✅ Environment Configuration
- ✅ Docker Support
- ✅ Error Handling

### Keuntungan Migrasi

- **Performance:** Go Fiber lebih cepat dibanding Express
- **Memory Usage:** Konsumsi memory yang lebih efisien
- **Concurrency:** Built-in goroutines untuk concurrent processing
- **Deployment:** Single binary deployment
- **Type Safety:** Static typing tanpa runtime overhead

## Contributing

1. Fork repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## License

Distributed under the MIT License. See `LICENSE` for more information.

## Contact

Your Name - your.email@example.com

Project Link: [https://github.com/yourusername/go-fiber-boilerplate](https://github.com/yourusername/go-fiber-boilerplate)

