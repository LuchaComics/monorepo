import React, { useState, useEffect }  from "react";
import { Link } from "react-router-dom";

import { GetFileViaIPFS } from "../../../wailsjs/go/main/App";

function FormRowYouTubeField(props) {

  ////
  //// Props.
  ////
  const { label, url, helpText, type = "text" } = props;

  ////
  //// Component states.
  ////

  ////
  //// Misc + API
  ////

  useEffect(() => {
   let mounted = true;

   if (mounted) {

   }
   return () => {
     mounted = false;
   };
 }, []);

  return (
    <div class="field pb-4">
      <label class="label has-text-black">{label}</label>
      <div class="control">
      <p>
        {url ? <>
            <iframe width="420" height="315"
                src={url}>
            </iframe>
        </> : <>
           No video
        </>}
      </p>
        {helpText !== undefined && helpText !== null && helpText !== "" && (
          <p class="help">{helpText}</p>
        )}
      </div>
    </div>
  );
}

export default FormRowYouTubeField;
