# Fixed Window Rate Limiting Strategy

## 1. Overview

Fixed Window là một chiến lược rate limiting đơn giản, trong đó thời gian được chia thành các khoảng (window) cố định. Trong mỗi window, hệ thống chỉ cho phép tối đa một số lượng request nhất định.

Khi bước sang window mới, bộ đếm (counter) sẽ được reset về 0.

---

## 2. How It Works

- Xác định:
  - `limit`: số request tối đa
  - `window_size`: kích thước cửa sổ thời gian (ví dụ: 1 phút, 10 giây)

- Với mỗi request:
  1. Xác định window hiện tại dựa trên timestamp
  2. Tăng counter của window đó
  3. Nếu counter vượt quá limit → reject request
  4. Nếu chưa vượt → allow request

---

## 3. Example

### Configuration

- Limit: 5 requests
- Window size: 10 seconds

### Timeline

#### Window [00:00 - 00:10]

| Time  | Request | Counter | Result |
|-------|--------|--------|--------|
| 00:01 | #1     | 1      | Allow  |
| 00:02 | #2     | 2      | Allow  |
| 00:03 | #3     | 3      | Allow  |
| 00:04 | #4     | 4      | Allow  |
| 00:05 | #5     | 5      | Allow  |
| 00:06 | #6     | 6      | Reject |

#### Window [00:10 - 00:20]

| Time  | Request | Counter | Result |
|-------|--------|--------|--------|
| 00:10 | #7     | 1      | Allow  |
| 00:11 | #8     | 2      | Allow  |

---

# 5. Ưu điểm
- Đơn giản để triển khai
- Overhead thấp (Chỉ cần 1 bộ đếm)
- Dễ dàng mở rộng (đặc biệt với Redis)

# 6. Nhược điểm

## Vấn đề về Burst

Fixed Window cho phép xảy ra hiện tượng "burst" tại ranh giới giữa 2 window.

Ví dụ:
5 requests tại 00:09 (cuối window 1)
5 requests tại 00:10 (đầu window 2)

→ Tổng cộng: 10 requests trong ~1 giây

Điều này vượt quá giới hạn thực tế nhưng vẫn hợp lệ theo Fixed Window.

## 7. Độ phức tạp

- Time complexity: O(1) cho mỗi request, vì chỉ làm lấy counter, tăng counter rồi so sánh với limit.
- Space complexity: O(1) cho mỗi IP