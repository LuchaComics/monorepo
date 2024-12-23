import getCustomAxios from "../Helpers/customAxios";
import { camelizeKeys, decamelizeKeys } from "humps";
import {
  COMICCOIN_FAUCET_PROFILE_API_ENDPOINT,
  COMICCOIN_FAUCET_PROFILE_CHANGE_PASSWORD_API_ENDPOINT,
  COMICCOIN_FAUCET_PROFILE_WALLET_ADDRESS_API_ENDPOINT,
  COMICCOIN_FAUCET_PROFILE_APPLY_FOR_VERIFICATION_API_ENDPOINT
} from "../Constants/API";

export function getProfileDetailAPI(
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);
  axios
    .get(COMICCOIN_FAUCET_PROFILE_API_ENDPOINT)
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

export function putProfileUpdateAPI(
  data,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);

  // To Snake-case for API from camel-case in React.
  let decamelizedData = decamelizeKeys(data);

  axios
    .put(COMICCOIN_FAUCET_PROFILE_API_ENDPOINT, decamelizedData)
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

export function putProfileChangePasswordAPI(
  decamelizedData,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);

  axios
    .put(COMICCOIN_FAUCET_PROFILE_CHANGE_PASSWORD_API_ENDPOINT, decamelizedData)
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

export function putProfileWalletAddressAPI(
  decamelizedData,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);

  axios
    .put(COMICCOIN_FAUCET_PROFILE_WALLET_ADDRESS_API_ENDPOINT, decamelizedData)
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      console.log("putProfileWalletAddressAPI: exception:", exception);
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function putProfileApplyForVerificationAPI(
  data,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);

  data.howDidYouHearAboutUs = parseInt(data.howDidYouHearAboutUs);
  data.howLongCollectingComicBooksForGrading = parseInt(data.howLongCollectingComicBooksForGrading);
  data.hasPreviouslySubmittedComicBookForGrading = parseInt(data.hasPreviouslySubmittedComicBookForGrading);
  data.hasOwnedGradedComicBooks = parseInt(data.hasOwnedGradedComicBooks);
  data.hasRegularComicBookShop = parseInt(data.hasRegularComicBookShop);
  data.hasPreviouslyPurchasedFromAuctionSite = parseInt(data.hasPreviouslyPurchasedFromAuctionSite);
  data.hasPreviouslyPurchasedFromFacebookMarketplace = parseInt(data.hasPreviouslyPurchasedFromFacebookMarketplace);
  data.hasRegularlyAttendedComicConsOrCollectibleShows = parseInt(data.hasRegularlyAttendedComicConsOrCollectibleShows);

  // To Snake-case for API from camel-case in React.
  let decamelizedData = decamelizeKeys(data);

  axios
    .put(COMICCOIN_FAUCET_PROFILE_APPLY_FOR_VERIFICATION_API_ENDPOINT, decamelizedData)
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      console.log("putProfileApplyForVerificationAPI: exception:", exception);
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}
