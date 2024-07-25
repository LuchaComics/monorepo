import React, { useState, useEffect } from "react";
import { Link, Navigate, useParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faDonate,
  faEye,
  faArrowLeft,
  faTasks,
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
  faBuilding,
  faCogs,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { putCreditUpdateAPI, getCreditDetailAPI } from "../../../../../API/Credit";
import FormErrorBox from "../../../../Reusable/FormErrorBox";
import FormInputField from "../../../../Reusable/FormInputField";
import FormTextareaField from "../../../../Reusable/FormTextareaField";
import FormRadioField from "../../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../../Reusable/FormSelectField";
import FormCheckboxField from "../../../../Reusable/FormCheckboxField";
import DataDisplayRowText from "../../../../Reusable/DataDisplayRowText";
import PageLoadingContent from "../../../../Reusable/PageLoadingContent";
import FormSelectFieldForOffer from "../../../../Reusable/FormSelectFieldForOffer";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../../AppState";
import {
  CREDIT_BUSINESS_FUNCTION_WITH_EMPTY_OPTIONS,
  CREDIT_STATUS_WITH_EMPTY_OPTIONS,
} from "../../../../../Constants/FieldOptions";

function AdminUserCreditUpdate() {
  ////
  //// URL Parameters.
  ////

  const { id, cid } = useParams();

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
  const [offerID, setOfferID] = useState("");
  const [businessFunction, setBusinessFunction] = useState(0);
  const [credit, setCredit] = useState({});
  const [status, setStatus] = useState(0);

  ////
  //// Event handling.
  ////

  ////
  //// API.
  ////

  const onSubmitClick = (e) => {
    console.log("onSubmitClick: Beginning...");
    setFetching(true);
    setErrors({});
    const credit = {
      id: cid,
      offer_id: offerID,
      business_function: businessFunction,
      user_id: id,
      status: status,
    };
    console.log("onSubmitClick, credit:", credit);
    putCreditUpdateAPI(
      credit,
      onAdminUserCreditUpdateSuccess,
      onAdminUserCreditUpdateError,
      onAdminUserCreditUpdateDone,
      onUnauthorized,
    );
  };

  function onAdminUserCreditUpdateSuccess(response) {
    // For debugging purposes only.
    console.log("onAdminUserCreditUpdateSuccess: Starting...");
    console.log(response);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Credit update");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onAdminUserCreditUpdateSuccess: Delayed for 2 seconds.");
      console.log(
        "onAdminUserCreditUpdateSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Redirect the user to a new page.
    setForceURL("/admin/user/" + id + "/credit/" + cid);
  }

  function onAdminUserCreditUpdateError(apiErr) {
    console.log("onAdminUserCreditUpdateError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onAdminUserCreditUpdateError: Delayed for 2 seconds.");
      console.log(
        "onAdminUserCreditUpdateError: topAlertMessage, topAlertStatus:",
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

  function onAdminUserCreditUpdateDone() {
    console.log("onAdminUserCreditUpdateDone: Starting...");
    setFetching(false);
  }

  // --- Detail --- //

  function onCreditDetailSuccess(response) {
    console.log("onCreditDetailSuccess: Starting...");
    setCredit(response);

    setOfferID(response.offerId);
    setBusinessFunction(response.businessFunction);
    setStatus(response.status);
  }

  function onCreditDetailError(apiErr) {
    console.log("onCreditDetailError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onCreditDetailDone() {
    console.log("onCreditDetailDone: Starting...");
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
      getCreditDetailAPI(
        cid,
        onCreditDetailSuccess,
        onCreditDetailError,
        onCreditDetailDone,
        onUnauthorized,
      );
    }

    return () => {
      mounted = false;
    };
  }, [cid]);

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
                <Link to="/admin/users" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faUsers} />
                  &nbsp;Users
                </Link>
              </li>
              <li class="">
                <Link to={`/admin/user/${id}/credits`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail (Credits)
                </Link>
              </li>
              <li class="">
                <Link
                  to={`/admin/user/${id}/credit/${cid}`}
                  aria-current="page"
                >
                  <FontAwesomeIcon className="fas" icon={faDonate} />
                  &nbsp;Credit
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faPencil} />
                  &nbsp;Update
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/admin/credits`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Credits
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faPencil} />
              &nbsp;Update Credit
            </p>

            {/* <p class="pb-4 has-text-grey">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Submitting..."} />
            ) : (
              <>
                <FormErrorBox errors={errors} />
                <div class="container">
                  {/*<p class="subtitle is-6"><FontAwesomeIcon className="fas" icon={faIdCard} />&nbsp;Identification</p>
                                    <hr />*/}

                  <DataDisplayRowText label="Credit ID #" value={credit.id} />

                  <FormSelectField
                    label="Business Function"
                    name="businessFunction"
                    placeholder="Pick"
                    selectedValue={businessFunction}
                    errorText={errors && errors.businessFunction}
                    helpText="Please select what this credit will do in our application"
                    onChange={(e) =>
                      setBusinessFunction(parseInt(e.target.value))
                    }
                    options={CREDIT_BUSINESS_FUNCTION_WITH_EMPTY_OPTIONS}
                  />

                  {businessFunction === 1 && (
                    <FormSelectFieldForOffer
                      label="Offer"
                      offerID={offerID}
                      setOfferID={setOfferID}
                      errorText={errors && errors.offerID}
                      helpText=""
                      maxWidth="310px"
                    />
                  )}

                  <FormSelectField
                    label="Status"
                    name="status"
                    placeholder="Pick"
                    selectedValue={status}
                    errorText={errors && errors.status}
                    helpText=""
                    onChange={(e) => setStatus(parseInt(e.target.value))}
                    options={CREDIT_STATUS_WITH_EMPTY_OPTIONS}
                  />

                  <DataDisplayRowText label="User ID #" value={credit.userId} />

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <Link
                        class="button is-medium is-fullwidth-mobile"
                        to={`/admin/user/${id}/credit/${cid}`}
                      >
                        <FontAwesomeIcon className="fas" icon={faTimesCircle} />
                        &nbsp;Cancel
                      </Link>
                    </div>
                    <div class="column is-half has-text-right">
                      <button
                        class="button is-medium is-primary is-fullwidth-mobile"
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

export default AdminUserCreditUpdate;
