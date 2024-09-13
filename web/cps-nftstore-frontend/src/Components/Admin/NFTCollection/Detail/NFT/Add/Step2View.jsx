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
  faArrowRight,
  faExclamationCircle,
  faCertificate
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import useLocalStorage from "../../../../../../Hooks/useLocalStorage";
import { getCollectionDetailAPI, getCollectionWalletBalanceAPI } from "../../../../../../API/NFTCollection";
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
  const [balance, setBalance] = useState("");

  // Form fields
  const [toAddress, setToAddress] = useState(addNFTState.toAddress);
  const [walletPassword, setWalletPassword] = useState(addNFTState.walletPassword);

  ////
  //// Event handling.
  ////

  const onSubmitClick = (e) => {
    console.log("onSubmitClick: Beginning...");
    setErrors({}); // Reset the errors in the GUI.

    // Variables used to create our new errors if we find them.
    let newErrors = {};
    let hasErrors = false;

    if (toAddress === undefined || toAddress === null || toAddress === "") {
      newErrors["toAddress"] = "missing value";
      hasErrors = true;
    }
    if (walletPassword === undefined || walletPassword === null || walletPassword === "") {
      newErrors["walletPassword"] = "missing value";
      hasErrors = true;
    }

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
    modifiedAddNFT.toAddress = toAddress;
    modifiedAddNFT.walletPassword = walletPassword;

    setAddNFT(modifiedAddNFT);
    setForceURL("/admin/collection/" + id + "/nfts/add/step-3");
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

  // --- BALANCE OPERATION --- //

  function onWalletBalanceSuccess(response) {
    console.log("onWalletBalanceSuccess: Starting...");
    console.log("onWalletBalanceSuccess: response:", response);
    setBalance(response.value);
  }

  function onWalletBalanceError(apiErr) {
    console.log("onWalletBalanceError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onWalletBalanceDone() {
    console.log("onWalletBalanceDone: Starting...");
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

      getCollectionWalletBalanceAPI(
          id,
          onWalletBalanceSuccess,
          onWalletBalanceError,
          onWalletBalanceDone,
          onUnauthorized
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
          <nav className="box has-background-light">
            <p className="subtitle is-5">Step 2 of 3</p>
            <progress
              class="progress is-success"
              value="66"
              max="100"
            >
              66%
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
                  {collection !== undefined && collection !== null && collection !== "" && <>

                      {balance != "" && <>
                        {balance <= 0 ? <>
                           <article class="message is-danger">
                             <div class="message-body">
                               <strong><FontAwesomeIcon className="fas" icon={faExclamationCircle} />&nbsp;Not enough funds:</strong>&nbsp;Please add some funds to your {collection.blockchain} wallet before proceeding with deployment.
                             </div>
                           </article>
                       </> : <>
                       <article class="message is-info">
                         <div class="message-body">
                           <strong><FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;Sufficient funds:</strong>&nbsp;You have enough funds in {collection.blockchain} wallet to proceed with minting.
                         </div>
                       </article>
                       </>}
                      </>}
                      <p class="subtitle is-6">
                        <FontAwesomeIcon className="fas" icon={faChain} />
                        &nbsp;Blockchain Information
                      </p>
                      <hr />

                      {collection !== undefined && collection !== null && collection !== "" && <>
                        <FormRadioField
                          label="Blockchain"
                          name="blockchain"
                          value={collection.blockchain}
                          opt1Value={"ethereum"}
                          opt1Label="Ethereum"
                          errorText={errors && errors.blockchain}
                          onChange={(e) =>{}}
                          helpText=""
                          maxWidth="180px"
                          disabled={true}
                        />

                        <FormInputField
                          label="Node URL"
                          name="nodeURL"
                          placeholder="Text input"
                          value={collection.nodeUrl}
                          errorText={errors && errors.nodeURL}
                          helpText={<>Please enter the url to connect to the blockchain node. For local developers use: <i>http://ganache:8545</i></>}
                          onChange={(e) => {}}
                          isRequired={true}
                          maxWidth="400px"
                          disabled={true}
                        />

                        <FormRadioField
                          label="Smart Contract"
                          name="smartContract"
                          value={collection.smartContract}
                          opt1Value={"Collectible Protection Service Submissions"}
                          opt1Label="Collectible Protection Service Submissions"
                          errorText={errors && errors.smartContract}
                          onChange={(e) =>{}
                          }
                          maxWidth="180px"
                          helpText="Please select the smart contract to use to build our NFT collection on."
                          disabled={true}
                        />

                        <FormInputField
                          label="Smart Contract Address"
                          name="smartContractAddress"
                          placeholder="Text input"
                          value={collection.smartContractAddress}
                          errorText={errors && errors.smartContractAddress}
                          helpText={<></>}
                          onChange={(e) => {}}
                          isRequired={true}
                          maxWidth="430px"
                          disabled={true}
                        />

                        <p class="subtitle is-6">
                          <FontAwesomeIcon className="fas" icon={faCertificate} />
                          &nbsp;Minting Information
                        </p>
                        <hr />

                        <FormInputField
                          label="To Address"
                          name="toAddress"
                          placeholder="Text input"
                          value={toAddress}
                          errorText={errors && errors.toAddress}
                          helpText="Please enter the address to be the owner of the new NFT."
                          onChange={(e) => setToAddress(e.target.value)}
                          isRequired={true}
                          maxWidth="425px"
                          disabled={balance != "" && balance <= 0}
                        />

                        <FormInputField
                          label="Wallet Password"
                          type="password"
                          name="walletPassword"
                          placeholder="Text input"
                          value={walletPassword}
                          errorText={errors && errors.walletPassword}
                          helpText="Please enter the password you set during NFT collection creation process."
                          onChange={(e) => setWalletPassword(e.target.value)}
                          isRequired={true}
                          maxWidth="380px"
                          disabled={balance != "" && balance <= 0}
                        />

                      </>}
                  </>}

                  <br />
                  <br />

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <Link
                        to={`/admin/collection/${id}/nfts/add/step-1`}
                        class="button is-fullwidth-mobile is-medium"
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back to Step 1
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
