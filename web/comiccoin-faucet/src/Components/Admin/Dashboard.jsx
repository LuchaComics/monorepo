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
  Users,
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
import {
  getUsersCountJoinedThisWeekAPI,
  getUserListAPI,
  postUserProfileVerificationJudgeOperationAPI
} from "../../API/user";
import {
  getFaucetBalanceAPI
} from "../../API/Faucet";
import AdminTopbar from "../Navigation/AdminTopbar";
import GalleryItem from './GalleryItem';
import UserDetailModal from "./UserDetailModal";

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

  // User list states with pagination
  const [selectedUser, setSelectedUser] = useState(null);
  const [users, setUsers] = useState([]);
  const [totalUsers, setTotalUsers] = useState(0);
  const [currentUserPage, setCurrentUserPage] = useState(1);
  const [currentSubmissionPage, setCurrentSubmissionPage] = useState(1);

  const itemsPerPage = 8; // For submissions
  const usersPerPage = 10; // For users list

  useEffect(() => {
    let mounted = true;

    const fetchInitialData = async () => {
      if (!mounted) return;

      setFetching(true);
      try {
        // Get total coins awarded
        await getComicSubmissionsTotalCoinsAwardedAPI(
          (resp) => {
            if (mounted) {
              console.log("getComicSubmissionsTotalCoinsAwardedAPI: Success", resp);
              setTotalCoinsAwarded(resp.count);
            }
          },
          (apiErr) => {
            if (mounted) {
              console.log("getComicSubmissionsTotalCoinsAwardedAPI: Error:", apiErr);
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

        // Get users joined this week
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
          }
        );

        // Get faucet balance
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
          }
        );
      } catch (error) {
        console.error("Failed to fetch initial data:", error);
        setFetching(false);
      }
    };

    const fetchPaginatedData = async () => {
      if (!mounted) return;

      setFetching(true);
      try {
        // Fetch users with pagination
        const userParams = new Map();
        userParams.set("page_size", usersPerPage);
        userParams.set("page", currentUserPage);
        userParams.set("profile_verification_status", 2);

        await getUserListAPI(
          userParams,
          (resp) => {
            if (mounted) {
              console.log("getUserListAPI: Success", resp);
              if (resp) {
                  setUsers(resp.users);
                  setTotalUsers(resp.users.length);
              }
            }
          },
          (apiErr) => {
            if (mounted) {
              console.log("getUserListAPI: Error:", apiErr);
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

        // Fetch submissions count
        const countParams = new Map();
        countParams.set("status", 1); // ComicSubmissionStatusInReview
        await getComicSubmissionsCountByFilterAPI(
          countParams,
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

        // Fetch paginated submissions
        const submissionParams = new Map();
        submissionParams.set("status", 1); // ComicSubmissionStatusInReview
        submissionParams.set("page_size", itemsPerPage);
        submissionParams.set("page", currentSubmissionPage);

        await getComicSubmissionListAPI(
          submissionParams,
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
        console.error("Failed to fetch paginated data:", error);
        setFetching(false);
      }
    };

    fetchInitialData();
    fetchPaginatedData();

    return () => {
      mounted = false;
    };
  }, [currentUserPage, currentSubmissionPage]);


  const handleUserClick = useCallback((user) => {
    setSelectedUser(user);
  }, []);

  const handleCloseModal = useCallback(() => {
    setSelectedUser(null);
  }, []);

  const handleAcceptUser = useCallback(async () => {
    console.log("handleAcceptUser: Beginning...");
    console.log("handleAcceptUser: selectedUser:", selectedUser);
    console.log("handleAcceptUser: status:", 3); // Approved
    try {
        setFetching(true);
        const judgeReq = {
            user_id: selectedUser.id,
            profile_verification_status: 3, // Approved
        };

        await postUserProfileVerificationJudgeOperationAPI(
          judgeReq,
          async (resp) => {
            console.log("Successfully approved user profile:", selectedUser.id);

            // Refresh submissions list
            const userParams = new Map();
            userParams.set("page_size", usersPerPage);
            userParams.set("page", currentUserPage);
            userParams.set("profile_verification_status", 2); // Submitted for review.

            await getUserListAPI(
              userParams,
              (resp) => {
                  setUsers(resp.users);
                  setTotalUsers(resp.users.length);
              },
              (err) => {
                  console.log("getUserListAPI: Error:", err);
                  setErrors(err);
              },
              () => setFetching(false),
              () => (window.location.href = "/login?unauthorized=true")
            );
          },
          (apiErr) => {
            console.error("Failed to approve submission:", apiErr);
            setErrors(apiErr);
            setFetching(false);
          },
          () => setFetching(false),
          () => (window.location.href = "/login?unauthorized=true")
        );

        //----
    } catch (error) {
      console.error("Error in handleApproveSubmission:", error);
      setErrors(error);
      setFetching(false);
    }
    setSelectedUser(null);
  }, [selectedUser, currentUserPage, usersPerPage]);

  const handleRejectUser = useCallback(async () => {
    console.log("handleRejectUser: Beginning...");
    console.log("handleRejectUser: selectedUser:", selectedUser);
    console.log("handleRejectUser: status:", 4); // Rejected
    try {
        setFetching(true);
        const judgeReq = {
            user_id: selectedUser.id,
            profile_verification_status: 4, // Rejected
        };

        await postUserProfileVerificationJudgeOperationAPI(
          judgeReq,
          async (resp) => {
            console.log("Successfully approved user profile:", selectedUser.id);

            // Refresh submissions list
            const userParams = new Map();
            userParams.set("page_size", usersPerPage);
            userParams.set("page", currentUserPage);
            userParams.set("profile_verification_status", 2); // Submitted for review.

            await getUserListAPI(
              userParams,
              (resp) => {
                  setUsers(resp.users);
                  setTotalUsers(resp.users.length);
              },
              (err) => {
                  console.log("getUserListAPI: Error:", err);
                  setErrors(err);
              },
              () => setFetching(false),
              () => (window.location.href = "/login?unauthorized=true")
            );
          },
          (apiErr) => {
            console.error("Failed to approve submission:", apiErr);
            setErrors(apiErr);
            setFetching(false);
          },
          () => setFetching(false),
          () => (window.location.href = "/login?unauthorized=true")
        );

        //----
    } catch (error) {
      console.error("Error in handleApproveSubmission:", error);
      setErrors(error);
      setFetching(false);
    }
    setSelectedUser(null);
  }, [selectedUser, currentUserPage, usersPerPage]);

  const handleApproveSubmission = useCallback(async (submissionId) => {
    try {
      setFetching(true);
      const submissionReq = {
        comic_submission_id: submissionId,
        status: 3, // Approved
        judgement_notes: "Approved by administrator",
      };

      await postComicSubmissionJudgementOperationAPI(
        submissionReq,
        async (resp) => {
          console.log("Successfully approved submission:", submissionId);

          // Refresh submissions list
          const params = new Map();
          params.set("status", 1);
          params.set("page_size", itemsPerPage);
          params.set("page", currentSubmissionPage);

          await getComicSubmissionListAPI(
            params,
            (resp) => setPendingSubmissions(resp.submissions),
            (err) => setErrors(err),
            () => setFetching(false),
            () => (window.location.href = "/login?unauthorized=true")
          );
        },
        (apiErr) => {
          console.error("Failed to approve submission:", apiErr);
          setErrors(apiErr);
          setFetching(false);
        },
        () => setFetching(false),
        () => (window.location.href = "/login?unauthorized=true")
      );
    } catch (error) {
      console.error("Error in handleApproveSubmission:", error);
      setErrors(error);
      setFetching(false);
    }
  }, [currentSubmissionPage]);

  const handleRejectSubmission = useCallback(async (submissionId) => {
    try {
      setFetching(true);
      const submissionReq = {
        comic_submission_id: submissionId,
        status: 2, // Rejected
        judgement_notes: "Rejected by administrator",
      };

      await postComicSubmissionJudgementOperationAPI(
        submissionReq,
        async (resp) => {
          console.log("Successfully rejected submission:", submissionId);

          // Refresh submissions list
          const params = new Map();
          params.set("status", 1);
          params.set("page_size", itemsPerPage);
          params.set("page", currentSubmissionPage);

          await getComicSubmissionListAPI(
            params,
            (resp) => setPendingSubmissions(resp.submissions),
            (err) => setErrors(err),
            () => setFetching(false),
            () => (window.location.href = "/login?unauthorized=true")
          );
        },
        (apiErr) => {
          console.error("Failed to reject submission:", apiErr);
          setErrors(apiErr);
          setFetching(false);
        },
        () => setFetching(false),
        () => (window.location.href = "/login?unauthorized=true")
      );
    } catch (error) {
      console.error("Error in handleRejectSubmission:", error);
      setErrors(error);
      setFetching(false);
    }
  }, [currentSubmissionPage]);

  const handleFlagSubmission = useCallback(async (submissionId, flagData) => {
    try {
      setFetching(true);
      const submissionReq = {
        comic_submission_id: submissionId,
        status: 6, // Flagged
        flag_issue: flagData.flagIssue,
        flag_issue_other: flagData.flagIssue === "other" ? flagData.flagIssueOther : "",
        flag_action: flagData.flagAction,
      };

      await postComicSubmissionJudgementOperationAPI(
        submissionReq,
        async (resp) => {
          console.log("Successfully flagged submission:", submissionId);

          // Refresh submissions list
          const params = new Map();
          params.set("status", 1);
          params.set("page_size", itemsPerPage);
          params.set("page", currentSubmissionPage);

          await getComicSubmissionListAPI(
            params,
            (resp) => setPendingSubmissions(resp.submissions),
            (err) => setErrors(err),
            () => setFetching(false),
            () => (window.location.href = "/login?unauthorized=true")
          );
        },
        (apiErr) => {
          console.error("Failed to flag submission:", apiErr);
          setErrors(apiErr);
          setFetching(false);
        },
        () => setFetching(false),
        () => (window.location.href = "/login?unauthorized=true")
      );
    } catch (error) {
      console.error("Error in handleFlagSubmission:", error);
      setFetching(false);
    }
  }, [currentSubmissionPage]);

  const handleUserPageChange = useCallback(
    (newPage) => {
      const maxPages = Math.ceil(totalUsers / usersPerPage);
      if (newPage >= 1 && newPage <= maxPages) {
        setCurrentUserPage(newPage);
      }
    },
    [totalUsers]
  );

  const handleSubmissionPageChange = useCallback(
    (newPage) => {
      const maxPages = Math.ceil(totalPendingSubmissions / itemsPerPage);
      if (newPage >= 1 && newPage <= maxPages) {
        setCurrentSubmissionPage(newPage);
      }
    },
    [totalPendingSubmissions]
  );

  if (isFetching) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-xl text-purple-600">Loading...</div>
      </div>
    );
  }

  if (Object.keys(errors).length > 0) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-xl text-red-600">Error loading data</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-purple-50">
      <AdminTopbar currentPage="Dashboard" />

      <main className="max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <h1 className="text-3xl font-bold text-purple-800 mb-8" style={{ fontFamily: "Comic Sans MS, cursive" }}>
          Admin Dashboard
        </h1>

        {/* Stats Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <div className="bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
            <div className="text-purple-600 text-lg font-semibold">
              Pending Reviews
            </div>
            <div className="text-3xl font-bold">{totalPendingSubmissions}</div>
          </div>
          <div className="bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
            <div className="text-purple-600 text-lg font-semibold">
              Total ComicCoins Paid
            </div>
            <div className="text-3xl font-bold">{totalCoinsAwarded} CC</div>
          </div>
          <div className="bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
            <div className="text-purple-600 text-lg font-semibold">
              Faucet Balance
            </div>
            <div className="text-3xl font-bold">{faucetBalance} CC</div>
          </div>
        </div>

        {/* Users Section */}
        <div className="bg-white rounded-xl shadow-lg p-6 mb-8 border-2 border-purple-200">
          <h2 className="text-2xl font-bold text-purple-800 mb-6" style={{ fontFamily: "Comic Sans MS, cursive" }}>
            Users Awaiting Review
          </h2>

          {users.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-12 text-gray-500">
              <Users className="w-12 h-12 mb-4 text-purple-300" />
              <p className="text-lg font-medium mb-2">No Pending Reviews</p>
              <p className="text-sm text-gray-400">
                There are currently no users waiting for review.
              </p>
            </div>
          ) : (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead>
                    <tr className="bg-purple-50">
                      <th className="px-4 py-2 text-left text-purple-600">Full Name</th>
                      <th className="px-4 py-2 text-left text-purple-600">Email</th>
                      <th className="px-4 py-2 text-left text-purple-600">Country</th>
                      <th className="px-4 py-2 text-left text-purple-600">Region</th>
                      <th className="px-4 py-2 text-left text-purple-600">City</th>
                      <th className="px-4 py-2 text-left text-purple-600">Join Date</th>
                    </tr>
                  </thead>
                  <tbody>
                    {users && users.map((user) => (
                      <tr key={user.id} className="border-t border-purple-100 hover:bg-purple-50 cursor-pointer" onClick={() => handleUserClick(user)}>
                        <td className="px-4 py-2">{user.name}</td>
                        <td className="px-4 py-2">{user.email}</td>
                        <td className="px-4 py-2">{user.country}</td>
                        <td className="px-4 py-2">{user.region}</td>
                        <td className="px-4 py-2">{user.city}</td>
                        <td className="px-4 py-2">{user.createdAt}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}

          {/* User Pagination */}
          <div className="mt-4 flex items-center justify-between">
            <div className="text-sm text-gray-600">
              Showing {((currentUserPage - 1) * usersPerPage) + 1} to {Math.min(currentUserPage * usersPerPage, totalUsers)} of {totalUsers} users
            </div>
            <div className="flex items-center space-x-2">
              <button
                onClick={() => handleUserPageChange(currentUserPage - 1)}
                disabled={currentUserPage === 1}
                className="p-2 rounded-lg border border-purple-200 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-purple-50"
              >
                <ChevronLeft className="w-5 h-5 text-purple-600" />
              </button>
              <span className="text-sm text-gray-600">
                Page {currentUserPage} of {Math.ceil(totalUsers / usersPerPage)}
              </span>
              <button
                onClick={() => handleUserPageChange(currentUserPage + 1)}
                disabled={currentUserPage >= Math.ceil(totalUsers / usersPerPage)}
                className="p-2 rounded-lg border border-purple-200 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-purple-50"
              >
                <ChevronRight className="w-5 h-5 text-purple-600" />
              </button>
            </div>
          </div>
        </div>

        {/* Submissions Section */}
        <div className="bg-white rounded-xl shadow-lg p-6 mb-8 border-2 border-purple-200">
          <h2 className="text-2xl font-bold text-purple-800 mb-6" style={{ fontFamily: "Comic Sans MS, cursive" }}>
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
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
              {pendingSubmissions.map((submission) => (
                <div key={submission.id} className="max-w-[300px] mx-auto w-full">
                  <GalleryItem
                    submission={submission}
                    onFlag={handleFlagSubmission}
                    handleApproveSubmission={handleApproveSubmission}
                    handleRejectSubmission={handleRejectSubmission}
                  />
                </div>
              ))}
            </div>
          )}

          {/* Submission Pagination */}
          <div className="mt-8 flex flex-col md:flex-row items-center justify-between gap-4">
            <div className="text-sm text-gray-600">
              {pendingSubmissions.length === 0
                ? "No submissions to display"
                : `Showing ${((currentSubmissionPage - 1) * itemsPerPage) + 1} to ${Math.min(
                    currentSubmissionPage * itemsPerPage,
                    totalPendingSubmissions
                  )} of ${totalPendingSubmissions} submissions`}
            </div>
            {totalPendingSubmissions > 0 && Math.ceil(totalPendingSubmissions / itemsPerPage) > 1 && (
              <div className="flex items-center space-x-2">
                <button
                  onClick={() => handleSubmissionPageChange(currentSubmissionPage - 1)}
                  disabled={currentSubmissionPage === 1}
                  className="p-2 rounded-lg border border-purple-200 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-purple-50"
                >
                  <ChevronLeft className="w-5 h-5 text-purple-600" />
                </button>
                <span className="text-sm text-gray-600">
                  Page {currentSubmissionPage} of {Math.ceil(totalPendingSubmissions / itemsPerPage)}
                </span>
                <button
                  onClick={() => handleSubmissionPageChange(currentSubmissionPage + 1)}
                  disabled={currentSubmissionPage >= Math.ceil(totalPendingSubmissions / itemsPerPage)}
                  className="p-2 rounded-lg border border-purple-200 disabled:opacity-50 disabled:cursor-not-allowed hover:bg-purple-50"
                >
                  <ChevronRight className="w-5 h-5 text-purple-600" />
                </button>
              </div>
            )}
          </div>
        </div>

        {/* Hidden image preload */}
        {pendingSubmissions && (
          <div style={{ display: "none" }}>
            {pendingSubmissions.map((submission, index) => (
              <React.Fragment key={submission.id}>
                <img
                  src={submission.frontCover?.objectUrl}
                  alt={`Preloading front ${index + 1}`}
                  onLoad={() => console.log(`Front image ${index + 1} loaded successfully`)}
                  onError={() => console.error(`Front image ${index + 1} failed to load`)}
                />
                <img
                  src={submission.backCover?.objectUrl}
                  alt={`Preloading back ${index + 1}`}
                  onLoad={() => console.log(`Back image ${index + 1} loaded successfully`)}
                  onError={() => console.error(`Back image ${index + 1} failed to load`)}
                />
              </React.Fragment>
            ))}
          </div>
        )}
      </main>
      {selectedUser && (
          <UserDetailModal
            user={selectedUser}
            onClose={handleCloseModal}
            onAccept={handleAcceptUser}
            onReject={handleRejectUser}
          />
       )}
    </div>
  );
};

export default AdminDashboard;
