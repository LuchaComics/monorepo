import React, { useState, useEffect } from 'react';
import {
  Coins, Home, Image, History, Wallet,
  Settings, HelpCircle, LogOut, Clock, CheckCircle, XCircle,
  Menu, X, ChevronLeft, ChevronRight, Archive, AlertTriangle
} from 'lucide-react';
import { Navigate, Link } from "react-router-dom";
import { useRecoilState } from "recoil";

import { currentUserState } from "../../../AppState";
import Topbar from "../../../Components/Navigation/Topbar";
import { getComicSubmissionListAPI } from "../../../API/ComicSubmission";

const getStatusInfo = (status) => {
  switch (status) {
    case 1: // ComicSubmissionStatusInReview
      return { icon: <Clock className="w-4 h-4 text-yellow-500" />, color: 'text-yellow-500', text: 'In Review' };
    case 2: // ComicSubmissionStatusRejected
      return { icon: <XCircle className="w-4 h-4 text-red-500" />, color: 'text-red-500', text: 'Rejected' };
    case 3: // ComicSubmissionStatusAccepted
      return { icon: <CheckCircle className="w-4 h-4 text-green-500" />, color: 'text-green-500', text: 'Accepted' };
    case 4: // ComicSubmissionStatusError
      return { icon: <AlertTriangle className="w-4 h-4 text-orange-500" />, color: 'text-orange-500', text: 'Error' };
    case 5: // ComicSubmissionStatusArchived
      return { icon: <Archive className="w-4 h-4 text-gray-500" />, color: 'text-gray-500', text: 'Archived' };
    default:
      return { icon: null, color: '', text: 'Unknown' };
  }
};

const SubmissionModal = ({ submission, onClose }) => {
  if (!submission) return null;

  const statusInfo = getStatusInfo(submission.status);

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className="bg-white rounded-xl max-w-2xl w-full p-6 relative max-h-[90vh] overflow-y-auto">
        <button
          onClick={onClose}
          className="absolute top-4 right-4 text-gray-500 hover:text-gray-700"
        >
          <X className="w-6 h-6" />
        </button>

        <div className="flex flex-col md:flex-row gap-6">
          <div className="w-full md:w-auto">
            <img
              src={submission.frontCover?.objectUrl || "/api/placeholder/256/320"}
              alt={submission.name}
              className="w-full md:w-64 h-80 object-cover rounded-lg"
            />
            {submission.backCover && (
              <img
                src={submission.backCover.objectUrl}
                alt="Back cover"
                className="w-full md:w-64 h-80 object-cover rounded-lg mt-4"
              />
            )}
          </div>

          <div className="flex-1">
            <h2 className="text-2xl font-bold text-purple-800 mb-4">
              {submission.name}
            </h2>

            <div className="space-y-4">
              <div className="flex flex-wrap items-center gap-2">
                <span className={`inline-flex items-center gap-1 px-3 py-1 rounded-full ${statusInfo.color} bg-opacity-10`}>
                  {statusInfo.icon}
                  <span className="font-medium">{statusInfo.text}</span>
                </span>
                {submission.coinsReward > 0 && (
                  <>{submission.status === 3 ?
                    <span className="inline-flex items-center gap-1 px-3 py-1 rounded-full bg-green-100 text-green-600">
                      <Coins className="w-4 h-4" />
                      {submission.coinsReward} Coins
                    </span>
                    :
                    <span className="inline-flex items-center gap-1 px-3 py-1 rounded-full bg-red-100 text-red-600">
                      <Coins className="w-4 h-4" />
                      0 Coins
                    </span>}
                  </>
                )}
              </div>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
                <div>
                  <p className="text-gray-500">Submitted By</p>
                  <p className="font-medium">{submission.createdByUserName}</p>
                </div>
                <div>
                  <p className="text-gray-500">Submission Date</p>
                  <p className="font-medium">
                    {new Date(submission.createdAt).toLocaleDateString()}
                  </p>
                </div>
                {submission.modifiedAt && (
                  <>
                    <div>
                      <p className="text-gray-500">Last Modified By</p>
                      <p className="font-medium">{submission.modifiedByUserName}</p>
                    </div>
                    <div>
                      <p className="text-gray-500">Modified Date</p>
                      <p className="font-medium">
                        {new Date(submission.modifiedAt).toLocaleDateString()}
                      </p>
                    </div>
                  </>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

const SubmissionCard = ({ submission, onClick }) => {
  const statusInfo = getStatusInfo(submission.status);

  return (
    <div
      className="w-64 bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow border border-purple-100 cursor-pointer"
      onClick={() => onClick(submission)}
    >
      <div className="relative w-full h-80">
        <img
          src={submission.frontCover?.objectUrl || "/api/placeholder/256/320"}
          alt={submission.name}
          className="w-full h-full object-cover rounded-t-lg"
        />
        <div className="absolute top-2 right-2 bg-white rounded-full p-1.5 shadow">
          {statusInfo.icon}
        </div>
      </div>
      <div className="p-4">
        <h3 className="font-medium text-sm truncate" title={submission.name}>
          {submission.name}
        </h3>
        <p className="text-sm mt-2">
          <span className={`font-medium ${statusInfo.color}`}>
            {statusInfo.text}
          </span>
        </p>
        {submission.coinsReward > 0 && (
          <>{submission.status === 3 ?
              <p className="text-sm text-green-600 mt-2 flex items-center gap-1">
                <Coins className="w-4 h-4" />
                {submission.coinsReward} Coins
              </p>
              :
              <p className="text-sm text-red-600 mt-2 flex items-center gap-1">
                <Coins className="w-4 h-4" />
                0 Coins
              </p>
          }
          </>
        )}
        <p className="text-xs text-gray-500 mt-2">
          Submitted by: {submission.createdByUserName}
        </p>
        <p className="text-xs text-gray-500 mt-1">
          {new Date(submission.createdAt).toLocaleDateString()}
        </p>
      </div>
    </div>
  );
};

const SubmissionsPage = () => {
  // Variable controls the global state of the app.
  const [currentUser] = useRecoilState(currentUserState);

  // Component state
  const [isFetching, setFetching] = useState(false);
  const [errors, setErrors] = useState({});
  const [submissions, setSubmissions] = useState([]);
  const [hasMore, setHasMore] = useState(false);
  const [lastID, setLastID] = useState(null);
  const [lastCreatedAt, setLastCreatedAt] = useState(null);
  const [selectedSubmission, setSelectedSubmission] = useState(null);

  const fetchSubmissions = (resetList = false) => {
    setFetching(true);
    const params = new Map();

    if (!resetList && lastID && lastCreatedAt) {
      params.set("last_id", lastID);
      params.set("last_created_at", lastCreatedAt);
    }

    params.set("limit", 12);
    params.set("user_id", currentUser.id);

    getComicSubmissionListAPI(
      params,
      (resp) => {
        setSubmissions(prev => resetList ? resp.submissions : [...prev, ...resp.submissions]);
        setHasMore(resp.hasMore);
        setLastID(resp.lastId);
        setLastCreatedAt(resp.lastCreatedAt);
      },
      setErrors,
      () => setFetching(false),
      () => window.location.href = "/login?unauthorized=true"
    );
  };

  useEffect(() => {
    window.scrollTo(0, 0);
    fetchSubmissions(true);
  }, [currentUser]);

  return (
    <div className="min-h-screen bg-purple-50">
      <Topbar currentPage="My Submissions" />

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <h1 className="text-3xl font-bold text-purple-800 mb-8" style={{fontFamily: 'Comic Sans MS, cursive'}}>
          My Submissions
        </h1>

        {submissions.length === 0 ? (
          <div className="bg-white rounded-xl shadow-lg p-8 text-center border-2 border-purple-200">
            <Image className="w-16 h-16 text-purple-300 mx-auto mb-4" />
            <p className="text-gray-600 mb-4">No submissions found</p>
            <Link to="/submit" className="inline-block px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors">
              Submit Your First Comic
            </Link>
          </div>
        ) : (
          <div className="bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
              {submissions.map(submission => (
                <SubmissionCard
                  key={submission.id}
                  submission={submission}
                  onClick={setSelectedSubmission}
                />
              ))}
            </div>

            {hasMore && (
              <div className="mt-8 text-center">
                <button
                  onClick={() => fetchSubmissions()}
                  disabled={isFetching}
                  className="px-6 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors disabled:opacity-50"
                >
                  {isFetching ? 'Loading...' : 'Load More'}
                </button>
              </div>
            )}
          </div>
        )}
      </main>

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
