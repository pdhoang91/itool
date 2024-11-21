// frontend/components/Upload/FaceRecognition.js

import { useState } from 'react';
import { recognizeFace } from '../../services/api';

export default function FaceRecognition({ setResult }) {
    const [imageFile, setImageFile] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    const handleFaceRecognition = async () => {
        if (!imageFile) {
            alert('Please upload an image file first.');
            return;
        }

        setLoading(true);
        setError(null);
        try {
            const response = await recognizeFace(imageFile);
            setResult(response.data);
        } catch (error) {
            console.error('Error in Face Recognition:', error);
            setError('Có lỗi xảy ra khi nhận diện khuôn mặt.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div>
            <h2>Face Recognition</h2>
            <input
                type="file"
                accept="image/*"
                onChange={(e) => setImageFile(e.target.files[0])}
            />
            <button onClick={handleFaceRecognition} disabled={loading}>
                {loading ? 'Đang xử lý...' : 'Recognize Face'}
            </button>
            {error && <p className="error">{error}</p>}
        </div>
    );
}
