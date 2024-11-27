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

import { getProfileDetailAPI } from "../../../API/Profile";
import { postDisableOTP } from "../../../API/Gateway";
import FormErrorBox from "../../Reusable/FormErrorBox";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
} from "../../../AppState";
import FormInputField from "../../Reusable/FormInputField";
import { COMMERCIAL_CUSTOMER_TYPE_OF_ID } from "../../../Constants/App";
import { USER_ROLE_ROOT, USER_ROLE_RETAILER, USER_ROLE_CUSTOMER } from "../../../Constants/App";


function AccountTwoFactorAuthenticationBackupCodeGenerate() {
  ////
  ////
  //// URL Parameters.
  ////

  const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140
  const backupCode = searchParams.get("v");

  ////
  //// Global state.
  ////

  const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

  ////
  //// Component states.
  ////

  // Page related states.
  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");

  ////
  //// Event handling.
  ////

  ////
  //// API.
  ////

  // --- Account Detail --- //

  const onAccountDetailSuccess = (response) => {
    console.log("onAccountDetailSuccess: Starting...");
    setCurrentUser(response);
  };

  const onAccountDetailError = (apiErr) => {
    console.log("onAccountDetailError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  };

  const onAccountDetailDone = () => {
    console.log("onAccountDetailDone: Starting...");
    setFetching(false);
  };

  // --- 2FA Disable --- //

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
      setFetching(true);
      setErrors({});
      getProfileDetailAPI(
        onAccountDetailSuccess,
        onAccountDetailError,
        onAccountDetailDone,
      );
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
            {/* Title + Options */}
            {currentUser && (
              <>
                <div className="columns">
                  <div className="column">
                    <p className="title is-4">
                      <FontAwesomeIcon className="fas" icon={faKey} />
                      &nbsp;2FA Backup Code
                    </p>
                  </div>
                  <div className="column has-text-right">
                    {/* Add buttons */}
                  </div>
                </div>

                {/* Notification */}
                <article class="message is-primary">
                  <div class="message-body">
                    <FontAwesomeIcon icon={faCheckCircle} />
                    &nbsp;Successfull 2FA Verification
                  </div>
                </article>

                {/* Content */}
                <div className="content">
                  <FormErrorBox errors={errors} />
                  <p class="has-text-grey">
                    You have successfully verified your 2FA code and now are
                    granted backup code which you can use in case you lose your
                    phone or experience data loss.
                  </p>
                  <p class="has-text-grey">
                    Do not share this code with anyone and store it in a safe
                    location, as it provides full access to your account.
                  </p>
                  <FormInputField
                    label="2FA Backup Code"
                    name="backupCode"
                    placeholder="See your authenticator app"
                    value={backupCode}
                    errorText={errors && errors.backupCode}
                    onChange={null}
                    isRequired={true}
                  />
                </div>
              </>
            )}
            <nav class="level">
              <div class="level-left">
                <div class="level-item">{/* Add button here */}</div>
              </div>
              <div class="level-right">
                <div class="level-item">
                  <Link
                    type="button"
                    class="button is-primary is-fullwidth-mobile"
                    to={`/account/2fa`}
                  >
                    <FontAwesomeIcon icon={faCheckCircle} />
                    &nbsp;I confirm go back to Detail (2FA)
                  </Link>
                </div>
              </div>
            </nav>
          </nav>
        </section>
      </div>
    </>
  );
}

export default AccountTwoFactorAuthenticationBackupCodeGenerate;
