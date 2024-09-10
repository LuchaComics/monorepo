import { camelizeKeys, decamelizeKeys, decamelize } from "humps";
import { DateTime } from "luxon";

import getCustomAxios from "../Helpers/customAxios";
import {
  CPS_NFT_METADATA_API_ENDPOINT,
  CPS_NFT_METADATUM_API_ENDPOINT,
  CPS_PIN_CONTENT_API_ENDPOINT
} from "../Constants/API";

export function getNFTMetadataListAPI(
  filtersMap = new Map(),
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);

  // The following code will generate the query parameters for the url based on the map.
  let aURL = CPS_NFT_METADATA_API_ENDPOINT;
  filtersMap.forEach((value, key) => {
    let decamelizedkey = decamelize(key);
    if (aURL.indexOf("?") > -1) {
      aURL += "&" + decamelizedkey + "=" + value;
    } else {
      aURL += "?" + decamelizedkey + "=" + value;
    }
  });

  axios
    .get(aURL)
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // Bugfixes.
      // console.log("getNFTListAPI | pre-fix | results:", data);
      if (
        data.results !== undefined &&
        data.results !== null &&
        data.results.length > 0
      ) {
        data.results.forEach((item, index) => {
          item.modifiedAt = DateTime.fromISO(
            item.modifiedAt,
          ).toLocaleString(DateTime.DATETIME_MED);
          item.createdAt = DateTime.fromISO(item.createdAt).toLocaleString(
            DateTime.DATETIME_MED,
          );
          // console.log(item, index);
        });
      }
      // console.log("getNFTListAPI | post-fix | results:", data);

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function postNFTMetadataCreateAPI(
  jsonData,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);

  axios
    .post(CPS_NFT_METADATA_API_ENDPOINT, jsonData)
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function getNFTMetadataDetailAPI(
  submissionID,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);
  axios
    .get(CPS_NFT_METADATUM_API_ENDPOINT.replace("{id}", submissionID))
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function putNFTMetadataUpdateAPI(
  id,
  data,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);

  axios
    .put(CPS_NFT_METADATUM_API_ENDPOINT.replace("{id}", id), data, {
      headers: {
        "Content-Type": "multipart/form-data",
        Accept: "application/json",
      },
    })
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function deleteNFTMetadataAPI(
  id,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);
  axios
    .delete(CPS_NFT_METADATUM_API_ENDPOINT.replace("{id}", id))
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function getNFTMetadataContentDetailAPI(
  requestID,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);
  axios
    .get(CPS_PIN_CONTENT_API_ENDPOINT.replace("{id}", requestID), { responseType: "blob", })
    .then((successResponse) => {
        console.log("getNFTContentDetailAPI: All response headers:", successResponse.headers);

        const contentDisposition = successResponse.headers['content-disposition'];
        let filename = ""; // Default filename is empty string - you will need to handle it in the upper functions that use this function.

        // Extract filename from Content-Disposition header if available
        if (contentDisposition && contentDisposition.indexOf('attachment') !== -1) {
            const filenameMatch = contentDisposition.match(/filename="(.+)"/);
            if (filenameMatch && filenameMatch.length === 2) {
                filename = filenameMatch[1];
            }
        }

        const contentType = successResponse.headers['content-type'];
        console.log("contentDisposition:", contentDisposition);
        console.log("contentType:", contentType);

        // Use `fileDownload` to download the file
        onSuccessCallback(successResponse.data, filename, contentType);
    })
    .catch((exception) => {
        let errors;
        if (exception.response) {
           errors = camelizeKeys(exception.response);
        } else {
            errors = exception.message ? { message: exception.message } : { message: 'An unknown error occurred' };
        }
        onErrorCallback(errors);
    })
    .then(onDoneCallback);
}
