import React, { useState } from 'react';
import { ChevronLeft, AlertCircle, Eye, EyeOff } from 'lucide-react';
import { Navigate } from 'react-router-dom';

import AdminTopbar from "../../../Components/Navigation/AdminTopbar";
import { putProfileChangePasswordAPI } from "../../../API/Profile";

const AdminChangePasswordPage = () => {
  const [forceURL, setForceURL] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [hasSubmitted, setHasSubmitted] = useState(false);
  const [errors, setErrors] = useState({});
  const [showPasswords, setShowPasswords] = useState({
    old: false,
    new: false,
    confirm: false
  });
  const [passwordStrength, setPasswordStrength] = useState(0);

  const [formData, setFormData] = useState({
    old_password: "",
    new_password: "",
    new_password_repeated: ""
  });

  const validateField = (name, value) => {
    switch (name) {
      case "old_password":
        if (!value) return "Current password is required";
        return "";

      case "new_password":
        const passwordErrors = [];
        if (!value) return "New password is required";
        if (value.length < 8) passwordErrors.push("at least 8 characters");
        if (!/[A-Z]/.test(value)) passwordErrors.push("one uppercase letter");
        if (!/[a-z]/.test(value)) passwordErrors.push("one lowercase letter");
        if (!/[0-9]/.test(value)) passwordErrors.push("one number");
        if (!/[!@#$%^&*]/.test(value)) passwordErrors.push("one special character");
        if (value === formData.old_password) return "New password must be different from current password";
        return passwordErrors.length ? `Password must contain ${passwordErrors.join(", ")}` : "";

      case "new_password_repeated":
        if (!value) return "Please confirm your new password";
        if (value !== formData.new_password) return "Passwords do not match";
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
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));

    if (name === "new_password") {
      setPasswordStrength(calculatePasswordStrength(value));
    }

    if (hasSubmitted) {
      const error = validateField(name, value);
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

  // Call the API with our formData
  putProfileChangePasswordAPI(
    formData,
    (resp) => {
      // For debugging purposes only.
      console.log("putProfileChangePasswordAPI: Starting...");
      console.log(resp);

      // Redirect the user to a new page.
      setForceURL("/admin/settings");
    },
    (apiErr) => {
      console.log("putProfileChangePasswordAPI: apiErr:", apiErr);
      setErrors(apiErr);
      setIsSubmitting(false);
    },
    () => {
      console.log("putProfileChangePasswordAPI: Done...");
      setIsSubmitting(false);
    }
  );
};

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  const PasswordInput = ({ id, name, value, placeholder, error, show }) => (
    <div className="relative">
      <input
        type={show ? "text" : "password"}
        id={id}
        name={name}
        value={value}
        onChange={handleChange}
        placeholder={placeholder}
        className={`w-full h-11 px-4 py-2 pr-10 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
          error ? "border-red-500" : "border-gray-300"
        }`}
      />
      <button
        type="button"
        onClick={() => setShowPasswords(prev => ({ ...prev, [name]: !prev[name] }))}
        className="absolute right-3 top-3 text-gray-400 hover:text-gray-600"
      >
        {show ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
      </button>
    </div>
  );

  return (
    <div className="min-h-screen bg-purple-50">
      <AdminTopbar currentPage="Settings" />

      <div className="p-8">
        {/* Header with Back Button */}
        <div className="max-w-2xl mx-auto">
          <nav className="mb-8">
            <button
              onClick={() => setForceURL("/admin/settings")}
              className="group flex items-center text-sm text-gray-600 hover:text-purple-600"
            >
              <ChevronLeft className="w-4 h-4 mr-1 transition-transform group-hover:-translate-x-1" />
              Back to Settings
            </button>
          </nav>

          <h1 className="text-3xl font-bold text-purple-800 mb-8" style={{fontFamily: 'Comic Sans MS, cursive'}}>
            Change Password
          </h1>
        </div>

        {/* Main Form */}
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
              {/* Current Password */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1" htmlFor="old_password">
                  Current Password *
                </label>
                <PasswordInput
                  id="old_password"
                  name="old_password"
                  value={formData.old_password}
                  placeholder="Enter your current password"
                  error={errors.old_password}
                  show={showPasswords.old}
                />
                {errors.old_password && (
                  <p className="mt-1 text-sm text-red-600">{errors.old_password}</p>
                )}
              </div>

              {/* New Password */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1" htmlFor="new_password">
                  New Password *
                </label>
                <PasswordInput
                  id="new_password"
                  name="new_password"
                  value={formData.new_password}
                  placeholder="Enter your new password"
                  error={errors.new_password}
                  show={showPasswords.new}
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
                {errors.new_password && (
                  <p className="mt-1 text-sm text-red-600">{errors.new_password}</p>
                )}
              </div>

              {/* Confirm New Password */}
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1" htmlFor="new_password_repeated">
                  Confirm New Password *
                </label>
                <PasswordInput
                  id="new_password_repeated"
                  name="new_password_repeated"
                  value={formData.new_password_repeated}
                  placeholder="Confirm your new password"
                  error={errors.new_password_repeated}
                  show={showPasswords.confirm}
                />
                {errors.new_password_repeated && (
                  <p className="mt-1 text-sm text-red-600">{errors.new_password_repeated}</p>
                )}
              </div>
            </div>

            {/* Submit and Cancel Buttons */}
            <div className="mt-8 flex gap-4">
              <button
                type="button"
                onClick={() => setForceURL("/admin/settings")}
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
                {isSubmitting ? "Changing Password..." : "Change Password"}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default AdminChangePasswordPage;
