import getCustomAxios from "../Helpers/customAxios";
import { camelizeKeys, decamelizeKeys, decamelize } from "humps";
import { DateTime } from "luxon";

import {
  COMICCOIN_FAUCET_ATTACHMENTS_API_ENDPOINT,
  COMICCOIN_FAUCET_ATTACHMENT_API_ENDPOINT,
} from "../Constants/API";

export function getAttachmentListAPI(
  filtersMap = new Map(),
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);

  // The following code will generate the query parameters for the url based on the map.
  let aURL = COMICCOIN_FAUCET_ATTACHMENTS_API_ENDPOINT;
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
      // console.log("getAttachmentListAPI | pre-fix | results:", data);
      if (
        data.results !== undefined &&
        data.results !== null &&
        data.results.length > 0
      ) {
        data.results.forEach((item, index) => {
          item.issueCoverDate = DateTime.fromISO(
            item.issueCoverDate,
          ).toLocaleString(DateTime.DATETIME_MED);
          item.createdAt = DateTime.fromISO(item.createdAt).toLocaleString(
            DateTime.DATETIME_MED,
          );
          // console.log(item, index);
        });
      }
      // console.log("getAttachmentListAPI | post-fix | results:", data);

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function postAttachmentCreateAPI(
  filename,
  mimeType,
  formData,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  // Defensive code.
  if (filename === undefined || filename === null || filename === "") {
    onErrorCallback({"filename": "does not exist: "+filename});
    return;
  }
  if (mimeType === undefined || mimeType === null || mimeType === "") {
    onErrorCallback({"mimeType": "does not exist: "+mimeType});
    return;
  }

  const axios = getCustomAxios(onUnauthorizedCallback);

  axios
    .post(COMICCOIN_FAUCET_ATTACHMENTS_API_ENDPOINT, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        'Content-Disposition': `form-data; name="file"; filename="${encodeURIComponent(filename)}"`,
      },
    })
    .then((successResponse) => {
      const responseData = successResponse.data;
      const data = camelizeKeys(responseData);
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function getAttachmentDetailAPI(
  submissionID,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);
  axios
    .get(COMICCOIN_FAUCET_ATTACHMENT_API_ENDPOINT.replace("{id}", submissionID))
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

export function putAttachmentUpdateAPI(
  id,
  data,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);

  axios
    .put(COMICCOIN_FAUCET_ATTACHMENT_API_ENDPOINT.replace("{id}", id), data, {
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

export function deleteAttachmentAPI(
  id,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);
  axios
    .delete(COMICCOIN_FAUCET_ATTACHMENT_API_ENDPOINT.replace("{id}", id))
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
