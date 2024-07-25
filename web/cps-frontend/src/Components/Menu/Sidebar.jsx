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
  faBuilding,
  faBarcode,
  faQuestionCircle
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { onHamburgerClickedState, currentUserState } from "../../AppState";
import { USER_ROLE_ROOT, USER_ROLE_RETAILER, USER_ROLE_CUSTOMER } from "../../Constants/App";

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
    "/register",
    "/register/user",
    "/register/store",
    "/register-successful",
    "/index",
    "/login",
    "/login/2fa",
    "/login/2fa/step-1",
    "/login/2fa/step-2",
    "/login/2fa/step-3",
    "/login/2fa/step-3/backup-code",
    "/login/2fa/backup-code",
    "/login/2fa/backup-code-recovery",
    "/logout",
    "/verify",
    "/forgot-password",
    "/password-reset",
    "/cpsrn-result",
    "/cpsrn-registry",
    "/terms",
    "/privacy",
    "/cpsrn"
  ];
  const location = useLocation();
  var arrayLength = ignorePathsArr.length;
  for (var i = 0; i < arrayLength; i++) {
    // console.log(location.pathname, "===", ignorePathsArr[i], " EQUALS ", location.pathname === ignorePathsArr[i]);
    if (location.pathname === ignorePathsArr[i]) {
      return null;
    }
  }

  //-------------//
  // CASE 2 OF 3 //
  //-------------//

  if (currentUser === null) {
    return null;
  }

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
      {currentUser.role === USER_ROLE_ROOT && (
        <div
          className={`column is-one-fifth has-background-black ${onHamburgerClicked ? "" : "is-hidden"}`}
        >
          <nav class="level is-hidden-mobile">
            <div class="level-item has-text-centered">
              <figure class="image">
                <img
                  src="/static/CPS logo 2023 GR.webp"
                  style={{ maxWidth: "200px" }}
                />
              </figure>
            </div>
          </nav>
          <aside class="menu p-4">
            <p class="menu-label has-text-grey-light">System Administrator</p>
            <ul class="menu-list">
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  to="/admin/dashboard"
                  class={`has-text-grey-light ${location.pathname.includes("dashboard") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faTachometer} />
                  &nbsp;Dashboard
                </Link>
              </li>
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  to="/admin/stores"
                  class={`has-text-grey-light ${location.pathname.includes("store") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faBuilding} />
                  &nbsp;Stores
                </Link>
              </li>
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  to="/admin/users"
                  class={`has-text-grey-light ${location.pathname.includes("user") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faUsers} />
                  &nbsp;All Users
                </Link>
              </li>
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  class={`has-text-grey-light`}
                >
                  <FontAwesomeIcon className="fas" icon={faTasks} />
                  &nbsp;Online Submissions
                </Link>
                <ul>
                  <li>
                    <Link
                      onClick={onLinkClickCloseHamburgerMenuIfMobile}
                      to="/admin/submissions/comics"
                      class={`has-text-grey-light ${location.pathname.includes("submissions/comic") && "is-active"}`}
                    >
                      <FontAwesomeIcon className="fas" icon={faBook} />
                      &nbsp;Comics
                    </Link>
                  </li>
                  {/*
                                    <li>
                                        <Link to="/admin/submissions/cards" class={`has-text-grey-light ${location.pathname.includes("card") && "is-active"}`}>
                                            <FontAwesomeIcon className="fas" icon={faTachometer} />&nbsp;Cards
                                        </Link>
                                    </li>
                                    */}
                </ul>
              </li>
            </ul>

            <p class="menu-label has-text-grey-light">System</p>
            <ul class="menu-list">
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  to="/admin/registry"
                  class={`has-text-grey-light ${location.pathname.includes("registry") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faBarcode} />
                  &nbsp;Registry
                </Link>
              </li>
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  to={`/admin/offers`}
                  class={`has-text-grey-light ${location.pathname.includes("offer") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faHandHolding} />
                  &nbsp;Offers
                </Link>
              </li>
            </ul>

            <p class="menu-label has-text-grey-light">Account</p>
            <ul class="menu-list">
              <li>
                <Link
                  to={`/account`}
                  class={`has-text-grey-light ${location.pathname.includes("account") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faUserCircle} />
                  &nbsp;Account
                </Link>
              </li>
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  onClick={(e) => setShowLogoutWarning(true)}
                  class={`has-text-grey-light ${location.pathname.includes("logout") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faSignOut} />
                  &nbsp;Sign Off
                </Link>
              </li>
            </ul>
          </aside>
        </div>
      )}
      {/*
        --------
        RETAILER
        --------
      */}
      {currentUser.role === USER_ROLE_RETAILER && (
        <div
          className={`column is-one-fifth has-background-black ${onHamburgerClicked ? "" : "is-hidden"}`}
        >
          <nav class="level is-hidden-mobile">
            <div class="level-item has-text-centered">
              <figure class="image">
                <img
                  src="/static/CPS logo 2023 GR.webp"
                  style={{ maxWidth: "200px" }}
                />
              </figure>
            </div>
          </nav>
          <aside class="menu p-4">
            <p class="menu-label has-text-grey-light">Store Owner/Manager</p>
            <ul class="menu-list">
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  to="/dashboard"
                  class={`has-text-grey-light ${location.pathname.includes("dashboard") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faTachometer} />
                  &nbsp;Dashboard
                </Link>
              </li>

              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  to="/customers"
                  class={`has-text-grey-light ${location.pathname.includes("customer") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faUsers} />
                  &nbsp;Customers
                </Link>
              </li>
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  to="/submissions"
                  class={`has-text-grey-light ${location.pathname.includes("submissions") && !location.pathname.includes("comic") && !location.pathname.includes("card") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faTasks} />
                  &nbsp;Online Submissions
                </Link>
                <ul>
                  <li>
                    <Link
                      onClick={onLinkClickCloseHamburgerMenuIfMobile}
                      to="/submissions/comics"
                      class={`has-text-grey-light ${location.pathname.includes("submissions/comic") && "is-active"}`}
                    >
                      <FontAwesomeIcon className="fas" icon={faBook} />
                      &nbsp;Comics
                    </Link>
                  </li>
                  {/*
                                    <li>
                                        <Link to="/submissions/cards" class={`has-text-grey-light ${location.pathname.includes("card") && "is-active"}`}>
                                            <FontAwesomeIcon className="fas" icon={faTachometer} />&nbsp;Cards
                                        </Link>
                                    </li>
                                    */}
                </ul>
              </li>
            </ul>

            <p class="menu-label has-text-grey-light">System</p>
            <ul class="menu-list">
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  to="/registry"
                  class={`has-text-grey-light ${location.pathname.includes("registry") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faBarcode} />
                  &nbsp;Registry
                </Link>
              </li>
            </ul>

            <p class="menu-label has-text-grey-light">Account</p>
            <ul class="menu-list">
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  to={`/account`}
                  class={`has-text-grey-light ${location.pathname.includes("account") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faUserCircle} />
                  &nbsp;Account
                </Link>
              </li>
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  to={`/store`}
                  class={`has-text-grey-light ${location.pathname.includes("store") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faBuilding} />
                  &nbsp;My Store
                </Link>
              </li>
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  onClick={(e) => setShowLogoutWarning(true)}
                  class={`has-text-grey-light ${location.pathname.includes("logout") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faSignOut} />
                  &nbsp;Sign Off
                </Link>
              </li>
            </ul>
          </aside>
        </div>
      )}
      {/*
        --------
        CUSTOMER
        --------
      */}
      {currentUser.role === USER_ROLE_CUSTOMER && (
        <div
          className={`column is-one-fifth has-background-black ${onHamburgerClicked ? "" : "is-hidden"}`}
        >
          <nav class="level is-hidden-mobile">
            <div class="level-item has-text-centered">
              <figure class="image">
                <img
                  src="/static/CPS logo 2023 GR.webp"
                  style={{ maxWidth: "200px" }}
                />
              </figure>
            </div>
          </nav>
          <aside class="menu p-4">
            <p class="menu-label has-text-grey-light">Regular User</p>
            <ul class="menu-list">
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  to="/c/dashboard"
                  class={`has-text-grey-light ${location.pathname.includes("dashboard") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faTachometer} />
                  &nbsp;Dashboard
                </Link>
              </li>
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  to="/c/submissions"
                  class={`has-text-grey-light ${location.pathname.includes("submissions") && !location.pathname.includes("comic") && !location.pathname.includes("card") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faTasks} />
                  &nbsp;My Submissions
                </Link>
                <ul>
                  <li>
                    <Link
                      onClick={onLinkClickCloseHamburgerMenuIfMobile}
                      to="/c/submissions/comics"
                      class={`has-text-grey-light ${location.pathname.includes("submissions/comic") && "is-active"}`}
                    >
                      <FontAwesomeIcon className="fas" icon={faBook} />
                      &nbsp;Comics
                    </Link>
                  </li>
                  {/*
                    <li>
                        <Link to="/submissions/cards" class={`has-text-grey-light ${location.pathname.includes("card") && "is-active"}`}>
                            <FontAwesomeIcon className="fas" icon={faTachometer} />&nbsp;Cards
                        </Link>
                    </li>
                    */}
                </ul>
              </li>
            </ul>

            <p class="menu-label has-text-grey-light">System</p>
            <ul class="menu-list">
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  to="/c/registry"
                  class={`has-text-grey-light ${location.pathname.includes("registry") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faBarcode} />
                  &nbsp;Registry
                </Link>
              </li>
            </ul>

            <p class="menu-label has-text-grey-light">Account</p>
            <ul class="menu-list">
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  to={`/account`}
                  class={`has-text-grey-light ${location.pathname.includes("account") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faUserCircle} />
                  &nbsp;Account
                </Link>
              </li>
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  to={`/help`}
                  class={`has-text-grey-light ${location.pathname.includes("help") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faQuestionCircle} />
                  &nbsp;Help
                </Link>
              </li>
              <li>
                <Link
                  onClick={onLinkClickCloseHamburgerMenuIfMobile}
                  onClick={(e) => setShowLogoutWarning(true)}
                  class={`has-text-grey-light ${location.pathname.includes("logout") && "is-active"}`}
                >
                  <FontAwesomeIcon className="fas" icon={faSignOut} />
                  &nbsp;Sign Off
                </Link>
              </li>
            </ul>
          </aside>
        </div>
      )}
    </>
  );
};
