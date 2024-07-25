import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faTachometer,
  faPlus,
  faEye,
  faArrowLeft,
  faCheckCircle,
  faPencil,
  faGauge,
  faBookOpen,
  faMagnifyingGlass,
  faBalanceScale,
  faUser,
  faArrowUpRightFromSquare,
  faArrowRight,
} from "@fortawesome/free-solid-svg-icons";
import Select from "react-select";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import useLocalStorage from "../../../../Hooks/useLocalStorage";
import { getComicSubmissionDetailAPI } from "../../../../API/ComicSubmission";
import { getUserDetailAPI } from "../../../../API/user";
import FormErrorBox from "../../../Reusable/FormErrorBox";
import FormInputField from "../../../Reusable/FormInputField";
import FormTextareaField from "../../../Reusable/FormTextareaField";
import FormRadioField from "../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../Reusable/FormMultiSelectField";
import FormCheckboxField from "../../../Reusable/FormCheckboxField";
import FormSelectField from "../../../Reusable/FormSelectField";
import FormDateField from "../../../Reusable/FormDateField";
import FormCountryField from "../../../Reusable/FormCountryField";
import FormRegionField from "../../../Reusable/FormRegionField";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import {
  FINDING_OPTIONS,
  OVERALL_NUMBER_GRADE_OPTIONS,
  PUBLISHER_NAME_OPTIONS,
  CPS_PERCENTAGE_GRADE_OPTIONS,
  HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS,
} from "../../../../Constants/FieldOptions";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../AppState";
import FormRowText from "../../../Reusable/FormRowText";
import FormTextYesNoRow from "../../../Reusable/FormRowTextYesNo";
import FormTextOptionRow from "../../../Reusable/FormRowTextOption";
import FormTextChoiceRow from "../../../Reusable/FormRowTextChoice";

function RetailerComicSubmissionDetailForCustomer() {
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
  const [submission, setComicSubmission] = useState({});
  const [user, setUser] = useState({});
  const [showCustomerEditOptions, setShowCustomerEditOptions] = useState(false);

  ////
  //// Event handling.
  ////

  ////
  //// API.
  ////

  // --- User --- //

  function onUserDetailSuccess(response) {
    console.log("onUserDetailSuccess: Starting...");
    setUser(response);
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

  // --- Comic Submissions --- //

  function onComicSubmissionDetailSuccess(response) {
    console.log("onComicSubmissionDetailSuccess: Starting...");
    setComicSubmission(response);

    if (
      response.customerId !== undefined &&
      response.customerId !== null &&
      response.customerId !== "" &&
      response.customerId !== "000000000000000000000000"
    ) {
      setFetching(true);
      getUserDetailAPI(
        response.customerId,
        onUserDetailSuccess,
        onUserDetailError,
        onUserDetailDone,
        onUnauthorized,
      );
    }
  }

  function onComicSubmissionDetailError(apiErr) {
    console.log("onComicSubmissionDetailError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onComicSubmissionDetailDone() {
    console.log("onComicSubmissionDetailDone: Starting...");
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
      window.scrollTo(0, 0); // Start the page at the top of the page.

      setFetching(true);
      getComicSubmissionDetailAPI(
        id,
        onComicSubmissionDetailSuccess,
        onComicSubmissionDetailError,
        onComicSubmissionDetailDone,
        onUnauthorized,
      );
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

  return (
    <>
      <div class={`modal ${showCustomerEditOptions ? "is-active" : ""}`}>
        <div class="modal-background"></div>
        <div class="modal-card">
          <header class="modal-card-head">
            <p class="modal-card-title">Customer Edit</p>
            <button
              class="delete"
              aria-label="close"
              onClick={(e) => setShowCustomerEditOptions(false)}
            ></button>
          </header>
          <section class="modal-card-body">
            To edit the customer, please select one of the following option:
            {/*
                            <br /><br />
                            <Link to={`/submissions/comic/${submission.id}/edit-customer`} class="button is-primary" disabled={true}>Edit Current Customer</Link> */}
            <br />
            <br />
            <Link
              to={`/submissions/comic/${submission.id}/cust/search`}
              class="button is-primary"
            >
              Pick a Different Customer
            </Link>
          </section>
          <footer class="modal-card-foot">
            <button
              class="button"
              onClick={(e) => setShowCustomerEditOptions(false)}
            >
              Close
            </button>
          </footer>
        </div>
      </div>

      <div class="container">
        <section class="section">
          {/* Desktop Breadcrumbs */}
          <nav class="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to="/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Dashboard
                </Link>
              </li>
              <li class="">
                <Link to="/submissions/comics" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faTasks} />
                  &nbsp;Online Submissions
                </Link>
              </li>
              <li class="">
                <Link to="/submissions/comics" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faBookOpen} />
                  &nbsp;Comics
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail (Customer)
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/submissions/comics`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Comics
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faTasks} />
              &nbsp;Comic Submission
            </p>
            <FormErrorBox errors={errors} />

            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                {submission && (
                  <div class="container">
                    <div class="tabs is-medium is-size-7-mobile">
                      <ul>
                        <li>
                          <Link to={`/submissions/comic/${id}`}>Detail</Link>
                        </li>
                        <li class="is-active">
                          <Link>
                            <b>Customer</b>
                          </Link>
                        </li>
                        <li>
                          <Link to={`/submissions/comic/${id}/comments`}>
                            Comments
                          </Link>
                        </li>
                        <li>
                          <Link to={`/submissions/comic/${id}/file`}>File</Link>
                        </li>
                        <li>
                          <Link to={`/submissions/comic/${id}/attachments`}>
                            Attachments
                          </Link>
                        </li>
                      </ul>
                    </div>
                    {submission.customerId !== "000000000000000000000000" &&
                    user !== undefined &&
                    user !== null &&
                    user !== "" ? (
                      <>
                        <p class="subtitle is-6 pt-4">
                          <FontAwesomeIcon className="fas" icon={faUser} />
                          &nbsp;Customer
                        </p>
                        <hr />
                        <p class="pb-5">
                          <Link
                            to={`/customer/${submission.customerId}`}
                            target="_blank"
                            rel="noreferrer"
                          >
                            Click here&nbsp;
                            <FontAwesomeIcon
                              className="fas"
                              icon={faArrowUpRightFromSquare}
                            />
                          </Link>{" "}
                          to view the customer.
                        </p>

                        <FormRowText
                          label="Name"
                          value={`${user.firstName} ${user.lastName}`}
                          maxWidth="280px"
                          helpText={""}
                        />
                        <FormRowText
                          label="Email"
                          type="email"
                          value={user.email}
                          helpText={""}
                        />
                        <FormRowText
                          label="Phone"
                          type="phone"
                          value={user.phone}
                          helpText={""}
                        />
                        <FormRowText
                          label="Country"
                          value={user.country}
                          helpText=""
                        />
                        <FormRowText
                          label="Province/Territory"
                          value={user.region}
                          helpText=""
                        />
                        <FormRowText
                          label="City"
                          value={user.city}
                          helpText={""}
                        />
                        <FormRowText
                          label="Address Line 1"
                          value={user.addressLine1}
                          helpText={""}
                        />
                        <FormRowText
                          label="Address Line 2"
                          value={user.addressLine2}
                          helpText={""}
                        />
                        <FormRowText
                          label="Postal Code"
                          value={user.postalCode}
                          helpText={""}
                        />

                        <FormTextOptionRow
                          label="How did you hear about us?"
                          selectedValue={user.howDidYouHearAboutUs}
                          helpText=""
                          options={HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS}
                        />
                        {user.howDidYouHearAboutUs === 1 && (
                          <FormRowText
                            label="Other (Please specify):"
                            value={user.howDidYouHearAboutUsOther}
                            helpText=""
                          />
                        )}
                        <FormTextYesNoRow
                          label="I agree to receive electronic updates from my local retailer and CPS"
                          checked={user.agreePromotionsEmail}
                        />

                        <div class="columns pt-5">
                          <div class="column is-half">
                            <Link
                              to={`/submissions/comics`}
                              class="button is-medium is-fullwidth-mobile"
                            >
                              <FontAwesomeIcon
                                className="fas"
                                icon={faArrowLeft}
                              />
                              &nbsp;Back to Comic Submissions
                            </Link>
                          </div>
                          <div class="column is-half has-text-right">
                            <Link
                              onClick={(e) => setShowCustomerEditOptions(true)}
                              class="button is-medium is-primary is-fullwidth-mobile"
                            >
                              <FontAwesomeIcon
                                className="fas"
                                icon={faPencil}
                              />
                              &nbsp;Edit Customer
                            </Link>
                          </div>
                        </div>
                      </>
                    ) : (
                      <>
                        <section class="hero is-medium is-warning">
                          <div class="hero-body">
                            <p class="title">
                              <FontAwesomeIcon className="fas" icon={faUser} />
                              &nbsp;No Customer
                            </p>
                            <p class="subtitle">
                              Assign a customer to this comic by clicking below:
                              <br />
                              <br />
                              <Link
                                to={`/submissions/comic/${submission.id}/cust/search`}
                              >
                                Assign&nbsp;
                                <FontAwesomeIcon
                                  className="fas"
                                  icon={faArrowRight}
                                />
                              </Link>
                            </p>
                          </div>
                        </section>
                      </>
                    )}
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

export default RetailerComicSubmissionDetailForCustomer;
