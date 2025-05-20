# Backend ED

---

# Vaani – Backend Architecture & Engineering Spec (v0.1)

**Scope:** MVP launch supporting Bulbul v2 (managed) & Orpheus‑TTS (self‑host) only

**Audience:** Backend engineers, DevOps, security, and SRE

**Revision Date:** 19 May 2025

---

## 1 Objectives

- Deliver < 0.5 s P90 preview latency for both Indic & English requests.
- Gross COGS ≤ ₹0.18 / k characters at 5 k paid users.
- 99.5 % availability (single‑region) for first 12 months.
- HIPAA‑lite, DPDP‑India & GDPR compliance for voice data.

---

## 2 High‑Level Component Map

```mermaid
graph TD
  subgraph Edge
    C1[Flutter / Web Client]
  end
  subgraph Cloud (ap-south-1)
    AGW(API‑Gateway Go)
    BA[Bulbul Adapter]
    OSP[Orpheus GPU Pod]
    MS[Media Service]
    FStore[(Firestore)]
    S3((Wasabi Bucket))
  end
  C1 -- HTTPS/WebSocket --> AGW
  AGW -- Indic TTS --> BA
  BA -- HTTPS --> BulbulAPI[(Sarvam Cloud)]
  AGW -- EN/Hinglish --> OSP
  OSP -- PCM chunks --> AGW
  AGW --> MS
  MS --> S3
  MS --> FStore
  AGW --> FStore

```

---

## 3 Component Detail

### 3.1 API‑Gateway (Go 1.22)

- **Endpoints** (external):
    - `POST /v1/generate` – generate full clip
    - `WS /v1/preview` – bidirectional stream for 3 s preview
    - `POST /v1/clone` – create/update user voice clone
    - `GET /v1/download/:id` – pre‑signed link redirect
- **Responsibilities**
    1. OAuth‑less **OTP auth** (Firebase Phone Auth) issuing JWT (`sub=userId`).
    2. **Language detector** – fastText 176‑lang binary in memory (< 5 µs).
    3. Rate‑limiter (redis‑less token bucket in memory per‑user → persisted every 10 s).
    4. Routing logic: Indic → BA, else EN → Orpheus.
    5. Quota decrement & billing meter (pub/sub to Firestore via batched writes).
    6. WebSocket multiplexer – first 3 s from OSP piped back to client for instant UX.
- **Scaling:** stateless; horizontal pod autoscaled on CPU>60 % (ECS or Fly.io). One t3a.medium per AZ initially.

### 3.2 Bulbul Adapter

- Thin Go client; signs each request with `X‑Sarvam‑Key`, forwards JSON with:
    
    ```json
    {"text":"…","lang":"hi","style":"excited","alpha":0.6}
    
    ```
    
- Caches successful responses SHA‑1(text+voice) → S3 (30‑day TTL) to reduce API spend by ~12 %.

### 3.3 Orpheus‑TTS Inference Pod

- Docker image: `ghcr.io/canopyai/orpheus:v0.32-cuda11`
- Runtime stack: **vLLM 0.5** + Triton gRPC gateway.
- GPU: `g5.xlarge` (Nvidia A10G 24 GB).
- **Concurrency:** 4 streams; dynamic batching every 30 ms.
- **Throughput:** ≈ 35 k char/min at 24 kHz (measured).
- Streams raw PCM (float32) frames over gRPC→Gateway (frame = 100 ms).
- Health check every 3 s; auto‑restart if 3 fails.

### 3.4 Media Service

- Encodes PCM → WAV (48 kHz) or LAME MP3 320 kbps.
- Uses **`go-audio`**; CPU‑bound – colocated with API‑GW.
- Uploads to Wasabi bucket with `private` ACL; returns pre‑signed URL (6 h validity).
- Writes Firestore doc: `{userId, projectId, chars, cost, lang, engine, durationMs}`.

### 3.5 Data Store – Firestore

```
users/{uid}
  credits: int  (chars)
  clones: map {voiceId: {engine, created}}
projects/{pid}
  owner: uid
  created, updated
  outputs: []
itx/{dailyId}
  totalChars
  bulbulChars
  orpheusChars

```

- Composite index: `(owner, created)` for dashboard list.

### 3.6 Auth & Identity

- Firebase PhoneAuth → JWT (30 d).
- Refresh via silent re‑verification.
- Admin users list stored in Firestore for dashboards.

### 3.7 Billing & Payments

- Firestore tally triggers **Cloud Function** each midnight:
    1. Sum `chars` per user.
    2. If > free cap, set `billableChars`.
    3. Razorpay Subscriptions API charges Creator / Pro.
- Weekend pass: one‑off Razorpay Order; adds temporary `extraCredits` valid 24 h.

---

## 4 API Contract (external)

### 4.1 Generate

```
POST /v1/generate
Authorization: Bearer <jwt>
{
  "text": "नमस्ते दुनिया",
  "voice_id": "bulbul_f_hindi_01",
  "emotion": 62,
  "format": "mp3"
}
→ 202 Accepted
{
  "job_id": "673a…",
  "estimate_ms": 3500
}

```

`GET /v1/jobs/{id}` → `{status, url?}`

### 4.2 Preview WebSocket

- Client sends `{text, voice_id, emotion}`.
- Server streams frames: `{seq, base64_pcm}`.
- Close when 3 s sent or error.

---

## 5 Security & Compliance

| Area | Measure |
| --- | --- |
| Data in transit | TLS 1.3, HSTS, ALPN required |
| Data at rest | Wasabi bucket encryption (AES‑256 SSE‑S3) |
| PII | Only phone number stored (hashed w/ bcrypt) |
| Voice consent | `/clone` requires spoken phrase "I consent to use my voice…" auto‑STT validated |
| Abuse detection | Celebrity list hashed via Bloom filter; rejects matches > 0.85 sim |

---

## 6 Observability

- **Logs** – Loki via Grafana Cloud; JSON structured.
- **Metrics** – Prometheus (api_latency_ms, gpu_util, chars_generated).
- **Tracing** – OpenTelemetry, Jaeger; sampling 1 %.
- **Alerts** – Grafana + Slack:
    - P95 API latency > 1 s for 5 min
    - GPU util <10 % for 1 h (over‑provision)

---

## 7 CI/CD & Environments

| Stage | Branch | Infra | Tests |
| --- | --- | --- | --- |
| *Dev* | `dev/*` | Fly.io tmp app | Unit (Go + Ginkgo) |
| *Staging* | `main` | ap‑south‑1 staging VPC | e2e (10 cases) |
| *Prod* | Tag vX.Y.Z | prod VPC | blue‑green deploy, smoke |
- GitHub Actions builds multi‑arch images; `docker‑buildx` pushes to GHCR.
- Terraform Cloud manages VPC, ECS, S3, security groups.

---

## 8 Capacity & Cost Model

- **Orpheus**: g5.xlarge @ ₹5.5 / h = ₹132 / day.
    
    *Break‑even* 1.6 M chars/day (₹288 rev > ₹132 cost).
    
- **Bulbul**: ₹0.20 / k → ₹200 / M chars.
- At 5 k Creator subs using 50 k chars/mo = 250 M chars → ₹5 L Bulbul COGS per month (covered by ₹19.9 L rev).
- Headroom: 4× GPU before margin < 80 %.

---

## 9 Disaster Recovery & Backup

- Wasabi bucket versioning + 30‑day lifecycle.
- Firestore daily export to GCS coldline.
- Terraform `apply` redeploy script < 15 min in new region (ap‑southeast‑1) using AMI snapshot of Orpheus image.

---

## 10 Future Extensions (post‑MVP)

1. GPU auto‑scaler (Karpenter) + multi‑GPU Orpheus serving.
2. Redis pub/sub for quota + realtime analytics.
3. Secondary TTS (OpenVoice) sidecar micro‑service.
4. On‑device TTS (Indic‑TTS) offline SDK.

---

**End of Backend Engineering Spec v0.1**