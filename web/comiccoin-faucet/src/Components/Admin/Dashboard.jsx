import React, { useState, useCallback, useEffect } from 'react';
import {
  Coins, Home, Settings, LogOut, Clock,
  CheckCircle, XCircle, Flag, ChevronLeft,
  ChevronRight, AlertTriangle, Menu, X
} from 'lucide-react';
import { Navigate, Link } from "react-router-dom";
import { useRecoilState } from "recoil";

import { currentUserState } from "../../AppState";
import Topbar from "../../Components/Navigation/Topbar";
import {
    getComicSubmissionsCountByFilterAPI,
    getComicSubmissionsCountCoinsRewardByFilterAPI,
    getComicSubmissionsCountTotalCreatedTodayByUserAPI,
    getComicSubmissionListAPI
} from "../../API/ComicSubmission";


const AdminDashboard = () => {

  // Variable controls the global state of the app.
  const [currentUser] = useRecoilState(currentUserState);

  // Data related.
  const [pendingSubmissions, setPendingSubmissions] = useState([]);
  const [isFetching, setFetching] = useState(false);
  const [errors, setErrors] = useState({});

  // GUI related
  const [currentPage, setCurrentPage] = useState(1);
  const [isNavOpen, setIsNavOpen] = useState(false);
  const [submissions, setSubmissions] = useState(Array.from({ length: 24 }, (_, i) => ({
    id: i + 1,
    title: `Comic #${i + 1}`,
    frontCover: "/api/placeholder/120/160",
    backCover: "/api/placeholder/120/160",
    submittedAt: new Date(Date.now() - i * 3600000).toISOString(),
    status: "review",
    submitter: `user_${i + 1}`,
    flagReason: null
  })));

  const itemsPerPage = 8;
  const pageCount = Math.ceil(submissions.length / itemsPerPage);
  const currentSubmissions = submissions.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

  // Navigation Handlers
  const toggleNav = () => {
    setIsNavOpen(!isNavOpen);
  };

  // Submission Action Handlers
  const handleApproveSubmission = useCallback((submissionId) => {
    setSubmissions(prevSubmissions =>
      prevSubmissions.map(submission => {
        if (submission.id === submissionId) {
          // In a real application, you would make an API call here
          console.log(`Approving submission ${submissionId}`);
          return {
            ...submission,
            status: 'approved',
            reviewedAt: new Date().toISOString(),
            reviewResult: 'approved'
          };
        }
        return submission;
      })
    );
  }, []);

  const handleRejectSubmission = useCallback((submissionId) => {
    setSubmissions(prevSubmissions =>
      prevSubmissions.map(submission => {
        if (submission.id === submissionId) {
          // In a real application, you would make an API call here
          console.log(`Rejecting submission ${submissionId}`);
          return {
            ...submission,
            status: 'rejected',
            reviewedAt: new Date().toISOString(),
            reviewResult: 'rejected'
          };
        }
        return submission;
      })
    );
  }, []);

  const handleFlagSubmission = useCallback((submissionId, flagReason) => {
    setSubmissions(prevSubmissions =>
      prevSubmissions.map(submission => {
        if (submission.id === submissionId) {
          // In a real application, you would make an API call here
          console.log(`Flagging submission ${submissionId} for: ${flagReason}`);
          return {
            ...submission,
            flagReason,
            status: submission.status === 'flagged' ? 'review' : 'flagged'
          };
        }
        return submission;
      })
    );
  }, []);

  // Pagination Handlers
  const handlePageChange = useCallback((newPage) => {
    if (newPage >= 1 && newPage <= pageCount) {
      setCurrentPage(newPage);
      // In a real application, you might want to fetch new data here
      console.log(`Navigating to page ${newPage}`);
    }
  }, [pageCount]);

  const handleNextPage = useCallback(() => {
    handlePageChange(currentPage + 1);
  }, [currentPage, handlePageChange]);

  const handlePrevPage = useCallback(() => {
    handlePageChange(currentPage - 1);
  }, [currentPage, handlePageChange]);

  // Flag Options Menu Component
  const FlagOptionsMenu = ({ submissionId, onClose }) => {
    const flagOptions = [
      "Duplicate submission",
      "Poor image quality",
      "Counterfeit",
      "Inappropriate content",
      "Other"
    ];

    return (
      <div className="absolute bottom-14 left-2 bg-white rounded-lg shadow-lg p-2 w-48 z-10">
        <div className="text-xs font-medium text-gray-600 mb-2">Flag Issue:</div>
        {flagOptions.map((option) => (
          <button
            key={option}
            onClick={() => {
              handleFlagSubmission(submissionId, option);
              onClose();
            }}
            className="w-full text-left text-xs px-2 py-1 hover:bg-purple-50 rounded"
          >
            {option}
          </button>
        ))}
      </div>
    );
  };

  // Gallery Item Component with Action Handlers
  const GalleryItem = ({ submission }) => {
    const [showBackCover, setShowBackCover] = useState(false);
    const [showFlagMenu, setShowFlagMenu] = useState(false);

    const toggleCover = () => setShowBackCover(prev => !prev);
    const toggleFlagMenu = () => setShowFlagMenu(prev => !prev);



    return (
      <div className="w-64 bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow border border-purple-100">
        <div className="relative w-full h-80">
          <img
            src={showBackCover ? submission.backCover : submission.frontCover}
            alt={`${submission.title} - ${showBackCover ? 'Back' : 'Front'} Cover`}
            className="w-full h-full object-cover rounded-t-lg"
          />
          <div className="absolute top-2 left-2 right-2 flex justify-between">
            <button
              onClick={toggleCover}
              className="bg-white rounded-md px-2 py-1 text-xs font-medium shadow hover:bg-gray-50"
            >
              {showBackCover ? 'View Front' : 'View Back'}
            </button>
            <div className="bg-white rounded-full p-1 shadow">
              <Clock className="w-4 h-4 text-yellow-500" />
            </div>
          </div>

          <div className="absolute bottom-2 left-2 right-2 flex justify-between">
            <div className="flex space-x-1">
              <button
                onClick={() => handleApproveSubmission(submission.id)}
                className="bg-white rounded-full p-2 shadow hover:bg-green-50"
                title="Approve Submission"
              >
                <CheckCircle className="w-5 h-5 text-green-500" />
              </button>
              <button
                onClick={() => handleRejectSubmission(submission.id)}
                className="bg-white rounded-full p-2 shadow hover:bg-red-50"
                title="Reject Submission"
              >
                <XCircle className="w-5 h-5 text-red-500" />
              </button>
              <button
                onClick={toggleFlagMenu}
                className="bg-white rounded-full p-2 shadow hover:bg-yellow-50"
                title="Flag for Review"
              >
                <Flag className={`w-5 h-5 ${submission.flagReason ? 'text-yellow-500' : 'text-gray-400'}`} />
              </button>
            </div>
          </div>

          {showFlagMenu && (
            <FlagOptionsMenu
              submissionId={submission.id}
              onClose={() => setShowFlagMenu(false)}
            />
          )}
        </div>

        <div className="p-3">
          <h3 className="font-medium text-sm truncate" title={submission.title}>
            {submission.title}
          </h3>
          <p className="text-xs text-gray-600 truncate">by {submission.submitter}</p>
          <p className="text-xs text-gray-500 mt-1">
            {new Date(submission.submittedAt).toLocaleDateString()}
          </p>
          {submission.flagReason && (
            <div className="mt-2 flex items-center space-x-1 text-yellow-600 bg-yellow-50 rounded-md px-2 py-1">
              <AlertTriangle className="w-3 h-3" />
              <span className="text-xs">{submission.flagReason}</span>
            </div>
          )}
        </div>
      </div>
    );
  };

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
      //------------------------------------------------------------------------

      let params = new Map();
      // params.set("page_size", limit); // Pagination
      // params.set("sort_field", "created_at"); // Sorting
      // params.set("sort_order", -1); // Sorting - descending, meaning most recent start date to oldest start date.
      params.set("status", 1); // ComicSubmissionStatusInReview
      //
      // params.set("store_id", sid);
      //
      // if (cur !== "") {
      //   // Pagination
      //   params.set("cursor", cur);
      // }
      //
      // // Filtering
      // if (keywords !== undefined && keywords !== null && keywords !== "") {
      //   // Searhcing
      //   params.set("search", keywords);
      // }
      getComicSubmissionListAPI(
        params,
        (resp) => {
          // For debugging purposes only.
          console.log("getComicSubmissionListAPI: Starting...");
          console.log(resp);
          setPendingSubmissions(resp.submissions);
        },
        (apiErr) => {
          console.log("getComicSubmissionListAPI: apiErr:", apiErr);
          setErrors(apiErr);
        },
        () => {
          console.log("getComicSubmissionListAPI: Starting...");
          setFetching(false);
        },
        () => {
          console.log("getComicSubmissionListAPI: unauthorized...");
          window.location.href = "/login?unauthorized=true";
        },
      );

      //------------------------------------------------------------------------
    }

    return () => {
      mounted = false;
    };
  }, [currentUser]);

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

            {/* Mobile menu button */}
            <div className="flex lg:hidden">
              <button
                onClick={toggleNav}
                className="inline-flex items-center justify-center p-2 rounded-md text-white hover:bg-purple-600 focus:outline-none"
              >
                {isNavOpen ? (
                  <X className="h-6 w-6" />
                ) : (
                  <Menu className="h-6 w-6" />
                )}
              </button>
            </div>

            {/* Desktop Navigation */}
            <div className="hidden lg:flex lg:items-center lg:space-x-4">
              <a href="#" className="flex items-center space-x-1 px-3 py-2 rounded-md bg-purple-600 bg-opacity-50">
                <Home className="h-4 w-4" />
                <span>Dashboard</span>
              </a>
              <a href="#" className="flex items-center space-x-1 px-3 py-2 rounded-md hover:bg-purple-600 hover:bg-opacity-25">
                <Settings className="h-4 w-4" />
                <span>Settings</span>
              </a>
              <button className="flex items-center space-x-1 px-3 py-2 rounded-md hover:bg-purple-600 hover:bg-opacity-25 text-purple-200 hover:text-white">
                <LogOut className="h-4 w-4" />
                <span>Logout</span>
              </button>
            </div>
          </div>
        </div>

        {/* Mobile Navigation Menu */}
        <div className={`lg:hidden ${isNavOpen ? 'block' : 'hidden'}`}>
          <div className="px-2 pt-2 pb-3 space-y-1">
            <a href="#" className="flex items-center space-x-1 px-3 py-2 rounded-md bg-purple-600 bg-opacity-50">
              <Home className="h-4 w-4" />
              <span>Dashboard</span>
            </a>
            <a href="#" className="flex items-center space-x-1 px-3 py-2 rounded-md hover:bg-purple-600 hover:bg-opacity-25">
              <Settings className="h-4 w-4" />
              <span>Settings</span>
            </a>
            <button className="w-full flex items-center space-x-1 px-3 py-2 rounded-md hover:bg-purple-600 hover:bg-opacity-25 text-purple-200 hover:text-white">
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
        <div className="flex flex-col md:flex-row justify-between items-stretch gap-6 mb-8">
          <div className="flex-1 bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
            <div className="text-purple-600 text-lg font-semibold">New Users This Week</div>
            <div className="text-3xl font-bold">47</div>
          </div>
          <div className="flex-1 bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
            <div className="text-purple-600 text-lg font-semibold">Pending Reviews</div>
            <div className="text-3xl font-bold">{submissions.length}</div>
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
          <div className="mt-8 flex flex-col md:flex-row items-center justify-between gap-4">
            <div className="text-sm text-gray-600">
              Showing {(currentPage - 1) * itemsPerPage + 1} to {Math.min(currentPage * itemsPerPage, submissions.length)} of {submissions.length} submissions
            </div>
            <div className="flex items-center space-x-2">
              <button
                onClick={handlePrevPage}
                disabled={currentPage === 1}
                className={`p-2 rounded-md ${currentPage === 1 ? 'text-gray-400 cursor-not-allowed' : 'text-purple-600 hover:bg-purple-50'}`}
              >
                <ChevronLeft className="w-5 h-5" />
              </button>
              {Array.from({ length: pageCount }, (_, i) => (
                <button
                  key={i + 1}
                  onClick={() => handlePageChange(i + 1)}
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
                onClick={handleNextPage}
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

export default AdminDashboard;
