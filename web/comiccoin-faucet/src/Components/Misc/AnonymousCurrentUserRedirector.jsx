import React from "react";
import { Navigate } from "react-router-dom";
import { useRecoilState } from "recoil";
import { useLocation } from "react-router-dom";

import { currentUserState } from "../../AppState";

/**
  The purpose of this component is to intercept anonymous users at our
  application URLs which require authorization.
 */
function AnonymousCurrentUserRedirector() {
  ////
  //// Global state.
  ////
  const [currentUser] = useRecoilState(currentUserState);

  ////
  //// Logic
  ////

  // Get the current location and if we are at specific URL paths then we
  // will not render this component.
  const ignorePathsArr = [
      "/",
      "/register",
      "/register/user",
      "/register/store",
      "/register-successful",
      "/add-my-wallet-to-faucet",
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
      "/cpsrn-result",
      "/cpsrn-registry",
      "/terms",
      "/privacy",
      "/cpsrn"
  ];
  const location = useLocation();
  var arrayLength = ignorePathsArr.length;
  for (var i = 0; i < arrayLength; i++) {
    // console.log(location.pathname, "===", ignorePathsArr[i], " EQUALS ", location.pathname === ignorePathsArr[i]); // For debugging purposes only.
    if (location.pathname === ignorePathsArr[i]) {
      return null;
    }
  }

  console.log("AnonymousCurrentUserRedirector | currentUser:", currentUser);

  if (currentUser === null) {
    console.log("No current user detected, redirecting back to login page.");
    return <Navigate to={`/login?unauthorized=true`} />;
  } else {
    return null;
  }
}

export default AnonymousCurrentUserRedirector;
