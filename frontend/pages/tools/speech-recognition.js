import React, { useState, useRef } from 'react';
import { Mic, StopCircle, Play, Download } from 'lucide-react';
import DashboardLayout from '../../components/DashboardLayout';

const SpeechRecognitionPage = () => {
  const [isRecording, setIsRecording] = useState(false);
  const [audioFile, setAudioFile] = useState(null);
  const [recognizedText, setRecognizedText] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const toggleRecording = () => {
    setIsRecording(!isRecording);
  };

  const handleRecognition = async () => {
    setLoading(true);
    try {
      await new Promise(resolve => setTimeout(resolve, 1500));
      setRecognizedText('Sample recognized text...');
      setError(null);
    } catch (err) {
      setError('Có lỗi xảy ra khi nhận dạng giọng nói.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <DashboardLayout title="Nhận Dạng Giọng Nói">
      <div className="bg-white rounded-lg shadow">
        <div className="p-6">
          {/* Recording Interface */}
          <div className="text-center py-8">
            <button
              onClick={toggleRecording}
              className={`
                p-8 rounded-full
                ${isRecording 
                  ? 'bg-red-100 hover:bg-red-200 text-red-600' 
                  : 'bg-blue-100 hover:bg-blue-200 text-blue-600'}
              `}
            >
              {isRecording ? (
                <StopCircle className="w-12 h-12" />
              ) : (
                <Mic className="w-12 h-12" />
              )}
            </button>
            <p className="mt-4 text-sm text-gray-600">
              {isRecording ? 'Đang ghi âm...' : 'Nhấn để bắt đầu ghi âm'}
            </p>
          </div>

          {/* Audio Upload */}
          <div className="mt-6">
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
                  <Play className="mx-auto h-12 w-12 text-gray-400" />
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

          {/* Recognition Result */}
          {recognizedText && (
            <div className="mt-6">
              <h3 className="text-sm font-medium text-gray-700 mb-2">
                Văn bản nhận dạng được
              </h3>
              <div className="bg-gray-50 p-4 rounded-lg">
                <p className="whitespace-pre-wrap">{recognizedText}</p>
              </div>
              <button className="mt-4 flex items-center text-blue-600 hover:text-blue-800">
                <Download className="w-4 h-4 mr-2" />
                Tải xuống văn bản
              </button>
            </div>
          )}

          {/* Action Buttons */}
          <div className="mt-6 flex justify-end">
            <button
              onClick={handleRecognition}
              disabled={loading || (!audioFile && !isRecording)}
              className={`
                flex items-center px-6 py-2 rounded-lg text-white
                ${loading || (!audioFile && !isRecording) 
                  ? 'bg-blue-400' 
                  : 'bg-blue-600 hover:bg-blue-700'}
              `}
            >
              <Mic className="w-4 h-4 mr-2" />
              {loading ? 'Đang xử lý...' : 'Nhận dạng giọng nói'}
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

export default SpeechRecognitionPage;