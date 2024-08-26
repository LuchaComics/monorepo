import getCustomAxios from "../Helpers/customAxios";
import { camelizeKeys, decamelizeKeys, decamelize } from "humps";
import { DateTime } from "luxon";

import {
  CPS_PROJECTS_API_ENDPOINT,
  CPS_PROJECTS_SELECT_OPTIONS_API_ENDPOINT,
  CPS_PROJECT_API_ENDPOINT,
  CPS_PROJECT_STAR_OPERATION_API_ENDPOINT,
  CPS_PROJECT_ARCHIVE_OPERATION_API_ENDPOINT,
  CPS_PROJECT_CREATE_COMMENT_OPERATION_API_ENDPOINT,
  CPS_PROJECT_UPGRADE_OPERATION_API_ENDPOINT,
  CPS_PROJECT_DOWNGRADE_OPERATION_API_ENDPOINT,
  CPS_PROJECT_AVATAR_OPERATION_API_ENDPOINT,
  CPS_PROJECT_CHANGE_PASSWORD_OPERATION_API_ENDPOINT,
  CPS_PROJECT_CHANGE_2FA_OPERATION_API_URL,
} from "../Constants/API";

export function getProjectListAPI(
  filtersMap = new Map(),
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);

  // The following code will generate the query parameters for the url based on the map.
  let aURL = CPS_PROJECTS_API_ENDPOINT;
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
      console.log("getProjectListAPI | pre-fix | results:", data);
      if (
        data.results !== undefined &&
        data.results !== null &&
        data.results.length > 0
      ) {
        data.results.forEach((item, index) => {
          item.createdAt = DateTime.fromISO(item.createdAt).toLocaleString(
            DateTime.DATETIME_MED,
          );
          console.log(item, index);
        });
      }

      data.tenantID = data.tenantId;

      console.log("getProjectListAPI | post-fix | results:", data);

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function getProjectSelectOptionListAPI(
  filtersMap = new Map(),
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);

  // The following code will generate the query parameters for the url based on the map.
  let aURL = CPS_PROJECTS_SELECT_OPTIONS_API_ENDPOINT;
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
      console.log("getProjectSelectOptionListAPI | pre-fix | results:", data);
      if (
        data.results !== undefined &&
        data.results !== null &&
        data.results.length > 0
      ) {
        data.results.forEach((item, index) => {
          item.createdAt = DateTime.fromISO(item.createdAt).toLocaleString(
            DateTime.DATETIME_MED,
          );
          console.log(item, index);
        });
      }
      console.log("getProjectSelectOptionListAPI | post-fix | results:", data);

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function postProjectCreateAPI(
  data,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);

  // To Snake-case for API from camel-case in React.
  let decamelizedData = decamelizeKeys(data);

  // Minor fix.
  decamelizedData.tenant_id = data.TenantID;
  if (
    decamelizedData.tenant_id !== undefined &&
    decamelizedData.tenant_id !== null &&
    decamelizedData.tenant_id !== "" &&
    decamelizedData.tenant_id.length < 14
  ) {
    delete decamelizedData.tenant_id;
  }
  decamelizedData.cps_partnership_reason = data.CPSPartnershipReason;

  axios
    .post(CPS_PROJECTS_API_ENDPOINT, decamelizedData)
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

export function getProjectDetailAPI(
  projectID,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);
  axios
    .get(CPS_PROJECT_API_ENDPOINT.replace("{id}", projectID))
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // Bugfix
      data.tenantID = data.tenantId;

      // For debugging purposeso pnly.
      console.log(data);

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function putProjectUpdateAPI(
  data,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);

  // To Snake-case for API from camel-case in React.
  let decamelizedData = decamelizeKeys(data);

  // Minor fix.
  decamelizedData.tenant_id = decamelizedData.tenant_i_d;
  delete decamelizedData.tenant_i_d;

  if (
    decamelizedData.tenant_id !== undefined &&
    decamelizedData.tenant_id !== null &&
    decamelizedData.tenant_id !== "" &&
    decamelizedData.tenant_id.length < 14
  ) {
    delete decamelizedData.tenant_id;
  }
  decamelizedData.cps_partnership_reason = data.CPSPartnershipReason;

  axios
    .put(
      CPS_PROJECT_API_ENDPOINT.replace("{id}", decamelizedData.id),
      decamelizedData,
    )
    .then((successResponse) => {
      const responseData = successResponse.data;

      // Snake-case from API to camel-case for React.
      const data = camelizeKeys(responseData);

      // Bugfix
      data.tenantID = data.tenantId;

      // Return the callback data.
      onSuccessCallback(data);
    })
    .catch((exception) => {
      let errors = camelizeKeys(exception);
      onErrorCallback(errors);
    })
    .then(onDoneCallback);
}

export function deleteProjectAPI(
  id,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);
  axios
    .delete(CPS_PROJECT_API_ENDPOINT.replace("{id}", id))
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

export function postProjectCreateCommentOperationAPI(
  projectID,
  content,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);
  const data = {
    project_id: projectID,
    content: content,
  };
  axios
    .post(CPS_PROJECT_CREATE_COMMENT_OPERATION_API_ENDPOINT, data)
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

export function postProjectStarOperationAPI(
  projectID,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);
  const data = {
    project_id: projectID,
  };
  axios
    .post(CPS_PROJECT_STAR_OPERATION_API_ENDPOINT, data)
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

export function postArchiveProjectAPI(
  id,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);
  const data = {
    project_id: id,
  };
  axios
    .post(CPS_PROJECT_ARCHIVE_OPERATION_API_ENDPOINT, data)
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

export function postUpgradeProjectAPI(
  decamelizedData,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);
  axios
    .post(CPS_PROJECT_UPGRADE_OPERATION_API_ENDPOINT, decamelizedData)
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

export function postDowngradeProjectAPI(
  id,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);
  const data = {
    project_id: id,
  };
  axios
    .post(CPS_PROJECT_DOWNGRADE_OPERATION_API_ENDPOINT, data)
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

export function postAvatarProjectAPI(
  formdata,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);

  axios
    .post(CPS_PROJECT_AVATAR_OPERATION_API_ENDPOINT, formdata, {
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

// export function getProjectSelectOptionListAPI(filtersMap=new Map(), onSuccessCallback, onErrorCallback, onDoneCallback, onUnauthorizedCallback) {
//     const axios = getCustomAxios(onUnauthorizedCallback);
//
//     // The following code will generate the query parameters for the url based on the map.
//     let aURL = CPS_PROJECTS_SELECT_OPTIONS_API_ENDPOINT;
//     filtersMap.forEach(
//         (value, key) => {
//             let decamelizedkey = decamelize(key)
//             if (aURL.indexOf('?') > -1) {
//                 aURL += "&"+decamelizedkey+"="+value;
//             } else {
//                 aURL += "?"+decamelizedkey+"="+value;
//             }
//         }
//     )
//
//     axios.get(aURL).then((successResponse) => {
//         const responseData = successResponse.data;
//
//         // Snake-case from API to camel-case for React.
//         const data = camelizeKeys(responseData);
//
//         // Bugfixes.
//         console.log("getProjectSelectOptionListAPI | pre-fix | results:", data);
//         if (data.results !== undefined && data.results !== null && data.results.length > 0) {
//             data.results.forEach(
//                 (item, index) => {
//                     item.createdAt = DateTime.fromISO(item.createdAt).toLocaleString(DateTime.DATETIME_MED);
//                     console.log(item, index);
//                 }
//             )
//         }
//         console.log("getProjectSelectOptionListAPI | post-fix | results:", data);
//
//         // Return the callback data.
//         onSuccessCallback(data);
//     }).catch( (exception) => {
//         let errors = camelizeKeys(exception);
//         onErrorCallback(errors);
//     }).then(onDoneCallback);
// }
//

export function postProjectChangePasswordAPI(
  data,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);
  axios
    .post(CPS_PROJECT_CHANGE_PASSWORD_OPERATION_API_ENDPOINT, data)
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

export function postProjectChangeTwoFactorAuthenticationAPI(
  data,
  onSuccessCallback,
  onErrorCallback,
  onDoneCallback,
  onUnauthorizedCallback,
) {
  const axios = getCustomAxios(onUnauthorizedCallback);
  axios
    .post(CPS_PROJECT_CHANGE_2FA_OPERATION_API_URL, data)
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
