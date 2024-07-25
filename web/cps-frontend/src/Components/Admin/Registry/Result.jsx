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
  faBook,
  faMagnifyingGlass,
  faBalanceScale,
  faUser,
  faArrowUpRightFromSquare,
  faDownload,
  faBarcode,
} from "@fortawesome/free-solid-svg-icons";
import Select from "react-select";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import useLocalStorage from "../../../Hooks/useLocalStorage";
import { getRegistryAPI } from "../../../API/registry";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormInputField from "../../Reusable/FormInputField";
import FormTextareaField from "../../Reusable/FormTextareaField";
import FormRadioField from "../../Reusable/FormRadioField";
import FormMultiSelectField from "../../Reusable/FormMultiSelectField";
import FormCheckboxField from "../../Reusable/FormCheckboxField";
import FormSelectField from "../../Reusable/FormSelectField";
import FormDateField from "../../Reusable/FormDateField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import {
  FINDING_OPTIONS,
  OVERALL_NUMBER_GRADE_OPTIONS,
  PUBLISHER_NAME_OPTIONS,
  CPS_PERCENTAGE_GRADE_OPTIONS,
  HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS,
  ISSUE_COVER_YEAR_OPTIONS,
  ISSUE_COVER_MONTH_WITH_EMPTY_OPTIONS,
} from "../../../Constants/FieldOptions";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";
import FormRowText from "../../Reusable/FormRowText";
import FormTextYesNoRow from "../../Reusable/FormRowTextYesNo";
import FormTextOptionRow from "../../Reusable/FormRowTextOption";
import FormTextChoiceRow from "../../Reusable/FormRowTextChoice";

function AdminRegistryResult() {
  ////
  //// URL Parameters.
  ////

  const { cpsn } = useParams();

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
  const [submission, setSubmission] = useState({});
  const [showCustomerEditOptions, setShowCustomerEditOptions] = useState(false);

  ////
  //// Event handling.
  ////

  ////
  //// API.
  ////

  function onSubmissionDetailSuccess(response) {
    console.log("onSubmissionDetailSuccess: Starting...");
    setSubmission(response);
  }

  function onSubmissionDetailError(apiErr) {
    console.log("onSubmissionDetailError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onSubmissionDetailDone() {
    console.log("onSubmissionDetailDone: Starting...");
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
      getRegistryAPI(
        cpsn,
        onSubmissionDetailSuccess,
        onSubmissionDetailError,
        onSubmissionDetailDone,
        onUnauthorized,
      );
    }

    return () => {
      mounted = false;
    };
  }, [cpsn]);

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
                <Link to="/admin/registry" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faBarcode} />
                  &nbsp;Registry
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/admin/registry`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Registry
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faTasks} />
              &nbsp;Submission
            </p>
            <FormErrorBox errors={errors} />

            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                {submission && (
                  <div class="container">
                    <p class="subtitle is-6 pt-4">
                      <FontAwesomeIcon className="fas" icon={faBook} />
                      &nbsp;Comic Book Information
                    </p>
                    <hr />

                    <FormRowText
                      label="Series Title"
                      value={submission.seriesTitle}
                      helpText=""
                    />

                    <FormRowText
                      label="Issue Vol"
                      value={submission.issueVol}
                      helpText=""
                    />

                    <FormRowText
                      label="Issue No"
                      value={submission.issueNo}
                      helpText=""
                    />

                    <FormTextOptionRow
                      label="Issue Cover Year"
                      selectedValue={submission.issueCoverYear}
                      helpText=""
                      options={ISSUE_COVER_YEAR_OPTIONS}
                    />

                    {submission.issueCoverYear !== 0 &&
                      submission.issueCoverYear !== 1 && (
                        <FormTextOptionRow
                          label="Issue Cover Month"
                          selectedValue={submission.issueCoverMonth}
                          helpText=""
                          options={ISSUE_COVER_MONTH_WITH_EMPTY_OPTIONS}
                        />
                      )}

                    <FormTextOptionRow
                      label="Publisher Name"
                      selectedValue={submission.publisherName}
                      helpText=""
                      options={PUBLISHER_NAME_OPTIONS}
                    />

                    {submission.publisherName === "Other" && (
                      <FormRowText
                        label="Publisher Name (Other)"
                        value={submission.publisherNameOther}
                        helpText=""
                        disabled={true}
                      />
                    )}

                    <FormTextYesNoRow
                      label="Is Key Issue?"
                      checked={submission.isKeyIssue}
                    />

                    {submission.primaryLabelDetails === 1 && (
                      <FormRowText
                        label="Primary Label Details (Other)"
                        value={submission.primaryLabelDetailsOther}
                        helpText=""
                      />
                    )}

                    <FormRowText
                      label="Special Notes (Optional)"
                      value={submission.specialNotes}
                    />

                    <p class="subtitle is-6 pt-4">
                      <FontAwesomeIcon
                        className="fas"
                        icon={faMagnifyingGlass}
                      />
                      &nbsp;Summary of Findings
                    </p>
                    <hr />

                    <FormTextChoiceRow
                      label="Shows signs of tampering/restoration"
                      value={parseInt(
                        submission.showsSignsOfTamperingOrRestoration,
                      )}
                      opt1Value={2}
                      opt1Label="No"
                      opt2Value={1}
                      opt2Label="Yes"
                    />

                    <FormRowText
                      label="Grading Notes (Optional)"
                      value={submission.gradingNotes}
                    />

                    <p class="subtitle is-6 pt-4">
                      <FontAwesomeIcon className="fas" icon={faBalanceScale} />
                      &nbsp;Grading
                    </p>
                    <hr />

                    <FormTextChoiceRow
                      label="Which type of grading scale would you prefer?"
                      value={parseInt(submission.gradingScale)}
                      opt1Value={1}
                      opt1Label="Letter Grade (Poor-Near Mint)"
                      opt2Value={2}
                      opt2Label="Numbers (0.5-10.0)"
                      opt3Value={3}
                      opt3Label="CPS Percentage (5%-100%)"
                    />

                    {submission && submission.gradingScale === 1 && (
                      <FormTextOptionRow
                        label="Overall Letter Grade"
                        selectedValue={submission.overallLetterGrade}
                        helpText=""
                        options={FINDING_OPTIONS}
                      />
                    )}

                    {submission && submission.gradingScale === 2 && (
                      <FormTextOptionRow
                        label="Overall Number Grade"
                        selectedValue={submission.overallNumberGrade}
                        helpText=""
                        options={OVERALL_NUMBER_GRADE_OPTIONS}
                      />
                    )}

                    {submission && submission.gradingScale === 3 && (
                      <FormTextOptionRow
                        label="CPS Percentage Grade"
                        selectedValue={submission.cpsPercentageGrade}
                        helpText=""
                        options={CPS_PERCENTAGE_GRADE_OPTIONS}
                      />
                    )}

                    <div class="columns pt-4">
                      <div class="column is-half">
                        <Link
                          to={`/admin/registry`}
                          class="button is-medium is-fullwidth-mobile"
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back to Search
                        </Link>
                      </div>
                      <div class="column is-half has-text-right"></div>
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

export default AdminRegistryResult;
