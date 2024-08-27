const HTTP_API_SERVER =
  process.env.REACT_APP_API_PROTOCOL + "://" + process.env.REACT_APP_API_DOMAIN;
export const CPS_API_BASE_PATH = "/api/v1";
export const CPS_VERSION_ENDPOINT = "version";
export const CPS_LOGIN_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/login";
export const CPS_REGISTER_BUSINESS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/register/business";
export const CPS_REGISTER_CUSTOMER_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/register/customer";
export const CPS_REFRESH_TOKEN_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/refresh-token";
export const CPS_EMAIL_VERIFICATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/verify";
export const CPS_LOGOUT_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/logout";
export const CPS_SUBMISSIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/submissions";
export const CPS_SUBMISSION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/submission/{id}";
export const CPS_SUBMISSION_CUSTOMER_SWAP_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/submissions/operation/set-user";
export const CPS_SUBMISSION_CREATE_COMMENT_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/submissions/operation/create-comment";
export const CPS_PROFILE_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/profile";
export const CPS_PROFILE_CHANGE_PASSWORD_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/profile/change-password";
export const CPS_CUSTOMERS_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/customers";
export const CPS_CUSTOMER_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/customer/{id}";
export const CPS_CUSTOMER_CREATE_COMMENT_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/customers/operation/create-comment";
export const CPS_CUSTOMER_STAR_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/customers/operation/star";
export const CPS_FORGOT_PASSWORD_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/forgot-password";
export const CPS_PASSWORD_RESET_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/password-reset";
export const CPS_REGISTRY_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/cpsrn/{id}";
export const CPS_ORGANIZATIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/tenants";
export const CPS_ORGANIZATIONS_SELECT_OPTIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/tenants/select-options";
export const CPS_ORGANIZATIONS_PUBLIC_SELECT_OPTIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/public/tenants-select-options";
export const CPS_ORGANIZATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/tenant/{id}";
export const CPS_ORGANIZATION_CREATE_COMMENT_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/tenants/operation/create-comment";
export const CPS_USERS_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/users";
export const CPS_USERS_SELECT_OPTIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/users/select-options";
export const CPS_USER_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/user/{id}";
export const CPS_USER_CREATE_COMMENT_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/users/operation/create-comment";
export const CPS_USER_STAR_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/users/operation/star";
export const CPS_COMIC_SUBMISSIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/comic-submissions";
export const CPS_COMIC_SUBMISSION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/comic-submission/{id}";
export const CPS_COMIC_SUBMISSION_CUSTOMER_SWAP_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/comic-submissions/operation/set-customer";
export const CPS_COMIC_SUBMISSION_CREATE_COMMENT_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/comic-submissions/operation/create-comment";
export const CPS_COMIC_SUBMISSION_FILE_ATTACHMENTS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/comic-submission/{id}/file-pinobjects";
export const CPS_ATTACHMENTS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/pins";
export const CPS_ATTACHMENT_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/pin/{id}";
export const CPS_OFFERS_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/offers";
export const CPS_OFFER_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/offer/{id}";
export const CPS_OFFER_BY_SERVICE_TYPE_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/offer/service-type/{serviceType}";
export const CPS_OFFER_SELECT_OPTIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/offers/select-options";
export const CPS_CREDITS_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/credits";
export const CPS_CREDIT_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/credit/{id}";
export const CPS_CREDIT_COMPLETE_STRIPE_CHECKOUT_SESSION_API_ENDPOINT =
  HTTP_API_SERVER +
  "/api/v1/stripe/complete-checkout-session?session_id={sessionID}";
export const CPS_CREDIT_CREATE_STRIPE_CHECKOUT_SESSION_FOR_COMIC_SUBMISSION_API_ENDPOINT =
  HTTP_API_SERVER +
  "/api/v1/stripe/create-checkout-session-for-comic-submission/{id}";
export const CPS_CREDIT_CANCEL_SUBSCRIPTION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/stripe/cancel-subscription";
export const CPS_CREDIT_PAYMENT_PROCESSOR_STRIPE_INVOICES_API_ENDPOINT =
  HTTP_API_SERVER +
  "/api/v1/stripe/invoices?user_id={userID}&cursor={cursor}&page_size={pageSize}";
export const CPS_USER_PURCHASES_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/user-purchases";
export const CPS_2FA_GENERATE_OTP_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/otp/generate";
export const CPS_2FA_GENERATE_OTP_AND_QR_CODE_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/otp/generate-qr-code";
export const CPS_2FA_VERIFY_OTP_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/otp/verify";
export const CPS_2FA_VALIDATE_OTP_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/otp/validate";
export const CPS_2FA_DISABLED_OTP_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/otp/disable";
export const CPS_2FA_RECOVERY_OTP_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/otp/recovery";
export const CPS_USER_ARCHIVE_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/users/operation/archive";
export const CPS_USER_UPGRADE_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/users/operation/upgrade";
export const CPS_USER_DOWNGRADE_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/users/operation/downgrade";
export const CPS_USER_AVATAR_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/users/operation/avatar";
export const CPS_USER_CHANGE_2FA_OPERATION_API_URL =
  HTTP_API_SERVER + "/api/v1/users/operations/change-2fa";
export const CPS_USER_CHANGE_PASSWORD_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/users/operations/change-password";

// Project
  export const CPS_PROJECTS_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/projects";
  export const CPS_PROJECTS_SELECT_OPTIONS_API_ENDPOINT =
    HTTP_API_SERVER + "/api/v1/projects/select-options";
  export const CPS_PROJECT_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/project/{id}";
  export const CPS_PROJECT_CREATE_COMMENT_OPERATION_API_ENDPOINT =
    HTTP_API_SERVER + "/api/v1/projects/operation/create-comment";
  export const CPS_PROJECT_STAR_OPERATION_API_ENDPOINT =
    HTTP_API_SERVER + "/api/v1/projects/operation/star";
  export const CPS_PROJECT_PURCHASES_API_ENDPOINT =
    HTTP_API_SERVER + "/api/v1/project-purchases";
  export const CPS_PROJECT_ARCHIVE_OPERATION_API_ENDPOINT =
    HTTP_API_SERVER + "/api/v1/projects/operation/archive";
  export const CPS_PROJECT_UPGRADE_OPERATION_API_ENDPOINT =
    HTTP_API_SERVER + "/api/v1/projects/operation/upgrade";
  export const CPS_PROJECT_DOWNGRADE_OPERATION_API_ENDPOINT =
    HTTP_API_SERVER + "/api/v1/projects/operation/downgrade";
  export const CPS_PROJECT_AVATAR_OPERATION_API_ENDPOINT =
    HTTP_API_SERVER + "/api/v1/projects/operation/avatar";
  export const CPS_PROJECT_CHANGE_2FA_OPERATION_API_URL =
    HTTP_API_SERVER + "/api/v1/projects/operations/change-2fa";
  export const CPS_PROJECT_CHANGE_PASSWORD_OPERATION_API_ENDPOINT =
    HTTP_API_SERVER + "/api/v1/projects/operations/change-password";

// Pin
export const CPS_PIN_CONTENT_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/pin/{id}/content";
