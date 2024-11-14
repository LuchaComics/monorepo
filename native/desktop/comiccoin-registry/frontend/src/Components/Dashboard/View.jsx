import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faGauge,
  faArrowRight,
  faUsers,
  faBarcode,
  faCubes,
  faFileInvoiceDollar,
  faCoins
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import logo from '../../assets/images/CPS-logo-2023-square.webp';
import { currentOpenWalletAtAddressState } from "../../AppState";
import PageLoadingContent from "../Reusable/PageLoadingContent";
import { toLower } from "lodash";


function DashboardView() {
    ////
    //// Global State
    ////

    const [currentOpenWalletAtAddress] = useRecoilState(currentOpenWalletAtAddressState);

    ////
    //// Component states.
    ////

    const [isLoading, setIsLoading] = useState(false);
    const [forceURL, setForceURL] = useState("");
    const [totalCoins, setTotalCoins] = useState(0);
    const [totalTokens, setTotalTokens] = useState(0);
    const [transactions, setTransactions] = useState([]);

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
            console.log("currentOpenWalletAtAddress:", currentOpenWalletAtAddress);


      }
      return () => {
        mounted = false;
      };
    }, [currentOpenWalletAtAddress]);

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
                      <FontAwesomeIcon className="fas" icon={faGauge} />
                      &nbsp;Overview
                    </h1>
                  </div>
                </div>

                {isLoading
                ? <>
                    <PageLoadingContent displayMessage="Fetching..." />

                </> : <>
                {/*
                <nav class="level">
                  <div class="level-item has-text-centered">
                    <div>
                      <p class="heading">Coins</p>
                      <p class="title">{totalCoins}</p>
                    </div>
                  </div>
                  <div class="level-item has-text-centered">
                    <div>
                      <p class="heading">Tokens Count</p>
                      <p class="title">{totalTokens}</p>
                    </div>
                  </div>
                </nav>
                */}

                <section class="hero is-info is-medium">
                  <div class="hero-body">
                    <p class="title"><FontAwesomeIcon className="fas" icon={faCubes} />&nbsp;Tokens</p>
                    <p class="subtitle">To view the tokens in the ComicCoin blockchain or create new tokens continue below:</p>
                    <p class="subtitle">
                        <Link to="/tokens">
                        View List&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} />
                        </Link>
                    </p>
                    <p class="subtitle">
                        <Link to="/tokens/new">
                        New Token&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} />
                        </Link>
                    </p>
                  </div>
                </section>
                </>}


            </nav>
            </section>
          </div>
        </>
    )
}

export default DashboardView
