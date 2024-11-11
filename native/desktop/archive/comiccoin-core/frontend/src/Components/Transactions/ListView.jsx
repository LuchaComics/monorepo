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
    faFileInvoiceDollar,
    faCoins,
    faEllipsis,
    faChevronRight
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { toLower } from "lodash";

import PageLoadingContent from "../Reusable/PageLoadingContent";
import {
    GetTransactions
} from "../../../wailsjs/go/main/App";
import { currentOpenWalletAtAddressState } from "../../AppState";


function ListTransactionsView() {
    ////
    //// Global State
    ////

    const [currentOpenWalletAtAddress] = useRecoilState(currentOpenWalletAtAddressState);

    ////
    //// Component states.
    ////

    const [isLoading, setIsLoading] = useState(false);
    const [forceURL, setForceURL] = useState("");
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

            setIsLoading(true);
            GetTransactions(currentOpenWalletAtAddress).then((txsResponse)=>{
                console.log("GetTransactions: results:", txsResponse);
                setTransactions(txsResponse);
            }).catch((errorRes)=>{
                console.log("GetTransactions: errors:", errorRes);
            }).finally(()=>{
                setIsLoading(false);
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

    if (isLoading) {
        return (
            <PageLoadingContent displayMessage="Loading..." />
        );
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
                    <Link to="/transactions" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faFileInvoiceDollar} />
                      &nbsp;Transactions
                    </Link>
                  </li>
                </ul>
              </nav>

              <nav class="box">
                <div class="columns">
                  <div class="column">
                    <h1 class="title is-4">
                      <FontAwesomeIcon className="fas" icon={faFileInvoiceDollar} />
                      &nbsp;Transactions
                    </h1>
                  </div>
                </div>

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
                          <th></th>
                        </tr>
                      </thead>
                      <tbody>
                        {transactions.map((transaction) => (
                          <tr key={transaction.timestamp}>
                            <td>{transaction.type === "coin" ? <><FontAwesomeIcon className="fas" icon={faCoins} /></> : <><FontAwesomeIcon className="fas" icon={faCubes} /></>}</td>
                            <td>{`${new Date(transaction.timestamp).toLocaleString()}`}</td>
                            <td>{transaction.from === toLower(currentOpenWalletAtAddress) ? "Sent" : "Received"}</td>
                            <td>{transaction.type === "coin" ? <>{transaction.value}</> : <>-</>}</td>
                            <td>{transaction.from}</td>
                            <td>{transaction.to}</td>
                            <td><Link to={`/more/transaction/${transaction.timestamp}`}><FontAwesomeIcon className="fas" icon={faChevronRight} /></Link></td>
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

export default ListTransactionsView
