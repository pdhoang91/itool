// components/Upload/TextToVoice.js

import React, { useState, useRef } from 'react';
import { Upload, Volume2, Play, Download } from 'lucide-react';

const TextToVoice = () => {
  const [text, setText] = useState('');
  const [selectedVoice, setSelectedVoice] = useState('Nam Minh');
  const [speed, setSpeed] = useState(1);
  const [pitch, setPitch] = useState(1);
  const [volume, setVolume] = useState(1);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [audioUrl, setAudioUrl] = useState(null);
  const audioRef = useRef(null);

  const voices = [
    { name: 'Hoài My', gender: 'female', lang: 'vi' },
    { name: 'Hương Thảo', gender: 'female', lang: 'vi' },
    { name: 'Thúy Tiên', gender: 'female', lang: 'vi' },
    { name: 'Thúy Linh', gender: 'female', lang: 'vi' },
    { name: 'Nam Minh', gender: 'male', lang: 'vi' },
  ];

  const handleConvert = async () => {
    if (!text.trim()) {
      setError('Vui lòng nhập văn bản cần chuyển đổi.');
      return;
    }

    setLoading(true);
    setError(null);
    setAudioUrl(null);

    try {
      const response = await fetch('/api/tts', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          text,
          voice: selectedVoice,
          language: 'vi',
          speed,
          pitch,
          volume,
        }),
      });

      if (!response.ok) {
        throw new Error('Conversion failed');
      }

      const data = await response.json();
      setAudioUrl(data.audio_url);

      if (audioRef.current) {
        audioRef.current.src = data.audio_url;
      }
    } catch (error) {
      console.error('Error in Text-to-Voice:', error);
      setError('Có lỗi xảy ra khi chuyển đổi văn bản thành giọng nói.');
    } finally {
      setLoading(false);
    }
  };

  const handleFileUpload = async (event) => {
    const file = event.target.files[0];
    if (!file) return;

    try {
      const text = await file.text();
      setText(text);
    } catch (error) {
      setError('Không thể đọc file. Vui lòng thử lại.');
      console.error('Không thể đọc file:', error);
    }
  };

  return (
    <div className="max-w-4xl mx-auto p-6 bg-white rounded-lg shadow">
      <h2 className="text-2xl font-bold mb-4">Chuyển văn bản thành giọng nói</h2>
      <p className="text-gray-600 mb-6">
        Hỗ trợ nhiều giọng đọc tự nhiên với khả năng chuyển đổi tối đa 10.000 ký tự mỗi lần
      </p>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
        <div>
          <label className="block text-sm font-medium mb-2">Tốc độ</label>
          <input
            type="range"
            min="0.5"
            max="2"
            step="0.1"
            value={speed}
            onChange={(e) => setSpeed(parseFloat(e.target.value))}
            className="w-full"
          />
        </div>
        <div>
          <label className="block text-sm font-medium mb-2">Cao độ</label>
          <input
            type="range"
            min="0.5"
            max="2"
            step="0.1"
            value={pitch}
            onChange={(e) => setPitch(parseFloat(e.target.value))}
            className="w-full"
          />
        </div>
        <div>
          <label className="block text-sm font-medium mb-2">Âm lượng</label>
          <input
            type="range"
            min="0.5"
            max="2"
            step="0.1"
            value={volume}
            onChange={(e) => setVolume(parseFloat(e.target.value))}
            className="w-full"
          />
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <textarea
            className="w-full h-48 p-3 border rounded-lg resize-none"
            placeholder="Nhập văn bản cần chuyển đổi..."
            value={text}
            onChange={(e) => setText(e.target.value)}
          />
          <div className="mt-4 flex justify-between items-center">
            <span className="text-sm text-gray-500">
              {text.length}/10,000 ký tự
            </span>
            <div className="space-x-2">
              <button className="px-4 py-2 border rounded hover:bg-gray-50">
                <input
                  type="file"
                  accept=".txt,.srt"
                  onChange={handleFileUpload}
                  className="hidden"
                  id="file-upload"
                />
                <label htmlFor="file-upload" className="cursor-pointer flex items-center">
                  <Upload className="w-4 h-4 mr-2" />
                  Upload
                </label>
              </button>
              <button 
                onClick={() => setText('')}
                className="px-4 py-2 border rounded hover:bg-gray-50"
              >
                Clear
              </button>
            </div>
          </div>
        </div>

        <div>
          <div className="mb-4">
            <label className="block text-sm font-medium mb-2">Giọng nói</label>
            <div className="space-y-2 max-h-48 overflow-y-auto">
              {voices.map((voice) => (
                <div
                  key={voice.name}
                  className={`p-3 border rounded-lg cursor-pointer flex items-center justify-between ${
                    selectedVoice === voice.name ? 'border-blue-500 bg-blue-50' : ''
                  }`}
                  onClick={() => setSelectedVoice(voice.name)}
                >
                  <div className="flex items-center">
                    {voice.gender === 'female' ? '👩' : '👨'} {voice.name}
                  </div>
                  <Play className="w-4 h-4" />
                </div>
              ))}
            </div>
          </div>

          <button 
            className={`w-full py-2 px-4 rounded-lg text-white ${
              loading ? 'bg-blue-400' : 'bg-blue-600 hover:bg-blue-700'
            }`}
            onClick={handleConvert}
            disabled={loading}
          >
            <Volume2 className="w-4 h-4 mr-2 inline" />
            {loading ? 'Đang chuyển đổi...' : 'Chuyển đổi Text thành Voice'}
          </button>
        </div>
      </div>

      {error && (
        <div className="mt-4 p-3 bg-red-100 text-red-700 rounded-lg">
          {error}
        </div>
      )}

      {audioUrl && (
        <div className="mt-6">
          <audio ref={audioRef} controls className="w-full" src={audioUrl} />
          <div className="mt-2 text-center">
            <a
              href={audioUrl}
              download="voice_output.mp3"
              className="inline-flex items-center text-blue-600 hover:text-blue-800"
            >
              <Download className="w-4 h-4 mr-2" />
              Download MP3
            </a>
          </div>
        </div>
      )}
    </div>
  );
};

export default TextToVoice;