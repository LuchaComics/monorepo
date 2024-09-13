import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";

import { getTenantSelectOptionListAPI } from "../../API/tenant";
import { getSelectedOptions } from "../../Helpers/selectHelper";


function DataDisplayRowTenant(props) {
  ////
  //// Props.
  ////

  const {
    label = "Please select",
    tenantID,
    helpText,
  } = props;

  ////
  //// Component states.
  ////

  // GUI related states.
  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);

  // Form states.
  const [tenantName, setTenantName] = useState("");

  ////
  //// API.
  ////

  // --- Get options --- //

  function onTenantOptionListSuccess(response) {
    console.log("onTenantOptionListSuccess: Starting...");
    if (response !== null) {


      response.forEach(function (item, index) {
          console.log(item);
          if (tenantID === item.value) {
              setTenantName(item.label);
          }
      });


    }
  }

  function onTenantOptionListError(apiErr) {
    console.log("onTenantOptionListError: Starting...");
    console.log("onTenantOptionListError: apiErr:", apiErr);
    setErrors(apiErr);
  }

  function onTenantOptionListDone() {
    console.log("onTenantOptionListDone: Starting...");
    setFetching(false);
  }

  // --- All --- //

  const onUnauthorized = () => {
    // Do nothing...
  };

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
        setFetching(true);

        let params = new Map();
        getTenantSelectOptionListAPI(
          params,
          onTenantOptionListSuccess,
          onTenantOptionListError,
          onTenantOptionListDone,
          onUnauthorized,
        );
    }

    return () => {
      mounted = false;
    };
  }, []);

  ////
  //// Component rendering.
  ////

  return (
    <div class="field pb-4">
      <label class="label">{label}</label>
      <div class="control">
        <p>
          {tenantName !== undefined && tenantName !== null && tenantName !== "" && (
            <>{tenantName}</>
          )}
        </p>
        {helpText !== undefined && helpText !== null && helpText !== "" && (
          <p class="help">{helpText}</p>
        )}
      </div>
    </div>
  );
}

export default DataDisplayRowTenant;
