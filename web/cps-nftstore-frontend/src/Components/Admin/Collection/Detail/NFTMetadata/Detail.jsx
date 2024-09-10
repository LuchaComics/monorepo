import React, { useState, useEffect } from "react";
import { Link, Navigate, useParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
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
  faEye,
  faArrowLeft,
  faLocationPin,
  faFile,
  faDownload,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import {
    getNFTMetadataDetailAPI,
    deleteNFTMetadataAPI,
    getNFTMetadataContentDetailAPI
} from "../../../../../API/NFTMetadata";
import FormErrorBox from "../../../../Reusable/FormErrorBox";
import FormInputField from "../../../../Reusable/FormInputField";
import FormTextareaField from "../../../../Reusable/FormTextareaField";
import FormRadioField from "../../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../../Reusable/FormSelectField";
import FormCheckboxField from "../../../../Reusable/FormCheckboxField";
import FormCountryField from "../../../../Reusable/FormCountryField";
import FormRegionField from "../../../../Reusable/FormRegionField";
import PageLoadingContent from "../../../../Reusable/PageLoadingContent";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../../AppState";
import FormRowText from "../../../../Reusable/FormRowText";
import FormTextDateTimeRow from "../../../../Reusable/FormRowTextDateTime";


function AdminCollectionNFTMetadataDetail() {
  ////
  //// URL Parameters.
  ////

  const { id, rid } = useParams();

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
  const [selectedFile, setSelectedFile] = useState(null);
  const [tokenID, setTokenID] = useState(0);
  const [name, setName] = useState("");
  const [image, setImage] = useState("");
  const [description, setDescription] = useState("");
  const [createdAt, setCreatedAt] = useState("");
  const [modifiedAt, setModifiedAt] = useState("");
  const [animationURL, setAnimationURL] = useState("");
  const [externalURL, setExternalURL] = useState("");
  const [backgroundColor, setBackgroundColor] = useState("");
  const [youtubeURL, setYoutubeURL] = useState("");
  const [ipnsPath, setIpnsPath] = useState("");
  const [selectedNFTMetadataRequestIDForDeletion, setSelectedNFTMetadataRequestIDForDeletion] =
    useState("");

  ////
  //// Event handling.
  ////

  const onDeleteConfirmButtonClick = (e) => {
    e.preventDefault();
    console.log("onDeleteConfirmButtonClick"); // For debugging purposes only.

    deleteNFTMetadataAPI(
      selectedNFTMetadataRequestIDForDeletion,
      onNFTMetadataDeleteSuccess,
      onNFTMetadataDeleteError,
      onNFTMetadataDeleteDone,
      onUnauthorized,
    );
  };

  const downloadFile = (data, filename, contentType) => {
    const blob = new Blob([data], { type: contentType });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.style.display = 'none';
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
};

  const onDownloadContentButtonClick = (e) => {
      e.preventDefault();
      console.log("onDownloadContentButtonClick: Started...");
      getNFTMetadataContentDetailAPI(
          rid,
          (data, filename, contentType) => {
              // // DEFENSIVE CODE: In case `filename` was not returned.
              // if (filename === undefined || filename === null || filename === "") {
              //     filename = meta["filename"];
              //     console.log("onDownloadContentButtonClick: `filename` not found, using meta:", filename);
              // }
              //
              //     // DEFENSIVE CODE: In case `contentType` was not returned.
              // if (contentType === undefined || contentType === null || contentType === "") {
              //     contentType = meta["contentType"];
              //     console.log("onDownloadContentButtonClick: `contentType` not found, using meta:", contentType);
              // }
              //
              // // Download the file.
              // console.log("onDownloadContentButtonClick: success:", data, filename, contentType);
              // downloadFile(data, filename, contentType);
          },
          (apiErr) => {
              console.log("onDownloadContentButtonClick: err:", apiErr);
              setErrors(apiErr);

              // Update notification.
              setTopAlertStatus("danger");
              setTopAlertMessage("Failed downloading file");
              setTimeout(() => {
                console.log(
                  "onDownloadContentButtonClick: topAlertMessage, topAlertStatus:",
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
          },
          () => {
              console.log("onDownloadContentButtonClick: done");
          },
          onUnauthorized
      );
  }

  ////
  //// API.
  ////

  // --- Get Details ---

  function onAdminCollectionNFTMetadataDetailSuccess(response) {
    // For debugging purposes only.
    console.log("onAdminCollectionNFTMetadataDetailSuccess: Starting...");
    console.log("onAdminCollectionNFTMetadataDetailSuccess: response:", response);
    setTokenID(response.tokenId);
    setName(response.name);
    setImage(response.image);
    setDescription(response.description);
    setCreatedAt(response.createdAt);
    setModifiedAt(response.modifiedAt);
    setAnimationURL(response.animationUrl);
    setExternalURL(response.externalUrl);
    setBackgroundColor(response.backgroundColor);
    setYoutubeURL(response.youtubeUrl);
    setIpnsPath(response.ipnsPath);
  }

  function onAdminCollectionNFTMetadataDetailError(apiErr) {
    console.log("onAdminCollectionNFTMetadataDetailError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onAdminCollectionNFTMetadataDetailError: Delayed for 2 seconds.");
      console.log(
        "onAdminCollectionNFTMetadataDetailError: topAlertMessage, topAlertStatus:",
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

  function onAdminCollectionNFTMetadataDetailDone() {
    console.log("onAdminCollectionNFTMetadataDetailDone: Starting...");
    setFetching(false);
  }

  // --- Deletion --- //

  function onNFTMetadataDeleteSuccess(response) {
    console.log("onNFTMetadataDeleteSuccess: Starting..."); // For debugging purposes only.

    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("NFT Metadata deleted");
    setTimeout(() => {
      console.log(
        "onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Redirect the collection to the collection NFT metadata page.
    setForceURL("/admin/collection/" + id + "/nft-metadata");
  }

  function onNFTMetadataDeleteError(apiErr) {
    console.log("onNFTMetadataDeleteError: Starting..."); // For debugging purposes only.
    setErrors(apiErr);

    // Update notification.
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed deleting");
    setTimeout(() => {
      console.log(
        "onNFTMetadataDeleteError: topAlertMessage, topAlertStatus:",
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

  function onNFTMetadataDeleteDone() {
    console.log("onNFTMetadataDeleteDone: Starting...");
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

      getNFTMetadataDetailAPI(
        rid,
        onAdminCollectionNFTMetadataDetailSuccess,
        onAdminCollectionNFTMetadataDetailError,
        onAdminCollectionNFTMetadataDetailDone,
        onUnauthorized,
      );
    }

    return () => {
      mounted = false;
    };
  }, [rid]);

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
                  &nbsp;Collections
                </Link>
              </li>
              <li class="">
                <Link to={`/admin/collection/${id}/nft-metadata`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail (NFT Metadata)
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faLocationPin} />
                  &nbsp;NFT Metadata
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/admin/collection/${id}/nft-metadata`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Detail (NFT Metadata)
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modals */}
          <div
            class={`modal ${selectedNFTMetadataRequestIDForDeletion ? "is-active" : ""}`}
          >
            <div class="modal-background"></div>
            <div class="modal-card">
              <header class="modal-card-head">
                <p class="modal-card-title">Are you sure?</p>
                <button
                  class="delete"
                  aria-label="close"
                  onClick={(e)=>{
                      e.preventDefault();
                      setSelectedNFTMetadataRequestIDForDeletion("");
                  }}
                ></button>
              </header>
              <section class="modal-card-body">
                You are about to <b>delete</b> this MFT metadata; it will no
                longer appear on your dashboard This action cannot be undone. Are you sure
                you would like to continue?
              </section>
              <footer class="modal-card-foot">
                <button
                  class="button is-success"
                  onClick={onDeleteConfirmButtonClick}
                >
                  Confirm
                </button>
                <button
                  class="button"
                  onClick={(e)=>{
                      e.preventDefault();
                      setSelectedNFTMetadataRequestIDForDeletion("");
                  }}
                >
                  Cancel
                </button>
              </footer>
            </div>
          </div>

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faLocationPin} />
              &nbsp;Pin
            </p>

            {/* <p class="pb-4 has-text-grey">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Submitting..."} />
            ) : (
              <>
                <FormErrorBox errors={errors} />
                <div class="container">
                  <p class="subtitle is-4 pt-4">
                    <FontAwesomeIcon className="fas" icon={faEye} />
                    &nbsp;Meta Information
                  </p>
                  <hr />

                  <FormRowText label="Token ID" value={tokenID} helpText="" type="text" />
                  <FormRowText label="Name" value={name} helpText="" />
                  <FormRowText label="Description" value={description} helpText="" />
                  <FormRowText label="Image" value={image} helpText="" />
                  <FormTextDateTimeRow label="Created At" value={createdAt} helpText="" />
                  <FormTextDateTimeRow label="Modified At" value={modifiedAt} helpText="" />
                  <FormRowText label="Animation URL" value={animationURL} helpText="" />
                  <FormRowText label="Background Color" value={backgroundColor} helpText="" />
                  <FormRowText label="YouTube URL" value={youtubeURL} helpText="" />
                  <FormRowText label="IPNS Path" value={ipnsPath} helpText="" />

                  <p class="subtitle is-4 pt-4">
                    <FontAwesomeIcon className="fas" icon={faFile} />
                    &nbsp;Data
                  </p>
                  <hr />
                  <p class="has-text-grey pb-4">
                    Click the following "Download File" button to start
                    downloading a copy of the NFT metadata to your computer.
                  </p>

                  <section class="hero has-background-white-ter">
                    <div class="hero-body">
                      <p class="subtitle">
                        <div class="has-text-centered">
                          <a
                            onClick={onDownloadContentButtonClick}
                            class="button is-large is-success is-hidden-touch"
                          >
                            <FontAwesomeIcon
                              className="fas"
                              icon={faDownload}
                            />
                            &nbsp;Download File
                          </a>
                          <a
                            onClick={onDownloadContentButtonClick}
                            rel="noreferrer"
                            class="button is-large is-success is-fullwidth is-hidden-desktop"
                          >
                            <FontAwesomeIcon
                              className="fas"
                              icon={faDownload}
                            />
                            &nbsp;Download File
                          </a>
                        </div>
                      </p>
                    </div>
                  </section>

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <Link
                        to={`/admin/collection/${id}/nft-metadata`}
                        class="button is-medium is-fullwidth-mobile"
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back
                      </Link>
                    </div>
                    <div class="column is-half has-text-right">
                      <Link
                        to={`/admin/collection/${id}/nft-metadatum/${rid}/edit`}
                        class="button is-medium is-warning is-fullwidth-mobile"
                      >
                        <FontAwesomeIcon className="fas" icon={faPencil} />
                        &nbsp;Edit
                      </Link>
                      &nbsp;
                      &nbsp;
                      <button
                          onClick={(e)=>{
                              e.preventDefault();
                              setSelectedNFTMetadataRequestIDForDeletion(rid);
                          }}
                          class="button is-medium is-danger is-fullwidth-mobile"
                        >
                          <FontAwesomeIcon className="fas" icon={faPencil} />
                          &nbsp;Delete
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

export default AdminCollectionNFTMetadataDetail;
