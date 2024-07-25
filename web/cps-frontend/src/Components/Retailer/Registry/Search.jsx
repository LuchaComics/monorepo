import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faTachometer,
  faPlus,
  faDownload,
  faArrowLeft,
  faArrowRight,
  faCheckCircle,
  faCheck,
  faGauge,
  faArrowUpRightFromSquare,
  faSearch,
  faFilter,
  faBarcode,
} from "@fortawesome/free-solid-svg-icons";
import Select from "react-select";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import useLocalStorage from "../../../Hooks/useLocalStorage";
import { getSubmissionDetailAPI } from "../../../API/ComicSubmission";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormInputField from "../../Reusable/FormInputField";
import FormTextareaField from "../../Reusable/FormTextareaField";
import FormRadioField from "../../Reusable/FormRadioField";
import FormMultiSelectField from "../../Reusable/FormMultiSelectField";
import FormSelectField from "../../Reusable/FormSelectField";
import FormInputFieldWithButton from "../../Reusable/FormInputFieldWithButton";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";

function RetailerRegistrySearch() {
  ////
  //// URL Parameters.
  ////

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
  const [customers, setCustomers] = useState({});
  const [hasCustomer, setHasCustomer] = useState(1);
  const [cpsrn, setCpsrn] = useState("");

  ////
  //// Event handling.
  ////

  const onSearchButtonClicked = (e) => {
    console.log("searchButtonClick: Starting...");
    let aURL = "/registry";
    let hasCPSRN = false;
    if (cpsrn !== "") {
      aURL += "/" + cpsrn;
      hasCPSRN = true;
    }

    // Validate before proceeding further by checkign to see if we've either
    // searched or filtered and if we did not then error.
    if (hasCPSRN === false) {
      setErrors({ cpsrn: "Please input data before submitting lookup." });

      // The following code will cause the screen to scroll to the top of
      // the page. Please see ``react-scroll`` for more information:
      // https://github.com/fisshy/react-scroll
      var scroll = Scroll.animateScroll;
      scroll.scrollToTop();
    } else {
      setForceURL(aURL);
    }
  };

  ////
  //// API.
  ////

  function onCustomerListSuccess(response) {
    console.log("onCustomerListSuccess: Starting...");
    if (response.results !== null) {
      setCustomers(response);
    }
  }

  function onCustomerListError(apiErr) {
    console.log("onCustomerListError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onCustomerListDone() {
    console.log("onCustomerListDone: Starting...");
    setFetching(false);
  }

  ////
  //// Misc.
  ////

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
                <Link to="/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Dashboard
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faBarcode} />
                  &nbsp;Registry
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
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faBarcode} />
              &nbsp;Registry
            </p>
            <FormErrorBox errors={errors} />
            <div class="container pb-5">
              <p class="subtitle is-6">
                <FontAwesomeIcon className="fas" icon={faSearch} />
                &nbsp;Lookup Submission
              </p>
              <hr />

              <FormInputField
                label="Lookup CPSRN"
                name="cpsrn"
                placeholder="Text input"
                value={cpsrn}
                errorText={errors && errors.cpsrn}
                helpText="MUST BE EXACTL VALUE AS FOUND ON RECORD"
                onChange={(e) => setCpsrn(e.target.value)}
                isRequired={true}
                maxWidth="380px"
              />
            </div>

            <div class="columns pt-5">
              <div class="column is-half">
                <Link
                  to={`/dashboard`}
                  class="button is-medium is-fullwidth-mobile"
                >
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Dashboard
                </Link>
              </div>
              <div class="column is-half has-text-right">
                <button
                  class="button is-medium is-primary is-fullwidth-mobile"
                  onClick={onSearchButtonClicked}
                >
                  <FontAwesomeIcon className="fas" icon={faSearch} />
                  &nbsp;Lookup
                </button>
              </div>
            </div>
          </nav>
        </section>
      </div>
    </>
  );
}

export default RetailerRegistrySearch;
