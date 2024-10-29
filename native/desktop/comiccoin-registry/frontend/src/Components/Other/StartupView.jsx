import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";

import PageLoadingContent from "../Reusable/PageLoadingContent";
import {GetIsIPFSRunning} from "../../../wailsjs/go/main/App";


function StartupView() {
    ////
    //// Component states.
    ////

    const [isLoading, setIsLoading] = useState(false);
    const [forceURL, setForceURL] = useState("");
    const [intervalId, setIntervalId] = useState(null);

    ////
    //// Event handling.
    ////

    // Function will make a call to check to see if our node is running
    // and if our backend says the node is running then we will redirect
    // to another page.
    const backgroundPollingTick = () => {
        // Update the GUI to let user know that the operation is under way.
        setIsLoading(true);

        GetIsIPFSRunning().then((isIPFSRunningRes)=>{
            console.log("isIPFSRunningRes:", isIPFSRunningRes);
            if (isIPFSRunningRes === true) {
              setForceURL("/dashboard");
            }
        }).finally(() => {
            // this will be executed after then or catch has been executed
            console.log("promise has been resolved or rejected");

            // Update the GUI to let user know that the operation is completed.
            setIsLoading(false);
        });
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
      let mounted = true;

      // Initalize the background timer to tick every second.
      const interval = setInterval(backgroundPollingTick, 1000);
      setIntervalId(interval);

      if (mounted) {
        window.scrollTo(0, 0); // Start the page at the top of the page.
      }

       // Cleanup the interval when the component unmounts
      return () => {
        console.log("unmounted")
        clearInterval(interval);
        mounted = false;
      };
    }, []); // The empty dependency array ensures the effect runs only once on mount

    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
        return <Navigate to={forceURL} />;
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
