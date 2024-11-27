import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faTachometer,
  faPlus,
  faArrowLeft,
  faCheckCircle,
  faUserCircle,
  faGauge,
  faPencil,
  faIdCard,
  faAddressBook,
  faContactCard,
  faChartPie,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { getProfileDetailAPI, putProfileUpdateAPI } from "../../../API/Profile";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormInputField from "../../Reusable/FormInputField";
import FormTextareaField from "../../Reusable/FormTextareaField";
import FormRadioField from "../../Reusable/FormRadioField";
import FormMultiSelectField from "../../Reusable/FormMultiSelectField";
import FormSelectField from "../../Reusable/FormSelectField";
import FormCheckboxField from "../../Reusable/FormCheckboxField";
import FormCountryField from "../../Reusable/FormCountryField";
import FormRegionField from "../../Reusable/FormRegionField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
} from "../../../AppState";
import { USER_ROLE_ROOT, USER_ROLE_RETAILER, USER_ROLE_CUSTOMER } from "../../../Constants/App";


function AccountUpdate() {
  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);
  const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [password, setPassword] = useState("");
  const [passwordRepeated, setPasswordRepeated] = useState("");
  const [companyName, setCompanyName] = useState("");
  const [postalCode, setPostalCode] = useState("");
  const [addressLine1, setAddressLine1] = useState("");
  const [addressLine2, setAddressLine2] = useState("");
  const [city, setCity] = useState("");
  const [region, setRegion] = useState("");
  const [country, setCountry] = useState("");
  const [agreePromotionsEmail, setHasPromotionalEmail] = useState(true);
  const [hasShippingAddress, setHasShippingAddress] = useState(false);
  const [shippingName, setShippingName] = useState("");
  const [shippingPhone, setShippingPhone] = useState("");
  const [shippingCountry, setShippingCountry] = useState("");
  const [shippingRegion, setShippingRegion] = useState("");
  const [shippingCity, setShippingCity] = useState("");
  const [shippingAddressLine1, setShippingAddressLine1] = useState("");
  const [shippingAddressLine2, setShippingAddressLine2] = useState("");
  const [shippingPostalCode, setShippingPostalCode] = useState("");

  ////
  //// Event handling.
  ////

  function onAgreePromotionsEmailChange(e) {
    setHasPromotionalEmail(!agreePromotionsEmail);
  }

  ////
  //// API.
  ////

  const onSubmitClick = (e) => {
    console.log("onSubmitClick: Beginning...");
    setFetching(true);
    setErrors({});

    const submission = {
      Email: email,
      Phone: phone,
      FirstName: firstName,
      LastName: lastName,
      Password: password,
      PasswordRepeated: passwordRepeated,
      CompanyName: companyName,
      PostalCode: postalCode,
      AddressLine1: addressLine1,
      AddressLine2: addressLine2,
      City: city,
      Region: region,
      Country: country,
      AgreePromotionsEmail: agreePromotionsEmail,
      HasShippingAddress: hasShippingAddress,
      ShippingName: shippingName,
      ShippingPhone: shippingPhone,
      ShippingCountry: shippingCountry,
      ShippingRegion: shippingRegion,
      ShippingCity: shippingCity,
      ShippingAddressLine1: shippingAddressLine1,
      ShippingAddressLine2: shippingAddressLine2,
      ShippingPostalCode: shippingPostalCode,
    };
    console.log("onSubmitClick, submission:", submission);
    putProfileUpdateAPI(
      submission,
      onProfileUpdateSuccess,
      onProfileUpdateError,
      onProfileUpdateDone,
      onUnauthorized,
    );
  };

  function onProfileDetailSuccess(response) {
    console.log("onProfileDetailSuccess: Starting...");
    setEmail(response.email);
    setPhone(response.phone);
    setFirstName(response.firstName);
    setLastName(response.lastName);
    setCompanyName(response.companyName);
    setPostalCode(response.postalCode);
    setAddressLine1(response.addressLine1);
    setAddressLine2(response.addressLine2);
    setCity(response.city);
    setRegion(response.region);
    setCountry(response.country);
    setHasPromotionalEmail(response.agreePromotionsEmail);
    setHasShippingAddress(response.hasShippingAddress);
    setShippingName(response.shippingName);
    setShippingPhone(response.shippingPhone);
    setShippingCountry(response.shippingCountry);
    setShippingRegion(response.shippingRegion);
    setShippingCity(response.shippingCity);
    setShippingAddressLine1(response.shippingAddressLine1);
    setShippingAddressLine2(response.shippingAddressLine2);
    setShippingPostalCode(response.shippingPostalCode);
  }

  function onProfileDetailError(apiErr) {
    console.log("onProfileDetailError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onProfileDetailDone() {
    console.log("onProfileDetailDone: Starting...");
    setFetching(false);
  }

  function onProfileUpdateSuccess(response) {
    // For debugging purposes only.
    console.log("onProfileUpdateSuccess: Starting...");
    console.log(response);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Profile updated");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onProfileUpdateSuccess: Delayed for 2 seconds.");
      console.log(
        "onProfileUpdateSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Redirect the user to a new page.
    setForceURL("/account");
  }

  function onProfileUpdateError(apiErr) {
    console.log("onProfileUpdateError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onProfileUpdateError: Delayed for 2 seconds.");
      console.log(
        "onProfileUpdateError: topAlertMessage, topAlertStatus:",
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

  function onProfileUpdateDone() {
    console.log("onProfileUpdateDone: Starting...");
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
      getProfileDetailAPI(
        onProfileDetailSuccess,
        onProfileDetailError,
        onProfileDetailDone,
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

  // Generate URL's based on user role.
  let dashboardURL = "/501";
  if (currentUser) {
      if (currentUser.role === USER_ROLE_ROOT) {
        dashboardURL = "/admin/dashboard";
      }
      if (currentUser.role === USER_ROLE_RETAILER) {
        dashboardURL = "/dashboard";
      }
      if (currentUser.role === USER_ROLE_RETAILER) {
        dashboardURL = "/dashboard";
      }
      if (currentUser.role === USER_ROLE_CUSTOMER) {
        dashboardURL = "/c/dashboard";
      }
  }

  return (
    <>
      <div class="container">
        <section class="section">
          {/* Desktop Breadcrumbs */}
          <nav class="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={dashboardURL} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Dashboard
                </Link>
              </li>
              <li class="">
                <Link to="/account" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faUserCircle} />
                  &nbsp;Account
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faPencil} />
                  &nbsp;Edit
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/account`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Account
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faUserCircle} />
              &nbsp;Account
            </p>
            <FormErrorBox errors={errors} />

            {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Submitting..."} />
            ) : (
              <>
                <div class="container">
                  <p class="subtitle is-6">
                    <FontAwesomeIcon className="fas" icon={faIdCard} />
                    &nbsp;Full Name
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

                  <p class="subtitle is-6">
                    <FontAwesomeIcon className="fas" icon={faContactCard} />
                    &nbsp;Contact Information
                  </p>
                  <hr />

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
                    maxWidth="150px"
                  />

                  <FormCheckboxField
                    label="Has shipping address different then billing address"
                    name="hasShippingAddress"
                    checked={hasShippingAddress}
                    errorText={errors && errors.hasShippingAddress}
                    onChange={(e) => setHasShippingAddress(!hasShippingAddress)}
                    maxWidth="180px"
                  />

                  <div class="columns">
                    <div class="column">
                      <p class="subtitle is-6">
                        {hasShippingAddress ? (
                          <p class="subtitle is-6">Billing Address</p>
                        ) : (
                          <p class="subtitle is-6">Address</p>
                        )}
                      </p>
                      <FormCountryField
                        priorityOptions={["CA", "US", "MX"]}
                        label="Country"
                        name="country"
                        placeholder="Text input"
                        selectedCountry={country}
                        errorText={errors && errors.country}
                        helpText=""
                        onChange={(value) => setCountry(value)}
                        isRequired={true}
                        maxWidth="160px"
                      />

                      <FormRegionField
                        label="Province/Territory"
                        name="region"
                        placeholder="Text input"
                        selectedCountry={country}
                        selectedRegion={region}
                        errorText={errors && errors.region}
                        helpText=""
                        onChange={(value) => setRegion(value)}
                        isRequired={true}
                        maxWidth="280px"
                      />

                      <FormInputField
                        label="City"
                        name="city"
                        placeholder="Text input"
                        value={city}
                        errorText={errors && errors.city}
                        helpText=""
                        onChange={(e) => setCity(e.target.value)}
                        isRequired={true}
                        maxWidth="380px"
                      />

                      <FormInputField
                        label="Address Line 1"
                        name="addressLine1"
                        placeholder="Text input"
                        value={addressLine1}
                        errorText={errors && errors.addressLine1}
                        helpText=""
                        onChange={(e) => setAddressLine1(e.target.value)}
                        isRequired={true}
                        maxWidth="380px"
                      />

                      <FormInputField
                        label="Address Line 2 (Optional)"
                        name="addressLine2"
                        placeholder="Text input"
                        value={addressLine2}
                        errorText={errors && errors.addressLine2}
                        helpText=""
                        onChange={(e) => setAddressLine2(e.target.value)}
                        isRequired={true}
                        maxWidth="380px"
                      />

                      <FormInputField
                        label="Postal Code"
                        name="postalCode"
                        placeholder="Text input"
                        value={postalCode}
                        errorText={errors && errors.postalCode}
                        helpText=""
                        onChange={(e) => setPostalCode(e.target.value)}
                        isRequired={true}
                        maxWidth="80px"
                      />
                    </div>
                    {hasShippingAddress && (
                      <div class="column">
                        <p class="subtitle is-6">Shipping Address</p>

                        <FormInputField
                          label="Name"
                          name="shippingName"
                          placeholder="Text input"
                          value={shippingName}
                          errorText={errors && errors.shippingName}
                          helpText="The name to contact for this shipping address"
                          onChange={(e) => setShippingName(e.target.value)}
                          isRequired={true}
                          maxWidth="350px"
                        />

                        <FormInputField
                          label="Phone"
                          name="shippingPhone"
                          placeholder="Text input"
                          value={shippingPhone}
                          errorText={errors && errors.shippingPhone}
                          helpText="The contact phone number for this shipping address"
                          onChange={(e) => setShippingPhone(e.target.value)}
                          isRequired={true}
                          maxWidth="150px"
                        />

                        <FormCountryField
                          priorityOptions={["CA", "US", "MX"]}
                          label="Country"
                          name="shippingCountry"
                          placeholder="Text input"
                          selectedCountry={shippingCountry}
                          errorText={errors && errors.shippingCountry}
                          helpText=""
                          onChange={(value) => setShippingCountry(value)}
                          isRequired={true}
                          maxWidth="160px"
                        />

                        <FormRegionField
                          label="Province/Territory"
                          name="shippingRegion"
                          placeholder="Text input"
                          selectedCountry={shippingCountry}
                          selectedRegion={shippingRegion}
                          errorText={errors && errors.shippingRegion}
                          helpText=""
                          onChange={(value) => setShippingRegion(value)}
                          isRequired={true}
                          maxWidth="280px"
                        />

                        <FormInputField
                          label="City"
                          name="shippingCity"
                          placeholder="Text input"
                          value={shippingCity}
                          errorText={errors && errors.shippingCity}
                          helpText=""
                          onChange={(e) => setShippingCity(e.target.value)}
                          isRequired={true}
                          maxWidth="380px"
                        />

                        <FormInputField
                          label="Address Line 1"
                          name="shippingAddressLine1"
                          placeholder="Text input"
                          value={shippingAddressLine1}
                          errorText={errors && errors.shippingAddressLine1}
                          helpText=""
                          onChange={(e) =>
                            setShippingAddressLine1(e.target.value)
                          }
                          isRequired={true}
                          maxWidth="380px"
                        />

                        <FormInputField
                          label="Address Line 2 (Optional)"
                          name="shippingAddressLine2"
                          placeholder="Text input"
                          value={shippingAddressLine2}
                          errorText={errors && errors.shippingAddressLine2}
                          helpText=""
                          onChange={(e) =>
                            setShippingAddressLine2(e.target.value)
                          }
                          isRequired={true}
                          maxWidth="380px"
                        />

                        <FormInputField
                          label="Postal Code"
                          name="shippingPostalCode"
                          placeholder="Text input"
                          value={shippingPostalCode}
                          errorText={errors && errors.shippingPostalCode}
                          helpText=""
                          onChange={(e) =>
                            setShippingPostalCode(e.target.value)
                          }
                          isRequired={true}
                          maxWidth="80px"
                        />
                      </div>
                    )}
                  </div>

                  <p class="subtitle is-6">
                    <FontAwesomeIcon className="fas" icon={faChartPie} />
                    &nbsp;Metrics
                  </p>
                  <hr />

                  <FormCheckboxField
                    label="I agree to receive electronic updates from my local retailer and COMICCOIN_FAUCET"
                    name="agreePromotionsEmail"
                    checked={agreePromotionsEmail}
                    errorText={errors && errors.agreePromotionsEmail}
                    onChange={onAgreePromotionsEmailChange}
                    maxWidth="180px"
                  />

                  <div class="columns pt-5">
                    <div class="column is-half">
                      <Link class="button is-fullwidth-mobile" to={"/account"}>
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back to Account
                      </Link>
                    </div>
                    <div class="column is-half has-text-right">
                      <button
                        class="button is-primary is-fullwidth-mobile"
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

export default AccountUpdate;
