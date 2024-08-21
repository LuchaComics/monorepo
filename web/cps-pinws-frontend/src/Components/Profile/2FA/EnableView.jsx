import React, { useState, useEffect } from "react";
import { Link, Navigate, useSearchParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faEye,
  faBuilding,
  faEnvelope,
  faSquarePhone,
  faTable,
  faHome,
  faLock,
  faTimesCircle,
  faArrowRight,
  faArrowUpRightFromSquare,
  faImage,
  faEllipsis,
  faRepeat,
  faTasks,
  faTachometer,
  faPlus,
  faArrowLeft,
  faCheckCircle,
  faUserCircle,
  faGauge,
  faPencil,
  faIdCard,
  faAddressBook,
  faContactCard,
  faChartPie,
  faKey,
  faUnlock,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import QRCode from "qrcode.react";

import { getProfileDetailAPI } from "../../../API/Profile";
import { postVertifyOTP } from "../../../API/Gateway";
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


function AccountEnableTwoFactorAuthentication() {
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

  const onVerificationTokenSubmitButtonClick = () => {
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
  };

  ////
  //// API.
  ////

  // --- Generate OTP --- //

  function onGenerateOPTSuccess(response) {
    console.log("onGenerateOPTSuccess: Starting...");
    console.log("response: ", response);
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
  //// BREADCRUMB
  ////

  const generateBreadcrumbItemLink = (currentUser) => {
    let dashboardLink;
    switch (currentUser.role) {
      case EXECUTIVE_ROLE_ID:
        dashboardLink = "/admin/dashboard";
        break;
      case MANAGEMENT_ROLE_ID:
        dashboardLink = "/admin/dashboard";
        break;
      case FRONTLINE_ROLE_ID:
        dashboardLink = "/admin/dashboard";
        break;
      case CUSTOMER_ROLE_ID:
        dashboardLink = "/dashboard";
        break;
      case USER_ROLE_CUSTOMER:
        dashboardLink = "/c/dashboard";
        break;
      default:
        dashboardLink = "/501"; // Default or error handling
        break;
    }
    return dashboardLink;
  };

  const breadcrumbItems = {
    items: [
      {
        text: "Dashboard",
        link: generateBreadcrumbItemLink(currentUser),
        isActive: false,
        icon: faGauge,
      },
      { text: "Account", link: "/account", icon: faUserCircle, isActive: true },
    ],
    mobileBackLinkItems: {
      link: generateBreadcrumbItemLink(currentUser),
      text: "Back to Dashboard",
      icon: faArrowLeft,
    },
  };

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

  return (
    <>
      <div className="container">
        <section className="section">
          {/* Desktop Breadcrumbs */}
          <nav
            className="breadcrumb has-background-light is-hidden-touch p-4"
            aria-label="breadcrumbs"
          >
            <ul>
              <li className="">
                <Link
                  to={generateBreadcrumbItemLink(currentUser)}
                  aria-current="page"
                >
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Dashboard
                </Link>
              </li>
              <li>
                <Link to={"/account/2fa"} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faUserCircle} />
                  &nbsp;Profile (2FA)
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
          <nav
            className="breadcrumb has-background-light is-hidden-desktop p-4"
            aria-label="breadcrumbs"
          >
            <ul>
              <li className="">
                <Link to="/account/2fa" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Detail (2FA)
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page Title */}
          <h1 className="title is-2">
            <FontAwesomeIcon className="fas" icon={faUserCircle} />
            &nbsp;Profile
          </h1>
          <h4 className="subtitle is-4">
            <FontAwesomeIcon className="fas" icon={faEye} />
            &nbsp;Detail (2FA)
          </h4>
          <hr />

          {/* Page */}
          <nav className="box">
            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                <FormErrorBox errors={errors} />
                {currentUser && (
                  <div className="container">
                    <div className="content">
                      <div className="columns">
                        {/*
                            ------
                            STEP 1
                            ------
                        */}
                        <div className="column">
                          <h1 className="title is-4">Step 1</h1>
                          <hr />
                          <p class="has-text-grey">
                            To begin, please download any of the following
                            applications for your mobile device.
                          </p>
                          {/* Google Authenticator */}
                          <div className="card">
                            <div className="card-content">
                              <div className="media">
                                <div className="media-content">
                                  <p className="title is-6">
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
                                  <p className="title is-6">
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

                        {/*
                            ------
                            STEP 2
                            ------
                        */}
                        <div className="column">
                          <h1 className="title is-4">Step 2</h1>
                          <hr />
                          <p class="has-text-grey">
                            With your 2FA application open, please scan the
                            following QR code with your device and click next
                            when ready.
                          </p>
                          <p>&nbsp;</p>
                          <div className="columns is-centered is-hidden-mobile">
                            <div class="column is-half">
                              <figure class="image">
                                {otpResponse && (
                                  <QRCode
                                    value={otpResponse.optAuthURL}
                                    size={100}
                                    className=""
                                  />
                                )}
                                <br />
                                <span className="is-centered">
                                  Scan with your app
                                </span>
                              </figure>
                            </div>
                          </div>

                          <div className="columns is-centered is-hidden-tablet">
                            <div class="column is-half">
                              <figure class="image">
                                {otpResponse && (
                                  <QRCode
                                    value={otpResponse.optAuthURL}
                                    size={200}
                                    className=""
                                  />
                                )}
                                <br />
                                <span className="is-centered">
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
                            Copy and paste the following values into your
                            device:
                          </p>
                          <FormTextareaField
                            label="Account Name"
                            name="accountName"
                            placeholder="-"
                            value={`${process.env.REACT_APP_WWW_DOMAIN}: ${currentUser.email}`}
                            errorText={null}
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
                        </div>
                        {/*
                            ------
                            STEP 3
                            ------
                        */}
                        <div className="column">
                          <h1 className="title is-4">Step 3</h1>
                          <hr />
                          <p class="has-text-grey">
                            Open the two-step verification app on your mobile
                            device to get your verification code.
                          </p>
                          <p>&nbsp;</p>

                          <FormInputField
                            label="Enter your Verification Token"
                            name="verificationToken"
                            placeholder="See your authenticator app"
                            value={verificationToken}
                            errorText={errors && errors.verificationToken}
                            helpText=""
                            onChange={(e) =>
                              setVerificationToken(e.target.value)
                            }
                            isRequired={true}
                            maxWidth="380px"
                          />
                          <button
                            className="button is-success is-fullwidth-mobile"
                            onClick={onVerificationTokenSubmitButtonClick}
                          >
                            <FontAwesomeIcon className="fas" icon={faLock} />
                            &nbsp;Submit Verification Token
                          </button>
                        </div>
                      </div>
                    </div>

                    <div className="columns pt-5">
                      <div className="column is-half">
                        <Link
                          className="button is-fullwidth-mobile"
                          to={`/account/2fa`}
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back to Detail (2FA)
                        </Link>
                      </div>
                      <div className="column is-half has-text-right"></div>
                    </div>
                  </div>
                )}
              </>
            )}
          </nav>
        </section>
      </div>
    </>
  );
}

export default AccountEnableTwoFactorAuthentication;
