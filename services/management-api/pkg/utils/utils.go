package utils

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

//// SaveUploadedFile lưu file từ form và trả về đường dẫn của file đã lưu
//func SaveUploadedFile(file multipart.File, header *multipart.FileHeader, uploadPath string) (string, error) {
//	// Tạo thư mục nếu chưa tồn tại
//	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
//		return "", err
//	}
//
//	// Đảm bảo tên file không chứa đường dẫn nguy hiểm
//	safeFilename := filepath.Base(header.Filename)
//	filePath := filepath.Join(uploadPath, safeFilename)
//
//	out, err := os.Create(filePath)
//	if err != nil {
//		return "", err
//	}
//	defer out.Close()
//
//	_, err = io.Copy(out, file)
//	if err != nil {
//		return "", err
//	}
//
//	return filePath, nil
//}
//
//// GenerateURL tạo URL tạm thời cho file đã tải lên
//func GenerateURL(filename, folder string) string {
//	return fmt.Sprintf("http://localhost:81/uploads/%s/%s", folder, filename)
//}

// SaveUploadedFile lưu file tải lên từ form vào đường dẫn cụ thể
func SaveUploadedFile(file io.Reader, header *multipart.FileHeader, uploadPath string) (string, error) {
	// Đảm bảo thư mục tồn tại
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		return "", err
	}

	// Đường dẫn file đầy đủ
	filePath := filepath.Join(uploadPath, header.Filename)

	// Tạo file và ghi dữ liệu
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		return "", err
	}

	return filePath, nil
}
