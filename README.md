# Vaani - Multilingual Voice AI Platform

Vaani is a multilingual voice AI platform that uses advanced ML models for voice synthesis, recognition, and transformation. This monorepo contains both backend (Go) and frontend (Flutter) code.

## Repository Structure

```
vaani/
├─ backend/                  # All server-side code & infra
│   ├─ services/             # Deployable Go binaries
│   ├─ libs/                 # Reusable Go pkgs
│   ├─ infra/                # Terraform, Dockerfiles, helpers
│   └─ Makefile
├─ frontend/                 # Single Flutter workspace
│   ├─ app/                  # Flutter client
│   └─ packages/             # Shared Dart pkgs
├─ shared/                   # Assets & specifications
├─ .github/                  # CI/CD workflows
├─ docs/                     # Documentation
└─ taskfiles/                # Taskfile for dev tasks
```

## Getting Started

### Prerequisites

- Go 1.22+
- Flutter 3.19+
- Docker & Docker Compose
- Task (taskfile.dev)
- Melos

### Setup

1. Clone the repository:
   ```
   git clone git@github.com:navoday/vaani.git
   cd vaani
   ```

2. Bootstrap the project:
   ```
   task bootstrap
   ```

3. Start the local development stack:
   ```
   cd backend && make dev-up
   ```

4. Run the Flutter web app:
   ```
   cd frontend/app && flutter run -d chrome
   ```

## Development Workflow

- Use `task` commands for unified workflows
- Backend: `cd backend && make test`
- Frontend: `cd frontend && melos test`

## Docker

The entire stack can be run using Docker Compose:

```
cd backend && make dev-up
```

## CI/CD

- PRs: Automated tests and builds
- Merge to main: Deploy to staging
- Tag release: Deploy to production

## Documentation

- See `/docs` directory for architecture decisions and diagrams
- Refer to `docs/vaani-monorepo.md` for monorepo structure details

## License

Proprietary - Navoday AI
