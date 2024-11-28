import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faChevronRight,
  faSearch,
  faTable,
  faPlus,
  faArrowLeft,
  faUserCircle,
  faGauge,
  faEye,
  faEllipsis,
  faKey,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { getProfileDetailAPI } from "../../../API/Profile";
import FormErrorBox from "../../Reusable/FormErrorBox";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { currentUserState } from "../../../AppState";
import BubbleLink from "../../Reusable/EveryPage/BubbleLink";
import AlertBanner from "../../Reusable/EveryPage/AlertBanner";
import { USER_ROLE_ROOT, USER_ROLE_RETAILER, USER_ROLE_CUSTOMER } from "../../../Constants/App";


function AccountMoreLaunchpad() {
  ////
  //// Global state.
  ////

  const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [client, setClient] = useState({});

  ////
  //// Event handling.
  ////

  // Do nothing...

  ////
  //// API.
  ////

  function onAccountDetailSuccess(response) {
    console.log("onAccountDetailSuccess: Starting...");
    setCurrentUser(response);
  }

  function onAccountDetailError(apiErr) {
    console.log("onAccountDetailError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onAccountDetailDone() {
    console.log("onAccountDetailDone: Starting...");
    setFetching(false);
  }

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
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faUserCircle} />
                  &nbsp;Account&nbsp;(More)
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={dashboardURL} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Dashboard
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page banner */}
          {client && client.status === 2 && (
            <AlertBanner message="Archived" status="info" />
          )}

          {/* Page */}
          <nav className="box">
            {/* Title + Options */}
            {client && (
              <div className="columns">
                <div className="column">
                  <p className="title is-4">
                    <FontAwesomeIcon className="fas" icon={faUserCircle} />
                    &nbsp;Account
                  </p>
                </div>
                <div className="column has-text-right"></div>
              </div>
            )}

            {/* <p className="pb-4">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                <FormErrorBox errors={errors} />

                {client && (
                  <div className="container">
                    {/* Tab Navigation */}
                    <div className="tabs is-medium is-size-7-mobile">
                      <ul>
                        <li>
                          <Link to={`/account`}>Detail</Link>
                        </li>
                        <li>
                          <Link to={`/account/2fa`}>2FA</Link>
                        </li>
                        <li className="is-active">
                          <Link>
                            <strong>
                              More&nbsp;&nbsp;
                              <FontAwesomeIcon
                                className="mdi"
                                icon={faEllipsis}
                              />
                            </strong>
                          </Link>
                        </li>
                      </ul>
                    </div>

                    {/* Page Menu Options (Tablet++) */}
                    <section className="hero is-hidden-mobile">
                      <div className="hero-body has-text-centered">
                        <div className="container">
                          <div className="columns is-vcentered is-multiline">
                            <div className="column">
                              <BubbleLink
                                title={`Change Password`}
                                subtitle={`Change the password credentials you use to log.`}
                                faIcon={faKey}
                                url={`/account/more/change-password`}
                                bgColour={`has-background-info-dark`}
                              />
                            </div>
                          </div>
                        </div>
                      </div>
                    </section>

                    {/* Page Menu Options (Mobile Only) */}
                    <div
                      className="has-background-white-ter is-hidden-tablet mb-6 p-5"
                      style={{ borderRadius: "15px" }}
                    >
                      <table className="is-fullwidth has-background-white-ter table">
                        <thead>
                          <tr>
                            <th colSpan="2">Menu</th>
                          </tr>
                        </thead>
                        <tbody>
                          <tr>
                            <td>
                              <FontAwesomeIcon className="fas" icon={faKey} />
                              &nbsp;Change Password
                            </td>
                            <td>
                              <div className="buttons is-right">
                                <Link
                                  to={`/account/more/change-password`}
                                  className="is-small"
                                >
                                  View&nbsp;
                                  <FontAwesomeIcon
                                    className="mdi"
                                    icon={faChevronRight}
                                  />
                                </Link>
                              </div>
                            </td>
                          </tr>
                          {/* End Associates */}
                        </tbody>
                      </table>
                    </div>
                    {/* END Page Menu Options (Mobile Only) */}

                    {/* Bottom Navigation */}
                    <div className="columns pt-5">
                      <div className="column is-half">
                        <Link
                          className="button is-medium is-fullwidth-mobile"
                          to={dashboardURL}
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back to Dashboard
                        </Link>
                      </div>
                      <div className="column is-medium is-half has-text-right"></div>
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

export default AccountMoreLaunchpad;
