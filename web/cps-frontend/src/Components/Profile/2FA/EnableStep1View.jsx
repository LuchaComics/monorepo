import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faEye,
  faLock,
  faArrowRight,
  faArrowUpRightFromSquare,
  faArrowLeft,
  faUserCircle,
  faGauge,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import QRCode from "qrcode.react";

import { postGenerateOTP } from "../../../API/Gateway";
import FormErrorBox from "../../Reusable/FormErrorBox";
import { currentUserState, currentOTPResponseState } from "../../../AppState";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { USER_ROLE_ROOT, USER_ROLE_RETAILER, USER_ROLE_CUSTOMER } from "../../../Constants/App";


function AccountEnableTwoFactorAuthenticationStep1() {
  ////
  //// Global state.
  ////

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
    console.log("onGenerateOPTSuccess: response: ", response);
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
                      <p className="subtitle is-5">Step 1 of 3</p>
                      <progress
                        class="progress is-success"
                        value="33"
                        max="100"
                      >
                        43%
                      </progress>
                    </nav>

                    {/* Content */}
                    <form className="">
                      <h1 className="title is-2 is-size-4-mobile has-text-centered">
                        Setup Two-Factor Authentication
                      </h1>
                      <FormErrorBox errors={errors} />
                      <div className="content">
                        <p class="has-text-grey">
                          To ensure your account stays secure, you need to sign
                          in using <i>two-factor Authentication (2FA)</i>. The
                          following wizard will help you get setup with 2FA.
                        </p>
                        <p class="has-text-grey">
                          To make initial 2FA setup easier, we encourage you to
                          login on a device BESIDES the mobile device with the
                          camera that you wish to use. We recommend the
                          following setup:
                          <ul>
                            <li>Login on a desktop device</li>
                            <li>
                              Use your mobile phone to scan the QR code and
                              complete setup
                            </li>
                          </ul>
                        </p>
                        <p class="has-text-grey">
                          To begin, please download any of the following
                          applications for your mobile device.
                        </p>
                        {/* Apple 2FA Authenticator */}
                        <div className="card">
                          <div className="card-content">
                            <div className="media">
                              <div className="media-content">
                                <p className="title is-5">
                                  <u>Apple 2FA</u>
                                </p>
                              </div>
                            </div>

                            <div className="content">
                              <p>
                                All iOS and Mac devices with a{" "}
                                <b>Safari Web Browser</b> come with build in a
                                2FA verification services. Sign in with your{" "}
                                <i>Apple ID</i> in Safari and you can take
                                advantage of this service.
                              </p>
                            </div>
                          </div>
                          {/*}
                        <footer className="card-footer">
                            <button className="card-footer-item button is-primary is-small">
                                Download&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} />
                            </button>
                        </footer>
                        */}
                        </div>

                        {/* Google Authenticator */}
                        <div className="card">
                          <div className="card-content">
                            <div className="media">
                              <div className="media-content">
                                <p className="title is-5">
                                  <u>Google Authenticator</u>
                                </p>
                              </div>
                            </div>

                            <div className="content">
                              <p>
                                This 2FA app is created by <b>Google, Inc.</b>
                              </p>
                              <b>Download for iOS:</b>&nbsp;
                              <Link
                                className=""
                                to="https://apps.apple.com/ca/app/google-authenticator/id388497605"
                                target="_blank"
                                rel="noreferrer"
                              >
                                Visit App Store&nbsp;
                                <FontAwesomeIcon
                                  className="fas"
                                  icon={faArrowUpRightFromSquare}
                                />
                              </Link>
                              <br />
                              <b>Download for Android:</b>&nbsp;
                              <Link
                                className=""
                                to="https://play.google.com/store/apps/details?id=com.google.android.apps.authenticator2&pli=1"
                                target="_blank"
                                rel="noreferrer"
                              >
                                Visit Google Play&nbsp;
                                <FontAwesomeIcon
                                  className="fas"
                                  icon={faArrowUpRightFromSquare}
                                />
                              </Link>
                              <br />
                            </div>
                          </div>
                          {/*}
                            <footer className="card-footer">
                                <button className="card-footer-item button is-primary is-small">
                                    Download&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} />
                                </button>
                            </footer>
                            */}
                        </div>

                        {/* Authenticator Chrome Extension */}
                        <div className="card">
                          <div className="card-content">
                            <div className="media">
                              <div className="media-content">
                                <p className="title is-5">
                                  <u>Authenticator</u>
                                </p>
                              </div>
                            </div>

                            <div className="content">
                              <p>
                                This 2FA app is created by{" "}
                                <b>authenticator.cc</b>
                              </p>
                              <b>Download for Chrome:</b>&nbsp;
                              <Link
                                className=""
                                to="https://chromewebstore.google.com/detail/authenticator/bhghoamapcdpbohphigoooaddinpkbai?pli=1"
                                target="_blank"
                                rel="noreferrer"
                              >
                                Visit Chrome web store&nbsp;
                                <FontAwesomeIcon
                                  className="fas"
                                  icon={faArrowUpRightFromSquare}
                                />
                              </Link>
                              <br />
                            </div>
                          </div>
                          {/*}
                        <footer className="card-footer">
                            <button className="card-footer-item button is-primary is-small">
                                Download&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} />
                            </button>
                        </footer>
                        */}
                        </div>
                      </div>
                    </form>

                    {/* Bottom Navigation */}
                    <br />
                    <nav class="level">
                      <div class="level-left">
                        <div class="level-item">
                          <Link
                            class="button is-link is-fullwidth-mobile"
                            to="/account/2fa"
                          >
                            <FontAwesomeIcon icon={faArrowLeft} />
                            &nbsp;Cancel
                          </Link>
                        </div>
                      </div>
                      <div class="level-right">
                        <div class="level-item">
                          <Link
                            type="button"
                            class="button is-primary is-fullwidth-mobile"
                            to="/account/2fa/setup/step-2"
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

export default AccountEnableTwoFactorAuthenticationStep1;
