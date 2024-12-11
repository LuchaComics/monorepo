import { useState, useEffect } from "react";
import { Link, useParams, Navigate } from "react-router-dom";
import { useRecoilState } from "recoil";
import { toLower } from "lodash";
import { AlertCircle, Info, Send, ArrowLeft } from "lucide-react";

import {
  GetNonFungibleToken,
  TransferToken,
} from "../../../../wailsjs/go/main/App";
import { currentOpenWalletAtAddressState } from "../../../AppState";

/**
 * TokenTransferView component allows users to transfer a specific non-fungible token (NFT) to another wallet address.
 *
 * It fetches NFT details using the tokenID from the URL parameters and the currently open wallet address from Recoil state.
 * The component handles form input validation, password-protected transaction authorization using TransferToken API call,
 * and provides visual feedback during the transfer process (loading indicator, error messages).
 * Success or failure is indicated through URL redirection.
 */
function TokenTransferView() {
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

  const [isLoading, setIsLoading] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [token, setToken] = useState([]);
  const [errors, setErrors] = useState({});
  const [showErrorBox, setShowErrorBox] = useState(false);
  const [formData, setFormData] = useState({
    transferTo: "",
    password: "",
  });
  const currentBalance = 105.0;

  ////
  //// Event handling.
  ////

  const validateForm = () => {
    const newErrors = {};

    if (!formData.transferTo.trim()) {
      newErrors.transferTo = "Transfer address is required";
    } else if (!/^[A-Za-z0-9]{42}$/.test(formData.transferTo)) {
      newErrors.transferTo = "Invalid wallet address format";
    }

    if (!formData.password) {
      newErrors.password = "Password is required to authorize transaction";
    }

    setErrors(newErrors);
    setShowErrorBox(Object.keys(newErrors).length > 0);
    return Object.keys(newErrors).length === 0;
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
    if (errors[name]) {
      setErrors((prev) => ({ ...prev, [name]: "" }));
      setShowErrorBox(false);
    }
  };

  // Handles the form submission. Validates the form data, then calls the TransferToken API to transfer the NFT.
  // On success, redirects to a success page; on failure, displays errors.  Manages loading state and error handling.
  const handleSubmit = (e) => {
    e.preventDefault();
    if (validateForm()) {
      // Form is valid, proceed with submission
      console.log("Form submitted:", formData);

      // Update the GUI to let user know that the operation is under way.
      setIsLoading(true);
      setShowErrorBox(false);

      TransferToken(
        formData.transferTo,
        parseInt(tokenID),
        currentOpenWalletAtAddress,
        formData.password,
      )
        .then(() => {
          console.log("TransferToken: Success");
          setForceURL("/more/token/" + tokenID + "/transfer-success");
        })
        .catch((errorJsonString) => {
          console.log("TransferToken: errorJsonString:", errorJsonString);

          let err = {};
          if (
            errorJsonString
              .toLowerCase()
              .includes("could not decrypt key with given password")
          ) {
            err.password = "Wrong password";
          }
          console.log("TransferToken: err:", err);
          window.scrollTo(0, 0); // Start the page at the top of the page.
          setShowErrorBox(true);
          setErrors(err);
        })
        .finally((errorRes) => {
          // Update the GUI to let user know that the operation is completed.
          setIsLoading(false);
        });
    } else {
      // Show error message
      setShowErrorBox(true);
      // Auto-hide error after 5 seconds
      setTimeout(() => setShowErrorBox(false), 5000);
    }
  };

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.

      // Update the GUI to let user know that the operation is under way.
      setIsLoading(true);

      GetNonFungibleToken(parseInt(tokenID))
        .then((nftokRes) => {
          console.log("GetNonFungibleToken: nftokRes:", nftokRes);
          setToken(nftokRes);
        })
        .catch((errorRes) => {
          console.log("GetNonFungibleToken: errorRes:", errorRes);
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
    console.log("URL redirecting to:", forceURL);
    return <Navigate to={forceURL} />;
  }

  if (isLoading) {
    return "------";
  }

  console.log("showErrorBox:", showErrorBox);
  console.log("errors:", errors);

  return (
    <div>
      <main className="max-w-2xl mx-auto px-6 py-12 mb-24">
        {showErrorBox && (
          <div className="mb-6 bg-red-50 border-l-4 border-red-500 rounded-r-lg p-4 flex items-start gap-3">
            <AlertCircle className="w-5 h-5 text-red-500 flex-shrink-0 mt-0.5" />
            <div className="flex-grow">
              <h3 className="font-semibold text-red-800">
                Unable to Transfer Token
              </h3>
              <div className="text-sm text-red-600 mt-1 space-y-1">
                {Object.values(errors).map((error, index) => (
                  <p key={index}>• {error}</p>
                ))}
              </div>
            </div>
          </div>
        )}

        <div className="bg-white rounded-xl border-2 border-gray-100">
          <div className="p-6">
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-yellow-100 rounded-xl">
                  <Send
                    className="w-5 h-5 text-yellow-600"
                    aria-hidden="true"
                  />
                </div>
                <h2 className="text-xl font-bold text-gray-900">
                  Transfer Token
                </h2>
              </div>
            </div>
            <p className="text-sm text-gray-500">
              Transfer NFT to another address. Please fill in all required
              fields.
            </p>
          </div>

          <div className="p-6 space-y-8">
            <div className="bg-amber-50 border border-amber-100 rounded-xl p-4 flex gap-3">
              <AlertCircle className="w-5 h-5 text-amber-600 flex-shrink-0 mt-0.5" />
              <div className="text-sm text-amber-800">
                <p className="font-semibold mb-1">Important Notice</p>
                <p>
                  All transactions are final and cannot be undone. Please verify
                  all details before transfering.
                </p>
              </div>
            </div>

            <div className="space-y-6">
              <label className="block">
                <span className="text-sm font-medium text-gray-700">
                  Transfer To <span className="text-red-500">*</span>
                </span>
                <input
                  type="text"
                  name="transferTo"
                  value={formData.transferTo}
                  onChange={handleInputChange}
                  className={`mt-1 block w-full px-4 py-3 bg-white border rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent transition-colors ${
                    errors.transferTo
                      ? "border-red-300 bg-red-50"
                      : "border-gray-200"
                  }`}
                  placeholder="Enter other address"
                />
                {errors.transferTo && (
                  <p className="mt-2 text-sm text-red-600 flex items-center gap-2">
                    <AlertCircle className="w-4 h-4" />
                    {errors.transferTo}
                  </p>
                )}
              </label>

              <label className="block">
                <div className="flex items-center justify-between mb-1">
                  <span className="text-sm font-medium text-gray-700">
                    Wallet Password <span className="text-red-500">*</span>
                  </span>
                </div>
                <input
                  type="password"
                  name="password"
                  value={formData.password}
                  onChange={handleInputChange}
                  className={`mt-1 block w-full px-4 py-3 bg-white border rounded-lg focus:ring-2 focus:ring-red-500 focus:border-transparent transition-colors ${
                    errors.password
                      ? "border-red-300 bg-red-50"
                      : "border-gray-200"
                  }`}
                  placeholder="Enter your wallet password"
                />
                {errors.password && (
                  <p className="mt-2 text-sm text-red-600 flex items-center gap-2">
                    <AlertCircle className="w-4 h-4" />
                    {errors.password}
                  </p>
                )}
                <p className="mt-2 text-sm text-gray-600 flex items-start gap-2">
                  <Info className="w-4 h-4 flex-shrink-0 mt-0.5" />
                  <span>
                    Your wallet is encrypted and stored locally. The password is
                    required to authorize this transaction.
                  </span>
                </p>
              </label>

              <div className="bg-blue-50 border border-blue-100 rounded-xl p-4 flex gap-3">
                <Info className="w-5 h-5 text-blue-600 flex-shrink-0 mt-0.5" />
                <div className="text-sm text-blue-800">
                  <p>
                    A network fee of 1 ComicCoin is required to support the
                    blockchain infrastructure and ensure secure transaction
                    processing.
                  </p>
                </div>
              </div>
            </div>

            <div className="flex gap-4 justify-end pt-4">
              <button
                onClick={() => setForceURL("/more/token/" + tokenID)}
                className="px-6 py-2.5 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors inline-flex items-center gap-2"
              >
                <ArrowLeft className="w-4 h-4" />
                Cancel
              </button>
              <button
                onClick={handleSubmit}
                className="px-6 py-2.5 bg-yellow-600 text-white rounded-lg hover:bg-yellow-700 transition-colors inline-flex items-center gap-2"
              >
                Transfer Token
                <Send className="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}

export default TokenTransferView;
