import React from "react";
import { Link } from "react-router-dom";
import { DateTime } from "luxon";

function DataDisplayRowText(props) {
  const { label, value, helpText, type = "text" } = props;
  return (
    <div class="field pb-4">
      <label class="label">{label}</label>
      <div class="control">
        <p>
          {value ? (
            <>
              {type === "text" && value}
              {type === "email" && <Link to={`mailto:${value}`}>{value}</Link>}
              {type === "phone" && <Link to={`tel:${value}`}>{value}</Link>}
              {type === "datetime" &&
                DateTime.fromISO(value).toLocaleString(DateTime.DATETIME_MED)}
              {type === "date" &&
                DateTime.fromISO(value).toLocaleString(DateTime.DATE_MED)}
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

export default DataDisplayRowText;
