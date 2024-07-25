import React, { useState, useEffect } from "react";
import { Link, Navigate, useSearchParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faBookOpen,
  faTachometer,
  faPlus,
  faDownload,
  faArrowLeft,
  faTable,
  faCheckCircle,
  faCheck,
  faGauge,
  faUsers,
  faArrowUpRightFromSquare,
  faArrowRight,
} from "@fortawesome/free-solid-svg-icons";
import Select from "react-select";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import useLocalStorage from "../../../../Hooks/useLocalStorage";
import { getUserListAPI } from "../../../../API/user";
import FormErrorBox from "../../../Reusable/FormErrorBox";
import FormInputField from "../../../Reusable/FormInputField";
import FormTextareaField from "../../../Reusable/FormTextareaField";
import FormRadioField from "../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../Reusable/FormSelectField";
import FormInputFieldWithButton from "../../../Reusable/FormInputFieldWithButton";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../AppState";

function AdminComicSubmissionAddStep1WithResult() {
  ////
  //// URL Parameters.
  ////

  const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140

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
  const [users, setUsers] = useState({});
  const [showCancelWarning, setShowCancelWarning] = useState(false);

  ////
  //// Event handling.
  ////

  ////
  //// API.
  ////

  function onUserListSuccess(response) {
    console.log("onUserListSuccess: Starting...");
    if (response.results !== null) {
      setUsers(response);
    }
  }

  function onUserListError(apiErr) {
    console.log("onUserListError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onUserListDone() {
    console.log("onUserListDone: Starting...");
    setFetching(false);
  }

  // --- All --- //

  const onUnauthorized = () => {
    setForceURL("/login?unauthorized=true"); // If token expired or user is not logged in, redirect back to login.
  };

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      console.log("useEffect: Starting.");

      window.scrollTo(0, 0); // Start the page at the top of the page.
      setFetching(true); // Let user knows that we are making an API endpoint.

      let queryParams = new Map(); // Create the URL map we'll be using when calling the backend.

      // CASE 1: Search
      const searchKeyword = searchParams.get("search");
      if (
        searchKeyword !== undefined &&
        searchKeyword !== null &&
        searchKeyword !== ""
      ) {
        queryParams.set("search", searchKeyword);
      }

      // CASE 2: First name.
      const firstName = searchParams.get("first_name");
      if (firstName !== undefined && firstName !== null && firstName !== "") {
        queryParams.set("first_name", firstName);
      }

      // CASE 3: Last name.
      const lastName = searchParams.get("last_name");
      if (lastName !== undefined && lastName !== null && lastName !== "") {
        queryParams.set("last_name", lastName);
      }

      // CASE 4: Phone.
      const phone = searchParams.get("phone");
      if (phone !== undefined && phone !== null && phone !== "") {
        queryParams.set("phone", phone);
      }

      // CASE 5: Email.
      const email = searchParams.get("email");
      if (email !== undefined && email !== null && email !== "") {
        queryParams.set("email", email);
      }

      // Submit the list request to our backend.
      getUserListAPI(
        queryParams,
        onUserListSuccess,
        onUserListError,
        onUserListDone,
        onUnauthorized,
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
      <div class="container">
        <section class="section">
          {/* Desktop Breadcrumbs */}
          <nav class="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to="/admin/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Admin Dashboard
                </Link>
              </li>
              <li class="">
                <Link to="/admin/submissions" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faTasks} />
                  &nbsp;Submissions
                </Link>
              </li>
              <li class="">
                <Link to="/admin/submissions/comics" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faBookOpen} />
                  &nbsp;Comics
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faPlus} />
                  &nbsp;New
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/admin/submissions/comics`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Comics
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modals */}
          <div class={`modal ${showCancelWarning ? "is-active" : ""}`}>
            <div class="modal-background"></div>
            <div class="modal-card">
              <header class="modal-card-head">
                <p class="modal-card-title">Are you sure?</p>
                <button
                  class="delete"
                  aria-label="close"
                  onClick={(e) => setShowCancelWarning(false)}
                ></button>
              </header>
              <section class="modal-card-body">
                Your submission will be cancelled and your work will be lost.
                This cannot be undone. Do you want to continue?
              </section>
              <footer class="modal-card-foot">
                <Link class="button is-success" to={`/dashboard`}>
                  Yes
                </Link>
                <button class="button" onClick={(e) => null}>
                  No
                </button>
              </footer>
            </div>
          </div>

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faPlus} />
              &nbsp;New Online Comic Submission
            </p>
            <p class="has-text-grey pb-4">
              Please select the user from the following results.
            </p>
            <FormErrorBox errors={errors} />

            <div class="container pb-5">
              <p class="subtitle is-6">
                <FontAwesomeIcon className="fas" icon={faUsers} />
                &nbsp;Results
              </p>
              <hr />

              {isFetching ? (
                <PageLoadingContent displayMessage={"Loading..."} />
              ) : (
                <>
                  {users && users.results && users.results.length > 0 ? (
                    <div class="columns is-multiline">
                      {users.results.map(function (user, i) {
                        return (
                          <div class="column is-one-quarter" key={user.id}>
                            <span className="has-text-grey-light is-size-7">
                              User found via{" "}
                              <a
                                href={`/admin/store/${user.storeId}`}
                                target="_blank"
                                rel="noreferrer"
                              >
                                {user.storeName}&nbsp;
                                <FontAwesomeIcon
                                  className="fas"
                                  icon={faArrowUpRightFromSquare}
                                />
                              </a>
                            </span>
                            <article class="message is-primary">
                              <div class="message-body">
                                <p>
                                  <Link
                                    to={`/admin/user/${user.id}`}
                                    target="_blank"
                                    rel="noreferrer"
                                  >
                                    <b>{user.name}</b>&nbsp;
                                    <FontAwesomeIcon
                                      className="fas"
                                      icon={faArrowUpRightFromSquare}
                                    />
                                  </Link>
                                </p>
                                <p>
                                  {user.country}&nbsp;{user.region}&nbsp;
                                  {user.city}
                                </p>
                                <p>
                                  {user.addressLine1}, {user.postalCode}
                                </p>
                                <p>
                                  <a href={`mailto:${user.email}`}>
                                    {user.email}
                                  </a>
                                </p>
                                <p>
                                  <a href={`tel:${user.phone}`}>{user.phone}</a>
                                </p>
                                <br />
                                <Link
                                  class="button is-medium is-primary"
                                  to={`/admin/submissions/comics/add?user_id=${user.id}`}
                                >
                                  <FontAwesomeIcon
                                    className="fas"
                                    icon={faCheckCircle}
                                  />
                                  &nbsp;Pick
                                </Link>
                              </div>
                            </article>
                          </div>
                        );
                      })}
                    </div>
                  ) : (
                    <section class="hero is-medium has-background-white-ter">
                      <div class="hero-body">
                        <p class="title">
                          <FontAwesomeIcon className="fas" icon={faTable} />
                          &nbsp;No Users
                        </p>
                        <p class="subtitle">
                          No results were found in the search.{" "}
                          <Link
                            class="is-medium is-warning"
                            to="/admin/customers/add"
                            target="_blank"
                            rel="noreferrer"
                          >
                            Click here&nbsp;
                            <FontAwesomeIcon
                              className="fas"
                              icon={faArrowUpRightFromSquare}
                            />
                          </Link>{" "}
                          to create a new customer or{" "}
                          <Link
                            class="is-medium is-danger"
                            to="/admin/submissions/comics/add"
                          >
                            click here&nbsp;
                            <FontAwesomeIcon
                              className="fas"
                              icon={faArrowRight}
                            />
                          </Link>{" "}
                          to continue without a customer.
                        </p>
                      </div>
                    </section>
                  )}
                </>
              )}
            </div>

            <div class="columns pt-5">
              <div class="column is-half">
                <Link
                  class="button is-medium is-fullwidth-mobile"
                  to="/admin/submissions/comics/add/search"
                >
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back
                </Link>
              </div>
              <div class="column is-half has-text-right">
                {/*
                                <button class="button is-primary is-hidden-touch" onClick={null}><FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;Next</button>
                                <button class="button is-primary is-fullwidth is-hidden-desktop" onClick={null}><FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;Next</button>
                                */}
              </div>
            </div>
          </nav>
        </section>
      </div>
    </>
  );
}

export default AdminComicSubmissionAddStep1WithResult;
