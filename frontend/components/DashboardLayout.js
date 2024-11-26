import React from 'react';
import Navbar from './Navbar';

const DashboardLayout = ({ children, title }) => {
  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar />
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div className="md:flex md:items-center md:justify-between mb-6">
          <div className="flex-1 min-w-0">
            <h1 className="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl">
              {title}
            </h1>
          </div>
        </div>
        {children}
      </div>
    </div>
  );
};

export default DashboardLayout;