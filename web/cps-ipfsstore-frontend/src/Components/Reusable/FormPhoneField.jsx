import React from "react";
import Input, {
  getCountries,
  getCountryCallingCode,
} from "react-phone-number-input/input";
const { getCode } = require("country-list");

function FormPhoneField({
  label,
  name,
  placeholder,
  selectedCountry,
  selectePhoneNumber,
  errorText,
  validationText,
  helpText,
  onChange,
  disabled,
  maxWidth,
}) {
  // DEVELOPERS NOTE:
  // https://github.com/country-regions/react-country-region-selector
  return (
    <div class="field pb-4">
      <label class="label">{label}</label>
      <div class="control" style={{ maxWidth: maxWidth }}>
        <span class="select">
          <Input
            country={getCode(selectedCountry)}
            name={name}
            placeholder={placeholder}
            disabled={disabled}
            class={`input ${errorText && "is-danger"} ${validationText && "is-success"} has-text-black`}
            value={selectePhoneNumber}
            onChange={onChange}
          />
        </span>
      </div>
      {helpText && <p class="help">{helpText}</p>}
      {errorText && <p class="help is-danger">{errorText}</p>}
    </div>
  );
}

export default FormPhoneField;
