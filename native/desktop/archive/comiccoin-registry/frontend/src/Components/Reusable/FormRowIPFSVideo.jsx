import React, { useState, useEffect }  from "react";
import { Link } from "react-router-dom";

import { GetFileViaIPFS } from "../../../wailsjs/go/main/App";

function FormRowIPFSVideo(props) {

  ////
  //// Props.
  ////
  const { label, ipfsGatewayDomain, ipfsPath, helpText, type = "text" } = props;

  ////
  //// Component states.
  ////

  const [fileURL, setFileURL] = useState("");
  const [contentType, setContentType] = useState("");

  ////
  //// Misc + API
  ////

  useEffect(() => {
   let mounted = true;

   if (mounted) {
     GetFileViaIPFS(ipfsPath).then((response) => {
       const bytes = new Uint8Array(response.data);
       const contentType = response.content_type;

       const fileUrl = URL.createObjectURL(new Blob([bytes], { type: contentType }));
       setFileURL(fileUrl);
       setContentType(contentType);
     }).catch((err) => {
       console.log("err:", err);
     });
   }
   return () => {
     mounted = false;
   };
 }, [ipfsPath]);

  return (
    <div class="field pb-4">
      <label class="label has-text-black">{label}</label>
      <div class="control">
        {fileURL && (
          <p>
            {ipfsPath && <>
                <video width="320" height="240" controls>
                    <source src={`${ipfsGatewayDomain}/ipfs/${ipfsPath.replace("ipfs://", "")}`} type={contentType} />
                    Your browser does not support the video tag.
                </video>
            </>}
          </p>
        )}
        {helpText !== undefined && helpText !== null && helpText !== "" && (
          <p class="help">{helpText}</p>
        )}
      </div>
    </div>
  );
}

export default FormRowIPFSVideo;
