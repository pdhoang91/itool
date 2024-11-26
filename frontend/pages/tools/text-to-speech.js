import React, { useState, useRef } from 'react';
import { Volume2, Upload, Download } from 'lucide-react';
import DashboardLayout from '../../components/DashboardLayout';
import { convertTextToVoice, uploadAudio } from '../../services/api';

const TextToSpeechPage = () => {
  const [text, setText] = useState('');
  const [selectedVoice, setSelectedVoice] = useState('female_1');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [audioUrl, setAudioUrl] = useState(null);
  const [speed, setSpeed] = useState(1.0);
  const [pitch, setPitch] = useState(1.0);
  const [previewPlaying, setPreviewPlaying] = useState(false);
  const audioRef = useRef(null);
  const fileInputRef = useRef(null);

  const voices = [
    { id: 'female_1', name: 'Hoài My', gender: 'female', lang: 'vi' },
    { id: 'male_1', name: 'Nam Minh', gender: 'male', lang: 'vi' },
    { id: 'female_2', name: 'Thu Hà', gender: 'female', lang: 'vi' },
    { id: 'male_2', name: 'Duy Anh', gender: 'male', lang: 'vi' }
  ];

  // Handle text file upload
  const handleFileUpload = async (e) => {
    const file = e.target.files[0];
    if (file) {
      if (file.type !== 'text/plain') {
        setError('Vui lòng chọn file .txt');
        return;
      }

      try {
        const reader = new FileReader();
        reader.onload = (e) => {
          setText(e.target.result);
        };
        reader.readAsText(file);
      } catch (err) {
        setError('Có lỗi khi đọc file.');
      }
    }
  };

  // Preview current text section
  const TextPreview = () => {
    if (!text) return null;

    const words = text.trim().split(/\s+/).length;
    const chars = text.length;

    return (
      <div className="mt-4 p-4 border rounded-lg bg-gray-50">
        <div className="flex justify-between items-center mb-2">
          <h3 className="text-sm font-medium text-gray-700">Preview</h3>
          <div className="text-sm text-gray-500">
            {words} từ | {chars} ký tự
          </div>
        </div>
        <div className="prose prose-sm max-w-none">
          {text.split('\n').map((paragraph, idx) => (
            <p key={idx} className="mb-2">
              {paragraph}
            </p>
          ))}
        </div>
      </div>
    );
  };

  // Handle conversion
  const handleConvert = async () => {
    if (!text.trim()) {
      setError('Vui lòng nhập văn bản cần chuyển đổi.');
      return;
    }

    setLoading(true);
    try {
      const response = await convertTextToVoice({
        text,
        voice: selectedVoice,
        language: 'vi',
        speed,
        pitch,
        volume: 1.0
      });
      
      setAudioUrl(response.data.audioUrl);
      setError(null);

      // Auto play preview
      if (audioRef.current) {
        audioRef.current.src = response.data.audioUrl;
        audioRef.current.play();
        setPreviewPlaying(true);
      }
    } catch (err) {
      setError('Có lỗi xảy ra khi chuyển đổi: ' + (err.response?.data?.message || err.message));
    } finally {
      setLoading(false);
    }
  };

  // Handle download
  const handleDownload = async () => {
    try {
      const response = await fetch(audioUrl);
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = 'audio.mp3';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    } catch (err) {
      setError('Có lỗi xảy ra khi tải xuống file.');
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
                  key={voice.id}
                  onClick={() => setSelectedVoice(voice.id)}
                  className={`
                    p-4 rounded-lg border flex items-center justify-between
                    ${selectedVoice === voice.id 
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

          {/* Speed & Pitch Controls */}
          <div className="mb-6 grid grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Tốc độ đọc
              </label>
              <input
                type="range"
                min="0.5"
                max="2"
                step="0.1"
                value={speed}
                onChange={(e) => setSpeed(parseFloat(e.target.value))}
                className="w-full"
              />
              <div className="text-sm text-gray-500 mt-1">{speed}x</div>
            </div>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Độ cao giọng
              </label>
              <input
                type="range"
                min="0.5"
                max="2"
                step="0.1"
                value={pitch}
                onChange={(e) => setPitch(parseFloat(e.target.value))}
                className="w-full"
              />
              <div className="text-sm text-gray-500 mt-1">{pitch}x</div>
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
              className="w-full h-40 p-3 border rounded-lg resize-none focus:ring-blue-500 focus:border-blue-500"
              placeholder="Nhập văn bản cần chuyển đổi..."
              maxLength={10000}
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

          {/* Text Preview */}
          <TextPreview />

          {/* Actions */}
          <div className="flex justify-between items-center mt-6">
            <button 
              onClick={() => fileInputRef.current?.click()}
              className="flex items-center px-4 py-2 border rounded-lg hover:bg-gray-50"
            >
              <Upload className="w-4 h-4 mr-2" />
              Upload File Text
            </button>
            <input
              ref={fileInputRef}
              type="file"
              accept=".txt"
              onChange={handleFileUpload}
              className="hidden"
            />
            <button
              onClick={handleConvert}
              disabled={loading || !text.trim()}
              className={`
                flex items-center px-6 py-2 rounded-lg text-white
                ${loading || !text.trim() ? 'bg-blue-400 cursor-not-allowed' : 'bg-blue-600 hover:bg-blue-700'}
              `}
            >
              {loading ? (
                <div className="flex items-center">
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                  Đang xử lý...
                </div>
              ) : (
                <>
                  <Volume2 className="w-4 h-4 mr-2" />
                  Chuyển đổi
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

          {/* Audio Preview */}
          {audioUrl && (
            <div className="mt-6 p-4 border rounded-lg">
              <h3 className="text-sm font-medium text-gray-700 mb-2">
                Preview Audio
              </h3>
              <audio 
                ref={audioRef}
                controls 
                className="w-full mb-4"
                onPlay={() => setPreviewPlaying(true)}
                onPause={() => setPreviewPlaying(false)}
                src={audioUrl}
              />
              <button 
                onClick={handleDownload}
                className="flex items-center text-blue-600 hover:text-blue-800"
              >
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