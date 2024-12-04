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
  faHome,
  faArrowUpRightFromSquare,
  faUpload,
  faDownload
} from "@fortawesome/free-solid-svg-icons";

import PageLoadingContent from "../Reusable/PageLoadingContent";
import { currentOpenWalletAtAddressState } from "../../AppState";
import {
    ListWallets,
    SetDefaultWalletAddress,
    ExportWalletUsingDialog,
    ImportWalletUsingDialog
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

    const onExportWalletClick = (walletAddress) => {
        console.log("currentOpenWalletAtAddress: Old:", currentOpenWalletAtAddress);
        ExportWalletUsingDialog(walletAddress).then(()=>{
            // // STEP 1: Adjust the wallet which is open.
            // setCurrentOpenWalletAtAddress(walletAddress);
            //
            // // STEP 2:
            // console.log("currentOpenWalletAtAddress: New:", walletAddress);
            //
            // // STEP 3: Redirect to the dashboard page.
            // setForceURL("/dashboard");
        });
    }

    const onImportWalletClick = (e) => {
        e.preventDefault();
            console.log("onImportWalletClick: Beginning...");
        ImportWalletUsingDialog()
            .then(() => {
                console.log("ImportWalletUsingDialog: Successfully imported wallet");
                return ListWallets();
            })
            .then((walletsResponse) => {
                console.log("ImportWalletUsingDialog: walletsResponse:", walletsResponse); // Debug log
                setWallets(walletsResponse);
            })
            .catch((errorJsonString) => {
                console.error("ImportWalletUsingDialog: Error importing wallet:", errorJsonString);
                const errorObject = JSON.parse(errorJsonString);
                let err = {};
                if (errorObject.filepath !== "") {
                    err.filepath = errorObject.filepath;
                }
                setErrors(err);
            })
            .finally(() => {
                console.log("ImportWalletUsingDialog: promise resolved or rejected");
            });
    };

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;
        setIsLoading(true);

        if (mounted) {
            window.scrollTo(0, 0);

            ListWallets()
                .then((walletsResponse) => {
                    console.log("useEffect: walletsResponse:", walletsResponse);
                    setWallets(walletsResponse);
                })
                .catch((error) => {
                    console.error("Error fetching wallets:", error);
                    setWallets([]);
                })
                .finally(() => {
                    setIsLoading(false);
            });
        }

        return () => {
            mounted = false;
        };
    }, []);

    ////
    //// Component rendering.
    ////

    if (isLoading) {
        return <PageLoadingContent displayMessage="Loading..." />;
    }

    if (forceURL !== "") {
        return <Navigate to={forceURL} />;
    }

    return (
        <>
        {wallets.length === 0 ? (
            <section className="hero is-fullheight-with-navbar is-info">
                <div className="hero-body">
                    <div className="container">
                        <div className="columns is-centered">
                            <div className="column is-6-tablet is-5-desktop is-4-widescreen">
                                <h1 className="title is-2">Welcome to ComicCoin Wallet</h1>
                                <h2 className="subtitle is-4">To begin, please create your wallet to get started.</h2>
                                <div className="buttons">
                                    <Link className="button is-link is-medium" to="/wallet/add">
                                        Create your wallet
                                    </Link>
                                    &nbsp;
                                    <Link
                                        className="button is-success is-medium"
                                        onClick={(e)=>onImportWalletClick(e)}
                                    >
                                        Import your wallet
                                    </Link>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </section>
        ) : (
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
                            <th></th>
                          </tr>
                        </thead>
                        <tbody>
                        {wallets.map((wallet) => (
                          <tr key={wallet.filepath}>
                              <td>{wallet.label}</td>
                              <td>{wallet.address}</td>
                              <td>
                                  <Link onClick={(e) => onExportWalletClick(wallet.address)}>
                                      <FontAwesomeIcon className="fas" icon={faDownload} />
                                      &nbsp;Export
                                  </Link>
                              </td>
                              <td>
                                  <Link onClick={(e) => onClick(wallet.address)}>
                                      Open&nbsp;<FontAwesomeIcon className="fas" icon={faChevronRight} />
                                  </Link>
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
                <div className="has-text-right">
                  <Link className="has-text-grey" onClick={(e)=>onImportWalletClick(e)}><FontAwesomeIcon className="fas" icon={faUpload} />&nbsp;Import Wallet from file</Link>
                </div>
              </section>
            </div>
        )}
        </>
    );

}

export default ListWalletsView
