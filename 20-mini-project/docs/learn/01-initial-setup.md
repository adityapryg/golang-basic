# 01 - Initial Setup

Complete project initialization from scratch.

---

## Step 1: Verify Prerequisites

### 1.1 Check Go Installation

```bash
$ go version
```

**Expected output:**

```
go version go1.22.x windows/amd64
```

⚠️ **If Go is not installed:**

1. Download from https://go.dev/dl/
2. Install for your OS
3. Restart terminal
4. Verify with `go version`

### 1.2 Check Git Installation

```bash
$ git --version
```

**Expected output:**

```
git version 2.x.x
```

---

## Step 2: Create Project Directory

### 2.1 Navigate to Your Workspace

```bash
# Example: Navigate to your projects folder
$ cd d:\vue-case-study\golang-demo
```

💡 **Adjust the path to where you want to create the project**

### 2.2 Create Project Folder

```bash
$ mkdir 20-mini-project
```

- [ ] Folder `20-mini-project` created

### 2.3 Enter Project Directory

```bash
$ cd 20-mini-project
```

- [ ] You are now inside `20-mini-project`

### 2.4 Verify Current Directory

```bash
$ pwd  # On Linux/Mac
$ cd   # On Windows shows current directory
```

**Expected output should end with:**

```
.../20-mini-project
```

---

## Step 3: Initialize Go Module

### 3.1 Create Go Module

```bash
$ go mod init github.com/adityapryg/golang-demo/20-mini-project
```

**Expected output:**

```
go: creating new go.mod: module github.com/adityapryg/golang-demo/20-mini-project
```

- [ ] Command executed without errors

💡 **What this does:**

- Creates `go.mod` file
- Sets module path for imports
- Enables dependency management

### 3.2 Verify go.mod Created

```bash
$ ls go.mod      # Linux/Mac
$ dir go.mod     # Windows
```

**Expected output:**

```
go.mod
```

- [ ] File `go.mod` exists

### 3.3 View go.mod Content

```bash
$ cat go.mod     # Linux/Mac
$ type go.mod    # Windows
```

**Expected content:**

```go
module github.com/adityapryg/golang-demo/20-mini-project

go 1.22
```

- [ ] Content matches above

---

## Step 4: Create Directory Structure

### 4.1 Create Main Directories

Run these commands **one by one**:

```bash
# Create cmd directory and subdirectory
$ mkdir cmd
$ mkdir cmd\api

# Create internal directory and all subdirectories
$ mkdir internal
$ mkdir internal\config
$ mkdir internal\dto
$ mkdir internal\handler
$ mkdir internal\middleware
$ mkdir internal\model
$ mkdir internal\repository
$ mkdir internal\route
$ mkdir internal\service
$ mkdir internal\utils

# Create supporting directories
$ mkdir migrations
$ mkdir tests
$ mkdir bin
$ mkdir .github
```

- [ ] All commands executed without errors

💡 **Why this structure?**

- `cmd/api/` - Application entry point
- `internal/` - Private application code
- `migrations/` - Database migration SQL files
- `tests/` - Test files
- `bin/` - Compiled binaries
- `.github/` - GitHub-specific files

### 4.2 Verify Directory Structure

**Linux/Mac:**

```bash
$ tree -L 2
```

**Windows (PowerShell):**

```powershell
$ Get-ChildItem -Recurse -Directory -Depth 1 | Select-Object FullName
```

**Windows (Command Prompt):**

```bash
$ dir /s /b /ad
```

**Expected structure:**

```
20-mini-project/
├── .github/
├── bin/
├── cmd/
│   └── api/
├── internal/
│   ├── config/
│   ├── dto/
│   ├── handler/
│   ├── middleware/
│   ├── model/
│   ├── repository/
│   ├── route/
│   ├── service/
│   └── utils/
├── migrations/
└── tests/
```

- [ ] All directories exist
- [ ] Structure matches above

### 4.3 Create .gitkeep Files (Optional)

To track empty directories in git:

```bash
# Linux/Mac
$ touch bin/.gitkeep tests/.gitkeep

# Windows
$ type nul > bin\.gitkeep
$ type nul > tests\.gitkeep
```

- [ ] .gitkeep files created (optional)

---

## Step 5: Initialize Git Repository

### 5.1 Initialize Git

```bash
$ git init
```

**Expected output:**

```
Initialized empty Git repository in .../20-mini-project/.git/
```

- [ ] Git repository initialized

### 5.2 Create .gitignore

📝 **Create file:** `.gitignore`

```bash
# Linux/Mac
$ cat > .gitignore << 'EOF'
# Binaries
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output of the go coverage tool
*.out
coverage.out

# Dependency directories
vendor/

# Environment variables
.env
.env.local

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Logs
*.log

# Go workspace file
go.work
EOF

# Windows (PowerShell)
$ @"
# Binaries
bin/
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary
*.test

# Output of the go coverage tool
*.out
coverage.out

# Dependency directories
vendor/

# Environment variables
.env
.env.local

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Logs
*.log

# Go workspace file
go.work
"@ | Out-File -FilePath .gitignore -Encoding UTF8
```

- [ ] File `.gitignore` created

### 5.3 Verify .gitignore

```bash
$ cat .gitignore    # Linux/Mac
$ type .gitignore   # Windows
```

- [ ] Content matches above

---

## Step 6: Create .env.example

📝 **Create file:** `.env.example`

**Method 1: Using editor**
Open your code editor and create `.env.example` with:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=todolist_db

# JWT Configuration
JWT_SECRET=your-super-secret-key-change-this-in-production

# Server Configuration
SERVER_PORT=8080

# Gin Mode (debug or release)
GIN_MODE=debug
```

**Method 2: Command line**

**Linux/Mac:**

```bash
$ cat > .env.example << 'EOF'
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=todolist_db

# JWT Configuration
JWT_SECRET=your-super-secret-key-change-this-in-production

# Server Configuration
SERVER_PORT=8080

# Gin Mode (debug or release)
GIN_MODE=debug
EOF
```

**Windows (PowerShell):**

```powershell
$ @"
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=todolist_db

# JWT Configuration
JWT_SECRET=your-super-secret-key-change-this-in-production

# Server Configuration
SERVER_PORT=8080

# Gin Mode (debug or release)
GIN_MODE=debug
"@ | Out-File -FilePath .env.example -Encoding UTF8
```

- [ ] File `.env.example` created
- [ ] Content verified

---

## Step 7: Verification Checklist

### 7.1 File Structure Verification

```bash
$ ls -la  # Linux/Mac
$ dir     # Windows
```

**You should see:**

```
.git/
.github/
bin/
cmd/
internal/
migrations/
tests/
.env.example
.gitignore
go.mod
```

- [ ] All files and folders present

### 7.2 Go Module Verification

```bash
$ go mod verify
```

**Expected output:**

```
all modules verified
```

- [ ] Verification passed

### 7.3 Directory Count

**Count directories created:**

```bash
# Linux/Mac
$ find . -type d | wc -l

# Windows (PowerShell)
$ (Get-ChildItem -Recurse -Directory).Count
```

**Expected:** Around 15-16 directories (including .git and subdirectories)

- [ ] Directory count reasonable

---

## Step 8: Initial Git Commit

### 8.1 Stage Files

```bash
$ git add .
```

- [ ] Files staged

### 8.2 Check Status

```bash
$ git status
```

**Expected output should show:**

```
Changes to be committed:
  new file:   .env.example
  new file:   .gitignore
  new file:   go.mod
```

- [ ] Status shows new files

### 8.3 Initial Commit

```bash
$ git commit -m "Initial project structure"
```

**Expected output:**

```
[main (root-commit) xxxxxxx] Initial project structure
 3 files changed, X insertions(+)
 create mode 100644 .env.example
 create mode 100644 .gitignore
 create mode 100644 go.mod
```

- [ ] Commit successful

### 8.4 Verify Commit

```bash
$ git log --oneline
```

**Expected output:**

```
xxxxxxx Initial project structure
```

- [ ] Commit appears in log

---

## ✅ Completion Checklist

Before proceeding to the next section, verify:

- [ ] Go version 1.22+ installed and verified
- [ ] Project directory created: `20-mini-project`
- [ ] Go module initialized with correct path
- [ ] All 15+ directories created
- [ ] `.gitignore` file created and populated
- [ ] `.env.example` file created with all config variables
- [ ] Git repository initialized
- [ ] Initial commit made
- [ ] No error messages in any step

---

## 🐛 Common Issues

### Issue: "go: command not found"

**Solution:** Go is not installed or not in PATH. Install Go from https://go.dev/dl/

### Issue: "mkdir: cannot create directory"

**Solution:** Check you have write permissions in current directory

### Issue: "git: command not found"

**Solution:** Install Git from https://git-scm.com/downloads

### Issue: Directory structure doesn't match

**Solution:** Delete all directories and repeat Step 4 carefully

---

## 📊 Current Progress

```
[████░░░░░░░░░░░░] 6% Complete
```

**Completed:** Project initialization, directory structure, git setup  
**Next:** Installing dependencies

---

**Previous:** [00-overview.md](00-overview.md)  
**Next:** [02-dependencies.md](02-dependencies.md)
