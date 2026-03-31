# Tài liệu GORM Tags

GORM Tags cho phép bạn cấu hình cách mapping giữa Golang Struct và Database Table trực tiếp trong code.

## 1. Cú pháp cơ bản
```go
type User struct {
    Name string `gorm:"primaryKey;column:user_name;type:varchar(100);not null"`
}
```
*   Dùng dấu chấm phẩy (`;`) để ngăn cách các tag.
*   Cấu trúc: `gorm:"tag_name:value"`.

## 2. Các tag phổ biến

### Primary Key
- `primaryKey`: Đánh dấu field này là Khóa chính.
    - VD: `ID uint `gorm:"primaryKey"``
- `autoIncrement`: Tự động tăng giá trị (mặc định nếu type là INT).
    - VD: `Number int `gorm:"autoIncrement"``

### Column Definition
- `column:NAME`: Đổi tên cột trong database.
    - VD: `FullName string `gorm:"column:user_full_name"``
- `type:TYPE`: Chỉ định kiểu dữ liệu SQL.
    - VD: `Data string `gorm:"type:text"`` hoặc `Metadata jsonb `gorm:"type:jsonb"``
- `size:SIZE`: Độ dài của cột (thường dùng cho string/varchar).
    - VD: `Password string `gorm:"size:255"``
- `precision:P;scale:S`: Độ chính xác cho số thập phân.
    - VD: `Price float64 `gorm:"precision:10;scale:2"`` (Tổng 10 chữ số, 2 chữ số sau dấu phẩy)

### Constraints & Validation
- `not null`: Cột không được để trống.
    - VD: `Email string `gorm:"not null"``
- `unique`: Giá trị duy nhất, không trùng lặp.
    - VD: `Username string `gorm:"unique"``
- `default:VALUE`: Giá trị mặc định khi INSERT nếu field bị rỗng.
    - VD: `Status string `gorm:"default:'active'"`` (Lưu ý dấu nháy đơn cho chuỗi)
- `index`: Tạo index đơn cho cột này để tăng tốc truy vấn.
    - VD: `Age int `gorm:"index"``
- `uniqueIndex`: Tạo unique index cho cột này.
    - VD: `Slug string `gorm:"uniqueIndex"``

### Trình quản lý (Permissions)
- `-`: Bỏ qua field này hoàn toàn, không map vào DB.
    - VD: `TempData string `gorm:"-"``
- `->:false`: Chế độ chỉ đọc (Read only), không lưu vào DB khi Create/Update.
    - VD: `DerivedField string `gorm:"->:false"``
- `<-:false`: Chỉ ghi (Write only), không lấy dữ liệu lên khi Select.
    - VD: `Password string `gorm:"column:password;<-:false"``

## 3. Association Tags (Quan hệ)
- `foreignKey`: Khóa ngoại trỏ đến Primary Key của bảng khác.
    - VD: `Orders []Order `gorm:"foreignKey:UserID"``
- `references`: Chỉ định field nào ở bảng cha được trỏ đến (nếu không dùng ID).
    - VD: `Profile Profile `gorm:"references:Email"``
- `many2many:TABLE_NAME`: Tên bảng trung gian cho quan hệ n-n.
    - VD: `Roles []Role `gorm:"many2many:user_roles"``
- `joinForeignKey`: Tên khóa ngoại của bảng gốc trong bảng trung gian.
    - VD: `many2many:user_roles;joinForeignKey:user_identifier`
- `joinReferences`: Tên khóa ngoại của bảng đích trong bảng trung gian.
    - VD: `many2many:user_roles;joinReferences:role_identifier`

## 4. Time Tracking
GORM có 3 field đặc biệt tự động cập nhật nếu bạn đặt đúng tên hoặc dùng tag:
- `CreatedAt`: Tự động điền thời gian lúc tạo.
    - VD: `Created time.Time `gorm:"autoCreateTime"``
- `UpdatedAt`: Tự động cập nhật thời gian mỗi khi save/update.
    - VD: `Updated time.Time `gorm:"autoUpdateTime"``
- `DeletedAt`: Dùng cho **Soft Delete** (Xóa mềm).
    - VD: `Deleted gorm.DeletedAt `gorm:"index"``

### Ví dụ tổng quát:
```go
type User struct {
    ID        uint           `gorm:"primaryKey"`
    Email     string         `gorm:"uniqueIndex;not null;size:100"`
    Password  string         `gorm:"column:password_hash;not null;<-:false"`
    Age       int            `gorm:"default:18;index"`
    Role      string         `gorm:"type:varchar(20);default:'member'"`
    CreatedAt time.Time      `gorm:"column:created_at"`
    UpdatedAt time.Time      `gorm:"column:updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```
