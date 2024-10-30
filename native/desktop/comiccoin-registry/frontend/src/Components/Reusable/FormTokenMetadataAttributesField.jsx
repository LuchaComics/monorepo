import React, { useState } from "react";
import { startCase } from "lodash";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPlus, faMinus } from "@fortawesome/free-solid-svg-icons";

import FormInputField from "./FormInputField";
import FormSelectField from "./FormSelectField";

export const DISPLAY_TYPE_OPTIONS = [
  { value: "boost_number", label: "Boost Number" },
  { value: "boost_percentage", label: "Boost Percentage" },
  { value: "number", label: "Number" },
  { value: "date", label: "Date" },
  { value: "string", label: "String" },
];

export const SIGNATURE_ROLES_WITH_EMPTY_OPTIONS = [
  { value: "", label: "Please select" }, // EMPTY OPTION
  ...DISPLAY_TYPE_OPTIONS,
];

/*
    DATA-STRUCTURE
    ---------------
    data - needs to be an array of dictionary objects. For example:
        [
            {
                "displayType": "Creator",
                "traitType": "Frank Herbert"
            }
        ].

    FUNCTIONS
    ---------------
    onDataChange - needs to look something like this in your JSX:
        onDataChange={(data)=>onDataChange(data)}

*/
function FormTokenMetadataAttributesField({
  data = [],
  onDataChange = null,
  disabled = false,
  helpText = "",
}) {
  ////
  //// Component states.
  ////

  const [showAddModal, setShowAddModal] = useState(false);
  const [errors, setErrors] = useState({});
  const [displayType, setDisplayType] = useState("");
  const [traitType, setName] = useState("");
  const [value, setValue] = useState("");

  ////
  //// Event handling.
  ////

  const onSubmitClick = () => {
    console.log("onSubmitClick...");
    let newErrors = {};
    if (displayType === undefined || displayType === null || displayType === "") {
      newErrors["displayType"] = "missing value";
    }
    if (traitType === undefined || traitType === null || traitType === "") {
      newErrors["traitType"] = "missing value";
    }
    if (value === undefined || value === null || value === "") {
      newErrors["value"] = "missing value";
    }
    if (Object.keys(newErrors).length === 0) {
      // Make a copy of the "array of strings" into a mutable array.
      let copyOfArr = [];
      if (data !== null) {
          copyOfArr = [...data];
      }

      // Update record.
      copyOfArr.push({ display_type: displayType, trait_type: traitType, value: value });

      // Run callback.
      onDataChange(copyOfArr);

      // Reset errors.
      setErrors({});

      // Reset fields and close the modal.
      setName("");
      setDisplayType("");
      setShowAddModal(false);
      return;
    }
    setErrors(newErrors);
  };

  const onRemoveRowClick = (i) => {
    // For debugging purposes.
    console.log(i);

    // Make a copy of the "array of strings" into a mutable array.
    const copyOfArr = [...data];

    // Delete record.
    const x = copyOfArr.splice(i, 1);

    // For debugging purposes.
    console.log(x);

    // Save
    onDataChange(copyOfArr);
  };

  ////
  //// Component rendering.
  ////

  // Render the JSX component with the data.
  return (
    <>
      <div class={`modal ${showAddModal && "is-active"}`}>
        <div class="modal-background"></div>
        <div class="modal-card">
          <header class="modal-card-head">
            <p class="modal-card-title">Add Token Metadata Attribute</p>
            <button
              class="delete"
              aria-label="close"
              onClick={(e) => setShowAddModal(!showAddModal)}
            ></button>
          </header>
          <section class="modal-card-body">
            <FormSelectField
                label="Display Type"
                name="displayType"
                placeholder="Pick displayType"
                selectedValue={displayType}
                errorText={errors && errors.displayType}
                helpText=""
                onChange={(e) => setDisplayType(e.target.value)}
                options={SIGNATURE_ROLES_WITH_EMPTY_OPTIONS}
                isRequired={true}
                maxWidth="220px"
            />
            <FormInputField
              label="Trait Type"
              name="traitType"
              placeholder="Text input"
              value={traitType}
              errorText={errors && errors.traitType}
              helpText=""
              onChange={(e) => setName(e.target.value)}
              isRequired={true}
              maxWidth="380px"
            />

            <FormInputField
              label="Value"
              name="value"
              placeholder="Text input"
              value={value}
              errorText={errors && errors.value}
              helpText=""
              onChange={(e) => setValue(e.target.value)}
              isRequired={true}
              maxWidth="380px"
            />
          </section>
          <footer class="modal-card-foot">
            <button class="button is-success" onClick={onSubmitClick}>
              Save changes
            </button>
            <button
              class="button"
              onClick={(e) => setShowAddModal(!showAddModal)}
            >
              Cancel
            </button>
          </footer>
        </div>
      </div>

      <div class="pb-4">
        <label class="label has-text-black">
         Attributes (Optional)
          {/*<button class="button is-success is-small" onClick={onAddListInputFieldClick} disabled={disabled}><FontAwesomeIcon className="fas" icon={faPlus} /></button>*/}
        </label>

        {data !== undefined &&
        data !== null &&
        data !== "" &&
        data.length > 0 ? (
          <>
            <table class="table">
              <thead className="is-size-7">
                <tr>
                  <th>
                    <abbr title="Display Type">Display Type</abbr>
                  </th>
                  <th>Trait Type</th>
                  <th>Value</th>
                  {disabled === false && (
                    <th>
                      <button
                        class="button is-success is-small"
                        onClick={(m) => setShowAddModal(!showAddModal)}
                      >
                        <FontAwesomeIcon className="fas" icon={faPlus} />
                        &nbsp;Add
                      </button>
                    </th>
                  )}
                </tr>
              </thead>
              <tbody className="is-size-7">
                {data &&
                  data.map(function (datum, i) {
                    return (
                      <tr>
                        <th>{datum.display_type}</th>
                        <td>{datum.trait_type}</td>
                        <td>{datum.value}</td>
                        {disabled === false && (
                          <td>
                            <button
                              class="button is-danger is-small"
                              onClick={(n) => onRemoveRowClick(i)}
                            >
                              <FontAwesomeIcon className="fas" icon={faMinus} />
                              &nbsp;Delete
                            </button>
                          </td>
                        )}
                      </tr>
                    );
                  })}
              </tbody>
            </table>
          </>
        ) : (
          <>
            <button
              class="button is-primary is-small"
              onClick={(e) => setShowAddModal(true)}
              disabled={disabled}
            >
              <FontAwesomeIcon className="fas" icon={faPlus} />
              &nbsp;Add Attribute
            </button>
          </>
        )}
        <p class="help">{helpText}</p>
      </div>
    </>
  );
}

export default FormTokenMetadataAttributesField;
