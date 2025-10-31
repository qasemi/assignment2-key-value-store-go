# assignment2-key-value-store-go
Design and Implementation of a Simple Persistent Key-Value Store Database Using Golang

هدف این پروژه طراحی و پیاده‌سازی یک **پایگاه‌داده کلید-مقدار (Key-Value Store)** ساده است.
در این سیستم، داده‌ها با استفاده از یک کلید ذخیره می‌شوند و با همان کلید قابل بازیابی هستند.  
ویژگی مهم این پروژه **پایداری داده‌ها (Persistence)** است، به این معنی که اطلاعات پس از بستن برنامه نیز از بین نمی‌روند.

---

## نحوه‌ی اجرا

### اجرای مستقیم (در صورت نصب Go)
```bash
go run main.go
```

### اجرای پروژه با Docker
```bash
docker build -t kvstore:02 .
docker run -p 18080:18080 kvstore:02
```

---

## ساختار پوشه‌ها
```
 ├── main.go         # کد اصلی برنامه
 ├── Store.json      # فایل ذخیره‌سازی داده‌ها
 ├── Dockerfile      # پیکربندی اجرای پروژه در Docker
 ├── go.mod          # اطلاعات ماژول و وابستگی‌های Go
 └── README.md       # فایل راهنما
```

---

## دستورات PowerShell
```powershell
cd C:\Users\pmiyo\Desktop\Excercise\02
docker build -t my-kvstore:latest .
docker run -d --name kv -p 18080:18080 -v kv_data:/app/data my-kvstore:latest
docker ps
docker volume create kv_data
```

---

## ذخیره‌سازی داده (درخواست PUT)
```powershell
Invoke-RestMethod -Method Put `
  -Uri "http://localhost:18080/store" `
  -ContentType "application/json" `
  -Body '{"key":"user:2001","value":{"name":"Sara Ahmadi","age":29,"email":"s.ahmadi@example.com"}}'

Invoke-RestMethod -Method Put `
  -Uri "http://localhost:18080/store" `
  -ContentType "application/json" `
  -Body '{"key":"user:2002","value":{"name":"Reza Karimi","age":35,"email":"r.karimi@example.com"}}'

Invoke-RestMethod -Method Put `
  -Uri "http://localhost:18080/store" `
  -ContentType "application/json" `
  -Body '{"key":"user:2003","value":{"name":"Ali Mohammadi","age":41,"email":"a.mohammadi@example.com"}}'
```
 **خروجی مورد انتظار:**  
در صورت موفقیت، هیچ خروجی بازگردانده نمی‌شود.

---

## بازیابی داده‌ها (درخواست GET)
```powershell
Invoke-RestMethod -Method Get -Uri "http://localhost:18080/objects/user:2001"
Invoke-RestMethod -Method Get -Uri "http://localhost:18080/objects/user:2002"
Invoke-RestMethod -Method Get -Uri "http://localhost:18080/objects/user:2003"
```

**نمونه خروجی:**
```
name          age  email
----          ---  -----
Ali Mohammadi 41   a.mohammadi@example.com
```

---

## مشاهده‌ی تمام داده‌ها
```powershell
Invoke-RestMethod -Method Get -Uri "http://localhost:18080/objects" | ConvertTo-Json -Depth 5
```

**نمونه خروجی:**
```json
{
    "user:2001": {"name":"Sara Ahmadi","age":29,"email":"s.ahmadi@example.com"},
    "user:2002": {"name":"Reza Karimi","age":35,"email":"r.karimi@example.com"},
    "user:2003": {"name":"Ali Mohammadi","age":41,"email":"a.mohammadi@example.com"}
}
```

---

## تست پایداری داده‌ها در Docker
```bash
# توقف و حذف کانتینر
docker stop kv
docker rm kv

# اجرای مجدد کانتینر با همان Volume
docker run -d --name kv -p 18080:18080 -v kv_data:/app/data my-kvstore:latest

# بررسی ماندگاری داده‌ها
Invoke-RestMethod -Method Get -Uri "http://localhost:18080/objects/user:2001"
Invoke-RestMethod -Method Get -Uri "http://localhost:18080/objects/user:2002"
Invoke-RestMethod -Method Get -Uri "http://localhost:18080/objects/user:2003"
```

---

## حذف کامل داده‌ها
```bash
docker volume rm kv_data
```

