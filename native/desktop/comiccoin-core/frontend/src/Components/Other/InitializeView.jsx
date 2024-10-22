import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";

import PageLoadingContent from "../Reusable/PageLoadingContent";
import {GetDataDirectoryFromPreferences} from "../../../wailsjs/go/main/App";


function InitializeView() {
    ////
    //// Component states.
    ////

    const [dataDirectory] = useState("");
    const [forceURL, setForceURL] = useState("");

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
            GetDataDirectoryFromPreferences().then( (dataDirResult) => {
                console.log("dataDirResult:", dataDirResult);
                if (dataDirResult === "") {
                    setForceURL("/pick-data-directory")
                } else {
                    setForceURL("/startup")
                }
            })
      }

      return () => {
        mounted = false;
      };
    }, []);

    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
      return <Navigate to={forceURL} />;
    }

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
            </div>
          </div>
        </section>
      </div>
    </div>
    )
}

export default InitializeView
