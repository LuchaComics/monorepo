import {useState, useEffect} from 'react';
import { Link, Navigate, useParams } from "react-router-dom";
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
    faEllipsis,
    faChevronRight,
    faArrowLeft,
    faLink,
    faFile
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { toLower } from "lodash";

import FormRowText from "../Reusable/FormRowText";
import {
    GetToken,
    GetNFTStoreRemoteAddressFromPreferences
} from "../../../wailsjs/go/main/App";
import PageLoadingContent from "../Reusable/PageLoadingContent";
import FormRowIPFSImage from "../Reusable/FormRowIPFSImage";
import FormRowIPFSVideo from "../Reusable/FormRowIPFSVideo";
import FormRowIPFSDownloadLink from "../Reusable/FormRowIPFSDownloadLink";
import FormRowYouTubeField from "../Reusable/FormRowYouTubeField";
import FormRowMetadataAttributes from "../Reusable/FormRowMetadataAttributes";


function TokenDetailView() {
    ////
    //// URL Parameters.
    ////

    const { id } = useParams();

    ////
    //// Component states.
    ////

    const [isLoading, setIsLoading] = useState(false);
    const [forceURL, setForceURL] = useState("");
    const [totalCoins, setTotalCoins] = useState(0);
    const [totalTokens, setTotalTokens] = useState(0);
    const [token, setToken] = useState(null);
    const [remoteAddress, setRemoteAddress] = useState("http://127.0.0.1:8080");

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

            GetToken(parseInt(id)).then((tokRes)=>{
                console.log("GetToken: results:", tokRes);
                setToken(tokRes);
            }).catch((errorRes)=>{
                console.log("GetToken: errors:", errorRes);
            }).finally(() => {
                // this will be executed after then or catch has been executed
                console.log("promise has been resolved or rejected");

                // Update the GUI to let user know that the operation is completed.
                setIsLoading(false);
            });

            // Get the remote address.
            GetNFTStoreRemoteAddressFromPreferences().then( (resp)=>{
                console.log("OnStartup: RemoteAddress:", resp);
                setRemoteAddress(resp);
            })
      }

      return () => {
          mounted = false;
      };
    }, [id]);

    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
        return <Navigate to={forceURL} />;
    }

    return (
        <>
          {isLoading ? <>
              <PageLoadingContent displayMessage="Fetching..." />
          </> : <>
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
                    <li>
                      <Link to="/tokens" aria-current="page">
                        <FontAwesomeIcon className="fas" icon={faCubes} />
                        &nbsp;Tokens
                      </Link>
                    </li>
                    <li class="is-active">
                      <Link to={`/token/${id}`} aria-current="page">
                        <FontAwesomeIcon className="fas" icon={faCube} />
                        &nbsp;Token ID {id}
                      </Link>
                    </li>
                  </ul>
                </nav>

                <nav class="box">
                  <div class="columns">
                      <div class="column">
                          <h1 class="title is-4">
                              <FontAwesomeIcon className="fas" icon={faCube} />
                              &nbsp;Token Detail
                          </h1>
                      </div>
                  </div>

                  {token !== undefined && token !== null && token !== "" && <>
                      <h1 class="title is-5">
                          <FontAwesomeIcon className="fas" icon={faLink} />
                          &nbsp;Blockchain
                      </h1>
                      <FormRowText label="ID" value={token.token_id} />
                      {remoteAddress && <>
                          <FormRowIPFSDownloadLink
                               label="Metadata URI"
                               ipfsGatewayDomain={remoteAddress}
                               ipfsPath={`${token.metadata_uri}`}
                          />
                      </>}
                      <h1 class="title is-5">
                          <FontAwesomeIcon className="fas" icon={faFile} />
                          &nbsp;Metadata
                      </h1>
                      <FormRowText label="Name" value={token.metadata.name} />
                      <FormRowText label="Description" value={token.metadata.description} />
                      <FormRowMetadataAttributes label="Attributes (Optional)" attributes={token.metadata.attributes} />
                      <FormRowText label="External URL (Optional)" value={token.metadata.external_url} />
                      <FormRowText label="Background Color" value={token.metadata.background_color} />
                      {remoteAddress && <>
                          <FormRowIPFSImage label="Image" ipfsGatewayDomain={remoteAddress} ipfsPath={token.metadata.image} />
                          <FormRowIPFSVideo label="Animation" ipfsGatewayDomain={remoteAddress} ipfsPath={token.metadata.animation_url} />
                      </>}
                      <FormRowYouTubeField label="YouTube URL (Optional)" url={token.metadata.youtube_url} />
                  </>}

                  <div class="columns pt-5" style={{alignSelf: "flex-start"}}>
                    <div class="column is-half ">
                      <Link
                        class="button is-fullwidth-mobile"
                        to={`/tokens`}
                      >
                        <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                        &nbsp;Back
                      </Link>
                    </div>
                    <div class="column is-half has-text-right">
                      {/*
                      <button
                        class="button is-primary is-fullwidth-mobile"
                        onClick={onSubmitClick}
                      >
                        <FontAwesomeIcon className="fas" icon={faCheckCircle} />
                        &nbsp;Save
                      </button>
                      */}
                    </div>
                  </div>

                </nav>
              </section>
            </div>
          </>}
        </>
    )
}

export default TokenDetailView
