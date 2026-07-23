# Identity Resolution & Fuzzy Name Matching Scorer (Expert Level)

Project ini mengimplementasikan **Intelligence-Grade Identity Resolution & Watchlist Screening Engine** (berstandar FBI NCIC, CIA HART, INTERPOL, & OFAC Compliance) dalam bahasa **Go** dan **Python**.

Engine ini dirancang untuk melakukan pencocokan identitas secara akurat dan cepat terhadap nama dan tanggal lahir (DOB), mengatasi kesalahan ejaan (*typo*), variasi transliterasi multibahasa, pembalikan urutan kata, kesalahan format tanggal, dan pengunaan nama alias.

---

## 🏛️ Komponen Expert Level yang Diimplementasikan

### 1. IDF Token Entropy Weighting (Pembobotan Entropi Kata)
Pencocokan nama tradisional sering gagal karena memberikan bobot sama rata pada setiap kata. Engine ini menghitung **Inverse Document Frequency (IDF)** untuk menentukan tingkat keunikan setiap token nama:
* Kata umum seperti `"MOHAMMAD"`, `"SANTO"`, `"AL"`, `"BIN"`, atau `"DE"` diberi bobot rendah sehingga tidak memicu *false positive*.
* Kata langka seperti `"ZULQARNAIN"` atau `"KAZIMIERZ"` diberi bobot tinggi yang mendominasi kalkulasi skor.

$$\text{Weight}(w) = \log\left(\frac{N + 1}{\text{DocFreq}(w)}\right)$$

### 2. Fuzzy Date of Birth (DOB) Decay Engine
Sistem tidak menggunakan pencocokan tanggal lahir secara murni *exact match*. Engine menangani variasi kesalahan tanggal lahir dengan:
* **Transposition & Month/Day Swap Detection**: Mengenali kesalahan input tukar bulan/tanggal (misal `1985-05-12` vs `1985-12-05`) dan memberikan penalti minimal.
* **Gaussian Decay Function**: Menghitung peluruhan kontinyu selisih hari ($\Delta d$) dengan $\sigma = 30\text{ hari}$:
  $$S_{\text{dob}} = \exp\left( -\frac{\Delta d^2}{2\sigma^2} \right)$$
* **Wildcard & Imputation Handling**: Mendukung format tanggal parsial seperti `1985-00-00` atau `1985-XX-XX`.

### 3. Probabilistic Fellegi-Sunter Decision Engine
Implementasi teori statistik **Fellegi-Sunter Methodology** untuk *record linkage*. Akumulasi nilai *log-odds* dihitung dari gabungan bobot seluruh atribut:

$$\text{LogOdds}_{\text{total}} = \sum_{i \in \{\text{Name}, \text{DOB}, \text{Passport}\}} \log_2\left(\frac{m_i}{u_i}\right)$$

* Menghasilkan keputusan deterministik dengan bobot kepercayaan (*confidence weights*).

### 4. Sub-Millisecond Candidate Blocking Engine (LSH / Inverted Index)
Untuk menangani basis data watchlist skala masif (jutaan *record*), pencarian *brute-force* dihindari:
* Menggunakan **Inverted Soundex & N-Gram Bucket Indexing**.
* Memangkas 10+ juta data menjadi $< 50$ kandidat paling relevan dalam waktu **$< 2\text{ milidetik}$**, sebelum dijalankan *multi-attribute rescoring*.

### 5. Onomastic Normalization & Prefix Stripping
Sistem linguistik *onomastik* untuk memotong prefiks kebudayaan/gelar secara otomatis sebelum pemprosesan fonetik:
* Menghapus prefiks seperti `AL-`, `EL-`, `ABDUL-`, `BIN`, `BINTI`, `VON`, `VAN`, `DE`, `DER`, `SAN`.
* Contoh: `"AL-RAHMAN"` dan `"ABDUL RAHMAN"` direduksi ke akar fonetik utama yang sama: `[RHMN]`.

### 6. Graph-Based Alias Network (Entity Resolution Graph)
Relasi nama alias (*Also Known As / AKA*) dipetakan ke dalam struktur *Graph Network*:
* Pencarian terhadap suatu nama secara otomatis menarik dan membandingkan node alias yang terhubung dalam satu kluster entitas.

---

## 📐 Arsitektur Sistem

```text
                           +------------------------+
                           |   Input Query Identitas |
                           +-----------+------------+
                                       |
                           +-----------v------------+
                           |  Candidate Blocking    |  (Inverted Soundex / LSH Index)
                           |  (Memangkas 10M -> 50) |  < 2 ms
                           +-----------+------------+
                                       |
                 +---------------------+---------------------+
                 |                     |                     |
        +--------v-------+    +--------v-------+    +--------v-------+
        | IDF Entropy    |    | Gaussian DOB   |    | Passport Match |
        | Name Matcher   |    | Decay Scorer   |    | Verification   |
        +--------+-------+    +--------+-------+    +--------+-------+
                 |                     |                     |
                 +---------------------+---------------------+
                                       |
                           +-----------v------------+
                           | Fellegi-Sunter Model   |
                           |  Log-Odds Weighting    |
                           +-----------+------------+
                                       |
                           +-----------v------------+
                           | Composite Score (%) &  |
                           |   Threshold Decision   |
                           +------------------------+
```

---

## 📊 Threshold & Ambang Batas Decision

| Composite Score % | Status Signal | Tindakan Operasional |
| :--- | :--- | :--- |
| **$\ge 88\%$** *(atau Paspor Match)* | **CRITICAL MATCH (RED)** | Hold / Block Otomatis / Red Flag Watchlist |
| **$70\% - 87\%$** | **POTENTIAL MATCH (YELLOW)** | Escalate ke Secondary Analyst Review |
| **$< 70\%$** | **NO MATCH (CLEAR)** | Pas Lolos Otomatis |

---

## 🛠️ Struktur Repository

```text
name-search/
├── README.md                 # Dokumentasi Sistem & Blueprint Arsitektur
├── .gitignore
├── python/
│   ├── matcher.py            # Basic Name Matcher (Python)
│   ├── expert_engine.py      # Expert-Level Identity Resolution Engine (Python)
│   ├── test_matcher.py       # Unit Test Suite (Python)
│   └── requirements.txt
└── go/
    ├── matcher.go            # Basic Name Matcher (Go)
    ├── expert_engine.go      # Expert-Level Identity Resolution Engine (Go)
    ├── matcher_test.go       # Unit & Benchmark Test Suite (Go)
    └── go.mod
```

---

## 🚀 Pengujian & Penggunaan

### 1. Engine Go (High-Performance Production)
```bash
cd go
# Menjalankan Expert Engine Demo
go run .

# Menjalankan Unit & Benchmark Testing
go test -v -bench=.
```

### 2. Engine Python (Riset & Analytics)
```bash
cd python
# Menjalankan Expert Engine Demo
python expert_engine.py

# Menjalankan Unit Testing
python -m unittest test_matcher.py
```

---

## 📄 Lisensi

Distributed under the MIT License.
