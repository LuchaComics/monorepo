import React, { useState, useEffect } from "react";
import { Link, Navigate, useSearchParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faEye,
  faLock,
  faArrowLeft,
  faCheckCircle,
  faUserCircle,
  faGauge,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import QRCode from "qrcode.react";

import { postVertifyOTP } from "../../../API/Gateway";
import FormErrorBox from "../../Reusable/FormErrorBox";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
  currentOTPResponseState,
} from "../../../AppState";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import FormInputField from "../../Reusable/FormInputField";
import { USER_ROLE_ROOT, USER_ROLE_RETAILER, USER_ROLE_CUSTOMER } from "../../../Constants/App";


function AccountEnableTwoFactorAuthenticationStep3() {
  ////
  //// URL Parameters.
  ////

  const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140
  const paramToken = searchParams.get("token");

  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);
  const [otpResponse, setOtpResponse] = useRecoilState(currentOTPResponseState);

  ////
  //// Component states.
  ////

  // Page related states.
  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

  // Modal related states.
  const [verificationToken, setVerificationToken] = useState("");
  const [submittedParamToken, setSubmittedParamToken] = useState(false);

  ////
  //// Event handling.
  ////

  function onButtonClick(e) {
    // Remove whitespace characters from verificationToken
    const cleanedVerificationToken = verificationToken.replace(/\s/g, "");

    const payload = {
      verification_token: cleanedVerificationToken,
    };
    postVertifyOTP(
      payload,
      onVerifyOPTSuccess,
      onVerifyOPTError,
      onVerifyOPTDone,
    );
  }

  ////
  //// API.
  ////

  // --- Enable 2FA --- //

  function onVerifyOPTSuccess(response) {
    console.log("onVerifyOPTSuccess: Starting...");
    if (response !== undefined && response !== null && response !== "") {
      console.log("response: ", response);
      if (
        response.user !== undefined &&
        response.user !== null &&
        response.user !== ""
      ) {
        console.log("response.user: ", response.user);
        console.log("response.otp_backup_code: ", response.otp_backup_code);

        // Clear errors.
        setErrors({});

        // Save our updated user account.
        setCurrentUser(response.user);

        // Delete the OTP code.
        console.log("deleting otp response", otpResponse);
        setOtpResponse("");

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("2FA Enabled");
        setTopAlertStatus("success");
        setTimeout(() => {
          console.log("onSuccess: Delayed for 2 seconds.");
          console.log(
            "onSuccess: topAlertMessage, topAlertStatus:",
            topAlertMessage,
            topAlertStatus,
          );
          setTopAlertMessage("");
        }, 2000);

        // Change page.
        setForceURL("/account/2fa/backup-code?v=" + response.otp_backup_code);
      }
    }
  }

  function onVerifyOPTError(apiErr) {
    console.log("onVerifyOPTError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed Enabling 2FA");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onSuccess: Delayed for 2 seconds.");
      console.log(
        "onSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onVerifyOPTDone() {
    console.log("onVerifyOPTDone: Starting...");
    setFetching(false);
  }

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.

      // DEVELOPERS NOTE:
      // It appears that `Apple Verification` service submits a `token` url
      // parameter to the page with the uniquely generated 2FA code; as a result,
      // the following code will check to see if this `token` url parameter
      // exists and whether it was submitted or not and if it wasn't submitted
      // then we submit for OTP verification and proceed.
      if (
        submittedParamToken === false &&
        paramToken !== undefined &&
        paramToken !== null &&
        paramToken !== ""
      ) {
        setFetching(true);
        setErrors({});

        const payload = {
          verification_token: paramToken,
        };
        postVertifyOTP(
          payload,
          onVerifyOPTSuccess,
          onVerifyOPTError,
          onVerifyOPTDone,
        );
        setSubmittedParamToken(true);
        setVerificationToken(paramToken);
      }
    }

    return () => {
      mounted = false;
    };
  }, []);

  ////
  //// Component rendering.
  ////

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  // Generate URL's based on user role.
  let dashboardURL = "/501";
  if (currentUser) {
      if (currentUser.role === USER_ROLE_ROOT) {
        dashboardURL = "/admin/dashboard";
      }
      if (currentUser.role === USER_ROLE_RETAILER) {
        dashboardURL = "/dashboard";
      }
      if (currentUser.role === USER_ROLE_RETAILER) {
        dashboardURL = "/dashboard";
      }
      if (currentUser.role === USER_ROLE_CUSTOMER) {
        dashboardURL = "/c/dashboard";
      }
  }

  return (
    <>
      <div className="container">
        <section className="section">
          {/* Desktop Breadcrumbs */}
          <nav class="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={dashboardURL} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Dashboard
                </Link>
              </li>
              <li class="">
                <Link aria-current="page" to="/account/2fa">
                  <FontAwesomeIcon className="fas" icon={faUserCircle} />
                  &nbsp;Account&nbsp;(2FA)
                </Link>
              </li>
              <li className="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faLock} />
                  &nbsp;Enable 2FA
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/account/2fa`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Account (2FA)
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page */}
          <nav className="box">
            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                {/* Progress Wizard */}
                <nav className="box has-background-success-light">
                  <p className="subtitle is-5">Step 3 of 3</p>
                  <progress class="progress is-success" value="100" max="100">
                    100%
                  </progress>
                </nav>
                {/* Content */}
                <form>
                  <h1 className="title is-size-2 is-size-4-mobile  has-text-centered">
                    Setup Two-Factor Authentication
                  </h1>
                  <FormErrorBox errors={errors} />
                  <p class="has-text-grey">
                    Open the two-step verification app on your mobile device to
                    get your verification code.
                  </p>
                  <p>&nbsp;</p>
                  <FormInputField
                    label="Enter your Verification Token"
                    name="verificationToken"
                    placeholder="See your authenticator app"
                    value={verificationToken}
                    errorText={errors && errors.verificationToken}
                    helpText=""
                    onChange={(e) => setVerificationToken(e.target.value)}
                    isRequired={true}
                    maxWidth="380px"
                  />
                </form>
                {/* Bottom Navigation */}
                <br />
                <nav class="level">
                  <div class="level-left">
                    <div class="level-item">
                      <Link
                        class="button is-link is-fullwidth-mobile"
                        to="/account/2fa/setup/step-2"
                      >
                        <FontAwesomeIcon icon={faArrowLeft} />
                        &nbsp;Back to Step 2
                      </Link>
                    </div>
                  </div>
                  <div class="level-right">
                    <div class="level-item">
                      <button
                        type="button"
                        class="button is-primary is-fullwidth-mobile"
                        onClick={onButtonClick}
                      >
                        <FontAwesomeIcon icon={faCheckCircle} />
                        &nbsp;Submit and Verify
                      </button>
                    </div>
                  </div>
                </nav>
              </>
            )}
          </nav>
        </section>
      </div>
    </>
  );
}

export default AccountEnableTwoFactorAuthenticationStep3;
