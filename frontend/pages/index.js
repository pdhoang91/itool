// frontend/pages/index.js
import React from 'react';
import { 
  Type, 
  Mic, 
  ImageIcon, 
  Camera, 
  Languages,
  Volume2,
  ArrowRight
} from 'lucide-react';

const FeatureCard = ({ icon: Icon, title, description, path }) => (
  <a 
    href={path}
    className="group relative p-6 bg-white rounded-xl shadow-sm hover:shadow-md transition-all duration-200 border border-gray-100"
  >
    <div className="flex items-start gap-4">
      <span className="rounded-lg bg-blue-50 p-3 text-blue-600">
        <Icon className="w-6 h-6" />
      </span>
      <div className="flex-1">
        <h3 className="font-semibold text-lg text-gray-900 group-hover:text-blue-600 transition-colors">
          {title}
        </h3>
        <p className="mt-2 text-gray-500 text-sm">
          {description}
        </p>
      </div>
      <ArrowRight className="w-5 h-5 text-gray-400 group-hover:text-blue-500 transition-colors" />
    </div>
  </a>
);

const HomePage = () => {
  const features = [
    {
      title: "Text to Speech",
      description: "Chuyển đổi văn bản thành giọng nói tự nhiên với nhiều giọng đọc khác nhau",
      icon: Volume2,
      path: "/tools/text-to-speech"
    },
    {
      title: "OCR - Nhận Dạng Văn Bản",
      description: "Trích xuất văn bản từ hình ảnh với độ chính xác cao",
      icon: Camera,
      path: "/tools/ocr"
    },
    {
      title: "Dịch Văn Bản",
      description: "Dịch văn bản qua lại giữa nhiều ngôn ngữ khác nhau",
      icon: Languages,
      path: "/tools/translation"
    },
    {
      title: "Nhận Dạng Giọng Nói",
      description: "Chuyển đổi giọng nói thành văn bản với nhiều ngôn ngữ",
      icon: Mic,
      path: "/tools/speech-recognition"
    },
    {
      title: "Xóa Nền Ảnh",
      description: "Tự động xóa nền ảnh với công nghệ AI tiên tiến",
      icon: ImageIcon,
      path: "/tools/background-remove"
    },
    {
      title: "Nhận Diện Khuôn Mặt",
      description: "Phát hiện và phân tích các đặc điểm khuôn mặt trong ảnh",
      icon: Camera,
      path: "/tools/face-recognition"
    }
  ];

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Hero Section */}
      <div className="bg-white border-b">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24">
          <div className="text-center">
            <h1 className="text-4xl font-bold text-gray-900 sm:text-5xl">
              AI Tools Dashboard
            </h1>
            <p className="mt-4 text-xl text-gray-600 max-w-2xl mx-auto">
              Bộ công cụ AI đa năng giúp bạn xử lý văn bản, âm thanh và hình ảnh một cách thông minh
            </p>
          </div>
        </div>
      </div>

      {/* Features Grid */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {features.map((feature, index) => (
            <FeatureCard key={index} {...feature} />
          ))}
        </div>
      </div>

      {/* Footer */}
      <footer className="bg-white border-t mt-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="text-center text-gray-500 text-sm">
            © 2024 AI Tools Dashboard. All rights reserved.
          </div>
        </div>
      </footer>
    </div>
  );
};

export default HomePage;
