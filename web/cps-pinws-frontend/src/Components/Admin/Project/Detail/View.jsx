import React, { useState, useEffect } from "react";
import { Link, Navigate, useSearchParams } from "react-router-dom";
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
  faProjectDiagram,
  faGauge,
  faPencil,
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
  getProjectDetailAPI,
  postProjectStarOperationAPI,
  deleteProjectAPI,
} from "../../../../API/Project";
import { getTenantSelectOptionListAPI } from "../../../../API/tenant";
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

function AdminProjectDetail() {
  ////
  ////
  //// URL Query Strings.
  ////

  const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140
  const secretCode = searchParams.get("secret");

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
  const [tenantSelectOptions, setTenantSelectOptions] = useState([]);
  const [selectedProjectForDeletion, setSelectedProjectForDeletion] = useState("");

  ////
  //// Event handling.
  ////

  const onStarClick = () => {
    setFetching(true);
    setErrors({});
    postProjectStarOperationAPI(
      id,
      onProjectDetailSuccess,
      onProjectDetailError,
      onProjectDetailDone,
      onUnauthorized,
    );
  };

  const onSelectProjectForDeletion = (e, project) => {
    console.log("onSelectProjectForDeletion", project);
    setSelectedProjectForDeletion(project);
  };

  const onDeselectProjectForDeletion = (e) => {
    console.log("onDeselectProjectForDeletion");
    setSelectedProjectForDeletion("");
  };

  const onDeleteConfirmButtonClick = (e) => {
    console.log("onDeleteConfirmButtonClick"); // For debugging purposes only.

    deleteProjectAPI(
      selectedProjectForDeletion.id,
      onProjectDeleteSuccess,
      onProjectDeleteError,
      onProjectDeleteDone,
      onUnauthorized,
    );
    setSelectedProjectForDeletion("");
  };

  ////
  //// API.
  ////

  // --- DETAIL --- //

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

  // --- STORE OPTIONS --- //

  function onTenantOptionListSuccess(response) {
    console.log("onTenantOptionListSuccess: Starting...");
    if (response !== null) {
      const selectOptions = [
        { value: 0, label: "Please select" }, // Add empty options.
        ...response,
      ];
      setTenantSelectOptions(selectOptions);
    }
  }

  function onTenantOptionListError(apiErr) {
    console.log("onTenantOptionListError: Starting...");
    console.log("onTenantOptionListError: apiErr:", apiErr);
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onTenantOptionListDone() {
    console.log("onTenantOptionListDone: Starting...");
    setFetching(false);
  }

  // --- DELETE --- //

  function onProjectDeleteSuccess(response) {
    console.log("onProjectDeleteSuccess: Starting..."); // For debugging purposes only.

    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("Project deleted");
    setTimeout(() => {
      console.log(
        "onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Fetch again an updated list.
    setForceURL("/admin/projects");
  }

  function onProjectDeleteError(apiErr) {
    console.log("onProjectDeleteError: Starting..."); // For debugging purposes only.
    setErrors(apiErr);

    // Update notification.
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed deleting");
    setTimeout(() => {
      console.log(
        "onProjectDeleteError: topAlertMessage, topAlertStatus:",
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

  function onProjectDeleteDone() {
    console.log("onProjectDeleteDone: Starting...");
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

      let params = new Map();
      getTenantSelectOptionListAPI(
        params,
        onTenantOptionListSuccess,
        onTenantOptionListError,
        onTenantOptionListDone,
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
                <Link to="/admin/projects" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faProjectDiagram} />
                  &nbsp;Projects
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
                <Link to={`/admin/projects`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Projects
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modals */}
          <div class={`modal ${selectedProjectForDeletion ? "is-active" : ""}`}>
            <div class="modal-background"></div>
            <div class="modal-card">
              <header class="modal-card-head">
                <p class="modal-card-title">Are you sure?</p>
                <button
                  class="delete"
                  aria-label="close"
                  onClick={onDeselectProjectForDeletion}
                ></button>
              </header>
              <section class="modal-card-body">
                You are about to <b>archive</b> this project; it will no longer
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
                <button class="button" onClick={onDeselectProjectForDeletion}>
                  Cancel
                </button>
              </footer>
            </div>
          </div>
          <div class={`modal ${secretCode ? "is-active" : ""}`}>
            <div class="modal-background"></div>
            <div class="modal-card">
              <header class="modal-card-head">
                <p class="modal-card-title">Please confirm the secret</p>
                <Link
                  class="delete"
                  aria-label="close"
                  to={`/admin/project/${id}`}
                ></Link>
              </header>
              <section class="modal-card-body" style={{wordWrap: "break-word"}}>
                You have successfully created a project! Here are your credentials. Please save the <strong>secret</strong> somewhere safe as you'll never get access to it again after you click <strong>confirm</strong>.
                <br />
                <br />
                <strong>Project ID:</strong>&nbsp;<br/><i>{id}</i>
                <br />
                <br />
                <strong>Project Secret:</strong>&nbsp;<br/><i>{secretCode}</i>
                <br />
                <br />
              </section>
              <footer class="modal-card-foot">
                <Link
                  class="button is-success"
                  to={`/admin/project/${id}`}
                >
                  I confirm & close
                </Link>
              </footer>
            </div>
          </div>

          {/* Page banner */}
          {project && project.status === 100 && (
            <AlertBanner message="Archived" status="info" />
          )}

          {/* Page */}
          <nav class="box">
            {project && (
              <div class="columns">
                <div class="column">
                  <p class="title is-4">
                    <FontAwesomeIcon className="fas" icon={faProjectDiagram} />
                    &nbsp;Project
                  </p>
                </div>
                {project && project.status === 1 && <div class="column has-text-right">
                  <Link
                    to={`/admin/pins/add?project_id=${project.id}&project_name=${project.name}&tenant_id=${project.tenantId}&from=projects&clear=true`}
                    class="button is-small is-success is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faPlus} />
                    &nbsp;Pin
                  </Link>
                  &nbsp;&nbsp;
                  <Link
                    to={`/admin/project/${project.id}/edit`}
                    class="button is-small is-warning is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faPencil} />
                    &nbsp;Edit
                  </Link>
                  &nbsp;&nbsp;
                  <button
                    onClick={(e, ses) => onSelectProjectForDeletion(e, project)}
                    class="button is-small is-danger is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faTrashCan} />
                    &nbsp;Delete
                  </button>
                  &nbsp;&nbsp;
                  {project.isStarred ? (
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
                {project && (
                  <div class="container">
                    <div class="tabs is-medium is-size-7-mobile">
                      <ul>
                        <li class="is-active">
                          <Link>
                            <b>Detail</b>
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/project/${project.id}/pins`}>
                            Pins
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

                    <p class="subtitle is-6">
                      <FontAwesomeIcon className="fas" icon={faCogs} />
                      &nbsp;Settings
                    </p>
                    <hr />

                    {tenantSelectOptions && tenantSelectOptions.length > 0 && (
                      <FormTextOptionRow
                        label="Tenant"
                        selectedValue={project.tenantID}
                        helpText=""
                        options={tenantSelectOptions}
                      />
                    )}
                    <FormTextChoiceRow
                      label="Status"
                      value={project.status}
                      opt1Value={1}
                      opt1Label="Active"
                      opt2Value={2}
                      opt2Label="Archived"
                    />

                    <p class="subtitle is-6">
                      <FontAwesomeIcon className="fas" icon={faIdCard} />
                      &nbsp;Detail
                    </p>
                    <hr />

                    <FormRowText
                      label="Name"
                      value={project.name}
                      helpText=""
                    />

                    <div class="columns pt-5">
                      <div class="column is-half">
                        <Link
                          class="button is-fullwidth-mobile"
                          to={`/admin/projects`}
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back to Projects
                        </Link>
                      </div>
                      <div class="column is-half has-text-right">
                        {project && project.status === 1 && <Link
                          to={`/admin/project/${id}/edit`}
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

export default AdminProjectDetail;
