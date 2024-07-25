import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faTachometer,
  faBuilding,
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
  faFile,
  faDownload,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";
import {
  ATTACHMENT_STATES,
  PAGE_SIZE_OPTIONS,
} from "../../../Constants/FieldOptions";

import { getStoreDetailAPI } from "../../../API/store";
import {
  getAttachmentListAPI,
  deleteAttachmentAPI,
} from "../../../API/Attachment";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormInputField from "../../Reusable/FormInputField";
import FormTextareaField from "../../Reusable/FormTextareaField";
import FormRadioField from "../../Reusable/FormRadioField";
import FormMultiSelectField from "../../Reusable/FormMultiSelectField";
import FormSelectField from "../../Reusable/FormSelectField";
import FormCheckboxField from "../../Reusable/FormCheckboxField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";
import AdminStoreDetailForAttachmentListDesktop from "./DetailForAttachmentListDektop";
import AdminStoreDetailForAttachmentListMobile from "./DetailForAttachmentListMobile";

function AdminStoreDetailForAttachmentList() {
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
  const [attachments, setAttachments] = useState("");
  const [selectedAttachmentForDeletion, setSelectedAttachmentForDeletion] =
    useState("");
  const [pageSize, setPageSize] = useState(10); // Pagination
  const [previousCursors, setPreviousCursors] = useState([]); // Pagination
  const [nextCursor, setNextCursor] = useState(""); // Pagination
  const [currentCursor, setCurrentCursor] = useState(""); // Pagination

  ////
  //// Event handling.
  ////

  const fetchAttachmentList = (cur, storeID, limit) => {
    setFetching(true);
    setErrors({});

    let params = new Map();
    params.set("ownership_id", id);
    params.set("ownership_role", 3); // 3=Store.
    params.set("page_size", limit);
    if (cur !== "") {
      params.set("cursor", cur);
    }

    getAttachmentListAPI(
      params,
      onAttachmentListSuccess,
      onAttachmentListError,
      onAttachmentListDone,
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

  const onSelectAttachmentForDeletion = (e, attachment) => {
    console.log("onSelectAttachmentForDeletion", attachment);
    setSelectedAttachmentForDeletion(attachment);
  };

  const onDeselectAttachmentForDeletion = (e) => {
    console.log("onDeselectAttachmentForDeletion");
    setSelectedAttachmentForDeletion("");
  };

  const onDeleteConfirmButtonClick = (e) => {
    console.log("onDeleteConfirmButtonClick"); // For debugging purposes only.

    deleteAttachmentAPI(
      selectedAttachmentForDeletion.id,
      onAttachmentDeleteSuccess,
      onAttachmentDeleteError,
      onAttachmentDeleteDone,
      onUnauthorized,
    );
    setSelectedAttachmentForDeletion("");
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

  // Attachment list.

  function onAttachmentListSuccess(response) {
    console.log("onAttachmentListSuccess: Starting...");
    if (response.results !== null) {
      setAttachments(response);
      if (response.hasNextPage) {
        setNextCursor(response.nextCursor); // For pagination purposes.
      }
    }
  }

  function onAttachmentListError(apiErr) {
    console.log("onAttachmentListError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onAttachmentListDone() {
    console.log("onAttachmentListDone: Starting...");
    setFetching(false);
  }

  // Attachment delete.

  function onAttachmentDeleteSuccess(response) {
    console.log("onAttachmentDeleteSuccess: Starting..."); // For debugging purposes only.

    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("Attachment deleted");
    setTimeout(() => {
      console.log(
        "onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Fetch again an updated list.
    fetchAttachmentList(currentCursor, id, pageSize);
  }

  function onAttachmentDeleteError(apiErr) {
    console.log("onAttachmentDeleteError: Starting..."); // For debugging purposes only.
    setErrors(apiErr);

    // Update notification.
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed deleting");
    setTimeout(() => {
      console.log(
        "onAttachmentDeleteError: topAlertMessage, topAlertStatus:",
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

  function onAttachmentDeleteDone() {
    console.log("onAttachmentDeleteDone: Starting...");
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
      fetchAttachmentList(currentCursor, id, pageSize);
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
      <div class={`modal ${selectedAttachmentForDeletion ? "is-active" : ""}`}>
        <div class="modal-background"></div>
        <div class="modal-card">
          <header class="modal-card-head">
            <p class="modal-card-title">Are you sure?</p>
            <button
              class="delete"
              aria-label="close"
              onClick={onDeselectAttachmentForDeletion}
            ></button>
          </header>
          <section class="modal-card-body">
            You are about to <b>archive</b> this attachment; it will no longer
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
            <button class="button" onClick={onDeselectAttachmentForDeletion}>
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
                  &nbsp;Detail (Attachments)
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
              {store && (
                <div class="column has-text-right">
                  <Link
                    to={`/admin/store/${id}/attachments/add`}
                    class="button is-small is-success is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faPlus} />
                    &nbsp;Add Attachment
                  </Link>
                </div>
              )}
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
                        <li>
                          <Link to={`/admin/store/${store.id}/comics`}>
                            Comics
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/store/${store.id}/comments`}>
                            Comments
                          </Link>
                        </li>
                        <li class="is-active">
                          <Link to={`/admin/store/${store.id}/attachments`}>
                            <b>Attachments</b>
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/store/${store.id}/purchases`}>
                            Purchases
                          </Link>
                        </li>
                      </ul>
                    </div>

                    {!isFetching &&
                    attachments &&
                    attachments.results &&
                    (attachments.results.length > 0 ||
                      previousCursors.length > 0) ? (
                      <div class="container">
                        {/*
                                                ##################################################################
                                                EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A DESKTOP SCREEN.
                                                ##################################################################
                                            */}
                        <div class="is-hidden-touch">
                          <AdminStoreDetailForAttachmentListDesktop
                            storeID={store.id}
                            listData={attachments}
                            setPageSize={setPageSize}
                            pageSize={pageSize}
                            previousCursors={previousCursors}
                            onPreviousClicked={onPreviousClicked}
                            onNextClicked={onNextClicked}
                            onSelectAttachmentForDeletion={
                              onSelectAttachmentForDeletion
                            }
                          />
                        </div>

                        {/*
                                                ###########################################################################
                                                EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A TABLET OR MOBILE SCREEN.
                                                ###########################################################################
                                            */}
                        <div class="is-fullwidth is-hidden-desktop">
                          <AdminStoreDetailForAttachmentListMobile
                            listData={attachments}
                            setPageSize={setPageSize}
                            pageSize={pageSize}
                            previousCursors={previousCursors}
                            onPreviousClicked={onPreviousClicked}
                            onNextClicked={onNextClicked}
                            onSelectAttachmentForDeletion={
                              onSelectAttachmentForDeletion
                            }
                          />
                        </div>
                      </div>
                    ) : (
                      <div class="container">
                        <article class="message is-dark">
                          <div class="message-body">
                            No attachments.{" "}
                            <b>
                              <Link to={`/admin/store/${id}/attachments/add`}>
                                Click here&nbsp;
                                <FontAwesomeIcon
                                  className="mdi"
                                  icon={faArrowRight}
                                />
                              </Link>
                            </b>{" "}
                            to get started creating a new attachment.
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
                          to={`/admin/store/${id}/attachments/add`}
                          class="button is-primary is-fullwidth-mobile"
                        >
                          <FontAwesomeIcon className="fas" icon={faPlus} />
                          &nbsp;Add Attachment
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

export default AdminStoreDetailForAttachmentList;
