# Rate Limit Middleware - Flow Documentation

## 1. Tổng quan
Middleware này dùng để giới hạn số lượng request theo từng IP.
Mỗi IP sẽ có một rate limiter riêng (token bucket).

---

## 2. Cấu trúc dữ liệu

### client
- limiter: Bộ giới hạn tốc độ (rate.Limiter)
- lastSeen: Thời điểm request gần nhất
- mu: Mutex để đảm bảo thread-safe khi cập nhật dữ liệu

### clients (global)
- Kiểu: sync.Map
- Key: IP (string)
- Value: *client

---

## 3. Luồng xử lý chính

### Bước 1: Lấy IP client
- Sử dụng `c.ClientIP()` để xác định IP gửi request

---

### Bước 2: Load hoặc tạo client mới
- Gọi `clients.LoadOrStore(ip, newClient)`
- Nếu IP chưa tồn tại:
  → Tạo mới client với:
    - limiter = rate.NewLimiter(requestsPerSecond, burst)
    - lastSeen = now
- Nếu đã tồn tại:
  → Lấy client cũ

---

### Bước 3: Cập nhật lastSeen (thread-safe)
- Lock mutex của client
- Cập nhật:
  - lastSeen = time.Now()
  - lấy limiter ra local variable
- Unlock mutex

---

### Bước 4: Kiểm tra rate limit
- Gọi `limiter.Allow()`

#### Nếu Allow = false:
- Trả về HTTP 429 (Too Many Requests)
- Body:
  {
    "message": "Rate limit exceeded. Please try again later."
  }
- Abort request

#### Nếu Allow = true:
- Cho request đi tiếp (`c.Next()`)

---

## 4. Cơ chế cleanup client cũ

### Goroutine chạy nền (init)

- Mỗi 1 phút:
  - Duyệt toàn bộ clients (Range)

### Với mỗi client:
- Lock mutex
- Kiểm tra:
  time.Since(lastSeen) > 3 phút ?
- Unlock mutex

#### Nếu inactive:
- Xóa khỏi map:
  `clients.Delete(ip)`

---

## 5. Tổng kết flow

Request đến
   ↓
Lấy IP
   ↓
LoadOrStore client
   ↓
Update lastSeen (lock)
   ↓
Check limiter.Allow()
   ↓
 ┌───────────────┬────────────────┐
 │ Allow = true  │ Allow = false  │
 │               │                │
 │ c.Next()      │ Return 429     │
 └───────────────┴────────────────┘

Background:
- Cleanup mỗi 1 phút
- Xóa client không hoạt động > 3 phút

---

## 6. Đặc điểm quan trọng

- Thread-safe:
  - sync.Map cho concurrent access
  - mutex cho từng client

- Memory optimization:
  - Cleanup client inactive để tránh leak

- Rate limiting strategy:
  - Token Bucket

- Granularity:
  - Theo từng IP

---