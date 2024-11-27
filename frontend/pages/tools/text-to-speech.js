import React, { useState, useRef, useEffect } from 'react';
import { Volume2, Upload, Download } from 'lucide-react';
import DashboardLayout from '../../components/DashboardLayout';
import { convertTextToVoice } from '../../services/api';
import axios from 'axios';

const TextToSpeechPage = () => {
  const [text, setText] = useState('');
  const [selectedModel, setSelectedModel] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [audioUrl, setAudioUrl] = useState(null);
  const [speed, setSpeed] = useState(1.0);
  const [pitch, setPitch] = useState(1.0);
  const [availableVoices, setAvailableVoices] = useState([]);
  const [languages, setLanguages] = useState([]);
  const [selectedLanguage, setSelectedLanguage] = useState('en');
  const audioRef = useRef(null);
  const fileInputRef = useRef(null);

  useEffect(() => {
    fetchLanguages();
  }, []);

  useEffect(() => {
    if (selectedLanguage) {
      fetchVoices(selectedLanguage);
    }
  }, [selectedLanguage]);

  const fetchLanguages = async () => {
    try {
      const response = await axios.get('http://localhost:81/tts/languages');
      setLanguages(response.data.languages);
      setSelectedLanguage('en');
    } catch (err) {
      setError('Cannot load languages');
    }
  };

  const fetchVoices = async (language) => {
    try {
      const response = await axios.get(`http://localhost:81/tts/voices/${language}`);
      setAvailableVoices(response.data.voices);
      if (response.data.voices.length > 0) {
        setSelectedModel(response.data.voices[0].model);
      }
    } catch (err) {
      setError('Cannot load voices');
    }
  };

  const handleFileUpload = async (e) => {
    const file = e.target.files[0];
    if (file && file.type === 'text/plain') {
      const reader = new FileReader();
      reader.onload = (e) => setText(e.target.result);
      reader.readAsText(file);
    } else {
      setError('Please select a .txt file');
    }
  };

  const handleConvert = async () => {
    if (!text.trim()) {
      setError('Please enter text to convert');
      return;
    }

    setLoading(true);
    try {
      const response = await axios.post('http://localhost:81/tts', {
        text,
        speed: 1.0, // Fixed value
        pitch: 1.0, // Fixed value
        model: selectedModel
      });
      
      setAudioUrl(`http://localhost:81${response.data.audio_url}`);
      setError(null);
    } catch (err) {
      setError('Conversion error: ' + (err.response?.data?.message || err.message));
    } finally {
      setLoading(false);
    }
  };

  const handleDownload = async () => {
    try {
      const response = await fetch(audioUrl);
      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = url;
      link.download = 'audio.wav';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    } catch (err) {
      setError('Download failed');
    }
  };

  return (
    <DashboardLayout title="Text to Speech">
      <div className="bg-white rounded-lg shadow p-6">
        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Language
          </label>
          <select
            value={selectedLanguage}
            onChange={(e) => setSelectedLanguage(e.target.value)}
            className="w-full p-2 border rounded-lg"
          >
            {languages.map((lang) => (
              <option key={lang} value={lang}>
                {lang.toUpperCase()}
              </option>
            ))}
          </select>
        </div>

        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Voice Model
          </label>
          <select
            value={selectedModel}
            onChange={(e) => setSelectedModel(e.target.value)}
            className="w-full p-2 border rounded-lg"
          >
            {availableVoices.map((voice) => (
              <option key={voice.id} value={voice.model}>
                {voice.name}
              </option>
            ))}
          </select>
        </div>

        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Speed ({speed}x)
          </label>
          <input
            type="range"
            min="0.5"
            max="2.0"
            step="0.1"
            value={speed}
            onChange={(e) => setSpeed(parseFloat(e.target.value))}
            className="w-full"
          />
        </div>

        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Pitch ({pitch})
          </label>
          <input
            type="range"
            min="0.5"
            max="2.0"
            step="0.1"
            value={pitch}
            onChange={(e) => setPitch(parseFloat(e.target.value))}
            className="w-full"
          />
        </div>

        <div className="mb-6">
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Text
          </label>
          <textarea
            value={text}
            onChange={(e) => setText(e.target.value)}
            className="w-full h-40 p-3 border rounded-lg resize-none focus:ring-blue-500 focus:border-blue-500"
            placeholder="Enter text to convert..."
          />
          <div className="mt-2 flex justify-between text-sm text-gray-500">
            <span>{text.length} characters</span>
            <button onClick={() => setText('')} className="text-gray-600 hover:text-gray-900">
              Clear
            </button>
          </div>
        </div>

        <div className="flex justify-between items-center">
          <button 
            onClick={() => fileInputRef.current?.click()}
            className="flex items-center px-4 py-2 border rounded-lg hover:bg-gray-50"
          >
            <Upload className="w-4 h-4 mr-2" />
            Upload Text File
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
            className={`flex items-center px-6 py-2 rounded-lg text-white
              ${loading || !text.trim() ? 'bg-blue-400 cursor-not-allowed' : 'bg-blue-600 hover:bg-blue-700'}`}
          >
            {loading ? (
              <div className="flex items-center">
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                Processing...
              </div>
            ) : (
              <>
                <Volume2 className="w-4 h-4 mr-2" />
                Convert
              </>
            )}
          </button>
        </div>

        {error && (
          <div className="mt-4 p-3 bg-red-50 text-red-700 rounded-lg">
            {error}
          </div>
        )}

        {audioUrl && (
          <div className="mt-6 p-4 border rounded-lg">
            <audio 
              ref={audioRef}
              controls 
              className="w-full mb-4"
              src={audioUrl}
            />
            <button 
              onClick={handleDownload}
              className="flex items-center text-blue-600 hover:text-blue-800"
            >
              <Download className="w-4 h-4 mr-2" />
              Download Audio
            </button>
          </div>
        )}
      </div>
    </DashboardLayout>
  );
};

export default TextToSpeechPage;