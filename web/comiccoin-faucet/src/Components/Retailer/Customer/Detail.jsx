import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTrashCan,
  faStar,
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
  postCustomerStarOperationAPI,
  deleteCustomerAPI,
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
import FormTextOptionRow from "../../Reusable/FormRowTextOption";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";
import FormRowText from "../../Reusable/FormRowText";
import FormTextYesNoRow from "../../Reusable/FormRowTextYesNo";
import {
  HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS,
  HOW_LONG_HAS_YOUR_STORE_BEEN_OPERATING_FOR_WITH_EMPTY_OPTIONS,
} from "../../../Constants/FieldOptions";
import FormTextChoiceRow from "../../Reusable/FormRowTextChoice";
import { USER_ROLE_CUSTOMER } from "../../../Constants/App";

function RetailerCustomerDetail() {
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
  const [customer, setCustomer] = useState({});
  const [selectedCustomerForDeletion, setSelectedCustomerForDeletion] =
    useState("");

  ////
  //// Event handling.
  ////

  const onStarClick = () => {
    setFetching(true);
    setErrors({});
    postCustomerStarOperationAPI(
      id,
      onCustomerDetailSuccess,
      onCustomerDetailError,
      onCustomerDetailDone,
      onUnauthorized,
    );
  };

  const onSelectCustomerForDeletion = (e, customer) => {
    console.log("onSelectCustomerForDeletion", customer);
    setSelectedCustomerForDeletion(customer);
  };

  const onDeselectCustomerForDeletion = (e) => {
    console.log("onDeselectCustomerForDeletion");
    setSelectedCustomerForDeletion("");
  };

  const onDeleteConfirmButtonClick = (e) => {
    console.log("onDeleteConfirmButtonClick"); // For debugging purposes only.

    deleteCustomerAPI(
      selectedCustomerForDeletion.id,
      onCustomerDeleteSuccess,
      onCustomerDeleteError,
      onCustomerDeleteDone,
      onUnauthorized,
    );
    setSelectedCustomerForDeletion("");
  };

  ////
  //// API.
  ////

  // --- Detail --- //

  function onCustomerDetailSuccess(response) {
    console.log("onCustomerDetailSuccess: Starting...");
    setCustomer(response);
  }

  function onCustomerDetailError(apiErr) {
    console.log("onCustomerDetailError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onCustomerDetailDone() {
    console.log("onCustomerDetailDone: Starting...");
    setFetching(false);
  }

  // --- DELETE --- //

  function onCustomerDeleteSuccess(response) {
    console.log("onCustomerDeleteSuccess: Starting..."); // For debugging purposes only.

    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("Customer deleted");
    setTimeout(() => {
      console.log(
        "onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    setForceURL("/customers");
  }

  function onCustomerDeleteError(apiErr) {
    console.log("onCustomerDeleteError: Starting..."); // For debugging purposes only.
    setErrors(apiErr);

    // Update notification.
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed deleting");
    setTimeout(() => {
      console.log(
        "onCustomerDeleteError: topAlertMessage, topAlertStatus:",
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

  function onCustomerDeleteDone() {
    console.log("onCustomerDeleteDone: Starting...");
    setFetching(false);
  }

  // --- ALL --- //

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
        onCustomerDetailSuccess,
        onCustomerDetailError,
        onCustomerDetailDone,
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
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/customers`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Customers
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modals */}
          <div
            class={`modal ${selectedCustomerForDeletion ? "is-active" : ""}`}
          >
            <div class="modal-background"></div>
            <div class="modal-card">
              <header class="modal-card-head">
                <p class="modal-card-title">Are you sure?</p>
                <button
                  class="delete"
                  aria-label="close"
                  onClick={onDeselectCustomerForDeletion}
                ></button>
              </header>
              <section class="modal-card-body">
                You are about to <b>archive</b> this customer; it will no longer
                appear on your dashboard This action can be undone but you'll
                need to contact the system administrator. Are you sure you would
                like to continue?
              </section>
              <footer class="modal-card-foot">
                <button
                  class="button is-success"
                  onClick={onDeleteConfirmButtonClick}
                >
                  Confirm
                </button>
                <button class="button" onClick={onDeselectCustomerForDeletion}>
                  Cancel
                </button>
              </footer>
            </div>
          </div>

          {/* Page */}
          <nav class="box">
            {customer && (
              <div class="columns">
                <div class="column">
                  <p class="title is-4">
                    <FontAwesomeIcon className="fas" icon={faUserCircle} />
                    &nbsp;Customer
                  </p>
                </div>
                <div class="column has-text-right">
                  <Link
                    to={`/submissions/pick-type-for-add?customer_id=${id}&customer_name=${customer.name}`}
                    class="button is-small is-success is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faPlus} />
                    &nbsp;COMICCOIN_FAUCET
                  </Link>
                  &nbsp;&nbsp;
                  <Link
                    to={`/customer/${customer.id}/edit`}
                    class="button is-small is-warning is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faPencil} />
                    &nbsp;Edit
                  </Link>
                  &nbsp;&nbsp;
                  <button
                    onClick={(e, ses) =>
                      onSelectCustomerForDeletion(e, customer)
                    }
                    class="button is-small is-danger is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faTrashCan} />
                    &nbsp;Delete
                  </button>
                  &nbsp;&nbsp;
                  {customer.isStarred ? (
                    <Link
                      class="button is-small is-fullwidth-mobile has-text-warning-dark has-background-warning"
                      type="button"
                      onClick={onStarClick}
                    >
                      <FontAwesomeIcon className="mdi" icon={faStar} />
                      &nbsp;Starred
                    </Link>
                  ) : (
                    <Link
                      class="button is-small is-fullwidth-mobile"
                      type="button"
                      onClick={onStarClick}
                    >
                      <FontAwesomeIcon className="mdi" icon={faStar} />
                      <span class="is-hidden-desktop is-hidden-tablet">
                        &nbsp;Unstarred
                      </span>
                    </Link>
                  )}
                </div>
              </div>
            )}
            <FormErrorBox errors={errors} />

            {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                {customer && (
                  <div class="container">
                    <div class="tabs is-medium is-size-7-mobile">
                      <ul>
                        <li class="is-active">
                          <Link>
                            <b>Detail</b>
                          </Link>
                        </li>
                        <li>
                          <Link to={`/customer/${customer.id}/comics`}>
                            Comics
                          </Link>
                        </li>
                        <li>
                          <Link to={`/customer/${customer.id}/comments`}>
                            Comments
                          </Link>
                        </li>
                        <li>
                          <Link to={`/customer/${id}/attachments`}>
                            Attachments
                          </Link>
                        </li>
                      </ul>
                    </div>

                    <p class="subtitle is-6">
                      <FontAwesomeIcon className="fas" icon={faIdCard} />
                      &nbsp;Full Name
                    </p>
                    <hr />

                    <FormRowText
                      label="First Name"
                      value={customer.firstName}
                      helpText=""
                    />
                    <FormRowText
                      label="Last Name"
                      value={customer.lastName}
                      helpText=""
                    />

                    <p class="subtitle is-6">
                      <FontAwesomeIcon className="fas" icon={faContactCard} />
                      &nbsp;Contact Information
                    </p>
                    <hr />

                    <FormRowText
                      label="Email"
                      value={customer.email}
                      helpText=""
                      type="email"
                    />

                    <FormRowText
                      label="Phone"
                      value={customer.phone}
                      helpText=""
                      type="phone"
                    />

                    <FormTextYesNoRow
                      label="Has shipping address different then billing address"
                      checked={customer.hasShippingAddress}
                    />

                    <div class="columns">
                      <div class="column">
                        <p class="subtitle is-6">
                          {customer.hasShippingAddress ? (
                            <p class="subtitle is-6">Billing Address</p>
                          ) : (
                            <p class="subtitle is-6">Address</p>
                          )}
                        </p>
                        <FormRowText
                          label="Country"
                          value={customer.country}
                          helpText=""
                        />

                        <FormRowText
                          label="Province/Territory"
                          value={customer.region}
                          helpText=""
                        />

                        <FormRowText
                          label="City"
                          value={customer.city}
                          helpText=""
                        />

                        <FormRowText
                          label="Address Line 1"
                          value={customer.addressLine1}
                          helpText=""
                        />

                        <FormRowText
                          label="Address Line 2 (Optional)"
                          value={customer.addressLine2}
                          helpText=""
                        />

                        <FormRowText
                          label="Postal Code"
                          value={customer.postalCode}
                          helpText=""
                        />
                      </div>
                      {customer.hasShippingAddress && (
                        <div class="column">
                          <p class="subtitle is-6">Shipping Address</p>

                          <FormRowText
                            label="Name"
                            value={customer.shippingName}
                            helpText="The name to contact for this shipping address"
                          />

                          <FormRowText
                            label="Phone"
                            value={customer.shippingPhone}
                            helpText="The contact phone number for this shipping address"
                          />

                          <FormRowText
                            label="Country"
                            value={customer.shippingCountry}
                            helpText=""
                          />

                          <FormRowText
                            label="Province/Territory"
                            value={customer.shippingRegion}
                            helpText=""
                          />

                          <FormRowText
                            label="City"
                            value={customer.shippingCity}
                            helpText=""
                          />

                          <FormRowText
                            label="Address Line 1"
                            value={customer.shippingAddressLine1}
                            helpText=""
                          />

                          <FormRowText
                            label="Address Line 2 (Optional)"
                            value={customer.shippingAddressLine2}
                            helpText=""
                          />

                          <FormRowText
                            label="Postal Code"
                            value={customer.shippingPostalCode}
                            helpText=""
                          />
                        </div>
                      )}
                    </div>

                    <p class="subtitle is-6">
                      <FontAwesomeIcon className="fas" icon={faChartPie} />
                      &nbsp;Metrics
                    </p>
                    <hr />

                    <FormTextOptionRow
                      label="How did you hear about us?"
                      selectedValue={customer.howDidYouHearAboutUs}
                      options={HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS}
                    />

                    {customer.role === USER_ROLE_CUSTOMER && <>
                        <FormTextChoiceRow
                          label="Question 1: How long have you been collecting for?"
                          name="howLongCollectingComicBooksForGrading"
                          value={customer.howLongCollectingComicBooksForGrading}
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
                          onChange={null}
                          maxWidth="180px"
                          hasOptPerLine={true}
                          disabled={true}
                          readonly={true}
                        />

                        <FormTextChoiceRow
                          label="Question 2: Have you ever submitted a comic book for grading?"
                          name="hasPreviouslySubmittedComicBookForGrading"
                          value={customer.hasPreviouslySubmittedComicBookForGrading}
                          opt1Value={1}
                          opt1Label="Yes"
                          opt2Value={2}
                          opt2Label="No"
                          errorText={
                            errors && errors.hasPreviouslySubmittedComicBookForGrading
                          }
                          onChange={null}
                          maxWidth="180px"
                        />

                        <FormTextChoiceRow
                          label="Question 3: Do you currently own any graded comic books>?"
                          name="hasOwnedGradedComicBooks"
                          value={customer.hasOwnedGradedComicBooks}
                          opt1Value={1}
                          opt1Label="Yes"
                          opt2Value={2}
                          opt2Label="No"
                          errorText={errors && errors.hasOwnedGradedComicBooks}
                          onChange={null}
                          maxWidth="180px"
                        />

                        <FormTextChoiceRow
                          label="Question 4: Do you have a regular comic book shop that you use?"
                          name="hasRegularComicBookShop"
                          value={customer.hasRegularComicBookShop}
                          opt1Value={1}
                          opt1Label="Yes"
                          opt2Value={2}
                          opt2Label="No"
                          errorText={errors && errors.hasRegularComicBookShop}
                          onChange={null}
                          maxWidth="180px"
                        />

                        <FormTextChoiceRow
                          label="Question 5: Have you ever purchase a comic book from an auction site such as eBay?"
                          name="hasPreviouslyPurchasedFromAuctionSite"
                          value={customer.hasPreviouslyPurchasedFromAuctionSite}
                          opt1Value={1}
                          opt1Label="Yes"
                          opt2Value={2}
                          opt2Label="No"
                          errorText={
                            errors && errors.hasPreviouslyPurchasedFromAuctionSite
                          }
                          onChange={null}
                          maxWidth="180px"
                        />

                        <FormTextChoiceRow
                          label="Question 6: Have you ever purchase a comic book from facebook marketplace?"
                          name="hasPreviouslyPurchasedFromFacebookMarketplace"
                          value={customer.hasPreviouslyPurchasedFromFacebookMarketplace}
                          opt1Value={1}
                          opt1Label="Yes"
                          opt2Value={2}
                          opt2Label="No"
                          errorText={
                            errors &&
                            errors.hasPreviouslyPurchasedFromFacebookMarketplace
                          }
                          onChange={null}
                          maxWidth="180px"
                        />

                        <FormTextChoiceRow
                          label="Question 7: Do you regularly attend comic cons or collectible shows?"
                          name="hasRegularlyAttendedComicConsOrCollectibleShows"
                          value={customer.hasRegularlyAttendedComicConsOrCollectibleShows}
                          opt1Value={1}
                          opt1Label="Yes"
                          opt2Value={2}
                          opt2Label="No"
                          errorText={
                            errors &&
                            errors.hasRegularlyAttendedComicConsOrCollectibleShows
                          }
                          onChange={null}
                          maxWidth="180px"
                        />
                    </>}

                    <FormTextYesNoRow
                      label="I agree to receive electronic updates from my local retailer and COMICCOIN_FAUCET"
                      checked={customer.agreePromotionsEmail}
                    />

                    <div class="columns pt-5">
                      <div class="column is-half">
                        <Link
                          class="button is-fullwidth-mobile"
                          to={`/customers`}
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back to Customers
                        </Link>
                      </div>
                      <div class="column is-half has-text-right">
                        <Link
                          to={`/customer/${id}/edit`}
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

export default RetailerCustomerDetail;
