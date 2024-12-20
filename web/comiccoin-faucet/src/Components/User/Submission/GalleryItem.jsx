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

const GalleryItem = ({ submission, onClick }) => {
  const statusInfo = getStatusInfo(submission.status);
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
      default:
        return 'border-purple-100';
    }
  };

  return (
    <div
      className={`w-full bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow cursor-pointer border-2 ${getBorderStyle()}`}
      onClick={() => onClick?.(submission)}
    >
      <div className="relative w-full aspect-[2/3]">
        <div className="relative h-full">
          <img
            src={submission.frontCover?.objectUrl || "/api/placeholder/256/384"}
            alt={submission.name}
            className={`w-full h-full object-cover rounded-t-lg ${isRejected ? 'opacity-50 grayscale' : ''}`}
          />
          {statusInfo.overlayClass && (
            <div className={`absolute inset-0 ${statusInfo.overlayClass} rounded-t-lg`} />
          )}
        </div>
        <div className="absolute top-2 right-2 bg-white rounded-full p-1.5 shadow-md">
          {statusInfo.icon}
        </div>
      </div>

      <div className="p-3">
        <h3 className="font-medium text-sm truncate" title={submission.name}>
          {submission.name}
        </h3>
        <p className="text-sm text-gray-600 truncate">
          by {submission.createdByUserName}
        </p>
        <div className="flex items-center justify-between mt-2">
          <span className={`text-sm font-medium ${statusInfo.color} inline-flex items-center gap-1`}>
            {statusInfo.text}
          </span>
          {isAccepted && submission.coinsAwarded > 0 && (
            <span className="text-sm text-green-600 flex items-center gap-1">
              <Coins className="w-4 h-4" />
              {submission.coinsAwarded}
            </span>
          )}
        </div>
        {isRejected && submission.reason && (
          <p className="text-sm text-red-600 mt-2 bg-red-50 p-2 rounded">
            {submission.reason}
          </p>
        )}
        <p className="text-sm text-gray-500 mt-2">
          {new Date(submission.createdAt).toLocaleDateString()}
        </p>
        {isAccepted && (
          <p className="text-sm text-green-600 mt-1">
            Approved: {new Date(submission.modifiedAt).toLocaleDateString()}
          </p>
        )}
      </div>
    </div>
  );
};

export default GalleryItem;
