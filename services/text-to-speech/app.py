from flask import Flask, request, send_file, jsonify
import numpy as np
import torch
from TTS.utils.synthesizer import Synthesizer
from TTS.utils.manage import ModelManager
from TTS.api import TTS
import io
import soundfile as sf

app = Flask(__name__)

# Khởi tạo model và synthesizer
model_manager = ModelManager()
global model_name, synthesizer

# Sử dụng model tiếng Anh mặc định
model_name = "tts_models/en/ljspeech/tacotron2-DDC"

try:
    model_path, config_path, _ = model_manager.download_model(model_name)
    synthesizer = Synthesizer(
        model_path,
        config_path,
        use_cuda=torch.cuda.is_available()
    )
except Exception as e:
    print(f"Error loading default model: {e}")
    # Fallback to another model if available
    available_models = TTS().list_models()
    if available_models:
        model_name = available_models[0]
        model_path, config_path, _ = model_manager.download_model(model_name)
        synthesizer = Synthesizer(
            model_path,
            config_path,
            use_cuda=torch.cuda.is_available()
        )

@app.route('/models', methods=['GET'])
def get_models():
    try:
        # Lấy danh sách tất cả các models có sẵn
        available_models = TTS().list_models()
        
        # Tổ chức models theo ngôn ngữ
        models_by_language = {}
        for model in available_models:
            parts = model.split('/')
            if len(parts) >= 3:
                lang = parts[1]
                if lang not in models_by_language:
                    models_by_language[lang] = []
                
                model_info = {
                    "model_id": model,
                    "type": parts[0],
                    "language": lang,
                    "dataset": parts[2],
                    "architecture": parts[3] if len(parts) > 3 else "unknown"
                }
                models_by_language[lang].append(model_info)
        
        return jsonify({
            "total_models": len(available_models),
            "models_by_language": models_by_language,
            "current_model": model_name
        })
        
    except Exception as e:
        return {'error': str(e)}, 500

@app.route('/models/languages', methods=['GET'])
def get_languages():
    try:
        available_models = TTS().list_models()
        languages = set()
        for model in available_models:
            parts = model.split('/')
            if len(parts) >= 2:
                languages.add(parts[1])
        
        return jsonify({
            "total_languages": len(languages),
            "languages": sorted(list(languages)),
            "current_model": model_name
        })
        
    except Exception as e:
        return {'error': str(e)}, 500

@app.route('/tts', methods=['POST'])
def text_to_speech():
    global model_name, synthesizer
    try:
        data = request.get_json()
        text = data.get('text', '')
        speed = float(data.get('speed', 1.0))
        pitch = float(data.get('pitch', 1.0))
        selected_model = data.get('model', model_name)
        
        # Kiểm tra và tải model mới nếu cần
        if selected_model != model_name:
            try:
                new_model_path, new_config_path, _ = model_manager.download_model(selected_model)
                synthesizer = Synthesizer(
                    new_model_path,
                    new_config_path,
                    use_cuda=torch.cuda.is_available()
                )
                model_name = selected_model
            except Exception as e:
                return {'error': f'Model not found: {str(e)}'}, 404
        
        # Tổng hợp giọng nói
        wav = synthesizer.tts(
            text,
            speaker_wav=None,
            speaker_name=None,
            speed=speed,
            pitch=pitch
        )
        
        # Chuyển đổi numpy array thành bytes
        wav_bytes = io.BytesIO()
        sf.write(wav_bytes, wav, synthesizer.output_sample_rate, format='WAV')
        wav_bytes.seek(0)
        
        return send_file(
            wav_bytes,
            mimetype='audio/wav',
            as_attachment=True,
            download_name='output.wav'
        )
        
    except Exception as e:
        return {'error': str(e)}, 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5001)