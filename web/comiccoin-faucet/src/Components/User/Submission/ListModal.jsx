import React from 'react';
import { X, Flag, Coins, Clock, XCircle, CheckCircle } from 'lucide-react';

const getStatusInfo = (status) => {
  switch (status) {
    case 1: // ComicSubmissionStatusInReview
      return {
        icon: <Clock className="w-4 h-4 text-yellow-500" />,
        color: 'text-yellow-500',
        text: 'In Review',
        overlayClass: 'bg-yellow-500 bg-opacity-10'
      };
    case 2: // ComicSubmissionStatusRejected
      return {
        icon: <XCircle className="w-4 h-4 text-red-500" />,
        color: 'text-red-500',
        text: 'Rejected',
        overlayClass: 'bg-red-500 bg-opacity-10'
      };
    case 3: // ComicSubmissionStatusAccepted
      return {
        icon: <CheckCircle className="w-4 h-4 text-green-500" />,
        color: 'text-green-500',
        text: 'Approved',
        overlayClass: 'bg-green-500 bg-opacity-20'
      };
    default:
      return {
        icon: null,
        color: '',
        text: 'Unknown',
        overlayClass: ''
      };
  }
};


const SubmissionModal = ({ submission, onClose }) => {
  if (!submission) return null;

  const statusInfo = getStatusInfo(submission.status);
  const isFlagged = submission.status === 6;
  const isAccepted = submission.status === 3;
  const isRejected = submission.status === 2;
  const isInReview = submission.status === 1;

  const getBorderStyle = () => {
    switch (submission.status) {
      case 1:
        return 'border-yellow-400';
      case 2:
        return 'border-red-500';
      case 3:
        return 'border-green-500';
      case 6:
        return 'border-red-500';
      default:
        return 'border-purple-100';
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div className={`bg-white rounded-xl max-w-2xl w-full relative max-h-[90vh] overflow-y-auto border ${getBorderStyle()}`}>
        <div className="p-6">
          <button
            onClick={onClose}
            className="absolute top-4 right-4 text-gray-500 hover:text-gray-700"
          >
            <X className="w-6 h-6" />
          </button>

          <div className="flex flex-col md:flex-row gap-6">
            {/* Image Section */}
            <div className="w-full md:w-auto">
              {isFlagged ? (
                <div className="w-full md:w-64 h-80 flex items-center justify-center bg-gray-100 rounded-lg">
                  <Flag className="w-32 h-32 text-red-500" />
                </div>
              ) : (
                <div className="relative">
                  <img
                    src={submission.frontCover?.objectUrl || "/api/placeholder/256/320"}
                    alt={submission.name}
                    className={`w-full md:w-64 h-80 object-cover rounded-lg ${isRejected ? 'opacity-50 grayscale' : ''}`}
                  />
                  {isAccepted && (
                    <div className="absolute inset-0 bg-green-500 bg-opacity-20 rounded-lg" />
                  )}
                  {isRejected && (
                    <div className="absolute inset-0 bg-red-500 bg-opacity-10 rounded-lg" />
                  )}
                  {isInReview && (
                    <div className="absolute inset-0 bg-yellow-500 bg-opacity-10 rounded-lg" />
                  )}
                </div>
              )}
              {submission.backCover && !isFlagged && (
                <div className="relative mt-4">
                  <img
                    src={submission.backCover.objectUrl}
                    alt="Back cover"
                    className={`w-full md:w-64 h-80 object-cover rounded-lg ${isRejected ? 'opacity-50 grayscale' : ''}`}
                  />
                  {isAccepted && (
                    <div className="absolute inset-0 bg-green-500 bg-opacity-20 rounded-lg" />
                  )}
                  {isRejected && (
                    <div className="absolute inset-0 bg-red-500 bg-opacity-10 rounded-lg" />
                  )}
                  {isInReview && (
                    <div className="absolute inset-0 bg-yellow-500 bg-opacity-10 rounded-lg" />
                  )}
                </div>
              )}
            </div>

            {/* Details Section */}
            <div className="flex-1">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-2xl font-bold text-purple-800">
                  {submission.name}
                </h2>
                <span className={`inline-flex items-center gap-1 px-3 py-1 rounded-full ${statusInfo.color} bg-opacity-10`}>
                  {statusInfo.icon}
                  <span className="font-medium">{statusInfo.text}</span>
                </span>
              </div>

              {isRejected && submission.reason && (
                <div className="mb-4 p-3 bg-red-50 rounded-lg border border-red-200">
                  <h3 className="text-sm font-medium text-red-800 mb-1">Rejection Reason:</h3>
                  <p className="text-sm text-red-600">{submission.reason}</p>
                </div>
              )}

              {isAccepted && submission.coinsAwarded > 0 && (
                <div className="mb-4 p-3 bg-green-50 rounded-lg border border-green-200">
                  <h3 className="text-sm font-medium text-green-800 mb-1">Reward Earned:</h3>
                  <div className="flex items-center gap-1 text-green-600">
                    <Coins className="w-5 h-5" />
                    <span className="text-lg font-medium">{submission.coinsAwarded} Coins</span>
                  </div>
                </div>
              )}

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-4">
                <div className="p-3 bg-gray-50 rounded-lg">
                  <p className="text-sm text-gray-500">Submitted By</p>
                  <p className="font-medium">{submission.createdByUserName}</p>
                  <p className="text-sm text-gray-500 mt-1">
                    {new Date(submission.createdAt).toLocaleDateString()}
                  </p>
                </div>

                {submission.modifiedAt && (
                  <div className="p-3 bg-gray-50 rounded-lg">
                    <p className="text-sm text-gray-500">Last Modified By</p>
                    <p className="font-medium">{submission.modifiedByUserName}</p>
                    <p className="text-sm text-gray-500 mt-1">
                      {new Date(submission.modifiedAt).toLocaleDateString()}
                    </p>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SubmissionModal;
