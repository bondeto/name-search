# Identity Resolution & Fuzzy Name Matching Scorer

Project ini mengimplementasikan **Enterprise-Grade Identity Resolution & Screening Engine** dalam bahasa **Go**, **Python**, dan **PostgreSQL SQL**.

Engine ini dirancang untuk pencocokan identitas secara akurat dan cepat terhadap nama, tanggal lahir (DOB), dan Nomor Identitas (ID Number), yang relevan untuk aplikasi publik, FinTech, E-Commerce, Banking AML/KYC, serta pencegahan penipuan (*Anti-Fraud System*).

---

## 🧮 Formulasi Matematika & Algoritma Utama

### 1. Hybrid Composite Name & Identity Score
Skor akhir pencocokan identitas total ($Score_{\text{total}}$) dihitung secara probabilistik berdasarkan akumulasi bobot atribut:

$$S_{\text{name}} = (w_1 \cdot S_{\text{phonetic}}) + (w_2 \cdot S_{\text{jarowinkler}}) + (w_3 \cdot S_{\text{token}})$$

$$Score_{\text{total}} = w_{\text{name}} \cdot S_{\text{name}} + w_{\text{dob}} \cdot S_{\text{dob}} + w_{\text{passport}} \cdot S_{\text{passport}}$$

Di mana $\sum w_i = 1.0$.

---

### 2. Inverse Document Frequency (IDF) Token Entropy Weighting
Menentukan tingkat keunikan setiap kata untuk mencegah kata umum (*MOHAMMAD, SANTO, DE*) memicu *false positive*:

$$\text{Weight}(w) = \log\left(\frac{N + 1}{\text{DocFreq}(w)}\right)$$

$$\text{SoftTFIDF}(s_1, s_2, \theta) = \sum_{w \in s_1 \cap_{\theta} s_2} \text{TFIDF}(w, s_1) \cdot \text{TFIDF}(w', s_2) \cdot \text{Sim}(w, w')$$

---

### 3. Gaussian Date of Birth (DOB) Decay Function
Peluruhan kontinyu selisih hari ($\Delta d = |T_{\text{query}} - T_{\text{candidate}}|$) menggunakan fungsi penalti Gaussian dengan $\sigma = 30\text{ hari}$:

$$S_{\text{dob}} = \exp\left( -\frac{\Delta d^2}{2\sigma^2} \right)$$

---

### 4. Fellegi-Sunter Probabilistic Decision Model (Log-Odds)
Teori pencocokan entitas berbasis akumulasi *log-likelihood ratio*:

$$\text{LogOdds}_{\text{total}} = \sum_{i \in \{\text{Name}, \text{DOB}, \text{Passport}\}} \log_2\left(\frac{m_i}{u_i}\right)$$

Di mana:
* $m_i = P(\text{Atribut } i \text{ cocok} \mid \text{Pasangan Entitas Sama})$
* $u_i = P(\text{Atribut } i \text{ cocok secara kebetulan} \mid \text{Pasangan Entitas Beda})$

---

## 🔒 Fitur Enterprise & Validasi Identitas (Production Ready)

### 1. Parsing & Verifikasi Identitas Nasional
* **Ekstraksi Demografi**: Tanggal lahir, jenis kelamin (pria/wanita), serta kode wilayah dari 16-digit ID.
* **Verifikasi Silang (Cross-Verification)**: Menguji keselarasan tanggal lahir hasil dekoding nomor identitas terhadap tanggal lahir query.

### 2. Cryptographic Audit Trail (Non-Repudiation Logging)
Setiap transaksi *screening* menghasilkan *log audit* terenkripsi SHA-256 yang *tamper-evident*:
* Menjamin ketiadaan penyangkalan (*non-repudiation*).
* Menyimpan hash PII query untuk menjaga kerahasiaan data privasi penggunanya.

### 3. Open API 3.0 Integration Standard
Spesifikasi OpenAPI 3.0 ([`openapi.yaml`](file:///e:/app-bw/name-search/openapi.yaml)) untuk integrasi *interoperabilitas* antar-sistem (*Enterprise Service Bus / API Gateway*).

---

## 🗄️ Implementasi Tingkat Basis Data (Database Level Implementation)

Selain dijalankan di layer *application microservice*, algoritma pencarian nama dan DOB dapat diterapkan langsung di berbagai *Database Management Engine*:

### 1. PostgreSQL (Rekomendasi RDBMS Utama)
* **Ekstensi `pg_trgm`**: Indeks **GIN (Generalized Inverted Index)** `gin_trgm_ops` untuk *candidate blocking* dan `similarity()`.
* **Ekstensi `fuzzystrmatch`**: Fungsi fonetik `dmetaphone()` (Double Metaphone) dan `soundex()`.
* **Stored Function Hybrid Scorer**: Script SQL lengkap tersedia di [`sql/postgres_watchlist_search.sql`](file:///e:/app-bw/name-search/sql/postgres_watchlist_search.sql).

```sql
-- Contoh Query Screening di PostgreSQL:
SELECT * FROM fn_screen_watchlist('MOHAMMED ALY KHAN', '1982-05-14', 0.70);
```

### 2. Elasticsearch / OpenSearch (Rekomendasi Distributed Search)
* **Candidate Blocking**: `n-gram` tokenizer & plugin `analysis-phonetic` (`beider_morse` & `double_metaphone`).
* **Rescoring Engine**: `script_score` (Painless Scripting) untuk kombinasi skor BM25 nama, selisih tanggal lahir, dan *log-odds*.

---

## 📚 Sumber & Referensi Ilmiah

1. **Model Probabilistik Record Linkage (Fellegi-Sunter)**
   * *Referensi*: Fellegi, I. P., & Sunter, A. B. (1969). *A Theory for Record Linkage*. Journal of the American Statistical Association, 64(328), 1183-1210.
2. **Perbandingan String & Edit Distance (Jaro-Winkler)**
   * *Referensi*: Winkler, W. E. (1990). *String Comparator Metrics and Deterministic Decision Rules in the Context of Record Linkage*. U.S. Bureau of the Census.
3. **Pencocokan Fonetik Multibahasa (Beider-Morse & Double Metaphone)**
   * *Referensi*: Beider, A., & Morse, S. P. (2008). *Beider-Morse Phonetic Matching System*. Avotaynu.
4. **Candidate Blocking & Indexing Masif (LSH / MinHash)**
   * *Referensi*: Broder, A. Z. (1997). *On the resemblance and containment of documents*. IEEE Proceedings of Compression and Complexity of Sequences.

---

## 🛠️ Struktur Repository

```text
name-search/
├── README.md                 # Dokumentasi & Formulasi Matematika Lengkap
├── openapi.yaml              # Spesifikasi REST API Enterprise (OpenAPI 3.0)
├── TODO.md                   # Roadmap & Backlog Features
├── .gitignore
├── sql/
│   └── postgres_watchlist_search.sql  # Schema & Function Hybrid Scorer PostgreSQL
├── python/
│   ├── matcher.py            # Basic Name Matcher (Python)
│   ├── expert_engine.py      # Identity Resolution Engine (Python)
│   ├── id_validation.py      # National ID Parser & Audit Logger (Python)
│   ├── test_matcher.py       # Unit Test Suite (Python)
│   └── requirements.txt
└── go/
    ├── matcher.go            # Basic Name Matcher (Go)
    ├── expert_engine.go      # Identity Resolution Engine (Go)
    ├── id_validation.go      # National ID Parser & Audit Logger (Go)
    ├── matcher_test.go       # Unit & Benchmark Test Suite (Go)
    └── go.mod
```

---

## 📄 Lisensi

Distributed under the MIT License.
