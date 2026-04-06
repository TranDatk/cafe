# Hướng dẫn sử dụng go-redis v9

Tài liệu này cung cấp hướng dẫn chi tiết về cách sử dụng thư viện `go-redis/v9`, thư viện Redis client phổ biến nhất cho ngôn ngữ lập trình Go.

## 1. Cài đặt
Để bắt đầu, hãy cài đặt thư viện bằng lệnh:
```bash
go get github.com/redis/go-redis/v9
```

## 2. Khởi tạo Kết nối
Trong phiên bản v9, mọi thao tác đều yêu cầu một `context.Context`.

```go
import (
    "context"
    "github.com/redis/go-redis/v9"
    "time"
)

var ctx = context.Background()

func Connect() *redis.Client {
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // Để trống nếu không có password
        DB:       0,  // Database mặc định
    })
    return rdb
}
```

## 3. Các thao tác cơ bản (String)

### 3.1 SET - Lưu giá trị
Lưu một giá trị vào Redis kèm theo thời gian hết hạn (TTL).
```go
// TTL = 0 nghĩa là không bao giờ hết hạn
err := rdb.Set(ctx, "user:1", "Antigravity", 1 * time.Hour).Err()
```

### 3.2 GET - Lấy giá trị
```go
val, err := rdb.Get(ctx, "user:1").Result()
if err == redis.Nil {
    fmt.Println("Key không tồn tại")
} else if err != nil {
    panic(err)
}
fmt.Println("Giá trị:", val)
```

### 3.3 DEL - Xóa Key
```go
err := rdb.Del(ctx, "user:1").Err()
```

### 3.4 EXISTS - Kiểm tra tồn tại
```go
count, err := rdb.Exists(ctx, "user:1").Result()
if count > 0 {
    fmt.Println("Key tồn tại")
}
```

## 4. Các cấu trúc dữ liệu khác

### 4.1 Hashes (HSET / HGET)
Phù hợp để lưu trữ đối tượng.
```go
// Lưu nhiều trường
err := rdb.HSet(ctx, "session:123", map[string]interface{}{
    "user_id": "1",
    "ip":      "127.0.0.1",
}).Err()

// Lấy 1 trường
userID, _ := rdb.HGet(ctx, "session:123", "user_id").Result()
```

### 4.2 Lists (LPUSH / RPOP)
Dùng làm Queue.
```go
// Thêm vào đầu list
rdb.LPush(ctx, "queue", "task1", "task2")

// Lấy từ cuối list (Blocking RPop)
result, err := rdb.BRPop(ctx, 5*time.Second, "queue").Result()
```

### 4.3 Sets (SADD / SMEMBERS)
Tập hợp các giá trị không trùng lặp.
```go
rdb.SAdd(ctx, "tags", "go", "redis", "coding")
members, _ := rdb.SMembers(ctx, "tags").Result()
```

### 4.4 Sorted Sets (ZADD / ZRANGE)
Tập hợp có sắp xếp theo Score.
```go
rdb.ZAdd(ctx, "leaderboard", redis.Z{Score: 100, Member: "Player1"})
rdb.ZAdd(ctx, "leaderboard", redis.Z{Score: 200, Member: "Player2"})

// Lấy top 10
users, _ := rdb.ZRangeWithScores(ctx, 0, 9).Result()
```

## 5. Pipelining & Transactions

### 5.1 Pipeline
Dùng để gửi nhiều lệnh cùng lúc nhằm giảm độ trễ mạng (Network Latency).
```go
pipe := rdb.Pipeline()

incr := pipe.Incr(ctx, "counter")
pipe.Expire(ctx, "counter", time.Hour)

_, err := pipe.Exec(ctx)
fmt.Println(incr.Val())
```

### 5.2 Transactions (Watch/Multi/Exec)
Đảm bảo tính nguyên tử (Atomicity).
```go
err := rdb.Watch(ctx, func(tx *redis.Tx) error {
    n, _ := tx.Get(ctx, "key").Int()
    
    _, err := tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
        pipe.Set(ctx, "key", n+1, 0)
        return nil
    })
    return err
}, "key")
```

## 6. Các hàm bổ trợ quan trọng
- `Expire(ctx, key, expiration)`: Đặt lại TTL cho key đã tồn tại.
- `TTL(ctx, key)`: Kiểm tra thời gian còn lại của key.
- `Keys(ctx, pattern)`: Tìm các key theo pattern (Cẩn thận khi dùng với DB lớn).
- `FlushDB(ctx)`: Xóa sạch dữ liệu trong DB hiện tại.

---
**Lưu ý**: Luôn kiểm tra lỗi `redis.Nil` khi thực hiện các lệnh `Get` để tránh nhầm lẫn giữa lỗi kết nối và việc không tìm thấy dữ liệu.
