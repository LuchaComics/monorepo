import React, { useState } from 'react';
import {
  Coins, Home, Settings, LogOut, Clock,
  CheckCircle, XCircle, Flag, ChevronLeft,
  ChevronRight, AlertTriangle
} from 'lucide-react';

const AdminDashboardPage = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 8;

  // Mock data for pending submissions with front and back covers
  const pendingSubmissions = Array.from({ length: 24 }, (_, i) => ({
    id: i + 1,
    title: `Comic #${i + 1}`,
    frontCover: "/api/placeholder/120/160",
    backCover: "/api/placeholder/120/160",
    submittedAt: new Date(Date.now() - i * 3600000).toISOString(),
    status: "review",
    submitter: `user_${i + 1}`,
    flagOptions: [
      "Duplicate submission",
      "Poor image quality",
      "Counterfeit",
      "Inappropriate content",
      "Other"
    ]
  }));

  const pageCount = Math.ceil(pendingSubmissions.length / itemsPerPage);
  const currentSubmissions = pendingSubmissions.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

  const GalleryItem = ({ submission }) => {
    const [showBackCover, setShowBackCover] = useState(false);
    const [showFlagMenu, setShowFlagMenu] = useState(false);
    const [selectedFlag, setSelectedFlag] = useState(null);

    return (
      <div className="w-64 bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow border border-purple-100">
        {/* Cover Image Section */}
        <div className="relative w-full h-80">
          <img
            src={showBackCover ? submission.backCover : submission.frontCover}
            alt={`${submission.title} - ${showBackCover ? 'Back' : 'Front'} Cover`}
            className="w-full h-full object-cover rounded-t-lg"
          />
          <div className="absolute top-2 left-2 right-2 flex justify-between">
            <button
              onClick={() => setShowBackCover(!showBackCover)}
              className="bg-white rounded-md px-2 py-1 text-xs font-medium shadow hover:bg-gray-50"
            >
              {showBackCover ? 'View Front' : 'View Back'}
            </button>
            <div className="bg-white rounded-full p-1 shadow">
              <Clock className="w-4 h-4 text-yellow-500" />
            </div>
          </div>

          {/* Action Buttons */}
          <div className="absolute bottom-2 left-2 right-2 flex justify-between">
            <div className="flex space-x-1">
              <button className="bg-white rounded-full p-2 shadow hover:bg-green-50">
                <CheckCircle className="w-5 h-5 text-green-500" />
              </button>
              <button className="bg-white rounded-full p-2 shadow hover:bg-red-50">
                <XCircle className="w-5 h-5 text-red-500" />
              </button>
              <button
                className="bg-white rounded-full p-2 shadow hover:bg-yellow-50"
                onClick={() => setShowFlagMenu(!showFlagMenu)}
              >
                <Flag className={`w-5 h-5 ${selectedFlag ? 'text-yellow-500' : 'text-gray-400'}`} />
              </button>
            </div>
          </div>

          {/* Flag Menu */}
          {showFlagMenu && (
            <div className="absolute bottom-14 left-2 bg-white rounded-lg shadow-lg p-2 w-48">
              <div className="text-xs font-medium text-gray-600 mb-2">Flag Issue:</div>
              {submission.flagOptions.map((option) => (
                <button
                  key={option}
                  onClick={() => {
                    setSelectedFlag(option);
                    setShowFlagMenu(false);
                  }}
                  className="w-full text-left text-xs px-2 py-1 hover:bg-purple-50 rounded"
                >
                  {option}
                </button>
              ))}
            </div>
          )}
        </div>

        {/* Submission Details */}
        <div className="p-3">
          <h3 className="font-medium text-sm truncate" title={submission.title}>
            {submission.title}
          </h3>
          <p className="text-xs text-gray-600 truncate">by {submission.submitter}</p>
          <p className="text-xs text-gray-500 mt-1">
            {new Date(submission.submittedAt).toLocaleDateString()}
          </p>
          {selectedFlag && (
            <div className="mt-2 flex items-center space-x-1 text-yellow-600 bg-yellow-50 rounded-md px-2 py-1">
              <AlertTriangle className="w-3 h-3" />
              <span className="text-xs">{selectedFlag}</span>
            </div>
          )}
        </div>
      </div>
    );
  };

  return (
    <div className="min-h-screen bg-purple-50">
      {/* Navigation */}
      <nav className="bg-gradient-to-r from-purple-700 to-indigo-800 text-white shadow-lg">
        <div className="max-w-full px-4">
          <div className="flex items-center justify-between h-16">
            {/* Logo */}
            <div className="flex items-center space-x-2">
              <Coins className="h-8 w-8" />
              <span className="text-xl font-bold" style={{fontFamily: 'Comic Sans MS, cursive'}}>
                ComicCoin Admin
              </span>
            </div>

            {/* Navigation Links */}
            <div className="flex items-center space-x-4">
              <a href="#" className="flex items-center space-x-1 px-3 py-2 rounded-md bg-purple-600 bg-opacity-50">
                <Home className="h-4 w-4" />
                <span>Dashboard</span>
              </a>
              <a href="#" className="flex items-center space-x-1 px-3 py-2 rounded-md hover:bg-purple-600 hover:bg-opacity-25">
                <Settings className="h-4 w-4" />
                <span>Settings</span>
              </a>
            </div>

            {/* Logout Button */}
            <button className="flex items-center space-x-1 px-3 py-2 rounded-md hover:bg-purple-600 hover:bg-opacity-25 text-purple-200 hover:text-white">
              <LogOut className="h-4 w-4" />
              <span>Logout</span>
            </button>
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main className="p-8">
        <h1 className="text-3xl font-bold text-purple-800 mb-8" style={{fontFamily: 'Comic Sans MS, cursive'}}>
          Admin Dashboard
        </h1>

        {/* Stats Row */}
        <div className="flex justify-between items-center mb-8 space-x-6">
          <div className="flex-1 bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
            <div className="text-purple-600 text-lg font-semibold">New Users This Week</div>
            <div className="text-3xl font-bold">47</div>
          </div>
          <div className="flex-1 bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
            <div className="text-purple-600 text-lg font-semibold">Pending Reviews</div>
            <div className="text-3xl font-bold">{pendingSubmissions.length}</div>
          </div>
          <div className="flex-1 bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
            <div className="text-purple-600 text-lg font-semibold">Total ComicCoins Paid</div>
            <div className="text-3xl font-bold">156,750</div>
          </div>
        </div>

        {/* Pending Submissions Gallery */}
        <div className="bg-white rounded-xl shadow-lg p-6 mb-8 border-2 border-purple-200">
          <h2 className="text-2xl font-bold text-purple-800 mb-6" style={{fontFamily: 'Comic Sans MS, cursive'}}>
            Submissions Awaiting Review
          </h2>
          <div className="flex flex-wrap gap-6">
            {currentSubmissions.map(submission => (
              <GalleryItem key={submission.id} submission={submission} />
            ))}
          </div>

          {/* Pagination */}
          <div className="mt-8 flex items-center justify-between">
            <div className="text-sm text-gray-600">
              Showing {(currentPage - 1) * itemsPerPage + 1} to {Math.min(currentPage * itemsPerPage, pendingSubmissions.length)} of {pendingSubmissions.length} submissions
            </div>
            <div className="flex items-center space-x-2">
              <button
                onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
                disabled={currentPage === 1}
                className={`p-2 rounded-md ${currentPage === 1 ? 'text-gray-400 cursor-not-allowed' : 'text-purple-600 hover:bg-purple-50'}`}
              >
                <ChevronLeft className="w-5 h-5" />
              </button>
              {Array.from({ length: pageCount }, (_, i) => (
                <button
                  key={i + 1}
                  onClick={() => setCurrentPage(i + 1)}
                  className={`px-3 py-1 rounded-md ${
                    currentPage === i + 1
                      ? 'bg-purple-600 text-white'
                      : 'text-purple-600 hover:bg-purple-50'
                  }`}
                >
                  {i + 1}
                </button>
              ))}
              <button 
                onClick={() => setCurrentPage(prev => Math.min(pageCount, prev + 1))}
                disabled={currentPage === pageCount}
                className={`p-2 rounded-md ${currentPage === pageCount ? 'text-gray-400 cursor-not-allowed' : 'text-purple-600 hover:bg-purple-50'}`}
              >
                <ChevronRight className="w-5 h-5" />
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
};

export default AdminDashboardPage;
