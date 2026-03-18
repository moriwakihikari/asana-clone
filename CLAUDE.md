# Asana Clone Project

## Overview
Full-stack Asana clone with Go DDD backend and Next.js frontend.

## Tech Stack
- **Frontend**: Next.js 14 (App Router) + TypeScript + TailwindCSS
- **Backend**: Go (chi router) with DDD architecture
- **Database**: PostgreSQL 16 + Redis 7
- **Infrastructure**: Docker Compose

## Project Structure
```
asana/
├── docker-compose.yml
├── backend/          # Go API (DDD: domain → application → infrastructure → interfaces)
├── frontend/         # Next.js + TailwindCSS
├── docker/           # Docker configs
└── docs/             # Design docs
```

## Backend Architecture (DDD)
- `internal/domain/` - Entities, Value Objects, Repository interfaces
- `internal/application/` - Use cases / Application services
- `internal/infrastructure/` - PostgreSQL repos, Redis cache, JWT
- `internal/interfaces/` - HTTP handlers, DTOs, middleware

## Ports
- Frontend: localhost:3000
- Backend API: localhost:8080
- PostgreSQL: localhost:5432
- Redis: localhost:6379

## CI/Build Rules
- Run `go vet ./...` and `go build ./...` before committing Go code
- Run `npm run lint` and `npm run build` before committing frontend code

## Languages & Conventions
- Go: Follow standard Go conventions, use chi router
- TypeScript: Strict mode, use TanStack Query for data fetching
- SQL: Use golang-migrate for migrations
