import React, { useState, useEffect } from 'react';
import {
  Coins, Home, Image, History, Wallet,
  Settings, HelpCircle, LogOut, Clock, CheckCircle, XCircle,
  Menu, X, Upload, ArrowRight, Sparkles
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


const DashboardPage = () => {

  // Variable controls the global state of the app.
  const [currentUser] = useRecoilState(currentUserState);

  // GUI related
  const [isNavOpen, setIsNavOpen] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [isFetching, setFetching] = useState(false);
  const [errors, setErrors] = useState({});

  // Data related.
  const [totalSubmissions, setTotalSubmissions] = useState(0);
  const [totalApprovedSubmissions, setTotalApprovedSubmissions] = useState(0);
  const [totalCoinsEarned, setTotalCoinsEarned] = useState(0);
  const [pendingSubmissions, setPendingSubmissions] = useState([]);
  const [recentSubmissions, setRecentSubmissions] = useState([]);

  const getStatusIcon = (status) => {
    switch (status) {
      case 1: // ComicSubmissionStatusInReview
        return <Clock className="w-4 h-4 text-yellow-500" />;
      case 3: // ComicSubmissionStatusAccepted
        return <CheckCircle className="w-4 h-4 text-green-500" />;
      case 2: // ComicSubmissionStatusRejected
        return <XCircle className="w-4 h-4 text-red-500" />;
      default:
        return null;
    }
  };

  const getStatusColor = (status) => {
    switch (status) {
      case 1: // ComicSubmissionStatusInReview
        return 'text-yellow-500';
      case 3: // ComicSubmissionStatusAccepted
        return 'text-green-500';
      case 2: // ComicSubmissionStatusRejected
        return 'text-red-500';
      default:
        return '';
    }
  };

  const GalleryItem = ({ submission }) => (
    <div className="w-32 bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow border border-purple-100">
      <div className="relative w-32 h-44">
        <img
          src={submission.frontCover.objectUrl}
          alt={submission.name}
          className="w-full h-full object-cover rounded-t-lg"
        />
        <div className="absolute top-1 right-1 bg-white rounded-full p-1 shadow">
          {getStatusIcon(submission.status)}
        </div>
      </div>
      <div className="p-2">
        <h3 className="font-medium text-xs truncate" title={submission.title}>
          {submission.name}
        </h3>
        <p className="text-xs text-gray-600 truncate">by {submission.createdByUserName}</p>
        <p className="text-xs mt-1">
          <span className={`font-medium ${getStatusColor(submission.status)}`}>
            {submission.status === 1 ? 'In Review' :
             submission.status === 2 ? 'Approved' : 'Rejected'}
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
          {new Date(submission.createdAt).toLocaleDateString()}
        </p>
      </div>
    </div>
  );

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.

      //------------------------------------------------------------------------

      setFetching(true);
      let params = new Map();
      params.set("user_id", currentUser.id);
      getComicSubmissionsCountByFilterAPI(
        params,
        (resp) => {
          // For debugging purposes only.
          console.log("getComicSubmissionsCountByFilterAPI: Starting...");
          console.log(resp);
          setTotalSubmissions(resp.count);
        },
        (apiErr) => {
          console.log("getComicSubmissionsCountByFilterAPI: apiErr:", apiErr);
          setErrors(apiErr);
        },
        () => {
          console.log("getComicSubmissionsCountByFilterAPI: Starting...");
          setFetching(false);
        },
        () => {
          console.log("getComicSubmissionsCountByFilterAPI: unauthorized...");
          window.location.href = "/login?unauthorized=true";
        },
      );

      //------------------------------------------------------------------------

      params = new Map();
      // params.set("page_size", limit); // Pagination
      // params.set("sort_field", "created_at"); // Sorting
      // params.set("sort_order", -1); // Sorting - descending, meaning most recent start date to oldest start date.
      params.set("status", 1); // ComicSubmissionStatusInReview
      params.set("user_id", currentUser.id);
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

      params = new Map();
      params.set("user_id", currentUser.id);
      params.set("status", 3); // Success
      getComicSubmissionsCountByFilterAPI(
        params,
        (resp) => {
          // For debugging purposes only.
          console.log("getComicSubmissionsCountByFilterAPI: ",resp);
          setTotalApprovedSubmissions(resp.count);
        },
        (apiErr) => {
          console.log("getComicSubmissionsCountByFilterAPI: apiErr:", apiErr);
          setErrors(apiErr);
        },
        () => {
          console.log("getComicSubmissionsCountByFilterAPI: Starting...");
          setFetching(false);
        },
        () => {
          console.log("getComicSubmissionsCountByFilterAPI: unauthorized...");
          window.location.href = "/login?unauthorized=true";
        },
      );

      //------------------------------------------------------------------------

      params = new Map();
      params.set("limit", 5);
      // params.set("sort_field", "created_at"); // Sorting
      // params.set("sort_order", -1); // Sorting - descending, meaning most recent start date to oldest start date.
      params.set("status", 3); // ComicSubmissionStatusAccepted
      params.set("user_id", currentUser.id);
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
          console.log("getComicSubmissionListAPI: Recent Submissions:", resp);
          setRecentSubmissions(resp.submissions);
        },
        (apiErr) => {
          console.log("getComicSubmissionListAPI: Recent Submissions: apiErr:", apiErr);
          setErrors(apiErr);
        },
        () => {
          console.log("getComicSubmissionListAPI: Recent Submissions: Starting...");
          setFetching(false);
        },
        () => {
          console.log("getComicSubmissionListAPI: Recent Submissions: unauthorized...");
          window.location.href = "/login?unauthorized=true";
        },
      );

      //------------------------------------------------------------------------

      setFetching(true);
      params = new Map();
      params.set("user_id", currentUser.id);
      params.set("status", 3); // ComicSubmissionStatusAccepted
      getComicSubmissionsCountCoinsRewardByFilterAPI(
        params,
        (resp) => {
          // For debugging purposes only.
          console.log("getComicSubmissionsCountCoinsRewardByFilterAPI: Results:", resp);
          setTotalCoinsEarned(resp.count);
        },
        (apiErr) => {
          console.log("getComicSubmissionsCountCoinsRewardByFilterAPI: apiErr:", apiErr);
          setErrors(apiErr);
        },
        () => {
          console.log("getComicSubmissionsCountCoinsRewardByFilterAPI: Starting...");
          setFetching(false);
        },
        () => {
          console.log("getComicSubmissionsCountCoinsRewardByFilterAPI: unauthorized...");
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
        <Topbar currentPage="Dashboard" />

        {/* Main Content */}
        <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 xl:px-12 2xl:px-24 py-8">
          {/* Dashboard Header */}
          <h1 className="text-3xl font-bold text-purple-800 mb-8" style={{fontFamily: 'Comic Sans MS, cursive'}}>
            Dashboard
          </h1>

          {/* Getting Started Section */}
          {totalSubmissions === 0 && <div className="bg-white rounded-xl shadow-lg p-6 mb-8 border-2 border-purple-200">
            <div className="mb-6">
              <h2 className="text-2xl font-bold text-purple-800 mb-2" style={{fontFamily: 'Comic Sans MS, cursive'}}>
                Welcome to ComicCoin! 👋
              </h2>
              <p className="text-gray-600">Let's get started with your comic book collection journey</p>
            </div>

            <div className="grid md:grid-cols-3 gap-6">
              <div className="flex flex-col items-center text-center p-6 bg-purple-50 rounded-lg">
                <Upload className="h-12 w-12 text-purple-600 mb-4" />
                <h3 className="font-semibold mb-2">1. Submit Your First Comic</h3>
                <p className="text-sm text-gray-600 mb-4">Upload photos of your comic book covers to start earning ComicCoins</p>
                <a href="#" className="text-purple-600 hover:text-purple-700 font-medium flex items-center">
                  Start Submission <ArrowRight className="h-4 w-4 ml-1" />
                </a>
              </div>
              <div className="flex flex-col items-center text-center p-6 bg-purple-50 rounded-lg">
                <Clock className="h-12 w-12 text-purple-600 mb-4" />
                <h3 className="font-semibold mb-2">2. Wait for Review</h3>
                <p className="text-sm text-gray-600 mb-4">Our community verifies submissions within 24 hours</p>
                <p className="text-purple-600 font-medium">Quick & Easy Process</p>
              </div>
              <div className="flex flex-col items-center text-center p-6 bg-purple-50 rounded-lg">
                <Sparkles className="h-12 w-12 text-purple-600 mb-4" />
                <h3 className="font-semibold mb-2">3. Earn Rewards</h3>
                <p className="text-sm text-gray-600 mb-4">Get ComicCoins for each approved submission</p>
                <p className="text-purple-600 font-medium">Instant Rewards</p>
              </div>
            </div>
          </div>}

          {/* Stats Row */}
          <div className="flex flex-row justify-between items-stretch gap-6 mb-8">
            <div className="flex-1 bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
              <div className="text-purple-600 text-lg font-semibold">Total Submissions</div>
              <div className="text-3xl font-bold">{totalSubmissions}</div>
            </div>
            <div className="flex-1 bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
              <div className="text-purple-600 text-lg font-semibold">Comics Approved</div>
              <div className="text-3xl font-bold">{totalApprovedSubmissions}</div>
            </div>
            <div className="flex-1 bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
              <div className="text-purple-600 text-lg font-semibold">ComicCoins Earned</div>
              <div className="text-3xl font-bold">{totalCoinsEarned}</div>
            </div>
          </div>

          {/* Pending Reviews Section */}
          <div className="bg-white rounded-xl shadow-lg p-6 mb-8 border-2 border-purple-200">
              <h2 className="text-xl lg:text-2xl font-bold text-purple-800 mb-4" style={{fontFamily: 'Comic Sans MS, cursive'}}>
                Pending Reviews
              </h2>
              {pendingSubmissions.length === 0 ? <div className="text-center py-12 bg-purple-50 rounded-lg">
                <Image className="h-16 w-16 text-purple-300 mx-auto mb-4" />
                <p className="text-gray-500 mb-4">No pending submissions yet</p>
                <Link to="/submit" className="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors">
                  Submit Your First Comic
                </Link>
              </div> :
              <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4">
                {pendingSubmissions.map(submission => (
                  <GalleryItem key={submission.id} submission={submission} />
                ))}
              </div>}
          </div>

          {/* Recent Submissions Section */}
          <div className="bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
            <h2 className="text-xl lg:text-2xl font-bold text-purple-800 mb-4" style={{fontFamily: 'Comic Sans MS, cursive'}}>
              Recent Approvals
            </h2>
            <div className="text-center py-12 bg-purple-50 rounded-lg">
              <History className="h-16 w-16 text-purple-300 mx-auto mb-4" />
              <p className="text-gray-500">Your approved submissions will appear here</p>
            </div>
          </div>
        </main>
      </div>
    );
  };


export default DashboardPage;
