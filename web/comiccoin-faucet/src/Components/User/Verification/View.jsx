import React, { useState } from 'react';
import {
  ArrowLeft,
  Info,
  MapPin,
  Truck,
  HelpCircle
} from 'lucide-react';
import { Link } from "react-router-dom";

import Topbar from "../../../Components/Navigation/Topbar";


const ApplyForVerificationPage = () => {
  const [formData, setFormData] = useState({
    phone: '',
    country: '',
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
    howDidYouHearAboutUs: '',
    howDidYouHearAboutUsOther: '',
    howLongCollectingComicBooksForGrading: '',
    hasPreviouslySubmittedComicBookForGrading: '',
    hasOwnedGradedComicBooks: '',
    hasRegularComicBookShop: '',
    hasPreviouslyPurchasedFromAuctionSite: '',
    hasPreviouslyPurchasedFromFacebookMarketplace: '',
    hasRegularlyAttendedComicConsOrCollectibleShows: ''
  });

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

  const howDidYouHearOptions = [
    { value: "social_media", label: 'Social Media' },
    { value: "friend_family", label: 'Friend/Family' },
    { value: "comic_shop", label: 'Comic Shop' },
    { value: "convention", label: 'Convention' },
    { value: "other", label: 'Other' }
  ];

  const experienceOptions = [
    { value: "less_than_1", label: 'Less than 1 year' },
    { value: "1_to_5", label: '1-5 years' },
    { value: "5_to_10", label: '5-10 years' },
    { value: "more_than_10", label: 'More than 10 years' }
  ];

  const yesNoOptions = [
    { value: "yes", label: 'Yes' },
    { value: "no", label: 'No' }
  ];

  return (
    <div className="min-h-screen bg-purple-50">
      <Topbar currentPage="Settings" />
      <main className="p-4 lg:p-8 max-w-5xl mx-auto">
        {/* Header */}
        <div className="mb-8">
          <Link to="/dashboard" className="flex items-center text-purple-600 hover:text-purple-700 mb-4">
            <ArrowLeft className="h-5 w-5 mr-1" />
            Back to Dashboard
          </Link>
          <h1 className="text-2xl lg:text-3xl font-bold text-purple-800 mb-2">
            Apply for Verification
          </h1>
          <p className="text-gray-600">Complete this form to get verified and start earning ComicCoins!</p>
        </div>

        {/* Main Form */}
        <div className="bg-white rounded-xl shadow-lg p-6 mb-8 border-2 border-purple-200 space-y-8">
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
                {/* Continue with other shipping address fields... */}
              </div>
            )}
          </section>

          {/* Experience Questions */}
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

              {/* Additional Experience Questions */}
              {[
                {
                  id: "hasPreviouslySubmitted",
                  name: "hasPreviouslySubmittedComicBookForGrading",
                  label: "Have you previously submitted comic books for grading?"
                },
                {
                  id: "hasOwnedGraded",
                  name: "hasOwnedGradedComicBooks",
                  label: "Have you owned graded comic books?"
                },
                {
                  id: "hasRegularShop",
                  name: "hasRegularComicBookShop",
                  label: "Do you have a regular comic book shop?"
                },
                {
                  id: "hasAuctionExperience",
                  name: "hasPreviouslyPurchasedFromAuctionSite",
                  label: "Have you previously purchased from auction sites?"
                },
                {
                  id: "hasFacebookExperience",
                  name: "hasPreviouslyPurchasedFromFacebookMarketplace",
                  label: "Have you previously purchased from Facebook Marketplace?"
                },
                {
                  id: "hasConventionExperience",
                  name: "hasRegularlyAttendedComicConsOrCollectibleShows",
                  label: "Have you regularly attended comic cons or collectible shows?"
                }
              ].map((question) => (
                <div key={question.id}>
                  <label htmlFor={question.id} className="block text-sm font-medium text-gray-700 mb-1">
                    {question.label} *
                  </label>
                  <select
                    id={question.id}
                    name={question.name}
                    value={formData[question.name]}
                    onChange={handleChange}
                    className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent bg-white ${
                      errors[question.name] ? "border-red-500" : "border-gray-300"
                    }`}
                  >
                    <option value="">Select an option</option>
                    {yesNoOptions.map(option => (
                      <option key={option.value} value={option.value}>{option.label}</option>
                    ))}
                  </select>
                  {errors[question.name] && (
                    <p className="mt-1 text-sm text-red-600">{errors[question.name]}</p>
                  )}
                </div>
              ))}
            </div>
          </section>

          {/* Submit Button */}
          <div className="flex justify-end space-x-4">
            <Link to="/dashboard" className="px-6 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors">
              Cancel
            </Link>
            <button className="px-6 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors">
              Submit Application
            </button>
          </div>
        </div>
      </main>
    </div>
  );
};

export default ApplyForVerificationPage;
