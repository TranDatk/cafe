# Continuous Integration (CI)
> [!NOTE]
> Tài liệu này mô tả quy trình CI tự động của hệ thống Cafe thông qua GitHub Actions. Nó giúp đảm bảo mã nguồn mới gửi lên không làm phá vỡ ứng dụng thông qua việc tự động tải package, build và chạy unit tests.

## Cấu Hình CI

Workflow được lưu trữ tại file: `.github/workflows/ci.yml`

### Triggers (Điều kiện kích hoạt)
Quy trình CI sẽ tự động chạy trong hai trường hợp:
1. **Push:** Mỗi khi có thay đổi được đẩy (push) lên nhánh `main`.
2. **Pull request:** Mỗi khi có một pull request mới chĩa vào nhánh `main`.

### Các Bước Thực Bước CI (Jobs: build-and-test)

Công việc CI này chạy trên môi trường giả lập mới nhất của Ubuntu (`ubuntu-latest`) và bao gồm các step chính sau:

1. **Checkout Code:** Lấy toàn bộ mã nguồn của kho lưu trữ về môi trường máy chủ chạy Github Actions bằng công cụ `actions/checkout@v4`.
2. **Setup Go Context:** Thiết lập môi trường Go lang để build. Thay vì hard-code một phiên bản cố định, nó gọi `actions/setup-go@v5` và tự động đọc phiên bản từ chính file `go.mod` (tham số `go-version-file: 'go.mod'`). Đảm bảo đồng bộ với máy dev.
3. **Quản lý Caches:** Lưu và sử dụng cache của Go module (`actions/cache@v4`). Việc này tiết kiệm được đáng kể thời gian tải module bằng lưu các metadata nếu file `go.sum` chưa hề thay đổi so với lần CI gần nhất.
4. **Cài Đặt Packages (Dependencies):** Tải về các thư viện cần thiết phục vụ cho quá trình build bằng `go mod download`.
5. **Build:** Chạy lệnh `go build -v ./...` để tiến hành biên dịch ứng dụng. Nhằm đảm bảo code không mắc bất kỳ lỗi cú pháp hay hỏng gói nào.
6. **Tests:** Chạy lệnh `go test -v -race -cover ./...` với chế độ data race detector. Việc gọi lệnh test đảm bảo mọi Unit Test hoặc Behavior Test được cài cắm trong repo đều hoạt động trơn tru.

## Tùy Biến

Nếu tương lai ứng dụng chuyển sang sử dụng Docker Image hoặc cần phải kết nối database trong CI, chúng ta nên bổ sung thêm Job (e.g. `docker build`) hoặc chèn script setup DB vào phần `Services` của GitHub Actions.
