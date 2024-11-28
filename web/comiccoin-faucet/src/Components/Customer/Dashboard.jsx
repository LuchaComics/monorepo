import React, { useEffect, useState } from "react";
import { Link, useSearchParams, Navigate } from "react-router-dom";
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
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { topAlertMessageState, topAlertStatusState } from "../../AppState";
import { currentUserState } from "../../AppState";
import FormInputField from "../Reusable/FormInputField";

function CustomerDashboard() {
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

  // Form State.
  const [walletAddress, setWalletAddress] = useState("");

  ////
  //// API.
  ////

  ////
  //// Event handling.
  ////

  const onSubmitClick = (e) => {
      e.preventDefault();


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
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Dashboard
                </Link>
              </li>
            </ul>
          </nav>
          <nav class="box">


            {currentUser && <>
                {currentUser.walletAddress ? <>
                    <div class="columns">
                      <div class="column">
                        <h1 class="title is-4">
                          <FontAwesomeIcon className="fas" icon={faGauge} />
                          &nbsp;Dashboard
                        </h1>
                      </div>
                    </div>
                    <section class="hero is-medium is-link">
                      <div class="hero-body">
                        <p class="title">
                          <FontAwesomeIcon className="fas" icon={faTasks} />
                          &nbsp;My Submissions
                        </p>
                        <p class="subtitle">
                          Submit a request to encapsulate your collectible or view existing collectibles by clicking
                          below:
                          <br />
                          <br />
                          <Link to={"/c/submissions/comics"}>
                            View Online Comic Submissions&nbsp;
                            <FontAwesomeIcon className="fas" icon={faArrowRight} />
                          </Link>
                          <br />
                          <br />
                          <Link to={"/c/submissions/pick-type-for-add"}>
                            Add&nbsp;
                            <FontAwesomeIcon className="fas" icon={faArrowRight} />
                          </Link>
                        </p>
                      </div>
                    </section>
                </> : <>
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
            </>}

          </nav>
        </section>
      </div>
    </>
  );
}

export default CustomerDashboard;
