# 🧾 Set Up Google Sheets API for Telegram Expense Bot

This guide walks you through creating a Google Sheet, enabling the Google Sheets API, and generating `credentials.json` for integration with the Go-based Telegram Expense Bot.

---

## ✅ Prerequisites

- A Google account
- Permission to create Google Sheets
- Permission to create projects in [Google Cloud Console](https://console.cloud.google.com)

---

## 🪜 Steps

### 1. Create a Google Sheet

1. Visit [Google Sheets](https://docs.google.com/spreadsheets/u/0/)
2. Click “Blank” to create a new sheet
3. Name it something like `telegram-expense-bot`
4. Rename the tab to `Sheet1` or `Expenses` (as used in the code)
5. Copy the **Google Sheet link**, e.g.:  
   `https://docs.google.com/spreadsheets/d/1AbCDEFGHIJKLmnOxyz1234567890abcdef`  
   > 🔹 The **Spreadsheet ID** is the part after `/d/`, e.g. `1AbCDEFGHIJKLmnOxyz1234567890abcdef`

---

### 2. Enable Google Sheets API

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Click “Select Project” → create a new project (e.g., `expense-bot`)
3. Go to `APIs & Services > Enabled APIs & services`
4. Click “+ ENABLE APIS AND SERVICES”
5. Search for `Google Sheets API` and click “Enable”

---

### 3. Create a Service Account

1. Go to `IAM & Admin > Service Accounts`
2. Click “+ Create Service Account”
3. Fill in:
   - Name: `expense-bot-service`
   - ID: `expense-bot-service`
   - Description: `For accessing Google Sheet`
4. Click “Create and Continue”
5. **On the permissions step**, select Role = `Editor`
6. Click “Done”

---

### 4. Generate Credentials (credentials.json)

1. From the Service Accounts list, click the one you created
2. Go to the `Keys` tab
3. Click “Add Key” → “Create new key”
4. Choose `JSON` → Click “Create”
5. A file named `credentials.json` will be downloaded — move it to your project folder  
   > ❗ **Never commit this file to Git**

---

### 5. Share the Google Sheet with the Service Account

1. Open the Google Sheet you created
2. Click “Share”
3. Copy the **email address** of the service account  
   (e.g., `expense-bot-service@your-project-id.iam.gserviceaccount.com`)
4. Paste it into the “Add people and groups” field
5. Grant `Editor` access and click “Send”

---

## ✅ Add to your `.env` file

Example:

```dotenv
SPREADSHEET_ID=1AbCDEFGHIJKLmnOxyz1234567890abcdef
SPREADSHEET_LINK=https://docs.google.com/spreadsheets/d/1AbCDEFGHIJKLmnOxyz1234567890abcdef
TELEGRAM_BOT_TOKEN=your-telegram-bot-token
