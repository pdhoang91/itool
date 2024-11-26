import React, { useState, useEffect } from 'react';
import { Upload, Image, Download } from 'lucide-react';
import DashboardLayout from '../../components/DashboardLayout';
import { removeBackground } from '../../services/api';

const BackgroundRemovePage = () => {
  const [imageFile, setImageFile] = useState(null);
  const [imagePreview, setImagePreview] = useState('');
  const [processedImage, setProcessedImage] = useState('');
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

      // Cleanup old preview URL
      return () => URL.revokeObjectURL(objectUrl);
    }
  };

  const handleRemoveBackground = async () => {
    if (!imageFile) {
      setError('Vui lòng tải lên hình ảnh trước.');
      return;
    }

    setLoading(true);
    try {
      const response = await removeBackground(imageFile);
      console.log("response", response)
      setProcessedImage(response.data.processedImageUrl);
      setError(null);
    } catch (err) {
      setError('Có lỗi xảy ra khi xử lý hình ảnh: ' + (err.response?.data?.message || err.message));
    } finally {
      setLoading(false);
    }
  };

  const handleDownload = async () => {
    try {
      const response = await fetch(processedImage);
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = 'processed-image.png';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    } catch (err) {
      setError('Có lỗi xảy ra khi tải xuống ảnh.');
    }
  };

  return (
    <DashboardLayout title="Xóa Nền Ảnh">
      <div className="bg-white rounded-lg shadow">
        <div className="p-6">
          {/* Drop Zone */}
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
                <Image className="mx-auto h-12 w-12 text-gray-400" />
                <div className="text-gray-600">
                  Click để tải lên hoặc kéo thả hình ảnh vào đây
                </div>
                <div className="text-sm text-gray-500">
                  PNG, JPG lên đến 10MB
                </div>
              </div>
            </label>
          </div>

          {/* Preview Section */}
          <div className="mt-8 grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Original Image Preview */}
            {imagePreview && (
              <div className="border rounded-lg p-4">
                <h3 className="text-sm font-medium text-gray-700 mb-2">
                  Ảnh gốc
                </h3>
                <div className="relative aspect-video">
                  <img 
                    src={imagePreview}
                    alt="Original"
                    className="w-full h-full object-contain rounded-lg"
                  />
                </div>
              </div>
            )}

            {/* Processed Image Preview */}
            {processedImage && (
              <div className="border rounded-lg p-4">
                <h3 className="text-sm font-medium text-gray-700 mb-2">
                  Ảnh đã xóa nền
                </h3>
                <div className="relative aspect-video bg-gray-100">
                  <img 
                    src={processedImage}
                    alt="Processed"
                    className="w-full h-full object-contain rounded-lg"
                  />
                </div>
                <button 
                  onClick={handleDownload}
                  className="mt-4 flex items-center text-blue-600 hover:text-blue-800"
                >
                  <Download className="w-4 h-4 mr-2" />
                  Tải xuống ảnh
                </button>
              </div>
            )}
          </div>

          {/* Action Button */}
          <div className="mt-6 flex justify-end">
            <button
              onClick={handleRemoveBackground}
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
                  <Image className="w-4 h-4 mr-2" />
                  Xóa nền ảnh
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

export default BackgroundRemovePage;