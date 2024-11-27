// Version 1.0.1
import React, { useState, useEffect } from "react";

/*
    // EXAMPLE USAGE:


    // Step 1: Create your state.
    const [birthday, setBirthday] = useState(null)

    // ...

    // Step 2: Use this componen.
    <FormAlternateDateField
        label="Birthday"
        name="birthday"
        placeholder="Text input"
        value={birthday}
        helpText=""
        onChange={(date)=>setBirthDay(date)}
        errorText={errors && errors.birthday}
        isRequired={true}
        maxWidth="180px"
        maxDate={new Date()}
    />

    // ...

    // NOTES:
    // Special thanks to the following URL:
    // https://www.smashingmagazine.com/2021/05/frustrating-design-patterns-birthday-picker/
 */
function FormAlternateDateField({
  label,
  name,
  placeholder,
  value,
  errorText,
  validationText,
  helpText,
  onChange,
  maxWidth,
  disabled = false,
  withPortal = true,
  filterDate = null,
  minDate = null,
  maxDate = null,
  monthsShown = 1,
}) {
  ////
  //// Component states.
  ////

  const [day, setDay] = useState(null);
  const [month, setMonth] = useState(null);
  const [year, setYear] = useState(null);

  useEffect(() => {
    if (value === "0001-01-01T00:00:00Z") {
      // Handle the case where 'value' is '0001-01-01T00:00:00Z'
      // You can avoid rendering or set the component to an initial state.
      return;
    }

    let dt = null;

    if (value !== undefined && value !== null && value !== "") {
      const valueDate = new Date(value);
      valueDate.setHours(0, 0, 0, 0);
      dt = valueDate;
    }

    setDay(dt ? dt.getDate() : null);
    setMonth(dt ? dt.getMonth() + 1 : null);
    setYear(dt ? dt.getFullYear() : null);
  }, [value]);

  ////
  //// Event handling.
  ////

  // Utility funciton to check if number is actually a number.
  const isNumeric = (str) => {
    // Use a regular expression to check if the string contains only numbers
    return /^[0-9]+$/.test(str);
  };

  const onDayChange = (d) => {
    // Convert the input value to a number
    const di = parseInt(d, 10);

    // If the input value is "0", set the day field to "0"
    if (d === "0") {
      setDay("0");
      return;
    }

    // If the input value is not a valid number or is empty, clear the day field
    if (isNaN(di) || d === "") {
      setDay("");
      return;
    }

    // Prevent future dates if max date exists
    if (
      maxDate &&
      year === maxDate.getFullYear() &&
      month === maxDate.getMonth() + 1
    ) {
      if (di > maxDate.getDate()) {
        setDay(maxDate.getDate());
        return;
      }
    }

    // If the input value is within the valid range, update the day field
    if (di >= 1 && di <= 31) {
      setDay(di);
      onDateChange(di, month, year);
    }
  };

  const onMonthChange = (m) => {
    // Convert the input value to a number
    const mi = parseInt(m, 10);

    // If the input value is "0", set the month field to "0"
    if (m === "0") {
      setMonth("0");
      return;
    }

    // If the input value is not a valid number or is empty, clear the month field
    if (isNaN(mi) || m === "") {
      setMonth("");
      return;
    }

    // Prevent future dates if max date exists
    if (maxDate && year === maxDate.getFullYear()) {
      if (mi > maxDate.getMonth() + 1) {
        setMonth(maxDate.getMonth() + 1);
        return;
      }
    }

    // If the input value is within the valid range, update the month field
    if (mi >= 1 && mi <= 12) {
      setMonth(mi);
      onDateChange(day, mi, year);
    }
  };

  const onYearChange = (y) => {
    // Convert into an integer.
    const yi = parseInt(y);

    // If user entered a value which is not a number, then remove.
    if (!isNumeric(y)) {
      setYear("");
      return;
    }

    // Prevent futuredates if max date exists.
    if (maxDate) {
      if (y > maxDate.getFullYear()) {
        setYear(maxDate.getFullYear());
        return;
      }
    }

    // Save the validated year.
    setYear(yi);

    // Update the main date.
    onDateChange(day, month, yi);
  };

  const onDateChange = (d, m, y) => {
    if (d > 0 && m > 0 && y > 0) {
      if (y > 1000) {
        // Convert day, month, and year into a JavaScript Date object
        const date = new Date(y, m - 1, d); // Months are 0-based, so subtract 1 from the month

        console.log("Date formatted successfully:", date);
        onChange(date); // Set the date up the component to the parent because we successfully generated a date.
      }
    }
  };

  ////
  //// Rendering
  ////

  let classNameText = "input";
  if (errorText) {
    classNameText = "input is-danger";
  }

  return (
    <div className="field pb-4">
      <label className="label">{label}</label>
      <div className="control">
        <div className="columns is-mobile">
          <div className="column" style={{ maxWidth: "75px" }}>
            <span className="is-fullwidth">
              Day
              <input
                className={classNameText}
                type="text"
                placeholder="DD"
                maxlength="2"
                value={day}
                onChange={(e) => {
                  onDayChange(e.target.value);
                }}
                disabled={disabled}
              />
            </span>
          </div>
          <div className="column" style={{ maxWidth: "75px" }}>
            <span className="is-fullwidth">
              Month
              <input
                className={classNameText}
                type="text"
                placeholder="MM"
                maxlength="2"
                value={month}
                onChange={(e) => {
                  onMonthChange(e.target.value);
                }}
                disabled={disabled}
              />
            </span>
          </div>
          <div className="column" style={{ maxWidth: "90px" }}>
            <span className="is-fullwidth">
              Year
              <input
                className={classNameText}
                type="text"
                placeholder="YYYY"
                maxlength="4"
                value={year}
                onChange={(e) => {
                  onYearChange(e.target.value);
                }}
                disabled={disabled}
              />
            </span>
          </div>
        </div>
      </div>
      {helpText && <p className="help">{helpText}</p>}
      {errorText && <p className="help is-danger">{errorText}</p>}
    </div>
  );
}

export default FormAlternateDateField;
