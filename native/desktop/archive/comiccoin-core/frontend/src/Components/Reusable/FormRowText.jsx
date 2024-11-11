import React from "react";
import { Link } from "react-router-dom";

function FormRowText(props) {
  const { label, value, helpText, type = "text" } = props;
  return (
    <div class="field pb-4">
      <label class="label has-text-black">{label}</label>
      <div class="control">
        <p>
          {value ? (
            <>
              {type === "text" && value}
              {type === "email" && <Link to={`mailto:${value}`}>{value}</Link>}
              {type === "phone" && <Link to={`tel:${value}`}>{value}</Link>}
            </>
          ) : (
            "-"
          )}
        </p>
        {helpText !== undefined && helpText !== null && helpText !== "" && (
          <p class="help">{helpText}</p>
        )}
      </div>
    </div>
  );
}

export default FormRowText;
