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


function RetailerComicSubmissionAddStep7() {
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

    // Generate the payload.
    const submission = {
      seriesTitle: seriesTitle,
      issueVol: issueVol,
      issueNo: issueNo,
      issueCoverYear: issueCoverYear,
      issueCoverMonth: issueCoverMonth,
      publisherName: publisherName,
      publisherNameOther: publisherNameOther,
      isKeyIssue: isKeyIssue,
      keyIssue: keyIssue,
      keyIssueOther: keyIssueOther,
      keyIssueDetail: keyIssueDetail,
      isInternationalEdition: isInternationalEdition,
      isVariantCover: isVariantCover,
      variantCoverDetail: variantCoverDetail,
      printing: printing,
      primaryLabelDetails: primaryLabelDetails,
      primaryLabelDetailsOther: primaryLabelDetailsOther,
      specialNotes: specialNotes,
      gradingNotes: gradingNotes,
      signatures: signatures,
      creasesFinding: creasesFinding,
      tearsFinding: tearsFinding,
      missingPartsFinding: missingPartsFinding,
      stainsFinding: stainsFinding,
      distortionFinding: distortionFinding,
      paperQualityFinding: paperQualityFinding,
      spineFinding: spineFinding,
      coverFinding: coverFinding,
      gradingScale: parseInt(gradingScale), // 2=Number Grading Scale
      overallLetterGrade: overallLetterGrade,
      isOverallLetterGradeNearMintPlus: isOverallLetterGradeNearMintPlus,
      overallNumberGrade: parseFloat(overallNumberGrade),
      cpsPercentageGrade: parseFloat(cpsPercentageGrade),
      showsSignsOfTamperingOrRestoration: parseInt(
        showsSignsOfTamperingOrRestoration,
      ), // 2=No
      status: 1, //1=Waiting to Receive, 7=Completed by Retail Partner
      serviceType: serviceType,
      storeID: currentUser.storeId,
      collectibleType: 1, // 1=Comic, 2=Card
      customerID: customerID,
    };

    // Submit to the backend.
    console.log("onSubmitClick: payload:", submission);
    postComicSubmissionCreateAPI(
      submission,
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
            <p className="subtitle is-5">Step 7 of 10</p>
            <progress
              class="progress is-success"
              value="70"
              max="100"
            >
              70%
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

                  {serviceType !== SERVICE_TYPE_CPS_CAPSULE_INDIE_MINT_GEM && (
                    <>
                      <p class="subtitle is-6">
                        <FontAwesomeIcon
                          className="fas"
                          icon={faMagnifyingGlass}
                        />
                        &nbsp;Summary of Findings
                      </p>
                      <hr />

                      <FormRadioField
                        label="Creases Finding"
                        name="creasesFinding"
                        value={creasesFinding}
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
                        errorText={errors && errors.creasesFinding}
                        onChange={(e) => setCreasesFinding(e.target.value)}
                        maxWidth="180px"
                      />

                      <FormRadioField
                        label="Tears Finding"
                        name="tearsFinding"
                        value={tearsFinding}
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
                        errorText={errors && errors.tearsFinding}
                        onChange={(e) => setTearsFinding(e.target.value)}
                        maxWidth="180px"
                      />

                      <FormRadioField
                        label="Missing Parts Finding"
                        name="missingPartsFinding"
                        value={missingPartsFinding}
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
                        errorText={errors && errors.missingPartsFinding}
                        onChange={(e) => setMissingPartsFinding(e.target.value)}
                        maxWidth="180px"
                      />

                      <FormRadioField
                        label="Stains/Marks/Substances"
                        name="stainsFinding"
                        value={stainsFinding}
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
                        errorText={errors && errors.stainsFinding}
                        onChange={(e) => setStainsFinding(e.target.value)}
                        maxWidth="180px"
                      />

                      <FormRadioField
                        label="Distortion/Colour"
                        name="distortionFinding"
                        value={distortionFinding}
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
                        errorText={errors && errors.distortionFinding}
                        onChange={(e) => setDistortionFinding(e.target.value)}
                        maxWidth="180px"
                      />

                      <FormRadioField
                        label="Paper Quality Finding"
                        name="paperQualityFinding"
                        value={paperQualityFinding}
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
                        errorText={errors && errors.paperQualityFinding}
                        onChange={(e) => setPaperQualityFinding(e.target.value)}
                        maxWidth="180px"
                      />

                      <FormRadioField
                        label="Spine/Staples"
                        name="spineFinding"
                        value={spineFinding}
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
                        errorText={errors && errors.spineFinding}
                        onChange={(e) => setSpineFinding(e.target.value)}
                        maxWidth="180px"
                      />

                      <FormRadioField
                        label="Cover (front and back)"
                        name="coverFinding"
                        value={coverFinding}
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
                        errorText={errors && errors.coverFinding}
                        onChange={(e) => setCoverFinding(e.target.value)}
                        maxWidth="180px"
                      />

                      <FormRadioField
                        label="Shows signs of tampering/restoration"
                        name="showsSignsOfTamperingOrRestoration"
                        value={showsSignsOfTamperingOrRestoration}
                        opt1Value={2}
                        opt1Label="No"
                        opt2Value={1}
                        opt2Label="Yes"
                        errorText={
                          errors && errors.showsSignsOfTamperingOrRestoration
                        }
                        onChange={(e) =>
                          setShowsSignsOfTamperingOrRestoration(
                            parseInt(e.target.value),
                          )
                        }
                        maxWidth="180px"
                      />

                      <FormTextareaField
                        label="Grading Notes (Optional)"
                        name="gradingNotes"
                        placeholder="Text input"
                        value={gradingNotes}
                        errorText={errors && errors.gradingNotes}
                        helpText=""
                        onChange={(e) => setGradingNotes(e.target.value)}
                        isRequired={true}
                        maxWidth="280px"
                        helpText={"Max 638 characters"}
                        rows={4}
                      />

                      <p class="subtitle is-6">
                        <FontAwesomeIcon
                          className="fas"
                          icon={faBalanceScale}
                        />
                        &nbsp;Grading
                      </p>
                      <hr />

                      <FormRadioField
                        label="Which type of grading scale would you prefer?"
                        name="gradingScale"
                        value={gradingScale}
                        opt1Value={1}
                        opt1Label="Letter Grade (Poor-Near Mint)"
                        opt2Value={2}
                        opt2Label="Numbers (0.5-10.0)"
                        opt3Value={3}
                        opt3Label="CPS Percentage (5%-100%)"
                        errorText={errors && errors.gradingScale}
                        onChange={(e) =>
                          setGradingScale(parseInt(e.target.value))
                        }
                        maxWidth="180px"
                      />

                      {gradingScale === 1 && (
                        <>
                          <FormSelectField
                            label="Overall Letter Grade"
                            name="overallLetterGrade"
                            placeholder="Overall Letter Grade"
                            selectedValue={overallLetterGrade}
                            errorText={errors && errors.overallLetterGrade}
                            helpText=""
                            onChange={(e) =>
                              setOverallLetterGrade(e.target.value)
                            }
                            options={FINDING_WITH_EMPTY_OPTIONS}
                          />
                          {isNMPlusOpen && (
                            <>
                              <FormCheckboxField
                                label="Is Near Mint plus?"
                                name="isOverallLetterGradeNearMintPlus"
                                checked={isOverallLetterGradeNearMintPlus}
                                errorText={
                                  errors &&
                                  errors.isOverallLetterGradeNearMintPlus
                                }
                                onChange={(e) =>
                                  setIsOverallLetterGradeNearMintPlus(
                                    !isOverallLetterGradeNearMintPlus,
                                  )
                                }
                                maxWidth="180px"
                              />
                            </>
                          )}
                        </>
                      )}

                      {gradingScale === 2 && (
                        <FormSelectField
                          label="Overall Number Grade"
                          name="overallNumberGrade"
                          placeholder="Overall Number Grade"
                          selectedValue={overallNumberGrade}
                          errorText={errors && errors.overallNumberGrade}
                          helpText=""
                          onChange={(e) =>
                            setOverallNumberGrade(e.target.value)
                          }
                          options={overallNumberGradeFilteredOptions}
                        />
                      )}

                      {gradingScale === 3 && (
                        <FormSelectField
                          label="CPS Percentage Grade"
                          name="cpsPercentageGrade"
                          placeholder="CPS Percentage Grade"
                          selectedValue={cpsPercentageGrade}
                          errorText={errors && errors.cpsPercentageGrade}
                          helpText=""
                          onChange={(e) =>
                            setCpsPercentageGrade(e.target.value)
                          }
                          options={cpsPercentageGradeFilteredOptions}
                        />
                      )}
                    </>
                  )}

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <button
                        class="button is-medium is-fullwidth-mobile"
                        onClick={(e) => {
                            e.preventDefault();
                            setForceURL("/submissions/comics/add/step-6");
                        }}
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back to Step 6
                      </button>
                    </div>
                    <div class="column is-half has-text-right">
                      <button
                        class="button is-medium is-primary is-fullwidth-mobile"
                        onClick={onSubmitClick}
                      >
                        <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                        &nbsp;Submit and Continue
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

export default RetailerComicSubmissionAddStep7;
