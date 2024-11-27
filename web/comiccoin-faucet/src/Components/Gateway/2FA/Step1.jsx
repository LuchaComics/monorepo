import React, { useState, useEffect } from "react";
import { Link, Navigate, useSearchParams } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faCheckCircle,
  faArrowLeft,
  faArrowRight,
  faEnvelope,
  faKey,
  faTriangleExclamation,
  faArrowUpRightFromSquare,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import FormErrorBox from "../../Reusable/FormErrorBox";
import useLocalStorage from "../../../Hooks/useLocalStorage";
import { postLoginAPI } from "../../../API/Gateway";
import { onHamburgerClickedState, currentUserState } from "../../../AppState";
import {
  EXECUTIVE_ROLE_ID,
  MANAGEMENT_ROLE_ID,
  FRONTLINE_ROLE_ID,
  ASSOCIATE_ROLE_ID,
  CUSTOMER_ROLE_ID,
} from "../../../Constants/App";

function TwoFactorAuthenticationWizardStep1() {
  ////
  //// Global state.
  ////

  const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [forceURL, setForceURL] = useState("");

  ////
  //// API.
  ////

  ////
  //// Event handling.
  ////

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
    }

    return () => (mounted = false);
  }, []);

  ////
  //// Component rendering.
  ////

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <>
      <div className="column is-12 container">
        <div className="section">
          <section className="hero is-fullheight">
            <div className="hero-body">
              <div className="container">
                <div className="columns is-centered">
                  <div className="column is-half-widescreen">
                    <div className="box is-rounded">
                      {/* Progress Wizard */}
                      <nav className="box has-background-light">
                        <p className="subtitle is-5">Step 1 of 3</p>
                        <progress
                          class="progress is-success"
                          value="33"
                          max="100"
                        >
                          43%
                        </progress>
                      </nav>

                      {/* Page */}
                      <form className="">
                        <h1 className="title is-2 is-size-4-mobile has-text-centered">
                          Setup Two-Factor Authentication
                        </h1>
                        <FormErrorBox errors={errors} />
                        <div className="content">
                          <p class="has-text-grey">
                            To ensure your account stays secure, you need to
                            sign in using <i>two-factor Authentication (2FA)</i>
                            . The following wizard will help you get setup with
                            2FA.
                          </p>
                          <p class="has-text-grey">
                            To make initial 2FA setup easier, we encourage you
                            to login on a device BESIDES the mobile device with
                            the camera that you wish to use. We recommend the
                            following setup:
                            <ul>
                              <li>Login on a desktop device</li>
                              <li>
                                Use your mobile phone to scan the QR code and
                                complete setup
                              </li>
                            </ul>
                          </p>
                          <p class="has-text-grey">
                            To begin, please download any of the following
                            applications for your mobile device.
                          </p>
                          {/* Apple 2FA Authenticator */}
                          <div className="card">
                            <div className="card-content">
                              <div className="media">
                                <div className="media-content">
                                  <p className="title is-5">
                                    <u>Apple 2FA</u>
                                  </p>
                                </div>
                              </div>

                              <div className="content">
                                <p>
                                  All iOS and Mac devices with a{" "}
                                  <b>Safari Web Browser</b> come with build in a
                                  2FA verification services. Sign in with your{" "}
                                  <i>Apple ID</i> in Safari and you can take
                                  advantage of this service.
                                </p>
                              </div>
                            </div>
                            {/*}
                              <footer className="card-footer">
                                  <button className="card-footer-item button is-primary is-small">
                                      Download&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} />
                                  </button>
                              </footer>
                              */}
                          </div>

                          {/* Google Authenticator */}
                          <div className="card">
                            <div className="card-content">
                              <div className="media">
                                <div className="media-content">
                                  <p className="title is-5">
                                    <u>Google Authenticator</u>
                                  </p>
                                </div>
                              </div>

                              <div className="content">
                                <p>
                                  This 2FA app is created by <b>Google, Inc.</b>
                                </p>
                                <b>Download for iOS:</b>&nbsp;
                                <Link
                                  className=""
                                  to="https://apps.apple.com/ca/app/google-authenticator/id388497605"
                                  target="_blank"
                                  rel="noreferrer"
                                >
                                  Visit App Store&nbsp;
                                  <FontAwesomeIcon
                                    className="fas"
                                    icon={faArrowUpRightFromSquare}
                                  />
                                </Link>
                                <br />
                                <b>Download for Android:</b>&nbsp;
                                <Link
                                  className=""
                                  to="https://play.google.com/store/apps/details?id=com.google.android.apps.authenticator2&pli=1"
                                  target="_blank"
                                  rel="noreferrer"
                                >
                                  Visit Google Play&nbsp;
                                  <FontAwesomeIcon
                                    className="fas"
                                    icon={faArrowUpRightFromSquare}
                                  />
                                </Link>
                                <br />
                              </div>
                            </div>
                            {/*}
                                                    <footer className="card-footer">
                                                        <button className="card-footer-item button is-primary is-small">
                                                           Download&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} />
                                                        </button>
                                                    </footer>
                                                    */}
                          </div>

                          {/* Authenticator Chrome Extension */}
                          <div className="card">
                            <div className="card-content">
                              <div className="media">
                                <div className="media-content">
                                  <p className="title is-5">
                                    <u>Authenticator</u>
                                  </p>
                                </div>
                              </div>

                              <div className="content">
                                <p>
                                  This 2FA app is created by{" "}
                                  <b>authenticator.cc</b>
                                </p>
                                <b>Download for Chrome:</b>&nbsp;
                                <Link
                                  className=""
                                  to="https://chromewebstore.google.com/detail/authenticator/bhghoamapcdpbohphigoooaddinpkbai?pli=1"
                                  target="_blank"
                                  rel="noreferrer"
                                >
                                  Visit Chrome web store&nbsp;
                                  <FontAwesomeIcon
                                    className="fas"
                                    icon={faArrowUpRightFromSquare}
                                  />
                                </Link>
                                <br />
                              </div>
                            </div>
                            {/*}
                                                    <footer className="card-footer">
                                                        <button className="card-footer-item button is-primary is-small">
                                                           Download&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} />
                                                        </button>
                                                    </footer>
                                                    */}
                          </div>
                        </div>
                      </form>
                      <br />

                      <nav class="level">
                        <div class="level-left">
                          <div class="level-item">
                            <Link
                              class="button is-link is-fullwidth-mobile"
                              to="/login"
                            >
                              <FontAwesomeIcon icon={faArrowLeft} />
                              &nbsp;Cancel
                            </Link>
                          </div>
                        </div>
                        <div class="level-right">
                          <div class="level-item">
                            <Link
                              type="button"
                              class="button is-primary is-fullwidth-mobile"
                              to="/login/2fa/step-2"
                            >
                              Next&nbsp;
                              <FontAwesomeIcon icon={faArrowRight} />
                            </Link>
                          </div>
                        </div>
                      </nav>
                    </div>
                    {/* End box */}

                    <div className="has-text-centered">
                      <p>© 2024 COMICCOIN_FAUCET</p>
                    </div>
                    {/* End suppoert text. */}
                  </div>
                  {/* End Column */}
                </div>
              </div>
              {/* End container */}
            </div>
            {/* End hero-body */}
          </section>
        </div>
      </div>
    </>
  );
}

export default TwoFactorAuthenticationWizardStep1;
