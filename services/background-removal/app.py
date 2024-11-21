import sys
from flask import Flask, request, jsonify, send_file
from rembg import remove
import os
import logging
from logging.handlers import RotatingFileHandler

app = Flask(__name__)

# Cấu hình logging
handler = RotatingFileHandler('/app/logs/app.log', maxBytes=1000000, backupCount=3)
formatter = logging.Formatter(
    '[%(asctime)s] %(levelname)s in %(module)s: %(message)s'
)
handler.setFormatter(formatter)
handler.setLevel(logging.DEBUG)  # Thay đổi mức log thành DEBUG để nhận nhiều thông tin hơn
app.logger.addHandler(handler)
app.logger.setLevel(logging.DEBUG)

@app.route('/remove-bg', methods=['POST'])
def remove_bg():
    app.logger.info('Received request to /remove-bg')

    if 'image' not in request.files:
        app.logger.warning('No image file provided in the request')
        return jsonify({"error": "No image file provided"}), 400

    image = request.files['image']
    image_path = f"/tmp/{image.filename}"
    output_image_path = f"/tmp/output_{image.filename}"

    app.logger.info(f"Image received: {image.filename}, saving to {image_path}")

    try:
        # Lưu ảnh gốc tạm thời
        image.save(image_path)
        app.logger.info(f"Image saved at {image_path}")

        # Xóa nền
        app.logger.info("Removing background...")
        with open(image_path, 'rb') as input_file:
            input_data = input_file.read()
            output_data = remove(input_data)
        app.logger.info("Background removed successfully")

        # Lưu ảnh đã xử lý
        with open(output_image_path, 'wb') as output_file:
            output_file.write(output_data)
        app.logger.info(f"Processed image saved at {output_image_path}")

        # Trả về đường dẫn file
        return jsonify({"processed_image_path": output_image_path}), 200
    except Exception as e:
        app.logger.error("Error during background removal", exc_info=True)
        return jsonify({"error": str(e)}), 500
    finally:
        # Xóa file gốc (nếu không cần nữa)
        if os.path.exists(image_path):
            os.remove(image_path)

if __name__ == '__main__':
    app.logger.info('Starting Flask server')
    app.run(host='0.0.0.0', port=5003)
