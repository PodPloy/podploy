# Contributing to Podploy 👨‍💻

First off, thanks for taking the time to contribute! Podploy follows a strict **Clean Architecture**, **Linux-First**, and **Monorepo** philosophy.

> **One-Line Setup (The Happy Path):**
> If you have [Nix](https://nixos.org/) installed, simply run `nix develop`. You will drop into a shell with Go, Task, Air, and all dependencies pre-configured. No installation required.

---

## 🛠️ The Developer Toolbelt

To ensure a reproducible environment, we require a specific set of tools.

### Option A: Nix (Recommended) ❄️
We provide a `flake.nix` that sets up the exact versions of all tools.
```bash
nix develop
# You are ready.
```

### Option B: Manual Installation
If you are not using Nix, please ensure you install these specific versions to avoid "It works on my machine" issues.

| Tool | Min Version | Purpose | Command |
| :--- | :--- | :--- | :--- |
| **[Go](https://go.dev/)** | `1.25+` | Language Runtime | `goenv install 1.25.0` |
| **[Task](https://taskfile.dev/)** | `v3.0+` | Task Runner (Replaces Make) | `sh -c "$(curl ...)"` |
| **[Air](https://github.com/air-verse/air)** | `Latest` | Live Reload | `go install github.com/air-verse/air@latest` |
| **[Lefthook](https://github.com/evilmartians/lefthook)** | `Latest` | Git Hooks Manager | `go install github.com/evilmartians/lefthook@latest` |
| **[Mockery](https://vektra.github.io/mockery/)** | `v2.0+` | Mock Generator | `go install github.com/vektra/mockery/v2@latest` |
| **[Atlas](https://atlasgo.io/)** | `Latest` | DB Migrations | `curl -sSf https://atlasgo.sh | sh` |
| **[GolangCI-Lint](https://golangci-lint.run/)** | `v1.55+` | Linter | *(See official docs)* |
| **[Govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck)** | `Latest` | Security Scanner | `go install golang.org/x/vuln/cmd/govulncheck@latest` |

---

## ⚡ Quick Start

1.  **Setup Environment:**
    Initialize git hooks and download dependencies.
    ```bash
    task setup
    ```

2.  **Generate Core Assets:**
    Generate Ent schemas, Protobuf/gRPC files, and Mocks. **Run this whenever you pull main.**
    ```bash
    task generate
    ```

3.  **Start Development:**
    We recommend running the Server and Agent in split terminals to distinguish logs.

    ```bash
    # Terminal 1: The Brain (API + UI)
    task dev:server

    # Terminal 2: The Muscle (Agent Runtime)
    task dev:agent

    # Terminal 3: CLI (Remote Control)
    task dev:cli
    ```

---

## 🧠 Development Workflows

### 1. Database Changes (Ent + Atlas)
We use **Ent** for ORM and **Atlas** for migrations. **Never edit SQL files manually.**

1.  Modify the schema in `internal/adapter/storage/db/schema/your_schema.go`.
2.  Generate the Go code:
    ```bash
    task generate
    ```
3.  Create a migration plan (diff):
    ```bash
    task migrate:diff name=add_user_settings
    ```
4.  Apply it locally (Optional, usually handled by dev server startup):
    ```bash
    task migrate:apply
    ```

### 2. Adding a New Feature (The "Where does it go?" Guide)
Podploy follows **Hexagonal Architecture**. Follow the dependency rule: **Dependencies only point INWARD.**

| Layer | Path | What goes here? | Can import... |
| :--- | :--- | :--- | :--- |
| **Domain** | `internal/domain` | Entities (Structs), Repository Interfaces, Custom Errors. **Pure Go.** | *Nothing* |
| **UseCase** | `internal/usecase` | Business Logic, "Service" methods. The "What" of the app. | `domain` |
| **Adapter** | `internal/adapter` | DB implementations (Ent), HTTP Handlers (Echo), gRPC Resolvers. | `usecase`, `domain` |
| **Drivers** | `cmd/` | `main.go`, Dependency Injection (Wiring), Config loading. | `adapter`, `usecase` |

---

## 📏 Standards & Pull Requests

### Conventional Commits (Strict)
We strictly enforce [Conventional Commits](https://www.conventionalcommits.org/) via `lefthook`.

* `feat`: New features (e.g., `feat(agent): add headless mode support`)
* `fix`: Bug fixes (e.g., `fix: resolve nil pointer in backup routine`)
* `chore`: Maintenance (e.g., `chore: update deps`, `chore: lint`)
* `refactor`: Code change that neither fixes a bug nor adds a feature
* `docs`: Documentation only changes

### Linting & Formatting
We use strict linting rules defined in `.golangci.yml`.
* **Formatter:** `gofumpt` (stricter than gofmt).
* **Imports:** Managed by `gci`. Your editor should group imports into: Standard Libs, 3rd Party, and `github.com/your-user/podploy`.

Run the linter manually:
```bash
task lint
```

### The Perfect PR Checklist
Before submitting a Pull Request, ensure:

1.  [ ] **Tests Pass:** Run `task test`.
2.  [ ] **No Lint Errors:** Run `task lint`.
3.  [ ] **Generated Code Updated:** Run `task generate` if you touched Protobufs or Ent schemas.
4.  [ ] **Clean History:** Squash your commits into logical units (no "wip" commits).

---

## 🔄 Useful Commands (`Taskfile`)

| Command | Description |
| :--- | :--- |
| `task setup` | Installs Git hooks and dependencies. |
| `task dev` | Starts full stack (Server + Agent + cli). |
| `task test` | Runs unit tests with race detection. |
| `task mocks` | Regenerates mocks for interfaces in `domain/ports`. |
| `task proto` | Compiles `.proto` files to Go code. |
| `task ent` | Generates Ent ORM code. |
| `task security` | Scans dependencies for known vulnerabilities (`govulncheck`). |
| `task build` | Builds static binaries for production. |
| `task migrate:diff name=foo` | Generates a new SQL migration file. |

---

## 🏛️ Architecture Guide

Podploy follows the **Hexagonal Architecture** (also known as Ports & Adapters) combined with **Clean Architecture** principles.

The Golden Rule of Dependency: **Dependencies only point INWARD.**
* Inner layers (`domain`) know nothing about outer layers (`adapter`).
* Outer layers (`adapter`) import inner layers to implement interfaces.

## 📂 Directory Structure Explained

```text
podploy/
├── cmd/                # Application Entrypoints
├── docs/               # documentation (API gRPC and Grahql, Examples, Design Project, Etc)
├── internal/           # Private Application Code (The Core)
│   ├── domain/         # 1. The Enterprise Logic (Pure Go)
│   ├── usecase/        # 2. The Application Logic (Orchestration)
│   └── adapter/        # 3. The Infrastructure (Implementation)
├── pkg/                # Public Shared Code (Generated)
├── proto/              # API Contracts (gRPC)
├── web/                # Frontend (SvelteKit)
├── migrations/         # Database Schema Changes (SQL)
├── scripts/            # utility scripts
├── templates/          # Containerfiles for Marketplace
```
## ❓ Need Help?
Check the docs/ folder for architectural diagrams or open a Discussion on GitHub.

Happy Hacking! 🚀