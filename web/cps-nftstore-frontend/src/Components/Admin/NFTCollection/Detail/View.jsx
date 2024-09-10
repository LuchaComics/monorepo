import React, { useState, useEffect } from "react";
import { Link, Navigate, useSearchParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faEllipsis,
  faTrashCan,
  faCog,
  faStar,
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
  faArrowRight,
  faArrowUpRightFromSquare
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import useLocalStorage from "../../../../Hooks/useLocalStorage";
import {
  getCollectionDetailAPI,
  deleteCollectionAPI,
} from "../../../../API/NFTCollection";
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
} from "../../../../AppState";
import FormRowText from "../../../Reusable/FormRowText";
import FormTextYesNoRow from "../../../Reusable/FormRowTextYesNo";
import FormTextOptionRow from "../../../Reusable/FormRowTextOption";
import FormTextChoiceRow from "../../../Reusable/FormRowTextChoice";
import {
  HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS,
  HOW_LONG_HAS_YOUR_STORE_BEEN_OPERATING_FOR_WITH_EMPTY_OPTIONS,
  USER_SPECIAL_COLLECTION_WITH_EMPTY_OPTIONS,
} from "../../../../Constants/FieldOptions";
import AlertBanner from "../../../Reusable/EveryPage/AlertBanner";
import {
  USER_ROLE_CUSTOMER
} from "../../../../Constants/App";

function AdminNFTCollectionDetail() {
  ////
  ////
  //// URL Query Strings.
  ////

  const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140
  const apiKey = searchParams.get("api_key");

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
  const [collection, setCollection] = useState({});
  const [tenantSelectOptions, setTenantSelectOptions] = useState([]);
  const [selectedCollectionForDeletion, setSelectedCollectionForDeletion] = useState("");

  ////
  //// Event handling.
  ////

  const onSelectCollectionForDeletion = (e, collection) => {
    console.log("onSelectCollectionForDeletion", collection);
    setSelectedCollectionForDeletion(collection);
  };

  const onDeselectCollectionForDeletion = (e) => {
    console.log("onDeselectCollectionForDeletion");
    setSelectedCollectionForDeletion("");
  };

  const onDeleteConfirmButtonClick = (e) => {
    console.log("onDeleteConfirmButtonClick"); // For debugging purposes only.

    deleteCollectionAPI(
      selectedCollectionForDeletion.id,
      onCollectionDeleteSuccess,
      onCollectionDeleteError,
      onCollectionDeleteDone,
      onUnauthorized,
    );
    setSelectedCollectionForDeletion("");
  };

  ////
  //// API.
  ////

  // --- DETAIL --- //

  function onCollectionDetailSuccess(response) {
    console.log("onCollectionDetailSuccess: Starting...");
    setCollection(response);
  }

  function onCollectionDetailError(apiErr) {
    console.log("onCollectionDetailError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onCollectionDetailDone() {
    console.log("onCollectionDetailDone: Starting...");
    setFetching(false);
  }

  // --- STORE OPTIONS --- //

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

  // --- DELETE --- //

  function onCollectionDeleteSuccess(response) {
    console.log("onCollectionDeleteSuccess: Starting..."); // For debugging purposes only.

    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("Collection deleted");
    setTimeout(() => {
      console.log(
        "onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Fetch again an updated list.
    setForceURL("/admin/collections");
  }

  function onCollectionDeleteError(apiErr) {
    console.log("onCollectionDeleteError: Starting..."); // For debugging purposes only.
    setErrors(apiErr);

    // Update notification.
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed deleting");
    setTimeout(() => {
      console.log(
        "onCollectionDeleteError: topAlertMessage, topAlertStatus:",
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

  function onCollectionDeleteDone() {
    console.log("onCollectionDeleteDone: Starting...");
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
        onCollectionDetailSuccess,
        onCollectionDetailError,
        onCollectionDetailDone,
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
                <Link to={`/admin/collections`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to NFT Collections
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modals */}
          <div class={`modal ${selectedCollectionForDeletion ? "is-active" : ""}`}>
            <div class="modal-background"></div>
            <div class="modal-card">
              <header class="modal-card-head">
                <p class="modal-card-title">Are you sure?</p>
                <button
                  class="delete"
                  aria-label="close"
                  onClick={onDeselectCollectionForDeletion}
                ></button>
              </header>
              <section class="modal-card-body">
              You are about to <b>delete</b> this collection; the data will be permanently deleted and no
              longer appear on your dashboard. This action cannot be undone. Are you sure
              you would like to continue?
              </section>
              <footer class="modal-card-foot">
                <button
                  class="button is-success"
                  onClick={onDeleteConfirmButtonClick}
                >
                  Confirm
                </button>
                <button class="button" onClick={onDeselectCollectionForDeletion}>
                  Cancel
                </button>
              </footer>
            </div>
          </div>
          <div class={`modal ${apiKey ? "is-active" : ""}`}>
            <div class="modal-background"></div>
            <div class="modal-card">
              <header class="modal-card-head">
                <p class="modal-card-title">Please confirm the API Key</p>
                <Link
                  class="delete"
                  aria-label="close"
                  to={`/admin/collection/${id}`}
                ></Link>
              </header>
              <section class="modal-card-body" style={{wordWrap: "break-word"}}>
                You have successfully created a collection! Here are your credentials. Please save the <strong>API Key</strong> somewhere safe as you'll never get access to it again after you click <strong>confirm</strong>. Note: If you forget your API key, you'll need to generate a new one.
                <br />
                <br />
                <strong>API Token:</strong>&nbsp;<br/><i>{apiKey}</i>
                <br />
                <br />
              </section>
              <footer class="modal-card-foot">
                <Link
                  class="button is-success"
                  to={`/admin/collection/${id}`}
                >
                  I confirm & close
                </Link>
              </footer>
            </div>
          </div>

          {/* Page banner */}
          {collection && collection.status === 100 && (
            <AlertBanner message="Archived" status="info" />
          )}

          {/* Page */}
          <nav class="box">
            {collection && (
              <div class="columns">
                <div class="column">
                  <p class="title is-4">
                    <FontAwesomeIcon className="fas" icon={faCubes} />
                    &nbsp;NFT Collection
                  </p>
                </div>
                {collection && collection.status === 1 && <div class="column has-text-right">
                  <Link
                    to={`/admin/collection/${collection.id}/nft-metadata/add`}
                    class="button is-small is-success is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faPlus} />
                    &nbsp;NFT Metadata
                  </Link>
                  &nbsp;&nbsp;
                  <Link
                    to={`/admin/collection/${collection.id}/edit`}
                    class="button is-small is-warning is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faPencil} />
                    &nbsp;Edit
                  </Link>
                  &nbsp;&nbsp;
                  <button
                    onClick={(e, ses) => onSelectCollectionForDeletion(e, collection)}
                    class="button is-small is-danger is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faTrashCan} />
                    &nbsp;Delete
                  </button>
                </div>}
              </div>
            )}
            <FormErrorBox errors={errors} />

            {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                {collection && (
                  <div class="container">
                    <div class="tabs is-medium is-size-7-mobile">
                      <ul>
                        <li class="is-active">
                          <Link>
                            <b>Detail</b>
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/collection/${collection.id}/nft-metadata`}>
                            NFT Metadata
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/collection/${collection.id}/more`}>
                            More&nbsp;&nbsp;
                            <FontAwesomeIcon
                              className="mdi"
                              icon={faEllipsis}
                            />
                          </Link>
                        </li>
                      </ul>
                    </div>

                    <p class="subtitle is-6">
                      <FontAwesomeIcon className="fas" icon={faCogs} />
                      &nbsp;Settings
                    </p>
                    <hr />

                    {tenantSelectOptions && tenantSelectOptions.length > 0 && (
                      <FormTextOptionRow
                        label="Tenant"
                        selectedValue={collection.tenantID}
                        helpText=""
                        options={tenantSelectOptions}
                      />
                    )}
                    <FormTextChoiceRow
                      label="Status"
                      value={collection.status}
                      opt1Value={1}
                      opt1Label="Active"
                      opt2Value={2}
                      opt2Label="Archived"
                    />

                    <p class="subtitle is-6">
                      <FontAwesomeIcon className="fas" icon={faIdCard} />
                      &nbsp;Detail
                    </p>
                    <hr />

                    <FormRowText
                      label="Name"
                      value={collection.name}
                      helpText=""
                    />

                    <FormRowText
                      label="IPNS"
                      value={`/ipns/${collection.ipnsName}`}
                      helpText={
                      <>
                      View file in <Link target="_blank" rel="noreferrer" to={`https://ipfs.io/ipns/${collection.ipnsName}`}>public IPFS gateway&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} /></Link>. Note: May take hours before image becomes available after submission.
                      </>} />

                    <FormRowText
                      label="Tokens Count"
                      type="number"
                      value={collection.tokensCount}
                      type="text"
                      helpText=""
                    />


                    <div class="columns pt-5">
                      <div class="column is-half">
                        <Link
                          class="button is-fullwidth-mobile"
                          to={`/admin/collections`}
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back to NFT Collections
                        </Link>
                      </div>
                      <div class="column is-half has-text-right">
                        {collection && collection.status === 1 && <Link
                          to={`/admin/collection/${id}/edit`}
                          class="button is-primary is-fullwidth-mobile"
                        >
                          <FontAwesomeIcon className="fas" icon={faPencil} />
                          &nbsp;Edit
                        </Link>}
                      </div>
                    </div>
                  </div>
                )}
              </>
            )}
          </nav>

          {/* Bottom Page Logout Link  */}
          {/*
          <div className="has-text-right has-text-grey">
            <Link to={`/admin/collection/${id}/nfts/add-via-ws`} className="has-text-grey">
              Add Pin via Web-Service API&nbsp;
              <FontAwesomeIcon className="mdi" icon={faArrowRight} />
            </Link>
          </div>
          */}
        </section>
      </div>
    </>
  );
}

export default AdminNFTCollectionDetail;
