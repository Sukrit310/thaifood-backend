# ผลการตรวจ ThaiFood Guide Backend

Latest additive data run: [23 more provinces and current 53-restaurant verification](TWENTY_THREE_NEW_PROVINCES_IMPORT_RESULTS.md). Earlier results below remain historical.

Later additive data run: [ten new provinces and current 30-restaurant verification](TEN_NEW_PROVINCES_IMPORT_RESULTS.md). The results below describe the earlier run.

รายงานด้านล่างเป็นผลตรวจชุดเดิมก่อนเปลี่ยนข้อมูลสาธิตเป็นร้านจริง ผลการนำเข้าร้านจริงและทดสอบล่าสุดวันที่ 2026-11-12 อยู่ใน [REAL_DATA_IMPORT_RESULTS.md](REAL_DATA_IMPORT_RESULTS.md)

ตรวจระบบ local วันที่ 2026-11-11 และ Supabase/Render วันที่ 2026-11-12 (Asia/Bangkok). ตรวจใน repository `Sukrit310/thaifood-backend` บน branch `feat/thaifood-guide-backend` โดยไม่มีการแก้ repository MovieHub ต้นแบบ

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

## Supabase และ Render จริง

- Supabase project: `xorsagxmphehnjurquen`, PostgreSQL17.11 ผ่าน Session pooler/TLS ตรวจแล้วว่า public schema ว่างก่อนใช้ schema.sql และ seed.sql
- หลัง seed มี77จังหวัด,7ประเภท,25ร้านสมมติ ข้อมูลอยู่ใน Supabase จริง
- รัน `go test -race -count=1 -v ./store ./integration` กับ Supabase ใน schema ทดสอบแยก ผ่าน8 top-level tests และลบ schema ทดสอบเมื่อเสร็จ
- Postman ก่อน deploy ผ่าน Go runtimeที่เชื่อม Supabase จริง64requests/296assertions/0failures
- [GitHub Actions Backend checks](https://github.com/Sukrit310/thaifood-backend/actions/runs/43065005725) ของ commit `53cdab8` ผ่าน formatting,vet,race testsและbuild
- Render service: `thaifood-guide-api-a7c7` (`srv-dbqaciqjnfac73dbf0f0`), Free/Docker/Singapore,health path `/health`,manual deploy
- Deploy `dep-dbqacjijnfac73dbf24g` ของ commit `53cdab8` ขึ้นสถานะliveแล้ว
- API URL: https://thaifood-guide-api-a7c7.onrender.com
- Live smoke ผ่านDB health,77จังหวัด,localhost CORS preflightและการปฏิเสธoriginที่ไม่อนุญาต

### SQL probe ที่ชั้น Cloudflare

PostmanรอบแรกบนRender:63requestsผ่านตามContract; request64ที่จำลองSQL injectionได้403HTMLจากCloudflare WAF ทำให้4assertionsที่คาด200JSONไม่ผ่าน หน้าresponseระบุว่าถูก web application firewall บล็อก พร้อมCF-Ray การป้องกันนี้อยู่ก่อนGo ([Renderอธิบายว่าทุกrequestผ่านCloudflare](https://render.com/docs/uptime-best-practices#log-the-cf-ray-id-of-each-request))

การค้นหาSQL metacharactersที่Go/Supabaseโดยตรงผ่านการทดสอบแล้ว ไม่มีการเปลี่ยนSQL queriesหรือปิดWAFเพื่อแก้ผลทดสอบ Collectionคงstrict200JSONเป็นค่าเริ่มต้น และเพิ่มตัวเลือก `allowEdgeSqlProbeBlock=true` สำหรับการตรวจRenderเท่านั้น ซึ่งต้องตรวจrequestสุดท้าย,HTTPS onrender.com,403,Server Cloudflare,CF-Ray,HTML Content-Typeและข้อความWAFที่ตรงกันก่อนยอมรับการบล็อก ทุกrequestอื่นยังตรวจContractเดิม

### ผลตรวจ live หลังระบุชั้นตอบกลับอย่างชัดเจน

- PostmanบนRenderผ่าน64requests/294assertions/0failures โดย63requestsตรวจContractของGoตามเดิม และ1SQLprobeยืนยัน403จากCloudflare WAFด้วยเงื่อนไขเฉพาะข้างต้น
- หลังแก้collection รันSQLprobeผ่านGo/Supabaseโดยตรงซ้ำด้วยค่าเริ่มต้นstrict ผ่าน1request/5assertions ยังคงต้องได้200JSONและผลค้นหาว่าง
- ไม่ได้เปลี่ยนGo API/SQLหรือปิดการป้องกันของRenderเพื่อให้ผลทดสอบผ่าน
- ลบเฉพาะบัญชีชั่วคราว7บัญชีที่บันทึกid/emailจากผลทดสอบแล้ว;เหลือ77จังหวัด/7ประเภท/25ร้านสมมติ และไม่มีสมาชิก/รีวิว/Favoritesทดสอบค้างในฐานข้อมูล
- ยืนยันว่าJWTของAdminทดสอบที่ลบแล้วได้401 และhealthหลังcleanupยังได้200

### Frontend ที่ยังรอ URL จริง

ผู้ใช้ยืนยันให้เพิ่ม Vercel origin ภายหลัง CORSปัจจุบันจึงอนุญาตเฉพาะ `http://localhost:3000` ดู[Frontend handoff](FRONTEND_HANDOFF.md)สำหรับการตั้ง `REACT_APP_API_URL` และappendโดเมนจริงในRender `ALLOWED_ORIGINS`
