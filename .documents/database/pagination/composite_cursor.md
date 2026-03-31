# Phân trang bằng con trỏ tổng hợp (Composite Cursor Pagination)

Đây là phiên bản nâng cao của Cursor Pagination, được sử dụng khi bạn cần sắp xếp dữ liệu theo các tiêu chí không duy nhất (như Ngày tạo, Điểm số, Tên) thay vì chỉ dùng ID.

## 1. Tại sao cần Composite Cursor?
Nếu bạn chỉ sử dụng con trỏ là `created_at` (Ngày tạo):
*   `SELECT * FROM posts WHERE created_at < '2023-01-01' ORDER BY created_at DESC LIMIT 10`
*   **Vấn đề:** Nếu có 20 bài viết cùng được tạo vào chính xác một thời điểm `2023-01-01`, con trỏ sẽ bị nhầm lẫn và gây ra hiện tượng mất dữ liệu hoặc trùng lặp vì điều kiện `<` sẽ bỏ qua các bài viết cùng giây/mili giây đó.

## 2. Cơ chế hoạt động
Chúng ta kết hợp cột cần sắp xếp (ví dụ `created_at`) với một cột duy nhất (thường là `ID`) để tạo thành một con trỏ phức hợp.

*   **Tham số đầu vào:** `last_created_at` và `last_id`.
*   **Logic (SQL):**
    ```sql
    SELECT * FROM posts
    WHERE (created_at < last_created_at)
       OR (created_at = last_created_at AND id < last_id)
    ORDER BY created_at DESC, id DESC
    LIMIT 10;
    ```

## 3. Ưu điểm
- **Sắp xếp linh hoạt:** Có thể sắp xếp theo bất kỳ tiêu chí nào mà vẫn giữ được ưu điểm của Cursor Pagination.
- **Độ chính xác tuyệt đối:** Đảm bảo không bao giờ bị sót bản ghi ngay cả khi các giá trị sắp xếp bị trùng nhau.

## 4. Nhược điểm
- **Query phức tạp hơn:** Câu lệnh SQL và logic xử lý ở code backend phức tạp hơn so với con trỏ đơn.
- **Yêu cầu Index phức hợp:** Để đạt hiệu năng tối ưu, bạn cần tạo một **Composite Index** trên Database cho cả 2 cột (ví dụ: `(created_at, id)`).

## 5. Trường hợp sử dụng
- Sắp xếp sản phẩm theo giá cả hoặc đánh giá.
- Sắp xếp danh sách người dùng theo tên hoặc điểm số.
- Bất kỳ danh sách nào cần sắp xếp theo thuộc tính không duy nhất.
