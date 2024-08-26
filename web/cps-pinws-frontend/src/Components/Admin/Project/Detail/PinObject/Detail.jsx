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
  faProjectCircle,
  faGauge,
  faPencil,
  faProjectDiagram,
  faIdCard,
  faAddressBook,
  faContactCard,
  faChartPie,
  faCogs,
  faEye,
  faArrowLeft,
  faLocationPin,
  faFile,
  faDownload,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { getPinObjectDetailAPI, deletePinObjectAPI } from "../../../../../API/PinObject";
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
import FormRowText from "../../../../Reusable/FormRowText";

function AdminProjectPinObjectDetail() {
  ////
  //// URL Parameters.
  ////

  const { id, rid } = useParams();

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
  const [name, setName] = useState("");
  const [cid, setCID] = useState("");
  const [description, setDescription] = useState("");
  const [objectUrl, setObjectUrl] = useState("");
  const [selectedPinObjectRequestIDForDeletion, setSelectedPinObjectRequestIDForDeletion] =
    useState("");

  ////
  //// Event handling.
  ////

  const onDeleteConfirmButtonClick = (e) => {
    e.preventDefault();
    console.log("onDeleteConfirmButtonClick"); // For debugging purposes only.

    deletePinObjectAPI(
      selectedPinObjectRequestIDForDeletion,
      onPinObjectDeleteSuccess,
      onPinObjectDeleteError,
      onPinObjectDeleteDone,
      onUnauthorized,
    );
  };

  ////
  //// API.
  ////

  // --- Get Details ---

  function onAdminProjectPinObjectDetailSuccess(response) {
    // For debugging purposes only.
    console.log("onAdminProjectPinObjectDetailSuccess: Starting...");
    console.log(response);
    setName(response.name);
    setCID(response.cid);
    setObjectUrl(response.objectUrl);
  }

  function onAdminProjectPinObjectDetailError(apiErr) {
    console.log("onAdminProjectPinObjectDetailError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onAdminProjectPinObjectDetailError: Delayed for 2 seconds.");
      console.log(
        "onAdminProjectPinObjectDetailError: topAlertMessage, topAlertStatus:",
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

  function onAdminProjectPinObjectDetailDone() {
    console.log("onAdminProjectPinObjectDetailDone: Starting...");
    setFetching(false);
  }

  // --- Deletion --- //

  function onPinObjectDeleteSuccess(response) {
    console.log("onPinObjectDeleteSuccess: Starting..."); // For debugging purposes only.

    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("Pin deleted");
    setTimeout(() => {
      console.log(
        "onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Redirect the project to the project pins page.
    setForceURL("/admin/project/" + id + "/pins");
  }

  function onPinObjectDeleteError(apiErr) {
    console.log("onPinObjectDeleteError: Starting..."); // For debugging purposes only.
    setErrors(apiErr);

    // Update notification.
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed deleting");
    setTimeout(() => {
      console.log(
        "onPinObjectDeleteError: topAlertMessage, topAlertStatus:",
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

  function onPinObjectDeleteDone() {
    console.log("onPinObjectDeleteDone: Starting...");
    setFetching(false);
  }

  // --- All --- //

  const onUnauthorized = () => {
    setForceURL("/login?unauthorized=true"); // If token expired or project is not logged in, redirect back to login.
  };

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.

      getPinObjectDetailAPI(
        rid,
        onAdminProjectPinObjectDetailSuccess,
        onAdminProjectPinObjectDetailError,
        onAdminProjectPinObjectDetailDone,
        onUnauthorized,
      );
    }

    return () => {
      mounted = false;
    };
  }, [rid]);

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
                <Link to="/admin/projects" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faProjectDiagram} />
                  &nbsp;Projects
                </Link>
              </li>
              <li class="">
                <Link to={`/admin/project/${id}/pins`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail (Pins)
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faLocationPin} />
                  &nbsp;Pin
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/admin/project/${id}/pins`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Detail (Pins)
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modals */}
          <div
            class={`modal ${selectedPinObjectRequestIDForDeletion ? "is-active" : ""}`}
          >
            <div class="modal-background"></div>
            <div class="modal-card">
              <header class="modal-card-head">
                <p class="modal-card-title">Are you sure?</p>
                <button
                  class="delete"
                  aria-label="close"
                  onClick={(e)=>{
                      e.preventDefault();
                      setSelectedPinObjectRequestIDForDeletion("");
                  }}
                ></button>
              </header>
              <section class="modal-card-body">
                You are about to <b>delete</b> this pin; it will no
                longer appear on your dashboard This action cannot be undone. Are you sure
                you would like to continue?
              </section>
              <footer class="modal-card-foot">
                <button
                  class="button is-success"
                  onClick={onDeleteConfirmButtonClick}
                >
                  Confirm
                </button>
                <button
                  class="button"
                  onClick={(e)=>{
                      e.preventDefault();
                      setSelectedPinObjectRequestIDForDeletion("");
                  }}
                >
                  Cancel
                </button>
              </footer>
            </div>
          </div>

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faLocationPin} />
              &nbsp;Pin
            </p>

            {/* <p class="pb-4 has-text-grey">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Submitting..."} />
            ) : (
              <>
                <FormErrorBox errors={errors} />
                <div class="container">
                  <p class="subtitle is-4 pt-4">
                    <FontAwesomeIcon className="fas" icon={faEye} />
                    &nbsp;Meta Information
                  </p>
                  <hr />

                  <FormRowText label="Name" value={name} helpText="" />
                  <FormRowText label="CID" value={cid} helpText="" />

                  <p class="subtitle is-4 pt-4">
                    <FontAwesomeIcon className="fas" icon={faFile} />
                    &nbsp;Data
                  </p>
                  <hr />
                  <p class="has-text-grey pb-4">
                    Click the following "Download File" button to start
                    downloading a copy of the pinobject to your computer.
                  </p>

                  <section class="hero has-background-white-ter">
                    <div class="hero-body">
                      <p class="subtitle">
                        <div class="has-text-centered">
                          <a
                            href={objectUrl}
                            target="_blank"
                            rel="noreferrer"
                            class="button is-large is-success is-hidden-touch"
                          >
                            <FontAwesomeIcon
                              className="fas"
                              icon={faDownload}
                            />
                            &nbsp;Download File
                          </a>
                          <a
                            href={objectUrl}
                            target="_blank"
                            rel="noreferrer"
                            class="button is-large is-success is-fullwidth is-hidden-desktop"
                          >
                            <FontAwesomeIcon
                              className="fas"
                              icon={faDownload}
                            />
                            &nbsp;Download File
                          </a>
                        </div>
                      </p>
                    </div>
                  </section>

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <Link
                        to={`/admin/project/${id}/pins`}
                        class="button is-medium is-fullwidth-mobile"
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back
                      </Link>
                    </div>
                    <div class="column is-half has-text-right">
                      <Link
                        to={`/admin/project/${id}/pin/${rid}/edit`}
                        class="button is-medium is-warning is-fullwidth-mobile"
                      >
                        <FontAwesomeIcon className="fas" icon={faPencil} />
                        &nbsp;Edit
                      </Link>
                      &nbsp;
                      &nbsp;
                      <button
                          onClick={(e)=>{
                              e.preventDefault();
                              setSelectedPinObjectRequestIDForDeletion(rid);
                          }}
                          class="button is-medium is-danger is-fullwidth-mobile"
                        >
                          <FontAwesomeIcon className="fas" icon={faPencil} />
                          &nbsp;Delete
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

export default AdminProjectPinObjectDetail;
