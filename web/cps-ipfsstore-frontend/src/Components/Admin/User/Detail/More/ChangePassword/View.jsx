import React, { useState, useEffect } from "react";
import { Link, Navigate, useParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faBuildingUser,
  faHomeUser,
  faKey,
  faImage,
  faPaperclip,
  faAddressCard,
  faSquarePhone,
  faTasks,
  faTachometer,
  faPlus,
  faArrowLeft,
  faCheckCircle,
  faUserCircle,
  faGauge,
  faPencil,
  faUsers,
  faEye,
  faIdCard,
  faAddressBook,
  faContactCard,
  faChartPie,
  faBuilding,
  faEllipsis,
  faCircleExclamation,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { getUserDetailAPI } from "../../../../../../API/user";
import FormErrorBox from "../../../../../Reusable/FormErrorBox";
import FormInputField from "../../../../../Reusable/FormInputField";
import FormSelectField from "../../../../../Reusable/FormSelectField";
import PageLoadingContent from "../../../../../Reusable/PageLoadingContent";
import AlertBanner from "../../../../../Reusable/EveryPage/AlertBanner";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
} from "../../../../../../AppState";
import { CLIENT_ORGANIZATION_TYPE_OPTIONS_WITH_EMPTY_OPTIONS } from "../../../../../../Constants/FieldOptions";
import {
  EXECUTIVE_ROLE_ID,
  MANAGEMENT_ROLE_ID,
  FRONTLINE_ROLE_ID,
  ASSOCIATE_ROLE_ID,
  CUSTOMER_ROLE_ID,
} from "../../../../../../Constants/App";
import { postUserChangePasswordAPI } from "../../../../../../API/user";

function AdminUserMoreOperationChangePassword() {
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
  const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [password, setPassword] = useState("");
  const [passwordRepeated, setPasswordRepeated] = useState("");

  ////
  //// Event handling.
  ////

  const onSubmitClick = () => {
    setErrors({});
    setFetching(true);
    postUserChangePasswordAPI(
      {
        user_id: id,
        password: password,
        password_repeated: passwordRepeated,
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

  function onUserDetailSuccess(response) {
    console.log("onUserDetailSuccess: Starting...");
    setCurrentUser(response);
  }

  function onUserDetailError(apiErr) {
    console.log("onUserDetailError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onUserDetailDone() {
    console.log("onUserDetailDone: Starting...");
    setFetching(false);
  }

  // --- ChangePassword --- //

  function onChangePasswordSuccess(response) {
    console.log("onChangePasswordSuccess: Starting...");

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Password changed");
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
      setErrors({});
      getUserDetailAPI(
        onUserDetailSuccess,
        onUserDetailError,
        onUserDetailDone,
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
                  <FontAwesomeIcon className="fas" icon={faKey} />
                  &nbsp;Password
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
          {currentUser && currentUser.status === 2 && (
            <AlertBanner message="Archived" status="info" />
          )}

          {/* Page Title */}
          <h1 className="title is-2">
            <FontAwesomeIcon className="fas" icon={faUserCircle} />
            &nbsp;User
          </h1>
          <h4 className="subtitle is-4">
            <FontAwesomeIcon className="fas" icon={faEye} />
            &nbsp;Detail
          </h4>
          <hr />

          {/* Page */}
          <nav className="box">
            {/* Title + Options */}
            {currentUser && (
              <div className="columns">
                <div className="column">
                  <p className="title is-4">
                    <FontAwesomeIcon className="fas" icon={faKey} />
                    &nbsp;Change Password
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

                {currentUser && (
                  <div className="container">
                    <article className="message is-warning">
                      <div className="message-body">
                        <p className="title is-4">
                          <FontAwesomeIcon
                            className="fas"
                            icon={faCircleExclamation}
                          />
                          &nbsp;Warning
                        </p>
                        <p>
                          You are about to <b>change the password</b> for this
                          user. Please make sure you enter it correctly or
                          else user will be locked out of the user and
                          requiring password reseting. Are you sure you want to
                          continue?
                        </p>
                      </div>
                    </article>

                    <FormInputField
                      label="Password"
                      name="password"
                      placeholder="Password input"
                      value={password}
                      errorText={errors && errors.password}
                      helpText=""
                      onChange={(e) => setPassword(e.target.value)}
                      isRequired={true}
                      maxWidth="380px"
                      type="password"
                    />

                    <FormInputField
                      label="Password Repeated"
                      name="passwordRepeated"
                      placeholder="Password input again"
                      value={passwordRepeated}
                      errorText={errors && errors.passwordRepeated}
                      helpText=""
                      onChange={(e) => setPasswordRepeated(e.target.value)}
                      isRequired={true}
                      maxWidth="380px"
                      type="password"
                    />

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
                          className="button is-danger is-fullwidth-mobile"
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

export default AdminUserMoreOperationChangePassword;
