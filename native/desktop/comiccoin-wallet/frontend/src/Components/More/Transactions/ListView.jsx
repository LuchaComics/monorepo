import { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import {
  WalletMinimal,
  Send,
  QrCode,
  MoreHorizontal,
  Clock,
  Wallet,
  ArrowUpRight,
  ArrowDownLeft,
  ArrowRight,
} from "lucide-react";

import { useRecoilState } from "recoil";
import { toLower } from "lodash";

import { GetTransactions } from "../../../../wailsjs/go/main/App";
import { currentOpenWalletAtAddressState } from "../../../AppState";

function ListTransactionsView() {
  ////
  //// Global State
  ////

  const [currentOpenWalletAtAddress] = useRecoilState(
    currentOpenWalletAtAddressState,
  );

  ////
  //// Component states.
  ////

  const [isLoading, setIsLoading] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [transactions, setTransactions] = useState([]);

  ////
  //// Event handling.
  ////

  const formatAddress = (address) => {
    return `${address.slice(0, 6)}...${address.slice(-4)}`;
  };

  const formatDate = (timestamp) => {
    const date = new Date(timestamp);
    return date.toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const handleTransactionClick = (tx) => {
    console.log(`Navigate to transaction details: ${tx}`);
    // Navigation logic would go here
    setForceURL("/more/transaction/" + tx.timestamp);
  };

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.

      setIsLoading(true);
      GetTransactions(currentOpenWalletAtAddress)
        .then((txsResponse) => {
          console.log("GetTransactions: results:", txsResponse);
          setTransactions(txsResponse);
        })
        .catch((errorRes) => {
          console.log("GetTransactions: errors:", errorRes);
        })
        .finally(() => {
          setIsLoading(false);
        });
    }

    return () => {
      mounted = false;
    };
  }, [currentOpenWalletAtAddress]);

  ////
  //// Component rendering.
  ////

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  // if (isLoading) {
  //     return (
  //         "----"
  //     );
  // }

  // Mock data with the specified fields
  // const transactions = [
  //   {
  //     id: 1,
  //     timestamp: '2024-04-10T14:30:00Z',
  //     type: 'receive',
  //     from: '0x1234567890abcdef1234567890abcdef12345678',
  //     to: '0x8765432109876543210987654321098765432109',
  //     value: '50.00'
  //   },
  //   {
  //     id: 2,
  //     timestamp: '2024-04-10T12:15:00Z',
  //     type: 'send',
  //     from: '0x8765432109876543210987654321098765432109',
  //     to: '0x9876543210123456789009876543210987654321',
  //     value: '25.50'
  //   },
  //   {
  //     id: 3,
  //     timestamp: '2024-04-09T18:45:00Z',
  //     type: 'receive',
  //     from: '0x1111222233334444555566667777888899990000',
  //     to: '0x8765432109876543210987654321098765432109',
  //     value: '100.00'
  //   },
  //   {
  //     id: 4,
  //     timestamp: '2024-04-09T10:20:00Z',
  //     type: 'send',
  //     from: '0x8765432109876543210987654321098765432109',
  //     to: '0x2222111133334444555566667777888899990000',
  //     value: '75.25'
  //   },
  //   {
  //     id: 5,
  //     timestamp: '2024-04-08T16:30:00Z',
  //     type: 'receive',
  //     from: '0x3333222211114444555566667777888899990000',
  //     to: '0x8765432109876543210987654321098765432109',
  //     value: '200.00'
  //   }
  // ];

  return (
    <div>
      <main className="max-w-2xl mx-auto px-6 py-12 mb-24">
        {/* Transactions List */}
        <div className="bg-white rounded-xl shadow-sm border border-gray-200">
          <div className="p-6 border-b border-gray-100">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-purple-100 rounded-xl">
                <Clock className="w-5 h-5 text-purple-600" aria-hidden="true" />
              </div>
              <h2 className="text-xl font-bold text-gray-900">
                All Transactions
              </h2>
            </div>
          </div>

          <div className="divide-y divide-gray-100">
            {transactions.map((tx) => {
              const isReceived =
                tx.to.toLowerCase() ===
                currentOpenWalletAtAddress.toLowerCase();

              // What's going on here? If we received coins/token then the transactional
              // fee was paid by the sender, so remove fee. If we are sender then show the
              // transactional fee we paid.
              const totalValue =
                parseFloat(tx.value) - (isReceived ? parseFloat(tx.fee) : 0);

              return (
                <button
                  key={`${tx.timestamp}-${tx.value}`}
                  onClick={() => handleTransactionClick(tx)}
                  className="w-full p-4 hover:bg-slate-50 transition-colors text-left"
                >
                  <div className="flex items-center justify-between gap-4">
                    <div className="flex items-center gap-4">
                      <div
                        className={`p-3 rounded-xl ${
                          isReceived ? "bg-green-100" : "bg-red-100"
                        }`}
                      >
                        {isReceived ? (
                          <ArrowDownLeft className="w-5 h-5 text-green-600" />
                        ) : (
                          <ArrowUpRight className="w-5 h-5 text-red-600" />
                        )}
                      </div>
                      <div>
                        <div className="flex items-center gap-2">
                          <p
                            className={`font-semibold ${
                              isReceived ? "text-green-600" : "text-red-600"
                            }`}
                          >
                            {isReceived ? "+" : "-"}
                            {totalValue} CC
                          </p>
                          <span className="px-2 py-1 text-xs font-medium text-gray-600 bg-gray-100 rounded-full">
                            {tx.type.toUpperCase()}
                          </span>
                          {!isReceived && (
                            <span className="text-xs text-gray-500">
                              (Fee: {tx.fee} CC)
                            </span>
                          )}
                        </div>
                        <div className="flex flex-col sm:flex-row sm:items-center gap-1 sm:gap-2 mt-1">
                          <p className="text-sm text-gray-600">
                            {isReceived ? "From:" : "To:"}{" "}
                            {formatAddress(isReceived ? tx.from : tx.to)}
                          </p>
                          <span className="hidden sm:inline text-gray-400">
                            •
                          </span>
                          <span className="text-sm text-gray-500">
                            {formatDate(tx.timestamp)}
                          </span>
                        </div>
                      </div>
                    </div>
                    <ArrowRight className="w-4 h-4 text-gray-400 flex-shrink-0" />
                  </div>
                </button>
              );
            })}
          </div>
        </div>
      </main>
    </div>
  );
}

export default ListTransactionsView;
