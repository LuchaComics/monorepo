import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faCog,
  faTasks,
  faTachometer,
  faPlus,
  faArrowLeft,
  faCheckCircle,
  faCubes,
  faGauge,
  faPencil,
  faEye,
  faIdCard,
  faAddressBook,
  faContactCard,
  faChartPie,
  faCogs,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import { getCollectionDetailAPI, putCollectionUpdateAPI } from "../../../../API/NFTCollection";
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
  USER_SPECIAL_COLLECTION_WITH_EMPTY_OPTIONS,
} from "../../../../Constants/FieldOptions";
import { topAlertMessageState, topAlertStatusState } from "../../../../AppState";
import { USER_ROLE_ROOT, USER_ROLE_RETAILER, USER_ROLE_CUSTOMER } from "../../../../Constants/App";


function AdminNFTCollectionUpdate() {
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
  const [tenantSelectOptions, setTenantSelectOptions] = useState([]);
  const [tenantID, setTenantID] = useState();
  const [name, setName] = useState("");
  const [status, setStatus] = useState();

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
      id: id,
      ID: id,
      Name: name,
      Status: status,
    };
    console.log("onSubmitClick, collection:", collection);
    putCollectionUpdateAPI(
      collection,
      onAdminNFTCollectionUpdateSuccess,
      onAdminNFTCollectionUpdateError,
      onAdminNFTCollectionUpdateDone,
      onUnauthorized,
    );
  };

  function onProfileDetailSuccess(response) {
    console.log("onProfileDetailSuccess: Starting...");
    setTenantID(response.tenantID);
    setName(response.name);
    setStatus(response.status);
  }

  function onProfileDetailError(apiErr) {
    console.log("onProfileDetailError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onProfileDetailDone() {
    console.log("onProfileDetailDone: Starting...");
    setFetching(false);
  }

  function onAdminNFTCollectionUpdateSuccess(response) {
    // For debugging purposes only.
    console.log("onAdminNFTCollectionUpdateSuccess: Starting...");
    console.log(response);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("NFT collection updated");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onAdminNFTCollectionUpdateSuccess: Delayed for 2 seconds.");
      console.log(
        "onAdminNFTCollectionUpdateSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Redirect the collection to a new page.
    setForceURL("/admin/collection/" + response.id);
  }

  function onAdminNFTCollectionUpdateError(apiErr) {
    console.log("onAdminNFTCollectionUpdateError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onAdminNFTCollectionUpdateError: Delayed for 2 seconds.");
      console.log(
        "onAdminNFTCollectionUpdateError: topAlertMessage, topAlertStatus:",
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

  function onAdminNFTCollectionUpdateDone() {
    console.log("onAdminNFTCollectionUpdateDone: Starting...");
    setFetching(false);
  }

  function onTenantOptionListSuccess(response) {
    console.log("onTenantOptionListSuccess: Starting...");
    console.log("onTenantOptionListSuccess: response:", response);
    if (response !== null) {
      const selectOptions = [
        { value: 0, label: "Please select" }, // Add empty options.
        ...response,
      ];
      console.log("onTenantOptionListSuccess: selectOptions:", selectOptions);
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

      setFetching(true);
      getCollectionDetailAPI(
        id,
        onProfileDetailSuccess,
        onProfileDetailError,
        onProfileDetailDone,
        onUnauthorized,
      );
      let params = new Map();
      getTenantSelectOptionListAPI(
        params,
        onTenantOptionListSuccess,
        onTenantOptionListError,
        onTenantOptionListDone,
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
              <li class="">
                <Link to={`/admin/collection/${id}`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faPencil} />
                  &nbsp;Update
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/admin/collection/${id}`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Detail
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modals */}
          {/* Do nothing... */}

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faCubes} />
              &nbsp;NFT Collection
            </p>
            <FormErrorBox errors={errors} />

            {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

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
                    errorText={errors && errors.tenantID}
                    helpText="Pick the tenant this collection belongs to and will be limited by"
                    isRequired={true}
                    onChange={(e) => setTenantID(e.target.value)}
                    options={tenantSelectOptions}
                    disabled={tenantSelectOptions.length === 0}
                  />

                  <FormInputField
                    label="Smart Contract"
                    name="name"
                    placeholder="Text input"
                    value={"Collectible Protection Service Submission Token"}
                    errorText={null}
                    helpText="All tokens in this collection will be implement this ERC721 based Smart Contract."
                    onChange={null}
                    isRequired={true}
                    maxWidth="400px"
                    disabled={true}
                  />

                  <FormRadioField
                    label="Status"
                    name="status"
                    value={status}
                    opt1Value={1}
                    opt1Label="Active"
                    opt2Value={2}
                    opt2Label="Archived"
                    errorText={errors && errors.status}
                    onChange={(e) => setStatus(parseInt(e.target.value))}
                    maxWidth="180px"
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
                    helpText=""
                    onChange={(e) => setName(e.target.value)}
                    isRequired={true}
                    maxWidth="380px"
                  />

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <Link
                        class="button is-hidden-touch"
                        to={`/admin/collection/${id}`}
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back
                      </Link>
                      <Link
                        class="button is-fullwidth is-hidden-desktop"
                        to={`/admin/collection/${id}`}
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back
                      </Link>
                    </div>
                    <div class="column is-half has-text-right">
                      <button
                        class="button is-primary is-hidden-touch"
                        onClick={onSubmitClick}
                      >
                        <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                        &nbsp;Save
                      </button>
                      <button
                        class="button is-primary is-fullwidth is-hidden-desktop"
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

export default AdminNFTCollectionUpdate;
