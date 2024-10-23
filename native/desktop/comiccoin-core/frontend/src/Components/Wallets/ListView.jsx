import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";
import { useRecoilState } from "recoil";

import PageLoadingContent from "../Reusable/PageLoadingContent";
import { currentOpenWalletAtAddressState } from "../../AppState";
import { DefaultWalletAddress, ListWallets } from "../../../wailsjs/go/main/App";


function ListWalletsView() {
    ////
    //// Component states.
    ////

    const [wallets, setWallets] = useState([]);
    const [errors, setErrors] = useState({});
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

            // Fetch the wallets we have in our app.
            ListWallets().then((walletsResponse, errResponse) => {
                console.log("walletsResponse:", walletsResponse);
                console.log("errResponse:", errResponse);
                setWallets(walletsResponse);
            })

            DefaultWalletAddress().then((addressResponse)=>{
                console.log("address:", addressResponse);
                if (addressResponse !== "") {
                    setForceURL("/dashboard");
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
        <>
            {wallets.length === 0 ? <>
                <section class="hero is-fullheight-with-navbar is-info">
                  <div class="hero-body">
                    <div class="container">
                      <div class="columns is-centered">
                        <div class="column is-6-tablet is-5-desktop is-4-widescreen">
                          <h1 class="title is-2">Welcome to ComicCoin Core</h1>
                          <h2 class="subtitle is-4">To begin, please create your wallet to get started.</h2>
                          <div class="buttons">
                            <Link class="button is-link is-medium" to="/wallet/add">Click here to create your wallet</Link>
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </section>
            </> : <>
            </>}
        </>
    )
}

export default ListWalletsView
