import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import Scroll from "react-scroll";
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
  faEllipsis
} from "@fortawesome/free-solid-svg-icons";

import FormErrorBox from "../Reusable/FormErrorBox";
import FormTextareaField from "../Reusable/FormTextareaField";
import FormRadioField from "../Reusable/FormRadioField";
import FormInputField from "../Reusable/FormInputField";
import FormInputFieldWithButton from "../Reusable/FormInputFieldWithButton";
import {
    GetNFTStoreRemoteAddressFromPreferences,
    GetNFTStoreAPIKeyFromPreferences,
    GetDataDirectoryFromDialog,
    SaveNFTStoreConfigVariables,
    ShutdownApp
} from "../../../wailsjs/go/main/App";
import PageLoadingContent from "../Reusable/PageLoadingContent";


function NFTSToreSetupView() {

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [isLoading, setIsLoading] = useState(false);
    const [useDefaultLocation, setUseDefaultLocation] = useState(1);
    const [forceURL, setForceURL] = useState("");
    const [remoteAddress, setRemoteAddress] = useState("http://127.0.0.1:8080");
    const [apiKey, setApiKey] = useState("");
    const [showCancelWarning, setShowCancelWarning] = useState(false);

    ////
    //// Event handling.
    ////

    const setRemoteAddressCallback = (result) => setRemoteAddress(result);

    ////
    //// API.
    ////

    const onSubmitClick = (e) => {
        e.preventDefault();
        setErrors({});
        setIsLoading(true);

        // Submit the `remoteAddress` value to our backend.
        SaveNFTStoreConfigVariables(apiKey, remoteAddress).then( (result) => {
            console.log("result:", result);
            setForceURL("/startup")
        }).catch((errorJsonString)=>{
            console.log("errRes:", errorJsonString);
            const errorObject = JSON.parse(errorJsonString);
            console.log("errorObject:", errorObject);

            let err = {};
            if (errorObject.nftStoreRemoteAddress != "") {
                err.remoteAddress = errorObject.nftStoreRemoteAddress;
            }
            if (errorObject.nftStoreAPIKey != "") {
                err.apiKey = errorObject.nftStoreAPIKey;
            }
            setErrors(err);

            // The following code will cause the screen to scroll to the top of
            // the page. Please see ``react-scroll`` for more information:
            // https://github.com/fisshy/react-scroll
            var scroll = Scroll.animateScroll;
            scroll.scrollToTop();
        }).finally(()=>{
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
            GetNFTStoreRemoteAddressFromPreferences().then( (resp)=>{
                console.log("OnStartup: RemoteAddress:", resp);
                setRemoteAddress(resp);
            })
            GetNFTStoreAPIKeyFromPreferences().then( (resp)=>{
                console.log("OnStartup: API Key:", resp);
                setApiKey(resp);
            })
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

    if (isLoading) {
        return (
            <PageLoadingContent displayMessage="Please wait..." />
        )
    }

    return (
        <>
          {/* Modals */}
          <div class={`modal ${showCancelWarning ? "is-active" : ""}`}>
            <div class="modal-background"></div>
            <div class="modal-card">
              <header class="modal-card-head">
                <p class="modal-card-title">Are you sure?</p>
                <button
                  class="delete"
                  aria-label="close"
                  onClick={(e) => setShowCancelWarning(false)}
                ></button>
              </header>
              <section class="modal-card-body">
                Your tenant record will be cancelled and your work will be lost. This
                cannot be undone. Do you want to continue?
              </section>
              <footer class="modal-card-foot">
                <Link class="button is-medium is-success" onClick={(e)=>{
                    ShutdownApp()
                }}>
                  Yes
                </Link>&nbsp;&nbsp;
                <button
                  class="button is-medium"
                  onClick={(e) => setShowCancelWarning(false)}
                >
                  No
                </button>
              </footer>
            </div>
          </div>

          <div class="container">
            <section class="section">
              {/* Page */}
              <nav class="box">
                <p class="title is-2">
                  <FontAwesomeIcon className="fas" icon={faHandHoldingHeart} />
                  &nbsp;Welcome to ComicCoin Registry.
                </p>

                <FormErrorBox errors={errors} />

                <p class="pb-4">Next you will need to configure how to connect to the <b>NFT Store</b>.</p>
                <p class="pb-4">ComicCoin Registry requires the following fields to be filled out.</p>

                <FormInputField
                  label="Remote Address"
                  name="remoteAddress"
                  placeholder=""
                  value={remoteAddress}
                  errorText={errors && errors.remoteAddress}
                  helpText="Please enter the address the NFT Store can be reached at."
                  onChange={(e) => setRemoteAddress(e.target.value)}
                  isRequired={true}
                  maxWidth="500px"
                />

                <FormTextareaField
                  label="API Key"
                  name="apiKey"
                  placeholder=""
                  value={apiKey}
                  errorText={errors && errors.apiKey}
                  helpText="Please keep this key safe!"
                  onChange={(e) => setApiKey(e.target.value)}
                  isRequired={true}
                  rows={5}
                />

                <div class="columns pt-5" style={{alignSelf: "flex-start"}}>
                  <div class="column is-half ">
                    <button
                      class="button is-fullwidth-mobile"
                      onClick={(e) => setShowCancelWarning(true)}
                    >
                      <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                      &nbsp;Cancel
                    </button>
                  </div>
                  <div class="column is-half has-text-right">
                    <button
                      class="button is-primary is-fullwidth-mobile"
                      onClick={onSubmitClick}
                    >
                      <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                      &nbsp;Save
                    </button>
                  </div>
                </div>

              </nav>
            </section>
          </div>
        </>
    )
}

export default NFTSToreSetupView
