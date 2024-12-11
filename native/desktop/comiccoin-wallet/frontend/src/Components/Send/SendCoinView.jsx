import { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import {
  WalletMinimal,
  Send,
  QrCode,
  MoreHorizontal,
  Wallet,
  AlertCircle,
  Info,
  Coins,
} from "lucide-react";
import { useRecoilState } from "recoil";

import { TransferCoin } from "../../../wailsjs/go/main/App";
import { currentOpenWalletAtAddressState } from "../../AppState";

function SendCoinView() {
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
  const [showErrorBox, setShowErrorBox] = useState(false);
  const [errors, setErrors] = useState({});
  const [forceURL, setForceURL] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const currentBalance = 1337.42;

  // Form Submission States.
  const [formData, setFormData] = useState({
    recipient: "",
    amount: "",
    message: "",
    password: "",
  });

  ////
  //// Event handling.
  ////

  ////
  //// API.
  ////

  const onSubmitClick = (e) => {
    e.preventDefault();

    // Update the GUI to let user know that the operation is under way.
    setIsLoading(true);

    TransferCoin(
      payTo,
      parseInt(coin),
      message,
      currentOpenWalletAtAddress,
      walletPassword,
    )
      .then(() => {
        console.log("onSubmitClick: Successful");
        setForceURL("/send-processing");
      })
      .catch((errorJsonString) => {
        console.log("onSubmitClick: errRes:", errorJsonString);
        const errorObject = JSON.parse(errorJsonString);
        console.log("onSubmitClick: errorObject:", errorObject);

        let err = {};
        if (errorObject.to != "") {
          err.payTo = errorObject.to;
        }
        if (errorObject.coin != "") {
          err.coin = errorObject.coin;
        }
        if (errorObject.value != "") {
          err.coin = errorObject.value;
        }
        if (errorObject.message != "") {
          err.message = errorObject.message;
        }
        if (errorObject.wallet_password != "") {
          err.walletPassword = errorObject.wallet_password;
        }
        console.log("onSubmitClick: err:", err);
        window.scrollTo(0, 0); // Start the page at the top of the page.
        setErrors(err);
      })
      .finally(() => {
        // Update the GUI to let user know that the operation is completed.
        setIsLoading(false);
      });
  };

  const validateForm = () => {
    const newErrors = {};

    if (!formData.recipient.trim()) {
      newErrors.recipient = "Recipient address is required";
    }

    if (!formData.amount) {
      newErrors.amount = "Amount is required";
    } else if (isNaN(formData.amount) || parseFloat(formData.amount) <= 0) {
      newErrors.amount = "Please enter a valid amount";
    } else if (parseFloat(formData.amount) > currentBalance) {
      newErrors.amount = "Insufficient balance";
    }

    if (!formData.password) {
      newErrors.password = "Password is required to authorize transaction";
    }

    setErrors(newErrors);
    setShowErrorBox(Object.keys(newErrors).length > 0);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = () => {
    if (validateForm()) {
      console.log("Form submitted");
    } else {
      window.scrollTo({ top: 0, behavior: "smooth" });
    }
  };

  const handleInputChange = (e) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
    if (errors[name]) {
      setErrors((prev) => ({ ...prev, [name]: "" }));
      setShowErrorBox(false);
    }
  };

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

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  const remainingBalance = formData.amount
    ? currentBalance - parseFloat(formData.amount)
    : currentBalance;

  return (
    <div>
      <main className="max-w-2xl mx-auto px-6 py-12 mb-24">
        {showErrorBox && (
          <div className="mb-6 bg-red-50 border-l-4 border-red-500 rounded-r-lg p-4 flex items-start gap-3">
            <AlertCircle className="w-5 h-5 text-red-500 flex-shrink-0 mt-0.5" />
            <div className="flex-grow">
              <h3 className="font-semibold text-red-800">
                Unable to Send Coins
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
                <div className="p-2 bg-purple-100 rounded-xl">
                  <Send
                    className="w-5 h-5 text-purple-600"
                    aria-hidden="true"
                  />
                </div>
                <h2 className="text-xl font-bold text-gray-900">
                  Send ComicCoins
                </h2>
              </div>
            </div>
            <p className="text-sm text-gray-500">
              Transfer coins to another wallet address. Please fill in all
              required fields.
            </p>
          </div>

          <div className="p-6 space-y-8">
            <div className="bg-amber-50 border border-amber-100 rounded-xl p-4 flex gap-3">
              <AlertCircle className="w-5 h-5 text-amber-600 flex-shrink-0 mt-0.5" />
              <div className="text-sm text-amber-800">
                <p className="font-semibold mb-1">Important Notice</p>
                <p>
                  All transactions are final and cannot be undone. Please verify
                  all details before sending.
                </p>
              </div>
            </div>

            <div className="space-y-6">
              <label className="block">
                <span className="text-sm font-medium text-gray-700">
                  Pay To <span className="text-red-500">*</span>
                </span>
                <input
                  type="text"
                  name="recipient"
                  value={formData.recipient}
                  onChange={handleInputChange}
                  className={`mt-1 block w-full px-4 py-3 bg-white border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-colors ${
                    errors.recipient
                      ? "border-red-300 bg-red-50"
                      : "border-gray-200"
                  }`}
                  placeholder="Enter recipient's wallet address"
                />
                {errors.recipient && (
                  <p className="mt-2 text-sm text-red-600 flex items-center gap-2">
                    <AlertCircle className="w-4 h-4" />
                    {errors.recipient}
                  </p>
                )}
              </label>

              <label className="block">
                <div className="flex items-center justify-between mb-1">
                  <span className="text-sm font-medium text-gray-700">
                    Coins <span className="text-red-500">*</span>
                  </span>
                  <span className="text-sm text-gray-500 flex items-center gap-1">
                    <Coins className="w-4 h-4" />
                    Balance: {currentBalance} CC
                  </span>
                </div>
                <input
                  type="number"
                  name="amount"
                  value={formData.amount}
                  onChange={handleInputChange}
                  max={currentBalance}
                  className={`mt-1 block w-full px-4 py-3 bg-white border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-colors ${
                    errors.amount
                      ? "border-red-300 bg-red-50"
                      : "border-gray-200"
                  }`}
                  placeholder="Enter amount of coins to send"
                />
                {errors.amount ? (
                  <p className="mt-2 text-sm text-red-600 flex items-center gap-2">
                    <AlertCircle className="w-4 h-4" />
                    {errors.amount}
                  </p>
                ) : (
                  formData.amount && (
                    <p className="mt-1.5 text-sm text-gray-500">
                      Remaining balance after transaction:{" "}
                      {remainingBalance.toFixed(2)} CC
                    </p>
                  )
                )}
              </label>

              <label className="block">
                <div className="flex items-center justify-between mb-1">
                  <span className="text-sm font-medium text-gray-700">
                    Message (Optional)
                  </span>
                  <span className="text-xs text-gray-500">
                    Include a note with your transaction
                  </span>
                </div>
                <textarea
                  name="message"
                  value={formData.message}
                  onChange={handleInputChange}
                  className="mt-1 block w-full px-4 py-3 bg-white border border-gray-200 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-colors"
                  placeholder="Add a message to this transaction"
                  rows={3}
                />
                <p className="mt-1 text-xs text-gray-500">
                  This message will be visible to the recipient of your
                  transaction.
                </p>
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
                  className={`mt-1 block w-full px-4 py-3 bg-white border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent transition-colors ${
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
            </div>

            <div className="flex justify-end pt-4">
              <button
                onClick={handleSubmit}
                className="px-6 py-2.5 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors inline-flex items-center gap-2"
              >
                Send Coins
                <Send className="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}

export default SendCoinView;
