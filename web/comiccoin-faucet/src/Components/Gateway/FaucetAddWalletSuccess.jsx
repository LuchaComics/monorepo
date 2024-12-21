import React, { useState } from "react";
import { Coins, AlertCircle, LogOut, CheckCircle2, ArrowRight } from "lucide-react";
import { Navigate, Link } from "react-router-dom";
import { useRecoilState } from "recoil";

import { currentUserState } from "../../AppState";


export default function UserAddWalletToFaucetSuccess() {
    // Variable controls the global state of the app.
    const [currentUser] = useRecoilState(currentUserState);
    console.log("UserAddWalletToFaucetSuccess: currentUser:", currentUser);

    const [forceURL, setForceURL] = useState("");

    if (forceURL !== "") {
      return <Navigate to={forceURL} />;
    }

  return (
    <div className="min-h-screen flex flex-col bg-gradient-to-b from-purple-100 to-white">
    <nav className="bg-gradient-to-r from-purple-700 to-indigo-800 text-white p-4">
      <div className="max-w-7xl mx-auto flex justify-between items-center">
        <div className="flex items-center space-x-2">
          <Coins className="h-8 w-8" />
          <span className="text-2xl font-bold">ComicCoin Faucet</span>
        </div>
        <button
          onClick={(e) => setForceURL("/logout")}
          className="flex items-center space-x-2 px-4 py-2 rounded-lg hover:bg-purple-600 transition-colors"
        >
          <LogOut className="h-5 w-5" />
          <span>Logout</span>
        </button>
      </div>
    </nav>

      <main className="flex-grow flex items-center justify-center">
        <div className="w-full max-w-2xl mx-4">
          <div className="bg-white rounded-xl p-8 shadow-lg border-2 border-purple-200">
            <div className="flex flex-col items-center space-y-6 text-center">
              {/* Success Icon */}
              <div className="rounded-full bg-green-100 p-3">
                <CheckCircle2 className="h-12 w-12 text-green-600" />
              </div>

              {/* Success Title */}
              <h1 className="text-3xl font-bold text-purple-800">
                Wallet Successfully Set!
              </h1>

              {/* Success Message */}
              <div className="space-y-4 text-gray-700">
                <p className="text-lg">
                  We've sent <span className="font-bold">10 ComicCoins</span> to your wallet!
                </p>
                <div className="bg-blue-50 border-l-4 border-blue-400 p-4 rounded-r text-left">
                  <p className="text-blue-800">
                    Please check your wallet balance. Note that transactions may take up to 5 minutes
                    to process through the network.
                  </p>
                </div>
              </div>

              {/* Action Button */}
              <button onClick={(e) => setForceURL("/dashboard")} className="mt-8 flex items-center space-x-2 px-8 py-3 bg-purple-600 hover:bg-purple-700 text-white font-bold rounded-lg transition-colors">
                <span>Go to Dashboard</span>
                <ArrowRight className="h-5 w-5" />
              </button>
            </div>
          </div>
        </div>
      </main>

      <footer className="bg-gradient-to-r from-purple-700 to-indigo-800 text-white py-8 mt-8">
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
