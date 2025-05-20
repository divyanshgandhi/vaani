# Vaani – Backend MVP Sprint Backlog

**Release target:** Public beta in 8 weeks (see PRD)
**Methodology:** 1‑week sprints, Kanban board in GitHub Projects
**Issue Label Conventions:** `epic`, `story`, `task`, `bug`, `infra`, `test`

---

## Epic 1 – Repository & Tooling

| ID        | Story / Task                                                   | Estimate | Acceptance Criteria                                       |
| --------- | -------------------------------------------------------------- | -------- | --------------------------------------------------------- |
| **B‑1.1** | Initialise monorepo with *backend* / *frontend* dirs           | 2 h      | `git clone`, `go run` passes, `flutter doctor` passes     |
| **B‑1.2** | Add root `Taskfile.yml` with `bootstrap`, `test:all`, `dev-up` | 4 h      | Running `task bootstrap` installs Go + Melos deps locally |
| **B‑1.3** | Configure direnv & sample `.env` files                         | 1 h      | `echo $SARVAM_KEY` surfaced when `.env.local` present     |
| **B‑1.4** | Setup GitHub Actions `ci-go.yml` skeleton                      | 3 h      | PR triggers Go test matrix, fails on lint error           |

---

## Epic 2 – API‑Gateway Service

### Story B‑2.0 *(Epic definition)* – A single Go HTTP/WS gateway that authenticates users, routes TTS requests, enforces quota, and streams previews.

| ID         | Task                                                                          | Est | Acceptance Criteria                                              |
| ---------- | ----------------------------------------------------------------------------- | --- | ---------------------------------------------------------------- |
| **B‑2.1**  | Scaffold Go `api-gateway` service w/ chi router                               | 4 h | `go run main.go` serves health endpoint `/healthz`               |
| **B‑2.2**  | Integrate Firebase PhoneAuth JWT middleware                                   | 6 h | Valid JWT in `Authorization` allows request; invalid returns 401 |
| **B‑2.3**  | Implement fastText language detection util (C binding or `alexcesaro/lingua`) | 4 h | Unit test > 95 % accuracy on 100‑sentence set                    |
| **B‑2.4**  | Add in‑memory token‑bucket rate‑limiter (per‑user)                            | 3 h | >5 requests/sec returns 429; reset after 60 s                    |
| **B‑2.5**  | REST `POST /v1/generate` stub with validation (<=5k chars)                    | 4 h | Returns 202 + job\_id JSON                                       |
| **B‑2.6**  | Firestore integration – record job metadata                                   | 4 h | Doc created under `projects/{pid}` with status `queued`          |
| **B‑2.7**  | WebSocket `/v1/preview` pass‑through (no backend yet)                         | 6 h | Echo text frames → client receives base64 dummy audio            |
| **B‑2.8**  | Route: Indic -> Bulbul adapter, English -> Orpheus                            | 4 h | Unit test using mocked adapters verifies correct path            |
| **B‑2.9**  | Quota decrement logic + 429 on exhaustion                                     | 4 h | Free user limited to 5 min/day indistinct                        |
| **B‑2.10** | Dockerfile multi‑stage build (scratch runtime)                                | 3 h | Image < 25 MB; `make docker` succeeds                            |
| **B‑2.11** | Unit tests (≥85 % cov) & bench for p95 latency <5 ms                          | 5 h | `go test ./... -cover` shows 85 %                                |

---

## Epic 3 – Bulbul Adapter

| ID        | Task                                                   | Est | Acceptance Criteria                                  |
| --------- | ------------------------------------------------------ | --- | ---------------------------------------------------- |
| **B‑3.1** | Build Go client pkg `sarvam` with request/resp structs | 4 h | Can POST sample payload to sandbox key & receive WAV |
| **B‑3.2** | Emotion slider mapping util                            | 2 h | Input 0–100 outputs `{style,alpha}` JSON fields      |
| **B‑3.3** | Add SHA‑1 text+voice caching layer (S3)                | 6 h | Cache hit returns in <50 ms; miss hits Sarvam        |
| **B‑3.4** | Integration tests using Sarvam sandbox env             | 4 h | Generates audio for 3 langs with >4 MOS              |

---

## Epic 4 – Orpheus Inference Pod

| ID        | Task                                                                          | Est | Acceptance Criteria                          |
| --------- | ----------------------------------------------------------------------------- | --- | -------------------------------------------- |
| **B‑4.1** | Create Docker build context w/ CUDA base, vLLM, Orpheus model download script | 6 h | `docker run ... orpheus serve` binds on 9000 |
| **B‑4.2** | Implement gRPC streaming server returning 100 ms PCM frames                   | 6 h | Demo client gets first frame <300 ms         |
| **B‑4.3** | Healthcheck endpoint `/healthz` probes GPU status                             | 3 h | Returns 200 JSON {gpu\_util:…}               |
| **B‑4.4** | Helm or Compose manifest for dev stack                                        | 3 h | `make dev-up` includes GPU‑enabled container |
| **B‑4.5** | Load‑test: 4 concurrent streams ≤60 % GPU util                                | 4 h | Locust report archived in `/docs/perf/`      |

---

## Epic 5 – Media Service

| ID        | Task                                                       | Est | Acceptance Criteria                                       |
| --------- | ---------------------------------------------------------- | --- | --------------------------------------------------------- |
| **B‑5.1** | Build Go `media-service`, expose gRPC `Encode(AudioChunk)` | 5 h | Receives PCM, writes WAV, uploads to Wasabi               |
| **B‑5.2** | Integrate LAME MP3 encoding via `go-lame`                  | 3 h | MP3 320 kbps output verified by ffprobe                   |
| **B‑5.3** | Implement S3 presigned URL generator (6 h TTL)             | 2 h | Returns HTTP 302 to signed URL                            |
| **B‑5.4** | Firestore `outputs` doc update on success                  | 3 h | Project status -> `complete` triggers front‑end poll pass |

---

## Epic 6 – Infra & Terraform

| ID        | Task                                        | Est | Acceptance Criteria                                |
| --------- | ------------------------------------------- | --- | -------------------------------------------------- |
| **B‑6.1** | Write VPC, subnets, IGW, SG modules         | 6 h | `terraform plan` shows zero drift                  |
| **B‑6.2** | ECS Fargate for api-gateway & media-service | 6 h | Staging URL reachable, ALB HTTPS cert              |
| **B‑6.3** | EC2 g5.xlarge launch template for Orpheus   | 4 h | GPU node up with user‑data auto‑pull │             |
| **B‑6.4** | Wasabi bucket + IAM access policy           | 2 h | Upload via Terraform proves ACL private            |
| **B‑6.5** | Firestore indexes & SA keys via TF          | 3 h | Index `projects(owner,created)` exists             |
| **B‑6.6** | GitHub OIDC -> AWS role for CD              | 2 h | GH Action can `sts:AssumeRole` without static keys |

---

## Epic 7 – Observability & Alerts

| ID        | Task                                                    | Est | Acceptance Criteria                     |
| --------- | ------------------------------------------------------- | --- | --------------------------------------- |
| **B‑7.1** | Integrate zap + OpenTelemetry trace exporter in Go libs | 4 h | Traces visible in Jaeger UI             |
| **B‑7.2** | Deploy Grafana Cloud agent for logs + metrics           | 3 h | `api_latency_ms` grafana panel imported |
| **B‑7.3** | Slack alert channel via PrometheusRules                 | 2 h | Latency alert >1 s fires test message   |

---

## Epic 8 – Billing & Quota

| ID        | Task                                            | Est | Acceptance Criteria               |
| --------- | ----------------------------------------------- | --- | --------------------------------- |
| **B‑8.1** | Firestore Cloud Function nightly job aggregator | 4 h | Updates `users.credits` correctly |
| **B‑8.2** | Razorpay Subscriptions integration              | 6 h | Test card charges ₹1 sandbox      |
| **B‑8.3** | Weekend pass one‑off order API                  | 4 h | Credits expire after 24 h cron    |

---

## Epic 9 – Security & Compliance

| ID        | Task                                             | Est | Acceptance Criteria                       |
| --------- | ------------------------------------------------ | --- | ----------------------------------------- |
| **B‑9.1** | Implement consent phrase STT validator (Whisper) | 4 h | Rejects if <90 % transcript match         |
| **B‑9.2** | Celebrity voice Bloom‑filter checker             | 6 h | Reject sample with ≥0.85 similarity to DB |
| **B‑9.3** | Enable TLS 1.3, HSTS on ALB                      | 2 h | Qualys SSL Labs grade A                   |

---

## Epic 10 – QA & Stabilisation

| ID         | Task                                                     | Est | Acceptance Criteria        |
| ---------- | -------------------------------------------------------- | --- | -------------------------- |
| **B‑10.1** | End‑to‑end tests (httpexpect) for `/generate` happy path | 5 h | Passes in CI, 3 assertions |
| **B‑10.2** | Load test 100 parallel generate calls                    | 4 h | Error rate <1 %, p95 <2 s  |
| **B‑10.3** | Pen‑test checklist (OWASP ASVS level 1)                  | 6 h | No criticals open          |

---

## Total Estimate

* ≈ **143 developer hours** (\~18 working days).  Buffer 30 % → 24 days ≈ 5 sprints to MVP.

---

## Next Steps

1. Copy backlog into GitHub Projects board; assign owners & sprint labels.
2. Sprint‑0 (setup) = Epic 1 + infra skeleton tasks B‑6.1/B‑6.4/B‑6.6.
3. Schedule daily stand‑ups and weekly sprint review/demo.

---

**End of Backend MVP Sprint Backlog**
