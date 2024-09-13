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
  faExclamationTriangle,
  faChain,
  faCertificate
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import useLocalStorage from "../../../../../../Hooks/useLocalStorage";
import { postNFTCreateAPI } from "../../../../../../API/NFT";
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
import DataDisplayRowText from "../../../../../Reusable/DataDisplayRowText";
import DataDisplayRowRadio from "../../../../../Reusable/DataDisplayRowRadio";
import DataDisplayRowCheckbox from "../../../../../Reusable/DataDisplayRowCheckbox";
import DataDisplayRowTenant from "../../../../../Reusable/DataDisplayRowTenant";
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

  // Form fields

  ////
  //// Event handling.
  ////

  const onSubmitClick = (e) => {
    console.log("onSubmitClick: Starting...");
    console.log("onSubmitClick: addNFT:", addNFT);
    setFetching(true);
    setErrors({});
    const jsonData = {
        collection_id: addNFT.collectionID,
        name: addNFT.name,
        image_id: addNFT.imageID,
        description: addNFT.description,
        animation_id: addNFT.animationID,
        external_url: addNFT.externalURL,
        background_color: addNFT.backgroundColor,
        youtube_url: addNFT.youtubeURL,
        attributes: decamelizeKeys(addNFT.attributes),
        to_address: addNFT.toAddress,
        wallet_password: addNFT.walletPassword,
    };
    postNFTCreateAPI(
      jsonData,
      onAdminNFTCollectionNFTAddStep1Success,
      onAdminNFTCollectionNFTAddStep1Error,
      onAdminNFTCollectionNFTAddStep1Done,
      onUnauthorized,
    );
    console.log("onSubmitClick: Finished.");
  };

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

  // --- Post NFT --- //

  function onAdminNFTCollectionNFTAddStep1Success(response) {
    // For debugging purposes only.
    console.log("onAdminNFTCollectionNFTAddStep1Success: Starting...");
    console.log(response);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("NFT created");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onAdminNFTCollectionNFTAddStep1Success: Delayed for 4 seconds.");
      console.log(
        "onAdminNFTCollectionNFTAddStep1Success: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 4000);

    // Reset our submission form.
    setAddNFT(ADD_NFT_STATE_DEFAULT);

    // Redirect the collection to the collection pinobjects page.
    setForceURL("/admin/collection/" + id + "/nft/" + response.id);
  }

  function onAdminNFTCollectionNFTAddStep1Error(apiErr) {
    console.log("onAdminNFTCollectionNFTAddStep1Error: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onAdminNFTCollectionNFTAddStep1Error: Delayed for 2 seconds.");
      console.log(
        "onAdminNFTCollectionNFTAddStep1Error: topAlertMessage, topAlertStatus:",
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

  function onAdminNFTCollectionNFTAddStep1Done() {
    console.log("onAdminNFTCollectionNFTAddStep1Done: Starting...");
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
          {/* None */}

          {/* Progress Wizard */}
          <nav className="box has-background-success-light">
            <p className="subtitle is-5">Step 3 of 3</p>
            <progress
              class="progress is-success"
              value="100"
              max="100"
            >
              100%
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
                  <p className="title is-5 mt-2">
                    <FontAwesomeIcon className="fas" icon={faDatabase} />
                    &nbsp;Metadata&nbsp;&nbsp;&nbsp;-&nbsp;&nbsp;&nbsp;
                    <Link to={`/admin/collection/${id}/nfts/add/step-1`}>
                      <FontAwesomeIcon className="fas" icon={faPencil} />
                      &nbsp;Edit
                    </Link>
                  </p>

                  <DataDisplayRowText
                    label="Name"
                    value={addNFT.name}
                  />

                  <DataDisplayRowText
                    label="Description"
                    value={addNFT.description}
                  />

                  <DataDisplayRowText
                    label="Image"
                    value={addNFT.imageFilename}
                  />

                  <DataDisplayRowText
                    label="Animation (Optional)"
                    value={addNFT.animationFilename}
                  />

                  <DataDisplayRowText
                    label="Background Color"
                    value={addNFT.backgroundColor}
                  />

                  <DataDisplayRowText
                    label="External URL (Optional)"
                    value={addNFT.externalURL}
                  />

                  <DataDisplayRowText
                    label="YouTube URL (Optional)"
                    value={addNFT.youtubeURL}
                  />

                  <FormNFTMetadataAttributesField
                    data={addNFT.attributes}
                    onDataChange={null}
                    disabled={true}
                  />


                  <p className="title is-5 mt-2">
                    <FontAwesomeIcon className="fas" icon={faChain} />
                    &nbsp;Blockchain&nbsp;&nbsp;&nbsp;-&nbsp;&nbsp;&nbsp;
                    <Link to={`/admin/collection/${id}/nfts/add/step-2`}>
                      <FontAwesomeIcon className="fas" icon={faPencil} />
                      &nbsp;Edit
                    </Link>
                  </p>

                  <DataDisplayRowRadio
                    label="Blockchain"
                    value={collection.blockchain}
                    opt1Value="ethereum"
                    opt1Label="Ethereum"
                  />

                  <DataDisplayRowText
                    label="Node URL"
                    value={collection.nodeUrl}
                  />

                  <DataDisplayRowRadio
                    label="Smart Contract"
                    value={collection.smartContract}
                    opt1Value={"Collectible Protection Service Submissions"}
                    opt1Label="Collectible Protection Service Submissions"
                  />

                  <DataDisplayRowText
                    label="Smart Contract Address"
                    value={collection.smartContractAddress}
                  />

                  <p className="title is-5 mt-2">
                    <FontAwesomeIcon className="fas" icon={faCertificate} />
                    &nbsp;Minting&nbsp;&nbsp;&nbsp;-&nbsp;&nbsp;&nbsp;
                    <Link to={`/admin/collection/${id}/nfts/add/step-2`}>
                      <FontAwesomeIcon className="fas" icon={faPencil} />
                      &nbsp;Edit
                    </Link>
                  </p>

                  <DataDisplayRowText
                    label="To Address"
                    value={addNFT.toAddress}
                  />

                  <DataDisplayRowCheckbox
                     label="Wallet Password"
                     checked={true}
                  />

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <Link
                        to={`/admin/collection/${id}/nfts/add/step-2`}
                        class="button is-fullwidth-mobile is-medium"
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back to Step 2
                      </Link>
                    </div>
                    <div class="column is-half has-text-right">
                      <button
                        class="button is-primary is-fullwidth-mobile is-medium"
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

export default AdminNFTCollectionNFTAddStep1;
