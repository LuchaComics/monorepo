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


function SendCoinSubmissionView() {
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

        // Update the GUI to let user know that the operation is under way.
        setIsLoading(true);

        TransferCoin(payTo, parseInt(coin), message, currentOpenWalletAtAddress, walletPassword).then(()=>{
            console.log("onSubmitClick: Successful")
            setForceURL("/send-success");
        }).catch((errorJsonString)=>{
            console.log("onSubmitClick: errRes:", errorJsonString);
            const errorObject = JSON.parse(errorJsonString);
            console.log("onSubmitClick: errorObject:", errorObject);

            let err = {};
            if (errorObject.to != "") {
                err.payTo = errorObject.to;
            }
            if (errorObject.coin != "") {
                err.coin = errorObject.coin;
            }
            if (errorObject.value != "") {
                err.coin = errorObject.value;
            }
            if (errorObject.message != "") {
                err.message = errorObject.message;
            }
            if (errorObject.wallet_password != "") {
                err.walletPassword = errorObject.wallet_password;
            }
            console.log("onSubmitClick: err:", err);
            window.scrollTo(0, 0); // Start the page at the top of the page.
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

                <FormErrorBox errors={errors} />

                {isLoading ? <>
                    <PageLoadingContent displayMessage="Sending..." />
                </> : <>
                    <p class="pb-4">Please fill out all required fields:</p>

                    <FormInputField
                      label="Pay To:"
                      name="payTo"
                      placeholder="0x000.."
                      value={payTo}
                      errorText={errors && errors.payTo}
                      helpText="Enter a ComicCoin address (e.g. 0x38e26e225a391ee497b63b90820a95eb36b5add6)."
                      onChange={(e) => setPayTo(e.target.value)}
                      isRequired={true}
                      maxWidth="400px"
                    />

                    <FormInputField
                      type="number"
                      label="Coin(s):"
                      name="coin"
                      placeholder="0"
                      value={coin}
                      errorText={errors && errors.coin}
                      helpText=""
                      onChange={(e) => setCoin(e.target.value)}
                      isRequired={true}
                      maxWidth="300px"
                    />

                    <FormTextareaField
                      label="Message (Optional)"
                      name="message"
                      placeholder="Enter your message here..."
                      value={message}
                      errorText={errors && errors.message}
                      onChange={(e) => setMessage(e.target.value)}
                      isRequired={true}
                      maxWidth="280px"
                      helpText={"Optional field you may use to write a message to the receipient."}
                      rows={4}
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
                        <button
                          class="button is-fullwidth-mobile"
                          onClick={(e) => setShowCancelWarning(true)}
                        >
                          <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                          &nbsp;Clear
                        </button>
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
    );
}

export default SendCoinSubmissionView
