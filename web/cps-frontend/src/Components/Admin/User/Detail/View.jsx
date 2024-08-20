import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faEllipsis,
  faTrashCan,
  faCog,
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
  faCogs,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import useLocalStorage from "../../../../Hooks/useLocalStorage";
import {
  getUserDetailAPI,
  postUserStarOperationAPI,
  deleteUserAPI,
} from "../../../../API/user";
import { getStoreSelectOptionListAPI } from "../../../../API/store";
import FormErrorBox from "../../../Reusable/FormErrorBox";
import FormInputField from "../../../Reusable/FormInputField";
import FormTextareaField from "../../../Reusable/FormTextareaField";
import FormRadioField from "../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../Reusable/FormSelectField";
import FormCheckboxField from "../../../Reusable/FormCheckboxField";
import FormCountryField from "../../../Reusable/FormCountryField";
import FormRegionField from "../../../Reusable/FormRegionField";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../AppState";
import FormRowText from "../../../Reusable/FormRowText";
import FormTextYesNoRow from "../../../Reusable/FormRowTextYesNo";
import FormTextOptionRow from "../../../Reusable/FormRowTextOption";
import FormTextChoiceRow from "../../../Reusable/FormRowTextChoice";
import {
  HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS,
  HOW_LONG_HAS_YOUR_STORE_BEEN_OPERATING_FOR_WITH_EMPTY_OPTIONS,
  USER_SPECIAL_COLLECTION_WITH_EMPTY_OPTIONS,
} from "../../../../Constants/FieldOptions";
import AlertBanner from "../../../Reusable/EveryPage/AlertBanner";
import {
  USER_ROLE_CUSTOMER
} from "../../../../Constants/App";

function AdminUserDetail() {
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
  const [user, setUser] = useState({});
  const [storeSelectOptions, setStoreSelectOptions] = useState([]);
  const [selectedUserForDeletion, setSelectedUserForDeletion] = useState("");

  ////
  //// Event handling.
  ////

  const onStarClick = () => {
    setFetching(true);
    setErrors({});
    postUserStarOperationAPI(
      id,
      onUserDetailSuccess,
      onUserDetailError,
      onUserDetailDone,
      onUnauthorized,
    );
  };

  const onSelectUserForDeletion = (e, user) => {
    console.log("onSelectUserForDeletion", user);
    setSelectedUserForDeletion(user);
  };

  const onDeselectUserForDeletion = (e) => {
    console.log("onDeselectUserForDeletion");
    setSelectedUserForDeletion("");
  };

  const onDeleteConfirmButtonClick = (e) => {
    console.log("onDeleteConfirmButtonClick"); // For debugging purposes only.

    deleteUserAPI(
      selectedUserForDeletion.id,
      onUserDeleteSuccess,
      onUserDeleteError,
      onUserDeleteDone,
      onUnauthorized,
    );
    setSelectedUserForDeletion("");
  };

  ////
  //// API.
  ////

  // --- DETAIL --- //

  function onUserDetailSuccess(response) {
    console.log("onUserDetailSuccess: Starting...");
    setUser(response);
  }

  function onUserDetailError(apiErr) {
    console.log("onUserDetailError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onUserDetailDone() {
    console.log("onUserDetailDone: Starting...");
    setFetching(false);
  }

  // --- STORE OPTIONS --- //

  function onStoreOptionListSuccess(response) {
    console.log("onStoreOptionListSuccess: Starting...");
    if (response !== null) {
      const selectOptions = [
        { value: 0, label: "Please select" }, // Add empty options.
        ...response,
      ];
      setStoreSelectOptions(selectOptions);
    }
  }

  function onStoreOptionListError(apiErr) {
    console.log("onStoreOptionListError: Starting...");
    console.log("onStoreOptionListError: apiErr:", apiErr);
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onStoreOptionListDone() {
    console.log("onStoreOptionListDone: Starting...");
    setFetching(false);
  }

  // --- DELETE --- //

  function onUserDeleteSuccess(response) {
    console.log("onUserDeleteSuccess: Starting..."); // For debugging purposes only.

    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("User deleted");
    setTimeout(() => {
      console.log(
        "onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Fetch again an updated list.
    setForceURL("/admin/users");
  }

  function onUserDeleteError(apiErr) {
    console.log("onUserDeleteError: Starting..."); // For debugging purposes only.
    setErrors(apiErr);

    // Update notification.
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed deleting");
    setTimeout(() => {
      console.log(
        "onUserDeleteError: topAlertMessage, topAlertStatus:",
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

  function onUserDeleteDone() {
    console.log("onUserDeleteDone: Starting...");
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
      getUserDetailAPI(
        id,
        onUserDetailSuccess,
        onUserDetailError,
        onUserDetailDone,
        onUnauthorized,
      );

      let params = new Map();
      getStoreSelectOptionListAPI(
        params,
        onStoreOptionListSuccess,
        onStoreOptionListError,
        onStoreOptionListDone,
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
                <Link to={`/admin/users`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Users
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modals */}
          <div class={`modal ${selectedUserForDeletion ? "is-active" : ""}`}>
            <div class="modal-background"></div>
            <div class="modal-card">
              <header class="modal-card-head">
                <p class="modal-card-title">Are you sure?</p>
                <button
                  class="delete"
                  aria-label="close"
                  onClick={onDeselectUserForDeletion}
                ></button>
              </header>
              <section class="modal-card-body">
                You are about to <b>archive</b> this user; it will no longer
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
                <button class="button" onClick={onDeselectUserForDeletion}>
                  Cancel
                </button>
              </footer>
            </div>
          </div>

          {/* Page banner */}
          {user && user.status === 100 && (
            <AlertBanner message="Archived" status="info" />
          )}

          {/* Page */}
          <nav class="box">
            {user && (
              <div class="columns">
                <div class="column">
                  <p class="title is-4">
                    <FontAwesomeIcon className="fas" icon={faUserCircle} />
                    &nbsp;User
                  </p>
                </div>
                {user && user.status === 1 && <div class="column has-text-right">
                  <Link
                    to={`/admin/submissions/pick-type-for-add?customer_id=${user.id}&customer_name=${user.name}&store_id=${user.storeId}&from=usercomics&clear=true`}
                    class="button is-small is-success is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faPlus} />
                    &nbsp;CPS
                  </Link>
                  &nbsp;&nbsp;
                  <Link
                    to={`/admin/user/${user.id}/edit`}
                    class="button is-small is-warning is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faPencil} />
                    &nbsp;Edit
                  </Link>
                  &nbsp;&nbsp;
                  <button
                    onClick={(e, ses) => onSelectUserForDeletion(e, user)}
                    class="button is-small is-danger is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faTrashCan} />
                    &nbsp;Delete
                  </button>
                  &nbsp;&nbsp;
                  {user.isStarred ? (
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
                </div>}
              </div>
            )}
            <FormErrorBox errors={errors} />

            {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                {user && (
                  <div class="container">
                    <div class="tabs is-medium is-size-7-mobile">
                      <ul>
                        <li class="is-active">
                          <Link>
                            <b>Detail</b>
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/user/${user.id}/comics`}>
                            Comics
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/user/${user.id}/comments`}>
                            Comments
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/user/${user.id}/attachments`}>
                            Attachments
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/user/${user.id}/credits`}>
                            Credits
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/user/${user.id}/more`}>
                            More&nbsp;&nbsp;
                            <FontAwesomeIcon
                              className="mdi"
                              icon={faEllipsis}
                            />
                          </Link>
                        </li>
                      </ul>
                    </div>

                    <p class="subtitle is-6">
                      <FontAwesomeIcon className="fas" icon={faCogs} />
                      &nbsp;Settings
                    </p>
                    <hr />

                    {storeSelectOptions && storeSelectOptions.length > 0 && (
                      <FormTextOptionRow
                        label="Store"
                        selectedValue={user.storeID}
                        helpText=""
                        options={storeSelectOptions}
                      />
                    )}
                    <FormTextChoiceRow
                      label="Role"
                      value={user.role}
                      opt1Value={1}
                      opt1Label="Admin"
                      opt2Value={2}
                      opt2Label="Store Owner/Manager"
                      opt3Value={3}
                      opt3Label="Customer"
                    />
                    <FormTextChoiceRow
                      label="Status"
                      value={user.status}
                      opt1Value={1}
                      opt1Label="Active"
                      opt2Value={2}
                      opt2Label="Archived"
                    />

                    <p class="subtitle is-6">
                      <FontAwesomeIcon className="fas" icon={faIdCard} />
                      &nbsp;Full Name
                    </p>
                    <hr />

                    <FormRowText
                      label="First Name"
                      value={user.firstName}
                      helpText=""
                    />
                    <FormRowText
                      label="Last Name"
                      value={user.lastName}
                      helpText=""
                    />

                    <p class="subtitle is-6">
                      <FontAwesomeIcon className="fas" icon={faContactCard} />
                      &nbsp;Contact Information
                    </p>
                    <hr />

                    <FormRowText
                      label="Email"
                      type="email"
                      value={user.email}
                      helpText=""
                    />

                    <FormRowText
                      label="Phone"
                      type="phone"
                      value={user.phone}
                      helpText=""
                    />

                    <FormTextYesNoRow
                      label="Has shipping address different then billing address"
                      checked={user.hasShippingAddress}
                    />

                    <div class="columns">
                      <div class="column">
                        <p class="subtitle is-6">
                          {user.hasShippingAddress ? (
                            <p class="subtitle is-6">Billing Address</p>
                          ) : (
                            <p class="subtitle is-6">Address</p>
                          )}
                        </p>
                        <FormRowText
                          label="Country"
                          value={user.country}
                          helpText=""
                        />

                        <FormRowText
                          label="Province/Territory"
                          value={user.region}
                          helpText=""
                        />

                        <FormRowText
                          label="City"
                          value={user.city}
                          helpText=""
                        />

                        <FormRowText
                          label="Address Line 1"
                          value={user.addressLine1}
                          helpText=""
                        />

                        <FormRowText
                          label="Address Line 2 (Optional)"
                          value={user.addressLine2}
                          helpText=""
                        />

                        <FormRowText
                          label="Postal Code"
                          value={user.postalCode}
                          helpText=""
                        />
                      </div>
                      {user.hasShippingAddress && (
                        <div class="column">
                          <p class="subtitle is-6">Shipping Address</p>

                          <FormRowText
                            label="Name"
                            value={user.shippingName}
                            helpText="The name to contact for this shipping address"
                          />

                          <FormRowText
                            label="Phone"
                            type="phone"
                            value={user.shippingPhone}
                            helpText="The contact phone number for this shipping address"
                          />

                          <FormRowText
                            name="shippingCountry"
                            value={user.shippingCountry}
                            helpText=""
                          />

                          <FormRowText
                            label="Province/Territory"
                            value={user.shippingRegion}
                            helpText=""
                          />

                          <FormRowText
                            label="City"
                            value={user.shippingCity}
                            helpText=""
                          />

                          <FormRowText
                            label="Address Line 1"
                            value={user.shippingAddressLine1}
                            helpText=""
                          />

                          <FormRowText
                            label="Address Line 2 (Optional)"
                            value={user.shippingAddressLine2}
                            helpText=""
                          />

                          <FormRowText
                            label="Postal Code"
                            value={user.shippingPostalCode}
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

                    {user.role === USER_ROLE_CUSTOMER && <>
                        <FormTextChoiceRow
                          label="Question 1: How long have you been collecting for?"
                          name="howLongCollectingComicBooksForGrading"
                          value={user.howLongCollectingComicBooksForGrading}
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
                          value={user.hasPreviouslySubmittedComicBookForGrading}
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
                          value={user.hasOwnedGradedComicBooks}
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
                          value={user.hasRegularComicBookShop}
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
                          value={user.hasPreviouslyPurchasedFromAuctionSite}
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
                          value={user.hasPreviouslyPurchasedFromFacebookMarketplace}
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
                          value={user.hasRegularlyAttendedComicConsOrCollectibleShows}
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

                    <FormTextOptionRow
                      label="How did you hear about us?"
                      selectedValue={user.howDidYouHearAboutUs}
                      options={HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS}
                    />

                    <FormTextYesNoRow
                      label="I agree to receive electronic updates from my local retailer and CPS"
                      checked={user.agreePromotionsEmail}
                    />

                    {/* <p class="subtitle is-6"><FontAwesomeIcon className="fas" icon={faCog} />&nbsp;Settings</p>
                                    <hr />
                                    */}

                    <div class="columns pt-5">
                      <div class="column is-half">
                        <Link
                          class="button is-fullwidth-mobile"
                          to={`/admin/users`}
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back to Users
                        </Link>
                      </div>
                      <div class="column is-half has-text-right">
                        {user && user.status === 1 && <Link
                          to={`/admin/user/${id}/edit`}
                          class="button is-primary is-fullwidth-mobile"
                        >
                          <FontAwesomeIcon className="fas" icon={faPencil} />
                          &nbsp;Edit
                        </Link>}
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

export default AdminUserDetail;
