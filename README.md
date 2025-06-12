# Telegram Expense Bot 💰

บันทึกรายรับรายจ่ายผ่าน Telegram + Google Sheets
เขียนด้วยภาษา Go รองรับคำสั่ง:

- `/summary` — สรุปรายวัน
- `/month` — สรุปรายเดือน
- `/export` — ส่งลิงก์และแนบ Excel

## วิธีใช้งาน

1. สร้าง Telegram Bot (ผ่าน BotFather)
2. สร้าง Google Sheet และเปิด API ตามขั้นตอนใน [docs]
3. สร้าง `.env` (จาก `.env.example`)
4. รัน:

```bash
go run main.go
```

## 🧪 ทดสอบ
1. รันบอทด้วย go run main.go
2. พิมพ์ข้อความเช่น กาแฟ 50 #กาแฟ ใน Telegram
3. ดูว่า Google Sheet มีแถวเพิ่ม ✅

---