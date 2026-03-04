# ME Bot — LINE OA Check-in/Check-out System

## Stack
- **Backend**: Go 1.22 + Gin
- **Database**: MySQL 8
- **Deploy**: Docker + Docker Compose บน Digital Ocean

---

## ขั้นตอนการรันบนเครื่องตัวเอง (Local Development)

### 1. ติดตั้ง Docker Desktop บน macOS
ดาวน์โหลดจาก https://www.docker.com/products/docker-desktop/

### 2. Clone โปรเจกต์และตั้งค่า
```bash
# Copy config
cp .env.example .env

# แก้ไข .env ใส่ค่าจาก LINE Developers Console
# LINE_CHANNEL_SECRET=...
# LINE_CHANNEL_ACCESS_TOKEN=...
```

### 3. รันด้วย Docker Compose
```bash
docker compose up --build
```

ตรวจสอบว่า server ขึ้นแล้ว:
```bash
curl http://localhost:8080/health
# ผลลัพธ์: {"service":"ME Bot","status":"ok"}
```

---

## ขั้นตอนการทดสอบ LINE Webhook บน Local (ใช้ ngrok)

LINE ต้องการ HTTPS เพื่อส่ง webhook มาหา server
ขณะ dev บน local ให้ใช้ ngrok เพื่อ tunnel

### 1. ติดตั้ง ngrok
```bash
brew install ngrok
ngrok config add-authtoken <your_ngrok_token>  # สมัครฟรีที่ ngrok.com
```

### 2. เปิด tunnel
```bash
ngrok http 8080
# จะได้ URL เช่น: https://xxxx.ngrok-free.app
```

### 3. ตั้ง Webhook URL ใน LINE Developers Console
- ไปที่ https://developers.line.biz
- เลือก Channel → Messaging API
- Webhook URL: `https://xxxx.ngrok-free.app/webhook`
- กด Verify ✅

---

## เพิ่มข้อมูลเริ่มต้น (Shop + Employee)

ตอนนี้ยังไม่มี admin UI ให้รันคำสั่ง SQL ตรงๆ ก่อน:

```bash
# เข้า MySQL container
docker compose exec mysql mysql -u mebot -p mebot_db

# เพิ่มร้าน (ใส่พิกัดและ LINE Group ID จริง)
INSERT INTO shops (name, lat, lng, radius_m, line_group_id)
VALUES ('ร้านกาแฟสาขา 1', 13.7563, 100.5018, 200, 'C_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx');

# ดู shop id ที่เพิ่งเพิ่ม
SELECT id FROM shops;

# เพิ่มพนักงาน (ใส่ LINE User ID ของพนักงานจริง)
INSERT INTO employees (line_user_id, name, shop_id)
VALUES ('Uxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx', 'สมชาย ใจดี', 1);
```

### วิธีหา LINE User ID ของพนักงาน
เปิด LINE Developers Console → Messaging API → เปิด "Use webhooks"
แล้วให้พนักงาน Add bot และส่งข้อความมา จะเห็น userId ใน Webhook log

### วิธีหา LINE Group ID
Add bot เข้า group แล้วดู log ใน server:
```bash
docker compose logs app | grep "Bot joined"
```

---

## Deploy บน Digital Ocean

### 1. สร้าง Droplet
- OS: Ubuntu 22.04 LTS
- Plan: Basic 2GB RAM ($12/mo)
- เปิด SSH key

### 2. ติดตั้ง Docker บน Droplet
```bash
ssh root@<droplet-ip>
curl -fsSL https://get.docker.com | sh
```

### 3. ติดตั้ง Nginx + SSL
```bash
apt install -y nginx certbot python3-certbot-nginx

# ตั้งค่า Nginx
cat > /etc/nginx/sites-available/mebot << 'EOF'
server {
    listen 80;
    server_name yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
EOF

ln -s /etc/nginx/sites-available/mebot /etc/nginx/sites-enabled/
nginx -t && systemctl reload nginx

# ออก SSL Certificate
certbot --nginx -d yourdomain.com
```

### 4. Deploy App
```bash
# บน Droplet
git clone https://github.com/yourrepo/me-bot.git
cd me-bot
cp .env.example .env
# แก้ไข .env ใส่ค่าจริง
nano .env

docker compose up -d --build
```

### 5. ตั้ง Webhook URL ใน LINE
```
https://yourdomain.com/webhook
```

---

## โครงสร้างโปรเจกต์

```
me-bot/
├── cmd/server/main.go          ← จุดเริ่มต้นโปรแกรม
├── internal/
│   ├── config/config.go        ← โหลด .env
│   ├── database/mysql.go       ← เชื่อม MySQL + auto migrate
│   ├── handler/webhook.go      ← รับ LINE events
│   ├── service/
│   │   ├── checkin.go          ← business logic check-in/out
│   │   └── location.go         ← คำนวณระยะห่าง
│   ├── repository/
│   │   ├── employee.go         ← query ข้อมูลพนักงาน
│   │   └── attendance.go       ← query ข้อมูลการเข้างาน
│   └── model/model.go          ← struct ของ database tables
├── .env.example
├── Dockerfile
├── docker-compose.yml
└── README.md
```
