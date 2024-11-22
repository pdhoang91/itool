from flask import Flask, request, jsonify
from PIL import Image
import pytesseract
import os
import cv2
import numpy as np

app = Flask(__name__)

def preprocess_image(image_path):
    # Đọc ảnh bằng OpenCV
    img = cv2.imread(image_path)
    
    # Chuyển sang ảnh xám
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    
    # Khử nhiễu bằng Gaussian Blur
    denoised = cv2.GaussianBlur(gray, (3, 3), 0)
    
    # Tăng cường độ tương phản bằng CLAHE
    clahe = cv2.createCLAHE(clipLimit=2.0, tileGridSize=(8,8))
    enhanced = clahe.apply(denoised)
    
    # Thresholding để tách biệt text với background
    _, binary = cv2.threshold(enhanced, 0, 255, cv2.THRESH_BINARY + cv2.THRESH_OTSU)
    
    return binary

@app.route('/ocr', methods=['POST'])
def ocr():
    if 'image' not in request.files:
        return jsonify({"error": "No image file provided"}), 400

    image = request.files['image']
    image_path = f"/tmp/{image.filename}"
    image.save(image_path)

    try:
        # Tiền xử lý ảnh
        processed_img = preprocess_image(image_path)
        
        # Cấu hình OCR
        custom_config = r'--oem 3 --psm 6 -l vie+eng'
        
        # Thực hiện OCR với ảnh đã xử lý
        text = pytesseract.image_to_string(
            Image.fromarray(processed_img),
            config=custom_config
        )
        
        # Xử lý text sau OCR
        text = text.strip()  # Loại bỏ khoảng trắng thừa
        
        response = {
            "text": text,
            "status": "success"
        }
        
    except Exception as e:
        response = {
            "error": str(e),
            "status": "error"
        }
    finally:
        # Xóa file tạm
        if os.path.exists(image_path):
            os.remove(image_path)

    return jsonify(response), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5006)