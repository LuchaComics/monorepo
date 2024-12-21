import React, { useState, useEffect } from "react";
import {
  Upload,
  X,
  AlertCircle,
  ArrowLeft,
  Camera,
  Info,
  Calendar,
  User,
} from "lucide-react";
import { Link, Navigate } from "react-router-dom";
import { useRecoilState } from "recoil";

import Topbar from "../../../Components/Navigation/Topbar";
import {
  postAttachmentCreateAPI,
  deleteAttachmentAPI,
} from "../../../API/Attachment";
import {
    postComicSubmissionCreateAPI,
    getComicSubmissionsCountByFilterAPI
} from "../../../API/ComicSubmission";
import { currentUserState } from "../../../AppState";


const SubmitComicPage = () => {
  // Variable controls the global state of the app.
  const [currentUser] = useRecoilState(currentUserState);

  const [totalSubmissions, setTotalSubmissions] = useState(0);
  const [frontCover, setFrontCover] = useState(null);
  const [backCover, setBackCover] = useState(null);
  const [frontCoverData, setFrontCoverData] = useState(null);
  const [backCoverData, setBackCoverData] = useState(null);
  const [comicName, setComicName] = useState("");
  const [showPhotoTips, setShowPhotoTips] = useState(false);
  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.

      if (currentUser === undefined || currentUser === null || currentUser === "") {
          console.log("DashboardPage: currentUser is null - not rendering view.");
          return () => {
            mounted = false;
          }
      }

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

    }

    return () => {
      mounted = false;
    };
  }, [currentUser]);

  const rules = [
    "You must only upload pictures of a physical comic book",
    "You must own the comic book you are submitting",
    "You must not have submitted this comic book previously",
    "Your submission must follow our terms of service",
    "All submissions will be reviewed for approval",
    "Upon successful review, you will receive 1 ComicCoin",
  ];

  // API handlers
  const handleFileUpload = (file, setFileData, coverType) => {
    setFetching(true);
    setErrors({});

    const formData = new FormData();
    formData.append("file", file, file.name);
    formData.append("filename", file.name);
    formData.append("mimeType", file.type || "application/octet-stream");

    postAttachmentCreateAPI(
      file.name,
      file.type || "application/octet-stream",
      formData,
      (response) => {
        console.log(`${coverType} upload success:`, response);
        console.log(`${coverType} objectUrl:`, response.objectUrl);
        setFileData(response); // Store the entire response
      },
      (apiErr) => {
        console.error(`${coverType} upload error:`, apiErr);
        setErrors((prev) => ({ ...prev, [coverType]: apiErr }));
      },
      () => {
        setFetching(false);
      },
      () => {
        window.location.href = "/login?unauthorized=true";
      },
    );
  };

  const handleDelete = (fileData, setFileData, coverType) => {
    if (!fileData?.id) return;

    setFetching(true);
    setErrors({});

    deleteAttachmentAPI(
      fileData.id,
      () => {
        setFileData(null);
      },
      (apiErr) => {
        setErrors((prev) => ({ ...prev, [coverType]: apiErr }));
      },
      () => {
        setFetching(false);
      },
      () => {
        window.location.href = "/login?unauthorized=true";
      },
    );
  };

  // File change handlers
  const handleFrontCoverChange = (event) => {
    const file = event.target.files[0];
    if (file) {
      setFrontCover(file);
      handleFileUpload(file, setFrontCoverData, "frontCover");
    }
  };

  const handleBackCoverChange = (event) => {
    const file = event.target.files[0];
    if (file) {
      setBackCover(file);
      handleFileUpload(file, setBackCoverData, "backCover");
    }
  };

  // Format date helper
  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  // Upload Preview Component
  const UploadPreview = ({
    title,
    fileData,
    onDelete,
    onChange,
    inputId,
    disabled,
  }) => (
    <div>
      <label className="block text-sm font-medium text-gray-700 mb-2">
        {title} <span className="text-gray-500">(required)</span>
      </label>
      <div className="border-2 border-dashed border-purple-200 rounded-lg p-6 hover:border-purple-400 transition-colors">
        <div className="flex flex-col items-center">
          {fileData ? (
            <div className="w-full space-y-4">
              {/* Image Preview */}
              <div className="relative w-full aspect-[3/4] bg-gray-100 rounded-lg overflow-hidden">
                {console.log("Rendering image with URL:", fileData?.objectUrl)}
                <img
                  src={fileData?.objectUrl}
                  alt={fileData?.filename}
                  className="w-full h-full object-contain"
                />
              </div>

              {/* File Info */}
              <div className="space-y-2 text-sm text-gray-600">
                <div className="flex items-center gap-2">
                  <span className="font-medium">File:</span>
                  {fileData.filename}
                </div>
                <div className="flex items-center gap-2">
                  <Calendar className="h-4 w-4" />
                  {formatDate(fileData.createdAt)}
                </div>
                <div className="flex items-center gap-2">
                  <User className="h-4 w-4" />
                  {fileData.createdByUserName}
                </div>
              </div>

              {/* Delete Button */}
              <button
                onClick={() => onDelete(fileData)}
                className="w-full px-4 py-2 bg-red-50 text-red-600 rounded-lg hover:bg-red-100 transition-colors flex items-center justify-center gap-2"
                disabled={disabled}
              >
                <X className="h-4 w-4" />
                Remove Image
              </button>
            </div>
          ) : (
            <>
              <Upload className="h-12 w-12 text-purple-400 mb-4" />
              <p className="text-sm text-gray-500 text-center mb-4">
                Click here to upload or drag and drop your photo
              </p>
              <input
                type="file"
                accept="image/*"
                onChange={onChange}
                className="hidden"
                id={inputId}
                disabled={disabled}
              />
              <label
                htmlFor={inputId}
                className="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors cursor-pointer disabled:opacity-50"
              >
                Choose Image
              </label>
            </>
          )}
          {errors[inputId] && (
            <p className="mt-2 text-sm text-red-600">{errors[inputId]}</p>
          )}
          <p className="mt-2 text-xs text-gray-500">
            Accepted formats: JPG, PNG (max 10MB)
          </p>
        </div>
      </div>
    </div>
  );

  const onSubmitClick = (e) => {
    console.log("onSubmitClick: Beginning...");
    setFetching(true);
    setErrors({});

    let comicSubmission = {
      name: comicName,
      frontCover: frontCoverData.id,
      backCover: backCoverData.id,
    };

    postComicSubmissionCreateAPI(
      comicSubmission,
      (response) => {
        console.log(`submit success:`, response);
        if (response) {
          setForceURL("/submit/success?id=" + response.id);
        } else {
          setForceURL("/submit/success");
        }
      },
      (apiErr) => {
        console.error(`submit error:`, apiErr);
        setErrors(prev => ({
          ...prev,
          submit: apiErr?.message || "An error occurred while submitting your comic."
        }));
        window.scrollTo({ top: 0, behavior: 'smooth' });
      },
      () => {
        setFetching(false);
      },
      () => {
        window.location.href = "/login?unauthorized=true";
      },
    );
  };

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <div className="min-h-screen bg-purple-50">
      <Topbar currentPage="Submit Comic" />
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="mb-8">
          {/*
          <Link
            to="/dashboard"
            className="flex items-center text-purple-600 hover:text-purple-700 mb-4"
          >
            <ArrowLeft className="h-5 w-5 mr-1" />
            Back to Dashboard
          </Link>
          */}
          <h1
            className="text-2xl lg:text-3xl font-bold text-purple-800 mb-2"
            style={{ fontFamily: "Comic Sans MS, cursive" }}
          >
            Submit a Comic
          </h1>
          {totalSubmissions === 0 &&<p className="text-gray-600">
            Follow the steps below to submit your comic and earn ComicCoins!
          </p>}
        </div>

        {/* Provide friendly instructions if user has never submitted before. */}
        {totalSubmissions === 0 && <>
            {/* Step-by-Step Guide */}
            <div className="mb-8 p-6 rounded-lg bg-white border-2 border-purple-200">
              <h2 className="text-xl font-bold text-purple-800 mb-4">
                How It Works
              </h2>
              <div className="grid md:grid-cols-3 gap-6">
                <div className="flex flex-col items-center p-4 bg-purple-50 rounded-lg">
                  <div className="bg-purple-100 p-3 rounded-full mb-3">
                    <Camera className="h-6 w-6 text-purple-600" />
                  </div>
                  <h3 className="font-semibold text-purple-800 mb-2">
                    1. Take Photos
                  </h3>
                  <p className="text-sm text-center text-gray-600">
                    Take clear photos of your comic's front and back covers in good
                    lighting
                  </p>
                </div>
                <div className="flex flex-col items-center p-4 bg-purple-50 rounded-lg">
                  <div className="bg-purple-100 p-3 rounded-full mb-3">
                    <Upload className="h-6 w-6 text-purple-600" />
                  </div>
                  <h3 className="font-semibold text-purple-800 mb-2">
                    2. Upload Photos
                  </h3>
                  <p className="text-sm text-center text-gray-600">
                    Upload both photos and fill in the comic's name below
                  </p>
                </div>
                <div className="flex flex-col items-center p-4 bg-purple-50 rounded-lg">
                  <div className="bg-purple-100 p-3 rounded-full mb-3">
                    <AlertCircle className="h-6 w-6 text-purple-600" />
                  </div>
                  <h3 className="font-semibold text-purple-800 mb-2">
                    3. Wait for Review
                  </h3>
                  <p className="text-sm text-center text-gray-600">
                    We'll review your submission and award your ComicCoins
                  </p>
                </div>
              </div>
            </div>

            {/* Rules Section */}
            <div className="mb-8 p-4 rounded-lg border-2 border-purple-200 bg-purple-50">
              <div className="flex items-start space-x-2">
                <Info className="h-5 w-5 text-purple-600 mt-1" />
                <div>
                  <h2 className="text-purple-800 font-bold text-lg mb-2">
                    Before You Start
                  </h2>
                  <p className="text-gray-600 mb-3">
                    Please make sure you meet all these requirements:
                  </p>
                  <ul className="list-disc pl-5 space-y-2">
                    {rules.map((rule, index) => (
                      <li key={index} className="text-gray-600">
                        {rule}
                      </li>
                    ))}
                  </ul>
                </div>
              </div>
            </div>
        </>}

        {/* Submission Form */}
        <div className="bg-white rounded-xl shadow-lg p-6 mb-8 border-2 border-purple-200">

          {errors.submit && <ErrorMessage message={errors.submit} />}

          <div className="space-y-6">
            {/* Comic Name */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Comic Book Name *{" "}
                <span className="text-gray-500">(as shown on the cover)</span>
              </label>
              <input
                type="text"
                value={comicName}
                onChange={(e) => setComicName(e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                placeholder="Example: Spider-Man #1 (2022)"
              />
              <p className="mt-2 text-sm text-gray-500">
                Include the issue number and year if available
              </p>
            </div>

            {/* Photo Tips Toggle */}
            <button
              onClick={() => setShowPhotoTips(!showPhotoTips)}
              className="flex items-center space-x-2 text-purple-600 hover:text-purple-700"
            >
              <Info className="h-4 w-4" />
              <span>
                Tips for taking good photos{" "}
                {showPhotoTips ? "(hide)" : "(show)"}
              </span>
            </button>

            {/* Photo Tips Section */}
            {showPhotoTips && (
              <div className="p-4 bg-purple-50 rounded-lg text-sm text-gray-600">
                <ul className="space-y-2">
                  <li>• Use good lighting - natural daylight works best</li>
                  <li>• Place comic on a flat, solid-colored surface</li>
                  <li>• Ensure the entire cover is visible in the frame</li>
                  <li>• Avoid glare or shadows on the cover</li>
                  <li>• Make sure the image is clear and not blurry</li>
                </ul>
              </div>
            )}

            {/* Upload Sections */}
            <div className="grid md:grid-cols-2 gap-6">
              <UploadPreview
                title="Front Cover"
                fileData={frontCoverData}
                onDelete={() =>
                  handleDelete(frontCoverData, setFrontCoverData, "frontCover")
                }
                onChange={handleFrontCoverChange}
                inputId="frontCover"
                disabled={isFetching}
              />

              <UploadPreview
                title="Back Cover"
                fileData={backCoverData}
                onDelete={() =>
                  handleDelete(backCoverData, setBackCoverData, "backCover")
                }
                onChange={handleBackCoverChange}
                inputId="backCover"
                disabled={isFetching}
              />
            </div>

            {/* Action Buttons */}
            <div className="flex justify-end space-x-4">
              <Link to="/dashboard"
                className="px-6 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors"
                disabled={isFetching}
              >
                Cancel
              </Link>
              <button
                onClick={onSubmitClick}
                className="px-6 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                disabled={
                  !comicName || !frontCoverData || !backCoverData || isFetching
                }
              >
                {isFetching
                  ? "Uploading..."
                  : !comicName || !frontCoverData || !backCoverData
                    ? "Please Complete All Fields"
                    : "Submit Comic"}
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
};

export default SubmitComicPage;

const ErrorMessage = ({ message }) => (
  <div className="mb-4 p-4 bg-red-50 border-2 border-red-200 rounded-lg">
    <div className="flex gap-2">
      <AlertCircle className="h-5 w-5 text-red-500 flex-shrink-0 mt-0.5" />
      <div>
        <h3 className="text-red-800 font-medium">Submission Error</h3>
        <p className="text-red-600 text-sm mt-1">{message}</p>
      </div>
    </div>
  </div>
);
