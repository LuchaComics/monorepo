export function prettyFormattedDateFromISO(isoDateString) {
  try {
    // Create a JavaScript Date object from the input string
    const date = new Date(isoDateString);

    // Extract the individual date components
    const year = date.getFullYear();
    const month = date.getMonth() + 1; // Adding 1 because getMonth() is zero-based
    const day = date.getDate();

    // Format the date as "MM/DD/YYYY"
    const formattedDate = `${month}/${day}/${year}`;

    console.log(formattedDate);
    return formattedDate;
  } catch (err) {
    return "Invalid ISO Date";
  }
}

export function prettyFormattedDateTimeFromISO(isoDateTimeString) {
  try {
    // Create a JavaScript Date object from the input string
    const localDate = new Date(isoDateTimeString);

    // Convert the local date and time to UTC
    const utcDate = new Date(
      Date.UTC(
        localDate.getUTCFullYear(),
        localDate.getUTCMonth(),
        localDate.getUTCDate(),
        localDate.getUTCHours(),
        localDate.getUTCMinutes(),
        localDate.getUTCSeconds(),
      ),
    );

    // Format the UTC date as "Mon DD, YYYY at hh:mm AM/PM"
    const formattedDate = new Intl.DateTimeFormat("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
      hour: "numeric",
      minute: "2-digit",
      hour12: true,
    }).format(utcDate);

    return formattedDate;
  } catch (err) {
    return "Invalid ISO Date/Time";
  }
}

export function prettyFormattedTimeFromISO(isoDateTimeString) {
  // Create a JavaScript Date object from the input string
  const date = new Date(isoDateTimeString);

  // Extract the hour and minute components
  const hours = date.getHours();
  const minutes = date.getMinutes();

  // Determine if it's AM or PM
  const period = hours >= 12 ? "PM" : "AM";

  // Convert hours from 24-hour format to 12-hour format
  const hours12 = hours % 12 || 12; // Handle 0 as 12 for 12:00 AM or 12:00 PM

  // Format the time as "h:mm AM/PM"
  const formattedTime = `${hours12}:${minutes < 10 ? "0" : ""}${minutes} ${period}`;

  return formattedTime;
}
