# Kheun Lift — แผนหลัก (Go + ภาพรวม)

งานหน้าเว็บอยู่ที่ [`../kheun-lift-frontend-next/PLAN.md`](../kheun-lift-frontend-next/PLAN.md) (เลข Phase ตรงกัน)

**เป้าหมาย:** deploy จริงบน domain ของตัวเอง ใช้เป็นผลงานใน resume

---

## สถาปัตยกรรม

```
app.<domain> (Next, Vercel) ──► api.<domain> (Go, Railway) ──► Neon Postgres (Singapore)
                                          └──► Resend (email), Google OAuth, Claude API
```

- Session cookie `kheun_session` (httpOnly, SameSite=Lax, Domain=.<domain>) ไม่ใช้ JWT
- CORS + CrossOriginProtection ยอมแค่ `APP_ORIGIN` origin เดียว
- login / register / forgot password ต้องเรียกจาก browser ตรง ๆ ไม่งั้น Go เห็น IP ของ Vercel และ rate limit นับทุกคนรวมกัน

## การตัดสินใจ

| เรื่อง | เลือก | เหตุผล |
|---|---|---|
| เชื่อม Next กับ Go | domain ของตัวเอง แยก app. / api. | same-site, cookie ส่งข้ามได้ |
| Host | Vercel / Railway Hobby / Neon prod ใหม่ | Railway ไม่หลับ (~$5/เดือน), region Singapore |
| Domain + DNS | Cloudflare Registrar | ราคาทุน, DNS อยู่ที่เดียวกับ R2 |
| Email | Resend ผ่าน `net/http` | ไม่ต้องลง SDK |
| Repost | ไม่ทำ | ซับซ้อน แต่ไม่โชว์ฝีมือเพิ่ม |
| Test | เฉพาะจุดเสี่ยง, ต้องอนุมัติ | auth, OTP, ranking |
| Reject | ต้องมีเหตุผล | ผู้ใช้ต้องรู้ว่าทำไมไม่ผ่าน |
| AI | ช่วย admin เท่านั้น ไม่ verify เอง | AI ตัดสินความลึก / lockout ไม่ได้ |
| หน้าเว็บ | mobile first, ต้อง login ก่อนทุกหน้า | ผู้ใช้หลักเข้าจากมือถือ ไม่มีหน้า public |

## กติกาโค้ด Go

- ลอกโครงจาก `internal/lifts/`: queries → model → repository → service → handler → routes
- ทุก request ที่แตะ DB ใส่ timeout 5 วิ, อ่าน body ด้วย `decodeBody`, id ใช้ `parseID`
- error ที่ผู้ใช้ควรรู้ → `var ErrXxx` ใน service; error อื่น → log แล้วตอบ 500 (ห้ามส่ง error ของ DB ออกไป)
- รูปแบบ response: `{"message": "...", "data": ...}`
- list ที่ยาวใช้ cursor (`?cursor=<id>`, ทีละ 20, ตอบ `nextCursor`) ไม่ใช้ OFFSET

## สถานะ (2026-11-15)

- ✅ register / login / logout / session, users (me, แก้, ลบ), exercises, lifts CRUD, leaderboards, weight-classes, CORS + trusted origin (ทดสอบจาก `localhost:3002` แล้ว), seed แยกใน `make seed`
- ❓ โฟลเดอร์ว่าง `internal/email`, `mailjob`, `verification` และ migration 00004 ที่จองไว้: ยังไม่ได้ตัดสินใจว่าใช้ทำอะไร

---

## Phase 1 — Deploy

- [ ] ซื้อ domain ที่ Cloudflare, สมัคร Railway (Hobby, Singapore) + Resend + Google Cloud, สร้าง Neon prod
  - DNS (เมฆสีเทา): `app` → Vercel, `api` → Railway, `img` → R2 prod
  - dev DB มี user ปลอม `password123` ห้ามใช้เป็น prod
- [x] รวมโค้ดตั้ง cookie เป็นฟังก์ชันเดียว ใช้ทั้ง login / logout / ลบบัญชี + เพิ่ม `COOKIE_DOMAIN` (`internal/middleware/cookie.go`)
  - ⚠ ต้องตั้ง `COOKIE_DOMAIN=.<domain>` บน prod ไม่งั้น Next มองไม่เห็น cookie; ยังไม่ได้ทดสอบบน domain จริง
- [ ] rate limit เห็น IP จริง: ตอนนี้ `SetTrustedProxies(nil)` → ตั้ง `TrustedPlatform` ตาม docs ของ Railway
- [ ] deploy + `make migrate` กับ prod (ขออนุมัติก่อน, ห้าม `make seed`)
- [ ] สร้าง bucket R2 prod + custom domain + ตั้ง env `R2_*` (bucket dev ใช้ r2.dev ทำแล้ว)
- [x] `POST /users/me/profile-image` (multipart ช่อง `profileImage`, ≤ 2 MB, jpeg/png/webp ตรวจจากเนื้อไฟล์) เก็บที่ R2 key `profile-images/<userId>` เขียนทับรูปเดิม, ลบบัญชีแล้วลบรูปด้วย
  - ทดสอบอัปรูปกับ bucket dev แล้ว เปิดลิงก์เห็นรูป; frontend ยังไม่ได้ทำ
- [x] `DELETE /users/me/profile-image` ลบใน R2 ก่อนแล้วตั้ง `profileImage` = NULL, กดซ้ำได้ 200
  - ⚠ ยังไม่ได้ทดสอบกับ R2 จริง

**เสร็จเมื่อ:** สมัคร → login → กรอกโปรไฟล์ → ส่ง lift ได้บน app.<domain>, logout แล้ว cookie หาย

## Phase 2 — Admin ตรวจ lift (ห้ามตัด)

ตอนนี้ verify ได้ทางเดียวคือแก้ DB เอง leaderboard จึงว่าง และ feature อื่นทั้งหมดต้องใช้ lift ที่ verified

- [ ] migration: `User.role` (user / admin), `LiftRecord.rejectReason`
  - ตั้ง admin คนแรกด้วย SQL ใน Neon console (ไม่ทำ endpoint ตั้ง admin)
- [ ] session และ `/users/me` ส่ง `role` กลับมา
- [ ] middleware `RequireAdmin` → 403
- [ ] `GET /admin/lifts?status=pending&cursor=` เก่าสุดก่อน
- [ ] `PATCH /admin/lifts/:id/status` `{status, reason}`
  - เปลี่ยนได้เฉพาะ lift ที่ยัง pending (ไม่ใช่ → 409)
  - reject ต้องมี reason 1–300 ตัว, verify ห้ามมี reason
  - ทำใน service + transaction (Phase 7 จะเพิ่ม notification ตรงนี้)

**เสร็จเมื่อ:** verify แล้ว lift ขึ้น leaderboard, reject แล้วเจ้าของเห็นเหตุผลใน `/lifts/me`

## Phase 3 — AI ช่วยตรวจ lift

ไม่ train model เอง ใช้ model สำเร็จรูปผ่าน API, admin กด verify เองเสมอ

- [ ] เช็กแบบไม่ใช้ AI ตอนสร้าง lift → ติดธง
  - วิดีโอมีอยู่จริง (YouTube oEmbed)
  - URL วิดีโอซ้ำกับ lift อื่น
  - น้ำหนักสูงผิดปกติเมื่อเทียบกับน้ำหนักตัว
- [ ] ส่งภาพจากวิดีโอให้ Claude ถาม: ท่าตรงไหม, ยกคนเดียวไหม, แผ่นน้ำหนักสอดคล้องกับที่กรอกไหม → เก็บ `aiVerdict` + `aiNote` ส่งไปกับ `/admin/lifts`
- ⚠ Claude รับรูป ไม่รับวิดีโอ, ดาวน์โหลดจาก YouTube ผิดเงื่อนไข → ต้องเลือก: ภาพปก YouTube หรือให้อัปโหลดไฟล์เอง
- รันแบบ background ไม่ให้การส่ง lift ช้า, AI ล่มต้องไม่ทำให้ส่ง lift ไม่ได้

**เสร็จเมื่อ:** admin เห็นธงและความเห็นของ AI ในการ์ดแต่ละอัน

## Phase 4 — Achievement

- [ ] `GET /users/:id/achievements` คำนวณสดทุกครั้ง ไม่สร้างตาราง
- badge: First Rep, Big Three, Bench 1×BW, Squat 1.5×BW, Deadlift 2×BW, 300/400/500 Club, Podium, Champion (นับเฉพาะ verified)
- ใช้น้ำหนักตัวปัจจุบัน (ข้อจำกัด: แก้น้ำหนักแล้ว badge อาจหาย)

**เสร็จเมื่อ:** lift ถูก verify แล้วเห็น badge ใหม่

## Phase 5 — Forgot password + Google login

- [ ] migration `PasswordReset` (1 user มีได้ 1 รหัส), ส่ง email ผ่าน Resend (ต้องยืนยัน domain ก่อน)
- [ ] `POST /auth/password/forgot`: ตอบข้อความเดิมเสมอ, ขอใหม่ได้ทุก 60 วิ, รหัส 6 หลักจาก `crypto/rand`, เก็บแค่ hash, หมดอายุ 10 นาที
- [ ] `POST /auth/password/reset`: ผิดได้ 5 ครั้ง; ถูกแล้วเปลี่ยนรหัส + ลบ session ทั้งหมด ใน transaction เดียว
- [ ] migration `User.googleSub` + password เป็น NULL ได้; login ที่ไม่มีรหัสให้ทำเหมือนไม่เจอ user
- [ ] `GET /auth/google` + `/callback` (`x/oauth2`, state + PKCE)
  - ⚠ ผูกบัญชีเดิมด้วย email ได้เฉพาะเมื่อ `email_verified = true`
- [ ] แยก `CreateSession` ใช้ร่วมกันระหว่าง login กับ Google

**เสร็จเมื่อ:** reset รหัสผ่านด้วย email จริงได้, login ด้วย Google ได้ทั้งบัญชีใหม่และบัญชีเดิม

## Phase 6 — Feed + Like + Comment

feed = lift ที่ verified เรียงใหม่สุดก่อน ไม่สร้างตาราง post แยก

- [ ] migration `LiftLike` (PK liftId+userId), `LiftComment` (1–500 ตัว), index สำหรับ feed
- [ ] `GET /feed?cursor=` มี likeCount, commentCount, likedByMe ใน query เดียว (ห้าม N+1)
- [ ] `PUT` / `DELETE /lifts/:id/like` กดซ้ำไม่ error, ตอบ likeCount ล่าสุด
- [ ] `GET` / `POST /lifts/:id/comments`, `DELETE /comments/:id` (เจ้าของหรือ admin)
- [ ] ตัดสินใจ: ลบบัญชีแล้วลบ lift ตามไปด้วยไหม (ตอนนี้ลบไม่ได้ถ้ายังมี lift)

**เสร็จเมื่อ:** ไถ feed, like, comment, ลบ comment ตัวเองได้

## Phase 7 — Notification

อันดับเปลี่ยนตอน admin verify เท่านั้น → สร้าง notification ตรงนั้นจุดเดียว ไม่ต้องมี cron

- [ ] migration `Notification` (เก็บ `type` + `data` JSON, หน้าเว็บสร้างข้อความเอง)
- [ ] ตอน verify / reject (transaction เดียวกัน): `lift_verified`, `lift_rejected`, `rank_dropped`
- [ ] `GET /notifications?cursor=`, `GET /notifications/unread-count`, `POST /notifications/read`

**เสร็จเมื่อ:** verify แล้วเจ้าของได้แจ้งเตือน, คนที่อันดับตกได้แจ้งเตือน

## ปิดงาน

- [ ] ลบ session ที่หมดอายุตอน login, รวม SQL หา weight class ให้เหลือชุดเดียว
- [ ] test (ถ้าอนุมัติ): password hash, เงื่อนไข badge, OTP
- [ ] บัญชี demo บน prod ที่มี lift verified + badge
- [ ] README: ลิงก์เว็บจริง, diagram, เหตุผลที่เลือก stack, ความปลอดภัยที่ทำ, วิธีรันบนเครื่อง

---

## Migration

- ไฟล์ใหม่เริ่มที่ `00007` (DB dev เคยรัน seed เป็นเลข 5–6 ถ้าใช้เลขน้อยกว่า goose จะ error)
- ลำดับ: เขียน → ขออนุมัติ → รันกับ dev → ทดสอบ → ขออนุมัติ → รันกับ prod

## Environment variables

| ตัวแปร | dev | prod |
|---|---|---|
| `DATABASE_URL` / `DIRECT_DATABASE_URL` | Neon dev | Neon prod (pooled / direct) |
| `PORT` | `3001` | host ตั้งให้ |
| `COOKIE_SECURE` | `false` | ไม่ใส่ |
| `COOKIE_DOMAIN` | ว่าง | `.<domain>` |
| `APP_ORIGIN` | `http://localhost:3002` | `https://app.<domain>` |
| `TRUSTED_PLATFORM_HEADER` | ว่าง | ตาม docs ของ host |
| `RESEND_API_KEY`, `EMAIL_FROM` | key ทดสอบ | key จริง, `noreply@<domain>` |
| `GOOGLE_CLIENT_ID` / `_SECRET` / `_REDIRECT_URL` | `http://localhost:3001/auth/google/callback` | `https://api.<domain>/auth/google/callback` |
| `ANTHROPIC_API_KEY` | key ทดสอบ | key จริง |
| `R2_ACCOUNT_ID`, `R2_ACCESS_KEY_ID`, `R2_SECRET_ACCESS_KEY`, `R2_BUCKET`, `R2_PUBLIC_URL` | bucket dev (ว่างได้ อัปรูปจะได้ 500) | bucket prod |

ทุกตัวต้องอยู่ใน `.env.example` (ค่าว่าง) — ห้าม commit `.env`

## ถ้าไม่ทัน ตัดจากบนลงล่าง

1. แจ้งเตือน like / comment
2. Notification
3. Comment
4. Google login
5. AI
6. Feed

**ห้ามตัด:** Phase 1, Phase 2, README + บัญชี demo

ต้องแก้ที่ Go (จดไว้ใน PLAN.md หัวข้อ "ต้องแก้ที่ Go" แล้ว)
GET /lifts/:id ควรส่ง lift ที่ยังไม่ verified ให้เฉพาะเจ้าของและ admin คนอื่นควรได้ 404
GET /lifts/users/:id ควรส่งให้คนอื่นเฉพาะ lift ที่ verified