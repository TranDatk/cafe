# Phân trang bằng Con trỏ (Cursor Pagination)

Cursor Pagination (còn gọi là Keyset Pagination) là phương pháp phân trang hiện đại, thường được các hệ thống lớn như Facebook, Slack hay Shopify sử dụng cho tính năng "Tải thêm" (Load more) hoặc "Cuộn vô tận" (Infinite Scroll).

## 1. Cơ chế hoạt động
Thay vì dựa vào số trang, nó dựa trên một "con trỏ" (cursor) duy nhất trỏ đến bản ghi cuối cùng của trang trước đó để lấy trang tiếp theo.

*   **Tham số đầu vào:** `cursor` (ID của bản ghi cuối trang trước) và `limit`.
*   **Logic:** `SELECT * FROM table WHERE id > last_seen_id ORDER BY id ASC LIMIT limit`.

## 2. Ưu điểm
- **Hiệu năng cực cao:** Database sử dụng Index trực tiếp trên cột con trỏ để lấy dữ liệu. Thời gian phản hồi trang 1 và trang 1000 là như nhau.
- **Dữ liệu nhất quán:** Không bị hiện tượng trùng lặp hoặc sót bản ghi khi có dữ liệu mới được thêm/xóa ở các trang trước, vì con trỏ luôn bắt đầu từ một vị trí xác định.
- **Tiết kiệm tài nguyên:** Không cần quét và bỏ qua các bản ghi cũ như Offset.

## 3. Nhược điểm
- **Không thể nhảy trang:** Người dùng buộc phải đi tuần tự từ trang 1 sang trang 2, không thể nhảy thẳng đến trang 50.
- **Không biết tổng số trang:** Khó khăn trong việc hiển thị tổng số bản ghi và tổng số trang.
- **Yêu cầu cột sắp xếp duy nhất:** Con trỏ phải dựa trên một cột có giá trị duy nhất (Unique) và có tính thứ tự (thường là ID tăng dần).

## 4. Trường hợp sử dụng
- Các ứng dụng mạng xã hội (Newfeed, Comment).
- Các hệ thống xử lý dữ liệu lớn (Big Data, Logs).
- API cho Mobile app cần tính năng Infinite Scroll.
