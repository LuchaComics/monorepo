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
  const walletAddress = "cc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh";

  const handleCopy = () => {
    navigator.clipboard.writeText(walletAddress);
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
                <svg width="240" height="240" viewBox="0 0 240 240" className="bg-white">
                  <rect x="0" y="0" width="240" height="240" fill="white" />

                  {/* QR Code Corner Markers */}
                  {/* Top Left */}
                  <g>
                    <rect x="20" y="20" width="60" height="60" fill="black" />
                    <rect x="30" y="30" width="40" height="40" fill="white" />
                    <rect x="40" y="40" width="20" height="20" fill="black" />
                  </g>

                  {/* Top Right */}
                  <g>
                    <rect x="160" y="20" width="60" height="60" fill="black" />
                    <rect x="170" y="30" width="40" height="40" fill="white" />
                    <rect x="180" y="40" width="20" height="20" fill="black" />
                  </g>

                  {/* Bottom Left */}
                  <g>
                    <rect x="20" y="160" width="60" height="60" fill="black" />
                    <rect x="30" y="170" width="40" height="40" fill="white" />
                    <rect x="40" y="180" width="20" height="20" fill="black" />
                  </g>

                  {/* QR Code Data Pattern (Simplified) */}
                  {[...Array(10)].map((_, i) => (
                    <rect
                      key={i}
                      x={90 + (i * 10)}
                      y={90}
                      width="8"
                      height="8"
                      fill={Math.random() > 0.5 ? "black" : "white"}
                    />
                  ))}

                  {/* Additional random patterns */}
                  {[...Array(8)].map((_, row) => (
                    [...Array(8)].map((_, col) => (
                      <rect
                        key={`${row}-${col}`}
                        x={90 + (col * 10)}
                        y={110 + (row * 10)}
                        width="8"
                        height="8"
                        fill={Math.random() > 0.5 ? "black" : "white"}
                      />
                    ))
                  ))}
                </svg>
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
                    value={walletAddress}
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
