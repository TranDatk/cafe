# Phân trang bằng Offset (Offset Pagination)

Đây là phương pháp phân trang phổ biến nhất và dễ triển khai nhất, thường thấy trong các ứng dụng web truyền thống có thanh chuyển trang (Pagination bar).

## 1. Cơ chế hoạt động
Phương pháp này dựa trên việc bỏ qua một số lượng bản ghi nhất định (`OFFSET`) và lấy ra một số lượng bản ghi tiếp theo (`LIMIT`).

*   **Tham số đầu vào:** `page` (số trang) và `limit` (số bản ghi mỗi trang).
*   **Công thức:** `OFFSET = (page - 1) * limit`.

## 2. Ưu điểm
- **Dễ triển khai:** Hầu hết các cơ sở dữ liệu quan hệ (SQL) đều hỗ trợ từ khóa `LIMIT` và `OFFSET`.
- **Nhảy trang trực tiếp:** Người dùng có thể nhảy thẳng từ trang 1 sang trang 100 dễ dàng.
- **Biết tổng số trang:** Thường đi kèm với câu lệnh `COUNT` để hiển thị tổng số trang cho người dùng.

## 3. Nhược điểm (Vấn đề lớn)
- **Hiệu năng giảm dần:** Khi ứng dụng có hàng triệu bản ghi, việc sử dụng `OFFSET 1,000,000` yêu cầu Database phải quét qua 1 triệu dòng trước đó và vứt bỏ chúng, gây tốn tài nguyên và tăng thời gian phản hồi.
- **Dữ liệu không nhất quán (Inconsistent Data):**
    - Nếu có một bản ghi mới được thêm vào trang 1 trong lúc người dùng đang ở trang 1, khi họ chuyển sang trang 2, bản ghi cuối cùng của trang 1 sẽ bị "đẩy" xuống trang 2, khiến người dùng nhìn thấy bản ghi đó hai lần.
    - Ngược lại, nếu xóa bản ghi ở trang 1, dữ liệu ở trang 2 sẽ bị "kéo" lên, khiến người dùng bị sót bản ghi khi chuyển trang.

## 4. Trường hợp sử dụng
- Các ứng dụng có lượng dữ liệu vừa và nhỏ.
- Các hệ thống admin cần phân trang rõ ràng và khả năng nhảy trang linh hoạt.
