import React, { useEffect, useState } from "react";
import { postLogoutAPI } from "../../API/Gateway";
import Scroll from "react-scroll";

import { Coins, Loader } from 'lucide-react'


function LogoutRedirector() {
  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});

  ////
  //// API.
  ////

  function onLogoutnSuccess(response) {
    console.log("onLogoutnSuccess: Starting...");
  }

  function onLogoutnError(apiErr) {
    console.log("onLogoutnError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onLogoutnDone() {
    console.log("onLogoutnDone: Starting...");
    function onRedirect(e) {
      // Clear the entire local storage.
      localStorage.clear();

      // Do not use `Link` but instead use the `window.location` change
      // to fix the issue with the `TopNavigation` component to restart.
      // If you use use `Link` then when you redirect to the navigation then
      // the menu will not update.
      window.location.href = "/login";
    }

    setTimeout(onRedirect, 250);
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
        postLogoutAPI(onLogoutnSuccess, onLogoutnError, onLogoutnDone);
    }

    return () => (mounted = false);
  }, []);

  return (
      <div className="min-h-screen flex flex-col bg-gradient-to-b from-purple-100 to-white">
        {/* Header */}
        <nav className="bg-gradient-to-r from-purple-700 to-indigo-800 text-white p-4">
          <div className="max-w-7xl mx-auto flex justify-between items-center">
            <div className="flex items-center space-x-2">
              <Coins className="h-8 w-8" />
              <span className="text-2xl font-bold">ComicCoin Faucet</span>
            </div>
          </div>
        </nav>

        {/* Main Content - Now properly centered */}
        <div className="flex-grow flex items-center justify-center">
          <div className="text-center">
            <div className="flex justify-center mb-6">
              <Loader className="w-12 h-12 text-purple-600 animate-spin" />
            </div>
            <h1 className="text-2xl font-bold text-purple-800 mb-4" style={{fontFamily: 'Comic Sans MS, cursive'}}>
              Logging Out...
            </h1>
            <p className="text-gray-600">
              Please wait while we securely log you out of your account
            </p>
          </div>
        </div>
      </div>
  );
}

export default LogoutRedirector;
