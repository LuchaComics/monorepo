import React, { useEffect, useState } from "react";
import { Link, useSearchParams, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faGauge,
  faArrowRight,
  faUsers,
  faBarcode,
  faQuestionCircle,
  faWallet,
  faDonate,
  faCheckCircle,
  faFaucet,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import FormErrorBox from "../Reusable/FormErrorBox";
import { putProfileWalletAddressAPI } from "../../API/Profile";
import { topAlertMessageState, topAlertStatusState } from "../../AppState";
import { currentUserState } from "../../AppState";
import FormInputField from "../Reusable/FormInputField";

function CustomerAddWalletToFaucet() {
  ////
  //// URL Parameters.
  ////

  const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140
  const cpsrn = searchParams.get("cpsrn");

  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);
  const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

  ////
  //// Component states.
  ////

  // GUI state.
  const [errors, setErrors] = useState({});
  const [isFetching, setIsFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [wasFaucetRecentlySet, setWasFaucetRecentlySet] = useState(false);

  // Form State.
  const [walletAddress, setWalletAddress] = useState(currentUser ? currentUser.walletAddress : "");

  ////
  //// API.
  ////

  function onRegisterSuccess(response) {
    // For debugging purposes only.
    console.log("onRegisterSuccess: Starting...");
    console.log(response);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Submission created");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onRegisterSuccess: Delayed for 2 seconds.");
      console.log(
        "onRegisterSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    setWasFaucetRecentlySet(true);

    setCurrentUser(response);

    // // Redirect the user to a new page.
    setForceURL("/added-my-wallet-to-faucet-successfully");
  }

  function onRegisterError(apiErr) {
    console.log("onRegisterError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onRegisterError: Delayed for 2 seconds.");
      console.log(
        "onRegisterError: topAlertMessage, topAlertStatus:",
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

  function onRegisterDone() {
    console.log("onRegisterDone: Starting...");
    setIsFetching(false);
  }

  ////
  //// Event handling.
  ////

  const onSubmitClick = (e) => {
      e.preventDefault();
      console.log("onSubmitClick: Beginning...");
      setIsFetching(true);
      setErrors({});
      const submission = {
        wallet_address: walletAddress,
      }
      console.log("onSubmitClick, submission:", submission);
      putProfileWalletAddressAPI(
        submission,
        onRegisterSuccess,
        onRegisterError,
        onRegisterDone,
      );

  }

  ////
  //// Misc.
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

  console.log("currentUser: ", currentUser);
  console.log("walletAddress: ", walletAddress);

  return (
    <>
      <div class="container">
        <section class="section">
          <nav class="breadcrumb" aria-label="breadcrumbs">
            <ul>
              <li class="is-active">
                <Link to="/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faFaucet} />
                  &nbsp;ComicCoin Faucet
                </Link>
              </li>
            </ul>
          </nav>
          <nav class="box">
            {currentUser && <>
                <div class="columns">
                  <div class="column">
                    <h1 class="title is-4">
                      <FontAwesomeIcon className="fas" icon={faDonate} />
                      &nbsp;Get ComicCoins
                    </h1>
                  </div>
                </div>
                Welcome to the <b>ComicCoin Faucet</b>! To begin, please download the latest <b>ComicCoin Wallet</b> and set your wallet address below:
                <br />
                <br />
                <FormErrorBox errors={errors} />
                <FormInputField
                  label="Wallet Address"
                  name="walletAddress"
                  placeholder="Text input"
                  value={walletAddress}
                  errorText={errors && errors.walletAddress}
                  helpText=""
                  onChange={(e) => setWalletAddress(e.target.value)}
                  isRequired={true}
                  maxWidth="380px"
                />
                <div class="columns pt-5">
                    <div class="column is-half">
                        <button
                          class="button is-medium is-block is-fullwidth is-primary"
                          type="button"
                          onClick={onSubmitClick}
                        >
                          Submit&nbsp;
                          <FontAwesomeIcon icon={faArrowRight} />
                        </button>
                    </div>
                    <div class="column is-half has-text-right">

                    </div>
                </div>
            </>}

          </nav>
        </section>
      </div>
    </>
  );
}

export default CustomerAddWalletToFaucet;
