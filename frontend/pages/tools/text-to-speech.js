import React, { useState } from 'react';
import { Volume2, Upload, Download } from 'lucide-react';

// Dashboard Layout Component
const DashboardLayout = ({ children, title }) => {
  return (
    <div className="min-h-screen bg-gray-50">
      <div className="flex h-screen overflow-hidden">
        <div className="flex flex-col flex-1 overflow-hidden">
          {/* Header */}
          <header className="bg-white shadow-sm z-10">
            <div className="px-4 sm:px-6 lg:px-8">
              <div className="flex justify-between h-16">
                <div className="flex">
                  <h1 className="flex items-center text-2xl font-semibold text-gray-900">
                    {title}
                  </h1>
                </div>
              </div>
            </div>
          </header>

          {/* Main Content */}
          <main className="flex-1 relative overflow-y-auto focus:outline-none">
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
};

// Text to Speech Page Component
const TextToSpeechPage = () => {
  const [text, setText] = useState('');
  const [selectedVoice, setSelectedVoice] = useState('Nam Minh');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [audioUrl, setAudioUrl] = useState(null);

  const voices = [
    { name: 'Hoài My', gender: 'female', lang: 'vi' },
    { name: 'Nam Minh', gender: 'male', lang: 'vi' },
  ];

  const handleConvert = async () => {
    if (!text.trim()) {
      setError('Vui lòng nhập văn bản cần chuyển đổi.');
      return;
    }
    // Handle conversion logic
    setLoading(true);
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1500));
      setAudioUrl('sample-audio-url.mp3');
      setError(null);
    } catch (err) {
      setError('Có lỗi xảy ra khi chuyển đổi.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <DashboardLayout title="Text to Speech">
      <div className="bg-white rounded-lg shadow">
        <div className="p-6">
          {/* Voice Selection */}
          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Chọn Giọng Đọc
            </label>
            <div className="grid grid-cols-2 gap-4">
              {voices.map((voice) => (
                <button
                  key={voice.name}
                  onClick={() => setSelectedVoice(voice.name)}
                  className={`
                    p-4 rounded-lg border flex items-center justify-between
                    ${selectedVoice === voice.name 
                      ? 'border-blue-500 bg-blue-50 text-blue-700' 
                      : 'border-gray-200 hover:border-gray-300'
                    }
                  `}
                >
                  <div className="flex items-center">
                    {voice.gender === 'female' ? '👩' : '👨'} 
                    <span className="ml-2">{voice.name}</span>
                  </div>
                </button>
              ))}
            </div>
          </div>

          {/* Text Input */}
          <div className="mb-6">
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Văn bản cần chuyển đổi
            </label>
            <textarea
              value={text}
              onChange={(e) => setText(e.target.value)}
              className="w-full h-40 p-3 border rounded-lg resize-none"
              placeholder="Nhập văn bản cần chuyển đổi..."
            />
            <div className="mt-2 flex justify-between text-sm text-gray-500">
              <span>{text.length}/10,000 ký tự</span>
              <button 
                onClick={() => setText('')}
                className="text-gray-600 hover:text-gray-900"
              >
                Xóa
              </button>
            </div>
          </div>

          {/* Actions */}
          <div className="flex justify-between items-center">
            <button className="flex items-center px-4 py-2 border rounded-lg hover:bg-gray-50">
              <Upload className="w-4 h-4 mr-2" />
              Upload File
            </button>
            <button
              onClick={handleConvert}
              disabled={loading}
              className={`
                flex items-center px-6 py-2 rounded-lg text-white
                ${loading ? 'bg-blue-400' : 'bg-blue-600 hover:bg-blue-700'}
              `}
            >
              <Volume2 className="w-4 h-4 mr-2" />
              {loading ? 'Đang xử lý...' : 'Chuyển đổi'}
            </button>
          </div>

          {/* Error Message */}
          {error && (
            <div className="mt-4 p-3 bg-red-50 text-red-700 rounded-lg">
              {error}
            </div>
          )}

          {/* Audio Player */}
          {audioUrl && (
            <div className="mt-6 p-4 border rounded-lg">
              <audio controls className="w-full mb-4" src={audioUrl} />
              <button className="flex items-center text-blue-600 hover:text-blue-800">
                <Download className="w-4 h-4 mr-2" />
                Tải xuống MP3
              </button>
            </div>
          )}
        </div>
      </div>
    </DashboardLayout>
  );
};

export default TextToSpeechPage;