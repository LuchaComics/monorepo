import React, { useState } from 'react';
import {
  Coins, Home, Image, History, Wallet,
  Settings, HelpCircle, LogOut, Menu, X,
  ArrowLeft, ArrowRight, ExternalLink
} from 'lucide-react';

import Topbar from "../../../Components/Navigation/Topbar";

const VerificationApplicationPage = () => {
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
      <Topbar currentPage="Settings" />

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 xl:px-12 2xl:px-24 py-8">
        TODO
      </main>
    </div>
  );
};

export default VerificationApplicationPage;
