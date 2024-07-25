import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faCogs,
  faChartPie,
  faIdCard,
  faTasks,
  faTachometer,
  faPlus,
  faArrowLeft,
  faCheckCircle,
  faGauge,
  faPencil,
  faBuilding,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { getStoreDetailAPI } from "../../API/store";
import FormErrorBox from "../Reusable/FormErrorBox";
import PageLoadingContent from "../Reusable/PageLoadingContent";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
} from "../../AppState";
import FormRowText from "../Reusable/FormRowText";
import FormTextOptionRow from "../Reusable/FormRowTextOption";
import FormTextChoiceRow from "../Reusable/FormRowTextChoice";
import {
  ESTIMATED_SUBMISSION_PER_MONTH_WITH_EMPTY_OPTIONS,
  HOW_LONG_HAS_YOUR_STORE_BEEN_OPERATING_FOR_WITH_EMPTY_OPTIONS,
  ORGANIZATION_LEVEL_WITH_EMPTY_OPTIONS,
  USER_SPECIAL_COLLECTION_WITH_EMPTY_OPTIONS,
} from "../../Constants/FieldOptions";

function StoreDetail() {
  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);
  const [currentUser] = useRecoilState(currentUserState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [store, setStore] = useState({});

  ////
  //// Event handling.
  ////

  //

  ////
  //// API.
  ////

  function onStoreDetailSuccess(response) {
    console.log("onStoreDetailSuccess: Starting...");
    setStore(response);
  }

  function onStoreDetailError(apiErr) {
    console.log("onStoreDetailError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onStoreDetailDone() {
    console.log("onStoreDetailDone: Starting...");
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
      getStoreDetailAPI(
        currentUser.storeId,
        onStoreDetailSuccess,
        onStoreDetailError,
        onStoreDetailDone,
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
                <Link to="/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Dashboard
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faBuilding} />
                  &nbsp;My Store
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/dashboard`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Dashboard
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faBuilding} />
              &nbsp;My Store
            </p>
            <FormErrorBox errors={errors} />

            {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                {store && (
                  <div class="container">
                    <div class="tabs is-medium is-size-7-mobile">
                      <ul>
                        <li class={`is-active`}>
                          <Link>
                            <b>Detail</b>
                          </Link>
                        </li>
                        <li>
                          <Link to={`/store/${store.id}/purchases`}>
                            Purchases
                          </Link>
                        </li>
                        <li>
                          <Link to={`/store/${store.id}/credits`}>Credits</Link>
                        </li>
                      </ul>
                    </div>

                    <p class="subtitle is-6 pt-4">
                      <FontAwesomeIcon className="fas" icon={faIdCard} />
                      &nbsp;Identification
                    </p>
                    <hr />

                    <FormRowText label="Name" value={store.name} helpText="" />
                    <FormRowText
                      label="Website URL"
                      value={store.websiteUrl}
                      helpText=""
                    />

                    <p class="subtitle is-6">
                      <FontAwesomeIcon className="fas" icon={faChartPie} />
                      &nbsp;Metrics
                    </p>
                    <hr />

                    <FormTextOptionRow
                      label="How many comic books are you planning to submit to us per month?"
                      selectedValue={store.estimatedSubmissionsPerMonth}
                      options={
                        ESTIMATED_SUBMISSION_PER_MONTH_WITH_EMPTY_OPTIONS
                      }
                    />

                    <FormTextChoiceRow
                      label="Are you currently submitting to any other grading companies?"
                      value={store.hasOtherGradingService}
                      opt1Value={1}
                      opt1Label="Yes"
                      opt2Value={2}
                      opt2Label="No"
                    />

                    {store.hasOtherGradingService === 1 && (
                      <FormRowText
                        label="Other Grading Service Name"
                        value={store.otherGradingServiceName}
                        helpText=""
                      />
                    )}

                    <FormTextChoiceRow
                      label="Would you like to receive a welcome package? This package includes promotional items and tools to help you improve your submissions, as well as our service terms and conditions.?"
                      value={store.requestWelcomePackage}
                      opt1Value={1}
                      opt1Label="Yes"
                      opt2Value={2}
                      opt2Label="No"
                    />

                    <FormTextOptionRow
                      label="How long has your store been operating for?"
                      selectedValue={store.howLongStoreOperating}
                      options={
                        HOW_LONG_HAS_YOUR_STORE_BEEN_OPERATING_FOR_WITH_EMPTY_OPTIONS
                      }
                    />

                    <FormRowText
                      label="Tell us about your level of experience with grading comics? (Optional)"
                      value={store.gradingComicsExperience}
                      helpText=""
                    />

                    <FormRowText
                      label="Please describe how you could become a good retail partner for CPS"
                      value={store.retailPartnershipReason}
                      helpText=""
                    />

                    <FormRowText
                      label="Please describe how CPS could help you grow your business"
                      value={store.cpsPartnershipReason}
                      helpText=""
                    />

                    <p class="subtitle is-6">
                      <FontAwesomeIcon className="fas" icon={faCogs} />
                      &nbsp;Settings
                    </p>
                    <hr />

                    <FormTextChoiceRow
                      label="Status"
                      value={store.status}
                      opt1Value={1}
                      opt1Label="Pending"
                      opt2Value={2}
                      opt2Label="Active"
                      opt3Value={3}
                      opt3Label="Rejected"
                      opt4Value={4}
                      opt4Label="Error"
                      opt5Value={5}
                      opt5Label="Archived"
                    />

                    <FormTextOptionRow
                      label="Level"
                      selectedValue={store.level}
                      options={ORGANIZATION_LEVEL_WITH_EMPTY_OPTIONS}
                    />

                    <FormTextOptionRow
                      label="Special Collection"
                      selectedValue={store.specialCollection}
                      options={USER_SPECIAL_COLLECTION_WITH_EMPTY_OPTIONS}
                    />

                    <div class="columns pt-5">
                      <div class="column is-half">
                        <Link
                          class="button is-fullwidth-mobile"
                          to={"/dashboard"}
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back to Dashboard
                        </Link>
                      </div>
                      <div class="column is-half has-text-right">
                        <Link
                          to={"/store/update"}
                          class="button is-primary is-fullwidth-mobile"
                        >
                          <FontAwesomeIcon className="fas" icon={faPencil} />
                          &nbsp;Edit
                        </Link>
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

export default StoreDetail;
