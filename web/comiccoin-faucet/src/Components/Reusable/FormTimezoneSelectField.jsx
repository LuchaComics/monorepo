import React from "react";
import { useTimezoneSelect, allTimezones } from "react-timezone-select";

const FormTimezoneSelectField = ({
  label = "Timezone",
  name = "timezone",
  placeholder,
  selectedTimezone,
  setSelectedTimezone,
  errorText,
  helpText = "Please select the timezone that your business operates in.",
  disabled,
  isRequired = true
}) => {
  const labelStyle = "original";
  const timezones = {
    ...allTimezones,
    "America/Toronto": "Toronto",
  };

  const { options, parseTimezone } = useTimezoneSelect({
    labelStyle,
    timezones,
  });

  return (
    <div>
      <label
        htmlFor={name}
        className="block text-sm font-medium text-gray-700 mb-1"
      >
        {label} {isRequired && "*"}
      </label>
      <select
        id={name}
        name={name}
        disabled={disabled}
        onChange={(e) => setSelectedTimezone(parseTimezone(e.currentTarget.value).value)}
        className={`w-full h-11 px-4 py-2 border rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent bg-white ${
          errorText ? "border-red-500" : "border-gray-300"
        }`}
      >
        {options.map((option) => (
          <option
            key={option.value}
            value={option.value}
            selected={selectedTimezone === option.value}
          >
            {option.label}
          </option>
        ))}
      </select>
      {helpText && (
        <p className="mt-1 text-sm text-gray-500">{helpText}</p>
      )}
      {errorText && (
        <p className="mt-1 text-sm text-red-600">{errorText}</p>
      )}
    </div>
  );
};

export default FormTimezoneSelectField;
