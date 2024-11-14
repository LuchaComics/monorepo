import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";

import PageLoadingContent from "../Reusable/PageLoadingContent";
import {
    StartupApp,
} from "../../../wailsjs/go/main/App";

/**
 * The purpose of this view is to display a loading message to the user while
 * the application loads up / finishes initializing. Once completed this view
 * will redirect the user to the dashboard.
 */
function StartupView() {
    ////
    //// Component states.
    ////

    const [isLoading, setIsLoading] = useState(true);
    const [forceURL, setForceURL] = useState("");
    const [intervalId, setIntervalId] = useState(null);

    ////
    //// Event handling.
    ////

    ////
    //// Misc.
    ////

    useEffect(() => {
      let mounted = true;

      const startBackend = async () => {
        try {
            const result = await StartupApp();
            if (mounted) {
                console.log("Backend is ready.");
                setForceURL("/dashboard");
            }
        } catch (errResp) {
            console.log("Backend has error.", errResp);
            if (mounted) {
                setErrors(errResp);
            }
        } finally {
            if (mounted) {
                console.log("Finished.");
            }
        }
      };

      startBackend();

      if (mounted) {
        window.scrollTo(0, 0); // Start the page at the top of the page.
      }

       // Cleanup the interval when the component unmounts
      return () => {
        console.log("unmounted")
        // clearInterval(interval);
        mounted = false;
      };
  }, []); // The empty dependency array ensures the effect runs only once on mount

    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
        return <Navigate to={forceURL} />;
    }

    if (isLoading) {
        return (
            <PageLoadingContent displayMessage="Please wait..." />
        )
    }

    return (
        <>
          <div class="container">
            <section class="section">
              <nav class="box">
                <div class="columns">
                  <div class="column">
                    <h1 class="title is-4">
                      &nbsp;Error Message
                    </h1>
                  </div>
                </div>



                <section class="hero is-warning is-medium">
                  <div class="hero-body">
                    <p class="title">Requires IPFS Node Running</p>
                    <p class="subtitle">Cannot start the application without IPFS node on your computer. Please load it up and when ready this page will be removed.</p>
                  </div>
                </section>



            </nav>
            </section>
          </div>
        </>
    )
}

export default StartupView
