import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faGauge,
  faArrowRight,
  faUsers,
  faBarcode,
  faCubes,
  faPaperPlane,
  faTimesCircle,
  faCheckCircle,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import logo from '../../assets/images/CPS-logo-2023-square.webp';
import FormErrorBox from "../Reusable/FormErrorBox";
import FormInputField from "../Reusable/FormInputField";
import FormRadioField from "../Reusable/FormRadioField";
import FormTextareaField from "../Reusable/FormTextareaField";
import {TransferCoin} from "../../../wailsjs/go/main/App";
import { currentOpenWalletAtAddressState } from "../../AppState";
import PageLoadingContent from "../Reusable/PageLoadingContent";


function SendCoinSuccessView() {
    ////
    //// Global State
    ////

    const [currentOpenWalletAtAddress] = useRecoilState(currentOpenWalletAtAddressState);

    ////
    //// Component states.
    ////

    // GUI States.
    const [errors, setErrors] = useState({});
    const [forceURL, setForceURL] = useState("");
    const [isLoading, setIsLoading] = useState(false);

    // Form Submission States.
    const [payTo, setPayTo] = useState("");
    const [coin, setCoin] = useState(0);
    const [message, setMessage] = useState("");
    const [walletPassword, setWalletPassword] = useState("");

    ////
    //// Event handling.
    ////

    ////
    //// API.
    ////

    const onSubmitClick = (e) => {
        e.preventDefault();


    }

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
              <nav class="box">
                <div class="columns">
                  <div class="column">
                    <h1 class="title is-4">
                      <FontAwesomeIcon className="fas" icon={faPaperPlane} />
                      &nbsp;Send ComicCoins
                    </h1>
                  </div>
                </div>

                <section class="hero is-success is-halfheight">
                  <div class="hero-body">
                    <div class="">
                      <p class="title"> <FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;Coins sent!</p>
                      <p class="subtitle">You have successfully sent coin(s) to the specified account. Please wait a few minutes for the transaction to get processed on the blockchain.</p>
                    </div>
                  </div>
                </section>

                <div class="columns pt-5" style={{alignSelf: "flex-start"}}>
                  <div class="column is-half">
                    {/*
                    <button
                      class="button is-fullwidth-mobile"
                      onClick={(e) => setShowCancelWarning(true)}
                    >
                      <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                      &nbsp;Clear
                    </button>
                    */}
                  </div>
                  <div class="column is-half has-text-right">
                    <Link
                      class="button is-primary is-fullwidth-mobile"
                      to="/more/transactions"
                    >
                      Go to transactions&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} />
                    </Link>
                  </div>
                </div>

              </nav>
            </section>
          </div>
        </>
    );
}

export default SendCoinSuccessView
