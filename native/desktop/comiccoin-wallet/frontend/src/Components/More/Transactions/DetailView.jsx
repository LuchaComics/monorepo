import {useState, useEffect} from 'react';
import { Navigate, useParams } from "react-router-dom";
import { WalletMinimal, Send, QrCode, MoreHorizontal, Wallet, ChevronDown, ChevronUp, Link, FileText, Shield, Copy } from 'lucide-react';
import { useRecoilState } from "recoil";
import { toLower } from "lodash";

import { GetBlockDataByBlockTransactionTimestamp } from "../../../../wailsjs/go/main/App";
import { currentOpenWalletAtAddressState } from "../../../AppState";

// Helper components moved outside main component
const InfoRow = ({ label, value, copyable = false }) => (
  <div className="py-3 flex flex-col sm:flex-row sm:items-center border-b border-gray-100 last:border-0">
    <span className="text-sm text-gray-600 sm:w-1/3">{label}</span>
    <div className="mt-1 sm:mt-0 sm:w-2/3 flex items-center gap-2">
      <span className="text-sm font-medium text-gray-900 break-all">{value}</span>
      {copyable && (
        <button
          onClick={() => navigator.clipboard.writeText(value)}
          className="p-1 hover:bg-gray-100 rounded-md"
        >
          <Copy className="w-4 h-4 text-gray-500" />
        </button>
      )}
    </div>
  </div>
);

const Section = ({ icon: Icon, title, children, expandable = false, expanded = false, onToggle = null }) => (
  <div className="mb-6 bg-white rounded-xl shadow-sm border border-gray-200">
    <div className="p-6 border-b border-gray-100">
      <div className="flex justify-between items-center">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-purple-100 rounded-xl">
            <Icon className="w-5 h-5 text-purple-600" />
          </div>
          <h2 className="text-xl font-bold text-gray-900">{title}</h2>
        </div>
        {expandable && (
          <button
            onClick={onToggle}
            className="flex items-center gap-1 text-sm text-purple-600 hover:text-purple-800"
          >
            {expanded ? (
              <>
                <ChevronUp className="w-4 h-4" />
                See Less
              </>
            ) : (
              <>
                <ChevronDown className="w-4 h-4" />
                See More
              </>
            )}
          </button>
        )}
      </div>
    </div>
    <div className="p-6">{children}</div>
  </div>
);

function TransactionDetailView() {
    const { timestamp } = useParams();

    // All state declarations consolidated at the top
    const [currentOpenWalletAtAddress] = useRecoilState(currentOpenWalletAtAddressState);
    const [isLoading, setIsLoading] = useState(false);
    const [blockData, setBlockData] = useState(null);
    const [showMoreBlockInfo, setShowMoreBlockInfo] = useState(false);
    const [showMoreTxInfo, setShowMoreTxInfo] = useState(false);
    const [forceURL, setForceURL] = useState("");

    useEffect(() => {
        let mounted = true;

        if (mounted) {
            window.scrollTo(0, 0);
            setIsLoading(true);

            GetBlockDataByBlockTransactionTimestamp(parseInt(timestamp))
                .then((res) => {
                    console.log("GetBlockDataByBlockTransactionTimestamp: res:", res);
                    if (mounted) {
                        setBlockData(res);
                    }
                })
                .catch((errorRes) => {
                    console.log("GetBlockDataByBlockTransactionTimestamp: errors:", errorRes);
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
    }, [timestamp]);

    if (forceURL !== "") {
        return null; // Or your navigation component
    }

    if (isLoading) {
        return "Loading...";
    }

    return (
        <div>
            {blockData && (
                <main className="max-w-3xl mx-auto px-4 sm:px-6 py-12 mb-24">
                    <Section
                        icon={Link}
                        title="Block Information"
                        expandable
                        expanded={showMoreBlockInfo}
                        onToggle={() => setShowMoreBlockInfo(!showMoreBlockInfo)}
                    >
                        <InfoRow label="Hash" value={blockData.hash} copyable />
                        <InfoRow label="Number" value={blockData.header.number} />
                        <InfoRow
                            label="Timestamp"
                            value={new Date(blockData.header.timestamp).toLocaleString()}
                        />

                        {showMoreBlockInfo && (
                            <>
                                <InfoRow label="Previous Block Hash" value={blockData.header.prev_block_hash} copyable />
                                <InfoRow label="Trans Root" value={blockData.header.trans_root} copyable />
                                <InfoRow label="Tokens Root" value={blockData.header.tokens_root} copyable />
                                <InfoRow label="Beneficiary" value={blockData.header.beneficiary} copyable />
                                <InfoRow label="Nonce" value={blockData.header.nonce} />
                                <InfoRow label="Mining Reward" value={blockData.header.mining_reward} />
                                <InfoRow label="Difficulty" value={blockData.header.difficulty} />
                                <InfoRow label="Header Signature" value={blockData.header_signature_bytes} copyable />
                            </>
                        )}
                    </Section>

                    <Section icon={Shield} title="Blockchain Validation Service">
                        <InfoRow label="Validator ID" value={blockData.validator.id} />
                        <InfoRow label="Public Key" value={blockData.validator.public_key_bytes} copyable />
                    </Section>

                    <Section
                        icon={FileText}
                        title="Transaction Information"
                        expandable
                        expanded={showMoreTxInfo}
                        onToggle={() => setShowMoreTxInfo(!showMoreTxInfo)}
                    >
                        {blockData.trans.map((tx, index) => (
                            <div key={index} className="space-y-3">
                                <InfoRow label="Purpose" value={tx.type === "coin" ? "Coin" : "Token"} />
                                <InfoRow label="Type" value={tx.type} />
                                <InfoRow label="Timestamp" value={new Date(tx.timestamp).toLocaleString()} />
                                <InfoRow label="Value" value={tx.value} />
                                <InfoRow label="Fee" value={tx.fee} />
                                <InfoRow label="Actual Value" value={(parseFloat(tx.value) - parseFloat(tx.fee)).toString()} />
                                <InfoRow label="From" value={tx.from} copyable />
                                <InfoRow label="To" value={tx.to} copyable />

                                {showMoreTxInfo && (
                                    <>
                                        <InfoRow label="Data" value={tx.data} />
                                        <InfoRow label="Nonce" value={tx.nonce} />
                                        <InfoRow label="Chain ID" value={tx.chain_id} />
                                        <InfoRow label="Gas Price" value={tx.gas_price} />
                                        <InfoRow label="Gas Unit" value={tx.gas_unit} />
                                        <InfoRow label="R" value={tx.r} />
                                        <InfoRow label="S" value={tx.s} />
                                        <InfoRow label="V" value={tx.v} />
                                        <InfoRow label="Tip" value={tx.tip} />
                                    </>
                                )}
                            </div>
                        ))}
                    </Section>
                </main>
            )}
        </div>
    );
}

export default TransactionDetailView;
