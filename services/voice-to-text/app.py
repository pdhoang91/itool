from flask import Flask, request, jsonify
from deepspeech import Model
import os

app = Flask(__name__)

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

    try:
        # Thực hiện chuyển đổi Voice-to-Text bằng DeepSpeech
        with open(audio_path, 'rb') as audio_file:
            audio_data = audio_file.read()
            converted_text = ds.stt(audio_data)

        response = {"text": converted_text}
    except Exception as e:
        response = {"error": str(e)}

    # Xóa file tạm
    os.remove(audio_path)

    return jsonify(response), 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5002)
