# Go E-Commerce API

## Prequisite

- Go 1.26 atau lebih baru
- PostgreSQL dan akses ke database
- `git`

Opsional:

- `sqlc` jika ingin menghasilkan ulang kode query SQL dari file `queries.sql`
- `goose` untuk migrasi tabel ke database

## Struktur Utama

- `cmd/main.go` - entrypoint server
- `internal/env/env.go` - helper untuk membaca variabel environment
- `internal/products/handlers.go` - HTTP handler untuk endpoint produk
- `internal/products/service.go` - service untuk logika produk
- `internal/adapters/postgresql/migrations/00002_create_products.sql` - skema tabel produk
- `internal/adapters/postgresql/sqlc/` - hasil generate SQLC untuk query database

## Langkah Persiapan

1. Clone repository:

```bash
git clone https://github.com/kodoktempur666/go-ecom-api.git
cd go-ecom-api
```

2. Install dependency Go:

```bash
go mod download
```

3. Siapkan database PostgreSQL.

Database default yang dipakai oleh aplikasi adalah:

```text
host=localhost user=postgres password=psql dbname=mydb sslmode=disable
```

Jika belum ada, buat database `mydb` dan user yang sesuai:

```bash
psql -U postgres
CREATE DATABASE mydb;
\q
```

4. Jalankan migrasi tabel produk.

Anda bisa mengimpor SQL secara manual:

```bash
psql "host=localhost user=postgres password=psql dbname=mydb sslmode=disable" -f internal/adapters/postgresql/migrations/00002_create_products.sql
```

Atau gunakan `goose` untuk mengelola migrasi:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
export GOOSE_DBSTRING="host=localhost user=postgres password=psql dbname=mydb sslmode=disable"
cd internal/adapters/postgresql
goose up
```

Jika Anda ingin membuat migrasi baru dengan `goose`, jalankan:

```bash
cd internal/adapters/postgresql
goose create create_products sql
```

Jika ingin menghapus seluruh migrasi tabel dengan `goose`, jalankan:

```bash
goose down
```

Lalu tambahkan SQL migrasi ke file yang dibuat.

5. Jika ingin menyesuaikan koneksi database, set environment variable `GOOSE_DBSTRING`:

```bash
export GOOSE_DBSTRING="host=localhost user=postgres password=psql dbname=mydb sslmode=disable"
```

## Menjalankan Aplikasi

Jalankan server dengan:

```bash
go run ./cmd/*.go
```

Server akan berjalan di `http://localhost:8080`.

## Endpoint

- `GET /health` — healthcheck sederhana
- `GET /products` — daftar produk (saat ini mengembalikan hasil kosong dari service)

## Mengembangkan dan Mengubah Query SQL

Jika Anda mengubah file `internal/adapters/postgresql/sqlc/queries.sql`, jalankan ulang generator `sqlc` untuk memperbarui kode:

```bash
go install github.com/kyleconroy/sqlc/cmd/sqlc@latest
sqlc generate
```

## Catatan

- API ini saat ini hanya menyediakan endpoint list produk.
- Service produk belum mengembalikan data nyata dari database; ini dapat dikembangkan lebih lanjut.
