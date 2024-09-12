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
  faArrowRight
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import * as bip39 from '@scure/bip39';
import { wordlist } from '@scure/bip39/wordlists/english';

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
import {
    topAlertMessageState,
    topAlertStatusState,
    addNFTCollectionState,
    ADD_NFT_COLLECTION_STATE_DEFAULT
} from "../../../../AppState";


function AdminNFTCollectionAddStep1() {
  ////
  //// URL Parameters.
  ////

  const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140
  const orgID = searchParams.get("tenant_id");
  const orgName = searchParams.get("tenant_name");

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

  // GUI states.
  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [showCancelWarning, setShowCancelWarning] = useState(false);

  // Form states.
  const [blockchain, setBlockchain] = useState(addNFTCollection.blockchain);
  const [nodeURL, setNodeURL] = useState(addNFTCollection.nodeURL);
  const [smartContract, setSmartContract] = useState(addNFTCollection.smartContract);
  const [walletMnemonic, setWalletMnemonic] = useState(addNFTCollection.walletMnemonic);
  const [walletPassword, setWalletPassword] = useState(addNFTCollection.walletPassword);

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

      if (blockchain === undefined || blockchain === null || blockchain === "") {
        newErrors["blockchain"] = "missing value";
        hasErrors = true;
      } else {
        if (nodeURL === undefined || nodeURL === null || nodeURL === "") {
            newErrors["nodeURL"] = "missing value";
            hasErrors = true;
        }
        if (smartContract === undefined || smartContract === null || smartContract === "") {
            newErrors["smartContract"] = "missing value";
            hasErrors = true;
        } else {
            if (walletMnemonic === undefined || walletMnemonic === null || walletMnemonic === "") {
                newErrors["walletMnemonic"] = "missing value";
                hasErrors = true;
            }
            if (walletPassword === undefined || walletPassword === null || walletPassword === "") {
                newErrors["walletPassword"] = "missing value";
                hasErrors = true;
            }
        }
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
      let modifiedAddNFTCollection = { ...addNFTCollection };
      modifiedAddNFTCollection.blockchain = blockchain;
      modifiedAddNFTCollection.nodeURL = nodeURL;
      modifiedAddNFTCollection.nodeUrl = nodeURL; // Bugfix when making api call.
      modifiedAddNFTCollection.smartContract = smartContract;
      modifiedAddNFTCollection.walletMnemonic = walletMnemonic;
      modifiedAddNFTCollection.walletPassword = walletPassword;
      setAddNFTCollection(modifiedAddNFTCollection);
      setForceURL("/admin/collections/add/step-2");
  };

  const onCancelClick = (e, url) => {
      e.preventDefault();
      setAddNFTCollection(ADD_NFT_COLLECTION_STATE_DEFAULT);
      setForceURL(url);
  };

  const onGenerateMnemonic = (e) => {
      e.preventDefault();
      // Generate x random words. Uses Cryptographically-Secure Random Number Generator.
      const mn = bip39.generateMnemonic(wordlist); // Special thanks to: https://github.com/paulmillr/scure-bip39
      console.log(mn);

      setWalletMnemonic(mn);
  }

  ////
  //// API.
  ////

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
                  {orgName !== undefined &&
                  orgName !== null &&
                  orgName !== "" ? (
                    <Link
                      class="button is-medium is-success"
                      onClick={(e)=>{
                          onCancelClick(e, "/admin/tenant/"+orgID+"/collections");
                      }}
                    >
                      Yes
                    </Link>
                  ) : (
                    <Link
                      class="button is-medium is-success"
                      onClick={(e)=>{
                          onCancelClick(e, "/admin/collections");
                      }}
                    >
                      Yes
                    </Link>
                  )}
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
              43%
            </progress>
          </nav>

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faPlus} />
              &nbsp;New NFT Collection
            </p>
            <FormErrorBox errors={errors} />

            {/* <p class="pb-4 has-text-grey">Please fill out all the required fields before submitting this form.</p> */}

            {walletMnemonic !== undefined && walletMnemonic !== null && walletMnemonic !== "" && <>
                <article class="message is-warning">
                  <div class="message-body">
                    <strong><FontAwesomeIcon className="fas" icon={faExclamationTriangle} />&nbsp;Warning:</strong>&nbsp;Please keep a copy of your <b>Wallet Mnemonic</b> somewhere safe. Loss of this will make you unable to get access to your wallet!
                  </div>
                </article>
            </>}

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
                  <p class="subtitle is-6">
                    <FontAwesomeIcon className="fas" icon={faChain} />
                    &nbsp;Blockchain Information
                  </p>
                  <hr />

                  <FormRadioField
                    label="Blockchain"
                    name="blockchain"
                    value={blockchain}
                    opt1Value={"ethereum"}
                    opt1Label="Ethereum"
                    errorText={errors && errors.blockchain}
                    onChange={(e) =>
                      setBlockchain(e.target.value)
                    }
                    helpText="Please select the blockchain to deploy this NFT collection to."
                    maxWidth="180px"
                  />

                  {blockchain !== undefined && blockchain !== null && blockchain !== "" && <>
                    <FormInputField
                      label="Node URL"
                      name="nodeURL"
                      placeholder="Text input"
                      value={nodeURL}
                      errorText={errors && errors.nodeURL}
                      helpText={<>Please enter the url to connect to the blockchain node. For local developers use: <i>http://localhost:8545</i></>}
                      onChange={(e) => setNodeURL(e.target.value)}
                      isRequired={true}
                      maxWidth="380px"
                    />

                    <FormRadioField
                      label="Smart Contract"
                      name="smartContract"
                      value={smartContract}
                      opt1Value={"Collectible Protection Service Submissions"}
                      opt1Label="Collectible Protection Service Submissions"
                      errorText={errors && errors.smartContract}
                      onChange={(e) =>
                        setSmartContract(e.target.value)
                      }
                      maxWidth="180px"
                      helpText="Please select the smart contract to use to build our NFT collection on."
                    />

                    {smartContract !== undefined && smartContract !== null && smartContract !== "" && <>
                        <FormTextareaField
                            label="Wallet Mnemonic"
                            name="walletMnemonic"
                            placeholder="Text input"
                            value={walletMnemonic}
                            errorText={errors && errors.walletMnemonic}
                            helpText="Write"
                            onChange={(e) => setWalletMnemonic(e.target.value)}
                            isRequired={true}
                            maxWidth="280px"
                            helpText={<>Please enter the phrase to <b>import your existing wallet</b> or <b>generate a new wallet</b> which will be tied to this NFT collection. <Link onClick={(e)=>onGenerateMnemonic(e)}>Click here</Link> to generate a new wallet mnemonic.</>}
                            rows={4}
                        />

                        <FormInputField
                          label="Wallet Password"
                          type="password"
                          name="walletPassword"
                          placeholder="Text input"
                          value={walletPassword}
                          errorText={errors && errors.walletPassword}
                          helpText="Choose a password to use every time you want mint an NFT in this collection."
                          onChange={(e) => setWalletPassword(e.target.value)}
                          isRequired={true}
                          maxWidth="380px"
                        />
                    </>}
                  </>}

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <button
                        class="button is-medium is-hidden-touch"
                        onClick={(e) => setShowCancelWarning(true)}
                      >
                        <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                        &nbsp;Cancel
                      </button>
                      <button
                        class="button is-medium is-fullwidth is-hidden-desktop"
                        onClick={(e) => setShowCancelWarning(true)}
                      >
                        <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                        &nbsp;Cancel
                      </button>
                    </div>
                    <div class="column is-half has-text-right">
                      <button
                        class="button is-medium is-primary is-hidden-touch"
                        onClick={onSubmitClick}
                      >
                        Save & Continue&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} />
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

export default AdminNFTCollectionAddStep1;
