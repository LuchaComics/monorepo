import React, { useState, useEffect } from "react";
import { Link, useLocation } from "react-router-dom";
import { WalletMinimal } from 'lucide-react';

function Topbar() {
  ////
  //// Global State
  ////

  // Do nothing.

  ////
  //// Local State
  ////

  // Do nothing.

  ////
  //// Events
  ////

  // Do nothing.

  ////
  //// Rendering.
  ////

  //-------------//
  // CASE 1 OF 3 //
  //-------------//

  // Get the current location and if we are at specific URL paths then we
  // will not render this component.
  const ignorePathsArr = ["/",];
  const location = useLocation();
  var arrayLength = ignorePathsArr.length;
  for (var i = 0; i < arrayLength; i++) {
    // console.log(location.pathname, "===", ignorePathsArr[i], " EQUALS ", location.pathname === ignorePathsArr[i]);
    if (location.pathname === ignorePathsArr[i]) {
      return null;
    }
  }

  // //-------------//
  // // CASE 2 OF 3 //
  // //-------------//
  //
  // if (currentUser === null) {
  //   return null;
  // }

  //-------------//
  // CASE 3 OF 3 //
  //-------------//

  // Render the following component GUI
  return (
      <header className="bg-gradient-to-r from-purple-700 to-indigo-800 px-8 py-6 text-white shadow-lg">
        <div className="flex items-center justify-center gap-4">
          <WalletMinimal className="w-8 h-8" aria-hidden="true" />
          <h1 className="text-3xl font-bold">ComicCoin Wallet</h1>
        </div>
      </header>
  );
}

export default Topbar;
