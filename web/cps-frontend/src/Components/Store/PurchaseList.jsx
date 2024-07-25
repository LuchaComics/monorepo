import React, { useState, useEffect } from "react";
import { Link, useParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faBuilding,
  faShoppingCart,
  faArrowLeft,
  faUsers,
  faEye,
  faPencil,
  faTrashCan,
  faPlus,
  faGauge,
  faArrowRight,
  faTable,
  faArrowUpRightFromSquare,
  faRefresh,
  faFilter,
  faSearch,
  faFilterCircleXmark,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import FormErrorBox from "../Reusable/FormErrorBox";
import { getUserPurchaseListAPI } from "../../API/UserPurchase";
import {
  topAlertMessageState,
  topAlertStatusState,
  currentUserState,
  offersFilterShowState,
  offersFilterTemporarySearchTextState,
  offersFilterActualSearchTextState,
  offersFilterStatusState,
} from "../../AppState";
import PageLoadingContent from "../Reusable/PageLoadingContent";
import FormSelectField from "../Reusable/FormSelectField";
import FormInputFieldWithButton from "../Reusable/FormInputFieldWithButton";
import {
  PAGE_SIZE_OPTIONS,
  OFFER_STATUS_OPTIONS,
} from "../../Constants/FieldOptions";
import StorePurchaseListDesktop from "./PurchaseListDesktop";
import StorePurchaseListMobile from "./PurchaseListMobile";

function StorePurchaseList() {
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
  const [currentUser] = useRecoilState(currentUserState);
  const [showFilter, setShowFilter] = useRecoilState(offersFilterShowState); // Filtering + Searching
  const [temporarySearchText, setTemporarySearchText] = useRecoilState(
    offersFilterTemporarySearchTextState,
  ); // Searching - The search field value as your writes their query.
  const [actualSearchText, setActualSearchText] = useRecoilState(
    offersFilterActualSearchTextState,
  ); // Searching - The actual search query value to submit to the API.
  const [status, setStatus] = useRecoilState(offersFilterStatusState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [listData, setListData] = useState("");
  const [isFetching, setFetching] = useState(false);
  const [pageSize, setPageSize] = useState(10); // Pagination
  const [previousCursors, setPreviousCursors] = useState([]); // Pagination
  const [nextCursor, setNextCursor] = useState(""); // Pagination
  const [currentCursor, setCurrentCursor] = useState(""); // Pagination
  const [sortField, setSortField] = useState("created_at"); // Sorting

  ////
  //// API.
  ////

  function onUserPurchaseListSuccess(response) {
    console.log("onUserPurchaseListSuccess: Starting...");
    if (response.results !== null) {
      setListData(response);
      if (response.hasNextPage) {
        setNextCursor(response.nextCursor); // For pagination purposes.
      }
    } else {
      setListData([]);
      setNextCursor("");
    }
  }

  function onUserPurchaseListError(apiErr) {
    console.log("onUserPurchaseListError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onUserPurchaseListDone() {
    console.log("onUserPurchaseListDone: Starting...");
    setFetching(false);
  }

  function onUserPurchaseListDone() {
    console.log("onUserPurchaseListDone: Starting...");
    setFetching(false);
  }

  ////
  //// Event handling.
  ////

  const fetchList = (cur, limit, keywords, status, sid) => {
    setFetching(true);
    setErrors({});

    let params = new Map();
    params.set("page_size", limit); // Pagination
    params.set("sort_field", "created_at"); // Sorting
    params.set("sort_order", -1); // Sorting - descending, meaning most recent start date to oldest start date.
    params.set("status", status);

    params.set("store_id", sid);

    if (cur !== "") {
      // Pagination
      params.set("cursor", cur);
    }

    // Filtering
    if (keywords !== undefined && keywords !== null && keywords !== "") {
      // Searhcing
      params.set("search", keywords);
    }

    getUserPurchaseListAPI(
      params,
      onUserPurchaseListSuccess,
      onUserPurchaseListError,
      onUserPurchaseListDone,
    );
  };

  const onNextClicked = (e) => {
    console.log("onNextClicked");
    let arr = [...previousCursors];
    arr.push(currentCursor);
    setPreviousCursors(arr);
    setCurrentCursor(nextCursor);
  };

  const onPreviousClicked = (e) => {
    console.log("onPreviousClicked");
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
    setStatus(2); // 1=Pending, 2=Active, 3=Archived
  };

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
      fetchList(currentCursor, pageSize, actualSearchText, status, id);
    }

    return () => {
      mounted = false;
    };
  }, [currentCursor, pageSize, actualSearchText, status, id]);

  ////
  //// Component rendering.
  ////

  return (
    <>
      <div className="container">
        <section className="section">
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
                <Link to="/store" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faBuilding} />
                  &nbsp;My Store
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faShoppingCart} />
                  &nbsp;Purchases
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/store`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to My Store
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page */}
          <nav className="box">
            <div className="columns">
              <div className="column">
                <h1 className="title is-4">
                  <FontAwesomeIcon className="fas" icon={faBuilding} />
                  &nbsp;My Store
                </h1>
              </div>
              <div className="column has-text-right">
                <button
                  onClick={() =>
                    fetchList(currentCursor, pageSize, actualSearchText, id)
                  }
                  class="is-fullwidth-mobile button is-link is-small"
                  type="button"
                >
                  <FontAwesomeIcon className="mdi" icon={faRefresh} />
                  &nbsp;
                  <span class="is-hidden-desktop is-hidden-tablet">
                    &nbsp;Refresh
                  </span>
                </button>
                {/*
                      &nbsp;
                      <button onClick={(e)=>setShowFilter(!showFilter)} class="is-fullwidth-mobile button is-small is-primary" type="button">
                          <FontAwesomeIcon className="mdi" icon={faFilter} />&nbsp;Filter
                      </button>
                      */}
              </div>
            </div>

            {/*
                {showFilter && (
                  <div class="has-background-white-bis" style={{ borderRadius: "15px", padding: "20px" }}>
                    <div class="columns">
                        <div class="column is-half">
                            <strong><u><FontAwesomeIcon className="mdi" icon={faFilter} />&nbsp;Filter</u></strong>
                        </div>
                        <div class="column is-half has-text-right">
                            <Link onClick={onClearFilterClick}><FontAwesomeIcon className="mdi" icon={faFilterCircleXmark} />&nbsp;Clear Filter</Link>
                        </div>
                    </div>
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
                            type="number"
                            placeholder="#"
                            selectedValue={status}
                            errorText={errors && errors.status}
                            helpText={
                              <ul class="content">
                                <li>pending - will not show up for members</li>
                                <li>active - will show up for everyone</li>
                                <li>archived - will be hidden from everyone</li>
                              </ul>
                            }
                            onChange={(e)=>setStatus(parseInt(e.target.value))}
                            isRequired={true}
                            options={OFFER_STATUS_OPTIONS}
                        />
                      </div>
                    </div>
                  </div>
                )}
                */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Please wait..."} />
            ) : (
              <>
                <div class="tabs is-medium is-size-7-mobile">
                  <ul>
                    <li>
                      <Link to={`/store`}>Detail</Link>
                    </li>
                    <li class={`is-active`}>
                      <Link>
                        <b>Purchases</b>
                      </Link>
                    </li>
                    <li>
                      <Link to={`/store/${id}/credits`}>Credits</Link>
                    </li>
                  </ul>
                </div>

                <FormErrorBox errors={errors} />
                {listData &&
                listData.results &&
                (listData.results.length > 0 || previousCursors.length > 0) ? (
                  <div className="container">
                    {/*
                            ##################################################################
                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A DESKTOP SCREEN.
                            ##################################################################
                        */}
                    <div class="is-hidden-touch">
                      <StorePurchaseListDesktop
                        listData={listData}
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
                      <StorePurchaseListMobile
                        listData={listData}
                        setPageSize={setPageSize}
                        pageSize={pageSize}
                        previousCursors={previousCursors}
                        onPreviousClicked={onPreviousClicked}
                        onNextClicked={onNextClicked}
                      />
                    </div>
                  </div>
                ) : (
                  <section className="hero is-medium has-background-white-ter">
                    <div className="hero-body">
                      <p className="title">
                        <FontAwesomeIcon className="fas" icon={faTable} />
                        &nbsp;No Purchases
                      </p>
                      <p className="subtitle">No purchases have been made.</p>
                    </div>
                  </section>
                )}
              </>
            )}

            <div class="columns pt-5">
              <div class="column is-half">
                <Link class="button is-fullwidth-mobile" to={`/dashboard`}>
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Dashboard
                </Link>
              </div>
              <div class="column is-half has-text-right">
                {/*
                        <Link to={`/admin/offers/add`} class="button is-success is-fullwidth-mobile"><FontAwesomeIcon className="fas" icon={faPlus} />&nbsp;New</Link>
                        */}
              </div>
            </div>
          </nav>
        </section>
      </div>
    </>
  );
}

export default StorePurchaseList;
