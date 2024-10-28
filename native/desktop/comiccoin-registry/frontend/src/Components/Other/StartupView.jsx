import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";

import PageLoadingContent from "../Reusable/PageLoadingContent";


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
        console.log("tick: done");
        // clearInterval(intervalId);
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
      let mounted = true;

      if (mounted) {
            window.scrollTo(0, 0); // Start the page at the top of the page.
            // const interval = setInterval(() => backgroundPollingTick(), 1000);
            // setIntervalId(interval);
            setForceURL("/dashboard");
      }

      return () => {
        mounted = false;
      };
    }, [forceURL]);

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
