# ECB Fyne Fullstack

Proyek demo fullstack (Go + Fyne + Python sensor) yang bisa dijalankan di Windows untuk keperluan demo dan mudah dipindahkan ke Raspberry Pi 3 untuk implementasi lapangan.

## Struktur Folder

```
ecb-fyne-fullstack/
├── backend/          # REST API (Gin + GORM)
├── frontend/         # Aplikasi GUI (Fyne)
├── sensor/           # Skrip pembacaan sensor (Python)
└── Makefile          # Shortcut perintah lintas OS
```

## 1. Persiapan Umum

- Go 1.21+ (`go version`)
- MySQL Server 8 / MariaDB 10.6+
- Python 3.8+
- Git (opsional untuk clone repo)

> **Konfigurasi Database**: Salin `backend/.env.example` menjadi `backend/.env` dan sesuaikan kredensial database jika perlu.
>
> ```ini
> API_PORT=8080
> DB_HOST=127.0.0.1
> DB_PORT=3306
> DB_USER=root
> DB_PASSWORD=your_password
> DB_NAME=ecb_test
> ```

## 2. Setup dan Demo di Windows

1. Pastikan Go, Python, dan MySQL sudah terpasang. Buat database `ecb_test`.
2. Jalankan instalasi dependencies:
   ```powershell
   cd ecb-fyne-fullstack
   make install
   ```
3. Jalankan backend:
   ```powershell
   make run-backend
   ```
4. Jalankan GUI (terminal baru):
   ```powershell
   make run-frontend
   ```
5. Jalankan sensor simulator (terminal baru):
   ```powershell
   make run-sensor
   ```
   Mode ini otomatis memakai data simulasi (`USE_SIMULATED_SENSOR=1`). Data dummy akan dikirim ke backend, dan GUI menampilkan status terbaru.

## 3. Setup di Raspberry Pi 3 (Raspbian OS)

1. Update sistem & install tool chain:
   ```bash
   sudo apt update && sudo apt upgrade -y
   sudo apt install -y golang python3 python3-pip git mysql-server gcc libgl1-mesa-dev xorg-dev
   ```
2. Clone/copy proyek dan masuk ke folder `ecb-fyne-fullstack`.
3. Instal dependencies:
   ```bash
   make install-pi
   ```
   Target ini memakai `sensor/requirements.pi.txt` sehingga paket `RPi.GPIO` ikut terpasang.
4. Konfigurasi database MySQL lalu salin `.env`:
   ```bash
   cp backend/.env.example backend/.env
   nano backend/.env   # sesuaikan DB_PASSWORD jika perlu
   ```
5. Jalankan service:
   ```bash
   make run-backend
   ```
   Terminal lain:
   ```bash
   make run-frontend
   ```
   Terminal lain (GPIO mode aktif secara default):
   ```bash
   make run-sensor-pi
   ```
   Jika ingin tetap menggunakan simulasi di Pi (misal belum pasang sensor) jalankan:
   ```bash
   USE_SIMULATED_SENSOR=1 make run-sensor-pi
   ```

## 4. Build Biner

- **Desktop/Windows**: `make build-all` menghasilkan `bin/ecb-backend.exe` dan `bin/ecb-frontend.exe`.
- **Raspberry Pi**: `make build-pi` menghasilkan biner ARM di `bin/`.

## 5. Variabel Lingkungan Penting

| Variabel            | Default                   | Keterangan                                      |
|---------------------|---------------------------|-------------------------------------------------|
| `API_PORT`          | `8080`                    | Port backend (set di `.env`)                    |
| `API_BASE_URL`      | `http://127.0.0.1:8080/api` | Dipakai frontend & sensor                       |
| `USE_SIMULATED_SENSOR` | `1` di Windows (`make run-sensor`) | Set `0` di Pi untuk aktifkan GPIO              |
| `POST_INTERVAL_SECONDS` | `3`                   | Jeda pengiriman data sensor                     |

## 6. Tips Migrasi Windows → Raspberry Pi

1. Kembangkan dan uji alur lengkap di Windows dengan mode simulasi.
2. Salin folder `bin/` hasil `make build-pi` atau deploy source code lengkap ke Pi.
3. Pastikan variabel lingkungan (`API_BASE_URL`, credential DB) sudah sesuai dengan jaringan produksi Pi.
4. Gunakan `systemd` service pada Pi jika ingin menjalankan backend/front-end/sensor sebagai service permanen (lihat contoh di `docs/systemd/` jika dibuat di kemudian hari).

---

Siap untuk demo di Windows dan siap dideploy ke Raspberry Pi 3 tanpa mengganti kode inti. Selamat mencoba! 🎉
