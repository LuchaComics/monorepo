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
  faChain
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
import DataDisplayRowText from "../../../Reusable/DataDisplayRowText";
import DataDisplayRowRadio from "../../../Reusable/DataDisplayRowRadio";
import {
    topAlertMessageState,
    topAlertStatusState,
    addNFTCollectionState,
    ADD_NFT_COLLECTION_STATE_DEFAULT
} from "../../../../AppState";


function AdminNFTCollectionAddStep3() {
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

  // GUI related states.
  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [showCancelWarning, setShowCancelWarning] = useState(false);
  const [tenantSelectOptions, setTenantSelectOptions] = useState([]);

  // Submission form.

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
    console.log("onSubmitClick, collection:", addNFTCollection);
    postCollectionCreateAPI(
      addNFTCollection,
      onAdminNFTCollectionAddStep3Success,
      onAdminNFTCollectionAddStep3Error,
      onAdminNFTCollectionAddStep3Done,
      onUnauthorized,
    );
  };

  function onAdminNFTCollectionAddStep3Success(response) {
    // For debugging purposes only.
    console.log("onAdminNFTCollectionAddStep3Success: Starting...");
    console.log(response);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("NFT collection created");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onAdminNFTCollectionAddStep3Success: Delayed for 2 seconds.");
      console.log(
        "onAdminNFTCollectionAddStep3Success: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    const apiKey = response.apiKey;

    setForceURL("/admin/collection/" + response.id);
  }

  function onAdminNFTCollectionAddStep3Error(apiErr) {
    console.log("onAdminNFTCollectionAddStep3Error: Starting...");
    console.log("onAdminNFTCollectionAddStep3Error: apiErr:", apiErr);
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      // console.log("onAdminNFTCollectionAddStep3Error: Delayed for 2 seconds.");
      // console.log("onAdminNFTCollectionAddStep3Error: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
      setTopAlertMessage("");
    }, 2000);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onAdminNFTCollectionAddStep3Done() {
    console.log("onAdminNFTCollectionAddStep3Done: Starting...");
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
          <nav className="box has-background-success-light">
            <p className="subtitle is-5">Step 3 of 3</p>
            <progress
              class="progress is-success"
              value="100"
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
            <p className="has-text-grey pb-4">
               Please review the following NFT collection summary before submitting into the system.
            </p>

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
                    <p className="title is-5 mt-2">
                      <FontAwesomeIcon className="fas" icon={faChain} />
                      &nbsp;Blockchain&nbsp;&nbsp;&nbsp;-&nbsp;&nbsp;&nbsp;
                      <Link to="/admin/collections/add/step-1">
                        <FontAwesomeIcon className="fas" icon={faPencil} />
                        &nbsp;Edit
                      </Link>
                    </p>

                    <DataDisplayRowRadio
                      label="Blockchain"
                      value={addNFTCollection.blockchain}
                      opt1Value="ethereum"
                      opt1Label="Ethereum"
                    />

                    <DataDisplayRowText
                      label="Node URL"
                      value={addNFTCollection.nodeURL}
                    />

                    <DataDisplayRowRadio
                      label="Smart Contract"
                      value={addNFTCollection.smartContract}
                      opt1Value={"Collectible Protection Service Submissions"}
                      opt1Label="Collectible Protection Service Submissions"
                    />

                    <DataDisplayRowText
                      label="Wallet Mnemonic"
                      value={addNFTCollection.walletMnemonic}
                    />

                    <p className="title is-5 mt-2">
                      <FontAwesomeIcon className="fas" icon={faIdCard} />
                      &nbsp;Collection Information&nbsp;&nbsp;&nbsp;-&nbsp;&nbsp;&nbsp;
                      <Link to="/admin/collections/add/step-2">
                        <FontAwesomeIcon className="fas" icon={faPencil} />
                        &nbsp;Edit
                      </Link>
                    </p>

                    <DataDisplayRowText
                      label="Tenant"
                      value={addNFTCollection.tenantName}
                    />

                    <DataDisplayRowText
                      label="Name"
                      value={addNFTCollection.name}
                    />

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <Link
                        class="button is-medium is-hidden-touch"
                        to="/admin/collections/add/step-2"
                      >
                        <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                        &nbsp;Back to Step 2
                      </Link>
                      <Link
                        class="button is-medium is-fullwidth is-hidden-desktop"
                        to="/admin/collections/add/step-2"
                      >
                        <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                        &nbsp;Back to Step 2
                      </Link>
                    </div>
                    <div class="column is-half has-text-right">
                      <button
                        class="button is-medium is-primary is-hidden-touch"
                        onClick={onSubmitClick}
                      >
                        <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                        &nbsp;Submit
                      </button>
                      <button
                        class="button is-medium is-primary is-fullwidth is-hidden-desktop"
                        onClick={onSubmitClick}
                      >
                        <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                        &nbsp;Submit
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

export default AdminNFTCollectionAddStep3;
