import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faShoppingCart,
  faTasks,
  faTachometer,
  faPlus,
  faEye,
  faArrowLeft,
  faCheckCircle,
  faPencil,
  faGauge,
  faBook,
  faBookOpen,
  faMagnifyingGlass,
  faBalanceScale,
  faUser,
  faArrowUpRightFromSquare,
  faDownload,
  faCogs,
} from "@fortawesome/free-solid-svg-icons";
import Select from "react-select";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import useLocalStorage from "../../../../Hooks/useLocalStorage";
import { getComicSubmissionDetailAPI } from "../../../../API/ComicSubmission";
import FormErrorBox from "../../../Reusable/FormErrorBox";
import FormComicSignaturesTable from "../../../Reusable/FormComicSignaturesTable";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import {
  FINDING_OPTIONS,
  OVERALL_NUMBER_GRADE_OPTIONS,
  PUBLISHER_NAME_OPTIONS,
  CPS_PERCENTAGE_GRADE_OPTIONS,
  HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS,
  ISSUE_COVER_YEAR_OPTIONS,
  ISSUE_COVER_MONTH_WITH_EMPTY_OPTIONS,
  SPECIAL_DETAILS_WITH_EMPTY_OPTIONS,
  SERVICE_TYPE_WITH_EMPTY_OPTIONS,
  SUBMISSION_KEY_ISSUE_WITH_EMPTY_OPTIONS,
  SUBMISSION_STATUS_WITH_EMPTY_OPTIONS,
  PAYMENT_PROCESSOR_WITH_EMPTY_OPTIONS,
} from "../../../../Constants/FieldOptions";
import { SERVICE_TYPE_CPS_CAPSULE_INDIE_MINT_GEM } from "../../../../Constants/App";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../AppState";
import FormRowText from "../../../Reusable/FormRowText";
import FormTextYesNoRow from "../../../Reusable/FormRowTextYesNo";
import DataDisplayRowSelect from "../../../Reusable/DataDisplayRowSelect";
import FormTextChoiceRow from "../../../Reusable/FormRowTextChoice";
import DataDisplayRowText from "../../../Reusable/DataDisplayRowText";
import DataDisplayRowURL from "../../../Reusable/DataDisplayRowURL";

function AdminComicSubmissionDetail() {
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
  const [showCustomerEditOptions, setShowCustomerEditOptions] = useState(false);

  ////
  //// Event handling.
  ////

  ////
  //// API.
  ////

  function onComicSubmissionDetailSuccess(response) {
    console.log("onComicSubmissionDetailSuccess: Starting...");
    setComicSubmission(response);
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

  // The following code will check to see if we need to grant the 'is NM+' option is available to the user.
  let isNMPlusOpen = false;
  if (submission !== undefined && submission !== null && submission !== "") {
    isNMPlusOpen =
      submission.gradingScale === 1 && submission.overallLetterGrade === "nm";
  }

  // Render the JSX content.
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
              to={`/admin/submissions/comic/${submission.id}/customer/search`}
              class="button is-medum is-menu is-primary"
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
                <Link to={`/admin/submissions/comics`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Comics
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page */}
          <nav class="box">
            {submission && (
              <div class="columns">
                <div class="column">
                  <p class="title is-4">
                    <FontAwesomeIcon className="fas" icon={faBookOpen} />
                    &nbsp;Online Comic Submission
                  </p>
                </div>
                <div class="column has-text-right">
                  <Link
                    to={`/admin/submissions/comic/${submission.id}/edit`}
                    class="button is-small is-warning is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faPencil} />
                    &nbsp;Edit
                  </Link>
                </div>
              </div>
            )}

            <FormErrorBox errors={errors} />

            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                {submission !== undefined &&
                  submission !== null &&
                  submission !== "" && (
                    <div class="container">
                      <div class="tabs is-medium is-size-7-mobile">
                        <ul>
                          <li class={`is-active`}>
                            <Link>
                              <b>Detail</b>
                            </Link>
                          </li>
                          <li>
                            <Link to={`/admin/submissions/comic/${id}/cust`}>
                              Customer
                            </Link>
                          </li>
                          <li>
                            <Link
                              to={`/admin/submissions/comic/${id}/comments`}
                            >
                              Comments
                            </Link>
                          </li>
                          <li>
                            <Link to={`/admin/submissions/comic/${id}/file`}>
                              File
                            </Link>
                          </li>
                          <li>
                            <Link
                              to={`/admin/submissions/comic/${id}/attachments`}
                            >
                              Attachments
                            </Link>
                          </li>
                        </ul>
                      </div>

                      <p class="subtitle is-6 pt-4">
                        <FontAwesomeIcon className="fas" icon={faBookOpen} />
                        &nbsp;Comic Book Information
                      </p>
                      <hr />

                      <FormRowText label="Store" value={submission.storeName} />

                      <DataDisplayRowSelect
                        label="Service Type"
                        selectedValue={submission.serviceType}
                        options={SERVICE_TYPE_WITH_EMPTY_OPTIONS}
                      />

                      <DataDisplayRowSelect
                        label="Primary Label Details"
                        selectedValue={submission.primaryLabelDetails}
                        helpText=""
                        options={SPECIAL_DETAILS_WITH_EMPTY_OPTIONS}
                      />

                      {submission.primaryLabelDetails === 1 && (
                        <FormRowText
                          label="Primary Label Details (Other)"
                          value={submission.primaryLabelDetailsOther}
                          helpText=""
                        />
                      )}

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

                      <DataDisplayRowSelect
                        label="Issue Cover Year"
                        selectedValue={submission.issueCoverYear}
                        helpText=""
                        options={ISSUE_COVER_YEAR_OPTIONS}
                      />

                      {submission.issueCoverYear !== 0 &&
                        submission.issueCoverYear !== 1 && (
                          <DataDisplayRowSelect
                            label="Issue Cover Month"
                            selectedValue={submission.issueCoverMonth}
                            helpText=""
                            options={ISSUE_COVER_MONTH_WITH_EMPTY_OPTIONS}
                          />
                        )}

                      <DataDisplayRowSelect
                        label="Publisher Name"
                        selectedValue={submission.publisherName}
                        helpText=""
                        options={PUBLISHER_NAME_OPTIONS}
                      />

                      {submission.publisherName === 1 && (
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

                      {submission.isKeyIssue && (
                        <>
                          <DataDisplayRowSelect
                            label="Key Issue"
                            selectedValue={submission.keyIssue}
                            helpText=""
                            options={SUBMISSION_KEY_ISSUE_WITH_EMPTY_OPTIONS}
                            disabled={true}
                          />
                          {submission.keyIssue === 1 && (
                            <>
                              <FormRowText
                                label="Key Issue Other"
                                value={submission.keyIssueOther}
                                helpText=""
                                helpText={"Max 638 characters"}
                                disabled={true}
                              />
                            </>
                          )}
                          <FormRowText
                            label="Key Issue Detail"
                            value={submission.keyIssueDetail}
                            helpText=""
                            disabled={true}
                          />
                        </>
                      )}

                      <FormTextYesNoRow
                        label="Is this an International Edition?"
                        checked={submission.isInternationalEdition}
                      />

                      <FormTextYesNoRow
                        label="Is variant cover?"
                        checked={submission.isVariantCover}
                      />

                      {submission.isVariantCover === true && (
                        <FormRowText
                          label="Variant cover detail"
                          value={submission.variantCoverDetail}
                        />
                      )}

                      <FormRowText
                        label="Special Notes (Optional)"
                        value={submission.specialNotes}
                      />

                      <FormComicSignaturesTable
                        data={submission.signatures}
                        disabled={true}
                      />

                      {submission.serviceType !==
                        SERVICE_TYPE_CPS_CAPSULE_INDIE_MINT_GEM && (
                        <>
                          <p class="subtitle is-6 pt-4">
                            <FontAwesomeIcon
                              className="fas"
                              icon={faMagnifyingGlass}
                            />
                            &nbsp;Summary of Findings
                          </p>
                          <hr />

                          <FormTextChoiceRow
                            label="Creases Finding"
                            value={submission.creasesFinding}
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

                          <FormTextChoiceRow
                            label="Tears Finding"
                            value={submission.tearsFinding}
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

                          <FormTextChoiceRow
                            label="Missing Parts Finding"
                            value={submission.missingPartsFinding}
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

                          <FormTextChoiceRow
                            label="Stains/Marks/Substances"
                            value={submission.stainsFinding}
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

                          <FormTextChoiceRow
                            label="Distortion/Colour"
                            value={submission.distortionFinding}
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

                          <FormTextChoiceRow
                            label="Paper Quality Finding"
                            value={submission.paperQualityFinding}
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

                          <FormTextChoiceRow
                            label="Spine/Staples"
                            name="spineFinding"
                            value={submission.spineFinding}
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
                            <FontAwesomeIcon
                              className="fas"
                              icon={faBalanceScale}
                            />
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

                          {submission.gradingScale === 1 && (
                            <>
                              <DataDisplayRowSelect
                                label="Overall Letter Grade"
                                selectedValue={submission.overallLetterGrade}
                                helpText=""
                                options={FINDING_OPTIONS}
                              />
                              {isNMPlusOpen && (
                                <>
                                  <FormTextChoiceRow
                                    label="Is Near Mint plus?"
                                    checked={
                                      submission.isOverallLetterGradeNearMintPlus
                                    }
                                  />
                                </>
                              )}
                            </>
                          )}

                          {submission.gradingScale === 2 && (
                            <DataDisplayRowSelect
                              label="Overall Number Grade"
                              selectedValue={submission.overallNumberGrade}
                              helpText=""
                              options={OVERALL_NUMBER_GRADE_OPTIONS}
                            />
                          )}

                          {submission.gradingScale === 3 && (
                            <DataDisplayRowSelect
                              label="CPS Percentage Grade"
                              selectedValue={submission.cpsPercentageGrade}
                              helpText=""
                              options={CPS_PERCENTAGE_GRADE_OPTIONS}
                            />
                          )}
                        </>
                      )}

                      {submission.paymentProcessor > 0 && (
                        <>
                          <p class="subtitle is-6">
                            <FontAwesomeIcon
                              className="fas"
                              icon={faShoppingCart}
                            />
                            &nbsp;Purchase Details
                          </p>
                          <hr />

                          <DataDisplayRowSelect
                            label="Payment Processor"
                            selectedValue={submission.paymentProcessor}
                            helpText=""
                            options={PAYMENT_PROCESSOR_WITH_EMPTY_OPTIONS}
                          />

                          <DataDisplayRowText
                            label="Purchase ID"
                            value={submission.paymentProcessorPurchaseId}
                            helpText="This is the purchase ID provided by the payment processor"
                          />

                          <DataDisplayRowText
                            label="Receipt ID"
                            value={submission.paymentProcessorReceiptId}
                            helpText="This is the receipt ID provided by the payment processor"
                          />

                          <DataDisplayRowURL
                            label="Receipt"
                            urlKey={`View Receipt`}
                            urlValue={submission.paymentProcessorReceiptUrl}
                            helpText="Click here to view your receipt"
                            type="external"
                          />

                          <DataDisplayRowText
                            label="Purchased At"
                            value={submission.paymentProcessorPurchasedAt}
                            helpText=""
                            type="date"
                          />

                          <DataDisplayRowText
                            label="Subtotal"
                            value={`$${submission.amountSubtotal}`}
                            helpText=""
                            type="text"
                          />

                          <DataDisplayRowText
                            label="Tax"
                            value={`$${submission.amountTax}`}
                            helpText=""
                            type="text"
                          />

                          <DataDisplayRowText
                            label="Total"
                            value={`$${submission.amountTotal}`}
                            helpText=""
                            type="text"
                          />
                        </>
                      )}

                      <p class="subtitle is-6">
                        <FontAwesomeIcon className="fas" icon={faCogs} />
                        &nbsp;Settings
                      </p>
                      <hr />

                      <FormRowText
                        label="CPSR #"
                        value={submission.cpsrn}
                        helpText="The unique identifier used by CPS for all submissions"
                      />

                      <DataDisplayRowSelect
                        label="Status"
                        selectedValue={submission.status}
                        helpText=""
                        options={SUBMISSION_STATUS_WITH_EMPTY_OPTIONS}
                      />

                      <div class="columns pt-4">
                        <div class="column is-half">
                          <Link
                            to={`/admin/submissions/comics`}
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
                            to={`/admin/submissions/comic/${id}/edit`}
                            class="button is-medium is-primary is-fullwidth-mobile"
                          >
                            <FontAwesomeIcon className="fas" icon={faPencil} />
                            &nbsp;Edit Comic Submission
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

export default AdminComicSubmissionDetail;
