// frontend/components/Upload/Translation.js

import { useState } from 'react';
import { translateText } from '../../services/api';

export default function Translation({ setResult }) {
    const [text, setText] = useState('');
    const [destLang, setDestLang] = useState('');
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    const handleTranslate = async () => {
        if (!text.trim()) {
            alert('Vui lòng nhập văn bản cần dịch.');
            return;
        }
        if (!destLang.trim()) {
            alert('Vui lòng nhập ngôn ngữ đích (ví dụ: "vi", "en").');
            return;
        }

        setLoading(true);
        setError(null);
        try {
            const response = await translateText(text, destLang);
            setResult(response.data);
        } catch (error) {
            console.error('Error in Translation:', error);
            setError('Có lỗi xảy ra khi dịch văn bản.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div>
            <h2>Translation</h2>
            <input
                type="text"
                value={text}
                onChange={(e) => setText(e.target.value)}
                placeholder="Enter text to translate"
            />
            <input
                type="text"
                value={destLang}
                onChange={(e) => setDestLang(e.target.value)}
                placeholder="Destination Language (e.g., 'vi', 'en')"
            />
            <button onClick={handleTranslate} disabled={loading}>
                {loading ? 'Đang dịch...' : 'Translate'}
            </button>
            {error && <p className="error">{error}</p>}
        </div>
    );
}
