// Dashboard.jsx
import React, { useState, useCallback, useEffect } from "react";
import {
  Coins,
  Home,
  Settings,
  LogOut,
  Clock,
  CheckCircle,
  XCircle,
  Flag,
  ChevronLeft,
  ChevronRight,
  AlertTriangle,
  Menu,
  X,
} from "lucide-react";
import { Navigate, Link } from "react-router-dom";
import { useRecoilState } from "recoil";

import { currentUserState } from "../../AppState";
import {
  getComicSubmissionListAPI,
  getComicSubmissionsCountByFilterAPI,
  getComicSubmissionsTotalCoinsAwardedAPI,
  postComicSubmissionJudgementOperationAPI,
} from "../../API/ComicSubmission";
import { getUsersCountJoinedThisWeekAPI } from "../../API/user";
import { getFaucetBalanceAPI } from "../../API/Faucet";
import AdminTopbar from "../Navigation/AdminTopbar";
import GalleryItem from './GalleryItem';

const AdminDashboard = () => {
  // Global state
  const [currentUser] = useRecoilState(currentUserState);

  // Data states
  const [totalPendingSubmissions, setTotalPendingSubmissions] = useState(0);
  const [totalCoinsAwarded, setTotalCoinsAwarded] = useState(0);
  const [totalUsersJoinedThisWeek, setTotalUsersJoinedThisWeek] = useState(0);
  const [faucetBalance, setFaucetBalance] = useState(0);
  const [pendingSubmissions, setPendingSubmissions] = useState([]);
  const [isFetching, setFetching] = useState(false);
  const [errors, setErrors] = useState({});
  const [isNavOpen, setIsNavOpen] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);

  const itemsPerPage = 8;
  const pageCount = Math.ceil(pendingSubmissions.length / itemsPerPage);
  const currentSubmissions = pendingSubmissions.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage,
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
        await getComicSubmissionsTotalCoinsAwardedAPI(
          (resp) => {
            if (mounted) {
              console.log("getComicSubmissionListAPI: Success", resp);
              setTotalCoinsAwarded(resp.count);
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
          },
        );

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
          },
        );

        await getUsersCountJoinedThisWeekAPI(
          (resp) => {
            if (mounted) {
              console.log("getUsersCountJoinedThisWeekAPI: Success", resp);
              setTotalUsersJoinedThisWeek(resp.count);
            }
          },
          (apiErr) => {
            if (mounted) {
              console.log("getUsersCountJoinedThisWeekAPI: Error:", apiErr);
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
          },
        );

        await getFaucetBalanceAPI(
          (resp) => {
            if (mounted) {
              console.log("getFaucetBalanceAPI: Success", resp);
              setFaucetBalance(resp.count);
            }
          },
          (apiErr) => {
            if (mounted) {
              console.log("getFaucetBalanceAPI: Error:", apiErr);
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
          },
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
              console.log(
                "getComicSubmissionsCountByFilterAPI: Error:",
                apiErr,
              );
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
          },
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
        status: 3, // 3 is the status code for "approved"
        judgement_notes: "Approved by administrator",
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
            () => (window.location.href = "/login?unauthorized=true"),
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
        () => (window.location.href = "/login?unauthorized=true"),
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
        status: 2, // 2 is the status code for "rejected"
        judgement_notes: "Approved by administrator",
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
            () => (window.location.href = "/login?unauthorized=true"),
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
        () => (window.location.href = "/login?unauthorized=true"),
      );
    } catch (error) {
      console.error("Failed to reject submission:", error);
    }
  }, []);

  const handleFlagSubmission = useCallback(async (submissionId, flagData) => {
    try {
      // Here you would call your flag API endpoint
      console.log(`Flagging submission ${submissionId} for:`, flagData);

      // Show we are processing
      setFetching(true);

      // Prepare request body for the approval
      const submissionReq = {
        comic_submission_id: submissionId,
        status: 6, // 6 is the status code for "flagged"
        flag_issue: flagData.flagIssue,
        flag_issue_other:
          flagData.flagIssue === "other" ? flagData.flagIssueOther : "",
        flag_action: flagData.flagAction,
      };

      await postComicSubmissionJudgementOperationAPI(
        submissionReq,
        // onSuccess callback
        async (resp) => {
          console.log("Successfully flagged submission:", submissionId);

          // Refresh the submissions list
          const params = new Map();
          params.set("status", 1); // Get pending submissions
          await getComicSubmissionListAPI(
            params,
            (resp) => setPendingSubmissions(resp.submissions),
            (err) => setErrors(err),
            () => setFetching(false),
            () => (window.location.href = "/login?unauthorized=true"),
          );
        },
        // onError callback
        (apiErr) => {
          console.error("Failed to flag submission:", apiErr);
          setErrors(apiErr);
          setFetching(false);
        },
        // onFinally callback
        () => setFetching(false),
        // onUnauthorized callback
        () => (window.location.href = "/login?unauthorized=true"),
      );
    } catch (error) {
      console.error("Failed to flag submission:", error);
      setFetching(false);
    }
  }, []);

  const handlePageChange = useCallback(
    (newPage) => {
      if (newPage >= 1 && newPage <= pageCount) {
        setCurrentPage(newPage);
      }
    },
    [pageCount],
  );

  const handleNextPage = useCallback(() => {
    handlePageChange(currentPage + 1);
  }, [currentPage, handlePageChange]);

  const handlePrevPage = useCallback(() => {
    handlePageChange(currentPage - 1);
  }, [currentPage, handlePageChange]);


  if (isFetching) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-xl text-purple-600">Loading submissions...</div>
      </div>
    );
  }

  if (Object.keys(errors).length > 0) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-xl text-red-600">Error loading submissions</div>
      </div>
    );
  }

  return (
  <div className="min-h-screen bg-purple-50">
    <AdminTopbar currentPage="Dashboard" />

    {/* Main Content with proper horizontal spacing */}
    <main className="max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <h1
        className="text-3xl font-bold text-purple-800 mb-8"
        style={{ fontFamily: "Comic Sans MS, cursive" }}
      >
        Admin Dashboard
      </h1>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div className="bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
          <div className="text-purple-600 text-lg font-semibold">
            New Users This Week
          </div>
          <div className="text-3xl font-bold">{totalUsersJoinedThisWeek}</div>
        </div>
        <div className="bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
          <div className="text-purple-600 text-lg font-semibold">
            Pending Reviews
          </div>
          <div className="text-3xl font-bold">
            {pendingSubmissions.length}
          </div>
        </div>
        <div className="bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
          <div className="text-purple-600 text-lg font-semibold">
            Total ComicCoins Paid
          </div>
          <div className="text-3xl font-bold">
            {totalCoinsAwarded}&nbsp;CC
          </div>
        </div>
        <div className="bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
          <div className="text-purple-600 text-lg font-semibold">
            Faucet Balance
          </div>
          <div className="text-3xl font-bold">{faucetBalance}&nbsp;CC</div>
        </div>
      </div>

      {/* Submissions Section */}
      <div className="bg-white rounded-xl shadow-lg p-6 mb-8 border-2 border-purple-200">
        <h2
          className="text-2xl font-bold text-purple-800 mb-6"
          style={{ fontFamily: "Comic Sans MS, cursive" }}
        >
          Submissions Awaiting Review
        </h2>

        {pendingSubmissions.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-12 text-gray-500">
            <Coins className="w-12 h-12 mb-4 text-purple-300" />
            <p className="text-lg font-medium mb-2">No Pending Reviews</p>
            <p className="text-sm text-gray-400">
              There are currently no comic submissions waiting for review.
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6 gap-6">
            {currentSubmissions.map((submission) => (
              <div className="max-w-[300px] mx-auto w-full">
                <GalleryItem
                  key={submission.id}
                  submission={submission}
                  onFlag={handleFlagSubmission}
                  handleApproveSubmission={handleApproveSubmission}
                  handleRejectSubmission={handleRejectSubmission}
                />
              </div>
            ))}
          </div>
        )}

        {/* Pagination */}
        <div className="mt-8 flex flex-col md:flex-row items-center justify-between gap-4">
          <div className="text-sm text-gray-600">
            {pendingSubmissions.length === 0
              ? "No submissions to display"
              : `Showing ${(currentPage - 1) * itemsPerPage + 1} to ${Math.min(
                  currentPage * itemsPerPage,
                  pendingSubmissions.length,
                )} of ${pendingSubmissions.length} submissions`}
          </div>
          {pendingSubmissions.length > 0 && pageCount > 1 && (
            <div className="flex items-center space-x-2">
              <button
                onClick={handlePrevPage}
                disabled={currentPage === 1}
                className="p-2 rounded-lg border border-purple-200 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-purple-50"
              >
                <ChevronLeft className="w-5 h-5 text-purple-600" />
              </button>
              <span className="text-sm text-gray-600">
                Page {currentPage} of {pageCount}
              </span>
              <button
                onClick={handleNextPage}
                disabled={currentPage === pageCount}
                className="p-2 rounded-lg border border-purple-200 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-purple-50"
              >
                <ChevronRight className="w-5 h-5 text-purple-600" />
              </button>
            </div>
          )}
        </div>
      </div>
    </main>
  </div>
);
}

export default AdminDashboard;
