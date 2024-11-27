from flask import Flask, request, jsonify
import numpy as np
import torch
from TTS.utils.synthesizer import Synthesizer
from TTS.utils.manage import ModelManager
from TTS.api import TTS
import soundfile as sf
import librosa
import os
from datetime import datetime
import logging
import threading
import time

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

app = Flask(__name__)
AUDIO_DIR = "/shared/static"
os.makedirs(AUDIO_DIR, exist_ok=True)

model_manager = ModelManager()
model_name = "tts_models/en/ljspeech/tacotron2-DDC"
all_models = {}
model_lock = threading.Lock()
default_model_ready = threading.Event()

def init_default_model():
    global synthesizer
    try:
        logger.info(f"Loading default model: {model_name}")
        model_path, config_path, model_item = model_manager.download_model(model_name)
        vocoder_name = model_item.get("default_vocoder")
        all_models[model_name] = {
            "model_path": model_path,
            "config_path": config_path,
            "vocoder_name": vocoder_name
        }
        synthesizer = Synthesizer(
            tts_checkpoint=model_path,
            tts_config_path=config_path,
            use_cuda=torch.cuda.is_available()
        )
        logger.info("Default model loaded successfully")
        default_model_ready.set()
    except Exception as e:
        logger.error(f"Error loading default model: {e}")
        raise

def load_model_in_background(model):
    try:
        logger.info(f"Loading model in background: {model}")
        model_path, config_path, model_item = model_manager.download_model(model)
        vocoder_name = model_item.get("default_vocoder")
        
        with model_lock:
            all_models[model] = {
                "model_path": model_path,
                "config_path": config_path,
                "vocoder_name": vocoder_name
            }
        logger.info(f"Successfully loaded model: {model}")
    except Exception as e:
        logger.error(f"Error loading model {model}: {e}")

def init_background_models():
    try:
        available_models = TTS().list_models()
        loaded_models = []
        for model in available_models:
            if model != model_name:  # Skip default model
                thread = threading.Thread(target=load_model_in_background, args=(model,))
                thread.daemon = True
                thread.start()
                loaded_models.append(model)
        logger.info(f"Started loading {len(loaded_models)} models in background")
    except Exception as e:
        logger.error(f"Error initializing background models: {e}")

def adjust_speed_pitch(audio, speed=1.0, pitch=1.0):
    if speed != 1.0:
        audio = librosa.effects.time_stretch(audio, rate=1/speed)
    
    if pitch != 1.0:
        n_steps = 12 * np.log2(pitch)
        audio = librosa.effects.pitch_shift(
            y=audio,
            sr=synthesizer.output_sample_rate,
            n_steps=n_steps
        )
    
    return audio

@app.route('/models', methods=['GET'])
def get_models():
    try:
        available_models = TTS().list_models()
        
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
    if not default_model_ready.is_set():
        return {'error': 'Default model not ready'}, 503
        
    global model_name, synthesizer
    try:
        data = request.get_json()
        text = data.get('text', '')
        speed = float(data.get('speed', 1.0))
        pitch = float(data.get('pitch', 1.0))
        selected_model = data.get('model', model_name)
        
        logger.info(f"Processing TTS request - Model: {selected_model}, Text: {text}")
        
        if selected_model != model_name:
            if selected_model not in all_models:
                return {'error': f'Model {selected_model} not yet loaded'}, 404
                
            model_data = all_models[selected_model]
            synthesizer = Synthesizer(
                tts_checkpoint=model_data["model_path"],
                tts_config_path=model_data["config_path"],
                use_cuda=torch.cuda.is_available()
            )
            model_name = selected_model
        
        wav = synthesizer.tts(text)
        wav = adjust_speed_pitch(wav, speed, pitch)
        
        timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
        filename = f"tts_{timestamp}.wav"
        filepath = os.path.join(AUDIO_DIR, filename)
        
        sf.write(filepath, wav, synthesizer.output_sample_rate)
        
        return jsonify({
            'success': True,
            'audio_url': filepath,
            'filename': filename
        })
        
    except Exception as e:
        logger.error(f"TTS error: {e}")
        return {'error': str(e)}, 500

if __name__ == '__main__':
    init_default_model()
    threading.Thread(target=init_background_models).start()
    app.run(host='0.0.0.0', port=5001)