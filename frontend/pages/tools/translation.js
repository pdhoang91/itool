import React, { useState } from 'react';
import { Languages, ArrowRight, Copy, ArrowLeftRight } from 'lucide-react';
import DashboardLayout from '../../components/DashboardLayout';


const TranslationPage = () => {
  const [sourceText, setSourceText] = useState('');
  const [translatedText, setTranslatedText] = useState('');
  const [sourceLang, setSourceLang] = useState('vi');
  const [targetLang, setTargetLang] = useState('en');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const languages = [
    { code: 'vi', name: 'Tiếng Việt' },
    { code: 'en', name: 'English' },
    { code: 'ja', name: '日本語' },
    { code: 'ko', name: '한국어' },
    { code: 'zh', name: '中文' },
  ];

  const handleTranslate = async () => {
    if (!sourceText.trim()) {
      setError('Vui lòng nhập văn bản cần dịch.');
      return;
    }

    setLoading(true);
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 1500));
      setTranslatedText('Translated text will appear here...');
      setError(null);
    } catch (err) {
      setError('Có lỗi xảy ra khi dịch văn bản.');
    } finally {
      setLoading(false);
    }
  };

  const swapLanguages = () => {
    const tempLang = sourceLang;
    setSourceLang(targetLang);
    setTargetLang(tempLang);
    setSourceText(translatedText);
    setTranslatedText(sourceText);
  };

  return (
    <DashboardLayout title="Dịch Văn Bản">
      <div className="bg-white rounded-lg shadow">
        <div className="p-6">
          {/* Language Selection */}
          <div className="flex items-center space-x-4 mb-6">
            <select
              value={sourceLang}
              onChange={(e) => setSourceLang(e.target.value)}
              className="block w-40 rounded-md border-gray-300 shadow-sm"
            >
              {languages.map(lang => (
                <option key={lang.code} value={lang.code}>
                  {lang.name}
                </option>
              ))}
            </select>

            <button 
              onClick={swapLanguages}
              className="p-2 rounded-full hover:bg-gray-100"
            >
              <ArrowLeftRight className="w-5 h-5 text-gray-500" />
            </button>

            <select
              value={targetLang}
              onChange={(e) => setTargetLang(e.target.value)}
              className="block w-40 rounded-md border-gray-300 shadow-sm"
            >
              {languages.map(lang => (
                <option key={lang.code} value={lang.code}>
                  {lang.name}
                </option>
              ))}
            </select>
          </div>

          {/* Translation Area */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {/* Source Text */}
            <div>
              <textarea
                value={sourceText}
                onChange={(e) => setSourceText(e.target.value)}
                placeholder="Nhập văn bản cần dịch..."
                className="w-full h-40 p-3 border rounded-lg resize-none"
              />
              <div className="mt-2 flex justify-between text-sm text-gray-500">
                <span>{sourceText.length} ký tự</span>
                <button 
                  onClick={() => setSourceText('')}
                  className="text-gray-600 hover:text-gray-900"
                >
                  Xóa
                </button>
              </div>
            </div>

            {/* Target Text */}
            <div>
              <div className="relative">
                <textarea
                  value={translatedText}
                  readOnly
                  placeholder="Bản dịch sẽ xuất hiện ở đây..."
                  className="w-full h-40 p-3 border rounded-lg resize-none bg-gray-50"
                />
                {translatedText && (
                  <button 
                    onClick={() => navigator.clipboard.writeText(translatedText)}
                    className="absolute top-2 right-2 p-2 text-gray-500 hover:text-gray-700 rounded-full hover:bg-gray-100"
                  >
                    <Copy className="w-4 h-4" />
                  </button>
                )}
              </div>
            </div>
          </div>

          {/* Action Button */}
          <div className="mt-6 flex justify-end">
            <button
              onClick={handleTranslate}
              disabled={loading || !sourceText}
              className={`
                flex items-center px-6 py-2 rounded-lg text-white
                ${loading || !sourceText ? 'bg-blue-400' : 'bg-blue-600 hover:bg-blue-700'}
              `}
            >
              <Languages className="w-4 h-4 mr-2" />
              {loading ? 'Đang dịch...' : 'Dịch văn bản'}
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

export default TranslationPage;