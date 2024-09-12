import React, { useState, useEffect } from "react";
import { Link, Navigate, useSearchParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faCog,
  faArrowLeft,
  faTasks,
  faTachometer,
  faPlus,
  faTimesCircle,
  faCheckCircle,
  faCollectionCircle,
  faGauge,
  faPencil,
  faCubes,
  faIdCard,
  faAddressBook,
  faContactCard,
  faChartPie,
  faCogs,
  faBuilding,
  faEye,
  faHourglassStart,
  faExclamationTriangle,
  faArrowRight
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import useLocalStorage from "../../../../Hooks/useLocalStorage";
import { postCollectionCreateAPI } from "../../../../API/NFTCollection";
import { getTenantSelectOptionListAPI } from "../../../../API/tenant";
import FormErrorBox from "../../../Reusable/FormErrorBox";
import FormInputField from "../../../Reusable/FormInputField";
import FormTextareaField from "../../../Reusable/FormTextareaField";
import FormRadioField from "../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../Reusable/FormSelectField";
import FormCheckboxField from "../../../Reusable/FormCheckboxField";
import FormCountryField from "../../../Reusable/FormCountryField";
import FormRegionField from "../../../Reusable/FormRegionField";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import {
    topAlertMessageState,
    topAlertStatusState,
    addNFTCollectionState,
    ADD_NFT_COLLECTION_STATE_DEFAULT
} from "../../../../AppState";


function AdminNFTCollectionAddStep2() {
  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);
  const [addNFTCollection, setAddNFTCollection] = useRecoilState(addNFTCollectionState);

  ////
  //// Component states.
  ////

  // GUI states.
  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [showCancelWarning, setShowCancelWarning] = useState(false);
  const [tenantSelectOptions, setTenantSelectOptions] = useState([]);

  // Form states.
  const [tenantID, setTenantID] = useState(addNFTCollection.tenantID);
  const [tenantName, setTenantName] = useState(addNFTCollection.tenantName);
  const [name, setName] = useState(addNFTCollection.name);

  ////
  //// Event handling.
  ////

  const onSubmitClick = (e) => {
    e.preventDefault();
    console.log("onSubmitClick: Beginning...");
    setErrors({}); // Reset the errors in the GUI.

    // Variables used to create our new errors if we find them.
    let newErrors = {};
    let hasErrors = false;

    if (tenantID === undefined || tenantID === null || tenantID === "") {
      newErrors["tenantID"] = "missing value";
      hasErrors = true;
    }
    if (name === undefined || name === null || name === "") {
      newErrors["name"] = "missing value";
      hasErrors = true;
    }

    if (hasErrors) {
      console.log("onSubmitClick: Aboring because of error(s)");

      // Set the associate based error validation.
      setErrors(newErrors);

      // The following code will cause the screen to scroll to the top of
      // the page. Please see ``react-scroll`` for more information:
      // https://github.com/fisshy/react-scroll
      var scroll = Scroll.animateScroll;
      scroll.scrollToTop();

      return;
  }

    console.log("onSubmitClick: Success");

    // Save to persistent storage.
    let modifiedAddNFTCollection = { ...addNFTCollection };
    modifiedAddNFTCollection.tenantID = tenantID;
    modifiedAddNFTCollection.tenantName = tenantName;
    modifiedAddNFTCollection.name = name;
    setAddNFTCollection(modifiedAddNFTCollection);
    setForceURL("/admin/collections/add/step-3");
  };

  ////
  //// API.
  ////

  function onTenantOptionListSuccess(response) {
    console.log("onTenantOptionListSuccess: Starting...");
    if (response !== null) {
      const selectOptions = [
        { value: 0, label: "Please select" }, // Add empty options.
        ...response,
      ];
      setTenantSelectOptions(selectOptions);

      response.forEach(function (item, index) {
          console.log(item);
          if (tenantID === item.value) {
              setTenantName(item.label);
          }
      });


    }
  }

  function onTenantOptionListError(apiErr) {
    console.log("onTenantOptionListError: Starting...");
    console.log("onTenantOptionListError: apiErr:", apiErr);
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onTenantOptionListDone() {
    console.log("onTenantOptionListDone: Starting...");
    setFetching(false);
  }

  // --- All --- //

  const onUnauthorized = () => {
    setForceURL("/login?unauthorized=true"); // If token expired or collection is not logged in, redirect back to login.
  };

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
      let params = new Map();
      getTenantSelectOptionListAPI(
        params,
        onTenantOptionListSuccess,
        onTenantOptionListError,
        onTenantOptionListDone,
        onUnauthorized,
      );
      setFetching(true);
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
                  <Link to="/admin/collections" aria-current="page">
                    <FontAwesomeIcon className="fas" icon={faCubes} />
                    &nbsp;NFT Collections
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
                  <Link to={`/admin/collections`} aria-current="page">
                    <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                    &nbsp;Back to NFT Collections
                  </Link>
                </li>
            </ul>
          </nav>

          {/* Modals */}
          {/* None */}

          {/* Progress Wizard */}
          <nav className="box has-background-light">
            <p className="subtitle is-5">Step 2 of 3</p>
            <progress
              class="progress is-success"
              value="66"
              max="100"
            >
              43%
            </progress>
          </nav>

          {/* Page */}
          <nav class="box">
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
                  Your collection record will be cancelled and your work will be lost.
                  This cannot be undone. Do you want to continue?
                </section>
                <footer class="modal-card-foot">
                <Link
                  class="button is-medium is-success"
                  to={`/admin/collections`}
                >
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

            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faPlus} />
              &nbsp;New NFT Collection
            </p>
            <FormErrorBox errors={errors} />

            {/* <p class="pb-4 has-text-grey">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <>
                 <article class="message is-warning">
                   <div class="message-body">
                     <strong><FontAwesomeIcon className="fas" icon={faExclamationTriangle} />&nbsp;Warning:</strong>&nbsp;Submitting to IPFS network may sometimes take 5 minutes or more, please wait until completion...
                   </div>
                 </article>
                 <PageLoadingContent displayMessage={"Submitting..."} />
              </>
            ) : (
              <>
                <div class="container">
                  <p class="subtitle is-6">
                    <FontAwesomeIcon className="fas" icon={faCogs} />
                    &nbsp;Settings
                  </p>
                  <hr />

                  <FormSelectField
                    label="Tenant ID"
                    name="tenantID"
                    placeholder="Pick"
                    selectedValue={tenantID}
                    errorText={errors && errors.tenantId}
                    helpText="Pick the tenant this collection belongs to and will be limited by"
                    isRequired={true}
                    onChange={(e) => setTenantID(e.target.value)}
                    options={tenantSelectOptions}
                    disabled={tenantSelectOptions.length === 0}
                  />

                  <p class="subtitle is-6">
                    <FontAwesomeIcon className="fas" icon={faIdCard} />
                    &nbsp;Detail
                  </p>
                  <hr />

                  <FormInputField
                    label="Name"
                    name="name"
                    placeholder="Text input"
                    value={name}
                    errorText={errors && errors.name}
                    helpText="This field will not be shown to NFT purchasers, only used for internal purposes"
                    onChange={(e) => setName(e.target.value)}
                    isRequired={true}
                    maxWidth="380px"
                  />

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <Link
                        class="button is-medium is-hidden-touch"
                        to="/admin/collections/add/step-1"
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back to Step 1
                      </Link>
                      <Link
                        class="button is-medium is-fullwidth is-hidden-desktop"
                        to="/admin/collections/add/step-1"
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back to Step 1
                      </Link>
                    </div>
                    <div class="column is-half has-text-right">
                      <button
                        class="button is-medium is-primary is-hidden-touch"
                        onClick={onSubmitClick}
                      >
                        &nbsp;Save & Continue&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} />
                      </button>
                      <button
                        class="button is-medium is-primary is-fullwidth is-hidden-desktop"
                        onClick={onSubmitClick}
                      >
                        &nbsp;Save & Continue&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} />
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

export default AdminNFTCollectionAddStep2;
