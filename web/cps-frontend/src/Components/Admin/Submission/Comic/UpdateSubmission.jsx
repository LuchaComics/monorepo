import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faTachometer,
  faPlus,
  faArrowLeft,
  faCheckCircle,
  faPencil,
  faEye,
  faGauge,
  faBookOpen,
  faBook,
  faMagnifyingGlass,
  faBalanceScale,
  faCogs,
} from "@fortawesome/free-solid-svg-icons";
import Select from "react-select";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";
import { DateTime } from "luxon";

import useLocalStorage from "../../../../Hooks/useLocalStorage";
import {
  getComicSubmissionDetailAPI,
  putComicSubmissionUpdateAPI,
} from "../../../../API/ComicSubmission";
import { getStoreSelectOptionListAPI } from "../../../../API/store";
import FormErrorBox from "../../../Reusable/FormErrorBox";
import FormInputField from "../../../Reusable/FormInputField";
import FormTextareaField from "../../../Reusable/FormTextareaField";
import FormRadioField from "../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../Reusable/FormSelectField";
import FormDateField from "../../../Reusable/FormDateField";
import FormCheckboxField from "../../../Reusable/FormCheckboxField";
import FormComicSignaturesTable from "../../../Reusable/FormComicSignaturesTable";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import {
  FINDING_OPTIONS,
  OVERALL_NUMBER_GRADE_OPTIONS,
  PUBLISHER_NAME_OPTIONS,
  CPS_PERCENTAGE_GRADE_OPTIONS,
  ISSUE_COVER_YEAR_OPTIONS,
  ISSUE_COVER_MONTH_WITH_EMPTY_OPTIONS,
  SPECIAL_DETAILS_WITH_EMPTY_OPTIONS,
  SERVICE_TYPE_WITH_EMPTY_OPTIONS,
  SUBMISSION_KEY_ISSUE_WITH_EMPTY_OPTIONS,
  SUBMISSION_PRINTING_WITH_EMPTY_OPTIONS,
  SUBMISSION_STATUS_WITH_EMPTY_OPTIONS,
} from "../../../../Constants/FieldOptions";
import { SERVICE_TYPE_CPS_CAPSULE_INDIE_MINT_GEM } from "../../../../Constants/App";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../AppState";

function AdminComicSubmissionUpdateForComicSubmission() {
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
  const [variantCoverDetail, setVariantCoverDetail] = useState(false);
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
  ] = useState(2); // 2=No
  const [
    isOverallLetterGradeNearMintPlus,
    setIsOverallLetterGradeNearMintPlus,
  ] = useState(false);
  const [status, setStatus] = useState(0);
  const [storeSelectOptions, setStoreSelectOptions] = useState([]);
  const [storeID, setStoreID] = useState("");
  const [serviceType, setServiceType] = useState(0);
  const [signatures, setSignatures] = useState([]);

  ////
  //// Event handling.
  ////

  const onSubmitClick = (e) => {
    console.log("onSubmitClick: Beginning...");
    setFetching(true);
    setErrors({});

    // Generate the payload.
    const submission = {
      id: id,
      series_title: seriesTitle,
      issue_vol: issueVol,
      issue_no: issueNo,
      issue_cover_year: issueCoverYear,
      issue_cover_month: issueCoverMonth,
      publisher_name: publisherName,
      publisher_name_other: publisherNameOther,
      is_key_issue: isKeyIssue,
      key_issue: keyIssue,
      key_issue_other: keyIssueOther,
      key_issue_detail: keyIssueDetail,
      is_international_edition: isInternationalEdition,
      is_variant_cover: isVariantCover,
      variant_cover_detail: variantCoverDetail,
      variantCoverDetail: variantCoverDetail,
      printing: printing,
      primary_label_details: primaryLabelDetails,
      primary_label_details_other: primaryLabelDetailsOther,
      special_notes: specialNotes,
      grading_notes: gradingNotes,
      creases_finding: creasesFinding,
      tears_finding: tearsFinding,
      missing_parts_finding: missingPartsFinding,
      stains_finding: stainsFinding,
      distortion_finding: distortionFinding,
      paper_quality_finding: paperQualityFinding,
      spine_finding: spineFinding,
      cover_finding: coverFinding,
      grading_scale: parseInt(gradingScale),
      overall_letter_grade: overallLetterGrade,
      is_overall_letter_grade_near_mint_plus: isOverallLetterGradeNearMintPlus,
      overall_number_grade: parseFloat(overallNumberGrade),
      cps_percentage_grade: parseFloat(cpsPercentageGrade),
      shows_signs_of_tampering_or_restoration: parseInt(
        showsSignsOfTamperingOrRestoration,
      ),
      status: status,
      service_type: serviceType,
      store_id: storeID,
      signatures: signatures,
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
    setIssueCoverYear(response.issueCoverYear);
    setIssueCoverMonth(response.issueCoverMonth);
    setPublisherName(response.publisherName);
    setPublisherNameOther(response.publisherNameOther);
    setIsKeyIssue(response.isKeyIssue);
    setKeyIssue(response.keyIssue);
    setKeyIssueDetail(response.keyIssueDetail);
    setKeyIssueOther(response.keyIssueOther);
    setIsInternationalEdition(response.isInternationalEdition);
    setIsVariantCover(response.isVariantCover);
    setVariantCoverDetail(response.variantCoverDetail);
    setPrimaryLabelDetails(response.primaryLabelDetails);
    setPrinting(response.printing);
    setPrimaryLabelDetailsOther(response.primaryLabelDetailsOther);
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
    setCpsPercentageGrade(response.cpsPercentageGrade);
    setSpecialNotes(response.specialNotes);
    setShowsSignsOfTamperingOrRestoration(
      response.showsSignsOfTamperingOrRestoration,
    );
    setGradingNotes(response.gradingNotes);
    setStatus(response.status);
    setServiceType(response.serviceType);
    setStoreID(response.storeId);
    setSignatures(response.signatures);
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
    setForceURL("/admin/submissions/comic/" + response.id);
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

  function onStoreOptionListSuccess(response) {
    console.log("onStoreOptionListSuccess: Starting...");
    if (response !== null) {
      const selectOptions = [
        { value: 0, label: "Please select" }, // Add empty options.
        ...response,
      ];
      setStoreSelectOptions(selectOptions);
    }
  }

  function onStoreOptionListError(apiErr) {
    console.log("onStoreOptionListError: Starting...");
    console.log("onStoreOptionListError: apiErr:", apiErr);
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onStoreOptionListDone() {
    console.log("onStoreOptionListDone: Starting...");
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

      let params = new Map();
      getStoreSelectOptionListAPI(
        params,
        onStoreOptionListSuccess,
        onStoreOptionListError,
        onStoreOptionListDone,
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
              <li class="">
                <Link to={`/admin/submissions/comic/${id}`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faPencil} />
                  &nbsp;Update (Online Comic Submission)
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/admin/submissions/comic/${id}`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Detail
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faTasks} />
              &nbsp;Online Comic Submission
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
                    <FontAwesomeIcon className="fas" icon={faCogs} />
                    &nbsp;Settings
                  </p>
                  <hr />

                  <FormSelectField
                    label="Store"
                    name="storeID"
                    placeholder="Pick"
                    selectedValue={storeID}
                    errorText={errors && errors.storeID}
                    helpText="Pick the store this user belongs to and will be limited by"
                    isRequired={true}
                    onChange={(e) => setStoreID(e.target.value)}
                    options={storeSelectOptions}
                    disabled={storeSelectOptions.length === 0}
                  />

                  <FormSelectField
                    label="Service Type"
                    name="serviceType"
                    selectedValue={serviceType}
                    errorText={errors && errors.serviceType}
                    onChange={(e) => {
                      setServiceType(parseInt(e.target.value));
                    }}
                    options={SERVICE_TYPE_WITH_EMPTY_OPTIONS}
                    maxWidth="400px"
                  />

                  <FormSelectField
                    label="Status"
                    name="status"
                    selectedValue={status}
                    errorText={errors && errors.status}
                    onChange={(e) => setStatus(parseInt(e.target.value))}
                    options={SUBMISSION_STATUS_WITH_EMPTY_OPTIONS}
                    maxWidth="400px"
                    disabledValues={[]}
                  />

                  <p class="subtitle is-6">
                    <FontAwesomeIcon className="fas" icon={faBook} />
                    &nbsp;Comic Book Information
                  </p>
                  <hr />

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
                    label="Issue Volume"
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
                    label="Issue No."
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
                    maxWidth="180px"
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

                  <FormTextareaField
                    label="Special Notes (Optional)"
                    name="specialNotes"
                    placeholder="Text input"
                    value={specialNotes}
                    errorText={errors && errors.specialNotes}
                    helpText=""
                    onChange={(e) => setSpecialNotes(e.target.value)}
                    isRequired={true}
                    maxWidth="280px"
                    helpText={"Max 638 characters"}
                    rows={4}
                  />

                  <FormComicSignaturesTable
                    data={signatures}
                    onDataChange={setSignatures}
                  />

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
                        value={parseInt(showsSignsOfTamperingOrRestoration)}
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
                            options={FINDING_OPTIONS}
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
                          placeholder="Overall Grade"
                          selectedValue={overallNumberGrade}
                          errorText={errors && errors.overallNumberGrade}
                          helpText=""
                          onChange={(e) =>
                            setOverallNumberGrade(e.target.value)
                          }
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
                          onChange={(e) =>
                            setCpsPercentageGrade(e.target.value)
                          }
                          options={CPS_PERCENTAGE_GRADE_OPTIONS}
                        />
                      )}
                    </>
                  )}

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <Link
                        to={`/admin/submissions/comic/${id}`}
                        class="button is-medium is-fullwidth-mobile"
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back
                      </Link>
                    </div>
                    <div class="column is-half has-text-right">
                      <button
                        class="button is-medium is-primary is-fullwidth-mobile"
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

export default AdminComicSubmissionUpdateForComicSubmission;
