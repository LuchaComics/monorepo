import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faTachometer,
  faPlus,
  faEye,
  faArrowLeft,
  faCheckCircle,
  faPencil,
  faGauge,
  faBook,
  faMagnifyingGlass,
  faBalanceScale,
  faUser,
  faArrowUpRightFromSquare,
  faComments,
  faUsers,
  faUserCircle,
  faBuilding,
} from "@fortawesome/free-solid-svg-icons";
import Select from "react-select";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";
import { DateTime } from "luxon";

import useLocalStorage from "../../../Hooks/useLocalStorage";
import {
  getTenantDetailAPI,
  postTenantCreateCommentOperationAPI,
} from "../../../API/tenant";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormInputField from "../../Reusable/FormInputField";
import FormTextareaField from "../../Reusable/FormTextareaField";
import FormRadioField from "../../Reusable/FormRadioField";
import FormMultiSelectField from "../../Reusable/FormMultiSelectField";
import FormCheckboxField from "../../Reusable/FormCheckboxField";
import FormSelectField from "../../Reusable/FormSelectField";
import FormDateField from "../../Reusable/FormDateField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import {
  FINDING_OPTIONS,
  OVERALL_NUMBER_GRADE_OPTIONS,
  PUBLISHER_NAME_OPTIONS,
  CPS_PERCENTAGE_GRADE_OPTIONS,
  HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS,
} from "../../../Constants/FieldOptions";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";

function AdminTenantDetailForCommentList() {
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
  const [tenant, setTenant] = useState({});
  const [showTenantEditOptions, setShowTenantEditOptions] = useState(false);
  const [content, setContent] = useState("");

  ////
  //// Event handling.
  ////

  const onSubmitClick = (e) => {
    console.log("onSubmitClick: Beginning..."); // Submit to the backend.
    console.log("onSubmitClick, tenant:", tenant);
    setErrors(null);
    postTenantCreateCommentOperationAPI(
      id,
      content,
      onTenantUpdateSuccess,
      onTenantUpdateError,
      onTenantUpdateDone,
      onUnauthorized,
    );
  };

  ////
  //// API.
  ////

  function onTenantDetailSuccess(response) {
    console.log("onTenantDetailSuccess: Starting...");
    setTenant(response);
  }

  function onTenantDetailError(apiErr) {
    console.log("onTenantDetailError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onTenantDetailDone() {
    console.log("onTenantDetailDone: Starting...");
    setFetching(false);
  }

  function onTenantUpdateSuccess(response) {
    // For debugging purposes only.
    console.log("onTenantUpdateSuccess: Starting...");
    console.log(response);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Comment created");
    setTopAlertStatus("success");
    setTimeout(() => {
      console.log("onTenantUpdateSuccess: Delayed for 2 seconds.");
      console.log(
        "onTenantUpdateSuccess: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Reset content.
    setContent("");

    // Fetch latest data.
    getTenantDetailAPI(
      id,
      onTenantDetailSuccess,
      onTenantDetailError,
      onTenantDetailDone,
      onUnauthorized,
    );
  }

  function onTenantUpdateError(apiErr) {
    console.log("onTenantUpdateError: Starting...");
    setErrors(apiErr);

    // Add a temporary banner message in the app and then clear itself after 2 seconds.
    setTopAlertMessage("Failed submitting");
    setTopAlertStatus("danger");
    setTimeout(() => {
      console.log("onTenantUpdateError: Delayed for 2 seconds.");
      console.log(
        "onTenantUpdateError: topAlertMessage, topAlertStatus:",
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

  function onTenantUpdateDone() {
    console.log("onTenantUpdateDone: Starting...");
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
      getTenantDetailAPI(
        id,
        onTenantDetailSuccess,
        onTenantDetailError,
        onTenantDetailDone,
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
      <div class={`modal ${showTenantEditOptions ? "is-active" : ""}`}>
        <div class="modal-background"></div>
        <div class="modal-card">
          <header class="modal-card-head">
            <p class="modal-card-title">Tenant Edit</p>
            <button
              class="delete"
              aria-label="close"
              onClick={(e) => setShowTenantEditOptions(false)}
            ></button>
          </header>
          <section class="modal-card-body">
            To edit the tenant, please select one of the following option:
            {/*
                            <br /><br />
                            <Link to={`/tenant/${tenant.id}/edit-tenant`} class="button is-primary" disabled={true}>Edit Current Tenant</Link> */}
            <br />
            <br />
            <Link
              to={`/admin/tenant/${tenant.id}/tenant/search`}
              class="button is-primary"
            >
              Pick a Different Tenant
            </Link>
          </section>
          <footer class="modal-card-foot">
            <button
              class="button"
              onClick={(e) => setShowTenantEditOptions(false)}
            >
              Close
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
                <Link to="/admin/tenants" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faBuilding} />
                  &nbsp;Tenants
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail (Comments)
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/admin/tenants`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Tenants
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
                  &nbsp;Tenant
                </p>
              </div>
              {/* HIDDEN */}
              <div class="is-hidden column has-text-right">
                {/* Mobile Specific */}
                <Link
                  to={`/admin/submissions/comics/add?tenant_id=${id}&tenant_name=${tenant.name}`}
                  class="button is-small is-success is-fullwidth is-hidden-desktop"
                  type="button"
                >
                  <FontAwesomeIcon className="mdi" icon={faPlus} />
                  &nbsp;CPS
                </Link>
                {/* Desktop Specific */}
                <Link
                  to={`/admin/submissions/comics/add?tenant_id=${id}&tenant_name=${tenant.name}`}
                  class="button is-small is-success is-hidden-touch"
                  type="button"
                >
                  <FontAwesomeIcon className="mdi" icon={faPlus} />
                  &nbsp;CPS
                </Link>
              </div>
            </div>
            <FormErrorBox errors={errors} />

            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                {tenant && (
                  <div class="container">
                    <div class="tabs is-medium is-size-7-mobile">
                      <ul>
                        <li>
                          <Link to={`/admin/tenant/${id}`}>Detail</Link>
                        </li>
                        <li>
                          <Link to={`/admin/tenant/${tenant.id}/users`}>
                            Users
                          </Link>
                        </li>
                        <li class="is-active">
                          <Link>
                            <b>Comments</b>
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/tenant/${tenant.id}/attachments`}>
                            Attachments
                          </Link>
                        </li>
                      </ul>
                    </div>

                    <p class="subtitle is-6 pt-4">
                      <FontAwesomeIcon className="fas" icon={faComments} />
                      &nbsp;Comments
                    </p>
                    <hr />

                    {tenant.comments && tenant.comments.length > 0 && (
                      <>
                        {tenant.comments.map(function (comment, i) {
                          console.log(comment); // For debugging purposes only.
                          return (
                            <div className="pb-3">
                              <span class="is-pulled-right has-text-grey-light">
                                {comment.createdByName} at{" "}
                                <b>
                                  {DateTime.fromISO(
                                    comment.createdAt,
                                  ).toLocaleString(DateTime.DATETIME_MED)}
                                </b>
                              </span>
                              <br />
                              <article class="message">
                                <div class="message-body">
                                  {comment.content}
                                </div>
                              </article>
                            </div>
                          );
                        })}
                      </>
                    )}

                    <div class="has-background-success-light mt-4 block p-3">
                      <FormTextareaField
                        label="Write your comment here:"
                        name="content"
                        placeholder="Text input"
                        value={content}
                        errorText={errors && errors.content}
                        helpText=""
                        onChange={(e) => setContent(e.target.value)}
                        isRequired={true}
                        maxWidth="180px"
                      />
                    </div>

                    <div class="columns pt-4">
                      <div class="column is-half">
                        <Link
                          to={`/admin/tenants`}
                          class="button is-fullwidth-mobile"
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back to Tenants
                        </Link>
                      </div>
                      <div class="column is-half has-text-right">
                        <button
                          onClick={onSubmitClick}
                          class="button is-primary is-fullwidth-mobile"
                        >
                          <FontAwesomeIcon className="fas" icon={faPlus} />
                          &nbsp;Add Comment
                        </button>
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

export default AdminTenantDetailForCommentList;
