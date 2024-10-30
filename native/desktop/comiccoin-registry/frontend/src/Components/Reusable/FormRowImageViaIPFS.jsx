import React, { useState, useEffect }  from "react";
import { Link } from "react-router-dom";

import { GetFileViaIPFS } from "../../../wailsjs/go/main/App";

function FormRowImageViaIPFS(props) {

  ////
  //// Props.
  ////
  const { label, ipfsPath, helpText, type = "text" } = props;

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
        GetFileViaIPFS(ipfsPath).then((response)=>{
            const bytes = new Uint8Array(response.data);

            const contentType = response.content_type;
            const contentLength = response.content_length;

            const decoder = new TextDecoder('utf-8');
            const fileContents = decoder.decode(bytes);

            const fileUrl = URL.createObjectURL(new Blob([bytes], { type: contentType }));
            console.log(fileUrl);
            setFileURL(fileUrl);

            setContentType(contentType);

            console.log("contentType:", contentType);
            console.log("contentLength:", contentLength);
            console.log("ipfsPath:", ipfsPath);

        }).catch((err)=>{
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
        <p>
          <img src={`fileURL`} />
        </p>
        {helpText !== undefined && helpText !== null && helpText !== "" && (
          <p class="help">{helpText}</p>
        )}
      </div>
    </div>
  );
}

export default FormRowImageViaIPFS;
