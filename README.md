# Every Nation Indonesia Backend

Backend untuk website Every Nation Indonesia yang dibangun menggunakan Go dan Gin Framework.

## Latar Belakang

Every Nation Indonesia adalah jaringan gereja yang tersebar di seluruh Indonesia. Website ini dibangun untuk memudahkan pengelolaan data dan informasi terkait gereja-gereja Every Nation di Indonesia, termasuk manajemen jemaat, kelompok kehidupan (LifeGroup), dan berbagai aktivitas gereja.

## Fitur-fitur

### 1. Manajemen Gereja
- Pendaftaran dan pengelolaan data gereja
- Informasi detail gereja (alamat, kontak, dll)
- Pengelompokan gereja berdasarkan wilayah (provinsi dan kota)

### 2. Manajemen Jemaat
- Pendaftaran dan pengelolaan data jemaat
- Informasi detail jemaat
- Pengelompokan jemaat berdasarkan gereja

### 3. Manajemen LifeGroup
- Pembuatan dan pengelolaan kelompok kehidupan
- Penugasan pemimpin LifeGroup
- Pengelolaan anggota LifeGroup
- Informasi jadwal dan lokasi pertemuan

### 4. Manajemen Departemen
- Pengelolaan departemen gereja
- Penugasan staff departemen
- Dokumentasi aktivitas departemen

### 5. Manajemen Pengguna
- Sistem autentikasi dan otorisasi
- Pengelolaan hak akses pengguna
- Profil pengguna

### 6. Manajemen Notifikasi
- Sistem notifikasi untuk berbagai aktivitas
- Pengiriman notifikasi ke pengguna
- Riwayat notifikasi

## Teknologi yang Digunakan

- Go (Golang)
- Gin Framework
- MySQL
- JWT untuk autentikasi
- GORM untuk ORM

## Struktur Proyek

```
.
├── config/         # Konfigurasi aplikasi
├── controller/     # Controller untuk handling request
├── dto/           # Data Transfer Object
├── entity/        # Model/Entity
├── middleware/    # Middleware (auth, logging, dll)
├── repository/    # Repository untuk akses database
├── routes/        # Definisi routes
├── service/       # Business logic
└── main.go        # Entry point aplikasi
```

## Cara Menjalankan

1. Clone repository
2. Install dependencies
3. Setup database
4. Jalankan aplikasi

## Kontribusi

Silakan berkontribusi dengan membuat pull request atau melaporkan issues.

## Lisensi

[MIT License](LICENSE)