# Tài liệu GORM Pagination (Phân trang)

Phân trang là kỹ thuật chia nhỏ tập dữ liệu lớn thành nhiều trang để cải thiện hiệu năng và trải nghiệm người dùng. Trong GORM, chúng ta sử dụng hai hàm chính: `Limit` và `Offset`.

## 1. Công thức cơ bản
*   **Limit**: Số lượng bản ghi muốn lấy trên một trang (Page Size).
*   **Offset**: Số lượng bản ghi cần bỏ qua để đến được trang mong muốn.

**Công thức tính Offset:**
`Offset = (Page - 1) * PageSize`

**Ví dụ:** Muốn lấy trang 2, mỗi trang 10 bản ghi:
*   `Limit = 10`
*   `Offset = (2 - 1) * 10 = 10`

```go
db.Limit(10).Offset(10).Find(&users)
```

## 2. Cách viết Reusable Scopes (Khuyên dùng)
Để tránh lặp lại code ở nhiều nơi, bạn nên viết một hàm **Scope** để tái sử dụng cho mọi Model.

```go
func Paginate(page, pageSize int) func(db *gorm.DB) *gorm.DB {
  return func(db *gorm.DB) *gorm.DB {
    if page <= 0 {
      page = 1
    }

    switch {
    case pageSize > 100:
      pageSize = 100 // Giới hạn tối đa để tránh kéo quá nhiều data
    case pageSize <= 0:
      pageSize = 10
    }

    offset := (page - 1) * pageSize
    return db.Offset(offset).Limit(pageSize)
  }
}

// Cách sử dụng:
db.Scopes(Paginate(page, pageSize)).Find(&users)
```

## 3. Lấy Tổng số bản ghi (Total Count)
Khi phân trang, frontend thường cần biết tổng số bản ghi (`total`) để hiển thị số trang. 

**Lưu ý:** Bạn phải gọi `Count` trước khi gọi `Limit` và `Offset` trên cùng một object query, hoặc dùng các biến chứa query riêng biệt.

```go
var total int64
var users []User

// 1. Tạo query cơ bản với các điều kiện lọc (nếu có)
query := db.Model(&User{}).Where("active = ?", true)

// 2. Đếm tổng số bản ghi khớp điều kiện
query.Count(&total)

// 3. Thực hiện phân trang và lấy dữ liệu
query.Scopes(Paginate(page, pageSize)).Find(&users)
```

## 4. Cấu trúc Response gợi ý cho API
Thông thường, một API phân trang nên trả về metadata đầy đủ:

```json
{
  "data": [...],
  "pagination": {
    "current_page": 1,
    "page_size": 10,
    "total_records": 50,
    "total_pages": 5
  }
}
```

**Cách tính Total Pages trong Go:**
`totalPages := int(math.Ceil(float64(total) / float64(pageSize)))`

## 5. Lưu ý quan trọng về Hiệu năng
*   **Deep Pagination:** Khi `Offset` quá lớn (ví dụ trang 1000, offset 10,000), database vẫn phải quét qua 10,000 bản ghi trước đó rồi mới lấy 10 bản ghi tiếp theo. Điều này làm query rất chậm.
*   **Giải pháp:** Với dữ liệu cực lớn, hãy cân nhắc sử dụng **Cursor Pagination** (Dựa trên ID của bản ghi cuối cùng của trang trước thay vì sử dụng Offset).
