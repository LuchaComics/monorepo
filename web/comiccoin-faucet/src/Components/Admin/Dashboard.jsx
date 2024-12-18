import React, { useState, useCallback, useEffect } from 'react';
import {
  Coins, Home, Settings, LogOut, Clock,
  CheckCircle, XCircle, Flag, ChevronLeft,
  ChevronRight, AlertTriangle, Menu, X
} from 'lucide-react';
import { Navigate, Link } from "react-router-dom";
import { useRecoilState } from "recoil";

import { currentUserState } from "../../AppState";
import {
    getComicSubmissionListAPI,
    getComicSubmissionsCountByFilterAPI,
    postComicSubmissionJudgementOperationAPI
} from "../../API/ComicSubmission";

const AdminDashboard = () => {
  // Global state
  const [currentUser] = useRecoilState(currentUserState);

  // Data states
  const [totalPendingSubmissions, setTotalPendingSubmissions] = useState(0);
  const [pendingSubmissions, setPendingSubmissions] = useState([]);
  const [isFetching, setFetching] = useState(false);
  const [errors, setErrors] = useState({});
  const [isNavOpen, setIsNavOpen] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);

  const itemsPerPage = 8;
  const pageCount = Math.ceil(pendingSubmissions.length / itemsPerPage);
  const currentSubmissions = pendingSubmissions.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

  useEffect(() => {
    let mounted = true;

    const fetchSubmissions = async () => {
      if (!mounted) return;

      setFetching(true);
      const params = new Map();
      params.set("status", 1); // ComicSubmissionStatusInReview
      params.set("page_size", itemsPerPage);
      params.set("page", currentPage);

      try {
        await getComicSubmissionListAPI(
          params,
          (resp) => {
            if (mounted) {
              console.log("getComicSubmissionListAPI: Success", resp);
              setPendingSubmissions(resp.submissions);
            }
          },
          (apiErr) => {
            if (mounted) {
              console.log("getComicSubmissionListAPI: Error:", apiErr);
              setErrors(apiErr);
            }
          },
          () => {
            if (mounted) {
              setFetching(false);
            }
          },
          () => {
            if (mounted) {
              window.location.href = "/login?unauthorized=true";
            }
          }
        );
      } catch (error) {
        console.error("Failed to fetch submissions:", error);
      }
    };

    fetchSubmissions();

    const fetchTotalPendingSubmissions = async () => {
      if (!mounted) return;

      setFetching(true);
      const params = new Map();
      params.set("status", 1); // ComicSubmissionStatusInReview

      try {
        await getComicSubmissionsCountByFilterAPI(
          params,
          (resp) => {
            if (mounted) {
              console.log("getComicSubmissionsCountByFilterAPI: Success", resp);
              setTotalPendingSubmissions(resp.submissions);
            }
          },
          (apiErr) => {
            if (mounted) {
              console.log("getComicSubmissionsCountByFilterAPI: Error:", apiErr);
              setErrors(apiErr);
              setTotalPendingSubmissions(0);
            }
          },
          () => {
            if (mounted) {
              setFetching(false);
            }
          },
          () => {
            if (mounted) {
              window.location.href = "/login?unauthorized=true";
            }
          }
        );
      } catch (error) {
        console.error("Failed to fetch total count submissions:", error);
      }
    };

    fetchTotalPendingSubmissions();

    return () => {
      mounted = false;
    };
  }, [currentPage, currentUser]);

  const handleApproveSubmission = useCallback(async (submissionId) => {
      try {
        // Show we are processing
        setFetching(true);

        // Prepare request body for the approval
        const submissionReq = {
          comic_submission_id: submissionId,
          status: 3,  // 3 is the status code for "approved"
          judgement_notes: "Approved by administrator"
        };

        await postComicSubmissionJudgementOperationAPI(
          submissionReq,
          // onSuccess callback
          async (resp) => {
            console.log("Successfully approved submission:", submissionId);

            // Refresh the submissions list
            const params = new Map();
            params.set("status", 1); // Get pending submissions
            await getComicSubmissionListAPI(
              params,
              (resp) => setPendingSubmissions(resp.submissions),
              (err) => setErrors(err),
              () => setFetching(false),
              () => window.location.href = "/login?unauthorized=true"
            );
          },
          // onError callback
          (apiErr) => {
            console.error("Failed to approve submission:", apiErr);
            setErrors(apiErr);
            setFetching(false);
          },
          // onFinally callback
          () => setFetching(false),
          // onUnauthorized callback
          () => window.location.href = "/login?unauthorized=true"
        );

      } catch (error) {
        console.error("Error in handleApproveSubmission:", error);
        setErrors(error);
        setFetching(false);
      }
  }, []);

  const handleRejectSubmission = useCallback(async (submissionId) => {
    try {
      // Here you would call your reject API endpoint
      console.log(`Rejecting submission ${submissionId}`);

      // Show we are processing
      setFetching(true);

      // Prepare request body for the approval
      const submissionReq = {
        comic_submission_id: submissionId,
        status: 2,  // 2 is the status code for "rejected"
        judgement_notes: "Approved by administrator"
      };

      await postComicSubmissionJudgementOperationAPI(
        submissionReq,
        // onSuccess callback
        async (resp) => {
          console.log("Successfully approved submission:", submissionId);

          // Refresh the submissions list
          const params = new Map();
          params.set("status", 1); // Get pending submissions
          await getComicSubmissionListAPI(
            params,
            (resp) => setPendingSubmissions(resp.submissions),
            (err) => setErrors(err),
            () => setFetching(false),
            () => window.location.href = "/login?unauthorized=true"
          );
        },
        // onError callback
        (apiErr) => {
          console.error("Failed to approve submission:", apiErr);
          setErrors(apiErr);
          setFetching(false);
        },
        // onFinally callback
        () => setFetching(false),
        // onUnauthorized callback
        () => window.location.href = "/login?unauthorized=true"
      );

    } catch (error) {
      console.error("Failed to reject submission:", error);
    }
  }, []);

  const handleFlagSubmission = useCallback(async (submissionId, flagReason) => {
    try {
      // Here you would call your flag API endpoint
      console.log(`Flagging submission ${submissionId} for: ${flagReason}`);

      // After successful API call, refresh the submissions list
      const params = new Map();
      params.set("status", 1);
      await getComicSubmissionListAPI(
        params,
        (resp) => setPendingSubmissions(resp.submissions),
        setErrors,
        () => setFetching(false),
        () => window.location.href = "/login?unauthorized=true"
      );
    } catch (error) {
      console.error("Failed to flag submission:", error);
    }
  }, []);

  const handlePageChange = useCallback((newPage) => {
    if (newPage >= 1 && newPage <= pageCount) {
      setCurrentPage(newPage);
    }
  }, [pageCount]);

  const handleNextPage = useCallback(() => {
    handlePageChange(currentPage + 1);
  }, [currentPage, handlePageChange]);

  const handlePrevPage = useCallback(() => {
    handlePageChange(currentPage - 1);
  }, [currentPage, handlePageChange]);

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

  const GalleryItem = ({ submission }) => {
    const [showBackCover, setShowBackCover] = useState(false);
    const [showFlagMenu, setShowFlagMenu] = useState(false);

    const toggleCover = () => setShowBackCover(prev => !prev);
    const toggleFlagMenu = () => setShowFlagMenu(prev => !prev);

    return (
      <div className="w-64 bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow border border-purple-100">
        <div className="relative w-full h-80">
          <img
            src={showBackCover ? submission.backCover.objectUrl : submission.frontCover.objectUrl}
            alt={`${submission.name} - ${showBackCover ? 'Back' : 'Front'} Cover`}
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
          <h3 className="font-medium text-sm truncate" title={submission.name}>
            {submission.name}
          </h3>
          <p className="text-xs text-gray-600 truncate">by {submission.submitter}</p>
          <p className="text-xs text-gray-500 mt-1">
            {new Date(submission.createdAt).toLocaleDateString()}
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

  if (isFetching) {
    return <div className="flex items-center justify-center min-h-screen">
      <div className="text-xl text-purple-600">Loading submissions...</div>
    </div>;
  }

  if (Object.keys(errors).length > 0) {
    return <div className="flex items-center justify-center min-h-screen">
      <div className="text-xl text-red-600">Error loading submissions</div>
    </div>;
  }

  return (
    <div className="min-h-screen bg-purple-50">
      <nav className="bg-gradient-to-r from-purple-700 to-indigo-800 text-white shadow-lg">
        <div className="max-w-full px-4">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center space-x-2">
              <Coins className="h-8 w-8" />
              <span className="text-xl font-bold" style={{fontFamily: 'Comic Sans MS, cursive'}}>
                ComicCoin Admin
              </span>
            </div>

            <div className="flex lg:hidden">
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

            <div className="hidden lg:flex lg:items-center lg:space-x-4">
              <Link to="/admin/dashboard" className="flex items-center space-x-1 px-3 py-2 rounded-md bg-purple-600 bg-opacity-50">
                <Home className="h-4 w-4" />
                <span>Dashboard</span>
              </Link>
              <Link to="/admin/settings" className="flex items-center space-x-1 px-3 py-2 rounded-md hover:bg-purple-600 hover:bg-opacity-25">
                <Settings className="h-4 w-4" />
                <span>Settings</span>
              </Link>
              <button className="flex items-center space-x-1 px-3 py-2 rounded-md hover:bg-purple-600 hover:bg-opacity-25 text-purple-200 hover:text-white">
                <LogOut className="h-4 w-4" />
                <span>Logout</span>
              </button>
            </div>
          </div>
        </div>

        <div className={`lg:hidden ${isNavOpen ? 'block' : 'hidden'}`}>
          <div className="px-2 pt-2 pb-3 space-y-1">
            <Link to="/admin/dashboard" className="flex items-center space-x-1 px-3 py-2 rounded-md bg-purple-600 bg-opacity-50">
              <Home className="h-4 w-4" />
              <span>Dashboard</span>
            </Link>
            <Link to="/admin/settings" className="flex items-center space-x-1 px-3 py-2 rounded-md hover:bg-purple-600 hover:bg-opacity-25">
              <Settings className="h-4 w-4" />
              <span>Settings</span>
            </Link>
            <button className="w-full flex items-center space-x-1 px-3 py-2 rounded-md hover:bg-purple-600 hover:bg-opacity-25 text-purple-200 hover:text-white">
              <LogOut className="h-4 w-4" />
              <span>Logout</span>
            </button>
          </div>
        </div>
      </nav>

      <main className="p-8">
        <h1 className="text-3xl font-bold text-purple-800 mb-8" style={{fontFamily: 'Comic Sans MS, cursive'}}>
          Admin Dashboard
        </h1>

        <div className="flex flex-col md:flex-row justify-between items-stretch gap-6 mb-8">
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

        <div className="bg-white rounded-xl shadow-lg p-6 mb-8 border-2 border-purple-200">
          <h2 className="text-2xl font-bold text-purple-800 mb-6" style={{fontFamily: 'Comic Sans MS, cursive'}}>
            Submissions Awaiting Review
          </h2>
          <div className="flex flex-wrap gap-6">
            {currentSubmissions.map(submission => (
                <GalleryItem key={submission.id} submission={submission} />
              ))}
            </div>

            <div className="mt-8 flex flex-col md:flex-row items-center justify-between gap-4">
              <div className="text-sm text-gray-600">
                Showing {(currentPage - 1) * itemsPerPage + 1} to {Math.min(currentPage * itemsPerPage, pendingSubmissions.length)} of {pendingSubmissions.length} submissions
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
