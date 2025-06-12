# 📲 Telegram Expense Bot 💰  
ระบบจดรายรับรายจ่ายผ่าน Telegram แล้วบันทึกอัตโนมัติลง Google Sheet

---

## 🔧 วิธีติดตั้ง

### 1. สร้าง Telegram Bot
- เปิด [BotFather](https://t.me/botfather)
- สร้างบอทใหม่และรับ `BOT_TOKEN`

### 2. สร้าง Google Sheet และเปิด API
ดูขั้นตอนใน [📄 docs/setup-google-sheet.md](docs/setup-google-sheet.md)

### 3. สร้างไฟล์ `.env` จาก `.env.example`
```bash
cp .env.example .env

4. install task
    - scoop install go-task (windows)
    - brew install go-task/tap/go-task (Mac)

5. รัน:

```bash
task dev
```

## 🧪 ทดสอบ
1. รันบอทด้วย task dev
2. ส่งข้อความใน Telegram เช่น: กาแฟ 50 #กาแฟ 
3. ดูว่า Google Sheet มีแถวเพิ่ม ✅

---