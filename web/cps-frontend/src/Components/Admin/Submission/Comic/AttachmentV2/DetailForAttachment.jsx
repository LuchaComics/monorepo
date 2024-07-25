import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faBook,
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
  faFile,
  faDownload,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";
import {
  ATTACHMENT_STATES,
  PAGE_SIZE_OPTIONS,
} from "../../../../Constants/FieldOptions";

import { getComicSubmissionDetailAPI } from "../../../../API/ComicSubmission";
import FormErrorBox from "../../../Reusable/FormErrorBox";
import FormInputField from "../../../Reusable/FormInputField";
import FormTextareaField from "../../../Reusable/FormTextareaField";
import FormRadioField from "../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../Reusable/FormSelectField";
import FormCheckboxField from "../../../Reusable/FormCheckboxField";
import PageLoadingContent from "../../../Reusable/PageLoadingContent";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../AppState";

function AdminComicSubmissionDetailForAttachment() {
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
  const [submission, setSubmission] = useState({});
  const [tabIndex, setTabIndex] = useState(1);
  const [attachments, setAttachments] = useState("");
  const [selectedAttachmentForDeletion, setSelectedAttachmentForDeletion] =
    useState("");

  ////
  //// Event handling.
  ////

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

  // Submission details.

  function onSubmissionDetailSuccess(response) {
    console.log("onSubmissionDetailSuccess: Starting...");
    setSubmission(response);
    setAttachments(response.fileAttachments);
  }

  function onSubmissionDetailError(apiErr) {
    console.log("onSubmissionDetailError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onSubmissionDetailDone() {
    console.log("onSubmissionDetailDone: Starting...");
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
      getComicSubmissionDetailAPI(
        id,
        onSubmissionDetailSuccess,
        onSubmissionDetailError,
        onSubmissionDetailDone,
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
          <nav class="breadcrumb" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to="/admin/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Admin Dashboard
                </Link>
              </li>
              <li class="">
                <Link to="/admin/submissions" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faTasks} />
                  &nbsp;Submissions
                </Link>
              </li>
              <li class="">
                <Link to="/admin/submissions/comics" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faBook} />
                  &nbsp;Comics
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
          <div
            class={`modal ${selectedAttachmentForDeletion ? "is-active" : ""}`}
          >
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
                You are about to <b>archive</b> this attachment; it will no
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
                  onClick={onDeselectAttachmentForDeletion}
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
                  <FontAwesomeIcon className="fas" icon={faTasks} />
                  &nbsp;Comic Submission
                </p>
              </div>
              {submission && (
                <div class="column has-text-right">
                  {/* Mobile Specific */}
                  <Link
                    to={`/admin/submissions/comic/${id}/attachments/add`}
                    class="button is-small is-success is-fullwidth is-hidden-desktop"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faPlus} />
                    &nbsp;Add Attachment
                  </Link>
                  {/* Desktop Specific */}
                  <Link
                    to={`/admin/submissions/comic/${id}/attachments/add`}
                    class="button is-small is-success is-hidden-touch"
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
                {submission && (
                  <div class="container">
                    <div class="tabs is-medium is-size-7-mobile">
                      <ul>
                        <li>
                          <Link to={`/admin/submissions/comic/${id}`}>
                            Detail
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/submissions/comic/${id}/cust`}>
                            Customer
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/submissions/comic/${id}/comments`}>
                            Comments
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/submissions/comic/${id}/file`}>
                            File
                          </Link>
                        </li>
                        <li class="is-active">
                          <Link
                            to={`/admin/submissions/comic/${id}/attachments`}
                          >
                            Attachments
                          </Link>
                        </li>
                      </ul>
                    </div>

                    {!isFetching && attachments && attachments.length > 0 ? (
                      <div class="container">
                        <p class="subtitle is-6 pt-4">
                          <FontAwesomeIcon className="fas" icon={faFile} />
                          &nbsp;Attachments
                        </p>
                        <hr />

                        <div class="b-table">
                          <div class="table-wrapper has-mobile-cards">
                            <table class="is-fullwidth is-striped is-hoverable is-fullwidth table">
                              <thead>
                                <tr>
                                  <th>Name</th>
                                  <th>Status</th>
                                  <th>Created</th>
                                  <th>File</th>
                                  <th></th>
                                </tr>
                              </thead>
                              <tbody>
                                {attachments &&
                                  attachments &&
                                  attachments.map(function (attachment, i) {
                                    return (
                                      <tr>
                                        <td data-label="Title">
                                          {attachment.name}
                                        </td>
                                        <td data-label="State">
                                          {ATTACHMENT_STATES[attachment.status]}
                                        </td>
                                        <td data-label="Created">
                                          {attachment.createdAt}
                                        </td>
                                        <td data-label="File">
                                          <a
                                            href={attachment.objectUrl}
                                            target="_blank"
                                            rel="noreferrer"
                                            class=""
                                          >
                                            <FontAwesomeIcon
                                              className="mdi"
                                              icon={faDownload}
                                            />
                                            &nbsp;Download File
                                          </a>
                                        </td>
                                        <td class="is-actions-cell">
                                          <div class="buttons is-right">
                                            <Link
                                              to={`/admin/submissions/comic/${submission.id}/attachment/${attachment.id}`}
                                              class="button is-small is-primary"
                                              type="button"
                                            >
                                              View
                                            </Link>
                                            <Link
                                              to={`/admin/submissions/comic/${submission.id}/attachment/${attachment.id}/edit`}
                                              class="button is-small is-warning"
                                              type="button"
                                            >
                                              Edit
                                            </Link>
                                            <button
                                              onClick={(e, ses) =>
                                                onSelectAttachmentForDeletion(
                                                  e,
                                                  attachment,
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
                                  })}
                              </tbody>
                            </table>
                          </div>
                        </div>
                      </div>
                    ) : (
                      <div class="container">
                        <article class="message is-dark">
                          <div class="message-body">
                            No attachments.{" "}
                            <b>
                              <Link
                                to={`/admin/submissions/comic/${id}/attachments/add`}
                              >
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
                          class="button is-hidden-touch"
                          to={`/admin/submissions`}
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back
                        </Link>
                        <Link
                          class="button is-fullwidth is-hidden-desktop"
                          to={`/submissions`}
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back
                        </Link>
                      </div>
                      <div class="column is-half has-text-right">
                        <Link
                          to={`/admin/submissions/comic/${id}/attachments/add`}
                          class="button is-primary is-hidden-touch"
                        >
                          <FontAwesomeIcon className="fas" icon={faPlus} />
                          &nbsp;Add Attachment
                        </Link>
                        <Link
                          to={`/admin/submissions/comic/${id}/attachments/add`}
                          class="button is-primary is-fullwidth is-hidden-desktop"
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

export default AdminComicSubmissionDetailForAttachment;
