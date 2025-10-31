## مرحله ی builder
#استفاده از نسخه ی سبک گو برای ایجاد برنامه
FROM golang:1.25 AS builder
#پوشه ی کاری داکر
WORKDIR /app
#کپی کردن فایل های پروژه به داخل داکر
COPY . .
#کامپایل برنامه و ساخت فایل اجرایی
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server .
## مرحله ی runing
# استفاده از نسخه ی سبک سیستم عامل برای اجرای برنامه
FROM alpine:3.18
#اتصال درست به سرور ها
RUN apk add --no-cache ca-certificates
#تنظیم مجدد پوشه ی کاری داکر
WORKDIR /app
#ایجاد پوشه ی دیتا با دسترسی کامل برای خواندن و نوشتن
RUN mkdir -p /app/data && chmod 777 /app/data
#کپی فایل اجرایی از مرحله ی builder به این مرحله
COPY --from=builder /app/server .
#اعلام اینکه برنامه روی پورت 8080 اجرا میشود
EXPOSE 18080
#اجرای فایل اجرایی
CMD ["./server"]