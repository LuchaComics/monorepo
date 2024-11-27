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
import { getAccountDetailAPI } from "../../../API/Profile";
import { currentOTPResponseState, currentUserState } from "../../../AppState";
import FormTextareaField from "../../Reusable/FormTextareaField";

function TwoFactorAuthenticationBackupCodeGenerate() {
  ////
  //// URL Parameters.
  ////

  const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140
  const backupCode = searchParams.get("v");

  ////
  //// Global state.
  ////

  const [currentUser] = useRecoilState(currentUserState);

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

  const onSubmitClick = () => {
    switch (currentUser.role) {
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
  };

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
                      {/* Notification */}
                      <article class="message is-primary">
                        <div class="message-body">
                          <FontAwesomeIcon icon={faCheckCircle} />
                          &nbsp;Successfull 2FA Verification
                        </div>
                      </article>

                      {/* Page */}
                      <div className="content">
                        <h1 className="title is-2 is-size-4-mobile has-text-centered">
                          2FA Backup Code
                        </h1>

                        <FormErrorBox errors={errors} />
                        <p class="has-text-grey">
                          You have successfully verified your 2FA code and now
                          are granted backup code which you can use in case you
                          lose your phone or experience data loss.
                        </p>
                        <p class="has-text-grey">
                          Please save this backup code in safe location. When
                          you have successfully saved this code, please
                          continue.
                        </p>
                        <FormTextareaField
                          label="Backup Code"
                          name="backupCode"
                          placeholder="See your authenticator app"
                          value={backupCode}
                          errorText={errors && errors.backupCode}
                          onChange={null}
                          isRequired={true}
                        />
                      </div>

                      <nav class="level">
                        <div class="level-left">
                          <div class="level-item">{/* Add button here */}</div>
                        </div>
                        <div class="level-right">
                          <div class="level-item">
                            <button
                              type="button"
                              class="button is-primary is-fullwidth-mobile"
                              onClick={onSubmitClick}
                            >
                              <FontAwesomeIcon icon={faCheckCircle} />
                              &nbsp;Confirm
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

export default TwoFactorAuthenticationBackupCodeGenerate;
