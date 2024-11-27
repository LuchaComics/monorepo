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
import QRCode from "qrcode.react";

import FormErrorBox from "../../Reusable/FormErrorBox";
import { postRecoveryOTP } from "../../../API/Gateway";
import { currentUserState } from "../../../AppState";
import FormInputField from "../../Reusable/FormInputField";

function TwoFactorAuthenticationBackupCodeRecovery() {
  ////
  //// Global state.
  ////

  const [currentUser, setCurrentUser] = useRecoilState(currentUserState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [forceURL, setForceURL] = useState("");
  const [backupCode, setBackupCode] = useState("");
  const [submittedParamToken, setSubmittedParamToken] = useState(false);

  ////
  //// API.
  ////

  function onRecoveryOPTSuccess(response) {
    console.log("onRecoveryOPTSuccess: Starting...");
    if (response !== undefined && response !== null && response !== "") {
      console.log("response: ", response);
      if (
        response.user !== undefined &&
        response.user !== null &&
        response.user !== ""
      ) {
        console.log("response.user: ", response.user);

        // Save our updated user account.
        setCurrentUser(response.user);

        switch (response.user.role) {
          case 1:
            setForceURL("/admin/dashboard");
            break;
          case 2:
            setForceURL("/dashboard");
            break;
          default:
            setForceURL("/501");
            break;
        }
      }
    }
  }

  function onRecoveryOPTError(apiErr) {
    console.log("onRecoveryOPTError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onRecoveryOPTDone() {
    console.log("onRecoveryOPTDone: Starting...");
  }

  ////
  //// Event handling.
  ////

  function onButtonClick(e) {
    // Remove whitespace characters from token
    const cleanedBackupCode = backupCode.replace(/\s/g, "");

    const payload = {
      backup_code: cleanedBackupCode,
    };
    postRecoveryOTP(
      payload,
      onRecoveryOPTSuccess,
      onRecoveryOPTError,
      onRecoveryOPTDone,
    );
  }

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
                      {/* Start Logo */}
                      <nav className="level">
                        <div className="level-item has-text-centered">
                          <figure className="image">
                            <Link to="/">
                              <img
                                src="/static/COMICCOIN_FAUCET logo 2023 GR.webp"
                                style={{ width: "256px" }}
                              />
                            </Link>
                          </figure>
                        </div>
                      </nav>
                      {/* End Logo */}

                      {/* Page */}
                      <form>
                        <h1 className="title is-2 is-size-4-mobile has-text-centered">
                          Two-Factor Authentication
                        </h1>
                        <h1 className="title is-4 is-size-6-mobile has-text-centered">
                          Backup Code Recovery
                        </h1>
                        <FormErrorBox errors={errors} />
                        <p class="has-text-grey">
                          Copy and paste your <b>2FA backup code</b> into the
                          following field and submit into the system to disable
                          2FA and log into your dashboard.
                        </p>
                        <p>&nbsp;</p>
                        <FormInputField
                          label="Enter your 2FA Backup Code:"
                          name="backupCode"
                          placeholder="Please enter here..."
                          value={backupCode}
                          errorText={errors && errors.backupCode}
                          helpText=""
                          onChange={(e) => setBackupCode(e.target.value)}
                          isRequired={true}
                          maxWidth="380px"
                        />
                        <br />
                      </form>

                      <nav class="level">
                        <div class="level-left">
                          <div class="level-item">
                            <Link
                              class="button is-link is-fullwidth-mobile"
                              to="/login/2fa"
                            >
                              <FontAwesomeIcon icon={faArrowLeft} />
                              &nbsp;Back to 2FA
                            </Link>
                          </div>
                        </div>
                        <div class="level-right">
                          <div class="level-item">
                            <button
                              type="button"
                              class="button is-primary is-fullwidth-mobile"
                              onClick={onButtonClick}
                            >
                              <FontAwesomeIcon icon={faCheckCircle} />
                              &nbsp;Submit
                            </button>
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

export default TwoFactorAuthenticationBackupCodeRecovery;
