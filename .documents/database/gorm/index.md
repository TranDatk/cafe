# Tài liệu GORM Queries

Tài liệu chi tiết về cách xây dựng truy vấn, xử lý quan hệ và tối ưu hiệu năng với GORM trong Go.

## 1. Nguyên lý Chaining (Chuỗi các hàm)
GORM sử dụng cơ chế **Method Chaining**. Mỗi hàm điều kiện (`Where`, `Select`, `Order`) sẽ trả về một đối tượng `*gorm.DB`, cho phép bạn viết tiếp các điều kiện khác.

Truy vấn chỉ thực sự được gửi đến Database khi bạn gọi một **Finisher Method** (như `First`, `Find`, `Create`, `Save`, `Delete`).

```go
// Ví dụ về chuỗi các hàm
db.Table("users").Select("name", "email").Where("age > ?", 20).Order("age desc").Find(&users)
```

## 2. Các hàm truy vấn cơ bản (Finisher Methods)
- **First(&dest)**: Lấy bản ghi đầu tiên theo primary key. Trả về lỗi `ErrRecordNotFound` nếu không tìm thấy.
- **Take(&dest)**: Lấy một bản ghi bất kỳ (không sắp xếp).
- **Find(&dest)**: Lấy danh sách bản ghi và gán vào slice. Nếu không tìm thấy, trả về slice rỗng (không báo lỗi).
- **Last(&dest)**: Lấy bản ghi cuối cùng theo primary key.

## 3. Filtering (Điều kiện lọc)
Sử dụng hàm `.Where()` với Placeholders (`?`) để tránh SQL Injection.

### Lọc đơn giản:
```go
db.Where("name = ?", "jinzhu").First(&user)
```

### Lọc phức tạp với Struct hoặc Map:
Khi truyền Struct vào `Where`, GORM chỉ lọc theo các field có giá trị khác "zero value" (0, "", false). Nếu bạn muốn lọc cả giá trị zero, hãy dùng Map.

```go
// Dùng Struct (Sẽ KHÔNG lọc theo age = 0)
db.Where(&User{Name: "jinzhu", Age: 0}).Find(&users)
// SQL: SELECT * FROM users WHERE name = "jinzhu";

// Dùng Map (Sẽ lọc chính xác age = 0)
db.Where(map[string]interface{}{"Name": "jinzhu", "Age": 0}).Find(&users)
// SQL: SELECT * FROM users WHERE name = "jinzhu" AND age = 0;
```

## 4. Xử lý Quan hệ (Associations)

### Preload (Eager Loading)
Dùng để nạp dữ liệu quan hệ trong các câu lệnh truy vấn riêng biệt (Tránh N+1 query).
```go
db.Preload("Orders").Find(&users)
// 1. SELECT * FROM users;
// 2. SELECT * FROM orders WHERE user_id IN (1,2,3...);
```

### Joins
Sử dụng SQL JOIN để lấy dữ liệu. Thường dùng khi bạn muốn lọc kết quả dựa trên dữ liệu của bảng liên quan.
```go
db.Joins("Company").Find(&users)
```

### Select
Chỉ định các trường (columns) cụ thể muốn lấy lên từ Database.
```go
db.Select("name", "age").Find(&users)
```

## 5. Xử lý lỗi
Luôn kiểm tra thuộc tính `.Error` sau mỗi truy vấn.
```go
result := db.First(&user)
if result.Error != nil {
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        // Xử lý không tìm thấy
    }
}
```

## 6. Raw SQL Query
Trong một số trường hợp cực kỳ phức tạp, bạn có thể viết SQL thuần:
```go
db.Raw("SELECT id, name, age FROM users WHERE name = ?", "jinzhu").Scan(&result)

// Thực thi lệnh ghi (Update/Delete/Insert)
db.Exec("UPDATE users SET name = ? WHERE id = ?", "new_name", 1)
```
