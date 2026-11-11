# ผลการตรวจ ThaiFood Guide Backend

Session date: 2026-11-11. ตรวจใน repository `Sukrit310/thaifood-backend` บน branch `feat/thaifood-guide-backend` โดยไม่มีการแก้ repository MovieHub ต้นแบบ

## สภาพแวดล้อมและผล

| รายการ | ผล |
|---|---|
| Go toolchain | 1.27.4, Linux amd64 |
| PostgreSQL | 17, container ทดสอบแยก ไม่มีการใช้ฐานข้อมูล MovieHub |
| `go test -race -count=1 -json ./...` | ผ่าน 71 test/subtest cases, 0 failed, 0 skipped test cases |
| `go vet ./...` | ผ่าน |
| `go build -trimpath -o bin/thaifood .` | ผ่าน และตรวจ binary เป็น Go1.27.4 |
| `govulncheck ./...` | No vulnerabilities found |
| Postman/Newman 6 | 64 requests, 296 assertions, 0 failures |
| `docker build -t thaifood-guide:test .` | ผ่าน |
| Container runtime | non-root `app`, PORT8081, health/DB/provinces/CORS ผ่าน |
| `render.yaml` | ผ่าน official Render JSON Schema; ไม่ได้ตรวจผ่าน Render CLI ซึ่งไม่ได้ติดตั้ง |
| `migrate` และ `seed` ผ่าน runtime binary | รันซ้ำสองครั้งผ่าน ข้อมูลไม่ซ้ำ |
| API Contract copy | `API_CONTRACT.md` ตรงกับไฟล์ v2 ที่ให้มา |

แพ็กเกจ root แสดง `[no test files]`; unit/config/store/HTTP integration ทุกแพ็กเกจมีการทดสอบจริงและไม่มี test case ถูก skip การรัน PostgreSQL tests ใช้ `TEST_DATABASE_URL` ที่ตั้งไว้และสร้าง schema สุ่มแยกกัน

## ขอบเขตที่ยืนยันแล้ว

- ครบ 20 endpoints ตาม Contract v2 และ GET /health
- Register/Login/Me/Profile/Change Password พร้อม bcrypt, token ผิด/หมดอายุ/ปลอม/algorithm ไม่ถูกต้อง
- Register/Profile ไม่รับ role; Member เข้า Admin ไม่ได้; role เปลี่ยนใน DB มีผลกับ token เดิมทันที
- Member/Owner/Admin ลบรีวิวตามสิทธิ์; Favorites แยกตามสมาชิก
- รีวิว integer1–5, ความยาว Unicode ไทย, รีวิวซ้ำ409, ค่าเฉลี่ย/จำนวนรีวิวเปลี่ยนหลังลบ; เมื่อหมดรีวิว rating เป็น null
- Pagination/default/cap, Thai province/category filters, literal search และ SQL-injection inputs
- Malformed/oversized/trailing JSON, unknown fields, NUL และ invalid UTF-8 filter ได้400
- อีเมลไม่สนตัวพิมพ์, FK constraints, missing resources404, missing province/category400
- Concurrent duplicate review/Favorite มีเพียงหนึ่ง201 ที่เหลือ409; store มี concurrent uniqueness tests เพิ่มเติม
- ลบร้านแล้ว cascade reviews/Favorites; empty items เป็น []
- Generic500/503 ไม่ส่ง credential/hash หรือ database connection string
- 77 จังหวัด, 7 ประเภท, 25 ร้านสมมติที่ผู้จัดทำสร้างเอง พร้อม source metadata; schema/seed รันซ้ำได้
- Runtime binary ใช้ pgx QueryExecModeExec และ statement timeout เช่นเดียวกับ production: Newman ตรวจทั้งอ่าน/เขียน/Profile/Reviews/Favorites/Admin CRUD ผ่าน binary นี้จริง
- CORS ผ่านทั้ง middleware tests และ HTTP requests จริงต่อ Docker container: localhost อนุญาต, preflight PATCH/Authorization ผ่าน, origin ที่ไม่อนุญาต403

Review แยกตรวจพบ invalid UTF-8 filter ซึ่งแก้และเพิ่ม regression tests แล้ว การสแกนครั้งแรกพบ GO-2026-6859 ใน standard library ของ Go1.27.3 จึงเปลี่ยน toolchain/CI/Docker เป็น Go1.27.4 และสแกนซ้ำผ่าน

Raw Newman report มี test credentials/JWT จึงเก็บนอก repository และไม่ commit; ไฟล์นี้บันทึกเฉพาะผลสรุป ไม่มี secrets

## สิ่งที่ยังไม่ได้ยืนยันบนระบบจริง

- ยังไม่มี `DATABASE_URL` ของ Supabase สำหรับ ThaiFood ใน environment; connected project ที่พบคือ `movie-hub-db` สถานะ INACTIVE จึงไม่เปลี่ยนข้อมูลหรือสถานะของโปรเจกต์นั้น
- ยังไม่ได้เลือก/ยืนยัน Render workspace และยังไม่ทราบ URL Vercel จริง จึงไม่ได้สร้าง service หรือ deploy
- ยังไม่มี Render API URL จริง และยังไม่ได้ทดสอบผ่าน Supabase/Render จริง

เมื่อกำหนดเป้าหมายและ credentials แล้ว ต้องรัน migrate/seed, ทดสอบ API ที่เชื่อม Supabase, deploy branch ที่ผ่าน checks, รอ Render สถานะ live และทดสอบ Postman ซ้ำก่อนถือว่าการ deploy เสร็จ
