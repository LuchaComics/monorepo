import {useState, useEffect} from 'react';
import { Link } from "react-router-dom";
import { WalletMinimal, Send, QrCode, MoreHorizontal, Wallet, Copy, CheckCircle2 } from 'lucide-react';
import { useRecoilState } from "recoil";

import {QRCodeSVG} from 'qrcode.react';
import { currentOpenWalletAtAddressState } from "../../AppState";


function ReceiveView() {
    ////
    //// Global State
    ////

    const [currentOpenWalletAtAddress] = useRecoilState(currentOpenWalletAtAddressState);

    ////
    //// Component states.
    ////

    ////
    //// Event handling.
    ////

    ////
    //// Misc.
    ////

    useEffect(() => {
      let mounted = true;

      if (mounted) {
            window.scrollTo(0, 0); // Start the page at the top of the page.
      }

      return () => {
        mounted = false;
      };
    }, []);

    ////
    //// Component rendering.
    ////

    const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(currentOpenWalletAtAddress);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div>
      <main className="max-w-2xl mx-auto px-6 py-12 mb-24">
        <div className="bg-white rounded-xl border-2 border-gray-100 overflow-hidden">
          <div className="p-6">
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-purple-100 rounded-xl">
                  <QrCode className="w-5 h-5 text-purple-600" aria-hidden="true" />
                </div>
                <h2 className="text-xl font-bold text-gray-900">Receive ComicCoins</h2>
              </div>
            </div>
            <p className="text-sm text-gray-500">
              Share your wallet address or QR code to receive coins and NFTs.
            </p>
          </div>

          <div className="p-6 space-y-8">
            {/* QR Code Placeholder */}
            <div className="flex justify-center">
              <div className="p-4 bg-white rounded-2xl border-2 border-gray-100">
                {/* https://www.npmjs.com/package/qrcode.react */}
                <QRCodeSVG value={currentOpenWalletAtAddress} size={240} />
              </div>
            </div>

            {/* Wallet Address */}
            <div className="space-y-2">
              <label className="block text-sm font-medium text-gray-700">
                Your Wallet Address
              </label>
              <div className="flex items-center gap-2">
                <div className="flex-grow relative">
                  <input
                    type="text"
                    readOnly
                    value={currentOpenWalletAtAddress}
                    className="w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-lg font-mono text-gray-800"
                  />
                </div>
                <button
                  onClick={handleCopy}
                  className="flex items-center gap-2 px-4 py-3 bg-purple-100 text-purple-700 rounded-lg hover:bg-purple-200 transition-colors"
                >
                  {copied ? (
                    <CheckCircle2 className="w-5 h-5" />
                  ) : (
                    <Copy className="w-5 h-5" />
                  )}
                  {copied ? 'Copied!' : 'Copy'}
                </button>
              </div>
            </div>

            {/* Promotional Message */}
            <div className="pt-4 border-t border-gray-100">
              <p className="text-center text-sm text-gray-400">
                Want to earn free ComicCoins? Visit{' '}
                <a
                  href="https://faucet.comiccoin.com"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-purple-500 hover:text-purple-600 transition-colors"
                >
                  faucet.comiccoin.com
                </a>
              </p>
            </div>
          </div>
        </div>
      </main>

      {/* Bottom Navigation */}
      <nav className="fixed bottom-0 w-full bg-white border-t border-gray-200 shadow-lg" aria-label="Primary navigation">
        <div className="grid grid-cols-4 h-20">
          <button className="flex flex-col items-center justify-center space-y-2">
            <Wallet className="w-7 h-7 text-gray-600" aria-hidden="true" />
            <span className="text-sm text-gray-600">Overview</span>
          </button>
          <button className="flex flex-col items-center justify-center space-y-2">
            <Send className="w-7 h-7 text-gray-600" aria-hidden="true" />
            <span className="text-sm text-gray-600">Send</span>
          </button>
          <button className="flex flex-col items-center justify-center space-y-2 bg-purple-50" aria-current="page">
            <QrCode className="w-7 h-7 text-purple-600" aria-hidden="true" />
            <span className="text-sm text-purple-600">Receive</span>
          </button>
          <button className="flex flex-col items-center justify-center space-y-2">
            <MoreHorizontal className="w-7 h-7 text-gray-600" aria-hidden="true" />
            <span className="text-sm text-gray-600">More</span>
          </button>
        </div>
      </nav>
    </div>
  );
};

export default ReceiveView
