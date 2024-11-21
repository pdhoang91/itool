// frontend/components/Upload/TextToVoice.js

import { useState } from 'react';
import { convertTextToVoice } from '../../services/api';

export default function TextToVoice({ setResult }) {
    const [text, setText] = useState('');
    const [language, setLanguage] = useState('en');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    const handleTTS = async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await convertTextToVoice({ text, language });
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
            />
            <div>
                <label htmlFor="language">Language:</label>
                <select
                    id="language"
                    value={language}
                    onChange={(e) => setLanguage(e.target.value)}
                >
                    <option value="en">English</option>
                    <option value="vi">Vietnamese</option>
                    <option value="fr">French</option>
                    {/* Thêm các ngôn ngữ khác tại đây */}
                </select>
            </div>
            <button onClick={handleTTS} disabled={loading}>
                {loading ? 'Đang chuyển đổi...' : 'Convert to Voice'}
            </button>
            {error && <p className="error">{error}</p>}
        </div>
    );
}
