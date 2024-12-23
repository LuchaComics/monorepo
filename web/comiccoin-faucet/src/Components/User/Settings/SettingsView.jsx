import React, { useState } from 'react';
import {
  Bell,
  Mail,
  Shield,
  Smartphone,
  Moon,
  Globe,
  ToggleLeft,
  ToggleRight,
  ChevronRight,
  Signature
} from 'lucide-react';
import { useRecoilState } from "recoil";
import { Navigate, Link } from "react-router-dom";


import Topbar from "../../../Components/Navigation/Topbar";
import { currentUserState } from "../../../AppState";


const SettingsPage = () => {
  // Variable controls the global state of the app.
  const [currentUser] = useRecoilState(currentUserState);

  const [notifications, setNotifications] = useState(true);
  const [darkMode, setDarkMode] = useState(false);
  const [twoFactor, setTwoFactor] = useState(false);

  const SettingToggle = ({ enabled, onToggle, icon: Icon, title, description }) => (
    <div className="flex items-center justify-between p-4 hover:bg-purple-50 rounded-lg transition-colors">
      <div className="flex items-start space-x-4">
        <div className="p-2 bg-purple-100 rounded-lg">
          <Icon className="w-5 h-5 text-purple-600" />
        </div>
        <div>
          <h3 className="font-medium text-gray-900">{title}</h3>
          <p className="text-sm text-gray-500">{description}</p>
        </div>
      </div>
      <button
        onClick={onToggle}
        className="text-purple-600"
      >
        {enabled ? (
          <ToggleRight className="w-6 h-6" />
        ) : (
          <ToggleLeft className="w-6 h-6" />
        )}
      </button>
    </div>
  );

  const SettingLink = ({ icon: Icon, title, description, to }) => (
    <Link className="flex items-center justify-between p-4 hover:bg-purple-50 rounded-lg transition-colors cursor-pointer" to={to}>
      <div className="flex items-start space-x-4">
        <div className="p-2 bg-purple-100 rounded-lg">
          <Icon className="w-5 h-5 text-purple-600" />
        </div>
        <div>
          <h3 className="font-medium text-gray-900">{title}</h3>
          <p className="text-sm text-gray-500">{description}</p>
        </div>
      </div>
      <ChevronRight className="w-5 h-5 text-gray-400" />
    </Link>
  );

  return (
    <div className="min-h-screen bg-purple-50">

      <Topbar currentPage="Settings" />

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">

        <h1 className="text-3xl font-bold text-purple-800 mb-8" style={{fontFamily: 'Comic Sans MS, cursive'}}>Settings</h1>

        {/* Account Settings */}
        <div className="bg-white rounded-xl shadow-lg border-2 border-purple-200">
          <div className="p-6 border-b border-purple-100">
            <h2 className="text-xl text-purple-800" style={{fontFamily: 'Comic Sans MS, cursive'}}>
              Account Settings
            </h2>
          </div>
          <div className="p-4 space-y-2">
            <SettingLink
              to="/settings/info"
              icon={Mail}
              title="Email Settings"
              description="Update your email address and communication preferences"
            />
            <SettingLink
              to="/settings/pass"
              icon={Shield}
              title="Password & Security"
              description="Manage your password and security settings"
            />
            {currentUser.profileVerificationStatus === 1 && <>
                <SettingLink
                  to="/apply-for-verification"
                  icon={Signature}
                  title="Apply for Verification"
                  description="Get increased daily submission limits and extra benefits"
                />
            </>}
            {currentUser.profileVerificationStatus === 2 && <>
                <div className="p-4">
                  <div className="p-4 bg-yellow-50 rounded-lg border border-yellow-200">
                    <div className="flex items-center space-x-2">
                      <div className="w-2 h-2 bg-yellow-500 rounded-full"></div>
                      <p className="text-yellow-700 font-medium">Pending Verification Review</p>
                    </div>
                    <p className="text-sm text-yellow-600 mt-1">Submitted Profile for verification, please wait 1 week</p>
                  </div>
                </div>
            </>}
            {currentUser.profileVerificationStatus === 3 && <>
                <div className="p-4">
                  <div className="p-4 bg-green-50 rounded-lg border border-green-200">
                    <div className="flex items-center space-x-2">
                      <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                      <p className="text-green-700 font-medium">Verified Profile</p>
                    </div>
                    <p className="text-sm text-green-600 mt-1">Your profile has been verified successfully</p>
                  </div>
                </div>
            </>}
            {currentUser.profileVerificationStatus === 4 && <>
                <div className="p-4">
                  <div className="p-4 bg-red-50 rounded-lg border border-red-200">
                    <div className="flex items-center space-x-2">
                      <div className="w-2 h-2 bg-red-500 rounded-full"></div>
                      <p className="text-red-700 font-medium">Profile Rejected</p>
                    </div>
                    <p className="text-sm text-red-600 mt-1">Your profile has been rejected</p>
                  </div>
                </div>
            </>}

            {/*
            <SettingToggle
              icon={Smartphone}
              title="Two-Factor Authentication"
              description="Add an extra layer of security to your account"
              enabled={twoFactor}
              onToggle={() => setTwoFactor(!twoFactor)}
            />
            */}
          </div>
        </div>

        {/* Notification Settings */}
        {/*
        <div className="bg-white rounded-xl shadow-lg border-2 border-purple-200">
          <div className="p-6 border-b border-purple-100">
            <h2 className="text-xl text-purple-800" style={{fontFamily: 'Comic Sans MS, cursive'}}>
              Notifications
            </h2>
          </div>
          <div className="p-4 space-y-2">
            <SettingToggle
              icon={Bell}
              title="Push Notifications"
              description="Receive notifications about your comic submissions and ComicCoins"
              enabled={notifications}
              onToggle={() => setNotifications(!notifications)}
            />
          </div>
        </div>
        */}

        {/* Appearance Settings */}
        {/*
        <div className="bg-white rounded-xl shadow-lg border-2 border-purple-200">
          <div className="p-6 border-b border-purple-100">
            <h2 className="text-xl text-purple-800" style={{fontFamily: 'Comic Sans MS, cursive'}}>
              Appearance
            </h2>
          </div>
          <div className="p-4 space-y-2">
            <SettingToggle
              icon={Moon}
              title="Dark Mode"
              description="Switch between light and dark theme"
              enabled={darkMode}
              onToggle={() => setDarkMode(!darkMode)}
            />
            <SettingLink
              icon={Globe}
              title="Language"
              description="Choose your preferred language"
            />
          </div>
        </div>
        */}

        {/* Connection Status */}
        <div className="bg-white rounded-xl shadow-lg border-2 border-purple-200">
          <div className="p-6 border-b border-purple-100">
            <h2 className="text-xl text-purple-800" style={{fontFamily: 'Comic Sans MS, cursive'}}>
              Wallet Connection
            </h2>
          </div>
          <div className="p-4">
            <div className="p-4 bg-green-50 rounded-lg border border-green-200">
              <div className="flex items-center space-x-2">
                <div className="w-2 h-2 bg-green-500 rounded-full"></div>
                <p className="text-green-700 font-medium">Connected to Wallet</p>
              </div>
              <p className="text-sm text-green-600 mt-1">{currentUser.walletAddress}</p>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
};

export default SettingsPage;
