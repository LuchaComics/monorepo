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


function AdminComicSubmissionAddStep7() {
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
  const [creasesFinding, setCreasesFinding] = useState(addComicSubmission.creasesFinding);
  const [tearsFinding, setTearsFinding] = useState(addComicSubmission.tearsFinding);
  const [missingPartsFinding, setMissingPartsFinding] = useState(addComicSubmission.missingPartsFinding);
  const [stainsFinding, setStainsFinding] = useState(addComicSubmission.stainsFinding);
  const [distortionFinding, setDistortionFinding] = useState(addComicSubmission.distortionFinding);
  const [paperQualityFinding, setPaperQualityFinding] = useState(addComicSubmission.paperQualityFinding);
  const [spineFinding, setSpineFinding] = useState(addComicSubmission.spineFinding);
  const [coverFinding, setCoverFinding] = useState(addComicSubmission.coverFinding);
  const [
    showsSignsOfTamperingOrRestoration,
    setShowsSignsOfTamperingOrRestoration,
  ] = useState(parseInt(addComicSubmission.showsSignsOfTamperingOrRestoration));

  ////
  //// Event handling.
  ////

  const onSubmitClick = (e) => {
      console.log("onSaveAndContinueClick: Beginning...");

      // Variables used to hold state if we got an error with validation.
      let newErrors = {};
      let hasErrors = false;

      // Perform validation.
      if (creasesFinding === undefined || creasesFinding === null || creasesFinding === 0 || creasesFinding === "") {
        newErrors["creasesFinding"] = "missing value";
        hasErrors = true;
      }
      if (tearsFinding === undefined || tearsFinding === null || tearsFinding === 0 || tearsFinding === "") {
        newErrors["tearsFinding"] = "missing value";
        hasErrors = true;
      }
      if (missingPartsFinding === undefined || missingPartsFinding === null || missingPartsFinding === 0 || missingPartsFinding === "") {
        newErrors["missingPartsFinding"] = "missing value";
        hasErrors = true;
      }
      if (stainsFinding === undefined || stainsFinding === null || stainsFinding === 0 || stainsFinding === "") {
        newErrors["stainsFinding"] = "missing value";
        hasErrors = true;
      }
      if (distortionFinding === undefined || distortionFinding === null || distortionFinding === 0 || distortionFinding === "") {
        newErrors["distortionFinding"] = "missing value";
        hasErrors = true;
      }
      if (paperQualityFinding === undefined || paperQualityFinding === null || paperQualityFinding === 0 || paperQualityFinding === "") {
        newErrors["paperQualityFinding"] = "missing value";
        hasErrors = true;
      }
      if (spineFinding === undefined || spineFinding === null || spineFinding === 0 || spineFinding === "") {
        newErrors["spineFinding"] = "missing value";
        hasErrors = true;
      }
      if (coverFinding === undefined || coverFinding === null || coverFinding === 0 || coverFinding === "") {
        newErrors["coverFinding"] = "missing value";
        hasErrors = true;
      }
      if (showsSignsOfTamperingOrRestoration === undefined || showsSignsOfTamperingOrRestoration === null || showsSignsOfTamperingOrRestoration === 0 || showsSignsOfTamperingOrRestoration === "") {
        newErrors["showsSignsOfTamperingOrRestoration"] = "missing value";
        hasErrors = true;
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

      console.log("onSaveAndContinueClick: Saving step 7 and redirecting to step 8.");

      // Variable holds a complete clone of the submission.
      let modifiedAddComicSubmission = { ...addComicSubmission };

      // Update our clone.
      modifiedAddComicSubmission.creasesFinding = creasesFinding;
      modifiedAddComicSubmission.tearsFinding = tearsFinding;
      modifiedAddComicSubmission.missingPartsFinding = missingPartsFinding;
      modifiedAddComicSubmission.stainsFinding = stainsFinding;
      modifiedAddComicSubmission.distortionFinding = distortionFinding;
      modifiedAddComicSubmission.paperQualityFinding = paperQualityFinding;
      modifiedAddComicSubmission.spineFinding = spineFinding;
      modifiedAddComicSubmission.coverFinding = coverFinding;
      modifiedAddComicSubmission.showsSignsOfTamperingOrRestoration = showsSignsOfTamperingOrRestoration;

      // Save to persistent storage.
      setAddComicSubmission(modifiedAddComicSubmission);

      // Redirect to the next page.
      setForceURL("/admin/submissions/comics/add/step-8")
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
                      &nbsp;Dashboard
                    </Link>
                  </li>
                  <li class="">
                    <Link to="/admin/customers" aria-current="page">
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
                  Please fill out all the required fields before continuing to the next step.
                </p>
                <div class="container">
                  <p class="subtitle is-6">
                    <FontAwesomeIcon
                      className="fas"
                      icon={faMagnifyingGlass}
                    />
                    &nbsp;Summary of Findings
                  </p>
                  <hr />
                  {addComicSubmission.serviceType !== SERVICE_TYPE_CPS_CAPSULE_INDIE_MINT_GEM ? (
                    <>


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
                    </>
                  ) : (<>
                      <article class="message">
                        <div class="message-body">
                          <p>
                            <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                            &nbsp;Completed, continue to next page.
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
                            setForceURL("/admin/submissions/comics/add/step-6");
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

export default AdminComicSubmissionAddStep7;
