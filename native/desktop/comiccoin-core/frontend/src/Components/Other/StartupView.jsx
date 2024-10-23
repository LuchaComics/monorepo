import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";

import PageLoadingContent from "../Reusable/PageLoadingContent";
import {GetIsBlockhainNodeRunning} from "../../../wailsjs/go/main/App";


function StartupView() {
    ////
    //// Component states.
    ////

    const [forceURL, setForceURL] = useState("");
    const [intervalId, setIntervalId] = useState(null);

    ////
    //// Event handling.
    ////

    // Function will make a call to check to see if our node is running
    // and if our backend says the node is running then we will redirect
    // to another page.
    const backgroundPollingTick = (e) => {
        GetIsBlockhainNodeRunning().then( (isNodeRunningResponse)=>{
            console.log("tick", new Date().getTime(), isNodeRunningResponse);
            if (isNodeRunningResponse) {
                clearInterval(intervalId);
                setForceURL("/wallets");
            }
        })
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
      let mounted = true;

      if (mounted) {
            window.scrollTo(0, 0); // Start the page at the top of the page.
            const interval = setInterval(() => backgroundPollingTick(), 1000);
            setIntervalId(interval);
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
                  Starting up...
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
