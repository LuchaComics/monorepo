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
import FormCheckboxField from "../Reusable/FormCheckboxField";
import {
    GetNFTStoreSettingsFromPreferences,
    SaveNFTStoreSettings,
    ShutdownApp
} from "../../../wailsjs/go/main/App";
import PageLoadingContent from "../Reusable/PageLoadingContent";


function NFTStoreSettingsSetupView() {

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [isLoading, setIsLoading] = useState(false);
    const [useDefaultLocation, setUseDefaultLocation] = useState(1);
    const [forceURL, setForceURL] = useState("");
    const [apiVersion, setApiVersion] = useState("2006-03-01");
    const [endpoint, setEndpoint] = useState("https://s3.filebase.com");
    const [secretAccessKey, setSecretAccessKey] = useState("");
    const [accessKeyId, setAccessKeyId] = useState("");
    const [region, setRegion] = useState("us-east-1");
    const [s3ForcePathStyle, setS3ForcePathStyle] = useState(true);
    const [showCancelWarning, setShowCancelWarning] = useState(false);

    ////
    //// Event handling.
    ////

    const EndpointCallback = (result) => Endpoint(result);

    ////
    //// API.
    ////

    const onSubmitClick = (e) => {
        e.preventDefault();
        setErrors({});
        setIsLoading(true);

        let config = {
            apiVersion: apiVersion,
            accessKeyId: accessKeyId,
            secretAccessKey: secretAccessKey,
            endpoint: endpoint,
            region: region,
            s3ForcePathStyle: s3ForcePathStyle ? "true" : "false",
        };

        // // Submit the `endpoint` value to our backend.
        SaveNFTStoreSettings(config).then( (result) => {
            console.log("result:", result);
            setForceURL("/startup")
        }).catch((errorJsonString)=>{
            console.log("errRes:", errorJsonString);
            const errorObject = JSON.parse(errorJsonString);
            console.log("errorObject:", errorObject);

            let err = {};
            if (errorObject.apiVersion != "") {
                err.apiVersion = errorObject.apiVersion;
            }
            if (errorObject.accessKeyId != "") {
                err.accessKeyId = errorObject.accessKeyId;
            }
            if (errorObject.secretAccessKey != "") {
                err.secretAccessKey = errorObject.secretAccessKey;
            }
            if (errorObject.endpoint != "") {
                err.endpoint = errorObject.endpoint;
            }
            if (errorObject.region != "") {
                err.region = errorObject.region;
            }
            if (errorObject.s3ForcePathStyle != "") {
                err.s3ForcePathStyle = errorObject.s3ForcePathStyle;
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
        GetNFTStoreSettingsFromPreferences().then( (resp)=>{
            console.log("GetNFTStoreSettingsFromPreferences: resp:", resp);
        });
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

                <p class="pb-4">Next you will need to configure how to connect to the <b>FileBase</b>.</p>
                <p class="pb-4">ComicCoin Registry requires the following fields to be filled out.</p>

                <FormInputField
                  label="API Version"
                  name="apiVersion"
                  placeholder="2006-03-01"
                  value={apiVersion}
                  errorText={errors && errors.apiVersion}
                  helpText=""
                  onChange={(e) => setApiVersion(e.target.value)}
                  isRequired={true}
                  maxWidth="500px"
                />


                <FormInputField
                  label="Access Key ID"
                  name="accessKeyId"
                  placeholder="Filebase-Access-Key"
                  value={accessKeyId}
                  errorText={errors && errors.accessKeyId}
                  helpText=""
                  onChange={(e) => setAccessKeyId(e.target.value)}
                  isRequired={true}
                  maxWidth="500px"
                />

                <FormTextareaField
                  label="Secret Access Key"
                  name="secretAccessKey"
                  placeholder="Filebase-Secret-Key"
                  value={secretAccessKey}
                  errorText={errors && errors.secretAccessKey}
                  helpText=""
                  onChange={(e) => setSecretAccessKey(e.target.value)}
                  isRequired={true}
                  rows={5}
                />

                <FormInputField
                  label="Endpoint"
                  name="endpoint"
                  placeholder="https://s3.filebase.com"
                  value={endpoint}
                  errorText={errors && errors.endpoint}
                  helpText=""
                  onChange={(e) => setEndpoint(e.target.value)}
                  isRequired={true}
                  maxWidth="500px"
                />

                <FormInputField
                  label="Region"
                  name="region"
                  placeholder="us-east-1"
                  value={region}
                  errorText={errors && errors.region}
                  helpText=""
                  onChange={(e) => setRegion(e.target.value)}
                  isRequired={true}
                  maxWidth="500px"
                />

                <FormCheckboxField
                  label="S3 Force Path Style"
                  name="s3ForcePathStyle"
                  checked={s3ForcePathStyle}
                  errorText={errors && errors.s3ForcePathStyle}
                  onChange={(e) => setS3ForcePathStyle(!s3ForcePathStyle)}
                  maxWidth="180px"
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

export default NFTStoreSettingsSetupView
