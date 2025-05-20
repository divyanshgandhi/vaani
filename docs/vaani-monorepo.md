# Vaani Monorepo – Developer Reference (restructured)

**Repo URL (private):** `git@github.com:navoday/vaani.git`
**Primary stacks:** Go 1.22 (back‑end), Flutter 3.19 (front‑end web + mobile)
**State management:** `flutter_bloc`
**Revision:** 0.2 – 19 May 2025

---

## 1 Top‑level Layout

We now keep **backend** and **frontend** as first‑class directories so teams can work mostly in isolation while still sharing scripts and CI.

```text
vaani/
├─ backend/                  # All server‑side code & infra
│   ├─ services/             # Deployable Go binaries
│   │   ├─ api-gateway/
│   │   ├─ media-service/
│   │   ├─ bulbul-adapter/
│   │   └─ orpheus-inference/
│   ├─ libs/                 # Reusable Go pkgs (logger, auth, sdk)
│   ├─ infra/                # Terraform, Dockerfiles, bash helpers
│   └─ Makefile              # `make test`, `make docker`, etc.
│
├─ frontend/                 # Single Flutter workspace
│   ├─ app/                  # Flutter client (Android/iOS/Web)
│   └─ packages/             # Shared Dart pkgs (ui_kit, data_models)
│
├─ shared/                   # Assets & specifications
│   ├─ proto/                # .proto or OpenAPI specs
│   └─ design/               # Figma exports, brand assets
│
├─ .github/                  # CI/CD workflows & templates
├─ docs/                     # ADRs, architecture diagrams
├─ taskfiles/                # Taskfile.yml for uniform CLI tasks
├─ .config/                  # Lint configs, commitlint, etc.
└─ README.md
```

### Key Rules

* **backend/services/** and **frontend/app/** create your deployables; everything else is libraries or infra.
* Cross‑import path rules:

  * Backend services can only `require` packages inside `backend/libs/…` or Go stdlib.
  * Flutter code imports only from `frontend/packages/…`.
  * Neither side reaches into the other’s private code (use **shared/** for canonical DTOs).

---

## 2 Build & Tooling by Area

### 2.1 Backend

* **Go Modules** pinned at `go 1.22`.  Each service has its own `go.mod` that points to `../../libs/go-common` via `replace` while developing.
* **Mage** tasks in root *backend* Makefile:

  * `make test` – unit tests with `go test ./...`.
  * `make docker` – multi‑arch build.
  * `make dev-up` – docker‑compose to spin up entire stack locally.
* **Infra/Terraform** lives inside `backend/infra/terraform`.  Separate workspaces: `dev`, `staging`, `prod`.

### 2.2 Frontend (Flutter)

* **Melos** manages the workspace: run `melos bs` at `frontend/`.
* Shared UI kit and models live in `frontend/packages/`. These packages are versioned via `melos version` and imported with `path:` during development.
* **Web build** output is deployed via Firebase Hosting configured in `frontend/app/firebase.json`.

---

## 3 Unified Developer Workflow

| Stage        | Command                                                                  | Notes                                         |
| ------------ | ------------------------------------------------------------------------ | --------------------------------------------- |
| Setup        | `task bootstrap`                                                         | runs `go install`, `melos bs`, `direnv allow` |
| Local stack  | `backend/Makefile dev-up`                                                | docker-compose (gateway + orpheus + mocks)    |
| Tests        | `task test:all`                                                          | Go + Dart tests + linters                     |
| Commit       | Conventional Commits enforced via Husky (pre‑commit lint & tests)        |                                               |
| PR CI        | GitHub Actions matrix for Go 1.22 & Flutter stable                       |                                               |
| Merge → main | auto‑deploy **staging** via Terraform + GHCR images                      |                                               |
| Tag `vX.Y.Z` | manual approval → deploy **prod** services + release Flutter AAB/IPA/PWA |                                               |

*`task` is powered by **Taskfile.yml** at repo root for uniform cross‑language dev UX.*

---

## 4 CI/CD Workflows (updated paths)

| Workflow file        | Trigger     | What it does                                                                                                               |
| -------------------- | ----------- | -------------------------------------------------------------------------------------------------------------------------- |
| `ci-go.yml`          | push PR     | `cd backend && make test` + lint, build docker images into cache                                                           |
| `ci-flutter.yml`     | push PR     | `cd frontend && melos test` + `flutter build web --release` artefact upload                                                |
| `deploy-staging.yml` | push → main | Builds/pushes docker images under `backend/services/*`; Terraform apply to staging; Firebase web deploy to `/staging` site |
| `release.yml`        | tag `v*`    | Goreleaser for CLI, Docker push `latest`, Play Store internal track, TestFlight beta, Firebase prod deploy                 |

---

## 5 Local Dev Scripts

```bash
# Bring up Orpheus + mock Bulbul locally
dev-stack up
# Re-generate shared protobuf models after change
proto/gen.sh
# Build Flutter web in watch mode
melos run webdev
```

---

## 6 Secrets & Env Files

```
backend/infra/.env.sample   # SARVAM_KEY, JWT_SECRET …
frontend/app/.env.sample    # API_BASE_URL, POSTHOG_KEY …
```

Developers copy to `.env.local` per directory; **direnv** auto‑loads.

---

## 7 Open Questions

1. Do we want a **top‑level Makefile** proxy that calls into `backend/` or `frontend/` make targets for convenience?
2. Should shared **protobuf/OpenAPI** definitions move to a separate `shared/contracts` Git submodule for other Navoday repos?

---

**End of Monorepo Structure v0.2**
