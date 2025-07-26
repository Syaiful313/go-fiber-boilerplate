# Forgot Password & Reset Password Feature Guide

## Overview

Fitur forgot password dan reset password telah berhasil diimplementasikan di template Go Fiber Boilerplate. Fitur ini memungkinkan pengguna untuk mereset password mereka melalui email verification.

## Features Added

### 1. Forgot Password
- **Endpoint**: `POST /api/v1/auth/forgot-password`
- **Functionality**: Mengirim email reset password ke pengguna
- **Security**: Menggunakan JWT token dengan expiry 1 jam

### 2. Reset Password
- **Endpoint**: `POST /api/v1/auth/reset-password`
- **Functionality**: Mereset password menggunakan token dari email
- **Security**: Validasi token dan update password dengan hash baru

## New Files Added

### 1. `utils/email.go`
- Fungsi untuk mengirim email menggunakan SMTP
- Template email untuk reset password
- Validasi email format

### 2. `utils/token.go`
- Generate JWT token khusus untuk reset password
- Validasi reset password token
- Generate random token untuk keamanan tambahan

### 3. Updated Files
- `config/config.go` - Tambah konfigurasi email dan frontend URL
- `models/user.go` - Tambah request/response models untuk forgot/reset password
- `handlers/auth.go` - Tambah handler untuk forgot/reset password
- `routes/routes.go` - Tambah rute baru
- `.env.example` dan `.env` - Tambah konfigurasi email

## Configuration

### Environment Variables

Tambahkan konfigurasi berikut ke file `.env`:

```env
# Email Configuration (for forgot password feature)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your_email@gmail.com
SMTP_PASSWORD=your_app_password
FROM_EMAIL=noreply@yourapp.com

# Frontend URL (for reset password links)
FRONTEND_URL=http://localhost:3000
```

### Gmail Setup (Recommended)

1. **Enable 2-Factor Authentication** di akun Gmail Anda
2. **Generate App Password**:
   - Buka Google Account settings
   - Security → 2-Step Verification → App passwords
   - Generate password untuk aplikasi
   - Gunakan app password ini sebagai `SMTP_PASSWORD`

## API Endpoints

### 1. Forgot Password

**Request:**
```bash
POST /api/v1/auth/forgot-password
Content-Type: application/json

{
  "email": "user@example.com"
}
```

**Response:**
```json
{
  "message": "Reset password link has been sent to your email"
}
```

### 2. Reset Password

**Request:**
```bash
POST /api/v1/auth/reset-password
Content-Type: application/json

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "new_password": "newpassword123"
}
```

**Response:**
```json
{
  "message": "Password has been reset successfully"
}
```

## Testing

### 1. Test Forgot Password

```bash
curl -X POST http://localhost:8000/api/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com"}'
```

### 2. Test Reset Password

```bash
curl -X POST http://localhost:8000/api/v1/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "token": "YOUR_RESET_TOKEN_HERE",
    "new_password": "newpassword123"
  }'
```

## Security Features

### 1. Token Security
- JWT token dengan expiry 1 jam
- Token type validation ("reset_password")
- User ID dan email validation

### 2. Email Security
- Tidak mengungkap apakah email terdaftar atau tidak
- HTML email template dengan styling
- Confirmation email setelah reset berhasil

### 3. Password Security
- Minimum 6 karakter
- Password di-hash menggunakan bcrypt
- Password lama otomatis invalid setelah reset

## Frontend Integration

### 1. Reset Password Page

Frontend harus menyediakan halaman `/reset-password` yang:
- Menerima parameter `token` dari URL query
- Menampilkan form input password baru
- Mengirim request ke API reset password

### 2. Example Frontend Flow

```javascript
// Extract token from URL
const urlParams = new URLSearchParams(window.location.search);
const token = urlParams.get('token');

// Reset password function
async function resetPassword(newPassword) {
  const response = await fetch('/api/v1/auth/reset-password', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      token: token,
      new_password: newPassword
    })
  });
  
  const result = await response.json();
  return result;
}
```

## Error Handling

### Common Error Responses

1. **Invalid Email Format**
```json
{
  "error": "Invalid email format"
}
```

2. **Invalid/Expired Token**
```json
{
  "error": "Invalid or expired reset token"
}
```

3. **Password Too Short**
```json
{
  "error": "Password must be at least 6 characters long"
}
```

4. **Email Send Failure**
```json
{
  "error": "Failed to send reset email"
}
```

## Production Considerations

### 1. Email Service
- Gunakan service email profesional (SendGrid, AWS SES, dll.)
- Setup SPF, DKIM, dan DMARC records
- Monitor email delivery rates

### 2. Rate Limiting
- Implementasi rate limiting untuk forgot password endpoint
- Batasi jumlah request per IP/user per waktu tertentu

### 3. Logging
- Log semua aktivitas reset password
- Monitor untuk aktivitas mencurigakan

### 4. Frontend Security
- Validasi token di frontend sebelum menampilkan form
- Implementasi CSRF protection
- Secure cookie handling

## Troubleshooting

### 1. Email Tidak Terkirim
- Periksa konfigurasi SMTP
- Pastikan app password Gmail benar
- Cek firewall/network restrictions

### 2. Token Invalid
- Periksa JWT secret consistency
- Pastikan token belum expired
- Validasi format token

### 3. Database Errors
- Pastikan koneksi database stabil
- Cek migrasi database
- Monitor database logs

## Next Steps

1. **Implementasi Rate Limiting**
2. **Tambah Email Templates yang Lebih Menarik**
3. **Integrasi dengan Service Email Profesional**
4. **Tambah Audit Logging**
5. **Implementasi Account Lockout Protection**

