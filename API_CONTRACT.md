# ThaiFood Guide - API Contract v2 (FE reference: MovieHub CRA)

เอกสารนี้ใช้ร่วมกันระหว่าง Codex เครื่อง Backend และเครื่อง Frontend ห้ามเปลี่ยนชื่อ Endpoint หรือ JSON Fields โดยไม่ตกลงกันก่อน

## FE Reference และข้อสังเกตเรื่อง Toolchain

- **Frontend reference:** https://github.com/670710721/moviehub-670710721
- **Backend reference:** https://github.com/Sukrit310/week6-lab-fork/tree/main/moviehub
- FE reference ใช้ Create React App (`react-scripts`) + React Router 6 + Tailwind CSS + Jest/React Testing Library; **ไม่ใช่ Vite**
- FE ให้ใช้ `process.env.REACT_APP_API_URL` และ `.env.example` ที่มี `REACT_APP_API_URL=https://YOUR-BACKEND.onrender.com` (อย่าใช้ `import.meta.env.VITE_API_URL` ถ้ายังเป็น CRA)
- FE dev server โดยปกติ `http://localhost:3000` ซึ่ง BE ควรเพิ่มใน `ALLOWED_ORIGINS`; URL ตัวอย่างไม่ใช่ Deploy จริง
- `src/api/backend.js` ใน MovieHub เป็นตัวอย่าง fetch/401/403/204/Authorization เท่านั้น **รายการ endpoints ของ ThaiFood ในเอกสารนี้เป็น Source of Truth**
- Auth storage keys, movie-specific types, TMDB, Wishlist, admin movie endpoints, score scale ต้องปรับเป็นร้านอาหาร และไม่แอบเปลี่ยนชื่อ JSON Field ฝ่ายเดียว

## 1. Base URL / Authentication

- Local (เมื่อเรียกจากเครื่องเดียวกับ BE): `http://localhost:8080`
- Remote (Render หลัง Deploy): `https://YOUR-BACKEND.onrender.com` **เป็น Placeholder**
- ทุก Request ใช้ JSON; ฝั่ง Protected API ใส่ `Authorization: Bearer <JWT_TOKEN>`
- JSON ใช้ `camelCase` เช่น `imageUrl`, `displayName`, `reviewCount`
- Response Error ทุก API: `{ "error": { "code": "ERROR_CODE", "message": "ข้อความอธิบาย" } }`
- `401` ไม่มี/Token ไม่ถูกต้อง, `403` ไม่มีสิทธิ์, `404` ไม่พบข้อมูล, `409` ข้อมูลซ้ำ, `400` Validation ไม่ผ่าน

## 2. Restaurant / Province / Category

| Method | Endpoint | Query/Body | Response |
|---|---|---|---|
| GET | `/api/provinces` | - | `{ "items": [{"id":1,"nameTh":"เชียงใหม่","region":"ภาคเหนือ"}] }` |
| GET | `/api/categories` | - | `{ "items": [{"id":1,"nameTh":"อาหารเหนือ"}] }` |
| GET | `/api/restaurants` | `province`, `category`, `q`, `page`, `limit` | `{ "items": Restaurant[], "page": 1, "limit": 12, "total": 0 }` |
| GET | `/api/restaurants/{id}` | - | `Restaurant` หรือ 404 |
| GET | `/api/restaurants/{id}/reviews` | - | `{ "items": Review[] }` |

ตัวอย่าง `Restaurant` (ข้อมูลสมมติ):

```json
{
  "id": 101,
  "name": "ร้านอาหารตัวอย่าง",
  "province": "เชียงใหม่",
  "category": "อาหารเหนือ",
  "address": "อำเภอเมืองเชียงใหม่",
  "description": "ร้านอาหารตัวอย่างสำหรับทดสอบ",
  "latitude": 18.788,
  "longitude": 98.985,
  "imageUrl": null,
  "rating": 4.5,
  "reviewCount": 2
}
```

- `rating` มาจากรีวิวบนเว็บ คะแนนเต็ม 5; **ไม่มีรีวิวให้ `rating: null`, `reviewCount: 0`**
- Query ค่าเริ่มต้น: `page=1`, `limit=12`, สูงสุด `limit=50`
- ฟิลเตอร์ `province` และ `category` ส่งเป็นชื่อภาษาไทย
- ไม่มีผลค้นหา ให้ส่ง 200 และ `{ "items": [], "page": 1, "limit": 12, "total": 0 }`

## 3. Authentication / Member

| Method | Endpoint | Request body | Response |
|---|---|---|---|
| POST | `/api/auth/register` | `{"email":"...","password":"...","displayName":"..."}` | `201` และ `Member` |
| POST | `/api/auth/login` | `{"email":"...","password":"..."}` | `200` และ `{ "token": "<JWT>", "member": Member }` |
| GET | `/api/me` | JWT | `Member` |
| PATCH | `/api/me` | `{"displayName":"..."}` | `Member` ที่แก้แล้ว |
| PUT | `/api/me/password` | `{"currentPassword":"...","newPassword":"..."}` | `204` |

`Member`: `{"id":1,"email":"demo@example.com","displayName":"Demo","role":"member"}` (ห้ามส่ง `password_hash`)

Logout สำหรับ MVP: FE ลบ Token และข้อมูลผู้ใช้ออกจาก Browser แล้ว Redirect; JWT เดิมยังใช้ได้จนหมดอายุ หากต้องยกเลิกทันทีต้องเพิ่ม Token Revocation ฝั่ง BE และปรับ Contract ร่วมกัน

## 4. Review / Favorite

| Method | Endpoint | Request body | Response |
|---|---|---|---|
| POST | `/api/restaurants/{id}/reviews` | `{"rating":5,"text":"อาหารอร่อย บริการดี"}` | `201` และ `Review` |
| DELETE | `/api/reviews/{id}` | JWT (Owner/Admin) | `204` |
| GET | `/api/favorites` | JWT | `{ "items": Restaurant[] }` |
| POST | `/api/favorites` | `{"restaurantId":101}` | `201` และ `{ "restaurantId": 101 }` |
| DELETE | `/api/favorites/{restaurantId}` | JWT | `204` |

ตัวอย่าง `Review`:

```json
{
  "id": 5,
  "restaurantId": 101,
  "rating": 5,
  "text": "อาหารอร่อย บริการดี",
  "createdAt": "2026-11-11T14:00:00Z",
  "member": { "id": 1, "displayName": "Demo" }
}
```

กติกา: รีวิวต้องมี `rating` ระหว่าง 1–5 และ `text` ยาวอย่างน้อย 10 ตัวอักษร; สมาชิก 1 คนต่อ 1 ร้านได้ 1 รีวิว, หากซ้ำให้ 409; รีวิวลบได้เฉพาะเจ้าของหรือ Admin

## 5. Admin API (ต้องตรวจ role ใน Go)

| Method | Endpoint | Request/หมายเหตุ | Response |
|---|---|---|---|
| GET | `/api/admin/members` | JWT role admin | `{ "items": Member[] }` |
| POST | `/api/admin/restaurants` | ชื่อร้าน, `provinceId`, `categoryId`, ที่อยู่, รายละเอียด, พิกัดและ `imageUrl` | `201` และ `Restaurant` |
| PUT | `/api/admin/restaurants/{id}` | ฟิลด์ร้านอาหารที่แก้ไข | `200` และ `Restaurant` |
| DELETE | `/api/admin/restaurants/{id}` | JWT role admin | `204` |
| DELETE | `/api/admin/reviews/{id}` | JWT role admin | `204` |

ห้ามสมัคร Admin ผ่าน Register; BE ต้องตรวจ Role จากข้อมูลที่เชื่อถือได้ และไม่รับค่า Role จากผู้สมัคร

## 6. Integration Checklist

- [ ] Go API มี Health Check (เลือก `GET /health`) และเชื่อม Supabase จริง
- [ ] BE ให้ URL Render จริง พร้อม `ALLOWED_ORIGINS` ที่รวม Vercel และ localhost FE
- [ ] FE ตั้ง `REACT_APP_API_URL` ให้ชี้ Backend จริงและเรียกผ่าน `fetch()`; Deploy Vercel แล้ว Rebuild
- [ ] ทดสอบ GET ร้าน, Register, Login, Me, Review, Favorites, Admin CRUD
- [ ] FE รองรับ Loading, Empty, Error / Retry รวมถึง 401 และ 403
- [ ] BE ไม่ส่ง Sensitive Fields ใน JSON และไม่ Commit Secrets

**หมายเหตุ:** ชุดข้อมูล ตัวเลข และ URL ที่แสดงเป็นตัวอย่าง ไม่ใช่ข้อมูลหรือ Deployment จริง
