import React, { useState } from 'react';
import {
  Coins, Home, Image, History, Wallet,
  Settings, HelpCircle, LogOut, Clock, CheckCircle, XCircle,
  Menu, X, ChevronLeft, ChevronRight
} from 'lucide-react';

const SubmissionsPage = () => {
  const [isNavOpen, setIsNavOpen] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedSubmission, setSelectedSubmission] = useState(null);
  const itemsPerPage = 12;

  const navigation = [
    { name: 'Dashboard', icon: Home, current: false },
    { name: 'Submit Comic', icon: Image, current: false },
    { name: 'My Submissions', icon: History, current: true },
    { name: 'My Wallet', icon: Wallet, current: false },
    { name: 'Help', icon: HelpCircle, current: false },
    { name: 'Settings', icon: Settings, current: false },
  ];

  // Mock submissions data
  const submissions = [
    {
      id: 1,
      title: "Amazing Spider-Man #300",
      coverImage: "/api/placeholder/120/160",
      submittedAt: "2024-12-12T10:30:00",
      status: "accepted",
      coinsAwarded: 150,
      description: "First appearance of Venom",
      grade: "9.8",
      publisher: "Marvel Comics",
      year: "1988"
    },
    {
      id: 2,
      title: "Detective Comics #27",
      coverImage: "/api/placeholder/120/160",
      submittedAt: "2024-12-11T15:20:00",
      status: "in-review",
      description: "First appearance of Batman",
      grade: "4.0",
      publisher: "DC Comics",
      year: "1939"
    },
    {
      id: 3,
      title: "Fantastic Four #1",
      coverImage: "/api/placeholder/120/160",
      submittedAt: "2024-12-10T09:45:00",
      status: "rejected",
      reason: "Insufficient documentation",
      description: "Origin of the Fantastic Four",
      grade: "6.5",
      publisher: "Marvel Comics",
      year: "1961"
    }
  ];

  // Generate more mock submissions
  for (let i = 4; i <= 50; i++) {
    const status = i % 3 === 0 ? "accepted" : i % 3 === 1 ? "in-review" : "rejected";
    submissions.push({
      id: i,
      title: `Comic #${i}`,
      coverImage: "/api/placeholder/120/160",
      submittedAt: new Date(Date.now() - i * 3600000).toISOString(),
      status,
      coinsAwarded: status === "accepted" ? Math.floor(Math.random() * 200) + 50 : null,
      reason: status === "rejected" ? "Verification failed" : null,
      description: `Sample comic book #${i}`,
      grade: `${Math.floor(Math.random() * 3) + 7}.${Math.floor(Math.random() * 10)}`,
      publisher: i % 2 === 0 ? "Marvel Comics" : "DC Comics",
      year: `${1960 + Math.floor(Math.random() * 60)}`
    });
  }

  const totalPages = Math.ceil(submissions.length / itemsPerPage);
  const paginatedSubmissions = submissions.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

  const getStatusIcon = (status) => {
    switch (status) {
      case 'in-review':
        return <Clock className="w-4 h-4 text-yellow-500" />;
      case 'accepted':
        return <CheckCircle className="w-4 h-4 text-green-500" />;
      case 'rejected':
        return <XCircle className="w-4 h-4 text-red-500" />;
      default:
        return null;
    }
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 'in-review':
        return 'text-yellow-500';
      case 'accepted':
        return 'text-green-500';
      case 'rejected':
        return 'text-red-500';
      default:
        return '';
    }
  };

  const SubmissionCard = ({ submission }) => (
    <div
      className="w-32 bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow border border-purple-100 cursor-pointer"
      onClick={() => setSelectedSubmission(submission)}
    >
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
        <p className="text-xs mt-1">
          <span className={`font-medium ${getStatusColor(submission.status)}`}>
            {submission.status.charAt(0).toUpperCase() + submission.status.slice(1)}
          </span>
        </p>
        {submission.coinsAwarded && (
          <p className="text-xs text-green-600 mt-1">
            +{submission.coinsAwarded} ComicCoins
          </p>
        )}
        {submission.reason && (
          <p className="text-xs text-red-500 mt-1 truncate" title={submission.reason}>
            {submission.reason}
          </p>
        )}
        <p className="text-xs text-gray-500 mt-1">
          {new Date(submission.submittedAt).toLocaleDateString()}
        </p>
      </div>
    </div>
  );

  const SubmissionModal = ({ submission, onClose }) => {
    if (!submission) return null;

    return (
      <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
        <div className="bg-white rounded-xl max-w-2xl w-full p-6 relative">
          <button
            onClick={onClose}
            className="absolute top-4 right-4 text-gray-500 hover:text-gray-700"
          >
            <X className="w-6 h-6" />
          </button>

          <div className="flex gap-6">
            <img
              src={submission.coverImage}
              alt={submission.title}
              className="w-48 h-64 object-cover rounded-lg"
            />

            <div className="flex-1">
              <h2 className="text-2xl font-bold text-purple-800 mb-4">
                {submission.title}
              </h2>

              <div className="space-y-3">
                <div className="flex items-center gap-2">
                  <span className={`inline-flex items-center gap-1 px-3 py-1 rounded-full ${
                    getStatusColor(submission.status)} bg-opacity-10`}>
                    {getStatusIcon(submission.status)}
                    <span className="font-medium">
                      {submission.status.charAt(0).toUpperCase() + submission.status.slice(1)}
                    </span>
                  </span>
                  {submission.coinsAwarded && (
                    <span className="inline-flex items-center gap-1 px-3 py-1 rounded-full bg-green-100 text-green-600">
                      <Coins className="w-4 h-4" />
                      {submission.coinsAwarded} ComicCoins
                    </span>
                  )}
                </div>

                <div className="grid grid-cols-2 gap-4 text-sm">
                  <div>
                    <p className="text-gray-500">Publisher</p>
                    <p className="font-medium">{submission.publisher}</p>
                  </div>
                  <div>
                    <p className="text-gray-500">Year</p>
                    <p className="font-medium">{submission.year}</p>
                  </div>
                  <div>
                    <p className="text-gray-500">Grade</p>
                    <p className="font-medium">{submission.grade}</p>
                  </div>
                  <div>
                    <p className="text-gray-500">Submitted</p>
                    <p className="font-medium">
                      {new Date(submission.submittedAt).toLocaleDateString()}
                    </p>
                  </div>
                </div>

                <div>
                  <p className="text-gray-500">Description</p>
                  <p className="font-medium">{submission.description}</p>
                </div>

                {submission.reason && (
                  <div>
                    <p className="text-red-500">Rejection Reason</p>
                    <p className="font-medium text-red-600">{submission.reason}</p>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  };

  return (
    <div className="min-h-screen bg-purple-50">
      {/* Navigation with full width gradient background */}
      <nav className="bg-gradient-to-r from-purple-700 to-indigo-800 text-white">
        {/* Navigation content container with responsive padding */}
        <div className="max-w-[2000px] mx-auto px-4 md:px-8 lg:px-12 xl:px-24 2xl:px-32">
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
                className="inline-flex items-center justify-center p-2 rounded-md text-white hover:bg-purple-600 focus:outline-none"
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

        {/* Mobile menu with matching padding */}
        <div className={`lg:hidden ${isNavOpen ? 'block' : 'hidden'}`}>
          <div className="max-w-[2000px] mx-auto px-4 md:px-8 lg:px-12 xl:px-24 2xl:px-32 pt-2 pb-3 space-y-1">
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

      {/* Main Content with responsive padding */}
      <main className="max-w-[2000px] mx-auto px-4 md:px-8 lg:px-12 xl:px-24 2xl:px-32 py-8">
        <h1 className="text-3xl font-bold text-purple-800 mb-8" style={{fontFamily: 'Comic Sans MS, cursive'}}>
          My Submissions
        </h1>

        {/* Submissions Grid */}
        <div className="bg-white rounded-xl shadow-lg p-6 mb-8 border-2 border-purple-200">
          <div className="flex flex-wrap gap-4 justify-center lg:justify-start">
            {paginatedSubmissions.map(submission => (
              <SubmissionCard key={submission.id} submission={submission} />
            ))}
          </div>

          {/* Pagination */}
          <div className="mt-8 flex justify-center items-center gap-4">
            <button
              onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
              disabled={currentPage === 1}
              className="p-2 rounded-lg hover:bg-purple-100 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <ChevronLeft className="w-5 h-5 text-purple-600" />
            </button>

            <span className="text-sm text-purple-600">
              Page {currentPage} of {totalPages}
            </span>

            <button
              onClick={() => setCurrentPage(prev => Math.min(totalPages, prev + 1))}
              disabled={currentPage === totalPages}
              className="p-2 rounded-lg hover:bg-purple-100 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <ChevronRight className="w-5 h-5 text-purple-600" />
            </button>
          </div>
        </div>
      </main>

      {/* Modal */}
      {selectedSubmission && (
         <SubmissionModal
           submission={selectedSubmission}
           onClose={() => setSelectedSubmission(null)}
         />
       )}
     </div>
   );
 };

 export default SubmissionsPage;
