import React from "react";

function FormTextChoiceRow(props) {
  const {
    label,
    value,
    helpText,
    opt1Value,
    opt1Label,
    opt2Value,
    opt2Label,
    opt3Value,
    opt3Label,
    opt4Value,
    opt4Label,
    opt5Value,
    opt5Label,
    opt6Value,
    opt6Label,
    opt7Value,
    opt7Label,
  } = props;
  return (
    <div class="field pb-4">
      <label class="label">{label}</label>
      <div class="control">
        {value === opt1Value && <span>{opt1Label}</span>}
        {value === opt2Value && <span>{opt2Label}</span>}
        {value === opt3Value && <span>{opt3Label}</span>}
        {value === opt4Value && <span>{opt4Label}</span>}
        {value === opt5Value && <span>{opt5Label}</span>}
        {value === opt6Value && <span>{opt6Label}</span>}
        {value === opt7Value && <span>{opt7Label}</span>}
        {helpText !== undefined && helpText !== null && helpText !== "" && (
          <p class="help">{helpText}</p>
        )}
      </div>
    </div>
  );
}

export default FormTextChoiceRow;
