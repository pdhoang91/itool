// frontend/components/Upload/TextToVoice.js

import { useState } from 'react';
import { convertTextToVoice } from '../../services/api';
import styles from '../../styles/UploadComponent.module.css';

export default function TextToVoice({ setResult }) {
    const [text, setText] = useState('');
    const [language, setLanguage] = useState('en');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [audioUrl, setAudioUrl] = useState(null);

    const handleTTS = async () => {
        if (!text.trim()) {
            alert('Vui lòng nhập văn bản cần chuyển đổi.');
            return;
        }

        setLoading(true);
        setError(null);
        setAudioUrl(null); // Reset audioUrl trước khi chuyển đổi mới
        try {
            const response = await convertTextToVoice({ text, language });
            // Giả định rằng response.data.audio_url chứa URL của tệp MP3
            const { audio_url } = response.data;
            setAudioUrl(audio_url);
            setResult(response.data);
        } catch (error) {
            console.error('Error in Text-to-Voice:', error);
            setError('Có lỗi xảy ra khi chuyển đổi văn bản thành giọng nói.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div>
            <h2>Text-to-Voice</h2>
            <input
                type="text"
                value={text}
                onChange={(e) => setText(e.target.value)}
                placeholder="Enter text"
                className={styles.inputText}
            />
            <div>
                <label htmlFor="language">Language:</label>
                <select
                    id="language"
                    value={language}
                    onChange={(e) => setLanguage(e.target.value)}
                    className={styles.select}
                >
                    <option value="en">English</option>
                    <option value="vi">Vietnamese</option>
                    <option value="fr">French</option>
                    {/* Thêm các ngôn ngữ khác tại đây */}
                </select>
            </div>
            <button onClick={handleTTS} disabled={loading} className={styles.button}>
                {loading ? 'Đang chuyển đổi...' : 'Convert to Voice'}
            </button>
            {error && <p className={styles.error}>{error}</p>}
            {audioUrl && (
                <div className={styles.audioSection}>
                    <h3>Result:</h3>
                    <audio controls src={audioUrl} className={styles.audioPlayer}>
                        Your browser does not support the audio element.
                    </audio>
                    <a href={audioUrl} download="voice_output.mp3" className={styles.downloadButton}>
                        Download MP3
                    </a>
                </div>
            )}
        </div>
    );
}
