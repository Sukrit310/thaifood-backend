# Admin role management — Backend handoff

งานนี้ทำเฉพาะ Backend ใน repository `thaifood-backend` ตามหัวข้อ 3 ของ `ThaiFood_Admin_Role_Management_Codex.md` ไม่มีการแก้ Frontend หรือฐานข้อมูล production และยังไม่ได้ deploy endpoint ใหม่นี้

## Contract สำหรับ Frontend

ใช้ `PATCH /api/admin/members/{id}/role` พร้อม `Authorization: Bearer <ADMIN_JWT>` และ `Content-Type: application/json`:

```json
{"role":"admin"}
```

ตัวอย่าง response `200 OK` (ข้อมูลสมมติ):

```json
{"id":2,"email":"member@example.com","displayName":"สมาชิกตัวอย่าง","role":"admin"}
```

Response เป็น Member โดยตรง ใช้ชื่อฟิลด์เดิม `id`, `email`, `displayName`, `role`; FE ใช้แทนข้อมูลแถวนั้นหรือ re-fetch `GET /api/admin/members` ได้ Error ใช้ `{"error":{"code":"...","message":"..."}}` ตาม [API Contract](../API_CONTRACT.md#เปลี่ยนสิทธิ์สมาชิก)

- `401 UNAUTHORIZED`: JWT ไม่มี/ผิด/หมดอายุ
- `403 FORBIDDEN`: member, แก้ตัวเอง (รวม role เดิม), หรือผู้สั่งถูกลดสิทธิ์ระหว่างรอ
- `400 VALIDATION_ERROR`: id, body หรือ role ไม่ถูกต้อง
- `404 NOT_FOUND`: ไม่มีสมาชิกเป้าหมาย
- `409 LAST_ADMIN`: last-admin guard แบบ defense in depth
- `500 INTERNAL_ERROR`: ข้อผิดพลาดทั่วไปที่ไม่เปิดเผย SQL/credentials

ส่ง role เดิมของคนอื่นตอบ 200 และไม่เปลี่ยน `updated_at` การส่ง role ผ่าน Register/Profile ยังถูกปฏิเสธ การเปลี่ยนสิทธิ์มีผลต่อ protected request ถัดไปด้วย JWT เดิม

## วิธีลองใน environment ทดสอบ

1. ใช้ PostgreSQL ทดสอบแยก ตั้ง `DATABASE_URL`, `JWT_SECRET`, `ALLOWED_ORIGINS` สำหรับ local ตาม README แล้วรัน `go run . migrate` และ `go run . serve` ไม่ต้องเปลี่ยน schema เพื่อเพิ่ม endpoint นี้
2. สมัครบัญชีทดสอบ A และ B ผ่าน `POST /api/auth/register` ด้วยอีเมลและรหัสผ่านที่ตั้งเอง จากนั้นเลื่อน A เป็น Admin คนแรกด้วย `go run . admin EMAIL_OF_A` โดยใช้ฐานข้อมูลทดสอบเดียวกัน
3. Login ทั้งสองบัญชี เก็บ JWT_A และ JWT_B ไว้ใน environment ส่วนตัว อย่า commit credentials
4. ใช้ JWT_A เรียก `GET /api/admin/members` แล้วเลือก id ของ B จาก response
5. เรียก PATCH ด้วย JWT_A ส่ง `{"role":"admin"}` ไปยัง id ของ B ต้องได้ 200 และ role ใหม่; ใช้ JWT_B เดิมเรียก `GET /api/me` และ `GET /api/admin/members` ต้องเห็น role ใหม่และเข้าถึง Admin API ได้
6. ใช้ JWT_A ลด B ด้วย `{"role":"member"}` ต้องได้ 200; JWT_B เดิมเรียก Admin API ต้องได้ 403 แต่ `GET /api/me` ยังได้ 200
7. ใช้ JWT_A เปลี่ยน role ของ A เองต้องได้ 403; ทดลอง Guest/Member, id ผิด, null, role ผิด หรือฟิลด์เกิน ต้องได้ error ตาม Contract
8. ฝั่ง FE เมื่อเพิ่ม action/dialog และเชื่อม endpoint นี้แล้ว ให้ Login A → Admin > สมาชิก → เลือก B → ยืนยัน role → ตรวจรายการโดยไม่ refresh และ Login B เพื่อตรวจสิทธิ์ใหม่ ขั้นตอนผ่านหน้าเว็บนี้ยังไม่ได้ทดสอบในงาน BE

Postman collection มี sequence 65–78 สำหรับ role management ใช้เฉพาะสมาชิกที่ collection สร้างเอง แล้วลดกลับเป็น member ก่อนจบ การรันค้างหลัง promotion อาจเหลือบัญชีทดสอบเป็น Admin ให้ Admin ลดกลับผ่าน endpoint เมื่อรันทดสอบเสร็จ

## การควบคุม concurrency

Store เปิด READ COMMITTED transaction แล้วล็อก `members` แบบ SHARE ROW EXCLUSIVE ก่อนอ่าน role ผู้สั่งและเป้าหมาย ตรวจจำนวน Admin และอัปเดตด้วย parameters ใน transaction เดียวกัน การล็อกอยู่ใน PostgreSQL จึงครอบคลุมหลาย API instances และมีผลระหว่างรอ concurrent requests ด้วย

การเปลี่ยน role ค้างการเขียน `members` อื่น ๆ ชั่วครู่ (รวม Register/Profile/Password) แต่ SELECT ปกติยังอ่านได้ เลือกวิธีนี้เพราะงานจัดการสิทธิ์มีขนาดเล็กและต้องการตรวจสิทธิ์กับจำนวน Admin อย่างสอดคล้องกัน ไม่มี schema migration/dependency เพิ่ม

Self-change และการตรวจ role ผู้สั่งหลัง lock ป้องกัน Admin สองคนลดสิทธิ์กันจนเหลือศูนย์ Last-admin count เป็นเงื่อนไขเสริมที่ไม่ควรถูกถึงผ่านคำขอปกติเมื่อสองเงื่อนไขแรกยังทำงานอยู่ การแก้ role ด้วย SQL โดยตรงหลัง transaction สิ้นสุดยังเป็นอำนาจของ database owner และอยู่นอก API นี้

## ไฟล์ที่เปลี่ยน

| ไฟล์ | หน้าที่ |
|---|---|
| `handler/handler.go` | ลงทะเบียน admin PATCH route และ mapping `LAST_ADMIN` |
| `handler/auth.go` | ตรวจ id/body/self-change และส่ง Member เดิมกลับ |
| `store/models.go` | เพิ่ม repository method และ error sentinel |
| `store/members.go` | transaction, current actor check, last-admin guard, parameterized UPDATE และ idempotence |
| `handler/member_roles_test.go` | unit tests ของ authorization, validation, response และ error |
| `store/member_roles_test.go` | ทดสอบ persistence/idempotence, guards และ concurrency กับ PostgreSQL จริง |
| `integration/member_roles_test.go` | HTTP → middleware → PostgreSQL, role ใหม่ด้วย JWT เดิม และ escalation rejection |
| `API_CONTRACT.md`, `README.md`, เอกสารนี้ | Contract, คู่มือทดสอบ และข้อจำกัด |
| `postman/ThaiFood_Guide.postman_collection.json`, `postman/README.md` | ตัวอย่างและขั้นตอนทดสอบ endpoint |
| `docs/superpowers/plans/2026-11-14-admin-member-roles.md` | แผนและผลตรวจตามขอบเขต BE |

## ผลตรวจ

รันด้วย Go 1.27.4 และ PostgreSQL 17 แบบ disposable local container ผ่าน `TEST_DATABASE_URL` แยกจาก production:

| การตรวจ | ผลจริง |
|---|---|
| `go test -race -count=1 -json ./...` | ผ่าน 39 top-level tests และ 90 subtests; fail 0, skip 0 |
| `go vet ./...` | ผ่าน (exit 0) |
| `go build ./...` | ผ่าน (exit 0) |
| gofmt และ `git diff --check` | ผ่าน |
| Parse Postman JSON และ compile JavaScript scripts | ผ่าน 78 requests / 82 scripts |
| Execute Admin Login Postman scripts ตามลำดับ prerequest → response → test ด้วย fixture | ผ่าน; เก็บ `adminId`/`adminToken` หลัง response |
| Independent code/test/contract review | ไม่พบข้อสำคัญค้างอยู่; แก้การเก็บ `adminId` ใน Postman ตามผล review แล้ว |

จำนวน top-level tests แยกตาม package: auth 3, handler 13, internal/config 1, mw 6, store 12, integration 4 มี package หลักที่ไม่มี test files จึงไม่รวมในจำนวน test; ไม่มี database integration test ถูก skip

หลักฐานพฤติกรรมใหม่มาจาก real HTTP integration tests และ real PostgreSQL store tests รวมทั้งกรณีคำขอลดสิทธิ์พร้อมกันและการเปลี่ยนสิทธิ์ผู้สั่งระหว่างรอ lock ยังไม่ได้รัน Postman collection ทั้งชุดผ่าน Newman, ยังไม่ได้ทดสอบกับ production API หรือหน้าเว็บ และไม่มีการ deploy/commit/push จึงไม่มี deployment revision หรือ commit SHA ใหม่สำหรับส่งต่อ FE ในรอบนี้
