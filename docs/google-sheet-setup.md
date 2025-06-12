# 🧾 ตั้งค่า Google Sheet API สำหรับ Telegram Expense Bot

เอกสารนี้สอนการสร้าง Google Sheet + เปิด Google Sheets API + สร้าง credentials เพื่อเชื่อมกับ Telegram Bot ที่ใช้ Go

---

## ✅ สิ่งที่ต้องเตรียม

- Google Account
- มีสิทธิ์สร้าง Google Sheet
- มีสิทธิ์สร้างโปรเจกต์ใน [Google Cloud Console](https://console.cloud.google.com)

---

## 🪜 ขั้นตอนทั้งหมด

### 1. สร้าง Google Sheet

1. ไปที่ [Google Sheets](https://docs.google.com/spreadsheets/u/0/)
2. คลิก "Blank" สร้างชีทใหม่
3. ตั้งชื่อเช่น `line-expense-bot`
4. เปลี่ยนชื่อแท็บด้านล่างเป็น `Sheet1` หรือ `Expenses` ตามที่ใช้ในโค้ด
5. คัดลอก **ลิงก์ Google Sheet** ไว้ เช่น: https://docs.google.com/spreadsheets/d/1AbCDEFGHIJKLmnOxyz1234567890abcdef 
> 🔹 ID คือส่วนหลัง `/d/` เช่น `1AbCDEFGHIJKLmnOxyz1234567890abcdef`

---

### 2. เปิด Google Sheets API

1. ไปที่ [Google Cloud Console](https://console.cloud.google.com/)
2. คลิก “Select Project” แล้วสร้าง Project ใหม่ เช่น `expense-bot`
3. ไปที่เมนู `APIs & Services > Enabled APIs & services`
4. คลิก “+ ENABLE APIS AND SERVICES”
5. ค้นหา “Google Sheets API” แล้วกด “Enable”

---

### 3. สร้าง Service Account

1. ไปที่ `IAM & Admin > Service Accounts`
2. คลิก “+ Create Service Account”
3. กรอก:
   - ชื่อ: `expense-bot-service`
   - ID: `expense-bot-service`
   - คำอธิบาย: `For accessing Google Sheet`
4. กด “Create and Continue”
5. **ในขั้นตอนให้สิทธิ์**: เลือก Role = `Editor`
6. คลิก “Done”

---

### 4. สร้าง Credentials (credentials.json)

1. ในหน้า Service Accounts > คลิกชื่อที่สร้างไว้
2. ไปที่แท็บ `Keys`
3. คลิก “Add Key” → “Create new key”
4. เลือก `JSON` → กด “Create”
5. จะได้ไฟล์ `credentials.json` ดาวน์โหลดมาไว้ในโปรเจกต์ของคุณ  
   > ❗ ห้าม push ขึ้น Git

---

### 5. แชร์ Google Sheet ให้กับ Service Account

1. เปิด Google Sheet ที่คุณสร้างไว้
2. คลิก “Share”
3. คัดลอก **email ของ Service Account** (เช่น `expense-bot-service@your-project-id.iam.gserviceaccount.com`)
4. วางลงในช่อง "Add people and groups" แล้วให้สิทธิ์เป็น `Editor`
5. กด “Send”

---

## ✅ เพิ่มค่าลงใน `.env` เช่น **ตัวอย่างด้านล่าง**

```dotenv
SPREADSHEET_ID=1AbCDEFGHIJKLmnOxyz1234567890abcdef
SPREADSHEET_LINK=https://docs.google.com/spreadsheets/d/1AbCDEFGHIJKLmnOxyz1234567890abcdef

---
