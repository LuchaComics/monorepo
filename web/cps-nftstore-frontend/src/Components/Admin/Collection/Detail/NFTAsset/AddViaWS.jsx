import React, { useState, useEffect } from "react";
import { Link, Navigate, useParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faTachometer,
  faPlus,
  faTimesCircle,
  faCheckCircle,
  faCollectionCircle,
  faGauge,
  faPencil,
  faCubes,
  faIdCard,
  faAddressBook,
  faContactCard,
  faChartPie,
  faCogs,
  faEye,
  faArrowLeft,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import useLocalStorage from "../../../../../Hooks/useLocalStorage";
import { postIpfsAddFileAPI } from "../../../../../API/IPFS";
import FormErrorBox from "../../../../Reusable/FormErrorBox";
import FormInputField from "../../../../Reusable/FormInputField";
import FormTextareaField from "../../../../Reusable/FormTextareaField";
import FormRadioField from "../../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../../Reusable/FormSelectField";
import FormCheckboxField from "../../../../Reusable/FormCheckboxField";
import FormCountryField from "../../../../Reusable/FormCountryField";
import FormRegionField from "../../../../Reusable/FormRegionField";
import PageLoadingContent from "../../../../Reusable/PageLoadingContent";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../../AppState";

function AdminCollectionNFTAssetAddViaWebService() {
  ////
  //// URL Parameters.
  ////

  const { id } = useParams();

  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [selectedFile, setSelectedFile] = useState(null);
  const [apiKey, setName] = useState("");

  ////
  //// Event handling.
  ////

  const onHandleFileChange = (event) => {
    setSelectedFile(event.target.files[0]);
  };

  const onSubmitClick = (e) => {
    e.preventDefault(); // Prevent default form submission behavior
    console.log("onSubmitClick: Starting...");
    setFetching(true);
    setErrors({});

    // Ensure that a file is selected
      if (!selectedFile) {
        console.error('No file selected');
        setErrors({ file: 'No file selected' });
        setFetching(false);
        return;
      }


    // Log the selected file to inspect its properties
    console.log('Selected file:', selectedFile);

    const formData = new FormData();
    formData.append("file", selectedFile);

    // Extract filename.
    const filename = selectedFile.name;
    console.log('Filename:', filename);

    const mimeType = selectedFile.type || "application/octet-stream"; // Fallback to octet-stream if type is not detected
    console.log('MIME Type:', mimeType);


    // Convert the FormData object to binary data and pass it directly
     // Note: You may need to adjust the handling depending on the type of file and how it needs to be processed
     const reader = new FileReader();
     reader.onload = () => {
       const fileBinaryData = reader.result;

       postIpfsAddFileAPI(
         apiKey,
         filename,
         fileBinaryData, // Pass binary data instead of FormData
         mimeType, // Pass the detected MIME type
         onAdminCollectionNFTAssetAddViaWebServiceSuccess,
         onAdminCollectionNFTAssetAddViaWebServiceError,
         onAdminCollectionNFTAssetAddViaWebServiceDone,
         onUnauthorized,
       );

       console.log("onSubmitClick: Finished.");
       setFetching(false); // Reset fetching state after API call
     };
     reader.onerror = () => {
       console.error('Error reading file');
       setErrors({ file: 'Error reading file' });
       setFetching(false);
     };

     reader.readAsArrayBuffer(selectedFile); // Read file as binary data
    console.log("onSubmitClick: Finished.");
  };

  ////
  //// API.
  ////

  function onAdminCollectionNFTAssetAddViaWebServiceSuccess(response) {
    // For debugging purposes only.
    console.log("onAdminCollectionNFTAssetAddViaWebServiceSuccess: Starting...");
    console.log(response);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Pin created");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onAdminCollectionNFTAssetAddViaWebServiceSuccess: Delayed for 2 seconds.");
      console.log(
        "onAdminCollectionNFTAssetAddViaWebServiceSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Redirect the collection to the collection pinobjects page.
    setForceURL("/admin/collection/" + id + "/pins");
  }

  function onAdminCollectionNFTAssetAddViaWebServiceError(apiErr) {
    console.log("onAdminCollectionNFTAssetAddViaWebServiceError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onAdminCollectionNFTAssetAddViaWebServiceError: Delayed for 2 seconds.");
      console.log(
        "onAdminCollectionNFTAssetAddViaWebServiceError: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onAdminCollectionNFTAssetAddViaWebServiceDone() {
    console.log("onAdminCollectionNFTAssetAddViaWebServiceDone: Starting...");
    setFetching(false);
  }

  // --- All --- //

  const onUnauthorized = () => {
    setForceURL("/login?unauthorized=true"); // If token expired or collection is not logged in, redirect back to login.
  };

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
          {/* Desktop Breadcrumbs */}
          <nav class="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to="/admin/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Admin Dashboard
                </Link>
              </li>
              <li class="">
                <Link to="/admin/collections" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faCubes} />
                  &nbsp;Collections
                </Link>
              </li>
              <li class="">
                <Link to={`/admin/collection/${id}/pins`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail (Pins)
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faPlus} />
                  &nbsp;New (via Web-Service)
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/admin/collection/${id}`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Detail
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modals */}
          {/* None */}

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faPlus} />
              &nbsp;New Pin (via Web-Service)
            </p>
            <FormErrorBox errors={errors} />

            {/* <p class="pb-4 has-text-grey">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Submitting..."} />
            ) : (
              <>
                <div class="container">
                  <FormTextareaField
                    label="API Key"
                    name="apiKey"
                    placeholder="Text input"
                    value={apiKey}
                    errorText={errors && errors.apiKey}
                    helpText=""
                    onChange={(e) => setName(e.target.value)}
                    isRequired={true}
                    maxWidth="150px"
                    rows={4}
                  />

                  <input
                    name="file"
                    type="file"
                    onChange={onHandleFileChange}
                  />
                  <br />
                  <br />

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <Link
                        to={`/admin/collection/${id}/pins`}
                        class="button is-hidden-touch"
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back
                      </Link>
                      <Link
                        to={`/admin/collection/${id}/pins`}
                        class="button is-fullwidth is-hidden-desktop"
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back
                      </Link>
                    </div>
                    <div class="column is-half has-text-right">
                      <button
                        class="button is-primary is-hidden-touch"
                        onClick={onSubmitClick}
                      >
                        <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                        &nbsp;Save
                      </button>
                      <button
                        class="button is-primary is-fullwidth is-hidden-desktop"
                        onClick={onSubmitClick}
                      >
                        <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                        &nbsp;Save
                      </button>
                    </div>
                  </div>
                </div>
              </>
            )}
          </nav>
        </section>
      </div>
    </>
  );
}

export default AdminCollectionNFTAssetAddViaWebService;
