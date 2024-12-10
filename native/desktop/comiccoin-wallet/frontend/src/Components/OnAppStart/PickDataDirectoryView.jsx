import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";
import { WalletMinimal, FolderOpen, Info, ChevronRight, Settings, XCircle, AlertCircle } from 'lucide-react';

import PageLoadingContent from "../Reusable/PageLoadingContent";
import FormErrorBox from "../Reusable/FormErrorBox";
import {
    GetDefaultDataDirectory,
    GetDataDirectoryFromDialog,
    SaveDataDirectory,
    ShutdownApp,
} from "../../../wailsjs/go/main/App";


function PickDataDirectoryView() {

    ////
    //// Component states.
    ////

    const [isLoading, setIsLoading] = useState(false);
    const [errors, setErrors] = useState({});
    const [useDefaultLocation, setUseDefaultLocation] = useState(1);
    const [forceURL, setForceURL] = useState("");
    const [dataDirectory, setDataDirectory] = useState("./ComicCoin");
    const [showCancelWarning, setShowCancelWarning] = useState(false);

    const defaultPath = "./ComicCoin";
    const [selectedPath, setSelectedPath] = useState('');
    const [locationType, setLocationType] = useState('automatic');
    const [showError, setShowError] = useState(false);

    ////
    //// Event handling.
    ////

    const onSubmitClick = (e) => {
        console.log("onSubmitClick: Start");

        if (locationType === 'custom' && !selectedPath) {
          setShowError(true);
          window.scrollTo({ top: 0, behavior: 'smooth' });
          return
        }

        setIsLoading(true);

        console.log("onSubmitClick: dataDirectory:", dataDirectory);

        // Submit the `dataDirectory` value to our backend.
        SaveDataDirectory(dataDirectory).then( (result) => {
            console.log("result:", result);
            setForceURL("/create-your-first-wallet")
        }).finally(()=>{
            setIsLoading(false);
        });
    }

    const handleLocationChange = (value) => {
        setLocationType(value);
        setShowError(false);
    };

    ////
    //// API.
    ////

    ////
    //// Misc.
    ////

    useEffect(() => {
      let mounted = true;

      if (mounted) {
            window.scrollTo(0, 0); // Start the page at the top of the page.
            GetDefaultDataDirectory().then( (defaultDataDirResponse)=>{
                setDataDirectory(defaultDataDirResponse);
            })
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

    if (isLoading) {
        return (
            <PageLoadingContent displayMessage="Saving..." />
        );
    }

    return (
    <div>
      <main className="max-w-2xl mx-auto px-6 py-12">
        {showError && (
          <div className="mb-6 bg-red-50 border-l-4 border-red-500 rounded-r-lg p-4 flex items-start gap-3">
            <XCircle className="w-5 h-5 text-red-500 flex-shrink-0 mt-0.5" />
            <div className="flex-grow">
              <h3 className="font-semibold text-red-800">Directory Selection Required</h3>
              <p className="text-sm text-red-600 mt-1">
                Please select a directory location for storing blockchain data before proceeding.
              </p>
            </div>
          </div>
        )}

        <div className="bg-white rounded-xl border-2 border-gray-100">
          <div className="p-6">
            <div className="flex items-center justify-between mb-2">
              <div className="flex items-center gap-3">
                <div className="p-2 bg-purple-100 rounded-xl">
                  <Settings className="w-5 h-5 text-purple-600" aria-hidden="true" />
                </div>
                <h2 className="text-xl font-bold text-gray-900">Choose Directory Location</h2>
              </div>
            </div>
            <p className="text-sm text-gray-500">
              As this is your first time launching the application, please select where you'd like to store your wallet data.
            </p>
          </div>

          <div className="p-6 space-y-8">
            {/* Key Information Box */}
            <div className="bg-gradient-to-r from-purple-50 to-indigo-50 rounded-xl">
              <div className="p-6 flex gap-4">
                <Info className="w-5 h-5 text-purple-600 flex-shrink-0 mt-1" />
                <div className="text-sm text-gray-700">
                  <p>
                    ComicCoin Wallet will download and maintain a local copy of the ComicCoin blockchain.
                    This will require approximately 1 MB of storage space. Your wallet data will also be
                    stored in this location.
                  </p>
                </div>
              </div>
            </div>

            {/* Selection Options */}
            <div className="space-y-4">
              {/* Automatic Option */}
              <label className="flex items-start gap-4 p-5 rounded-xl cursor-pointer border-2 border-transparent hover:border-purple-200 hover:bg-purple-50 transition-all">
                <input
                  type="radio"
                  name="location"
                  value="automatic"
                  checked={locationType === 'automatic'}
                  onChange={(e) => handleLocationChange(e.target.value)}
                  className="mt-1"
                />
                <div>
                  <div className="flex items-center gap-3">
                    <p className="font-semibold text-gray-900">Automatic Location</p>
                    <span className="text-xs bg-purple-100 text-purple-700 px-2 py-0.5 rounded-full">Recommended</span>
                  </div>
                  <p className="text-sm text-gray-600 mt-2">
                    Let ComicCoin choose the best location to store blockchain data
                  </p>
                </div>
              </label>

              {/* Custom Option */}
              <label className={`flex items-start gap-4 p-5 rounded-xl cursor-pointer border-2 transition-all ${
                showError ? 'border-red-200 bg-red-50' : 'border-transparent hover:border-purple-200 hover:bg-purple-50'
              }`}>
                <input
                  type="radio"
                  name="location"
                  value="custom"
                  checked={locationType === 'custom'}
                  onChange={(e) => handleLocationChange(e.target.value)}
                  className="mt-1"
                />
                <div className="w-full">
                  <p className="font-semibold text-gray-900">Custom Location</p>
                  <p className="text-sm text-gray-600 mt-2">
                    Choose your preferred location for storing blockchain data
                  </p>

                  {locationType === 'custom' && (
                    <div className="mt-4 space-y-3">
                      <button
                        className={`inline-flex items-center px-4 py-2 rounded-lg transition-colors ${
                          showError
                            ? 'bg-red-100 text-red-700 hover:bg-red-200'
                            : 'bg-purple-100 text-purple-700 hover:bg-purple-200'
                        }`}
                        onClick={(e) =>
                          GetDataDirectoryFromDialog().then((dataDirectoryResult) => {
                              if (dataDirectoryResult !== "") {
                                  setDataDirectory(dataDirectoryResult);
                              }
                          })
                         }
                      >
                        <FolderOpen className="w-5 h-5 mr-2" />
                        Browse...
                      </button>
                      <div className="relative">
                        <input
                          type="text"
                          value={selectedPath}
                          readOnly
                          placeholder="No directory selected"
                          className={`w-full p-3 border rounded-lg text-sm ${
                            showError
                              ? 'border-red-200 bg-red-50 text-red-600 placeholder-red-300'
                              : 'border-gray-200 bg-gray-50 text-gray-600 placeholder-gray-400'
                          }`}
                        />
                        {showError && (
                          <div className="flex items-center gap-2 mt-2 text-red-600 text-sm">
                            <AlertCircle className="w-4 h-4" />
                            <span>Please select a directory location</span>
                          </div>
                        )}
                      </div>
                    </div>
                  )}
                </div>
              </label>
            </div>

            <div className="space-y-6">
              <p className="text-sm text-gray-600">
                After clicking Save, ComicCoin Wallet will begin downloading and processing the complete
                ComicCoin blockchain (1 MB), starting with the earliest transactions from 2024 when
                ComicCoin was launched.
              </p>

              <div className="bg-amber-50 border border-amber-100 rounded-xl p-6 text-sm text-amber-800">
                <p>
                  Please note: The initial synchronization process is resource-intensive and may reveal
                  previously unnoticed hardware issues with your computer. Each time you launch ComicCoin
                  Wallet, it will resume synchronization from where it last stopped.
                </p>
              </div>
            </div>

            <div className="flex justify-end gap-4 pt-4">
              <button
                className="px-6 py-2.5 text-gray-700 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-colors"
                onClick={(e)=>{
                    ShutdownApp();
                }}
              >
                Cancel & Close Wallet
              </button>
              <button
                className="px-6 py-2.5 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors inline-flex items-center gap-2"
                onClick={(e)=>{
                    onSubmitClick();
                }}
              >
                Save
                <ChevronRight className="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
};

export default PickDataDirectoryView
