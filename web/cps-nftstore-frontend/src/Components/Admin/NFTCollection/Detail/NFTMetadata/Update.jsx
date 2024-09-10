import React, { useState, useEffect } from "react";
import { Link, Navigate, useParams } from "react-router-dom";
import Scroll from "react-scroll";
import { decamelizeKeys } from "humps";
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
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import useLocalStorage from "../../../../../Hooks/useLocalStorage";
import {
  putNFTMetadataUpdateAPI,
  getNFTMetadataDetailAPI,
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
import FormNFTAssetField from "../../../../Reusable/FormNFTAssetField";
import FormRowText from "../../../../Reusable/FormRowText";
import FormNFTMetadataAttributesField from "../../../../Reusable/FormNFTMetadataAttributesField";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../../AppState";

function AdminNFTCollectionNFTMetadataUpdate() {
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

  // Form state.
  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");

  // Form fields
  const [tokenID, setTokenID] = useState(0);
  const [name, setName] = useState("");
  const [imageID, setImageID] = useState("");
  const [imageFilename, setImageFilename] = useState("");
  const [description, setDescription] = useState("");
  const [createdAt, setCreatedAt] = useState("");
  const [modifiedAt, setModifiedAt] = useState("");
  const [animationID, setAnimationID] = useState("");
  const [animationFilename, setAnimationFilename] = useState("");
  const [externalURL, setExternalURL] = useState("");
  const [backgroundColor, setBackgroundColor] = useState("");
  const [youtubeURL, setYoutubeURL] = useState("");
  const [fileCID, setFileCID] = useState("");
  const [fileIpnsPath, setFileIpnsPath] = useState("");
  const [status, setStatus] = useState("");
  const [attributes, setAttributes] = useState([]);

  ////
  //// Event handling.
  ////

  const onSubmitClick = (e) => {
      console.log("onSubmitClick: Starting...");
      setFetching(true);
      setErrors({});

      const jsonData = {
          id: rid,
          token_id: tokenID,
          collection_id: id,
          name: name,
          image_id: imageID,
          description: description,
          animation_id: animationID,
          external_url: externalURL,
          background_color: backgroundColor,
          youtube_url: youtubeURL,
          status: status,
          attributes: decamelizeKeys(attributes),
      };

    putNFTMetadataUpdateAPI(
      rid,
      jsonData,
      onAdminNFTCollectionNFTMetadataUpdateSuccess,
      onAdminNFTCollectionNFTMetadataUpdateError,
      onAdminNFTCollectionNFTMetadataUpdateDone,
      onUnauthorized,
    );
    console.log("onSubmitClick: Finished.");
  };

  ////
  //// API.
  ////

  function onAdminNFTCollectionNFTMetadataUpdateSuccess(response) {
    // For debugging purposes only.
    console.log("onAdminNFTCollectionNFTMetadataUpdateSuccess: Starting...");
    console.log(response);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("NFT Metadata updated");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onAdminNFTCollectionNFTMetadataUpdateSuccess: Delayed for 2 seconds.");
      console.log(
        "onAdminNFTCollectionNFTMetadataUpdateSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Redirect the collection to the collection NFT metadatum page.
    setForceURL("/admin/collection/" + id + "/nft-metadatum/" + rid);
  }

  function onAdminNFTCollectionNFTMetadataUpdateError(apiErr) {
    console.log("onAdminNFTCollectionNFTMetadataUpdateError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onAdminNFTCollectionNFTMetadataUpdateError: Delayed for 2 seconds.");
      console.log(
        "onAdminNFTCollectionNFTMetadataUpdateError: topAlertMessage, topAlertStatus:",
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

  function onAdminNFTCollectionNFTMetadataUpdateDone() {
    console.log("onAdminNFTCollectionNFTMetadataUpdateDone: Starting...");
    setFetching(false);
  }

  function onAdminNFTCollectionNFTMetadataDetailSuccess(response) {
      // For debugging purposes only.
      console.log("onAdminNFTCollectionNFTMetadataDetailSuccess: Starting...");
      console.log("onAdminNFTCollectionNFTMetadataDetailSuccess: response:", response);
      setTokenID(response.tokenId);
      setName(response.name);
      setImageID(response.imageId);
      setImageFilename(response.imageFilename);
      setDescription(response.description);
      setCreatedAt(response.createdAt);
      setModifiedAt(response.modifiedAt);
      setAnimationID(response.animationId);
      setAnimationFilename(response.animationFilename);
      setExternalURL(response.externalUrl);
      setBackgroundColor(response.backgroundColor);
      setYoutubeURL(response.youtubeUrl);
      setFileCID(response.fileCid);
      setFileIpnsPath(response.fileIpnsPath);
      setStatus(response.status);
      if (response.attributes !== undefined && response.attributes !== null && response.attributes === "") {
        setAttributes([]);
      } else {
        setAttributes(response.attributes);
      }
  }

  function onAdminNFTCollectionNFTMetadataDetailError(apiErr) {
    console.log("onAdminNFTCollectionNFTMetadataDetailError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onAdminNFTCollectionNFTMetadataDetailError: Delayed for 2 seconds.");
      console.log(
        "onAdminNFTCollectionNFTMetadataDetailError: topAlertMessage, topAlertStatus:",
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

  function onAdminNFTCollectionNFTMetadataDetailDone() {
    console.log("onAdminNFTCollectionNFTMetadataDetailDone: Starting...");
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
        onAdminNFTCollectionNFTMetadataDetailSuccess,
        onAdminNFTCollectionNFTMetadataDetailError,
        onAdminNFTCollectionNFTMetadataDetailDone,
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
                <Link to={`/admin/collection/${id}/nft-metadata`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail (NFT Metadata)
                </Link>
              </li>
              <li class="">
                <Link
                  to={`/admin/collection/${id}/nft-metadatum/${rid}`}
                  aria-current="page"
                >
                  <FontAwesomeIcon className="fas" icon={faCube} />
                  &nbsp;NFT Metadata
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faPencil} />
                  &nbsp;Edit
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link
                  to={`/admin/collection/${id}/pin/${rid}`}
                  aria-current="page"
                >
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to NFT Metadata
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modals */}
          {/* None */}

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faPencil} />
              &nbsp;Edit NFT Metadata
            </p>

            {/* <p class="pb-4 has-text-grey">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Submitting..."} />
            ) : (
              <>
                <FormErrorBox errors={errors} />
                <div class="container">
                  <article class="message is-warning">
                    <div class="message-body">
                      <strong>Warning:</strong> Changing <b>Token ID, Image</b> and <b>Animation</b> fields has been disabled.
                    </div>
                  </article>

                  <FormRowText label="Token ID" value={tokenID} helpText="" type="text" />

                  <FormInputField
                    label="Name"
                    name="name"
                    placeholder="Text input"
                    value={name}
                    errorText={errors && errors.name}
                    helpText=""
                    onChange={(e) => setName(e.target.value)}
                    isRequired={true}
                    maxWidth="150px"
                  />

                  <FormTextareaField
                    label="Description"
                    name="description"
                    placeholder="Text input"
                    value={description}
                    errorText={errors && errors.description}
                    helpText=""
                    onChange={(e) => setDescription(e.target.value)}
                    isRequired={true}
                    maxWidth="150px"
                    rows={4}
                  />

                <FormRowText label="Image" value={imageFilename} helpText="" type="text" />
                <FormRowText label="Animation" value={animationFilename} helpText="" type="text" />
{/*
                  <FormNFTAssetField
                    label="Image"
                    name="imageId"
                    filename={imageFilename}
                    setFilename={setImageFilename}
                    nftAssetID={imageID}
                    setNFTAssetID={setImageID}
                    helpText={`Upload the image for this NFT. This should be the submission that was reviewed by CPS.`}
                    errorText={errors && errors.imageId}
                  />


                  <FormNFTAssetField
                    label="Animation"
                    name="animationId"
                    filename={animationFilename}
                    setFilename={setAnimationFilename}
                    nftAssetID={animationID}
                    setNFTAssetID={setAnimationID}
                    helpText={`Upload the submission review video for this NFT. This should be the submission that was reviewed by CPS.`}
                    errorText={errors && errors.animationId}
                  />
*/}
                  <FormInputField
                    label="Background Color"
                    name="backgroundColor"
                    placeholder="Text input"
                    value={backgroundColor}
                    errorText={errors && errors.backgroundColor}
                    helpText="Please use hexadecimal values"
                    onChange={(e) => setBackgroundColor(e.target.value)}
                    isRequired={true}
                    maxWidth="150px"
                  />

                  <FormInputField
                    label="External URL (Optional)"
                    name="externalURL"
                    placeholder="Text input"
                    value={externalURL}
                    errorText={errors && errors.externalURL}
                    helpText={<>
                        <p>If you do not fill this then system will set its own value.</p>
                    </>}
                    onChange={(e) => setExternalURL(e.target.value)}
                    isRequired={true}
                    maxWidth="250px"
                  />

                  <FormInputField
                    label="YouTube URL"
                    name="youtubeURL"
                    placeholder="Text input"
                    value={youtubeURL}
                    errorText={errors && errors.youtubeUrl}
                    helpText=""
                    onChange={(e) => setYoutubeURL(e.target.value)}
                    isRequired={true}
                    maxWidth="275px"
                  />

                  <FormNFTMetadataAttributesField
                    data={attributes}
                    onDataChange={setAttributes}
                  />

                  <br />
                  <br />

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <Link
                        to={`/admin/collection/${id}/nft-metadatum/${rid}`}
                        class="button is-hidden-touch"
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back
                      </Link>
                      <Link
                        to={`/admin/collection/${id}/nft-metadatum/${rid}`}
                        class="button is-fullwidth is-hidden-desktop"
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back
                      </Link>
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

export default AdminNFTCollectionNFTMetadataUpdate;
