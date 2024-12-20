import React, { useState, useEffect } from 'react';
import {
  Coins, Home, Image, History, Wallet, Flag,
  Settings, HelpCircle, LogOut, Clock, CheckCircle, XCircle,
  Menu, X, ChevronLeft, ChevronRight, Archive, AlertTriangle
} from 'lucide-react';
import { Navigate, Link } from "react-router-dom";
import { useRecoilState } from "recoil";

import { currentUserState } from "../../../AppState";
import Topbar from "../../../Components/Navigation/Topbar";
import { getComicSubmissionListAPI } from "../../../API/ComicSubmission";
import SubmissionModal from './ListModal';
import GalleryItem from './GalleryItem';

const ITEMS_PER_PAGE = 12;

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
    case 6: // ComicSubmissionStatusFlagged
        return { icon: <Flag className="w-4 h-4 text-gray-500" />, color: 'text-red-500', text: 'Flagged' };
    default:
      return { icon: null, color: '', text: 'Unknown' };
  }
};

const PaginationControls = ({ currentPage, totalPages, onPageChange, disabled }) => {
  return (
    <div className="flex items-center justify-center gap-2 mt-6">
      <button
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage === 1 || disabled}
        className="p-2 rounded-lg hover:bg-purple-100 disabled:opacity-50 disabled:hover:bg-transparent"
      >
        <ChevronLeft className="w-5 h-5 text-purple-600" />
      </button>

      <div className="flex items-center gap-1">
        {Array.from({ length: totalPages }, (_, i) => i + 1).map((page) => (
          <button
            key={page}
            onClick={() => onPageChange(page)}
            disabled={disabled}
            className={`w-8 h-8 rounded-lg ${
              currentPage === page
                ? 'bg-purple-600 text-white'
                : 'hover:bg-purple-100'
            } disabled:opacity-50`}
          >
            {page}
          </button>
        ))}
      </div>

      <button
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage === totalPages || disabled}
        className="p-2 rounded-lg hover:bg-purple-100 disabled:opacity-50 disabled:hover:bg-transparent"
      >
        <ChevronRight className="w-5 h-5 text-purple-600" />
      </button>
    </div>
  );
};

const SubmissionsPage = () => {
  const [currentUser] = useRecoilState(currentUserState);

  // Component state
  const [isFetching, setFetching] = useState(false);
  const [errors, setErrors] = useState({});
  const [submissions, setSubmissions] = useState([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [pageStates, setPageStates] = useState([]);
  const [selectedSubmission, setSelectedSubmission] = useState(null);

  const fetchSubmissions = async (page) => {
    setFetching(true);
    const params = new Map();
    params.set("limit", ITEMS_PER_PAGE);
    params.set("user_id", currentUser.id);

    // Get the last ID and created_at from previous page state
    if (page > 1 && pageStates[page - 2]) {
      const prevState = pageStates[page - 2];
      params.set("last_id", prevState.lastId);
      params.set("last_created_at", prevState.lastCreatedAt);
    }

    try {
      getComicSubmissionListAPI(
        params,
        (resp) => {
          // Update page states
          const newPageStates = [...pageStates];
          newPageStates[page - 1] = {
            submissions: resp.submissions,
            lastId: resp.lastId,
            lastCreatedAt: resp.lastCreatedAt,
            hasMore: resp.hasMore
          };
          setPageStates(newPageStates);

          // Update current page submissions
          setSubmissions(resp.submissions);

          // Update total pages if this is first page
          if (page === 1) {
            const calculatedTotalPages = resp.hasMore ?
              Math.ceil((ITEMS_PER_PAGE * 2) / ITEMS_PER_PAGE) :
              Math.ceil(resp.submissions.length / ITEMS_PER_PAGE);
            setTotalPages(calculatedTotalPages);
          }
        },
        setErrors,
        () => setFetching(false),
        () => window.location.href = "/login?unauthorized=true"
      );
    } catch (error) {
      setErrors(error);
      setFetching(false);
    }
  };

  const handlePageChange = (newPage) => {
    if (newPage === currentPage || isFetching) return;

    setCurrentPage(newPage);

    // If we already have the page data, use it
    if (pageStates[newPage - 1]) {
      setSubmissions(pageStates[newPage - 1].submissions);
    } else {
      // Otherwise fetch new data
      fetchSubmissions(newPage);
    }

    // Update total pages if we're moving to a new page and there's more data
    if (newPage === totalPages && pageStates[newPage - 1]?.hasMore) {
      setTotalPages(prev => prev + 1);
    }
  };

  useEffect(() => {
    window.scrollTo(0, 0);
    fetchSubmissions(1);
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
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6 gap-6">
              {submissions.map(submission => (
                <div key={submission.id} className="w-full">
                  <GalleryItem
                    submission={submission}
                    onClick={setSelectedSubmission}
                  />
                </div>
              ))}
            </div>

            <PaginationControls
              currentPage={currentPage}
              totalPages={totalPages}
              onPageChange={handlePageChange}
              disabled={isFetching}
            />
          </div>
        )}
      </main>

      {selectedSubmission && (
        <SubmissionModal
          submission={selectedSubmission}
          onClose={() => setSelectedSubmission(null)}
        />
      )}

      {/* Hidden image elements to preload */}
      {submissions && <div style={{ display: "none" }}>
          {submissions.map((submission, index) => (
            <>
                <img
                  key={`front_${index}`}
                  src={submission.frontCover?.objectUrl}
                  alt={`Preloading ${index + 1}`}
                  onLoad={() => console.log(`Image ${index + 1} loaded successfully`)}
                  onError={() => console.error(`Image ${index + 1} failed to load`)}
                />
                <img
                  key={`back_${index}`}
                  src={submission.backCover?.objectUrl}
                  alt={`Preloading ${index + 1}`}
                  onLoad={() => console.log(`Image ${index + 1} loaded successfully`)}
                  onError={() => console.error(`Image ${index + 1} failed to load`)}
                />
            </>
          ))}
      </div>}
    </div>
  );
};

export default SubmissionsPage;
