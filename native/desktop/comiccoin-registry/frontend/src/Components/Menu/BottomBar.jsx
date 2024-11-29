import React, { useState, useEffect } from "react";
import { Link, useLocation } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faBars,
  faRightFromBracket,
  faTachometer,
  faTasks,
  faSignOut,
  faUserCircle,
  faUsers,
  faBuilding,
  faGauge,
  faPaperPlane,
  faEllipsis,
  faInbox,
  faCubes
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import logo from '../../assets/images/CPS-logo-2023-square.webp';
import { onHamburgerClickedState, currentUserState } from "../../AppState";
import { USER_ROLE_ROOT, USER_ROLE_RETAILER, USER_ROLE_CUSTOMER } from "../../Constants/App";

function BottomTabBar() {
  ////
  //// Global State
  ////



  ////
  //// Local State
  ////

  ////
  //// Events
  ////

  // Do nothing.

  ////
  //// Rendering.
  ////

  //-------------//
  // CASE 1 OF 3 //
  //-------------//

  // Get the current location and if we are at specific URL paths then we
  // will not render this component.
  const ignorePathsArr = [
    "/",
    "/pick-data-directory",
    "/config-nftstore",
    "/startup",
    "/wallets",
    "/wallet/add",
  ];
  const location = useLocation();
  var arrayLength = ignorePathsArr.length;
  for (var i = 0; i < arrayLength; i++) {
    // console.log(location.pathname, "===", ignorePathsArr[i], " EQUALS ", location.pathname === ignorePathsArr[i]);
    if (location.pathname === ignorePathsArr[i]) {
      return null;
    }
  }

  // //-------------//
  // // CASE 2 OF 3 //
  // //-------------//
  //
  // if (currentUser === null) {
  //   return null;
  // }

  //-------------//
  // CASE 3 OF 3 //
  //-------------//

  // Render the following component GUI
  return (
    <nav className="tabs is-bottom is-fullwidth">
        <ul>
            <li className={`has-text-grey-light ${location.pathname.includes("dashboard") && "is-active"}`}>
              <Link to="/dashboard">
                <span className="icon is-small">
                  <FontAwesomeIcon icon={faGauge} />
                </span>
                <span>Dashboard</span>
              </Link>
            </li>
            <li className={`has-text-grey-light ${location.pathname.includes("token") && "is-active"}`}>
              <Link to="/tokens">
                <span className="icon is-small">
                  <FontAwesomeIcon icon={faCubes} />
                </span>
                <span>Tokens</span>
              </Link>
            </li>
            {/*
            <li className={`has-text-grey-light ${location.pathname.includes("send") && "is-active"}`}>
              <Link to="/send">
                <span className="icon is-small">
                  <FontAwesomeIcon icon={faPaperPlane} />
                </span>
                <span>Send</span>
              </Link>
            </li>
            <li className={`has-text-grey-light ${location.pathname.includes("receive") && "is-active"}`}>
              <Link to="/receive">
                <span className="icon is-small">
                  <FontAwesomeIcon icon={faInbox} />
                </span>
                <span>Receive</span>
              </Link>
            </li>
            <li className={`has-text-grey-light ${!location.pathname.includes("dashboard") && !location.pathname.includes("send") && !location.pathname.includes("receive") && "is-active"}`}>
              <Link to="/more">
                <span className="icon is-small">
                  <FontAwesomeIcon icon={faEllipsis} />
                </span>
                <span>More</span>
              </Link>
            </li>
            */}
        </ul>
    </nav>
  );
}

export default BottomTabBar;
