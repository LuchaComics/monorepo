// GalleryItem.jsx
import React, { useState } from "react";
import { X, Flag, Coins, Clock, XCircle, CheckCircle, AlertTriangle } from 'lucide-react';

const GalleryItem = ({ submission, onFlag, handleApproveSubmission, handleRejectSubmission }) => {
    const [showBackCover, setShowBackCover] = useState(false);
    const [showFlagModal, setShowFlagModal] = useState(false);

    const toggleCover = () => setShowBackCover((prev) => !prev);
    const toggleFlagModal = () => setShowFlagModal((prev) => !prev);

    // Comic books typically have a 2:3 aspect ratio
    return (
      <div className="w-full bg-white rounded-lg shadow-sm hover:shadow-md transition-shadow border-2 border-purple-200">
        <div className="relative w-full aspect-[2/3]">
          <img
            src={
              showBackCover
                ? submission.backCover?.objectUrl || "/api/placeholder/192/288"
                : submission.frontCover?.objectUrl || "/api/placeholder/192/288"
            }
            alt={`${submission.name} - ${showBackCover ? "Back" : "Front"} Cover`}
            className="w-full h-full object-cover rounded-t-lg"
          />
          <div className="absolute top-2 left-2 right-2 flex justify-between">
            <button
              onClick={toggleCover}
              className="bg-white/90 backdrop-blur-sm rounded-md px-2 py-1 text-xs font-medium shadow-sm hover:bg-white transition-colors"
            >
              {showBackCover ? "View Front" : "View Back"}
            </button>
            <div className="bg-white/90 backdrop-blur-sm rounded-full p-1.5 shadow-sm">
              <Clock className="w-4 h-4 text-yellow-500" />
            </div>
          </div>

          <div className="absolute bottom-2 left-2 right-2 flex justify-between">
            <div className="flex space-x-1.5">
              <button
                onClick={() => handleApproveSubmission(submission.id)}
                className="bg-white/90 backdrop-blur-sm rounded-full p-2 shadow-sm hover:bg-green-50 transition-colors"
                title="Approve Submission"
              >
                <CheckCircle className="w-5 h-5 text-green-500" />
              </button>
              <button
                onClick={() => handleRejectSubmission(submission.id)}
                className="bg-white/90 backdrop-blur-sm rounded-full p-2 shadow-sm hover:bg-red-50 transition-colors"
                title="Reject Submission"
              >
                <XCircle className="w-5 h-5 text-red-500" />
              </button>
              <button
                onClick={toggleFlagModal}
                className="bg-white/90 backdrop-blur-sm rounded-full p-2 shadow-sm hover:bg-yellow-50 transition-colors"
                title="Flag for Review"
              >
                <Flag
                  className={`w-5 h-5 ${submission.flagReason ? "text-yellow-500" : "text-gray-400"}`}
                />
              </button>
            </div>
          </div>
        </div>

        <div className="p-3">
          <h3 className="font-medium text-sm truncate" title={submission.name}>
            {submission.name}
          </h3>
          <p className="text-xs text-gray-600 truncate">
            by {submission.submitter || submission.createdByUserName}
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

        {showFlagModal && (
          <FlagModal
            isOpen={showFlagModal}
            onClose={toggleFlagModal}
            onSubmit={onFlag}
            submissionId={submission.id}
          />
        )}
      </div>
    );
};

export default GalleryItem;


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
