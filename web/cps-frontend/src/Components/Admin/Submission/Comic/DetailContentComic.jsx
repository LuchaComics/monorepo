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
  faCogs,
} from "@fortawesome/free-solid-svg-icons";
import Select from "react-select";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import useLocalStorage from "../../../../Hooks/useLocalStorage";
import FormErrorBox from "../../../Reusable/FormErrorBox";
import FormInputField from "../../../Reusable/FormInputField";
import FormTextareaField from "../../../Reusable/FormTextareaField";
import FormRadioField from "../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../Reusable/FormMultiSelectField";
import FormCheckboxField from "../../../Reusable/FormCheckboxField";
import FormSelectField from "../../../Reusable/FormSelectField";
import FormDateField from "../../../Reusable/FormDateField";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import {
  FINDING_OPTIONS,
  OVERALL_NUMBER_GRADE_OPTIONS,
  PUBLISHER_NAME_OPTIONS,
  CPS_PERCENTAGE_GRADE_OPTIONS,
  HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS,
  ISSUE_COVER_YEAR_OPTIONS,
  ISSUE_COVER_MONTH_WITH_EMPTY_OPTIONS,
} from "../../../../Constants/FieldOptions";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../AppState";

function AdminComicSubmissionDetailContentComic({
  id,
  errors,
  isFetching,
  submission,
  showCustomerEditOptions,
  setShowCustomerEditOptions,
}) {
  // The following code will check to see if we need to grant the 'is NM+' option is available to the user.
  let isNMPlusOpen = false;
  if (submission !== undefined && submission !== null && submission !== "") {
    isNMPlusOpen =
      submission.gradingScale === 1 && submission.overallLetterGrade === "nm";
  }

  // Render the JSX content.
  return (
    <>
      {/* Modals */}
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
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faTasks} />
              &nbsp;Online Comic Submission
            </p>
            <FormErrorBox errors={errors} />

            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                {submission && (
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
                          <Link to={`/admin/submissions/comic/${id}/comments`}>
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
                      <FontAwesomeIcon className="fas" icon={faBook} />
                      &nbsp;Comic Book Information
                    </p>
                    <hr />

                    {submission && (
                      <FormInputField
                        label="Series Title"
                        name="seriesTitle"
                        placeholder="Text input"
                        value={submission.seriesTitle}
                        helpText=""
                        isRequired={true}
                        maxWidth="380px"
                        disabled={true}
                      />
                    )}

                    {submission && (
                      <FormInputField
                        label="Issue Vol"
                        name="issueVol"
                        placeholder="Text input"
                        value={submission.issueVol}
                        helpText=""
                        isRequired={true}
                        maxWidth="180px"
                        disabled={true}
                      />
                    )}

                    {submission && (
                      <FormInputField
                        label="Issue No"
                        name="issueNo"
                        placeholder="Text input"
                        value={submission.issueNo}
                        helpText=""
                        isRequired={true}
                        maxWidth="180px"
                        disabled={true}
                      />
                    )}

                    <FormSelectField
                      label="Issue Cover Year"
                      name="issueCoverYear"
                      placeholder="Issue Cover Year"
                      selectedValue={submission.issueCoverYear}
                      helpText=""
                      options={ISSUE_COVER_YEAR_OPTIONS}
                      isRequired={true}
                      maxWidth="110px"
                      disabled={true}
                    />

                    {submission.issueCoverYear !== 0 &&
                      submission.issueCoverYear !== 1 && (
                        <FormSelectField
                          label="Issue Cover Month"
                          name="issueCoverMonth"
                          placeholder="Issue Cover Month"
                          selectedValue={submission.issueCoverMonth}
                          helpText=""
                          options={ISSUE_COVER_MONTH_WITH_EMPTY_OPTIONS}
                          isRequired={true}
                          maxWidth="110px"
                          disabled={true}
                        />
                      )}

                    <FormSelectField
                      label="Publisher Name"
                      name="publisherName"
                      placeholder="Publisher Name"
                      selectedValue={submission.publisherName}
                      helpText=""
                      options={PUBLISHER_NAME_OPTIONS}
                      disabled={true}
                    />

                    {submission.publisherName === "Other" && (
                      <FormInputField
                        label="Publisher Name (Other)"
                        name="publisherNameOther"
                        placeholder="Text input"
                        value={submission.publisherNameOther}
                        helpText=""
                        isRequired={true}
                        maxWidth="280px"
                        disabled={true}
                      />
                    )}

                    <FormTextareaField
                      label="Special Notes (Optional)"
                      name="specialNotes"
                      placeholder="Text input"
                      value={submission.specialNotes}
                      isRequired={true}
                      maxWidth="280px"
                      helpText={"Max 638 characters"}
                      disabled={true}
                      rows={4}
                    />

                    <p class="subtitle is-6 pt-4">
                      <FontAwesomeIcon
                        className="fas"
                        icon={faMagnifyingGlass}
                      />
                      &nbsp;Summary of Findings
                    </p>
                    <hr />

                    {submission && (
                      <FormRadioField
                        label="Creases Finding"
                        name="creasesFinding"
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
                        maxWidth="180px"
                        disabled={true}
                      />
                    )}

                    {submission && (
                      <FormRadioField
                        label="Tears Finding"
                        name="tearsFinding"
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
                        maxWidth="180px"
                        disabled={true}
                      />
                    )}

                    {submission && (
                      <FormRadioField
                        label="Missing Parts Finding"
                        name="missingPartsFinding"
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
                        maxWidth="180px"
                        disabled={true}
                      />
                    )}

                    {submission && (
                      <FormRadioField
                        label="Stains/Marks/Substances"
                        name="stainsFinding"
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
                        maxWidth="180px"
                        disabled={true}
                      />
                    )}

                    {submission && (
                      <FormRadioField
                        label="Distortion/Colour"
                        name="distortionFinding"
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
                        maxWidth="180px"
                        disabled={true}
                      />
                    )}

                    {submission && (
                      <FormRadioField
                        label="Paper Quality Finding"
                        name="paperQualityFinding"
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
                        maxWidth="180px"
                        disabled={true}
                      />
                    )}

                    {submission && (
                      <FormRadioField
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
                        maxWidth="180px"
                        disabled={true}
                      />
                    )}

                    <FormRadioField
                      label="Shows signs of tampering/restoration"
                      name="showsSignsOfTamperingOrRestoration"
                      value={parseInt(
                        submission.showsSignsOfTamperingOrRestoration,
                      )}
                      opt1Value={2}
                      opt1Label="No"
                      opt2Value={1}
                      opt2Label="Yes"
                      maxWidth="180px"
                      disabled={true}
                    />

                    <FormTextareaField
                      label="Grading Notes (Optional)"
                      name="gradingNotes"
                      placeholder="Text input"
                      value={submission.gradingNotes}
                      isRequired={true}
                      maxWidth="280px"
                      helpText={"Max 638 characters"}
                      disabled={true}
                      rows={4}
                    />

                    <p class="subtitle is-6 pt-4">
                      <FontAwesomeIcon className="fas" icon={faBalanceScale} />
                      &nbsp;Grading
                    </p>
                    <hr />

                    <FormRadioField
                      label="Which type of grading scale would you prefer?"
                      name="gradingScale"
                      value={parseInt(submission.gradingScale)}
                      opt1Value={1}
                      opt1Label="Letter Grade (Poor-Near Mint)"
                      opt2Value={2}
                      opt2Label="Numbers (0.5-10.0)"
                      opt3Value={3}
                      opt3Label="CPS Percentage (5%-100%)"
                      maxWidth="180px"
                    />

                    {submission && submission.gradingScale === 1 && (
                      <>
                        <FormSelectField
                          label="Overall Letter Grade"
                          name="overallLetterGrade"
                          placeholder="Overall Letter Grade"
                          selectedValue={submission.overallLetterGrade}
                          helpText=""
                          options={FINDING_OPTIONS}
                          disabled={true}
                        />
                        {isNMPlusOpen && (
                          <>
                            <FormCheckboxField
                              label="Is Near Mint plus?"
                              name="isOverallLetterGradeNearMintPlus"
                              checked={
                                submission.isOverallLetterGradeNearMintPlus
                              }
                              disabled={true}
                              maxWidth="180px"
                            />
                          </>
                        )}
                      </>
                    )}

                    {submission && submission.gradingScale === 2 && (
                      <FormSelectField
                        label="Overall Number Grade"
                        name="overallNumberGrade"
                        placeholder="Overall Number Grade"
                        selectedValue={submission.overallNumberGrade}
                        helpText=""
                        options={OVERALL_NUMBER_GRADE_OPTIONS}
                        disabled={true}
                      />
                    )}

                    {submission && submission.gradingScale === 3 && (
                      <FormSelectField
                        label="CPS Percentage Grade"
                        name="cpsPercentageGrade"
                        placeholder="CPS Percentage Grade"
                        selectedValue={submission.cpsPercentageGrade}
                        helpText=""
                        options={CPS_PERCENTAGE_GRADE_OPTIONS}
                        disabled={true}
                      />
                    )}

                    <p class="subtitle is-6">
                      <FontAwesomeIcon className="fas" icon={faCogs} />
                      &nbsp;Settings
                    </p>
                    <hr />

                    {submission && (
                      <FormInputField
                        label="CPSR #"
                        name="cpsrn"
                        placeholder="Text input"
                        value={submission.cpsrn}
                        helpText="The unique identifier used by CPS for all submissions"
                        isRequired={true}
                        maxWidth="200px"
                        disabled={true}
                      />
                    )}

                    <FormSelectField
                      label="Store ID"
                      name="storeID"
                      placeholder="Pick"
                      selectedValue={submission.storeId}
                      helpText="Pick the store this user belongs to and will be limited by"
                      isRequired={true}
                      options={[
                        {
                          value: submission.storeId,
                          label: submission.storeName,
                        },
                      ]}
                      disabled={true}
                    />
                    <FormRadioField
                      label="Service Type"
                      name="role"
                      value={submission.serviceType}
                      opt1Value={1}
                      opt1Label="Pre-Screening Service"
                      opt2Value={2}
                      opt2Label="Pedigree Service"
                      maxWidth="180px"
                      disabled={true}
                    />
                    <FormRadioField
                      label="Status"
                      name="status"
                      value={submission.status}
                      opt1Value={1}
                      opt1Label="Pending"
                      opt2Value={2}
                      opt2Label="Active"
                      opt3Value={3}
                      opt3Label="Error"
                      opt4Value={4}
                      opt4Label="Archived"
                      maxWidth="180px"
                      disabled={true}
                    />

                    <div class="columns pt-4">
                      <div class="column is-half">
                        <Link
                          to={`/admin/submissions/comics`}
                          class="button is-medium is-hidden-touch"
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back
                        </Link>
                        <Link
                          to={`/admin/submissions/comics`}
                          class="button is-medium is-fullwidth is-hidden-desktop"
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back
                        </Link>
                      </div>
                      <div class="column is-half has-text-right">
                        <Link
                          to={`/admin/submissions/comic/${id}/edit`}
                          class="button is-medium is-primary is-hidden-touch"
                        >
                          <FontAwesomeIcon className="fas" icon={faPencil} />
                          &nbsp;Edit Comic Submission
                        </Link>
                        <Link
                          to={`/admin/submissions/comic/${id}/edit`}
                          class="button is-medium is-primary is-fullwidth is-hidden-desktop"
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

export default AdminComicSubmissionDetailContentComic;
