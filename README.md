# Identity Resolution & Fuzzy Name Matching Scorer

Project ini mengimplementasikan **Identity Resolution & Watchlist Screening Engine** dalam bahasa **Go** dan **Python**.

Engine ini dirancang untuk melakukan pencocokan identitas secara akurat dan cepat terhadap nama dan tanggal lahir (DOB), mengatasi kesalahan ejaan (*typo*), variasi transliterasi multibahasa, pembalikan urutan kata, kesalahan format tanggal, dan penggunaan nama alias.

---

## 📚 Sumber & Referensi Ilmiah / Standar Industri

Arsitektur dan algoritma yang diimplementasikan pada project ini dibangun berdasarkan standar industri dan literatur ilmiah berikut:

1. **Model Probabilistik Record Linkage (Fellegi-Sunter)**
   * *Referensi*: Fellegi, I. P., & Sunter, A. B. (1969). *A Theory for Record Linkage*. Journal of the American Statistical Association, 64(328), 1183-1210.
   * *Penerapan*: Akumulasi bobot *log-odds* antar-atribut (Nama, DOB, Paspor) untuk menentukan probabilitas pencocokan entitas.

2. **Perbandingan String & Edit Distance (Jaro-Winkler)**
   * *Referensi*: Winkler, W. E. (1990). *String Comparator Metrics and Deterministic Decision Rules in the Context of Record Linkage*. U.S. Bureau of the Census.
   * *Penerapan*: Pengukuran kemiripan struktur nama dengan pembobotan prefiks.

3. **Pencocokan Fonetik Multibahasa (Beider-Morse & Double Metaphone)**
   * *Referensi*: Beider, A., & Morse, S. P. (2008). *Beider-Morse Phonetic Matching System*. Avotaynu.
   * *Penerapan*: Pencocokan fonetik lintas etnis/bahasa (Arab, Slavia, Eropa, Asia) dan normalisasi prefiks onomastik.

4. **Candidate Blocking & Indexing Masif (LSH / MinHash)**
   * *Referensi*: Broder, A. Z. (1997). *On the resemblance and containment of documents*. IEEE Proceedings of Compression and Complexity of Sequences.
   * *Penerapan*: Pemangkasan pencarian kandidat dari puluhan juta data menjadi $< 50$ kandidat dalam latensi sub-milidetik ($< 2\text{ ms}$).

5. **Standar Watchlist & Sanction Screening (OFAC & DHS)**
   * *Referensi*: U.S. Department of the Treasury - Office of Foreign Assets Control (OFAC). *Sanctions List Search Engine Architecture & Name Matching Guidelines*.
   * *Referensi*: U.S. Department of Homeland Security (DHS). *IDENT/HART Identity Resolution Technical Framework*.

---

## 🏛️ Komponen Algoritma yang Diimplementasikan

### 1. IDF Token Entropy Weighting
Menghitung **Inverse Document Frequency (IDF)** untuk menentukan tingkat keunikan setiap token nama:
$$\text{Weight}(w) = \log\left(\frac{N + 1}{\text{DocFreq}(w)}\right)$$

### 2. Fuzzy Date of Birth (DOB) Decay Engine
* **Transposition & Month/Day Swap Detection**: Deteksi kesalahan input tukar bulan/tanggal (misal `1985-05-12` vs `1985-12-05`).
* **Gaussian Decay Function**: Peluruhan kontinyu selisih hari ($\Delta d$) dengan $\sigma = 30\text{ hari}$:
  $$S_{\text{dob}} = \exp\left( -\frac{\Delta d^2}{2\sigma^2} \right)$$

### 3. Sub-Millisecond Candidate Blocking Engine
Menggunakan **Inverted Soundex & N-Gram Bucket Indexing** untuk memangkas basis data skala masif sebelum *multi-attribute rescoring*.

### 4. Onomastic Normalization & Prefix Stripping
Menghapus prefiks gelar/kebudayaan (`AL-`, `EL-`, `ABDUL-`, `BIN`, `BINTI`, `VON`, `VAN`, `DE`, `DER`, `SAN`).

### 5. Graph-Based Alias Network
Pemetaan relasi nama alias (*Also Known As / AKA*) dalam struktur *Graph Network*.

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
├── README.md                 # Dokumentasi & Sumber Referensi Ilmiah
├── TODO.md                   # Roadmap & Backlog Features
├── .gitignore
├── python/
│   ├── matcher.py            # Basic Name Matcher (Python)
│   ├── expert_engine.py      # Identity Resolution Engine (Python)
│   ├── test_matcher.py       # Unit Test Suite (Python)
│   └── requirements.txt
└── go/
    ├── matcher.go            # Basic Name Matcher (Go)
    ├── expert_engine.go      # Identity Resolution Engine (Go)
    ├── matcher_test.go       # Unit & Benchmark Test Suite (Go)
    └── go.mod
```

---

## 🚀 Pengujian & Penggunaan

### 1. Engine Go (High-Performance Production)
```bash
cd go
# Menjalankan Engine Demo
go run .

# Menjalankan Unit & Benchmark Testing
go test -v -bench=.
```

### 2. Engine Python (Riset & Analytics)
```bash
cd python
# Menjalankan Engine Demo
python expert_engine.py

# Menjalankan Unit Testing
python -m unittest test_matcher.py
```

---

## 📄 Lisensi

Distributed under the MIT License.
