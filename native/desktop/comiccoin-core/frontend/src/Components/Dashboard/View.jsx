import {useState, useEffect} from 'react';
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faGauge,
  faArrowRight,
  faUsers,
  faBarcode,
  faCubes
} from "@fortawesome/free-solid-svg-icons";

function DashboardView() {


    useEffect(() => {
      let mounted = true;

      if (mounted) {
            window.scrollTo(0, 0); // Start the page at the top of the page.
      }


      return () => {
        mounted = false;
      };
    }, []);



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
                      &nbsp;NFT Collections
                    </p>
                    <p class="subtitle">
                      Submit a request to encapsulate your collectible by clicking
                      below:
                      <br />
                      <br />
                      <Link to={"/admin/collections"}>
                        View&nbsp;
                        <FontAwesomeIcon className="fas" icon={faArrowRight} />
                      </Link>
                      <br />
                      <br />
                      <Link to={"/admin/collections/add/step-1"}>
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
    )
}

export default DashboardView
