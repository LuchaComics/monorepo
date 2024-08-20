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
import {
  addComicSubmissionState,
  ADD_COMIC_SUBMISSION_STATE_DEFAULT,
} from "../../../../AppState";
import { parseBool } from "../../../../Helpers/boolUtility";


function AdminComicSubmissionAddStep9() {
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
  const [gradingScale, setGradingScale] = useState(parseInt(addComicSubmission.gradingScale));
  const [overallLetterGrade, setOverallLetterGrade] = useState(addComicSubmission.overallLetterGrade);
  const [overallNumberGrade, setOverallNumberGrade] = useState(parseFloat(addComicSubmission.overallNumberGrade));
  const [cpsPercentageGrade, setCpsPercentageGrade] = useState(parseFloat(addComicSubmission.cpsPercentageGrade));
  const [showCancelWarning, setShowCancelWarning] = useState(false);
  const [
    isOverallLetterGradeNearMintPlus,
    setIsOverallLetterGradeNearMintPlus,
  ] = useState(parseBool(addComicSubmission.isOverallLetterGradeNearMintPlus));
  const [serviceType, setServiceType] = useState(addComicSubmission.serviceType);

  ////
  //// Event handling.
  ////

  const onSaveAndContinueClick = (e) => {
    console.log("onSaveAndContinueClick: Beginning...");

    // Variables used to hold state if we got an error with validation.
    let newErrors = {};
    let hasErrors = false;

    // Perform validation.
    if (gradingScale === undefined || gradingScale === null || gradingScale === 0 || gradingScale === "" || isNaN(gradingScale)) {
      newErrors["gradingScale"] = "missing value";
      hasErrors = true;
    } else {
        // CASE 1 of 3: Letter Grade
        if (gradingScale === 1) {
            if (overallLetterGrade === undefined || overallLetterGrade === null || overallLetterGrade === 0 || overallLetterGrade === "") {
              newErrors["overallLetterGrade"] = "missing value";
              hasErrors = true;
            }
        }

        // CASE 2 of 3: Numbers
        if (gradingScale === 2) {
            if (overallNumberGrade === undefined || overallNumberGrade === null || overallNumberGrade === 0 || overallNumberGrade === "" || isNaN(overallNumberGrade)) {
              newErrors["overallNumberGrade"] = "missing value";
              hasErrors = true;
            }
        }

        // CASE 3 of 3: Numbers
        if (gradingScale === 3) {
            if (cpsPercentageGrade === undefined || cpsPercentageGrade === null || cpsPercentageGrade === 0 || cpsPercentageGrade === "" || isNaN(cpsPercentageGrade)) {
              newErrors["cpsPercentageGrade"] = "missing value";
              hasErrors = true;
            }
        }
    }

    // For debuggin purposes only.
    console.log("onSaveAndContinueClick: gradingScale:", gradingScale);
    console.log("onSaveAndContinueClick: overallLetterGrade:", overallLetterGrade);
    console.log("onSaveAndContinueClick: overallNumberGrade:", overallNumberGrade);
    console.log("onSaveAndContinueClick: cpsPercentageGrade:", cpsPercentageGrade)
    console.log("onSaveAndContinueClick: hasErrors:", hasErrors);

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

    console.log("onSaveAndContinueClick: Saving step 9 and redirecting to step 10.");

    // Variable holds a complete clone of the submission.
    let modifiedAddComicSubmission = { ...addComicSubmission };

    // Update our clone.
    modifiedAddComicSubmission.gradingScale = parseInt(gradingScale); // 2=Number Grading Scale
    modifiedAddComicSubmission.overallLetterGrade = overallLetterGrade;
    modifiedAddComicSubmission.isOverallLetterGradeNearMintPlus = isOverallLetterGradeNearMintPlus;
    modifiedAddComicSubmission.overallNumberGrade = parseFloat(overallNumberGrade);
    modifiedAddComicSubmission.cpsPercentageGrade = parseFloat(cpsPercentageGrade);

    // Save to persistent storage.
    setAddComicSubmission(modifiedAddComicSubmission);

    // Redirect to the next page.
    setForceURL("/admin/submissions/comics/add/step-10")
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
                      &nbsp;Admin Dashboard
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
          <nav className="box has-background-light">
            <p className="subtitle is-5">Step 9 of 10</p>
            <progress
              class="progress is-success"
              value="90"
              max="100"
            >
              90%
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
                  Please fill out all the required fields before continuing to the last step.
                </p>
                <div class="container">
                  <p class="subtitle is-6">
                    <FontAwesomeIcon
                      className="fas"
                      icon={faBalanceScale}
                    />
                    &nbsp;Grading
                  </p>
                  <hr />

                  {serviceType !== SERVICE_TYPE_CPS_CAPSULE_INDIE_MINT_GEM ? (
                    <>
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
                  ) : (<>
                      <article class="message">
                        <div class="message-body">
                          <p>
                            <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                            &nbsp;Auto-set by <strong>Indie Mint Gem</strong> service type, continue to next page.
                          </p>
                        </div>
                      </article>
                  </>)}

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <button
                        class="button is-medium is-fullwidth-mobile"
                        onClick={(e) => {
                            e.preventDefault();
                            setForceURL("/admin/submissions/comics/add/step-8");
                        }}
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back to Step 8
                      </button>
                    </div>
                    <div class="column is-half has-text-right">
                      <button
                        class="button is-medium is-primary is-fullwidth-mobile"
                        onClick={onSaveAndContinueClick}
                      >
                        Save and Review&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} />
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

export default AdminComicSubmissionAddStep9;
