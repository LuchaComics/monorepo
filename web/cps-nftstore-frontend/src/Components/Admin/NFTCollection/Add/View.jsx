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
  HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS,
  HOW_LONG_HAS_YOUR_STORE_BEEN_OPERATING_FOR_WITH_EMPTY_OPTIONS,
} from "../../../../Constants/FieldOptions";
import { topAlertMessageState, topAlertStatusState } from "../../../../AppState";
import { USER_ROLE_RETAILER, USER_ROLE_CUSTOMER } from "../../../../Constants/App";


function AdminNFTCollectionAdd() {
  ////
  //// URL Parameters.
  ////

  const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140
  const orgID = searchParams.get("tenant_id");
  const orgName = searchParams.get("tenant_name");

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
  const [tenantSelectOptions, setTenantSelectOptions] = useState([]);
  const [tenantID, setTenantID] = useState(orgID);

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
    const collection = {
      TenantID: tenantID,
      Status: 1, // 1 = CollectionActiveStatus
      Name: name,
    };
    console.log("onSubmitClick, collection:", collection);
    postCollectionCreateAPI(
      collection,
      onAdminNFTCollectionAddSuccess,
      onAdminNFTCollectionAddError,
      onAdminNFTCollectionAddDone,
      onUnauthorized,
    );
  };

  function onAdminNFTCollectionAddSuccess(response) {
    // For debugging purposes only.
    console.log("onAdminNFTCollectionAddSuccess: Starting...");
    console.log(response);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Collection created");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onAdminNFTCollectionAddSuccess: Delayed for 2 seconds.");
      console.log(
        "onAdminNFTCollectionAddSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    const apiKey = response.apiKey;

    if (orgName !== undefined && orgName !== null && orgName !== "") {
      // Redirect the collection to a new page.
      setForceURL("/admin/tenant/" + orgID);
    } else {
      // Redirect the collection to a new page.
      setForceURL("/admin/collection/" + response.id);
    }
  }

  function onAdminNFTCollectionAddError(apiErr) {
    console.log("onAdminNFTCollectionAddError: Starting...");
    console.log("onAdminNFTCollectionAddError: apiErr:", apiErr);
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      // console.log("onAdminNFTCollectionAddError: Delayed for 2 seconds.");
      // console.log("onAdminNFTCollectionAddError: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
      setTopAlertMessage("");
    }, 2000);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onAdminNFTCollectionAddDone() {
    console.log("onAdminNFTCollectionAddDone: Starting...");
    setFetching(false);
  }

  function onTenantOptionListSuccess(response) {
    console.log("onTenantOptionListSuccess: Starting...");
    if (response !== null) {
      const selectOptions = [
        { value: 0, label: "Please select" }, // Add empty options.
        ...response,
      ];
      setTenantSelectOptions(selectOptions);
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
            {orgName !== undefined && orgName !== null && orgName !== "" ? (
              <ul>
                <li class="">
                  <Link to="/admin/dashboard" aria-current="page">
                    <FontAwesomeIcon className="fas" icon={faGauge} />
                    &nbsp;Admin Dashboard
                  </Link>
                </li>
                <li class="">
                  <Link to="/admin/tenants" aria-current="page">
                    <FontAwesomeIcon className="fas" icon={faBuilding} />
                    &nbsp;Tenants
                  </Link>
                </li>
                <li class="">
                  <Link to={`/admin/tenant/${orgID}/collections`} aria-current="page">
                    <FontAwesomeIcon className="fas" icon={faEye} />
                    &nbsp;Detail (Collections)
                  </Link>
                </li>
                <li class="is-active">
                  <Link aria-current="page">
                    <FontAwesomeIcon className="fas" icon={faPlus} />
                    &nbsp;New Collection
                  </Link>
                </li>
              </ul>
            ) : (
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
                    &nbsp;Collections
                  </Link>
                </li>
                <li class="is-active">
                  <Link aria-current="page">
                    <FontAwesomeIcon className="fas" icon={faPlus} />
                    &nbsp;New
                  </Link>
                </li>
              </ul>
            )}
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              {orgName !== undefined && orgName !== null && orgName !== "" ? (
                <li class="">
                  <Link to={`/admin/tenant/${orgID}/collections`} aria-current="page">
                    <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                    &nbsp;Back to Detail (Collections)
                  </Link>
                </li>
              ) : (
                <li class="">
                  <Link to={`/admin/collections`} aria-current="page">
                    <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                    &nbsp;Back to Collections
                  </Link>
                </li>
              )}
            </ul>
          </nav>

          {/* Modals */}
          {/* None */}

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
                  {orgName !== undefined &&
                  orgName !== null &&
                  orgName !== "" ? (
                    <Link
                      class="button is-medium is-success"
                      to={`/admin/tenant/${orgID}/collections`}
                    >
                      Yes
                    </Link>
                  ) : (
                    <Link
                      class="button is-medium is-success"
                      to={`/admin/collections`}
                    >
                      Yes
                    </Link>
                  )}
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
              &nbsp;New Collection
            </p>
            <FormErrorBox errors={errors} />

            {/* <p class="pb-4 has-text-grey">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Submitting..."} />
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
                      <button
                        class="button is-medium is-hidden-touch"
                        onClick={(e) => setShowCancelWarning(true)}
                      >
                        <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                        &nbsp;Cancel
                      </button>
                      <button
                        class="button is-medium is-fullwidth is-hidden-desktop"
                        onClick={(e) => setShowCancelWarning(true)}
                      >
                        <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                        &nbsp;Cancel
                      </button>
                    </div>
                    <div class="column is-half has-text-right">
                      <button
                        class="button is-medium is-primary is-hidden-touch"
                        onClick={onSubmitClick}
                      >
                        <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                        &nbsp;Save
                      </button>
                      <button
                        class="button is-medium is-primary is-fullwidth is-hidden-desktop"
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

export default AdminNFTCollectionAdd;
