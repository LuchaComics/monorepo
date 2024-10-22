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
  faEllipsis
} from "@fortawesome/free-solid-svg-icons";

import FormErrorBox from "../../Reusable/FormErrorBox";
import FormRadioField from "../../Reusable/FormRadioField";
import FormInputField from "../../Reusable/FormInputField";
import FormInputFieldWithButton from "../../Reusable/FormInputFieldWithButton";
import {GetDataDirectoryFromDialog, SaveDataDirectory} from "../../../../wailsjs/go/main/App";


function PickStorageLocationOnStartupView() {

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [useDefaultLocation, setUseDefaultLocation] = useState(1);
    const [forceURL, setForceURL] = useState("");
    const [dataDirectory, setDataDirectory] = useState("./ComicCoin");
    const [showCancelWarning, setShowCancelWarning] = useState(false);

    ////
    //// Event handling.
    ////

    const setDataDirectoryCallback = (result) => setDataDirectory(result);

    ////
    //// API.
    ////

    const onSubmitClick = (e) => {
        e.preventDefault();

        // Submit the `dataDirectory` value to our backend.
        SaveDataDirectory(dataDirectory).then( (result) => {
            setForceURL("/dashboard")
        })
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
                <Link class="button is-medium is-success" to={`/admin/tenants`}>
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
                  &nbsp;Welcome to ComicCoin Core.
                </p>

                <FormErrorBox errors={errors} />

                <p class="pb-4">As this is the first time the program is launched, you can choose where ComicCoin Core will store its data</p>
                <p class="pb-4">ComicCoin Core will download and store a copy of the ComicCoin block chain. Approximately 1 MB of data will be stored in this directory. The wallet will also be stored in this directory.</p>

                <FormRadioField
                  label="Are you currently submitting to any other grading companies?"
                  name="hasOtherGradingService"
                  value={useDefaultLocation}
                  opt1Value={1}
                  opt1Label="Use the default data directory."
                  opt2Value={2}
                  opt2Label="Use a custom data directory."
                  errorText={errors && errors.useDefaultLocation}
                  onChange={(e) =>
                    setUseDefaultLocation(parseInt(e.target.value))
                  }
                  maxWidth="180px"
                  hasOptPerLine={true}
                />

                <FormInputFieldWithButton
                  label="Data Directory"
                  name="dataDirectory"
                  placeholder="Data Directory"
                  value={dataDirectory}
                  errorText={errors && errors.dataDirectory}
                  helpText=""
                  onChange={(e) => setDataDirectory(e.target.value)}
                  isRequired={true}
                  maxWidth="500px"
                  disabled={useDefaultLocation == 1}
                  buttonLabel={<><FontAwesomeIcon className="fas" icon={faEllipsis} /></>}
                  onButtonClick={(e) =>
                    GetDataDirectoryFromDialog().then(setDataDirectoryCallback)
                  }
                />

                <p class="pb-4">When you dick OK, ComicCoin Core will begin to download and process the full ComicCoin block chain (1 MB) starting with the earliest transactions in 2024 when ComicCoin initially launched.</p>

                <p class="pb-4">This initial synchronisation is very demanding, and may expose hardware problems with your computer that had previously gone unnoticed. Each time you run ComicCoin Core, it will continue downloading where it left off.</p>

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

export default PickStorageLocationOnStartupView
