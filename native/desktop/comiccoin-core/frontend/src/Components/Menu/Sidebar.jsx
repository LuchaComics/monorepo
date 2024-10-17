import React, { useState, useEffect } from "react";
import { Link, useLocation } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faHandHolding,
  faBars,
  faBook,
  faRightFromBracket,
  faTachometer,
  faTasks,
  faSignOut,
  faUserCircle,
  faUsers,
  faCubes,
  faBuilding,
  faBarcode,
  faQuestionCircle,
  faCloud,
  faPlus,
  faCogs,
  faInbox,
  faPaperPlane,
  faFileInvoiceDollar,
  faSignIn
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import logo from '../../assets/images/CPS-logo-2023-GR.webp';
import { onHamburgerClickedState, currentUserState } from "../../AppState";
import { USER_ROLE_ROOT } from "../../Constants/App";

export default (props) => {
  ////
  //// Global State
  ////

  const [onHamburgerClicked, setOnHamburgerClicked] = useRecoilState(
    onHamburgerClickedState,
  );
  const [currentUser] = useRecoilState(currentUserState);

  ////
  //// Local State
  ////

  const [showLogoutWarning, setShowLogoutWarning] = useState(false);

  ////
  //// Events
  ////

  // onLinkClick function will check to see if we are on a mobile device and if we are then we will close the hanburger menu.
  const onLinkClickCloseHamburgerMenuIfMobile = (e) => {
    // Special thanks to: https://dev.to/timhuang/a-simple-way-to-detect-if-browser-is-on-a-mobile-device-with-javascript-44j3
    if (
      /Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(
        navigator.userAgent,
      )
    ) {
      // document.write("mobile");
      setOnHamburgerClicked(false);
    } else {
      // document.write("not mobile");
    }
  };

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
    "/pick-storage-location-on-startup",
  ];
  const location = useLocation();
  var arrayLength = ignorePathsArr.length;
  for (var i = 0; i < arrayLength; i++) {
    // console.log(location.pathname, "===", ignorePathsArr[i], " EQUALS ", location.pathname === ignorePathsArr[i]);
    if (location.pathname === ignorePathsArr[i]) {
      return null;
    }
  }
  //
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

  return (
    <>
      <div class={`modal ${showLogoutWarning ? "is-active" : ""}`}>
        <div class="modal-background"></div>
        <div class="modal-card">
          <header class="modal-card-head">
            <p class="modal-card-title">Are you sure?</p>
            <button
              class="delete"
              aria-label="close"
              onClick={(e) => setShowLogoutWarning(false)}
            ></button>
          </header>
          <section class="modal-card-body">
            You are about to log out of the system and you'll need to log in
            again next time. Are you sure you want to continue?
          </section>
          <footer class="modal-card-foot">
            <Link class="button is-success" to={`/logout`}>
              Yes
            </Link>
            <button class="button" onClick={(e) => setShowLogoutWarning(false)}>
              No
            </button>
          </footer>
        </div>
      </div>
      {/*
        -----
        STAFF
        -----
      */}
        <div
          className={`column is-one-fifth has-background-black ${onHamburgerClicked ? "" : "is-hidden"}`}
        >
          <nav class="level is-hidden-mobile">
            <div class="level-item has-text-centered">
              <figure class="image">
                <img
                  src={logo}
                  style={{ maxWidth: "200px" }}
                />
              </figure>
            </div>
          </nav>
          <aside class="menu p-4">
            <p class="menu-label has-text-grey-light">Menu</p>
            <ul class="menu-list">
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  to="/dashboard"
                  class={`has-text-grey-light ${location.pathname.includes("overview") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faTachometer} />
                  &nbsp;Overview
                </Link>
              </li>
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  to="/send"
                  class={`has-text-grey-light ${location.pathname.includes("tenant") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faPaperPlane} />
                  &nbsp;Send
                </Link>
              </li>
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  to="/receive"
                  class={`has-text-grey-light ${location.pathname.includes("user") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faInbox} />
                  &nbsp;Receive
                </Link>
              </li>
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  to="/admin/collections"
                  class={`has-text-grey-light ${location.pathname.includes("collection") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faCubes} />
                  &nbsp;Tokens
                </Link>
              </li>

              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  to="/transactions"
                  class={`has-text-grey-light ${location.pathname.includes("ipfs") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faFileInvoiceDollar} />
                  &nbsp;Transactions
                </Link>
              </li>
            </ul>
            <p class="menu-label has-text-grey-light">System</p>
            <ul class="menu-list">
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  onClick={(e) => setShowLogoutWarning(true)}
                  class={`has-text-grey-light ${location.pathname.includes("logout") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faSignIn} />
                  &nbsp;Create Wallet
                </Link>
              </li>
              <li>
                <Link
                  to={`/settings`}
                  class={`has-text-grey-light ${location.pathname.includes("settings") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faCogs} />
                  &nbsp;Settings
                </Link>
              </li>
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  onClick={(e) => setShowLogoutWarning(true)}
                  class={`has-text-grey-light ${location.pathname.includes("logout") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faSignOut} />
                  &nbsp;Close Wallet
                </Link>
              </li>
            </ul>
          </aside>
        </div>
    </>
  );
};
