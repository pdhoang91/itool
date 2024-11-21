// frontend/components/Upload/RemoveBackground.js

import { useState } from 'react';
import { removeBackground } from '../../services/api';

export default function RemoveBackground() {
    const [imageFile, setImageFile] = useState(null);
    const [result, setResult] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    const handleRemoveBg = async () => {
        if (!imageFile) {
            alert('Please upload an image file first.');
            return;
        }

        setLoading(true);
        setError(null);
        try {
            const response = await removeBackground(imageFile);
            setResult(response.data);
        } catch (error) {
            console.error('Error in Background Removal:', error);
            setError('Có lỗi xảy ra khi loại bỏ nền ảnh.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div>
            <h2>Background Removal</h2>
            <input
                type="file"
                accept="image/*"
                onChange={(e) => setImageFile(e.target.files[0])}
            />
            <button onClick={handleRemoveBg} disabled={loading}>
                {loading ? 'Đang xử lý...' : 'Remove Background'}
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
