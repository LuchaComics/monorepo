import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";
import { WalletMinimal, Send, QrCode, MoreHorizontal, Wallet, CheckCircle2, ArrowRight } from 'lucide-react';
import { useRecoilState } from "recoil";

import { currentOpenWalletAtAddressState } from "../../AppState";


function SendCoinSuccessView() {
    ////
    //// Global State
    ////

    const [currentOpenWalletAtAddress] = useRecoilState(currentOpenWalletAtAddressState);

    ////
    //// Component states.
    ////

    // GUI States.
    const [errors, setErrors] = useState({});
    const [forceURL, setForceURL] = useState("");
    const [isLoading, setIsLoading] = useState(false);

    // Form Submission States.
    ////
    //// Event handling.
    ////

    ////
    //// API.
    ////

    const onSubmitClick = (e) => {
        e.preventDefault();


    }

    ////
    //// Misc.
    ////

    useEffect(() => {
      let mounted = true;

      if (mounted) {
          window.scrollTo(0, 0); // Start the page at the top of the page.
      }

      return () => {
          mounted = false;
      };
    }, []);

    ////
    //// Component rendering.
    ////

    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
      return <Navigate to={forceURL} />;
    }

    return (
        <div>
            <main className="max-w-2xl mx-auto px-6 py-12 mb-24">
                <div className="bg-white rounded-xl border-2 border-gray-100 overflow-hidden">
                 <div className="p-12 text-center">
                   <div className="flex justify-center mb-6">
                     <div className="w-20 h-20 bg-green-100 rounded-full flex items-center justify-center">
                       <CheckCircle2 className="w-12 h-12 text-green-600" />
                     </div>
                   </div>

                   <h2 className="text-2xl font-bold text-gray-900 mb-4">
                     Transaction Submitted!
                   </h2>

                   <p className="text-gray-600 max-w-md mx-auto mb-8">
                     Your coins have been successfully sent to the specified account. Please allow a few minutes for the transaction to be processed on the blockchain.
                   </p>

                   <Link className="inline-flex items-center justify-center gap-2 px-6 py-3 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors" to={`/dashboard`}>
                     Return to Overview
                     <ArrowRight className="w-4 h-4" />
                   </Link>
                 </div>
                </div>
            </main>
       </div>
    );
};

export default SendCoinSuccessView
