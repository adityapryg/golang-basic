# 06 - Utility Functions

Create helper functions for password hashing and JWT token management.

---

## Overview

We'll create:

1. `internal/utils/password.go` - Password hashing with bcrypt
2. `internal/utils/jwt.go` - JWT token generation and validation

---

## Step 1: Create Password Utilities

### 1.1 Create password.go File

📝 **Create file:** `internal/utils/password.go`

```go
package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword membuat hash dari password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword memverifikasi password dengan hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
```

- [ ] File created at `internal/utils/password.go`
- [ ] Content copied exactly
- [ ] File saved

### 1.2 Understanding Bcrypt

💡 **What is bcrypt?**

- Industry-standard password hashing algorithm
- Slow by design (prevents brute force attacks)
- Automatically handles salt generation
- Output always different even for same password

**Example:**

```go
hash1, _ := HashPassword("password123")
hash2, _ := HashPassword("password123")

fmt.Println(hash1) // $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
fmt.Println(hash2) // $2a$10$bJbAZ5M/lZpjRJXlLYj8qeYqYWZNmXfKW4HJiMpXaH... (different!)
```

💡 **bcrypt.DefaultCost:**

- Cost = 10 (default)
- Higher cost = more secure but slower
- Each increment doubles the time
- Cost 10 = ~100ms on modern hardware

### 1.3 Password Hashing Flow

```
Registration:
1. User provides plain password: "mypassword"
2. HashPassword() called
3. Bcrypt generates salt and hash
4. Store hash in database: "$2a$10$..."
5. NEVER store plain password!

Login:
1. User provides plain password: "mypassword"
2. Get stored hash from database
3. CheckPassword(plain, hash) called
4. Bcrypt verifies password
5. Returns true if match, false if not
```

---

## Step 2: Create JWT Utilities

### 2.1 Create jwt.go File

📝 **Create file:** `internal/utils/jwt.go`

```go
package utils

import (
	"errors"
	"time"

	"github.com/adityapryg/golang-demo/20-mini-project/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

// Claims struktur JWT claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken membuat JWT token untuk user
func GenerateToken(userID uint, username string) (string, error) {
	// Token berlaku 24 jam
	expirationTime := time.Now().Add(24 * time.Hour)

	// Buat claims
	claims := &Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Buat token dengan claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token dengan secret key
	cfg := config.LoadConfig()
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken memvalidasi JWT token dan mengembalikan claims
func ValidateToken(tokenString string) (*Claims, error) {
	// Load config untuk secret key
	cfg := config.LoadConfig()

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validasi signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Ambil claims dari token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
```

- [ ] File created at `internal/utils/jwt.go`
- [ ] Content copied exactly
- [ ] File saved

### 2.2 Understanding JWT Structure

💡 **JWT Token Format:**

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ1c2VybmFtZSI6InRlc3R1c2VyIiwiZXhwIjoxNzAwMDAwMDAwfQ.signature

HEADER.PAYLOAD.SIGNATURE
```

**Header:**

```json
{
  "alg": "HS256",
  "typ": "JWT"
}
```

**Payload (Claims):**

```json
{
  "user_id": 1,
  "username": "testuser",
  "exp": 1700000000,
  "iat": 1699913600
}
```

**Signature:**

- HMAC SHA256(header + payload + secret)
- Prevents tampering
- Only server knows secret

### 2.3 JWT Claims Explained

```go
type Claims struct {
	UserID   uint   `json:"user_id"`      // Custom claim
	Username string `json:"username"`     // Custom claim
	jwt.RegisteredClaims                  // Standard claims
}
```

**Custom Claims:**

- `UserID` - Identifies the user
- `Username` - User's username (optional, for convenience)

**Standard Claims (RegisteredClaims):**

- `ExpiresAt` - When token expires
- `IssuedAt` - When token was created
- `NotBefore` - Token not valid before this time
- `Issuer` - Who issued the token
- `Subject` - Subject of the token
- `Audience` - Who token is intended for

### 2.4 Token Generation Flow

```
1. User logs in successfully
2. GenerateToken(userID, username) called
3. Create Claims with user info
4. Set expiration (24 hours from now)
5. Create JWT token with HS256 algorithm
6. Sign with JWT_SECRET
7. Return token string
8. Client stores token (localStorage/cookies)
```

### 2.5 Token Validation Flow

```
1. Client sends request with token in header:
   Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

2. ValidateToken(tokenString) called
3. Parse token structure
4. Verify signature with JWT_SECRET
5. Check expiration time
6. Extract claims (userID, username)
7. Return claims or error
```

---

## Step 3: Security Considerations

### 3.1 Password Security

✅ **Do:**

- Always hash passwords before storing
- Use bcrypt (or argon2, scrypt)
- Never log passwords
- Never expose passwords in responses
- Use HTTPS in production

❌ **Don't:**

- Store plain text passwords
- Use MD5 or SHA1 for passwords
- Send passwords in URL parameters
- Include passwords in error messages

### 3.2 JWT Security

✅ **Do:**

- Use strong, random JWT secret (32+ characters)
- Store secret in environment variables
- Use HTTPS to prevent token interception
- Set reasonable expiration times
- Validate tokens on every protected route

❌ **Don't:**

- Hardcode JWT secret in code
- Use weak secrets (e.g., "secret", "123456")
- Store sensitive data in JWT payload
- Trust client-provided tokens without validation
- Use same secret for different environments

### 3.3 Token Expiration

```go
expirationTime := time.Now().Add(24 * time.Hour)
```

**24 hours = good balance:**

- Short enough: Limits damage if stolen
- Long enough: Good user experience

**Adjust based on needs:**

```go
15 * time.Minute   // High security apps
1 * time.Hour      // Standard web apps
24 * time.Hour     // Mobile apps (current)
7 * 24 * time.Hour // Long-lived tokens (risky!)
```

---

## Step 4: Verify Utilities

### 4.1 Check File Structure

```bash
$ ls -la internal/utils/    # Linux/Mac
$ dir internal\utils\       # Windows
```

**Expected output:**

```
jwt.go
password.go
```

- [ ] Both files present

### 4.2 Verify Utilities Compile

```bash
$ go build internal/utils/password.go internal/utils/jwt.go
```

**Expected:** No errors

- [ ] Utilities compile successfully

---

## Step 5: Test Password Functions

📝 **Create temporary file:** `test_password.go`

```go
package main

import (
	"fmt"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/utils"
)

func main() {
	password := "mySecretPassword123"

	// Hash password
	fmt.Println("Original password:", password)
	hash, err := utils.HashPassword(password)
	if err != nil {
		fmt.Println("Error hashing:", err)
		return
	}
	fmt.Println("Hashed password:", hash)

	// Verify correct password
	if utils.CheckPassword(password, hash) {
		fmt.Println("✅ Password verification: SUCCESS")
	} else {
		fmt.Println("❌ Password verification: FAILED")
	}

	// Verify wrong password
	if utils.CheckPassword("wrongPassword", hash) {
		fmt.Println("❌ Wrong password accepted (BAD!)")
	} else {
		fmt.Println("✅ Wrong password rejected: SUCCESS")
	}

	// Test that same password gives different hash
	hash2, _ := utils.HashPassword(password)
	fmt.Println("\nSecond hash:", hash2)
	fmt.Println("Hashes different?", hash != hash2)
	fmt.Println("Both valid?", utils.CheckPassword(password, hash2))
}
```

### Run Test

```bash
$ go run test_password.go
```

**Expected output:**

```
Original password: mySecretPassword123
Hashed password: $2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy
✅ Password verification: SUCCESS
✅ Wrong password rejected: SUCCESS

Second hash: $2a$10$bJbAZ5M/lZpjRJXlLYj8qeYqYWZNmXfKW4HJiMpXaH...
Hashes different? true
Both valid? true
```

- [ ] Password hashed successfully
- [ ] Correct password verified
- [ ] Wrong password rejected
- [ ] Same password produces different hashes

### Clean Up

```bash
$ rm test_password.go    # Linux/Mac
$ del test_password.go   # Windows
```

---

## Step 6: Test JWT Functions

📝 **Create temporary file:** `test_jwt.go`

```go
package main

import (
	"fmt"
	"time"
	"github.com/adityapryg/golang-demo/20-mini-project/internal/utils"
)

func main() {
	userID := uint(1)
	username := "testuser"

	// Generate token
	fmt.Println("Generating token for:", username)
	token, err := utils.GenerateToken(userID, username)
	if err != nil {
		fmt.Println("Error generating token:", err)
		return
	}
	fmt.Println("Generated token:", token[:50]+"...") // Print first 50 chars

	// Validate token
	fmt.Println("\nValidating token...")
	claims, err := utils.ValidateToken(token)
	if err != nil {
		fmt.Println("Error validating token:", err)
		return
	}

	fmt.Println("✅ Token validated successfully!")
	fmt.Println("User ID from token:", claims.UserID)
	fmt.Println("Username from token:", claims.Username)
	fmt.Println("Expires at:", claims.ExpiresAt.Time)
	fmt.Println("Issued at:", claims.IssuedAt.Time)
	fmt.Println("Valid for:", claims.ExpiresAt.Time.Sub(time.Now()).Round(time.Hour))

	// Test invalid token
	fmt.Println("\nTesting invalid token...")
	_, err = utils.ValidateToken("invalid.token.here")
	if err != nil {
		fmt.Println("✅ Invalid token rejected:", err)
	} else {
		fmt.Println("❌ Invalid token accepted (BAD!)")
	}
}
```

### Run Test

```bash
$ go run test_jwt.go
```

**Expected output:**

```
Generating token for: testuser
Generated token: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2...

Validating token...
✅ Token validated successfully!
User ID from token: 1
Username from token: testuser
Expires at: 2025-11-15 14:30:00 +0700 WIB
Issued at: 2025-11-14 14:30:00 +0700 WIB
Valid for: 24h0m0s

Testing invalid token...
✅ Invalid token rejected: token is malformed: token contains an invalid number of segments
```

- [ ] Token generated successfully
- [ ] Token validated successfully
- [ ] Claims extracted correctly
- [ ] Expiration time is 24 hours
- [ ] Invalid token rejected

### Clean Up

```bash
$ rm test_jwt.go    # Linux/Mac
$ del test_jwt.go   # Windows
```

---

## Step 7: Commit Changes

### 7.1 Check Status

```bash
$ git status
```

- [ ] Utils files shown as untracked

### 7.2 Stage Utils

```bash
$ git add internal/utils/
```

- [ ] Files staged

### 7.3 Commit

```bash
$ git commit -m "Add password hashing and JWT utilities"
```

- [ ] Commit successful

---

## ✅ Completion Checklist

- [ ] `internal/utils/password.go` created
- [ ] `HashPassword()` function implemented
- [ ] `CheckPassword()` function implemented
- [ ] Bcrypt with DefaultCost (10)
- [ ] `internal/utils/jwt.go` created
- [ ] `Claims` struct with custom fields
- [ ] `GenerateToken()` function implemented
- [ ] `ValidateToken()` function implemented
- [ ] 24-hour token expiration
- [ ] HMAC SHA256 signing method
- [ ] JWT secret from config
- [ ] Password functions tested
- [ ] JWT functions tested
- [ ] Understanding of bcrypt security
- [ ] Understanding of JWT structure
- [ ] Utilities compile without errors
- [ ] Changes committed to git

---

## 🐛 Common Issues

### Issue: "crypto/bcrypt not found"

**Solution:** Run `go get golang.org/x/crypto@v0.17.0`

### Issue: "jwt package not found"

**Solution:** Run `go get github.com/golang-jwt/jwt/v5@v5.2.0`

### Issue: "token signature is invalid"

**Solution:** JWT secret must match between generation and validation

### Issue: Password check always fails

**Solution:** Make sure you're comparing plain password with hash, not hash with hash

---

## 📊 Current Progress

```
[██████████████████████] 37% Complete
```

**Completed:** Setup, dependencies, config, models, DTOs, utilities  
**Next:** Middleware layer

---

**Previous:** [05-dtos.md](05-dtos.md)  
**Next:** [07-middleware.md](07-middleware.md)
