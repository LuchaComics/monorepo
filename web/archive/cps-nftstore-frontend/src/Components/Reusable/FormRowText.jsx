import React from "react";
import { Link } from "react-router-dom";

function FormRowText(props) {
  const { label, value, helpText, type = "text" } = props;
  return (
    <div className="field pb-4">
      <label className="label">{label}</label>
      <div className="control">
        <p>
          {/* Check if value is not undefined or null */}
          {value !== undefined && value !== null && value !== "" ? (
            <>
              {type === "text" && value}
              {type === "email" && <Link to={`mailto:${value}`}>{value}</Link>}
              {type === "phone" && <Link to={`tel:${value}`}>{value}</Link>}
            </>
          ) : (
            "-"
          )}
        </p>
        {helpText && <p className="help">{helpText}</p>}
      </div>
    </div>
  );
}

export default FormRowText;
