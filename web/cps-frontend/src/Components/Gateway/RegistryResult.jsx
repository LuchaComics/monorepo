import React, { useState, useEffect } from "react";
import { Link, Navigate, useSearchParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faTachometer,
  faEye,
  faPencil,
  faTrashCan,
  faPlus,
  faGauge,
  faArrowRight,
  faBarcode,
  faArrowLeft,
  faMagnifyingGlass,
  faBook,
  faBalanceScale,
  faTriangleExclamation
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import FormErrorBox from "../Reusable/FormErrorBox";
import FormInputField from "../Reusable/FormInputField";
import FormDateField from "../Reusable/FormDateField";
import FormSelectField from "../Reusable/FormSelectField";
import FormRadioField from "../Reusable/FormRadioField";
import FormTextareaField from "../Reusable/FormTextareaField";
import { getRegistryAPI } from "../../API/registry";
import {
  FINDING_OPTIONS,
  OVERALL_NUMBER_GRADE_OPTIONS,
  PUBLISHER_NAME_OPTIONS,
  CPS_PERCENTAGE_GRADE_OPTIONS,
  HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS,
  ISSUE_COVER_YEAR_OPTIONS,
  ISSUE_COVER_MONTH_WITH_EMPTY_OPTIONS,
} from "../../Constants/FieldOptions";
import { currentUserState } from "../../AppState";


function PublicRegistryResult() {
  ////
  //// URL Parameters.
  ////

  const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140
  const cpsrn = searchParams.get("v");

  ////
  //// Global State
  ////

  const [currentUser] = useRecoilState(currentUserState);

  ////
  //// Component states.
  ////

  const [submission, setSubmission] = useState(null);
  const [errors, setErrors] = useState({});
  const [forceURL, setForceURL] = useState("");

  ////
  //// API.
  ////

  function onRegistrySuccess(response) {
    console.log("onRegistrySuccess: Starting...");
    console.log("registry:", response);
    setSubmission(response);
  }

  function onRegistryError(apiErr) {
    console.log("onRegistryError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onRegistryDone() {
    console.log("onRegistryDone: Starting...");
  }

  // --- All --- //

  const onUnauthorized = () => {
    setForceURL("/login?unauthorized=true"); // If token expired or user is not logged in, redirect back to login.
  };

  ////
  //// Event handling.
  ////

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      getRegistryAPI(
        cpsrn,
        onRegistrySuccess,
        onRegistryError,
        onRegistryDone,
        onUnauthorized,
      );
    }

    return () => (mounted = false);
  }, [cpsrn]);

  ////
  //// Component rendering.
  ////

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <>
      <div class="column is-12 container">
        <div class="section">
          <section class="hero is-fullheight">
            <div class="hero-body">
              <div class="container">
                <div class="columns is-centered">
                  <div class="column is-half-tablet">
                    <div class="box is-rounded">
                      {/* Start Logo */}
                      <nav class="level">
                        <div class="level-item has-text-centered">
                          <figure class="image">
                            <Link to="/">
                              <img
                                src="/static/CPS logo 2023 GR.webp"
                                style={{ width: "256px" }}
                              />
                            </Link>
                          </figure>
                        </div>
                      </nav>
                      {/* End Logo */}

                      {submission ? (
                        <form>
                          <h1 className="title is-4 has-text-centered">
                            Registry:
                          </h1>
                          {currentUser === null && (
                            <article className="message is-danger">
                              <div className="message-body">
                                <FontAwesomeIcon
                                  className="fas"
                                  icon={faTriangleExclamation}
                                />
                                &nbsp;An account is required for lookup. <Link to="/login"><b>Clich here&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} /></b></Link> to sign into your account.</div>
                            </article>
                          )}

                          <p class="subtitle is-6 has-text-centered pt-4">
                            <FontAwesomeIcon className="fas" icon={faBook} />
                            &nbsp;Comic Book Information
                          </p>
                          <hr />

                          <div class="field pb-4">
                            <label class="label">CPS #</label>
                            <div class="control has-icons-left has-icons-right">
                              <input
                                class={`input`}
                                name="cpsrn"
                                type="text"
                                placeholder="Enter CPS identification number."
                                value={submission.cpsrn}
                                disabled={true}
                              />
                              <span class="icon is-small is-left">
                                <FontAwesomeIcon
                                  className="fas"
                                  icon={faBarcode}
                                />
                              </span>
                            </div>
                            {errors && errors.cpsrn && (
                              <p class="help is-danger">{errors.cpsrn}</p>
                            )}
                            <p class="help">
                              The registry number found in our database.
                            </p>
                          </div>
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

                          <p class="subtitle is-6 has-text-centered pt-4">
                            <FontAwesomeIcon
                              className="fas"
                              icon={faMagnifyingGlass}
                            />
                            &nbsp;Summary of Findings
                          </p>
                          <hr />

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

                          <p class="subtitle is-6 has-text-centered pt-4">
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
                            <FormSelectField
                              label="Overall Letter Grade"
                              name="overallLetterGrade"
                              placeholder="Overall Letter Grade"
                              selectedValue={submission.overallLetterGrade}
                              helpText=""
                              options={FINDING_OPTIONS}
                              disabled={true}
                            />
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
                        </form>
                      ) : (
                        <>
                          <h1 className="title is-4 has-text-centered">
                            Registry:
                          </h1>
                          <article class="message is-danger">
                            <div class="message-body">
                              No registry record found for this CSPR #.
                            </div>
                          </article>
                        </>
                      )}

                      <br />
                      <br />
                      <nav class="level">
                        <div class="level-item has-text-centered">
                          <div>
                            <Link
                              to="/cpsrn-registry"
                              className="is-size-7-tablet"
                            >
                              <FontAwesomeIcon icon={faArrowLeft} />
                              &nbsp;Back
                            </Link>
                          </div>
                        </div>
                      </nav>
                    </div>
                    {/* End box */}

                    <div className="has-text-centered">
                      <p>Need help?</p>
                      <p>
                        <Link to="Support@cpscapsule.com">
                          Support@cpscapsule.com
                        </Link>
                      </p>
                      <p>
                        <a href="tel:+15199142685">(519) 914-2685</a>
                      </p>
                    </div>
                    {/* End suppoert text. */}
                  </div>
                  {/* End Column */}
                </div>
              </div>
              {/* End container */}
            </div>
            {/* End hero-body */}
          </section>
        </div>
      </div>
    </>
  );
}

export default PublicRegistryResult;
