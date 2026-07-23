# Identity Resolution & Fuzzy Name Matching Scorer

Dokumentasi dan implementasi *Intelligence-Grade Identity Resolution & Fuzzy Name Matching Engine*. Project ini mengimplementasikan pencocokan nama multi-layer (Phonetic, Edit Distance, dan Token Permutation) untuk sistem watchlist, cekal imigrasi, dan verifikasi identitas.

---

## 💡 Rekomendasi Bahasa Pemrograman

| Bahasa | Penggunaan Ideal | Keunggulan |
| :--- | :--- | :--- |
| **Python** *(Rekomendasi Utama)* | Prototyping, Data Science, & Service Umum | Ekosistem perpustakaan pencocokan string/NLP paling lengkap (`RapidFuzz`, `Jellyfish`, `Abydos`, `Metaphone`). |
| **Go (Golang)** | Production Microservice High-Throughput | Performa eksekusi ekstrem, penggunaan memori sangat hemat, konkurensi (goroutines) untuk jutaan *lookup*/detik. |
| **TypeScript / Node.js** | Integration Gateway / Fullstack Web App | Mudah diintegrasikan langsung ke API Gateway atau backend Node.js. |

---

## 🏛️ Arsitektur Hybrid Scoring

Sistem ini tidak mengandalkan 1 algoritma murni, melainkan mengombinasikan 3 pendekatan dengan bobot terkalkulasi:

$$\text{Composite Score} = (w_{\text{phonetic}} \times S_{\text{phonetic}}) + (w_{\text{distance}} \times S_{\text{distance}}) + (w_{\text{token}} \times S_{\text{token}})$$

```
                                +-------------------+
                                | Input Nama Search |
                                +---------+---------+
                                          |
                        +-----------------+-----------------+
                        |                 |                 |
               +--------v-------+  +------v--------+  +-----v---------+
               | DoubleMetaphone|  | Jaro-Winkler  |  |  Token Sort   |
               | (Sound Match)  |  | (Edit Dist)   |  | (Word Order)  |
               +--------+-------+  +------+--------+  +-----+---------+
                        |                 |                 |
                        +-----------------+-----------------+
                                          |
                               +----------v----------+
                               | Composite Score %   |
                               +---------------------+
```

---

## 🛠️ Struktur Project

```text
name-search/
├── README.md             # Dokumentasi Sistem & Artikel
├── python/
│   ├── matcher.py        # Implementasi Engine Python
│   └── requirements.txt  # Dependency Python
└── go/
    ├── matcher.go        # Implementasi Engine Go (High Performance)
    └── go.mod
```

---

## 📊 Threshold & Ambang Batas Decision

| Skor % | Status Action | Tindakan Sistem |
| :--- | :--- | :--- |
| **$\ge 90\%$** | **CRITICAL MATCH** | Hold otomatis / Flag Red Light (Cekal/Watchlist Hit) |
| **$75\% - 89\%$** | **POTENTIAL MATCH** | Antrean peninjauan manual (Secondary Inspection) |
| **$< 75\%$** | **NO MATCH** | Pas lolos verifikasi |

---

## 🚀 Quick Start

### 1. Python Implementation
```bash
cd python
pip install -r requirements.txt
python matcher.py
```

### 2. Go Implementation
```bash
cd go
go run matcher.go
```
