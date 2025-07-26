# Perbaikan Error CORS

## Masalah

Error yang terjadi:
```
panic: [CORS] Insecure setup, 'AllowCredentials' is set to true, and 'AllowOrigins' is set to a wildcard.
```

## Penyebab

Browser modern tidak mengizinkan kombinasi:
- `AllowCredentials: true` (mengizinkan pengiriman cookies/credentials)
- `AllowOrigins: "*"` (mengizinkan semua origin)

Kombinasi ini dianggap tidak aman karena dapat menyebabkan serangan CSRF.

## Solusi yang Diterapkan

### 1. Konfigurasi CORS yang Fleksibel

File `config/config.go` sekarang mendukung konfigurasi CORS melalui environment variables:

```go
type Config struct {
    // ... other fields
    AllowedOrigins   string
    AllowCredentials bool
}
```

### 2. Middleware CORS yang Aman

File `middleware/cors.go` sekarang memiliki logika untuk mencegah kombinasi yang tidak aman:

```go
func CORSMiddleware(cfg *config.Config) fiber.Handler {
    // If credentials are allowed, origins cannot be wildcard
    if cfg.AllowCredentials && cfg.AllowedOrigins == "*" {
        // For development, disable credentials when using wildcard
        return cors.New(cors.Config{
            AllowOrigins:     "*",
            AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
            AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
            AllowCredentials: false,
        })
    }

    return cors.New(cors.Config{
        AllowOrigins:     cfg.AllowedOrigins,
        AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
        AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
        AllowCredentials: cfg.AllowCredentials,
    })
}
```

### 3. Environment Variables

Tambahkan ke file `.env`:

```env
# CORS Configuration
CORS_ALLOWED_ORIGINS=*
CORS_ALLOW_CREDENTIALS=false
```

## Opsi Konfigurasi

### Development (Default)
```env
CORS_ALLOWED_ORIGINS=*
CORS_ALLOW_CREDENTIALS=false
```

### Production dengan Frontend Spesifik
```env
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com
CORS_ALLOW_CREDENTIALS=true
```

### Production dengan Multiple Domains
```env
CORS_ALLOWED_ORIGINS=https://app.yourdomain.com,https://admin.yourdomain.com
CORS_ALLOW_CREDENTIALS=true
```

## Testing

Setelah perbaikan, aplikasi dapat dijalankan tanpa error:

```bash
go run cmd/main.go
```

Output yang diharapkan:
```
2025/07/25 22:22:31 Database connected successfully
2025/07/25 22:22:31 Database migration completed
2025/07/25 22:22:31 Server starting on port 8000
```

## Rekomendasi

1. **Development**: Gunakan `CORS_ALLOW_CREDENTIALS=false` dengan `CORS_ALLOWED_ORIGINS=*`
2. **Production**: Gunakan `CORS_ALLOW_CREDENTIALS=true` dengan origin spesifik
3. **Testing**: Selalu test CORS dengan frontend yang sebenarnya

## Security Notes

- Jangan pernah menggunakan `AllowOrigins: "*"` dengan `AllowCredentials: true` di production
- Selalu specify origin yang spesifik di production
- Gunakan HTTPS di production untuk keamanan tambahan

