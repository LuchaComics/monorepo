import {useState, useEffect} from 'react';
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faGauge,
  faArrowRight,
  faUsers,
  faBarcode,
  faCubes,
  faInbox
} from "@fortawesome/free-solid-svg-icons";


import logo from '../../assets/images/CPS-logo-2023-square.webp';

function ReceiveView() {


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
                  <li>
                    <Link to="/dashboard" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faGauge} />
                      &nbsp;Overview
                    </Link>
                  </li>
                  <li class="is-active">
                    <Link to="/dashboard" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faInbox} />
                      &nbsp;Receive
                    </Link>
                  </li>
                </ul>
              </nav>

              <nav class="box">
                <div class="columns">
                  <div class="column">
                    <h1 class="title is-4">
                      <FontAwesomeIcon className="fas" icon={faInbox} />
                      &nbsp;Receive ComicCoins
                    </h1>
                  </div>
                </div>



              </nav>
            </section>
          </div>
        </>
    )
}

export default ReceiveView
