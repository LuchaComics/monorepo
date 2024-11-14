import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";

import PageLoadingContent from "../Reusable/PageLoadingContent";
import {
    GetDataDirectoryFromPreferences,
    GetNFTStoreSettingsFromPreferences
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
            GetDataDirectoryFromPreferences().then( (dataDirResp) => {
                console.log("dataDirResp:", dataDirResp);
                if (dataDirResp === "") {
                    setForceURL("/pick-data-directory")
                } else {
                    GetNFTStoreSettingsFromPreferences().then((config)=>{
                        console.log("config", config);
                        if (config === undefined || config === null || config === "") {
                             setForceURL("/config-nftstore");
                        } else {
                            if (config.accessKeyId === "" || config.apiVersion === "" || config.endpoint === "" || config.region === "" || config.s3ForcePathStyle === "" || config.secretAccessKey === "") {
                                console.log("config:", config);
                                setForceURL("/config-nftstore");
                            } else {
                                console.log("config:", config);
                                setForceURL("/startup");
                            }
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
