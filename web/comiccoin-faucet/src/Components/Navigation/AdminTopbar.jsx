import React, { useState } from 'react';
import {
  Coins, Home, Image, History, Wallet,
  Settings, HelpCircle, LogOut, Clock, CheckCircle, XCircle,
  Menu, X, Upload, ArrowRight, Sparkles
} from 'lucide-react';
import { Link, Navigate } from "react-router-dom";

const AdminTopbar = (props) => {
  const { currentPage } = props;
  const [isNavOpen, setIsNavOpen] = useState(false);
  const [forceURL, setForceURL] = useState("");

  const navigation = [
    { name: 'Dashboard', icon: Home, current: currentPage === "Dashboard", url: '/admin/dashboard' },
    { name: 'Settings', icon: Settings, current: currentPage === "Settings", url: '/admin/settings' },
  ];

  return (
    <nav className="bg-gradient-to-r from-purple-700 to-indigo-800 text-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 xl:px-12 2xl:px-24">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <div className="flex items-center space-x-2">
            <Coins className="h-8 w-8" />
            <span className="text-xl font-bold" style={{fontFamily: 'Comic Sans MS, cursive'}}>
              ComicCoin Admin
            </span>
          </div>

          {/* Mobile menu button */}
          <div className="flex items-center lg:hidden">
            <button
              onClick={() => setIsNavOpen(!isNavOpen)}
              className="inline-flex items-center justify-center p-2 rounded-md text-white hover:bg-purple-600 focus:outline-none"
            >
              {isNavOpen ? (
                <X className="h-6 w-6" />
              ) : (
                <Menu className="h-6 w-6" />
              )}
            </button>
          </div>

          {/* Desktop navigation and logout combined */}
          <div className="hidden lg:flex lg:items-center lg:space-x-4">
            <div className="flex items-center space-x-4">
              {navigation.map((item) => (
                <Link
                  key={item.name}
                  to={item.url}
                  className={`flex items-center space-x-1 px-3 py-2 rounded-md text-sm font-medium ${
                    item.current
                      ? 'bg-purple-600 bg-opacity-50'
                      : 'hover:bg-purple-600 hover:bg-opacity-25'
                  }`}
                >
                  <item.icon className="h-4 w-4" />
                  <span>{item.name}</span>
                </Link>
              ))}
              <Link
                to="/logout"
                className="flex items-center space-x-1 px-3 py-2 rounded-md hover:bg-purple-600 hover:bg-opacity-25 text-purple-200 hover:text-white"
              >
                <LogOut className="h-4 w-4" />
                <span>Logout</span>
              </Link>
            </div>
          </div>
        </div>
      </div>

      {/* Mobile menu */}
      <div className={`lg:hidden ${isNavOpen ? 'block' : 'hidden'}`}>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 xl:px-12 2xl:px-24 pt-2 pb-3 space-y-1">
          {navigation.map((item) => (
            <Link
              key={item.name}
              to={item.url}
              className={`flex items-center space-x-2 px-3 py-2 rounded-md text-base font-medium ${
                item.current
                  ? 'bg-purple-600 bg-opacity-50'
                  : 'hover:bg-purple-600 hover:bg-opacity-25'
              }`}
            >
              <item.icon className="h-5 w-5" />
              <span>{item.name}</span>
            </Link>
          ))}
          <Link
            to="/logout"
            className="w-full flex items-center space-x-2 px-3 py-2 rounded-md text-base font-medium text-purple-200 hover:text-white hover:bg-purple-600 hover:bg-opacity-25"
          >
            <LogOut className="h-5 w-5" />
            <span>Logout</span>
          </Link>
        </div>
      </div>
    </nav>
  );
};

export default AdminTopbar;
