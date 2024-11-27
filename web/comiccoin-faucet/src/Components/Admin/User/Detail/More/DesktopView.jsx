import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faChevronRight,
  faMobile,
  faKey,
  faBuildingUser,
  faImage,
  faPaperclip,
  faAddressCard,
  faSquarePhone,
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
  faBuilding,
  faEllipsis,
  faArchive,
  faBoxOpen,
  faTrashCan,
  faHomeUser,
} from "@fortawesome/free-solid-svg-icons";

import BubbleLink from "../../../../Reusable/EveryPage/BubbleLink";
import {
  USER_ROLE_ROOT
} from "../../../../../Constants/App";


const AdminClientDetailMoreDesktop = ({ id, user, currentUser }) => {
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
                        subtitle={`Upload a photo of the user`}
                        faIcon={faImage}
                        url={`/admin/user/${id}/avatar`}
                        bgColour={`has-background-danger-dark`}
                    />
                </div>
                */}

              {/* ---------------------------------------------------------------------- */}

              {user.status === 100 ? (
                <div className="column">
                  <BubbleLink
                    title={`Unarchive`}
                    subtitle={`Make user visible in list and search results`}
                    faIcon={faBoxOpen}
                    url={`/admin/user/${id}/more/unarchive`}
                    bgColour={`has-background-success-dark`}
                  />
                </div>
              ) : (
                <div className="column">
                  <BubbleLink
                    title={`Archive`}
                    subtitle={`Make user hidden from list and search results`}
                    faIcon={faArchive}
                    url={`/admin/user/${id}/more/archive`}
                    bgColour={`has-background-success-dark`}
                  />
                </div>
              )}

              {/* ---------------------------------------------------------------------- */}

              {user && user.status === 1 && <>
                  {user.role !== USER_ROLE_ROOT && <div className="column">
                    <BubbleLink
                      title={`Delete`}
                      subtitle={`Permanently delete this user and all associated data`}
                      faIcon={faTrashCan}
                      url={`/admin/user/${id}/more/permadelete`}
                      bgColour={`has-background-danger`}
                    />
                  </div>}
                  <div className="column">
                    <BubbleLink
                      title={`Password`}
                      subtitle={`Change or reset the user\'s password`}
                      faIcon={faKey}
                      url={`/admin/user/${id}/more/change-password`}
                      bgColour={`has-background-danger-dark`}
                    />
                  </div>
                  <div className="column">
                    <BubbleLink
                      title={`2FA`}
                      subtitle={`Enable or disable two-factor authentication`}
                      faIcon={faMobile}
                      url={`/admin/user/${id}/more/change-2fa`}
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
