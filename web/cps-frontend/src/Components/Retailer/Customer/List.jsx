import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faFilter,
  faFilterCircleXmark,
  faArrowLeft,
  faUsers,
  faTachometer,
  faEye,
  faPencil,
  faTrashCan,
  faPlus,
  faGauge,
  faArrowRight,
  faTable,
  faRefresh,
  faSearch,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { getCustomerListAPI, deleteCustomerAPI } from "../../../API/customer";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
  customerFilterShowState,
  customerFilterTemporarySearchTextState,
  customerFilterActualSearchTextState,
  customerFilterStatusState,
  customerFilterJoinedAfterState,
} from "../../../AppState";
import {
  PAGE_SIZE_OPTIONS,
  USER_STATUS_LIST_OPTIONS,
} from "../../../Constants/FieldOptions";
import FormErrorBox from "../../Reusable/FormErrorBox";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import FormInputFieldWithButton from "../../Reusable/FormInputFieldWithButton";
import FormSelectField from "../../Reusable/FormSelectField";
import FormDateField from "../../Reusable/FormDateField";
import RetailerCustomerListDesktop from "./ListDesktop";
import RetailerCustomerListMobile from "./ListMobile";

function RetailerCustomerList() {
  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);
  const [currentUser] = useRecoilState(currentUserState);
  const [showFilter, setShowFilter] = useRecoilState(customerFilterShowState); // Filtering + Searching
  const [temporarySearchText, setTemporarySearchText] = useRecoilState(
    customerFilterTemporarySearchTextState,
  ); // Searching - The search field value as your writes their query.
  const [actualSearchText, setActualSearchText] = useRecoilState(
    customerFilterActualSearchTextState,
  ); // Searching - The actual search query value to submit to the API.
  const [status, setStatus] = useRecoilState(customerFilterStatusState); // Filtering
  const [createdAtGTE, setCreatedAtGTE] = useRecoilState(
    customerFilterJoinedAfterState,
  ); // Filtering

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [forceURL, setForceURL] = useState("");
  const [customers, setCustomers] = useState("");
  const [isFetching, setFetching] = useState(false);
  const [pageSize, setPageSize] = useState(10); // Pagination
  const [previousCursors, setPreviousCursors] = useState([]); // Pagination
  const [nextCursor, setNextCursor] = useState(""); // Pagination
  const [currentCursor, setCurrentCursor] = useState(""); // Pagination
  const [sortField, setSortField] = useState("created_at"); // Sorting

  ////
  //// API.
  ////

  // --- LIST --- //

  function onCustomerListSuccess(response) {
    console.log("onCustomerListSuccess: Starting...");
    if (response.results !== null) {
      setCustomers(response);
      if (response.hasNextPage) {
        setNextCursor(response.nextCursor); // For pagination purposes.
      }
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

  // --- ALL --- //

  const onUnauthorized = () => {
    setForceURL("/login?unauthorized=true"); // If token expired or user is not logged in, redirect back to login.
  };

  ////
  //// Event handling.
  ////

  const fetchList = (cur, limit, keywords, s, cagte) => {
    setFetching(true);
    setErrors({});

    let params = new Map();
    params.set("page_size", limit); // Pagination
    params.set("sort_field", "created_at"); // Sorting

    if (cur !== "") {
      // Pagination
      params.set("cursor", cur);
    }

    // Filtering
    if (keywords !== undefined && keywords !== null && keywords !== "") {
      // Searhcing
      params.set("search", keywords);
    }
    if (s !== undefined && s !== null && s !== "") {
      params.set("status", s);
    }
    if (cagte !== undefined && cagte !== null && cagte !== "") {
      // DEVELOPERS NOTE:
      // If the value is a string, assume it's the correct format which will
      // translate into a `Date` object.
      if (cagte instanceof Date === false) {
        const cagteStr = new Date(cagte).getTime();
        params.set("created_at_gte", cagteStr);
      } else {
        const cagteStr = cagte.getTime();
        params.set("created_at_gte", cagteStr);
      }
    }

    console.log("fetchList | params:", params);

    getCustomerListAPI(
      params,
      onCustomerListSuccess,
      onCustomerListError,
      onCustomerListDone,
      onUnauthorized,
    );
  };

  const onNextClicked = (e) => {
    let arr = [...previousCursors];
    arr.push(currentCursor);
    setPreviousCursors(arr);
    setCurrentCursor(nextCursor);
  };

  const onPreviousClicked = (e) => {
    let arr = [...previousCursors];
    const previousCursor = arr.pop();
    setPreviousCursors(arr);
    setCurrentCursor(previousCursor);
  };

  const onSearchButtonClick = (e) => {
    // Searching
    console.log("Search button clicked...");
    setActualSearchText(temporarySearchText);
  };

  // Function resets the filter state to its default state.
  const onClearFilterClick = (e) => {
    setShowFilter(false);
    setActualSearchText("");
    setTemporarySearchText("");
    setStatus("");
    setCreatedAtGTE("");
  };

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
      fetchList(
        currentCursor,
        pageSize,
        actualSearchText,
        status,
        createdAtGTE,
      );
    }

    return () => {
      mounted = false;
    };
  }, [currentCursor, pageSize, actualSearchText, status, createdAtGTE]);

  ////
  //// Component rendering.
  ////

  console.log("state | previousCursors:", previousCursors);
  console.log("state | currentCursor:", currentCursor);
  console.log("state | nextCursor:", nextCursor);
  console.log();

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
                  <FontAwesomeIcon className="fas" icon={faUsers} />
                  &nbsp;Customers
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
            <div class="columns">
              <div class="column is-8">
                <h1 class="title is-4">
                  <FontAwesomeIcon className="fas" icon={faUsers} />
                  &nbsp;Customer List for:{" "}
                  {currentUser && currentUser.storeName}
                </h1>
              </div>

              {/* Mobile only */}
              <div className="column has-text-right is-hidden-desktop is-hidden-tablet">
                <button
                  onClick={() =>
                    fetchList(
                      currentCursor,
                      pageSize,
                      actualSearchText,
                      status,
                      createdAtGTE,
                    )
                  }
                  class="is-fullwidth button is-small is-info"
                  type="button"
                >
                  <FontAwesomeIcon className="mdi" icon={faRefresh} />
                  &nbsp;
                  <span class="is-hidden-desktop is-hidden-tablet">
                    &nbsp;Refresh
                  </span>
                </button>
                &nbsp;
                <button
                  onClick={(e) => setShowFilter(!showFilter)}
                  class="is-fullwidth button is-small is-success"
                  type="button"
                >
                  <FontAwesomeIcon className="mdi" icon={faFilter} />
                  &nbsp;Filter
                </button>
                &nbsp;
                <Link
                  to={`/customers/add`}
                  className="is-fullwidth button is-small is-primary"
                  type="button"
                >
                  <FontAwesomeIcon className="mdi" icon={faPlus} />
                  &nbsp;New
                </Link>
              </div>
              {/* Tablet and Desktop only */}
              <div className="column has-text-right is-hidden-mobile">
                <button
                  onClick={() =>
                    fetchList(
                      currentCursor,
                      pageSize,
                      actualSearchText,
                      status,
                      createdAtGTE,
                    )
                  }
                  class="button is-small is-info"
                  type="button"
                >
                  <FontAwesomeIcon className="mdi" icon={faRefresh} />
                </button>
                &nbsp;
                <button
                  onClick={(e) => setShowFilter(!showFilter)}
                  class="button is-small is-success"
                  type="button"
                >
                  <FontAwesomeIcon className="mdi" icon={faFilter} />
                  &nbsp;Filter
                </button>
                &nbsp;
                <Link
                  to={`/customers/add`}
                  className="button is-small is-primary"
                  type="button"
                >
                  <FontAwesomeIcon className="mdi" icon={faPlus} />
                  &nbsp;New
                </Link>
              </div>
            </div>

            {/* FILTER */}
            {showFilter && (
              <div
                class="has-background-white-bis"
                style={{ borderRadius: "15px", padding: "20px" }}
              >
                {/* Filter Title + Clear Button */}
                <div class="columns">
                  <div class="column is-half">
                    <strong>
                      <u>
                        <FontAwesomeIcon className="mdi" icon={faFilter} />
                        &nbsp;Filter
                      </u>
                    </strong>
                  </div>
                  <div class="column is-half has-text-right">
                    <Link onClick={onClearFilterClick}>
                      <FontAwesomeIcon
                        className="mdi"
                        icon={faFilterCircleXmark}
                      />
                      &nbsp;Clear Filter
                    </Link>
                  </div>
                </div>

                {/* Filter Options */}
                <div class="columns">
                  <div class="column">
                    <FormInputFieldWithButton
                      label={"Search"}
                      name="temporarySearchText"
                      type="text"
                      placeholder="Search by name"
                      value={temporarySearchText}
                      helpText=""
                      onChange={(e) => setTemporarySearchText(e.target.value)}
                      isRequired={true}
                      maxWidth="100%"
                      buttonLabel={
                        <>
                          <FontAwesomeIcon className="fas" icon={faSearch} />
                        </>
                      }
                      onButtonClick={onSearchButtonClick}
                    />
                  </div>
                  <div class="column">
                    <FormSelectField
                      label="Status"
                      name="status"
                      placeholder="Pick status"
                      selectedValue={status}
                      helpText=""
                      onChange={(e) => setStatus(parseInt(e.target.value))}
                      options={USER_STATUS_LIST_OPTIONS}
                      isRequired={true}
                    />
                  </div>
                  <div class="column">
                    <FormDateField
                      label="Joined After"
                      name="createdAtGTE"
                      placeholder="Text input"
                      value={createdAtGTE}
                      helpText=""
                      onChange={(date) => setCreatedAtGTE(date)}
                      isRequired={true}
                      maxWidth="120px"
                    />
                  </div>
                </div>
              </div>
            )}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                <FormErrorBox errors={errors} />
                {customers &&
                customers.results &&
                (customers.results.length > 0 || previousCursors.length > 0) ? (
                  <div class="container">
                    {/*
                                            ##################################################################
                                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A DESKTOP SCREEN.
                                            ##################################################################
                                        */}
                    <div class="is-hidden-touch">
                      <RetailerCustomerListDesktop
                        listData={customers}
                        setPageSize={setPageSize}
                        pageSize={pageSize}
                        previousCursors={previousCursors}
                        onPreviousClicked={onPreviousClicked}
                        onNextClicked={onNextClicked}
                      />
                    </div>

                    {/*
                                            ###########################################################################
                                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A TABLET OR MOBILE SCREEN.
                                            ###########################################################################
                                        */}
                    <div class="is-fullwidth is-hidden-desktop">
                      <RetailerCustomerListMobile
                        listData={customers}
                        setPageSize={setPageSize}
                        pageSize={pageSize}
                        previousCursors={previousCursors}
                        onPreviousClicked={onPreviousClicked}
                        onNextClicked={onNextClicked}
                      />
                    </div>
                  </div>
                ) : (
                  <section class="hero is-medium has-background-white-ter">
                    <div class="hero-body">
                      <p class="title">
                        <FontAwesomeIcon className="fas" icon={faTable} />
                        &nbsp;No Customers
                      </p>
                      <p class="subtitle">
                        No customers.{" "}
                        <b>
                          <Link to="/customers/add">
                            Click here&nbsp;
                            <FontAwesomeIcon
                              className="mdi"
                              icon={faArrowRight}
                            />
                          </Link>
                        </b>{" "}
                        to add your first customer.
                      </p>
                    </div>
                  </section>
                )}
              </>
            )}

            {/* Bottom navigation */}
            <div class="columns pt-4">
              <div class="column is-half">
                <Link
                  to={`/dashboard`}
                  class="button is-medium is-fullwidth-mobile"
                >
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Dashboard
                </Link>
              </div>
              <div class="column is-half has-text-right"></div>
            </div>
          </nav>
        </section>
      </div>
    </>
  );
}

export default RetailerCustomerList;
