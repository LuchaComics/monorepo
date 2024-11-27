import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faEye,
  faLock,
  faArrowRight,
  faArrowLeft,
  faUserCircle,
  faGauge,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import QRCode from "qrcode.react";

import { postGenerateOTP } from "../../../API/Gateway";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormTextareaField from "../../Reusable/FormTextareaField";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
  currentOTPResponseState,
} from "../../../AppState";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import FormInputField from "../../Reusable/FormInputField";
import { USER_ROLE_ROOT, USER_ROLE_RETAILER, USER_ROLE_CUSTOMER } from "../../../Constants/App";


function AccountEnableTwoFactorAuthenticationStep2() {
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

  ////
  //// Event handling.
  ////

  ////
  //// API.
  ////

  // --- Generate OTP --- //

  function onGenerateOPTSuccess(response) {
    console.log("onGenerateOPTSuccess: Starting...");
    console.log("onGenerateOPTSuccess: response:", response);
    setOtpResponse(response);
  }

  function onGenerateOPTError(apiErr) {
    console.log("onGenerateOPTError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onGenerateOPTDone() {
    console.log("onGenerateOPTDone: Starting...");
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
      // If no verification code exists then generate a new OTP code now.
      if (
        otpResponse === undefined ||
        otpResponse === null ||
        otpResponse === ""
      ) {
        setFetching(true);
        setErrors({});
        postGenerateOTP(
          onGenerateOPTSuccess,
          onGenerateOPTError,
          onGenerateOPTDone,
        );
      }
    }

    return () => {
      mounted = false;
    };
  }, [otpResponse, onGenerateOPTSuccess]);

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

  console.log("otpResponse.accountName:", otpResponse.accountName);

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
                <FormErrorBox errors={errors} />
                {currentUser && (
                  <>
                    {/* Progress Wizard */}
                    <nav className="box has-background-light">
                      <p className="subtitle is-5">Step 2 of 3</p>
                      <progress
                        class="progress is-success"
                        value="66"
                        max="100"
                      >
                        66%
                      </progress>
                    </nav>
                    {/* Content */}
                    <form>
                      <h1 className="title is-size-2 is-size-4-mobile  has-text-centered">
                        Setup Two-Factor Authentication
                      </h1>
                      <FormErrorBox errors={errors} />
                      <p class="has-text-grey">
                        With your 2FA application open, please scan the
                        following QR code with your device and click next when
                        ready.
                      </p>
                      <p>&nbsp;</p>
                      <div className="columns is-centered has-text-centered">
                        <div class="column is-half-tablet">
                          <figure class="image">
                            {otpResponse && (
                              <QRCode
                                value={otpResponse.optAuthURL}
                                size={250}
                                className=""
                              />
                            )}
                            <br />
                            <span className="has-text-centered">
                              Scan with your app
                            </span>
                          </figure>
                        </div>
                      </div>

                      <br />
                      <p className="has-text-grey title is-3 has-text-centered">
                        - OR -
                      </p>
                      <br />

                      <p class="has-text-grey">
                        Copy and paste the following values into your device:
                      </p>
                      <br />
                      <FormTextareaField
                        label="Account Name"
                        name="accountName"
                        placeholder="-"
                        value={otpResponse.accountName}
                        errorText={errors && errors.accountName}
                        onChange={null}
                        isRequired={true}
                        maxWidth="380px"
                        readonly={true}
                      />
                      <FormTextareaField
                        label="Your Key"
                        name="yourKey"
                        placeholder="-"
                        value={otpResponse && otpResponse.base32}
                        errorText={errors && errors.base32}
                        onChange={null}
                        isRequired={true}
                        maxWidth="380px"
                        readonly={true}
                      />
                      <FormInputField
                        label="Type of Key"
                        name="typeOfKey"
                        placeholder="-"
                        value={"Time based"}
                        errorText={null}
                        helpText=""
                        onChange={null}
                        isRequired={true}
                        maxWidth="380px"
                        readonly={true}
                      />
                    </form>
                    {/* Bottom Navigation */}
                    <br />
                    <nav class="level">
                      <div class="level-left">
                        <div class="level-item">
                          <Link
                            class="button is-link is-fullwidth-mobile"
                            to="/account/2fa/setup/step-1"
                          >
                            <FontAwesomeIcon icon={faArrowLeft} />
                            &nbsp;Back to Step 1
                          </Link>
                        </div>
                      </div>
                      <div class="level-right">
                        <div class="level-item">
                          <Link
                            type="button"
                            class="button is-primary is-fullwidth-mobile"
                            to="/account/2fa/setup/step-3"
                          >
                            Next&nbsp;
                            <FontAwesomeIcon icon={faArrowRight} />
                          </Link>
                        </div>
                      </div>
                    </nav>
                  </>
                )}
              </>
            )}
          </nav>
        </section>
      </div>
    </>
  );
}

export default AccountEnableTwoFactorAuthenticationStep2;
