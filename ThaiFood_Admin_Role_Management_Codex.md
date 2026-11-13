# ThaiFood Guide — Admin Role Management

> **ไฟล์ข้อกำหนดสำหรับส่งให้ Codex**  
> โครงการ: ThaiFood Guide (React Frontend + Go REST API + PostgreSQL/Supabase)  
> งาน: เพิ่มความสามารถให้ผู้ดูแลระบบเปลี่ยนสิทธิ์สมาชิกจากหน้าเว็บได้โดยตรง  
> หมายเหตุ: **ต้องแก้ทั้ง Backend และ Frontend** โดยทำงานตาม API Contract เดียวกัน

## 1. Repository ที่เกี่ยวข้อง

- Backend (BE): https://github.com/Sukrit310/thaifood-backend
- Frontend (FE): https://github.com/Sukrit310/thaifood-frontend
- Frontend source อยู่ภายในโฟลเดอร์ `thaifood/`

**สำคัญ:** ให้ Codex สำรวจไฟล์จริงใน Repository ที่เปิดอยู่ก่อนแก้โค้ด อย่าสร้างระบบ Authentication/Authorization ใหม่ซ้ำกับของเดิม และห้ามแก้ Repository อีกฝั่งโดยไม่ได้รับมอบหมาย หากใช้ Codex คนละเครื่อง ให้ทำ BE และ FE แยกตามหัวข้อ 3–4 ด้านล่าง

## 2. เป้าหมายและขอบเขต

ปัจจุบันหน้า **Admin > สมาชิก** แสดงรายชื่อ อีเมล และสิทธิ์ (`member` / `admin`) อยู่แล้ว ต้องเพิ่มการจัดการสิทธิ์แบบใช้งานได้จริงดังนี้

1. Admin เห็นปุ่ม **เปลี่ยนสิทธิ์** ในแถวสมาชิกที่สามารถแก้ไขได้
2. กดแล้วเปิดหน้าต่างยืนยัน (modal/dialog) โดยแสดงชื่อ อีเมล สิทธิ์ปัจจุบัน และตัวเลือกสิทธิ์ใหม่ `member` หรือ `admin`
3. กด **ยืนยัน** จึงส่งคำขอเปลี่ยนสิทธิ์ไปยัง Go Backend จริง
4. เมื่อสำเร็จ รายการสมาชิกอัปเดตสิทธิ์บนหน้าเว็บทันทีโดยไม่ต้องรีโหลดหน้า (จะ re-fetch รายการก็ได้)
5. แสดง Loading / Success / Error จากผลลัพธ์จริง และป้องกันการกดบันทึกซ้ำระหว่างส่งคำขอ
6. ห้ามเปลี่ยนสิทธิ์บัญชี Admin ที่ล็อกอินอยู่เอง และห้ามทำให้ระบบไม่มี Admin เหลืออยู่
7. Member/Guest ต้องไม่สามารถใช้ Endpoint เปลี่ยนสิทธิ์ได้ แม้เรียก API ตรงด้วย Postman/curl
8. คง UI/UX, Design System, Routing, การทำงานของ Register/Login/Profile และฟีเจอร์เดิมทั้งหมด

**ไม่อยู่ในขอบเขตงาน:** การเพิ่ม role ใหม่, การลบสมาชิก, การแก้รหัสผ่านของผู้อื่น, การแก้ฐานข้อมูล production ด้วยมือ, การสร้างระบบสิทธิ์ชุดใหม่

## 3. งานฝั่ง Backend — `thaifood-backend`

### 3.1 ตรวจโครงสร้างเดิมก่อน

ตรวจและนำรูปแบบเดิมมาใช้ในไฟล์ที่เกี่ยวข้อง เช่น:

- `handler/handler.go` — ลงทะเบียน route และตัวช่วยส่ง JSON/error
- `handler/auth.go` — API ผู้ใช้/รูปแบบ Response
- `mw/auth.go` — `RequireLogin`, `RequireRole` (ตรวจ role ปัจจุบันจาก DB)
- `store/members.go`, `store/models.go`, `store/store.go` — Model, Repository, SQL
- `schema.sql` — ตาราง `members` และ constraints
- `API_CONTRACT.md` — เอกสาร API
- `integration/`, `handler/*_test.go`, `store/*_test.go` — รูปแบบ Test

### 3.2 เพิ่ม Endpoint

```http
PATCH /api/admin/members/{id}/role
Authorization: Bearer <ADMIN_JWT>
Content-Type: application/json

{
  "role": "admin"
}
```

**พฤติกรรมที่ต้องมี**

- ใช้ middleware เดิม: `RequireLogin` + `RequireRole("admin")` หรือ helper `admin(...)` ที่มีอยู่แล้ว
- อ่าน `id` เป็นจำนวนเต็มบวกและตรวจว่าพบบัญชีจริง
- รับ JSON object ที่มี `role` เพียงฟิลด์เดียวและเป็นค่า `member` หรือ `admin` เท่านั้น โดยใช้ validation แนวเดียวกับโค้ดเดิม (ไม่รับฟิลด์แปลก, JSON ผิดรูปแบบ, ค่า null ฯลฯ)
- อัปเดต `members.role` ใน PostgreSQL ผ่านชั้น Store/Repository ด้วย parameterized SQL; ไม่ต่อ SQL string จาก input
- ให้การตรวจเงื่อนไขและการอัปเดตข้อมูลเป็นการดำเนินการที่ปลอดภัยต่อ concurrent requests ตามสมควร
- ห้าม Admin เปลี่ยนสิทธิ์บัญชีตัวเอง โดยตรวจจาก identity ของ JWT ที่ Middleware ยืนยันแล้ว เทียบกับ `{id}` ไม่ใช่จากข้อมูลที่ Frontend ส่งมา
- ห้ามลดสิทธิ์ Admin คนสุดท้าย (defense in depth) และห้ามปล่อยให้ระบบไม่มี Admin ที่ใช้งานได้
- หลังเปลี่ยนสิทธิ์ **สิทธิ์ใหม่ต้องมีผลกับ API request ถัดไป** เพราะ middleware อ่าน role ปัจจุบันจากฐานข้อมูลอยู่แล้ว
- ไม่เปิด API สำหรับการกำหนด role ผ่าน Register หรือ Profile และไม่ส่ง password hash, secret หรือ token กลับไป
- กรณีเปลี่ยนเป็น role เดิมให้กำหนดพฤติกรรมที่ชัดเจน (แนะนำตอบสำเร็จแบบ idempotent) และบันทึกใน Contract

**รูปแบบ Response**: ใช้ JSON member model/serialization ที่ระบบใช้ใน `GET /api/admin/members` หรือ `GET /api/me` อยู่แล้ว เพื่อไม่สร้าง schema สองแบบโดยไม่จำเป็น ระบุให้แน่ชัดใน API Contract

**HTTP status ที่ควรครอบคลุม** (ให้คงรูปแบบ `{"error":{"code":"...","message":"..."}}` เดิม):

| สถานการณ์ | ค่าที่คาดหวัง |
|---|---|
| Admin เปลี่ยนสิทธิ์สำเร็จ | `200 OK` + ข้อมูลสมาชิกที่อัปเดตแล้ว |
| ไม่ส่ง/ส่ง JWT ไม่ถูกต้อง | `401 Unauthorized` |
| Member เรียก Admin API | `403 Forbidden` |
| Admin เปลี่ยนสิทธิ์ตัวเอง / ขัดกับกฎความปลอดภัย | `403 Forbidden` หรือ `409 Conflict` โดยเลือกให้สอดคล้องกับ Contract เดิม |
| `id` ผิดรูปแบบ / `role` ผิดค่า / JSON ผิด | `400 Bad Request` |
| ไม่พบสมาชิกเป้าหมาย | `404 Not Found` |
| Database/Server ขัดข้อง | ใช้รูปแบบ error ทั่วไปของระบบเดิมและไม่เผยข้อมูลลับ |

> ห้ามเชื่อ UI เพียงอย่างเดียว: ทุกข้อจำกัดสำคัญต้องบังคับใช้ใน Backend ด้วย

### 3.3 Tests ที่ต้องเพิ่ม

- Admin เปลี่ยน `member -> admin` สำเร็จ และอ่านข้อมูลซ้ำได้ role ใหม่
- Admin เปลี่ยน `admin -> member` ของ **บัญชี Admin คนอื่น** ได้เมื่อยังมี Admin เหลือ
- Member/Guest เรียก Endpoint แล้วถูกปฏิเสธ
- ไม่อนุญาต Admin ลด/เปลี่ยนสิทธิ์ตัวเอง
- ไม่อนุญาตให้ระบบเหลือ Admin = 0 รวมถึงกรณี concurrent requests ที่เกี่ยวข้อง
- `id` ไม่ถูกต้อง, ไม่พบ user, `role` ไม่ถูกต้อง, ฟิลด์เกิน, JSON ผิด
- Registration/Profile ไม่สามารถยกระดับสิทธิ์ได้
- Role ใน Database ที่เปลี่ยนแล้วมีผลกับ protected endpoint ถัดไป
- รักษา Tests เดิมให้ผ่านทั้งหมด

คำสั่งตรวจสอบตามโครงการ (ต้องรายงานผลจริง ไม่อ้างว่าผ่านถ้าไม่ได้รัน):

```bash
go test -race -count=1 ./...
go vet ./...
go build ./...
```

สำหรับ Integration Tests ต้องใช้ฐานข้อมูลทดสอบแยก (`TEST_DATABASE_URL`) ไม่ใช่ production; แจ้งให้ชัดถ้าบางชุดถูก Skip

### 3.4 เอกสารที่ต้องอัปเดต

- `API_CONTRACT.md`: path, method, auth, request, response, error, ข้อจำกัดเรื่อง self/last admin
- README หรือ docs API ที่เกี่ยวข้อง
- เพิ่มตัวอย่าง API ทดสอบใน Postman Collection ถ้าตรงกับโครงสร้างโครงการ โดยไม่ใส่รหัสผ่าน/token จริง

## 4. งานฝั่ง Frontend — `thaifood-frontend/thaifood`

### 4.1 ตรวจโครงสร้างเดิมก่อน

- `src/pages/admin/AdminMembers.jsx` — หน้ารายชื่อสมาชิก
- `src/api/services.js` — `adminService`
- `src/api/backend.js` — fetch helper, JWT, error handling
- `src/auth/AuthContext.jsx` — user/session ปัจจุบัน
- `src/pages/admin/AdminUI.jsx` และไฟล์ CSS ของ Admin — ใช้ Design System เดิม
- `src/pages/admin/AdminInteractions.test.jsx` และ Tests ที่เกี่ยวข้อง

### 4.2 เปลี่ยน UI เฉพาะหน้า Admin > สมาชิก

- คง Layout, สี, Typography และ Responsive Style เดิม
- เพิ่มคอลัมน์ **จัดการ** หรือ action ในแต่ละแถว
- เพิ่มปุ่ม **เปลี่ยนสิทธิ์** เฉพาะสมาชิกที่แก้สิทธิ์ได้
- สำหรับบัญชี Admin ที่กำลังล็อกอิน ให้แสดง **บัญชีของคุณ** หรือปุ่ม disabled พร้อมเหตุผล
- เมื่อกด action ให้แสดง dialog ที่มี:
  - ชื่อและอีเมลบัญชีเป้าหมาย
  - Role ปัจจุบัน
  - ตัวเลือก `Member` / `Admin` (บันทึกค่า API เป็น `member` / `admin`)
  - คำเตือนชัดเจนหากกำลังเพิ่มสิทธิ์เป็น Admin
  - ปุ่ม **ยกเลิก** และ **ยืนยันการเปลี่ยนสิทธิ์**
- ห้ามส่งคำขอถ้าไม่ได้เปลี่ยน Role หรือกำลังบันทึก
- หลัง API สำเร็จ แสดง Role ใหม่ในรายการทันทีโดย update state หรือ re-fetch (`GET /api/admin/members`) และแสดง Success feedback
- ถ้า API ตอบ error ให้แสดงข้อความใน UI อย่างเหมาะสม ไม่เปลี่ยน Role บนจอเป็นค่าที่ไม่ได้บันทึกสำเร็จ
- รองรับการใช้งานบนมือถือและ keyboard เบื้องต้น: dialog มี title/label, focus management, ปิดด้วยยกเลิกหรือ Escape ตามแนวทางที่ใช้ในโครงการ

### 4.3 เชื่อม API จริง

เพิ่มเมธอดใน `adminService` โดยคงรูปแบบ request helper เดิม เช่น:

```js
updateMemberRole: (id, role) =>
  request(`/api/admin/members/${encodeURIComponent(id)}/role`, {
    method: 'PATCH',
    body: { role },
  }),
```

- ห้ามสร้างเส้นทาง API ไม่ตรงกับ Backend
- ต้องใช้งานร่วมกับ Auth/JWT ปัจจุบันของโปรเจกต์
- Production ต้องเรียก Go API จริง ไม่ใช้ mock แทน API ใหม่
- ถ้า BE ยังไม่ได้ Deploy ให้พัฒนา UI/Test ด้วย mock เฉพาะชุดทดสอบ และระบุว่า **ยังไม่ทดสอบกับ Production API**

### 4.4 Frontend Tests

- Admin มองเห็นและเปิด dialog ได้
- เลือก role ใหม่, ยืนยัน, เรียก `PATCH` ด้วย path/body ที่ถูกต้อง
- เมื่อสำเร็จ Role แสดงผลใหม่ทันที
- ระหว่าง Loading ไม่ส่ง request ซ้ำ
- กรณี Error (เช่น 400/403/404/500) แสดงข้อผิดพลาดและไม่อัปเดตข้อมูลผิด ๆ
- บัญชี Admin ของตัวเองเปลี่ยนสิทธิ์ไม่ได้ใน UI
- Member ไม่สามารถเข้าหน้า Admin ตาม Route Guard เดิม
- ตรวจไม่ให้ Test เดิมเสีย

คำสั่งจากโฟลเดอร์ `thaifood/`:

```bash
CI=true npm test -- --watchAll=false --runInBand
CI=true npm run build
```

## 5. ลำดับทำงานเมื่อต้องใช้ Codex คนละเครื่อง

**เครื่อง/คนทำ BE**
1. Pull โค้ด Backend ล่าสุดและตรวจ API/Store/Auth เดิม
2. เพิ่ม endpoint + validation + authorization + SQL + tests
3. รัน Tests และ Build พร้อมสรุปผล
4. Commit และ Push เข้า Backend Repository ตาม workflow ปกติ
5. ส่ง API Contract, sample request/response และ Commit SHA ให้ฝั่ง FE
6. Deploy Backend หลังตรวจทุกอย่างและได้รับอนุญาตตาม workflow ของทีม

**เครื่อง/คนทำ FE**
1. Pull โค้ด Frontend ล่าสุด
2. นำ Contract ที่ BE ส่งมาใช้ (endpoint, body, response และ error)
3. เพิ่ม action/dialog + API service + tests
4. รัน Tests และ Build พร้อมสรุปผล
5. ทดสอบ UI กับ Backend ที่มี endpoint ใหม่แล้ว
6. Commit/Push และ Deploy Vercel หลังยืนยันการเชื่อมต่อจริง

> ไม่ต้องใช้ GitHub บัญชีแยกเสมอไป แต่ห้ามทำงานทับกันโดยไม่มีการ Pull/Push/Review ที่ชัดเจน และไม่ commit `.env`, database URL, JWT secret, token หรือ credential ลง Repository

## 6. เกณฑ์รับงาน (Acceptance Criteria)

- [ ] Admin เปิดหน้า **สมาชิก** และเห็นรายการสมาชิกตามเดิม
- [ ] Admin สามารถเลื่อน Member คนอื่นเป็น Admin ผ่าน UI ได้จริง
- [ ] Admin สามารถลดสิทธิ์ Admin คนอื่นเป็น Member ได้ตามกฎความปลอดภัย
- [ ] Role ใหม่ถูกเก็บใน PostgreSQL จริงและแสดงบน UI โดยไม่ refresh
- [ ] บัญชีที่ถูกเปลี่ยนสิทธิ์ถูก Backend บังคับใช้ Role ใหม่กับ request ถัดไป
- [ ] สมาชิกทั่วไป/Guest ไม่สามารถใช้ API เปลี่ยนสิทธิ์ได้
- [ ] ไม่สามารถเปลี่ยน Role ของ Admin ที่กำลังใช้งานเอง
- [ ] ไม่สามารถทำให้ระบบไม่มี Admin เหลืออยู่
- [ ] มี Loading, Success, Error และ Confirm Dialog ที่ใช้งานได้
- [ ] React Unit Tests และ Go Tests ที่เกี่ยวข้องผ่าน
- [ ] Build FE/BE ผ่าน และมีหลักฐานการทดสอบ API จริง
- [ ] ฟีเจอร์เดิม เช่น Register, Login, Profile, Restaurants, Favorites, Reviews, Admin CRUD ไม่เสีย
- [ ] อัปเดต API Contract/README ที่จำเป็นแล้ว

## 7. สิ่งที่ Codex ต้องรายงานเมื่อทำเสร็จ

1. รายชื่อไฟล์ที่แก้ไขใน Repository ที่รับผิดชอบ พร้อมคำอธิบายสั้น ๆ
2. API Contract ใหม่ และตัวอย่าง Request/Response ที่ **ตรงกับโค้ดที่ implement จริง**
3. ผลทดสอบที่ **รันจริง** (จำนวนผ่าน/ไม่ผ่าน/Skip) และผล Build
4. วิธีลองใช้งานทีละขั้นในหน้าเว็บ พร้อม test account แนวทางสร้างเอง (ไม่เผยรหัสผ่านจริง)
5. สิ่งที่ยังทำไม่ได้หรือยังไม่ได้ทดสอบ เช่น BE ยังไม่ deploy, ไม่ได้ทดสอบ Admin จริง
6. หากมี Deploy ให้รายงาน URL/Revision ที่ Deploy สำเร็จ; ห้ามอ้างว่า Deploy สำเร็จหากไม่ได้ตรวจ
7. ระบุความเปลี่ยนแปลงด้านความปลอดภัย เช่น การห้าม Self-demotion และตรวจ Role จากฐานข้อมูล

---

**คำสั่งสุดท้ายสำหรับ Codex:** ลงมือแก้ไฟล์จริงเฉพาะ Repository ที่เปิดอยู่ตามขอบเขตด้านบน ไม่เพียงอธิบายวิธีทำ ตรวจสอบผล Tests/Build ที่รันได้ และสรุปงานพร้อมข้อจำกัดอย่างตรงไปตรงมา
