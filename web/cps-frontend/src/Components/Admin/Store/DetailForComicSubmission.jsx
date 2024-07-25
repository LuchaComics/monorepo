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
  faArrowRight,
  faTrashCan,
  faArrowUpRightFromSquare,
  faBuilding,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";
import {
  SUBMISSION_STATES,
  PAGE_SIZE_OPTIONS,
} from "../../../Constants/FieldOptions";

import useLocalStorage from "../../../Hooks/useLocalStorage";
import { getStoreDetailAPI } from "../../../API/store";
import {
  getComicSubmissionListAPI,
  deleteComicSubmissionAPI,
} from "../../../API/ComicSubmission";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormInputField from "../../Reusable/FormInputField";
import FormTextareaField from "../../Reusable/FormTextareaField";
import FormRadioField from "../../Reusable/FormRadioField";
import FormMultiSelectField from "../../Reusable/FormMultiSelectField";
import FormSelectField from "../../Reusable/FormSelectField";
import FormCheckboxField from "../../Reusable/FormCheckboxField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";

function AdminStoreDetailForComicSubmission() {
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
  const [store, setStore] = useState({});
  const [tabIndex, setTabIndex] = useState(1);
  const [submissions, setComicSubmissions] = useState("");
  const [
    selectedComicSubmissionForDeletion,
    setSelectedComicSubmissionForDeletion,
  ] = useState("");
  const [pageSize, setPageSize] = useState(10); // Pagination
  const [previousCursors, setPreviousCursors] = useState([]); // Pagination
  const [nextCursor, setNextCursor] = useState(""); // Pagination
  const [currentCursor, setCurrentCursor] = useState(""); // Pagination

  ////
  //// Event handling.
  ////

  const fetchSubmissionList = (cur, storeID, limit) => {
    setFetching(true);
    setErrors({});

    let params = new Map();
    params.set("store_id", id);
    params.set("page_size", limit);
    if (cur !== "") {
      params.set("cursor", cur);
    }

    getComicSubmissionListAPI(
      params,
      onComicSubmissionListSuccess,
      onComicSubmissionListError,
      onComicSubmissionListDone,
      onUnauthorized,
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

  const onSelectComicSubmissionForDeletion = (e, submission) => {
    console.log("onSelectComicSubmissionForDeletion", submission);
    setSelectedComicSubmissionForDeletion(submission);
  };

  const onDeselectComicSubmissionForDeletion = (e) => {
    console.log("onDeselectComicSubmissionForDeletion");
    setSelectedComicSubmissionForDeletion("");
  };

  const onDeleteConfirmButtonClick = (e) => {
    console.log("onDeleteConfirmButtonClick"); // For debugging purposes only.

    deleteComicSubmissionAPI(
      selectedComicSubmissionForDeletion.id,
      onComicSubmissionDeleteSuccess,
      onComicSubmissionDeleteError,
      onComicSubmissionDeleteDone,
      onUnauthorized,
    );
    setSelectedComicSubmissionForDeletion("");
  };

  ////
  //// API.
  ////

  // Store details.

  function onStoreDetailSuccess(response) {
    console.log("onStoreDetailSuccess: Starting...");
    setStore(response);
  }

  function onStoreDetailError(apiErr) {
    console.log("onStoreDetailError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onStoreDetailDone() {
    console.log("onStoreDetailDone: Starting...");
    setFetching(false);
  }

  // ComicSubmission list.

  function onComicSubmissionListSuccess(response) {
    console.log("onComicSubmissionListSuccess: Starting...");
    if (response.results !== null) {
      setComicSubmissions(response);
      if (response.hasNextPage) {
        setNextCursor(response.nextCursor); // For pagination purposes.
      }
    }
  }

  function onComicSubmissionListError(apiErr) {
    console.log("onComicSubmissionListError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onComicSubmissionListDone() {
    console.log("onComicSubmissionListDone: Starting...");
    setFetching(false);
  }

  // ComicSubmission delete.

  function onComicSubmissionDeleteSuccess(response) {
    console.log("onComicSubmissionDeleteSuccess: Starting..."); // For debugging purposes only.

    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("ComicSubmission deleted");
    setTimeout(() => {
      console.log(
        "onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Fetch again an updated list.
    fetchSubmissionList(currentCursor, id, pageSize);
  }

  function onComicSubmissionDeleteError(apiErr) {
    console.log("onComicSubmissionDeleteError: Starting..."); // For debugging purposes only.
    setErrors(apiErr);

    // Update notification.
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed deleting");
    setTimeout(() => {
      console.log(
        "onComicSubmissionDeleteError: topAlertMessage, topAlertStatus:",
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

  function onComicSubmissionDeleteDone() {
    console.log("onComicSubmissionDeleteDone: Starting...");
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
      getStoreDetailAPI(
        id,
        onStoreDetailSuccess,
        onStoreDetailError,
        onStoreDetailDone,
        onUnauthorized,
      );

      fetchSubmissionList(currentCursor, id, pageSize);
    }

    return () => {
      mounted = false;
    };
  }, [currentCursor, id, pageSize]);

  ////
  //// Component rendering.
  ////

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <>
      {/* Modals */}
      <div
        class={`modal ${selectedComicSubmissionForDeletion ? "is-active" : ""}`}
      >
        <div class="modal-background"></div>
        <div class="modal-card">
          <header class="modal-card-head">
            <p class="modal-card-title">Are you sure?</p>
            <button
              class="delete"
              aria-label="close"
              onClick={onDeselectComicSubmissionForDeletion}
            ></button>
          </header>
          <section class="modal-card-body">
            You are about to <b>archive</b> this submission; it will no longer
            appear on your dashboard This action can be undone but you'll need
            to contact the system administrator. Are you sure you would like to
            continue?
          </section>
          <footer class="modal-card-foot">
            <button
              class="button is-success"
              onClick={onDeleteConfirmButtonClick}
            >
              Confirm
            </button>
            <button
              class="button"
              onClick={onDeselectComicSubmissionForDeletion}
            >
              Cancel
            </button>
          </footer>
        </div>
      </div>

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
                <Link to="/admin/stores" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faBuilding} />
                  &nbsp;Stores
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail (Comics)
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/admin/stores`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Stores
                </Link>
              </li>
            </ul>
          </nav>

          {/* Page */}
          <nav class="box">
            <div class="columns">
              <div class="column">
                <p class="title is-4">
                  <FontAwesomeIcon className="fas" icon={faBuilding} />
                  &nbsp;Store
                </p>
              </div>
              <div class="column has-text-right">
                <Link
                  to={`/admin/submissions/comics/add/search`}
                  target="_blank"
                  rel="noreferrer"
                  class="button is-small is-success is-fullwidth-mobile"
                  type="button"
                >
                  <FontAwesomeIcon className="mdi" icon={faPlus} />
                  &nbsp;CPS&nbsp;
                  <FontAwesomeIcon
                    className="fas"
                    icon={faArrowUpRightFromSquare}
                  />
                </Link>
              </div>
            </div>
            <FormErrorBox errors={errors} />

            {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                {store && (
                  <div class="container">
                    <div class="tabs is-medium is-size-7-mobile">
                      <ul>
                        <li>
                          <Link to={`/admin/store/${store.id}`}>Detail</Link>
                        </li>
                        <li>
                          <Link to={`/admin/store/${store.id}/users`}>
                            Users
                          </Link>
                        </li>
                        <li class="is-active">
                          <Link>
                            <b>Comics</b>
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/store/${store.id}/comments`}>
                            Comments
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/store/${store.id}/attachments`}>
                            Attachments
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/store/${store.id}/purchases`}>
                            Purchases
                          </Link>
                        </li>
                      </ul>
                    </div>

                    <p class="subtitle is-6 pt-4">
                      <FontAwesomeIcon className="fas" icon={faTasks} />
                      &nbsp;Online Comic Submissions
                    </p>
                    <hr />

                    {!isFetching &&
                    submissions &&
                    submissions.results &&
                    (submissions.results.length > 0 ||
                      previousCursors.length > 0) ? (
                      <div class="container">
                        <div class="b-table">
                          <div class="table-wrapper has-mobile-cards">
                            <table class="is-fullwidth is-striped is-hoverable is-fullwidth table">
                              <thead>
                                <tr>
                                  <th>Title</th>
                                  <th>Vol</th>
                                  <th>No</th>
                                  <th>Status</th>
                                  <th>Created</th>
                                  <th></th>
                                </tr>
                              </thead>
                              <tbody>
                                {submissions &&
                                  submissions.results &&
                                  submissions.results.map(
                                    function (submission, i) {
                                      return (
                                        <tr>
                                          <td data-label="Title">
                                            {submission.seriesTitle}
                                          </td>
                                          <td data-label="Vol">
                                            {submission.issueVol}
                                          </td>
                                          <td data-label="No">
                                            {submission.issueNo}
                                          </td>
                                          <td data-label="State">
                                            {
                                              SUBMISSION_STATES[
                                                submission.status
                                              ]
                                            }
                                          </td>
                                          <td data-label="Created">
                                            {submission.createdAt}
                                          </td>
                                          <td class="is-actions-cell">
                                            <div class="buttons is-right">
                                              <Link
                                                to={`/admin/submissions/comic/${submission.id}`}
                                                target="_blank"
                                                rel="noreferrer"
                                                class="button is-small is-primary"
                                                type="button"
                                              >
                                                View&nbsp;
                                                <FontAwesomeIcon
                                                  className="fas"
                                                  icon={
                                                    faArrowUpRightFromSquare
                                                  }
                                                />
                                              </Link>
                                              <Link
                                                to={`/admin/submissions/comic/${submission.id}/edit`}
                                                target="_blank"
                                                rel="noreferrer"
                                                class="button is-small is-warning"
                                                type="button"
                                              >
                                                Edit&nbsp;
                                                <FontAwesomeIcon
                                                  className="fas"
                                                  icon={
                                                    faArrowUpRightFromSquare
                                                  }
                                                />
                                              </Link>
                                              <button
                                                onClick={(e, ses) =>
                                                  onSelectComicSubmissionForDeletion(
                                                    e,
                                                    submission,
                                                  )
                                                }
                                                class="button is-small is-danger"
                                                type="button"
                                              >
                                                <FontAwesomeIcon
                                                  className="mdi"
                                                  icon={faTrashCan}
                                                />
                                                &nbsp;Delete
                                              </button>
                                            </div>
                                          </td>
                                        </tr>
                                      );
                                    },
                                  )}
                              </tbody>
                            </table>

                            <div class="columns">
                              <div class="column is-half">
                                <span class="select">
                                  <select
                                    class={`input has-text-grey-light`}
                                    name="pageSize"
                                    onChange={(e) =>
                                      setPageSize(parseInt(e.target.value))
                                    }
                                  >
                                    {PAGE_SIZE_OPTIONS.map(
                                      function (option, i) {
                                        return (
                                          <option
                                            selected={pageSize === option.value}
                                            value={option.value}
                                          >
                                            {option.label}
                                          </option>
                                        );
                                      },
                                    )}
                                  </select>
                                </span>
                              </div>
                              <div class="column is-half has-text-right">
                                {previousCursors.length > 0 && (
                                  <button
                                    class="button"
                                    onClick={onPreviousClicked}
                                  >
                                    Previous
                                  </button>
                                )}
                                {submissions.hasNextPage && (
                                  <>
                                    <button
                                      class="button"
                                      onClick={onNextClicked}
                                    >
                                      Next
                                    </button>
                                  </>
                                )}
                              </div>
                            </div>
                          </div>
                        </div>
                      </div>
                    ) : (
                      <div class="container">
                        <article class="message is-dark">
                          <div class="message-body">
                            No submissions.{" "}
                            <b>
                              <Link to="/admin/submissions/comics/add/search">
                                Click here&nbsp;
                                <FontAwesomeIcon
                                  className="mdi"
                                  icon={faArrowRight}
                                />
                              </Link>
                            </b>{" "}
                            to get started creating a new submission.
                          </div>
                        </article>
                      </div>
                    )}

                    <div class="columns pt-5">
                      <div class="column is-half">
                        <Link
                          class="button is-fullwidth-mobile"
                          to={`/admin/stores`}
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back to Stores
                        </Link>
                      </div>
                      <div class="column is-half has-text-right">
                        <Link
                          to={`/admin/submissions/comics/add/search`}
                          target="_blank"
                          rel="noreferrer"
                          class="button is-primary is-fullwidth-mobile"
                        >
                          <FontAwesomeIcon className="fas" icon={faPlus} />
                          &nbsp;CPS&nbsp;
                          <FontAwesomeIcon
                            className="fas"
                            icon={faArrowUpRightFromSquare}
                          />
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

export default AdminStoreDetailForComicSubmission;
