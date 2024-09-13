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
  faCertificate,
  faArrowRight,
  faMoneyBillAlt,
  faExclamationCircle
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import useLocalStorage from "../../../../../Hooks/useLocalStorage";
import {
  putNFTUpdateAPI,
  getNFTDetailAPI,
} from "../../../../../API/NFT";
import {
  getCollectionDetailAPI,
  getCollectionWalletBalanceAPI,
  postCollectionMintOperationAPI,
} from "../../../../../API/NFTCollection";
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


function AdminNFTCollectionNFTMint() {
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

  // GUI related states.
  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [collection, setCollection] = useState({});
  const [balance, setBalance] = useState("");

  // Form submission states.
  const [tokenID, setTokenID] = useState(0);
  const [toAddress, setToAddress] = useState("");
  const [walletPassword, setWalletPassword] = useState("");

  ////
  //// Event handling.
  ////

  const onSubmitClick = (e) => {
      console.log("onSubmitClick: Starting...");
      setFetching(true);
      setErrors({});

      const jsonData = {
          collection_id: id,
          token_id: tokenID,
          to_address: toAddress,
          wallet_password: walletPassword,
      };

    postCollectionMintOperationAPI(
      jsonData,
      onAdminNFTCollectionNFTMintSuccess,
      onAdminNFTCollectionNFTMintError,
      onAdminNFTCollectionNFTMintDone,
      onUnauthorized,
    );
    console.log("onSubmitClick: Finished.");
  };

  ////
  //// API.
  ////

  // --- MINT --- //

  function onAdminNFTCollectionNFTMintSuccess(response) {
    // For debugging purposes only.
    console.log("onAdminNFTCollectionNFTMintSuccess: Starting...");
    console.log(response);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("NFT updated");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onAdminNFTCollectionNFTMintSuccess: Delayed for 2 seconds.");
      console.log(
        "onAdminNFTCollectionNFTMintSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Redirect the collection to the collection NFT metadatum page.
    setForceURL("/admin/collection/" + id + "/nft/" + rid);
  }

  function onAdminNFTCollectionNFTMintError(apiErr) {
    console.log("onAdminNFTCollectionNFTMintError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onAdminNFTCollectionNFTMintError: Delayed for 2 seconds.");
      console.log(
        "onAdminNFTCollectionNFTMintError: topAlertMessage, topAlertStatus:",
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

  function onAdminNFTCollectionNFTMintDone() {
    console.log("onAdminNFTCollectionNFTMintDone: Starting...");
    setFetching(false);
  }

  // --- DETAIL --- //

  function onAdminNFTCollectionNFTDetailSuccess(response) {
      // For debugging purposes only.
      console.log("onAdminNFTCollectionNFTDetailSuccess: Starting...");
      console.log("onAdminNFTCollectionNFTDetailSuccess: response:", response);
      setTokenID(response.tokenId);
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

      getNFTDetailAPI(
        rid,
        onAdminNFTCollectionNFTDetailSuccess,
        onAdminNFTCollectionNFTDetailError,
        onAdminNFTCollectionNFTDetailDone,
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
              <li class="">
                <Link
                  to={`/admin/collection/${id}/nft/${rid}`}
                  aria-current="page"
                >
                  <FontAwesomeIcon className="fas" icon={faCube} />
                  &nbsp;NFT
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faCertificate} />
                  &nbsp;Mint
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
                  &nbsp;Back to NFT
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modals */}
          {/* None */}

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faCertificate} />
              &nbsp;Mint
            </p>

            {/* <p class="pb-4 has-text-grey">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Submitting..."} />
            ) : (
              <>
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
                     <strong><FontAwesomeIcon className="fas" icon={faExclamationCircle} />&nbsp;Sufficient funds:</strong>&nbsp;You have enough funds in {collection.blockchain} wallet to proceed with minting.
                   </div>
                 </article>
                 </>}
                </>}
                {(balance !== "" && balance <= 0) ?
                    <>
                        <section class="hero is-medium has-background-white-ter">
                          <div class="hero-body">
                            <p class="title">
                              <FontAwesomeIcon className="fas" icon={faMoneyBillAlt} />
                              &nbsp;Not enough funds
                            </p>
                            <p class="subtitle">
                              It appears the wallet account for this smart contract does not have enough funds to mint a token. Please add some funds to the wallent account.
                            </p>
                          </div>
                        </section>
                    </>
                    :
                    <>
                        <FormErrorBox errors={errors} />
                        <div class="container content">

                          <p>
                            You are about to <b>mint</b> this NFT to the <i>{collection.blockchain} blockchain</i>; as a result, this will cost you funds from your wallet. This action cannot be undone and the NFT will exist permanently on the blochain.
                          </p>

                          <FormInputField
                            label="To Address"
                            name="toAddress"
                            placeholder="Text input"
                            value={toAddress}
                            errorText={errors && errors.toAddress}
                            helpText="Please enter the address to be the owner of the new NFT."
                            onChange={(e) => setToAddress(e.target.value)}
                            isRequired={true}
                            maxWidth="380px"
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

                          <div class="columns pt-5">
                            <div class="column is-half">
                              <Link
                                to={`/admin/collection/${id}/nft/${rid}`}
                                class="button is-hidden-touch"
                              >
                                <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                                &nbsp;Back
                              </Link>
                              <Link
                                to={`/admin/collection/${id}/nft/${rid}`}
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
                                &nbsp;Submit for NFT Minting
                              </button>
                              <button
                                class="button is-medium is-primary is-fullwidth is-hidden-desktop"
                                onClick={onSubmitClick}
                              >
                                <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                                &nbsp;Submit for NFT Minting
                              </button>
                            </div>
                          </div>
                        </div>
                    </>
                }
              </>
            )}
          </nav>
        </section>
      </div>
    </>
  );
}

export default AdminNFTCollectionNFTMint;
