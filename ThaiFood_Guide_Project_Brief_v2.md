# ThaiFood Guide
## เอกสารแบ่งงาน Backend (BE) และ Frontend (FE)

**สำหรับ Term Project วิชา BE-FE | React + Fetch API + Go + Supabase**  
**ฉบับปรับปรุง V2: เพิ่ม Repository ตัวอย่าง FE | วันที่ 11 พฤศจิกายน 2569**

> **เป้าหมาย:** สร้างเว็บแนะนำร้านอาหารในประเทศไทย โดยค้นหาและกรองตามจังหวัด/ประเภทอาหาร มี Guest, Member, Admin ใช้งานจริงผ่าน React → Go REST API → Supabase PostgreSQL พร้อม Deploy FE บน Vercel และ BE บน Render

---

## 1) ข้อตกลงสำคัญที่ทั้งสองคนต้องอ่าน

1. แยกทำงาน 2 เครื่อง / 2 บัญชี GitHub / 2 Repository ได้ โดยไม่ต้องรอโค้ดอีกฝ่ายเสร็จ
2. **BE เป็นเจ้าของ Database, Authentication, Authorization, Business Logic และ API**; FE เป็นเจ้าของ React UI, Form, State, Fetch, Routing และ React Unit Test
3. ทั้งสองฝ่ายใช้ **API Contract ฉบับเดียวกัน**: ชื่อ Endpoint, Method, JSON, HTTP Status, Authentication และ Error Format ห้ามเปลี่ยนเองโดยไม่แจ้งกัน
4. FE เริ่มพัฒนาด้วย **Mock Data** ที่มีรูปแบบ JSON ตรงกับ API Contract และเปลี่ยนเป็น URL จริงหลัง BE Deploy
5. ไม่ Commit `.env`, password, `DATABASE_URL`, `JWT_SECRET` หรือ API Key ลง GitHub; Commit เฉพาะ `.env.example`
6. Backend **ต้องตรวจสิทธิ์จริง** ไม่ใช้แค่ซ่อนปุ่มบน FE; ระบบที่มีสิทธิ์ Admin ต้องตรวจ JWT และ Role ทาง Server
7. ใช้ร้านอาหารจริงเฉพาะข้อมูลที่มีสิทธิ์นำมาใช้ และแสดงที่มาข้อมูลตามเงื่อนไข ไม่คัดลอกรูป/คะแนนจาก Google โดยพลการ
8. รองรับชื่อจังหวัดไทยทั้ง 77 จังหวัด แต่ชุดข้อมูลร้านอาหารเริ่มต้นอาจมีเฉพาะจังหวัดตัวอย่างก่อน เพื่อให้ระบบทำงานครบ

### 1.1 เทคโนโลยีที่ต้องใช้

| ส่วน | เทคโนโลยี | Deploy/ใช้งาน |
|---|---|---|
| Frontend | React 19 + Create React App (react-scripts) + React Router 6 + Tailwind CSS 3 + Fetch API | Vercel |
| Backend | Go + net/http + pgxpool + JWT/bcrypt | Render |
| Database | PostgreSQL บน Supabase | Supabase |
| API Testing | Postman | เครื่อง BE |
| React Unit Test | Jest (react-scripts test) + React Testing Library | เครื่อง FE |
| Version Control | Git + GitHub | Repo แยก BE/FE |

**โค้ดอ้างอิง BE:** https://github.com/Sukrit310/week6-lab-fork/tree/main/moviehub  
ใช้แนวทาง `auth/`, `handler/`, `store/`, `mw/` และ `main.go` จาก MovieHub เป็นตัวอย่าง **ไม่ใช่คัดลอกโมเดลหนังมาใช้ตรง ๆ**

**โค้ดอ้างอิง FE (เพิ่มใหม่):** https://github.com/670710721/moviehub-670710721  
โค้ด FE ต้นแบบที่ตรวจพบเป็น **Create React App (CRA)** ใช้ `react-scripts`, **ไม่ใช่ Vite**. มี `React 19`, `react-router-dom 6`, `Tailwind CSS 3`, `Jest` + `React Testing Library`, `src/api/backend.js`, `src/auth/AuthContext.jsx`, Route Guard, Admin pages และ tests. หากนำโครงสร้างนี้มาปรับ ให้คง toolchain CRA ในระยะแรก (หรือย้าย Vite อย่างชัดเจนเป็นงานแยก)

### 1.2 แนวทางนำ Repository FE มาอ้างอิง (ตรวจโค้ดแล้ว)

| จาก FE MovieHub | ปรับเป็น ThaiFood Guide | สิ่งที่ควรทำ |
|---|---|---|
| `src/App.jsx` | `src/App.jsx` | แก้ Routes `/movies` → `/restaurants` และ Admin routes |
| `src/pages/Movies.jsx` | `src/pages/Restaurants.jsx` | ค้นหา/กรองร้านตามจังหวัดและประเภท; API Pagination |
| `src/pages/MovieDetail.jsx` | `src/pages/RestaurantDetail.jsx` | รายละเอียดร้าน, รีวิว, ให้คะแนน, ร้านโปรด |
| `src/components/MovieCard.jsx`, `MovieGrid.jsx` | `RestaurantCard.jsx`, `RestaurantGrid.jsx` | แสดงรูปอาหาร ชื่อร้าน ที่ตั้ง คะแนน |
| `src/api/backend.js` | `src/api/backend.js` | ปรับ functions ให้ใช้ endpoint `/api/restaurants` และ API Contract ฉบับนี้ |
| `src/auth/AuthContext.jsx` | `src/auth/AuthContext.jsx` | เปลี่ยนชื่อ localStorage จาก moviehub, ปรับ profile/password |
| `src/components/ProtectedRoute.jsx`, `RequireRole.jsx` | ใช้แนวทางเดิม | Role Guard เป็นเพียง UI; BE ตรวจ JWT และ Role จริง |
| `src/pages/Wishlist.jsx` | `src/pages/Favorites.jsx` | ร้านโปรด Member |
| `src/pages/admin/AdminMovies.jsx` และ `AdminMovieForm.jsx` | `AdminRestaurants.jsx` และ `AdminRestaurantForm.jsx` | CRUD ร้านอาหาร |
| `src/components/ReviewForm.jsx`, `ReviewList.jsx` | ปรับใช้กับร้าน | รีวิวพร้อม rating 1–5 |
| `src/**/*.test.jsx` และ `src/api/backend.test.js` | React tests ของร้าน | ใช้ Jest + React Testing Library ต่อจากของเดิม |
| `src/api/tmdb.js`, `mock/movies.json`, `src/data/data.js` | ไม่ใช่ข้อมูลจริงของ ThaiFood | ไม่พึ่ง TMDB; แทนด้วย mock restaurants ที่ตรง Contract |

**ระวังโค้ดต้นแบบยังมี TODO, ข้อความที่เกี่ยวกับ MovieHub/TMDB และฟังก์ชัน API ที่ไม่ตรงกับ Contract ร้านอาหาร** จึงห้ามคัดลอกไปแล้วถือว่าใช้งานได้ทันที ต้องแก้เป็นร้านอาหารและทดสอบ End-to-End จริง


## 2) ฟังก์ชันและสิทธิ์ผู้ใช้

| ฟังก์ชัน | Guest | Member | Admin |
|---|:---:|:---:|:---:|
| ค้นหา/ดูรายการร้านและรายละเอียด | ✓ | ✓ | ✓ |
| ดูรีวิวและคะแนนจากสมาชิกเว็บ | ✓ | ✓ | ✓ |
| สมัครสมาชิก / เข้าสู่ระบบ | ✓ | ✓ | ✓ |
| ดู/แก้ไข Profile และเปลี่ยนรหัสผ่าน | - | ✓ | ✓ |
| เพิ่มรีวิว/ให้คะแนนร้าน (1–5) | - | ✓ | ✓ |
| ลบรีวิวของตัวเอง | - | ✓ | ✓ |
| บันทึก/ลบร้านโปรด | - | ✓ | ✓ |
| เพิ่ม/แก้ไข/ลบร้านอาหาร | - | - | ✓ |
| จัดการรีวิว/ดูรายชื่อสมาชิก | - | - | ✓ |

**ขอบเขตที่ต้องทำ:** ระบบสมาชิกครบ Register, Login, Logout, Profile, Change Password; CRUD ร้านโดย Admin; ร้านและรีวิวแสดงข้อมูลจริงจาก Database; Loading/Empty/Error ครบ; React Unit Test ผ่าน

**ถ้ามีเวลา:** แผนที่, แนะนำร้านยอดนิยม, ฟิลเตอร์ขั้นสูง  
**ยังไม่ต้องทำ:** จองโต๊ะ, แชต, AI Recommendation, Payment

<div class="pagebreak"></div>

## 3) รายละเอียดงาน Backend (ผู้รับผิดชอบคนที่ 1)

### BE-01 โครงสร้างโปรเจกต์

- สร้าง Go Module ชื่อโปรเจกต์ใหม่ และจัดแพ็กเกจ `auth`, `handler`, `store`, `mw`, `cmd` หรือ `main.go`
- ใช้ Environment Variables: `DATABASE_URL`, `JWT_SECRET`, `ALLOWED_ORIGINS`, `PORT`
- เชื่อม Supabase PostgreSQL ด้วย `pgxpool` และทดสอบการเชื่อมต่อ
- มี JSON response/error แบบเดียวกันทุก Endpoint, validation และ logging ที่ไม่เปิดเผยข้อมูลลับ

### BE-02 Database Schema (PostgreSQL / Supabase)

| ตาราง | คอลัมน์หลัก | กติกา/ความสัมพันธ์ |
|---|---|---|
| `provinces` | `id`, `name_th`, `region` | ชื่อจังหวัดไม่ซ้ำ มี 77 จังหวัด |
| `categories` | `id`, `name_th` | เช่น อาหารเหนือ, อีสาน, คาเฟ่ |
| `restaurants` | `id`, `name`, `province_id`, `category_id`, `address`, `description`, `latitude`, `longitude`, `image_url`, `created_at`, `updated_at` | FK ไป provinces/categories, ชื่อร้านไม่ว่าง |
| `members` | `id`, `email`, `password_hash`, `display_name`, `role`, `created_at`, `updated_at` | email unique แบบไม่สนตัวพิมพ์, role member/admin |
| `reviews` | `id`, `restaurant_id`, `member_id`, `rating`, `text`, `created_at`, `updated_at` | rating 1–5, 1 สมาชิกรีวิว 1 ร้านได้ 1 รีวิว (unique) |
| `favorites` | `member_id`, `restaurant_id`, `created_at` | composite PK กันบันทึกซ้ำ |

- ทำ Foreign Keys และ `ON DELETE` ให้เหมาะสม เช่น ลบร้านแล้วลบรีวิว/รายการโปรดที่อ้างถึงร้าน
- ทำ Index ที่ค้นหาบ่อย เช่น `restaurants(province_id, category_id)`, `reviews(restaurant_id)`, `favorites(member_id)`
- **คะแนนเฉลี่ย (rating) และจำนวนรีวิว (reviewCount) ต้องคำนวณจาก reviews ของเว็บเอง**; ร้านยังไม่มีรีวิวให้ใช้ `rating: null` และ `reviewCount: 0`
- มี `schema.sql` และ `seed.sql` หรือสคริปต์ Import ที่รันซ้ำโดยไม่สร้างข้อมูลซ้ำ

### BE-03 Import ข้อมูลร้านอาหาร

- เตรียมชุดข้อมูลเริ่มต้นจากแหล่งที่ได้รับอนุญาต เช่น ททท. Open Data / OpenStreetMap หรือข้อมูลที่ผู้จัดทำกรอกเอง
- Mapping จังหวัด, ประเภทอาหาร, ชื่อร้าน, ที่อยู่, พิกัด และภาพที่มีสิทธิ์ใช้งาน
- ตรวจข้อมูลซ้ำและบันทึก `source_name` / `source_url` หากเหมาะสม เพื่อแสดงที่มา
- แนะนำให้เริ่ม 3–5 จังหวัดตัวอย่าง (จังหวัดละ 5–10 ร้าน) แล้วค่อยเพิ่ม; UI ต้องรองรับจังหวัดทั้งหมด

### BE-04 ระบบสมาชิกและความปลอดภัย

- `POST /api/auth/register`: ตรวจฟิลด์, hash password ด้วย bcrypt, ป้องกันอีเมลซ้ำ
- `POST /api/auth/login`: ตรวจรหัสผ่าน, ส่ง Bearer JWT พร้อมข้อมูลสมาชิกที่ไม่รวม hash
- Logout แบบพื้นฐาน: FE ลบ token ฝั่ง Client (JWT ยังใช้ได้จนหมดอายุ); ถ้าต้องยกเลิกทันทีต้องเพิ่มระบบ revoke/denylist ฝั่ง BE
- `GET/PATCH /api/me`: อ่านและแก้ไขเฉพาะข้อมูลของผู้ใช้ที่ login
- `PUT /api/me/password`: ตรวจรหัสผ่านเดิมก่อนเปลี่ยน และ hash ใหม่
- JWT middleware ตอบ `401` ถ้าไม่ได้ล็อกอิน, role middleware ตอบ `403` หากสิทธิ์ไม่ถึง
- Admin CRUD ต้องยืนยัน role ฝั่ง Server; Member ลบได้เฉพาะรีวิวของตนเอง

### BE-05 API ร้านอาหาร/จังหวัด/ประเภท

- ดูจังหวัดและประเภทอาหาร
- ค้นหาและแบ่งหน้าร้านอาหารตาม `province`, `category`, `q`, `page`, `limit`
- ดูรายละเอียดร้านและรีวิว (รวมคะแนนเฉลี่ยจากระบบเว็บ)
- กำหนด Validation, 404, Limit/Pagination และ Error เมื่อ Query ผิด

### BE-06 API รีวิว/ร้านโปรด/Admin

- เพิ่มรีวิวพร้อม rating 1–5; กันรีวิวซ้ำหรืออนุญาตให้แก้ไขตาม Contract เดียวกัน
- ลบรีวิวได้เฉพาะเจ้าของหรือ Admin
- Member เพิ่ม/ลบ/ดู Favorites ของตัวเอง
- Admin เพิ่ม แก้ไข ลบร้าน และจัดการรีวิว; อาจให้ดูรายชื่อสมาชิกแบบไม่เผย Password Hash

### BE-07 ทดสอบและ Deploy

- Postman: ทดสอบ Happy Path และ 400/401/403/404/409/500 ที่เกี่ยวข้อง
- ทดสอบ SQL Injection, token ผิด/หมดอายุ, สิทธิ์ Member เข้าหน้า Admin, รีวิวร้านที่ไม่มีอยู่
- Deploy Render ให้ใช้ `PORT` จาก Environment และเชื่อม Supabase ได้
- ตั้ง `ALLOWED_ORIGINS` ให้มี URL localhost FE และ Vercel ที่ใช้จริง
- ส่งมอบ `API_CONTRACT.md`, `README.md`, `.env.example`, `schema.sql`, seed/import และ URL Render

**Definition of Done ของ BE:** เรียก API จริงบน Render ผ่าน Postman ได้, Auth/Role ถูกต้อง, ข้อมูลอยู่ใน Supabase, Endpoint ตรง Contract, ไม่มี Secret ใน Repo

<div class="pagebreak"></div>

## 4) รายละเอียดงาน Frontend (ผู้รับผิดชอบคนที่ 2)

### FE-01 ตั้งค่าและออกแบบระบบ

- ใช้ Repository `670710721/moviehub-670710721` เป็น **ต้นแบบโครงสร้าง FE**; สร้างโปรเจกต์ ThaiFood ใหม่โดยดัดแปลง CRA + React Router + Tailwind CSS จากต้นแบบ (หากย้ายไป Vite ให้ระบุแผนย้ายแยกต่างหาก)
- ออกแบบ Responsive สำหรับ Desktop/มือถือ; มี Navbar, Footer, Search, Filters, Cards, Forms
- กำหนด Theme สี, ฟอนต์, Empty State, Error Message ให้สม่ำเสมอ
- ใช้ `REACT_APP_API_URL` (เพราะต้นแบบเป็น CRA) ใน `.env` และ `.env.example`; ห้ามเก็บ Secret ไว้ในตัวแปรฝั่ง FE ที่ถูกฝังใน Bundle

### FE-02 หน้าที่ต้องทำ

| Page/Route ตัวอย่าง | เนื้อหา | สิทธิ์ |
|---|---|---|
| `/` | Home / จังหวัดยอดนิยม / ร้านแนะนำจากฐานข้อมูล | ทุกคน |
| `/restaurants` | ค้นหา, กรองจังหวัด/ประเภท, Pagination, Cards | ทุกคน |
| `/restaurants/:id` | ภาพ/รายละเอียดร้าน/คะแนน/รีวิว/ปุ่มโปรด | ทุกคน; เขียนรีวิวต้อง Login |
| `/login` และ `/register` | ฟอร์ม Authentication | Guest |
| `/profile` | แก้ไขข้อมูลส่วนตัว | Member |
| `/change-password` | เปลี่ยนรหัสผ่าน | Member |
| `/favorites` | ดูและลบร้านโปรด | Member |
| `/my-reviews` | ดูและลบรีวิวของตัวเอง | Member (เสริมถ้ามี API) |
| `/admin` | Dashboard | Admin |
| `/admin/restaurants` | Table + Form เพิ่ม/แก้ไข/ลบร้าน | Admin |
| `/admin/reviews` | จัดการรีวิว | Admin |
| `/admin/members` | ดูรายชื่อสมาชิกแบบจำกัดข้อมูล | Admin |
| `*` | Not Found (404) | ทุกคน |

### FE-03 State / Forms / Fetch API

- ใช้ Controlled Forms (`value`, `onChange`) พร้อม Validation และข้อความเตือน
- ดัดแปลง `src/api/backend.js` ของต้นแบบสำหรับ `fetch()` ผ่าน URL จาก `process.env.REACT_APP_API_URL`; ทำ service functions ให้ตรงกับ Contract
- แยก `authService`, `restaurantService`, `reviewService`, `favoriteService`, `adminService`
- แสดง **Loading / Empty / Error** ทุกจุดที่ดึงข้อมูล เช่น ร้าน, รีวิว, ร้านโปรด, Dashboard
- React state ต้องไม่แก้ object/array เดิมตรง ๆ; ควบคุม Submit Loading และป้องกันกดซ้ำ
- Protected Route / Role UI มีไว้ให้ใช้งานง่าย แต่ **BE ต้องตรวจสิทธิ์จริงอีกชั้น**
- ระหว่าง BE ยังไม่เสร็จ ใช้ Mock JSON ตาม API Contract; เปลี่ยนไปใช้ API จริงก่อนส่งงาน

### FE-04 Unit Test (จำเป็นตาม Rubric)

- ใช้ Jest + React Testing Library (`react-scripts test`) ตามต้นแบบ และ Mock `fetch()`/services
- ตัวอย่างอย่างน้อย 5 Test:
  1. Restaurant List แสดงรายการจาก API ได้
  2. ไม่พบร้าน แสดง Empty State
  3. API Error แสดงข้อความผิดพลาดและปุ่ม Retry
  4. Register/Login ฟอร์มไม่ส่งเมื่อกรอกไม่ครบ
  5. Guest ถูกกันออกจากหน้าที่ต้อง Login; Member ถูกกันออกจาก Admin UI
- ตรวจ Assertions สำคัญจริงและ Test ต้องผ่านใน CI หรือ Local

### FE-05 Deploy และตรวจรับ

- Deploy Vercel ตั้ง `REACT_APP_API_URL` เป็น Render URL จริง (CRA จำเป็นต้อง Rebuild/ Redeploy เมื่อแก้ตัวแปร)
- ทดสอบเส้นทาง Guest → Register → Login → Profile → Favorite → Review → Logout
- ตรวจ Responsive, ลิงก์, CORS, Error State, Form Validation, Browser Console
- ส่งมอบ `README.md`, `.env.example`, Unit Tests และ URL Vercel

**Definition of Done ของ FE:** หน้าเว็บทำงานได้จริงกับ Go API / Supabase, ครบ 3 Role, Loading/Empty/Error, Test ผ่าน, Deploy Vercel สำเร็จ

<div class="pagebreak"></div>

## 5) API Contract เริ่มต้นที่ต้องใช้ร่วมกัน

**Base URL เมื่อพัฒนา Local:** `http://localhost:8080` (ใช้ได้เมื่อ BE อยู่เครื่องเดียวกับผู้เรียก)  
**Base URL เมื่อแยกเครื่อง:** URL Render จริง เช่น `https://YOUR-BACKEND.onrender.com` (ตัวอย่าง Placeholder ไม่ใช่ URL ที่สร้างแล้ว)  
**Auth Header:** `Authorization: Bearer <JWT_TOKEN>`  
**Content-Type:** `application/json`

| Method | Endpoint | Auth | หน้าที่ |
|---|---|---|---|
| GET | `/api/provinces` | Public | รายชื่อจังหวัด |
| GET | `/api/categories` | Public | ประเภทอาหาร |
| GET | `/api/restaurants` | Public | ค้นหา/กรอง/แบ่งหน้าร้าน |
| GET | `/api/restaurants/{id}` | Public | รายละเอียดร้าน |
| GET | `/api/restaurants/{id}/reviews` | Public | รีวิวของร้าน |
| POST | `/api/auth/register` | Public | สมัครสมาชิก |
| POST | `/api/auth/login` | Public | Login รับ JWT |
| GET | `/api/me` | Member | ดูโปรไฟล์ |
| PATCH | `/api/me` | Member | แก้ไขโปรไฟล์ |
| PUT | `/api/me/password` | Member | เปลี่ยนรหัสผ่าน |
| GET | `/api/favorites` | Member | ร้านโปรดของตัวเอง |
| POST | `/api/favorites` | Member | เพิ่มร้านโปรด |
| DELETE | `/api/favorites/{restaurantId}` | Member | ลบร้านโปรด |
| POST | `/api/restaurants/{id}/reviews` | Member | เขียนรีวิว/ให้คะแนน |
| DELETE | `/api/reviews/{id}` | Owner/Admin | ลบรีวิว |
| GET | `/api/admin/members` | Admin | ดูรายชื่อสมาชิก |
| POST | `/api/admin/restaurants` | Admin | เพิ่มร้าน |
| PUT | `/api/admin/restaurants/{id}` | Admin | แก้ไขร้าน |
| DELETE | `/api/admin/restaurants/{id}` | Admin | ลบร้าน |
| DELETE | `/api/admin/reviews/{id}` | Admin | Moderation รีวิว |

**Logout:** สำหรับ MVP ให้ FE ลบ JWT และข้อมูล session ฝั่ง Browser พร้อม Redirect ไป Login/Home (ไม่ต้องมี Logout Endpoint หากไม่มี Server-side Revocation); ระบุข้อจำกัดนี้ใน README

### 5.1 Request/Response ที่ตกลงร่วมกัน

**GET `/api/restaurants?province=เชียงใหม่&category=อาหารเหนือ&q=ข้าวซอย&page=1&limit=12`**

```json
{
  "items": [
    {
      "id": 101,
      "name": "ร้านอาหารตัวอย่าง",
      "province": "เชียงใหม่",
      "category": "อาหารเหนือ",
      "address": "อำเภอเมืองเชียงใหม่",
      "description": "ข้อมูลร้านอาหารตัวอย่าง",
      "latitude": 18.788,
      "longitude": 98.985,
      "imageUrl": null,
      "rating": 4.5,
      "reviewCount": 2
    }
  ],
  "page": 1,
  "limit": 12,
  "total": 1
}
```

> JSON นี้เป็น **ข้อมูลสมมติสำหรับทดสอบ Contract** ไม่ใช่ข้อมูลร้านอาหารจริง และ `rating` คำนวณจาก Review ในระบบ

**POST `/api/auth/login` Request:** `{"email":"demo@example.com","password":"ExamplePassword123"}`  
**Response 200:** `{"token":"<JWT>","member":{"id":1,"displayName":"Demo","email":"demo@example.com","role":"member"}}`  
**POST `/api/restaurants/{id}/reviews` Request:** `{"rating":5,"text":"อาหารอร่อย บริการดี"}`  
**Error 400/401/403/404/409/500:** `{"error":{"code":"ERROR_CODE","message":"ข้อความอธิบาย"}}`

**Status ที่ต้องใช้ให้เหมาะสม:** `200` OK, `201` Created, `204` No Content, `400` Bad Request, `401` Unauthorized, `403` Forbidden, `404` Not Found, `409` Conflict, `500` Internal Error

### 5.2 ประเด็นที่ต้องยืนยันก่อน Coding

- ฟิลด์ใน `Restaurant` และ `Member` ใช้ camelCase ใน JSON (เช่น `imageUrl`, `displayName`, `reviewCount`)
- Rating ใช้ **1–5** จากรีวิวผู้ใช้ภายในเว็บเท่านั้น; เมื่อไม่มีรีวิว `rating: null`
- `GET /api/restaurants` มี `items`, `page`, `limit`, `total` เสมอ แม้ไม่มีข้อมูล
- เลือกฟิลเตอร์จังหวัดและประเภทโดยใช้ **ชื่อไทย** ใน Query ก่อน (BE ต้อง map กับ DB ให้ถูก)
- แก้ไขรีวิว: **ไม่อยู่ในขอบเขต MVP** แต่สามารถเพิ่ม PATCH ภายหลังเมื่อทั้งสองฝ่ายตกลง
- บัญชี Admin ไม่ให้สมัครได้เองจากหน้า Register; ตั้ง role ผ่านกระบวนการ Admin/DB ที่ควบคุมได้เท่านั้น

<div class="pagebreak"></div>

## 6) การเชื่อมงานสองเครื่อง / สอง Codex

### 6.1 GitHub

- Repo A: `thaifood-backend` (Go, Database schema, API_CONTRACT.md)
- Repo B: `thaifood-frontend` (React, Tests, README)
- ใช้ GitHub คนละบัญชีได้; เจ้าของ Repo เพิ่มอีกคนเป็น Collaborator ตามความจำเป็น
- หลีกเลี่ยง Push ตรงเข้า `main` เมื่อทำงานพร้อมกัน; ใช้ Branch, Pull Request และ Review
- **เอกสาร Contract มี Source of Truth ที่เดียว**: Backend Repo; FE ดึง/คัดลอก version เดียวกันก่อนพัฒนา

### 6.2 ขั้นตอนเชื่อม

1. BE เขียน `schema.sql` และ API Contract ตามเอกสารนี้ พร้อมส่ง Mock JSON ให้ FE
2. FE ทำหน้าจอและ Services จาก Mock JSON ไม่ต้องรอ Backend เสร็จ
3. BE ทดสอบทุก Endpoint ด้วย Postman และ Deploy บน Render
4. FE ใส่ Render URL ใน `.env` เช่น `REACT_APP_API_URL=https://YOUR-BACKEND.onrender.com`
5. BE ตั้ง `ALLOWED_ORIGINS` ให้รองรับ `http://localhost:3000` (CRA) และโดเมน Vercel จริง
6. FE เปลี่ยน Mock → Fetch API จริง ทดสอบร่วมกันทุก Role
7. BE/FE ร่วมกันตรวจ Network (Request, Response, Authorization), CORS, 400/401/403/404 และ Database

**สำคัญ:** `localhost` ของ FE ที่อยู่คนละเครื่องไม่ได้หมายถึงเครื่อง BE ดังนั้นใช้ Render URL หรือ IP ภายใน LAN ที่กำหนดให้เข้าถึงได้

### 6.3 ตัวอย่าง React Fetch

```javascript
const API_URL = process.env.REACT_APP_API_URL;

export async function getRestaurants(params = {}) {
  const query = new URLSearchParams(params).toString();
  const response = await fetch(`${API_URL}/api/restaurants?${query}`);
  if (!response.ok) throw new Error(`HTTP ${response.status}`);
  return response.json();
}
```

### 6.4 Environment Variables

**BE (`.env.example`):**

```dotenv
DATABASE_URL=postgresql://USER:PASSWORD@HOST:5432/DBNAME
JWT_SECRET=replace-with-long-random-secret-at-least-32-characters
ALLOWED_ORIGINS=http://localhost:3000,https://YOUR-FRONTEND.vercel.app
PORT=8080
```

**FE (`.env.example`):**

```dotenv
REACT_APP_API_URL=https://YOUR-BACKEND.onrender.com
```

**หมายเหตุ:** ค่าด้านบนเป็น Placeholder ห้ามใช้ Password/Secret ตัวอย่างกับระบบจริง และห้าม Commit `.env`

<div class="pagebreak"></div>

## 7) แผนพัฒนาและ Checklist ตรวจรับร่วมกัน

| ระยะ | Backend (BE) | Frontend (FE) | ผลลัพธ์ |
|---|---|---|---|
| 1. ตกลงงาน | Schema + API Contract | Wireframe + Pages + JSON Contract | ตกลงฟิลด์เดียวกัน |
| 2. ทำระบบแกนหลัก | DB + Auth + Restaurant API | Layout + Home + List + Login (Mock) | ทำพร้อมกันได้ |
| 3. ฟีเจอร์ครบ | Review/Favorite/Admin/Role | Detail/Review/Favorite/Admin/Forms | ทั้งสองฝั่งพร้อมเชื่อม |
| 4. Integration | Render, Postman, CORS | Fetch จริง, Loading/Empty/Error | ทำงานกับ DB จริง |
| 5. ตรวจรับ | Security + API Tests + Docs | React Unit Tests + Responsive + Vercel | พร้อมสาธิตและส่งงาน |

### Checklist ระบบหลักก่อนส่ง

- [ ] ทั้ง FE และ BE มี Repo/README/.env.example และโครงสร้างไฟล์ชัดเจน
- [ ] Database บน Supabase มีตารางและข้อมูล Seed ที่ทดสอบได้
- [ ] Guest ดูรายการ/กรองตามจังหวัด/ดูรายละเอียดและรีวิวได้
- [ ] Register/Login/Logout/Profile/Change Password ใช้งานได้จริง
- [ ] Member เพิ่มรีวิวและคะแนน / ลบรีวิวตัวเอง / เก็บ Favorites ได้
- [ ] Admin เพิ่ม/แก้ไข/ลบร้านและจัดการรีวิวได้ (BE ตรวจ Role)
- [ ] FE เรียก Go API ผ่าน Fetch จริงและจัดการ Loading/Empty/Error
- [ ] React Unit Tests มี Assertions ถูกต้องและผ่าน
- [ ] Render / Vercel / Supabase เชื่อมจริง; CORS และ Environment Variables ถูกต้อง
- [ ] ทดสอบ 401/403/404 และข้อมูลผิดพลาดอย่างน้อยตามจุดสำคัญ
- [ ] ไม่มีรหัสผ่าน/Token/Secret จริงใน GitHub
- [ ] สาธิตระบบและอธิบายโค้ดของแต่ละคนได้

## 8) เกณฑ์อาจารย์และแผนคลิป

อ้างอิงจากภาพเกณฑ์ที่ผู้ทำโปรเจกต์ส่งมา: **คะแนนกลุ่ม 13 คะแนน + คลิป 7 คะแนน = 20 คะแนนต่อคน**

| เกณฑ์ผลงานกลุ่ม | คะแนน | จุดที่ต้องมี |
|---|---:|---|
| Tech Stack / Deploy | 2 | React, Go, Vercel, Render, Supabase |
| ระบบสมาชิก | 2 | Register/Login/Logout/Profile/Change Password |
| Guest / Member / Admin | 4 | ฟังก์ชันและ Backend Authorization |
| API / Database / UI | 2 | FE Fetch → Go → Supabase + Loading/Empty/Error |
| React Unit Test | 1 | Test สำคัญมี Assertions ถูกและผ่าน |
| Project Structure | 2 | แยก FE/BE, ชัดเจน, แยกค่าตั้งต้น |

**คลิปไม่เกิน 15 นาทีต่อทีม ทุกคนต้องอธิบายโค้ดที่ตนรับผิดชอบ**: สาธิตจริง 2 คะแนน + อธิบายโค้ด/เส้นทางข้อมูลรายบุคคล 3 คะแนน + เหตุผลและปัญหา/การแก้ไขรายบุคคล 2 คะแนน

**ตัวอย่างแผนคลิป 14 นาที:** Demo 2 นาที → BE อธิบาย 5 นาที → FE อธิบาย 5 นาที → ปัญหาและสิ่งที่แก้ 2 นาที; ปรับให้ทุกคนได้อธิบายตามงานจริง และระบุชื่อผู้พูด

<div class="pagebreak"></div>

## 9) คำสั่งสำหรับ Codex - เครื่อง Backend

คัดลอกส่วนนี้ส่งให้ Codex ที่เปิด Repository `thaifood-backend`:

> คุณเป็น Backend Developer ของโปรเจกต์ ThaiFood Guide เว็บแนะนำร้านอาหารไทยแยกตามจังหวัด ให้อ่านเอกสาร ThaiFood_Guide_Project_Brief_v2.md ทั้งหมด และใช้ https://github.com/Sukrit310/week6-lab-fork/tree/main/moviehub เป็นตัวอย่างสถาปัตยกรรม Go เท่านั้น สร้าง Backend Go + net/http + pgxpool + PostgreSQL Supabase + bcrypt + JWT แยกแพ็กเกจ auth/handler/store/mw; ทำ schema.sql, seed/import, API ทั้งหมดตาม Contract, Validation, Error JSON, Authorization Admin/Member/Owner, Pagination/Search, Postman Collection, README และ .env.example; ทดสอบการทำงานก่อนสรุป ห้าม hardcode Secrets ห้ามสร้าง API หรือเปลี่ยนชื่อฟิลด์โดยไม่อัปเดต API_CONTRACT.md; ให้แสดงรายการไฟล์ที่สร้างและวิธี Run/Test/Deploy Render อย่างชัดเจน

**ลำดับให้ Codex ทำ:** (1) ออกแบบ Schema และ Contract → (2) สร้างโครงสร้าง Go → (3) Auth/Role → (4) GET Restaurant/Search → (5) Review/Favorite/Admin → (6) Tests/Postman → (7) Deploy/README

## 10) คำสั่งสำหรับ Codex - เครื่อง Frontend

คัดลอกส่วนนี้ส่งให้ Codex ที่เปิด Repository `thaifood-frontend`:

> คุณเป็น Frontend Developer ของโปรเจกต์ ThaiFood Guide เว็บแนะนำร้านอาหารไทยแยกจังหวัด อ่าน ThaiFood_Guide_Project_Brief_v2.md และ ThaiFood_Guide_API_CONTRACT_v2.md ก่อนเริ่มงาน ใช้ https://github.com/670710721/moviehub-670710721 เป็น **ต้นแบบ Frontend** โดยตรวจ `package.json`, `src/App.jsx`, `src/api/backend.js`, `src/auth/AuthContext.jsx`, `src/components/ProtectedRoute.jsx`, `src/components/RequireRole.jsx`, หน้า Home/Movies/MovieDetail/Wishlist/Admin และ tests. ต้นแบบใช้ React + Create React App (react-scripts) + Tailwind + Jest ไม่ใช่ Vite: ให้คง CRA เมื่อเริ่มดัดแปลงและใช้ `process.env.REACT_APP_API_URL`, รัน `npm start`, `npm test`, `npm run build`. สร้างโปรเจกต์ ThaiFood ใหม่ ไม่แก้ต้นฉบับ MovieHub โดยตรง ปรับ Movie -> Restaurant, Movies -> Restaurants, Wishlist -> Favorites, Genre -> Category, คะแนนเป็น 1–5 และเพิ่ม Province filter; เปลี่ยนทุก endpoint `/api/movies` เป็น `/api/restaurants` ให้ตรง Contract, ไม่ใช้ TMDB จริง. หน้าเว็บต้องมี Home, Restaurant List/Detail, Register/Login/Logout, Profile/Change Password, Favorites, Reviews, Admin Dashboard/Admin Restaurants/Admin Reviews/Admin Members, 404; รองรับ Guest/Member/Admin ผ่าน UI Route Guard (BE ตรวจสิทธิ์จริง), Forms/Validation, Loading/Empty/Error/Retry, Responsive. ระหว่าง BE ยังไม่เสร็จ ใช้ mock restaurants JSON ตาม Contract แต่ก่อนส่งต้องเปลี่ยนเป็น fetch Go API จริง; เขียน React Unit Tests ด้วย Jest + React Testing Library โดย mock API; เตรียม README, .env.example และ Deploy Vercel ห้าม Commit Secrets หรือแก้ Contract ฝ่ายเดียว และสรุปไฟล์ที่เปลี่ยนพร้อมผลทดสอบ

**ลำดับให้ Codex ทำ:** (1) Wireframe/Routes + Mock JSON → (2) Layout/Home/List/Detail → (3) Auth/Profile → (4) Review/Favorite/Admin → (5) Fetch จริง/สถานะต่าง ๆ → (6) Unit Tests → (7) Deploy/README

---

## 11) รายการสิ่งที่ทั้งสองฝ่ายต้องส่งสุดท้าย

**BE:** GitHub Repo, Render API URL, Supabase Schema+Seed/Import, API Contract, Postman Collection, README, ตัวอย่าง Response และผลทดสอบ  
**FE:** GitHub Repo, Vercel URL, หน้าจอครบ Guest/Member/Admin, React Unit Tests, README และหลักฐานว่าดึงข้อมูล API จริง  
**ร่วมกัน:** Integration Demo, รายการปัญหา/วิธีแก้, คลิปอธิบายโค้ดแต่ละคนไม่เกิน 15 นาที

**ข้อย้ำ:** เอกสารนี้คือข้อตกลงตั้งต้นของทีม ให้ปรับตามคำสั่งอาจารย์หรือ API Contract เวอร์ชันล่าสุดที่ทีมตกลงกัน โดยต้องแจ้งทั้ง BE และ FE ทุกครั้งที่แก้รูปแบบข้อมูล
