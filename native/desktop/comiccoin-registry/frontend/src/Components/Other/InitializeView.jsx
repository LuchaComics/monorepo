import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";

import PageLoadingContent from "../Reusable/PageLoadingContent";
import {
    GetDataDirectoryFromPreferences,
    GetNFTStoreRemoteAddressFromPreferences,
    GetNFTStoreAPIKeyFromPreferences
} from "../../../wailsjs/go/main/App";


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

            //
            // STEP 1
            // Confirm our data directory was set.
            //

            GetDataDirectoryFromPreferences().then( (dataDirResp) => {
                console.log("dataDirResp:", dataDirResp);
                if (dataDirResp === "") {
                    setForceURL("/pick-data-directory")
                } else {
                    //
                    // STEP 2
                    // Confirm our NFT Storage remote address was set.
                    //

                    GetNFTStoreRemoteAddressFromPreferences().then((remoteAddressResp)=>{
                        console.log("remoteAddressResp", remoteAddressResp);
                        if (remoteAddressResp != "") {

                            //
                            // STEP
                            // Confirm our NFT Storage API key was set.
                            //

                            GetNFTStoreAPIKeyFromPreferences().then((apiKeyResp)=>{
                                console.log("apiKeyResp", apiKeyResp);
                                if (apiKeyResp != "") {
                                    setForceURL("/startup");
                                } else {
                                    setForceURL("/config-nftstore");
                                }
                            })

                        } else {
                            setForceURL("/config-nftstore");
                        }
                    })
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
        <PageLoadingContent displayMessage="Initializing..." />
    )
}

export default InitializeView
