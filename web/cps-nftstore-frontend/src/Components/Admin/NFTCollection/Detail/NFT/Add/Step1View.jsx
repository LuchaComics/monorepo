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
  faDatabase,
  faEye,
  faArrowLeft,
  faArrowRight,
  faExclamationTriangle,
  faChain
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import useLocalStorage from "../../../../../../Hooks/useLocalStorage";
import { getCollectionDetailAPI } from "../../../../../../API/NFTCollection";
import FormErrorBox from "../../../../../Reusable/FormErrorBox";
import FormInputField from "../../../../../Reusable/FormInputField";
import FormTextareaField from "../../../../../Reusable/FormTextareaField";
import FormRadioField from "../../../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../../../Reusable/FormSelectField";
import FormCheckboxField from "../../../../../Reusable/FormCheckboxField";
import FormCountryField from "../../../../../Reusable/FormCountryField";
import FormRegionField from "../../../../../Reusable/FormRegionField";
import FormNFTAssetField from "../../../../../Reusable/FormNFTAssetField";
import PageLoadingContent from "../../../../../Reusable/PageLoadingContent";
import FormNFTMetadataAttributesField from "../../../../../Reusable/FormNFTMetadataAttributesField";
import {
  topAlertMessageState,
  topAlertStatusState,
  ADD_NFT_STATE_DEFAULT,
  addNFTState
} from "../../../../../../AppState";

function AdminNFTCollectionNFTAddStep1() {
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
  const [addNFT, setAddNFT] = useRecoilState(addNFTState);

  ////
  //// Component states.
  ////

  // Form GUI related.
  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [collection, setCollection] = useState({});
  const [showCancelWarning, setShowCancelWarning] = useState(false);

  // Form fields
  const [name, setName] = useState(addNFT.name);
  const [imageID, setImageID] = useState(addNFT.imageID);
  const [imageFilename, setImageFilename] = useState(addNFT.imageFilename);
  const [description, setDescription] = useState(addNFT.description);
  const [animationID, setAnimationID] = useState(addNFT.animationID);
  const [animationFilename, setAnimationFilename] = useState(addNFT.animationFilename);
  const [externalURL, setExternalURL] = useState(addNFT.externalURL);
  const [backgroundColor, setBackgroundColor] = useState(addNFT.backgroundColor);
  const [youtubeURL, setYoutubeURL] = useState(addNFT.youtubeURL);
  const [attributes, setAttributes] = useState(addNFT.attributes);

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

      if (name === undefined || name === null || name === "") {
        newErrors["name"] = "missing value";
        hasErrors = true;
      }
      if (imageID === undefined || imageID === null || imageID === "") {
        newErrors["imageId"] = "missing value";
        hasErrors = true;
      }
      // if (imageFilename === undefined || imageFilename === null || imageFilename === "") {
      //   newErrors["imageFilename"] = "missing value";
      //   hasErrors = true;
      // }
      if (description === undefined || description === null || description === "") {
        newErrors["description"] = "missing value";
        hasErrors = true;
      }
      // if (animationID === undefined || animationID === null || animationID === "") {
      //   newErrors["animationId"] = "missing value";
      //   hasErrors = true;
      // }
      // if (animationFilename === undefined || animationFilename === null || animationFilename === "") {
      //   newErrors["animationFilename"] = "missing value";
      //   hasErrors = true;
      // }
      // if (externalURL === undefined || externalURL === null || externalURL === "") {
      //   newErrors["externalURL"] = "missing value";
      //   hasErrors = true;
      // }
      if (backgroundColor === undefined || backgroundColor === null || backgroundColor === "") {
        newErrors["backgroundColor"] = "missing value";
        hasErrors = true;
      }
      // if (youtubeURL === undefined || youtubeURL === null || youtubeURL === "") {
      //   newErrors["youtubeURL"] = "missing value";
      //   hasErrors = true;
      // }
      // if (attributes === undefined || attributes === null || attributes === "") {
      //   newErrors["attributes"] = "missing value";
      //   hasErrors = true;
      // }

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
      let modifiedAddNFT = { ...addNFT };
      modifiedAddNFT.collectionID = id;
      modifiedAddNFT.collectionId = id; // Bugfix for js api call.
      modifiedAddNFT.name = name;
      modifiedAddNFT.imageID = imageID;
      modifiedAddNFT.imageFilename = imageFilename;
      modifiedAddNFT.description = description;
      modifiedAddNFT.animationID = animationID;
      modifiedAddNFT.animationFilename = animationFilename;
      modifiedAddNFT.externalURL = externalURL;
      modifiedAddNFT.backgroundColor = backgroundColor;
      modifiedAddNFT.youtubeURL = youtubeURL;
      modifiedAddNFT.attributes = attributes;
      setAddNFT(modifiedAddNFT);

      setForceURL("/admin/collection/" + id + "/nfts/add/step-2");
  };

  const onCancelClick = (e, url) => {
    e.preventDefault();
    setAddNFT(ADD_NFT_STATE_DEFAULT);
    setForceURL(url);
  }

  ////
  //// API.
  ////

  // --- Get Collection --- //

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
              <li class="">
                <Link to={`/admin/collection/${id}/nfts`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail (NFTs)
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
                <Link to={`/admin/collection/${id}/nfts`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Detail (NFTs)
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
                  Your NFT record will be cancelled and your work will be lost.
                  This cannot be undone. Do you want to continue?
                </section>
                <footer class="modal-card-foot">
                <Link
                  class="button is-medium is-success"
                  onClick={(e)=>{
                      onCancelClick(e, "/admin/collection/" + id + "/nfts");
                  }}
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

          {/* Progress Wizard */}
          <nav className="box has-background-light">
            <p className="subtitle is-5">Step 1 of 3</p>
            <progress
              class="progress is-success"
              value="33"
              max="100"
            >
              33%
            </progress>
          </nav>

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faPlus} />
              &nbsp;New NFT
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
                  <>
                      <p class="subtitle is-6">
                        <FontAwesomeIcon className="fas" icon={faDatabase} />
                        &nbsp;Metadata
                      </p>
                      <hr />

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
                        label="Animation (Optional)"
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
                        label="YouTube URL (Optional)"
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
                  </>

                  <br />
                  <br />

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <Link
                        onClick={(e) => setShowCancelWarning(true)}
                        class="button is-fullwidth-mobile is-medium"
                      >
                        <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                        &nbsp;Cancel
                      </Link>
                    </div>
                    <div class="column is-half has-text-right">
                      <button
                        class="button is-primary is-fullwidth-mobile is-medium"
                        onClick={onSubmitClick}
                      >

                        Save & Continue&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} />
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

export default AdminNFTCollectionNFTAddStep1;
