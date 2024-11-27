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
import FormInputField from "../../Reusable/FormInputField";
import FormTextareaField from "../../Reusable/FormTextareaField";
import {
  postGenerateOTP,
  postGenerateOTPAndQRCodeImage,
} from "../../../API/Gateway";
import { currentOTPResponseState, currentUserState } from "../../../AppState";
import {
  EXECUTIVE_ROLE_ID,
  MANAGEMENT_ROLE_ID,
  FRONTLINE_ROLE_ID,
  ASSOCIATE_ROLE_ID,
  CUSTOMER_ROLE_ID,
} from "../../../Constants/App";

function TwoFactorAuthenticationWizardStep2() {
  ////
  //// Global state.
  ////

  const [currentUser, setCurrentUser] = useRecoilState(currentUserState);
  const [otpResponse, setOtpResponse] = useRecoilState(currentOTPResponseState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [forceURL, setForceURL] = useState("");

  ////
  //// API.
  ////

  function onGenerateOPTSuccess(response) {
    console.log("onGenerateOPTSuccess: Starting...");
    console.log("response: ", response);
    setOtpResponse(response);
  }

  function onGenerateOPTError(apiErr) {
    console.log("onGenerateOPTError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onGenerateOPTDone() {
    console.log("onGenerateOPTDone: Starting...");
  }

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

      if (
        otpResponse === undefined ||
        otpResponse === null ||
        otpResponse === ""
      ) {
        postGenerateOTP(
          onGenerateOPTSuccess,
          onGenerateOPTError,
          onGenerateOPTDone,
        );
      }
    }

    return () => (mounted = false);
  }, [otpResponse]);

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
                        <p className="subtitle is-5">Step 2 of 3</p>
                        <progress
                          class="progress is-success"
                          value="66"
                          max="100"
                        >
                          66%
                        </progress>
                      </nav>

                      {/* Page */}
                      <form>
                        <h1 className="title is-size-2 is-size-4-mobile  has-text-centered">
                          Setup Two-Factor Authentication
                        </h1>
                        <FormErrorBox errors={errors} />
                        <p class="has-text-grey">
                          With your 2FA application open, please scan the
                          following QR code with your device and click next when
                          ready.
                        </p>
                        <p>&nbsp;</p>
                        <div className="columns is-centered has-text-centered">
                          <div class="column is-half-tablet">
                            <figure class="image">
                              {otpResponse && (
                                <QRCode
                                  value={otpResponse.optAuthURL}
                                  size={250}
                                  className=""
                                />
                              )}
                              <br />
                              <span className="has-text-centered">
                                Scan with your app
                              </span>
                            </figure>
                          </div>
                        </div>

                        <br />
                        <p className="has-text-grey title is-3 has-text-centered">
                          - OR -
                        </p>
                        <br />

                        <p class="has-text-grey">
                          Copy and paste the following values into your device:
                        </p>
                        <br />
                        <FormTextareaField
                          label="Account Name"
                          name="accountName"
                          placeholder="-"
                          value={`${process.env.REACT_APP_WWW_DOMAIN}: ${currentUser.email}`}
                          errorText={null}
                          onChange={null}
                          isRequired={true}
                          maxWidth="380px"
                          readonly={true}
                        />
                        <FormTextareaField
                          label="Your Key"
                          name="yourKey"
                          placeholder="-"
                          value={otpResponse && otpResponse.base32}
                          errorText={errors && errors.base32}
                          onChange={null}
                          isRequired={true}
                          maxWidth="380px"
                          readonly={true}
                        />
                        <FormInputField
                          label="Type of Key"
                          name="typeOfKey"
                          placeholder="-"
                          value={"Time based"}
                          errorText={null}
                          helpText=""
                          onChange={null}
                          isRequired={true}
                          maxWidth="380px"
                          readonly={true}
                        />
                      </form>
                      <br />

                      <nav class="level">
                        <div class="level-left">
                          <div class="level-item">
                            <Link
                              class="button is-link is-fullwidth-mobile"
                              to="/login/2fa/step-1"
                            >
                              <FontAwesomeIcon icon={faArrowLeft} />
                              &nbsp;Back
                            </Link>
                          </div>
                        </div>
                        <div class="level-right">
                          <div class="level-item">
                            <Link
                              type="button"
                              class="button is-primary is-fullwidth-mobile"
                              to="/login/2fa/step-3"
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

export default TwoFactorAuthenticationWizardStep2;
