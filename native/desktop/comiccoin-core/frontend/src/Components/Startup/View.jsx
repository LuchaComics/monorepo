import {useState, useEffect} from 'react';
import { Link } from "react-router-dom";


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
        <div class="column is-12 container">
          <div class="section">
            <section class="hero is-fullheight">
              <div class="hero-body">
                <div class="container">
                  <div class="columns is-centered">
                    <div class="column is-half-tablet">
                      <h1 className="is-size-1">Loading ...</h1>
                      <Link to="/pick-storage-location-on-startup">Go to next page</Link>
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
