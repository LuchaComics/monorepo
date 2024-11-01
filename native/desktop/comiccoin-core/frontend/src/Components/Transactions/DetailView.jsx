import {useState, useEffect} from 'react';
import { Link, useParams } from "react-router-dom";
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
    faChevronRight,
    faEye,
    faLink,
    faBullhorn
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { toLower } from "lodash";

import PageLoadingContent from "../Reusable/PageLoadingContent";
import { GetBlockDataByBlockTransactionTimestamp } from "../../../wailsjs/go/main/App";
import { currentOpenWalletAtAddressState } from "../../AppState";
import FormRowText from "../Reusable/FormRowText";


function TransactionDetailView() {
    ////
    //// URL Parameters.
    ////

    const { timestamp } = useParams();

    ////
    //// Global State
    ////

    const [currentOpenWalletAtAddress] = useRecoilState(currentOpenWalletAtAddressState);

    ////
    //// Component states.
    ////

    const [isLoading, setIsLoading] = useState(false);
    const [forceURL, setForceURL] = useState("");
    const [blockData, setBlockData] = useState(null);

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

            GetBlockDataByBlockTransactionTimestamp(parseInt(timestamp)).then((res)=>{
                console.log("GetBlockDataByBlockTransactionTimestamp: res:", res);
                setBlockData(res);
            }).catch((errorRes)=>{
                console.log("GetBlockDataByBlockTransactionTimestamp: errors:", errorRes);
            }).finally(() => {
                // Update the GUI to let user know that the operation is completed.
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
            <PageLoadingContent displayMessage="Please wait..." style={{ marginBottom: "100px" }} />
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
                  <li>
                    <Link to="/more/transactions" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faFileInvoiceDollar} />
                      &nbsp;Transactions
                    </Link>
                  </li>
                  <li class="is-active">
                    <Link to="/transactions" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faEye} />
                      &nbsp;Detail
                    </Link>
                  </li>
                </ul>
              </nav>

              <nav class="box">
                <div class="columns">
                  <div class="column">
                    <h1 class="title is-4">
                      <FontAwesomeIcon className="fas" icon={faEye} />
                      &nbsp;Transaction Detail
                    </h1>
                  </div>
                </div>

                {blockData !== undefined && blockData !== null && blockData !== "" && <>
                    {blockData.trans.map((transaction) => (
                        <>
                            {transaction.type === "token" && <>
                                <article class="message is-primary">
                                  <div class="message-header">
                                    <p><FontAwesomeIcon className="fas" icon={faBullhorn} />&nbsp;Non-Fungible Token (NFT) Detected</p>
                                  </div>
                                  <div class="message-body">
                                    This is a NFT! To view the contents of it, please <Link to={`/more/token/${transaction.token_id}`}>click here&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} /></Link>.
                                  </div>
                                </article>
                            </>}
                        </>
                    ))}

                    <h1 class="title is-5">
                        <FontAwesomeIcon className="fas" icon={faLink} />
                        &nbsp;Block Information
                    </h1>
                    <FormRowText label="ID" value={blockData.hash} />
                    <FormRowText label="Number" value={blockData.header.number} />

                    <h1 class="title is-5">
                        <FontAwesomeIcon className="fas" icon={faFileInvoiceDollar} />
                        &nbsp;Transaction Information
                    </h1>
                    {blockData.trans.map((transaction) => (
                      <div key={transaction.timestamp}>
                        <FormRowText label="Purpose" value={transaction.type === "coin" ? "Coin" : "Token"} />
                        <FormRowText label="Type" value={transaction.from === toLower(currentOpenWalletAtAddress) ? "Sent" : "Received"} />
                        <FormRowText label="Timestamp" value={`${new Date(transaction.timestamp).toLocaleString()}`} />
                        {transaction.type === "coin" ? <>
                            <FormRowText label="Value" value={transaction.value} />
                        </> : <>

                        </>}
                            <FormRowText label="From" value={transaction.from} />
                            <FormRowText label="To" value={transaction.to} />
                       </div>
                    ))}

                </>}

              </nav>
            </section>
          </div>
        </>
    )
}

export default TransactionDetailView
