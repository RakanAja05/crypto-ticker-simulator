# Crypto Ticker Simulator

Aplikasi simulasi streaming harga mata uang kripto (BTC, ETH, SOL) secara real-time. Project ini menggunakan **Go WebSocket Server** dengan payload biner **Protocol Buffers (Protobuf)**, serta antarmuka web berbasis **Vue 3 + TypeScript (Vite)**.

---

## Daftar Isi
- [Struktur Project](#struktur-project)
- [Prasyarat Sistem](#prasyarat-sistem)
- [Instalasi Dependencies](#instalasi-dependencies)
- [Regenerasi Protobuf (Opsional)](#regenerasi-protobuf-opsional)
- [Cara Menjalankan Project](#cara-menjalankan-project)
  - [1. Menjalankan Backend Server](#1-menjalankan-backend-server)
  - [2. Menjalankan CLI Client (Terminal)](#2-menjalankan-cli-client-terminal)
  - [3. Menjalankan Frontend Web](#3-menjalankan-frontend-web)
- [Arsitektur & Alur Kerja](#arsitektur--alur-kerja)

---

## Struktur Project

```text
crypto-ticker-simulator/
├── backend/
│   ├── cmd/
│   │   ├── client/main.go      # Go CLI client untuk konsumsi WebSocket
│   │   └── server/main.go      # Go WebSocket server & generator ticker
│   ├── gen/
│   │   └── tickerpb/           # Hasil kompilasi Protobuf untuk Go
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── public/                 # Aset statis web
│   ├── src/
│   │   ├── composables/        # useTickerSocket.ts (WebSocket client logic)
│   │   ├── gen/                # Hasil kompilasi Protobuf untuk TypeScript
│   │   ├── App.vue             # Tampilan dashboard ticker
│   │   └── main.ts             # Entry point Vue 3
│   ├── package.json
│   ├── tsconfig.json
│   └── vite.config.ts
├── proto/
│   └── ticker.proto            # Definisi schema data Protocol Buffers
└── README.md
```

---

## Prasyarat Sistem

Sebelum memulai, pastikan perangkat Anda telah terpasang:
- **Go** (versi 1.22 atau lebih baru)
- **Node.js** (versi 20 atau 22+) dan **npm**
- (Opsional) **protoc** (Protocol Buffer Compiler) jika ingin memodifikasi schema `.proto`

---

## Instalasi Dependencies

### 1. Backend (Go)
Masuk ke direktori `backend` dan unduh semua module Go yang dibutuhkan:
```bash
cd backend
go mod download
go mod tidy
cd ..
```

### 2. Frontend (Vue 3)
Masuk ke direktori `frontend` dan install seluruh dependensi npm:
```bash
cd frontend
npm install
cd ..
```

---

## Regenerasi Protobuf (Opsional)

Jika Anda melakukan perubahan pada schema `proto/ticker.proto`, jalankan perintah berikut dari root project:

1. **Untuk Go (Backend):**
   ```bash
   protoc --go_out=backend/gen --go_opt=module=github.com/RakanAja05/crypto-ticker-simulator/gen proto/ticker.proto
   ```

2. **Untuk TypeScript (Frontend):**
   ```bash
   npx protoc --ts_proto_out=frontend/src/gen --proto_path=proto proto/ticker.proto --ts_proto_opt=outputServices=generic-definitions,useOptionals=messages,env=browser
   ```

*(Catatan: File hasil kompilasi proto sudah tersedia di dalam direktori `backend/gen` dan `frontend/src/gen`, sehingga langkah ini hanya diperlukan jika schema diubah).*

---

## Cara Menjalankan Project

### 1. Menjalankan Backend Server
Server bertindak sebagai WebSocket server yang memancarkan perubahan harga secara berkala (setiap 100ms) dalam bentuk biner Protobuf.

Buka terminal pertama:
```bash
cd backend
go run cmd/server/main.go
```
Output jika berhasil:
```text
listening on :8080
```
Server kini aktif di `ws://localhost:8080/ws`.

---

### 2. Menjalankan CLI Client (Terminal)
Sebagai alternatif atau pengujian backend tanpa web, Anda dapat menjalankan client CLI berbasis Go.

Buka terminal kedua:
```bash
cd backend
go run cmd/client/main.go
```
Output stream di terminal:
```text
BTC 67120.45 (+0.18%)
ETH 3498.20 (-0.05%)
SOL 150.75 (+0.50%)
...
```

---

### 3. Menjalankan Frontend Web
Frontend menampilkan dashboard harga kripto real-time dengan update dinamis.

Buka terminal ketiga:
```bash
cd frontend
npm run dev
```
Setelah Vite berjalan, buka browser pada URL yang ditampilkan (biasanya [http://localhost:5173](http://localhost:5173)). Status pada web akan berubah menjadi **● Connected (Binary)** dan kartu harga akan bergerak secara dinamis.

---

## Arsitektur & Alur Kerja

1. **Generasi Harga**: Server Go menjalankan `time.Ticker` setiap 100ms dan mensimulasikan pergerakan harga (*random walk* ±0.5%).
2. **Serialisasi Biner**: Data harga dibungkus ke dalam pesan Protobuf `Tick` dan diserialisasi ke format biner (`proto.Marshal`).
3. **Transmisi WebSocket**: Pesan dikirimkan melalui WebSocket dengan tipe frame `websocket.BinaryMessage`.
4. **Deserialisasi Client**:
   - **Frontend (Browser)**: Menerima frame biner dengan `binaryType = 'arraybuffer'`, lalu didecode menggunakan fungsi `Tick.decode()` dari library protobufjs/ts-proto.
   - **CLI Client**: Menerima slice byte biner dan didecode menggunakan `proto.Unmarshal()`.
5. **Render Reaktif**: State diperbarui secara reaktif menggunakan Vue 3 `shallowRef` untuk performa tinggi tanpa lag.
