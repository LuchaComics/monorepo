import { camelizeKeys, decamelizeKeys, decamelize } from "humps";
import { DateTime } from "luxon";
import axios from "axios";

import {
  CPS_IPFS_ADDFILE_API_ENDPOINT
} from "../Constants/API";


export function postIPFSAddFileAPI(
  apiKey,
  filename,
  file, // This should be a File object or Blob
  mimeType,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback
) {
  // Create a new Axios instance
  const customAxios = axios.create({
    headers: {
      "Authorization": `JWT ${apiKey}`,
      "Accept": "application/json",
    },
  });

  // Defensive code.
  if (filename === undefined || filename === null || filename === "") {
    onErrorCallback({"filename": "does not exist: "+filename});
    return;
  }
  if (mimeType === undefined || mimeType === null || mimeType === "") {
    onErrorCallback({"mimeType": "does not exist: "+mimeType});
    return;
  }



  customAxios.post(CPS_IPFS_ADDFILE_API_ENDPOINT, file, {
    headers: {
      "Content-Type": mimeType,
      "Content-Disposition": `attachment; filename=${filename}`, // Add filename here
    },
  })
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Convert snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      console.log('API Success Data:', data); // Debug statement

      // Return the callback data.
      if (onSuccessCallback) {
        onSuccessCallback(data);
      }
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception.response ? exception.response.data : exception);
      console.error('API Error:', errors); // Debug statement

      if (onErrorCallback) {
        onErrorCallback(errors);
      }
    })
    .finally(() => {
      if (onDoneCallback) {
        onDoneCallback();
      }
    });
}
