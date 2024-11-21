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

    app.logger.info(f'Image received: {image.filename}')
    app.logger.info(f'Saving image to {image_path}')

    try:
        # Lưu tạm ảnh gốc
        image.save(image_path)
        app.logger.info(f'Image saved successfully at {image_path}')

        # Xóa nền
        app.logger.info('Starting background removal process')
        with open(image_path, 'rb') as input_file:
            input_data = input_file.read()
            app.logger.debug('Input image data read successfully')
            output_data = remove(input_data)
            app.logger.info('Background removal completed successfully')

        # Lưu ảnh đã xóa nền
        app.logger.info(f'Saving output image to {output_image_path}')
        with open(output_image_path, 'wb') as output_file:
            output_file.write(output_data)
        app.logger.info(f'Output image saved successfully at {output_image_path}')

        # Trả ảnh đã xóa nền
        app.logger.info('Sending output image back to client')
        return send_file(output_image_path, mimetype='image/png')

    except Exception as e:
        app.logger.error('Error during background removal', exc_info=True)
        return jsonify({"error": str(e)}), 500

    finally:
        # Dọn dẹp file tạm
        try:
            if os.path.exists(image_path):
                os.remove(image_path)
                app.logger.info(f'Removed temporary file {image_path}')
            if os.path.exists(output_image_path):
                os.remove(output_image_path)
                app.logger.info(f'Removed temporary file {output_image_path}')
        except Exception as cleanup_error:
            app.logger.error('Error during cleanup', exc_info=True)

if __name__ == '__main__':
    app.logger.info('Starting Flask server')
    app.run(host='0.0.0.0', port=5003)
