import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";
import { useRecoilState } from "recoil";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faWallet,
  faArrowRight,
  faUsers,
  faBarcode,
  faCubes,
  faFileInvoiceDollar,
  faCoins,
  faPlus,
  faChevronRight,
  faChevronLeft,
  faLock,
  faHome
} from "@fortawesome/free-solid-svg-icons";

import PageLoadingContent from "../Reusable/PageLoadingContent";
import { currentOpenWalletAtAddressState } from "../../AppState";
import {
    ListWallets,
    SetDefaultWalletAddress
} from "../../../wailsjs/go/main/App";


function ListWalletsView() {
    ////
    //// Global State
    ////

    const [currentOpenWalletAtAddress, setCurrentOpenWalletAtAddress] = useRecoilState(currentOpenWalletAtAddressState);

    ////
    //// Component states.
    ////

    const [isLoading, setIsLoading] = useState(false);
    const [wallets, setWallets] = useState([]);
    const [errors, setErrors] = useState({});
    const [forceURL, setForceURL] = useState("");

    ////
    //// Event handling.
    ////

    const onClick = (walletAddress) => {
        console.log("currentOpenWalletAtAddress: Old:", currentOpenWalletAtAddress);
        SetDefaultWalletAddress(walletAddress).then(()=>{
            // STEP 1: Adjust the wallet which is open.
            setCurrentOpenWalletAtAddress(walletAddress);

            // STEP 2:
            console.log("currentOpenWalletAtAddress: New:", walletAddress);

            // STEP 3: Redirect to the dashboard page.
            setForceURL("/dashboard");
        });
    }

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
            {isLoading
            ? <>
                <PageLoadingContent displayMessage="Loading..." />
            </> : <>
                {wallets.length === 0 ? <>
                    <section class="hero is-fullheight-with-navbar is-info">
                      <div class="hero-body">
                        <div class="container">
                          <div class="columns is-centered">
                            <div class="column is-6-tablet is-5-desktop is-4-widescreen">
                              <h1 class="title is-2">Welcome to ComicCoin Wallet</h1>
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
                    <div class="container">
                      <section class="section">
                        <nav class="breadcrumb" aria-label="breadcrumbs">
                            <ul>
                              <li>
                                <Link to="/more" aria-current="page">
                                  <FontAwesomeIcon className="fas" icon={faChevronLeft} />
                                  &nbsp;Back to More
                                </Link>
                              </li>
                            </ul>
                        </nav>
                        <nav class="box">
                          <div class="columns">
                            <div class="column">
                              <h1 class="title is-2">
                                <FontAwesomeIcon className="fas" icon={faWallet} />
                                &nbsp;Wallets
                              </h1>
                              <p class="subtitle has-text-grey">These are all wallets currently residing on your local computer. Pick from any to continue:</p>

                            </div>
                          </div>

                          <div className="has-background-white-ter is-round p-3">

                              <table className="table is-fullwidth is-size-6 has-background-white-ter">
                                <thead>
                                  <tr>
                                    <th>Label</th>
                                    <th>Address</th>
                                    <th></th>
                                  </tr>
                                </thead>
                                <tbody>
                                  {wallets.map((wallet) => (
                                    <tr key={wallet.filepath}>
                                      <td>{wallet.label}</td>
                                      <td>{wallet.address}</td>
                                      <td>
                                          <Link onClick={(e)=>onClick(wallet.address)}>Open Wallet&nbsp;<FontAwesomeIcon className="fas" icon={faChevronRight} /></Link>
                                      </td>
                                    </tr>
                                  ))}
                                </tbody>
                              </table>
                          </div>
                          <br />
                          <p className="is-size-7 has-text-grey"><b><FontAwesomeIcon className="fas" icon={faLock} />&nbsp;Secure Storage</b>: Your wallet is stored encrypted at rest to protect your coins and tokens. To access your wallet and perform transactions, you'll need to enter the password you created when you set up your wallet.</p>
                          <br />
                          <p className="is-size-7 has-text-grey"><b><FontAwesomeIcon className="fas" icon={faHome} />&nbsp;Local First</b>: Your wallet is stored exclusively on your computer and not in the cloud.</p>
                          <br />
                          <div className="has-text-right">
                            <Link className="button is-primary is-fullwidth-mobile" to={`/wallet/add`}><FontAwesomeIcon className="fas" icon={faPlus} />&nbsp;New Wallet</Link>
                          </div>
                      </nav>
                      </section>
                    </div>
                </>}
            </>}
        </>
    )
}

export default ListWalletsView
