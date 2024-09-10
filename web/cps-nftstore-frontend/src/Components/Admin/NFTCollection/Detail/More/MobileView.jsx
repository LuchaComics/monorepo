import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faChevronRight,
  faMobile,
  faKey,
  faBuildingCollection,
  faImage,
  faPaperclip,
  faAddressCard,
  faSquarePhone,
  faTasks,
  faTachometer,
  faPlus,
  faArrowLeft,
  faCheckCircle,
  faCollectionCircle,
  faGauge,
  faPencil,
  faCollections,
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
  faHomeCollection,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";
import {
  USER_ROLE_ROOT
} from "../../../../../Constants/App";


const AdminClientDetailMoreMobile = ({ id, collection, currentCollection }) => {
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
             Coming soon...
            {/* ---------------------------------------------------------------------- */}

          </tbody>
        </table>
      </div>
      {/* END Page Menu Options (Mobile Only) */}
    </>
  );
};

export default AdminClientDetailMoreMobile;
