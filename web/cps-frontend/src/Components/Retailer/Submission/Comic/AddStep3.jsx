import React, { useState, useEffect } from "react";
import { Link, Navigate, useSearchParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faShoppingCart,
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

import { postCreateStripeCheckoutSessionURLForComicSubmissionAPI } from "../../../../API/PaymentProcessor";
import { getComicSubmissionDetailAPI } from "../../../../API/ComicSubmission";
import { getOfferDetailByServiceTypeAPI } from "../../../../API/Offer";
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
import { SERVICE_TYPE_PRE_SCREENING_SERVICE } from "../../../../Constants/App";

function RetailerComicSubmissionAddStep3() {
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
  const [offer, setOffer] = useState({});

  ////
  //// Event handling.
  ////

  const onClick = (comicSubmissionID) => {
    // action={`${process.env.REACT_APP_API_PROTOCOL}://${process.env.REACT_APP_API_DOMAIN}/api/v1/stripe/create-subscription-checkout-session`}
    setFetching(true);
    postCreateStripeCheckoutSessionURLForComicSubmissionAPI(
      comicSubmissionID,
      (response) => {
        console.log("onSuccess: Starting...");
        console.log("onSuccess: response:", response);
        console.log("onSuccess: Redirecting at", response.checkoutSessionURL);

        // Force the user's browser to a different domain address.
        window.location.href = response.checkoutSessionURL;
      },
      (apiErr) => {
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
      },
      () => {
        setFetching(false);
      },
      onUnauthorized,
    );
  };

  ////
  //// API.
  ////

  function onComicSubmissionDetailSuccess(response) {
    console.log("onComicSubmissionDetailSuccess: Starting...");

    // ------ CASE 1 ------ //

    // DEVELOPERS NOTE:
    // We will check the response to see if there is a `creditId` field and
    // if there is that means the user burned a credit so we don't need to
    // show this checkout screen but instead redirect the user to the
    // success page.
    if (
      response.creditId !== undefined &&
      response.creditId !== null &&
      response.creditId !== "" &&
      response.creditId !== "000000000000000000000000"
    ) {
      console.log(
        "onComicSubmissionDetailSuccess: user already paid, redirecting user to success page now",
      );
      setForceURL("/submissions/comics/add/" + response.id + "/confirmation");
      return;
    }

    // ------ CASE 2 ------ //

    // DEVELOPERS NOTE:
    // We will check to see if the user has already purchased this comic
    // book submission and then redirect the user to the correct location.
    if (
      response.paymentProcessorPurchaseId !== undefined &&
      response.paymentProcessorPurchaseId !== null &&
      response.paymentProcessorPurchaseId !== "" &&
      response.paymentProcessorPurchaseId.length > 0
    ) {
      console.log(
        "onComicSubmissionDetailSuccess: user already paid, redirecting user to success page now",
      );
      setForceURL("/submissions/comics/add/" + response.id + "/confirmation");
      return;
    }

    // ------ CASE 3 ------ //

    if (response.serviceType === SERVICE_TYPE_PRE_SCREENING_SERVICE) {
      console.log("onComicSubmissionDetailSuccess: user selected cbff");
      setForceURL("/submissions/comics/add/" + response.id + "/confirmation");
      return;
    }

    // ------ CASE 4 ------ //

    console.log(
      "onComicSubmissionDetailSuccess: user must purchase",
      response.serviceType,
    );

    setComicSubmission(response);

    getOfferDetailByServiceTypeAPI(
      response.serviceType,
      (response) => {
        console.log("onSuccess: Starting...");
        console.log("onSuccess: response:", response);
        console.log("onSuccess: Redirecting at", response.checkoutSessionURL);
        setOffer(response);
      },
      (apiErr) => {
        setErrors(apiErr);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
      },
      () => {
        setFetching(false);
      },
      onUnauthorized,
    );
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
              &nbsp;New Online Comic Submission (Checkout)
            </p>
            <FormErrorBox errors={errors} />

            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                <div class="container">
                  <article>
                    <p class="pb-3">
                      You are now ready to checkout and purchase this comic
                      submission. Please review before proceeding:
                    </p>

                    {offer !== undefined && offer !== null && offer !== "" && (
                      <article class="message pb-3" style={{ width: "300px" }}>
                        <div class="message-body">
                          <p class="pb-1">
                            <b>{offer.name}</b>
                          </p>
                          <p class="pb-1">Comic: {submission.seriesTitle}</p>
                          <p class="pb-1">
                            Price: ${offer.price}&nbsp;{offer.priceCurrency}
                          </p>
                        </div>
                      </article>
                    )}
                  </article>

                  <div class="columns pt-5">
                    <div class="column is-half">
                      {customerName === null ? (
                        <div class="">
                          <Link
                            to={`/submissions/comics`}
                            class="button is-medium is-secondary is-fullwidth-mobile"
                          >
                            <FontAwesomeIcon
                              className="fas"
                              icon={faArrowLeft}
                            />
                            &nbsp;Back to Comic Submissions
                          </Link>
                        </div>
                      ) : (
                        <div class="">
                          <Link
                            to={`/customer/${customerID}/comics`}
                            class="button is-medium is-secondary is-fullwidth-mobile"
                          >
                            <FontAwesomeIcon
                              className="fas"
                              icon={faArrowLeft}
                            />
                            &nbsp;Back to Customer
                          </Link>
                        </div>
                      )}
                    </div>
                    <div class="column is-half has-text-right">
                      <button
                        onClick={(e) => {
                          onClick(id);
                        }}
                        class="button is-medium is-primary is-fullwidth-mobile"
                        type="button"
                      >
                        <FontAwesomeIcon
                          className="fas"
                          icon={faShoppingCart}
                        />
                        &nbsp;Checkout
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

export default RetailerComicSubmissionAddStep3;
