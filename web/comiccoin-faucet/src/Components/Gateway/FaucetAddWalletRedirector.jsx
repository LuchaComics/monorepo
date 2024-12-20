import React from "react";
import { Navigate } from "react-router-dom";
import { useRecoilState } from "recoil";
import { useLocation } from "react-router-dom";

import { currentUserState } from "../../AppState";

/**
  The purpose of this component is to intercept users whom have no wallet
  address set and then redirect them to the `/add-my-wallet-to-faucet` url.
 */
function FaucetAddWalletRedirector() {
  ////
  //// Global state.
  ////
  const [currentUser] = useRecoilState(currentUserState);

  ////
  //// Logic
  ////

  // Developers Note:
  // This needs to always be above any GUI rendering, or else react.js will error.
  const location = useLocation();

  //
  // Step 1:
  // If user is not authenticated nor authorized then skip this redirector.
  //

  if (currentUser === undefined || currentUser === null || currentUser === "") {
      console.log("FaucetAddWalletRedirector: No user, exiting now...");
      return null;
  }

  //
  // Step 2:
  // Get the current location and if we are at specific URL paths then we
  // will not render this component.
  //

  const ignorePathsArr = [
      "/",
      "/register",
      "/register/user",
      "/register/store",
      "/register-successful",
      "/added-my-wallet-to-faucet-successfully",
      "/index",
      "/login",
      "/login/2fa",
      "/login/2fa/step-1",
      "/login/2fa/step-2",
      "/login/2fa/step-3",
      "/login/2fa/step-3/backup-code",
      "/login/2fa/backup-code",
      "/login/2fa/backup-code-recovery",
      "/logout",
      "/verify",
      "/forgot-password",
      "/password-reset",
      "/terms",
      "/privacy",
      "/add-my-wallet-to-faucet",
  ];
  var arrayLength = ignorePathsArr.length;
  for (var i = 0; i < arrayLength; i++) {
    // console.log(location.pathname, "===", ignorePathsArr[i], " EQUALS ", location.pathname === ignorePathsArr[i]); // For debugging purposes only.
    if (location.pathname === ignorePathsArr[i]) {
      console.log("FaucetAddWalletRedirector: Ignoring path, exiting now...");
      return null;
    }
  }

  //
  // Step 3:
  // Check the role of the user and skip if system administrator.
  //

  if (currentUser.role === 1) {
      console.log("FaucetAddWalletRedirector: Detected admin role - done.");
      return null;
  }

  //
  // Step 4:
  // Check if wallet address set if not then redirect immediately!
  //

  if (currentUser.walletAddress === undefined || currentUser.walletAddress === null || currentUser.walletAddress === "") {
      console.log("FaucetAddWalletRedirector | Redirecting user immediately to: /add-my-wallet-to-faucet");
      return <Navigate to={`/add-my-wallet-to-faucet`} />;
  }

  console.log("FaucetAddWalletRedirector: Done.");
  return null;
}

export default FaucetAddWalletRedirector;
