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
    faArrowLeft
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { toLower } from "lodash";

import { GetTokens } from "../../../wailsjs/go/main/App";
import FormRowText from "../Reusable/FormRowText";
import { GetToken } from "../../../wailsjs/go/main/App";
import PageLoadingContent from "../Reusable/PageLoadingContent";
import FormRowIPFSImageField from "../Reusable/FormRowIPFSImageField";
import FormRowIPFSVideoField from "../Reusable/FormRowIPFSVideoField";
import FormRowYouTubeField from "../Reusable/FormRowYouTubeField";
import FormRowMetadataAttributesField from "../Reusable/FormRowMetadataAttributesField";

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
                      <FormRowText label="ID" value={token.token_id} />
                      <FormRowText label="Metadata URI" value={token.metadata_uri} />
                      <FormRowText label="Name" value={token.metadata.name} />
                      <FormRowText label="Description" value={token.metadata.description} />
                      <FormRowMetadataAttributesField label="Attributes" attributes={token.metadata.attributes} />
                      <FormRowText label="External URL" value={token.metadata.external_url} />
                      <FormRowText label="Background Color" value={token.metadata.background_color} />
                      <FormRowIPFSImageField label="Image" ipfsPath={token.metadata.image} />
                      <FormRowIPFSVideoField label="Animation" ipfsPath={token.metadata.animation_url} />
                      <FormRowYouTubeField label="YouTube URL" url={token.metadata.youtube_url} />
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
