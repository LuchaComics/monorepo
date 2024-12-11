import { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import {
  WalletMinimal,
  Send,
  QrCode,
  MoreHorizontal,
  Wallet,
  Settings,
  Coins,
  Clock,
  Ticket,
  Tickets,
} from "lucide-react";

function MoreView() {
  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
    }
    return () => {
      mounted = false;
    };
  }, []);

  return (
    <div>
      <main className="max-w-2xl mx-auto px-6 py-12 mb-24">
        <div className="grid gap-4 sm:grid-cols-2">
          <Link
            className="bg-white p-6 rounded-xl border-2 border-gray-100 hover:border-purple-200 transition-colors flex items-start gap-4 text-left"
            to={`/more/transactions`}
          >
            <div className="p-3 bg-purple-100 rounded-xl">
              <Coins className="w-6 h-6 text-purple-600" />
            </div>
            <div>
              <h2 className="font-semibold text-gray-900 mb-1">Transactions</h2>
              <p className="text-sm text-gray-500">
                View your recent transactions
              </p>
            </div>
          </Link>

          <Link
            className="bg-white p-6 rounded-xl border-2 border-gray-100 hover:border-purple-200 transition-colors flex items-start gap-4 text-left"
            to={`/more/wallets`}
          >
            <div className="p-3 bg-purple-100 rounded-xl">
              <Wallet className="w-6 h-6 text-purple-600" />
            </div>
            <div>
              <h2 className="font-semibold text-gray-900 mb-1">Wallets</h2>
              <p className="text-sm text-gray-500">
                View and sign into different wallets on your local computer
              </p>
            </div>
          </Link>

          <Link
            className="bg-white p-6 rounded-xl border-2 border-gray-100 hover:border-purple-200 transition-colors flex items-start gap-4 text-left"
            to={`/more/tokens`}
          >
            <div className="p-3 bg-purple-100 rounded-xl">
              <Tickets className="w-6 h-6 text-purple-600" />
            </div>
            <div>
              <h2 className="font-semibold text-gray-900 mb-1">Tokens</h2>
              <p className="text-sm text-gray-500">
                View and manage your tokens
              </p>
            </div>
          </Link>

          <Link
            className="bg-white p-6 rounded-xl border-2 border-gray-100 hover:border-purple-200 transition-colors flex items-start gap-4 text-left"
            to={`/more/settings`}
          >
            <div className="p-3 bg-purple-100 rounded-xl">
              <Settings className="w-6 h-6 text-purple-600" />
            </div>
            <div>
              <h2 className="font-semibold text-gray-900 mb-1">Settings</h2>
              <p className="text-sm text-gray-500">
                Configure your wallet settings
              </p>
            </div>
          </Link>
        </div>
      </main>
    </div>
  );
}

export default MoreView;
