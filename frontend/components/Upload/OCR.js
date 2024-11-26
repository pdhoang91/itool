import { useState } from 'react';
import { performOCR } from '../../services/api';
import Image from 'next/image';

export default function OCR({ setResult }) {
    const [imageFile, setImageFile] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);
    const [ocrText, setOcrText] = useState('');
    const [imagePreview, setImagePreview] = useState('');

    const handleImageChange = (e) => {
        const file = e.target.files[0];
        if (file) {
            setImageFile(file);
            // Tạo preview cho hình ảnh
            const reader = new FileReader();
            reader.onload = () => {
                setImagePreview(reader.result);
            };
            reader.readAsDataURL(file);
        }
    };

    const handleOCR = async () => {
        if (!imageFile) {
            alert('Vui lòng tải lên tệp hình ảnh trước.');
            return;
        }

        setLoading(true);
        setError(null);
        try {
            const response = await performOCR(imageFile);
            
            // Kiểm tra cấu trúc response và xử lý dữ liệu phù hợp
            let extractedText = '';
            if (response && response.data) {
                if (typeof response.data === 'string') {
                    extractedText = response.data;
                } else if (response.data.text) {
                    extractedText = response.data.text;
                } else if (Array.isArray(response.data)) {
                    extractedText = response.data.join('\n');
                }
            }
            
            setResult(extractedText);
            setOcrText(extractedText);
        } catch (error) {
            console.error('Error in OCR:', error);
            setError('Có lỗi xảy ra khi thực hiện OCR.');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="p-4">
            <h2 className="text-xl font-bold mb-4">OCR (Optical Character Recognition)</h2>
            
            {/* Phần upload và xử lý */}
            <div className="mb-4">
                <input
                    type="file"
                    accept="image/*"
                    onChange={handleImageChange}
                    className="mb-2"
                />
                <button 
                    onClick={handleOCR} 
                    disabled={loading}
                    className="bg-blue-500 text-white px-4 py-2 rounded disabled:bg-gray-400"
                >
                    {loading ? 'Đang xử lý...' : 'Thực hiện OCR'}
                </button>
            </div>

            {error && <p className="text-red-500 mb-4">{error}</p>}

            {/* Hiển thị preview hình ảnh và kết quả */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {/* Preview hình ảnh */}
                {imagePreview && (
                    <div className="border p-4 rounded">
                        <h3 className="font-semibold mb-2">Hình ảnh đã chọn:</h3>
                        <Image 
                            src={imagePreview} 
                            alt="Preview" 
                            className="max-w-full h-auto"
                        />
                    </div>
                )}

                {/* Kết quả OCR */}
                {ocrText && (
                    <div className="border p-4 rounded">
                        <h3 className="font-semibold mb-2">Kết quả OCR:</h3>
                        <div className="whitespace-pre-wrap bg-gray-50 p-3 rounded">
                            {ocrText}
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}