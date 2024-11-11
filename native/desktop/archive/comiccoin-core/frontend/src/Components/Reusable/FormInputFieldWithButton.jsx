import React from "react";
import { startCase } from "lodash";

function FormInputFieldWithButton({
  label,
  name,
  placeholder,
  value,
  type = "text",
  errorText,
  validationText,
  helpText,
  onChange,
  maxWidth,
  disabled = false,
  onButtonClick,
  buttonLabel,
}) {
  let classNameText = "input";
  if (errorText) {
    classNameText = "input is-danger is-medium";
  }
  return (
    <>
      <label class="label has-text-black">{label}:</label>
      <div class="field has-addons pb-4">
        <div class="control is-expanded" style={{ maxWidth: maxWidth }}>
          <input
            class={classNameText}
            name={name}
            type={type}
            placeholder={placeholder}
            value={value}
            onChange={onChange}
            disabled={disabled}
            autoComplete="off"
          />
        </div>
        <div class="control" style={{ width: maxWidth }}>
          <button
            class="button is-info is-medium"
            onClick={onButtonClick}
            disabled={disabled}
          >
            <div style={{ display: "flex", justifyContent: "center", alignItems: "center", height: "100%" }}>
              {buttonLabel}
            </div>
          </button>
        </div>
        {errorText && <p class="help is-danger">{errorText}</p>}
        {helpText && <p class="help">{helpText}</p>}
      </div>
    </>
  );
}

export default FormInputFieldWithButton;
