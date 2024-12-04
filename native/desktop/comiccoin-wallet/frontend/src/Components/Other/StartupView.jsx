import { useState, useEffect, useRef } from "react";
import { Navigate } from "react-router-dom";
import { useRecoilState } from "recoil";

import PageLoadingContent from "../Reusable/PageLoadingContent";
import { GetIsBlockhainNodeRunning, DefaultWalletAddress } from "../../../wailsjs/go/main/App";
import { currentOpenWalletAtAddressState } from "../../AppState";

function StartupView() {
  // Global State
  const [currentOpenWalletAtAddress, setCurrentOpenWalletAtAddress] = useRecoilState(currentOpenWalletAtAddressState);

  // Component States
  const [forceURL, setForceURL] = useState("");

   // Developers note:
   // React `useRef` provides a way to persist a mutable reference across
   // renders without triggering re-renders. You can use it to store the
   // interval ID and reference it directly when calling `clearInterval`.

  // Ref to store the interval ID
  const intervalIdRef = useRef(null);

  // Function to poll the blockchain node status
  const backgroundPollingTick = () => {
    GetIsBlockhainNodeRunning().then((isNodeRunningResponse) => {
      console.log("tick", new Date().getTime(), isNodeRunningResponse);

      if (isNodeRunningResponse) {
        console.log("tick: done");

        // Stop the interval
        if (intervalIdRef.current) {
          clearInterval(intervalIdRef.current);
          intervalIdRef.current = null; // Reset the ref to avoid unintended usage
        }

        // Check default wallet address and redirect accordingly
        DefaultWalletAddress().then((addressResponse) => {
          console.log("default wallet address:", addressResponse);
          if (addressResponse) {
            console.log("currentOpenWalletAtAddress:", currentOpenWalletAtAddress);
            setCurrentOpenWalletAtAddress(addressResponse);
            setForceURL("/dashboard");
          } else {
            setForceURL("/wallets");
          }
        });
      }
    });
  };

  useEffect(() => {
    // Start the interval when the component mounts
    intervalIdRef.current = setInterval(() => backgroundPollingTick(), 1000);

    // Cleanup: Clear the interval when the component unmounts
    return () => {
      if (intervalIdRef.current) {
        clearInterval(intervalIdRef.current);
      }
    };
  }, []); // Empty dependency array ensures this runs only on mount

  // Redirect if `forceURL` is set
  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  // Render loading message
  return <PageLoadingContent displayMessage="Starting up..." />;
}

export default StartupView;
