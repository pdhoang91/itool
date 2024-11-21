from flask import Flask, request, jsonify, send_file
from rembg import remove
import os

app = Flask(__name__)

@app.route('/remove-bg', methods=['POST'])
def remove_bg():
    if 'image' not in request.files:
        return jsonify({"error": "No image file provided"}), 400

    image = request.files['image']
    image_path = f"/tmp/{image.filename}"
    output_image_path = f"/tmp/output_{image.filename}"

    try:
        # Lưu tạm ảnh gốc
        image.save(image_path)

        # Xóa nền
        with open(image_path, 'rb') as input_file:
            input_data = input_file.read()
            output_data = remove(input_data)

        # Lưu ảnh đã xóa nền
        with open(output_image_path, 'wb') as output_file:
            output_file.write(output_data)

        # Trả ảnh đã xóa nền
        return send_file(output_image_path, mimetype='image/png')
    except Exception as e:
        return jsonify({"error": str(e)}), 500
    finally:
        # Dọn dẹp file tạm
        if os.path.exists(image_path):
            os.remove(image_path)
        if os.path.exists(output_image_path):
            os.remove(output_image_path)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5003)
