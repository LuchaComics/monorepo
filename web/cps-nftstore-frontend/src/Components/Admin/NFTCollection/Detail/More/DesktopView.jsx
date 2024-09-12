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
  faRocket
} from "@fortawesome/free-solid-svg-icons";

import BubbleLink from "../../../../Reusable/EveryPage/BubbleLink";
import {
  USER_ROLE_ROOT
} from "../../../../../Constants/App";


const AdminClientDetailMoreDesktop = ({ collection }) => {
  return (
    <>
      <section className="hero is-hidden-mobile">
        <div className="hero-body has-text-centered">
          <div className="container">
            <div className="columns is-vcentered is-multiline">
              {/* ---------------------------------------------------------------------- */}
              {collection.smartContractStatus === 1 && (
                <div className="column">
                  <BubbleLink
                    title={`Deploy`}
                    subtitle={`Deploy smart contract to blockchain.`}
                    faIcon={faRocket}
                    url={`/admin/collection/${collection.id}/more/deploy`}
                    bgColour={`has-background-danger-dark`}
                  />
                </div>
              )}
              {/* ---------------------------------------------------------------------- */}
            </div>
          </div>
        </div>
      </section>
    </>
  );
};

export default AdminClientDetailMoreDesktop;
