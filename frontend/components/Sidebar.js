import React, { useState } from 'react';
import { 
  Type, 
  Mic, 
  Image, 
  Camera, 
  Languages,
  Volume2
} from 'lucide-react';

const Sidebar = () => {
  // Simulate current path - in real app this would come from your routing system
  const [currentPath, setCurrentPath] = useState('/tools/text-to-speech');

  const tools = [
    {
      category: "Xử Lý Văn Bản",
      icon: <Type size={20} />,
      tools: [
        { name: "Text to Speech", path: "/tools/text-to-speech", icon: <Volume2 size={16} /> },
        { name: "OCR", path: "/tools/ocr", icon: <Camera size={16} /> },
        { name: "Translation", path: "/tools/translation", icon: <Languages size={16} /> }
      ]
    },
    {
      category: "Xử Lý Âm Thanh", 
      icon: <Mic size={20} />,
      tools: [
        { name: "Speech Recognition", path: "/tools/speech-recognition", icon: <Mic size={16} /> },
        { name: "Voice to Text", path: "/tools/voice-to-text", icon: <Mic size={16} /> }
      ]
    },
    {
      category: "Xử Lý Hình Ảnh",
      icon: <Image size={20} />,
      tools: [
        { name: "Background Remove", path: "/tools/background-remove", icon: <Image size={16} /> },
        { name: "Face Recognition", path: "/tools/face-recognition", icon: <Camera size={16} /> }
      ]
    }
  ];

  const handleNavigate = (path) => {
    setCurrentPath(path);
    // In a real app, you would handle navigation here
    // history.push(path) or similar
  };

  return (
    <div className="w-64 bg-white h-screen flex flex-col border-r">
      {/* Logo */}
      <div className="flex items-center px-6 py-4">
        <img 
          src="/api/placeholder/32/32" 
          alt="AI Tools"
          className="w-8 h-8"
        />
        <span className="ml-2 text-lg font-semibold text-gray-900">AI Tools</span>
      </div>

      {/* Navigation */}
      <nav className="flex-1 overflow-y-auto px-4 py-2">
        {tools.map((category, idx) => (
          <div key={idx} className="mb-6">
            {/* Category Header */}
            <h4 className="flex items-center px-2 mb-2 text-xs font-semibold text-gray-600 uppercase tracking-wider">
              <span className="text-gray-500 mr-2">{category.icon}</span>
              {category.category}
            </h4>

            {/* Tools List */}
            <ul className="space-y-1">
              {category.tools.map((tool, toolIdx) => (
                <li key={toolIdx}>
                  <a 
                    href={tool.path}
                    onClick={(e) => {
                      e.preventDefault();
                      handleNavigate(tool.path);
                    }}
                    className={`
                      flex items-center px-2 py-2 text-sm rounded-md cursor-pointer
                      ${currentPath === tool.path 
                        ? 'bg-blue-50 text-blue-700 font-medium' 
                        : 'text-gray-600 hover:bg-gray-50 hover:text-gray-900'
                      }
                    `}
                  >
                    <span className="mr-2 text-gray-500">{tool.icon}</span>
                    {tool.name}
                  </a>
                </li>
              ))}
            </ul>
          </div>
        ))}
      </nav>

      {/* Footer */}
      <div className="p-4 border-t">
        <div className="text-xs text-gray-500 text-center">
          Version 1.0.0
        </div>
      </div>
    </div>
  );
};

export default Sidebar;