import React, { useState, useEffect } from "react";
import { Link, Navigate, useParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
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

import { postCreditCreateAPI } from "../../../../../API/Credit";
import FormErrorBox from "../../../../Reusable/FormErrorBox";
import FormInputField from "../../../../Reusable/FormInputField";
import FormTextareaField from "../../../../Reusable/FormTextareaField";
import FormRadioField from "../../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../../Reusable/FormSelectField";
import FormCheckboxField from "../../../../Reusable/FormCheckboxField";
import PageLoadingContent from "../../../../Reusable/PageLoadingContent";
import FormSelectFieldForOffer from "../../../../Reusable/FormSelectFieldForOffer";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../../AppState";
import {
  CREDIT_BUSINESS_FUNCTION_WITH_EMPTY_OPTIONS,
  NUMBER_OF_CREDITS_WITH_EMPTY_OPTIONS,
} from "../../../../../Constants/FieldOptions";

function AdminUserCreditAdd() {
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
  const [offerID, setOfferID] = useState("");
  const [businessFunction, setBusinessFunction] = useState(0);
  const [numberOfCredits, setNumberOfCredits] = useState(0);

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
      offer_id: offerID,
      business_function: businessFunction,
      user_id: id,
      number_of_credits: numberOfCredits,
    };
    console.log("onSubmitClick, credit:", credit);
    postCreditCreateAPI(
      credit,
      onAdminUserCreditAddSuccess,
      onAdminUserCreditAddError,
      onAdminUserCreditAddDone,
      onUnauthorized,
    );
  };

  function onAdminUserCreditAddSuccess() {
    // For debugging purposes only.
    console.log("onAdminUserCreditAddSuccess: Starting...");

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Credit created");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onAdminUserCreditAddSuccess: Delayed for 2 seconds.");
      console.log(
        "onAdminUserCreditAddSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Redirect the user to a new page.
    setForceURL("/admin/user/" + id + "/credits");
  }

  function onAdminUserCreditAddError(apiErr) {
    console.log("onAdminUserCreditAddError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onAdminUserCreditAddError: Delayed for 2 seconds.");
      console.log(
        "onAdminUserCreditAddError: topAlertMessage, topAlertStatus:",
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

  function onAdminUserCreditAddDone() {
    console.log("onAdminUserCreditAddDone: Starting...");
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

      setFetching(false);
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
                <Link to="/admin/users" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faUsers} />
                  &nbsp;Users
                </Link>
              </li>
              <li class="">
                <Link to={`/admin/user/${id}`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail (Credits)
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
              <FontAwesomeIcon className="fas" icon={faPlus} />
              &nbsp;New Credit
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
                    label="Number of Credits"
                    name="numberOfCredits"
                    placeholder="Pick"
                    selectedValue={numberOfCredits}
                    errorText={errors && errors.numberOfCredits}
                    helpText="Please select how many credits to grant to the user"
                    onChange={(e) =>
                      setNumberOfCredits(parseInt(e.target.value))
                    }
                    options={NUMBER_OF_CREDITS_WITH_EMPTY_OPTIONS}
                  />

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <Link
                        class="button is-medium is-fullwidth-mobile"
                        to={`/admin/user/${id}/credits`}
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

export default AdminUserCreditAdd;
