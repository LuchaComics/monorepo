import React, { useState, useEffect } from "react";
import { Link, useLocation } from "react-router-dom";
import { WalletMinimal, Send, QrCode, MoreHorizontal, Clock, Coins, Wallet, ArrowRight, ArrowUpRight, ArrowDownLeft, Ticket } from 'lucide-react';


function BottomTabBar() {
  ////
  //// Global State
  ////

  // Do nothing.

  ////
  //// Local State
  ////

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
  const ignorePathsArr = [
    "/",
    "/pick-data-directory",
    "/startup",
    "/create-your-first-wallet",
    "/wallets",
    "/wallet/add",
    "/send-processing",
    "/send-success",
  ];
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
      <nav className="fixed bottom-0 w-full bg-white border-t border-gray-200 shadow-lg" aria-label="Primary navigation">
        <div className="grid grid-cols-4 h-20">
          <Link
            to="/dashboard"
            className={`flex flex-col items-center justify-center space-y-2 ${location.pathname.includes("dashboard") && "bg-purple-50"}`}
            aria-label="Overview tab, currently selected"
            aria-current="page"
          >
            <Wallet className="w-7 h-7 text-purple-600" aria-hidden="true" />
            <span className="text-sm text-purple-600">Overview</span>
          </Link>
          <Link
            to="/send"
            className={`flex flex-col items-center justify-center space-y-2 ${location.pathname.includes("send") && "bg-purple-50"}`}
            aria-label="Send tab"
          >
            <Send className="w-7 h-7 text-gray-600" aria-hidden="true" />
            <span className="text-sm text-gray-600">Send</span>
          </Link>
          <Link
            to="/receive"
            className={`flex flex-col items-center justify-center space-y-2 ${location.pathname.includes("receive") && "bg-purple-50"}`}
            aria-label="Receive tab"
          >
            <QrCode className="w-7 h-7 text-gray-600" aria-hidden="true" />
            <span className="text-sm text-gray-600">Receive</span>
          </Link>
          <Link
            to="/more"
            className={`flex flex-col items-center justify-center space-y-2 ${!location.pathname.includes("dashboard") && !location.pathname.includes("send") && !location.pathname.includes("receive") && "bg-purple-50"}`}
            aria-label="More options tab"
          >
            <MoreHorizontal className="w-7 h-7 text-gray-600" aria-hidden="true" />
            <span className="text-sm text-gray-600">More</span>
          </Link>
        </div>
      </nav>
  );
}

export default BottomTabBar;
