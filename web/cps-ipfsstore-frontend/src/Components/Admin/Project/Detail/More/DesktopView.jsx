import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faChevronRight,
  faMobile,
  faKey,
  faBuildingProject,
  faImage,
  faPaperclip,
  faAddressCard,
  faSquarePhone,
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
  faIdCard,
  faAddressBook,
  faContactCard,
  faChartPie,
  faBuilding,
  faEllipsis,
  faArchive,
  faBoxOpen,
  faTrashCan,
  faHomeProject,
} from "@fortawesome/free-solid-svg-icons";

import BubbleLink from "../../../../Reusable/EveryPage/BubbleLink";
import {
  USER_ROLE_ROOT
} from "../../../../../Constants/App";


const AdminClientDetailMoreDesktop = ({ id, project, currentProject }) => {
  return (
    <>
      <section className="hero is-hidden-mobile">
        <div className="hero-body has-text-centered">
          <div className="container">
            <div className="columns is-vcentered is-multiline">
              {/*
                <div className="column">
                    <BubbleLink
                        title={`Photo`}
                        subtitle={`Upload a photo of the project`}
                        faIcon={faImage}
                        url={`/admin/project/${id}/avatar`}
                        bgColour={`has-background-danger-dark`}
                    />
                </div>
                */}

              {/* ---------------------------------------------------------------------- */}

              {project.status === 100 ? (
                <div className="column">
                  <BubbleLink
                    title={`Unarchive`}
                    subtitle={`Make project visible in list and search results`}
                    faIcon={faBoxOpen}
                    url={`/admin/project/${id}/more/unarchive`}
                    bgColour={`has-background-success-dark`}
                  />
                </div>
              ) : (
                <div className="column">
                  <BubbleLink
                    title={`Archive`}
                    subtitle={`Make project hidden from list and search results`}
                    faIcon={faArchive}
                    url={`/admin/project/${id}/more/archive`}
                    bgColour={`has-background-success-dark`}
                  />
                </div>
              )}

              {/* ---------------------------------------------------------------------- */}

              {project && project.status === 1 && <>
                  {project.role !== USER_ROLE_ROOT && <div className="column">
                    <BubbleLink
                      title={`Delete`}
                      subtitle={`Permanently delete this project and all associated data`}
                      faIcon={faTrashCan}
                      url={`/admin/project/${id}/more/permadelete`}
                      bgColour={`has-background-danger`}
                    />
                  </div>}
                  <div className="column">
                    <BubbleLink
                      title={`Password`}
                      subtitle={`Change or reset the project\'s password`}
                      faIcon={faKey}
                      url={`/admin/project/${id}/more/change-password`}
                      bgColour={`has-background-danger-dark`}
                    />
                  </div>
                  <div className="column">
                    <BubbleLink
                      title={`2FA`}
                      subtitle={`Enable or disable two-factor authentication`}
                      faIcon={faMobile}
                      url={`/admin/project/${id}/more/change-2fa`}
                      bgColour={`has-background-dark`}
                    />
                  </div>
               </>}

              {/* ---------------------------------------------------------------------- */}
            </div>
          </div>
        </div>
      </section>
    </>
  );
};

export default AdminClientDetailMoreDesktop;
