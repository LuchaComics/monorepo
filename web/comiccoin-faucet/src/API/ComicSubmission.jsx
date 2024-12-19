import getCustomAxios from "../Helpers/customAxios";
import { camelizeKeys, decamelizeKeys, decamelize } from "humps";
import { DateTime } from "luxon";

import {
  COMICCOIN_FAUCET_COMIC_SUBMISSIONS_API_ENDPOINT,
  COMICCOIN_FAUCET_COMIC_SUBMISSION_API_ENDPOINT,
  COMICCOIN_FAUCET_COMIC_SUBMISSIONS_COUNT_BY_FILTER_API_ENDPOINT,
  COMICCOIN_FAUCET_COMIC_SUBMISSIONS_COUNT_COINS_REWARD_BY_FILTER_API_ENDPOINT,
  COMICCOIN_FAUCET_COMIC_SUBMISSIONS_COUNT_TOTAL_CREATED_TODAY_BY_USER_API_ENDPOINT,
  COMICCOIN_FAUCET_COMIC_SUBMISSION_CUSTOMER_SWAP_OPERATION_API_ENDPOINT,
  COMICCOIN_FAUCET_COMIC_SUBMISSION_CREATE_COMMENT_OPERATION_API_ENDPOINT,
  COMICCOIN_FAUCET_COMIC_SUBMISSION_FILE_ATTACHMENTS_API_ENDPOINT,
  COMICCOIN_FAUCET_COMIC_SUBMISSIONS_JUDGE_OPERATION_API_ENDPOINT,
  COMICCOIN_FAUCET_COMIC_SUBMISSIONS_TOTAL_COINS_AWARDED_API_ENDPOINT
} from "../Constants/API";

export function getComicSubmissionListAPI(
  filtersMap = new Map(),
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);

  // The following code will generate the query parameters for the url based on the map.
  let aURL = COMICCOIN_FAUCET_COMIC_SUBMISSIONS_API_ENDPOINT;
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
      // console.log("getComicSubmissionListAPI | pre-fix | results:", data);
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
      // console.log("getComicSubmissionListAPI | post-fix | results:", data);

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function postComicSubmissionCreateAPI(
  data,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);

  // console.log("postComicSubmissionCreateAPI | pre-modified | data:", data);

  // To Snake-case for API from camel-case in React.
  let decamelizedData = decamelizeKeys(data);

  // Minor bugfixes.
  decamelizedData.store_id = decamelizedData.store_i_d;
  delete decamelizedData.store_i_d;
  decamelizedData.customer_id = data.customerID;
  delete decamelizedData.customer_i_d;

  // if (data.issueCoverDate !==undefined && data.issueCoverDate !==null && data.issueCoverDate !=="") {
  //     decamelizedData.issue_cover_date = new Date(data.issueCoverDate).toISOString();
  // }

  // console.log("postComicSubmissionCreateAPI | post-modified | data:", decamelizedData);

  axios
    .post(COMICCOIN_FAUCET_COMIC_SUBMISSIONS_API_ENDPOINT, decamelizedData)
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // Minor bugfix.
      data.showsSignsOfTamperingOrRestoration = parseInt(
        data.showsSignsOfTamperingOrRestoration,
      );
      // data.issueCoverDate = DateTime.fromISO(data.issueCoverDate).toJSDate();

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function getComicSubmissionDetailAPI(
  submissionID,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);
  axios
    .get(COMICCOIN_FAUCET_COMIC_SUBMISSION_API_ENDPOINT.replace("{id}", submissionID))
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // Minor bugfix.
      data.showsSignsOfTamperingOrRestoration = parseInt(
        data.showsSignsOfTamperingOrRestoration,
      );
      // data.issueCoverDate = DateTime.fromISO(data.issueCoverDate).toJSDate();

      // // For debugging purposeso pnly.
      // console.log(data);

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function getComicSubmissionsCountByFilterAPI(
  filtersMap = new Map(),
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  // The following code will generate the query parameters for the url based on the map.
  let aURL = COMICCOIN_FAUCET_COMIC_SUBMISSIONS_COUNT_BY_FILTER_API_ENDPOINT;
  filtersMap.forEach((value, key) => {
    let decamelizedkey = decamelize(key);
    if (aURL.indexOf("?") > -1) {
      aURL += "&" + decamelizedkey + "=" + value;
    } else {
      aURL += "?" + decamelizedkey + "=" + value;
    }
  });

  const axios = getCustomAxios(onUnauthorizedCallback);
  axios
    .get(aURL)
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // // For debugging purposeso pnly.
      // console.log(data);

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}


export function getComicSubmissionsCountCoinsRewardByFilterAPI(
  filtersMap = new Map(),
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  // The following code will generate the query parameters for the url based on the map.
  let aURL = COMICCOIN_FAUCET_COMIC_SUBMISSIONS_COUNT_COINS_REWARD_BY_FILTER_API_ENDPOINT;
  filtersMap.forEach((value, key) => {
    let decamelizedkey = decamelize(key);
    if (aURL.indexOf("?") > -1) {
      aURL += "&" + decamelizedkey + "=" + value;
    } else {
      aURL += "?" + decamelizedkey + "=" + value;
    }
  });

  const axios = getCustomAxios(onUnauthorizedCallback);
  axios
    .get(aURL)
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // // For debugging purposeso pnly.
      // console.log(data);

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function getComicSubmissionsCountTotalCreatedTodayByUserAPI(
  userID,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);
  axios
    .get(COMICCOIN_FAUCET_COMIC_SUBMISSIONS_COUNT_TOTAL_CREATED_TODAY_BY_USER_API_ENDPOINT.replace("{user_id}", userID))
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // // For debugging purposeso pnly.
      // console.log(data);

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function putComicSubmissionUpdateAPI(
  decamelizedData,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);

  axios
    .put(
      COMICCOIN_FAUCET_COMIC_SUBMISSION_API_ENDPOINT.replace("{id}", decamelizedData.id),
      decamelizedData,
    )
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // Minor bugfix.
      data.showsSignsOfTamperingOrRestoration = parseInt(
        data.showsSignsOfTamperingOrRestoration,
      );
      // data.issueCoverDate = DateTime.fromISO(data.issueCoverDate).toJSDate();

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function deleteComicSubmissionAPI(
  id,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);
  axios
    .delete(COMICCOIN_FAUCET_COMIC_SUBMISSION_API_ENDPOINT.replace("{id}", id))
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

export function postComicSubmissionCustomerSwapOperationAPI(
  submissionID,
  customerID,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);
  const data = {
    submission_id: submissionID,
    customer_id: customerID,
  };
  axios
    .post(COMICCOIN_FAUCET_COMIC_SUBMISSION_CUSTOMER_SWAP_OPERATION_API_ENDPOINT, data)
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

export function postComicSubmissionCreateCommentOperationAPI(
  submissionID,
  content,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);
  const data = {
    submission_id: submissionID,
    content: content,
  };
  axios
    .post(COMICCOIN_FAUCET_COMIC_SUBMISSION_CREATE_COMMENT_OPERATION_API_ENDPOINT, data)
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

export function postComicSubmissionFileAttachmentCreateAPI(
  submissionID,
  formdata,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);

  axios
    .post(
      COMICCOIN_FAUCET_COMIC_SUBMISSION_FILE_ATTACHMENTS_API_ENDPOINT.replace(
        "{id}",
        submissionID,
      ),
      formdata,
      {
        headers: {
          "Content-Type": "multipart/form-data",
          Accept: "application/json",
        },
      },
    )
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


export function postComicSubmissionJudgementOperationAPI(
  data,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);

  // console.log("postComicSubmissionCreateAPI | pre-modified | data:", data);

  // To Snake-case for API from camel-case in React.
  let decamelizedData = decamelizeKeys(data);

  // Minor bugfixes.
  decamelizedData.store_id = decamelizedData.store_i_d;
  delete decamelizedData.store_i_d;
  decamelizedData.customer_id = data.customerID;
  delete decamelizedData.customer_i_d;

  // if (data.issueCoverDate !==undefined && data.issueCoverDate !==null && data.issueCoverDate !=="") {
  //     decamelizedData.issue_cover_date = new Date(data.issueCoverDate).toISOString();
  // }

  // console.log("postComicSubmissionCreateAPI | post-modified | data:", decamelizedData);

  axios
    .post(COMICCOIN_FAUCET_COMIC_SUBMISSIONS_JUDGE_OPERATION_API_ENDPOINT, decamelizedData)
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // Minor bugfix.
      data.showsSignsOfTamperingOrRestoration = parseInt(
        data.showsSignsOfTamperingOrRestoration,
      );
      // data.issueCoverDate = DateTime.fromISO(data.issueCoverDate).toJSDate();

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function getComicSubmissionsTotalCoinsAwardedAPI(
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  // The following code will generate the query parameters for the url based on the map.
  let aURL = COMICCOIN_FAUCET_COMIC_SUBMISSIONS_TOTAL_COINS_AWARDED_API_ENDPOINT;
  const axios = getCustomAxios(onUnauthorizedCallback);
  axios
    .get(aURL)
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // // For debugging purposeso pnly.
      // console.log(data);

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}
