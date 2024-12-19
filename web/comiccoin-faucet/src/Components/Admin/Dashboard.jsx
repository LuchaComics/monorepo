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

  const GalleryItem = ({ submission, onFlag }) => {
    const [showBackCover, setShowBackCover] = useState(false);
    const [showFlagModal, setShowFlagModal] = useState(false);

    const toggleCover = () => setShowBackCover((prev) => !prev);
    const toggleFlagModal = () => setShowFlagModal((prev) => !prev);

    return (
      <div className="w-64 bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow border border-purple-100">
        <div className="relative w-full h-80">
          <img
            src={
              showBackCover
                ? submission.backCover.objectUrl
                : submission.frontCover.objectUrl
            }
            alt={`${submission.name} - ${showBackCover ? "Back" : "Front"} Cover`}
            className="w-full h-full object-cover rounded-t-lg"
          />
          <div className="absolute top-2 left-2 right-2 flex justify-between">
            <button
              onClick={toggleCover}
              className="bg-white rounded-md px-2 py-1 text-xs font-medium shadow hover:bg-gray-50"
            >
              {showBackCover ? "View Front" : "View Back"}
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
                onClick={toggleFlagModal}
                className="bg-white rounded-full p-2 shadow hover:bg-yellow-50"
                title="Flag for Review"
              >
                <Flag
                  className={`w-5 h-5 ${submission.flagReason ? "text-yellow-500" : "text-gray-400"}`}
                />
              </button>
            </div>
          </div>

          {showFlagModal && (
            <FlagModal
              isOpen={showFlagModal}
              onClose={toggleFlagModal}
              onSubmit={onFlag}
              submissionId={submission.id}
            />
          )}
        </div>

        <div className="p-3">
          <h3 className="font-medium text-sm truncate" title={submission.name}>
            {submission.name}
          </h3>
          <p className="text-xs text-gray-600 truncate">
            by {submission.submitter}
          </p>
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

      <main className="p-8">
        <h1
          className="text-3xl font-bold text-purple-800 mb-8"
          style={{ fontFamily: "Comic Sans MS, cursive" }}
        >
          Admin Dashboard
        </h1>

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

        <div className="bg-white rounded-xl shadow-lg p-6 mb-8 border-2 border-purple-200">
          <h2
            className="text-2xl font-bold text-purple-800 mb-6"
            style={{ fontFamily: "Comic Sans MS, cursive" }}
          >
            Submissions Awaiting Review
          </h2>
          <div className="flex flex-wrap gap-6">
            {currentSubmissions.map((submission) => (
              <GalleryItem
                key={submission.id}
                submission={submission}
                onFlag={handleFlagSubmission} // Add this line
              />
            ))}
          </div>

          <div className="mt-8 flex flex-col md:flex-row items-center justify-between gap-4">
            <div className="text-sm text-gray-600">
              Showing {(currentPage - 1) * itemsPerPage + 1} to{" "}
              {Math.min(currentPage * itemsPerPage, pendingSubmissions.length)}{" "}
              of {pendingSubmissions.length} submissions
            </div>
            <div className="flex items-center space-x-2">
              <button
                onClick={handlePrevPage}
                disabled={currentPage === 1}
                className={`p-2 rounded-md ${currentPage === 1 ? "text-gray-400 cursor-not-allowed" : "text-purple-600 hover:bg-purple-50"}`}
              >
                <ChevronLeft className="w-5 h-5" />
              </button>
              {Array.from({ length: pageCount }, (_, i) => (
                <button
                  key={i + 1}
                  onClick={() => handlePageChange(i + 1)}
                  className={`px-3 py-1 rounded-md ${
                    currentPage === i + 1
                      ? "bg-purple-600 text-white"
                      : "text-purple-600 hover:bg-purple-50"
                  }`}
                >
                  {i + 1}
                </button>
              ))}
              <button
                onClick={handleNextPage}
                disabled={currentPage === pageCount}
                className={`p-2 rounded-md ${currentPage === pageCount ? "text-gray-400 cursor-not-allowed" : "text-purple-600 hover:bg-purple-50"}`}
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

const FlagModal = ({ isOpen, onClose, onSubmit, submissionId }) => {
  const [flagIssue, setFlagIssue] = useState("");
  const [flagAction, setFlagAction] = useState("");

  if (!isOpen) return null;

  const handleSubmit = () => {
    onSubmit(submissionId, { flagIssue, flagAction });
    onClose();
    setFlagIssue("");
    setFlagAction("");
  };

  // Values need to be exactly as is in the backends `domain/comicsubmission.go` file.
  const flagIssueOptions = [
    { value: 2, label: "Duplicate submission" },
    { value: 3, label: "Poor image quality" },
    { value: 4, label: "Counterfeit" },
    { value: 5, label: "Inappropriate Content" },
    { value: 1, label: "Other" },
  ];

  const flagActionOptions = [
    { value: 1, label: "Do nothing" },
    { value: 2, label: "Lockout User" },
    { value: 3, label: "Lockout User and Ban IP Address" },
  ];

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center">
      <div className="bg-white rounded-lg max-w-md w-full mx-4">
        <div className="p-6">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-xl font-semibold">Flag Content</h2>
            <button
              onClick={onClose}
              className="text-gray-500 hover:text-gray-700"
            >
              <X className="w-5 h-5" />
            </button>
          </div>

          <p className="text-red-500 text-sm mb-6">
            Warning: Once submission is flagged, uploads are deleted
          </p>

          <div className="space-y-6">
            <div>
              <h3 className="font-medium mb-3">Flag Issue</h3>
              <div className="space-y-2">
                {flagIssueOptions.map((option) => (
                  <div
                    key={option.value}
                    className="flex items-center space-x-2"
                  >
                    <input
                      type="radio"
                      id={`issue-${option.value}`}
                      name="flagIssue"
                      value={parseInt(option.value)}
                      checked={flagIssue === option.value}
                      onChange={(e) => setFlagIssue(parseInt(e.target.value))}
                      className="text-purple-600 focus:ring-purple-500"
                    />
                    <label
                      htmlFor={`issue-${option.value}`}
                      className="text-sm text-gray-700 cursor-pointer"
                    >
                      {option.label}
                    </label>
                  </div>
                ))}
              </div>
            </div>

            <div>
              <h3 className="font-medium mb-3">Flag Action</h3>
              <div className="space-y-2">
                {flagActionOptions.map((option) => (
                  <div
                    key={option.value}
                    className="flex items-center space-x-2"
                  >
                    <input
                      type="radio"
                      id={`action-${option.value}`}
                      name="flagAction"
                      value={parseInt(option.value)}
                      checked={flagAction === option.value}
                      onChange={(e) => setFlagAction(parseInt(e.target.value))}
                      className="text-purple-600 focus:ring-purple-500"
                    />
                    <label
                      htmlFor={`action-${option.value}`}
                      className="text-sm text-gray-700 cursor-pointer"
                    >
                      {option.label}
                    </label>
                  </div>
                ))}
              </div>
            </div>
          </div>

          <div className="mt-8 flex justify-end space-x-3">
            <button
              onClick={onClose}
              className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-purple-500"
            >
              Cancel
            </button>
            <button
              onClick={handleSubmit}
              disabled={!flagIssue || !flagAction}
              className={`px-4 py-2 text-sm font-medium text-white rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500
                  ${
                    !flagIssue || !flagAction
                      ? "bg-red-300 cursor-not-allowed"
                      : "bg-red-600 hover:bg-red-700"
                  }`}
            >
              Submit Flag
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
