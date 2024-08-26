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
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";
import {
  USER_ROLE_ROOT
} from "../../../../../Constants/App";


const AdminClientDetailMoreMobile = ({ id, project, currentProject }) => {
  return (
    <>
      <div
        className="has-background-white-ter is-hidden-tablet mb-6 p-5"
        style={{ borderRadius: "15px" }}
      >
        <table className="is-fullwidth has-background-white-ter table">
          <thead>
            <tr>
              <th colSpan="2">Menu</th>
            </tr>
          </thead>
          <tbody>
            {/* ---------------------------------------------------------------------- */}

            {project && project.status === 100 ? (
              <tr>
                <td>
                  <FontAwesomeIcon className="fas" icon={faBoxOpen} />
                  &nbsp;Unarchive
                </td>
                <td>
                  <div className="buttons is-right">
                    <Link
                      to={`/admin/project/${id}/more/unarchive`}
                      className="is-small"
                    >
                      View&nbsp;
                      <FontAwesomeIcon className="mdi" icon={faChevronRight} />
                    </Link>
                  </div>
                </td>
              </tr>
            ) : (
              <tr>
                <td>
                  <FontAwesomeIcon className="fas" icon={faArchive} />
                  &nbsp;Archive
                </td>
                <td>
                  <div className="buttons is-right">
                    <Link
                      to={`/admin/project/${id}/more/archive`}
                      className="is-small"
                    >
                      View&nbsp;
                      <FontAwesomeIcon className="mdi" icon={faChevronRight} />
                    </Link>
                  </div>
                </td>
              </tr>
            )}

            {/* ---------------------------------------------------------------------- */}

            {project.role !== USER_ROLE_ROOT && <tr>
              <td>
                <FontAwesomeIcon className="fas" icon={faTrashCan} />
                &nbsp;Delete
              </td>
              <td>
                <div className="buttons is-right">
                  <Link
                    to={`/admin/project/${id}/more/permadelete`}
                    className="is-small"
                  >
                    View&nbsp;
                    <FontAwesomeIcon
                      className="mdi"
                      icon={faChevronRight}
                    />
                  </Link>
                </div>
              </td>
            </tr>}

            {/* ---------------------------------------------------------------------- */}
            <tr>
              <td>
                <FontAwesomeIcon className="fas" icon={faKey} />
                &nbsp;Password
              </td>
              <td>
                <div className="buttons is-right">
                  <Link
                    to={`/admin/project/${id}/more/change-password`}
                    className="is-small"
                  >
                    View&nbsp;
                    <FontAwesomeIcon
                      className="mdi"
                      icon={faChevronRight}
                    />
                  </Link>
                </div>
              </td>
            </tr>
            {/* ---------------------------------------------------------------------- */}
            <tr>
              <td>
                <FontAwesomeIcon className="fas" icon={faMobile} />
                &nbsp;2FA
              </td>
              <td>
                <div className="buttons is-right">
                  <Link
                    to={`/admin/project/${id}/more/change-2fa`}
                    className="is-small"
                  >
                    View&nbsp;
                    <FontAwesomeIcon
                      className="mdi"
                      icon={faChevronRight}
                    />
                  </Link>
                </div>
              </td>
            </tr>
            {/* ---------------------------------------------------------------------- */}
          </tbody>
        </table>
      </div>
      {/* END Page Menu Options (Mobile Only) */}
    </>
  );
};

export default AdminClientDetailMoreMobile;
