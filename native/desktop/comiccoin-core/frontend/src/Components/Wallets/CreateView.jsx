import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArrowLeft,
  faTasks,
  faTachometer,
  faHandHoldingHeart,
  faTimesCircle,
  faCheckCircle,
  faUserCircle,
  faGauge,
  faPencil,
  faUsers,
  faIdCard,
  faAddressBook,
  faContactCard,
  faChartPie,
  faBuilding,
  faCogs,
  faEllipsis,
  faPlus
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import FormErrorBox from "../Reusable/FormErrorBox";
import FormRadioField from "../Reusable/FormRadioField";
import FormInputField from "../Reusable/FormInputField";
import FormInputFieldWithButton from "../Reusable/FormInputFieldWithButton";
import PageLoadingContent from "../Reusable/PageLoadingContent";
import { CreateWallet } from "../../../wailsjs/go/main/App";
import { currentOpenWalletAtAddressState } from "../../AppState";


function CreateWalletView() {
    ////
    //// Global State
    ////

    const [currentOpenWalletAtAddress, setCurrentOpenWalletAtAddress] = useRecoilState(currentOpenWalletAtAddressState);

    ////
    //// Component states.
    ////

    const [label, setLabel] = useState("");
    const [password, setPassword] = useState("");
    const [passwordRepeated, setPasswordRepeated] = useState("");
    const [errors, setErrors] = useState({});
    const [forceURL, setForceURL] = useState("");

    ////
    //// Event handling.
    ////

    const onSubmitClick = (e) => {
        CreateWallet(password, passwordRepeated, label).then((addressRes)=>{
            console.log("address:", addressRes);
            console.log("currentOpenWalletAtAddress:", currentOpenWalletAtAddress);
            setCurrentOpenWalletAtAddress(addressRes);
            setForceURL("/dashboard");
        }).catch((errorJsonString)=>{
            console.log("errRes:", errorJsonString);
            const errorObject = JSON.parse(errorJsonString);
            let err = {};
            if (errorObject.wallet_password != "") {
                err.password = errorObject.wallet_password;
            }
            if (errorObject.wallet_password_repeated != "") {
                err.passwordRepeated = errorObject.wallet_password_repeated;
            }
            setErrors(err);
        }).finally(() => {
            // this will be executed after then or catch has been executed
            console.log("CreateWallet promise has been resolved or rejected");
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

    if (forceURL !== "") {
      return <Navigate to={forceURL} />;
    }

    return (
        <>
        <div class="container">
          <section class="section">
            {/* Page */}
            <nav class="box">
              <p class="title is-2">
                <FontAwesomeIcon className="fas" icon={faPlus} />
                &nbsp;Add Wallet
              </p>

              <FormErrorBox errors={errors} />

              <p class="pb-4">Please pick a secure password:</p>

              <FormInputField
                type="text"
                label="Label (Optional)"
                name="label"
                placeholder=""
                value={label}
                errorText={errors && errors.label}
                helpText="Give this wallet a label to help you describe your accounts better."
                onChange={(e) => setLabel(e.target.value)}
                isRequired={true}
                maxWidth="500px"
              />

              <FormInputField
                type="password"
                label="Password"
                name="password"
                placeholder=""
                value={password}
                errorText={errors && errors.password}
                helpText=""
                onChange={(e) => setPassword(e.target.value)}
                isRequired={true}
                maxWidth="500px"
              />

              <FormInputField
                type="password"
                label="Password Repeated"
                name="passwordRepeated"
                placeholder=""
                value={passwordRepeated}
                errorText={errors && errors.passwordRepeated}
                helpText=""
                onChange={(e) => setPasswordRepeated(e.target.value)}
                isRequired={true}
                maxWidth="500px"
              />

              <div class="columns pt-5" style={{alignSelf: "flex-start"}}>
                <div class="column is-half ">
                  <Link
                    class="button is-fullwidth-mobile"
                    to={`/wallets`}
                  >
                    <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                    &nbsp;Cancel & Go Back
                  </Link>
                </div>
                <div class="column is-half has-text-right">
                  <button
                    class="button is-primary is-fullwidth-mobile"
                    onClick={onSubmitClick}
                  >
                    <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                    &nbsp;Submit
                  </button>
                </div>
              </div>
            </nav>
          </section>
          </div>
        </>
    )
}

export default CreateWalletView
