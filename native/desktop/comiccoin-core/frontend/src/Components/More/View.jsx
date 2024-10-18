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
  faPaperPlane,
  faEllipsis
} from "@fortawesome/free-solid-svg-icons";


import logo from '../../assets/images/CPS-logo-2023-square.webp';

function MoreView() {


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
                    <Link to="/more" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faEllipsis} />
                      &nbsp;More
                    </Link>
                  </li>
                </ul>
              </nav>

              <nav class="box">
                <div class="columns">
                  <div class="column">
                    <h1 class="title is-4">
                      <FontAwesomeIcon className="fas" icon={faEllipsis} />
                      &nbsp;More
                    </h1>
                  </div>
                </div>

                <div className="section">
      <div className="container">
        <div className="columns is-multiline">
          <div className="column is-4">
            <div className="box">
              <article className="media">
                <div className="media-left">
                  <i className="fas fa-list-alt fa-2x" />
                </div>
                <div className="media-content">
                  <div className="content">
                    <h4>Transactions</h4>
                    <p>View your recent transactions</p>
                  </div>
                </div>
              </article>
              <Link to="/transactions" className="button is-fullwidth is-link">Go to Transactions</Link>
            </div>
          </div>
          <div className="column is-4">
            <div className="box">
              <article className="media">
                <div className="media-left">
                  <i className="fas fa-coins fa-2x" />
                </div>
                <div className="media-content">
                  <div className="content">
                    <h4>Tokens</h4>
                    <p>View and manage your tokens</p>
                  </div>
                </div>
              </article>
              <Link to="/tokens" className="button is-fullwidth is-link">Go to Tokens</Link>
            </div>
          </div>
          <div className="column is-4">
            <div className="box">
              <article className="media">
                <div className="media-left">
                  <i className="fas fa-cog fa-2x" />
                </div>
                <div className="media-content">
                  <div className="content">
                    <h4>Settings</h4>
                    <p>Configure your account settings</p>
                  </div>
                </div>
              </article>
              <Link to="/settings" className="button is-fullwidth is-link">Go to Settings</Link>
            </div>
          </div>

          
        </div>
      </div>
    </div>


              </nav>
            </section>
          </div>
        </>
    )
}

export default MoreView
