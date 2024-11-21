from flask import Flask, request, send_file, jsonify
from rembg import remove
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

@app.route('/remove-bg', methods=['POST'])
def remove_bg():
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
        ("background-removal", "processing", json.dumps({"image_file": image.filename}))
    )
    task_id = cur.fetchone()[0]
    conn.commit()

    # Thực hiện xóa nền
    try:
        with open(image_path, 'rb') as i:
            input_data = i.read()
            output_data = remove(input_data)

        output_image_path = f"/output/output_{task_id}.png"
        os.makedirs(os.path.dirname(output_image_path), exist_ok=True)
        with open(output_image_path, 'wb') as o:
            o.write(output_data)

        # Cập nhật task
        output_url = f"http://localhost:5003/output/output_{task_id}.png"
        cur.execute(
            "UPDATE tasks SET status=%s, output_data=%s, updated_at=NOW() WHERE id=%s",
            ("completed", json.dumps({"image_url": output_url}), task_id)
        )
        conn.commit()
    except Exception as e:
        # Cập nhật task với trạng thái lỗi
        cur.execute(
            "UPDATE tasks SET status=%s, output_data=%s, updated_at=NOW() WHERE id=%s",
            ("failed", json.dumps({"error": str(e)}), task_id)
        )
        conn.commit()

    cur.close()
    conn.close()

    # Xóa file tạm
    os.remove(image_path)

    return jsonify({"image_url": output_url}), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5003)
