import React, { useState } from 'react';
import { Upload, Camera, Download } from 'lucide-react';
import DashboardLayout from '../../components/DashboardLayout';


const FaceRecognitionPage = () => {
  const [imageFile, setImageFile] = useState(null);
  const [imagePreview, setImagePreview] = useState('');
  const [recognitionResults, setRecognitionResults] = useState(null);
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

  const handleRecognition = async () => {
    if (!imageFile) {
      setError('Vui lòng tải lên hình ảnh trước.');
      return;
    }

    setLoading(true);
    try {
      await new Promise(resolve => setTimeout(resolve, 1500));
      setRecognitionResults({
        faces: [
          { confidence: 98.5, age: 25, gender: 'Nam' },
          { confidence: 96.2, age: 30, gender: 'Nữ' }
        ]
      });
      setError(null);
    } catch (err) {
      setError('Có lỗi xảy ra khi nhận diện khuôn mặt.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <DashboardLayout title="Nhận Diện Khuôn Mặt">
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
                  <Camera className="mx-auto h-12 w-12 text-gray-400" />
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

          {/* Preview and Results */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* Image Preview */}
            {imagePreview && (
              <div className="border rounded-lg p-4">
                <h3 className="text-sm font-medium text-gray-700 mb-2">
                  Hình ảnh đã chọn
                </h3>
                <img 
                  src={imagePreview}
                  alt="Preview"
                  className="w-full rounded-lg"
                />
              </div>
            )}

            {/* Recognition Results */}
            {recognitionResults && (
              <div className="border rounded-lg p-4">
                <h3 className="text-sm font-medium text-gray-700 mb-2">
                  Kết quả nhận diện
                </h3>
                <div className="space-y-4">
                  {recognitionResults.faces.map((face, index) => (
                    <div key={index} className="bg-gray-50 p-4 rounded-lg">
                      <h4 className="font-medium">Khuôn mặt #{index + 1}</h4>
                      <ul className="mt-2 space-y-1 text-sm">
                        <li>Độ chính xác: {face.confidence}%</li>
                        <li>Độ tuổi: {face.age}</li>
                        <li>Giới tính: {face.gender}</li>
                      </ul>
                    </div>
                  ))}
                </div>
                <button className="mt-4 flex items-center text-blue-600 hover:text-blue-800">
                  <Download className="w-4 h-4 mr-2" />
                  Tải xuống kết quả
                </button>
              </div>
            )}
          </div>

          {/* Action Button */}
          <div className="mt-6 flex justify-end">
            <button
              onClick={handleRecognition}
              disabled={loading || !imageFile}
              className={`
                flex items-center px-6 py-2 rounded-lg text-white
                ${loading || !imageFile ? 'bg-blue-400' : 'bg-blue-600 hover:bg-blue-700'}
              `}
            >
              <Camera className="w-4 h-4 mr-2" />
              {loading ? 'Đang xử lý...' : 'Nhận diện khuôn mặt'}
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

export default FaceRecognitionPage;