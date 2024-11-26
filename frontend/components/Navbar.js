import React, { useState } from 'react';
import { useRouter } from 'next/router';
import Link from 'next/link';
import { 
  Menu,
  X,
  Volume2, 
  Camera,
  Languages,
  Mic,
  ImageIcon,
  User,
  LogOut,
  Settings
} from 'lucide-react';

const Navbar = () => {
  const router = useRouter();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  const navigation = [
    { name: 'Text to Speech', href: '/tools/text-to-speech', icon: Volume2 },
    { name: 'OCR - Nhận Dạng Văn Bản', href: '/tools/ocr', icon: Camera },
    { name: 'Dịch Văn Bản', href: '/tools/translation', icon: Languages },
    { name: 'Nhận Dạng Giọng Nói', href: '/tools/speech-recognition', icon: Mic },
    { name: 'Xóa Nền Ảnh', href: '/tools/background-remove', icon: ImageIcon },
  ];

  const isActive = (path) => router.pathname === path;

  return (
    <nav className="bg-white border-b">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between h-16">
          {/* Logo & Desktop Navigation */}
          <div className="flex">
            {/* Logo */}
            <div className="flex-shrink-0 flex items-center">
              <Link href="/" className="text-xl font-bold text-blue-600">
                AI Tools
              </Link>
            </div>

            {/* Desktop Navigation */}
            <div className="hidden md:ml-6 md:flex md:space-x-4">
              {navigation.map((item) => (
                <Link
                  key={item.name}
                  href={item.href}
                  className={`
                    inline-flex items-center px-3 py-2 text-sm font-medium rounded-md
                    ${isActive(item.href)
                      ? 'text-blue-600 bg-blue-50'
                      : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
                    }
                  `}
                >
                  <item.icon className="h-4 w-4 mr-2" />
                  {item.name}
                </Link>
              ))}
            </div>
          </div>

          {/* User Menu & Mobile Menu Button */}
          <div className="flex items-center">
            {/* User Menu */}
            <div className="hidden md:ml-4 md:flex md:items-center">
              <button className="p-1 rounded-full text-gray-600 hover:text-gray-900">
                <Settings className="h-6 w-6" />
              </button>
              
              <div className="ml-4 relative flex items-center">
                <button className="flex text-sm rounded-full focus:outline-none">
                  <User className="h-8 w-8 p-1 rounded-full bg-gray-100 text-gray-600" />
                </button>
                <button className="ml-4 p-1 rounded-full text-gray-600 hover:text-gray-900">
                  <LogOut className="h-6 w-6" />
                </button>
              </div>
            </div>

            {/* Mobile menu button */}
            <div className="flex md:hidden">
              <button
                onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
                className="inline-flex items-center justify-center p-2 rounded-md text-gray-600 hover:text-gray-900 hover:bg-gray-100"
              >
                {isMobileMenuOpen ? (
                  <X className="block h-6 w-6" />
                ) : (
                  <Menu className="block h-6 w-6" />
                )}
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Mobile menu */}
      {isMobileMenuOpen && (
        <div className="md:hidden">
          <div className="px-2 pt-2 pb-3 space-y-1">
            {navigation.map((item) => (
              <Link
                key={item.name}
                href={item.href}
                className={`
                  block px-3 py-2 rounded-md text-base font-medium
                  ${isActive(item.href)
                    ? 'text-blue-600 bg-blue-50'
                    : 'text-gray-600 hover:text-gray-900 hover:bg-gray-50'
                  }
                `}
                onClick={() => setIsMobileMenuOpen(false)}
              >
                <div className="flex items-center">
                  <item.icon className="h-5 w-5 mr-3" />
                  {item.name}
                </div>
              </Link>
            ))}
            
            {/* Mobile User Menu Items */}
            <div className="border-t border-gray-200 pt-4 pb-3">
              <div className="flex items-center px-3">
                <User className="h-8 w-8 p-1 rounded-full bg-gray-100 text-gray-600" />
                <div className="ml-3">
                  <div className="text-base font-medium text-gray-800">User Name</div>
                  <div className="text-sm font-medium text-gray-500">user@example.com</div>
                </div>
              </div>
              <div className="mt-3 space-y-1">
                <Link
                  href="/settings"
                  className="block px-3 py-2 rounded-md text-base font-medium text-gray-600 hover:text-gray-900 hover:bg-gray-50"
                  onClick={() => setIsMobileMenuOpen(false)}
                >
                  <div className="flex items-center">
                    <Settings className="h-5 w-5 mr-3" />
                    Cài đặt
                  </div>
                </Link>
                <button
                  className="w-full text-left block px-3 py-2 rounded-md text-base font-medium text-gray-600 hover:text-gray-900 hover:bg-gray-50"
                >
                  <div className="flex items-center">
                    <LogOut className="h-5 w-5 mr-3" />
                    Đăng xuất
                  </div>
                </button>
              </div>
            </div>
          </div>
        </div>
      )}
    </nav>
  );
};

export default Navbar;