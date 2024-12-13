import React, { useState } from 'react';
import {
  Upload, X, AlertCircle,
  ArrowLeft, Camera, Info
} from 'lucide-react';
import { Link } from "react-router-dom";

import Topbar from "../../../Components/Navigation/Topbar";
import { postAttachmentCreateAPI, deleteAttachmentAPI } from "../../../API/Attachment";


const SubmitComicPage = () => {
    const [frontCover, setFrontCover] = useState(null);
    const [backCover, setBackCover] = useState(null);
    const [frontCoverFilename, setFrontCoverFilename] = useState('');
    const [backCoverFilename, setBackCoverFilename] = useState('');
    const [frontCoverAssetID, setFrontCoverAssetID] = useState('');
    const [backCoverAssetID, setBackCoverAssetID] = useState('');
    const [comicName, setComicName] = useState('');
    const [agreed, setAgreed] = useState(false);
    const [showPhotoTips, setShowPhotoTips] = useState(false);
    const [errors, setErrors] = useState({});
    const [isFetching, setFetching] = useState(false);

    // API handlers
    const handleFileUpload = (file, setAssetID, setFilename, coverType) => {
  setFetching(true);
  setErrors({});

  const formData = new FormData();

  // Add the file with explicit content type
  formData.append('file', file, file.name);

  // Add additional metadata if needed
  formData.append('filename', file.name);
  formData.append('mimeType', file.type || 'application/octet-stream');

  postAttachmentCreateAPI(
    file.name,
    file.type || "application/octet-stream",
    formData,
    (response) => {
      console.log(`${coverType} upload success:`, response);
      setAssetID(response.id);
      setFilename(response.filename);
    },
    (apiErr) => {
      console.error(`${coverType} upload error:`, apiErr);
      setErrors(prev => ({ ...prev, [coverType]: apiErr }));
    },
    () => {
      setFetching(false);
    },
    () => {
      window.location.href = "/login?unauthorized=true";
    }
  );
};

    const handleDelete = (assetID, setAssetID, setFilename, coverType) => {
      setFetching(true);
      setErrors({});

      deleteAttachmentAPI(
        assetID,
        () => {
          setAssetID("");
          setFilename("");
        },
        (apiErr) => {
          setErrors(prev => ({ ...prev, [coverType]: apiErr }));
        },
        () => {
          setFetching(false);
        },
        () => {
          window.location.href = "/login?unauthorized=true";
        }
      );
    };

    // File change handlers
  const handleFrontCoverChange = (event) => {
    const file = event.target.files[0];
    if (file) {
      setFrontCover(file);
      handleFileUpload(file, setFrontCoverAssetID, setFrontCoverFilename, 'frontCover');
    }
  };

  const handleBackCoverChange = (event) => {
    const file = event.target.files[0];
    if (file) {
      setBackCover(file);
      handleFileUpload(file, setBackCoverAssetID, setBackCoverFilename, 'backCover');
    }
  };


  const rules = [
    "You must only upload pictures of a physical comic book",
    "You must own the comic book you are submitting",
    "You must not have submitted this comic book previously",
    "Your submission must follow our terms of service",
    "All submissions will be reviewed for approval",
    "Upon successful review, you will receive 1 ComicCoin"
  ];

  return (
    <div className="min-h-screen bg-purple-50">
      <Topbar currentPage="Submit Comic" />
      <main className="p-4 lg:p-8 max-w-5xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <Link to="/dashboard" className="flex items-center text-purple-600 hover:text-purple-700 mb-4">
            <ArrowLeft className="h-5 w-5 mr-1" />
            Back to Dashboard
          </Link>
          <h1 className="text-2xl lg:text-3xl font-bold text-purple-800 mb-2" style={{fontFamily: 'Comic Sans MS, cursive'}}>
            Submit a Comic
          </h1>
          <p className="text-gray-600">Follow the steps below to submit your comic and earn ComicCoins!</p>
        </div>

        {/* Step-by-Step Guide */}
        <div className="mb-8 p-6 rounded-lg bg-white border-2 border-purple-200">
          <h2 className="text-xl font-bold text-purple-800 mb-4">How It Works</h2>
          <div className="grid md:grid-cols-3 gap-6">
            <div className="flex flex-col items-center p-4 bg-purple-50 rounded-lg">
              <div className="bg-purple-100 p-3 rounded-full mb-3">
                <Camera className="h-6 w-6 text-purple-600" />
              </div>
              <h3 className="font-semibold text-purple-800 mb-2">1. Take Photos</h3>
              <p className="text-sm text-center text-gray-600">Take clear photos of your comic's front and back covers in good lighting</p>
            </div>
            <div className="flex flex-col items-center p-4 bg-purple-50 rounded-lg">
              <div className="bg-purple-100 p-3 rounded-full mb-3">
                <Upload className="h-6 w-6 text-purple-600" />
              </div>
              <h3 className="font-semibold text-purple-800 mb-2">2. Upload Photos</h3>
              <p className="text-sm text-center text-gray-600">Upload both photos and fill in the comic's name below</p>
            </div>
            <div className="flex flex-col items-center p-4 bg-purple-50 rounded-lg">
              <div className="bg-purple-100 p-3 rounded-full mb-3">
                <AlertCircle className="h-6 w-6 text-purple-600" />
              </div>
              <h3 className="font-semibold text-purple-800 mb-2">3. Wait for Review</h3>
              <p className="text-sm text-center text-gray-600">We'll review your submission and award your ComicCoins</p>
            </div>
          </div>
        </div>

        {/* Rules Section */}
        <div className="mb-8 p-4 rounded-lg border-2 border-purple-200 bg-purple-50">
          <div className="flex items-start space-x-2">
            <Info className="h-5 w-5 text-purple-600 mt-1" />
            <div>
              <h2 className="text-purple-800 font-bold text-lg mb-2">Before You Start</h2>
              <p className="text-gray-600 mb-3">Please make sure you meet all these requirements:</p>
              <ul className="list-disc pl-5 space-y-2">
                {rules.map((rule, index) => (
                  <li key={index} className="text-gray-600">{rule}</li>
                ))}
              </ul>
            </div>
          </div>
        </div>

        {/* Submission Form */}
        <div className="bg-white rounded-xl shadow-lg p-6 mb-8 border-2 border-purple-200">
          <div className="space-y-6">
            {/* Comic Name */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Comic Book Name * <span className="text-gray-500">(as shown on the cover)</span>
              </label>
              <input
                type="text"
                value={comicName}
                onChange={(e) => setComicName(e.target.value)}
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                placeholder="Example: Spider-Man #1 (2022)"
              />
              <p className="mt-2 text-sm text-gray-500">Include the issue number and year if available</p>
            </div>

            {/* Photo Tips Toggle */}
            <button
              onClick={() => setShowPhotoTips(!showPhotoTips)}
              className="flex items-center space-x-2 text-purple-600 hover:text-purple-700"
            >
              <Info className="h-4 w-4" />
              <span>Tips for taking good photos {showPhotoTips ? '(hide)' : '(show)'}</span>
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
            {/* Front Cover Upload */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Front Cover * <span className="text-gray-500">(required)</span>
        </label>
        <div className="border-2 border-dashed border-purple-200 rounded-lg p-6 hover:border-purple-400 transition-colors">
          <div className="flex flex-col items-center">
            {frontCoverFilename ? (
              <div className="w-full">
                <div className="flex items-center justify-between mb-4">
                  <span className="text-sm text-gray-600">{frontCoverFilename}</span>
                  <button
                    onClick={() => handleDelete(frontCoverAssetID, setFrontCoverAssetID, setFrontCoverFilename, 'frontCover')}
                    className="text-red-500 hover:text-red-700"
                    disabled={isFetching}
                  >
                    <X className="h-5 w-5" />
                  </button>
                </div>
              </div>
            ) : (
              <>
                <Upload className="h-12 w-12 text-purple-400 mb-4" />
                <p className="text-sm text-gray-500 text-center mb-4">
                  Click here to upload or drag and drop your front cover photo
                </p>
                <input
                  type="file"
                  accept="image/*"
                  onChange={handleFrontCoverChange}
                  className="hidden"
                  id="frontCoverUpload"
                  disabled={isFetching}
                />
                <label
                  htmlFor="frontCoverUpload"
                  className="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors cursor-pointer disabled:opacity-50"
                >
                  Choose Front Cover
                </label>
              </>
            )}
            {errors.frontCover && (
              <p className="mt-2 text-sm text-red-600">{errors.frontCover}</p>
            )}
            <p className="mt-2 text-xs text-gray-500">Accepted formats: JPG, PNG (max 10MB)</p>
          </div>
        </div>
      </div>

      {/* Back Cover Upload */}
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-2">
          Back Cover * <span className="text-gray-500">(required)</span>
        </label>
        <div className="border-2 border-dashed border-purple-200 rounded-lg p-6 hover:border-purple-400 transition-colors">
          <div className="flex flex-col items-center">
            {backCoverFilename ? (
              <div className="w-full">
                <div className="flex items-center justify-between mb-4">
                  <span className="text-sm text-gray-600">{backCoverFilename}</span>
                  <button
                    onClick={() => handleDelete(backCoverAssetID, setBackCoverAssetID, setBackCoverFilename, 'backCover')}
                    className="text-red-500 hover:text-red-700"
                    disabled={isFetching}
                  >
                    <X className="h-5 w-5" />
                  </button>
                </div>
              </div>
            ) : (
              <>
                <Upload className="h-12 w-12 text-purple-400 mb-4" />
                <p className="text-sm text-gray-500 text-center mb-4">
                  Click here to upload or drag and drop your back cover photo
                </p>
                <input
                  type="file"
                  accept="image/*"
                  onChange={handleBackCoverChange}
                  className="hidden"
                  id="backCoverUpload"
                  disabled={isFetching}
                />
                <label
                  htmlFor="backCoverUpload"
                  className="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors cursor-pointer disabled:opacity-50"
                >
                  Choose Back Cover
                </label>
              </>
            )}
            {errors.backCover && (
              <p className="mt-2 text-sm text-red-600">{errors.backCover}</p>
            )}
            <p className="mt-2 text-xs text-gray-500">Accepted formats: JPG, PNG (max 10MB)</p>
          </div>
        </div>
      </div>
            </div>

            {/* Terms Agreement */}
            <div className="flex items-start space-x-3 bg-purple-50 p-4 rounded-lg">
              <input
                type="checkbox"
                checked={agreed}
                onChange={(e) => setAgreed(e.target.checked)}
                className="mt-1 h-4 w-4 text-purple-600 border-gray-300 rounded focus:ring-purple-500"
              />
              <div>
                <label className="text-sm text-gray-600">
                  I confirm that:
                </label>
                <ul className="text-sm text-gray-600 mt-1 list-disc pl-5">
                  <li>I own this comic book</li>
                  <li>I haven't submitted it before</li>
                  <li>I agree to the submission rules and terms of service</li>
                </ul>
              </div>
            </div>

            {/* Action Buttons */}
            <div className="flex justify-end space-x-4">
              <button className="px-6 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors">
                Cancel
              </button>
              <button
                className="px-6 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                disabled={!comicName || !frontCover || !backCover || !agreed}
              >
                {!comicName || !frontCover || !backCover || !agreed ?
                  'Please Complete All Fields' : 'Submit Comic'}
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
};

export default SubmitComicPage;
