import React, { useState, useEffect } from "react";
import { Link, Navigate, useSearchParams } from "react-router-dom";
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
  faArrowLeft,
  faTriangleExclamation
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";


import FormErrorBox from "../Reusable/FormErrorBox";
import { currentUserState } from "../../AppState";
import { USER_ROLE_ROOT, USER_ROLE_RETAILER, USER_ROLE_CUSTOMER } from "../../Constants/App";


function CPSRNRedirector() {
  ////
  //// For debugging purposes only.
  ////

  console.log("REACT_APP_WWW_PROTOCOL:", process.env.REACT_APP_WWW_PROTOCOL);
  console.log("REACT_APP_WWW_DOMAIN:", process.env.REACT_APP_WWW_DOMAIN);
  console.log("REACT_APP_API_PROTOCOL:", process.env.REACT_APP_API_PROTOCOL);
  console.log("REACT_APP_API_DOMAIN:", process.env.REACT_APP_API_DOMAIN);

  ////
  //// URL Parameters.
  ////

  const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140
  const cpsrn = searchParams.get("v");


    ////
    //// Global State
    ////

    const [currentUser] = useRecoilState(currentUserState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [version, setVersion] = useState("");

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
      // getVersionAPI(onVersionSuccess, onVersionError, onVersionDone);
    }

    return () => (mounted = false);
  }, []);

  ////
  //// Component rendering.
  ////

  if (currentUser) {
      var forceURL = "";
      if (currentUser.otpEnabled === false) {
        console.log("onLoginSuccess | redirecting to dashboard");
        switch (currentUser.role) {
          case USER_ROLE_ROOT:
            forceURL = "/admin/registry/"+cpsrn;
            break;
          case USER_ROLE_RETAILER:
            forceURL = "/registry/"+cpsrn;
            break;
          case USER_ROLE_CUSTOMER:
            forceURL = "/c/registry/"+cpsrn;
            break;
          default:
            forceURL = "/501";
            break;
        }
      } else {
        if (currentUser.otpVerified === false) {
          console.log("onLoginSuccess | redirecting to 2fa setup wizard");
          forceURL = "/login/2fa/step-1?cpsrn="+cpsrn;
        } else {
          console.log("onLoginSuccess | redirecting to 2fa validation");
          forceURL = "/login/2fa?cpsrn="+cpsrn;
        }
      }
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
                            <img
                              src="/static/CPS logo 2023 GR.webp"
                              style={{ width: "256px" }}
                            />
                          </figure>
                        </div>
                      </nav>
                      {/* End Logo */}
                      <form>
                        <h1 className="title is-4 has-text-centered">
                          Collection Lookup
                        </h1>

                        {(!currentUser) ? (
                          <article className="message is-danger">
                            <div className="message-body">
                              <FontAwesomeIcon
                                className="fas"
                                icon={faTriangleExclamation}
                              />
                              &nbsp;You need an account.
                              <br /><br />
                              Please <b><Link to={`/login?cpsrn=${cpsrn}`}>login with your existing account&nbsp;<FontAwesomeIcon icon={faArrowRight} /></Link></b> or if you don't have an account then you will need to <b><Link to={`/register?cpsrn=${cpsrn}`}>create a new account&nbsp;<FontAwesomeIcon icon={faArrowRight} /></Link></b> to get started.
                            </div>
                          </article>
                      ) : (
                          <>
                              <i className="level-item has-text-centered">Loading, please wait...</i>
                          </>
                      )}

                        {/*
                        <Link
                          class="button is-medium is-block is-fullwidth is-primary"
                          type="button"
                          to="/login"
                          style={{ backgroundColor: "#FF0000" }}
                        >
                          Login <FontAwesomeIcon icon={faArrowRight} />
                        </Link>
                        <br />
                        <Link
                          class="button is-medium is-block is-fullwidth is-info"
                          type="button"
                          to="/register"
                        >
                          Register <FontAwesomeIcon icon={faArrowRight} />
                        </Link>
                        */}

                      </form>
                      <br />
                      {/*
                      <nav class="level">
                        <div class="level-item has-text-centered">
                          <div>
                            <a
                              href="https://cpscapsule.com"
                              className="is-size-7-tablet"
                            >
                              <FontAwesomeIcon icon={faArrowLeft} /> Back Home
                            </a>
                          </div>
                        </div>
                        <div class="level-item has-text-centered">
                          <div>
                            <Link
                              to="/cpsrn-registry"
                              className="is-size-7-tablet"
                            >
                              CPSRN Registry{" "}
                              <FontAwesomeIcon icon={faArrowRight} />
                            </Link>
                          </div>
                        </div>
                      </nav>
                      */}
                    </div>
                    {/* End box */}

                    <div className="has-text-centered">
                      <p>Need help?</p>
                      <p>
                        <a href="mail:support@cpscapsule.com">
                          support@cpscapsule.com
                        </a>
                      </p>
                      <p>
                        <a href="tel:+15199142685">(519) 914-2685</a>
                      </p>
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

export default CPSRNRedirector;
