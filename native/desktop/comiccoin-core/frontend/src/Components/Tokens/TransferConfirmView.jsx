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
    TransferToken
} from "../../../wailsjs/go/main/App";
import { currentOpenWalletAtAddressState } from "../../AppState";
import PageLoadingContent from "../Reusable/PageLoadingContent";
import FormErrorBox from "../Reusable/FormErrorBox";
import FormInputField from "../Reusable/FormInputField";
import FormRadioField from "../Reusable/FormRadioField";
import FormRowYouTubeField from "../Reusable/FormRowYouTubeField";


function TransferConfirmView() {
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
    const [walletPassword, setWalletPassword] = useState("");

    ////
    //// Event handling.
    ////

    const onSubmitClick = (e) => {
        e.preventDefault();

        // Update the GUI to let user know that the operation is under way.
        setIsLoading(true);

        TransferToken(transferTo, parseInt(tokenID), currentOpenWalletAtAddress, walletPassword).then(()=>{
            console.log("Successful")
            setForceURL("/more/token/"+ tokenID + "/transfer-success");
        }).catch((errorJsonString)=>{
            console.log("errRes:", errorJsonString);
            const errorObject = JSON.parse(errorJsonString);
            let err = {};
            if (errorObject.recipient_address != "") {
                err.transferTo = errorObject.recipient_address;
            }
            if (errorObject.tokenID != "") {
                err.tokenID = errorObject.tokenID;
            }
            if (errorObject.value != "") {
                err.tokenID = errorObject.value;
            }
            if (errorObject.walletPassword != "") {
                err.walletPassword = errorObject.token_owner_password;
            }
            setErrors(err);
        }).finally(() => {
            // Update the GUI to let user know that the operation is completed.
            setIsLoading(false);
        });
    }

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
            <PageLoadingContent displayMessage="Sending..." />
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
                        &nbsp;Transfer
                    </Link>
                  </li>
                </ul>
              </nav>

              <nav class="box">
                <div class="columns">
                  <div class="column">
                    <h1 class="title is-4">
                      <FontAwesomeIcon className="fas" icon={faExchange} />
                      &nbsp;Transfer
                    </h1>
                  </div>
                </div>

                <FormErrorBox errors={errors} />

                <p class="pb-4">Please fill out all required fields before submitting.</p>

                {token !== undefined && token !== null && token != "" && <>
                    <FormInputField
                      label="Transfer To:"
                      name="transferTo"
                      placeholder="0x000.."
                      value={transferTo}
                      errorText={errors && errors.transferTo}
                      helpText="Enter a ComicCoin address (e.g. 0x38e26e225a391ee497b63b90820a95eb36b5add6)."
                      onChange={(e) => setTransferTo(e.target.value)}
                      isRequired={true}
                      maxWidth="400px"
                    />
                    <FormInputField
                      type="password"
                      label="Wallet Password:"
                      name="walletPassword"
                      placeholder=""
                      value={walletPassword}
                      errorText={errors && errors.walletPassword}
                      helpText="Your wallet is safely stored on only your computer in encrypted format and as result you'll need to submit a password to unlock the wallet to send with."
                      onChange={(e) => setWalletPassword(e.target.value)}
                      isRequired={true}
                      maxWidth="300px"
                    />

                    <div class="columns pt-5" style={{alignSelf: "flex-start"}}>
                      <div class="column is-half">
                        <Link
                          class="button is-fullwidth-mobile"
                          to={`/more/token/${tokenID}`}
                        >
                          <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                          &nbsp;Cancel
                        </Link>
                      </div>
                      <div class="column is-half has-text-right">
                        <button
                          class="button is-primary is-fullwidth-mobile"
                          onClick={onSubmitClick}
                        >
                          <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                          &nbsp;Send
                        </button>
                      </div>
                    </div>
                </>}
              </nav>
            </section>
          </div>
        </>
    )
}

export default TransferConfirmView
