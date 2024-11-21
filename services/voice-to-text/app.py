from flask import Flask, request, jsonify
from deepspeech import Model
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

# Load mô hình DeepSpeech
MODEL_PATH = 'deepspeech-0.9.3-models.pbmm'
SCORER_PATH = 'deepspeech-0.9.3-models.scorer'
ds = Model(MODEL_PATH)
ds.enableExternalScorer(SCORER_PATH)

@app.route('/convert', methods=['POST'])
def convert():
    if 'audio' not in request.files:
        return jsonify({"error": "No audio file provided"}), 400

    audio = request.files['audio']
    audio_path = f"/tmp/{audio.filename}"
    audio.save(audio_path)

    # Insert task vào cơ sở dữ liệu
    conn = get_db_connection()
    cur = conn.cursor()
    cur.execute(
        "INSERT INTO tasks (service_name, status, input_data) VALUES (%s, %s, %s) RETURNING id",
        ("voice-to-text", "processing", json.dumps({"audio_file": audio.filename}))
    )
    task_id = cur.fetchone()[0]
    conn.commit()

    # TODO: Thực hiện chuyển đổi Voice-to-Text
    # Sử dụng DeepSpeech để chuyển đổi file audio thành text
    # Ở đây, chúng ta giả lập quá trình chuyển đổi
    time.sleep(2)
    converted_text = "Converted text from audio"

    # Cập nhật task
    audio_url = f"http://localhost:5002/audio/{audio.filename}"
    cur.execute(
        "UPDATE tasks SET status=%s, output_data=%s, updated_at=NOW() WHERE id=%s",
        ("completed", json.dumps({"text": converted_text}), task_id)
    )
    conn.commit()
    cur.close()
    conn.close()

    # Xóa file tạm
    os.remove(audio_path)

    return jsonify({"text": converted_text}), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5002)
