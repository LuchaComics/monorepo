import React from "react";
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
    faArrowUpRightFromSquare,
} from "@fortawesome/free-solid-svg-icons";


function FormRowIPFSDownloadLink(props) {
  const { label, ipfsGatewayDomain, ipfsPath, helpText, type = "text" } = props;
  return (
    <div class="field pb-4">
      <label class="label has-text-black">{label}</label>
      <div class="control">
        <p>
          <Link to={`${ipfsGatewayDomain}/ipfs/${ipfsPath.replace("ipfs://", "")}`}>{ipfsPath}&nbsp;<FontAwesomeIcon className="fas" icon={faArrowUpRightFromSquare} /></Link>
        </p>
        {helpText !== undefined && helpText !== null && helpText !== "" && (
          <p class="help">{helpText}</p>
        )}
      </div>
    </div>
  );
}

export default FormRowIPFSDownloadLink;
