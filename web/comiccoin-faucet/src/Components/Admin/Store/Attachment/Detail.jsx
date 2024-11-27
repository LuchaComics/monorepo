import React, { useState, useEffect } from "react";
import { Link, Navigate, useParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faBuilding,
  faTachometer,
  faPlus,
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
  faCogs,
  faEye,
  faArrowLeft,
  faFile,
  faDownload,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { getAttachmentDetailAPI } from "../../../../API/Attachment";
import FormErrorBox from "../../../Reusable/FormErrorBox";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../AppState";
import FormRowText from "../../../Reusable/FormRowText";

function AdminStoreAttachmentDetail() {
  ////
  //// URL Parameters.
  ////

  const { id, aid } = useParams();

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
  const [description, setDescription] = useState("");
  const [objectUrl, setObjectUrl] = useState("");

  ////
  //// Event handling.
  ////

  ////
  //// API.
  ////

  function onAdminStoreAttachmentDetailSuccess(response) {
    // For debugging purposes only.
    console.log("onAdminStoreAttachmentDetailSuccess: Starting...");
    console.log(response);
    setName(response.name);
    setDescription(response.description);
    setObjectUrl(response.objectUrl);
  }

  function onAdminStoreAttachmentDetailError(apiErr) {
    console.log("onAdminStoreAttachmentDetailError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onAdminStoreAttachmentDetailError: Delayed for 2 seconds.");
      console.log(
        "onAdminStoreAttachmentDetailError: topAlertMessage, topAlertStatus:",
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

  function onAdminStoreAttachmentDetailDone() {
    console.log("onAdminStoreAttachmentDetailDone: Starting...");
    setFetching(false);
  }

  // --- All --- //

  const onUnauthorized = () => {
    setForceURL("/login?unauthorized=true"); // If token expired or user is not logged in, redirect back to login.
  };

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.

      getAttachmentDetailAPI(
        aid,
        onAdminStoreAttachmentDetailSuccess,
        onAdminStoreAttachmentDetailError,
        onAdminStoreAttachmentDetailDone,
        onUnauthorized,
      );
    }

    return () => {
      mounted = false;
    };
  }, [aid]);

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
                <Link to="/admin/stores" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faBuilding} />
                  &nbsp;Stores
                </Link>
              </li>
              <li class="">
                <Link to={`/admin/store/${id}/attachments`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail (Attachments)
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faFile} />
                  &nbsp;Attachment
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/admin/store/${id}/attachments`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Detail (Attachments)
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faFile} />
              &nbsp;Attachment
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

                  <FormRowText
                    label="Description"
                    value={description}
                    helpText=""
                  />

                  <p class="subtitle is-4 pt-4">
                    <FontAwesomeIcon className="fas" icon={faFile} />
                    &nbsp;Data
                  </p>
                  <hr />
                  <p class="has-text-grey pb-4">
                    Click the following "Download File" button to start
                    downloading a copy of the attachment to your computer.
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
                        to={`/admin/store/${id}/attachments`}
                        class="button is-medium is-hidden-touch"
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back
                      </Link>
                      <Link
                        to={`/admin/store/${id}/attachments`}
                        class="button is-medium is-fullwidth is-hidden-desktop"
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back
                      </Link>
                    </div>
                    <div class="column is-half has-text-right">
                      <Link
                        to={`/admin/store/${id}/attachment/${aid}/edit`}
                        class="button is-medium is-warning is-hidden-touch"
                      >
                        <FontAwesomeIcon className="fas" icon={faPencil} />
                        &nbsp;Edit
                      </Link>
                      <Link
                        to={`/admin/store/${id}/attachment/${aid}/edit`}
                        class="button is-medium is-warning is-fullwidth is-hidden-desktop"
                      >
                        <FontAwesomeIcon className="fas" icon={faPencil} />
                        &nbsp;Edit
                      </Link>
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

export default AdminStoreAttachmentDetail;
