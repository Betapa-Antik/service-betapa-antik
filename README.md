# service-betapa-antik

Service backend untuk proyek Aplikasi-Jentik — modul `betapa-antik-service`.

## Ringkasan

README ini berisi langkah cepat untuk menginstal dependensi, membangun, dan menjalankan service pada lingkungan pengembangan (Windows / PowerShell maupun Linux/macOS).

## Prasyarat

- Go >= 1.24 (toolchain yang tercantum: `go1.24.11`).
- Database dan/atau broker pesan (mis. RabbitMQ) jika service membutuhkan koneksi eksternal — konfigurasi ada di folder `configs/`.

## Instalasi & Dependensi

1. Clone repo dan masuk ke direktori modul:

```bash
git clone <repo-url>
cd service-betapa-antik
```

2. Unduh dan sinkronkan dependensi modul:

```bash
go mod download   # (opsional) unduh dependensi sesuai go.sum
go mod tidy       # tambahkan dependensi yang diperlukan dan perbarui go.sum
```

3. Menambahkan dependensi spesifik (jika perlu):

```bash
go get github.com/joho/godotenv@latest  # contoh: menambah godotenv
```

4. Periksa daftar modul terpasang:

```bash
go list -m all
go mod graph
```

Catatan: saat ini `go.mod` hanya berisi deklarasi modul dan versi Go. Jalankan `go mod tidy` untuk mengisi `go.sum` berdasarkan kode aktual.

## Build & Run

- Build seluruh paket:

```bash
go build ./...
```

- Jalankan langsung (development):

```bash
go run .
```

- Atau jalankan binary yang dihasilkan (Windows):

```powershell
.\\service.exe
```

## Konfigurasi lingkungan

- Periksa folder `configs/` untuk contoh file konfigurasi.
- Jika proyek menggunakan variabel lingkungan, buat file `.env` atau atur env vars sebelum menjalankan. Contoh variabel umum: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASS`, `RABBITMQ_URL`.

## Testing

Jalankan unit test (jika ada):

```bash
go test ./...
```

## Troubleshooting singkat

- Error dependency: jalankan `go mod tidy` lalu `go mod download`.
- Ingin memperbarui semua modul: `go get -u ./...` lalu `go mod tidy`.

---

Jika mau, saya bisa:

- Menjalankan `go mod tidy` sekarang untuk mengisi `go.sum` dan kemudian mengekstrak daftar dependensi.
- Membuat contoh `.env` berdasarkan `configs/`.
