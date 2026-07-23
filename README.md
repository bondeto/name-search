# Identity Resolution & Fuzzy Name Matching Scorer (Expert Level)

Project ini berisi implementasi **Intelligence-Grade Identity Resolution Engine** (setara sistem FBI NCIC, CIA HART, INTERPOL, & OFAC Sanctions Screening) untuk pencocokan nama dan tanggal lahir (DOB).

---

## 🔬 Fitur Utama Engine (Expert Level)

1. **IDF Token Entropy Weighting**
   - Menghitung bobot keunikan kata (*Inverse Document Frequency*). Kata umum seperti `"MOHAMMAD"`, `"AL"`, atau `"DE"` diberi penalti bobot rendah, sedangkan kata langka seperti `"ZULQARNAIN"` mendominasi kalkulasi skor.
2. **Fuzzy Date of Birth (DOB) Decay Engine**
   - Deteksi *Month/Day Swapping* (`1985-05-12` vs `1985-12-05`) dan kesalahan ketik tanggal.
   - Fungsi peluruhan kontinu **Gaussian Decay** ($\sigma = 30\text{ hari}$) untuk mengukur deviasi selisih hari.
   - Dukungan pencocokan parsial (*Wildcard DOB* `1985-00-00`).
3. **Fellegi-Sunter Probabilistic Decision Engine**
   - Mengakumulasi nilai *log-odds* antar-atribut (Nama, DOB, Paspor, Alias) berdasarkan teori statistik Fellegi-Sunter untuk *record linkage*.
4. **Onomastic Prefix Stripping & Stemming**
   - Pembersihan otomatis prefiks kebudayaan (*Al-*, *El-*, *Abdul-*, *Bin*, *Binti*, *Von*, *De*).
5. **Sub-Millisecond Candidate Blocking (LSH / Inverted Index)**
   - Menggunakan *Inverted Soundex/N-gram Bucket Indexing* untuk memangkas pencarian kandidat dari jutaan data menjadi $< 50$ data kandidat dalam waktu $< 2\text{ ms}$.
6. **Alias Graph Network**
   - Pemetaan relasi node alias (*Also Known As / AKA*).

---

## 🛠️ Struktur Project

```text
name-search/
├── README.md                 # Blueprint & Dokumentasi Sistem
├── .gitignore
├── python/
│   ├── matcher.py            # Basic Engine (Python)
│   ├── expert_engine.py      # Expert-Level Engine (Python)
│   ├── test_matcher.py       # Unit Tests
│   └── requirements.txt
└── go/
    ├── matcher.go            # Basic Engine (Go)
    ├── expert_engine.go      # Expert-Level Engine (Go)
    ├── matcher_test.go       # Unit & Benchmark Tests
    └── go.mod
```

---

## 📊 Threshold & Ambang Batas Decision

| Composite Score % | Status Signal | Tindakan Operasional |
| :--- | :--- | :--- |
| **$\ge 88\%$** *(atau Paspor Match)* | **CRITICAL MATCH (RED)** | Hold / Block Otomatis / Red Flag Watchlist |
| **$70\% - 87\%$** | **POTENTIAL MATCH (YELLOW)** | Escalate ke Secondary Analyst Review |
| **$< 70\%$** | **NO MATCH (CLEAR)** | Pas Lolos Otomatis |

---

## 🚀 Quick Start

### Run Python Expert Engine:
```bash
cd python
python expert_engine.py
```

### Run Go Expert Engine:
```bash
cd go
go run .
```
