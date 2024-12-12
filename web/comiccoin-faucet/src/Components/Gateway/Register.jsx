import React, { useState } from 'react';
import { Coins, BookOpen, Camera, Gift, Github, ArrowRight } from 'lucide-react';
import { AlertCircle, ArrowLeft, CheckCircle2 } from 'lucide-react';
import { Navigate } from "react-router-dom";

const IndexPage = () => {
    const [formData, setFormData] = useState({
   firstName: '',
   lastName: '',
   email: '',
   phone: '',
   password: '',
   passwordConfirm: '',
   agreeTos: false,
   agreePromotional: false
 });

 const [errors, setErrors] = useState({});
 const [touched, setTouched] = useState({});
 const [isSubmitting, setIsSubmitting] = useState(false);
 const [passwordStrength, setPasswordStrength] = useState(0);
 const [forceURL, setForceURL] = useState("");

 if (forceURL !== "") {
   return <Navigate to={forceURL} />;
 }

 const validateField = (name, value) => {
   switch (name) {
     case 'firstName':
       if (!value.trim()) return 'First name is required';
       if (value.length < 2) return 'First name must be at least 2 characters';
       if (value.length > 50) return 'First name must be less than 50 characters';
       if (!/^[a-zA-Z\s-']+$/.test(value)) return 'First name can only contain letters, spaces, hyphens, and apostrophes';
       return '';

     case 'lastName':
       if (!value.trim()) return 'Last name is required';
       if (value.length < 2) return 'Last name must be at least 2 characters';
       if (value.length > 50) return 'Last name must be less than 50 characters';
       if (!/^[a-zA-Z\s-']+$/.test(value)) return 'Last name can only contain letters, spaces, hyphens, and apostrophes';
       return '';

     case 'email':
       if (!value) return 'Email is required';
       if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) return 'Please enter a valid email address';
       if (value.length > 100) return 'Email must be less than 100 characters';
       return '';

     case 'phone':
       if (!value) return 'Phone number is required';
       const phoneRegex = /^\+?([0-9]{1,3})?[-. ]?\(?([0-9]{3})\)?[-. ]?([0-9]{3})[-. ]?([0-9]{4})$/;
       if (!phoneRegex.test(value)) return 'Please enter a valid phone number';
       return '';

     case 'password':
       const passwordErrors = [];
       if (!value) return 'Password is required';
       if (value.length < 8) passwordErrors.push('at least 8 characters');
       if (!/[A-Z]/.test(value)) passwordErrors.push('one uppercase letter');
       if (!/[a-z]/.test(value)) passwordErrors.push('one lowercase letter');
       if (!/[0-9]/.test(value)) passwordErrors.push('one number');
       if (!/[!@#$%^&*]/.test(value)) passwordErrors.push('one special character');
       return passwordErrors.length ? `Password must contain ${passwordErrors.join(', ')}` : '';

     case 'passwordConfirm':
       if (!value) return 'Please confirm your password';
       if (value !== formData.password) return 'Passwords do not match';
       return '';

     case 'agreeTos':
       if (!value) return 'You must agree to the Terms of Service';
       return '';

     default:
       return '';
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
   const newValue = type === 'checkbox' ? checked : value;

   setFormData(prev => ({
     ...prev,
     [name]: newValue
   }));

   setTouched(prev => ({
     ...prev,
     [name]: true
   }));

   if (name === 'password') {
     setPasswordStrength(calculatePasswordStrength(value));
   }

   const error = validateField(name, newValue);
   setErrors(prev => ({
     ...prev,
     [name]: error
   }));
 };

 const handleSubmit = async (e) => {
   e.preventDefault();
   setIsSubmitting(true);

   const newErrors = {};
   Object.keys(formData).forEach(field => {
     const error = validateField(field, formData[field]);
     if (error) newErrors[field] = error;
   });

   setErrors(newErrors);
   setTouched(Object.keys(formData).reduce((acc, field) => ({...acc, [field]: true}), {}));

   if (Object.keys(newErrors).length === 0) {
     try {
       await new Promise(resolve => setTimeout(resolve, 1000));
       console.log('Form submitted successfully:', formData);
     } catch (error) {
       console.error('Submission error:', error);
       setErrors(prev => ({
         ...prev,
         submit: 'Failed to submit form. Please try again.'
       }));
     }
   }

   setIsSubmitting(false);
 };

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

        <main className="flex-grow container mx-auto px-4 py-8 max-w-2xl">
          <h1 className="text-4xl font-bold mb-8 text-purple-800 text-center">
            Register for ComicCoin
          </h1>

          {/* Custom Error Alert */}
            {Object.keys(errors).length > 0 && (
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
                        {Object.values(errors).map((error, index) => (
                          <li key={index}>{error}</li>
                        ))}
                      </ul>
                    </div>
                  </div>
                </div>
              </div>
            )}

          <form onSubmit={handleSubmit} className="bg-white rounded-xl p-8 shadow-lg border-2 border-purple-200">
            <div className="grid md:grid-cols-2 gap-6">
              <div>
                <label htmlFor="firstName" className="block text-sm font-medium text-gray-700 mb-1">
                  First Name *
                </label>
                <input
                  type="text"
                  id="firstName"
                  name="firstName"
                  maxLength="50"
                  value={formData.firstName}
                  onChange={handleChange}
                  className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                    errors.firstName ? 'border-red-500' : 'border-gray-300'
                  }`}
                />
                {errors.firstName && (
                  <p className="mt-1 text-sm text-red-600">{errors.firstName}</p>
                )}
              </div>

              <div>
                <label htmlFor="lastName" className="block text-sm font-medium text-gray-700 mb-1">
                  Last Name *
                </label>
                <input
                  type="text"
                  id="lastName"
                  name="lastName"
                  maxLength="50"
                  value={formData.lastName}
                  onChange={handleChange}
                  className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                    errors.lastName ? 'border-red-500' : 'border-gray-300'
                  }`}
                />
                {errors.lastName && (
                  <p className="mt-1 text-sm text-red-600">{errors.lastName}</p>
                )}
              </div>

              <div className="md:col-span-2">
                <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-1">
                  Email Address *
                </label>
                <input
                  type="email"
                  id="email"
                  name="email"
                  maxLength="100"
                  value={formData.email}
                  onChange={handleChange}
                  className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                    errors.email ? 'border-red-500' : 'border-gray-300'
                  }`}
                />
                {errors.email && (
                  <p className="mt-1 text-sm text-red-600">{errors.email}</p>
                )}
              </div>

              <div className="md:col-span-2">
                <label htmlFor="phone" className="block text-sm font-medium text-gray-700 mb-1">
                  Phone Number *
                </label>
                <input
                  type="tel"
                  id="phone"
                  name="phone"
                  placeholder="+1 (555) 555-5555"
                  value={formData.phone}
                  onChange={handleChange}
                  className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                    errors.phone ? 'border-red-500' : 'border-gray-300'
                  }`}
                />
                {errors.phone && (
                  <p className="mt-1 text-sm text-red-600">{errors.phone}</p>
                )}
              </div>

              <div>
                <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-1">
                  Password *
                </label>
                <input
                  type="password"
                  id="password"
                  name="password"
                  value={formData.password}
                  onChange={handleChange}
                  className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                    errors.password ? 'border-red-500' : 'border-gray-300'
                  }`}
                />
                {passwordStrength > 0 && (
                  <div className="mt-2">
                    <div className="h-2 bg-gray-200 rounded-full">
                      <div
                        className={`h-full rounded-full transition-all ${
                          passwordStrength <= 40 ? 'bg-red-500' :
                          passwordStrength <= 80 ? 'bg-yellow-500' :
                          'bg-green-500'
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
                <label htmlFor="passwordConfirm" className="block text-sm font-medium text-gray-700 mb-1">
                  Confirm Password *
                </label>
                <input
                  type="password"
                  id="passwordConfirm"
                  name="passwordConfirm"
                  value={formData.passwordConfirm}
                  onChange={handleChange}
                  className={`w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent ${
                    errors.passwordConfirm ? 'border-red-500' : 'border-gray-300'
                  }`}
                />
                {errors.passwordConfirm && (
                  <p className="mt-1 text-sm text-red-600">{errors.passwordConfirm}</p>
                )}
              </div>

              <div className="md:col-span-2 space-y-4">
                <div className="flex items-center">
                  <input
                    type="checkbox"
                    id="agreeTos"
                    name="agreeTos"
                    checked={formData.agreeTos}
                    onChange={handleChange}
                    className="h-4 w-4 text-purple-600 focus:ring-purple-500 border-gray-300 rounded"
                  />
                  <label htmlFor="agreeTos" className="ml-2 block text-sm text-gray-700">
                    I agree to the <a href="#" className="text-purple-600 hover:text-purple-500 underline">Terms of Service</a> *
                  </label>
                </div>
                {errors.agreeTos && (
                  <p className="mt-1 text-sm text-red-600">{errors.agreeTos}</p>
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
                  <label htmlFor="agreePromotional" className="ml-2 block text-sm text-gray-700">
                    I would like to receive promotional emails
                  </label>
                </div>
              </div>
            </div>

            <div className="mt-8">
              <button
                type="submit"
                disabled={isSubmitting}
                className={`w-full px-6 py-3 bg-purple-600 hover:bg-purple-700 text-white font-bold rounded-lg transition-colors ${
                  isSubmitting ? 'opacity-50 cursor-not-allowed' : ''
                }`}
              >
                {isSubmitting ? 'Creating Account...' : 'Create Account'}
                </button>

               {/* New Login Link */}
               <div className="text-center">
                 <br/>
                 <p className="text-gray-600">
                   Already have an account?{' '}
                   <a
                     onClick={(e)=>setForceURL("/login")}
                     className="text-purple-600 hover:text-purple-700 font-medium underline"
                   >
                     Click here to login
                   </a>
                 </p>
               </div>
             </div>
           </form>
         </main>

      <footer className="bg-gradient-to-r from-purple-700 to-indigo-800 text-white py-8">
        <div className="max-w-7xl mx-auto px-4 text-center">
          <p className="mb-4">© 2024 ComicCoin Faucet. All rights reserved.</p>
          <p>
            <a href="#" className="underline hover:text-purple-200">Accessibility Statement</a>
            {' '} | {' '}
            <a href="#" className="underline hover:text-purple-200">Terms of Service</a>
            {' '} | {' '}
            <a href="#" className="underline hover:text-purple-200">Privacy Policy</a>
          </p>
        </div>
      </footer>
    </div>

  );
};

export default IndexPage;
