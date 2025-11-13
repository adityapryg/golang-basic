# 18. Best Practice Testing di Go

## Tujuan Pembelajaran

Setelah mempelajari materi ini, Anda akan memahami:

1. Cara menulis unit tests di Go
2. Table-driven tests pattern
3. Test suites dengan testify
4. Test coverage dan benchmark
5. Mocking dan dependency injection
6. Best practices testing di Go

## Penjelasan

### Struktur Testing di Go

```
main.go        # Kode yang akan ditest
main_test.go   # Test cases
```

**Konvensi:**

- File test: `*_test.go`
- Function test: `Test*` (harus dimulai dengan `Test`)
- Benchmark: `Benchmark*`
- Example: `Example*`

### Table-Driven Tests

Pattern umum di Go untuk testing berbagai kasus:

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 5, 3, 8},
        {"negative", -5, -3, -8},
        {"zero", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Test Suite dengan testify

```go
type UserServiceTestSuite struct {
    suite.Suite
    service *UserService
}

func (suite *UserServiceTestSuite) SetupTest() {
    // Setup sebelum setiap test
    suite.service = NewUserService()
}

func (suite *UserServiceTestSuite) TestCreateUser() {
    user, err := suite.service.CreateUser("john", "john@example.com", 25)
    suite.NoError(err)
    suite.NotNil(user)
}
```

### Assertions dengan testify

```go
assert.Equal(t, expected, actual)
assert.NoError(t, err)
assert.Error(t, err)
assert.Nil(t, value)
assert.NotNil(t, value)
assert.True(t, condition)
assert.False(t, condition)
assert.Len(t, slice, 5)
assert.Contains(t, slice, element)
```

## Prasyarat

```bash
cd 18-testing-docs
go mod download
```

Dependencies:

- `github.com/stretchr/testify` - Testing toolkit

## Cara Menjalankan Tests

### 1. Run All Tests

```bash
go test
```

Output:

```
PASS
ok      github.com/adityapryg/golang-demo/18-testing-docs       0.123s
```

### 2. Run Tests dengan Verbose

```bash
go test -v
```

Output detail untuk setiap test case:

```
=== RUN   TestCalculator_Add
=== RUN   TestCalculator_Add/positive_numbers
=== RUN   TestCalculator_Add/negative_numbers
--- PASS: TestCalculator_Add (0.00s)
    --- PASS: TestCalculator_Add/positive_numbers (0.00s)
    --- PASS: TestCalculator_Add/negative_numbers (0.00s)
```

### 3. Run Specific Test

```bash
go test -run TestCalculator_Add
go test -run "TestCalculator_Add/positive"
```

### 4. Test Coverage

```bash
# Lihat coverage percentage
go test -cover

# Generate coverage report
go test -coverprofile=coverage.out

# Lihat coverage di browser
go tool cover -html=coverage.out
```

Output:

```
PASS
coverage: 85.7% of statements
ok      github.com/adityapryg/golang-demo/18-testing-docs       0.123s
```

### 5. Run Benchmarks

```bash
go test -bench=.
```

Output:

```
BenchmarkCalculator_Add-8               1000000000       0.3054 ns/op
BenchmarkCalculator_Multiply-8          1000000000       0.3065 ns/op
BenchmarkUserService_CreateUser-8       2000000          725 ns/op
```

### 6. Continuous Testing (Watch Mode)

Install `reflex`:

```bash
go install github.com/cespare/reflex@latest
reflex -r '\.go$' -- go test -v
```

## Test Cases yang Diimplementasikan

### Calculator Tests

- ✅ Add: positive, negative, mixed, zero
- ✅ Subtract
- ✅ Multiply
- ✅ Divide: valid division, divide by zero

### UserService Tests (Test Suite)

- ✅ CreateUser: success, empty username, empty email, invalid age, duplicate username
- ✅ GetUser: success, not found
- ✅ GetAllUsers
- ✅ UpdateUser: success, not found
- ✅ DeleteUser: success, not found

### User Methods Tests

- ✅ IsAdult: child, teenager, adult, senior
- ✅ GetEmailDomain: gmail, company, subdomain, invalid

### Benchmarks

- ✅ Calculator operations
- ✅ UserService CreateUser

## Output Test Run

```bash
$ go test -v -cover
```

```
=== RUN   TestCalculator_Add
=== RUN   TestCalculator_Add/positive_numbers
=== RUN   TestCalculator_Add/negative_numbers
=== RUN   TestCalculator_Add/mixed_numbers
=== RUN   TestCalculator_Add/zero
--- PASS: TestCalculator_Add (0.00s)
    --- PASS: TestCalculator_Add/positive_numbers (0.00s)
    --- PASS: TestCalculator_Add/negative_numbers (0.00s)
    --- PASS: TestCalculator_Add/mixed_numbers (0.00s)
    --- PASS: TestCalculator_Add/zero (0.00s)

=== RUN   TestCalculator_Divide
=== RUN   TestCalculator_Divide/valid_division
=== RUN   TestCalculator_Divide/divide_by_zero
--- PASS: TestCalculator_Divide (0.00s)
    --- PASS: TestCalculator_Divide/valid_division (0.00s)
    --- PASS: TestCalculator_Divide/divide_by_zero (0.00s)

=== RUN   TestUserServiceTestSuite
=== RUN   TestUserServiceTestSuite/TestCreateUser_Success
=== RUN   TestUserServiceTestSuite/TestCreateUser_EmptyUsername
=== RUN   TestUserServiceTestSuite/TestCreateUser_EmptyEmail
=== RUN   TestUserServiceTestSuite/TestCreateUser_InvalidAge
=== RUN   TestUserServiceTestSuite/TestCreateUser_InvalidAge/negative_age
=== RUN   TestUserServiceTestSuite/TestCreateUser_InvalidAge/too_old
=== RUN   TestUserServiceTestSuite/TestCreateUser_DuplicateUsername
=== RUN   TestUserServiceTestSuite/TestGetUser_Success
=== RUN   TestUserServiceTestSuite/TestGetUser_NotFound
=== RUN   TestUserServiceTestSuite/TestGetAllUsers
=== RUN   TestUserServiceTestSuite/TestUpdateUser_Success
=== RUN   TestUserServiceTestSuite/TestUpdateUser_NotFound
=== RUN   TestUserServiceTestSuite/TestDeleteUser_Success
=== RUN   TestUserServiceTestSuite/TestDeleteUser_NotFound
--- PASS: TestUserServiceTestSuite (0.00s)

=== RUN   TestUser_IsAdult
=== RUN   TestUser_IsAdult/child
=== RUN   TestUser_IsAdult/teenager
=== RUN   TestUser_IsAdult/adult_(18)
=== RUN   TestUser_IsAdult/adult_(25)
=== RUN   TestUser_IsAdult/senior
--- PASS: TestUser_IsAdult (0.00s)

=== RUN   TestUser_GetEmailDomain
=== RUN   TestUser_GetEmailDomain/gmail
=== RUN   TestUser_GetEmailDomain/company
=== RUN   TestUser_GetEmailDomain/subdomain
=== RUN   TestUser_GetEmailDomain/no_at_sign
--- PASS: TestUser_GetEmailDomain (0.00s)

PASS
coverage: 92.5% of statements
ok      github.com/adityapryg/golang-demo/18-testing-docs       0.145s
```

## Best Practices

### 1. Naming Conventions

```go
// ✅ Good
func TestUserService_CreateUser(t *testing.T) {}
func TestCalculator_Add(t *testing.T) {}

// ❌ Bad
func TestCreate(t *testing.T) {}
func Test1(t *testing.T) {}
```

### 2. Table-Driven Tests

```go
tests := []struct {
    name     string
    input    int
    expected int
}{
    {"case 1", 5, 10},
    {"case 2", 10, 20},
}
```

### 3. Setup & Teardown

```go
func (suite *MyTestSuite) SetupTest() {
    // Setup sebelum setiap test
}

func (suite *MyTestSuite) TearDownTest() {
    // Cleanup setelah setiap test
}
```

### 4. Test Isolation

Setiap test harus independent:

```go
// ✅ Good - create fresh service
func (suite *UserServiceTestSuite) SetupTest() {
    suite.service = NewUserService()
}

// ❌ Bad - shared state
var globalService = NewUserService()
```

### 5. Descriptive Test Names

```go
// ✅ Good
t.Run("returns error when username is empty", func(t *testing.T) {})

// ❌ Bad
t.Run("test 1", func(t *testing.T) {})
```

### 6. Coverage Target

Aim for **80-90%** code coverage, tapi jangan obsessed dengan 100%.

### 7. Test Errors

Selalu test error cases:

```go
result, err := DoSomething()
assert.Error(t, err)
assert.EqualError(t, err, "expected error message")
```

## Testing Pyramid

```
        /\
       /  \       E2E Tests (5-10%)
      /    \
     /------\     Integration Tests (20-30%)
    /        \
   /          \   Unit Tests (60-70%)
  /____________\
```

Fokus pada unit tests, kemudian integration, baru E2E.

## Mocking (Advanced)

Untuk dependency injection dan mocking, gunakan interface:

```go
type UserRepository interface {
    GetUser(id int) (*User, error)
    CreateUser(user *User) error
}

type UserService struct {
    repo UserRepository
}

// Dalam test, gunakan mock repository
type MockUserRepository struct {
    mock.Mock
}
```

## Continuous Integration

Contoh GitHub Actions:

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.22
      - run: go test -v -cover ./...
```

## Referensi

- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Effective Go - Testing](https://go.dev/doc/effective_go#testing)
- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)
