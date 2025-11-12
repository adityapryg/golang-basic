# Fungsi dan Struct di Go

Go adalah bahasa pemrograman yang mendukung pemrograman fungsional dengan fitur struct yang powerful. Dokumen ini menjelaskan konsep-konsep penting terkait fungsi dan structs.

## 📚 Daftar Isi
1. [Fungsi Sederhana](#fungsi-sederhana)
2. [Beberapa Nilai Kembali](#beberapa-nilai-kembali)
3. [Nilai yang Diberi Nama](#nilai-yang-diberi-nama)
4. [Definisi Struct](#definisi-struct)
5. [Receiver Value vs Pointer](#receiver-value-vs-pointer) ⭐
6. [Pola Konstruktor](#pola-konstruktor)
7. [Fungsi Variadic](#fungsi-variadic) ⭐
8. [Praktik Terbaik](#praktik-terbaik)

---

## Fungsi Sederhana

Fungsi sederhana di Go didefinisikan dengan kata kunci `func`:

```go
func HelloWorld() {
    fmt.Println("Hello, World!")
}
```

**Karakteristik:**
- Bisa tidak menerima parameter
- Bisa tidak mengembalikan nilai
- Nama fungsi diawali huruf kapital = exported (bisa diakses dari package lain)

---

## Beberapa Nilai Kembali

Go memungkinkan fungsi mengembalikan lebih dari satu nilai:

```go
func GetValues() (int, string) {
    return 1, "Hello"
}

// Penggunaan
num, text := GetValues()
fmt.Println(num, text) // Output: 1 Hello

// Mengabaikan salah satu nilai dengan _
num, _ := GetValues()  // Hanya ambil angka
```

**Use Case:**
- Mengembalikan nilai dan error
- Mengembalikan multiple hasil perhitungan
- Parsing data dengan status

**Contoh Real-World:**
```go
func Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

result, err := Divide(10, 2)
if err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println("Result:", result) // Output: 5
```

---

## Nilai yang Diberi Nama

Kita bisa memberi nama pada nilai kembali untuk lebih mudah dipahami:

```go
func GetDetails() (id int, name string) {
    id = 1
    name = "Aditya"
    return // naked return
}

// Atau dengan nilai langsung
func GetUser() (id int, name string, active bool) {
    return 1, "Aditya", true
}
```

**Keuntungan:**
- Dokumentasi lebih jelas
- Bisa menggunakan naked return
- Nilai sudah ter-inisialisasi dengan zero value

**Catatan:** Naked return sebaiknya hanya untuk fungsi pendek!

---

## Definisi Struct

Struct digunakan untuk mengelompokkan data:

```go
type Person struct {
    Name string
    Age  int
}

// Membuat instance - cara 1
p1 := Person{Name: "Budi", Age: 25}

// Cara 2 - positional
p2 := Person{"Ani", 30}

// Cara 3 - partial initialization
p3 := Person{Name: "Citra"} // Age akan 0

// Cara 4 - zero value
var p4 Person // Name: "", Age: 0
```

**Anonymous Struct:**
```go
// Berguna untuk data temporary
person := struct {
    Name string
    Age  int
}{
    Name: "John",
    Age:  25,
}
```

---

## Receiver Value vs Pointer ⭐

### 🤔 Apa itu Receiver?

Receiver adalah parameter khusus yang membuat fungsi menjadi **method** dari sebuah struct. Method adalah fungsi yang "attached" ke tipe data tertentu.

```go
// Ini adalah FUNCTION (tidak ada receiver)
func Greet(p Person) {
    fmt.Println("Hello", p.Name)
}

// Ini adalah METHOD (ada receiver)
func (p Person) Greet() {
    fmt.Println("Hello", p.Name)
}

// Pemanggilan
person := Person{Name: "Budi"}
Greet(person)    // Function call
person.Greet()   // Method call (lebih natural!)
```

---

### 📊 Perbedaan Value Receiver vs Pointer Receiver

#### 1️⃣ **Value Receiver** (Receiver by Value)

```go
type Person struct {
    Name string
    Age  int
}

// Value receiver - menerima COPY dari struct
func (p Person) Greet() {
    fmt.Printf("Hello, my name is %s and I'm %d years old\n", p.Name, p.Age)
}

// Value receiver - mencoba mengubah age (TIDAK AKAN BERPENGARUH!)
func (p Person) HaveBirthdayWrong() {
    p.Age++ // Ini hanya mengubah COPY, bukan original
    fmt.Println("Inside method:", p.Age)
}
```

**Karakteristik Value Receiver:**
- ✅ Menerima **copy** dari struct
- ✅ Tidak bisa mengubah nilai original
- ✅ Aman dari race condition (concurrent safe)
- ✅ Cocok untuk struct kecil (< 16 bytes)
- ✅ Digunakan untuk method yang hanya membaca data
- ✅ Memory: Copy seluruh struct setiap kali dipanggil

**Contoh Penggunaan:**
```go
p := Person{Name: "Budi", Age: 25}
p.Greet()                  // Output: Hello, my name is Budi and I'm 25 years old

p.HaveBirthdayWrong()     // Output: Inside method: 26
fmt.Println(p.Age)         // Output: 25 (TIDAK BERUBAH!)
```

**Visual Diagram:**
```
Original Person          Copy of Person
┌──────────────┐        ┌──────────────┐
│ Name: "Budi" │   →    │ Name: "Budi" │
│ Age:  25     │ copy   │ Age:  25     │
└──────────────┘        └──────────────┘
                               ↓
                        Method mengubah copy
                               ↓
                        ┌──────────────┐
                        │ Age:  26     │
                        └──────────────┘
                        (original tetap 25)
```

---

#### 2️⃣ **Pointer Receiver** (Receiver by Pointer)

```go
// Pointer receiver - menerima POINTER ke struct
func (p *Person) HaveBirthday() {
    p.Age++ // Ini mengubah nilai ORIGINAL
    fmt.Println("Happy birthday! Now", p.Age, "years old")
}

// Pointer receiver - method yang mengubah data
func (p *Person) ChangeName(newName string) {
    p.Name = newName
}

// Pointer receiver - method yang mengubah multiple fields
func (p *Person) UpdateProfile(name string, age int) {
    p.Name = name
    p.Age = age
}
```

**Karakteristik Pointer Receiver:**
- ✅ Menerima **pointer** ke struct (address memory)
- ✅ Bisa mengubah nilai original
- ✅ Efisien untuk struct besar (tidak copy seluruh data)
- ✅ Digunakan untuk method yang memodifikasi data
- ✅ Memory: Hanya copy pointer (8 bytes di 64-bit system)
- ⚠️ Harus hati-hati dengan concurrency (bisa race condition)

**Contoh Penggunaan:**
```go
p := Person{Name: "Budi", Age: 25}

p.HaveBirthday()           // Go otomatis convert ke (&p).HaveBirthday()
// Output: Happy birthday! Now 26 years old
fmt.Println(p.Age)         // Output: 26 (BERUBAH!)

p.ChangeName("Andi")
fmt.Println(p.Name)        // Output: Andi (BERUBAH!)

p.UpdateProfile("Citra", 30)
fmt.Printf("%+v\n", p)     // Output: {Name:Citra Age:30}
```

**Visual Diagram:**
```
Original Person          Pointer to Person
┌──────────────┐        ┌──────────────┐
│ Name: "Budi" │ ←──────│  *Person     │
│ Age:  25     │  points└──────────────┘
└──────────────┘               ↓
      ↓                Method mengubah
Method langsung              original
mengubah ini                   ↓
      ↓                 ┌──────────────┐
┌──────────────┐        │ Name: "Budi" │
│ Name: "Budi" │        │ Age:  26     │
│ Age:  26     │        └──────────────┘
└──────────────┘        (original berubah)
```

---

### 📋 Tabel Perbandingan Lengkap

| Aspek | Value Receiver | Pointer Receiver |
|-------|----------------|------------------|
| **Sintaks** | `func (p Person)` | `func (p *Person)` |
| **Menerima** | Copy dari struct | Pointer ke struct |
| **Mengubah data** | ❌ Tidak bisa | ✅ Bisa |
| **Memory Usage** | Copy seluruh struct | Hanya copy pointer (8 bytes) |
| **Performa (struct kecil)** | ⚡ Cepat | ⚡ Cepat |
| **Performa (struct besar)** | 🐌 Lambat | ⚡ Cepat |
| **Thread Safety** | ✅ Thread-safe | ⚠️ Perlu synchronization |
| **Use Case** | Read-only operations | Modify operations |
| **Contoh Method** | `GetName()`, `String()` | `SetName()`, `Update()` |

---

### 🎯 Kapan Menggunakan Masing-Masing?

#### Gunakan **Value Receiver** ketika:

1. ✅ **Struct berukuran kecil** (< 16 bytes atau 2-3 field primitif)
2. ✅ **Method hanya membaca data** (getter, formatter)
3. ✅ **Tidak perlu mengubah nilai struct**
4. ✅ **Struct adalah immutable** (tidak boleh berubah)
5. ✅ **Method bersifat "pure function"** (no side effects)

**Contoh Struct Kecil:**
```go
type Point struct {
    X, Y int // Hanya 16 bytes total
}

// Value receiver - struct kecil, hanya baca
func (p Point) Distance() float64 {
    return math.Sqrt(float64(p.X*p.X + p.Y*p.Y))
}

// Value receiver - hanya membaca
func (p Point) String() string {
    return fmt.Sprintf("(%d, %d)", p.X, p.Y)
}

// Value receiver - immutable operation
func (p Point) Add(other Point) Point {
    return Point{X: p.X + other.X, Y: p.Y + other.Y}
}
```

**Penggunaan:**
```go
p1 := Point{X: 3, Y: 4}
dist := p1.Distance()        // 5.0
fmt.Println(p1.String())     // (3, 4)

p2 := Point{X: 1, Y: 2}
p3 := p1.Add(p2)             // p1 dan p2 tidak berubah
fmt.Println(p3)              // (4, 6)
```

---

#### Gunakan **Pointer Receiver** ketika:

1. ✅ **Perlu mengubah nilai struct**
2. ✅ **Struct berukuran besar** (banyak field atau ada slice/map)
3. ✅ **Untuk konsistensi** (jika ada 1 method pointer, sebaiknya semua pointer)
4. ✅ **Struct mengandung sync.Mutex** atau field yang tidak boleh di-copy
5. ✅ **Method memodifikasi state**

**Contoh Struct Besar:**
```go
type User struct {
    ID        int
    Username  string
    Email     string
    Password  string
    Profile   UserProfile  // nested struct
    Settings  UserSettings // nested struct
    CreatedAt time.Time
    UpdatedAt time.Time
} // Struct besar, > 100 bytes

// Pointer receiver - struct besar
func (u *User) UpdateEmail(newEmail string) {
    u.Email = newEmail
    u.UpdatedAt = time.Now()
}

// Pointer receiver - modifikasi
func (u *User) SetPassword(password string) error {
    if len(password) < 8 {
        return errors.New("password too short")
    }
    u.Password = hashPassword(password)
    return nil
}
```

**Contoh dengan Mutex (HARUS Pointer):**
```go
type BankAccount struct {
    Balance float64
    mu      sync.Mutex // TIDAK BOLEH DI-COPY!
}

// HARUS pointer receiver karena ada mutex
func (b *BankAccount) Deposit(amount float64) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.Balance += amount
}

// HARUS pointer receiver
func (b *BankAccount) Withdraw(amount float64) error {
    b.mu.Lock()
    defer b.mu.Unlock()
    
    if b.Balance < amount {
        return errors.New("insufficient balance")
    }
    b.Balance -= amount
    return nil
}
```

**Penggunaan:**
```go
account := &BankAccount{Balance: 1000}
account.Deposit(500)
fmt.Println(account.Balance) // 1500

err := account.Withdraw(200)
if err != nil {
    fmt.Println("Error:", err)
}
fmt.Println(account.Balance) // 1300
```

---

### ⚠️ Catatan Penting & Best Practices

#### 1. **Go Akan Otomatis Convert**

Go sangat fleksibel dalam calling method:

```go
p := Person{Name: "Budi", Age: 25}

// Ini semua valid!
p.HaveBirthday()      // value.PointerMethod() → Go convert ke (&p).HaveBirthday()

ptr := &Person{Name: "Ani", Age: 30}
ptr.Greet()           // pointer.ValueMethod() → Go convert ke (*ptr).Greet()
```

**Tapi hati-hati dengan addressability:**
```go
// ❌ Ini akan ERROR!
Person{Name: "Budi", Age: 25}.HaveBirthday()
// Error: cannot take address of Person{...}

// ✅ Solusinya:
p := Person{Name: "Budi", Age: 25}
p.HaveBirthday() // OK!
```

---

#### 2. **Konsistensi Itu Penting**

**BAIK - Semua pointer receiver:**
```go
type User struct {
    Name string
    Age  int
}

func (u *User) GetName() string { return u.Name }
func (u *User) SetName(name string) { u.Name = name }
func (u *User) GetAge() int { return u.Age }
func (u *User) SetAge(age int) { u.Age = age }
```

**KURANG BAIK - Mixed (inconsistent):**
```go
func (u User) GetName() string { return u.Name }  // value
func (u *User) SetName(name string) { u.Name = name } // pointer
func (u User) GetAge() int { return u.Age }  // value
func (u *User) SetAge(age int) { u.Age = age } // pointer
```

**Best Practice:** Pilih satu style dan konsisten!

---

#### 3. **Interface Satisfaction**

Implementasi interface berbeda untuk value dan pointer:

```go
type Stringer interface {
    String() string
}

type Person struct {
    Name string
}

// Value receiver
func (p Person) String() string {
    return p.Name
}

// ✅ Kedua ini satisfy interface Stringer
var s1 Stringer = Person{Name: "Budi"}
var s2 Stringer = &Person{Name: "Ani"}

// Tapi kalau pointer receiver:
func (p *Person) String() string {
    return p.Name
}

// ❌ Ini ERROR!
var s1 Stringer = Person{Name: "Budi"} // ERROR: Person does not implement Stringer

// ✅ Ini OK
var s2 Stringer = &Person{Name: "Ani"} // OK
```

**Rule:** Pointer receiver method hanya bisa dipanggil via pointer!

---

## Pola Konstruktor

Kita dapat menggunakan fungsi untuk membuat struct baru (konstruktor):

```go
// Basic constructor
func NewPerson(name string, age int) *Person {
    return &Person{
        Name: name,
        Age:  age,
    }
}

// Constructor dengan validasi
func NewPersonValidated(name string, age int) (*Person, error) {
    if name == "" {
        return nil, errors.New("name cannot be empty")
    }
    if age < 0 || age > 150 {
        return nil, errors.New("invalid age")
    }
    return &Person{
        Name: name,
        Age:  age,
    }, nil
}

// Constructor dengan default values
func NewPersonWithDefaults(name string) *Person {
    return &Person{
        Name: name,
        Age:  18, // default age
    }
}
```

**Penggunaan:**
```go
p1 := NewPerson("Budi", 25)

p2, err := NewPersonValidated("", 25)
if err != nil {
    fmt.Println("Error:", err) // Error: name cannot be empty
}

p3 := NewPersonWithDefaults("Citra")
fmt.Println(p3.Age) // 18
```

**Keuntungan Constructor Pattern:**
- ✅ Validasi saat pembuatan object
- ✅ Inisialisasi field kompleks
- ✅ Encapsulation (bisa private struct fields)
- ✅ Consistent object creation

---

## Fungsi Variadic ⭐

### 🤔 Apa itu Fungsi Variadic?

Fungsi variadic adalah fungsi yang dapat menerima **jumlah argument yang variable (flexible)** dengan tipe yang sama. Di dalam fungsi, parameter variadic akan menjadi **slice**.

### 📝 Sintaks

```go
func NamaFungsi(param ...TipeData) ReturnType {
    // param adalah slice: []TipeData
    // bisa loop, akses index, dll
}
```

---

### 💡 Contoh Sederhana

```go
// Fungsi variadic untuk menjumlahkan angka
func Sum(numbers ...int) int {
    total := 0
    for _, number := range numbers {
        total += number
    }
    return total
}

// Penggunaan dengan berbagai jumlah argument
result1 := Sum(1, 2, 3)              // Output: 6
result2 := Sum(1, 2, 3, 4, 5)        // Output: 15
result3 := Sum(10)                    // Output: 10
result4 := Sum()                      // Output: 0 (tidak ada argument)

// Dengan slice
numbers := []int{1, 2, 3, 4, 5}
result5 := Sum(numbers...)            // Output: 15 (expand slice)
```

---

### 🎯 Mengapa Menggunakan Fungsi Variadic?

#### 1️⃣ **Fleksibilitas - Tidak Perlu Overloading**

Tanpa variadic (buruk):
```go
func Sum2(a, b int) int { 
    return a + b 
}

func Sum3(a, b, c int) int { 
    return a + b + c 
}

func Sum4(a, b, c, d int) int { 
    return a + b + c + d 
}

// Terus bagaimana kalau butuh 5, 6, 7 parameter? 😵
```

Dengan variadic (baik):
```go
func Sum(numbers ...int) int {
    total := 0
    for _, n := range numbers {
        total += n
    }
    return total
}

// Satu fungsi untuk semua case! ✨
```

---

#### 2️⃣ **API yang Lebih Bersih**

Contoh dari standard library Go:

```go
// fmt.Println adalah variadic
fmt.Println("Hello", "World", 123, true)

// fmt.Printf juga variadic
fmt.Printf("Name: %s, Age: %d, Active: %v\n", name, age, active)

// strings.Join
result := strings.Join([]string{"a", "b", "c"}, ",") // perlu slice

// Dengan variadic custom:
func Join(separator string, items ...string) string {
    return strings.Join(items, separator)
}

result := Join(",", "a", "b", "c") // lebih natural!
```

---

#### 3️⃣ **Mudah Diperluas**

Code lebih maintainable dan tidak perlu refactor saat requirement berubah:

```go
// Versi 1: Print 1 message
func Log(message string) {
    fmt.Println("[LOG]", message)
}

// Versi 2: Print banyak messages - cukup ubah signature!
func Log(messages ...string) {
    for _, msg := range messages {
        fmt.Println("[LOG]", msg)
    }
}

// Existing code tetap jalan:
Log("Server started")

// New code bisa multiple:
Log("Server started", "Port: 8080", "Environment: production")
```

---

### 📚 Contoh-Contoh Penggunaan Real-World

#### Contoh 1: Menggabungkan String dengan Separator

```go
func JoinStrings(separator string, texts ...string) string {
    if len(texts) == 0 {
        return ""
    }
    return strings.Join(texts, separator)
}

// Penggunaan
result1 := JoinStrings(", ", "Go", "Python", "Java", "JavaScript")
fmt.Println(result1) // Output: Go, Python, Java, JavaScript

result2 := JoinStrings(" | ", "Apple", "Banana")
fmt.Println(result2) // Output: Apple | Banana

result3 := JoinStrings(", ")
fmt.Println(result3) // Output: (empty string)
```

---

#### Contoh 2: Mencari Nilai Maksimum

```go
func Max(numbers ...int) (int, error) {
    if len(numbers) == 0 {
        return 0, errors.New("no numbers provided")
    }
    
    max := numbers[0]
    for _, n := range numbers {
        if n > max {
            max = n
        }
    }
    return max, nil
}

// Penggunaan
maxNum, err := Max(3, 7, 2, 9, 1, 5)
if err != nil {
    fmt.Println("Error:", err)
} else {
    fmt.Println("Max:", maxNum) // Output: Max: 9
}

// Error case
_, err = Max()
fmt.Println(err) // Output: no numbers provided
```

---

#### Contoh 3: Logger dengan Multiple Messages dan Timestamp

```go
func LogInfo(messages ...string) {
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    for _, msg := range messages {
        fmt.Printf("[%s] INFO: %s\n", timestamp, msg)
    }
}

func LogError(messages ...string) {
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    for _, msg := range messages {
        fmt.Printf("[%s] ERROR: %s\n", timestamp, msg)
    }
}

// Penggunaan
LogInfo("Server started", "Port: 8080", "Environment: production")
// Output:
// [2025-01-12 01:40:55] INFO: Server started
// [2025-01-12 01:40:55] INFO: Port: 8080
// [2025-01-12 01:40:55] INFO: Environment: production

LogError("Database connection failed", "Retrying in 5 seconds")
```

---

#### Contoh 4: Validasi Multiple Conditions

```go
type ValidationFunc func() error

func ValidateAll(validators ...ValidationFunc) error {
    for i, validator := range validators {
        if err := validator(); err != nil {
            return fmt.Errorf("validation %d failed: %w", i+1, err)
        }
    }
    return nil
}

// Penggunaan
func validateName() error {
    if name == "" {
        return errors.New("name is required")
    }
    return nil
}

func validateAge() error {
    if age < 0 {
        return errors.New("age must be positive")
    }
    return nil
}

func validateEmail() error {
    if !strings.Contains(email, "@") {
        return errors.New("invalid email")
    }
    return nil
}

// Run all validations
err := ValidateAll(validateName, validateAge, validateEmail)
if err != nil {
    fmt.Println("Validation failed:", err)
}
```

---

#### Contoh 5: Calculate Average

```go
func Average(numbers ...float64) (float64, error) {
    if len(numbers) == 0 {
        return 0, errors.New("cannot calculate average of empty set")
    }
    
    sum := 0.0
    for _, n := range numbers {
        sum += n
    }
    return sum / float64(len(numbers)), nil
}

// Penggunaan
avg, _ := Average(10.5, 20.3, 15.7, 25.9)
fmt.Printf("Average: %.2f\n", avg) // Output: Average: 18.10
```

---

### 🔧 Passing Slice ke Fungsi Variadic

Jika Anda sudah memiliki slice, gunakan `...` (spread operator) untuk expand:

```go
numbers := []int{1, 2, 3, 4, 5}

// ✅ BENAR - expand slice
result := Sum(numbers...)

// ❌ SALAH - ini akan error type mismatch
// result := Sum(numbers)  // Error: cannot use []int as int

// Penjelasan:
// numbers... → expand menjadi: Sum(1, 2, 3, 4, 5)
// numbers    → tetap slice:     Sum([]int{1,2,3,4,5})
```

**Contoh Lebih Lengkap:**
```go
func PrintAll(items ...string) {
    for i, item := range items {
        fmt.Printf("%d: %s\n", i+1, item)
    }
}

// Individual arguments
PrintAll("apple", "banana", "cherry")

// From slice
fruits := []string{"apple", "banana", "cherry"}
PrintAll(fruits...) // expand the slice

// Mixed? NO!
// PrintAll("apple", fruits...)  // ❌ ERROR: cannot mix
```

---

### ⚙️ Kombinasi Parameter Regular dan Variadic

**Aturan Penting:** Parameter variadic **HARUS di posisi terakhir**

```go
// ✅ BENAR - variadic di akhir
func PrintWithPrefix(prefix string, messages ...string) {
    for _, msg := range messages {
        fmt.Println(prefix + msg)
    }
}

// ✅ BENAR - multiple regular params, variadic terakhir
func LogWithLevel(level string, category string, messages ...string) {
    for _, msg := range messages {
        fmt.Printf("[%s] [%s] %s\n", level, category, msg)
    }
}

// ❌ SALAH - variadic tidak di akhir
func InvalidFunc(messages ...string, suffix string) {} 
// Compile Error: syntax error: cannot use ... with non-final parameter

// ❌ SALAH - multiple variadic
func AlsoInvalid(nums1 ...int, nums2 ...int) {}
// Compile Error: can only use ... with final parameter
```

**Contoh Penggunaan yang Benar:**
```go
PrintWithPrefix("[INFO] ", "Server started", "Port 8080", "Ready to accept connections")
// Output:
// [INFO] Server started
// [INFO] Port 8080
// [INFO] Ready to accept connections

LogWithLevel("ERROR", "Database", "Connection timeout", "Retrying...")
// Output:
// [ERROR] [Database] Connection timeout
// [ERROR] [Database] Retrying...
```

---

### 🎯 Kapan Menggunakan Fungsi Variadic?

#### ✅ Gunakan Variadic Ketika:

1. **Jumlah parameter tidak pasti/bervariasi**
   ```go
   func Max(numbers ...int) int
   ```

2. **Semua parameter memiliki tipe yang sama**
   ```go
   func Concat(strings ...string) string
   ```

3. **Membuat utility functions** (Sum, Max, Min, Average)
   ```go
   func Sum(nums ...int) int
   func Min(nums ...float64) float64
   ```

4. **Logging atau printing dengan multiple arguments**
   ```go
   func Log(messages ...string)
   func Print(items ...interface{})
   ```

5. **Builder pattern atau chainable functions**
   ```go
   func CreateUser(attrs ...UserAttribute) *User
   ```

6. **Validation dengan multiple checkers**
   ```go
   func Validate(checks ...ValidationFunc) error
   ```

---

#### ❌ Jangan Gunakan Variadic Ketika:

1. **Parameter memiliki tipe berbeda**
   ```go
   // ❌ JANGAN ini
   func CreateUser(args ...interface{}) // name, age, email?
   
   // ✅ LAKUKAN ini - gunakan struct
   type CreateUserParams struct {
       Name  string
       Age   int
       Email string
   }
   func CreateUser(params CreateUserParams) *User
   ```

2. **Parameter memiliki makna berbeda**
   ```go
   // ❌ JANGAN ini - bingung parameter mana yang mana
   func ProcessData(data ...string) // filename? encoding? mode?
   
   // ✅ LAKUKAN ini - named parameters
   func ProcessData(filename, encoding, mode string)
   ```

3. **Hanya ada 2-3 parameter fixed**
   ```go
   // ❌ JANGAN ini - overkill
   func Add(numbers ...int) int // kalau cuma mau 2 angka
   
   // ✅ LAKUKAN ini
   func Add(a, b int) int
   ```

4. **Perlu validasi parameter spesifik**
   ```go
   // ❌ JANGAN ini - sulit validasi inside
   func Connect(params ...string) // host, port, username, password?
   
   // ✅ LAKUKAN ini
   func Connect(host string, port int, username, password string)
   ```

---

### 📊 Contoh Real-World: Database Query Builder

```go
type QueryBuilder struct {
    table      string
    conditions []string
}

func NewQuery(table string) *QueryBuilder {
    return &QueryBuilder{table: table}
}

// Variadic untuk multiple conditions
func (q *QueryBuilder) Where(conditions ...string) *QueryBuilder {
    q.conditions = append(q.conditions, conditions...)
    return q
}

func (q *QueryBuilder) Build() string {
    query := fmt.Sprintf("SELECT * FROM %s", q.table)
    if len(q.conditions) > 0 {
        query += " WHERE " + strings.Join(q.conditions, " AND ")
    }
    return query
}

// Penggunaan
sql := NewQuery("users").
    Where("age > 18", "active = true", "country = 'ID'").
    Build()
    
fmt.Println(sql)
// Output: SELECT * FROM users WHERE age > 18 AND active = true AND country = 'ID'

// Atau tanpa condition
sql2 := NewQuery("products").Build()
fmt.Println(sql2)
// Output: SELECT * FROM products
```

---

### 🎨 Contoh Real-World: Event Handler System

```go
type Event struct {
    Name string
    Data interface{}
}

type EventHandler func(Event)

type EventEmitter struct {
    handlers map[string][]EventHandler
}

func NewEventEmitter() *EventEmitter {
    return &EventEmitter{
        handlers: make(map[string][]EventHandler),
    }
}

// Variadic untuk register multiple handlers sekaligus
func (e *EventEmitter) On(eventName string, handlers ...EventHandler) {
    e.handlers[eventName] = append(e.handlers[eventName], handlers...)
}

func (e *EventEmitter) Emit(event Event) {
    if handlers, ok := e.handlers[event.Name]; ok {
        for _, handler := range handlers {
            handler(event)
        }
    }
}

// Penggunaan
emitter := NewEventEmitter()

// Register multiple handlers sekaligus
emitter.On("user_registered",
    func(e Event) { fmt.Println("Log:", e.Name) },
    func(e Event) { fmt.Println("Send email to:", e.Data) },
    func(e Event) { fmt.Println("Update analytics") },
)

// Emit event
emitter.Emit(Event{
    Name: "user_registered",
    Data: "user@example.com",
})
// Output:
// Log: user_registered
// Send email to: user@example.com
// Update analytics
```

---

### 💡 Tips & Tricks dengan Variadic

#### 1. Empty Variadic Check
```go
func ProcessItems(items ...string) {
    if len(items) == 0 {
        fmt.Println("No items to process")
        return
    }
    // process items...
}
```

#### 2. Variadic sebagai Optional Parameters
```go
func Connect(host string, options ...ConnectOption) *Connection {
    conn := &Connection{Host: host}
    
    // Apply options if provided
    for _, opt := range options {
        opt(conn)
    }
    
    return conn
}

type ConnectOption func(*Connection)

func WithTimeout(timeout time.Duration) ConnectOption {
    return func(c *Connection) {
        c.Timeout = timeout
    }
}

func WithRetry(maxRetries int) ConnectOption {
    return func(c *Connection) {
        c.MaxRetries = maxRetries
    }
}

// Penggunaan
conn := Connect("localhost:8080")
conn := Connect("localhost:8080", WithTimeout(30*time.Second))
conn := Connect("localhost:8080", WithTimeout(30*time.Second), WithRetry(3))
```

#### 3. Combine Multiple Slices
```go
func CombineAll(slices ...[]int) []int {
    var result []int
    for _, slice := range slices {
        result = append(result, slice...)
    }
    return result
}

a := []int{1, 2, 3}
b := []int{4, 5, 6}
c := []int{7, 8, 9}

combined := CombineAll(a, b, c)
fmt.Println(combined) // [1 2 3 4 5 6 7 8 9]
```

---

## Praktik Terbaik

### 🎯 Functions
1. ✅ Gunakan nama fungsi yang **deskriptif** dan **action-oriented**
   ```go
   // ❌ Bad
   func Do(x int) int
   
   // ✅ Good
   func CalculateTax(income int) int
   ```

2. ✅ Jaga agar fungsi tetap **kecil** dan fokus pada **satu tugas** (Single Responsibility)
   ```go
   // ❌ Bad - fungsi melakukan terlalu banyak
   func ProcessUserAndSendEmailAndLog(user User) error {
       // validate
       // save to DB
       // send email
       // log activity
   }
   
   // ✅ Good - pisahkan responsibility
   func ValidateUser(user User) error
   func SaveUser(user User) error
   func SendWelcomeEmail(user User) error
   func LogUserActivity(user User)
   ```

3. ✅ Return error sebagai nilai **terakhir**
   ```go
   // ✅ Good
   func GetUser(id int) (*User, error)
   func SaveData(data []byte) (int, error)
   ```

4. ✅ Gunakan **named return values** untuk dokumentasi yang lebih baik
   ```go
   func Divide(a, b float64) (result float64, err error) {
       if b == 0 {
           err = errors.New("division by zero")
           return
       }
       result = a / b
       return
   }
   ```

---

### 🎯 Structs
1. ✅ Gunakan **pointer receiver** untuk method yang mengubah state
   ```go
   func (u *User) SetName(name string) {
       u.Name = name
   }
   ```

2. ✅ Gunakan **value receiver** untuk method read-only
   ```go
   func (u User) GetName() string {
       return u.Name
   }
   ```

3. ✅ **Konsisten** dalam penggunaan receiver (semua pointer atau semua value)
   ```go
   // ✅ Good - konsisten semua pointer
   func (u *User) GetName() string { return u.Name }
   func (u *User) SetName(name string) { u.Name = name }
   
   // ❌ Bad - mixed
   func (u User) GetName() string { return u.Name }
   func (u *User) SetName(name string) { u.Name = name }
   ```

4. ✅ Gunakan **constructor pattern** untuk validasi saat pembuatan struct
   ```go
   func NewUser(name, email string) (*User, error) {
       if name == "" {
           return nil, errors.New("name required")
       }
       if !isValidEmail(email) {
           return nil, errors.New("invalid email")
       }
       return &User{Name: name, Email: email}, nil
   }
   ```

---

### 🎯 Variadic Functions
1. ✅ **Dokumentasikan** dengan jelas bahwa fungsi menerima variadic args
   ```go
   // Sum calculates the sum of all provided numbers.
   // If no numbers are provided, returns 0.
   func Sum(numbers ...int) int
   ```

2. ✅ **Handle case** ketika tidak ada argument yang diberikan
   ```go
   func Max(numbers ...int) (int, error) {
       if len(numbers) == 0 {
           return 0, errors.New("no numbers provided")
       }
       // ...
   }
   ```

3. ✅ Gunakan **validation** untuk memastikan minimal argument jika diperlukan
   ```go
   func Average(numbers ...float64) (float64, error) {
       if len(numbers) == 0 {
           return 0, errors.New("need at least one number")
       }
       // ...
   }
   ```

4. ✅ Parameter variadic **selalu di posisi terakhir**
   ```go
   // ✅ Good
   func Log(level string, messages ...string)
   
   // ❌ Bad
   func Log(messages ...string, level string) // COMPILE ERROR
   ```

5. ✅ Pertimbangkan **menggunakan struct** jika parameter memiliki makna berbeda
   ```go
   // ❌ Bad - parameter meaning tidak jelas
   func CreateUser(data ...string) // name? email? age?
   
   // ✅ Good - jelas dan type-safe
   type UserData struct {
       Name  string
       Email string
       Age   int
   }
   func CreateUser(data UserData) *User
   ```

---

## 🔍 Ringkasan Cepat

### Receiver Value vs Pointer

| Aspek | Value Receiver | Pointer Receiver |
|-------|----------------|------------------|
| **Gunakan untuk** | Struct kecil (< 16 bytes) | Struct besar |
| **Operation** | Read-only | Modify data |
| **Data** | Immutable | Mutable |
| **Side effects** | No | Yes |
| **Example** | `func (p Person) GetName()` | `func (p *Person) SetName()` |
| **Memory** | Copy entire struct | Copy pointer only |
| **Thread-safe** | ✅ Yes | ⚠️ Need sync |

**Quick Decision:**
```go
// Perlu ubah data? → Pointer receiver
func (p *Person) HaveBirthday() { p.Age++ }

// Cuma baca? → Value receiver
func (p Person) GetName() string { return p.Name }

// Struct besar? → Pointer receiver (performance)
func (u *User) UpdateProfile() { /* ... */ }
```

---

### Fungsi Variadic

**Sintaks:**
```go
func Name(params ...Type) ReturnType {
    // params is a slice: []Type
}
```

**Penggunaan:**
```go
Name(a, b, c)           // Multiple individual args
Name(slice...)          // Expand slice
Name()                  // Zero args (valid)
```

**Aturan:**
- ✅ Variadic parameter **harus di akhir**
- ✅ Hanya bisa ada **1 variadic parameter**
- ✅ Inside function, params adalah **slice**
- ✅ Gunakan `...` untuk expand slice

**Quick Decision:**
```go
// Jumlah parameter tidak pasti? → Variadic
func Sum(numbers ...int) int

// Parameter berbeda makna? → Regular params atau struct
func Connect(host string, port int, timeout time.Duration)

// Parameter beda tipe? → Struct
type Config struct { Host string; Port int; }
func Connect(config Config)
```

---

## 🎓 Latihan Praktis

### Latihan 1: Value vs Pointer Receiver
Buat struct `Rectangle` dengan fields `Width` dan `Height`, kemudian implementasikan:
- `Area() float64` menggunakan **value receiver** (menghitung luas)
- `Scale(factor float64)` menggunakan **pointer receiver** (mengubah ukuran)
- `Perimeter() float64` menggunakan **value receiver** (menghitung keliling)

```go
type Rectangle struct {
    Width  float64
    Height float64
}

// Implementasikan method-method di atas!
```

**Expected Output:**
```go
rect := Rectangle{Width: 10, Height: 5}
fmt.Println(rect.Area())       // 50
fmt.Println(rect.Perimeter())  // 30

rect.Scale(2)
fmt.Println(rect.Width)        // 20
fmt.Println(rect.Height)       // 10
```

---

### Latihan 2: Fungsi Variadic - Average Calculator
Buat fungsi `Average(numbers ...float64) (float64, error)` yang:
- Menghitung rata-rata dari semua angka
- Return error jika tidak ada angka yang diberikan
- Return error jika ada angka negatif

```go
func Average(numbers ...float64) (float64, error) {
    // Implementasi di sini
}
```

**Expected Output:**
```go
avg, _ := Average(10, 20, 30)
fmt.Println(avg)  // 20.0

avg, _ = Average(5.5, 6.5, 7.5, 8.5)
fmt.Println(avg)  // 7.0

_, err := Average()
fmt.Println(err)  // Error: no numbers provided

_, err = Average(10, -5, 20)
fmt.Println(err)  // Error: negative numbers not allowed
```

---

### Latihan 3: Kombinasi - Calculator Struct
Buat struct `Calculator` dengan history dan implementasikan:
- Value receiver: `GetHistory() []float64`
- Pointer receiver: `Clear()`
- Variadic pointer receiver: `AddAll(numbers ...int) int`
- Variadic pointer receiver: `MultiplyAll(numbers ...int) int`

```go
type Calculator struct {
    History []float64
}

// Implementasikan method-method di atas!
```

**Expected Output:**
```go
calc := &Calculator{}

sum := calc.AddAll(1, 2, 3, 4, 5)
fmt.Println(sum)  // 15

product := calc.MultiplyAll(2, 3, 4)
fmt.Println(product)  // 24

history := calc.GetHistory()
fmt.Println(history)  // [15, 24]

calc.Clear()
fmt.Println(calc.History)  // []
```

---

### Latihan 4: Variadic dengan Validation
Buat fungsi `CreateTags(tags ...string) ([]string, error)` yang:
- Menerima multiple tag strings
- Validasi: tidak boleh ada tag kosong
- Validasi: tidak boleh ada tag duplikat
- Return unique tags dalam lowercase
- Return error jika ada validasi gagal

```go
func CreateTags(tags ...string) ([]string, error) {
    // Implementasi di sini
}
```

**Expected Output:**
```go
tags, _ := CreateTags("Go", "Python", "Java")
fmt.Println(tags)  // [go python java]

tags, _ = CreateTags("Go", "go", "GO")
fmt.Println(tags)  // [go]

_, err := CreateTags("Go", "", "Python")
fmt.Println(err)  // Error: empty tag not allowed

_, err = CreateTags()
fmt.Println(err)  // Error: at least one tag required
```

---

## 📚 Referensi

### Official Documentation
- [Go Tour - Methods](https://go.dev/tour/methods)
- [Effective Go - Methods](https://go.dev/doc/effective_go#methods)
- [Go by Example - Variadic Functions](https://gobyexample.com/variadic-functions)
- [Go Specification - Function Types](https://go.dev/ref/spec#Function_types)

### Best Practices Articles
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- [Practical Go: Real World Advice](https://dave.cheney.net/practical-go/presentations/qcon-china.html)

### Performance
- [Go Performance - Value vs Pointer](https://segment.com/blog/allocation-efficiency-in-high-performance-go-services/)
- [Understanding Allocation in Go](https://medium.com/a-journey-with-go/go-understand-the-design-of-sync-pool-2dde3024e277)

---

## 🎯 Kesimpulan

**Functions dan Structs** adalah fondasi penting dalam Go:

1. **Functions** memberikan modularitas dan reusability
2. **Structs** mengorganisir data dengan efektif
3. **Receiver** menghubungkan behavior (methods) dengan data (structs)
4. **Variadic** memberikan fleksibilitas dalam jumlah parameter

**Key Takeaways:**
- Gunakan **value receiver** untuk read-only, **pointer receiver** untuk modify
- **Konsistensi** receiver type dalam satu struct
- **Variadic** untuk flexibility, tapi jangan overuse
- **Constructor pattern** untuk validasi object creation
- Selalu prioritaskan **readability** dan **maintainability**

Dengan memahami konsep-konsep ini dengan baik, Anda akan bisa menulis Go code yang clean, efficient, dan idiomatic! 🚀
