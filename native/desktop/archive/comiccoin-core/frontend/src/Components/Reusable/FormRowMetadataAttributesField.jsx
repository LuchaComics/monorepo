import React from "react";
import { Link } from "react-router-dom";

function FormRowMetadataAttributesField(props) {
  const { label, attributes, helpText, type = "text" } = props;
  return (
    <div class="field pb-4">
      <label class="label has-text-black">{label}</label>
      <div class="control">
        <p>

        {attributes !== undefined &&
        attributes !== null &&
        attributes !== "" &&
        attributes.length > 0 ? (
          <>
            <table class="table">
              <thead className="is-size-7">
                <tr>
                  <th>
                    <abbr title="Display Type">Display Type</abbr>
                  </th>
                  <th>Trait Type</th>
                  <th>Value</th>
                </tr>
              </thead>
              <tbody className="is-size-7">
                {attributes &&
                  attributes.map(function (datum, i) {
                    return (
                      <tr>
                        <th>{datum.display_type}</th>
                        <td>{datum.trait_type}</td>
                        <td>{datum.value}</td>
                      </tr>
                    );
                  })}
              </tbody>
            </table>
          </>
        ) : (
          <>
          </>
        )}

        </p>
        {helpText !== undefined && helpText !== null && helpText !== "" && (
          <p class="help">{helpText}</p>
        )}
      </div>
    </div>
  );
}

export default FormRowMetadataAttributesField;
