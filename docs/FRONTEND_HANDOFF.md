# ส่งต่องาน Frontend: ThaiFood Guide

Backend URL: `https://thaifood-guide-api-a7c7.onrender.com`

Render ขึ้นสถานะ live และทดสอบ Auth/Profile/Reviews/Favorites/Admin CRUD บน URL นี้ผ่านแล้ว รายละเอียด endpoint และ JSON ให้ยึด [API_CONTRACT.md](../API_CONTRACT.md)

## ตั้งค่า CRA และ Vercel

Frontend ใช้ Create React App ให้ตั้ง environment ต่อไปนี้ทั้งเครื่องพัฒนาและ Vercel:

```dotenv
REACT_APP_USE_MOCK=false
REACT_APP_API_URL=https://thaifood-guide-api-a7c7.onrender.com
```

อ่านผ่าน `process.env.REACT_APP_API_URL`; restart dev server หลังแก้ `.env` และ rebuild/redeploy Vercel หลังแก้ environment เพราะ CRA ฝังค่านี้ใน bundle

## Public API

| Endpoint | ผลลัพธ์สำคัญ |
|---|---|
| `GET /api/provinces` | `{ items: [{ id, nameTh, region }] }` ครบ 77 จังหวัด |
| `GET /api/categories` | `{ items: [{ id, nameTh }] }` |
| `GET /api/restaurants` | `{ items: Restaurant[], page, limit, total }` |
| `GET /api/restaurants/{id}` | `Restaurant` หรือ 404 |
| `GET /api/restaurants/{id}/reviews` | `{ items: Review[] }` |

ฟิลเตอร์ `province` และ `category` ส่งชื่อภาษาไทย; ค้นหาด้วย `q` และแบ่งหน้าด้วย `page`/`limit` ค่าเริ่มต้นคือ 1/12 และ limit สูงสุด 50

```javascript
const API_URL = process.env.REACT_APP_API_URL;

export async function getRestaurants(province) {
  const query = new URLSearchParams({ province, page: "1", limit: "12" });
  const response = await fetch(`${API_URL}/api/restaurants?${query}`);
  const data = await response.json();
  if (!response.ok) throw new Error(data.error?.message || `HTTP ${response.status}`);
  return data;
}
```

JSON ใช้ camelCase เช่น `imageUrl`, `displayName`, `reviewCount` ร้านไม่มีรีวิวส่ง `rating: null`, `reviewCount: 0`; อย่าแปลง null เป็นคะแนนศูนย์ คะแนนจากรีวิวในเว็บเต็ม 5 และ `imageUrl` อาจเป็น null ผลค้นหาว่างยังตอบ 200 พร้อม `items: []`

ข้อมูลเริ่มต้นปัจจุบันเป็นร้านจริง 20 ร้านใน 9 จังหวัดพร้อม `imageUrl` ของภาพที่ตรงกับร้าน ใช้ข้อมูล ททท. และผู้ประกอบการ พร้อมเครดิต/สัญญาอนุญาตภาพใน `description` ให้แสดงรายละเอียดส่วนนี้ต่อไป ดู [แหล่งข้อมูลและใบอนุญาตภาพ](RESTAURANT_DATA_SOURCES.md) คะแนนยังมาจากรีวิวในระบบเอง ไม่ใช่คะแนน Google และข้อมูล snapshot ไม่รับรองสถานะเปิดทำการล่าสุด

FE repository ที่ผู้ใช้ส่งมาอ่าน `restaurant.imageUrl` อยู่แล้ว และใช้ SVG เมื่อรูปเป็น null หรือโหลดไม่ได้ จึงไม่ต้องเปลี่ยนชื่อฟิลด์ หากต้องการข้อมูลจาก Backend จริง ให้ตั้ง `REACT_APP_USE_MOCK=false` คู่กับ `REACT_APP_API_URL` แล้ว restart/rebuild CRA; ค่าเริ่มต้นของ FE ยังเป็น mock หากไม่ตั้งตัวแปรนี้

## สมาชิกและข้อผิดพลาด

- Register ส่ง `email`, `password`, `displayName`; Login ส่ง `email`, `password` และรับ `{ token, member }`
- `Member` มี `{ id, email, displayName, role }` โดย role เป็น `member` หรือ `admin`; ไม่ส่ง role ใน Register
- Protected API ใช้ header `Authorization: Bearer <token>` และ request JSON ใช้ `Content-Type: application/json`
- Logout ให้ลบ token และข้อมูลผู้ใช้จาก browser แล้ว redirect; JWT อายุ 1 ชั่วโมงและยังใช้ได้จนหมดอายุ ไม่มี logout endpoint
- Error ใช้ `{ error: { code, message } }`; รองรับ 400/401/403/404/409 และ server errors พร้อม Loading/Empty/Error/Retry
- 204 ไม่มี response body จึงต้องตรวจ status ก่อนเรียก `response.json()`; เมื่อได้ 401 ให้เคลียร์ session และกลับหน้า Login ส่วน 403 ให้แสดงว่าไม่มีสิทธิ์

## เพิ่มโดเมน Frontend ใน CORS

ขณะนี้ `ALLOWED_ORIGINS` ตั้งเป็น `http://localhost:3000` ผู้ใช้จะส่ง URL Vercel จริงภายหลัง จึงต้องเพิ่ม origin ก่อนทดสอบจาก Vercel:

1. เปิด [Render service Environment](https://dashboard.render.com/web/srv-dbqaciqjnfac73dbf0f0/env)
2. แก้ `ALLOWED_ORIGINS` เป็น `http://localhost:3000,https://YOUR-ACTUAL-FRONTEND.vercel.app` โดยแทนตัวอย่างด้วยโดเมนจริง ไม่มี path หรือ `/` ท้าย URL
3. Save และ deploy ให้ค่าใหม่มีผล แล้วตรวจ preflight และ flow Register → Login → Profile → Favorite → Review → Logout

คง localhost ไว้และใช้รายการ origin ที่ระบุชัดเจน ห้ามใช้ wildcard `*` ไม่ต้องใส่ `DATABASE_URL`, `JWT_SECRET` หรือ Supabase database password ใน Frontend
