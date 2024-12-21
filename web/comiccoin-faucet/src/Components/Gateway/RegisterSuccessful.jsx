import React, { useState } from "react";
import { Coins, Mail, ArrowRight, ArrowLeft } from 'lucide-react';
import { Navigate, Link } from "react-router-dom";

export default function EmailConfirmationPage() {
  const [forceURL, setForceURL] = useState("");

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <div className="min-h-screen flex flex-col bg-gradient-to-b from-purple-100 to-white">
      {/* Navigation */}
      <nav className="bg-gradient-to-r from-purple-700 to-indigo-800 text-white p-4">
        <div className="max-w-7xl mx-auto flex justify-between items-center">
          <div className="flex items-center space-x-2">
            <Coins className="h-8 w-8" />
            <span className="text-2xl font-bold">ComicCoin Faucet</span>
          </div>
          <button onClick={(e)=>setForceURL("/")} className="flex items-center space-x-2 px-4 py-2 rounded-lg hover:bg-purple-600 transition-colors">
            <ArrowLeft className="h-5 w-5" />
            <span>Back to Home</span>
          </button>
        </div>
      </nav>

      {/* Main Content - Centered both vertically and horizontally */}
      <main className="flex-grow flex items-center justify-center">
        <div className="w-full max-w-2xl mx-4">
          <div className="bg-white rounded-xl p-8 shadow-lg border-2 border-purple-200">
            <div className="flex flex-col items-center space-y-6">
              {/* Header */}
              <div className="flex items-center space-x-2 text-purple-600">
                <Mail className="h-8 w-8" />
                <h1 className="text-2xl font-bold">Email Sent</h1>
              </div>

              {/* Message */}
              <div className="text-center text-gray-700 space-y-6">
                <p>
                  Thank you for registering - an <span className="font-bold">activation email</span>{" "}
                  has been sent to you. Please be sure to check your
                  social, promotions and spam folders if it does not
                  arrive within 5 minutes.
                </p>

                {/* Back to index link */}
                <div>
                  <button
                    onClick={(e)=>setForceURL("/")}
                    className="inline-flex items-center space-x-2 text-purple-600 hover:text-purple-700 font-medium"
                  >
                    <span>Back to index</span>
                    <ArrowRight className="h-4 w-4" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>

      {/* Footer */}
      <footer className="bg-gradient-to-r from-purple-700 to-indigo-800 text-white py-8">
        <div className="max-w-7xl mx-auto px-4 text-center">
          <p className="mb-4">© 2024 ComicCoin Faucet. All rights reserved.</p>
          <p>
            <Link to="/terms" className="underline hover:text-purple-200">Terms of Service</Link>
            {' | '}
            <Link to="/privacy" className="underline hover:text-purple-200">Privacy Policy</Link>
          </p>
        </div>
      </footer>
    </div>
  );
}
