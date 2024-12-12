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
  faHand,
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
             <section class="hero is-medium is-success">
               <div class="hero-body">
                 <p class="title">
                   <FontAwesomeIcon className="fas" icon={faHand} />
                   &nbsp;Welcome to ComicCoin Faucet
                 </p>
                 <p class="subtitle">
                   Do you have a question or concern? Contact us below.
                   <br />
                   <br />
                   <Link to={"/help"}>
                     Go to Help&nbsp;
                     <FontAwesomeIcon className="fas" icon={faArrowRight} />
                   </Link>
                 </p>
               </div>
             </section>


          </nav>
        </section>
      </div>
    </>
  );
}

export default CustomerDashboard;
