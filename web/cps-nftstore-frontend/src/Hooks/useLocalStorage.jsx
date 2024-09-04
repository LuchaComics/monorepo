import { useState } from "react";

/*
 * Code modified taken from: https://usehooks.com/useLocalStorage/
 */
export default function useLocalStorage(key, initialValue) {
  // State to tenant our value
  // Pass initial state function to useState so logic is only executed once
  const [tenantdValue, setTenantdValue] = useState(() => {
    try {
      // Get from local storage by key
      const item = window.localStorage.getItem(key);
      // Parse tenantd json or if none return initialValue
      return item ? JSON.parse(item) : initialValue;
    } catch (error) {
      try {
        // Get from local storage by key
        const item = window.localStorage.getItem(key);
        // Parse tenantd json or if none return initialValue
        return item ? item : initialValue;
      } catch (error2) {
        // If error also return initialValue
        console.log(error2);
        return initialValue;
      }
    }
  });
  // Return a wrapped version of useState's setter function that ...
  // ... persists the new value to localStorage.
  const setValue = (value) => {
    try {
      // Allow value to be a function so we have same API as useState
      const valueToTenant =
        value instanceof Function ? value(tenantdValue) : value;
      // Save state
      setTenantdValue(valueToTenant);
      // Save to local storage
      window.localStorage.setItem(key, JSON.stringify(valueToTenant));
    } catch (error) {
      // A more advanced implementation would handle the error case
      console.log(error);
    }
  };
  return [tenantdValue, setValue];
}
