import React, { useState } from 'react';
import { ChevronLeft, Mail, Globe, Clock, AlertCircle } from 'lucide-react';
import { Navigate } from 'react-router-dom';
import { useRecoilState } from "recoil";

import Topbar from "../../../Components/Navigation/Topbar";
import { currentUserState } from "../../../AppState";
import { putProfileUpdateAPI } from "../../../API/Profile";


const EmailSettingsPage = () => {

  // Variable controls the global state of the app.
  const [currentUser] = useRecoilState(currentUserState);

  const [forceURL, setForceURL] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [hasSubmitted, setHasSubmitted] = useState(false);
  const [errors, setErrors] = useState({});

  const [formData, setFormData] = useState({
    email: currentUser.email,
    firstName: currentUser.firstName,
    lastName: currentUser.lastName,
    country: ['United States', 'Canada', 'Mexico'].includes(currentUser.country) ? currentUser.country : 'Other',
    countryOther: currentUser.country,
    timezone: currentUser.timezone,
    agreeTermsOfService: currentUser.agreeTermsOfService,
    agreePromotional: currentUser.agreePromotional
  });

  const validateField = (name, value) => {
    switch (name) {
      case "firstName":
        if (!value.trim()) return "First name is required";
        if (value.length < 2) return "First name must be at least 2 characters";
        if (!/^[a-zA-Z\s-']+$/.test(value))
          return "First name can only contain letters, spaces, hyphens, and apostrophes";
        return "";

      case "lastName":
        if (!value.trim()) return "Last name is required";
        if (value.length < 2) return "Last name must be at least 2 characters";
        if (!/^[a-zA-Z\s-']+$/.test(value))
          return "Last name can only contain letters, spaces, hyphens, and apostrophes";
        return "";

      case "email":
        if (!value) return "Email is required";
        if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value))
          return "Please enter a valid email address";
        return "";

      case "country":
        if (!value) return "Please select your country";
        return "";

      case "agreeTermsOfService":
        if (!value) return "You must agree to the Terms of Service";
        return "";

      default:
        return "";
    }
  };

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    const newValue = type === "checkbox" ? checked : value;

    setFormData(prev => ({
      ...prev,
      [name]: newValue
    }));

    if (hasSubmitted) {
      const error = validateField(name, newValue);
      setErrors(prev => ({
        ...prev,
        [name]: error
      }));
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    setHasSubmitted(true);

    // Validate all fields
    const newErrors = {};
    Object.keys(formData).forEach(key => {
      const error = validateField(key, formData[key]);
      if (error) newErrors[key] = error;
    });

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      setIsSubmitting(false);
      return;
    }

    putProfileUpdateAPI(
      formData,
      (resp) => {
        // For debugging purposes only.
        console.log("onRegisterSuccess: Starting...");
        console.log(resp);

        // Redirect the user to a new page.
        setForceURL("/settings");
      },
      (apiErr) => {
        console.log("onRegisterError: apiErr:", apiErr);
        setErrors(apiErr);
      },
      () => {
        console.log("onRegisterDone: Starting...");

      },
    );
    setHasSubmitted(true);
  };

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <div className="min-h-screen bg-purple-50">
      <Topbar currentPage="Settings" />

      <div className="p-8">
        <div className="max-w-2xl mx-auto">
          <nav className="mb-8">
            <button
              onClick={() => setForceURL("/settings")}
              className="group flex items-center text-sm text-gray-600 hover:text-purple-600"
            >
              <ChevronLeft className="w-4 h-4 mr-1 transition-transform group-hover:-translate-x-1" />
              Back to Settings
            </button>
          </nav>

          <h1 className="text-3xl font-bold text-purple-800 mb-8" style={{fontFamily: 'Comic Sans MS, cursive'}}>
            Email Settings
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

          <form onSubmit={handleSubmit} className="bg-white rounded-xl p-8 shadow-lg border-2 border-purple-200">
            <div className="space-y-6">
              {/* Name fields */}
              <div className="grid md:grid-cols-2 gap-6">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1" htmlFor="firstName">
                    First Name *
                  </label>
                  <input
                    type="text"
                    id="firstName"
                    name="firstName"
                    value={formData.firstName}
                    onChange={handleChange}
                    className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                      errors.firstName ? "border-red-500" : "border-gray-300"
                    }`}
                  />
                  {errors.firstName && (
                    <p className="mt-1 text-sm text-red-600">{errors.firstName}</p>
                  )}
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1" htmlFor="lastName">
                    Last Name *
                  </label>
                  <input
                    type="text"
                    id="lastName"
                    name="lastName"
                    value={formData.lastName}
                    onChange={handleChange}
                    className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                      errors.lastName ? "border-red-500" : "border-gray-300"
                    }`}
                  />
                  {errors.lastName && (
                    <p className="mt-1 text-sm text-red-600">{errors.lastName}</p>
                  )}
                </div>
              </div>

              {/* Email field */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1" htmlFor="email">
                  Email Address *
                </label>
                <input
                  type="email"
                  id="email"
                  name="email"
                  value={formData.email}
                  onChange={handleChange}
                  className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                    errors.email ? "border-red-500" : "border-gray-300"
                  }`}
                />
                {errors.email && (
                  <p className="mt-1 text-sm text-red-600">{errors.email}</p>
                )}
              </div>

              {/* Country Selection */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1" htmlFor="country">
                  Country *
                </label>
                <select
                  id="country"
                  name="country"
                  value={formData.country}
                  onChange={handleChange}
                  className={`w-full h-11 px-4 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent bg-white leading-tight ${
                    errors.country ? "border-red-500" : "border-gray-300"
                  }`}
                  style={{ paddingTop: '0px', paddingBottom: '0px' }}
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

              {/* Country Other */}
              {formData.country === "Other" && (
                <div>
                  <label
                    htmlFor="countryOther"
                    className="block text-sm font-medium text-gray-700 mb-1"
                  >
                    Specify Country *
                  </label>
                  <input
                    type="text"
                    id="countryOther"
                    name="countryOther"
                    value={formData.countryOther}
                    onChange={handleChange}
                    className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                      errors.countryOther ? "border-red-500" : "border-gray-300"
                    }`}
                  />
                  {errors.countryOther && (
                    <p className="mt-1 text-sm text-red-600">
                      {errors.countryOther}
                    </p>
                  )}
                </div>
              )}

              {/* Timezone Selection */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1" htmlFor="timezone">
                  Timezone *
                </label>
                <select
                  id="timezone"
                  name="timezone"
                  value={formData.timezone}
                  onChange={handleChange}
                  className="w-full h-11 px-4 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent bg-white leading-tight"
                  style={{ paddingTop: '0px', paddingBottom: '0px' }}
                >
                  <option value="America/New_York">Eastern Time (ET)</option>
                  <option value="America/Chicago">Central Time (CT)</option>
                  <option value="America/Denver">Mountain Time (MT)</option>
                  <option value="America/Los_Angeles">Pacific Time (PT)</option>
                  <option value="Asia/Tokyo">Japan Standard Time (JST)</option>
                </select>
              </div>

              {/* Checkboxes */}
              <div className="space-y-4">
                <div className="flex items-center">
                  <input
                    type="checkbox"
                    id="agreeTermsOfService"
                    name="agreeTermsOfService"
                    checked={formData.agreeTermsOfService}
                    onChange={handleChange}
                    className={`h-4 w-4 text-purple-600 focus:ring-purple-500 border-gray-300 rounded ${
                      errors.agreeTermsOfService ? "border-red-500" : ""
                    }`}
                  />
                  <label className="ml-2 block text-sm text-gray-700" htmlFor="agreeTermsOfService">
                    I agree to the{" "}
                    <a href="#" className="text-purple-600 hover:text-purple-500 underline">
                      Terms of Service
                    </a>{" "}
                    *
                  </label>
                </div>
                {errors.agreeTermsOfService && (
                  <p className="mt-1 text-sm text-red-600">{errors.agreeTermsOfService}</p>
                )}

                <div className="flex items-center">
                  <input
                    type="checkbox"
                    id="agreePromotional"
                    name="agreePromotional"
                    checked={formData.agreePromotional}
                    onChange={handleChange}
                    className="h-4 w-4 text-purple-600 focus:ring-purple-500 border-gray-300 rounded"
                  />
                  <label className="ml-2 block text-sm text-gray-700" htmlFor="agreePromotional">
                    I would like to receive promotional communications
                  </label>
                </div>
              </div>
            </div>

            {/* Submit and Cancel Buttons */}
            <div className="mt-8 flex gap-4">
              <button
                type="button"
                onClick={() => setForceURL("/settings")}
                className="w-full px-6 py-3 bg-white hover:bg-gray-50 text-gray-700 font-bold rounded-lg transition-colors border border-gray-300"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={isSubmitting}
                className={`w-full px-6 py-3 bg-purple-600 hover:bg-purple-700 text-white font-bold rounded-lg transition-colors ${
                  isSubmitting ? "opacity-50 cursor-not-allowed" : ""
                }`}
              >
                {isSubmitting ? "Saving Changes..." : "Save Changes"}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default EmailSettingsPage;
