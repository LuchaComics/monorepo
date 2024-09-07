import React, { useEffect } from "react";
import { Link, useSearchParams, Navigate } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faGauge,
  faArrowRight,
  faUsers,
  faBarcode,
  faCubes
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import { topAlertMessageState, topAlertStatusState } from "../../AppState";

function AdminDashboard() {
  ////
  //// URL Parameters.
  ////

  const [searchParams] = useSearchParams(); // Special thanks via https://stackoverflow.com/a/65451140
  const cpsrn = searchParams.get("cpsrn");

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
  //// API.
  ////

  ////
  //// Event handling.
  ////

  ////
  //// Misc.
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

  // If the `cpsrn` url parameter exists then redirect the user to the registry page.
  if (cpsrn) {
      return <Navigate to={`/admin/registry/${cpsrn}`} />;
  }

  return (
    <>
      <div class="container">
        <section class="section">
          <nav class="breadcrumb" aria-label="breadcrumbs">
            <ul>
              <li class="is-active">
                <Link to="/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Admin Dashboard
                </Link>
              </li>
            </ul>
          </nav>
          <nav class="box">
            <div class="columns">
              <div class="column">
                <h1 class="title is-4">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Admin Dashboard
                </h1>
              </div>
            </div>

              {/*
            <section class="hero is-medium is-link">
              <div class="hero-body">
                <p class="title">
                  <FontAwesomeIcon className="fas" icon={faTasks} />
                  &nbsp;Online Submissions
                </p>
                <p class="subtitle">
                  Submit a request to encapsulate your collectible by clicking
                  below:
                  <br />
                  <br />
                  <Link to={"/admin/submissions/comics"}>
                    View Online Comic Submissions&nbsp;
                    <FontAwesomeIcon className="fas" icon={faArrowRight} />
                  </Link>
                  <br />
                  <br />
                  <Link to={"/admin/submissions/pick-type-for-add"}>
                    Add&nbsp;
                    <FontAwesomeIcon className="fas" icon={faArrowRight} />
                  </Link>
                </p>
              </div>
            </section>

            <section class="hero is-medium is-primary">
              <div class="hero-body">
                <p class="title">
                  <FontAwesomeIcon className="fas" icon={faBarcode} />
                  &nbsp;Registry
                </p>
                <p class="subtitle">
                  Have a CPS registry number? Use the following to lookup
                  existing records:
                  <br />
                  <br />
                  <Link to={"/admin/registry"}>
                    View&nbsp;
                    <FontAwesomeIcon className="fas" icon={faArrowRight} />
                  </Link>
                </p>
              </div>
            </section>
            */}
            <section class="hero is-medium is-success">
              <div class="hero-body">
                <p class="title">
                  <FontAwesomeIcon className="fas" icon={faTasks} />
                  &nbsp;Tenants
                </p>
                <p class="subtitle">
                  Manage the tenants that belong to your system.
                  <br />
                  <br />
                  <Link to={"/admin/tenants"}>
                    View&nbsp;
                    <FontAwesomeIcon className="fas" icon={faArrowRight} />
                  </Link>
                </p>
              </div>
            </section>
            <section class="hero is-medium is-link">
              <div class="hero-body">
                <p class="title">
                  <FontAwesomeIcon className="fas" icon={faCubes} />
                  &nbsp;Collections
                </p>
                <p class="subtitle">
                  Submit a request to encapsulate your collectible by clicking
                  below:
                  <br />
                  <br />
                  <Link to={"/admin/collections"}>
                    View Collections&nbsp;
                    <FontAwesomeIcon className="fas" icon={faArrowRight} />
                  </Link>
                  <br />
                  <br />
                  <Link to={"/admin/collections/add"}>
                    Add&nbsp;
                    <FontAwesomeIcon className="fas" icon={faArrowRight} />
                  </Link>
                </p>
              </div>
            </section>
            <section class="hero is-medium is-info">
              <div class="hero-body">
                <p class="title">
                  <FontAwesomeIcon className="fas" icon={faUsers} />
                  &nbsp;All Users
                </p>
                <p class="subtitle">
                  Manage all the users that belong to your system.
                  <br />
                  <br />
                  <Link to={"/admin/users"}>
                    View&nbsp;
                    <FontAwesomeIcon className="fas" icon={faArrowRight} />
                  </Link>
                </p>
              </div>
            </section>
          </nav>
        </section>
      </div>
    </>
  );
}

export default AdminDashboard;
