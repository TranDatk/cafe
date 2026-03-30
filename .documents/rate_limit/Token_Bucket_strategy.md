# Token Bucket Strategy (Rate Limiting)

## 1. Overview
Token Bucket là một thuật toán dùng để giới hạn số lượng request theo thời gian. 
Mỗi request cần “lấy” một token để được xử lý.

## 2. Core Idea
- Hệ thống có một "bucket" chứa các token
- Token được nạp vào bucket theo tốc độ cố định (refill rate)
- Bucket có dung lượng tối đa (capacity)
- Mỗi request tiêu tốn 1 token

Nếu còn token → request được xử lý  
Nếu hết token → request bị từ chối (hoặc chờ)

## 3. Key Components
- Capacity (burst): số token tối đa trong bucket
- Refill rate: số token được thêm vào mỗi giây
- Tokens hiện tại: số token đang có trong bucket

## 4. How It Works
1. Khi có request đến:
   - Cập nhật số token dựa trên thời gian trôi qua (refill)
2. Kiểm tra:
   - Nếu tokens >= 1:
       → trừ 1 token
       → cho phép request
   - Ngược lại:
       → từ chối request

## 5. Example
Cấu hình:
- Capacity = 10
- Refill rate = 5 tokens/second
- token = min(capacity, tokens + refill_rate * Δt) mà Δt = now - last_refill_time

Flow:
- Ban đầu: 10 tokens
- 10 request đến cùng lúc → tất cả được xử lý (tokens = 0)
- Request thứ 11 → bị từ chối
- Sau 1 giây → +5 tokens → có thể xử lý tiếp 5 request

## 6. Advantages
- Cho phép burst (xử lý nhiều request cùng lúc)
- Kiểm soát tốc độ request hiệu quả
- Trải nghiệm người dùng tốt hơn so với Fixed Window

## 7. Comparison
- Fixed Window: dễ bị dồn request ở cuối window
- Sliding Window: chính xác hơn nhưng phức tạp
- Leaky Bucket: xử lý đều, không hỗ trợ burst
- Token Bucket: cân bằng giữa burst và kiểm soát tốc độ

## 8. Độ phức tạp

- Time complexity: O(1) cho mỗi request
- Space complexity: O(1), chỉ cần lưu tokens và last_refill_timestamp

## 9. Summary
Token Bucket cho phép request bùng nổ tại 1 thời điểm nhưng vẫn kiểm soát được tốc độ tổng thể giữa các request (Ở đây burst là thể hiện cho sức chứa tối đa, còn refill rate là thể hiện cho tốc độ nạp token). Nó tốt hơn fixed window ở chỗ là nếu có 2 t1 và t2 sát nhau với t1 < t2 thì số request trong khoảng [t1, t2] sẽ không vượt quá capacity + refill_rate * (t2 - t1).