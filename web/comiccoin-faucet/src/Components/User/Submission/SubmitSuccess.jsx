import React, { useState, useEffect } from "react";
import {
  CheckCircle,
  Clock,
  ArrowRight,
  AlertCircle,
  Home,
  Upload,
  BadgeCheck,
} from "lucide-react";
import { Link, useSearchParams } from "react-router-dom";
import { useRecoilState } from "recoil";

import Topbar from "../../../Components/Navigation/Topbar";
import {
  getComicSubmissionDetailAPI,
  getComicSubmissionsCountTotalCreatedTodayByUserAPI,
} from "../../../API/ComicSubmission";
import { currentUserState } from "../../../AppState";

const SubmitComicSuccessPage = () => {
  const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140
  const comicSubmissionID = searchParams.get("id");

  // Variable controls the global state of the app.
  const [currentUser] = useRecoilState(currentUserState);

  const [totalSubmissionsToday, setTotalSubmissionsToday] = useState(0);
  const [submission, setSubmission] = useState({});

  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [errors, setErrors] = useState({});

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.

      setFetching(true);
      getComicSubmissionsCountTotalCreatedTodayByUserAPI(
        currentUser.id,
        (resp) => {
          // For debugging purposes only.
          console.log(
            "getComicSubmissionsCountTotalCreatedTodayByUserAPI: Starting...",
          );
          console.log(resp);
          setTotalSubmissionsToday(resp.count);
        },
        (apiErr) => {
          console.log(
            "getComicSubmissionsCountTotalCreatedTodayByUserAPI: apiErr:",
            apiErr,
          );
          setErrors(apiErr);
        },
        () => {
          console.log(
            "getComicSubmissionsCountTotalCreatedTodayByUserAPI: Starting...",
          );
          setFetching(false);
        },
        () => {
          console.log(
            "getComicSubmissionsCountTotalCreatedTodayByUserAPI: unauthorized...",
          );
          window.location.href = "/login?unauthorized=true";
        },
      );

      //

      setFetching(true);
      getComicSubmissionDetailAPI(
        comicSubmissionID,
        (resp) => {
          // For debugging purposes only.
          console.log("getComicSubmissionDetailAPI: Starting...");
          console.log(resp);
          setSubmission(resp);
        },
        (apiErr) => {
          console.log("getComicSubmissionDetailAPI: apiErr:", apiErr);
          setErrors(apiErr);
        },
        () => {
          console.log("getComicSubmissionDetailAPI: Starting...");
          setFetching(false);
        },
        () => {
          console.log("getComicSubmissionDetailAPI: unauthorized...");
          window.location.href = "/login?unauthorized=true";
        },
      );

      //
    }

    return () => {
      mounted = false;
    };
  }, [currentUser, comicSubmissionID]);

  const totalSubmissionsTodayPercent = 100 * (totalSubmissionsToday / 6);
  const totalSubmissionsTodayPercentStr = totalSubmissionsTodayPercent + "%";

  return (
    <div className="min-h-screen bg-purple-50">
      <Topbar currentPage="Submit Comic" />
      <main className="p-4 lg:p-8 max-w-4xl mx-auto">
        {/* Success Message */}
        <div className="bg-white rounded-xl shadow-lg p-8 mb-8 border-2 border-green-200 text-center">
          <div className="flex justify-center mb-6">
            <CheckCircle className="h-16 w-16 text-green-500" />
          </div>
          <h1 className="text-2xl lg:text-3xl font-bold text-green-600 mb-4">
            Submission Successful!
          </h1>
          <p className="text-gray-600 text-lg mb-6">
            Your comic "
            <span className="font-semibold">
              {submission && <>{submission.name}</>}
            </span>
            " has been submitted for review
          </p>
          <div className="flex items-center justify-center space-x-2 text-purple-600">
            <Clock className="h-5 w-5" />
            <span>Estimated review time: 24-48 hours</span>
          </div>
        </div>

        {/* Pending Reward Card */}
        <div className="bg-white rounded-xl shadow-lg p-6 mb-8 border-2 border-purple-200">
          <h2
            className="text-xl font-bold text-purple-800 mb-4"
            style={{ fontFamily: "Comic Sans MS, cursive" }}
          >
            Pending Reward
          </h2>
          <div className="bg-purple-50 rounded-lg p-6">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center space-x-3">
                <Clock className="h-6 w-6 text-purple-600" />
                <span className="text-lg font-semibold text-purple-800">
                  {submission && <>{submission.coinsReward}</>} ComicCoins
                </span>
              </div>
              <span className="px-3 py-1 bg-purple-100 text-purple-600 rounded-full text-sm">
                Pending Approval
              </span>
            </div>
            <p className="text-gray-600 text-sm mb-4">
              Your ComicCoins will be automatically transferred to your wallet
              once your submission is approved.
            </p>
            <div className="flex items-center space-x-2 text-sm text-purple-600">
              <AlertCircle className="h-4 w-4" />
              <span>Rewards are subject to our content review guidelines</span>
            </div>
          </div>
        </div>

        {/* Daily Limit Notice */}
        <div className="bg-white rounded-xl shadow-lg p-6 mb-8 border-2 border-purple-200">
          <h2
            className="text-xl font-bold text-purple-800 mb-4"
            style={{ fontFamily: "Comic Sans MS, cursive" }}
          >
            Daily Submission Limit
          </h2>
          <div className="grid md:grid-cols-2 gap-6">
            <div className="p-4 bg-purple-50 rounded-lg">
              <h3 className="font-semibold text-purple-800 mb-2 flex items-center space-x-2">
                <Upload className="h-5 w-5" />
                <span>Standard User</span>
              </h3>
              <p className="text-gray-600 mb-3">
                You have used {totalSubmissionsToday}/6 submissions today
              </p>
              <div className="w-full bg-purple-200 rounded-full h-2">
                <div
                  className="bg-purple-600 h-2 rounded-full"
                  style={{ width: totalSubmissionsTodayPercentStr }}
                ></div>
              </div>
              <p className="mt-3 text-sm text-gray-500">
                Limit resets daily at midnight
              </p>
            </div>

            <div className="p-4 bg-purple-50 rounded-lg">
              <h3 className="font-semibold text-purple-800 mb-2 flex items-center space-x-2">
                <BadgeCheck className="h-5 w-5" />
                <span>Become Verified</span>
              </h3>
              <p className="text-gray-600 mb-3">
                Get increased daily submission limits and extra benefits
              </p>
              <Link
                to="/apply-for-verification"
                className="w-full px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors flex items-center justify-center space-x-2"
              >
                <span>Apply for Verification</span>
                <ArrowRight className="h-4 w-4" />
              </Link>
            </div>
          </div>
        </div>

        {/* Action Buttons */}
        <div className="flex flex-col sm:flex-row justify-center space-y-4 sm:space-y-0 sm:space-x-6">
          <Link
            to="/submit"
            className="flex items-center justify-center space-x-2 px-6 py-3 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors"
          >
            <Upload className="h-5 w-5" />
            <span>Submit Another Comic</span>
          </Link>
          <Link
            to="/dashboard"
            className="flex items-center justify-center space-x-2 px-6 py-3 border-2 border-purple-200 text-purple-600 rounded-lg hover:bg-purple-50 transition-colors"
          >
            <Home className="h-5 w-5" />
            <span>Return to Dashboard</span>
          </Link>
        </div>
      </main>
    </div>
  );
};

export default SubmitComicSuccessPage;
