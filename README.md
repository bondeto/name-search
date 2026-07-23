# Identity Resolution & Fuzzy Name Matching Scorer (Government Ready)

Project ini mengimplementasikan **Government-Grade Identity Resolution & Watchlist Screening Engine** dalam bahasa **Go** dan **Python**.

Engine ini dirancang untuk pencocokan identitas secara akurat dan cepat terhadap nama, tanggal lahir (DOB), dan Nomor Induk Kependudukan (NIK/KTP), berstandar instansi pemerintah (Imigrasi Cekal, Polri, PPATK, Bea Cukai).

---

## 🏛️ Fitur Kesiapan Instansi Pemerintah (Government-Readiness)

### 1. Parsing & Verifikasi Struktural NIK KTP Indonesia
Sistem dilengkapi modul otomatis untuk mengekstrak dan memverifikasi integritas 16 digit NIK KTP:
* **Ekstraksi Demografi**: Tanggal lahir, jenis kelamin (pria/wanita), serta kode wilayah (provinsi, kabupaten/kota, kecamatan).
* **Verifikasi Silang (Cross-Verification)**: Menguji keselarasan tanggal lahir hasil dekoding NIK terhadap tanggal lahir query.

### 2. ISO 27001 Cryptographic Audit Trail (Non-Repudiation Logging)
Setiap transaksi *screening* menghasilkan *log audit* terenkripsi SHA-256 yang *tamper-evident*:
* Menjamin ketiadaan penyangkalan (*non-repudiation*).
* Menyimpan hash PII query untuk menjaga kerahasiaan data privasi sesuai regulasi PDP (Pelindungan Data Pribadi).

### 3. Open API 3.0 Integration Standard
Dilengkapi spesifikasi OpenAPI 3.0 ([`openapi.yaml`](file:///e:/app-bw/name-search/openapi.yaml)) untuk mempermudah integrasi *interoperabilitas* antar-instansi (*Enterprise Service Bus / API Gateway*).

---

## 📚 Sumber & Referensi Ilmiah / Standar Industri

1. **Model Probabilistik Record Linkage (Fellegi-Sunter)**
   * *Referensi*: Fellegi, I. P., & Sunter, A. B. (1969). *A Theory for Record Linkage*. Journal of the American Statistical Association, 64(328), 1183-1210.
2. **Perbandingan String & Edit Distance (Jaro-Winkler)**
   * *Referensi*: Winkler, W. E. (1990). *String Comparator Metrics and Deterministic Decision Rules in the Context of Record Linkage*. U.S. Bureau of the Census.
3. **Pencocokan Fonetik Multibahasa (Beider-Morse & Double Metaphone)**
   * *Referensi*: Beider, A., & Morse, S. P. (2008). *Beider-Morse Phonetic Matching System*. Avotaynu.
4. **Candidate Blocking & Indexing Masif (LSH / MinHash)**
   * *Referensi*: Broder, A. Z. (1997). *On the resemblance and containment of documents*. IEEE Proceedings of Compression and Complexity of Sequences.
5. **Standar Watchlist & Sanction Screening (OFAC & DHS)**
   * *Referensi*: U.S. Department of the Treasury - Office of Foreign Assets Control (OFAC). *Sanctions List Search Engine Architecture Guidelines*.
   * *Referensi*: U.S. Department of Homeland Security (DHS). *IDENT/HART Identity Resolution Technical Framework*.

---

## 🛠️ Struktur Repository

```text
name-search/
├── README.md                 # Dokumentasi & Sumber Referensi Ilmiah
├── openapi.yaml              # Spesifikasi REST API Enterprise (OpenAPI 3.0)
├── TODO.md                   # Roadmap & Backlog Features
├── .gitignore
├── python/
│   ├── matcher.py            # Basic Name Matcher (Python)
│   ├── expert_engine.py      # Identity Resolution Engine (Python)
│   ├── gov_readiness.py      # KTP NIK Parser & ISO 27001 Audit Logger (Python)
│   ├── test_matcher.py       # Unit Test Suite (Python)
│   └── requirements.txt
└── go/
    ├── matcher.go            # Basic Name Matcher (Go)
    ├── expert_engine.go      # Identity Resolution Engine (Go)
    ├── gov_readiness.go      # KTP NIK Parser & ISO 27001 Audit Logger (Go)
    ├── matcher_test.go       # Unit & Benchmark Test Suite (Go)
    └── go.mod
```

---

## 🚀 Pengujian & Penggunaan Modul Government

### 1. Engine Go (KTP Parser & Audit Trail)
```bash
cd go
go run gov_readiness.go
```

### 2. Engine Python (KTP Parser & Audit Trail)
```bash
cd python
python gov_readiness.py
```

---

## 📄 Lisensi

Distributed under the MIT License.
