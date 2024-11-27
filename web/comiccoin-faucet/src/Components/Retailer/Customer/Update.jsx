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
  faUsers,
  faEye,
  faIdCard,
  faAddressBook,
  faContactCard,
  faChartPie,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import useLocalStorage from "../../../Hooks/useLocalStorage";
import {
  getCustomerDetailAPI,
  putCustomerUpdateAPI,
} from "../../../API/customer";
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
  HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS,
  HOW_LONG_HAS_YOUR_STORE_BEEN_OPERATING_FOR_WITH_EMPTY_OPTIONS,
} from "../../../Constants/FieldOptions";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";

function RetailerCustomerUpdate() {
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
  const [email, setEmail] = useState("");
  const [phone, setPhone] = useState("");
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [companyName, setCompanyName] = useState("");
  const [postalCode, setPostalCode] = useState("");
  const [addressLine1, setAddressLine1] = useState("");
  const [addressLine2, setAddressLine2] = useState("");
  const [city, setCity] = useState("");
  const [region, setRegion] = useState("");
  const [country, setCountry] = useState("");
  const [agreePromotionsEmail, setHasPromotionalEmail] = useState(true);
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

  ////
  //// API.
  ////

  const onSubmitClick = (e) => {
    console.log("onSubmitClick: Beginning...");
    setFetching(true);
    setErrors({});
    const customer = {
      id: id,
      Email: email,
      Phone: phone,
      FirstName: firstName,
      LastName: lastName,
      CompanyName: companyName,
      PostalCode: postalCode,
      AddressLine1: addressLine1,
      AddressLine2: addressLine2,
      City: city,
      Region: region,
      Country: country,
      AgreePromotionsEmail: agreePromotionsEmail,
      HowDidYouHearAboutUs: howDidYouHearAboutUs,
      HowDidYouHearAboutUsOther: howDidYouHearAboutUsOther,
      HasShippingAddress: hasShippingAddress,
      ShippingName: shippingName,
      ShippingPhone: shippingPhone,
      ShippingCountry: shippingCountry,
      ShippingRegion: shippingRegion,
      ShippingCity: shippingCity,
      ShippingAddressLine1: shippingAddressLine1,
      ShippingAddressLine2: shippingAddressLine2,
      ShippingPostalCode: shippingPostalCode,
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
        hasRegularlyAttendedComicConsOrCollectibleShows
    };
    console.log("onSubmitClick, customer:", customer);
    putCustomerUpdateAPI(
      customer,
      onRetailerCustomerUpdateSuccess,
      onRetailerCustomerUpdateError,
      onRetailerCustomerUpdateDone,
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
    setHowDidYouHearAboutUs(response.howDidYouHearAboutUs);
    setHowDidYouHearAboutUsOther(response.howDidYouHearAboutUsOther);
    setHasShippingAddress(response.hasShippingAddress);
    setShippingName(response.shippingName);
    setShippingPhone(response.shippingPhone);
    setShippingCountry(response.shippingCountry);
    setShippingRegion(response.shippingRegion);
    setShippingCity(response.shippingCity);
    setShippingAddressLine1(response.shippingAddressLine1);
    setShippingAddressLine2(response.shippingAddressLine2);
    setShippingPostalCode(response.shippingPostalCode);
    setHowLongCollectingComicBooksForGrading(response.howLongCollectingComicBooksForGrading);
    setHasPreviouslySubmittedComicBookForGrading(response.hasPreviouslySubmittedComicBookForGrading);
    setHasOwnedGradedComicBooks(response.hasOwnedGradedComicBooks);
    setHasRegularComicBookShop(response.hasRegularComicBookShop);
    setHasPreviouslyPurchasedFromAuctionSite(response.hasPreviouslyPurchasedFromAuctionSite);
    setHasPreviouslyPurchasedFromFacebookMarketplace(response.hasPreviouslyPurchasedFromFacebookMarketplace);
    setHasRegularlyAttendedComicConsOrCollectibleShows(response.hasRegularlyAttendedComicConsOrCollectibleShows);
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

  function onRetailerCustomerUpdateSuccess(response) {
    // For debugging purposes only.
    console.log("onRetailerCustomerUpdateSuccess: Starting...");
    console.log(response);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Customer updated");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onRetailerCustomerUpdateSuccess: Delayed for 2 seconds.");
      console.log(
        "onRetailerCustomerUpdateSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Redirect the user to a new page.
    setForceURL("/customer/" + response.id);
  }

  function onRetailerCustomerUpdateError(apiErr) {
    console.log("onRetailerCustomerUpdateError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onRetailerCustomerUpdateError: Delayed for 2 seconds.");
      console.log(
        "onRetailerCustomerUpdateError: topAlertMessage, topAlertStatus:",
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

  function onRetailerCustomerUpdateDone() {
    console.log("onRetailerCustomerUpdateDone: Starting...");
    setFetching(false);
  }

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
      getCustomerDetailAPI(
        id,
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
              <li class="">
                <Link to="/customers" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faUsers} />
                  &nbsp;Customers
                </Link>
              </li>
              <li class="">
                <Link to={`/customer/${id}`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail
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
                <Link to={`/customer/${id}`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Detail
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page */}
          <nav class="box">
            <p class="title is-4">
              <FontAwesomeIcon className="fas" icon={faUserCircle} />
              &nbsp;Customer
            </p>
            <FormErrorBox errors={errors} />

            {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
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
                        label="Country (Optional)"
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
                        label="Province/Territory (Optional)"
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
                        label="City (Optional)"
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
                        label="Address Line 1 (Optional)"
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
                        label="Postal Code (Optional)"
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
                          label="Name (Optional)"
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
                          label="Phone (Optional)"
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
                          label="Country (Optional)"
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
                          label="Province/Territory (Optional)"
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
                          label="City (Optional)"
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
                          label="Address Line 1 (Optional)"
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
                          label="Postal Code (Optional)"
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
                      <Link
                        class="button is-fullwidth-mobile"
                        to={`/customer/${id}`}
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back
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

export default RetailerCustomerUpdate;
