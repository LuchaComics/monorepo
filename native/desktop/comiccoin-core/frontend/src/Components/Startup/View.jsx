import {useState, useEffect} from 'react';
import { Link } from "react-router-dom";

import PageLoadingContent from "../Reusable/PageLoadingContent";


function StartupView() {
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
        <div className="column is-12">
      <div className="section">
        <section className="hero is-fullheight">
          <div className="hero-body">
            <div className="container">
              <div className="columns is-centered">
                <div className="column is-4 has-text-centered">
                  Loading...
                </div>
              </div>

              <div className="columns is-centered" style={{ paddingTop: "20px" }}>
                <div className="column is-4 has-text-centered">
                  <Link to="/pick-storage-location-on-startup" className="button is-primary is-large">
                    Go to next page
                  </Link>
                </div>
              </div>

            </div>
          </div>
        </section>
      </div>
    </div>
    )
}

export default StartupView
