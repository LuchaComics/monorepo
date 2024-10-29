import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";

import PageLoadingContent from "../Reusable/PageLoadingContent";
import {GetIsIPFSRunning} from "../../../wailsjs/go/main/App";


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
    const backgroundPollingTick = () => {
        GetIsIPFSRunning().then((isIPFSRunningRes)=>{
            console.log("isIPFSRunningRes:", isIPFSRunningRes);
            if (isIPFSRunningRes === true) {
              setForceURL("/dashboard");
            }
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
        <PageLoadingContent displayMessage="Starting up..." />
    )
}

export default StartupView
