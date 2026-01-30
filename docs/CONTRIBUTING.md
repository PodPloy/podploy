# Contributing to Podploy 👨‍💻

First off, thanks for taking the time to contribute! Podploy follows a strict **Clean Architecture**, **Linux-First**, and **Monorepo** philosophy.

## 🛠️ The Developer Toolbelt

To ensure a reproducible environment, we require a specific set of tools. You can install them globally or use Nix (if available).

### Required Tools

| Tool | Min Version | Purpose | Command |
| :--- | :--- | :--- | :--- |
| **[Go](https://go.dev/)** | `1.23+` | Language Runtime | `goenv install 1.23.0` |
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
    After cloning the repo, install the Git Hooks (Lefthook). This is mandatory to pass CI checks.
    ```bash
    task setup
    ```

2.  **Generate Code:**
    Generate Ent schemas, Protobuf files, and Mocks:
    ```bash
    task generate
    ```

3.  **Start Development Server:**
    You have two options for running the environment:

    * **Option A: Parallel (Single Terminal)**
        Runs both Server and Agent in the same window. Logs will be mixed.
        ```bash
        task dev
        ```

    * **Option B: Split (Recommended)**
        Run these in separate terminals for cleaner logs:
        ```bash
        task dev:server   # Terminal 1
        task dev:agent    # Terminal 2
        ```

---

## 📏 Coding Standards & Workflow

We use **Lefthook** to enforce standards locally before you push.

### 1. Conventional Commits (Strict)
We strictly enforce [Conventional Commits](https://www.conventionalcommits.org/). Your commit message **must** match the regex:
`^(feat|fix|docs|style|refactor|perf|test|build|ci|chore|revert)(\(.+\))?!?: .+$`

* ✅ `feat(server): add graceful shutdown`
* ✅ `fix: resolve nil pointer in agent`
* ❌ `wip`
* ❌ `fixed bug`

### 2. Linting & Formatting
We use strict linting rules defined in `.golangci.yml`.
* **Formatter:** `gofumpt` (stricter than gofmt).
* **Imports:** Managed by `gci`. Your editor should group imports into: Standard Libs, 3rd Party, and `github.com/your-user/podploy`.

Run the linter manually:
```bash
task lint
```

## 🔄 Useful Commands (`Taskfile`)

| Command | Description |
| :--- | :--- |
| `task setup` | Installs Git hooks (Run once). |
| `task dev` | Starts Server & Agent with hot-reload (Air). |
| `task test` | Runs unit tests (race detector enabled). |
| `task mocks` | Regenerates mocks for interfaces in `domain/ports`. |
| `task security` | Scans dependencies for known vulnerabilities. |
| `task migrate:diff name=foo` | Generates a new SQL migration file. |

---

## 🏛️ Architecture Guide

Podploy follows the **Hexagonal Architecture** (also known as Ports & Adapters) combined with **Clean Architecture** principles.

The Golden Rule of Dependency: **Dependencies only point INWARD.**
* Inner layers (`domain`) know nothing about outer layers (`adapter`).
* Outer layers (`adapter`) import inner layers to implement interfaces.

### 📂 Directory Structure Explained

```text
podploy/
├── cmd/                # Application Entrypoints
├── internal/           # Private Application Code (The Core)
│   ├── domain/         # 1. The Enterprise Logic (Pure Go)
│   ├── usecase/        # 2. The Application Logic (Orchestration)
│   └── adapter/        # 3. The Infrastructure (Implementation)
├── pkg/                # Public Shared Code (Generated)
├── proto/              # API Contracts (gRPC)
├── web/                # Frontend (SvelteKit)
├── migrations/         # Database Schema Changes (SQL)
└── design/             # Visual Documentation
```