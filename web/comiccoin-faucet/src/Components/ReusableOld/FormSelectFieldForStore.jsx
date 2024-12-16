import React, { useState, useEffect } from "react";
import { getStoreSelectOptionListAPI } from "../../API/store";

/**
EXAMPLE USAGE:

    <FormStoreField
      storeID={storeID}
      setStoreID={setStoreID}
      storeName={storeName}
      setStoreName={setStoreName}
      errorText={errors && errors.storeID}
      helpText="Please select the store"
      maxWidth="310px"
    />
*/
function FormSelectFieldForStore({
  label = "Store",
  storeID,
  setStoreID,
  storeName = null,
  setStoreName = null,
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
  const [workoutProgramTypeSelectOptions, setStoreSelectOptions] = useState([]);

  ////
  //// Event handling.
  ////

  const setStoreIDAndStoreName = (oid, on) => {
    setStoreID(oid);
    setStoreName(on);
  };

  ////
  //// API.
  ////

  function onStoreSelectOptionsSuccess(response) {
    console.log("onStoreSelectOptionsSuccess: Starting...");
    let b = [{ value: "", label: "Please select" }, ...response];
    setStoreSelectOptions(b);
  }

  function onStoreSelectOptionsError(apiErr) {
    console.log("onStoreSelectOptionsError: Starting...");
    setErrors(apiErr);
  }

  function onStoreSelectOptionsDone() {
    console.log("onStoreSelectOptionsDone: Starting...");
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
      getStoreSelectOptionListAPI(
        new Map(),
        onStoreSelectOptionsSuccess,
        onStoreSelectOptionsError,
        onStoreSelectOptionsDone,
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
            name={`storeID`}
            placeholder={`Pick the workout program type`}
            onChange={(e) =>
              setStoreIDAndStoreName(
                e.target.value,
                e.target.options[e.target.selectedIndex].text,
              )
            }
            disabled={disabled}
          >
            {workoutProgramTypeSelectOptions &&
              workoutProgramTypeSelectOptions.length > 0 &&
              workoutProgramTypeSelectOptions.map(function (option, i) {
                // console.log("storeID", storeID);
                // console.log("option.value", option.value);
                // console.log(storeID, "===", option.value, "->>>", storeID === option.value);
                // console.log("");
                return (
                  <option
                    selected={storeID === option.value}
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

export default FormSelectFieldForStore;
