# itool
Sau khi đã cấu hình tất cả các service, bạn có thể khởi chạy toàn bộ hệ thống bằng Docker Compose.

Khởi Chạy Hệ Thống
Chạy lệnh sau trong thư mục gốc của dự án:


docker-compose up --build
Kiểm Tra Kết Nối Các Service
Frontend: Truy cập http://localhost:3000 để xem giao diện người dùng.
Management API: Truy cập http://localhost:81 để kiểm tra API.
Cơ sở Dữ liệu: Sử dụng công cụ quản lý PostgreSQL như pgAdmin hoặc DBeaver để truy cập vào cơ sở dữ liệu tại localhost:5432 với thông tin đăng nhập đã cấu hình.
Kiểm Tra Các Service: Bạn có thể sử dụng Postman hoặc curl để gửi yêu cầu tới các endpoint của từng service.
Ví dụ Sử Dụng curl
Text-to-Voice:


curl -X POST http://localhost:5001/convert -H "Content-Type: application/json" -d '{"text": "Hello World"}'
Voice-to-Text:


curl -X POST http://localhost:5002/convert -F "audio=@/path/to/audio/file.mp3"
Background Removal:


curl -X POST http://localhost:5003/remove-bg -F "image=@/path/to/image/file.png"
Speech Recognition:


curl -X POST http://localhost:5004/recognize -H "Content-Type: application/json" -d '{"audio_url": "http://example.com/audio.mp3"}'
Face Recognition:

curl -X POST http://localhost:5005/recognize-face -F "image=@/path/to/image/file.png"
OCR:

curl -X POST http://localhost:5006/ocr -F "image=@/path/to/image/file.png"
Translation:

curl -X POST http://localhost:5007/translate -H "Content-Type: application/json" -d '{"text": "Hello", "dest_lang": "vi"}'
Management API (Text-to-Voice):

curl -X POST http://localhost:81/tts -H "Content-Type: application/json" -d '{"text": "Hello World"}'
Kết Luận
Bạn đã có một hệ thống microservices hoàn chỉnh với các service chính như Text-to-Voice, Voice-to-Text, Background Removal, Speech Recognition, Face Recognition, OCR và Translation. Hệ thống được điều phối thông qua Management API và giao diện người dùng được xây dựng bằng Next.js. Mỗi service được triển khai riêng biệt, dễ dàng mở rộng và bảo trì.