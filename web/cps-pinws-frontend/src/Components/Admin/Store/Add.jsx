import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArrowLeft,
  faTasks,
  faTachometer,
  faPlus,
  faTimesCircle,
  faCheckCircle,
  faUserCircle,
  faGauge,
  faPencil,
  faUsers,
  faIdCard,
  faAddressBook,
  faContactCard,
  faChartPie,
  faBuilding,
  faCogs,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import useLocalStorage from "../../../Hooks/useLocalStorage";
import { postStoreCreateAPI } from "../../../API/store";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormInputField from "../../Reusable/FormInputField";
import FormTextareaField from "../../Reusable/FormTextareaField";
import FormRadioField from "../../Reusable/FormRadioField";
import FormMultiSelectField from "../../Reusable/FormMultiSelectField";
import FormSelectField from "../../Reusable/FormSelectField";
import FormCheckboxField from "../../Reusable/FormCheckboxField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import FormTimezoneSelectField from "../../Reusable/FormTimezoneField";
import {
  HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS,
  ESTIMATED_SUBMISSION_PER_MONTH_WITH_EMPTY_OPTIONS,
  HOW_LONG_HAS_YOUR_STORE_BEEN_OPERATING_FOR_WITH_EMPTY_OPTIONS,
  ORGANIZATION_LEVEL_WITH_EMPTY_OPTIONS,
  USER_SPECIAL_COLLECTION_WITH_EMPTY_OPTIONS,
} from "../../../Constants/FieldOptions";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";

function AdminStoreAdd() {
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
  const [name, setName] = useState("");
  const [showCancelWarning, setShowCancelWarning] = useState(false);
  const [websiteURL, setWebsiteURL] = useState("");
  const [estimatedSubmissionsPerMonth, setEstimatedSubmissionsPerMonth] =
    useState("");
  const [hasOtherGradingService, setHasOtherGradingService] = useState(0);
  const [otherGradingServiceName, setOtherGradingServiceName] = useState("");
  const [requestWelcomePackage, setRequestWelcomePackage] = useState(0);
  const [howLongStoreOperating, setHowLongStoreOperating] = useState(0);
  const [gradingComicsExperience, setGradingComicsExperience] = useState("");
  const [retailPartnershipReason, setRetailPartnershipReason] = useState("");
  const [cpsPartnershipReason, setCPSPartnershipReason] = useState("");
  const [status, setStatus] = useState(0);
  const [level, setLevel] = useState(0);
  const [specialCollection, setSpecialCollection] = useState(0);
  const [timezone, setTimezone] = useState(
    Intl.DateTimeFormat().resolvedOptions().timeZone,
  );

  ////
  //// Event handling.
  ////

  ////
  //// API.
  ////

  const onSubmitClick = (e) => {
    console.log("onSubmitClick: Beginning...");
    setFetching(true);
    setErrors({});
    const store = {
      Name: name,
      WebsiteUrl: websiteURL,
      EstimatedSubmissionsPerMonth: parseInt(estimatedSubmissionsPerMonth),
      HasOtherGradingService: parseInt(hasOtherGradingService),
      OtherGradingServiceName: otherGradingServiceName,
      RequestWelcomePackage: parseInt(requestWelcomePackage),
      HowLongStoreOperating: howLongStoreOperating,
      GradingComicsExperience: gradingComicsExperience,
      RetailPartnershipReason: retailPartnershipReason,
      CPSPartnershipReason: cpsPartnershipReason,
      Status: status,
      Level: level,
      SpecialCollection: specialCollection,
      Timezone: timezone,
    };
    console.log("onSubmitClick, store:", store);
    postStoreCreateAPI(
      store,
      onAdminStoreAddSuccess,
      onAdminStoreAddError,
      onAdminStoreAddDone,
      onUnauthorized,
    );
  };

  function onAdminStoreAddSuccess(response) {
    // For debugging purposes only.
    console.log("onAdminStoreAddSuccess: Starting...");
    console.log(response);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Store created");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onAdminStoreAddSuccess: Delayed for 2 seconds.");
      console.log(
        "onAdminStoreAddSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Redirect the user to a new page.
    setForceURL("/admin/store/" + response.id);
  }

  function onAdminStoreAddError(apiErr) {
    console.log("onAdminStoreAddError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onAdminStoreAddError: Delayed for 2 seconds.");
      console.log(
        "onAdminStoreAddError: topAlertMessage, topAlertStatus:",
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

  function onAdminStoreAddDone() {
    console.log("onAdminStoreAddDone: Starting...");
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

      setFetching(false);
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

  return (
    <>
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
            Your store record will be cancelled and your work will be lost. This
            cannot be undone. Do you want to continue?
          </section>
          <footer class="modal-card-foot">
            <Link class="button is-medium is-success" to={`/admin/stores`}>
              Yes
            </Link>
            <button
              class="button is-medium"
              onClick={(e) => setShowCancelWarning(false)}
            >
              No
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
                <Link to="/admin/stores" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faBuilding} />
                  &nbsp;Stores
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
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/admin/stores`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Stores
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faPlus} />
              &nbsp;New Store
            </p>

            {/* <p class="pb-4 has-text-grey">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Submitting..."} />
            ) : (
              <>
                <FormErrorBox errors={errors} />
                <div class="container">
                  <p class="subtitle is-6">
                    <FontAwesomeIcon className="fas" icon={faIdCard} />
                    &nbsp;Identification
                  </p>
                  <hr />

                  <FormInputField
                    label="Name"
                    name="name"
                    placeholder="Text input"
                    value={name}
                    errorText={errors && errors.name}
                    helpText=""
                    onChange={(e) => setName(e.target.value)}
                    isRequired={true}
                    maxWidth="380px"
                  />

                  <FormInputField
                    label="What is your website address?"
                    name="websiteURL"
                    placeholder="URL input"
                    value={websiteURL}
                    errorText={errors && errors.websiteUrl}
                    helpText=""
                    onChange={(e) => setWebsiteURL(e.target.value)}
                    isRequired={true}
                    maxWidth="100%"
                  />

                  <p class="subtitle is-6">
                    <FontAwesomeIcon className="fas" icon={faChartPie} />
                    &nbsp;Metrics
                  </p>
                  <hr />

                  <FormSelectField
                    label="How many comic books are you planning to submit to us per month?"
                    name="estimatedSubmissionsPerMonth"
                    placeholder="Pick"
                    selectedValue={estimatedSubmissionsPerMonth}
                    errorText={errors && errors.estimatedSubmissionsPerMonth}
                    helpText=""
                    onChange={(e) =>
                      setEstimatedSubmissionsPerMonth(parseInt(e.target.value))
                    }
                    options={ESTIMATED_SUBMISSION_PER_MONTH_WITH_EMPTY_OPTIONS}
                  />

                  <FormRadioField
                    label="Are you currently submitting to any other grading companies?"
                    name="hasOtherGradingService"
                    value={hasOtherGradingService}
                    opt1Value={1}
                    opt1Label="Yes"
                    opt2Value={2}
                    opt2Label="No"
                    errorText={errors && errors.hasOtherGradingService}
                    onChange={(e) =>
                      setHasOtherGradingService(parseInt(e.target.value))
                    }
                    maxWidth="180px"
                  />

                  {hasOtherGradingService === 1 && (
                    <FormInputField
                      label="Other Grading Service Name (Optional)"
                      name="otherGradingServiceName"
                      placeholder="Text input"
                      value={otherGradingServiceName}
                      errorText={errors && errors.otherGradingServiceName}
                      helpText=""
                      onChange={(e) =>
                        setOtherGradingServiceName(e.target.value)
                      }
                      isRequired={true}
                      maxWidth="380px"
                    />
                  )}

                  <FormRadioField
                    label="Would you like to receive a welcome package? This package includes promotional items and tools to help you improve your submissions, as well as our service terms and conditions."
                    name="requestWelcomePackage"
                    value={requestWelcomePackage}
                    opt1Value={1}
                    opt1Label="Yes"
                    opt2Value={2}
                    opt2Label="No"
                    errorText={errors && errors.requestWelcomePackage}
                    onChange={(e) =>
                      setRequestWelcomePackage(parseInt(e.target.value))
                    }
                    maxWidth="180px"
                  />

                  <FormSelectField
                    label="How long has your store been operating for?"
                    name="howLongStoreOperating"
                    placeholder="Pick"
                    selectedValue={howLongStoreOperating}
                    errorText={errors && errors.howLongStoreOperating}
                    helpText=""
                    onChange={(e) =>
                      setHowLongStoreOperating(parseInt(e.target.value))
                    }
                    options={
                      HOW_LONG_HAS_YOUR_STORE_BEEN_OPERATING_FOR_WITH_EMPTY_OPTIONS
                    }
                  />

                  <FormTextareaField
                    label="Tell us about your level of experience with grading comics? (Optional)"
                    name="gradingComicsExperience"
                    placeholder="Text input"
                    value={gradingComicsExperience}
                    errorText={errors && errors.gradingComicsExperience}
                    helpText=""
                    onChange={(e) => setGradingComicsExperience(e.target.value)}
                    isRequired={true}
                    maxWidth="280px"
                    helpText={""}
                    rows={4}
                  />

                  <FormTextareaField
                    label="Please describe how you could become a good retail partner for CPS"
                    name="retailPartnershipReason"
                    placeholder="Text input"
                    value={retailPartnershipReason}
                    errorText={errors && errors.retailPartnershipReason}
                    helpText=""
                    onChange={(e) => setRetailPartnershipReason(e.target.value)}
                    isRequired={true}
                    maxWidth="280px"
                    helpText={""}
                    rows={4}
                  />

                  <FormTextareaField
                    label="Please describe how CPS could help you grow your business"
                    name="cpsPartnershipReason"
                    placeholder="Text input"
                    value={cpsPartnershipReason}
                    errorText={errors && errors.cpsPartnershipReason}
                    helpText=""
                    onChange={(e) => setCPSPartnershipReason(e.target.value)}
                    isRequired={true}
                    maxWidth="280px"
                    helpText={""}
                    rows={4}
                  />

                  <p class="subtitle is-6">
                    <FontAwesomeIcon className="fas" icon={faCogs} />
                    &nbsp;Settings
                  </p>
                  <hr />

                  <FormRadioField
                    label="Status"
                    name="status"
                    value={status}
                    opt1Value={1}
                    opt1Label="Pending"
                    opt2Value={2}
                    opt2Label="Active"
                    opt3Value={3}
                    opt3Label="Rejected"
                    opt4Value={5}
                    opt4Label="Archived"
                    errorText={errors && errors.status}
                    onChange={(e) => setStatus(parseInt(e.target.value))}
                    maxWidth="180px"
                  />

                  <FormSelectField
                    label="Level"
                    name="level"
                    placeholder="Pick"
                    selectedValue={level}
                    errorText={errors && errors.level}
                    helpText=""
                    onChange={(e) => setLevel(parseInt(e.target.value))}
                    options={ORGANIZATION_LEVEL_WITH_EMPTY_OPTIONS}
                  />

                  <FormSelectField
                    label="Special Collection (Optional)"
                    name="specialCollection"
                    placeholder="Pick"
                    selectedValue={specialCollection}
                    errorText={errors && errors.specialCollection}
                    helpText=""
                    onChange={(e) =>
                      setSpecialCollection(parseInt(e.target.value))
                    }
                    options={USER_SPECIAL_COLLECTION_WITH_EMPTY_OPTIONS}
                  />

                  <FormTimezoneSelectField
                    label="Timezone"
                    name="timezone"
                    placeholder="Text input"
                    selectedTimezone={timezone}
                    setSelectedTimezone={(value) => setTimezone(value)}
                    errorText={errors && errors.timezone}
                    helpText="Please select the timezone that your business operates in."
                    isRequired={true}
                    maxWidth="550px"
                  />

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <button
                        class="button is-medium is-fullwidth-mobile"
                        onClick={(e) => setShowCancelWarning(true)}
                      >
                        <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                        &nbsp;Cancel
                      </button>
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

export default AdminStoreAdd;
