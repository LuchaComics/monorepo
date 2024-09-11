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
  faCube,
  faFile,
  faDownload,
  faArrowUpRightFromSquare,
  faBox,
  faChain,
  faCertificate
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import {
    getNFTDetailAPI,
    deleteNFTAPI,
    getNFTContentDetailAPI
} from "../../../../../API/NFT";
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
import FormNFTMetadataAttributesField from "../../../../Reusable/FormNFTMetadataAttributesField";


function AdminNFTCollectionNFTDetail() {
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
  const [imageCID, setImageCID] = useState("");
  const [description, setDescription] = useState("");
  const [createdAt, setCreatedAt] = useState("");
  const [modifiedAt, setModifiedAt] = useState("");
  const [animationURL, setAnimationURL] = useState("");
  const [animationCID, setAnimationCID] = useState("");
  const [externalURL, setExternalURL] = useState("");
  const [backgroundColor, setBackgroundColor] = useState("");
  const [youtubeURL, setYoutubeURL] = useState("");
  const [fileCID, setFileCID] = useState("");
  const [fileIpnsPath, setFileIpnsPath] = useState("");
  const [attributes, setAttributes] = useState([]);
  const [selectedNFTRequestIDForDeletion, setSelectedNFTRequestIDForDeletion] =
    useState("");

  ////
  //// Event handling.
  ////

  const onArchiveConfirmButtonClick = (e) => {
    e.preventDefault();
    console.log("onArchiveConfirmButtonClick"); // For debugging purposes only.

    deleteNFTAPI(
      selectedNFTRequestIDForDeletion,
      onNFTArchiveSuccess,
      onNFTArchiveError,
      onNFTArchiveDone,
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
      getNFTContentDetailAPI(
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

  function onAdminNFTCollectionNFTDetailSuccess(response) {
    // For debugging purposes only.
    console.log("onAdminNFTCollectionNFTDetailSuccess: Starting...");
    console.log("onAdminNFTCollectionNFTDetailSuccess: response:", response);
    setTokenID(response.tokenId);
    setName(response.name);
    setImage(response.image);
    setImageCID(response.imageCid);
    setDescription(response.description);
    setCreatedAt(response.createdAt);
    setModifiedAt(response.modifiedAt);
    setAnimationURL(response.animationUrl);
    setAnimationCID(response.animationCid);
    setExternalURL(response.externalUrl);
    setBackgroundColor(response.backgroundColor);
    setYoutubeURL(response.youtubeUrl);
    setFileCID(response.fileCid);
    setFileIpnsPath(response.fileIpnsPath);
    setAttributes(response.attributes);
  }

  function onAdminNFTCollectionNFTDetailError(apiErr) {
    console.log("onAdminNFTCollectionNFTDetailError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onAdminNFTCollectionNFTDetailError: Delayed for 2 seconds.");
      console.log(
        "onAdminNFTCollectionNFTDetailError: topAlertMessage, topAlertStatus:",
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

  function onAdminNFTCollectionNFTDetailDone() {
    console.log("onAdminNFTCollectionNFTDetailDone: Starting...");
    setFetching(false);
  }

  // --- Deletion --- //

  function onNFTArchiveSuccess(response) {
    console.log("onNFTArchiveSuccess: Starting..."); // For debugging purposes only.

    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("NFT archived");
    setTimeout(() => {
      console.log(
        "onArchiveConfirmButtonClick: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Redirect the collection to the collection NFT metadata page.
    setForceURL("/admin/collection/" + id + "/nfts");
  }

  function onNFTArchiveError(apiErr) {
    console.log("onNFTArchiveError: Starting..."); // For debugging purposes only.
    setErrors(apiErr);

    // Update notification.
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed deleting");
    setTimeout(() => {
      console.log(
        "onNFTArchiveError: topAlertMessage, topAlertStatus:",
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

  function onNFTArchiveDone() {
    console.log("onNFTArchiveDone: Starting...");
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

      getNFTDetailAPI(
        rid,
        onAdminNFTCollectionNFTDetailSuccess,
        onAdminNFTCollectionNFTDetailError,
        onAdminNFTCollectionNFTDetailDone,
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
                  &nbsp;NFT Collections
                </Link>
              </li>
              <li class="">
                <Link to={`/admin/collection/${id}/nfts`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail (NFTs)
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faCube} />
                  &nbsp;NFT
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/admin/collection/${id}/nfts`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Detail (NFTs)
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modals */}
          <div
            class={`modal ${selectedNFTRequestIDForDeletion ? "is-active" : ""}`}
          >
            <div class="modal-background"></div>
            <div class="modal-card">
              <header class="modal-card-head">
                <p class="modal-card-title">Are you sure?</p>
                <button
                  class="archive"
                  aria-label="close"
                  onClick={(e)=>{
                      e.preventDefault();
                      setSelectedNFTRequestIDForDeletion("");
                  }}
                ></button>
              </header>
              <section class="modal-card-body">
                You are about to <b>archive</b> this MFT metadata; it will no
                longer appear on your dashboard. This action cannot be undone. Are you sure
                you would like to continue?
              </section>
              <footer class="modal-card-foot">
                <button
                  class="button is-success"
                  onClick={onArchiveConfirmButtonClick}
                >
                  Confirm
                </button>
                <button
                  class="button"
                  onClick={(e)=>{
                      e.preventDefault();
                      setSelectedNFTRequestIDForDeletion("");
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
              <FontAwesomeIcon className="fas" icon={faCube} />
              &nbsp;NFT
            </p>

            {/* <p class="pb-4 has-text-grey">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Submitting..."} />
            ) : (
              <>
                <FormErrorBox errors={errors} />
                <div class="container">
                  <p class="subtitle is-4 pt-4">
                     <FontAwesomeIcon className="fas" icon={faChain} />
                     &nbsp;Blockchain
                  </p>
                  <hr />

                  <FormRowText label="Blockchain Membership" value={`Ethereum`} helpText="" type="text" />
                  <FormRowText label="Smart Contract" value="Collectible Protection Service Submission Token" helpText="" type="text" />
                  <FormRowText label="Token ID" value={tokenID} helpText="" type="text" />

                  <p class="subtitle is-4 pt-4">
                    <FontAwesomeIcon className="fas" icon={faEye} />
                    &nbsp;Meta Information
                  </p>
                  <hr />

                  <FormRowText label="Name" value={name} helpText="" />
                  <FormRowText label="Description" value={description} helpText="" />
                  <FormRowText label="Image" value={image} helpText={
                  <>
                      View file in <Link target="_blank" rel="noreferrer" to={`${process.env.REACT_APP_API_PROTOCOL}://${process.env.REACT_APP_API_DOMAIN}/ipfs/${imageCID}`}>local gateway&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} /></Link>&nbsp;or&nbsp;
                      <Link target="_blank" rel="noreferrer" to={`https://ipfs.io/ipfs/${imageCID}`}>public gateway&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} /></Link>.
                  </>} />
                  <FormRowText label="Animation URL" value={animationURL} helpText={
                  <>
                      View file in  <Link target="_blank" rel="noreferrer" to={`${process.env.REACT_APP_API_PROTOCOL}://${process.env.REACT_APP_API_DOMAIN}/ipfs/${animationCID}`}>local gateway&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} /></Link>&nbsp;or&nbsp;
                      <Link target="_blank" rel="noreferrer" to={`https://ipfs.io/ipfs/${animationCID}`}>public gateway&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} /></Link>.
                  </>} />
                  <FormRowText label="Background Color" value={backgroundColor} helpText="" />
                  <FormRowText label="YouTube URL" value={youtubeURL} helpText="" />
                  <FormTextDateTimeRow label="Created At" value={createdAt} helpText="" />
                  <FormTextDateTimeRow label="Modified At" value={modifiedAt} helpText="" />
                  <FormNFTMetadataAttributesField
                      data={attributes}
                      onDataChange={setAttributes}
                      disabled={true}
                   />

                  <p class="subtitle is-4 pt-4">
                    <FontAwesomeIcon className="fas" icon={faFile} />
                    &nbsp;File Information
                  </p>
                  <hr />

                  <FormRowText label="Metadata File CID" value={`/ipfs/${fileCID}`} helpText={
                  <>
                      View file in <Link target="_blank" rel="noreferrer" to={`${process.env.REACT_APP_API_PROTOCOL}://${process.env.REACT_APP_API_DOMAIN}/ipfs/${fileCID}`}>local gateway&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} /></Link>&nbsp;or&nbsp;
                      <Link target="_blank" rel="noreferrer" to={`https://ipfs.io/ipfs/${fileCID}`}>public gateway&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} /></Link>.
                  </>} />
                  <FormRowText label="Metadata File IPNS Path" value={fileIpnsPath} helpText={
                  <>
                      View file in <Link target="_blank" rel="noreferrer" to={`${process.env.REACT_APP_API_PROTOCOL}://${process.env.REACT_APP_API_DOMAIN}${fileIpnsPath}`}>local gateway&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} /></Link>&nbsp;or&nbsp;
                      <Link target="_blank" rel="noreferrer" to={`https://ipfs.io${fileIpnsPath}`}>public gateway&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} /></Link>.
                  </>} />

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <Link
                        to={`/admin/collection/${id}/nfts`}
                        class="button is-medium is-fullwidth-mobile"
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back
                      </Link>
                    </div>
                    <div class="column is-half has-text-right">
                      <Link
                        to={`/admin/collection/${id}/nft/${rid}/mint`}
                        class="button is-medium is-success is-fullwidth-mobile"
                      >
                        <FontAwesomeIcon className="fas" icon={faCertificate} />
                        &nbsp;Mint
                      </Link>
                      &nbsp;
                      &nbsp;
                      <Link
                        to={`/admin/collection/${id}/nft/${rid}/edit`}
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
                              setSelectedNFTRequestIDForDeletion(rid);
                          }}
                          class="button is-medium is-danger is-fullwidth-mobile"
                        >
                          <FontAwesomeIcon className="fas" icon={faBox} />
                          &nbsp;Archive
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

export default AdminNFTCollectionNFTDetail;
