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
    faCube,
    faCoins,
    faEllipsis
} from "@fortawesome/free-solid-svg-icons";
import logo from '../../assets/images/CPS-logo-2023-square.webp';
import { useRecoilState } from "recoil";
import { toLower } from "lodash";

import { GetNonFungibleToken } from "../../../wailsjs/go/main/App";
import { currentOpenWalletAtAddressState } from "../../AppState";
import PageLoadingContent from "../Reusable/PageLoadingContent";
import FormRowText from "../Reusable/FormRowText";
import FormRowMetadataAttributesField from "../Reusable/FormRowMetadataAttributesField";
import FormRowYouTubeField from "../Reusable/FormRowYouTubeField";


function TokenDetailView() {
    ////
    //// URL Parameters.
    ////

    const { tokenID } = useParams();

    ////
    //// Global State
    ////

    const [currentOpenWalletAtAddress] = useRecoilState(currentOpenWalletAtAddressState);

    ////
    //// Component states.
    ////

    const [isLoading, setIsLoading] = useState(false);
    const [forceURL, setForceURL] = useState("");
    const [token, setToken] = useState([]);

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

            // Update the GUI to let user know that the operation is under way.
            setIsLoading(true);

            GetNonFungibleToken(parseInt(tokenID)).then((nftokRes)=>{
                console.log("GetNonFungibleToken: nftokRes:", nftokRes);
                setToken(nftokRes);
            }).catch((errorRes)=>{
                console.log("GetNonFungibleToken: errorRes:", errorRes);
            }).finally((errorRes)=>{
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
            <PageLoadingContent displayMessage="Fetching..." />
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
                  <li class="">
                    <Link to="/more/tokens" aria-current="page">
                        <FontAwesomeIcon className="fas" icon={faCubes} />
                        &nbsp;Tokens
                    </Link>
                  </li>
                  <li class="is-active">
                    <Link to={`/more/token/${tokenID}`} aria-current="page">
                        <FontAwesomeIcon className="fas" icon={faCube} />
                        &nbsp;Token ID {tokenID}
                    </Link>
                  </li>
                </ul>
              </nav>

              <nav class="box">
                <div class="columns">
                  <div class="column">
                    <h1 class="title is-4">
                      <FontAwesomeIcon className="fas" icon={faCube} />
                      &nbsp;Token ID {tokenID}
                    </h1>
                  </div>
                </div>
                {token !== undefined && token !== null && token != "" && <>
                    <FormRowText label="Name" value={token.metadata.name} />
                    <FormRowText label="Description" value={token.metadata.description} />
                    <FormRowMetadataAttributesField label="Attributes (Optional)" attributes={token.metadata.attributes} />
                    <FormRowText label="External URL (Optional)" value={token.metadata.external_url} />
                    {/*<FormRowText label="Background Color" value={token.metadata.background_color} />*/}
                    <FormRowYouTubeField label="YouTube URL (Optional)" url={token.metadata.youtube_url} />
                    <figure class="image is-4by3">
                        <img
                            src={token.metadata.image}
                            alt={token.metadata.name}
                        />
                    </figure>
                </>}
              </nav>
            </section>
          </div>
        </>
    )
}

export default TokenDetailView
