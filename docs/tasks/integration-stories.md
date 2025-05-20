# Vaani – Backend ⇆ Frontend Integration Sprint Backlog

**Purpose:** Track the bridging work once core backend & frontend MVP features are in place.
**Duration:** 1‑week dedicated integration sprint **(Sprint‑5)** immediately after FE/BE MVP sprints.
**Lead Engineer:** Full‑stack integrator
**QA:** SDET

---

## Epic INT‑0 – Contract Freeze & Shared Schemas

| ID          | Task                                                                                | Est | Acceptance Criteria                                                |
| ----------- | ----------------------------------------------------------------------------------- | --- | ------------------------------------------------------------------ |
| **INT‑0.1** | Freeze OpenAPI spec `v1` under `shared/contracts/openapi.yaml`                      | 2 h | Review signed by FE & BE leads; lint passes (`speccy lint`)        |
| **INT‑0.2** | Generate Go server stubs + Dart API client via `oapi-codegen` / `openapi-generator` | 3 h | Builds with no manual edits; clients imported in Gateway & Flutter |
| **INT‑0.3** | Publish versioned `@vaani/api` Dart package to GitHub Packages                      | 2 h | Flutter switches from raw dio calls to typed client                |

---

## Epic INT‑1 – Auth & Environment Wiring

| ID          | Task                                                            | Est | Acceptance Criteria                                             |
| ----------- | --------------------------------------------------------------- | --- | --------------------------------------------------------------- |
| **INT‑1.1** | Configure CORS & ALB headers for web origin                     | 1 h | Web PWA fetch succeeds w/o preflight error                      |
| **INT‑1.2** | Pass Firebase JWT from Flutter to Gateway via `dio` interceptor | 2 h | Protected endpoint returns 200 in integration test              |
| **INT‑1.3** | Ensure JWT refresh flow interoperates with Gateway 401 retry    | 3 h | Cypress test logs in → token expiration → silent refresh passes |

---

## Epic INT‑2 – Generate & Preview Flow

| ID          | Task                                                               | Est | Acceptance Criteria                                         |
| ----------- | ------------------------------------------------------------------ | --- | ----------------------------------------------------------- |
| **INT‑2.1** | Implement `/v1/preview` WS binary stream framing agreement         | 4 h | 3‑sec audio plays in Flutter <800 ms; frame ordering intact |
| **INT‑2.2** | Error propagation: Gateway error codes → Flutter SnackBar messages | 2 h | Simulated 429 displays quota toast                          |
| **INT‑2.3** | Editor cost‑meter calculation uses live `/pricing` endpoint        | 2 h | Cost matches backend billing ±1 % variance                  |

---

## Epic INT‑3 – File Upload & Download Pipeline

| ID          | Task                                                                           | Est | Acceptance Criteria                              |
| ----------- | ------------------------------------------------------------------------------ | --- | ------------------------------------------------ |
| **INT‑3.1** | Flutter downloads presigned URL & streams to File System                       | 2 h | 30‑sec clip downloads in <5 s on 50 Mbps network |
| **INT‑3.2** | Implement resumable download for flaky mobile (dio `download` w/ range header) | 3 h | Kill network mid‑download → resume succeeds      |
| **INT‑3.3** | Playback of remote MP3 via cache; fallback to streaming if >20 MB              | 2 h | Large file plays while downloading               |

---

## Epic INT‑4 – Quota & Billing Sync

| ID          | Task                                                                                | Est | Acceptance Criteria                             |
| ----------- | ----------------------------------------------------------------------------------- | --- | ----------------------------------------------- |
| **INT‑4.1** | WebSocket push or SSE feed for quota reduction events                               | 3 h | Quota chip decrements without app reload        |
| **INT‑4.2** | Razorpay web checkout deep‑link returns `payment_id` → Gateway verifies via webhook | 4 h | Payment success immediately adds credits in app |
| **INT‑4.3** | Weekend pass expiry timer in Flutter syncs with backend cron                        | 2 h | Pass credits disappear exactly at 24 h ±1 min   |

---

## Epic INT‑5 – End‑to‑End (E2E) Test Automation

| ID          | Task                                                                  | Est | Acceptance Criteria                |
| ----------- | --------------------------------------------------------------------- | --- | ---------------------------------- |
| **INT‑5.1** | Cypress web E2E: login → write script → preview → generate → download | 5 h | Test green in CI (Chrome)          |
| **INT‑5.2** | Flutter integration test (Android emulator GitHub Action) same path   | 6 h | Runs headless, passes under 4 min  |
| **INT‑5.3** | k6 load test hitting `/generate` + WebSocket preview concurrency 50   | 4 h | Error rate <1 %, P95 <2 s recorded |

---

## Epic INT‑6 – Observability Dashboards

| ID          | Task                                                                          | Est | Acceptance Criteria                            |
| ----------- | ----------------------------------------------------------------------------- | --- | ---------------------------------------------- |
| **INT‑6.1** | Frontend PostHog events: `preview_start`, `generate_success`, `download`      | 2 h | Events visible, linked to `userId`             |
| **INT‑6.2** | Grafana dashboard joins FE performance (`custom_metrics.js`) & BE API latency | 3 h | Single view shows correlation graph            |
| **INT‑6.3** | Alert: Preview start→first sound >1 s for >5 % users triggers Slack           | 2 h | Synthetic delay triggers alert in test channel |

---

## Epic INT‑7 – Security & Privacy Alignment

| ID          | Task                                                                         | Est | Acceptance Criteria                                    |
| ----------- | ---------------------------------------------------------------------------- | --- | ------------------------------------------------------ |
| **INT‑7.1** | Ensure Flutter only sends minimal PII; verify network logs via Charles Proxy | 2 h | No phone number or name in any API payload             |
| **INT‑7.2** | Sentry breadcrumbs include trace‑id to correlate with backend OTLP           | 2 h | Clicking FE error links to matching BE trace in Jaeger |

---

## Total Estimate

**≈ 71 dev hours** (\~9 working days)

> With two engineers (1 BE, 1 FE) pair‑working, sprint‑‑5 is realistic to fully integrate and pass all E2E tests.

---

## Sprint‑5 Calendar (suggested)

| Day   | Focus                                         |
| ----- | --------------------------------------------- |
| Mon   | Contract freeze, auth wiring (INT‑0, INT‑1)   |
| Tue   | Preview flow, error mapping (INT‑2)           |
| Wed   | File download & quota sync (INT‑3, INT‑4.1)   |
| Thu   | Billing callbacks, weekend pass (INT‑4.2‑4.3) |
| Fri   | E2E tests, dashboards, alerts (INT‑5, INT‑6)  |
| Sat\* | Buffer / bug‑bash                             |

\*optional stretch day if issues found during QA.

---

## Exit Criteria for Integration Sprint

1. **Cypress + Flutter integration pipelines green** in CI.
2. **MOS ≥4.5** in internal dog‑food session (10 users).
3. **P95 preview latency <0.8 s** measured in Play Console vitals.
4. End‑to‑end workflow usable on Android, iOS TestFlight, and Web.

---

**End of Backend‑Frontend Integration Backlog v0.1**
