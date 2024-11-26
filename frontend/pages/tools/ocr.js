import React, { useState } from 'react';
import { Upload, Camera, Download } from 'lucide-react';
import DashboardLayout from '../../components/DashboardLayout';
import { performOCR } from '../../services/api';

const OCRPage = () => {
  const [imageFile, setImageFile] = useState(null);
  const [imagePreview, setImagePreview] = useState('');
  const [ocrText, setOcrText] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const validateFile = (file) => {
    const validTypes = ['image/jpeg', 'image/png', 'image/gif'];
    const maxSize = 10 * 1024 * 1024; // 10MB

    if (!validTypes.includes(file.type)) {
      setError('Định dạng file không được hỗ trợ. Vui lòng sử dụng JPG, PNG hoặc GIF.');
      return false;
    }

    if (file.size > maxSize) {
      setError('File không được vượt quá 10MB.');
      return false;
    }

    return true;
  };

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      if (!validateFile(file)) {
        return;
      }

      setImageFile(file);
      const objectUrl = URL.createObjectURL(file);
      setImagePreview(objectUrl);

      return () => URL.revokeObjectURL(objectUrl);
    }
  };

  const handleOCR = async () => {
    if (!imageFile) {
      setError('Vui lòng tải lên hình ảnh trước.');
      return;
    }

    setLoading(true);
    try {
      const response = await performOCR(imageFile);
      setOcrText(response.data.text);
      setError(null);
    } catch (err) {
      setError('Có lỗi xảy ra khi xử lý OCR: ' + (err.response?.data?.message || err.message));
    } finally {
      setLoading(false);
    }
  };

  const handleDownloadText = () => {
    const blob = new Blob([ocrText], { type: 'text/plain' });
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = 'ocr-result.txt';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
  };

  return (
    <DashboardLayout title="OCR - Nhận Dạng Văn Bản">
      <div className="bg-white rounded-lg shadow">
        <div className="p-6">
          {/* Upload Section */}
          <div 
            className="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center"
            onDragOver={(e) => e.preventDefault()}
            onDrop={(e) => {
              e.preventDefault();
              const file = e.dataTransfer.files[0];
              if (file && validateFile(file)) {
                setImageFile(file);
                const objectUrl = URL.createObjectURL(file);
                setImagePreview(objectUrl);
              }
            }}
          >
            <input
              type="file"
              accept="image/*"
              onChange={handleImageChange}
              className="hidden"
              id="image-upload"
            />
            <label htmlFor="image-upload" className="cursor-pointer">
              <div className="space-y-2">
                <Camera className="mx-auto h-12 w-12 text-gray-400" />
                <div className="text-gray-600">
                  Click để tải lên hoặc kéo thả hình ảnh vào đây
                </div>
                <div className="text-sm text-gray-500">
                  PNG, JPG, GIF up to 10MB
                </div>
              </div>
            </label>
          </div>

          {/* Preview and Result */}
          <div className="mt-8 grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Image Preview */}
            {imagePreview && (
              <div className="border rounded-lg p-4">
                <h3 className="text-sm font-medium text-gray-700 mb-2">
                  Hình ảnh đã chọn
                </h3>
                <div className="relative aspect-video">
                  <img 
                    src={imagePreview}
                    alt="Preview"
                    className="w-full h-full object-contain rounded-lg"
                  />
                </div>
              </div>
            )}

            {/* OCR Result */}
            {ocrText && (
              <div className="border rounded-lg p-4">
                <h3 className="text-sm font-medium text-gray-700 mb-2">
                  Kết quả OCR
                </h3>
                <div className="bg-gray-50 p-3 rounded-lg max-h-[300px] overflow-y-auto">
                  <pre className="whitespace-pre-wrap text-sm">
                    {ocrText}
                  </pre>
                </div>
                <button 
                  onClick={handleDownloadText}
                  className="mt-4 flex items-center text-blue-600 hover:text-blue-800"
                >
                  <Download className="w-4 h-4 mr-2" />
                  Tải xuống văn bản
                </button>
              </div>
            )}
          </div>

          {/* Action Button */}
          <div className="mt-6 flex justify-end">
            <button
              onClick={handleOCR}
              disabled={loading || !imageFile}
              className={`
                flex items-center px-6 py-2 rounded-lg text-white
                ${loading || !imageFile ? 'bg-blue-400 cursor-not-allowed' : 'bg-blue-600 hover:bg-blue-700'}
              `}
            >
              {loading ? (
                <div className="flex items-center">
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                  Đang xử lý...
                </div>
              ) : (
                <>
                  <Camera className="w-4 h-4 mr-2" />
                  Nhận dạng văn bản
                </>
              )}
            </button>
          </div>

          {/* Error Message */}
          {error && (
            <div className="mt-4 p-3 bg-red-50 text-red-700 rounded-lg">
              {error}
            </div>
          )}
        </div>
      </div>
    </DashboardLayout>
  );
};

export default OCRPage;