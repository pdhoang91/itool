import React, { useState } from 'react';
import { Mic, Upload, Download, Play } from 'lucide-react';
import DashboardLayout from '../../components/DashboardLayout';


const VoiceToTextPage = () => {
  const [audioFile, setAudioFile] = useState(null);
  const [convertedText, setConvertedText] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleConversion = async () => {
    if (!audioFile) {
      setError('Vui lòng tải lên file âm thanh trước.');
      return;
    }

    setLoading(true);
    try {
      await new Promise(resolve => setTimeout(resolve, 1500));
      setConvertedText('Sample converted text from voice...');
      setError(null);
    } catch (err) {
      setError('Có lỗi xảy ra khi chuyển đổi.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <DashboardLayout title="Chuyển Giọng Nói Thành Văn Bản">
      <div className="bg-white rounded-lg shadow">
        <div className="p-6">
          {/* Audio Upload */}
          <div className="mb-6">
            <div className="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center">
              <input
                type="file"
                accept="audio/*"
                onChange={(e) => setAudioFile(e.target.files[0])}
                className="hidden"
                id="audio-upload"
              />
              <label htmlFor="audio-upload" className="cursor-pointer">
                <div className="space-y-2">
                  <Upload className="mx-auto h-12 w-12 text-gray-400" />
                  <div className="text-gray-600">
                    Click để tải lên hoặc kéo thả file âm thanh vào đây
                  </div>
                  <div className="text-sm text-gray-500">
                    MP3, WAV up to 10MB
                  </div>
                </div>
              </label>
            </div>
          </div>

          {/* Audio Preview */}
          {audioFile && (
            <div className="mb-6">
              <h3 className="text-sm font-medium text-gray-700 mb-2">
                File âm thanh đã chọn
              </h3>
              <div className="bg-gray-50 p-4 rounded-lg">
                <div className="flex items-center justify-between">
                  <div className="flex items-center">
                    <Play className="w-4 h-4 mr-2 text-gray-500" />
                    <span className="text-sm text-gray-600">
                      {audioFile.name}
                    </span>
                  </div>
                  <span className="text-sm text-gray-500">
                    {(audioFile.size / (1024 * 1024)).toFixed(2)} MB
                  </span>
                </div>
              </div>
            </div>
          )}

          {/* Conversion Result */}
          {convertedText && (
            <div className="mb-6">
              <h3 className="text-sm font-medium text-gray-700 mb-2">
                Văn bản chuyển đổi
              </h3>
              <div className="bg-gray-50 p-4 rounded-lg">
                <p className="whitespace-pre-wrap">{convertedText}</p>
              </div>
              <button className="mt-4 flex items-center text-blue-600 hover:text-blue-800">
                <Download className="w-4 h-4 mr-2" />
                Tải xuống văn bản
              </button>
            </div>
          )}

          {/* Action Button */}
          <div className="mt-6 flex justify-end">
            <button
              onClick={handleConversion}
              disabled={loading || !audioFile}
              className={`
                flex items-center px-6 py-2 rounded-lg text-white
                ${loading || !audioFile ? 'bg-blue-400' : 'bg-blue-600 hover:bg-blue-700'}
              `}
            >
              <Mic className="w-4 h-4 mr-2" />
              {loading ? 'Đang xử lý...' : 'Chuyển đổi giọng nói'}
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

export default VoiceToTextPage;