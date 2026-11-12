# ThaiFood Guide Backend

Go REST API ตาม [Project Brief v2](ThaiFood_Guide_Project_Brief_v2.md) และ [API Contract v2](API_CONTRACT.md) ใช้ `net/http`, `pgxpool`, bcrypt และ JWT เชื่อม PostgreSQL บน Supabase

อ้างอิงการแยกแพ็กเกจจาก [week6-lab-fork/moviehub](https://github.com/Sukrit310/week6-lab-fork/tree/main/moviehub) โดยพัฒนาใน repository นี้เท่านั้น ไม่มีการแก้ไข repository ต้นแบบ

## สิ่งที่ทำได้

- Guest ดูจังหวัด 77 จังหวัด ประเภทอาหาร ร้าน รายละเอียด รีวิว และค้นหา/กรอง/แบ่งหน้า
- Member สมัคร/เข้าสู่ระบบ แก้ชื่อ เปลี่ยนรหัสผ่าน เขียน/ลบรีวิวของตนเอง และจัดการ Favorites
- Admin ดูสมาชิก เพิ่ม/แก้/ลบร้าน และลบรีวิว โดย Go ตรวจ role ปัจจุบันจากฐานข้อมูลทุก protected request
- JWT อายุ 1 ชั่วโมง; Register รับเฉพาะ member; ไม่ส่ง password hash ใน JSON
- คะแนน 1–5 คำนวณจากรีวิวในฐานข้อมูล ร้านไม่มีรีวิวส่ง `rating: null`, `reviewCount: 0`
- Error รูปแบบเดียวกัน: `{"error":{"code":"...","message":"..."}}`

## โครงสร้าง

```text
auth/               bcrypt และ HS256 JWT พร้อม tests
handler/            routes, validation, JSON, API tests
mw/                 JWT/role, CORS, logging, panic recovery
store/              PostgreSQL queries และ database integration tests
internal/config/    environment validation
integration/        HTTP → Go → PostgreSQL ทดสอบทุก endpoint
main.go             server, pgxpool, graceful shutdown, DB/admin commands
schema.sql          ตาราง, FK, CHECK, indexes, RLS (รันซ้ำได้)
seed.sql            77 จังหวัด, 7 ประเภท, 30 ร้านจริงพร้อมรูปและเครดิต (รันซ้ำได้)
postman/            collection และ environment ตัวอย่าง
render.yaml         Render Blueprint
Dockerfile          build แบบหลาย stage ใช้ non-root runtime
.github/workflows/  CI พร้อม PostgreSQL และ race detector
```

## Run กับ Supabase

ต้องมี Go 1.27.4+ (หรือ Go ที่รองรับการดาวน์โหลด toolchain อัตโนมัติตาม go.mod) และ Supabase project สำหรับ ThaiFood โดยเฉพาะ

1. คัดลอก `.env.example` เป็น `.env` แล้วตั้ง `DATABASE_URL`, `JWT_SECRET`, `ALLOWED_ORIGINS`, `PORT`
2. เอา connection string จาก Supabase Dashboard → Connect → **Session pooler** เพื่อรองรับ IPv4; URL-encode อักขระพิเศษในรหัสผ่าน และใช้ TLS (`sslmode=require` เป็นขั้นต่ำ; `verify-full` พร้อม CA ของโปรเจกต์สำหรับตรวจ certificate)
3. สร้าง JWT secret เช่น `openssl rand -hex 32`; ตั้งค่าเฉพาะ environment ห้าม commit
4. รัน:

```bash
set -a
source .env
set +a
go mod download
go run . migrate
go run . seed
go run . serve
```

ตรวจ `curl http://localhost:8080/health` ต้องได้ `{"status":"ok"}`; health ตรวจฐานข้อมูลจริงและตอบ 503 เมื่อใช้ฐานข้อมูลไม่ได้ โปรแกรมไม่โหลด `.env` ให้อัตโนมัติ จึงต้อง export ตามตัวอย่าง

Go เป็นผู้เข้าถึงฐานข้อมูลเพียงจุดเดียว ตารางเปิด RLS และไม่ให้ `anon`/`authenticated` อ่านหรือเขียนผ่าน Supabase Data API; `DATABASE_URL` จึงต้องเป็น database owner ที่เชื่อถือได้ Schema/seed ต้องใช้กับฐานข้อมูลของ ThaiFood เท่านั้น

ชุดร้านเริ่มต้นเป็น **ร้านจริง 30 ร้านใน 19 จังหวัด พร้อมรูปที่ตรงกับร้านและเครดิตภาพ** ใช้ข้อมูล ททท. จาก data.go.th 23 ร้าน หน่วยงานรัฐอื่น 3 ร้าน ข้อมูลสถานที่ใน Wikimedia Commons 3 ร้าน และผู้ประกอบการโดยตรง 1 ร้าน ภาพแต่ละรูปมาจาก Wikimedia Commons พร้อมสัญญาอนุญาตของรูปนั้น ที่มาและเครดิตอยู่ใน `description` และมี `source_name`/`source_url` ในฐานข้อมูล ดู [แหล่งข้อมูลและใบอนุญาตภาพ](docs/RESTAURANT_DATA_SOURCES.md) ข้อมูลเป็น snapshot ไม่รับรองสถานะเปิดทำการล่าสุด และไม่สร้างที่อยู่/พิกัดที่แหล่งข้อมูลไม่มี คะแนนยังคำนวณจากรีวิวของสมาชิกในเว็บเท่านั้น

นำเข้าฐานข้อมูลที่ Render API ใช้แล้วเมื่อ 2026-11-12 ดู [ผลนำเข้า 20 ร้านแรก](docs/REAL_DATA_IMPORT_RESULTS.md) และ [ผลเพิ่ม 10 ร้านในจังหวัดใหม่](docs/TEN_NEW_PROVINCES_IMPORT_RESULTS.md)

การรัน `go run . seed` จาก revision นี้จะนำเข้าร้านจริงและลบเฉพาะร้านสมมติชุดเดิมที่ไม่เคยถูกแก้ไขและไม่มีรีวิวหรือรายการโปรดอ้างถึง ร้านที่สมาชิกใช้งานแล้วและร้านที่ผู้ดูแลเพิ่มหรือแก้ไขเองยังคงอยู่ รันซ้ำไม่เพิ่มร้านซ้ำและไม่ทับข้อมูลที่ Admin แก้ภายหลัง ใช้ source/build ของ revision นี้ เพราะ binary เก่าฝัง seed สมมติไว้ ก่อนใช้กับฐานข้อมูลจริงให้สำรองข้อมูลและตรวจ `DATABASE_URL` ว่าเป็นโปรเจกต์ ThaiFood

## สร้าง Admin

สมัครสมาชิกปกติผ่าน `POST /api/auth/register` ก่อน แล้วผู้ดูแลที่มีสิทธิ์เข้าถึงฐานข้อมูลรัน:

```bash
go run . admin your-admin@example.com
```

คำสั่งนี้เลื่อนสิทธิ์ของบัญชีที่มีอยู่แล้วเท่านั้น ไม่มีรหัสผ่าน Admin เริ่มต้นและไม่มี public API ที่เปลี่ยน role ได้ การเพิ่ม `role` ใน Register หรือ Profile จะถูกปฏิเสธ

## API และ Validation

รายละเอียด request/response ทุก endpoint อยู่ใน [API_CONTRACT.md](API_CONTRACT.md) ซึ่งตรงกับ v2 ที่ให้มา

- Password 8–72 bytes ตามขีดจำกัด bcrypt; displayName 1–100 ตัวอักษร
- อีเมล trim และ lowercase มี unique index แบบ case-insensitive
- รีวิวใช้ integer 1–5 และข้อความ trim 10–5000 ตัวอักษร Unicode; หนึ่งสมาชิกต่อหนึ่งร้านได้หนึ่งรีวิว ซ้ำตอบ 409
- Filters `province`/`category` ใช้ชื่อไทย; `q` ค้นหาชื่อ รายละเอียด และที่อยู่แบบ literal (SQL parameters)
- `page=1`, `limit=12`; limit เกิน 50 ใช้ 50; page/limit ไม่ใช่จำนวนเต็มบวกตอบ 400
- Admin POST ร้านต้องมี name, provinceId, categoryId, address, description; พิกัดเลือกใส่ได้แต่ต้องเป็นคู่และอยู่ในช่วงที่ถูกต้อง; imageUrl รับ HTTP(S) หรือ null
- Admin PUT ร้านส่งเฉพาะฟิลด์ที่แก้ได้ โดยฟิลด์ที่ละไว้คงเดิม; `imageUrl: null` ล้างรูป
- JSON ต้องเป็น object เดียว ไม่รับฟิลด์ไม่รู้จัก ขนาดไม่เกิน 1 MiB
- ไม่พบข้อมูลตอบ 404, ข้อมูลซ้ำ 409, ไม่มี/Token ผิด/หมดอายุ 401, สิทธิ์ไม่พอ 403, database failure 500 แบบไม่เปิดเผยข้อมูลลับ

Logout ของ MVP ให้ Frontend ลบ JWT และข้อมูล session จาก browser; Token เดิมยังใช้ได้จนหมดอายุ รวมถึงหลังเปลี่ยนรหัสผ่าน ไม่มี logout/revocation endpoint ตาม Contract

## ทดสอบ

Unit tests รันได้โดยไม่ใช้ DB; integration tests จะระบุ skip หากไม่มี `TEST_DATABASE_URL`:

```bash
go test -race -count=1 ./...
go vet ./...
go build -trimpath -o bin/thaifood .
```

ทดสอบครบกับ PostgreSQL จริง ใช้ฐานข้อมูลทดสอบแยกจาก production:

```bash
export POSTGRES_PASSWORD='choose-a-local-test-password'
docker compose up -d --wait
export TEST_DATABASE_URL="postgresql://postgres:${POSTGRES_PASSWORD}@127.0.0.1:5432/thaifood_test?sslmode=disable"
make integration
```

Tests สร้าง schema ชื่อสุ่มและลบทิ้งเมื่อเสร็จ ตรวจ replay schema/seed, ทุก endpoint, Guest/Member/Admin, password/profile, malformed/expired/forged JWT, Thai filters, validation, duplicate requests พร้อมกัน, SQL injection, ownership, ค่าเฉลี่ยรีวิวและ cascade delete CI ใช้ PostgreSQL 17 และรันทุก test พร้อม race detector

Postman: Import collection และ environment ใน `postman/`, ตั้ง `baseUrl`, `adminEmail`, `adminPassword` เป็นบัญชีที่เลื่อนสิทธิ์แล้ว แล้ว Run collection ตามลำดับ Collection สมัครสมาชิกทดสอบเองและทำ CRUD เฉพาะร้านที่สร้างในรอบนั้น จึงควรรันบนฐานข้อมูลทดสอบ

```bash
npx newman run postman/ThaiFood_Guide.postman_collection.json \
  -e /path/to/private-postman-environment.json
```

เก็บ environment ที่มีรหัสผ่านจริงไว้นอก Git; ดูผลการตรวจล่าสุดใน [docs/TEST_RESULTS.md](docs/TEST_RESULTS.md)

เมื่อรันผ่าน Render ให้ตั้ง `allowEdgeSqlProbeBlock=true` ใน private Postman environment เพื่อระบุผลเฉพาะ SQL probe สุดท้ายที่ Cloudflare WAF บล็อกด้วย403 ค่าเริ่มต้นยังตรวจ200JSONอย่างเคร่งครัด และทุก API ปกติยังใช้ Contract เดิม ดูรายละเอียดใน [postman/README.md](postman/README.md)

## Deploy Render

`render.yaml` ตั้ง Docker web service, free plan, Singapore, branch `feat/thaifood-guide-backend`, health check `/health` และปิด auto deploy เพื่อให้ทดสอบก่อนแต่ละ release ไม่มีการสร้างฐานข้อมูล Render เพราะใช้ Supabase

1. ทดสอบตามหัวข้อด้านบนก่อน push branch นี้ไป repository ThaiFood
2. ตั้ง Supabase project และรัน `migrate`/`seed` ด้วย `DATABASE_URL` ของโปรเจกต์นั้น จากเครื่องที่เชื่อมต่อได้ (free Render ไม่มี pre-deploy command)
3. สร้าง Render Web Service จาก repo/branch นี้ด้วย Docker หรือใช้ Blueprint `render.yaml`
4. ตั้ง `DATABASE_URL` ใน secret environment, สร้าง `JWT_SECRET` แบบสุ่ม และตั้ง `ALLOWED_ORIGINS=http://localhost:3000,https://ACTUAL-FRONTEND.vercel.app` โดยแทน URL ด้วยโดเมนจริงที่ยืนยันแล้ว
5. Render ส่ง `PORT` มาให้โปรแกรมฟังบน `0.0.0.0:$PORT` อัตโนมัติ
6. รอ deployment สถานะ live แล้วตรวจ `/health`, ร้าน, Register/Login/Profile/Review/Favorites และ Admin ผ่าน Postman อีกครั้ง
7. ส่ง Render URL จริงให้ Frontend ตั้ง `REACT_APP_API_URL` และ rebuild CRA

API ที่ deploy แล้ว: **https://thaifood-guide-api-a7c7.onrender.com**

Health: https://thaifood-guide-api-a7c7.onrender.com/health

ใช้งาน Supabase project `xorsagxmphehnjurquen` ผ่าน Session pooler พร้อม TLS; Render service อยู่บน Free plan ที่ Singapore และตั้ง `/health` แล้ว ไม่มีการเปลี่ยน MovieHub ต้นแบบ

CORS ปัจจุบันอนุญาต `http://localhost:3000` เท่านั้น ผู้ใช้จะเพิ่ม URL Vercel ภายหลัง โดย append origin จริงใน `ALLOWED_ORIGINS` แล้ว deploy ให้ค่าใหม่มีผล ดู [คู่มือส่งต่องาน Frontend](docs/FRONTEND_HANDOFF.md)

สำหรับ FE ตั้ง `REACT_APP_API_URL=https://thaifood-guide-api-a7c7.onrender.com` และ rebuild CRA หลังแก้ค่า
