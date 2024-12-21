import React, { useState } from "react";
import { Coins, AlertCircle, ArrowLeft } from "lucide-react";
import { Navigate, Link } from "react-router-dom";

import { postRegisterAPI } from "../../API/Gateway";
import FormTimezoneSelectField from "../Reusable/FormTimezoneSelectField";

const RegisterPage = () => {
  const [timezone, setTimezone] = useState(Intl.DateTimeFormat().resolvedOptions().timeZone);

  const [formData, setFormData] = useState({
    firstName: "",
    lastName: "",
    email: "",
    country: "",
    countryOther: "",
    timezone: timezone,
    password: "",
    passwordConfirm: "",
    agreeTermsOfService: false,
    agreePromotional: false,
  });

  const [errors, setErrors] = useState({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [passwordStrength, setPasswordStrength] = useState(0);
  const [hasSubmitted, setHasSubmitted] = useState(false);
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");

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

      case "countryOther":
        if (formData.country === "Other" && !value.trim()) {
          return "Please specify your country";
        }
        return "";

      case "password":
        const passwordErrors = [];
        if (!value) return "Password is required";
        if (value.length < 8) passwordErrors.push("at least 8 characters");
        if (!/[A-Z]/.test(value)) passwordErrors.push("one uppercase letter");
        if (!/[a-z]/.test(value)) passwordErrors.push("one lowercase letter");
        if (!/[0-9]/.test(value)) passwordErrors.push("one number");
        if (!/[!@#$%^&*]/.test(value))
          passwordErrors.push("one special character");
        return passwordErrors.length
          ? `Password must contain ${passwordErrors.join(", ")}`
          : "";

      case "passwordConfirm":
        if (!value) return "Please confirm your password";
        if (value !== formData.password) return "Passwords do not match";
        return "";

      case "agreeTermsOfService":
        if (!value) return "You must agree to the Terms of Service";
        return "";

      default:
        return "";
    }
  };

  const calculatePasswordStrength = (password) => {
    let strength = 0;
    if (password.length >= 8) strength++;
    if (/[A-Z]/.test(password)) strength++;
    if (/[a-z]/.test(password)) strength++;
    if (/[0-9]/.test(password)) strength++;
    if (/[!@#$%^&*]/.test(password)) strength++;
    return (strength / 5) * 100;
  };

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    const newValue = type === "checkbox" ? checked : value;

    setFormData((prev) => ({
      ...prev,
      [name]: newValue,
      ...(name === "country" && value !== "Other" && { countryOther: "" }),
    }));

    if (name === "password") {
      setPasswordStrength(calculatePasswordStrength(value));
    }

    if (hasSubmitted) {
      const error = validateField(name, newValue);
      setErrors((prev) => ({
        ...prev,
        [name]: error,
      }));
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setFetching(true);
    const submission = {
      Email: formData.email,
      FirstName: formData.firstName,
      LastName: formData.lastName,
      Password: formData.password,
      PasswordConfirm: formData.passwordConfirm,
      Country: formData.country,
      AgreeTermsOfService: formData.agreeTermsOfService,
      AgreePromotional: formData.agreePromotional,
    };

    postRegisterAPI(
      formData,
      (resp) => {
        // For debugging purposes only.
        console.log("onRegisterSuccess: Starting...");
        console.log(resp);

        // Redirect the user to a new page.
        setForceURL("/register-successful");
      },
      (apiErr) => {
        console.log("onRegisterError: apiErr:", apiErr);
        setErrors(apiErr);
      },
      () => {
        console.log("onRegisterDone: Starting...");
        setFetching(false);
      },
    );
    setHasSubmitted(true);
  };

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <div className="min-h-screen flex flex-col bg-gradient-to-b from-purple-100 to-white">
      <nav className="bg-gradient-to-r from-purple-700 to-indigo-800 text-white p-4">
        <div className="max-w-7xl mx-auto flex justify-between items-center">
          <div className="flex items-center space-x-2">
            <Coins className="h-8 w-8" />
            <span className="text-2xl font-bold">ComicCoin Faucet</span>
          </div>
          <button onClick={(e)=>setForceURL("/")} className="flex items-center space-x-2 px-4 py-2 rounded-lg hover:bg-purple-600 transition-colors">
            <ArrowLeft className="h-5 w-5" />
            <span>Back to Home</span>
          </button>
        </div>
      </nav>

      <main className="flex-grow flex items-center justify-center">
        <div className="w-full max-w-md mx-4">
          <h1 className="text-4xl font-bold mb-8 text-purple-800 text-center">
          Register for ComicCoin
        </h1>

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

        <form
          onSubmit={handleSubmit}
          className="bg-white rounded-xl p-8 shadow-lg border-2 border-purple-200"
        >
          <div className="space-y-6">
            {/* Name fields */}
            <div className="grid md:grid-cols-2 gap-6">
              <div>
                <label
                  htmlFor="firstName"
                  className="block text-sm font-medium text-gray-700 mb-1"
                >
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
                <label
                  htmlFor="lastName"
                  className="block text-sm font-medium text-gray-700 mb-1"
                >
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
              <label
                htmlFor="email"
                className="block text-sm font-medium text-gray-700 mb-1"
              >
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

            {/* Country selection */}
            <div>
              <label
                htmlFor="country"
                className="block text-sm font-medium text-gray-700 mb-1"
              >
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

            {/* Other country field */}
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

            <FormTimezoneSelectField
                label="Timezone"
                name="timezone"
                selectedTimezone={timezone}
                setSelectedTimezone={setTimezone}
                errorText={errors.timezone}
                isRequired={true}
            />

            {/* Password fields */}
            <div className="grid md:grid-cols-2 gap-6">
              <div>
                <label
                  htmlFor="password"
                  className="block text-sm font-medium text-gray-700 mb-1"
                >
                  Password *
                </label>
                <input
                  type="password"
                  id="password"
                  name="password"
                  value={formData.password}
                  onChange={handleChange}
                  className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                    errors.password ? "border-red-500" : "border-gray-300"
                  }`}
                />
                {passwordStrength > 0 && (
                  <div className="mt-2">
                    <div className="h-2 bg-gray-200 rounded-full">
                      <div
                        className={`h-full rounded-full transition-all ${
                          passwordStrength <= 40
                            ? "bg-red-500"
                            : passwordStrength <= 80
                              ? "bg-yellow-500"
                              : "bg-green-500"
                        }`}
                        style={{ width: `${passwordStrength}%` }}
                      />
                    </div>
                  </div>
                )}
                {errors.password && (
                  <p className="mt-1 text-sm text-red-600">{errors.password}</p>
                )}
              </div>

              <div>
                <label
                  htmlFor="passwordConfirm"
                  className="block text-sm font-medium text-gray-700 mb-1"
                >
                  Confirm Password *
                </label>
                <input
                  type="password"
                  id="passwordConfirm"
                  name="passwordConfirm"
                  value={formData.passwordConfirm}
                  onChange={handleChange}
                  className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                    errors.passwordConfirm ? "border-red-500" : "border-gray-300"
                  }`}
                />
                {errors.passwordConfirm && (
                  <p className="mt-1 text-sm text-red-600">
                    {errors.passwordConfirm}
                  </p>
                )}
              </div>
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
                <label
                  htmlFor="agreeTermsOfService"
                  className="ml-2 block text-sm text-gray-700"
                >
                  I agree to the{" "}
                  <a
                    href="#"
                    className="text-purple-600 hover:text-purple-500 underline"
                  >
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
                <label
                  htmlFor="agreePromotional"
                  className="ml-2 block text-sm text-gray-700"
                >
                  I would like to receive promotional communications
                </label>
              </div>
            </div>
          </div>

          {/* Submit button and login link */}
          <div className="mt-8 space-y-4">
            <button
              type="submit"
              disabled={isSubmitting}
              className={`w-full px-6 py-3 bg-purple-600 hover:bg-purple-700 text-white font-bold rounded-lg transition-colors ${
                isSubmitting ? "opacity-50 cursor-not-allowed" : ""
              }`}
            >
              {isSubmitting ? "Creating Account..." : "Create Account"}
            </button>

            <div className="text-center">
              <p className="text-gray-600">
                Already have an account?{" "}
                <a
                  href="/login"
                  className="text-purple-600 hover:text-purple-700 font-medium underline"
                >
                  Click here to login
                </a>
              </p>
            </div>
          </div>
        </form>
        </div>
      </main>

      <footer className="bg-gradient-to-r from-purple-700 to-indigo-800 text-white py-8">
        <div className="max-w-7xl mx-auto px-4 text-center">
          <p className="mb-4">© 2024 ComicCoin Faucet. All rights reserved.</p>
          <p>
            <Link to="/terms" className="underline hover:text-purple-200">
              Terms of Service
            </Link>{" "}
            |{" "}
            <Link to="/privacy" className="underline hover:text-purple-200">
              Privacy Policy
            </Link>
          </p>
        </div>
      </footer>
    </div>
  );
};

export default RegisterPage;
