import React, { useEffect } from "react";
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArrowLeft,
  faTasks,
  faGauge,
  faArrowRight,
  faUsers,
  faBarcode,
  faNewspaper,
  faBookOpen,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { topAlertMessageState, topAlertStatusState } from "../../../AppState";

function RetailerSubmissionLaunchpad() {
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

  ////
  //// API.
  ////

  ////
  //// Event handling.
  ////

  ////
  //// Misc.
  ////

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

  return (
    <>
      <div class="container">
        <section class="section">
          {/* Desktop Breadcrumbs */}
          <nav class="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to="/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Dashboard
                </Link>
              </li>
              <li class="is-active">
                <Link to="" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faTasks} />
                  &nbsp;Online Submissions
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/dashboard`} aria-current="page">
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
                  <FontAwesomeIcon className="fas" icon={faTasks} />
                  &nbsp;Online Submissions
                </h1>
              </div>
            </div>

            <section class="hero is-medium is-link">
              <div class="hero-body">
                <p class="title">
                  <FontAwesomeIcon className="fas" icon={faBookOpen} />
                  &nbsp;Comics
                </p>
                <p class="subtitle">
                  Submit a request to encapsulate your comics by clicking below:
                  <br />
                  <br />
                  <Link to={"/submissions/comics/add/search"}>
                    New&nbsp;
                    <FontAwesomeIcon className="fas" icon={faArrowRight} />
                  </Link>
                  <br />
                  <br />
                  <Link to={"/submissions/comics"}>
                    View&nbsp;
                    <FontAwesomeIcon className="fas" icon={faArrowRight} />
                  </Link>
                </p>
              </div>
            </section>

            <section class="hero is-medium is-info">
              <div class="hero-body">
                <p class="title">
                  <FontAwesomeIcon className="fas" icon={faNewspaper} />
                  &nbsp;Cards
                </p>
                <p class="subtitle">
                  Submit a request to encapsulate your cards by clicking below:
                  <br />
                  <br />
                  <Link to={"/submissions/cards/add/search"}>
                    New&nbsp;
                    <FontAwesomeIcon className="fas" icon={faArrowRight} />
                  </Link>
                  <br />
                  <br />
                  <Link to={"/submissions/cards"}>
                    View&nbsp;
                    <FontAwesomeIcon className="fas" icon={faArrowRight} />
                  </Link>
                </p>
              </div>
            </section>
          </nav>
        </section>
      </div>
    </>
  );
}

export default RetailerSubmissionLaunchpad;
