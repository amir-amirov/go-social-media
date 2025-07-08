# Go Social Media

Social Media backend written in Go — featuring authentication, role-based authorization, post feeds, search, and real-time metrics monitoring.

---

## 🚀 Features

- **Registration** with email verification (MailTrap)
- **JWT-based Authorization** with roles: `user`, `moderator`, `admin`
- **Role-based Access Control**

  - Moderators can update others' posts/comments
  - Admins can update and delete others' posts/comments

- **Post & Comment**: create, read, update, delete
- **Follow Users** and receive a paginated post feed
- **Full-text Search** with PostgreSQL `pg_trgm` and GIN indexing
- **Rate Limiting** (fixed window)
- **Redis Caching** for user data
- **Prometheus Metrics** and **Grafana Visualization**
- **Swagger API Docs**: [View Documentation](https://gosocial.amir-amirov.kz/v1/swagger/index.html)
- **Frontend**: [React App](https://habar.amir-amirov.kz)

---

## 🛠 Tech Stack

- **Backend**: Go 1.24.1, [Chi](https://github.com/go-chi/chi)
- **Database**: PostgreSQL, `pg_trgm` and GIN for search
- **Caching**: Redis
- **Migrations**: [`golang-migrate`](https://github.com/golang-migrate/migrate)
- **Email Service**: MailTrap (SMTP)
- **Containerization**: Docker, Docker Compose
- **Monitoring**: Prometheus + Grafana
- **Deployment**: AWS EC2 (Backend), PS.kz (Frontend)

---

## 📂 Project Structure

```plaintext
go-social-media/
├── .github/workflows   # Automation files
├── bin/                # Binary files
├── cmd/                # Entry points: API, migrations
│   ├── api/            # HTTP transport layer
│   └── migrate/        # Migration files
│
├── internal/           # Core app logic (cannot import cmd)
│   ├── db/
│   ├── storage/        # Repositories and their methods
│   ├────────── cache/  # Redis caching
│   ├── mailer/         # Email service
│   ├── ratelimiter/    # Rate limiter module
│   └── utils/          # Helpers and utilities
│
├── scripts/            # Scripts
├── docs/               # Swagger docs
├── web/                # React frontend (deployed separately)
├── Makefile            # CLI helpers for migration, tests, etc.
├── docker-compose.yml  # Dev containers
└── ...

```

---

## ⚙️ Getting Started

### Prerequisites

- Go 1.24.1+
- Docker & Docker Compose
- `make` utility

### Clone & Run

```bash
git clone https://github.com/amir-amirov/go-social-media.git
cd go-social-media

# Start services
docker compose up -d

# Run DB migrations
make migrate-up

# (Optional) Seed initial data
make seed
```

---

## 🧪API Documentation

Swagger is available at:
📘 [**https://gosocial.amir-amirov.kz/v1/swagger/index.html**](https://gosocial.amir-amirov.kz/v1/swagger/index.html)

---

## 🧱 Migrations

Uses [golang-migrate](https://github.com/golang-migrate/migrate).

[EXAMPLE] Create a new migration:

```bash
make migration create_users_table
```

Apply migrations:

```bash
make migrate-up
```

Rollback:

```bash
make migrate-down
```

Migrations are located in: `cmd/migrate/migrations`

---

## 🧰 Middleware

- **Auth Middleware** – Validates JWT token
- **Role Middleware** – Allows actions based on `user`, `moderator`, `admin`
- **Rate Limiter** – Fixed-window strategy (20 requests/min)
- **Redis Cache** – Caches frequently accessed user data
- **User Middleware** – Middleware to fetch user (cached)

---

## 📊 Monitoring

- **/metrics** endpoint exposed for Prometheus
- **Grafana** dashboards included in Docker Compose
- Prometheus + Grafana run in separate containers

---

## ✅ Tests

- Located alongside application logic
- Use standard Go testing library

```bash
go test -race ./...
```

---

## 🧑‍💻 CI / CD

**GitHub Actions** enforce quality via the `audit.yaml` workflow:

- Runs on PRs to `main` or `dev`
- Steps:

  - `go vet`
  - `staticcheck`
  - `go test`
  - `go build`

Also includes **release automation** via `release-please`:

- Generates changelogs and version bumps based on **Conventional Commits**

Automation files are located in: `.github/workflows`

---

## 📇 Contacts

**Amir Amirov**

💼 [LinkedIn](https://www.linkedin.com/in/amir-amirov-6b2527338/)

✉️ [cgb.zko@gmail.com](mailto:cgb.zko@gmail.com)
