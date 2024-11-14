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
    faEllipsis,
    faChevronRight,
    faPlus
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { toLower } from "lodash";

import { GetTokens } from "../../../wailsjs/go/main/App";


function ListTokensView() {
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

            GetTokens().then((txsResponse)=>{
                console.log("GetTokens: results:", txsResponse);
                setTokens(txsResponse);
            }).catch((errorRes)=>{
                console.log("GetTokens: errors:", errorRes);
            });
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
          <div class="container">
            <section class="section">
              <nav class="breadcrumb" aria-label="breadcrumbs">
                <ul>
                  <li>
                    <Link to="/dashboard" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faGauge} />
                      &nbsp;Dashboard
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

                {tokens !== undefined && tokens !== null && tokens !== "" && tokens.length === 0 ? <>
                    <section class="hero is-warning is-medium">
                      <div class="hero-body">
                        <p class="title"><FontAwesomeIcon className="fas" icon={faCubes} />&nbsp;No recent tokens</p>
                        <p class="subtitle">ComicCoin currently does not have any tokens, to get started creating your first token please <Link to="/tokens/new">click here&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} />.</Link></p>
                      </div>
                    </section>
                </> : <>
                    <table className="table is-fullwidth is-size-7">
                      <thead>
                        <tr>
                          <th>Token ID</th>
                          <th>Name</th>
                          <th>Description</th>
                          <th>Created</th>
                          <th><Link to="/tokens/new"><FontAwesomeIcon className="fas" icon={faPlus} />&nbsp;New Token</Link></th>
                        </tr>
                      </thead>
                      <tbody>
                        {tokens.map((token) => (
                          <tr key={token.token_id}>
                            <td>{token.token_id}</td>
                            <td>{token.metadata.name}</td>
                            <td>{token.metadata.description}</td>
                            <td>{`${new Date(token.timestamp).toLocaleString()}`}</td>
                            <td>
                                <Link to={`/token/${token.token_id}`}>View&nbsp;<FontAwesomeIcon className="fas" icon={faChevronRight} /></Link>
                            </td>
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
