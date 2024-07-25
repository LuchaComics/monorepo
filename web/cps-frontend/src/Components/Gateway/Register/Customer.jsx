import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faExclamationTriangle,
  faTasks,
  faTachometer,
  faPlus,
  faArrowLeft,
  faCheckCircle,
  faArrowRight,
  faArrowUpRightFromSquare,
} from "@fortawesome/free-solid-svg-icons";
import Select from "react-select";
import { useRecoilState } from "recoil";

import { postRegisterCustomerAPI } from "../../../API/Gateway";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormInputField from "../../Reusable/FormInputField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import FormTextareaField from "../../Reusable/FormTextareaField";
import FormRadioField from "../../Reusable/FormRadioField";
import FormMultiSelectField from "../../Reusable/FormMultiSelectField";
import FormSelectField from "../../Reusable/FormSelectField";
import FormCheckboxField from "../../Reusable/FormCheckboxField";
import FormCountryField from "../../Reusable/FormCountryField";
import FormRegionField from "../../Reusable/FormRegionField";
import FormTimezoneSelectField from "../../Reusable/FormTimezoneField";
import FormSelectFieldForPublicStore from "../../Reusable/FormSelectFieldForPublicStore";
import {
  HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS,
  HOW_LONG_HAS_YOUR_STORE_BEEN_OPERATING_FOR_WITH_EMPTY_OPTIONS,
  ESTIMATED_SUBMISSION_PER_MONTH_WITH_EMPTY_OPTIONS,
} from "../../../Constants/FieldOptions";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";


function RegisterAsCustomer() {
  ////
  ////
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
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [password, setPassword] = useState("");
  const [passwordRepeated, setPasswordRepeated] = useState("");
  const [postalCode, setPostalCode] = useState("");
  const [addressLine1, setAddressLine1] = useState("");
  const [addressLine2, setAddressLine2] = useState("");
  const [city, setCity] = useState("");
  const [region, setRegion] = useState("");
  const [country, setCountry] = useState("");
  const [agreePromotionsEmail, setHasPromotionalEmail] = useState(true);
  const [agreeTOS, setAgreeTOS] = useState();
  const [howDidYouHearAboutUs, setHowDidYouHearAboutUs] = useState(0);
  const [howDidYouHearAboutUsOther, setHowDidYouHearAboutUsOther] =
    useState("");
  const [hasShippingAddress, setHasShippingAddress] = useState(false);
  const [shippingName, setShippingName] = useState("");
  const [shippingPhone, setShippingPhone] = useState("");
  const [shippingCountry, setShippingCountry] = useState("");
  const [shippingRegion, setShippingRegion] = useState("");
  const [shippingCity, setShippingCity] = useState("");
  const [shippingAddressLine1, setShippingAddressLine1] = useState("");
  const [shippingAddressLine2, setShippingAddressLine2] = useState("");
  const [shippingPostalCode, setShippingPostalCode] = useState("");
  const [storeID, setStoreID] = useState("");
  const [storeName, setStoreName] = useState("");
  const [
    howLongCollectingComicBooksForGrading,
    setHowLongCollectingComicBooksForGrading,
  ] = useState(0);
  const [
    hasPreviouslySubmittedComicBookForGrading,
    setHasPreviouslySubmittedComicBookForGrading,
  ] = useState(0);
  const [hasOwnedGradedComicBooks, setHasOwnedGradedComicBooks] = useState(0);
  const [hasRegularComicBookShop, setHasRegularComicBookShop] = useState(0);
  const [
    hasPreviouslyPurchasedFromAuctionSite,
    setHasPreviouslyPurchasedFromAuctionSite,
  ] = useState(0);
  const [
    hasPreviouslyPurchasedFromFacebookMarketplace,
    setHasPreviouslyPurchasedFromFacebookMarketplace,
  ] = useState(0);
  const [
    hasRegularlyAttendedComicConsOrCollectibleShows,
    setHasRegularlyAttendedComicConsOrCollectibleShows,
  ] = useState(0);

  ////
  //// Event handling.
  ////

  function onAgreePromotionsEmailChange(e) {
    setHasPromotionalEmail(!agreePromotionsEmail);
  }

  function onAgreeTOSChange(e) {
    setAgreeTOS(!agreeTOS);
  }

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
      PostalCode: postalCode,
      AddressLine1: addressLine1,
      AddressLine2: addressLine2,
      City: city,
      Region: region,
      Country: country,
      AgreeTOS: agreeTOS,
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
      StoreId: storeID,
      HowDidYouHearAboutUs: howDidYouHearAboutUs,
      HowDidYouHearAboutUsOther: howDidYouHearAboutUsOther,
      HowLongCollectingComicBooksForGrading:
        howLongCollectingComicBooksForGrading,
      HasPreviouslySubmittedComicBookForGrading:
        hasPreviouslySubmittedComicBookForGrading,
      HasOwnedGradedComicBooks: hasOwnedGradedComicBooks,
      HasRegularComicBookShop: hasRegularComicBookShop,
      HasPreviouslyPurchasedFromAuctionSite:
        hasPreviouslyPurchasedFromAuctionSite,
      HasPreviouslyPurchasedFromFacebookMarketplace:
        hasPreviouslyPurchasedFromFacebookMarketplace,
      HasRegularlyAttendedComicConsOrCollectibleShows:
        hasRegularlyAttendedComicConsOrCollectibleShows,
    };
    console.log("onSubmitClick, submission:", submission);
    postRegisterCustomerAPI(
      submission,
      onRegisterSuccess,
      onRegisterError,
      onRegisterDone,
    );
  };

  ////
  //// API.
  ////

  function onRegisterSuccess(response) {
    // For debugging purposes only.
    console.log("onRegisterSuccess: Starting...");
    console.log(response);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Submission created");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onRegisterSuccess: Delayed for 2 seconds.");
      console.log(
        "onRegisterSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Redirect the user to a new page.
    setForceURL("/register-successful");
  }

  function onRegisterError(apiErr) {
    console.log("onRegisterError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onRegisterError: Delayed for 2 seconds.");
      console.log(
        "onRegisterError: topAlertMessage, topAlertStatus:",
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

  function onRegisterDone() {
    console.log("onRegisterDone: Starting...");
    setFetching(false);
  }

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

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <>
      <div class="container">
        <section class="section">
          <nav class="box">
            {/*
              <article class="message is-danger">
                <div class="message-body">
                  <FontAwesomeIcon className="mdi" icon={faExclamationTriangle} />&nbsp;This site is in active development and all data will be lost. Use at your own risk.
                </div>
              </article>
            */}
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faTasks} />
              &nbsp;Register
            </p>
            <p class="is-italic">
              Please note: you are about register a customer of a store
            </p>
            <br />
            <FormErrorBox errors={errors} />

            {isFetching && (
              <PageLoadingContent displayMessage={"Submitting..."} />
            )}

            {!isFetching && (
              <div class="container">
                <p class="subtitle is-3">
                  <u>
                    <b>Collectibles Dealer / Retailer</b>
                  </u>
                </p>

                <FormSelectFieldForPublicStore
                  label="Store"
                  storeID={storeID}
                  setStoreID={setStoreID}
                  storeName={storeName}
                  setStoreName={setStoreName}
                  errorText={errors && errors.storeID}
                  helpText="Please select the store to filter by"
                  maxWidth="310px"
                />

                <p class="subtitle is-3">
                  <u>
                    <b>Your Information</b>
                  </u>
                </p>

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
                  label="Password"
                  name="password"
                  type="password"
                  placeholder="Text input"
                  value={password}
                  errorText={errors && errors.password}
                  helpText=""
                  onChange={(e) => setPassword(e.target.value)}
                  isRequired={true}
                  maxWidth="380px"
                />

                <FormInputField
                  label="Password Repeated"
                  name="passwordRepeated"
                  type="password"
                  placeholder="Text input"
                  value={passwordRepeated}
                  errorText={errors && errors.passwordRepeated}
                  helpText=""
                  onChange={(e) => setPasswordRepeated(e.target.value)}
                  isRequired={true}
                  maxWidth="380px"
                />

                <p class="subtitle is-3">
                  <u>
                    <b>Contact Information</b>
                  </u>
                </p>

                <FormInputField
                  label="Email"
                  name="email"
                  type="email"
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
                        <p class="subtitle is-3">
                          <u>
                            <b>Billing Address</b>
                          </u>
                        </p>
                      ) : (
                        <p class="subtitle is-3">
                          <u>
                            <b>Address</b>
                          </u>
                        </p>
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
                      <p class="subtitle is-3">
                        <u>
                          <b>Shipping Address</b>
                        </u>
                      </p>

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
                        onChange={(e) => setShippingPostalCode(e.target.value)}
                        isRequired={true}
                        maxWidth="80px"
                      />
                    </div>
                  )}
                </div>

                <p class="subtitle is-3">
                  <u>
                    <b>About You</b>
                  </u>
                </p>

                <FormSelectField
                  label="How did you hear about us?"
                  name="howDidYouHearAboutUs"
                  placeholder="Pick"
                  selectedValue={howDidYouHearAboutUs}
                  errorText={errors && errors.howDidYouHearAboutUs}
                  helpText=""
                  onChange={(e) =>
                    setHowDidYouHearAboutUs(parseInt(e.target.value))
                  }
                  options={HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS}
                />

                {howDidYouHearAboutUs === 1 && (
                  <FormInputField
                    label="Other (Please specify):"
                    name="howDidYouHearAboutUsOther"
                    placeholder="Text input"
                    value={howDidYouHearAboutUsOther}
                    errorText={errors && errors.howDidYouHearAboutUsOther}
                    helpText=""
                    onChange={(e) =>
                      setHowDidYouHearAboutUsOther(e.target.value)
                    }
                    isRequired={true}
                    maxWidth="380px"
                  />
                )}

                <FormRadioField
                  label="Question 1: How long have you been collecting for?"
                  name="howLongCollectingComicBooksForGrading"
                  value={howLongCollectingComicBooksForGrading}
                  opt1Value={1}
                  opt1Label="Less than 1 year"
                  opt2Value={2}
                  opt2Label="1-3 years"
                  opt3Value={3}
                  opt3Label="3-5 years"
                  opt4Value={4}
                  opt4Label="5-10 years"
                  opt5Value={5}
                  opt5Label="10+ years"
                  errorText={
                    errors && errors.howLongCollectingComicBooksForGrading
                  }
                  onChange={(e) =>
                    setHowLongCollectingComicBooksForGrading(
                      parseInt(e.target.value),
                    )
                  }
                  maxWidth="180px"
                  hasOptPerLine={true}
                />

                <FormRadioField
                  label="Question 2: Have you ever submitted a comic book for grading?"
                  name="hasPreviouslySubmittedComicBookForGrading"
                  value={hasPreviouslySubmittedComicBookForGrading}
                  opt1Value={1}
                  opt1Label="Yes"
                  opt2Value={2}
                  opt2Label="No"
                  errorText={
                    errors && errors.hasPreviouslySubmittedComicBookForGrading
                  }
                  onChange={(e) =>
                    setHasPreviouslySubmittedComicBookForGrading(
                      parseInt(e.target.value),
                    )
                  }
                  maxWidth="180px"
                />

                <FormRadioField
                  label="Question 3: Do you currently own any graded comic books>?"
                  name="hasOwnedGradedComicBooks"
                  value={hasOwnedGradedComicBooks}
                  opt1Value={1}
                  opt1Label="Yes"
                  opt2Value={2}
                  opt2Label="No"
                  errorText={errors && errors.hasOwnedGradedComicBooks}
                  onChange={(e) =>
                    setHasOwnedGradedComicBooks(parseInt(e.target.value))
                  }
                  maxWidth="180px"
                />

                <FormRadioField
                  label="Question 4: Do you have a regular comic book shop that you use?"
                  name="hasRegularComicBookShop"
                  value={hasRegularComicBookShop}
                  opt1Value={1}
                  opt1Label="Yes"
                  opt2Value={2}
                  opt2Label="No"
                  errorText={errors && errors.hasRegularComicBookShop}
                  onChange={(e) =>
                    setHasRegularComicBookShop(parseInt(e.target.value))
                  }
                  maxWidth="180px"
                />

                <FormRadioField
                  label="Question 5: Have you ever purchase a comic book from an auction site such as eBay?"
                  name="hasPreviouslyPurchasedFromAuctionSite"
                  value={hasPreviouslyPurchasedFromAuctionSite}
                  opt1Value={1}
                  opt1Label="Yes"
                  opt2Value={2}
                  opt2Label="No"
                  errorText={
                    errors && errors.hasPreviouslyPurchasedFromAuctionSite
                  }
                  onChange={(e) =>
                    setHasPreviouslyPurchasedFromAuctionSite(
                      parseInt(e.target.value),
                    )
                  }
                  maxWidth="180px"
                />

                <FormRadioField
                  label="Question 6: Have you ever purchase a comic book from facebook marketplace?"
                  name="hasPreviouslyPurchasedFromFacebookMarketplace"
                  value={hasPreviouslyPurchasedFromFacebookMarketplace}
                  opt1Value={1}
                  opt1Label="Yes"
                  opt2Value={2}
                  opt2Label="No"
                  errorText={
                    errors &&
                    errors.hasPreviouslyPurchasedFromFacebookMarketplace
                  }
                  onChange={(e) =>
                    setHasPreviouslyPurchasedFromFacebookMarketplace(
                      parseInt(e.target.value),
                    )
                  }
                  maxWidth="180px"
                />

                <FormRadioField
                  label="Question 7: Do you regularly attend comic cons or collectible shows?"
                  name="hasRegularlyAttendedComicConsOrCollectibleShows"
                  value={hasRegularlyAttendedComicConsOrCollectibleShows}
                  opt1Value={1}
                  opt1Label="Yes"
                  opt2Value={2}
                  opt2Label="No"
                  errorText={
                    errors &&
                    errors.hasRegularlyAttendedComicConsOrCollectibleShows
                  }
                  onChange={(e) =>
                    setHasRegularlyAttendedComicConsOrCollectibleShows(
                      parseInt(e.target.value),
                    )
                  }
                  maxWidth="180px"
                />

                <p class="subtitle is-3">
                  <u>
                    <b>Legalese</b>
                  </u>
                </p>

                <FormCheckboxField
                  label="I agree to receive electronic updates from my local retailer and CPS"
                  name="agreePromotionsEmail"
                  checked={agreePromotionsEmail}
                  errorText={errors && errors.agreePromotionsEmail}
                  onChange={onAgreePromotionsEmailChange}
                  maxWidth="180px"
                />

                <FormCheckboxField
                  label={
                    <>
                      I agree to{" "}
                      <Link to="/terms" target="_blank" rel="noreferrer">
                        terms of service&nbsp;
                        <FontAwesomeIcon
                          className="fas"
                          icon={faArrowUpRightFromSquare}
                        />
                      </Link>{" "}
                      and{" "}
                      <Link to="/privacy" target="_blank" rel="noreferrer">
                        privacy policy&nbsp;
                        <FontAwesomeIcon
                          className="fas"
                          icon={faArrowUpRightFromSquare}
                        />
                      </Link>
                    </>
                  }
                  name="agreeTOS"
                  checked={agreeTOS}
                  errorText={errors && errors.agreeTos}
                  onChange={onAgreeTOSChange}
                  maxWidth="180px"
                />

                <div class="columns">
                  <div class="column is-half">
                    <Link
                      to={`/register`}
                      class="button is-medium is-fullwidth-mobile"
                    >
                      <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                      &nbsp;Back
                    </Link>
                  </div>
                  <div class="column is-half has-text-right">
                    <button
                      class="button is-medium is-primary is-fullwidth-mobile"
                      onClick={onSubmitClick}
                    >
                      <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                      &nbsp;Register
                    </button>
                  </div>
                </div>
              </div>
            )}
          </nav>
          <span className="is-pulled-right has-text-grey">
            Already have an account?{" "}
            <Link to="/login">
              Click here&nbsp;
              <FontAwesomeIcon className="fas" icon={faArrowRight} />
            </Link>{" "}
            to sign in.
          </span>
        </section>
        <div className="has-text-centered">
          <br />
          <p>Need help?</p>
          <p>
            <Link to="Support@cpscapsule.com">Support@cpscapsule.com</Link>
          </p>
          <p>
            <a href="tel:+15199142685">(519) 914-2685</a>
          </p>
          <p>
            <Link to="/cpsrn-registry" className="">
              CPSRN Registry&nbsp;
              <FontAwesomeIcon className="fas" icon={faArrowRight} />
            </Link>
          </p>
        </div>
      </div>
    </>
  );
}

export default RegisterAsCustomer;
