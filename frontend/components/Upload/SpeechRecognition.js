// frontend/components/Upload/SpeechRecognition.js

import { useState } from 'react';
import { recognizeSpeech, uploadAudio } from '../../services/api';

export default function SpeechRecognition({ setResult }) {
    const [audioFile, setAudioFile] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    const handleSpeechRecognition = async () => {
        if (!audioFile) {
            alert('Vui lòng tải lên tệp âm thanh trước.');
            return;
        }

        setLoading(true);
        setError(null);
        try {
            // Tải tệp âm thanh lên để lấy URL
            const uploadResponse = await uploadAudio(audioFile);
            const audioUrl = uploadResponse.data.audio_url;

            // Gọi API nhận diện giọng nói
            const response = await recognizeSpeech(audioUrl);
            setResult(response.data);
        } catch (error) {
            console.error('Error in Speech Recognition:', error);
            setError('Có lỗi xảy ra khi nhận diện giọng nói.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div>
            <h2>Speech Recognition</h2>
            <input
                type="file"
                accept="audio/*"
                onChange={(e) => setAudioFile(e.target.files[0])}
            />
            <button onClick={handleSpeechRecognition} disabled={loading}>
                {loading ? 'Đang xử lý...' : 'Recognize Speech'}
            </button>
            {error && <p className="error">{error}</p>}
        </div>
    );
}
