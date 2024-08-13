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

function RetailerComicSubmissionAddStep3() {
  ////
  //// URL Parameters.
  ////

  const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140
  const customerID = searchParams.get("customer_id");
  const customerName = searchParams.get("customer_name");

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
  const [seriesTitle, setSeriesTitle] = useState("");
  const [issueVol, setIssueVol] = useState("");
  const [issueNo, setIssueNo] = useState("");
  const [issueCoverYear, setIssueCoverYear] = useState(0);
  const [issueCoverMonth, setIssueCoverMonth] = useState(0);
  const [publisherName, setPublisherName] = useState(0);
  const [publisherNameOther, setPublisherNameOther] = useState("");
  const [isKeyIssue, setIsKeyIssue] = useState(false);
  const [keyIssue, setKeyIssue] = useState(0);
  const [keyIssueOther, setKeyIssueOther] = useState("");
  const [keyIssueDetail, setKeyIssueDetail] = useState("");
  const [isInternationalEdition, setIsInternationalEdition] = useState(false);
  const [isVariantCover, setIsVariantCover] = useState(false);
  const [variantCoverDetail, setVariantCoverDetail] = useState("");
  const [printing, setPrinting] = useState(1);
  const [primaryLabelDetails, setPrimaryLabelDetails] = useState(2); // 2=Regular Edition
  const [primaryLabelDetailsOther, setPrimaryLabelDetailsOther] = useState("");
  const [creasesFinding, setCreasesFinding] = useState("");
  const [tearsFinding, setTearsFinding] = useState("");
  const [missingPartsFinding, setMissingPartsFinding] = useState("");
  const [stainsFinding, setStainsFinding] = useState("");
  const [distortionFinding, setDistortionFinding] = useState("");
  const [paperQualityFinding, setPaperQualityFinding] = useState("");
  const [spineFinding, setSpineFinding] = useState("");
  const [coverFinding, setCoverFinding] = useState("");
  const [gradingScale, setGradingScale] = useState(0);
  const [overallLetterGrade, setOverallLetterGrade] = useState("");
  const [overallNumberGrade, setOverallNumberGrade] = useState("");
  const [cpsPercentageGrade, setCpsPercentageGrade] = useState("");
  const [specialNotes, setSpecialNotes] = useState("");
  const [gradingNotes, setGradingNotes] = useState("");
  const [
    showsSignsOfTamperingOrRestoration,
    setShowsSignsOfTamperingOrRestoration,
  ] = useState(2); // 2=no
  const [showCancelWarning, setShowCancelWarning] = useState(false);
  const [
    isOverallLetterGradeNearMintPlus,
    setIsOverallLetterGradeNearMintPlus,
  ] = useState(false);
  const [serviceType, setServiceType] = useState(0);
  const [signatures, setSignatures] = useState([]);

  ////
  //// Event handling.
  ////

  const onSubmitClick = (e) => {
    console.log("onSubmitClick: Beginning...");
    console.log("onSubmitClick: Generating payload for submission.");
    setFetching(true);
    setErrors({});

    setForceURL("/submissions/comics/add/step-4");
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
    setTopAlertMessage("ComicSubmission created");
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

    let urlParams = "";
    if (customerName !== null) {
      urlParams +=
        "?customer_id=" + customerID + "&customer_name=" + customerName;
    }

    // Redirect the user to a new page.
    setForceURL("/submissions/comics/add/" + response.id + urlParams);
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
          {customerName === null ? (
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
                {customerName === null ? (
                  <Link
                    class="button is-medium is-success"
                    to={`/submissions/comics/add/step-1/search`}
                  >
                    Yes
                  </Link>
                ) : (
                  <Link
                    class="button is-medium is-success"
                    to={`/customer/${customerID}/sub`}
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
            <p className="subtitle is-5">Step 3 of 4</p>
            <progress
              class="progress is-success"
              value="75"
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
                    <FontAwesomeIcon className="fas" icon={faBookOpen} />
                    &nbsp;Book Information
                  </p>
                  <hr />

                  <FormInputField
                    label="Series Title"
                    name="seriesTitle"
                    placeholder="Text input"
                    value={seriesTitle}
                    errorText={errors && errors.seriesTitle}
                    helpText=""
                    onChange={(e) => setSeriesTitle(e.target.value)}
                    isRequired={true}
                    maxWidth="380px"
                  />

                  <FormInputField
                    label="Issue Vol"
                    name="issueVol"
                    placeholder="Text input"
                    value={issueVol}
                    errorText={errors && errors.issueVol}
                    helpText=""
                    onChange={(e) => setIssueVol(e.target.value)}
                    isRequired={true}
                    maxWidth="180px"
                  />

                  <FormInputField
                    label="Issue No"
                    name="issueNo"
                    placeholder="Text input"
                    value={issueNo}
                    errorText={errors && errors.issueNo}
                    helpText=""
                    onChange={(e) => setIssueNo(e.target.value)}
                    isRequired={true}
                    maxWidth="180px"
                  />

                  <FormSelectField
                    label="Issue Cover Year"
                    name="issueCoverYear"
                    placeholder="Issue Cover Year"
                    selectedValue={issueCoverYear}
                    errorText={errors && errors.issueCoverYear}
                    helpText=""
                    onChange={(e) =>
                      setIssueCoverYear(parseInt(e.target.value))
                    }
                    options={ISSUE_COVER_YEAR_OPTIONS}
                    isRequired={true}
                    maxWidth="200px"
                  />

                  {issueCoverYear !== 0 && (
                    <FormSelectField
                      label="Issue Cover Month"
                      name="issueCoverMonth"
                      placeholder="Issue Cover Month"
                      selectedValue={issueCoverMonth}
                      errorText={errors && errors.issueCoverMonth}
                      helpText=""
                      onChange={(e) =>
                        setIssueCoverMonth(parseInt(e.target.value))
                      }
                      options={ISSUE_COVER_MONTH_WITH_EMPTY_OPTIONS}
                      isRequired={true}
                      maxWidth="210px"
                    />
                  )}

                  <FormSelectField
                    label="Publisher Name"
                    name="publisherName"
                    placeholder="Publisher Name"
                    selectedValue={publisherName}
                    errorText={errors && errors.publisherName}
                    helpText=""
                    onChange={(e) => setPublisherName(parseInt(e.target.value))}
                    options={PUBLISHER_NAME_WITH_EMPTY_OPTIONS}
                  />

                  {publisherName === 1 && (
                    <FormInputField
                      label="Publisher Name (Other)"
                      name="publisherNameOther"
                      placeholder="Text input"
                      value={publisherNameOther}
                      errorText={errors && errors.publisherNameOther}
                      helpText=""
                      onChange={(e) => setPublisherNameOther(e.target.value)}
                      isRequired={true}
                      maxWidth="280px"
                    />
                  )}

                  <FormCheckboxField
                    label="Is Key Issue?"
                    name="isKeyIssue"
                    checked={isKeyIssue}
                    errorText={errors && errors.isKeyIssue}
                    onChange={(e) => setIsKeyIssue(!isKeyIssue)}
                    maxWidth="180px"
                  />

                  {isKeyIssue && (
                    <>
                      <FormSelectField
                        label="Key Issue"
                        name="keyIssue"
                        placeholder="Text input"
                        selectedValue={keyIssue}
                        errorText={errors && errors.keyIssue}
                        helpText=""
                        onChange={(e) => setKeyIssue(parseInt(e.target.value))}
                        options={SUBMISSION_KEY_ISSUE_WITH_EMPTY_OPTIONS}
                      />
                      {keyIssue === 1 && (
                        <>
                          <FormTextareaField
                            label="Key Issue Other"
                            name="keyIssueOther"
                            placeholder="Text input"
                            value={keyIssueOther}
                            errorText={errors && errors.keyIssueOther}
                            helpText=""
                            onChange={(e) => setKeyIssueOther(e.target.value)}
                            isRequired={true}
                            maxWidth="280px"
                            helpText={"Max 638 characters"}
                            rows={4}
                          />
                        </>
                      )}
                      {keyIssue !== 1 && (
                        <FormInputField
                          label="Key Issue Detail"
                          name="keyIssueDetail"
                          placeholder="Text input"
                          value={keyIssueDetail}
                          errorText={errors && errors.keyIssueDetail}
                          helpText=""
                          onChange={(e) => setKeyIssueDetail(e.target.value)}
                          isRequired={true}
                          maxWidth="280px"
                        />
                      )}
                    </>
                  )}

                  <FormCheckboxField
                    label="Is this an International Edition?"
                    name="isInternationalEdition"
                    checked={isInternationalEdition}
                    errorText={errors && errors.isInternationalEdition}
                    onChange={(e) =>
                      setIsInternationalEdition(!isInternationalEdition)
                    }
                    maxWidth="180px"
                  />

                  <FormCheckboxField
                    label="Is variant cover?"
                    name="isVariantCover"
                    checked={isVariantCover}
                    errorText={errors && errors.isVariantCover}
                    onChange={(e) => setIsVariantCover(!isVariantCover)}
                    maxWidth="180px"
                  />

                  {isVariantCover === true && (
                    <FormTextareaField
                      label="Variant cover detail"
                      name="variantCoverDetail"
                      placeholder="Text input"
                      value={variantCoverDetail}
                      errorText={errors && errors.variantCoverDetail}
                      helpText=""
                      onChange={(e) => setVariantCoverDetail(e.target.value)}
                      isRequired={true}
                      maxWidth="280px"
                      helpText={"Max 638 characters"}
                      rows={4}
                    />
                  )}

                  <FormSelectField
                    label="Which printing is this?"
                    name="printing"
                    placeholder="Text input"
                    selectedValue={printing}
                    errorText={errors && errors.printing}
                    helpText=""
                    onChange={(e) => setPrinting(parseInt(e.target.value))}
                    options={SUBMISSION_PRINTING_WITH_EMPTY_OPTIONS}
                  />

                  <FormSelectField
                    label="Primary Label Details"
                    name="primaryLabelDetails"
                    placeholder="Text input"
                    selectedValue={primaryLabelDetails}
                    errorText={errors && errors.primaryLabelDetails}
                    helpText=""
                    onChange={(e) =>
                      setPrimaryLabelDetails(parseInt(e.target.value))
                    }
                    options={SPECIAL_DETAILS_WITH_EMPTY_OPTIONS}
                  />

                  {primaryLabelDetails === 1 && (
                    <FormInputField
                      label="Primary Label Details (Other)"
                      name="primaryLabelDetailsOther"
                      placeholder="Text input"
                      value={primaryLabelDetailsOther}
                      errorText={errors && errors.primaryLabelDetailsOther}
                      helpText=""
                      onChange={(e) =>
                        setPrimaryLabelDetailsOther(e.target.value)
                      }
                      isRequired={true}
                      maxWidth="280px"
                    />
                  )}

                  <FormTextareaField
                    label="Special Note (Optional)"
                    name="specialNotes"
                    placeholder="Text input"
                    value={specialNotes}
                    errorText={errors && errors.specialNotesLine1}
                    helpText=""
                    onChange={(e) => setSpecialNotes(e.target.value)}
                    isRequired={true}
                    maxWidth="280px"
                    helpText={"Max 638 characters"}
                    rows={4}
                  />

                  {primaryLabelDetails === 1 && (
                    <FormInputField
                      label="Primary Label Details (Other)"
                      name="primaryLabelDetailsOther"
                      placeholder="Text input"
                      value={primaryLabelDetailsOther}
                      errorText={errors && errors.primaryLabelDetailsOther}
                      helpText=""
                      onChange={(e) =>
                        setPrimaryLabelDetailsOther(e.target.value)
                      }
                      isRequired={true}
                      maxWidth="280px"
                    />
                  )}

                  <FormComicSignaturesTable
                    data={signatures}
                    onDataChange={setSignatures}
                    disabled={true}
                    helpText={
                      <>
                        Only available to CPS,{" "}
                        <a href="mailto:support@cpscapsule.com">
                          contact us&nbsp;
                          <FontAwesomeIcon
                            className="mdi"
                            icon={faArrowUpRightFromSquare}
                          />
                        </a>{" "}
                        if you are interested in learning more.
                      </>
                    }
                  />

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
                        onClick={onSubmitClick}
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

export default RetailerComicSubmissionAddStep3;
