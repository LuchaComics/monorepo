import React from "react";
import { useTimezoneSelect, allTimezones } from "react-timezone-select";

/*
    HOW TO USE:
    ------------

    // This code in the state declaration.
    const [timezone, setTimezone] = useState(Intl.DateTimeFormat().resolvedOptions().timeZone)

    ...

    // This code goes in your renderer.
    <FormTimezoneSelectField
        label="Timezone"
        name="timezone"
        placeholder="Text input"
        selectedTimezone={timezone}
        errorText={errors && errors.timezone}
        helpText="Please select the timezone that your business operates in."
        setTimezone={(value)=>setTimezone(value)}
        isRequired={true}
        maxWidth="280px"
/>

*/
function FormTimezoneSelectField({
  label,
  name,
  placeholder,
  selectedTimezone,
  setSelectedTimezone,
  errorText,
  validationText,
  helpText,
  disabled,
  maxWidth,
}) {
  // DEVELOPERS NOTE:

  const labelStyle = "original";
  const timezones = {
    ...allTimezones,
    "America/Toronto": "Toronto", // Add not included timezone which is still a valid timezone.
  };

  console.log("FormTimezoneSelectField | Input:", selectedTimezone);

  const { options, parseTimezone } = useTimezoneSelect({
    labelStyle,
    timezones,
  });

  return (
    <div class="field pb-4">
      <label class="label">{label}</label>
      <div class="control" style={{ maxWidth: maxWidth }}>
        <span class="select">
          <select
            onChange={(e) =>
              setSelectedTimezone(parseTimezone(e.currentTarget.value).value)
            }
          >
            {options.map((option) => (
              <option
                value={option.value}
                selected={selectedTimezone === option.value}
              >
                {option.label}
              </option>
            ))}
          </select>
        </span>
      </div>
      {helpText && <p class="help">{helpText}</p>}
      {errorText && <p class="help is-danger">{errorText}</p>}
    </div>
  );
}

export default FormTimezoneSelectField;
