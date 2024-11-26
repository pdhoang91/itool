// frontend/pages/dashboard.js
import React from 'react';
import { 
  Home, Upload, BarChart, Settings, 
  Type, Mic, Image, Camera, Languages
} from 'lucide-react';

export default function DashboardLayout() {
  const tools = [
    {
      category: "Text Processing",
      icon: <Type size={20} />,
      tools: [
        { name: "Text to Speech", path: "/tools/text-to-speech" },
        { name: "OCR", path: "/tools/ocr" },
        { name: "Translation", path: "/tools/translation" }
      ]
    },
    {
      category: "Audio Processing", 
      icon: <Mic size={20} />,
      tools: [
        { name: "Speech Recognition", path: "/tools/speech-recognition" },
        { name: "Voice to Text", path: "/tools/voice-to-text" }
      ]
    },
    {
      category: "Image Processing",
      icon: <Image size={20} />,
      tools: [
        { name: "Background Remove", path: "/tools/background-remove" },
        { name: "Face Recognition", path: "/tools/face-recognition" }
      ]
    }
  ];

  return (
    <div className="flex h-screen bg-gray-50">
      {/* Sidebar */}
      <div className="hidden md:flex md:flex-shrink-0">
        <div className="flex flex-col w-64">
          <div className="flex flex-col flex-grow pt-5 overflow-y-auto bg-white border-r">
            <div className="flex items-center flex-shrink-0 px-4">
              <img
                className="w-auto h-8"
                src="/api/placeholder/32/32"
                alt="AI Tools"
              />
              <span className="ml-2 text-xl font-semibold text-gray-900">AI Tools</span>
            </div>
            
            <div className="mt-6">
              {tools.map((category, idx) => (
                <div key={idx} className="px-3 mb-6">
                  <h3 className="flex items-center px-3 text-xs font-semibold text-gray-600 uppercase tracking-wider">
                    {category.icon}
                    <span className="ml-2">{category.category}</span>
                  </h3>
                  <div className="mt-3 space-y-1">
                    {category.tools.map((tool, toolIdx) => (
                      <a
                        key={toolIdx}
                        href={tool.path}
                        className="flex items-center px-3 py-2 text-sm font-medium text-gray-600 rounded-md hover:bg-gray-50 hover:text-gray-900"
                      >
                        {tool.name}
                      </a>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex flex-col flex-1 overflow-hidden">
        {/* Header */}
        <header className="bg-white shadow-sm">
          <div className="flex items-center justify-between h-16 px-4">
            <div className="flex items-center">
              <button className="p-1 text-gray-400 rounded-full hover:bg-gray-100 md:hidden">
                <Settings size={24} />
              </button>
              <h1 className="ml-3 text-2xl font-semibold text-gray-900">Dashboard</h1>
            </div>
            <div className="flex items-center space-x-4">
              <button className="p-1 text-gray-400 rounded-full hover:bg-gray-100">
                <Settings size={24} />
              </button>
            </div>
          </div>
        </header>

        {/* Main Content Area */}
        <main className="flex-1 relative overflow-y-auto focus:outline-none">
          <div className="py-6">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 md:px-8">
              {/* Content placeholder */}
              <div className="bg-white rounded-lg shadow p-6">
                <h2 className="text-lg font-medium mb-4">Welcome to AI Tools Dashboard</h2>
                <p className="text-gray-600">Select a tool from the sidebar to get started</p>
              </div>
            </div>
          </div>
        </main>
      </div>
    </div>
  );
}