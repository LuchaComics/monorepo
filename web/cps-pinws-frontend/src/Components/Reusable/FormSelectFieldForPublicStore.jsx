import React, { useState, useEffect } from "react";
import { getPublicTenantSelectOptionListAPI } from "../../API/tenant";

/**
EXAMPLE USAGE:

    <FormTenantField
      tenantID={tenantID}
      setTenantID={setTenantID}
      tenantName={tenantName}
      setTenantName={setTenantName}
      errorText={errors && errors.tenantID}
      helpText="Please select the tenant"
      maxWidth="310px"
    />
*/
function FormSelectFieldForPublicTenant({
  label = "Tenant",
  tenantID,
  setTenantID,
  tenantName = null,
  setTenantName = null,
  errorText,
  validationText,
  helpText,
  disabled,
}) {
  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setIsFetching] = useState(false);
  const [workoutProgramTypeSelectOptions, setTenantSelectOptions] = useState([]);

  ////
  //// Event handling.
  ////

  const setTenantIDAndTenantName = (oid, on) => {
    setTenantID(oid);
    setTenantName(on);
  };

  ////
  //// API.
  ////

  function onTenantSelectOptionsSuccess(response) {
    console.log("onTenantSelectOptionsSuccess: Starting...");
    let b = [{ value: "", label: "Please select" }, ...response];
    setTenantSelectOptions(b);
  }

  function onTenantSelectOptionsError(apiErr) {
    console.log("onTenantSelectOptionsError: Starting...");
    setErrors(apiErr);
  }

  function onTenantSelectOptionsDone() {
    console.log("onTenantSelectOptionsDone: Starting...");
    setIsFetching(false);
  }

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      setIsFetching(true);
      setErrors({});
      getPublicTenantSelectOptionListAPI(
        new Map(),
        onTenantSelectOptionsSuccess,
        onTenantSelectOptionsError,
        onTenantSelectOptionsDone,
      );
    }

    return () => {
      mounted = false;
    };
  }, []);

  ////
  //// Component rendering.
  ////

  // Render the JSX component.
  return (
    <div class="field pb-4">
      <label class="label">{label}</label>
      <div class="control">
        <span class="select">
          <select
            class={`input ${errorText && "is-danger"} ${validationText && "is-success"} has-text-black`}
            name={`tenantID`}
            placeholder={`Pick the workout program type`}
            onChange={(e) =>
              setTenantIDAndTenantName(
                e.target.value,
                e.target.options[e.target.selectedIndex].text,
              )
            }
            disabled={disabled}
          >
            {workoutProgramTypeSelectOptions &&
              workoutProgramTypeSelectOptions.length > 0 &&
              workoutProgramTypeSelectOptions.map(function (option, i) {
                // console.log("tenantID", tenantID);
                // console.log("option.value", option.value);
                // console.log(tenantID, "===", option.value, "->>>", tenantID === option.value);
                // console.log("");
                return (
                  <option
                    selected={tenantID === option.value}
                    value={option.value}
                  >
                    {option.label}
                  </option>
                );
              })}
          </select>
        </span>
      </div>
      {helpText && <p class="help">{helpText}</p>}
      {errorText && <p class="help is-danger">{errorText}</p>}
    </div>
  );
}

export default FormSelectFieldForPublicTenant;
