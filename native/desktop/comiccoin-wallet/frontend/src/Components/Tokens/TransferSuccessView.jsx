import {useState, useEffect} from 'react';
import { Link, useParams, Navigate } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
    faTasks,
    faGauge,
    faArrowRight,
    faUsers,
    faBarcode,
    faCubes,
    faCube,
    faCoins,
    faEllipsis,
    faExchange,
    faCheckCircle,
    faTimesCircle
} from "@fortawesome/free-solid-svg-icons";
import logo from '../../assets/images/CPS-logo-2023-square.webp';
import { useRecoilState } from "recoil";
import { toLower } from "lodash";

import {
    GetNonFungibleToken,
} from "../../../wailsjs/go/main/App";
import { currentOpenWalletAtAddressState } from "../../AppState";
import PageLoadingContent from "../Reusable/PageLoadingContent";
import FormErrorBox from "../Reusable/FormErrorBox";
import FormInputField from "../Reusable/FormInputField";
import FormRadioField from "../Reusable/FormRadioField";
import FormTextareaField from "../Reusable/FormTextareaField";
import FormRowYouTubeField from "../Reusable/FormRowYouTubeField";


function TokenTransferSuccessView() {
    ////
    //// URL Parameters.
    ////

    const { tokenID } = useParams();

    ////
    //// Global State
    ////

    const [currentOpenWalletAtAddress] = useRecoilState(currentOpenWalletAtAddressState);

    ////
    //// Component states.
    ////

    // GUI States.
    const [isLoading, setIsLoading] = useState(false);
    const [forceURL, setForceURL] = useState("");
    const [token, setToken] = useState([]);
    const [errors, setErrors] = useState({});

    // Form Submission States.
    const [transferTo, setTransferTo] = useState("");
    const [message, setMessage] = useState("");
    const [walletPassword, setWalletPassword] = useState("");

    ////
    //// Event handling.
    ////



    ////
    //// Misc.
    ////

    useEffect(() => {
      let mounted = true;

      if (mounted) {
            window.scrollTo(0, 0); // Start the page at the top of the page.

            // Update the GUI to let user know that the operation is under way.
            setIsLoading(true);

            GetNonFungibleToken(parseInt(tokenID)).then((nftokRes)=>{
                console.log("GetNonFungibleToken: nftokRes:", nftokRes);
                setToken(nftokRes);
            }).catch((errorRes)=>{
                console.log("GetNonFungibleToken: errorRes:", errorRes);
            }).finally((errorRes)=>{
                // Update the GUI to let user know that the operation is completed.
                setIsLoading(false);
            });
      }

      return () => {
          mounted = false;
      };
    }, [currentOpenWalletAtAddress]);

    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
        return <Navigate to={forceURL} />;
    }

    if (isLoading) {
        return (
            <PageLoadingContent displayMessage="Loading..." />
        );
    }

    return (
        <>
          <div class="container">
            <section class="section">
              <nav class="breadcrumb" aria-label="breadcrumbs">
                <ul>
                  <li>
                    <Link to="/dashboard" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faGauge} />
                      &nbsp;Overview
                    </Link>
                  </li>
                  <li>
                    <Link to="/more" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faEllipsis} />
                      &nbsp;More
                    </Link>
                  </li>
                  <li class="">
                    <Link to="/more/tokens" aria-current="page">
                        <FontAwesomeIcon className="fas" icon={faCubes} />
                        &nbsp;Tokens
                    </Link>
                  </li>
                  <li>
                    <Link to={`/more/token/${tokenID}`} aria-current="page">
                        <FontAwesomeIcon className="fas" icon={faCube} />
                        &nbsp;Token ID {tokenID}
                    </Link>
                  </li>
                  <li class="is-active">
                    <Link to={`/more/token/${tokenID}/transfer`} aria-current="page">
                        <FontAwesomeIcon className="fas" icon={faExchange} />
                        &nbsp;Transfer - Success
                    </Link>
                  </li>
                </ul>
              </nav>

              <nav class="box">
                <div class="columns">
                  <div class="column">
                    <h1 class="title is-4">
                      <FontAwesomeIcon className="fas" icon={faExchange} />
                      &nbsp;Transfer - Success
                    </h1>
                  </div>
                </div>

                <section class="hero is-success is-halfheight">
                  <div class="hero-body">
                    <div class="">
                      <p class="title"> <FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;Token transfered!</p>
                      <p class="subtitle">You have successfully transfered a token to the specified account. Please wait a few minutes for the transaction to get processed on the blockchain.</p>
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
                      to="/more/tokens"
                    >
                      Go to tokens&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} />
                    </Link>
                  </div>
                </div>

              </nav>
            </section>
          </div>
        </>
    )
}

export default TokenTransferSuccessView
