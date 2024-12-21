import React, { useState } from "react";
import { Coins, AlertCircle, LogOut } from "lucide-react";
import { Navigate, Link } from "react-router-dom";
import { useRecoilState } from "recoil";

import { putProfileWalletAddressAPI } from "../../API/Profile";
import { currentUserState } from "../../AppState";

export default function UserAddWalletToFaucet() {
  // Variable controls the global state of the app.
  const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

  // Variable controls the login form.
  const [walletAddress, setWalletAddress] = useState("");
  const [errors, setErrors] = useState({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [forceURL, setForceURL] = useState("");

  const handleSubmit = (e) => {
    e.preventDefault();

    // Store in persistence storage in the browser by simply updating
    // the existing logged in user.
    const updatedCurrentUser = {...currentUser};
    updatedCurrentUser.walletAddress = walletAddress
    setCurrentUser(updatedCurrentUser);
    console.log("putProfileWalletAddressAPI: currentUser:", currentUser);

    setIsSubmitting(true);

    const formData = {
      wallet_address: walletAddress,
    };

    putProfileWalletAddressAPI(
      formData,
      (resp) => {
        // SUCCESS HANDLER
        // ---------------

        // For debugging purposes only.
        console.log("putProfileWalletAddressAPI: Starting...");
        console.log(resp);

        setForceURL("/added-my-wallet-to-faucet-successfully");

        // -------------------
        // end SUCCESS HANDLER
      },
      (apiErr) => {
        console.log("putProfileWalletAddressAPI: apiErr:", apiErr);
        if (JSON.stringify(apiErr).includes("Wallet address already set")) {
          setForceURL("/added-my-wallet-to-faucet-successfully");
          return;
        }
        setErrors(apiErr);
      },
      () => {
        console.log("putProfileWalletAddressAPI: Starting...");
        setIsSubmitting(false);
      },
      () => {
        // console.log("putProfileWalletAddressAPI: unauthorized...");
        // window.location.href = "/login?unauthorized=true";
      },
    );
  };

  if (forceURL !== "") {
    console.log("UserAddWalletToFaucet: Redirecting to:", forceURL);
    return <Navigate to={forceURL} />;
  }

 console.log("UserAddWalletToFaucet: Returning");

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
          <div className="space-y-8">
            {/* Welcome Section */}
            <div className="text-center space-y-4">
              <h1 className="text-4xl font-bold text-purple-800">
                Welcome to ComicCoin Faucet!
              </h1>
              <p className="text-xl text-gray-700">
                You're about to receive{" "}
                <span className="font-bold">10 free ComicCoins</span> to start
                your journey!
              </p>
            </div>

            <div className="bg-white rounded-xl p-8 shadow-lg border-2 border-purple-200 space-y-8">
              {/* Key Information First */}
              <div className="space-y-6">
                {/* Important Security Notice */}
                <div className="bg-yellow-50 border-l-4 border-yellow-400 p-4 rounded-r">
                  <h3 className="font-bold text-yellow-800 mb-2">
                    🔒 Important: Secure Your Wallet
                  </h3>
                  <p className="text-yellow-800">
                    Your wallet is your key to accessing your ComicCoins. Keep
                    it safe:
                  </p>
                  <ul className="list-disc pl-5 mt-2 space-y-1 text-yellow-800">
                    <li>
                      Back up your wallet credentials in multiple secure
                      locations
                    </li>
                    <li>Never share your private keys with anyone</li>
                    <li>
                      Lost wallets cannot be recovered - keep your backup
                      phrases safe!
                    </li>
                  </ul>
                </div>

                {/* Benefits Section */}
                <div className="bg-purple-50 p-4 rounded-lg">
                  <h3 className="font-bold text-purple-800 mb-2">
                    🚀 With ComicCoin, you'll be able to:
                  </h3>
                  <ul className="grid grid-cols-1 md:grid-cols-2 gap-3 mt-2">
                    <li className="flex items-center space-x-2">
                      <span className="text-purple-600">•</span>
                      <span>Trade digital comics securely</span>
                    </li>
                    <li className="flex items-center space-x-2">
                      <span className="text-purple-600">•</span>
                      <span>Join exclusive comic communities</span>
                    </li>
                    <li className="flex items-center space-x-2">
                      <span className="text-purple-600">•</span>
                      <span>Earn trading rewards</span>
                    </li>
                    <li className="flex items-center space-x-2">
                      <span className="text-purple-600">•</span>
                      <span>Access limited editions</span>
                    </li>
                  </ul>
                </div>
              </div>

              {/* Error Messages */}
              {Object.keys(errors).length > 0 && (
                <div className="bg-red-50 border-l-4 border-red-500 p-4 rounded-r-lg">
                  <div className="flex">
                    <div className="flex-shrink-0">
                      <AlertCircle className="h-5 w-5 text-red-400" />
                    </div>
                    <div className="ml-3">
                      <h3 className="text-sm font-medium text-red-800">
                        Please correct the following errors:
                      </h3>
                      <div className="mt-2 text-sm text-red-700">
                        <ul className="list-disc space-y-1 pl-5">
                          {Object.values(errors).map((error, index) => (
                            <li key={index}>{error}</li>
                          ))}
                        </ul>
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {/* Wallet Setup Form */}
              <div className="border-t pt-6">
                <h2 className="text-xl font-bold text-purple-800 mb-6">
                  Set Up Your Wallet
                </h2>
                <form onSubmit={handleSubmit} className="space-y-6">
                  <div className="space-y-1">
                    <label
                      htmlFor="walletAddress"
                      className="block text-sm font-medium text-gray-700"
                    >
                      Wallet Address *
                    </label>
                    <input
                      type="text"
                      id="walletAddress"
                      value={walletAddress}
                      onChange={(e) => setWalletAddress(e.target.value)}
                      className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                        errors.walletAddress
                          ? "border-red-500 bg-red-50"
                          : "border-gray-300"
                      }`}
                      placeholder="Enter your 42-character wallet address"
                    />
                    {errors.walletAddress && (
                      <p className="mt-1 text-sm text-red-600">
                        {errors.walletAddress}
                      </p>
                    )}
                  </div>

                  <div className="flex space-x-4">
                    <button
                      type="submit"
                      disabled={isSubmitting}
                      className={`w-full px-6 py-3 bg-purple-600 hover:bg-purple-700 text-white font-bold rounded-lg transition-colors ${
                        isSubmitting ? "opacity-50 cursor-not-allowed" : ""
                      }`}
                    >
                      {isSubmitting ? "Submitting..." : "Submit Wallet Address"}
                    </button>
                  </div>
                </form>
              </div>
            </div>
          </div>
        </div>
      </main>

      <footer className="bg-gradient-to-r from-purple-700 to-indigo-800 text-white py-8 mt-8">
        <div className="max-w-7xl mx-auto px-4 text-center">
          <p className="mb-4">© 2024 ComicCoin Faucet. All rights reserved.</p>
          <p>
            <Link to="/terms" className="underline hover:text-purple-200">
              Terms of Service
            </Link>
            {" | "}
            <Link to="/privacy" className="underline hover:text-purple-200">
              Privacy Policy
            </Link>
          </p>
        </div>
      </footer>
    </div>
  );
}
