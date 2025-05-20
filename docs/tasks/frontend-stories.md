# Vaani – Front‑end MVP Sprint Backlog

**Stack:** Flutter 3.19 + Dart 3.4, `flutter_bloc` 9.x, `hydrated_bloc` *(PWA enabled)*
**Design system:** Material 3 + custom brand swatch
**Release window:** 8 weeks (aligned with backend schedule)
**Cadence:** 1‑week sprints; Kanban in GitHub Projects with labels `fe-epic`, `fe-story`, `fe-task`, `bug`, `ui`, `test`

---

## Epic FE‑0 – Repo & Tooling

| ID         | Story / Task                                                                 | Est | Acceptance Criteria                                                              | Status |
| ---------- | ---------------------------------------------------------------------------- | --- | -------------------------------------------------------------------------------- | ------ |
| **FE‑0.1** | Scaffold Flutter workspace under `frontend/` with `melos`                    | 3 h | `melos bs` completes; `flutter run` boots sample app on Android, iOS sim, Chrome | ✅ Done |
| **FE‑0.2** | Add `makefile` / Taskfile targets `flutter:run`, `flutter:test`, `web:serve` | 2 h | Running `task web:serve` opens hot‑reload dev server                             | ✅ Done |
| **FE‑0.3** | Configure `flutter_intl`, i18n ARB files (en, hi)                            | 2 h | `Intl.message` strings generate; language switch works in debug                  | ✅ Done |
| **FE‑0.4** | Add CI (`ci-flutter.yml`) to run `flutter analyze` + tests                   | 2 h | PR with lint error fails CI                                                      | ✅ Done |

---

## Epic FE‑1 – Core Architecture & State

| ID         | Task                                                                                | Est | Acceptance Criteria                             | Status |
| ---------- | ----------------------------------------------------------------------------------- | --- | ----------------------------------------------- | ------ |
| **FE‑1.1** | Implement folder structure (`presentation/`, `blocs/`, `data/`, `domain/`, `core/`) | 2 h | Analyzer shows zero import‑cycle warnings       | ✅ Done |
| **FE‑1.2** | Setup `get_it` service locator + global `RepositoryProvider`                        | 3 h | `getIt<T>()` retrieves singleton in widget test | ✅ Done |
| **FE‑1.3** | Create base `AppBlocObserver` logging transitions                                   | 1 h | Bloc logs visible in debug console              | ✅ Done |
| **FE‑1.4** | Install `hydrated_bloc` storage (Hive)                                              | 2 h | Toggle dark‑mode persists across restart        | ✅ Done |

---

## Epic FE‑2 – Authentication Flow

| ID         | Task                                                  | Est | Acceptance Criteria                             | Status |
| ---------- | ----------------------------------------------------- | --- | ----------------------------------------------- | ------ |
| **FE‑2.1** | Phone auth UI (country code picker, OTP input)        | 4 h | Accepts 10‑digit IN number, moves to OTP screen | ✅ Done |
| **FE‑2.2** | Implement `AuthBloc` (events: SubmitPhone, VerifyOtp) | 4 h | Success stores JWT; failure shows SnackBar      | ✅ Done |
| **FE‑2.3** | JWT refresh flow with dio `AuthInterceptor`           | 3 h | 401 triggers silent refresh & retry once        | ✅ Done |
| **FE‑2.4** | Unit tests for `AuthBloc` (happy & error)             | 3 h | `bloc_test` coverage ≥ 90 %                     | ✅ Done |

---

## Epic FE‑3 – Home & Project Dashboard

| ID         | Task                                                  | Est | Acceptance Criteria                                  | Status |
| ---------- | ----------------------------------------------------- | --- | ---------------------------------------------------- | ------ |
| **FE‑3.1** | Implement adaptive `HomePage` with list of projects   | 4 h | On <900 px width switches to list‑only mobile layout | ✅ Done |
| **FE‑3.2** | `ProjectBloc` fetches projects from API (`/projects`) | 3 h | Pull‑to‑refresh reloads list                         | ✅ Done |
| **FE‑3.3** | FAB "New Project" navigates to Editor screen          | 2 h | Route push uses `go_router` named route              | ✅ Done |

---

## Epic FE‑4 – Script Editor

| ID         | Task                                                  | Est | Acceptance Criteria                                 | Status |
| ---------- | ----------------------------------------------------- | --- | --------------------------------------------------- | ------ |
| **FE‑4.1** | Rich text field with auto‑punctuation button          | 5 h | Pressing ⚙️ formats text via local NLP util         |  |
| **FE‑4.2** | Voice picker bottom‑sheet (stock + cloned)            | 4 h | Displays avatar, language flag, sample play         |  |
| **FE‑4.3** | Emotion slider widget (0‑100) with live emoji preview | 2 h | Slider emits `EmotionChanged` to bloc               |  |
| **FE‑4.4** | `EditorBloc` manages text, voice, emotion, charCount  | 3 h | Char counter updates realtime, turns red over limit |  |
| **FE‑4.5** | Estimate cost banner (chars × ₹/char from `/pricing`) | 2 h | Shows ₹ value ≤ 0.01 precision                      |  |

---

## Epic FE‑5 – Preview Player (WebSocket)

| ID         | Task                                                      | Est | Acceptance Criteria                                   | Status |
| ---------- | --------------------------------------------------------- | --- | ----------------------------------------------------- | ------ |
| **FE‑5.1** | Connect WebSocket `/v1/preview` with `web_socket_channel` | 3 h | Sends JSON payload, receives first frame under 800 ms |  |
| **FE‑5.2** | PCM stream decoder to `AudioSource` (just\_audio)         | 4 h | Audio plays smoothly with <100 ms jitter              |  |
| **FE‑5.3** | Loading shimmer & error state views                       | 2 h | Network drop shows retry CTA                          |  |
| **FE‑5.4** | Preview toolbar: play/stop, regenerate                    | 3 h | Regenerate cancels prev WS & starts new               |  |

---

## Epic FE‑6 – Lite Mixer & Timeline

| ID         | Task                                              | Est | Acceptance Criteria                       |
| ---------- | ------------------------------------------------- | --- | ----------------------------------------- |
| **FE‑6.1** | Implement `ReorderableListView` of voice segments | 4 h | Long‑press drag reorders segments         |
| **FE‑6.2** | Add background‑music track (pick local file)      | 3 h | Volume slider 0–100 %                     |
| **FE‑6.3** | Waveform thumbnail using `audio_waveforms` pkg    | 4 h | Renders within 200 ms for 30 s clip       |
| **FE‑6.4** | `MixerBloc` manages tracks, offsets, durations    | 4 h | Dropping new clip updates timeline length |

---

## Epic FE‑7 – Export & Share

| ID         | Task                                   | Est | Acceptance Criteria                       |
| ---------- | -------------------------------------- | --- | ----------------------------------------- |
| **FE‑7.1** | POST `/generate` + progress overlay    | 3 h | Progress bar updates via SSE (job status) |
| **FE‑7.2** | Download MP3 to device using `dio`     | 2 h | File visible in Files app (Android 13)    |
| **FE‑7.3** | Share‑sheet integration (`share_plus`) | 1 h | Tap share opens native dialog             |
| **FE‑7.4** | Success Lottie + CTA to Rate App       | 1 h | Google in‑app review triggered            |

---

## Epic FE‑8 – Billing & Quota UI

| ID         | Task                                   | Est | Acceptance Criteria                              |
| ---------- | -------------------------------------- | --- | ------------------------------------------------ |
| **FE‑8.1** | Quota indicator chip in AppBar         | 2 h | Shows remaining minutes/day                      |
| **FE‑8.2** | Razorpay WebView checkout screen       | 4 h | Successful payment triggers `CreditsAdded` event |
| **FE‑8.3** | "Weekend Pass" banner when quota <10 % | 1 h | Deep‑links to one‑time pay flow                  |

---

## Epic FE‑9 – Theming & Localization

| ID         | Task                                               | Est | Acceptance Criteria                       |
| ---------- | -------------------------------------------------- | --- | ----------------------------------------- |
| **FE‑9.1** | Define brand `ColorScheme` seed & Material 3 tones | 2 h | Light/dark toggles with animation         |
| **FE‑9.2** | Hindi translations for all strings                 | 2 h | Setting language reloads app (Bloc event) |
| **FE‑9.3** | Accessibility: textScaleFactor test to 2.0         | 2 h | No overflow errors in golden run          |

---

## Epic FE‑10 – Testing & QA

| ID          | Task                                            | Est | Acceptance Criteria                                     |
| ----------- | ----------------------------------------------- | --- | ------------------------------------------------------- |
| **FE‑10.1** | Bloc unit tests for Editor, Preview, Mixer      | 6 h | ≥90 % coverage each bloc                                |
| **FE‑10.2** | Golden tests for Key screens (phone & web‑wide) | 4 h | Goldens committed; diff fails on change                 |
| **FE‑10.3** | Integration test: auth → generate → export path | 6 h | Runs in GitHub CI using `flutter test integration_test` |

---

## Epic FE‑11 – Performance & PWA

| ID          | Task                                             | Est | Acceptance Criteria                 |
| ----------- | ------------------------------------------------ | --- | ----------------------------------- |
| **FE‑11.1** | Optimise asset pre‑cache & deferred font loading | 2 h | 20 % faster cold‑start per DevTools |
| **FE‑11.2** | Add `flutter_service_worker.js` cache strategy   | 3 h | PWA offline shows cached HomePage   |
| **FE‑11.3** | Lighthouse audit ≥ 90 % perf on web build        | 2 h | Scores screenshot saved to docs     |

---

## Total Estimate

* **≈ 128 dev hours** → \~16 work days + 30 % buffer ⇒ \~21 days (3 × 1‑week sprints + hardening).

---

## Sprint Allocation (suggestion)

| Sprint | Focus                                                  |
| ------ | ------------------------------------------------------ |
| 0      | Epic FE‑0, FE‑1 setup                                  |
| 1      | Auth (FE‑2) + Dash (FE‑3)                              |
| 2      | Editor & Preview (FE‑4, FE‑5)                          |
| 3      | Mixer & Export (FE‑6, FE‑7)                            |
| 4      | Billing, Theming, QA polish (FE‑8, FE‑9, FE‑10, FE‑11) |

---

**End of Front‑end MVP Sprint Backlog**
