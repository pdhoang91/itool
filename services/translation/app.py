from flask import Flask, request, jsonify
from googletrans import Translator
import os
import psycopg2
import json
import time

app = Flask(__name__)
translator = Translator()

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

@app.route('/translate', methods=['POST'])
def translate():
    data = request.get_json()
    if not data or 'text' not in data or 'dest_lang' not in data:
        return jsonify({"error": "Invalid input"}), 400

    text = data['text']
    dest_lang = data['dest_lang']

    # Insert task vào cơ sở dữ liệu
    conn = get_db_connection()
    cur = conn.cursor()
    cur.execute(
        "INSERT INTO tasks (service_name, status, input_data) VALUES (%s, %s, %s) RETURNING id",
        ("translation", "processing", json.dumps({"text": text, "dest_lang": dest_lang}))
    )
    task_id = cur.fetchone()[0]
    conn.commit()

    try:
        # Thực hiện dịch văn bản
        translated = translator.translate(text, dest=dest_lang).text

        # Cập nhật task
        cur.execute(
            "UPDATE tasks SET status=%s, output_data=%s, updated_at=NOW() WHERE id=%s",
            ("completed", json.dumps({"translated_text": translated}), task_id)
        )
        conn.commit()

        response = {"translated_text": translated}
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

    return jsonify(response), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5007)
