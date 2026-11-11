# Kheun Lift — แผนพัฒนา Backend (Go) + ภาพรวมโปรเจกต์

> ไฟล์นี้เป็น **แผนหลัก** ของโปรเจกต์ เก็บทั้งงานฝั่ง Go/DB และเรื่องที่ใช้ร่วมกันทั้งสองฝั่ง
> (การตัดสินใจ, สถาปัตยกรรม, migration, ความเสี่ยง)
>
> งานฝั่งหน้าเว็บอยู่ที่ [`kheun-lift-frontend-next/PLAN.md`](https://github.com/K-Bamrungsuk/Kheun-Lift-front-NEXT/blob/main/PLAN.md)
> (บนเครื่อง: `../kheun-lift-frontend-next/PLAN.md`)

**เป้าหมาย:** deploy จริง ใช้เป็นผลงานใน resume ภายใน 5 สัปดาห์

**หลักการ:**
1. ทุกสัปดาห์ต้องจบด้วยของที่ **deploy แล้วใช้ได้จริง** บนเว็บจริง ไม่ใช่แค่บนเครื่อง
2. ถ้าเวลาไม่พอ **ตัดจากท้ายแผน** ไม่ใช่ทำทุกอย่างครึ่งๆ (ลำดับการตัดอยู่ท้ายไฟล์)
3. เขียนเองก่อน หรือให้ AI เขียนแล้ว **อ่านจนอธิบายได้ทุกบรรทัด** — ตอนสัมภาษณ์จะโดนถาม

---

## สารบัญ

- [วิธีอ่านแผนนี้](#วิธีอ่านแผนนี้)
- [การตัดสินใจ D1–D6](#การตัดสินใจ-d1d6)
- [สถาปัตยกรรมตอน deploy](#สถาปัตยกรรมตอน-deploy)
- [วิธีเพิ่มฟีเจอร์ใน Go ตามโครงที่มีอยู่](#วิธีเพิ่มฟีเจอร์ใน-go-ตามโครงที่มีอยู่)
- [สัปดาห์ 0 — เตรียมตัว](#สัปดาห์-0--เตรียมตัว)
- [สัปดาห์ 1 — พร้อมให้ Next ต่อ + deploy](#สัปดาห์-1--พร้อมให้-next-ต่อ--deploy)
- [สัปดาห์ 2 — Admin + Achievement](#สัปดาห์-2--admin--achievement)
- [สัปดาห์ 3 — Forgot password + Google OAuth](#สัปดาห์-3--forgot-password-otp--google-oauth)
- [สัปดาห์ 4 — Feed + Like + Comment](#สัปดาห์-4--feed--like--comment)
- [สัปดาห์ 5 — Notification + เก็บงาน](#สัปดาห์-5--notification--เก็บงาน)
- [Migration ทั้งหมด](#migration-ทั้งหมด)
- [Environment variables](#environment-variables)
- [Endpoint ใหม่ทั้งหมด](#endpoint-ใหม่ทั้งหมด)
- [ความเสี่ยง](#ความเสี่ยง)
- [ถ้าไม่ทัน ตัดอะไรก่อน](#ถ้าไม่ทัน-ตัดอะไรก่อน)

---

## วิธีอ่านแผนนี้

- งานแต่ละข้อมีรหัส เช่น **B1.3** = Backend สัปดาห์ 1 ข้อ 3 ส่วนฝั่ง Next ใช้ **F** (เช่น F1.4)
- ถ้าเขียนว่า **"ต้องรอ B2.2"** แปลว่าต้องทำข้อนั้นให้เสร็จก่อน
- แต่ละข้อแบ่งเป็น **ทำอะไร / ทำไม / ทำยังไง / เสร็จเมื่อ**
- ติ๊ก `[x]` เมื่อเสร็จ แล้ว commit ไฟล์นี้ไปพร้อมโค้ด
- กฎของ repo ยังใช้: เสนอก่อนแก้, **ห้ามรัน migration/seed กับ DB โดยไม่ขออนุมัติ**, ห้าม build เอง, test ได้เฉพาะจุดเสี่ยงและต้องอนุมัติ (ดู `AGENTS.md`)

---

## การตัดสินใจ D1–D6

✅ ตัดสินใจแล้ว (2026-11-11)

| # | เรื่อง | เลือก | เหตุผล |
|---|---|---|---|
| D1 | เชื่อม Next กับ Go ยังไง | **ซื้อ domain** แยก `app.<domain>` (Next) กับ `api.<domain>` (Go) | 2 subdomain ของ domain เดียวกันนับเป็น "same-site" → cookie `SameSite=Lax` ส่งข้ามกันได้, Go เห็น IP จริงของผู้ใช้, และ Resend ต้องมี domain สำหรับส่ง OTP อยู่แล้ว |
| D2 | Host | Next → **Vercel**, Go → **Render** หรือ **Railway**, DB → **Neon project ใหม่** ทุกตัว region **Singapore** | ฟรีหรือถูก, deploy อัตโนมัติจาก GitHub, อยู่ใกล้ผู้ใช้ไทย |
| D3 | ส่ง email | **Resend** เรียกผ่าน HTTP ด้วย `net/http` | ไม่ต้องลง SDK, API ง่าย |
| D4 | Repost | **ตัดออก** | ต้องจัดการ feed ซ้ำ, ลบต้นฉบับแล้ว repost เป็นยังไง — ซับซ้อนแต่ไม่โชว์ฝีมือเพิ่ม |
| D5 | Test | **มีได้** เฉพาะจุดเสี่ยงและต้องอนุมัติ | คนรีวิว resume ดูว่ามี test ไหม |
| D6 | เหตุผลตอน reject | **ใส่** | ผู้ใช้ต้องรู้ว่าทำไม lift ไม่ผ่าน |

---

## สถาปัตยกรรมตอน deploy

```
                    ผู้ใช้ (browser)
                   │                │
     หน้าเว็บ (HTML/JS)            API (fetch + cookie)
                   ▼                ▼
   app.<domain>  ── Vercel        api.<domain>  ── Render/Railway
   Next.js 16                     Go (Gin)
   (Server Component ดึงข้อมูล ──►  │
    โดยส่ง cookie ต่อให้ Go)         ▼
                                  Neon PostgreSQL (prod project, Singapore)
                                  Resend (ส่ง email)    Google OAuth
```

**เรื่องที่ต้องเข้าใจก่อนเริ่ม:**

- **Cookie:** Go ตั้ง cookie `kheun_session` ด้วย `Domain=.<domain>` → browser ส่ง cookie นี้ไปทั้ง `app.` และ `api.`
- **CORS:** browser ที่อยู่หน้า `app.` เรียก `api.` = คนละ origin → Go ต้องตอบ header CORS ว่ายอมให้ `https://app.<domain>` ส่ง cookie มาได้ (`Access-Control-Allow-Credentials: true`)
- **CrossOriginProtection:** Go บล็อก POST/PATCH/DELETE ที่มาจาก origin อื่นอยู่ (กัน CSRF) → ต้องเพิ่ม `app.<domain>` เป็น trusted origin
- **IP จริง:** host ของ Go มี proxy อยู่ข้างหน้า → ต้องบอก Gin ว่าเชื่อ header ไหนที่บอก IP จริง ไม่งั้น rate limit จะนับทุกคนเป็น IP เดียว
- **⚠ Login / register / forgot password ต้องเรียกจาก browser ตรงๆ** ไม่ใช่ผ่าน server ของ Next — ไม่งั้น Go จะเห็น IP ของ Vercel แทน IP ผู้ใช้ และ rate limit จะพัง

---

## วิธีเพิ่มฟีเจอร์ใน Go ตามโครงที่มีอยู่

ทุกฟีเจอร์ใหม่ให้ **ลอกโครงจาก package ที่คล้ายที่สุด** (ส่วนใหญ่คือ `internal/lifts/`)

| ชั้น | ไฟล์ | ทำอะไร | ตัวอย่างที่มีอยู่ |
|---|---|---|---|
| queries | `internal/queries/<feature>.go` | เก็บ SQL เป็น `const` | `queries/lifts.go` |
| model | `internal/<feature>/model/` | struct ของข้อมูลเข้า/ออก + tag `json` | `lifts/model/lifts.go` |
| repository | `internal/<feature>/repository/` | รัน SQL แล้ว `Scan` ใส่ struct — ไม่มี logic | `lifts/repository/lifts.go` |
| service | `internal/<feature>/service/` | logic + transaction (ถ้าง่ายมากข้ามได้ เหมือน `ListMine` ที่ handler เรียก repository ตรง) | `lifts/service/lifts.go` |
| handler | `internal/<feature>/handler/` | อ่าน request, validate, เรียก service, ตอบ JSON | `lifts/handler/lifts.go` |
| routes | `internal/routes/routes.go` | ผูก path กับ handler, public หรือ `protected` | — |

**รูปแบบที่ต้องทำตามทุก endpoint (ดูใน `lifts/handler/lifts.go`):**
- ใส่ `context.WithTimeout(..., 5*time.Second)` ทุก request ที่แตะ DB
- อ่าน body ด้วย `decodeBody` (เช็ก Content-Type, จำกัด 16 KB, ไม่รับ field แปลก)
- id จาก URL ใช้ `parseID` (ต้องเป็นจำนวนเต็มบวก 32-bit)
- error ที่ผู้ใช้ควรรู้ → ประกาศเป็น `var ErrXxx = errors.New(...)` ใน service แล้วแปลงเป็น status ใน `respondError`
- error อื่น → `log` แล้วตอบ 500 ข้อความกลางๆ **ห้ามส่ง error จาก DB ให้ผู้ใช้เห็น**
- ตอบรูปแบบเดียวกัน: `{"message": "...", "data": ...}`

---

## สัปดาห์ 0 — เตรียมตัว

- [ ] **B0.1 Commit งานค้าง** — ย้าย seed ไป `seeds/` + `Makefile` + ไฟล์ `PLAN.md` นี้
- [ ] **B0.2 ซื้อ domain** — ร้านไหนก็ได้ (Cloudflare Registrar, Namecheap, Porkbun) ราคาประมาณ 300–500 บาท/ปี
  - **ทำไมรีบ:** การยืนยัน DNS ของ Resend / Vercel ใช้เวลาตั้งแต่ไม่กี่นาทีถึงหลายชั่วโมง
- [ ] **B0.3 สมัครบริการ** — Render หรือ Railway (เลือกตัวเดียว), Resend, Google Cloud (สร้าง OAuth client ไว้ก่อน)
- [ ] **B0.4 สร้าง Neon project สำหรับ production** — region Singapore
  - project ปัจจุบันถือเป็น **dev** (มี user ปลอม `password123`) **ห้ามใช้เป็น prod**
  - เก็บ connection string 2 แบบ: **pooled** (host มี `-pooler`) ให้แอปใช้, **direct** ให้ goose ใช้

---

## สัปดาห์ 1 — พร้อมให้ Next ต่อ + deploy

**เป้าหมาย:** Go บน `api.<domain>` รับ request จาก `app.<domain>` ได้ทั้ง cookie, CORS และเห็น IP จริง

### B1.1 Spike: deploy Go ขึ้น host จริง (วันที่ 1)
- **ทำอะไร:** deploy Go ตัวปัจจุบันขึ้น Render/Railway ชี้ `api.<domain>` มาที่มัน ต่อ DB prod
- **ทำไม:** ปัญหา proxy / cookie / HTTPS เจอเฉพาะบน host จริง เจอวันแรกแก้ง่ายกว่าเจอสัปดาห์ 5
- **ทำยังไง:**
  1. ตั้ง build command `go build -o app ./cmd/api` และ start command `./app`
  2. ตั้ง env ตาม [ตาราง env](#environment-variables)
  3. รัน migration กับ DB prod **(ขออนุมัติก่อน)**: `make migrate` โดยใช้ `DIRECT_DATABASE_URL` ของ prod — **ห้าม `make seed`**
  4. ชี้ DNS `api.<domain>` ไปที่ host (CNAME ตามที่ host บอก)
  5. ทำพร้อมกับ F1.1 (deploy Next หน้าเปล่า)
- **เสร็จเมื่อ:** `https://api.<domain>/health` ตอบ `{"status":"ok"}`

### B1.2 แก้ rate limit ให้เห็น IP จริง (review รอบแรกข้อ 1)
- **ทำอะไร:** บอก Gin ว่าเชื่อ header IP จาก proxy ของ host
- **ทำไม:** ตอนนี้ `main.go` ตั้ง `SetTrustedProxies(nil)` → `c.ClientIP()` ได้ IP ของ proxy → ทุกคนใช้โควตา login 5 ครั้ง/นาทีร่วมกัน
- **ทำยังไง:**
  1. อ่าน docs ของ host ว่าส่ง IP จริงมาใน header ไหน (เช่น `X-Forwarded-For`, `True-Client-IP`)
  2. ถ้าเป็น header ที่มี IP เดียว → `router.TrustedPlatform = "<ชื่อ header>"`
  3. ถ้าเป็น `X-Forwarded-For` → `router.SetTrustedProxies(<IP range ของ proxy ตาม docs>)`
  4. อ่านค่าจาก env (`TRUSTED_PLATFORM_HEADER`) เพื่อให้บนเครื่องยังทำงานแบบเดิม
- **เสร็จเมื่อ:** เพิ่ม log ชั่วคราว พิมพ์ `c.ClientIP()` ตอน login บน host จริง → ต้องเป็น IP บ้านเรา (เช็กกับ whatismyip) แล้วลบ log ออก

### B1.3 CORS + trusted origin
- **ทำอะไร:** ยอมให้ `https://app.<domain>` เรียก API พร้อม cookie
- **ทำไม:** browser จะบล็อก response ถ้า Go ไม่ตอบ header CORS, และ `CrossOriginProtection` จะตอบ 403 กับ POST ที่มาจาก `app.`
- **ทำยังไง:**
  1. middleware เล็กๆ ใน `internal/middleware/cors.go` (ไม่ต้องลง library):
     - ถ้า `Origin` == `APP_ORIGIN` → ตอบ `Access-Control-Allow-Origin: <origin>`, `Access-Control-Allow-Credentials: true`, `Vary: Origin`
     - ถ้าเป็น `OPTIONS` (preflight) → ตอบ `Allow-Methods: GET, POST, PUT, PATCH, DELETE`, `Allow-Headers: Content-Type` แล้ว 204
  2. ใส่ middleware นี้ที่ `router.Use(...)` ก่อน route ทั้งหมด
  3. ใน `routes.go` → `crossOrigin.AddTrustedOrigin(os.Getenv("APP_ORIGIN"))`
- **เสร็จเมื่อ:** จากหน้า `app.<domain>` ใน DevTools console → `fetch("https://api.<domain>/users/me", {credentials: "include"})` ไม่ติด CORS error

### B1.4 รวมโค้ด cookie + ใส่ `Domain` (review รอบแรกข้อ 2)
- **ทำอะไร:** ฟังก์ชันเดียว `setSessionCookie(w, value, expires)` ใน `internal/auth/handler/` ใช้ทั้ง login, logout, delete account
- **ทำไม:** ต้องเพิ่ม `Domain` อยู่แล้ว ถ้าแก้ไม่ครบ 3 ที่ logout จะลบ cookie ไม่ได้ (browser ลบได้ต่อเมื่อ Domain/Path ตรงกับตอนตั้ง)
- **ทำยังไง:** อ่าน `COOKIE_SECURE` และ `COOKIE_DOMAIN` ครั้งเดียวใน `routes.go` แล้วส่งเข้า handler (ตอนนี้อ่าน `COOKIE_SECURE` 3 ครั้ง)
  - `user/handler/delete.go` ต้องเรียกฟังก์ชันนี้ได้ด้วย → วางไว้ที่ที่ทั้งสอง package import ได้ (เช่น `internal/middleware/` ซึ่ง `user/handler` import อยู่แล้ว)
- **เสร็จเมื่อ:** login → logout บนเว็บจริง แล้ว DevTools → Application → Cookies ไม่มี `kheun_session` เหลือ

### B1.5 เปิด `/weight-classes` ให้ดูได้โดยไม่ login
- **ทำไม:** หน้า leaderboard ต้องดูได้โดยไม่ login (HR ที่เข้ามาดูจะไม่สมัคร) แต่ตัวเลือกรุ่นน้ำหนักต้องใช้ endpoint นี้
- **ทำยังไง:** ย้ายบรรทัด `protected.GET("/weight-classes", ...)` ไปเป็น `router.GET(...)` ใน `routes.go`
- **เสร็จเมื่อ:** เรียกโดยไม่มี cookie แล้วได้ 200

### B1.6 เก็บกวาดเล็กๆ (review รอบแรกข้อ 4)
- [ ] แก้ข้อความ `"model.Exercise not found"` → `"Exercise not found."` ใน `exercises/handler/exercises.go`
- [ ] ลบ `type Exercise` ที่ไม่ได้ใช้ใน `auth/model/auth.go`
- [ ] `.env.example` เพิ่มตัวแปรใหม่ทั้งหมดของสัปดาห์นี้

**สัปดาห์ 1 เสร็จเมื่อ:** (ร่วมกับฝั่ง Next) สมัคร → login → กรอกโปรไฟล์ → ส่ง lift → เห็นใน profile **บนเว็บจริง**

---

## สัปดาห์ 2 — Admin + Achievement

**ทำไมต้องทำก่อนฟีเจอร์อื่น:** ตอนนี้ไม่มีทาง verify lift นอกจากแก้ DB เอง → leaderboard ว่างตลอด และ achievement, feed, notification ต้องใช้ lift ที่ verified ทั้งหมด

### B2.1 Migration `00007_add_user_role.sql` (ต้องอนุมัติก่อนรัน)
```sql
-- +goose Up
CREATE TYPE "UserRole" AS ENUM ('user', 'admin');
ALTER TABLE "User" ADD COLUMN "role" "UserRole" NOT NULL DEFAULT 'user';
ALTER TABLE "LiftRecord" ADD COLUMN "rejectReason" TEXT;

-- +goose Down
ALTER TABLE "LiftRecord" DROP COLUMN "rejectReason";
ALTER TABLE "User" DROP COLUMN "role";
DROP TYPE "UserRole";
```
- **ตั้ง admin คนแรก:** รัน SQL เองใน Neon console → `UPDATE "User" SET "role" = 'admin' WHERE "email" = '<email ของคุณ>';`
- **ทำไมไม่มี endpoint ตั้ง admin:** ถ้ามี ก็ต้องป้องกันไม่ให้คนแอบเรียก — ไม่มีเลยปลอดภัยที่สุด

### B2.2 Session รู้ role
- **ทำยังไง:**
  1. `queries/auth.go` → `FindSessionUser` เพิ่ม `u."role"` ใน SELECT
  2. `auth/model/auth.go` → `User` เพิ่ม `Role string \`json:"role"\``
  3. `auth/repository/session.go` → เพิ่ม `&session.User.Role` ใน `Scan` (ลำดับต้องตรงกับ SELECT)
  4. `CreateUser`, `UpdateUser` ที่ `RETURNING` ข้อมูล user → เพิ่ม `"role"` ด้วย ไม่งั้น `/users/me` หลังแก้โปรไฟล์จะได้ role ว่าง
- **เสร็จเมื่อ:** `GET /users/me` มี `"role": "user"` หรือ `"admin"`

### B2.3 `RequireAdmin` middleware
- **ทำยังไง:** `internal/middleware/admin.go` → ดึง session จาก context (`SessionFromContext`) ถ้า `Role != "admin"` → 403 `{"message": "Admin only."}`
- ใน `routes.go`: `admin := protected.Group("/admin"); admin.Use(appmiddleware.RequireAdmin())`

### B2.4 `GET /admin/lifts?status=pending&cursor=<id>`
- **ทำอะไร:** รายการ lift ที่รอตรวจ ทีละ 20 รายการ **เก่าสุดก่อน** (คนที่รอนานได้ตรวจก่อน)
- **ทำยังไง:**
  - SQL: ใช้ `liftColumns` + `liftJoins` ที่มีอยู่ใน `queries/lifts.go` แล้วเพิ่ม `WHERE l."status" = $1 AND l."id" > $2 ORDER BY l."id" LIMIT 20`
  - **cursor** = id ตัวสุดท้ายของหน้าก่อน (หน้าแรกส่ง 0) — เร็วกว่า `OFFSET` และไม่ข้าม/ซ้ำเมื่อมีข้อมูลใหม่เข้ามา
  - ตอบ `{"data": [...], "nextCursor": <id สุดท้าย หรือ null ถ้าไม่ครบ 20>}`
- **เสร็จเมื่อ:** user ธรรมดาได้ 403, admin ได้รายการ

### B2.5 `PATCH /admin/lifts/:id/status`
- **body:** `{"status": "verified" | "rejected", "reason": "..."}`
- **กฎ:**
  - เปลี่ยนได้เฉพาะ lift ที่ยัง `pending` → ถ้าไม่ใช่ตอบ 409 (กันกดซ้ำ และทำให้ notification ในสัปดาห์ 5 ง่าย)
  - `rejected` ต้องมี `reason` 1–300 ตัวอักษร, `verified` ห้ามมี reason
- **SQL:** `UPDATE "LiftRecord" SET "status" = $2, "rejectReason" = $3, "updatedAt" = CURRENT_TIMESTAMP WHERE "id" = $1 AND "status" = 'pending' RETURNING ...` → ไม่มีแถวกลับมา = ไม่พบหรือไม่ใช่ pending
- **ทำไมเป็น transaction ใน service:** สัปดาห์ 5 จะเพิ่มการสร้าง notification ในจุดนี้ — วางโครงเป็น service + transaction ไว้ตั้งแต่ตอนนี้
- `rejectReason` ต้องโผล่ใน `LiftRecord` model (`lifts/model`) ด้วย ให้เจ้าของเห็นเหตุผล
- **เสร็จเมื่อ:** verify แล้ว lift โผล่ใน leaderboard, reject แล้วเจ้าของเห็นเหตุผลใน `/lifts/me`

### B2.6 Achievement — `GET /users/:id/achievements`
- **หลักคิด:** **คำนวณสดทุกครั้ง ไม่สร้างตาราง** — ข้อมูลน้อย คำนวณเร็ว และไม่มีปัญหาข้อมูลไม่ตรงกัน (เช่น lift ถูกลบแต่ badge ยังอยู่)
- **รายการ:**

| key | ชื่อ | เงื่อนไข (นับเฉพาะ verified) |
|---|---|---|
| `first_rep` | First Rep | มี lift อย่างน้อย 1 |
| `big_three` | Big Three | มีครบทั้ง Squat, Bench Press, Deadlift |
| `bench_1x` | Bench Bodyweight | Bench ≥ 1 × น้ำหนักตัว |
| `squat_1_5x` | Squat 1.5× | Squat ≥ 1.5 × น้ำหนักตัว |
| `deadlift_2x` | Deadlift 2× | Deadlift ≥ 2 × น้ำหนักตัว |
| `club_300` / `club_400` / `club_500` | 300 / 400 / 500 Club | ผลรวม best ของ 3 ท่า ≥ 300 / 400 / 500 kg |
| `podium` | Podium | อันดับ ≤ 3 ในรุ่นไหนก็ได้ |
| `champion` | Champion | อันดับ 1 ในรุ่นไหนก็ได้ |

- **ทำยังไง:**
  1. Query 1 — best ของแต่ละท่า: `SELECT e."name", MAX(l."weight") FROM "LiftRecord" l JOIN "Exercise" e ... WHERE l."userId" = $1 AND l."status" = 'verified' GROUP BY e."name"`
  2. Query 2 — อันดับดีที่สุด: ใช้ window function `RANK() OVER (PARTITION BY "exerciseId", "weightClassId" ORDER BY best_weight DESC, ...)` กับ best ต่อ user แล้วหา `MIN(rank)` ของ user นี้ (เรียงแบบเดียวกับ `ListLeaderboard`)
  3. ใน Go: วนเช็กเงื่อนไขทีละ achievement → ตอบทุกตัวพร้อม `achieved: true/false`
  4. รายการ achievement เก็บเป็น slice ใน Go (ไม่ต้องอยู่ใน DB)
- **ข้อจำกัด:** ใช้น้ำหนักตัว **ปัจจุบัน** → ถ้าผู้ใช้แก้น้ำหนักขึ้น badge `_1x/_2x` อาจหาย ใส่ comment `// shortcut: uses current bodyWeight; store bodyWeight on LiftRecord if badges must stay`
- **Test (ถ้าอนุมัติ):** ฟังก์ชันเช็กเงื่อนไขแยกจาก DB → test ได้ง่ายด้วย table-driven test
- **เสร็จเมื่อ:** ผู้ใช้ที่มี lift verified เห็น badge ที่ถูกต้อง

---

## สัปดาห์ 3 — Forgot password (OTP) + Google OAuth

### B3.1 Migration `00008_create_password_reset.sql`
```sql
-- +goose Up
CREATE TABLE "PasswordReset" (
    "userId"    INTEGER PRIMARY KEY REFERENCES "User"("id") ON DELETE CASCADE,
    "codeHash"  BYTEA NOT NULL,
    "attempts"  INTEGER NOT NULL DEFAULT 0,
    "expiresAt" TIMESTAMPTZ NOT NULL,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE "PasswordReset";
```
- `userId` เป็น PK → **1 user มีรหัสที่ใช้ได้แค่ 1 รหัส** ขอใหม่ = แทนที่ตัวเก่า (`ON CONFLICT ("userId") DO UPDATE`)

### B3.2 `internal/email/` — ส่ง email ผ่าน Resend
- ฟังก์ชันเดียว `Send(ctx, to, subject, text string) error`
- `POST https://api.resend.com/emails` ด้วย header `Authorization: Bearer <RESEND_API_KEY>` และ JSON `{from, to, subject, text}`
- ใช้ `net/http` + timeout 10 วินาที ไม่ต้องลง SDK
- **ต้องยืนยัน domain ใน Resend ก่อน** ไม่งั้นส่งได้แค่หา email ของตัวเอง

### B3.3 `POST /auth/password/forgot` body `{email}`
- **ตอบ 200 ข้อความเดิมเสมอ** ("If this email exists, we sent a code.") ไม่ว่ามี email นี้หรือไม่ → กันคนเดาว่า email ไหนมีบัญชี
- ถ้ามี user:
  1. ถ้าขอรหัสไปไม่ถึง **60 วินาที** → ไม่ส่งใหม่ (กันคนสแปม email เหยื่อ)
  2. สร้างรหัส 6 หลักด้วย `crypto/rand` (`rand.Int(rand.Reader, big.NewInt(1_000_000))` → `fmt.Sprintf("%06d", n)`) — **ห้ามใช้ `math/rand`**
  3. เก็บแค่ `sha256` ของรหัส, `attempts = 0`, หมดอายุ 10 นาที
  4. ส่ง email
- ใส่ `authRateLimit` (ตัวเดียวกับ login)

### B3.4 `POST /auth/password/reset` body `{email, code, newPassword}`
- ตรวจตามลำดับ:
  1. หา user + แถว reset ไม่เจอ / หมดอายุ / `attempts >= 5` → 400 "Invalid or expired code."
  2. hash รหัสที่ส่งมาเทียบด้วย `subtle.ConstantTimeCompare` → ผิด: `attempts + 1` แล้ว 400
  3. ถูก: **ใน transaction เดียว** → อัปเดตรหัสผ่าน (ใช้ `hashPassword` ที่มีอยู่), **ลบ session ทั้งหมดของ user** (คนที่แอบ login อยู่จะหลุด), ลบแถว reset
- กฎรหัสผ่านใหม่เหมือนตอน register (≥ 6 ตัว, ≤ 1024 byte) → ใช้ฟังก์ชัน validate ร่วมกัน
- **Test (ถ้าอนุมัติ):** หมดอายุ, ใส่ผิดครบ 5 ครั้ง, ใช้แล้วใช้ซ้ำไม่ได้
- **เสร็จเมื่อ:** ขอรหัส → ได้ email จริง → ตั้งรหัสใหม่ → login ด้วยรหัสใหม่ได้, รหัสเก่าใช้ไม่ได้

### B3.5 Migration `00009_add_google_login.sql`
```sql
-- +goose Up
ALTER TABLE "User" ADD COLUMN "googleSub" TEXT UNIQUE;
ALTER TABLE "User" ALTER COLUMN "password" DROP NOT NULL;

-- +goose Down
ALTER TABLE "User" ALTER COLUMN "password" SET NOT NULL;
ALTER TABLE "User" DROP COLUMN "googleSub";
```
- ⚠ Down จะ fail ถ้ามี user ที่ไม่มีรหัสผ่านแล้ว — ยอมรับได้ (ไม่ควร rollback หลังมีผู้ใช้จริง)
- `"googleSub"` = id ถาวรของบัญชี Google (ใช้แทน email เพราะ email เปลี่ยนได้)

### B3.6 แก้ login ให้รองรับคนไม่มีรหัสผ่าน
- `FindLoginUser` สแกน `password` เป็น `*string` → ถ้า `nil` ให้ทำเหมือนไม่เจอ user (ยังรัน `verifyPassword` กับ dummy hash เพื่อให้เวลาตอบเท่ากัน)
- คนที่สมัครด้วย Google แล้วอยากใช้รหัสผ่าน → ใช้ forgot password ตั้งรหัสได้เลย ไม่ต้องทำหน้าใหม่

### B3.7 `GET /auth/google` และ `GET /auth/google/callback`
- **library:** `golang.org/x/oauth2` + `golang.org/x/oauth2/google` (ทีม Go ดูแลเอง)
- **`/auth/google`:**
  1. สร้าง `state` สุ่ม 32 byte + PKCE verifier (`oauth2.GenerateVerifier()`)
  2. เก็บทั้งสองใน cookie อายุ 10 นาที (`HttpOnly`, `Secure`, `SameSite=Lax`)
  3. redirect ไป `conf.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier))`
- **`/auth/google/callback`:**
  1. `state` ใน query ต้องตรงกับ cookie → ไม่ตรง = มีคนพยายามปลอม → 400
  2. `conf.Exchange(ctx, code, oauth2.VerifierOption(verifier))`
  3. เรียก `https://openidconnect.googleapis.com/v1/userinfo` ด้วย token → ได้ `sub`, `email`, `email_verified`
  4. หา user:
     - มี `googleSub` นี้ → login
     - ไม่มี แต่มี email นี้ **และ `email_verified = true`** → ผูก `googleSub` ให้บัญชีเดิม แล้ว login
     - ไม่มีเลย → สร้าง user ใหม่ ไม่มีรหัสผ่าน, username = ส่วนหน้า email (เอาแค่ a-z0-9 ไม่เกิน 15 ตัว) + เลขสุ่ม 4 หลัก, ถ้าชน unique ให้สุ่มใหม่ (สูงสุด 5 ครั้ง)
  5. สร้าง session (ดู B3.8) → ตั้ง cookie → redirect ไป `APP_ORIGIN + "/dashboard"`
  6. error ใดๆ → redirect ไป `APP_ORIGIN + "/login?error=google"` (ไม่ทิ้งผู้ใช้ไว้ที่หน้า JSON)
- **⚠ ห้ามผูกบัญชีถ้า `email_verified` เป็น false** — ไม่งั้นใครก็สร้าง Google account ด้วย email คนอื่นแล้วเข้าบัญชีเขาได้
- **ใน Google Cloud Console:** redirect URI = `https://api.<domain>/auth/google/callback` (+ `http://localhost:3001/auth/google/callback` สำหรับ dev)

### B3.8 แยกการสร้าง session ออกมาใช้ร่วม
- ตอนนี้โค้ดสร้าง token + hash + insert session อยู่ใน `service.Login`
- แยกเป็น `service.CreateSession(ctx, pool, userID) (token string, expiresAt time.Time, error)` → `Login` และ Google callback เรียกตัวเดียวกัน
- **เสร็จเมื่อ (B3.5–B3.8):** login ด้วย Google ได้ทั้ง (ก) บัญชีใหม่ (ข) บัญชีที่เคยสมัครด้วย email เดียวกัน

---

## สัปดาห์ 4 — Feed + Like + Comment

**หลักคิด:** feed = lift ที่ verified แล้ว เรียงใหม่สุดก่อน **ไม่สร้างตาราง post แยก** — lift record มีทุกอย่างที่โพสต์ต้องมีอยู่แล้ว (ใคร, ท่าไหน, หนักเท่าไร, caption, วิดีโอ)

### B4.1 Migration `00010_create_likes_comments.sql`
```sql
-- +goose Up
CREATE TABLE "LiftLike" (
    "liftId"    INTEGER NOT NULL REFERENCES "LiftRecord"("id") ON DELETE CASCADE,
    "userId"    INTEGER NOT NULL REFERENCES "User"("id") ON DELETE CASCADE,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("liftId", "userId")
);

CREATE TABLE "LiftComment" (
    "id"        SERIAL PRIMARY KEY,
    "liftId"    INTEGER NOT NULL REFERENCES "LiftRecord"("id") ON DELETE CASCADE,
    "userId"    INTEGER NOT NULL REFERENCES "User"("id") ON DELETE CASCADE,
    "body"      TEXT NOT NULL CHECK (char_length("body") BETWEEN 1 AND 500),
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX "LiftComment_liftId_id_idx" ON "LiftComment"("liftId", "id");
CREATE INDEX "LiftRecord_feed_idx" ON "LiftRecord"("status", "id" DESC);

-- +goose Down
DROP INDEX "LiftRecord_feed_idx";
DROP TABLE "LiftComment";
DROP TABLE "LiftLike";
```
- **PK (`liftId`, `userId`)** → like ซ้ำไม่ได้โดยธรรมชาติ ไม่ต้องเช็กเอง
- **`ON DELETE CASCADE`** → ลบ lift แล้ว like/comment หายตาม, ไม่ทำให้การลบ lift พัง

### B4.2 `GET /feed?cursor=<id>`
- ทีละ 20, `WHERE l."status" = 'verified' AND ($cursor = 0 OR l."id" < $cursor) ORDER BY l."id" DESC`
- แต่ละรายการมี: ข้อมูล lift + `user` (id, username, profileImage) + `likeCount` + `commentCount` + `likedByMe`
- **ต้องดึงทั้งหมดใน query เดียว** ด้วย subquery:
  ```sql
  (SELECT COUNT(*) FROM "LiftLike" k WHERE k."liftId" = l."id") AS "likeCount",
  (SELECT COUNT(*) FROM "LiftComment" c WHERE c."liftId" = l."id") AS "commentCount",
  EXISTS (SELECT 1 FROM "LiftLike" k WHERE k."liftId" = l."id" AND k."userId" = $me) AS "likedByMe"
  ```
  **ห้าม** query นับ like ทีละแถวใน Go loop (N+1 query = 20 รายการกลายเป็น 41 query)
- ตอบ `{"data": [...], "nextCursor": ...}` แบบเดียวกับ B2.4
- ต้อง login (ต้องรู้ `likedByMe`)

### B4.3 Like
- `PUT /lifts/:id/like` → `INSERT ... ON CONFLICT DO NOTHING` (กดซ้ำไม่ error)
- `DELETE /lifts/:id/like` → ลบ (ไม่มีอยู่แล้วก็ไม่ error)
- like ได้เฉพาะ lift ที่ verified → ไม่ใช่ตอบ 404
- ตอบ `{"likeCount": N, "likedByMe": true/false}` ให้หน้าเว็บอัปเดตตัวเลขได้เลย

### B4.4 Comment
- `GET /lifts/:id/comments?cursor=<id>` — เก่าสุดก่อน ทีละ 20, มี username + profileImage
- `POST /lifts/:id/comments` body `{body}` — trim แล้วต้องยาว 1–500 ตัวอักษร (`utf8.RuneCountInString`)
- `DELETE /comments/:id` — **เจ้าของ comment หรือ admin เท่านั้น** (`WHERE "id" = $1 AND ("userId" = $2 OR $3)`)

### B4.5 ลบบัญชีเมื่อมี lift (ต้องตัดสินใจ)
- ตอนนี้ลบบัญชีไม่ได้ถ้ายังมี lift (FK `RESTRICT`) แต่ like/comment จะ `CASCADE`
- ตัวเลือก: (ก) ให้ลบ lift ตามไปด้วย (เปลี่ยน FK เป็น CASCADE ใน migration) (ข) คงเดิม
- **แนะนำ (ก)** — ผู้ใช้ควรลบข้อมูลตัวเองได้ และ leaderboard จะคำนวณใหม่เองอยู่แล้ว
- **เสร็จเมื่อ (สัปดาห์ 4):** ไถ feed ได้, like/unlike ได้, comment และลบ comment ตัวเองได้, admin ลบ comment ใครก็ได้

---

## สัปดาห์ 5 — Notification + เก็บงาน

**หลักคิด:** อันดับใน leaderboard **เปลี่ยนได้ตอนเดียวคือตอน admin verify** (เพราะนับเฉพาะ verified) → สร้างแจ้งเตือนอันดับที่จุดนั้นจุดเดียว ไม่ต้องมี job มาคอยเช็ก

### B5.1 Migration `00011_create_notifications.sql`
```sql
-- +goose Up
CREATE TABLE "Notification" (
    "id"        SERIAL PRIMARY KEY,
    "userId"    INTEGER NOT NULL REFERENCES "User"("id") ON DELETE CASCADE,
    "type"      TEXT NOT NULL,
    "data"      JSONB NOT NULL DEFAULT '{}',
    "readAt"    TIMESTAMPTZ,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX "Notification_userId_id_idx" ON "Notification"("userId", "id" DESC);

-- +goose Down
DROP TABLE "Notification";
```
- เก็บ **ข้อมูล** (`data`) ไม่ใช่ข้อความ → หน้าเว็บเป็นคนสร้างข้อความเอง จะเปลี่ยนคำหรือทำหลายภาษาได้โดยไม่แก้ DB

### B5.2 สร้าง notification ตอน verify / reject (ต่อจาก B2.5)
| type | ใครได้ | data |
|---|---|---|
| `lift_verified` | เจ้าของ lift | `{liftId, exercise, weightClass, rank}` |
| `lift_rejected` | เจ้าของ lift | `{liftId, exercise, reason}` |
| `rank_dropped` | คนที่อันดับตกเพราะ lift นี้ | `{exercise, weightClass, oldRank, newRank}` |

- **ทำยังไง (ใน transaction เดียวกับการเปลี่ยน status):**
  1. ก่อน UPDATE: ดึงอันดับของ board นั้น (`exerciseId`, `weightClassId`) → `map[userId]rank` (ใช้ SQL แบบเดียวกับ `ListLeaderboard` แต่ SELECT แค่ `userId` กับ rank)
  2. UPDATE status
  3. หลัง UPDATE: ดึงอันดับอีกครั้ง
  4. เทียบ: ใครที่ rank ใหม่ > rank เก่า (ตกลง) และไม่ใช่เจ้าของ lift → insert `rank_dropped`
  5. เจ้าของ lift → insert `lift_verified` พร้อม rank ใหม่ (ถ้า lift นี้ไม่ใช่ best ของเขา ก็ยังได้แจ้งว่าผ่าน แต่อันดับเท่าเดิม)
- **ทำไมต้อง transaction:** ถ้า insert notification fail แต่ status เปลี่ยนไปแล้ว คนจะไม่ได้แจ้งเตือนตลอดไป
- ⚠ ข้อจำกัด: ผู้ใช้ลบ lift ที่ verified เอง → อันดับคนอื่นขยับขึ้นโดยไม่มีแจ้งเตือน (ยอมรับได้ ขึ้นอันดับไม่ต้องรีบบอก)

### B5.3 Endpoint
- `GET /notifications?cursor=<id>` — ใหม่สุดก่อน ทีละ 20
- `GET /notifications/unread-count` → `{"count": N}` (หน้าเว็บเรียกทุก 60 วิ ต้องเบา: `COUNT(*) WHERE "userId" = $1 AND "readAt" IS NULL`)
- `POST /notifications/read` → ตั้ง `readAt = now()` ให้ทุกอันที่ยังไม่อ่าน
- **(ถ้าเหลือเวลา)** `lift_liked`, `lift_commented` — insert ใน handler ของ like/comment (ไม่แจ้งเมื่อ like ของตัวเอง)

### B5.4 เก็บงานจาก review รอบแรก
- [ ] **ลบ session หมดอายุ** (ข้อ 5): ตอน login เพิ่ม `DELETE FROM "Session" WHERE "userId" = $1 AND "expiresAt" < CURRENT_TIMESTAMP`
- [ ] **รวม SQL หา weight class** (ข้อ 3): `FindLiftWeightClass` กับ `FindUserWeightClass` ใช้เงื่อนไขเดียวกัน → แยก WHERE เป็น const ตัวเดียว

### B5.5 Test (เฉพาะที่อนุมัติ)
| ไฟล์ | test อะไร |
|---|---|
| `auth/service/password_test.go` | hash แล้ว verify ผ่าน, รหัสผิดไม่ผ่าน, hash รูปแบบเสียได้ error |
| `achievements/.._test.go` | เงื่อนไข badge แต่ละตัว (ไม่ต้องใช้ DB) |
| `auth/service/reset_test.go` | OTP หมดอายุ, ผิดครบ 5, ใช้ซ้ำ (ต้องมี DB test → ใช้ Neon branch แยก) |

### B5.6 บัญชี demo + README
- [ ] สร้างบัญชี demo **1 บัญชี** บน prod ผ่านหน้า register ปกติ ส่ง lift จริง ให้ admin verify → มี badge, มีอันดับ
- [ ] README (repo นี้):
  - ลิงก์เว็บจริง + บัญชี demo
  - diagram จากหัวข้อ [สถาปัตยกรรม](#สถาปัตยกรรมตอน-deploy)
  - **เหตุผลที่เลือก** Go, session cookie แทน JWT (ยกเลิก session ได้ทันที เช่นตอน reset password), pgx + SQL ตรงแทน ORM, Neon
  - ความปลอดภัยที่ทำ: argon2id, rate limit, CSRF protection, OTP hash + จำกัดครั้ง, ไม่บอกว่า email มีในระบบ
  - วิธีรันบนเครื่อง: `.env`, `make migrate`, `make seed`, `make run`
  - ตาราง endpoint

---

## Migration ทั้งหมด

**DB dev บน Neon รันถึง version 6 แล้ว → ไฟล์ใหม่เริ่มที่ `00007`** (ถ้าตั้งเลขน้อยกว่านี้ goose จะ error)

| ไฟล์ | สัปดาห์ | เนื้อหา |
|---|---|---|
| `00007_add_user_role.sql` | 2 | enum `UserRole`, `User.role`, `LiftRecord.rejectReason` |
| `00008_create_password_reset.sql` | 3 | ตาราง `PasswordReset` |
| `00009_add_google_login.sql` | 3 | `User.googleSub`, `User.password` เป็น NULL ได้ |
| `00010_create_likes_comments.sql` | 4 | `LiftLike`, `LiftComment`, index feed |
| `00011_create_notifications.sql` | 5 | `Notification` |
| (ถ้าเลือก B4.5 ก) `00012_cascade_user_lifts.sql` | 4–5 | FK `LiftRecord.userId` → `ON DELETE CASCADE` |

**ขั้นตอนทุกครั้ง:** เขียนไฟล์ → รีวิว → **ขออนุมัติ** → รันกับ DB dev → ทดสอบ → **ขออนุมัติอีกครั้ง** → รันกับ DB prod

---

## Environment variables

| ตัวแปร | ใช้ตอน | dev | prod |
|---|---|---|---|
| `DATABASE_URL` | มีอยู่แล้ว | Neon dev (pooled) | Neon prod (pooled) |
| `DIRECT_DATABASE_URL` | มีอยู่แล้ว (goose) | Neon dev (direct) | Neon prod (direct) |
| `PORT` | มีอยู่แล้ว | `3001` | host ตั้งให้ |
| `COOKIE_SECURE` | มีอยู่แล้ว | `false` | ไม่ใส่ (default = secure) |
| `COOKIE_DOMAIN` | B1.4 | ว่าง | `.<domain>` |
| `APP_ORIGIN` | B1.3 | `http://localhost:3000` | `https://app.<domain>` |
| `TRUSTED_PLATFORM_HEADER` | B1.2 | ว่าง | ตาม docs ของ host |
| `RESEND_API_KEY` | B3.2 | key ทดสอบ | key จริง |
| `EMAIL_FROM` | B3.2 | — | `Kheun Lift <noreply@<domain>>` |
| `GOOGLE_CLIENT_ID` / `GOOGLE_CLIENT_SECRET` | B3.7 | จาก Google Console | ตัวเดียวกันได้ |
| `GOOGLE_REDIRECT_URL` | B3.7 | `http://localhost:3001/auth/google/callback` | `https://api.<domain>/auth/google/callback` |

ทุกตัวต้องอยู่ใน `.env.example` (ค่าว่าง) — **ห้าม commit `.env`**

---

## Endpoint ใหม่ทั้งหมด

| Method | Path | สิทธิ์ | สัปดาห์ |
|---|---|---|---|
| GET | `/weight-classes` | public (ย้ายจาก protected) | 1 |
| GET | `/admin/lifts?status=&cursor=` | admin | 2 |
| PATCH | `/admin/lifts/:id/status` | admin | 2 |
| GET | `/users/:id/achievements` | login | 2 |
| POST | `/auth/password/forgot` | public + rate limit | 3 |
| POST | `/auth/password/reset` | public + rate limit | 3 |
| GET | `/auth/google` | public | 3 |
| GET | `/auth/google/callback` | public | 3 |
| GET | `/feed?cursor=` | login | 4 |
| PUT / DELETE | `/lifts/:id/like` | login | 4 |
| GET / POST | `/lifts/:id/comments` | login | 4 |
| DELETE | `/comments/:id` | เจ้าของหรือ admin | 4 |
| GET | `/notifications?cursor=` | login | 5 |
| GET | `/notifications/unread-count` | login | 5 |
| POST | `/notifications/read` | login | 5 |

---

## ความเสี่ยง

| ความเสี่ยง | ผลกระทบ | รับมือ |
|---|---|---|
| Cookie / CORS / proxy ไม่ทำงานบน host จริง | login ไม่ได้ทั้งเว็บ | B1.1 spike วันแรก |
| เขียน Go ช้ากว่าที่คิด | ไม่ทันแผน | ตัดตาม[ลำดับด้านล่าง](#ถ้าไม่ทัน-ตัดอะไรก่อน), ลอกโครงจาก `lifts/` |
| Host ของ Go แบบฟรีหลับเมื่อไม่มีคนใช้ | request แรกช้าหลายวินาที | ยอมรับได้สำหรับ portfolio หรือจ่ายแผนถูกสุด |
| Neon Free compute หมดกลางเดือน | DB หยุด | ดู usage ใน console ทุกสัปดาห์ |
| Resend ยังไม่ยืนยัน domain | OTP ส่งได้แค่หาตัวเอง | ตั้ง DNS ตั้งแต่ B0.2 |
| ลืมขออนุมัติก่อนรัน migration กับ prod | ข้อมูลจริงเสีย | ทำตามขั้นตอนใน[หัวข้อ migration](#migration-ทั้งหมด) |

---

## ถ้าไม่ทัน ตัดอะไรก่อน

ตัดจากบนลงล่าง (อันบนสุดตัดก่อน):

1. แจ้งเตือน like / comment (B5.3 ส่วนท้าย)
2. Notification ทั้งหมด (B5.1–B5.3)
3. Comment (เหลือแค่ like)
4. Google OAuth (B3.5–B3.8)
5. Feed ทั้งหมด

**ห้ามตัด:** สัปดาห์ 1 (ต่อ API + deploy), Admin (B2.1–B2.5), README + บัญชี demo
