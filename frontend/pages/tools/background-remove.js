import React, { useState } from 'react';
import { Upload, Image, Download } from 'lucide-react';

const DashboardLayout = ({ children, title }) => (
  <div className="min-h-screen bg-gray-50">
    <div className="flex h-screen overflow-hidden">
      <div className="flex flex-col flex-1 overflow-hidden">
        <header className="bg-white shadow-sm z-10">
          <div className="px-4 sm:px-6 lg:px-8">
            <div className="flex justify-between h-16">
              <h1 className="flex items-center text-2xl font-semibold text-gray-900">
                {title}
              </h1>
            </div>
          </div>
        </header>
        <main className="flex-1 relative overflow-y-auto">
          <div className="py-6">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 md:px-8">
              {children}
            </div>
          </div>
        </main>
      </div>
    </div>
  </div>
);

const BackgroundRemovePage = () => {
  const [imageFile, setImageFile] = useState(null);
  const [imagePreview, setImagePreview] = useState('');
  const [processedImage, setProcessedImage] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      setImageFile(file);
      const reader = new FileReader();
      reader.onload = () => setImagePreview(reader.result);
      reader.readAsDataURL(file);
    }
  };

  const handleRemoveBackground = async () => {
    if (!imageFile) {
      setError('Vui lòng tải lên hình ảnh trước.');
      return;
    }

    setLoading(true);
    try {
      await new Promise(resolve => setTimeout(resolve, 2000));
      setProcessedImage('/api/placeholder/400/400'); // Using placeholder for demo
      setError(null);
    } catch (err) {
      setError('Có lỗi xảy ra khi xử lý hình ảnh.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <DashboardLayout title="Xóa Nền Ảnh">
      <div className="bg-white rounded-lg shadow">
        <div className="p-6">
          {/* Image Upload */}
          <div className="mb-6">
            <div className="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center">
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
          </div>

          {/* Preview and Result */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Original Image */}
            {imagePreview && (
              <div className="border rounded-lg p-4">
                <h3 className="text-sm font-medium text-gray-700 mb-2">
                  Ảnh gốc
                </h3>
                <img 
                  src={imagePreview}
                  alt="Original"
                  className="w-full rounded-lg"
                />
              </div>
            )}

            {/* Processed Image */}
            {processedImage && (
              <div className="border rounded-lg p-4">
                <h3 className="text-sm font-medium text-gray-700 mb-2">
                  Ảnh đã xóa nền
                </h3>
                <img 
                  src={processedImage}
                  alt="Processed"
                  className="w-full rounded-lg bg-gray-100"
                />
                <button className="mt-4 flex items-center text-blue-600 hover:text-blue-800">
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
                ${loading || !imageFile ? 'bg-blue-400' : 'bg-blue-600 hover:bg-blue-700'}
              `}
            >
              <Image className="w-4 h-4 mr-2" />
              {loading ? 'Đang xử lý...' : 'Xóa nền ảnh'}
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