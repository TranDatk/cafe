# Sliding Window Log Rate Limiting

## 1. Tổng quan

Sliding Window Log là một chiến lược rate limiting dùng để giới hạn số lượng request trong một khoảng thời gian nhất định.

Thay vì chia thời gian thành các cửa sổ cố định như Fixed Window, phương pháp này sử dụng một "cửa sổ trượt" theo thời gian thực, giúp kiểm soát chính xác hơn số lượng request.

---

## 2. Ý tưởng cốt lõi

- Lưu lại timestamp của từng request
- Khi có request mới:
  - Loại bỏ các request đã quá thời gian window
  - Đếm số request còn lại
  - So sánh với giới hạn (limit)

Nếu vượt quá limit → từ chối request

---

## 3. Cách hoạt động

Giả sử:
- Limit: 3 request
- Window: 10 giây

### Quy trình:

1. Nhận request tại thời điểm hiện tại (now)
2. Lấy danh sách các timestamp trước đó
3. Xóa các timestamp nhỏ hơn (now - window)
4. Kiểm tra:
   - Nếu số lượng còn lại >= limit → reject
   - Ngược lại → accept và thêm timestamp mới

---

## 4. Ví dụ minh họa

Giả sử các request xảy ra tại:

[1s, 3s, 8s]

Tại thời điểm 9s:
- Các request trong 10s gần nhất: [1s, 3s, 8s]
- Tổng = 3 → đạt limit → request mới bị từ chối

Tại thời điểm 12s:
- Window: từ 2s đến 12s
- Các request hợp lệ: [3s, 8s]
- Có thể accept request mới → [3s, 8s, 12s]

---

## 5. Ưu điểm

- Độ chính xác cao (tính theo từng request)
- Không bị hiện tượng burst ở đầu window
- Công bằng hơn giữa các client

---

## 6. Nhược điểm

- Tốn bộ nhớ do phải lưu toàn bộ timestamp
- Tốn CPU do cần lọc dữ liệu mỗi lần request
- Khó scale trong hệ thống phân tán có lưu lượng lớn

---

## 7. Độ phức tạp

- Time complexity: O(n) cho mỗi request (lọc log)
- Space complexity: O(n) với n là số request trong window

