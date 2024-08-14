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
  const [addComicSubmission, setAddComicSubmission] = useRecoilState(addComicSubmissionState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false); // Bool
  const [forceURL, setForceURL] = useState("");
  const [seriesTitle, setSeriesTitle] = useState(addComicSubmission.seriesTitle);
  const [issueVol, setIssueVol] = useState(addComicSubmission.issueVol);
  const [issueNo, setIssueNo] = useState(addComicSubmission.issueNo);
  const [issueCoverYear, setIssueCoverYear] = useState(parseInt(addComicSubmission.issueCoverYear));
  const [issueCoverMonth, setIssueCoverMonth] = useState(parseInt(addComicSubmission.issueCoverMonth));
  const [publisherName, setPublisherName] = useState(parseInt(addComicSubmission.publisherName));
  const [publisherNameOther, setPublisherNameOther] = useState(addComicSubmission.publisherNameOther);
  const [isKeyIssue, setIsKeyIssue] = useState(false); // Bool
  const [keyIssue, setKeyIssue] = useState(parseInt(addComicSubmission.keyIssue));
  const [keyIssueOther, setKeyIssueOther] = useState(addComicSubmission.serviceType);
  const [keyIssueDetail, setKeyIssueDetail] = useState(addComicSubmission.serviceType);
  const [isInternationalEdition, setIsInternationalEdition] = useState(false); // Bool
  const [isVariantCover, setIsVariantCover] = useState(false); // Bool
  const [variantCoverDetail, setVariantCoverDetail] = useState(addComicSubmission.serviceType);
  const [printing, setPrinting] = useState(parseInt(addComicSubmission.printing));
  const [primaryLabelDetails, setPrimaryLabelDetails] = useState(parseInt(addComicSubmission.primaryLabelDetails)); // 2=Regular Edition
  const [primaryLabelDetailsOther, setPrimaryLabelDetailsOther] = useState(addComicSubmission.primaryLabelDetailsOther);
  const [creasesFinding, setCreasesFinding] = useState(addComicSubmission.serviceType);
  const [tearsFinding, setTearsFinding] = useState(addComicSubmission.serviceType);
  const [missingPartsFinding, setMissingPartsFinding] = useState(addComicSubmission.serviceType);
  const [stainsFinding, setStainsFinding] = useState(addComicSubmission.serviceType);
  const [distortionFinding, setDistortionFinding] = useState(addComicSubmission.serviceType);
  const [paperQualityFinding, setPaperQualityFinding] = useState(addComicSubmission.serviceType);
  const [specialNotes, setSpecialNotes] = useState(addComicSubmission.specialNotes);
  const [
    showsSignsOfTamperingOrRestoration,
    setShowsSignsOfTamperingOrRestoration,
  ] = useState(2); // 2=no  // Bool
  const [showCancelWarning, setShowCancelWarning] = useState(false); // Bool
  const [
    isOverallLetterGradeNearMintPlus,
    setIsOverallLetterGradeNearMintPlus,
  ] = useState(false); // Bool
  const [serviceType, setServiceType] = useState(parseInt(addComicSubmission.serviceType));
  const [signatures, setSignatures] = useState([]);

  ////
  //// Event handling.
  ////

  const onSaveAndContinueClick = (e) => {
      console.log("onSaveAndContinueClick: Beginning...");

      // Variables used to hold state if we got an error with validation.
      let newErrors = {};
      let hasErrors = false;

      // Perform validation.
      if (seriesTitle === undefined || seriesTitle === null || seriesTitle === 0 || seriesTitle === "") {
        newErrors["seriesTitle"] = "missing value";
        hasErrors = true;
      }
      if (issueVol === undefined || issueVol === null || issueVol === "") {
        newErrors["issueVol"] = "missing value";
        hasErrors = true;
      }
      if (issueNo === undefined || issueNo === null || issueNo === "") {
        newErrors["issueNo"] = "missing value";
        hasErrors = true;
      }
      if (issueCoverYear === undefined || issueCoverYear === null || issueCoverYear === 0 || issueCoverYear === "") {
        newErrors["issueCoverYear"] = "missing value";
        hasErrors = true;
      }
      if (issueCoverMonth === undefined || issueCoverMonth === null || issueCoverMonth === 0 || issueCoverMonth === "") {
        newErrors["issueCoverMonth"] = "missing value";
        hasErrors = true;
      }
      if (publisherName === undefined || publisherName === null || publisherName === 0 || publisherName === "") {
        newErrors["publisherName"] = "missing value";
        hasErrors = true;
      } else if (publisherName === 1) { // Is other.
          if (publisherNameOther === undefined || publisherNameOther === null || publisherNameOther === "") {
            newErrors["publisherNameOther"] = "missing value";
            hasErrors = true;
          }
      }
      if (printing === undefined || printing === null || printing === 0 || printing === "") {
        newErrors["printing"] = "missing value";
        hasErrors = true;
      }
      if (issueVol === undefined || issueVol === null || issueVol === "") {
        newErrors["issueVol"] = "missing value";
        hasErrors = true;
      }
      if (primaryLabelDetails === undefined || primaryLabelDetails === null || primaryLabelDetails === 0 || primaryLabelDetails === "") {
        newErrors["primaryLabelDetails"] = "missing value";
        hasErrors = true;
      } else if (primaryLabelDetails === 1) {
          if (primaryLabelDetailsOther === undefined || primaryLabelDetailsOther === null || primaryLabelDetailsOther === 0 || primaryLabelDetailsOther === "") {
            newErrors["primaryLabelDetailsOther"] = "missing value";
            hasErrors = true;
          }
      }

      //
      // CASE 1 of 2: Has errors.
      //

      if (hasErrors) {
        console.log("onSaveAndContinueClick: Aboring because of error(s)");

        // Set the associate based error validation.
        setErrors(newErrors);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();

        return;
      }

      //
      // CASE 2 of 2: Has no errors.
      //

      console.log("onSaveAndContinueClick: Saving step 3 and redirecting to step 4.");

      // Variable holds a complete clone of the submission.
      let modifiedAddComicSubmission = { ...addComicSubmission };

      // Update our clone.
      modifiedAddComicSubmission.seriesTitle = seriesTitle;
      modifiedAddComicSubmission.issueVol = issueVol;
      modifiedAddComicSubmission.issueNo = issueNo;
      modifiedAddComicSubmission.issueCoverYear = issueCoverYear;
      modifiedAddComicSubmission.issueCoverMonth = issueCoverMonth;
      modifiedAddComicSubmission.publisherName = parseInt(publisherName);
      modifiedAddComicSubmission.publisherNameOther = publisherNameOther;
      modifiedAddComicSubmission.printing = parseInt(printing);
      modifiedAddComicSubmission.issueVol = issueVol;
      modifiedAddComicSubmission.primaryLabelDetails = primaryLabelDetails;
      modifiedAddComicSubmission.primaryLabelDetailsOther = primaryLabelDetailsOther

      // Save to persistent storage.
      setAddComicSubmission(modifiedAddComicSubmission);

      // Redirect to the next page.
      setForceURL("/submissions/comics/add/step-4")
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
          {/* ------ */}

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
                        onClick={(e) => {
                            e.preventDefault();
                            setForceURL("/submissions/comics/add/step-2")
                        }}
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back to Step 2
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

export default RetailerComicSubmissionAddStep3;
