import { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import {
  WalletMinimal,
  Send,
  QrCode,
  MoreHorizontal,
  Clock,
  Coins,
  Wallet,
  ArrowRight,
  ArrowUpRight,
  ArrowDownLeft,
  Ticket,
} from "lucide-react";
import { useRecoilState } from "recoil";

import logo from "../../assets/images/CPS-logo-2023-square.webp";
import {
  GetTotalCoins,
  GetTotalTokens,
  GetRecentTransactions,
} from "../../../wailsjs/go/main/App";
import { currentOpenWalletAtAddressState } from "../../AppState";
import PageLoadingContent from "../Reusable/PageLoadingContent";
import { toLower } from "lodash";

function DashboardView() {
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
  const [totalCoins, setTotalCoins] = useState(0);
  const [totalTokens, setTotalTokens] = useState(0);
  const [transactions, setTransactions] = useState([]);

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
      console.log("currentOpenWalletAtAddress:", currentOpenWalletAtAddress);

      // Update the GUI to let user know that the operation is under way.
      setIsLoading(true);

      Promise.all([
        GetTotalCoins(currentOpenWalletAtAddress)
          .then((totalCoinsResult) => {
            console.log("GetTotalCoins: results:", totalCoinsResult);
            setTotalCoins(totalCoinsResult);
          })
          .catch((errorRes) => {
            console.log("GetTotalCoins: errors:", errorRes);
            if (errorRes.includes("address is null")) {
              setForceURL("/wallets");
            }
          }),
        GetTotalTokens(currentOpenWalletAtAddress)
          .then((totalTokensResult) => {
            console.log("GetTotalTokens: results:", totalTokensResult);
            setTotalTokens(totalTokensResult);
          })
          .catch((errorRes) => {
            console.log("GetTotalTokens: errors:", errorRes);
          }),
        GetRecentTransactions(currentOpenWalletAtAddress)
          .then((txsResponse) => {
            console.log("GetRecentTransactions: results:", txsResponse);
            setTransactions(txsResponse);
          })
          .catch((errorRes) => {
            console.log("GetRecentTransactions: errors:", errorRes);
          }),
      ])
        .then(() => {
          // Update the GUI to let user know that the operation is completed.
          setIsLoading(false);
        })
        .catch((error) => {
          console.log("Error:", error);
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

  // const recentTransactions = [
  //   {
  //     id: 1,
  //     type: 'Received',
  //     assetType: 'Coin',
  //     amount: '+50',
  //     from: '0x1234...5678',
  //     timestamp: '2 mins ago',
  //     description: 'Payment received from friend'
  //   },
  //   {
  //     id: 2,
  //     type: 'Sent',
  //     assetType: 'NFT',
  //     amount: 'Hero #123',
  //     to: '0x8765...4321',
  //     timestamp: '1 hour ago',
  //     description: 'Sent NFT to marketplace'
  //   },
  //   {
  //     id: 3,
  //     type: 'Received',
  //     assetType: 'Coin',
  //     amount: '+100',
  //     from: '0x9876...5432',
  //     timestamp: '3 hours ago',
  //     description: 'Salary payment'
  //   },
  //   {
  //     id: 4,
  //     type: 'Sent',
  //     assetType: 'Coin',
  //     amount: '-10',
  //     to: '0x2468...1357',
  //     timestamp: '5 hours ago',
  //     description: 'Coffee purchase'
  //   },
  //   {
  //     id: 5,
  //     type: 'Received',
  //     assetType: 'NFT',
  //     amount: 'Villain #456',
  //     from: '0x1357...2468',
  //     timestamp: 'Yesterday',
  //     description: 'NFT reward received'
  //   }
  // ];

  const handleTransactionClick = (tx) => {
    console.log(`Navigate to transaction ${tx}`);
    setForceURL("/more/transaction/" + tx.timestamp);
  };

  const formatAddress = (address) => {
    return `${address.slice(0, 6)}...${address.slice(-4)}`;
  };

  const getTransactionTime = (timestamp) => {
    const now = Date.now();
    const diff = now - timestamp;
    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);

    if (days > 0) return `${days} day${days > 1 ? "s" : ""} ago`;
    if (hours > 0) return `${hours} hour${hours > 1 ? "s" : ""} ago`;
    if (minutes > 0) return `${minutes} min ago`;
    return "Just now";
  };

  const isReceived = (tx) =>
    tx.to.toLowerCase() === currentOpenWalletAtAddress.toLowerCase();

  return (
    <div>
      <main className="max-w-2xl mx-auto px-6 py-12 mb-24">
        {/* Balance Cards */}
        <div className="grid grid-cols-2 gap-6 mb-6">
          {/* Coin Balance Card */}
          <div className="relative overflow-hidden rounded-2xl bg-gradient-to-br from-purple-500 to-purple-700 text-white p-1">
            <div className="absolute top-0 right-0 w-32 h-32 transform translate-x-16 -translate-y-8">
              <div className="w-full h-full rounded-full bg-purple-400 opacity-20"></div>
            </div>
            <div className="relative p-6 space-y-6">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-white/20 rounded-lg">
                  <Coins className="w-6 h-6" aria-hidden="true" />
                </div>
                <span className="text-sm font-medium">Total Coins</span>
              </div>
              <div className="space-y-1">
                <p
                  className="text-3xl font-bold"
                  aria-label={`${totalCoins} ComicCoins`}
                >
                  {totalCoins}
                </p>
                <p className="text-sm text-purple-200">Available Balance</p>
              </div>
            </div>
          </div>

          {/* NFT Balance Card */}
          <div className="relative overflow-hidden rounded-2xl bg-gradient-to-br from-indigo-500 to-indigo-700 text-white p-1">
            <div className="absolute top-0 right-0 w-32 h-32 transform translate-x-16 -translate-y-8">
              <div className="w-full h-full rounded-full bg-indigo-400 opacity-20"></div>
            </div>
            <div className="relative p-6 space-y-6">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-white/20 rounded-lg">
                  <Ticket className="w-6 h-6" aria-hidden="true" />
                </div>
                <span className="text-sm font-medium">Total NFTs</span>
              </div>
              <div className="space-y-1">
                <p
                  className="text-3xl font-bold"
                  aria-label={`${totalTokens} NFTs`}
                >
                  {totalTokens}
                </p>
                <p className="text-sm text-indigo-200">Collected Items</p>
              </div>
            </div>
          </div>
        </div>

        {/* Recent Transactions */}
        {transactions && (
          <div className="bg-white rounded-xl shadow-sm border border-gray-200">
            {transactions.length <= 0 ? (
              <>
                <div className="bg-white rounded-xl shadow-sm border border-gray-200">
                  <div className="p-6 border-b border-gray-100">
                    <div className="flex justify-between items-center">
                      <div className="flex items-center gap-3">
                        <div className="p-2 bg-purple-100 rounded-xl">
                          <Clock
                            className="w-5 h-5 text-purple-600"
                            aria-hidden="true"
                          />
                        </div>
                        <h2 className="text-xl font-bold text-gray-900">
                          Recent Transactions
                        </h2>
                      </div>
                    </div>
                  </div>

                  {/* Empty State Message */}
                  <div className="py-16 px-6">
                    <div className="text-center">
                      <div className="w-16 h-16 mx-auto mb-4 bg-purple-100 rounded-full flex items-center justify-center">
                        <Clock
                          className="w-8 h-8 text-purple-600"
                          aria-hidden="true"
                        />
                      </div>
                      <h3 className="text-lg font-medium text-gray-900 mb-2">
                        No Transactions Yet
                      </h3>
                      <p className="text-gray-500 max-w-sm mx-auto">
                        Start your journey by sending or receiving ComicCoins or
                        NFTs. Your transaction history will appear here.
                      </p>
                    </div>
                  </div>
                </div>
              </>
            ) : (
              <>
                <div className="p-6 border-b border-gray-100">
                  <div className="flex justify-between items-center">
                    <div className="flex items-center gap-3">
                      <div className="p-2 bg-purple-100 rounded-xl">
                        <Clock
                          className="w-5 h-5 text-purple-600"
                          aria-hidden="true"
                        />
                      </div>
                      <h2 className="text-xl font-bold text-gray-900">
                        Recent Transactions
                      </h2>
                    </div>
                    <button className="inline-flex items-center gap-2 text-purple-600 hover:text-purple-800 transition-colors text-lg">
                      See More
                      <ArrowRight className="w-5 h-5" />
                    </button>
                  </div>
                </div>

                {/* Recent Transactions */}
                <div className="bg-white rounded-xl shadow-sm border border-gray-200">
                  <div className="divide-y divide-gray-100">
                    {transactions.map((tx) => {
                      const received = isReceived(tx);

                      return (
                        <button
                          key={tx.timestamp}
                          onClick={() => handleTransactionClick(tx)}
                          className="w-full p-4 hover:bg-slate-50 transition-colors flex items-center justify-between gap-4"
                          aria-label={`${received ? "Received" : "Sent"} ${tx.value} coins ${
                            received ? "from" : "to"
                          } ${received ? tx.from : tx.to}, ${getTransactionTime(tx.timestamp)}. Click for details.`}
                        >
                          <div className="flex items-center gap-4">
                            <div
                              className={`p-3 rounded-xl ${
                                received ? "bg-green-100" : "bg-red-100"
                              }`}
                            >
                              {received ? (
                                <ArrowDownLeft className="w-5 h-5 text-green-600" />
                              ) : (
                                <ArrowUpRight className="w-5 h-5 text-red-600" />
                              )}
                            </div>
                            <div>
                              <div className="flex items-center gap-2">
                                <p
                                  className={`font-semibold ${
                                    received ? "text-green-600" : "text-red-600"
                                  }`}
                                >
                                  {received ? "+" : "-"}
                                  {/* What's going on here? If we received coins/token then the transactional fee was paid by the sender, so remove fee. If we are sender then show the transactional fee we paid. */}
                                  {`${received ? tx.value - 1 : tx.value}`} CC
                                </p>
                                <span>
                                  <span className="px-2 py-1 text-xs font-medium text-gray-600 bg-gray-100 rounded-full">
                                    {tx.type.toUpperCase()}
                                  </span>
                                </span>
                              </div>
                              <p className="text-sm text-gray-600 mt-1">
                                {received ? "From: " : "To: "}
                                {formatAddress(received ? tx.from : tx.to)}
                              </p>
                            </div>
                          </div>
                          <div className="flex items-center gap-3">
                            <span className="text-sm text-gray-500">
                              {getTransactionTime(tx.timestamp)}
                            </span>
                            <ArrowRight className="w-4 h-4 text-gray-400" />
                          </div>
                        </button>
                      );
                    })}
                  </div>
                </div>
              </>
            )}
          </div>
        )}
      </main>
    </div>
  );
}

export default DashboardView;
