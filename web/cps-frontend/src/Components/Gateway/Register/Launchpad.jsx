import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faExclamationTriangle,
  faTasks,
  faTachometer,
  faPlus,
  faArrowLeft,
  faCheckCircle,
  faArrowRight,
  faArrowUpRightFromSquare,
} from "@fortawesome/free-solid-svg-icons";
import Select from "react-select";
import { useRecoilState } from "recoil";

import useLocalStorage from "../../../Hooks/useLocalStorage";
import { postRegisterAPI } from "../../../API/Gateway";
import FormErrorBox from "../../Reusable/FormErrorBox";
import FormInputField from "../../Reusable/FormInputField";
import PageLoadingContent from "../../Reusable/PageLoadingContent";
import FormTextareaField from "../../Reusable/FormTextareaField";
import FormRadioField from "../../Reusable/FormRadioField";
import FormMultiSelectField from "../../Reusable/FormMultiSelectField";
import FormSelectField from "../../Reusable/FormSelectField";
import FormCheckboxField from "../../Reusable/FormCheckboxField";
import FormCountryField from "../../Reusable/FormCountryField";
import FormRegionField from "../../Reusable/FormRegionField";
import FormTimezoneSelectField from "../../Reusable/FormTimezoneField";
import {
  HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS,
  HOW_LONG_HAS_YOUR_STORE_BEEN_OPERATING_FOR_WITH_EMPTY_OPTIONS,
  ESTIMATED_SUBMISSION_PER_MONTH_WITH_EMPTY_OPTIONS,
} from "../../../Constants/FieldOptions";
import { topAlertMessageState, topAlertStatusState } from "../../../AppState";

function RegisterLaunchpad() {
  ////
  ////
  ////

  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);

  ////
  //// Component states.
  ////

  ////
  //// Event handling.
  ////

  ////
  //// API.
  ////

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.
    }

    return () => {
      mounted = false;
    };
  }, []);

  ////
  //// Component rendering.
  ////

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
                                src="/static/CPS logo 2023 GR.webp"
                                style={{ width: "256px" }}
                              />
                            </Link>
                          </figure>
                        </div>
                      </nav>
                      {/* End Logo */}
                      <form>
                        <h1 className="title is-4 has-text-centered">
                          What do you want to register as?
                        </h1>
                        <FormErrorBox errors={null} />
                        <Link
                          class="button is-medium is-block is-fullwidth is-primary"
                          type="button"
                          to="/register/store"
                          style={{ backgroundColor: "#FF0000" }}
                        >
                          Register as business user{" "}
                          <FontAwesomeIcon icon={faArrowRight} />
                        </Link>
                        &nbsp;
                        <Link
                          class="button is-medium is-block is-fullwidth is-info"
                          type="button"
                          to="/register/user"
                        >
                          Register as regular user{" "}
                          <FontAwesomeIcon icon={faArrowRight} />
                        </Link>
                      </form>
                      <br />
                      <nav class="level">
                        <div class="level-item has-text-centered">
                          <div>
                          <Link to="/" className="is-size-7-tablet">
                             <FontAwesomeIcon icon={faArrowLeft} />{" "}Back
                          </Link>
                          </div>
                        </div>
                        <div class="level-item has-text-centered">
                          <div>

                            <Link to="/login" className="is-size-7-tablet">
                              Login{" "}
                              <FontAwesomeIcon icon={faArrowRight} />
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
                      <p>(App version: 11)</p>
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

export default RegisterLaunchpad;
