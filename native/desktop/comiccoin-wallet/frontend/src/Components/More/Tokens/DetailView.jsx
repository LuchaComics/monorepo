import {useState, useEffect} from 'react';
import { Link, useParams } from "react-router-dom";
import { useRecoilState } from "recoil";
import { toLower } from "lodash";

import { GetNonFungibleToken } from "../../../../wailsjs/go/main/App";
import { currentOpenWalletAtAddressState } from "../../../AppState";


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

    // if (isLoading) {
    //     return (
    //         "------"
    //     );
    // }

 return (
   <div>
    {token && token.metadata && <div className="min-h-screen bg-gray-100">
      <div className="container mx-auto px-4 py-6">
        {/* Token Image and Basic Info */}
        <div className="bg-white rounded-lg shadow-md p-6 mb-6">
          <div className="flex flex-col md:flex-row gap-6">
            <div className="w-full md:w-1/2">
              <img
                src={token.metadata.image}
                alt={token.metadata.name}
                className="w-full rounded-lg"
              />
            </div>
            <div className="w-full md:w-1/2">
              <h2 className="text-3xl font-bold mb-4">{token.metadata.name}</h2>
              <p className="text-gray-600 mb-6">{token.metadata.description}</p>

              {/* Action Buttons */}
              <div className="flex gap-4 mb-6">
                <Link className="bg-red-500 text-white px-6 py-2 rounded-lg hover:bg-red-600" to={`/more/token/${tokenID}/burn`}>
                  Burn
                </Link>
                <Link className="bg-blue-500 text-white px-6 py-2 rounded-lg hover:bg-blue-600" to={`/more/token/${tokenID}/transfer`}>
                  Transfer to Another Address
                </Link>
              </div>

              {/* External Links */}
              <div className="space-y-2">
                <p className="text-sm">
                  <span className="font-semibold">External URL: </span>
                  <a href={token.metadata.external_url} className="text-blue-500 hover:underline">
                    {token.metadata.external_url}
                  </a>
                </p>
                <p className="text-sm">
                  <span className="font-semibold">Background Color: </span>
                  {token.metadata.background_color}
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Animation Video */}
        {token.metadata.animation_url && (
          <div className="bg-white rounded-lg shadow-md p-6 mb-6">
            <h3 className="text-xl font-bold mb-4">Animation</h3>
            <div className="flex justify-center">
              <video
                controls
                className="rounded-lg max-w-full"
                preload="metadata"
                style={{ maxHeight: '500px' }}
              >
                <source src={token.metadata.animation_url} type="video/webm" />
                Your browser does not support the video tag.
              </video>
            </div>
          </div>
        )}

        {/* YouTube Video */}
        {token.metadata.youtube_url && (
          <div className="bg-white rounded-lg shadow-md p-6 mb-6">
            <h3 className="text-xl font-bold mb-4">YouTube Content</h3>
            <iframe
              width="420"
              height="315"
              src={token.metadata.youtube_url}
              className="mx-auto"
            />
          </div>
        )}

        {/* Attributes */}
        <div className="bg-white rounded-lg shadow-md p-6 mb-20">
          <h3 className="text-xl font-bold mb-4">Attributes</h3>
          {token.metadata.attributes !== undefined &&
            token.metadata.attributes !== null &&
            token.metadata.attributes !== "" &&
            token.metadata.attributes.length > 0 ? (
              <div className="overflow-x-auto">
                <table className="w-full">
                  <thead className="text-sm">
                    <tr className="bg-gray-50">
                      <th className="px-4 py-2 text-left">Display Type</th>
                      <th className="px-4 py-2 text-left">Trait Type</th>
                      <th className="px-4 py-2 text-left">Value</th>
                    </tr>
                  </thead>
                  <tbody className="text-sm">
                    {token.metadata.attributes.map((datum, i) => (
                      <tr key={i} className="border-t">
                        <td className="px-4 py-2">{datum.display_type}</td>
                        <td className="px-4 py-2">{datum.trait_type}</td>
                        <td className="px-4 py-2">{datum.value}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            ) : null}
        </div>
      </div>

      {/* Bottom Tab Navigation */}
      <div className="fixed bottom-0 left-0 right-0 bg-white shadow-lg">
        <div className="flex justify-around items-center">
          <div className="px-4 py-3 text-gray-500">Overview</div>
          <div className="px-4 py-3 text-gray-500">Send</div>
          <div className="px-4 py-3 text-gray-500">Receive</div>
          <div className="px-4 py-3 text-blue-500 border-t-2 border-blue-500">More</div>
        </div>
      </div>
    </div>
    }
   </div>
 );
};


export default TokenDetailView
