import React, { useState, useEffect } from "react";
import { getOfferSelectOptionListAPI } from "../../API/Offer";

/**
EXAMPLE USAGE:

    <FormSelectFieldForOffer
      label="Offer"
      offerID={offerID}
      setOfferID={setOfferID}
      errorText={errors && errors.offerID}
      helpText=""
      maxWidth="310px"
    />
*/
function FormSelectFieldForOffer({
  label = "Offer",
  offerID,
  setOfferID,
  errorText,
  validationText,
  helpText,
  disabled,
  extraOptions = [],
}) {
  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [offerSelectOptions, setOfferSelectOptions] = useState([]);

  ////
  //// Event handling.
  ////

  // Do nothing...

  ////
  //// API.
  ////

  function onOfferSelectOptionsSuccess(response) {
    console.log("onOfferSelectOptionsSuccess: Starting...");
    let b = [
      { value: "", label: "Please select" },
      ...extraOptions,
      ...response,
    ];
    setOfferSelectOptions(b);
  }

  function onOfferSelectOptionsError(apiErr) {
    console.log("onOfferSelectOptionsError: Starting...");
    setErrors(apiErr);
  }

  function onOfferSelectOptionsDone() {
    console.log("onOfferSelectOptionsDone: Starting...");
    setFetching(false);
  }

  ////
  //// Misc.
  ////

  useEffect(() => {
    let mounted = true;

    if (mounted) {
      setFetching(true);
      getOfferSelectOptionListAPI(
        onOfferSelectOptionsSuccess,
        onOfferSelectOptionsError,
        onOfferSelectOptionsDone,
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
            name={`offerID`}
            placeholder={`Pick the workout program type`}
            onChange={(e) => setOfferID(e.target.value)}
            disabled={disabled}
          >
            {offerSelectOptions &&
              offerSelectOptions.length > 0 &&
              offerSelectOptions.map(function (option, i) {
                // console.log("offerID", offerID);
                // console.log("option.value", option.value);
                // console.log(offerID, "===", option.value, "->>>", offerID === option.value);
                // console.log("");
                return (
                  <option
                    selected={offerID === option.value}
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

export default FormSelectFieldForOffer;
