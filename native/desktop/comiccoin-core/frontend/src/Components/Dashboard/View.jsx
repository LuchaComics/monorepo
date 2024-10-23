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
import { GetTotalCoins, GetTotalTokens, GetRecentTransactions } from "../../../wailsjs/go/main/App";
import { currentOpenWalletAtAddressState } from "../../AppState";


function DashboardView() {
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

            GetTotalCoins(currentOpenWalletAtAddress).then((totalCoinsResult)=>{
                console.log("GetTotalCoins: results:", totalCoinsResult);
                setTotalCoins(totalCoinsResult);
            }).catch((errorRes)=>{
                console.log("GetTotalCoins: errors:", errorRes);
            });

            GetTotalTokens(currentOpenWalletAtAddress).then((totalTokensResult)=>{
                console.log("GetTotalTokens: results:", totalTokensResult);
                setTotalTokens(totalTokensResult);
            }).catch((errorRes)=>{
                console.log("GetTotalTokens: errors:", errorRes);
            });

            GetRecentTransactions(currentOpenWalletAtAddress).then((txsResponse)=>{
                console.log("GetRecentTransactions: results:", txsResponse);
                setTransactions(txsResponse);
            }).catch((errorRes)=>{
                console.log("GetRecentTransactions: errors:", errorRes);
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
              <nav class="box">
                <div class="columns">
                  <div class="column">
                    <h1 class="title is-4">
                      <FontAwesomeIcon className="fas" icon={faGauge} />
                      &nbsp;Overview
                    </h1>
                  </div>
                </div>

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

            <h1 className="subtitle is-4 pt-5"><FontAwesomeIcon className="fas" icon={faFileInvoiceDollar} />&nbsp;Recent Transactions</h1>
            {transactions.length === 0 ? <>
                <section class="hero is-warning is-medium">
                  <div class="hero-body">
                    <p class="title"><FontAwesomeIcon className="fas" icon={faFileInvoiceDollar} />&nbsp;No recent transactions</p>
                    <p class="subtitle">This wallet currently does not have any transactions.</p>
                  </div>
                </section>
            </> : <>
                <table className="table is-fullwidth is-size-7">
                  <thead>
                    <tr>
                      <th></th>
                      <th>Date</th>
                      <th>Type</th>
                      <th>Coin(s)</th>
                      <th>Sender</th>
                      <th>Receiver</th>
                    </tr>
                  </thead>
                  <tbody>
                    {transactions.map((transaction) => (
                      <tr key={transaction.hash}>
                        <td>{transaction.type === "coin" ? <><FontAwesomeIcon className="fas" icon={faCoins} /></> : <><FontAwesomeIcon className="fas" icon={faCubes} /></>}</td>
                        <td>{transaction.timestamp}</td>
                        <td>{transaction.from === currentOpenWalletAtAddress ? "Sent" : "Received"}</td>
                        <td>{transaction.value}</td>
                        <td>{transaction.from}</td>
                        <td>{transaction.to}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
                <div className="has-text-right">
                    <Link to={`/transactions`}>See More&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} /></Link>
                </div>
            </>}
            </nav>
            </section>
          </div>
        </>
    )
}

export default DashboardView
