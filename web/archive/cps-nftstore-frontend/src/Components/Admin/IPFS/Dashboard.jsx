import React, { useEffect, useState } from "react";
import { Link, useSearchParams, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faGauge,
  faArrowRight,
  faUsers,
  faBarcode,
  faCubes,
  faCloud,
  faArrowLeft,
  faIdCard
} from "@fortawesome/free-solid-svg-icons";
import FormRowText from "../../Reusable/FormRowText";
import { useRecoilState } from "recoil";

import { topAlertMessageState, topAlertStatusState } from "../../../AppState";
import { getIpfsInfoAPI } from "../../../API/IPFS";


function AdminIPFSDashboard() {
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
  const [ipfsInfo, setIpfsInfo] = useState("");


  ////
  //// API.
  ////

  function onIpfsInfoSuccess(response) {
    console.log("onIpfsInfoSuccess: Starting...");
      console.log("onIpfsInfoSuccess: response:", response);
    setIpfsInfo(response);
  }

  function onIpfsInfoError(apiErr) {
    console.log("onIpfsInfoError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onIpfsInfoDone() {
    console.log("onIpfsInfoDone: Starting...");
    setFetching(false);
  }

  // --- All --- //

  const onUnauthorized = () => {
    setForceURL("/login?unauthorized=true"); // If token expired or user is not logged in, redirect back to login.
  };

  ////
  //// Event handling.
  ////

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.

      setFetching(true);
      getIpfsInfoAPI(
        onIpfsInfoSuccess,
        onIpfsInfoError,
        onIpfsInfoDone,
        onUnauthorized,
      );
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
            <li class="is-active">
              <Link aria-current="page">
                <FontAwesomeIcon className="fas" icon={faCloud} />
                &nbsp;IPFS
              </Link>
            </li>
          </ul>
        </nav>

        {/* Mobile Breadcrumbs */}
        <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
          <ul>
            <li class="">
              <Link to={`/admin/dashboard`} aria-current="page">
                <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                &nbsp;Back to Dashboard
              </Link>
            </li>
          </ul>
        </nav>

        {/* Page */}
        <nav class="box">
          <div class="columns">
            <div class="column">
              <h1 class="title is-4">
                <FontAwesomeIcon className="fas" icon={faCloud} />
                &nbsp;IPFS
              </h1>
            </div>
            </div>
              {ipfsInfo && <div class="container">
                  {/* Title */}
                  <p class="subtitle is-6">
                    <FontAwesomeIcon className="fas" icon={faIdCard} />
                    &nbsp;Information
                  </p>
                  <hr />

                  <FormRowText
                    label="ID"
                    value={ipfsInfo.id}
                    helpText=""
                  />

                  <FormRowText
                    label="Public Key"
                    value={ipfsInfo.publicKey}
                    helpText=""
                  />

                  <FormRowText
                    label="Protocol Version"
                    value={ipfsInfo.protocolVersion}
                    helpText=""
                  />

                  <FormRowText
                    label="Agent Version"
                    value={ipfsInfo.agentVersion}
                    helpText=""
                  />

                  <FormRowText
                    label="Addresses"
                    value={ipfsInfo.addresses}
                    helpText=""
                  />
              </div>}
        </nav>


        </section>
      </div>
    </>
  );
}

export default AdminIPFSDashboard;
