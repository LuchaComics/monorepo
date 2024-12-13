import React, { useState } from 'react';
import {
  Coins, Home, Image, History, Wallet,
  Settings, HelpCircle, LogOut, Menu, X,
  ArrowLeft, ArrowRight, ExternalLink
} from 'lucide-react';

const MyWalletView = () => {
  const [isNavOpen, setIsNavOpen] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 10;

  const navigation = [
    { name: 'Dashboard', icon: Home, current: false },
    { name: 'Submit Comic', icon: Image, current: false },
    { name: 'My Submissions', icon: History, current: false },
    { name: 'My Wallet', icon: Wallet, current: true },
    { name: 'Help', icon: HelpCircle, current: false },
    { name: 'Settings', icon: Settings, current: false },
  ];

  // Mock transaction data
  const transactions = Array.from({ length: 50 }, (_, i) => ({
    id: i + 1,
    date: new Date(Date.now() - i * 86400000).toISOString(),
    amount: Math.floor(Math.random() * 76) + 25,
    type: 'Comic Submission Reward',
    comicTitle: `Comic #${Math.floor(Math.random() * 1000)}`,
    txHash: `0x${Array.from({ length: 64 }, () => Math.floor(Math.random() * 16).toString(16)).join('')}`
  }));

  const totalBalance = transactions.reduce((sum, tx) => sum + tx.amount, 0);
  const totalPages = Math.ceil(transactions.length / itemsPerPage);
  const currentTransactions = transactions.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  return (
    <div className="min-h-screen bg-purple-50">
      {/* Navigation - Using the same nav from the dashboard */}
      <nav className="bg-gradient-to-r from-purple-700 to-indigo-800 text-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 xl:px-12 2xl:px-24">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center space-x-2">
              <Coins className="h-8 w-8" />
              <span className="text-xl font-bold" style={{fontFamily: 'Comic Sans MS, cursive'}}>
                ComicCoin
              </span>
            </div>

            <div className="flex items-center lg:hidden">
              <button
                onClick={() => setIsNavOpen(!isNavOpen)}
                className="p-2 rounded-md text-white hover:bg-purple-600 focus:outline-none"
              >
                {isNavOpen ? <X className="h-6 w-6" /> : <Menu className="h-6 w-6" />}
              </button>
            </div>

            <div className="hidden lg:flex lg:items-center lg:space-x-4">
              {navigation.map((item) => (
                <a
                  key={item.name}
                  href="#"
                  className={`flex items-center space-x-1 px-3 py-2 rounded-md text-sm font-medium ${
                    item.current
                      ? 'bg-purple-600 bg-opacity-50'
                      : 'hover:bg-purple-600 hover:bg-opacity-25'
                  }`}
                >
                  <item.icon className="h-4 w-4" />
                  <span>{item.name}</span>
                </a>
              ))}
            </div>

            <div className="hidden lg:flex">
              <button className="flex items-center space-x-1 px-3 py-2 rounded-md hover:bg-purple-600 hover:bg-opacity-25 text-purple-200 hover:text-white">
                <LogOut className="h-4 w-4" />
                <span>Logout</span>
              </button>
            </div>
          </div>
        </div>

        <div className={`lg:hidden ${isNavOpen ? 'block' : 'hidden'}`}>
          <div className="px-4 pt-2 pb-3 space-y-1">
            {navigation.map((item) => (
              <a
                key={item.name}
                href="#"
                className={`flex items-center space-x-2 px-3 py-2 rounded-md text-base font-medium ${
                  item.current
                    ? 'bg-purple-600 bg-opacity-50'
                    : 'hover:bg-purple-600 hover:bg-opacity-25'
                }`}
              >
                <item.icon className="h-5 w-5" />
                <span>{item.name}</span>
              </a>
            ))}
            <button className="w-full flex items-center space-x-2 px-3 py-2 rounded-md text-base font-medium text-purple-200 hover:text-white hover:bg-purple-600 hover:bg-opacity-25">
              <LogOut className="h-5 w-5" />
              <span>Logout</span>
            </button>
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 xl:px-12 2xl:px-24 py-8">
        <h1 className="text-3xl font-bold text-purple-800 mb-8" style={{fontFamily: 'Comic Sans MS, cursive'}}>
          My Wallet
        </h1>

        {/* Wallet Balance Card */}
        <div className="bg-white rounded-xl shadow-lg p-6 mb-8 border-2 border-purple-200">
          <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4">
            <div>
              <h2 className="text-xl font-bold text-purple-800 mb-2">Total Balance</h2>
              <div className="flex items-center space-x-2">
                <Coins className="h-6 w-6 text-purple-600" />
                <span className="text-3xl font-bold">{totalBalance}</span>
                <span className="text-gray-600">ComicCoins</span>
              </div>
            </div>
            <div className="flex flex-col space-y-2">
              <div className="text-sm text-gray-600">Connected Wallet</div>
              <div className="flex items-center space-x-2">
                <span className="font-mono text-sm bg-purple-50 px-3 py-1 rounded">
                  0x1234...5678
                </span>
                <button className="text-purple-600 hover:text-purple-700">
                  <ExternalLink className="h-4 w-4" />
                </button>
              </div>
            </div>
          </div>
        </div>

        {/* Transaction History */}
        <div className="bg-white rounded-xl shadow-lg border-2 border-purple-200">
          <div className="p-6 border-b border-purple-100">
            <h2 className="text-xl font-bold text-purple-800">Transaction History</h2>
          </div>

          <div className="overflow-x-auto">
            <table className="w-full">
              <thead className="bg-purple-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-purple-800 uppercase tracking-wider">Date</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-purple-800 uppercase tracking-wider">Type</th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-purple-800 uppercase tracking-wider">Comic</th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-purple-800 uppercase tracking-wider">Amount</th>
                  <th className="px-6 py-3 text-right text-xs font-medium text-purple-800 uppercase tracking-wider">Transaction</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-purple-100">
                {currentTransactions.map((tx) => (
                  <tr key={tx.id} className="hover:bg-purple-50">
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-600">
                      {formatDate(tx.date)}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {tx.type}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {tx.comicTitle}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-right text-green-600 font-medium">
                      +{tx.amount}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-right text-sm">
                      <a href="#" className="text-purple-600 hover:text-purple-700 font-mono">
                        {`${tx.txHash.slice(0, 6)}...${tx.txHash.slice(-4)}`}
                      </a>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {/* Pagination */}
          <div className="px-6 py-4 flex items-center justify-between border-t border-purple-100">
            <div className="text-sm text-gray-600">
              Showing {(currentPage - 1) * itemsPerPage + 1} to {Math.min(currentPage * itemsPerPage, transactions.length)} of {transactions.length} transactions
            </div>
            <div className="flex space-x-2">
              <button
                onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
                disabled={currentPage === 1}
                className="p-2 rounded-md text-purple-600 hover:bg-purple-50 disabled:text-purple-300 disabled:hover:bg-transparent"
              >
                <ArrowLeft className="h-5 w-5" />
              </button>
              <button
                onClick={() => setCurrentPage(prev => Math.min(totalPages, prev + 1))}
                disabled={currentPage === totalPages}
                className="p-2 rounded-md text-purple-600 hover:bg-purple-50 disabled:text-purple-300 disabled:hover:bg-transparent"
              >
                <ArrowRight className="h-5 w-5" />
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
};

export default MyWalletView;
