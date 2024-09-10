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
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import useLocalStorage from "../../../../../Hooks/useLocalStorage";
import { postNFTMetadataCreateAPI } from "../../../../../API/NFTMetadata";
import FormErrorBox from "../../../../Reusable/FormErrorBox";
import FormInputField from "../../../../Reusable/FormInputField";
import FormTextareaField from "../../../../Reusable/FormTextareaField";
import FormRadioField from "../../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../../Reusable/FormSelectField";
import FormCheckboxField from "../../../../Reusable/FormCheckboxField";
import FormCountryField from "../../../../Reusable/FormCountryField";
import FormRegionField from "../../../../Reusable/FormRegionField";
import FormNFTAssetField from "../../../../Reusable/FormNFTAssetField";
import PageLoadingContent from "../../../../Reusable/PageLoadingContent";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../../AppState";

function AdminCollectionNFTMetadataAdd() {
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

  // Form state.
  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");

  // Form fields
  const [name, setName] = useState("");
  const [imageID, setImageID] = useState("");
  const [imageFilename, setImageFilename] = useState("");
  const [description, setDescription] = useState("");
  const [animationID, setAnimationID] = useState("");
  const [animationFilename, setAnimationFilename] = useState("");
  const [externalURL, setExternalURL] = useState("");
  const [backgroundColor, setBackgroundColor] = useState("");
  const [youtubeURL, setYoutubeURL] = useState("");

  ////
  //// Event handling.
  ////



  const onSubmitClick = (e) => {
    console.log("onSubmitClick: Starting...");
    setFetching(true);
    setErrors({});

    const jsonData = {
        collection_id: id,
        name: name,
        image_id: imageID,
        description: description,
        animation_id: animationID,
        external_url: externalURL,
        background_color: backgroundColor,
        youtube_url: youtubeURL,
    };
    // formData.append("file", selectedFile);
    // formData.append("name", name);
    // formData.append("collection_id", id);

    postNFTMetadataCreateAPI(
      jsonData,
      onAdminCollectionNFTMetadataAddSuccess,
      onAdminCollectionNFTMetadataAddError,
      onAdminCollectionNFTMetadataAddDone,
      onUnauthorized,
    );
    console.log("onSubmitClick: Finished.");
  };

  ////
  //// API.
  ////

  function onAdminCollectionNFTMetadataAddSuccess(response) {
    // For debugging purposes only.
    console.log("onAdminCollectionNFTMetadataAddSuccess: Starting...");
    console.log(response);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("NFTMetadata created");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onAdminCollectionNFTMetadataAddSuccess: Delayed for 2 seconds.");
      console.log(
        "onAdminCollectionNFTMetadataAddSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Redirect the collection to the collection pinobjects page.
    setForceURL("/admin/collection/" + id + "/nft-metadatum/" + response.id);
  }

  function onAdminCollectionNFTMetadataAddError(apiErr) {
    console.log("onAdminCollectionNFTMetadataAddError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onAdminCollectionNFTMetadataAddError: Delayed for 2 seconds.");
      console.log(
        "onAdminCollectionNFTMetadataAddError: topAlertMessage, topAlertStatus:",
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

  function onAdminCollectionNFTMetadataAddDone() {
    console.log("onAdminCollectionNFTMetadataAddDone: Starting...");
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
                <Link to={`/admin/collection/${id}/nft-metadata`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Detail (NFT Metadata)
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modals */}
          {/* None */}

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faPlus} />
              &nbsp;New NFT Metadata
            </p>
            <FormErrorBox errors={errors} />

            {/* <p class="pb-4 has-text-grey">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Submitting..."} />
            ) : (
              <>
                <div class="container">
                  <FormInputField
                    label="Name"
                    name="name"
                    placeholder="Text input"
                    value={name}
                    errorText={errors && errors.name}
                    helpText="Optional"
                    onChange={(e) => setName(e.target.value)}
                    isRequired={true}
                    maxWidth="450px"
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

                  <br />
                  <br />

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <Link
                        to={`/admin/collection/${id}/nft-metadata`}
                        class="button is-hidden-touch"
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back
                      </Link>
                      <Link
                        to={`/admin/collection/${id}/nft-metadata`}
                        class="button is-fullwidth is-hidden-desktop"
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

export default AdminCollectionNFTMetadataAdd;
