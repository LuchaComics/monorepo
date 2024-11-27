import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faBook,
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
  faEye,
  faPencil,
} from "@fortawesome/free-solid-svg-icons";
import Select from "react-select";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import useLocalStorage from "../../../../Hooks/useLocalStorage";
import FormErrorBox from "../../../Reusable/FormErrorBox";
import FormInputField from "../../../Reusable/FormInputField";
import FormTextareaField from "../../../Reusable/FormTextareaField";
import FormRadioField from "../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../Reusable/FormSelectField";
import FormInputFieldWithButton from "../../../Reusable/FormInputFieldWithButton";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../AppState";

function AdminComicSubmissionUpdatePickCustomerWithSearch() {
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
  const [customers, setCustomers] = useState({});
  const [hasCustomer, setHasCustomer] = useState(1);
  const [searchKeyword, setSearchKeyword] = useState("");
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");

  ////
  //// Event handling.
  ////

  const onSearchButtonClicked = (e) => {
    console.log("searchButtonClick: Starting...");
    let aURL = "/admin/submissions/comic/" + id + "/cust/results";
    if (searchKeyword !== "") {
      aURL += "?search=" + searchKeyword;
    }
    if (firstName !== "") {
      if (aURL.indexOf("?") > -1) {
        aURL += "&first_name=" + firstName;
      } else {
        aURL += "?first_name=" + firstName;
      }
    }
    if (lastName !== "") {
      if (aURL.indexOf("?") > -1) {
        aURL += "&last_name=" + lastName;
      } else {
        aURL += "?last_name=" + lastName;
      }
    }
    if (email !== "") {
      if (aURL.indexOf("?") > -1) {
        aURL += "&email=" + email;
      } else {
        aURL += "?email=" + email;
      }
    }
    if (phone !== "") {
      if (aURL.indexOf("?") > -1) {
        aURL += "&phone=" + phone;
      } else {
        aURL += "?phone=" + phone;
      }
    }

    // Validate before proceeding further by checkign to see if we've either
    // searched or filtered and if we did not then error.
    if (aURL.indexOf("?") <= -1) {
      setErrors({ Validation: "Please input data before submitting search." });

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
                <Link to="/admin/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Admin Dashboard
                </Link>
              </li>
              <li class="">
                <Link to="/admin/submissions" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faTasks} />
                  &nbsp;Online Submissions
                </Link>
              </li>
              <li class="">
                <Link to="/admin/submissions/comics" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faBook} />
                  &nbsp;Comics
                </Link>
              </li>
              <li class="">
                <Link
                  to={`/admin/submissions/comic/${id}/cust`}
                  aria-current="page"
                >
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Details
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faPencil} />
                  &nbsp;Update (Customer)
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link
                  to={`/admin/submissions/comic/${id}/cust`}
                  aria-current="page"
                >
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Details
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faPencil} />
              &nbsp;Update Online Comic Submission
            </p>
            <FormErrorBox errors={errors} />

            <div class="container pb-5">
              <p class="subtitle is-6">
                <FontAwesomeIcon className="fas" icon={faSearch} />
                &nbsp;Search
              </p>
              <hr />

              <FormInputField
                label="Search Keywords"
                name="searchKeyword"
                placeholder="Text input"
                value={searchKeyword}
                errorText={errors && errors.searchKeyword}
                helpText=""
                onChange={(e) => setSearchKeyword(e.target.value)}
                isRequired={true}
                maxWidth="380px"
              />
            </div>

            <div class="container pb-5">
              <p class="subtitle is-6">
                <FontAwesomeIcon className="fas" icon={faFilter} />
                &nbsp;Filter
              </p>
              <hr />

              <FormInputField
                label="First Name"
                name="firstName"
                placeholder="Text input"
                value={firstName}
                errorText={errors && errors.firstName}
                helpText=""
                onChange={(e) => setFirstName(e.target.value)}
                isRequired={true}
                maxWidth="380px"
              />

              <FormInputField
                label="Last Name"
                name="lastName"
                placeholder="Text input"
                value={lastName}
                errorText={errors && errors.lastName}
                helpText=""
                onChange={(e) => setLastName(e.target.value)}
                isRequired={true}
                maxWidth="380px"
              />

              <FormInputField
                label="Email"
                name="email"
                placeholder="Text input"
                value={email}
                errorText={errors && errors.email}
                helpText=""
                onChange={(e) => setEmail(e.target.value)}
                isRequired={true}
                maxWidth="380px"
              />

              <FormInputField
                label="Phone"
                name="phone"
                placeholder="Text input"
                value={phone}
                errorText={errors && errors.phone}
                helpText=""
                onChange={(e) => setPhone(e.target.value)}
                isRequired={true}
                maxWidth="380px"
              />
            </div>

            <div class="columns pt-5">
              <div class="column is-half">
                <Link
                  to={`/admin/submissions/comic/${id}/cust`}
                  class="button is-medium is-hidden-touch"
                >
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back
                </Link>
                <Link
                  to={`/admin/submissions/comic/${id}/cust`}
                  class="button is-medium is-fullwidth is-hidden-desktop"
                >
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back
                </Link>
              </div>
              <div class="column is-half has-text-right">
                <button
                  class="button is-medium is-primary is-hidden-touch"
                  onClick={onSearchButtonClicked}
                >
                  <FontAwesomeIcon className="fas" icon={faSearch} />
                  &nbsp;Search
                </button>
                <button
                  class="button is-medium is-primary is-fullwidth is-hidden-desktop"
                  onClick={onSearchButtonClicked}
                >
                  <FontAwesomeIcon className="fas" icon={faSearch} />
                  &nbsp;Search
                </button>
              </div>
            </div>
          </nav>
        </section>
      </div>
    </>
  );
}

export default AdminComicSubmissionUpdatePickCustomerWithSearch;
