import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";
import { useRecoilState } from "recoil";
import { WalletMinimal, Send, QrCode, MoreHorizontal, Wallet, ArrowRight, Plus, Download, ExternalLink, Lock, HardDrive, Copy, Upload } from 'lucide-react';

import { currentOpenWalletAtAddressState } from "../../../AppState";
import {
    ListWallets,
    SetDefaultWalletAddress,
    ExportWalletUsingDialog,
    ImportWalletUsingDialog
} from "../../../../wailsjs/go/main/App";

const ListWalletsView = () => {
    const [currentOpenWalletAtAddress, setCurrentOpenWalletAtAddress] = useRecoilState(currentOpenWalletAtAddressState);
    const [isLoading, setIsLoading] = useState(false);
    const [wallets, setWallets] = useState([]);
    const [errors, setErrors] = useState({});
    const [forceURL, setForceURL] = useState("");

    const loadWallets = () => {
        setIsLoading(true);
        ListWallets()
            .then((walletsResponse) => {
                console.log("loadWallets: walletsResponse:", walletsResponse);
                setWallets(walletsResponse);
            })
            .catch((error) => {
                console.error("Error fetching wallets:", error);
                setWallets([]);
            })
            .finally(() => {
                setIsLoading(false);
            });
    };

    useEffect(() => {
        window.scrollTo(0, 0);
        loadWallets();
    }, []);

    const onWalletOpen = (walletAddress) => {
        console.log("currentOpenWalletAtAddress: Old:", currentOpenWalletAtAddress);
        SetDefaultWalletAddress(walletAddress).then(()=>{
            setCurrentOpenWalletAtAddress(walletAddress);
            console.log("currentOpenWalletAtAddress: New:", walletAddress);
            setForceURL("/dashboard");
        }).catch((error) => {
            console.error("Error setting default wallet:", error);
        });
    };

    const onExportWallet = (walletAddress, e) => {
        e.preventDefault(); // Prevent navigation
        console.log("Exporting wallet:", walletAddress);
        ExportWalletUsingDialog(walletAddress).catch((error) => {
            console.error("Error exporting wallet:", error);
        });
    };

    const onImportWallet = (e) => {
        e.preventDefault();
        console.log("onImportWallet: Beginning...");
        ImportWalletUsingDialog()
            .then(() => {
                console.log("ImportWalletUsingDialog: Successfully imported wallet");
                loadWallets(); // Reload the wallet list after import
            })
            .catch((errorJsonString) => {
                console.error("ImportWalletUsingDialog: Error importing wallet:", errorJsonString);
                try {
                    const errorObject = JSON.parse(errorJsonString);
                    let err = {};
                    if (errorObject.filepath !== "") {
                        err.filepath = errorObject.filepath;
                    }
                    setErrors(err);
                } catch (e) {
                    console.error("Error parsing error response:", e);
                }
            });
    };

    if (forceURL !== "") {
        return <Navigate to={forceURL} />;
    }

    return (
        <div>
            <main className="max-w-4xl mx-auto px-4 sm:px-6 py-12 mb-24">
                {/* Add/Import Wallet Buttons */}
                <div className="mb-8 flex flex-col sm:flex-row gap-4">
                    <h2 className="text-2xl font-bold text-gray-900">Your Wallets</h2>
                    <div className="flex gap-4 sm:ml-auto">
                        <button
                            onClick={onImportWallet}
                            className="inline-flex items-center gap-2 border border-purple-600 text-purple-600 px-4 py-2 rounded-lg hover:bg-purple-50 transition-colors"
                            aria-label="Import wallet"
                        >
                            <Upload className="w-5 h-5" />
                            Import Wallet
                        </button>
                        <Link
                            to="/more/wallets/add"
                            className="inline-flex items-center gap-2 bg-purple-600 text-white px-4 py-2 rounded-lg hover:bg-purple-700 transition-colors"
                            aria-label="Create new wallet"
                        >
                            <Plus className="w-5 h-5" />
                            Add New Wallet
                        </Link>
                    </div>
                </div>

                {/* Security Information Cards - unchanged */}
                <div className="grid gap-6 mb-8">
                    <div className="bg-blue-50 border border-blue-200 rounded-xl p-6">
                        <div className="flex items-start gap-4">
                            <div className="p-2 bg-blue-100 rounded-lg shrink-0">
                                <Lock className="w-6 h-6 text-blue-600" />
                            </div>
                            <div>
                                <h3 className="font-semibold text-blue-900 mb-2">Encrypted Storage</h3>
                                <p className="text-blue-800">
                                    Your wallet is stored encrypted at rest to protect your coins and tokens.
                                    To access your wallet and perform transactions, you'll need to enter the
                                    password you created when you set up your wallet.
                                </p>
                            </div>
                        </div>
                    </div>

                    <div className="bg-green-50 border border-green-200 rounded-xl p-6">
                        <div className="flex items-start gap-4">
                            <div className="p-2 bg-green-100 rounded-lg shrink-0">
                                <HardDrive className="w-6 h-6 text-green-600" />
                            </div>
                            <div>
                                <h3 className="font-semibold text-green-900 mb-2">Local Storage</h3>
                                <p className="text-green-800">
                                    Your wallet is stored exclusively on your computer and not in the cloud.
                                    Make sure to keep regular backups and store your recovery phrase in a
                                    safe place.
                                </p>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Wallet List */}
                <div className="bg-white rounded-xl shadow-sm border border-gray-200">
                    <div className="p-6 border-b border-gray-100">
                        <div className="grid grid-cols-12 gap-4 text-sm font-medium text-gray-500">
                            <div className="col-span-4 sm:col-span-3">Wallet Label</div>
                            <div className="col-span-8 sm:col-span-6">Address</div>
                            <div className="hidden sm:block sm:col-span-3 text-right">Actions</div>
                        </div>
                    </div>

                    <div className="divide-y divide-gray-100">
                        {wallets.map((wallet) => (
                            <div key={wallet.address} className="p-6">
                                {/* Top row: Label and Address */}
                                <div className="grid grid-cols-12 gap-4 items-center mb-4 sm:mb-0">
                                    <div className="col-span-4 sm:col-span-3">
                                        <div className="font-medium text-gray-900 truncate" title={wallet.label}>
                                            {wallet.label}
                                        </div>
                                    </div>
                                    <div className="col-span-8 sm:col-span-6 flex items-center gap-2">
                                        <div className="font-mono text-gray-600 truncate" title={wallet.address}>
                                            {wallet.address}
                                        </div>
                                        <button
                                            onClick={() => navigator.clipboard.writeText(wallet.address)}
                                            className="p-1 hover:bg-gray-100 rounded-md shrink-0"
                                            aria-label="Copy address"
                                        >
                                            <Copy className="w-4 h-4 text-gray-500" />
                                        </button>
                                    </div>

                                    {/* Buttons - Desktop */}
                                    <div className="hidden sm:flex sm:col-span-3 items-center justify-end gap-2">
                                        <button
                                            onClick={(e) => onExportWallet(wallet.address, e)}
                                            className="inline-flex items-center gap-1 px-3 py-1.5 text-sm border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
                                            aria-label={`Export ${wallet.label}`}
                                        >
                                            <Download className="w-4 h-4" />
                                            Export
                                        </button>
                                        <button
                                            onClick={() => onWalletOpen(wallet.address)}
                                            className="inline-flex items-center gap-1 px-3 py-1.5 text-sm bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors"
                                            aria-label={`Open ${wallet.label}`}
                                        >
                                            <ExternalLink className="w-4 h-4" />
                                            Open
                                        </button>
                                    </div>
                                </div>

                                {/* Bottom row: Buttons (Mobile only) */}
                                <div className="sm:hidden flex justify-end gap-2">
                                    <button
                                        onClick={(e) => onExportWallet(wallet.address, e)}
                                        className="inline-flex items-center gap-1 px-3 py-1.5 text-sm border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
                                        aria-label={`Export ${wallet.label}`}
                                    >
                                        <Download className="w-4 h-4" />
                                        Export
                                    </button>
                                    <button
                                        onClick={() => onWalletOpen(wallet.address)}
                                        className="inline-flex items-center gap-1 px-3 py-1.5 text-sm bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors"
                                        aria-label={`Open ${wallet.label}`}
                                    >
                                        <ExternalLink className="w-4 h-4" />
                                        Open
                                    </button>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            </main>
        </div>
    );
};

export default ListWalletsView;
