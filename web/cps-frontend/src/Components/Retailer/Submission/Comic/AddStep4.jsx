import React, { useState, useEffect } from "react";
import { Link, Navigate, useSearchParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faBookOpen,
  faTachometer,
  faPlus,
  faDownload,
  faArrowLeft,
  faCheckCircle,
  faCheck,
  faGauge,
  faUsers,
  faEye,
} from "@fortawesome/free-solid-svg-icons";
import Select from "react-select";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import useLocalStorage from "../../../../Hooks/useLocalStorage";
import { getComicSubmissionDetailAPI } from "../../../../API/ComicSubmission";
import FormErrorBox from "../../../Reusable/FormErrorBox";
import FormInputField from "../../../Reusable/FormInputField";
import FormTextareaField from "../../../Reusable/FormTextareaField";
import FormRadioField from "../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../Reusable/FormSelectField";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../AppState";

function RetailerComicSubmissionAddStep4() {
  ////
  //// URL Arguments.
  ////

  const { id } = useParams();

  ////
  //// URL Parameters.
  ////

  const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140
  const customerID = searchParams.get("customer_id");
  const customerName = searchParams.get("customer_name");

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
  const [submission, setComicSubmission] = useState({});

  ////
  //// Event handling.
  ////

  ////
  //// API.
  ////

  function onComicSubmissionDetailSuccess(response) {
    console.log("onComicSubmissionDetailSuccess: Starting...");
    setComicSubmission(response);
  }

  function onComicSubmissionDetailError(apiErr) {
    console.log("onComicSubmissionDetailError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onComicSubmissionDetailDone() {
    console.log("onComicSubmissionDetailDone: Starting...");
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

      setFetching(true);
      getComicSubmissionDetailAPI(
        id,
        onComicSubmissionDetailSuccess,
        onComicSubmissionDetailError,
        onComicSubmissionDetailDone,
        onUnauthorized,
      );
    }

    return () => {
      mounted = false;
    };
  }, [id]);

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
          {/* Conditional breadcrumbs */}
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
                  <li class="">
                    <Link to="/submissions/comics" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faBookOpen} />
                      &nbsp;Comics
                    </Link>
                  </li>
                  <li class="is-active">
                    <Link aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faPlus} />
                      &nbsp;New
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
                    <Link to={`/submissions/comics`} aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                      &nbsp;Back to Comics
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
                      &nbsp;New
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
                    <Link to={`/submissions/comics`} aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                      &nbsp;Back to Comics
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
              &nbsp;New Online Comic Submission (Confirmation)
            </p>
            <FormErrorBox errors={errors} />

            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                <div class="container">
                  <article class="message is-success">
                    <div class="message-body">
                      <p>
                        <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                        &nbsp;PDF is ready for download.
                      </p>
                    </div>
                  </article>
                  <article>
                    <p>
                      We're excited to inform you that your PDF is now ready for
                      download. Simply click the provided link, and you'll have
                      access to your PDF file.
                    </p>

                    <p class="pb-3">
                      Once you've downloaded the PDF, please sign it and keep it
                      with the comic. This adds a personal touch and ensures the
                      authenticity of the document.
                    </p>

                    <p class="pb-3">
                      After signing, we ask you to attach the signed PDF to the
                      comic book you wish to have encapsulated. Safely packaging
                      your comic book helps protect it during transit and
                      ensures its safe arrival at our facility.
                    </p>

                    <p class="pb-3">
                      If you will be submitting this comic book for grading as
                      part of a pedigree or encapsulation order, please include
                      the signed PDF, along with{" "}
                      <Link>this submission order form</Link>, and send your
                      order to the address provided below:
                    </p>

                    <article class="message pb-3" style={{ width: "300px" }}>
                      <div class="message-body">
                        <p class="pb-1">
                          <b>CPS</b>
                        </p>
                        <p class="pb-1">
                          <a href="tel:5199142685">(519) 914-2685</a>
                        </p>
                        <p class="pb-1">
                          <a href="mailto:info@cpscapsule.com">
                            info@cpscapsule.com
                          </a>
                        </p>
                        <p class="pb-1">
                          8-611 Wonderland Road North, P.M.B. 125
                        </p>
                        <p class="pb-1">London, Ontario</p>
                        <p class="pb-1">N6H1T6</p>
                        <p class="pb-1">Canada</p>
                      </div>
                    </article>

                    <p class="pb-3">
                      Once completed, please wait a few weeks for us to receive
                      and process your request.
                    </p>

                    <section class="hero has-background-white-ter">
                      <div class="hero-body">
                        <p class="subtitle">
                          <div class="has-text-centered">
                            <a
                              href={submission.findingsFormObjectUrl}
                              target="_blank"
                              rel="noreferrer"
                              class="button is-large is-success is-fullwidth-mobile"
                            >
                              <FontAwesomeIcon
                                className="fas"
                                icon={faDownload}
                              />
                              &nbsp;Download PDF
                            </a>
                          </div>
                        </p>
                      </div>
                    </section>
                  </article>

                  <div class="columns pt-5">
                    <div class="column is-half"></div>
                    {customerName === null ? (
                      <div class="column is-half has-text-right">
                        <Link
                          to={`/submissions/comics`}
                          class="button is-medium is-primary is-fullwidth-mobile"
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back to List
                        </Link>
                      </div>
                    ) : (
                      <div class="column is-half has-text-right">
                        <Link
                          to={`/customer/${customerID}/comics`}
                          class="button is-medium is-primary is-fullwidth-mobile"
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back to Customer
                        </Link>
                      </div>
                    )}
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

export default RetailerComicSubmissionAddStep4;
