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
  faChain,
  faBoxOpen,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { postNFTCollectionRestoreOperationAPI } from "../../../../API/NFTCollection";
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
import DataDisplayRowCheckbox from "../../../Reusable/DataDisplayRowCheckbox";
import DataDisplayRowTenant from "../../../Reusable/DataDisplayRowTenant";
import {
    topAlertMessageState,
    topAlertStatusState,
    addNFTCollectionState,
    ADD_NFT_COLLECTION_STATE_DEFAULT
} from "../../../../AppState";


function AdminNFTCollectionAddViaBackupfile() {
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
  const [filename, setFilename] = useState("");
  const [selectedFile, setSelectedFile] = useState("");

  ////
  //// Event handling.
  ////

  const onDeleteClick = (e) => {

  }

  const onHandleFileChange = (event) => {
      console.log("onHandleFileChange: Starting...")
      // setFetching(true);
      setErrors({});


      const sf = event.target.files[0];
      setSelectedFile(sf);
  };

  ////
  //// API.
  ////

  const onSubmitClick = (event) => {
    console.log("onSubmitClick: Beginning...");
    event.preventDefault();
    setFetching(true);
    setErrors({});


    const formData = new FormData();
    formData.append('file', selectedFile);

    // Extract filename.
    const filename = selectedFile.name;
    console.log('Filename:', filename);

    const mimeType = selectedFile.type || "application/octet-stream"; // Fallback to octet-stream if type is not detected
    console.log('MIME Type:', mimeType);


    // Convert the FormData object to binary data and pass it directly
     // Note: You may need to adjust the handling depending on the type of file and how it needs to be processed
     const reader = new FileReader();
     reader.onload = () => {
       const fileBinaryData = reader.result;

       postNFTCollectionRestoreOperationAPI(
           filename,
           mimeType, // Pass the detected MIME type
           fileBinaryData, // Pass binary data instead of FormData
           onAdminNFTCollectionAddViaBackupfileSuccess,
           onAdminNFTCollectionAddViaBackupfileError,
           onAdminNFTCollectionAddViaBackupfileDone,
           onUnauthorized
       );

       console.log("onSubmitClick: Finished.");
       setFetching(false); // Reset fetching state after API call
     };
     reader.onerror = () => {
       console.error('Error reading file');
       setErrors({ file: 'Error reading file' });
       setFetching(false);
     };

     reader.readAsArrayBuffer(selectedFile); // Read file as binary data
    console.log("onSubmitClick: Finished.");

  };

  function onAdminNFTCollectionAddViaBackupfileSuccess(response) {
    // For debugging purposes only.
    console.log("onAdminNFTCollectionAddViaBackupfileSuccess: Starting...");
    console.log(response);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("NFT collection created");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onAdminNFTCollectionAddViaBackupfileSuccess: Delayed for 2 seconds.");
      console.log(
        "onAdminNFTCollectionAddViaBackupfileSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    const apiKey = response.apiKey;

    setForceURL("/admin/collection/" + response.id);
  }

  function onAdminNFTCollectionAddViaBackupfileError(apiErr) {
    console.log("onAdminNFTCollectionAddViaBackupfileError: Starting...");
    console.log("onAdminNFTCollectionAddViaBackupfileError: apiErr:", apiErr);
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      // console.log("onAdminNFTCollectionAddViaBackupfileError: Delayed for 2 seconds.");
      // console.log("onAdminNFTCollectionAddViaBackupfileError: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
      setTopAlertMessage("");
    }, 2000);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onAdminNFTCollectionAddViaBackupfileDone() {
    console.log("onAdminNFTCollectionAddViaBackupfileDone: Starting...");
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
                      <FontAwesomeIcon className="fas" icon={faBoxOpen} />
                      &nbsp;Restore
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

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faBoxOpen} />
              &nbsp;Restore a NFT Collection
            </p>
            <FormErrorBox errors={errors} />
            <p className="has-text-grey pb-4">
               Please upload backupfile and hit <i>Submit & Restore</i> to restore your NFT colleciton.
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
                  <input
                          name={`backupfile`}
                          type={"file"}
                   placeholder={`Please upload your backup file here`}
                      onChange={onHandleFileChange}
                  autoComplete="off" />




                  <div class="columns pt-5">
                    <div class="column is-half">
                      <Link
                        class="button is-medium is-fullwidth-mobile"
                        to="/admin/collections"
                      >
                        <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                        &nbsp;Cancel
                      </Link>
                    </div>
                    <div class="column is-half has-text-right">
                      <button
                        class="button is-medium is-primary is-fullwidth-mobile"
                        onClick={onSubmitClick}
                      >
                        <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                        &nbsp;Submit & Restore
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

export default AdminNFTCollectionAddViaBackupfile;
