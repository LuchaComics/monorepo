import React, { useState } from 'react';
import { Coins, AlertCircle, ArrowLeft } from 'lucide-react';
import { Navigate, Link } from "react-router-dom";
import { useRecoilState } from "recoil";

import { postLoginAPI } from "../../API/Gateway";
import { currentUserState } from "../../AppState";

export default function LoginPage() {
  // Variable controls the global state of the app.
  const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

  // Variable controls the login form.
  const [formData, setFormData] = useState({
    email: "",
    password: "",
    rememberMe: false,
  });

  const [errors, setErrors] = useState({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [hasSubmitted, setHasSubmitted] = useState(false);

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === "checkbox" ? checked : value
    }));
  };

  const handleSubmit = (e) => {
    e.preventDefault();
    setIsSubmitting(true);

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

    postLoginAPI(
      formData,
      (resp) => {
        // SUCCESS HANDLER

        // For debugging purposes only.
        console.log("onLoginSuccess: Starting...");
        console.log(resp);

        // Store in persistance storage in the browser.
        setCurrentUser(resp.user);

        if (resp.user.otpEnabled === false) {
          if (resp.user.role === 1) { // Administrator
              console.log("onLoginSuccess | redirecting to dashboard");
              setForceURL("/admin/dashboard");
          } else { // Non-administrator
              console.log("onLoginSuccess | redirecting to dashboard");
              setForceURL("/dashboard");
          }

        } else {
          if (resp.user.otpVerified === false) {
            console.log("onLoginSuccess | redirecting to 2fa setup wizard");
            setForceURL("/login/2fa/step-1");
          } else {
            console.log("onLoginSuccess | redirecting to 2fa validation");
            setForceURL("/login/2fa");
          }
        }
        // end SUCCESS HANDLER
      },
      (apiErr) => {
        console.log("onLoginError: apiErr:", apiErr);
        setErrors(apiErr);
      },
      () => {
        console.log("onLoginDone: Starting...");
        setIsSubmitting(false);
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
            Login to ComicCoin Faucet
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
            className="bg-white rounded-xl p-8 shadow-lg border-2 border-purple-200 space-y-6"
          >
            {/* Email field */}
            <div className="space-y-1">
              <label
                htmlFor="email"
                className="block text-sm font-medium text-gray-700"
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
                  errors.email ? "border-red-500 bg-red-50" : "border-gray-300"
                }`}
                placeholder="Enter your email"
              />
              {errors.email && (
                <p className="mt-1 text-sm text-red-600">{errors.email}</p>
              )}
            </div>

            {/* Password field */}
            <div className="space-y-1">
              <label
                htmlFor="password"
                className="block text-sm font-medium text-gray-700"
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
                  errors.password ? "border-red-500 bg-red-50" : "border-gray-300"
                }`}
                placeholder="Enter your password"
              />
              {errors.password && (
                <p className="mt-1 text-sm text-red-600">{errors.password}</p>
              )}
            </div>

            {/* Remember Me checkbox and Forgot Password link */}
            <div className="flex items-center justify-between">
              <div className="flex items-center">
                <input
                  type="checkbox"
                  id="rememberMe"
                  name="rememberMe"
                  checked={formData.rememberMe}
                  onChange={handleChange}
                  className="h-4 w-4 text-purple-600 focus:ring-purple-500 border-gray-300 rounded"
                />
                <label
                  htmlFor="rememberMe"
                  className="ml-2 block text-sm text-gray-700"
                >
                  Remember me
                </label>
              </div>

              <button
                onClick={(e)=>setForceURL("/forgot-password")}
                type="button"
                className="text-sm text-purple-600 hover:text-purple-700 font-medium"
              >
                Forgot password?
              </button>
            </div>

            {/* Submit button */}
            <div className="space-y-4">
              <button
                type="submit"
                disabled={isSubmitting}
                className={`w-full px-6 py-3 bg-purple-600 hover:bg-purple-700 text-white font-bold rounded-lg transition-colors ${
                  isSubmitting ? "opacity-50 cursor-not-allowed" : ""
                }`}
              >
                {isSubmitting ? "Signing in..." : "Sign in"}
              </button>

              <div className="text-center">
                <p className="text-gray-600">
                  Don't have an account?{" "}
                  <button
                    onClick={(e)=>setForceURL("/register")}
                    type="button"
                    className="text-purple-600 hover:text-purple-700 font-medium"
                  >
                    Create an account
                  </button>
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
            <Link to="/terms" className="underline hover:text-purple-200">Terms of Service</Link>
            {' | '}
            <Link to="/privacy" className="underline hover:text-purple-200">Privacy Policy</Link>
          </p>
        </div>
      </footer>
    </div>
  );
}
