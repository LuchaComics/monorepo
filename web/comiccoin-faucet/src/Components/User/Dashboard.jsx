import React, { useState } from 'react';
import {
  Coins, Home, Image, History, Wallet,
  Settings, HelpCircle, LogOut, Clock, CheckCircle, XCircle,
  Menu, X
} from 'lucide-react';

const DashboardPage = () => {
  const [isNavOpen, setIsNavOpen] = useState(false);

  const navigation = [
    { name: 'Dashboard', icon: Home, current: true },
    { name: 'Submit Comic', icon: Image, current: false },
    { name: 'My Submissions', icon: History, current: false },
    { name: 'My Wallet', icon: Wallet, current: false },
    { name: 'Help', icon: HelpCircle, current: false },
    { name: 'Settings', icon: Settings, current: false },
  ];

  // Mock data for pending submissions
  const pendingSubmissions = [
    {
      id: 1,
      title: "Amazing Spider-Man #300",
      coverImage: "/api/placeholder/120/160",
      submittedAt: "2024-12-12T10:30:00",
      status: "review",
      submitter: "peter_parker"
    },
    {
      id: 2,
      title: "Batman: Origins",
      coverImage: "/api/placeholder/120/160",
      submittedAt: "2024-12-12T09:45:00",
      status: "review",
      submitter: "bruce_wayne"
    }
  ];

  // Add more mock pending submissions
  for (let i = 3; i <= 12; i++) {
    pendingSubmissions.push({
      id: i,
      title: `Comic Book #${i}`,
      coverImage: "/api/placeholder/120/160",
      submittedAt: new Date(Date.now() - i * 3600000).toISOString(),
      status: "review",
      submitter: `user_${i}`
    });
  }

  // Mock data for completed submissions
  const completedSubmissions = [
    {
      id: 101,
      title: "X-Men #141",
      coverImage: "/api/placeholder/120/160",
      submittedAt: "2024-12-11T15:20:00",
      status: "approved",
      submitter: "charles_xavier",
      coinsAwarded: 50
    },
    {
      id: 102,
      title: "Superman #75",
      coverImage: "/api/placeholder/120/160",
      submittedAt: "2024-12-11T14:10:00",
      status: "approved",
      submitter: "clark_kent",
      coinsAwarded: 75
    },
    {
      id: 103,
      title: "Amazing Spider-Man #300",
      coverImage: "/api/placeholder/120/160",
      submittedAt: "2024-12-11T13:00:00",
      status: "rejected",
      submitter: "eddie_brock",
      reason: "Duplicate submission"
    }
  ];

  // Add more mock completed submissions
  for (let i = 4; i <= 20; i++) {
    completedSubmissions.push({
      id: 100 + i,
      title: `Completed Comic #${i}`,
      coverImage: "/api/placeholder/120/160",
      submittedAt: new Date(Date.now() - i * 3600000).toISOString(),
      status: i % 5 === 0 ? "rejected" : "approved",
      submitter: `user_${i}`,
      coinsAwarded: i % 5 === 0 ? null : Math.floor(Math.random() * 50) + 25,
      reason: i % 5 === 0 ? "Duplicate submission" : null
    });
  }

  const getStatusIcon = (status) => {
    switch (status) {
      case 'review':
        return <Clock className="w-4 h-4 text-yellow-500" />;
      case 'approved':
        return <CheckCircle className="w-4 h-4 text-green-500" />;
      case 'rejected':
        return <XCircle className="w-4 h-4 text-red-500" />;
      default:
        return null;
    }
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'review':
        return 'text-yellow-500';
      case 'approved':
        return 'text-green-500';
      case 'rejected':
        return 'text-red-500';
      default:
        return '';
    }
  };

  const GalleryItem = ({ submission }) => (
    <div className="w-32 bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow border border-purple-100">
      <div className="relative w-32 h-44">
        <img
          src={submission.coverImage}
          alt={submission.title}
          className="w-full h-full object-cover rounded-t-lg"
        />
        <div className="absolute top-1 right-1 bg-white rounded-full p-1 shadow">
          {getStatusIcon(submission.status)}
        </div>
      </div>
      <div className="p-2">
        <h3 className="font-medium text-xs truncate" title={submission.title}>
          {submission.title}
        </h3>
        <p className="text-xs text-gray-600 truncate">by {submission.submitter}</p>
        <p className="text-xs mt-1">
          <span className={`font-medium ${getStatusColor(submission.status)}`}>
            {submission.status === 'review' ? 'In Review' :
             submission.status === 'approved' ? 'Approved' : 'Rejected'}
          </span>
        </p>
        {submission.coinsAwarded && (
          <p className="text-xs text-green-600 mt-1">
            +{submission.coinsAwarded} ComicCoins
          </p>
        )}
        {submission.reason && (
          <p className="text-xs text-red-500 mt-1">
            {submission.reason}
          </p>
        )}
        <p className="text-xs text-gray-500 mt-1">
          {new Date(submission.submittedAt).toLocaleDateString()}
        </p>
      </div>
    </div>
  );

  return (
    <div className="min-h-screen bg-purple-50">
      {/* Mobile-first Navigation */}
      <nav className="bg-gradient-to-r from-purple-700 to-indigo-800 text-white">
        {/* Mobile navigation header */}
        <div className="px-4">
          <div className="flex items-center justify-between h-16">
            {/* Logo */}
            <div className="flex items-center space-x-2">
              <Coins className="h-8 w-8" />
              <span className="text-xl font-bold" style={{fontFamily: 'Comic Sans MS, cursive'}}>
                ComicCoin
              </span>
            </div>

            {/* Mobile menu button */}
            <div className="flex items-center lg:hidden">
              <button
                onClick={() => setIsNavOpen(!isNavOpen)}
                className="inline-flex items-center justify-center p-2 rounded-md text-white hover:bg-purple-600 focus:outline-none"
              >
                {isNavOpen ? (
                  <X className="h-6 w-6" />
                ) : (
                  <Menu className="h-6 w-6" />
                )}
              </button>
            </div>

            {/* Desktop navigation */}
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

            {/* Desktop Logout button */}
            <div className="hidden lg:flex">
              <button className="flex items-center space-x-1 px-3 py-2 rounded-md hover:bg-purple-600 hover:bg-opacity-25 text-purple-200 hover:text-white">
                <LogOut className="h-4 w-4" />
                <span>Logout</span>
              </button>
            </div>
          </div>
        </div>

        {/* Mobile menu, show/hide based on menu state */}
        <div className={`lg:hidden ${isNavOpen ? 'block' : 'hidden'}`}>
          <div className="px-2 pt-2 pb-3 space-y-1">
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
            {/* Mobile Logout button */}
            <button
              className="w-full flex items-center space-x-2 px-3 py-2 rounded-md text-base font-medium text-purple-200 hover:text-white hover:bg-purple-600 hover:bg-opacity-25"
            >
              <LogOut className="h-5 w-5" />
              <span>Logout</span>
            </button>
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main className="p-8">
        <h1 className="text-3xl font-bold text-purple-800 mb-8" style={{fontFamily: 'Comic Sans MS, cursive'}}>Dashboard</h1>

        {/* Stats Row */}
        <div className="flex justify-between items-center mb-8 space-x-6">
          <div className="flex-1 bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
            <div className="text-purple-600 text-lg font-semibold">Total Submissions</div>
            <div className="text-3xl font-bold">127</div>
          </div>
          <div className="flex-1 bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
            <div className="text-purple-600 text-lg font-semibold">Comics Approved</div>
            <div className="text-3xl font-bold">98</div>
          </div>
          <div className="flex-1 bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
            <div className="text-purple-600 text-lg font-semibold">ComicCoins Earned</div>
            <div className="text-3xl font-bold">4,750</div>
          </div>
        </div>

        {/* Pending Submissions Gallery */}
        <div className="bg-white rounded-xl shadow-lg p-6 mb-8 border-2 border-purple-200">
          <h2 className="text-2xl font-bold text-purple-800 mb-6" style={{fontFamily: 'Comic Sans MS, cursive'}}>
            Pending Reviews
          </h2>
          <div className="flex flex-wrap gap-4">
            {pendingSubmissions.map(submission => (
              <GalleryItem key={submission.id} submission={submission} />
            ))}
          </div>
        </div>

        {/* Completed Submissions Gallery */}
        <div className="bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
          <h2 className="text-2xl font-bold text-purple-800 mb-6" style={{fontFamily: 'Comic Sans MS, cursive'}}>
            Recent Submissions
          </h2>
          <div className="flex flex-wrap gap-4">
            {completedSubmissions.map(submission => (
              <GalleryItem key={submission.id} submission={submission} />
            ))}
          </div>
        </div>
      </main>
    </div>
  );
};

export default DashboardPage;
