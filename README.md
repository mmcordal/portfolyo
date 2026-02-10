# portfolyo-api
## Portföy & Varlık Yönetim API (Go)

Clean Architecture prensiplerine uygun olarak Go ile geliştirilmiş bu servis; kullanıcıların finansal varlıklarını (altın, döviz vb.) yönetebildiği, işlemlerini kaydedebildiği, toplam portföy değerini hesaplayabildiği ve rapor çıktıları (PDF / Excel) üretebildiği bir backend uygulamasıdır.

---

### Bu projenin amacı:
- Gerçek bir senaryo üzerinden servis katmanı kurgusunu öğrenmek
- İş mantığını framework bağımsız şekilde tasarlamak
- Katmanlı mimari (Clean Architecture) ile sürdürülebilir backend geliştirmek
- Kullanıcı bazlı finansal işlem yönetimini modellemek
- Raporlama (PDF/Excel) üretme mantığını öğrenmek

---

### Ana Özellikler:
- Kullanıcı bazlı varlık işlemleri (ekleme / çıkarma)
- Altın, döviz ve benzeri varlık tipleri desteği
- Hedef para birimine göre toplam portföy hesaplama
- İşlem geçmişi listeleme
- İşlem bazlı dekont (PDF) üretme
- Tüm işlemleri Excel olarak dışa aktarma
- JWT tabanlı kimlik doğrulama
- Clean Architecture ile uyumlu katmanlı mimari
- Merkezi hata yönetimi yapısı

---

### Mimari Genel Bakış:
```pgsql
cmd/                    → Uygulama giriş noktası
internal/
   ├── handler/         → HTTP katmanı (endpointler)
   ├── service/         → İş mantığı katmanı
   ├── repository/      → Veri erişim katmanı
   ├── model/           → Domain modelleri
   ├── viewmodel/       → API response modelleri
   └── document/        → PDF & Excel üretimi
```
---

### Mimari Prensibi:
Bağımlılıklar içeri doğru akar. İş mantığı; framework, veritabanı ve dış kütüphanelerden bağımsız kalır.

---

### Kullanılan Teknolojiler:
- Go
- Fiber (HTTP framework)
- PostgreSQL
- GORM
- JWT (Authentication)
- gofpdf (PDF üretimi)
- excelize (Excel üretimi)

---

### Uygulama Akışı:

- Kullanıcı sisteme kayıt olur ve giriş yapar. 
- JWT token ile kimlik doğrulaması yapılır. 
- Kullanıcı varlık işlemleri ekler veya çıkarır. 
- Servis katmanı işlemleri işler ve portföy toplamını hesaplar. 
- Kullanıcı:
  - Tüm işlemlerini JSON olarak alabilir 
  - İşlem bazlı dekontu PDF olarak indirebilir 
  - Tüm işlemleri Excel olarak dışa aktarabilir

---

## Kurulum ve Çalıştırma
### Gereksinimler

- Go 1.20+ 
- PostgreSQL

### Uygulamayı çalıştırma
```bash
go mod tidy
go run cmd/main.go
```

---

### Varsayılan Yapılandırma
Uygulama iç yapılandırma dosyaları üzerinden veritabanı ve sunucu ayarlarını kullanır.

---

### Örnek Endpointler
```bash
POST   /auth/register
POST   /auth/login
GET    /transactions
GET    /transactions/excel
GET    /transactions/pdf/:id
GET    /portfolyo
```

---

> **Notlar:**
- Proje öğrenme ve pratik amaçlı geliştirilmiştir.
- Finansal hesaplamalar örnek senaryo içindir, gerçek yatırım aracı değildir.
- Raporlama çıktıları kullanıcı işlemlerine göre dinamik üretilir.

---

---
# EN

---

# portfolyo-api
## Portfolio & Asset Management API (Go)

This service is developed in Go following Clean Architecture principles. It allows users to manage financial assets (gold, currencies, etc.), record transactions, calculate total portfolio value, and generate reports in PDF and Excel formats.

---

### Purpose of this project:
- Learn service layer design through a real-world scenario
- Design business logic independently from frameworks
- Build a sustainable backend using layered architecture (Clean Architecture)
- Model user-based financial transaction management
- Learn report generation (PDF/Excel)

---

### Key Features:
- User-based asset transactions (add / subtract)
- Support for gold, currency, and similar asset types
- Portfolio total calculation based on target currency
- Transaction history listing
- Transaction-based receipt generation (PDF)
- Export all transactions to Excel
- JWT-based authentication
- Layered architecture aligned with Clean Architecture
- Centralized error handling structure

---

### Architecture Overview:
```pgsql
cmd/                    → Application entry point
internal/
   ├── handler/         → HTTP layer (endpoints)
   ├── service/         → Business logic layer
   ├── repository/      → Data access layer
   ├── model/           → Domain models
   ├── viewmodel/       → API response models
   └── document/        → PDF & Excel generation
   ```

---

### Architectural Principle:
Dependencies flow inward. Business logic remains independent from frameworks, databases, and external libraries.

---

### Technologies Used:
- Go  
- Fiber (HTTP framework)  
- PostgreSQL  
- GORM  
- JWT (Authentication)  
- gofpdf (PDF generation)  
- excelize (Excel generation)

---

### Application Flow:

- User registers and logs in  
- JWT authentication is applied  
- User adds or removes asset transactions  
- Service layer processes transactions and calculates portfolio totals  
- User can:
  - Retrieve all transactions as JSON  
  - Download a transaction receipt as PDF  
  - Export all transactions as Excel  

---

## Setup & Run
### Requirements

- Go 1.20+  
- PostgreSQL  

---

### Run the application
```bash
go mod tidy
go run cmd/main.go
```

---

### Default Configuration
The application uses internal configuration files for database and server settings.

---

### Example Endpoints
```bash
POST   /auth/register
POST   /auth/login
GET    /transactions
GET    /transactions/excel
GET    /transactions/pdf/:id
GET    /portfolio
```

---

> **Notes:**
- The project is built for learning and practice purposes. 
- Financial calculations are for demonstration only and not investment advice. 
- Reports are dynamically generated based on user transactions.

---