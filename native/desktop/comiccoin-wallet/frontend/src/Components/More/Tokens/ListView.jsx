import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";
import { WalletMinimal, Send, QrCode, MoreHorizontal, Clock, Coins, Wallet, ArrowRight, ArrowUpRight, ArrowDownLeft, Tickets } from 'lucide-react';

import { useRecoilState } from "recoil";
import { toLower } from "lodash";

import { GetNonFungibleTokensByOwnerAddress } from "../../../../wailsjs/go/main/App";
import { currentOpenWalletAtAddressState } from "../../../AppState";


function ListTokensView() {
    ////
    //// Global State
    ////

    const [currentOpenWalletAtAddress] = useRecoilState(currentOpenWalletAtAddressState);

    ////
    //// Component states.
    ////

    const [isLoading, setIsLoading] = useState(false);
    const [forceURL, setForceURL] = useState("");
    const [tokens, setTokens] = useState([]);

    ////
    //// Event handling.
    ////

    const handleTokenClick = (tokID) => {
        console.log("tokID->", tokID);
        setForceURL("/more/token/"+tokID);
    }

    ////
    //// Misc.
    ////

    useEffect(() => {
      let mounted = true;

      if (mounted) {
            window.scrollTo(0, 0); // Start the page at the top of the page.

            // Update the GUI to let user know that the operation is under way.
            setIsLoading(true);

            GetNonFungibleTokensByOwnerAddress(currentOpenWalletAtAddress).then((nftoksRes)=>{
                console.log("GetNonFungibleTokensByOwnerAddress: nftoksRes:", nftoksRes);
                setTokens(nftoksRes);
            }).catch((errorRes)=>{
                console.log("GetNonFungibleTokensByOwnerAddress: errorRes:", errorRes);
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
            "---"
        );
    }

    if (tokens.length === 0) {
        return (
            <div>
                <main className="max-w-2xl mx-auto px-6 py-12 mb-24">
                <div className="bg-white rounded-xl shadow-sm border border-gray-200">
                 <div className="p-6 border-b border-gray-100">
                   <div className="flex justify-between items-center">
                     <div className="flex items-center gap-3">
                       <div className="p-2 bg-purple-100 rounded-xl">
                         <Tickets className="w-5 h-5 text-purple-600" aria-hidden="true" />
                       </div>
                       <h2 className="text-xl font-bold text-gray-900">Tokens</h2>
                     </div>
                   </div>
                 </div>

                 {/* Empty State Message */}
                 <div className="py-16 px-6">
                   <div className="text-center">
                     <div className="w-16 h-16 mx-auto mb-4 bg-purple-100 rounded-full flex items-center justify-center">
                       <Tickets className="w-8 h-8 text-purple-600" aria-hidden="true" />
                     </div>
                     <h3 className="text-lg font-medium text-gray-900 mb-2">No Tokens Yet</h3>
                     <p className="text-gray-500 max-w-sm mx-auto">
                       Start your journey by sending or receiving ComicCoins or NFTs. Your transaction history will appear here.
                     </p>
                   </div>
                 </div>
               </div>
               </main>
            </div>
        );
    }

    return (
        <>
            <div className="container mx-auto px-4 py-8">
                 <h1 className="text-3xl font-bold mb-6">My NFT Collection</h1>

                 <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                   {tokens.map((token) => (
                     <div
                       key={token.token_id}
                       className="bg-white rounded-lg shadow-md hover:shadow-lg transition-shadow duration-200 cursor-pointer"
                       onClick={() => handleTokenClick(token.token_id)}
                     >
                       <div className="p-4">
                         <div className="aspect-square relative mb-4">
                           <img
                             src={`/${token.metadata.image}`}
                             alt={token.metadata.name}
                             className="rounded-lg object-cover w-full h-full"
                           />
                         </div>
                         <div className="space-y-2">
                           <h2 className="text-xl font-semibold">{token.metadata.name}</h2>
                           <p className="text-sm text-gray-500">Token ID: {token.token_id}</p>
                         </div>
                       </div>
                     </div>
                   ))}
                 </div>

                 {tokens.length === 0 && (
                   <div className="text-center py-12">
                     <p className="text-gray-500 text-lg">No NFTs found in your wallet</p>
                   </div>
                 )}
            </div>
        </>
    )
}

export default ListTokensView
