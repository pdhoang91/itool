from flask import Flask, request, jsonify
from PIL import Image
import pytesseract
import os
import psycopg2
import json
import time

app = Flask(__name__)

# Kết nối tới cơ sở dữ liệu
def get_db_connection():
    conn = psycopg2.connect(
        host=os.getenv('DB_HOST', 'localhost'),
        port=os.getenv('DB_PORT', '5432'),
        database=os.getenv('DB_NAME', 'ai_tools'),
        user=os.getenv('DB_USER', 'admin'),
        password=os.getenv('DB_PASSWORD', 'password')
    )
    return conn

@app.route('/ocr', methods=['POST'])
def ocr():
    if 'image' not in request.files:
        return jsonify({"error": "No image file provided"}), 400

    image = request.files['image']
    image_path = f"/tmp/{image.filename}"
    image.save(image_path)

    # Insert task vào cơ sở dữ liệu
    conn = get_db_connection()
    cur = conn.cursor()
    cur.execute(
        "INSERT INTO tasks (service_name, status, input_data) VALUES (%s, %s, %s) RETURNING id",
        ("ocr", "processing", json.dumps({"image_file": image.filename}))
    )
    task_id = cur.fetchone()[0]
    conn.commit()

    try:
        # Thực hiện OCR
        img = Image.open(image_path)
        text = pytesseract.image_to_string(img)

        # Cập nhật task
        cur.execute(
            "UPDATE tasks SET status=%s, output_data=%s, updated_at=NOW() WHERE id=%s",
            ("completed", json.dumps({"text": text}), task_id)
        )
        conn.commit()

        response = {"text": text}
    except Exception as e:
        # Cập nhật task với trạng thái lỗi
        cur.execute(
            "UPDATE tasks SET status=%s, output_data=%s, updated_at=NOW() WHERE id=%s",
            ("failed", json.dumps({"error": str(e)}), task_id)
        )
        conn.commit()
        response = {"error": str(e)}

    cur.close()
    conn.close()

    # Xóa file tạm
    os.remove(image_path)

    return jsonify(response), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5006)
