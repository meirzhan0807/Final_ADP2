# 🏢 Job Finder Platform

Production-grade **microservices** job platform built with **Go**, **gRPC**, **NATS**, **PostgreSQL**, **Redis** and full observability.

## ✅ Requirements Coverage

| Requirement | Points | Status |
|---|---|---|
| Clean Architecture (Repository → Service → Handler) | 20% | ✅ |
| **36 gRPC Endpoints** (12 × 3 services) | 20% | ✅ |
| **NATS** Message Queue (events: register, apply, status) | 20% | ✅ |
| **PostgreSQL** × 3 + migrations + transactions + **Redis** cache | 20% | ✅ |
| **Email** HTML templates (Gmail / Microsoft SMTP) | 10% | ✅ |
| **Unit + Integration** Tests (mock + real) | 10% | ✅ |
| **Grafana + Jaeger + Prometheus + Loki** (BONUS) | 10% | ✅ |

**Total: 110%**

---

## 🏗️ Architecture

```
                    ┌─────────────────────┐
   HTTP :8080       │     API Gateway     │   Gin + JWT Auth
                    │   (api-gateway)     │   Rate Limit, CORS
                    └──────────┬──────────┘
                               │ gRPC
           ┌───────────────────┼───────────────────┐
           │                   │                   │
    ┌──────▼──────┐    ┌───────▼──────┐   ┌────────▼───────┐
    │ User Service│    │ Job Service  │   │ Notif Service  │
    │  :50051     │    │   :50052     │   │    :50053      │
    └──────┬──────┘    └──────┬───────┘   └────────┬───────┘
           │                   │                   │
           └───────────────────┼───────────────────┘
                               │
                    ┌──────────▼──────────┐
                    │   NATS :4222        │  Message Queue
                    │ (user.registered,   │  Event-driven
                    │  job.applied,       │  notifications
                    │  application.status)│
                    └─────────────────────┘
           │                   │                   │
    ┌──────▼──────┐    ┌───────▼──────┐   ┌────────▼───────┐
    │ PostgreSQL  │    │ PostgreSQL   │   │  PostgreSQL    │
    │  (userdb)   │    │  (jobdb)     │   │ (notificationdb)│
    └─────────────┘    └──────────────┘   └────────────────┘
                               │
                    ┌──────────▼──────────┐
                    │    Redis :6379      │  Cache + Sessions
                    └─────────────────────┘
```

---

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose

### 1. Start
```bash
chmod +x scripts/start.sh
./scripts/start.sh
```

Or manually:
```bash
cp .env.example .env
docker compose up -d --build
```

### 2. Test
```bash
# Health
curl http://localhost:8080/health

# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H 'Content-Type: application/json' \
  -d '{"email":"employer@test.com","password":"password123","first_name":"John","last_name":"Doe","role":"employer"}'

# Login → get token
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"employer@test.com","password":"password123"}' | python3 -c "import sys,json; print(json.load(sys.stdin)['access_token'])")

# Create job
curl -X POST http://localhost:8080/api/v1/employer/jobs \
  -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"title":"Senior Go Developer","description":"Backend microservices role","location":"Almaty","job_type":"full-time","category":"Technology","salary_min":300000,"salary_max":500000,"currency":"KZT","experience_level":"senior"}'

# Search jobs
curl "http://localhost:8080/api/v1/jobs/search?q=developer&location=Almaty"
```

---

## 🌐 API Endpoints

### Public
| Method | URL | Description |
|---|---|---|
| GET | `/health` | Health check |
| POST | `/api/v1/auth/register` | Register (role: jobseeker/employer) |
| POST | `/api/v1/auth/login` | Login → JWT tokens |
| POST | `/api/v1/auth/refresh` | Refresh access token |
| GET | `/api/v1/auth/verify-email?token=` | Verify email |
| GET | `/api/v1/jobs` | List jobs |
| GET | `/api/v1/jobs/search` | Search: `?q=&location=&category=&job_type=&salary_min=` |
| GET | `/api/v1/jobs/categories` | Job categories |
| GET | `/api/v1/jobs/:id` | Job details |

### Authenticated (Bearer token)
| Method | URL | Description |
|---|---|---|
| GET | `/api/v1/me` | My account |
| PUT | `/api/v1/me` | Update account |
| PUT | `/api/v1/me/password` | Change password |
| GET | `/api/v1/me/profile` | My profile (bio, skills, resume) |
| PUT | `/api/v1/me/profile` | Update profile |
| GET | `/api/v1/me/applications` | My job applications |
| POST | `/api/v1/jobs/:id/apply` | Apply for job |
| GET | `/api/v1/notifications` | Notifications |
| GET | `/api/v1/notifications/unread-count` | Unread count |
| PUT | `/api/v1/notifications/:id/read` | Mark read |
| DELETE | `/api/v1/notifications/:id` | Delete notification |

### Employer (role: employer)
| Method | URL | Description |
|---|---|---|
| POST | `/api/v1/employer/jobs` | Post a job |
| PUT | `/api/v1/employer/jobs/:id` | Update job |
| DELETE | `/api/v1/employer/jobs/:id` | Delete job |
| GET | `/api/v1/employer/jobs` | My job postings |
| GET | `/api/v1/employer/jobs/:id/applications` | Job applicants |
| PUT | `/api/v1/employer/applications/:id/status` | Update status (pending/reviewed/accepted/rejected) |

---

## 📊 gRPC Endpoints (36 total)

### User Service (12)
`Register` · `Login` · `GetUser` · `UpdateUser` · `DeleteUser` · `ListUsers` · `ChangePassword` · `VerifyEmail` · `RefreshToken` · `GetUserProfile` · `UpdateUserProfile` · `ValidateToken`

### Job Service (12)
`CreateJob` · `GetJob` · `UpdateJob` · `DeleteJob` · `ListJobs` · `SearchJobs` · `ApplyJob` · `GetApplications` · `UpdateApplicationStatus` · `GetJobsByEmployer` · `GetApplicationsByUser` · `GetJobCategories`

### Notification Service (12)
`SendEmail` · `SendWelcomeEmail` · `SendJobApplicationEmail` · `SendApplicationStatusEmail` · `SendPasswordResetEmail` · `GetNotifications` · `MarkNotificationRead` · `GetUnreadCount` · `DeleteNotification` · `SendBulkEmail` · `SendJobAlertEmail` · `CreateNotification`

---

## 📧 Email Setup

**Works without email** — if `SMTP_USER` is empty, emails are logged to console (mock mode).

### Gmail
1. Enable 2FA on Google account
2. Go to: https://myaccount.google.com/apppasswords
3. Create App Password for "Mail"
4. Add to `.env`:
```env
SMTP_USER=your@gmail.com
SMTP_PASSWORD=xxxx-xxxx-xxxx-xxxx
```

### Microsoft/Outlook
```env
SMTP_HOST=smtp.office365.com
SMTP_PORT=587
SMTP_USER=your@outlook.com
SMTP_PASSWORD=your-password
```

---

## 📈 Observability

| Service | URL | Credentials |
|---|---|---|
| Grafana (metrics + logs) | http://localhost:3000 | admin / admin123 |
| Prometheus | http://localhost:9090 | — |
| Jaeger (tracing) | http://localhost:16686 | — |
| NATS Monitor | http://localhost:8222 | — |

---

## 🧪 Running Tests

```bash
# Unit tests (no external services needed)
cd user-service && go test ./tests/... -short -v
cd job-service && go test ./tests/... -short -v
cd notification-service && go test ./tests/... -short -v

# Integration tests (requires running docker-compose)
cd user-service && go test ./tests/... -v -run TestIntegration
```

---

## 🗂️ Project Structure

```
job-finder-platform/
├── api-gateway/              # HTTP → gRPC proxy (Gin)
│   ├── cmd/main.go
│   └── internal/
│       ├── config/
│       ├── handlers/         # All HTTP handlers
│       ├── middleware/        # Auth, RateLimit, CORS
│       └── pb/               # gRPC client stubs
│           ├── user/
│           ├── job/
│           └── notification/
│
├── user-service/             # User management (gRPC :50051)
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── db/               # PostgreSQL + Redis setup + migrations
│   │   ├── handlers/         # gRPC handlers
│   │   ├── messaging/        # NATS publisher
│   │   ├── models/
│   │   ├── pb/               # gRPC server definitions
│   │   ├── repositories/     # Data access layer
│   │   └── services/         # Business logic + JWT
│   └── tests/
│
├── job-service/              # Jobs & Applications (gRPC :50052)
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── db/
│   │   ├── handlers/
│   │   ├── messaging/
│   │   ├── models/
│   │   ├── pb/
│   │   ├── repositories/
│   │   └── services/
│   └── tests/
│
├── notification-service/     # Email + In-app (gRPC :50053)
│   ├── cmd/main.go
│   ├── internal/
│   │   ├── db/
│   │   ├── email/            # SMTP + HTML templates
│   │   ├── handlers/
│   │   ├── models/
│   │   ├── pb/
│   │   ├── repositories/
│   │   └── services/         # NATS event listeners
│   └── tests/
│
├── docker/
│   ├── prometheus/
│   ├── grafana/
│   ├── loki/
│   └── promtail/
│
├── scripts/start.sh
├── docker-compose.yml
├── .env.example
└── README.md
```

---

## 🛑 Stop / Reset

```bash
docker compose down          # stop
docker compose down -v       # stop + delete all data
```

---

## 🌐 Frontend (Бонус +10%)

Веб-интерфейс `frontend/index.html` файлында орналасқан.

### Іске қосу
1. `docker compose up -d` арқылы backend-ті іске қосыңыз
2. `frontend/index.html` файлын браузерде ашыңыз

### Функционалдылық
- 🔐 Тіркелу / Кіру (jobseeker & employer)
- 📋 Жұмыстар тізімі, іздеу, фильтр
- 📝 Жұмысқа өтініш беру
- 👤 Профиль басқару
- 🏢 Жұмыс беруші панелі (жариялау, өтініштер)
- 🔔 Хабарламалар
