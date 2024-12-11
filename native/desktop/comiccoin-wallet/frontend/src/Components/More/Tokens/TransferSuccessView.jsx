import { useState, useEffect } from "react";
import { Link, useParams, Navigate } from "react-router-dom";
import { useRecoilState } from "recoil";
import { toLower } from "lodash";
import { CheckCircle2, Wallet, Clock, Info } from "lucide-react";

import { GetTotalCoins } from "../../../../wailsjs/go/main/App";
import { currentOpenWalletAtAddressState } from "../../../AppState";

function TokenTransferSuccessView() {
  ////
  //// URL Parameters.
  ////

  const { tokenID } = useParams();

  ////
  //// Global State
  ////

  const [currentOpenWalletAtAddress] = useRecoilState(
    currentOpenWalletAtAddressState,
  );

  ////
  //// Component states.
  ////

  // GUI States.
  const [isLoading, setIsLoading] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [totalCoins, setTotalCoins] = useState(0);
  const [errors, setErrors] = useState({});

  // Form Submission States.
  const [transferTo, setTransferTo] = useState("");
  const [message, setMessage] = useState("");
  const [walletPassword, setWalletPassword] = useState("");

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

      GetTotalCoins(currentOpenWalletAtAddress)
        .then((totalCoinsResult) => {
          console.log("GetTotalCoins: results:", totalCoinsResult);
          setTotalCoins(totalCoinsResult);
        })
        .catch((errorRes) => {
          console.log("GetTotalCoins: errors:", errorRes);
          if (errorRes.includes("address is null")) {
            setForceURL("/wallets");
          }
        })
        .finally((errorRes) => {
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
    return "----";
  }

  const transactionDetails = {
    oldBalance: totalCoins + 1,
    newBalance: totalCoins,
    fee: "1",
  };

  return (
    <div>
      <main className="max-w-2xl mx-auto px-6 py-12 mb-24">
        <div className="bg-white rounded-xl border-2 border-gray-100">
          <div className="p-6">
            <div className="flex flex-col items-center text-center">
              <div className="p-3 bg-green-100 rounded-full mb-4">
                <CheckCircle2
                  className="w-8 h-8 text-green-600"
                  aria-hidden="true"
                />
              </div>
              <h2 className="text-2xl font-bold text-gray-900 mb-2">
                Transaction Successful!
              </h2>
              <p className="text-sm text-gray-500">
                Your token has been successfully transferred.
              </p>
            </div>
          </div>

          <div className="p-6 space-y-8">
            {/* Processing Notice */}
            <div className="bg-blue-50 border border-blue-100 rounded-xl p-4 flex gap-3">
              <Clock className="w-5 h-5 text-blue-600 flex-shrink-0 mt-0.5" />
              <div className="text-sm text-blue-800">
                <p className="font-semibold mb-1">Processing Time</p>
                <p>
                  Please allow a few minutes for the blockchain network to
                  process and confirm your transaction.
                </p>
              </div>
            </div>

            {/* Balance Information */}
            <div className="bg-gray-50 rounded-xl p-6">
              <h3 className="text-sm font-medium text-gray-700 mb-4 flex items-center gap-2">
                <Wallet className="w-4 h-4" />
                Balance Update
              </h3>

              <div className="space-y-3">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-600">Previous Balance</span>
                  <span className="font-medium">
                    {transactionDetails.oldBalance} CC
                  </span>
                </div>
                <div className="flex justify-between text-sm text-red-600">
                  <span>Network Fee</span>
                  <span>- {transactionDetails.fee} CC</span>
                </div>
                <div className="flex justify-between text-sm font-medium pt-3 border-t border-gray-200">
                  <span>New Balance</span>
                  <span>{transactionDetails.newBalance} CC</span>
                </div>
              </div>
            </div>

            {/* Fee Information */}
            <div className="flex gap-3 items-start">
              <Info className="w-5 h-5 text-gray-400 flex-shrink-0 mt-0.5" />
              <p className="text-sm text-gray-600">
                Transaction fees support the blockchain network's security and
                processing capabilities.
              </p>
            </div>

            {/* Return Button */}
            <div className="flex justify-center pt-4">
              <Link
                className="px-6 py-2.5 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors"
                to="/dashboard"
              >
                Return to Overview
              </Link>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}

export default TokenTransferSuccessView;
