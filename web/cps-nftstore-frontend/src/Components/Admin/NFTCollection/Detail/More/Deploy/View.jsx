import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faBuildingCollection,
  faImage,
  faPaperclip,
  faAddressCard,
  faSquarePhone,
  faTasks,
  faTachometer,
  faPlus,
  faArrowLeft,
  faCheckCircle,
  faCubes,
  faGauge,
  faPencil,
  faEye,
  faIdCard,
  faAddressBook,
  faContactCard,
  faChartPie,
  faBuilding,
  faEllipsis,
  faRocket,
  faExclamationCircle
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import {
  getCollectionDetailAPI,
  getCollectionWalletBalanceAPI,
  postCollectionDeploySmartContractAPI,
} from "../../../../../../API/NFTCollection";
import FormErrorBox from "../../../../../Reusable/FormErrorBox";
import PageLoadingContent from "../../../../../Reusable/PageLoadingContent";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../../../AppState";
import FormInputField from "../../../../../Reusable/FormInputField";


function AdminNFTCollectionDetailMoreDeployOperation() {
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

  // GUI related states.
  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [collection, setCollection] = useState({});
  const [balance, setBalance] = useState("");

  // Form submission states.
  const [walletPassword, setWalletPassword] = useState("");

  ////
  //// Event handling.
  ////

  const onSubmitClick = () => {
    setErrors({});
    setFetching(true);
    postCollectionDeploySmartContractAPI(
      {
          collection_id: id,
          wallet_password: walletPassword,
      },
      onDeploySuccess,
      onDeployError,
      onDeployDone,
      onUnauthorized,
    );
  };

  ////
  //// API.
  ////

  // --- Detail --- //

  function onSuccess(response) {
    console.log("onSuccess: Starting...");
    setCollection(response);
  }

  function onError(apiErr) {
    console.log("onError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onDone() {
    console.log("onDone: Starting...");
    setFetching(false);
  }

  // --- Wallet Balance --- //

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


  // --- Deploy --- //

  function onDeploySuccess(response) {
    console.log("onDeploySuccess: Starting...");

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Collection archived");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onSuccess: Delayed for 2 seconds.");
      console.log(
        "onSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    setForceURL("/admin/collection/" + id + "/more");
  }

  function onDeployError(apiErr) {
    console.log("onDeployError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onDeployDone() {
    console.log("onDeployDone: Starting...");
    setFetching(false);
  }

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
      getCollectionDetailAPI(id, onSuccess, onError, onDone, onUnauthorized);
      getCollectionWalletBalanceAPI(id, onWalletBalanceSuccess, onWalletBalanceError, onWalletBalanceDone, onUnauthorized);
    }

    return () => {
      mounted = false;
    };
  }, [id]);

  ////
  //// Component rendering.
  ////

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <>
      <div className="container">
        <section className="section">
          {/* Desktop Breadcrumbs */}
          <nav
            className="breadcrumb is-hidden-touch"
            aria-label="breadcrumbs"
          >
            <ul>
              <li className="">
                <Link to="/admin/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Admin Dashboard
                </Link>
              </li>
              <li className="">
                <Link to="/admin/collections" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faCubes} />
                  &nbsp;NFT Collections
                </Link>
              </li>
              <li className="">
                <Link to={`/admin/collection/${id}/more`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail (More)
                </Link>
              </li>
              <li className="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faRocket} />
                  &nbsp;Deploy
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav
            className="breadcrumb has-background-light is-hidden-desktop p-4"
            aria-label="breadcrumbs"
          >
            <ul>
              <li className="">
                <Link to={`/admin/collection/${id}/more`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Detail
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page Title */}
          <h1 className="title is-2">
            <FontAwesomeIcon className="fas" icon={faCubes} />
            &nbsp;Collection
          </h1>
          <h4 className="subtitle is-4">
            <FontAwesomeIcon className="fas" icon={faEye} />
            &nbsp;Detail
          </h4>
          <hr />

          {/* Page */}
          <nav className="box">
            {/* Title + Options */}
            {collection && (
              <div className="columns">
                <div className="column">
                  <p className="title is-4">
                    <FontAwesomeIcon className="fas" icon={faRocket} />
                    &nbsp;Deploy Smart Contract - Are you sure?
                  </p>
                </div>
                <div className="column has-text-right"></div>
              </div>
            )}

            {/* <p className="pb-4">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                <FormErrorBox errors={errors} />

                {collection && (
                  <>
                      {balance != "" && <>
                        {balance <= 0 ? <>
                           <article class="message is-danger">
                             <div class="message-body">
                               <strong><FontAwesomeIcon className="fas" icon={faExclamationCircle} />&nbsp;Not enough coins:</strong>&nbsp;Please add some coins to your {collection.blockchain} wallet before proceeding with deployment.
                             </div>
                           </article>
                       </> : <>
                       <article class="message is-info">
                         <div class="message-body">
                           <strong><FontAwesomeIcon className="fas" icon={faExclamationCircle} />&nbsp;Sufficient coins amount:</strong>&nbsp;You have enough coins in {collection.blockchain} wallet to proceed with deployment.
                         </div>
                       </article>
                       </>}
                      </>}
                      <div className="container content">
                        <p>
                          You are about to <b>deploy</b> the smart contract, called <u>{collection.smartContract}</u> in this collection, to the <i>{collection.blockchain} blockchain</i>; as a result, this will cost you coins from your wallet. This action cannot be undone and the smart contract will exist permanently on the blochain.
                        </p>
                        <p>Before proceeding, it is recommended you have done the following:</p>
                        <ul>
                            <li>Safely backed up your wallet mnemonic</li>
                            <li>Safely backed up your wallet password</li>
                        </ul>
                        <p>Are you sure you would like to continue? If so, please enter your <b>wallet password</b>:</p>

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

                        {/* Bottom Navigation */}
                        <div className="columns pt-5">
                          <div className="column is-half">
                            <Link
                              className="button is-fullwidth-mobile"
                              to={`/admin/collection/${id}/more`}
                            >
                              <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                              &nbsp;Back to Detail
                            </Link>
                          </div>
                          <div className="column is-half has-text-right">
                            <button
                              className="button is-success is-fullwidth-mobile"
                              onClick={onSubmitClick}
                              type="button"
                              disabled={balance != "" && balance <= 0}
                            >
                              <FontAwesomeIcon
                                className="fas"
                                icon={faCheckCircle}
                                type="button"
                              />
                              &nbsp;Confirm and Deploy
                            </button>
                          </div>
                        </div>
                      </div>
                  </>
                )}
              </>
            )}
          </nav>
        </section>
      </div>
    </>
  );
}

export default AdminNFTCollectionDetailMoreDeployOperation;
