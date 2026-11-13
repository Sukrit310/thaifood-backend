// ไฟล์นี้ระบุชื่อ Go module รุ่นภาษา/toolchain และเวอร์ชัน dependencies ของโปรเจกต์
module thaifood-backend

// เวอร์ชัน Go ขั้นต่ำของ module และ toolchain ที่ต้องการใช้ build/test
go 1.26.0

toolchain go1.27.4

// Dependencies ที่โค้ดเรียกโดยตรง: JWT, PostgreSQL driver/pool และ bcrypt
require (
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/jackc/pgx/v5 v5.12.0
	golang.org/x/crypto v0.59.0
)

// Dependencies ทางอ้อมที่ไลบรารีด้านบนต้องใช้; Go ใช้ร่วมกับ go.sum เพื่อตรวจสอบ dependency
require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.3 // indirect
	golang.org/x/sync v0.25.0 // indirect
	golang.org/x/text v0.44.0 // indirect
)
