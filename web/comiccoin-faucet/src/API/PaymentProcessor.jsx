import getCustomAxios from "../Helpers/customAxios";
import { camelizeKeys, decamelizeKeys, decamelize } from "humps";
import { DateTime } from "luxon";

import {
  COMICCOIN_FAUCET_CREDIT_CREATE_STRIPE_CHECKOUT_SESSION_FOR_COMIC_SUBMISSION_API_ENDPOINT,
  COMICCOIN_FAUCET_CREDIT_COMPLETE_STRIPE_CHECKOUT_SESSION_API_ENDPOINT,
  COMICCOIN_FAUCET_CREDIT_PAYMENT_PROCESSOR_STRIPE_INVOICES_API_ENDPOINT,
  COMICCOIN_FAUCET_CREDIT_PAYMENT_PROCESSOR_SEND_SUBSCRIPTION_REQUEST_EMAIL_API_ENDPOINT,
} from "../Constants/API";

export function postCreateStripeCheckoutSessionURLForComicSubmissionAPI(
  comicSubmissionID,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
) {
  const axios = getCustomAxios();
  const postData = {};

  axios
    .post(
      COMICCOIN_FAUCET_CREDIT_CREATE_STRIPE_CHECKOUT_SESSION_FOR_COMIC_SUBMISSION_API_ENDPOINT.replace(
        "{id}",
        comicSubmissionID,
      ),
      postData,
    )
    .then((successResponse) => {
      const responseData = successResponse.data;

      console.log(
        "postCreateStripeCheckoutSessionURLForComicSubmissionAPI | response:",
        responseData,
      );

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

// export function getCompleteStripeSubscriptionCheckoutSessionAPI(sessionID, onSuccessCallback, onErrorCallback, onDoneCallback) {
//     const axios = getCustomAxios();
//     axios.get(COMICCOIN_FAUCET_CREDIT_COMPLETE_STRIPE_CHECKOUT_SESSION_API_ENDPOINT.replace("{sessionID}", sessionID)).then((successResponse) => {
//         const responseData = successResponse.data;
//
//         // Snake-case from API to camel-case for React.
//         const data = camelizeKeys(responseData);
//
//         // For debugging purposeso pnly.
//         console.log("completeStripeSubscriptionCheckoutSession: Response Data: ", data);
//
//         // Return the callback data.
//         onSuccessCallback(data);
//     }).catch( (exception) => {
//         let errors = camelizeKeys(exception);
//         onErrorCallback(errors);
//     }).then(onDoneCallback);
// }
//
// export function getPaymentProcessorStripeInvoiceListAPI(userID, cursor, pageSize, onSuccessCallback, onErrorCallback, onDoneCallback) {
//     const axios = getCustomAxios();
//
//     let aURL = COMICCOIN_FAUCET_CREDIT_PAYMENT_PROCESSOR_STRIPE_INVOICES_API_ENDPOINT.replace("{userID}",userID).replace("{cursor}",cursor).replace("{pageSize}",pageSize);
//
//     axios.get(aURL).then((successResponse) => {
//         const responseData = successResponse.data;
//
//         // Snake-case from API to camel-case for React.
//         const data = camelizeKeys(responseData);
//
//         // Bugfixes.
//         // console.log("getMemberListAPI | pre-fix | results:", data);
//         // if (data.results !== undefined && data.results !== null && data.results.length > 0) {
//         //     data.results.forEach(
//         //         (item, index) => {
//         //             item.createdAt = DateTime.fromISO(item.createdAt).toLocaleString(DateTime.DATETIME_MED);
//         //             console.log(item, index);
//         //         }
//         //     )
//         // }
//         // console.log("getMemberListAPI | post-fix | results:", data);
//
//         // Return the callback data.
//         onSuccessCallback(data);
//     }).catch( (exception) => {
//         let errors = camelizeKeys(exception);
//         onErrorCallback(errors);
//     }).then(onDoneCallback);
// }
//
// export function postPaymentProcessorSendSubscriptionRequestEmailAPI(userID, offerId, onSuccessCallback, onErrorCallback, onDoneCallback) {
//     const axios = getCustomAxios();
//     const data = {
//         member_id: userID,
//         offer_id: offerId,
//     };
//     axios.post(COMICCOIN_FAUCET_CREDIT_PAYMENT_PROCESSOR_SEND_SUBSCRIPTION_REQUEST_EMAIL_API_ENDPOINT, data).then((successResponse) => {
//         const responseData = successResponse.data;
//
//         // Snake-case from API to camel-case for React.
//         const data = camelizeKeys(responseData);
//
//         // Return the callback data.
//         onSuccessCallback(data);
//     }).catch( (exception) => {
//         let errors = camelizeKeys(exception);
//         onErrorCallback(errors);
//     }).then(onDoneCallback);
// }
