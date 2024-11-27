import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faTachometer,
  faEye,
  faPencil,
  faTrashCan,
  faPlus,
  faGauge,
  faArrowRight,
  faBarcode,
  faTriangleExclamation
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import FormErrorBox from "../Reusable/FormErrorBox";
import { getVersionAPI } from "../../API/Gateway";
import { currentUserState } from "../../AppState";


function PublicRegistrySearch() {
  // For debugging purposes only.
  console.log("REACT_APP_WWW_PROTOCOL:", process.env.REACT_APP_WWW_PROTOCOL);
  console.log("REACT_APP_WWW_DOMAIN:", process.env.REACT_APP_WWW_DOMAIN);
  console.log("REACT_APP_API_PROTOCOL:", process.env.REACT_APP_API_PROTOCOL);
  console.log("REACT_APP_API_DOMAIN:", process.env.REACT_APP_API_DOMAIN);

  ////
  //// Global State
  ////

  const [currentUser] = useRecoilState(currentUserState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [validation, setValidation] = useState({
    cpsrn: false,
  });
  const [version, setVersion] = useState("");
  const [cpsrn, setCpsn] = useState("");
  const [forceURL, setForceURL] = useState("");

  ////
  //// API.
  ////

  function onVersionSuccess(response) {
    console.log("onVersionSuccess: Starting...");
    setVersion(response);
  }

  function onVersionError(apiErr) {
    console.log("onVersionError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onVersionDone() {
    console.log("onVersionDone: Starting...");
  }

  ////
  //// Event handling.
  ////

  function onButtonClick(e) {
    var newErrors = {};
    var newValidation = {};
    if (cpsrn === undefined || cpsrn === null || cpsrn === "") {
      newErrors["cpsrn"] = "value is missing";
    } else {
      newValidation["cpsrn"] = true;
    }

    /// Save to state.
    setErrors(newErrors);
    setValidation(newValidation);

    if (Object.keys(newErrors).length > 0) {
      //
      // Handle errors.
      //

      console.log("failed validation");

      // window.scrollTo(0, 0);  // Start the page at the top of the page.

      // The following code will cause the screen to scroll to the top of
      // the page. Please see ``react-scroll`` for more information:
      // https://github.com/fisshy/react-scroll
      var scroll = Scroll.animateScroll;
      scroll.scrollToTop();
    } else {
      //
      // Submit to server.
      //

      console.log("successful validation, submitting to API server.");
      setForceURL("/cpsrn-result?v=" + cpsrn);
    }
  }

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      getVersionAPI(onVersionSuccess, onVersionError, onVersionDone);
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
      <div class="column is-12 container">
        <div class="section">
          <section class="hero is-fullheight">
            <div class="hero-body">
              <div class="container">
                <div class="columns is-centered">
                  <div class="column is-half-tablet">
                    <div class="box is-rounded">
                      {/* Start Logo */}
                      <nav class="level">
                        <div class="level-item has-text-centered">
                          <figure class="image">
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
                      <form>
                        <h1 className="title is-4 has-text-centered">
                          Check your Grading
                        </h1>
                        {currentUser === null && (
                          <article className="message is-danger">
                            <div className="message-body">
                              <FontAwesomeIcon
                                className="fas"
                                icon={faTriangleExclamation}
                              />
                              &nbsp;An account is required for lookup. <Link to="/login"><b>Clich here&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} /></b></Link> to sign into your account.</div>
                          </article>
                        )}
                        <FormErrorBox errors={errors} />
                        <div class="field">
                          <label class="label is-small has-text-grey-light">
                            COMICCOIN_FAUCET Registry #
                          </label>
                          <div class="control has-icons-left has-icons-right">
                            <input
                              class={`input ${errors && errors.cpsrn && "is-danger"} ${validation && validation.cpsrn && "is-success"}`}
                              name="cpsrn"
                              type="text"
                              placeholder="Enter COMICCOIN_FAUCET registry number."
                              value={cpsrn}
                              onChange={(e) => setCpsn(e.target.value)}
                              disabled={currentUser === null}
                            />
                            <span class="icon is-small is-left">
                              <FontAwesomeIcon
                                className="fas"
                                icon={faBarcode}
                              />
                            </span>
                          </div>
                          {errors && errors.cpsrn && (
                            <p class="help is-danger">{errors.cpsrn}</p>
                          )}
                          <p class="help">
                            Enter the identification number found with your
                            collectible.
                          </p>
                        </div>
                        <br />
                        <button
                          class="button is-medium is-block is-fullwidth is-primary"
                          type="button"
                          onClick={onButtonClick}
                          style={{ backgroundColor: "#FF0000" }}
                          disabled={currentUser === null}
                        >
                          Lookup COMICCOIN_FAUCETRN <FontAwesomeIcon icon={faArrowRight} />
                        </button>
                      </form>
                      <br />
                      <nav class="level">
                        <div class="level-item has-text-centered">
                          <div>
                            <Link to="/login" className="is-size-7-tablet">
                              Login
                            </Link>
                          </div>
                        </div>
                        <div class="level-item has-text-centered">
                          <div>
                            <Link to="/register" className="is-size-7-tablet">
                              Create an Account
                            </Link>
                          </div>
                        </div>
                      </nav>
                    </div>
                    {/* End box */}

                    <div className="has-text-centered">
                      <p>Need help?</p>
                      <p>
                        <Link to="Support@cpscapsule.com">
                          Support@cpscapsule.com
                        </Link>
                      </p>
                      <p>
                        <a href="tel:+15199142685">(519) 914-2685</a>
                      </p>
                      <p>(App version: {version})</p>
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

export default PublicRegistrySearch;
