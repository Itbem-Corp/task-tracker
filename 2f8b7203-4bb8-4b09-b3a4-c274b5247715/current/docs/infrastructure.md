# Infrastructure Documentation

## System Architecture

The Professional Services SaaS Platform consists of three services:

| Service | Technology | Port (Local) | Port (Container) | Directory |
|---------|-----------|-------------|-------------------|-----------|
| Landing Page | Astro.js 4.x + Tailwind | 3000 | 3000 | `saas-landing/` |
| Dashboard | Next.js 14 + React 18 | 3001 | 3000 | `saas-dashboard/` |
| API | Express 4 + Prisma 5 | 4000 | 3001 | `saas-api/` |

Supporting services:
- **PostgreSQL 16** - Primary database (port 5432)
- **Redis 7** - Session storage and caching (port 6379)

## Local Development

### Prerequisites
- Docker & Docker Compose
- Node.js 20+ (for running services outside Docker)

### Quick Start (Docker)

```bash
# From the project root (this repo)
docker compose up --build

# Services available at:
# Landing:   http://localhost:3000
# Dashboard: http://localhost:3001
# API:       http://localhost:4000
```

### Individual Service Development

```bash
# Landing (Astro.js)
cd saas-landing && npm install && npm run dev

# Dashboard (Next.js)
cd saas-dashboard && npm install && npm run dev

# API (Express)
cd saas-api && npm install && npm run dev
```

### Environment Variables

Each service has an `.env.example` file. Copy it to `.env` before running:

```bash
cp .env.example .env
```

#### API Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `3001` |
| `NODE_ENV` | Environment | `development` |
| `DATABASE_URL` | PostgreSQL connection string | (see .env.example) |
| `JWT_SECRET` | JWT signing secret | (required) |
| `JWT_EXPIRES_IN` | Token expiry | `7d` |
| `CORS_ORIGIN` | Allowed CORS origin | `http://localhost:3000` |
| `STRIPE_SECRET_KEY` | Stripe secret key | (required for payments) |
| `STRIPE_WEBHOOK_SECRET` | Stripe webhook secret | (required for webhooks) |

#### Dashboard Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `NEXT_PUBLIC_API_URL` | Backend API URL | `http://localhost:3001` |
| `NEXT_PUBLIC_STRIPE_PUBLISHABLE_KEY` | Stripe publishable key | (required for payments) |

#### Landing Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PUBLIC_SITE_URL` | Landing page URL | `http://localhost:3000` |
| `PUBLIC_DASHBOARD_URL` | Dashboard URL (for CTAs) | `http://localhost:3001` |
| `PUBLIC_API_URL` | API URL (for forms) | `http://localhost:4000` |

## Docker

### Build Strategy

All services use multi-stage Docker builds:

1. **saas-landing**: `node:20-alpine` (build) -> `nginx:1.27-alpine` (serve static)
2. **saas-dashboard**: `node:20-alpine` (build) -> `node:20-alpine` (standalone server, non-root user)
3. **saas-api**: `node:20-alpine` (build) -> `node:20-alpine` (production deps only, non-root user)

### Health Checks

| Service | Endpoint | Method |
|---------|----------|--------|
| Landing | `/health` | nginx returns JSON |
| Dashboard | `/` | Next.js serves page |
| API | `/health` | Express endpoint |
| PostgreSQL | `pg_isready` | CLI check |
| Redis | `redis-cli ping` | CLI check |

## CI/CD

Each service has a GitHub Actions workflow at `.github/workflows/ci.yml`.

### Pipeline Stages

All services follow the same pattern:

1. **Lint** - ESLint / Astro check
2. **Test** - Vitest (API includes PostgreSQL service container)
3. **Build** - TypeScript compilation / static generation
4. **Docker Build** - Multi-stage image build (main branch only)

### Triggers

- **Push** to `main` or `develop` branches
- **Pull requests** targeting `main`

### API-Specific CI

The API pipeline spins up a PostgreSQL 16 service container for integration tests, with automatic Prisma migration.

## Database

- **Engine**: PostgreSQL 16 (Alpine)
- **ORM**: Prisma 5.10+
- **Schema**: `saas-api/prisma/schema.prisma`

### Database Commands

```bash
# Generate Prisma client
cd saas-api && npx prisma generate

# Run migrations
npx prisma migrate deploy

# Push schema changes (dev only)
npx prisma db push

# Open Prisma Studio
npx prisma studio
```
