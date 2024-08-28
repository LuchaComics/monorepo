import React, { useState, useEffect } from "react";
import { Link, Navigate, useParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faBuildingUser,
  faHomeUser,
  faMobile,
  faImage,
  faPaperclip,
  faAddressCard,
  faSquarePhone,
  faTasks,
  faTachometer,
  faPlus,
  faArrowLeft,
  faCheckCircle,
  faUserTie,
  faGauge,
  faPencil,
  faUsers,
  faEye,
  fidCard,
  faAddressBook,
  faContactCard,
  faChartPie,
  faBuilding,
  faEllipsis,
  faCircleExclamation,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import FormErrorBox from "../../../../../Reusable/FormErrorBox";
import FormInputField from "../../../../../Reusable/FormInputField";
import FormSelectField from "../../../../../Reusable/FormSelectField";
import PageLoadingContent from "../../../../../Reusable/PageLoadingContent";
import AlertBanner from "../../../../../Reusable/EveryPage/AlertBanner";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../../../AppState";
import { postUserChangeTwoFactorAuthenticationAPI } from "../../../../../../API/user";
import { getUserDetailAPI } from "../../../../../../API/user";

function AdminUserMoreOperation2FAToggle() {
  ////
  //// URL Parameters.
  ////

  const { id } = useParams();

  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [User, setUser] = useState({});

  ////
  //// Event handling.
  ////

  const onSubmitClick = () => {
    setErrors({});
    setFetching(true);
    postUserChangeTwoFactorAuthenticationAPI(
      {
        user_id: id,
        otp_enabled: !User.otpEnabled,
      },
      onChangePasswordSuccess,
      onChangePasswordError,
      onChangePasswordDone,
      onUnauthorized,
    );
  };

  ////
  //// API.
  ////

  // --- Detail --- //

  function onSuccess(response) {
    console.log("onSuccess: Starting...");
    setUser(response);
  }

  function onError(apiErr) {
    console.log("onError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onDone() {
    console.log("onDone: Starting...");
    setFetching(false);
  }

  // --- ChangePassword --- //

  function onChangePasswordSuccess(response) {
    console.log("onChangePasswordSuccess: Starting...");

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("2FA changed");
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

    setForceURL("/admin/user/" + id + "/more");
  }

  function onChangePasswordError(apiErr) {
    console.log("onChangePasswordError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onChangePasswordDone() {
    console.log("onChangePasswordDone: Starting...");
    setFetching(false);
  }

  const onUnauthorized = () => {
    setForceURL("/login?unauthorized=true"); // If token expired or user is not logged in, redirect back to login.
  };

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
      setFetching(true);
      getUserDetailAPI(id, onSuccess, onError, onDone);
    }

    return () => {
      mounted = false;
    };
  }, [id]);

  ////
  //// Component rendering.
  ////

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  console.log("User --->", User);
  console.log("User.otpEnabled --->", User.otpEnabled);

  return (
    <>
      <div className="container">
        <section className="section">
          {/* Desktop Breadcrumbs */}
          <nav
            className="breadcrumb is-hidden-touch"
            aria-label="breadcrumbs"
          >
            <ul>
              <li className="">
                <Link to="/admin/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Admin Dashboard
                </Link>
              </li>
              <li className="">
                <Link to="/admin/users" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faUsers} />
                  &nbsp;Users
                </Link>
              </li>
              <li className="">
                <Link to={`/admin/user/${id}/more`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail (More)
                </Link>
              </li>
              <li className="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faMobile} />
                  &nbsp;Two-Factor Authentication
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
                <Link to={`/admin/user/${id}/more`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Detail (More)
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page banner */}
          {User && User.status === 100 && (
            <AlertBanner message="Archived" status="info" />
          )}

          {/* Page Title */}
          <h1 className="title is-2">
            <FontAwesomeIcon className="fas" icon={faUserTie} />
            &nbsp;User
          </h1>
          <h4 className="subtitle is-4">
            <FontAwesomeIcon className="fas" icon={faEye} />
            &nbsp;Detail
          </h4>
          <hr />

          {/* Page */}
          <nav className="box">
            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                {/* Title + Options */}
                {User && (
                  <div className="columns">
                    <div className="column is-8">
                      <p className="title is-4">
                        <FontAwesomeIcon className="fas" icon={faMobile} />
                        &nbsp;Change Two-Factor Authentication
                      </p>
                    </div>
                    <div className="column has-text-right"></div>
                  </div>
                )}

                {/* <p className="pb-4">Please fill out all the required fields before submitting this form.</p> */}

                <FormErrorBox errors={errors} />

                {User && (
                  <div className="container">
                    {!User.otpEnabled ? (
                      <article className="message is-success">
                        <div className="message-body">
                          <p className="title is-4">
                            <FontAwesomeIcon
                              className="fas"
                              icon={faCheckCircle}
                            />
                            &nbsp;Enable 2FA
                          </p>
                          <p>
                            You are about to <b>enable 2FA</b> for this User
                            member. This operation will force User member on
                            next successfull login to be taken through a{" "}
                            <b>3-step wizard</b> to setup 2FA. Afterwords every
                            time the User member logs in, they will be asked
                            to carry out a 2FA process. Are you sure you want to
                            continue?
                          </p>
                        </div>
                      </article>
                    ) : (
                      <article className="message is-warning">
                        <div className="message-body">
                          <p className="title is-4">
                            <FontAwesomeIcon
                              className="fas"
                              icon={faCircleExclamation}
                            />
                            &nbsp;Remove 2FA
                          </p>
                          <p>
                            You are about to <b>remove 2FA</b> for this User
                            member. This operation will remove previous 2FA
                            settup codes and disable 2FA on login for this
                            User. This is recommended if the user lost their
                            2FA codes from their device. Are you sure you want
                            to continue?
                          </p>
                        </div>
                      </article>
                    )}

                    {/* Bottom Navigation */}
                    <div className="columns pt-5">
                      <div className="column is-half">
                        <Link
                          className="button is-fullwidth-mobile"
                          to={`/admin/user/${id}/more`}
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back to Detail (More)
                        </Link>
                      </div>
                      <div className="column is-half has-text-right">
                        <button
                          className={`button ${!User.otpEnabled ? "is-success" : "is-warning"} is-fullwidth-mobile`}
                          onClick={onSubmitClick}
                        >
                          <FontAwesomeIcon
                            className="fas"
                            icon={faCheckCircle}
                            type="button"
                          />
                          &nbsp;Confirm and Submit
                        </button>
                      </div>
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

export default AdminUserMoreOperation2FAToggle;
