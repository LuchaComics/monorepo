import { useState, useEffect } from "react";
import { Navigate } from "react-router-dom";
import {
  Settings,
  AlertCircle,
  Info,
  Save,
  RefreshCw,
  ArrowLeft,
} from "lucide-react";

import {
  GetPreferences,
  SavePreferences,
} from "../../../../wailsjs/go/main/App";

function SettingsView() {
  const [isLoading, setIsLoading] = useState(false);
  const [initialFormData, setInitialFormData] = useState({});
  const [formData, setFormData] = useState({});
  const [errors, setErrors] = useState({});
  const [showErrorBox, setShowErrorBox] = useState(false);
  const [forceURL, setForceURL] = useState("");

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      window.scrollTo(0, 0);
      setIsLoading(true);

      GetPreferences()
        .then((res) => {
          console.log("GetPreferences: res:", res);
          if (mounted) {
            setInitialFormData(res);
            setFormData(res);
          }
        })
        .catch((errorRes) => {
          console.log("GetPreferences: errors:", errorRes);
        })
        .finally(() => {
          if (mounted) {
            setIsLoading(false);
          }
        });
    }

    return () => {
      mounted = false;
    };
  }, []);

  // Check if editable fields have been modified
  const hasChanges = () => {
    return (
      formData.nft_storage_address !== initialFormData.nft_storage_address ||
      formData.authority_address !== initialFormData.authority_address
    );
  };

  const validateForm = () => {
    const newErrors = {};

    if (!formData.nft_storage_address) {
      newErrors.nft_storage_address = "NFT Storage Address is required";
    } else if (!formData.nft_storage_address.startsWith("http")) {
      newErrors.nft_storage_address =
        "NFT Storage Address must be a valid HTTP URL";
    }

    if (!formData.authority_address) {
      newErrors.authority_address = "Authority Address is required";
    } else if (!formData.authority_address.startsWith("http")) {
      newErrors.authority_address =
        "Authority Address must be a valid HTTP URL";
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

  const handleSubmit = () => {
    if (validateForm()) {
      console.log("Settings saved");
      setIsLoading(true);
      SavePreferences(formData)
        .then((res) => {
          console.log("SavePreferences: res:", res);
          // setForceURL("/more");
        })
        .catch((errorRes) => {
          console.log("SavePreferences: errors:", errorRes);
        })
        .finally(() => {
          setIsLoading(false);
        });
    } else {
      window.scrollTo({ top: 0, behavior: "smooth" });
    }
  };

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <div>
      <main className="max-w-2xl mx-auto px-6 py-12 mb-24">
        {showErrorBox && (
          <div className="mb-6 bg-red-50 border-l-4 border-red-500 rounded-r-lg p-4 flex items-start gap-3">
            <AlertCircle className="w-5 h-5 text-red-500 flex-shrink-0 mt-0.5" />
            <div className="flex-grow">
              <h3 className="font-semibold text-red-800">
                Unable to Save Settings
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
                <div className="p-2 bg-blue-100 rounded-xl">
                  <Settings
                    className="w-5 h-5 text-blue-600"
                    aria-hidden="true"
                  />
                </div>
                <h2 className="text-xl font-bold text-gray-900">Settings</h2>
              </div>
            </div>
            <p className="text-sm text-gray-500">
              Configure your wallet settings and connection parameters.
            </p>
          </div>

          <div className="p-6 space-y-8">
            {/* Show restart notice only if there are changes */}
            {hasChanges() && (
              <div className="bg-amber-50 border border-amber-100 rounded-xl p-4 flex gap-3">
                <RefreshCw className="w-5 h-5 text-amber-600 flex-shrink-0 mt-0.5" />
                <div className="text-sm text-amber-800">
                  <p className="font-semibold mb-1">
                    Application Restart Required
                  </p>
                  <p>
                    The application will restart automatically after saving
                    these settings to apply the changes.
                  </p>
                </div>
              </div>
            )}

            <div className="space-y-6">
              <label className="block">
                <span className="text-sm font-medium text-gray-700">
                  Chain ID
                </span>
                <input
                  type="text"
                  value={formData.chain_id}
                  className="mt-1 block w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-lg text-gray-500"
                  readOnly
                />
                <p className="mt-2 text-sm text-gray-600 flex items-start gap-2">
                  <Info className="w-4 h-4 flex-shrink-0 mt-0.5" />
                  <span>
                    This is the blockchain network your wallet is subscribed to.
                  </span>
                </p>
              </label>

              <label className="block">
                <span className="text-sm font-medium text-gray-700">
                  Authority Address <span className="text-red-500">*</span>
                </span>
                <input
                  type="text"
                  name="authority_address"
                  value={formData.authority_address}
                  onChange={handleInputChange}
                  className={`mt-1 block w-full px-4 py-3 bg-white border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors ${
                    errors.authority_address
                      ? "border-red-300 bg-red-50"
                      : "border-gray-200"
                  }`}
                  placeholder="Enter blockchain authority address"
                />
                {errors.authority_address && (
                  <p className="mt-2 text-sm text-red-600 flex items-center gap-2">
                    <AlertCircle className="w-4 h-4" />
                    {errors.authority_address}
                  </p>
                )}
                <p className="mt-2 text-sm text-gray-600 flex items-start gap-2">
                  <Info className="w-4 h-4 flex-shrink-0 mt-0.5" />
                  <span>
                    The authority address is used to connect your wallet to the
                    Blockchain Network Service Provider.
                  </span>
                </p>
              </label>

              <label className="block">
                <span className="text-sm font-medium text-gray-700">
                  Data Directory
                </span>
                <input
                  type="text"
                  value={formData.data_directory}
                  className="mt-1 block w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-lg text-gray-500"
                  readOnly
                />
              </label>

              <label className="block">
                <span className="text-sm font-medium text-gray-700">
                  NFT Storage Address <span className="text-red-500">*</span>
                </span>
                <input
                  type="text"
                  name="nft_storage_address"
                  value={formData.nft_storage_address}
                  onChange={handleInputChange}
                  className={`mt-1 block w-full px-4 py-3 bg-white border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent transition-colors ${
                    errors.nft_storage_address
                      ? "border-red-300 bg-red-50"
                      : "border-gray-200"
                  }`}
                  placeholder="Enter IPFS gateway address"
                />
                {errors.nft_storage_address && (
                  <p className="mt-2 text-sm text-red-600 flex items-center gap-2">
                    <AlertCircle className="w-4 h-4" />
                    {errors.nft_storage_address}
                  </p>
                )}
                <p className="mt-2 text-sm text-gray-600 flex items-start gap-2">
                  <Info className="w-4 h-4 flex-shrink-0 mt-0.5" />
                  <span>
                    The IPFS gateway address is used to connect to the IPFS
                    network for accessing NFT metadata and content.
                  </span>
                </p>
              </label>

              <label className="block">
                <span className="text-sm font-medium text-gray-700">
                  Default Wallet Address
                </span>
                <input
                  type="text"
                  value={formData.default_wallet_address}
                  className="mt-1 block w-full px-4 py-3 bg-gray-50 border border-gray-200 rounded-lg text-gray-500"
                  readOnly
                />
              </label>
            </div>

            <div className="flex justify-between pt-4">
              <button
                onClick={() => setForceURL("/more")}
                className="px-6 py-2.5 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors inline-flex items-center gap-2"
              >
                <ArrowLeft className="w-4 h-4" />
                Back
              </button>
              <button
                onClick={handleSubmit}
                className="px-6 py-2.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors inline-flex items-center gap-2"
              >
                Save Settings
                <Save className="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
}

export default SettingsView;
