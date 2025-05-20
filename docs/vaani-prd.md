# PRD

---

# Vaani – MVP Product Requirements Document (PRD)

**Revision:** 0.2 | **Date:** 19 May 2025

---

## 1 Purpose & Vision

Ship, in **≤ 8 weeks**, a cross‑platform voice‑AI studio that:

- Generates **studio‑quality TTS** for English and 11 Indian languages.
- Offers **instant voice cloning** and an **emotion slider**.
- Costs creators **≤ ₹399 / month** at launch.

> Model strategy (finalised)
> 
> 
> • **Bulbul v2** (Sarvam AI managed API) for all Indic‑language synthesis.
> 
> • **Orpheus‑TTS** (Apache‑2.0, self‑hosted) for English & Hinglish.
> 

No other engines/models are in scope for MVP.

---

## 2 Success Metrics (90 days post‑launch)

| Metric | Target |
| --- | --- |
| Weekly Active Users | **30 k** |
| Free → Paid conversion | **≥ 8 %** |
| Average preview latency (P90) | **≤ 0.5 s** |
| COGS per 1 M chars | **≤ ₹500** |
| CSAT | **≥ 4.5 / 5** |

---

## 3 Target Personas & JTBD

| Persona | Job‑to‑Be‑Done | Pain Today |
| --- | --- | --- |
| Regional YouTuber | Narrate Hindi/English videos fast | Expensive VO, robotic Hindi voices |
| Student / Accessibility user | Listen to study notes in own voice | Speechify ₹11 k/yr, no Hindi emotion |
| Reels Creator | Make catchy VO on phone | No end‑to‑end mobile tool |
| Tutor | Dub lessons to other Indian languages | Manual dubbing slow/costly |

---

## 4 MVP Functional Requirements

| # | Feature | Details |
| --- | --- | --- |
| F‑1 | **Neural TTS – Indic** | Bulbul API. Languages: hi, ta, te, ml, kn, mr, gu, bn, pa, or, en-IN. |
| F‑2 | **Neural TTS – English** | Orpheus streaming w/ zero‑shot cloning. |
| F‑3 | **Instant Voice Cloning** | ≤ 30 s user sample → Orpheus prompt. 1 clone on Free, 3 on Creator plan. |
| F‑4 | **Emotion Slider** | 0‑100 mapped: calm→excited (Bulbul `style`, Orpheus `<laugh>` tokens). |
| F‑5 | **Script Editor + Auto‑punctuate** | AI fixes commas &  tags; shows cost meter. |
| F‑6 | **Real‑time Preview** | Stream first 3 s audio; full clip after render. |
| F‑7 | **Lite Mixer** | 5 tracks: multiple voices + 1 music FX. Drag & drop, snap alignment. |
| F‑8 | **Export** | WAV/MP3 and caption SRT. |
| F‑9 | **Auth & Sync** | OTP login; projects & quotas in Firestore. |

*Out of MVP:* cross‑lang dubbing, lip‑sync, marketplace, podcast mode, multi‑GPU autoscale.

---

## 5 Simplified Architecture

```
Flutter / Next.js Clients
        │  HTTPS / WebSocket
        ▼
┌──────────────────┐
│  API‑Gateway (Go)│  – auth, rate‑limit, lang‑detect
└───────┬──────────┘
        │
  Indic │          English/Hinglish
        ▼                      ▼
┌────────────────┐    ┌────────────────────┐
│Bulbul Adapter   │    │Orpheus Inference   │
│(HTTPS → Sarvam) │    │Docker + vLLM + GPU │
└────────┬───────┘    └─────────┬──────────┘
         │                       │
         ▼                       ▼
      Media Service (Go) – encode, upload
         │
  S3‑compatible storage (Wasabi)
         │
      Firestore – users, projects, quotas

```

### Deployment BOM

| Component | Runtime | Instance | Cost/mo* |
| --- | --- | --- | --- |
| API & Media svc | Go | t3a.medium | $20 |
| Orpheus GPU | Docker | g5.xlarge (A10G) | $250 |
| Bulbul API | SaaS | pay‑as‑you‑go | ₹0.20 /k‑char |
| Storage (Wasabi) | — | 100 GB | $6 |
| Firebase | — | free tier | $0 |
| *initial traffic (<50 k previews/day) |  |  |  |

---

## 6 Timeline (8‑week burst)

| Week | Milestone |
| --- | --- |
| 1 | Orpheus container live; local streaming demo |
| 2 | Bulbul Adapter wired; emotion mapping done |
| 3 | API Gateway routing + auth; basic Flutter screen |
| 4 | Real‑time preview end‑to‑end; S3 uploads |
| 5 | Quota meter + Razorpay sandbox payments |
| 6 | Closed beta (50 users) – MOS + bug fix sprint |
| 7 | Creator plan billing; polish UX, Hindi UI strings |
| 8 | Play Store (Android) + Web launch |

---

## 7 Risks & Mitigations

| Risk | Impact | Mitigation |
| --- | --- | --- |
| Sarvam rate‑limit / outage | Indic TTS downtime | Cache common phrases; display fallback TTS busy msg |
| GPU cost spike | COGS ↑ | Spot backup node, nightly scale‑to‑zero |
| Celebrity clone abuse | Legal | Consent gate, watermark, ban list |
| High STT error on noisy phones | Poor clone quality | On‑device NR pre‑filter, user tip overlay |

---

## 8 Open Questions

1. Does Bulbul contract allow local caching of audio for 30 days?
2. Should Free tier include STT dictation minutes?
3. Minimum iOS release or Android‑first is okay?

---

**Ready for Dev Kick‑off** – sign‑off pending product + engineering leads.