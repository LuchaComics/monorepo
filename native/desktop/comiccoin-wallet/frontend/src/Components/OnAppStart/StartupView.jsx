import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";
import { useRecoilState } from "recoil";
import { WalletMinimal } from 'lucide-react';

import { GetDataDirectoryFromPreferences, GetIsBlockhainNodeRunning, DefaultWalletAddress } from "../../../wailsjs/go/main/App";
import { currentOpenWalletAtAddressState } from "../../AppState";


function StartupView() {
    // Global State
    const [currentOpenWalletAtAddress] = useRecoilState(currentOpenWalletAtAddressState);

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
                    if (currentOpenWalletAtAddress) {
                        console.log("currentOpenWalletAtAddress:", currentOpenWalletAtAddress);
                        setForceURL("/dashboard")
                    } else {
                        setForceURL("/create-your-first-wallet")
                    }
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
      console.log("Redirected to URL:", forceURL);
      return <Navigate to={forceURL} />;
    }

    return (
        <div>
          <main className="flex-grow flex items-center justify-center">
            <div className="text-center space-y-8 px-6">
              <div className="flex justify-center">
                <div className="relative w-20 h-20">
                  <div className="absolute inset-0 border-4 border-purple-200 rounded-full"></div>
                  <div className="absolute inset-0 border-4 border-purple-600 rounded-full border-t-transparent animate-spin"></div>
                </div>
              </div>

              <div className="space-y-3">
                <h2 className="text-2xl font-bold text-gray-900">Starting Wallet</h2>
                <p className="text-gray-600 max-w-sm mx-auto">
                  Please wait while we start up. This may take a few moments.
                </p>
              </div>

              <div className="flex items-center justify-center gap-2">
                <div className="w-2 h-2 bg-purple-600 rounded-full animate-bounce"></div>
                <div className="w-2 h-2 bg-purple-600 rounded-full animate-bounce [animation-delay:0.2s]"></div>
                <div className="w-2 h-2 bg-purple-600 rounded-full animate-bounce [animation-delay:0.4s]"></div>
              </div>
            </div>
          </main>
        </div>
    );
}

export default StartupView
