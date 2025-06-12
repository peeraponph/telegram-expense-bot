# 📲 Telegram Expense Bot 💰  
A Telegram bot for recording income and expenses, automatically logging them into Google Sheets.

---

## 🔧 Installation Guide

1. **Create a Telegram Bot**
- Open [BotFather](https://t.me/botfather)
- Create a new bot and obtain the `BOT_TOKEN`

2. **Create a Google Sheet and Enable the API**  
Follow the steps in [📄 docs/google-sheet-setup.md](docs/google-sheet-setup.md)

3. **Create a `.env` file from the example**
```bash
cp .env.example .env
```
4. install task
    - scoop install go-task (Windows)
    - brew install go-task/tap/go-task (Mac)

5. run:

```bash
task dev
```

## 🧪 Testing
1. Run the bot using task dev
2. Send a message to the bot in Telegram, e.g., coffee 50 #drink
3. Check that a new row is added to your Google Sheet ✅

---