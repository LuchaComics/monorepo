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
  faIdCard,
  faCog,
  faFileSignature,
  faPencil,
  faClipboardCheck
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import useLocalStorage from "../../../../Hooks/useLocalStorage";
import { postComicSubmissionCreateAPI } from "../../../../API/ComicSubmission";
import FormErrorBox from "../../../Reusable/FormErrorBox";
import FormInputField from "../../../Reusable/FormInputField";
import FormDateField from "../../../Reusable/FormDateField";
import FormTextareaField from "../../../Reusable/FormTextareaField";
import FormRadioField from "../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../Reusable/FormSelectField";
import FormCheckboxField from "../../../Reusable/FormCheckboxField";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import DataDisplayRowText from "../../../Reusable/DataDisplayRowText";
import DataDisplayRowSelect from "../../../Reusable/DataDisplayRowSelect";
import DataDisplayRowCheckbox from "../../../Reusable/DataDisplayRowCheckbox";
import DataDisplayRowRadio from "../../../Reusable/DataDisplayRowRadio";
import FormComicSignaturesTable from "../../../Reusable/FormComicSignaturesTable";
import {
  FINDING_WITH_EMPTY_OPTIONS,
  OVERALL_NUMBER_GRADE_WITH_EMPTY_OPTIONS,
  PUBLISHER_NAME_WITH_EMPTY_OPTIONS,
  CPS_PERCENTAGE_GRADE_WITH_EMPTY_OPTIONS,
  ISSUE_COVER_YEAR_OPTIONS,
  ISSUE_COVER_MONTH_WITH_EMPTY_OPTIONS,
  SPECIAL_DETAILS_WITH_EMPTY_OPTIONS,
  SUBMISSION_KEY_ISSUE_WITH_EMPTY_OPTIONS,
  SUBMISSION_PRINTING_WITH_EMPTY_OPTIONS,
  SERVICE_TYPE_WITH_EMPTY_OPTIONS
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


function AdminComicSubmissionAddStep10() {
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
  const [gradingScale, setGradingScale] = useState(0);
  const [overallLetterGrade, setOverallLetterGrade] = useState("");
  const [overallNumberGrade, setOverallNumberGrade] = useState("");
  const [cpsPercentageGrade, setCpsPercentageGrade] = useState("");

  ////
  //// Event handling.
  ////

  const onSubmitClick = (e) => {
    console.log("onSubmitClick: Beginning...");
    console.log("onSubmitClick: Generating payload for submission.");
    setFetching(true);
    setErrors({});

    // Variable holds a complete clone of the submission.
    let modifiedAddComicSubmission = { ...addComicSubmission };
    
    // Submit to the backend.
    console.log("onSubmitClick: payload:", addComicSubmission);
    postComicSubmissionCreateAPI(
      modifiedAddComicSubmission,
      onComicSubmissionCreateSuccess,
      onComicSubmissionCreateError,
      onComicSubmissionCreateDone,
      onUnauthorized,
    );
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

  function onComicSubmissionCreateSuccess(response) {
    // For debugging purposes only.
    console.log("onComicSubmissionCreateSuccess: Starting...");
    console.log(response);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Comic submission created");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onComicSubmissionCreateSuccess: Delayed for 2 seconds.");
      console.log(
        "onComicSubmissionCreateSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    let urlParams = "?from="+addComicSubmission.fromPage + "&submission_id=" + response.id;

    // Redirect the user to a new page.
    setForceURL("/admin/submissions/comics/add/checkout" + urlParams);
  }

  function onComicSubmissionCreateError(apiErr) {
    console.log("onComicSubmissionCreateError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onComicSubmissionCreateError: Delayed for 2 seconds.");
      console.log(
        "onComicSubmissionCreateError: topAlertMessage, topAlertStatus:",
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

  function onComicSubmissionCreateDone() {
    console.log("onComicSubmissionCreateDone: Starting...");
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

  // The following code will check to see if we need to grant the 'is NM+' option is available to the user.
  let isNMPlusOpen = gradingScale === 1 && overallLetterGrade === "nm";

  // Apply the custom function to your options
  const cpsPercentageGradeFilteredOptions = cpsPercentageGradeFilterOptions(
    CPS_PERCENTAGE_GRADE_WITH_EMPTY_OPTIONS,
    currentUser.storeLevel,
  );
  const overallNumberGradeFilteredOptions = overallNumberGradeFilterOptions(
    OVERALL_NUMBER_GRADE_WITH_EMPTY_OPTIONS,
    currentUser.storeLevel,
  );

  // Render the JSX content.
  return (
    <>
      <div class="container">
        <section class="section">
          {addComicSubmission.fromPage !== "customer" ? (
            <>
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
                      &nbsp;Online Submissions
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
              <nav
                class="breadcrumb is-hidden-desktop"
                aria-label="breadcrumbs"
              >
                <ul>
                  <li class="">
                    <Link to={`/admin/submissions/comics`} aria-current="page">
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
                    <Link to="/admin/dashboard" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faGauge} />
                      &nbsp;Dashboard
                    </Link>
                  </li>
                  <li class="">
                    <Link to="/admin/customers" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faUsers} />
                      &nbsp;Customers
                    </Link>
                  </li>
                  <li class="">
                    <Link
                      to={`/admin/customer/${addComicSubmission.customerID}/comics`}
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
                    <Link to={`/admin/submissions/comics`} aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                      &nbsp;Back to Comics
                    </Link>
                  </li>
                </ul>
              </nav>
            </>
          )}

          {/* Modals */}
          {/* ------ */}

          {/* Progress Wizard */}
          <nav className="box has-background-success-light">
            <p className="subtitle is-5">Step 10 of 10</p>
            <progress
              class="progress is-success"
              value="100"
              max="100"
            >
              75%
            </progress>
          </nav>

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faPlus} />
              &nbsp;New Online Comic Submission
            </p>

            <p className="title is-4 pb-2">
              <FontAwesomeIcon className="fas" icon={faFileSignature} />
              &nbsp;Review
            </p>
            <p className="has-text-grey pb-4">
              Please review the following comic submission summary before submitting into the system.
            </p>

            {isFetching ? (
              <PageLoadingContent displayMessage={"Submitting..."} />
            ) : (
              <>
                <FormErrorBox errors={errors} />
                <div class="container">

                {/* STEP 2 of 10: OWNERSHIP  */}
                <p className="title is-5 mt-2">
                  <FontAwesomeIcon className="fas" icon={faIdCard} />
                  &nbsp;Ownership&nbsp;&nbsp;&nbsp;-&nbsp;&nbsp;&nbsp;
                  <Link to="/admin/submissions/comics/add/step-2">
                    <FontAwesomeIcon className="fas" icon={faPencil} />
                    &nbsp;Edit
                  </Link>
                </p>

                <DataDisplayRowText
                  label="Store"
                  value={addComicSubmission.storeName}
                />

                {addComicSubmission.customerId &&
                    <DataDisplayRowText
                      label="Customer"
                      value={addComicSubmission.customerName}
                    />
                }

                {/* STEP 3 OF 10: SETTINGS */}
                <p className="title is-5 mt-2">
                  <FontAwesomeIcon className="fas" icon={faCog} />
                  &nbsp;Settings&nbsp;&nbsp;&nbsp;-&nbsp;&nbsp;&nbsp;
                  <Link to="/admin/submissions/comics/add/step-3">
                    <FontAwesomeIcon className="fas" icon={faPencil} />
                    &nbsp;Edit
                  </Link>
                </p>

                <DataDisplayRowSelect
                  label="Service Type"
                  selectedValue={addComicSubmission.serviceType}
                  options={SERVICE_TYPE_WITH_EMPTY_OPTIONS}
                />

                {/* STEP 4 OF 10: BOOK INFORMATION */}
                <p className="title is-5 mt-2">
                  <FontAwesomeIcon className="fas" icon={faBookOpen} />
                  &nbsp;Book Information&nbsp;&nbsp;&nbsp;-&nbsp;&nbsp;&nbsp;
                  <Link to="/admin/submissions/comics/add/step-4">
                    <FontAwesomeIcon className="fas" icon={faPencil} />
                    &nbsp;Edit
                  </Link>
                </p>

                <DataDisplayRowText
                  label="Series Title"
                  value={addComicSubmission.seriesTitle}
                />

                <DataDisplayRowText
                  label="Issue Vol"
                  value={addComicSubmission.issueVol}
                />

                <DataDisplayRowText
                  label="Issue No"
                  value={addComicSubmission.issueNo}
                />

                <DataDisplayRowSelect
                  label="Issue Cover Year"
                  selectedValue={addComicSubmission.issueCoverYear}
                  options={ISSUE_COVER_YEAR_OPTIONS}
                />

                <DataDisplayRowSelect
                  label="Issue Cover Month"
                  selectedValue={addComicSubmission.issueCoverMonth}
                  options={ISSUE_COVER_MONTH_WITH_EMPTY_OPTIONS}
                />

                <DataDisplayRowSelect
                  label="Publisher Name"
                  selectedValue={addComicSubmission.publisherName}
                  options={PUBLISHER_NAME_WITH_EMPTY_OPTIONS}
                />
                {addComicSubmission.publisherName === 1 && (
                    <DataDisplayRowText
                       label="Publisher Name (Other)"
                       value={addComicSubmission.publisherNameOther}
                    />
                )}

                {/* STEP 5 OF 10: ADDITONAL BOOK INFORMATION */}
                <p className="title is-5 mt-2">
                  <FontAwesomeIcon className="fas" icon={faBookOpen} />
                  &nbsp;Additional Book Information&nbsp;&nbsp;&nbsp;-&nbsp;&nbsp;&nbsp;
                  <Link to="/admin/submissions/comics/add/step-5">
                    <FontAwesomeIcon className="fas" icon={faPencil} />
                    &nbsp;Edit
                  </Link>
                </p>

                <DataDisplayRowCheckbox
                  label="Is Key Issue?"
                  checked={addComicSubmission.isKeyIssue}
                />

                {addComicSubmission.isKeyIssue && (
                    <>
                        <DataDisplayRowSelect
                          label="Key Issue"
                          selectedValue={addComicSubmission.keyIssue}
                          options={SUBMISSION_KEY_ISSUE_WITH_EMPTY_OPTIONS}
                        />
                        {addComicSubmission.keyIssue === 1 && (
                            <DataDisplayRowText
                              label="Key Issue (Other)"
                              value={addComicSubmission.keyIssueOther}
                            />
                        )}
                        <DataDisplayRowText
                          label="Key Issue Detail"
                          value={addComicSubmission.keyIssueDetail}
                        />
                    </>
                )}

                <DataDisplayRowCheckbox
                  label="Is variant cover?"
                  checked={addComicSubmission.isVariantCover}
                />
                {addComicSubmission.isKeyIssue && (
                    <DataDisplayRowText
                      label="Variant cover detail"
                      value={addComicSubmission.variantCoverDetail}
                    />
                )}

                <DataDisplayRowSelect
                  label="Which printing is this?"
                  selectedValue={addComicSubmission.printing}
                  options={SUBMISSION_PRINTING_WITH_EMPTY_OPTIONS}
                />

                <DataDisplayRowSelect
                  label="Primary Label Details"
                  selectedValue={addComicSubmission.primaryLabelDetails}
                  options={SPECIAL_DETAILS_WITH_EMPTY_OPTIONS}
                />
                {addComicSubmission.primaryLabelDetails === 1 && (
                    <DataDisplayRowText
                      label="Primary Label Details (Other)"
                      value={addComicSubmission.primaryLabelDetailsOther}
                    />
                )}

                <DataDisplayRowText
                  label="Special Note (Optional)"
                  value={addComicSubmission.specialNotes}
                />

                {/* STEP 6 OF 10: BOOK SIGNATURES*/}
                <p className="title is-5 mt-2">
                  <FontAwesomeIcon className="fas" icon={faFileSignature} />
                  &nbsp;Book Signatures&nbsp;&nbsp;&nbsp;-&nbsp;&nbsp;&nbsp;
                  <Link to="/admin/submissions/comics/add/step-6">
                    <FontAwesomeIcon className="fas" icon={faPencil} />
                    &nbsp;Edit
                  </Link>
                </p>

                <DataDisplayRowText
                  label="Comic Signatures (Optional)"
                  value={`None`}
                />

                {/* STEP 7 OF 10: SUMMARY OF FINDINGS */}
                <p className="title is-5 mt-2">
                  <FontAwesomeIcon className="fas" icon={faMagnifyingGlass} />
                  &nbsp;Summary of Findings&nbsp;&nbsp;&nbsp;-&nbsp;&nbsp;&nbsp;
                  <Link to="/admin/submissions/comics/add/step-7">
                    <FontAwesomeIcon className="fas" icon={faPencil} />
                    &nbsp;Edit
                  </Link>
                </p>

                <DataDisplayRowRadio
                  label="Creases Finding"
                  value={addComicSubmission.creasesFinding}
                  opt1Value="pr"
                  opt1Label="Poor"
                  opt2Value="fr"
                  opt2Label="Fair"
                  opt3Value="gd"
                  opt3Label="Good"
                  opt4Value="vg"
                  opt4Label="Very good"
                  opt5Value="fn"
                  opt5Label="Fine"
                  opt6Value="vf"
                  opt6Label="Very Fine"
                  opt7Value="nm"
                  opt7Label="Near Mint"
                />

                <DataDisplayRowRadio
                  label="Tears Finding"
                  value={addComicSubmission.tearsFinding}
                  opt1Value="pr"
                  opt1Label="Poor"
                  opt2Value="fr"
                  opt2Label="Fair"
                  opt3Value="gd"
                  opt3Label="Good"
                  opt4Value="vg"
                  opt4Label="Very good"
                  opt5Value="fn"
                  opt5Label="Fine"
                  opt6Value="vf"
                  opt6Label="Very Fine"
                  opt7Value="nm"
                  opt7Label="Near Mint"
                />

                <DataDisplayRowRadio
                  label="Missing Parts Finding"
                  value={addComicSubmission.missingPartsFinding}
                  opt1Value="pr"
                  opt1Label="Poor"
                  opt2Value="fr"
                  opt2Label="Fair"
                  opt3Value="gd"
                  opt3Label="Good"
                  opt4Value="vg"
                  opt4Label="Very good"
                  opt5Value="fn"
                  opt5Label="Fine"
                  opt6Value="vf"
                  opt6Label="Very Fine"
                  opt7Value="nm"
                  opt7Label="Near Mint"
                />

                <DataDisplayRowRadio
                  label="Stains/Marks/Substances"
                  value={addComicSubmission.stainsFinding}
                  opt1Value="pr"
                  opt1Label="Poor"
                  opt2Value="fr"
                  opt2Label="Fair"
                  opt3Value="gd"
                  opt3Label="Good"
                  opt4Value="vg"
                  opt4Label="Very good"
                  opt5Value="fn"
                  opt5Label="Fine"
                  opt6Value="vf"
                  opt6Label="Very Fine"
                  opt7Value="nm"
                  opt7Label="Near Mint"
                />

                <DataDisplayRowRadio
                  label="Distortion/Colour"
                  value={addComicSubmission.distortionFinding}
                  opt1Value="pr"
                  opt1Label="Poor"
                  opt2Value="fr"
                  opt2Label="Fair"
                  opt3Value="gd"
                  opt3Label="Good"
                  opt4Value="vg"
                  opt4Label="Very good"
                  opt5Value="fn"
                  opt5Label="Fine"
                  opt6Value="vf"
                  opt6Label="Very Fine"
                  opt7Value="nm"
                  opt7Label="Near Mint"
                />

                <DataDisplayRowRadio
                  label="Paper Quality Finding"
                  value={addComicSubmission.paperQualityFinding}
                  opt1Value="pr"
                  opt1Label="Poor"
                  opt2Value="fr"
                  opt2Label="Fair"
                  opt3Value="gd"
                  opt3Label="Good"
                  opt4Value="vg"
                  opt4Label="Very good"
                  opt5Value="fn"
                  opt5Label="Fine"
                  opt6Value="vf"
                  opt6Label="Very Fine"
                  opt7Value="nm"
                  opt7Label="Near Mint"
                />

                <DataDisplayRowRadio
                  label="Spine/Staples"
                  value={addComicSubmission.spineFinding}
                  opt1Value="pr"
                  opt1Label="Poor"
                  opt2Value="fr"
                  opt2Label="Fair"
                  opt3Value="gd"
                  opt3Label="Good"
                  opt4Value="vg"
                  opt4Label="Very good"
                  opt5Value="fn"
                  opt5Label="Fine"
                  opt6Value="vf"
                  opt6Label="Very Fine"
                  opt7Value="nm"
                  opt7Label="Near Mint"
                />

                <DataDisplayRowRadio
                  label="Cover (front and back)"
                  value={addComicSubmission.coverFinding}
                  opt1Value="pr"
                  opt1Label="Poor"
                  opt2Value="fr"
                  opt2Label="Fair"
                  opt3Value="gd"
                  opt3Label="Good"
                  opt4Value="vg"
                  opt4Label="Very good"
                  opt5Value="fn"
                  opt5Label="Fine"
                  opt6Value="vf"
                  opt6Label="Very Fine"
                  opt7Value="nm"
                  opt7Label="Near Mint"
                />

                <DataDisplayRowRadio
                  label="Shows signs of tampering/restoration"
                  value={addComicSubmission.showsSignsOfTamperingOrRestoration}
                  opt1Value={2}
                  opt1Label="No"
                  opt2Value={1}
                  opt2Label="Yes"
                />

                {/* STEP 8 OF 10: NOTES ON GRADING */}
                <p className="title is-5 mt-2">
                  <FontAwesomeIcon className="fas" icon={faClipboardCheck} />
                  &nbsp;Notes on Grading&nbsp;&nbsp;&nbsp;-&nbsp;&nbsp;&nbsp;
                  <Link to="/admin/submissions/comics/add/step-8">
                    <FontAwesomeIcon className="fas" icon={faPencil} />
                    &nbsp;Edit
                  </Link>
                </p>

                <DataDisplayRowText
                  label="Grading Notes (Optional)"
                  value={addComicSubmission.gradingNotes}
                />

                {/* STEP 9 OF 10: GRADING */}
                <p className="title is-5 mt-2">
                  <FontAwesomeIcon className="fas" icon={faBalanceScale} />
                  &nbsp;Grading&nbsp;&nbsp;&nbsp;-&nbsp;&nbsp;&nbsp;
                  <Link to="/admin/submissions/comics/add/step-9">
                    <FontAwesomeIcon className="fas" icon={faPencil} />
                    &nbsp;Edit
                  </Link>
                </p>

                <DataDisplayRowRadio
                  label="Which type of grading scale would you prefer?"
                  value={addComicSubmission.gradingScale}
                  opt1Value={1}
                  opt1Label="Letter Grade (Poor-Near Mint)"
                  opt2Value={2}
                  opt2Label="Numbers (0.5-10.0)"
                  opt3Value={3}
                  opt3Label="Yes"
                />

                {addComicSubmission.gradingScale === 1 && (
                    <>
                        <DataDisplayRowSelect
                            label="Overall Letter Grade"
                            selectedValue={addComicSubmission.overallLetterGrade}
                            options={FINDING_WITH_EMPTY_OPTIONS}
                        />
                        {isNMPlusOpen && (
                            <DataDisplayRowCheckbox
                              label="Is Near Mint plus?"
                              checked={addComicSubmission.isOverallLetterGradeNearMintPlus}
                            />
                        )}
                    </>
                )}

                {addComicSubmission.gradingScale === 2 && (
                    <DataDisplayRowSelect
                        label="Overall Number Grade"
                        selectedValue={addComicSubmission.overallNumberGrade}
                        options={overallNumberGradeFilteredOptions}
                    />
                )}

                {addComicSubmission.gradingScale === 3 && (
                    <DataDisplayRowSelect
                        label="CPS Percentage Grade"
                        selectedValue={addComicSubmission.cpsPercentageGrade}
                        options={cpsPercentageGradeFilteredOptions}
                    />
                )}

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <button
                        class="button is-medium is-fullwidth-mobile"
                        onClick={(e) => {
                            e.preventDefault();
                            setForceURL("/admin/submissions/comics/add/step-9");
                        }}
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back to Step 9
                      </button>
                    </div>
                    <div class="column is-half has-text-right">
                      <button
                        class="button is-medium is-primary is-fullwidth-mobile"
                        onClick={onSubmitClick}
                      >
                        <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                        &nbsp;Submit
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

export default AdminComicSubmissionAddStep10;
