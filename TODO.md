# 📌 Open Backlog & Unfinished Job Flags (Roadmap)

Dokumen ini mencatat daftar *Unfinished Jobs*, *Feature Flags*, dan rencana pengembangan lanjutan (*Roadmap*) untuk `name-search`.

---

## 🚩 Pending / Unfinished Features Flag

### 1. High-Priority (Phase 2)
- [ ] **[TODO-001] Distributed LSH Indexing via Redis Cluster**:
  - *Status*: PENDING
  - *Description*: Mengubah *in-memory inverted index* lokal menjadi *Redis Search / Vector Store Cluster* untuk mendukung scalability 100M+ identitas.
- [ ] **[TODO-002] Live OFAC & Interpol API Ingestion Pipeline**:
  - *Status*: PENDING
  - *Description*: Automasi penarikan feed data *Sanctions List* (OFAC XML/CSV & Interpol API) secara otomatis per jam via Cron / Background Worker.
- [ ] **[TODO-003] Neural Multi-Lingual Transliteration (Transformer-based)**:
  - *Status*: PENDING
  - *Description*: Mengganti *rule-based Double Metaphone* dengan model transformer transliterasi nama (Arabic/Cyrillic/Hanzi $\rightarrow$ Latin).

### 2. Medium-Priority (Phase 3)
- [ ] **[TODO-004] gRPC & REST API Gateway Wrapper**:
  - *Status*: PENDING
  - *Description*: Membungkus engine Go ke dalam server gRPC (`:50051`) dan REST HTTP (`:8080`) dengan middleware Prometheus metrics.
- [ ] **[TODO-005] Biometric Vector Embedding Integration**:
  - *Status*: PENDING
  - *Description*: Mengintegrasikan skor kemiripan vektor wajah (DeepFace / ArcFace 512-dim embedding) ke dalam kalkulasi *Fellegi-Sunter*.

---

## 🏷️ Code-Level TODO Flags

* **Go Engine** (`go/expert_engine.go`):
  * `// TODO(perf): Implement SIMD-accelerated Jaro-Winkler distance calculation.`
  * `// TODO(scalability): Migrate map[string]set to distributed Redis bitsets for multi-node deployments.`
* **Python Engine** (`python/expert_engine.py`):
  * `# TODO(ml): Integrate Sentence-Transformers (all-MiniLM-L6-v2) for semantic name embedding comparison.`
