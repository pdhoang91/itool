import { useState } from 'react';
import { removeBackground } from '../../services/api';

export default function RemoveBackground() {
    const [imageFile, setImageFile] = useState(null);
    const [resultUrl, setResultUrl] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    const handleRemoveBg = async () => {
        if (!imageFile) {
            alert('Please upload an image file first.');
            return;
        }

        setLoading(true);
        setError(null);
        setResultUrl(null); // Reset result before processing
        try {
            const response = await removeBackground(imageFile);
            const { processed_image_path } = response.data; // Đường dẫn ảnh từ backend
            setResultUrl(`${process.env.NEXT_PUBLIC_API_BASE_URL}/images/${processed_image_path}`);
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
            {resultUrl && (
                <div>
                    <h3>Result:</h3>
                    <img
                        src={resultUrl}
                        alt="Processed"
                        style={{ maxWidth: '100%', marginBottom: '1em' }}
                    />
                    <a href={resultUrl} download="processed_image.jpg">
                        <button>Download Image</button>
                    </a>
                </div>
            )}
        </div>
    );
}
