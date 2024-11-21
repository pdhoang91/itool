// frontend/components/Upload/VoiceToText.js

import { useState } from 'react';
import { uploadAudio, convertVoiceToText } from '../../services/api';

export default function VoiceToText() {
    const [audioFile, setAudioFile] = useState(null);
    const [result, setResult] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    const handleVTS = async () => {
        if (!audioFile) {
            alert('Please upload an audio file first.');
            return;
        }

        setLoading(true);
        setError(null);
        try {
            const uploadResponse = await uploadAudio(audioFile);
            const audioUrl = uploadResponse.data.audio_url;
            const response = await convertVoiceToText(audioUrl);
            setResult(response.data);
        } catch (error) {
            console.error('Error in Voice-to-Text:', error);
            setError('Có lỗi xảy ra khi chuyển đổi giọng nói thành văn bản.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div>
            <h2>Voice-to-Text</h2>
            <input
                type="file"
                accept="audio/*"
                onChange={(e) => setAudioFile(e.target.files[0])}
            />
            <button onClick={handleVTS} disabled={loading}>
                {loading ? 'Đang chuyển đổi...' : 'Convert to Text'}
            </button>
            {error && <p className="error">{error}</p>}
            {result && (
                <div>
                    <h3>Result:</h3>
                    <pre>{JSON.stringify(result, null, 2)}</pre>
                </div>
            )}
        </div>
    );
}
