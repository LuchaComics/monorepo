import React, { useState, useEffect } from "react";
import { Link, Navigate, useSearchParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArrowLeft,
  faTasks,
  faTachometer,
  faPlus,
  faTimesCircle,
  faCheckCircle,
  faGauge,
  faUsers,
  faEye,
  faBookOpen,
  faMagnifyingGlass,
  faBalanceScale,
  faUser,
  faArrowUpRightFromSquare,
  faCog,
  faArrowRight
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import useLocalStorage from "../../../../Hooks/useLocalStorage";
import FormErrorBox from "../../../Reusable/FormErrorBox";
import FormInputField from "../../../Reusable/FormInputField";
import FormDateField from "../../../Reusable/FormDateField";
import FormTextareaField from "../../../Reusable/FormTextareaField";
import FormRadioField from "../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../Reusable/FormSelectField";
import FormCheckboxField from "../../../Reusable/FormCheckboxField";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import FormComicSignaturesTable from "../../../Reusable/FormComicSignaturesTable";
import {
  FINDING_WITH_EMPTY_OPTIONS,
  OVERALL_NUMBER_GRADE_WITH_EMPTY_OPTIONS,
  PUBLISHER_NAME_WITH_EMPTY_OPTIONS,
  CPS_PERCENTAGE_GRADE_WITH_EMPTY_OPTIONS,
  ISSUE_COVER_YEAR_OPTIONS,
  ISSUE_COVER_MONTH_WITH_EMPTY_OPTIONS,
  SPECIAL_DETAILS_WITH_EMPTY_OPTIONS,
  RETAILER_AVAILABLE_SERVICE_TYPE_WITH_EMPTY_OPTIONS,
  SUBMISSION_KEY_ISSUE_WITH_EMPTY_OPTIONS,
  SUBMISSION_PRINTING_WITH_EMPTY_OPTIONS,
} from "../../../../Constants/FieldOptions";
import {
  SERVICE_TYPE_PRE_SCREENING_SERVICE,
  SERVICE_TYPE_CPS_CAPSULE_INDIE_MINT_GEM,
  SERVICE_TYPE_CPS_CAPSULE_U_GRADE_SIGNATURE_COLLECTION,
} from "../../../../Constants/App";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
} from "../../../../AppState";
import {
  addComicSubmissionState,
  ADD_COMIC_SUBMISSION_STATE_DEFAULT,
} from "../../../../AppState";


function RetailerComicSubmissionAddStep2() {
  ////
  //// URL Parameters.
  ////

  const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140
  const customerID = searchParams.get("customer_id");
  const customerName = searchParams.get("customer_name");
  const fromPage = searchParams.get("from");
  const shouldClear = searchParams.get("clear");

  console.log("customer_id:", customerID, "customer_name:", customerName, "from:", fromPage);

  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);
  const [currentUser] = useRecoilState(currentUserState);
  const [addComicSubmission, setAddComicSubmission] = useRecoilState(addComicSubmissionState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [showCancelWarning, setShowCancelWarning] = useState(false);

  ////
  //// Event handling.
  ////

  const onSaveAndContinueClick = (e) => {
    console.log("onSaveAndContinueClick: Beginning...");

    let newErrors = {};
    let hasErrors = false;

    // if (serviceType === undefined || serviceType === null || serviceType === 0 || serviceType === "") {
    //   newErrors["serviceType"] = "missing value";
    //   hasErrors = true;
    // }
    //
    // //
    // // CASE 1 of 2: Has errors.
    // //
    //
    // if (hasErrors) {
    //   console.log("onSaveAndContinueClick: Aboring because of error(s)");
    //
    //   // Set the associate based error validation.
    //   setErrors(newErrors);
    //
    //   // The following code will cause the screen to scroll to the top of
    //   // the page. Please see ``react-scroll`` for more information:
    //   // https://github.com/fisshy/react-scroll
    //   var scroll = Scroll.animateScroll;
    //   scroll.scrollToTop();
    //
    //   return;
    // }

    //
    // CASE 2 of 2: Has no errors.
    //

    console.log("onSaveAndContinueClick: Saving step 2 and redirecting to step 3.");

    // Variable holds a complete clone of the submission.
    let modifiedAddComicSubmission = { ...addComicSubmission };

    // storeID: currentUser.storeId,
    // collectibleType: 1, // 1=Comic, 2=Card
    // customerID: customerID,

    // // Update our clone.
    modifiedAddComicSubmission.customerID = customerID;
    modifiedAddComicSubmission.customerId = customerID;
    modifiedAddComicSubmission.customerName = customerName;
    modifiedAddComicSubmission.storeID = currentUser.storeId;
    modifiedAddComicSubmission.fromPage = fromPage

    // Save to persistent storage.
    setAddComicSubmission(modifiedAddComicSubmission);

    // Redirect to the next page.
    setForceURL("/submissions/comics/add/step-3")
  };

  // Function will filter the available options based on user's organization level.
  // Special thanks via:
  // https://github.com/LuchaComics/cps-frontend/issues/160
  const cpsPercentageGradeFilterOptions = (options, storeLevel) => {
    return options.filter((option) => {
      if (storeLevel === 1) {
        return option.value <= 96;
      }
      if (storeLevel === 2 || storeLevel === 3) {
        return option.value <= 98;
      }
      return false;
    });
  };

  // Function will filter the available options based on user's organization level.
  // Special thanks via:
  // https://github.com/LuchaComics/cps-frontend/issues/160
  const overallNumberGradeFilterOptions = (options, storeLevel) => {
    return options.filter((option) => {
      if (storeLevel === 1) {
        return option.value <= 9.6;
      }
      if (storeLevel === 2 || storeLevel === 3) {
        return option.value <= 9.8;
      }
      return false;
    });
  };

  ////
  //// API.
  ////

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.

      // Developer notes: If we have `clear=true` as URL argument then we clear
      // the previous inputted submissions.
      if (shouldClear === "true") {
          // Delete the previous submission filling details.
          console.log("deleting previous addComicSubmission:", addComicSubmission);
          setAddComicSubmission(ADD_COMIC_SUBMISSION_STATE_DEFAULT);
      }
    }

    return () => {
      mounted = false;
    };
  }, [shouldClear]);

  ////
  //// Component rendering.
  ////

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  // Apply service type limitation based on the retailer store's level.
  const conditionalServiceTypeOptions = ((currentUser) => {
    if (currentUser.storeLevel === 1 || currentUser.storeLevel === 2) {
      return RETAILER_AVAILABLE_SERVICE_TYPE_WITH_EMPTY_OPTIONS;
    } else {
      // DEVELOPERS NOTE: Level 3 retailer stores are allowed to add a
      // new type of service type.
      const newServiceTypeOptions = [
        ...RETAILER_AVAILABLE_SERVICE_TYPE_WITH_EMPTY_OPTIONS,
        {
          value: SERVICE_TYPE_CPS_CAPSULE_U_GRADE_SIGNATURE_COLLECTION,
          label: "CPS Capsule U-Grade Signature Collection",
        },
      ];
      return newServiceTypeOptions;
    }
  })(currentUser);

  // Render the JSX content.
  return (
    <>
      <div class="container">
        <section class="section">
          {fromPage !== "customer" ? (
            <>
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
                    <Link to="/submissions" aria-current="page">
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
                      <FontAwesomeIcon className="fas" icon={faPlus} />
                      &nbsp;New
                    </Link>
                  </li>
                </ul>
              </nav>

              {/* Mobile Breadcrumbs */}
              <nav
                class="breadcrumb is-hidden-desktop"
                aria-label="breadcrumbs"
              >
                <ul>
                  <li class="">
                    <Link to={`/submissions/comics`} aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                      &nbsp;Back to Comics
                    </Link>
                  </li>
                </ul>
              </nav>
            </>
          ) : (
            <>
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
                    <Link to="/customers" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faUsers} />
                      &nbsp;Customers
                    </Link>
                  </li>
                  <li class="">
                    <Link
                      to={`/customer/${customerID}/comics`}
                      aria-current="page"
                    >
                      <FontAwesomeIcon className="fas" icon={faEye} />
                      &nbsp;Detail (Comics)
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
              <nav
                class="breadcrumb is-hidden-desktop"
                aria-label="breadcrumbs"
              >
                <ul>
                  <li class="">
                    <Link to={`/submissions/comics`} aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                      &nbsp;Back to Comics
                    </Link>
                  </li>
                </ul>
              </nav>
            </>
          )}

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
                {fromPage !== "customer" ? (
                  <Link
                    class="button is-medium is-success"
                    to={`/submissions/comics/add/step-1/search`}
                  >
                    Yes
                  </Link>
                ) : (
                  <Link
                    class="button is-medium is-success"
                    to={`/customer/${customerID}/comics`}
                  >
                    Yes
                  </Link>
                )}
                <button
                  class="button is-medium "
                  onClick={(e) => setShowCancelWarning(false)}
                >
                  No
                </button>
              </footer>
            </div>
          </div>

          {/* Progress Wizard */}
          <nav className="box has-background-light">
            <p className="subtitle is-5">Step 2 of 10</p>
            <progress
              class="progress is-success"
              value="20"
              max="100"
            >
              20%
            </progress>
          </nav>

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faPlus} />
              &nbsp;New Online Comic Submission
            </p>
            {isFetching ? (
              <PageLoadingContent displayMessage={"Submitting..."} />
            ) : (
              <>
                <FormErrorBox errors={errors} />
                <p class="has-text-grey pb-4">
                  Please fill out all the required fields before submitting this
                  form.
                </p>
                <div class="container">
                  <p class="subtitle is-6">
                    <FontAwesomeIcon className="fas" icon={faCog} />
                    &nbsp;Settings
                  </p>
                  <hr />

                  <FormInputField
                    label="Store"
                    name="store"
                    placeholder="Text input"
                    value={currentUser.storeName}
                    helpText="Your organization that you will submit under this submission."
                    isRequired={true}
                    maxWidth="380px"
                    disabled={true}
                  />

                  {((customerID !== undefined && customerID !== null && customerID !== "") || addComicSubmission.customerName) && <>
                      <FormInputField
                        label="Customer"
                        name="customerName"
                        placeholder="Text input"
                        value={customerName || addComicSubmission.customerName}
                        helpText="The name of the customer for this submission."
                        isRequired={true}
                        maxWidth="380px"
                        disabled={true}
                      />
                  </>}

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <button
                        class="button is-medium is-fullwidth-mobile"
                        onClick={(e) => setShowCancelWarning(true)}
                      >
                        <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                        &nbsp;Cancel
                      </button>
                    </div>
                    <div class="column is-half has-text-right">
                      <button
                        class="button is-medium is-primary is-fullwidth-mobile"
                        onClick={onSaveAndContinueClick}
                      >
                        Save and Continue&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} />
                      </button>
                    </div>
                  </div>
                </div>
              </>
            )}
          </nav>
        </section>
      </div>
    </>
  );
}

export default RetailerComicSubmissionAddStep2;
