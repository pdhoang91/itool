// frontend/components/Upload/OCR.js

import { useState } from 'react';
import { performOCR } from '../../services/api';

export default function OCR({ setResult }) {
    const [imageFile, setImageFile] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    const handleOCR = async () => {
        if (!imageFile) {
            alert('Vui lòng tải lên tệp hình ảnh trước.');
            return;
        }

        setLoading(true);
        setError(null);
        try {
            const response = await performOCR(imageFile);
            setResult(response.data);
        } catch (error) {
            console.error('Error in OCR:', error);
            setError('Có lỗi xảy ra khi thực hiện OCR.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div>
            <h2>OCR (Optical Character Recognition)</h2>
            <input
                type="file"
                accept="image/*"
                onChange={(e) => setImageFile(e.target.files[0])}
            />
            <button onClick={handleOCR} disabled={loading}>
                {loading ? 'Đang xử lý...' : 'Perform OCR'}
            </button>
            {error && <p className="error">{error}</p>}
        </div>
    );
}
