import {useState, useEffect} from 'react';
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
    faTasks,
    faGauge,
    faArrowRight,
    faUsers,
    faBarcode,
    faCubes,
    faCoins,
    faEllipsis
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { toLower } from "lodash";

import { GetTokens } from "../../../wailsjs/go/main/App";
import { currentOpenWalletAtAddressState } from "../../AppState";


function ListTokensView() {
    ////
    //// Global State
    ////

    const [currentOpenWalletAtAddress] = useRecoilState(currentOpenWalletAtAddressState);

    ////
    //// Component states.
    ////

    const [forceURL, setForceURL] = useState("");
    const [totalCoins, setTotalCoins] = useState(0);
    const [totalTokens, setTotalTokens] = useState(0);
    const [tokens, setTokens] = useState([]);

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

            GetTokens(currentOpenWalletAtAddress).then((txsResponse)=>{
                console.log("GetTokens: results:", txsResponse);
                setTokens(txsResponse);
            }).catch((errorRes)=>{
                console.log("GetTokens: errors:", errorRes);
            });
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
              <nav class="breadcrumb" aria-label="breadcrumbs">
                <ul>
                  <li>
                    <Link to="/dashboard" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faGauge} />
                      &nbsp;Overview
                    </Link>
                  </li>
                  <li>
                    <Link to="/more" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faEllipsis} />
                      &nbsp;More
                    </Link>
                  </li>
                  <li class="is-active">
                    <Link to="/tokens" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faCubes} />
                      &nbsp;Tokens
                    </Link>
                  </li>
                </ul>
              </nav>

              <nav class="box">
                <div class="columns">
                  <div class="column">
                    <h1 class="title is-4">
                      <FontAwesomeIcon className="fas" icon={faCubes} />
                      &nbsp;Tokens
                    </h1>
                  </div>
                </div>

                {tokens.length === 0 ? <>
                    <section class="hero is-warning is-medium">
                      <div class="hero-body">
                        <p class="title"><FontAwesomeIcon className="fas" icon={faCubes} />&nbsp;No recent tokens</p>
                        <p class="subtitle">This wallet currently does not have any tokens.</p>
                      </div>
                    </section>
                </> : <>
                    <table className="table is-fullwidth is-size-7">
                      <thead>
                        <tr>
                          <th></th>
                          {/*
                          <th>Date</th>
                          <th>Type</th>
                          <th>Coin(s)</th>
                          <th>Sender</th>
                          <th>Receiver</th>
                          */}
                        </tr>
                      </thead>
                      <tbody>
                        {tokens.map((Token) => (
                          <tr key={Token.hash}>
                            <td>TODO</td>

                            {/*
                            <td>{`${new Date(Token.timestamp).toLocaleString()}`}</td>
                            <td>{Token.from === toLower(currentOpenWalletAtAddress) ? "Sent" : "Received"}</td>
                            <td>{Token.value}</td>
                            <td>{Token.from}</td>
                            <td>{Token.to}</td>
                            */}
                          </tr>
                        ))}
                      </tbody>
                    </table>
                </>}

              </nav>
            </section>
          </div>
        </>
    )
}

export default ListTokensView
