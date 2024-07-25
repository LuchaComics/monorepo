import React, { useState, useEffect } from "react";
import { Link, Navigate, useSearchParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faArrowLeft,
  faTasks,
  faTachometer,
  faPlus,
  faTimesCircle,
  faCheckCircle,
  faGauge,
  faUsers,
  faEye,
  faCube,
  faMagnifyingGlass,
  faBalanceScale,
  faUser,
  faCogs,
  faBookOpen,
  faNewspaper,
  faArrowRight,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import FormErrorBox from "../../Reusable/FormErrorBox";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";

function RetailerSubmissionPickTypeForAdd() {
  ////
  //// URL Parameters.
  ////

  const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140
  const customerID = searchParams.get("customer_id");
  const customerName = searchParams.get("customer_name");
  const orgID = searchParams.get("store_id");
  // const orgName = searchParams.get("store_name");

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
  //// Event handling.
  ////

  ////
  //// API.
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

  // Render the JSX content.
  return (
    <>
      <div class="container">
        <section class="section">
          {customerName === null ? (
            <>
              {/* Desktop Breadcrumbs */}
              <nav class="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
                <ul>
                  <li class="">
                    <Link to="/dashboard" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faGauge} />
                      &nbsp;Dashboard
                    </Link>
                  </li>
                  <li class="">
                    <Link to="/submissions" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faTasks} />
                      &nbsp;Online Submissions
                    </Link>
                  </li>
                  <li class="is-active">
                    <Link aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faPlus} />
                      &nbsp;Add
                    </Link>
                  </li>
                </ul>
              </nav>

              {/* Mobile Breadcrumbs */}
              <nav
                class="breadcrumb is-hidden-desktop"
                aria-label="breadcrumbs"
              >
                <ul>
                  <li class="">
                    <Link to={`/submissions`} aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                      &nbsp;Back to Online Submissions
                    </Link>
                  </li>
                </ul>
              </nav>
            </>
          ) : (
            <>
              {/* Desktop Breadcrumbs */}
              <nav class="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
                <ul>
                  <li class="">
                    <Link to="/dashboard" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faGauge} />
                      &nbsp;Dashboard
                    </Link>
                  </li>
                  <li class="">
                    <Link to="/customers" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faUsers} />
                      &nbsp;Customers
                    </Link>
                  </li>
                  <li class="">
                    <Link
                      to={`/customer/${customerID}/comics`}
                      aria-current="page"
                    >
                      <FontAwesomeIcon className="fas" icon={faEye} />
                      &nbsp;Detail (Comics)
                    </Link>
                  </li>
                  <li class="is-active">
                    <Link aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faPlus} />
                      &nbsp;Add
                    </Link>
                  </li>
                </ul>
              </nav>

              {/* Mobile Breadcrumbs */}
              <nav
                class="breadcrumb is-hidden-desktop"
                aria-label="breadcrumbs"
              >
                <ul>
                  <li class="">
                    <Link
                      to={`/customer/${customerID}/comics`}
                      aria-current="page"
                    >
                      <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                      &nbsp;Back to Detail (Comics)
                    </Link>
                  </li>
                </ul>
              </nav>
            </>
          )}

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faPlus} />
              &nbsp;New Online Submission
            </p>
            <div class="container">
              <p class="has-text-grey pb-4">
                Please select the type of collectible product you would like to
                submit to CPS.
              </p>

              <p class="subtitle is-6">
                <FontAwesomeIcon className="fas" icon={faCube} />
                &nbsp;Product Type
              </p>
              <hr />

              <section class="hero is-medium is-link">
                <div class="hero-body">
                  <p class="title">
                    <FontAwesomeIcon className="fas" icon={faBookOpen} />
                    &nbsp;Comics
                  </p>
                  <p class="subtitle">
                    Currently we accept submissions of comics up to 64 pages.
                    Comics must be standard US Golden Age, Silver Age, Bronze
                    Age or Modern Age sizes (no oversize submissions can be
                    processed at this time).
                    <br />
                    <br />
                    {customerName === null ? (
                      <Link to={`/submissions/comics/add/search`}>
                        Select&nbsp;
                        <FontAwesomeIcon className="fas" icon={faArrowRight} />
                      </Link>
                    ) : (
                      <Link
                        to={`/submissions/comics/add?customer_id=${customerID}&customer_name=${customerName}`}
                      >
                        Select&nbsp;
                        <FontAwesomeIcon className="fas" icon={faArrowRight} />
                      </Link>
                    )}
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
                    Currently we accept of standard size, non-sports cards.
                    <br />
                    <br />
                    <Link>
                      <i>Coming soon</i>
                    </Link>
                  </p>
                </div>
              </section>

              <div class="columns pt-5">
                <div class="column is-half">
                  {customerName === null ? (
                    <>
                      <Link
                        to={`/submissions/comics`}
                        class="button is-medium is-hidden-touch"
                      >
                        <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                        &nbsp;Cancel
                      </Link>
                      <Link
                        to={`/submissions/comics`}
                        class="button is-medium is-fullwidth is-hidden-desktop"
                      >
                        <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                        &nbsp;Cancel
                      </Link>
                    </>
                  ) : (
                    <>
                      <Link
                        to={`/customer/${customerID}/sub`}
                        class="button is-medium is-hidden-touch"
                      >
                        <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                        &nbsp;Cancel
                      </Link>
                      <Link
                        to={`/customer/${customerID}/sub`}
                        class="button is-medium is-fullwidth is-hidden-desktop"
                      >
                        <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                        &nbsp;Cancel
                      </Link>
                    </>
                  )}
                </div>
                <div class="column is-half has-text-right">
                  {/*
                                    <button class="button is-medium is-primary is-hidden-touch" onClick={onSubmitClick}><FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;Save</button>
                                    <button class="button is-medium is-primary is-fullwidth is-hidden-desktop" onClick={onSubmitClick}><FontAwesomeIcon className="fas" icon={faCheckCircle} />&nbsp;Save</button>
                                    */}
                </div>
              </div>
            </div>
          </nav>
        </section>
      </div>
    </>
  );
}

export default RetailerSubmissionPickTypeForAdd;
