import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faEllipsis,
  faTasks,
  faTachometer,
  faPlus,
  faArrowLeft,
  faCheckCircle,
  faProjectCircle,
  faGauge,
  faPencil,
  faProjects,
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
} from "../../../../../Constants/FieldOptions";

import { getProjectDetailAPI } from "../../../../../API/Project";
import {
  getPinObjectListAPI,
  deletePinObjectAPI,
} from "../../../../../API/PinObject";
import FormErrorBox from "../../../../Reusable/FormErrorBox";
import FormInputField from "../../../../Reusable/FormInputField";
import FormTextareaField from "../../../../Reusable/FormTextareaField";
import FormRadioField from "../../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../../Reusable/FormSelectField";
import FormCheckboxField from "../../../../Reusable/FormCheckboxField";
import PageLoadingContent from "../../../../Reusable/PageLoadingContent";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../../AppState";
import AdminProjectDetailForPinObjectListDesktop from "./DetailForPinObjectListDektop";
import AdminProjectDetailForPinObjectListMobile from "./DetailForPinObjectListMobile";
import AlertBanner from "../../../../Reusable/EveryPage/AlertBanner";

function AdminProjectDetailForPinObjectList() {
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
  const [project, setProject] = useState({});
  const [tabIndex, setTabIndex] = useState(1);
  const [pinobjects, setPinObjects] = useState("");
  const [selectedPinObjectForDeletion, setSelectedPinObjectForDeletion] =
    useState("");
  const [pageSize, setPageSize] = useState(10); // Pagination
  const [previousCursors, setPreviousCursors] = useState([]); // Pagination
  const [nextCursor, setNextCursor] = useState(""); // Pagination
  const [currentCursor, setCurrentCursor] = useState(""); // Pagination

  ////
  //// Event handling.
  ////

  const fetchPinObjectList = (cur, projectID, limit) => {
    setFetching(true);
    setErrors({});

    let params = new Map();
    params.set("ownership_id", id);
    params.set("ownership_role", 1); // 1=Project or Project.
    params.set("page_size", limit);
    if (cur !== "") {
      params.set("cursor", cur);
    }

    getPinObjectListAPI(
      params,
      onPinObjectListSuccess,
      onPinObjectListError,
      onPinObjectListDone,
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

  const onSelectPinObjectForDeletion = (e, pinobject) => {
    console.log("onSelectPinObjectForDeletion", pinobject);
    setSelectedPinObjectForDeletion(pinobject);
  };

  const onDeselectPinObjectForDeletion = (e) => {
    console.log("onDeselectPinObjectForDeletion");
    setSelectedPinObjectForDeletion("");
  };

  const onDeleteConfirmButtonClick = (e) => {
    console.log("onDeleteConfirmButtonClick"); // For debugging purposes only.

    deletePinObjectAPI(
      selectedPinObjectForDeletion.id,
      onPinObjectDeleteSuccess,
      onPinObjectDeleteError,
      onPinObjectDeleteDone,
      onUnauthorized,
    );
    setSelectedPinObjectForDeletion("");
  };

  ////
  //// API.
  ////

  // Project details.

  function onProjectDetailSuccess(response) {
    console.log("onProjectDetailSuccess: Starting...");
    setProject(response);
  }

  function onProjectDetailError(apiErr) {
    console.log("onProjectDetailError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onProjectDetailDone() {
    console.log("onProjectDetailDone: Starting...");
    setFetching(false);
  }

  // PinObject list.

  function onPinObjectListSuccess(response) {
    console.log("onPinObjectListSuccess: Starting...");
    if (response.results !== null) {
      setPinObjects(response);
      if (response.hasNextPage) {
        setNextCursor(response.nextCursor); // For pagination purposes.
      }
    }
  }

  function onPinObjectListError(apiErr) {
    console.log("onPinObjectListError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onPinObjectListDone() {
    console.log("onPinObjectListDone: Starting...");
    setFetching(false);
  }

  // PinObject delete.

  function onPinObjectDeleteSuccess(response) {
    console.log("onPinObjectDeleteSuccess: Starting..."); // For debugging purposes only.

    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("PinObject deleted");
    setTimeout(() => {
      console.log(
        "onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Fetch again an updated list.
    fetchPinObjectList(currentCursor, id, pageSize);
  }

  function onPinObjectDeleteError(apiErr) {
    console.log("onPinObjectDeleteError: Starting..."); // For debugging purposes only.
    setErrors(apiErr);

    // Update notification.
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed deleting");
    setTimeout(() => {
      console.log(
        "onPinObjectDeleteError: topAlertMessage, topAlertStatus:",
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

  function onPinObjectDeleteDone() {
    console.log("onPinObjectDeleteDone: Starting...");
    setFetching(false);
  }

  // --- All --- //

  const onUnauthorized = () => {
    setForceURL("/login?unauthorized=true"); // If token expired or project is not logged in, redirect back to login.
  };

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.

      setFetching(true);
      getProjectDetailAPI(
        id,
        onProjectDetailSuccess,
        onProjectDetailError,
        onProjectDetailDone,
        onUnauthorized,
      );
      fetchPinObjectList(currentCursor, id, pageSize);
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
                <Link to="/admin/projects" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faProjects} />
                  &nbsp;Projects
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail (PinObjects)
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/admin/projects`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Projects
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modals */}
          {/* None */}

          {/* Page banner */}
          {project && project.status === 100 && (
            <AlertBanner message="Archived" status="info" />
          )}

          {/* Page */}
          <div
            class={`modal ${selectedPinObjectForDeletion ? "is-active" : ""}`}
          >
            <div class="modal-background"></div>
            <div class="modal-card">
              <header class="modal-card-head">
                <p class="modal-card-title">Are you sure?</p>
                <button
                  class="delete"
                  aria-label="close"
                  onClick={onDeselectPinObjectForDeletion}
                ></button>
              </header>
              <section class="modal-card-body">
                You are about to <b>archive</b> this pinobject; it will no
                longer appear on your dashboard This action can be undone but
                you'll need to contact the system administrator. Are you sure
                you would like to continue?
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
                  onClick={onDeselectPinObjectForDeletion}
                >
                  Cancel
                </button>
              </footer>
            </div>
          </div>
          <nav class="box">
            <div class="columns">
              <div class="column">
                <p class="title is-4">
                  <FontAwesomeIcon className="fas" icon={faProjectCircle} />
                  &nbsp;Project
                </p>
              </div>
              {project && project.status === 1 && (
                <div class="column has-text-right">
                  <Link
                    to={`/admin/project/${id}/pinobjects/add`}
                    class="button is-small is-success is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faPlus} />
                    &nbsp;Add PinObject
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
                {project && (
                  <div class="container">
                    <div class="tabs is-medium is-size-7-mobile">
                      <ul>
                        <li>
                          <Link to={`/admin/project/${project.id}`}>Detail</Link>
                        </li>
                        <li>
                          <Link to={`/admin/project/${project.id}/comments`}>
                            Comments
                          </Link>
                        </li>
                        <li class="is-active">
                          <Link to={`/admin/project/${project.id}/pinobjects`}>
                            <b>PinObjects</b>
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/project/${project.id}/more`}>
                            More&nbsp;&nbsp;
                            <FontAwesomeIcon
                              className="mdi"
                              icon={faEllipsis}
                            />
                          </Link>
                        </li>
                      </ul>
                    </div>

                    {!isFetching &&
                    pinobjects &&
                    pinobjects.results &&
                    (pinobjects.results.length > 0 ||
                      previousCursors.length > 0) ? (
                      <div class="container">
                        {/*
                                                ##################################################################
                                                EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A DESKTOP SCREEN.
                                                ##################################################################
                                            */}
                        <div class="is-hidden-touch">
                          <AdminProjectDetailForPinObjectListDesktop
                            projectID={id}
                            listData={pinobjects}
                            setPageSize={setPageSize}
                            pageSize={pageSize}
                            previousCursors={previousCursors}
                            onPreviousClicked={onPreviousClicked}
                            onNextClicked={onNextClicked}
                            onSelectPinObjectForDeletion={
                              onSelectPinObjectForDeletion
                            }
                          />
                        </div>

                        {/*
                                                ###########################################################################
                                                EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A TABLET OR MOBILE SCREEN.
                                                ###########################################################################
                                            */}
                        <div class="is-fullwidth is-hidden-desktop">
                          <AdminProjectDetailForPinObjectListMobile
                            projectID={id}
                            listData={pinobjects}
                            setPageSize={setPageSize}
                            pageSize={pageSize}
                            previousCursors={previousCursors}
                            onPreviousClicked={onPreviousClicked}
                            onNextClicked={onNextClicked}
                            onSelectPinObjectForDeletion={
                              onSelectPinObjectForDeletion
                            }
                          />
                        </div>
                      </div>
                    ) : (
                      <div class="container">
                        <article class="message is-dark">
                          <div class="message-body">
                            No pinobjects.{" "}
                            <b>
                              <Link to={`/admin/project/${id}/pinobjects/add`}>
                                Click here&nbsp;
                                <FontAwesomeIcon
                                  className="mdi"
                                  icon={faArrowRight}
                                />
                              </Link>
                            </b>{" "}
                            to get started creating a new pinobject.
                          </div>
                        </article>
                      </div>
                    )}

                    <div class="columns pt-5">
                      <div class="column is-half">
                        <Link class="button is-fullwidth-mobile" to={`/projects`}>
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back to Projects
                        </Link>
                      </div>
                      <div class="column is-half has-text-right">
                        {project && project.status === 1 && <Link
                          to={`/admin/project/${id}/pinobjects/add`}
                          class="button is-primary is-fullwidth-mobile"
                        >
                          <FontAwesomeIcon className="fas" icon={faPlus} />
                          &nbsp;Add PinObject
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

export default AdminProjectDetailForPinObjectList;
