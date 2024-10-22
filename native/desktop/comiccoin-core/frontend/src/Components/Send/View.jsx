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

import logo from '../../assets/images/CPS-logo-2023-square.webp';
import FormErrorBox from "../Reusable/FormErrorBox";
import FormInputField from "../Reusable/FormInputField";
import FormRadioField from "../Reusable/FormRadioField";


function SendView() {
    ////
    //// Component states.
    ////

    // GUI States.
    const [errors, setErrors] = useState({});
    const [forceURL, setForceURL] = useState("");

    // Form Submission States.
    const [toAddress, setToAddress] = useState("");
    const [amount, setAmount] = useState("");

    ////
    //// Event handling.
    ////

    ////
    //// API.
    ////

    const onSubmitClick = (e) => {
        e.preventDefault();
        setForceURL("/dashboard")
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
                <p class="pb-4">Please fill out required fields:</p>

                <FormInputField
                  label="Pay To:"
                  name="toAddress"
                  placeholder="0x000.."
                  value={toAddress}
                  errorText={errors && errors.toAddress}
                  helpText=""
                  onChange={(e) => setToAddress(e.target.value)}
                  isRequired={true}
                  maxWidth="300px"
                />

                <FormInputField
                  label="Amount:"
                  name="amount"
                  placeholder="0"
                  value={amount}
                  errorText={errors && errors.amount}
                  helpText=""
                  onChange={(e) => setAmount(e.target.value)}
                  isRequired={true}
                  maxWidth="300px"
                />

                <FormInputField
                  label="Transaction Fee:"
                  name="transactionFee"
                  placeholder=""
                  value={`0`}
                  errorText={errors && errors.transactionFee}
                  helpText="ComicCoin has no transaction fees."
                  onChange={null}
                  isRequired={true}
                  maxWidth="300px"
                  disabled={true}
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



              </nav>
            </section>
          </div>
        </>
    )
}

export default SendView
