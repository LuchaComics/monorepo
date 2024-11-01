import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";
import { useRecoilState } from "recoil";

import PageLoadingContent from "../Reusable/PageLoadingContent";
import { GetIsBlockhainNodeRunning, DefaultWalletAddress } from "../../../wailsjs/go/main/App";
import { currentOpenWalletAtAddressState } from "../../AppState";


function StartupView() {
    ////
    //// Global State
    ////

    const [currentOpenWalletAtAddress, setCurrentOpenWalletAtAddress] = useRecoilState(currentOpenWalletAtAddressState);

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

                console.log("tick: done");
                clearInterval(intervalId);

                // Check to see if we already have an address set, else
                // the user needs to log in again.
                DefaultWalletAddress().then((addressResponse)=>{
                    console.log("default wallet address:", addressResponse);
                    if (addressResponse !== undefined && addressResponse !== null && addressResponse !== "") {
                        console.log("currentOpenWalletAtAddress:", currentOpenWalletAtAddress);
                        setCurrentOpenWalletAtAddress(addressResponse);
                        setForceURL("/dashboard");
                    } else {
                        setForceURL("/wallets");
                    }
                })
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
        <PageLoadingContent displayMessage="Starting up..." />
    )
}

export default StartupView
