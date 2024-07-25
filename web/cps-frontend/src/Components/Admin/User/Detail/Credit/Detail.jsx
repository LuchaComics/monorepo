import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faDonate,
  faTasks,
  faTachometer,
  faPlus,
  faArrowLeft,
  faCheckCircle,
  faUserCircle,
  faGauge,
  faPencil,
  faUsers,
  faEye,
  faIdCard,
  faAddressBook,
  faContactCard,
  faChartPie,
  faBuilding,
  faCogs,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import { getCreditDetailAPI } from "../../../../../API/Credit";
import FormErrorBox from "../../../../Reusable/FormErrorBox";
import PageLoadingContent from "../../../../Reusable/PageLoadingContent";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../../AppState";
import DataDisplayRowSelect from "../../../../Reusable/DataDisplayRowSelect";
import DataDisplayRowOffer from "../../../../Reusable/DataDisplayRowOffer";
import DataDisplayRowText from "../../../../Reusable/DataDisplayRowText";
import {
  CREDIT_BUSINESS_FUNCTION_WITH_EMPTY_OPTIONS,
  CREDIT_STATUS_WITH_EMPTY_OPTIONS,
} from "../../../../../Constants/FieldOptions";

function AdminUserCreditDetail() {
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
  const [credit, setCredit] = useState({});

  ////
  //// Event handling.
  ////

  //

  ////
  //// API.
  ////

  function onCreditDetailSuccess(response) {
    console.log("onCreditDetailSuccess: Starting...");
    setCredit(response);
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
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faDonate} />
                  &nbsp;Credit
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
            {credit && (
              <div class="columns">
                <div class="column">
                  <p class="title is-4">
                    <FontAwesomeIcon className="fas" icon={faDonate} />
                    &nbsp;Credit
                  </p>
                </div>
                {/* HIDDEN */}
                <div class="is-hidden column has-text-right">
                  {/* Mobile Specific */}
                  <Link
                    to={`/admin/submissions/comics/add?credit_id=${id}&credit_name=${credit.name}`}
                    class="button is-small is-success is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faPlus} />
                    &nbsp;CPS
                  </Link>
                </div>
              </div>
            )}
            <FormErrorBox errors={errors} />

            {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                {credit && (
                  <div class="container">
                    {/*<p class="subtitle is-6 pt-4"><FontAwesomeIcon className="fas" icon={faIdCard} />&nbsp;Identification</p>
                                    <hr />*/}

                    <DataDisplayRowText label="Credit ID #" value={credit.id} />

                    <DataDisplayRowSelect
                      label="Business Function"
                      selectedValue={credit.businessFunction}
                      options={CREDIT_BUSINESS_FUNCTION_WITH_EMPTY_OPTIONS}
                    />

                    <DataDisplayRowOffer
                      label="Offer"
                      offerID={credit.offerId}
                      helpText={`ID #${credit.offerId}`}
                    />

                    <DataDisplayRowSelect
                      label="Status"
                      selectedValue={credit.status}
                      options={CREDIT_STATUS_WITH_EMPTY_OPTIONS}
                    />

                    <DataDisplayRowText
                      label="User ID #"
                      value={credit.userId}
                    />

                    <div class="columns pt-5">
                      <div class="column is-half">
                        <Link
                          class="button is-fullwidth-mobile"
                          to={`/admin/user/${id}/credits`}
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back
                        </Link>
                      </div>
                      <div class="column is-half has-text-right">
                        <Link
                          to={`/admin/user/${id}/credit/${cid}/edit`}
                          class="button is-primary is-fullwidth-mobile"
                        >
                          <FontAwesomeIcon className="fas" icon={faPencil} />
                          &nbsp;Edit
                        </Link>
                      </div>
                    </div>
                  </div>
                )}
              </>
            )}
          </nav>
        </section>
      </div>
    </>
  );
}

export default AdminUserCreditDetail;
