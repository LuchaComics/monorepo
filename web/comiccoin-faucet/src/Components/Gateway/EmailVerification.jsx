import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import { useSearchParams } from "react-router-dom";
import Scroll from "react-scroll";
import { Coins, Mail, ArrowRight, ArrowLeft, CheckCircle } from 'lucide-react';

import { postEmailVerificationAPI } from "../../API/Gateway";
function EmailVerification() {
  ////
  //// URL Parameters
  ////

  let [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [emailVerification, setEmailVerification] = useState("");
  const [forceURL, setForceURL] = useState("");

  ////
  //// API.
  ////

  function onEmailVerificationSuccess(response) {
      console.log("onEmailVerificationSuccess: Starting...");
      console.log("onEmailVerificationSuccess: response:", response);
      setEmailVerification(response);
  }

  function onEmailVerificationError(apiErr) {
    console.log("onEmailVerificationError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onEmailVerificationDone() {
    console.log("onEmailVerificationDone: Starting...");
  }

  ////
  //// Event handling.
  ////

  // (Do nothing)

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      postEmailVerificationAPI(
        searchParams.get("q"), // Extract the verification code from the query parameter in the URL.
        onEmailVerificationSuccess,
        onEmailVerificationError,
        onEmailVerificationDone,
      );
    }

    return () => (mounted = false);
  }, []);

  ////
  //// Component rendering.
  ////

  return (
      <div className="min-h-screen flex flex-col bg-gradient-to-b from-purple-100 to-white">
     {/* Navigation */}
     <nav className="bg-gradient-to-r from-purple-700 to-indigo-800 text-white p-4">
       <div className="max-w-7xl mx-auto flex justify-between items-center">
         <div className="flex items-center space-x-2">
           <Coins className="h-8 w-8" />
           <span className="text-2xl font-bold">ComicCoin Faucet</span>
         </div>
         <button
           onClick={() => setForceURL("/")}
           className="flex items-center space-x-2 px-4 py-2 rounded-lg hover:bg-purple-600 transition-colors"
         >
           <ArrowLeft className="h-5 w-5" />
           <span>Back to Home</span>
         </button>
       </div>
     </nav>

     {/* Main Content */}
     <main className="flex-grow flex items-center justify-center">
       <div className="w-full max-w-2xl mx-4">
         <div className="bg-white rounded-xl p-8 shadow-lg border-2 border-purple-200">
           <div className="flex flex-col items-center space-y-6">
             {/* Success Icon */}
             <div className="text-green-500">
               <CheckCircle className="h-16 w-16" />
             </div>

             {/* Header */}
             <div className="flex flex-col items-center space-y-2">
               <div className="flex items-center space-x-2 text-purple-600">
                 <Mail className="h-8 w-8" />
                 <h1 className="text-2xl font-bold">Email Verified Successfully!</h1>
               </div>
               <p className="text-gray-500">Your email address has been confirmed</p>
             </div>

             {/* Message */}
             <div className="text-center space-y-4">
               <p className="text-gray-700">
                 Thank you for verifying your email address. You can now access all features of ComicCoin Faucet.
               </p>
               <p className="text-gray-600">
                 Start exploring the platform and earn ComicCoins by submitting your comic books!
               </p>
             </div>

             {/* Action Buttons */}
             <div className="flex flex-col sm:flex-row items-center space-y-4 sm:space-y-0 sm:space-x-4 pt-4">
               <Link
                 to={"/login"}
                 className="w-full sm:w-auto px-6 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors"
               >
                 Login
               </Link>
               <Link
                 to={"/"}
                 className="w-full sm:w-auto px-6 py-2 border border-purple-600 text-purple-600 rounded-lg hover:bg-purple-50 transition-colors flex items-center justify-center space-x-2"
               >
                 <span>Back to Home</span>
                 <ArrowRight className="h-4 w-4" />
               </Link>
             </div>
           </div>
         </div>
       </div>
     </main>

     {/* Footer */}
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

export default EmailVerification;
