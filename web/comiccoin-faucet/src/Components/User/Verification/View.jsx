import React, { useState } from 'react';
import {
  ArrowLeft,
  Info,
  MapPin,
  Truck,
  HelpCircle,
  ChevronLeft,
  AlertCircle
} from 'lucide-react';
import { Link, Navigate } from "react-router-dom";
import { useRecoilState } from "recoil";

import Topbar from "../../../Components/Navigation/Topbar";
import { currentUserState } from "../../../AppState";
import { putProfileApplyForVerificationAPI } from "../../../API/Profile";


const ApplyForVerificationPage = () => {

  // Variable controls the global state of the app.
  const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

  const [formData, setFormData] = useState({
    phone: '',
    country: currentUser.country,
    region: '',
    city: '',
    postalCode: '',
    addressLine1: '',
    addressLine2: '',
    hasShippingAddress: false,
    shippingName: '',
    shippingPhone: '',
    shippingCountry: '',
    shippingRegion: '',
    shippingCity: '',
    shippingPostalCode: '',
    shippingAddressLine1: '',
    shippingAddressLine2: '',
    howDidYouHearAboutUs: 0,
    howDidYouHearAboutUsOther: '',
    howLongCollectingComicBooksForGrading: 0,
    hasPreviouslySubmittedComicBookForGrading: 0,
    hasOwnedGradedComicBooks: 0,
    hasRegularComicBookShop: 0,
    hasPreviouslyPurchasedFromAuctionSite: 0,
    hasPreviouslyPurchasedFromFacebookMarketplace: 0,
    hasRegularlyAttendedComicConsOrCollectibleShows: 0
  });

  const [forceURL, setForceURL] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [hasSubmitted, setHasSubmitted] = useState(false);
  const [errors, setErrors] = useState({});

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
    if (errors[name]) {
      setErrors(prev => ({
        ...prev,
        [name]: ''
      }));
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    setHasSubmitted(true);

    putProfileApplyForVerificationAPI(
      formData,
      (resp) => {
        // For debugging purposes only.
        console.log("putProfileApplyForVerificationAPI: Starting...");
        console.log(resp);

        // Update the user profile.
        setCurrentUser(resp);

        // Redirect the user to a new page.
        setForceURL("/settings");
      },
      (apiErr) => {
        console.log("putProfileApplyForVerificationAPI: apiErr:", apiErr);
        setErrors(apiErr);
        window.scrollTo({ top: 0, behavior: 'smooth' });
      },
      () => {
        console.log("putProfileApplyForVerificationAPI: Starting...");

      },
    );
    setHasSubmitted(true);
  };


  const howDidYouHearOptions = [
    { value: 2, label: "My local comic book shop" },
    { value: 3, label: "CPS website" },
    { value: 4, label: "Comic Con booth" },
    { value: 5, label: "Friend" },
    { value: 6, label: "Social media" },
    { value: 7, label: "Blog post article" },
    { value: 1, label: "Other (Please specify)" },
  ];

  const experienceOptions = [
    { value: 1, label: "1 year" },
    { value: 2, label: "2-5 years" },
    { value: 3, label: "5-9 years" },
    { value: 4, label: "10+ years" },
  ];

  const yesNoOptions = [
    { value: 1, label: 'Yes' },
    { value: 2, label: 'No' }
  ];

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <div className="min-h-screen bg-purple-50">
      <Topbar currentPage="Settings" />
      <div className="p-8">
          <div className="max-w-2xl mx-auto">
            <nav className="mb-8">
              <Link
                to={"/settings"}
                className="group flex items-center text-sm text-gray-600 hover:text-purple-600"
              >
                <ChevronLeft className="w-4 h-4 mr-1 transition-transform group-hover:-translate-x-1" />
                Back to Settings
              </Link>
            </nav>

            <h1 className="text-3xl font-bold text-purple-800 mb-8" style={{fontFamily: 'Comic Sans MS, cursive'}}>
              Apply for Verification
            </h1>
          </div>

          <div className="max-w-2xl mx-auto">
            {hasSubmitted && Object.keys(errors).length > 0 && (
              <div className="mb-6 bg-red-50 border-l-4 border-red-500 p-4 rounded-r-lg">
                <div className="flex">
                  <div className="flex-shrink-0">
                    <AlertCircle className="h-5 w-5 text-red-400" />
                  </div>
                  <div className="ml-3">
                    <h3 className="text-sm font-medium text-red-800">
                      Please correct the following errors:
                    </h3>
                    <div className="mt-2 text-sm text-red-700">
                      <ul className="list-disc space-y-1 pl-5">
                        {Object.values(errors)
                          .filter(Boolean)
                          .map((error, index) => (
                            <li key={index}>{error}</li>
                          ))}
                      </ul>
                    </div>
                  </div>
                </div>
              </div>
            )}
          </div>

          <div className="max-w-2xl mx-auto">
            {/* Main Form */}
            <div className="bg-white rounded-xl shadow-lg p-6 mb-8 border-2 border-purple-200 space-y-8">
            <p className="text-gray-600">Complete this form to get verified and start earning ComicCoins!</p>
              {/* Contact Information */}
              <section>
                <h2 className="text-xl font-bold text-purple-800 mb-4 flex items-center">
                  <Info className="h-5 w-5 mr-2" />
                  Contact Information
                </h2>
                <div className="grid md:grid-cols-2 gap-4">
                  <div>
                    <label htmlFor="phone" className="block text-sm font-medium text-gray-700 mb-1">
                      Phone Number *
                    </label>
                    <input
                      id="phone"
                      name="phone"
                      type="tel"
                      className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                        errors.phone ? "border-red-500" : "border-gray-300"
                      }`}
                      placeholder="+1 (555) 555-5555"
                      value={formData.phone}
                      onChange={handleChange}
                    />
                    {errors.phone && (
                      <p className="mt-1 text-sm text-red-600">{errors.phone}</p>
                    )}
                  </div>
                </div>
              </section>

              {/* Primary Address */}
              <section>
                <h2 className="text-xl font-bold text-purple-800 mb-4 flex items-center">
                  <MapPin className="h-5 w-5 mr-2" />
                  Primary Address
                </h2>
                <div className="grid md:grid-cols-2 gap-4">
                  <div>
                    <label htmlFor="country" className="block text-sm font-medium text-gray-700 mb-1">
                      Country *
                    </label>
                    <select
                      id="country"
                      name="country"
                      value={formData.country}
                      onChange={handleChange}
                      className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent bg-white ${
                        errors.country ? "border-red-500" : "border-gray-300"
                      }`}
                    >
                      <option value="">Select your country</option>
                      <option value="Canada">Canada</option>
                      <option value="United States">United States</option>
                      <option value="Mexico">Mexico</option>
                      <option value="Other">Other</option>
                    </select>
                    {errors.country && (
                      <p className="mt-1 text-sm text-red-600">{errors.country}</p>
                    )}
                  </div>
                  <div>
                    <label htmlFor="region" className="block text-sm font-medium text-gray-700 mb-1">
                      Region/State *
                    </label>
                    <input
                      id="region"
                      name="region"
                      type="text"
                      className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                        errors.region ? "border-red-500" : "border-gray-300"
                      }`}
                      value={formData.region}
                      onChange={handleChange}
                    />
                    {errors.region && (
                      <p className="mt-1 text-sm text-red-600">{errors.region}</p>
                    )}
                  </div>
                  <div>
                    <label htmlFor="city" className="block text-sm font-medium text-gray-700 mb-1">
                      City *
                    </label>
                    <input
                      id="city"
                      name="city"
                      type="text"
                      className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                        errors.city ? "border-red-500" : "border-gray-300"
                      }`}
                      value={formData.city}
                      onChange={handleChange}
                    />
                    {errors.city && (
                      <p className="mt-1 text-sm text-red-600">{errors.city}</p>
                    )}
                  </div>
                  <div>
                    <label htmlFor="postalCode" className="block text-sm font-medium text-gray-700 mb-1">
                      Postal Code *
                    </label>
                    <input
                      id="postalCode"
                      name="postalCode"
                      type="text"
                      className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                        errors.postalCode ? "border-red-500" : "border-gray-300"
                      }`}
                      value={formData.postalCode}
                      onChange={handleChange}
                    />
                    {errors.postalCode && (
                      <p className="mt-1 text-sm text-red-600">{errors.postalCode}</p>
                    )}
                  </div>
                  <div className="md:col-span-2">
                    <label htmlFor="addressLine1" className="block text-sm font-medium text-gray-700 mb-1">
                      Address Line 1 *
                    </label>
                    <input
                      id="addressLine1"
                      name="addressLine1"
                      type="text"
                      className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                        errors.addressLine1 ? "border-red-500" : "border-gray-300"
                      }`}
                      value={formData.addressLine1}
                      onChange={handleChange}
                    />
                    {errors.addressLine1 && (
                      <p className="mt-1 text-sm text-red-600">{errors.addressLine1}</p>
                    )}
                  </div>
                  <div className="md:col-span-2">
                    <label htmlFor="addressLine2" className="block text-sm font-medium text-gray-700 mb-1">
                      Address Line 2
                    </label>
                    <input
                      id="addressLine2"
                      name="addressLine2"
                      type="text"
                      className="w-full h-11 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                      value={formData.addressLine2}
                      onChange={handleChange}
                    />
                  </div>
                </div>
              </section>

              {/* Shipping Address */}
              <section>
                <div className="flex items-center justify-between mb-4">
                  <h2 className="text-xl font-bold text-purple-800 flex items-center">
                    <Truck className="h-5 w-5 mr-2" />
                    Shipping Address
                  </h2>
                  <label className="flex items-center space-x-2">
                    <input
                      type="checkbox"
                      name="hasShippingAddress"
                      checked={!formData.hasShippingAddress}
                      onChange={(e) => handleChange({
                        target: {
                          name: 'hasShippingAddress',
                          type: 'checkbox',
                          checked: !e.target.checked
                        }
                      })}
                      className="h-4 w-4 text-purple-600 border-gray-300 rounded focus:ring-purple-500"
                    />
                    <span className="text-sm text-gray-600">Same as primary address</span>
                  </label>
                </div>
                {formData.hasShippingAddress && (
                  <div className="grid md:grid-cols-2 gap-4">
                    <div>
                      <label htmlFor="shippingName" className="block text-sm font-medium text-gray-700 mb-1">
                        Full Name *
                      </label>
                      <input
                        id="shippingName"
                        name="shippingName"
                        type="text"
                        className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                          errors.shippingName ? "border-red-500" : "border-gray-300"
                        }`}
                        value={formData.shippingName}
                        onChange={handleChange}
                      />
                      {errors.shippingName && (
                        <p className="mt-1 text-sm text-red-600">{errors.shippingName}</p>
                      )}
                    </div>
                    <div>
                      <label htmlFor="shippingPhone" className="block text-sm font-medium text-gray-700 mb-1">
                        Phone Number *
                      </label>
                      <input
                        id="shippingPhone"
                        name="shippingPhone"
                        type="tel"
                        className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                          errors.shippingPhone ? "border-red-500" : "border-gray-300"
                        }`}
                        value={formData.shippingPhone}
                        onChange={handleChange}
                      />
                      {errors.shippingPhone && (
                        <p className="mt-1 text-sm text-red-600">{errors.shippingPhone}</p>
                      )}
                    </div>
                    {/* Replicate the same address fields as primary address */}
                    <div>
                      <label htmlFor="shippingCountry" className="block text-sm font-medium text-gray-700 mb-1">
                        Country *
                      </label>
                      <select
                        id="shippingCountry"
                        name="shippingCountry"
                        value={formData.shippingCountry}
                        onChange={handleChange}
                        className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent bg-white ${
                          errors.shippingCountry ? "border-red-500" : "border-gray-300"
                        }`}
                      >
                        <option value="">Select your country</option>
                        <option value="Canada">Canada</option>
                        <option value="United States">United States</option>
                        <option value="Mexico">Mexico</option>
                        <option value="Other">Other</option>
                      </select>
                      {errors.shippingCountry && (
                        <p className="mt-1 text-sm text-red-600">{errors.shippingCountry}</p>
                      )}
                    </div>
                    <div>
                      <label htmlFor="shippingRegion" className="block text-sm font-medium text-gray-700 mb-1">
                        Region/State *
                      </label>
                      <input
                        id="shippingRegion"
                        name="shippingRegion"
                        type="text"
                        className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                          errors.shippingRegion ? "border-red-500" : "border-gray-300"
                        }`}
                        value={formData.shippingRegion}
                        onChange={handleChange}
                      />
                      {errors.shippingRegion && (
                        <p className="mt-1 text-sm text-red-600">{errors.shippingRegion}</p>
                      )}
                    </div>
                    <div>
                      <label htmlFor="shippingCity" className="block text-sm font-medium text-gray-700 mb-1">
                        City *
                      </label>
                      <input
                        id="shippingCity"
                        name="shippingCity"
                        type="text"
                        className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                          errors.shippingCity ? "border-red-500" : "border-gray-300"
                        }`}
                        value={formData.shippingCity}
                        onChange={handleChange}
                      />
                      {errors.shippingCity && (
                        <p className="mt-1 text-sm text-red-600">{errors.shippingCity}</p>
                      )}
                    </div>
                    <div>
                      <label htmlFor="shippingPostalCode" className="block text-sm font-medium text-gray-700 mb-1">
                        Postal Code *
                      </label>
                      <input
                        id="shippingPostalCode"
                        name="shippingPostalCode"
                        type="text"
                        className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                          errors.shippingPostalCode ? "border-red-500" : "border-gray-300"
                        }`}
                        value={formData.shippingPostalCode}
                        onChange={handleChange}
                      />
                      {errors.shippingPostalCode && (
                        <p className="mt-1 text-sm text-red-600">{errors.shippingPostalCode}</p>
                      )}
                    </div>
                    <div>
                      <label htmlFor="shippingAddressLine1" className="block text-sm font-medium text-gray-700 mb-1">
                        Address Line 1 *
                      </label>
                      <input
                        id="shippingAddressLine1"
                        name="shippingAddressLine1"
                        type="text"
                        className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                          errors.shippingAddressLine1 ? "border-red-500" : "border-gray-300"
                        }`}
                        value={formData.shippingAddressLine1}
                        onChange={handleChange}
                      />
                      {errors.shippingAddressLine1 && (
                        <p className="mt-1 text-sm text-red-600">{errors.shippingAddressLine1}</p>
                      )}
                    </div>
                    <div>
                      <label htmlFor="shippingAddressLine2" className="block text-sm font-medium text-gray-700 mb-1">
                        Address Line 2
                      </label>
                      <input
                        id="shippingAddressLine2"
                        name="shippingAddressLine2"
                        type="text"
                        className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                          errors.shippingAddressLine2 ? "border-red-500" : "border-gray-300"
                        }`}
                        value={formData.shippingAddressLine2}
                        onChange={handleChange}
                      />
                      {errors.shippingAddressLine2 && (
                        <p className="mt-1 text-sm text-red-600">{errors.shippingAddressLine2}</p>
                      )}
                    </div>

                  </div>
                )}
              </section>

              <section>
                <h2 className="text-xl font-bold text-purple-800 mb-4 flex items-center">
                  <HelpCircle className="h-5 w-5 mr-2" />
                  Comic Collecting Experience
                </h2>
                <div className="space-y-6">
                  <div>
                    <label htmlFor="howDidYouHearAboutUs" className="block text-sm font-medium text-gray-700 mb-1">
                      How did you hear about us? *
                    </label>
                    <select
                      id="howDidYouHearAboutUs"
                      name="howDidYouHearAboutUs"
                      value={formData.howDidYouHearAboutUs}
                      onChange={handleChange}
                      className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent bg-white ${
                        errors.howDidYouHearAboutUs ? "border-red-500" : "border-gray-300"
                      }`}
                    >
                      <option value="">Select an option</option>
                      {howDidYouHearOptions.map(option => (
                        <option key={option.value} value={option.value}>{option.label}</option>
                      ))}
                    </select>
                    {errors.howDidYouHearAboutUs && (
                      <p className="mt-1 text-sm text-red-600">{errors.howDidYouHearAboutUs}</p>
                    )}
                  </div>

                  {formData.howDidYouHearAboutUs === "other" && (
                    <div>
                      <label htmlFor="howDidYouHearAboutUsOther" className="block text-sm font-medium text-gray-700 mb-1">
                        Please specify how you heard about us *
                      </label>
                      <input
                        id="howDidYouHearAboutUsOther"
                        name="howDidYouHearAboutUsOther"
                        type="text"
                        className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                          errors.howDidYouHearAboutUsOther ? "border-red-500" : "border-gray-300"
                        }`}
                        value={formData.howDidYouHearAboutUsOther}
                        onChange={handleChange}
                      />
                      {errors.howDidYouHearAboutUsOther && (
                        <p className="mt-1 text-sm text-red-600">{errors.howDidYouHearAboutUsOther}</p>
                      )}
                    </div>
                  )}

                  <div>
                    <label htmlFor="howLongCollecting" className="block text-sm font-medium text-gray-700 mb-1">
                      How long have you been collecting comic books for grading? *
                    </label>
                    <select
                      id="howLongCollecting"
                      name="howLongCollectingComicBooksForGrading"
                      value={formData.howLongCollectingComicBooksForGrading}
                      onChange={handleChange}
                      className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent bg-white ${
                        errors.howLongCollectingComicBooksForGrading ? "border-red-500" : "border-gray-300"
                      }`}
                    >
                      <option value="">Select an option</option>
                      {experienceOptions.map(option => (
                        <option key={option.value} value={option.value}>{option.label}</option>
                      ))}
                    </select>
                    {errors.howLongCollectingComicBooksForGrading && (
                      <p className="mt-1 text-sm text-red-600">{errors.howLongCollectingComicBooksForGrading}</p>
                    )}
                  </div>

                  <div>
                    <label htmlFor="hasPreviouslySubmitted" className="block text-sm font-medium text-gray-700 mb-1">
                      Have you previously submitted comic books for grading? *
                    </label>
                    <select
                      id="hasPreviouslySubmitted"
                      name="hasPreviouslySubmittedComicBookForGrading"
                      value={formData.hasPreviouslySubmittedComicBookForGrading}
                      onChange={handleChange}
                      className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent bg-white ${
                        errors.hasPreviouslySubmittedComicBookForGrading ? "border-red-500" : "border-gray-300"
                      }`}
                    >
                      <option value="">Select an option</option>
                      {yesNoOptions.map(option => (
                        <option key={option.value} value={option.value}>{option.label}</option>
                      ))}
                    </select>
                    {errors.hasPreviouslySubmittedComicBookForGrading && (
                      <p className="mt-1 text-sm text-red-600">{errors.hasPreviouslySubmittedComicBookForGrading}</p>
                    )}
                  </div>

                  <div>
                    <label htmlFor="hasOwnedGradedComicBooks" className="block text-sm font-medium text-gray-700 mb-1">
                      Have you owned graded comic books? *
                    </label>
                    <select
                      id="hasOwnedGradedComicBooks"
                      name="hasOwnedGradedComicBooks"
                      value={formData.hasOwnedGradedComicBooks}
                      onChange={handleChange}
                      className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent bg-white ${
                        errors.hasOwnedGradedComicBooks ? "border-red-500" : "border-gray-300"
                      }`}
                    >
                      <option value="">Select an option</option>
                      {yesNoOptions.map(option => (
                        <option key={option.value} value={option.value}>{option.label}</option>
                      ))}
                    </select>
                    {errors.hasOwnedGradedComicBooks && (
                      <p className="mt-1 text-sm text-red-600">{errors.hasOwnedGradedComicBooks}</p>
                    )}
                  </div>

                  <div>
                    <label htmlFor="hasRegularComicBookShop" className="block text-sm font-medium text-gray-700 mb-1">
                      Do you have a regular comic book shop? *
                    </label>
                    <select
                      id="hasRegularShop"
                      name="hasRegularComicBookShop"
                      value={formData.hasRegularComicBookShop}
                      onChange={handleChange}
                      className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent bg-white ${
                        errors.hasRegularComicBookShop ? "border-red-500" : "border-gray-300"
                      }`}
                    >
                      <option value="">Select an option</option>
                      {yesNoOptions.map(option => (
                        <option key={option.value} value={option.value}>{option.label}</option>
                      ))}
                    </select>
                    {errors.hasRegularComicBookShop && (
                      <p className="mt-1 text-sm text-red-600">{errors.hasRegularComicBookShop}</p>
                    )}
                  </div>

                  <div>
                    <label htmlFor="hasPreviouslyPurchasedFromAuctionSite" className="block text-sm font-medium text-gray-700 mb-1">
                      Have you previously purchased from auction sites? *
                    </label>
                    <select
                      id="hasAuctionExperience"
                      name="hasPreviouslyPurchasedFromAuctionSite"
                      value={formData.hasPreviouslyPurchasedFromAuctionSite}
                      onChange={handleChange}
                      className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent bg-white ${
                        errors.hasPreviouslyPurchasedFromAuctionSite ? "border-red-500" : "border-gray-300"
                      }`}
                    >
                      <option value="">Select an option</option>
                      {yesNoOptions.map(option => (
                        <option key={option.value} value={option.value}>{option.label}</option>
                      ))}
                    </select>
                    {errors.hasPreviouslyPurchasedFromAuctionSite && (
                      <p className="mt-1 text-sm text-red-600">{errors.hasPreviouslyPurchasedFromAuctionSite}</p>
                    )}
                  </div>

                  <div>
                    <label htmlFor="hasPreviouslyPurchasedFromFacebookMarketplace" className="block text-sm font-medium text-gray-700 mb-1">
                      Have you previously purchased from Facebook Marketplace? *
                    </label>
                    <select
                      id="hasFacebookExperience"
                      name="hasPreviouslyPurchasedFromFacebookMarketplace"
                      value={formData.hasPreviouslyPurchasedFromFacebookMarketplace}
                      onChange={handleChange}
                      className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent bg-white ${
                        errors.hasPreviouslyPurchasedFromFacebookMarketplace ? "border-red-500" : "border-gray-300"
                      }`}
                    >
                      <option value="">Select an option</option>
                      {yesNoOptions.map(option => (
                        <option key={option.value} value={option.value}>{option.label}</option>
                      ))}
                    </select>
                    {errors.hasPreviouslyPurchasedFromFacebookMarketplace && (
                      <p className="mt-1 text-sm text-red-600">{errors.hasPreviouslyPurchasedFromFacebookMarketplace}</p>
                    )}
                  </div>

                  <div>
                    <label htmlFor="hasRegularlyAttendedComicConsOrCollectibleShows" className="block text-sm font-medium text-gray-700 mb-1">
                      Have you regularly attended comic cons or collectible shows? *
                    </label>
                    <select
                      id="hasConventionExperience"
                      name="hasRegularlyAttendedComicConsOrCollectibleShows"
                      value={formData.hasRegularlyAttendedComicConsOrCollectibleShows}
                      onChange={handleChange}
                      className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent bg-white ${
                        errors.hasRegularlyAttendedComicConsOrCollectibleShows ? "border-red-500" : "border-gray-300"
                      }`}
                    >
                      <option value="">Select an option</option>
                      {yesNoOptions.map(option => (
                        <option key={option.value} value={option.value}>{option.label}</option>
                      ))}
                    </select>
                    {errors.hasRegularlyAttendedComicConsOrCollectibleShows && (
                      <p className="mt-1 text-sm text-red-600">{errors.hasRegularlyAttendedComicConsOrCollectibleShows}</p>
                    )}
                  </div>


                </div>
              </section>

              {/* Submit Button */}
              <div className="flex justify-end space-x-4">
                <Link to="/settings" className="px-6 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors">
                  Cancel
                </Link>
                <button onClick={handleSubmit} className="px-6 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors">
                  Submit Application
                </button>
              </div>
            </div>
          </div>
      </div>
    </div>
  );
};

export default ApplyForVerificationPage;
