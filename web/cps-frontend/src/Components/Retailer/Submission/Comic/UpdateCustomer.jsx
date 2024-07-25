import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faBookOpen,
  faTachometer,
  faPlus,
  faArrowLeft,
  faCheckCircle,
  faPencil,
  faEye,
  faGauge,
} from "@fortawesome/free-solid-svg-icons";
import Select from "react-select";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import useLocalStorage from "../../../../Hooks/useLocalStorage";
import {
  getComicSubmissionDetailAPI,
  putComicSubmissionUpdateAPI,
} from "../../../../API/ComicSubmission";
import FormErrorBox from "../../../Reusable/FormErrorBox";
import FormInputField from "../../../Reusable/FormInputField";
import FormTextareaField from "../../../Reusable/FormTextareaField";
import FormRadioField from "../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../Reusable/FormSelectField";
import FormCheckboxField from "../../../Reusable/FormCheckboxField";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import {
  FINDING_OPTIONS,
  OVERALL_NUMBER_GRADE_OPTIONS,
  PUBLISHER_NAME_OPTIONS,
  CPS_PERCENTAGE_GRADE_OPTIONS,
} from "../../../../Constants/FieldOptions";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../AppState";

function RetailerComicSubmissionUpdateForCustomer() {
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
  const [seriesTitle, setSeriesTitle] = useState("");
  const [issueVol, setIssueVol] = useState("");
  const [issueNo, setIssueNo] = useState("");
  const [issueCoverDate, setIssueCoverDate] = useState("");
  const [publisherName, setPublisherName] = useState(0);
  const [publisherNameOther, setPublisherNameOther] = useState("");
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
  const [specialNotesLine1, setSpecialNotesLine1] = useState("");
  const [specialNotesLine2, setSpecialNotesLine2] = useState("");
  const [specialNotesLine3, setSpecialNotesLine3] = useState("");
  const [specialNotesLine4, setSpecialNotesLine4] = useState("");
  const [specialNotesLine5, setSpecialNotesLine5] = useState("");
  const [gradingNotesLine1, setGradingNotesLine1] = useState("");
  const [gradingNotesLine2, setGradingNotesLine2] = useState("");
  const [gradingNotesLine3, setGradingNotesLine3] = useState("");
  const [gradingNotesLine4, setGradingNotesLine4] = useState("");
  const [gradingNotesLine5, setGradingNotesLine5] = useState("");
  const [
    showsSignsOfTamperingOrRestoration,
    setShowsSignsOfTamperingOrRestoration,
  ] = useState("");
  const [
    isOverallLetterGradeNearMintPlus,
    setIsOverallLetterGradeNearMintPlus,
  ] = useState(false);

  ////
  //// Event handling.
  ////

  const onSeriesTitleChange = (e) => {
    setSeriesTitle(e.target.value);
  };

  const onIssueVolChange = (e) => {
    setIssueVol(e.target.value);
  };

  const onIssueNoChange = (e) => {
    setIssueNo(e.target.value);
  };

  const onIssueCoverDateChange = (e) => {
    setIssueCoverDate(e.target.value);
  };

  const onPublisherNameChange = (e) => {
    setPublisherName(parseInt(e.target.value));
  };

  const onPublisherNameOtherChange = (e) => {
    setPublisherNameOther(e.target.value);
  };

  const onCreasesFindingChange = (e) => {
    setCreasesFinding(e.target.value);
  };

  const onTearsFindingChange = (e) => {
    setTearsFinding(e.target.value);
  };

  const onMissingPartsFindingChange = (e) => {
    setMissingPartsFinding(e.target.value);
  };

  const onStainsFindingdChange = (e) => {
    setStainsFinding(e.target.value);
  };

  const onStainsFindingChange = (e) => {
    setStainsFinding(e.target.value);
  };

  const onDistortionFindingChange = (e) => {
    setDistortionFinding(e.target.value);
  };

  const onPaperQualityFindingChange = (e) => {
    setPaperQualityFinding(e.target.value);
  };

  const onSpineFindingChange = (e) => {
    setSpineFinding(e.target.value);
  };

  const onCoverFindingChange = (e) => {
    setCoverFinding(e.target.value);
  };

  const onGradingScaleChange = (e) => {
    setGradingScale(parseInt(e.target.value));
  };

  const onOverallLetterGradeChange = (e) => {
    setOverallLetterGrade(e.target.value);
  };

  const onOverallNumberGradeChange = (e) => {
    setOverallNumberGrade(e.target.value);
  };

  const onCpsPercentageGradeChange = (e) => {
    setCpsPercentageGrade(e.target.value);
  };

  const onSpecialNotesLine1Change = (e) => {
    setSpecialNotesLine1(e.target.value);
  };

  const onSpecialNotesLine2Change = (e) => {
    setSpecialNotesLine2(e.target.value);
  };

  const onSpecialNotesLine3Change = (e) => {
    setSpecialNotesLine3(e.target.value);
  };

  const onSpecialNotesLine4Change = (e) => {
    setSpecialNotesLine4(e.target.value);
  };

  const onSpecialNotesLine5Change = (e) => {
    setSpecialNotesLine5(e.target.value);
  };

  const onGradingNotesLine1Change = (e) => {
    setGradingNotesLine1(e.target.value);
  };

  const onGradingNotesLine2Change = (e) => {
    setGradingNotesLine2(e.target.value);
  };

  const onGradingNotesLine3Change = (e) => {
    setGradingNotesLine3(e.target.value);
  };

  const onGradingNotesLine4Change = (e) => {
    setGradingNotesLine4(e.target.value);
  };

  const onGradingNotesLine5Change = (e) => {
    setGradingNotesLine5(e.target.value);
  };

  const onShowsSignsOfTamperingOrRestorationChange = (e) => {
    setShowsSignsOfTamperingOrRestoration(e.target.value);
  };

  const onSubmitClick = (e) => {
    console.log("onSubmitClick: Beginning...");
    setErrors({});
    setFetching(true);

    // Generate the payload.
    const submission = {
      ID: id,
      SeriesTitle: seriesTitle,
      IssueVol: issueVol,
      IssueNo: issueNo,
      IssueCoverDate: issueCoverDate,
      PublisherName: publisherName,
      PublisherNameOther: publisherNameOther,
      SpecialNotesLine1: specialNotesLine1,
      SpecialNotesLine2: specialNotesLine2,
      SpecialNotesLine3: specialNotesLine3,
      SpecialNotesLine4: specialNotesLine4,
      SpecialNotesLine5: specialNotesLine5,
      GradingNotesLine1: gradingNotesLine1,
      GradingNotesLine2: gradingNotesLine2,
      GradingNotesLine3: gradingNotesLine3,
      GradingNotesLine4: gradingNotesLine4,
      GradingNotesLine5: gradingNotesLine5,
      CreasesFinding: creasesFinding,
      TearsFinding: tearsFinding,
      MissingPartsFinding: missingPartsFinding,
      StainsFinding: stainsFinding,
      DistortionFinding: distortionFinding,
      PaperQualityFinding: paperQualityFinding,
      SpineFinding: spineFinding,
      CoverFinding: coverFinding,
      GradingScale: parseInt(gradingScale),
      OverallLetterGrade: overallLetterGrade,
      IsOverallLetterGradeNearMintPlus: isOverallLetterGradeNearMintPlus,
      OverallNumberGrade: parseFloat(overallNumberGrade),
      CpsPercentageGrade: parseFloat(cpsPercentageGrade),
      ShowsSignsOfTamperingOrRestoration: parseInt(
        showsSignsOfTamperingOrRestoration,
      ),
    };

    // Submit to the backend.
    console.log("onSubmitClick, submission:", submission);
    putComicSubmissionUpdateAPI(
      submission,
      onComicSubmissionUpdateSuccess,
      onComicSubmissionUpdateError,
      onComicSubmissionUpdateDone,
      onUnauthorized,
    );
  };

  ////
  //// API.
  ////

  function onComicSubmissionDetailSuccess(response) {
    console.log("onComicSubmissionDetailSuccess: Starting...");
    setSeriesTitle(response.seriesTitle);
    setIssueVol(response.issueVol);
    setIssueNo(response.issueNo);
    setIssueCoverDate(response.issueCoverDate);
    setPublisherName(response.publisherName);
    setPublisherNameOther(response.publisherNameOther);
    setCreasesFinding(response.creasesFinding);
    setTearsFinding(response.tearsFinding);
    setMissingPartsFinding(response.missingPartsFinding);
    setStainsFinding(response.stainsFinding);
    setDistortionFinding(response.distortionFinding);
    setPaperQualityFinding(response.paperQualityFinding);
    setSpineFinding(response.spineFinding);
    setCoverFinding(response.coverFinding);
    setGradingScale(parseInt(response.gradingScale));
    setOverallLetterGrade(response.overallLetterGrade);
    setIsOverallLetterGradeNearMintPlus(
      response.isOverallLetterGradeNearMintPlus,
    );
    setOverallNumberGrade(response.overallNumberGrade);
    setSpecialNotesLine1(response.specialNotesLine1);
    setSpecialNotesLine2(response.specialNotesLine2);
    setSpecialNotesLine3(response.specialNotesLine3);
    setSpecialNotesLine4(response.specialNotesLine4);
    setSpecialNotesLine5(response.specialNotesLine5);
    setShowsSignsOfTamperingOrRestoration(
      response.showsSignsOfTamperingOrRestoration,
    );
    setGradingNotesLine1(response.gradingNotesLine1);
    setGradingNotesLine2(response.gradingNotesLine2);
    setGradingNotesLine3(response.gradingNotesLine3);
    setGradingNotesLine4(response.gradingNotesLine4);
    setGradingNotesLine5(response.gradingNotesLine5);
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

  function onComicSubmissionUpdateSuccess(response) {
    // For debugging purposes only.
    console.log("onComicSubmissionUpdateSuccess: Starting...");
    console.log(response);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("ComicSubmission created");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onComicSubmissionUpdateSuccess: Delayed for 2 seconds.");
      console.log(
        "onComicSubmissionUpdateSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Redirect the user to a new page.
    setForceURL("/submissions/comic/" + response.id);
  }

  function onComicSubmissionUpdateError(apiErr) {
    console.log("onComicSubmissionUpdateError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onComicSubmissionUpdateError: Delayed for 2 seconds.");
      console.log(
        "onComicSubmissionUpdateError: topAlertMessage, topAlertStatus:",
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

  function onComicSubmissionUpdateDone() {
    console.log("onComicSubmissionUpdateDone: Starting...");
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

  // The following code will check to see if we need to grant the 'is NM+' option is available to the user.
  let isNMPlusOpen = gradingScale === 1 && overallLetterGrade === "nm";

  // Render the JSX content.
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
              <li class="">
                <Link to={`/submissions/comic/${id}`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Details
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faPencil} />
                  &nbsp;Update (Customer)
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/submissions/comic/${id}`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Details
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

            <p class="has-text-grey pb-4">
              Please fill out all the required fields before submitting this
              form.
            </p>

            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                <div class="container">
                  <p class="subtitle is-6">Comic Book Information</p>
                  <hr />

                  <FormInputField
                    label="Series Title"
                    name="seriesTitle"
                    placeholder="Text input"
                    value={seriesTitle}
                    errorText={errors && errors.seriesTitle}
                    helpText=""
                    onChange={onSeriesTitleChange}
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
                    onChange={onIssueVolChange}
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
                    onChange={onIssueNoChange}
                    isRequired={true}
                    maxWidth="180px"
                  />

                  <FormInputField
                    label="Issue Cover Date"
                    name="issueCoverDate"
                    placeholder="Text input"
                    value={issueCoverDate}
                    errorText={errors && errors.issueCoverDate}
                    helpText=""
                    onChange={onIssueCoverDateChange}
                    isRequired={true}
                    maxWidth="180px"
                  />

                  <FormSelectField
                    label="Publisher Name"
                    name="publisherName"
                    placeholder="Publisher Name"
                    selectedValue={publisherName}
                    errorText={errors && errors.publisherName}
                    helpText=""
                    onChange={onPublisherNameChange}
                    options={PUBLISHER_NAME_OPTIONS}
                  />

                  {publisherName === 1 && (
                    <FormInputField
                      label="Publisher Name (Other)"
                      name="publisherNameOther"
                      placeholder="Text input"
                      value={publisherNameOther}
                      errorText={errors && errors.publisherNameOther}
                      helpText=""
                      onChange={onPublisherNameOtherChange}
                      isRequired={true}
                      maxWidth="280px"
                    />
                  )}

                  <FormInputField
                    label="Special Note - Line 1 (Optional)"
                    name="specialNotesLine1"
                    placeholder="Text input"
                    value={specialNotesLine1}
                    errorText={errors && errors.specialNotesLine1}
                    helpText=""
                    onChange={onSpecialNotesLine1Change}
                    isRequired={true}
                    maxWidth="280px"
                    helpText={"Max 17 characters"}
                  />

                  <FormInputField
                    label="Special Note - Line 2 (Optional)"
                    name="specialNotesLine2"
                    placeholder="Text input"
                    value={specialNotesLine2}
                    errorText={errors && errors.specialNotesLine2}
                    helpText=""
                    onChange={onSpecialNotesLine2Change}
                    isRequired={true}
                    maxWidth="280px"
                    helpText={"Max 17 characters"}
                  />

                  <FormInputField
                    label="Special Note - Line 3 (Optional)"
                    name="specialNotesLine3"
                    placeholder="Text input"
                    value={specialNotesLine3}
                    errorText={errors && errors.specialNotesLine3}
                    helpText=""
                    onChange={onSpecialNotesLine3Change}
                    isRequired={true}
                    maxWidth="280px"
                    helpText={"Max 17 characters"}
                  />

                  <FormInputField
                    label="Special Note - Line 4 (Optional)"
                    name="specialNotesLine4"
                    placeholder="Text input"
                    value={specialNotesLine4}
                    errorText={errors && errors.specialNotesLine4}
                    helpText=""
                    onChange={onSpecialNotesLine4Change}
                    isRequired={true}
                    maxWidth="280px"
                    helpText={"Max 17 characters"}
                  />

                  <FormInputField
                    label="Special Note - Line 5 (Optional)"
                    name="specialNotesLine5"
                    placeholder="Text input"
                    value={specialNotesLine5}
                    errorText={errors && errors.specialNotesLine5}
                    helpText=""
                    onChange={onSpecialNotesLine5Change}
                    isRequired={true}
                    maxWidth="280px"
                    helpText={"Max 17 characters"}
                  />

                  <p class="subtitle is-6">Summary of Findings</p>
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
                    onChange={onCreasesFindingChange}
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
                    onChange={onTearsFindingChange}
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
                    onChange={onMissingPartsFindingChange}
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
                    onChange={onStainsFindingChange}
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
                    onChange={onDistortionFindingChange}
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
                    onChange={onPaperQualityFindingChange}
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
                    onChange={onSpineFindingChange}
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
                    onChange={onCoverFindingChange}
                    maxWidth="180px"
                  />

                  <FormRadioField
                    label="Shows signs of tampering/restoration"
                    name="showsSignsOfTamperingOrRestoration"
                    value={showsSignsOfTamperingOrRestoration}
                    opt1Value={"2"}
                    opt1Label="No"
                    opt2Value={"1"}
                    opt2Label="Yes"
                    errorText={
                      errors && errors.showsSignsOfTamperingOrRestoration
                    }
                    onChange={onShowsSignsOfTamperingOrRestorationChange}
                    maxWidth="180px"
                  />

                  <FormInputField
                    label="Grading Note - Line 1 (Optional)"
                    name="gradingNotesLine1"
                    placeholder="Text input"
                    value={gradingNotesLine1}
                    errorText={errors && errors.gradingNotesLine1}
                    helpText=""
                    onChange={onGradingNotesLine1Change}
                    isRequired={true}
                    maxWidth="280px"
                    helpText={"Max 17 characters"}
                  />

                  <FormInputField
                    label="Grading Note - Line 2 (Optional)"
                    name="gradingNotesLine2"
                    placeholder="Text input"
                    value={gradingNotesLine2}
                    errorText={errors && errors.gradingNotesLine2}
                    helpText=""
                    onChange={onGradingNotesLine2Change}
                    isRequired={true}
                    maxWidth="280px"
                    helpText={"Max 17 characters"}
                  />

                  <FormInputField
                    label="Grading Note - Line 3 (Optional)"
                    name="gradingNotesLine3"
                    placeholder="Text input"
                    value={gradingNotesLine3}
                    errorText={errors && errors.gradingNotesLine3}
                    helpText=""
                    onChange={onGradingNotesLine3Change}
                    isRequired={true}
                    maxWidth="280px"
                    helpText={"Max 17 characters"}
                  />

                  <FormInputField
                    label="Grading Note - Line 4 (Optional)"
                    name="gradingNotesLine4"
                    placeholder="Text input"
                    value={gradingNotesLine4}
                    errorText={errors && errors.gradingNotesLine4}
                    helpText=""
                    onChange={onGradingNotesLine4Change}
                    isRequired={true}
                    maxWidth="280px"
                    helpText={"Max 17 characters"}
                  />

                  <FormInputField
                    label="Grading Note - Line 5 (Optional)"
                    name="gradingNotesLine5"
                    placeholder="Text input"
                    value={gradingNotesLine5}
                    errorText={errors && errors.gradingNotesLine5}
                    helpText=""
                    onChange={onGradingNotesLine5Change}
                    isRequired={true}
                    maxWidth="280px"
                    helpText={"Max 17 characters"}
                  />

                  <p class="subtitle is-6">Grading</p>
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
                    onChange={onGradingScaleChange}
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
                        onChange={(e) => setOverallLetterGrade(e.target.value)}
                        options={FINDING_OPTIONS}
                      />
                      {isNMPlusOpen && (
                        <>
                          <FormCheckboxField
                            label="Is Near Mint plus?"
                            name="isOverallLetterGradeNearMintPlus"
                            checked={isOverallLetterGradeNearMintPlus}
                            errorText={
                              errors && errors.isOverallLetterGradeNearMintPlus
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
                      placeholder="Overall Grade"
                      selectedValue={overallNumberGrade}
                      errorText={errors && errors.overallNumberGrade}
                      helpText=""
                      onChange={onOverallNumberGradeChange}
                      options={OVERALL_NUMBER_GRADE_OPTIONS}
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
                      onChange={onCpsPercentageGradeChange}
                      options={CPS_PERCENTAGE_GRADE_OPTIONS}
                    />
                  )}
                  <div class="columns pt-5">
                    <div class="column is-half">
                      <Link
                        to={`/submissions/comic/${id}`}
                        class="button is-medium is-fullwidth-mobile"
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back
                      </Link>
                    </div>
                    <div class="column is-half has-text-right">
                      <button
                        class="button is-primary is-medium is-fullwidth-mobile"
                        onClick={onSubmitClick}
                      >
                        <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                        &nbsp;Save
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

export default RetailerComicSubmissionUpdateForCustomer;
