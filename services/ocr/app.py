from flask import Flask, request, jsonify
from PIL import Image
import pytesseract
import os

app = Flask(__name__)

@app.route('/ocr', methods=['POST'])
def ocr():
    if 'image' not in request.files:
        return jsonify({"error": "No image file provided"}), 400

    image = request.files['image']
    image_path = f"/tmp/{image.filename}"
    image.save(image_path)

    try:
        # Thực hiện OCR
        img = Image.open(image_path)
        text = pytesseract.image_to_string(img)

        response = {"text": text}
    except Exception as e:
        response = {"error": str(e)}

    # Xóa file tạm
    os.remove(image_path)

    return jsonify(response), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5006)
