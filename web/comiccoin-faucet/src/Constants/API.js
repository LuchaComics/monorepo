const HTTP_API_SERVER =
  process.env.REACT_APP_API_PROTOCOL + "://" + process.env.REACT_APP_API_DOMAIN;
export const COMICCOIN_FAUCET_API_BASE_PATH = "/api/v1";
export const COMICCOIN_FAUCET_VERSION_ENDPOINT = "version";
export const COMICCOIN_FAUCET_LOGIN_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/login";
export const COMICCOIN_FAUCET_REGISTER_BUSINESS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/register/business";
export const COMICCOIN_FAUCET_REGISTER_USER_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/register/user";
export const COMICCOIN_FAUCET_REFRESH_TOKEN_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/refresh-token";
export const COMICCOIN_FAUCET_EMAIL_VERIFICATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/verify";
export const COMICCOIN_FAUCET_LOGOUT_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/logout";
export const COMICCOIN_FAUCET_SUBMISSIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/submissions";
export const COMICCOIN_FAUCET_SUBMISSION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/submission/{id}";
export const COMICCOIN_FAUCET_SUBMISSION_CUSTOMER_SWAP_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/submissions/operation/set-user";
export const COMICCOIN_FAUCET_SUBMISSION_CREATE_COMMENT_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/submissions/operation/create-comment";
export const COMICCOIN_FAUCET_PROFILE_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/profile";
export const COMICCOIN_FAUCET_PROFILE_CHANGE_PASSWORD_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/profile/change-password";
export const COMICCOIN_FAUCET_PROFILE_WALLET_ADDRESS_API_ENDPOINT =
    HTTP_API_SERVER + "/api/v1/profile/wallet-password";

export const COMICCOIN_FAUCET_PROFILE_APPLY_FOR_VERIFICATION_API_ENDPOINT =
    HTTP_API_SERVER + "/api/v1/profile/apply-for-verification";

export const COMICCOIN_FAUCET_CUSTOMERS_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/customers";
export const COMICCOIN_FAUCET_CUSTOMER_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/customer/{id}";
export const COMICCOIN_FAUCET_CUSTOMER_CREATE_COMMENT_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/customers/operation/create-comment";
export const COMICCOIN_FAUCET_CUSTOMER_STAR_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/customers/operation/star";
export const COMICCOIN_FAUCET_FORGOT_PASSWORD_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/forgot-password";
export const COMICCOIN_FAUCET_PASSWORD_RESET_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/password-reset";
export const COMICCOIN_FAUCET_REGISTRY_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/cpsrn/{id}";
export const COMICCOIN_FAUCET_ORGANIZATIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/stores";
export const COMICCOIN_FAUCET_ORGANIZATIONS_SELECT_OPTIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/stores/select-options";
export const COMICCOIN_FAUCET_ORGANIZATIONS_PUBLIC_SELECT_OPTIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/public/stores-select-options";
export const COMICCOIN_FAUCET_ORGANIZATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/store/{id}";
export const COMICCOIN_FAUCET_ORGANIZATION_CREATE_COMMENT_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/stores/operation/create-comment";
export const COMICCOIN_FAUCET_USERS_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/users";
export const COMICCOIN_FAUCET_USERS_SELECT_OPTIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/users/select-options";
export const COMICCOIN_FAUCET_USER_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/user/{id}";
export const COMICCOIN_FAUCET_USER_CREATE_COMMENT_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/users/operation/create-comment";
export const COMICCOIN_FAUCET_USER_STAR_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/users/operation/star";

export const COMICCOIN_FAUCET_USER_PROFILE_VERIFICATION_JUDGE_OPERATION_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/users/operations/profile-verification-judge";

export const COMICCOIN_FAUCET_COMIC_SUBMISSIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/comic-submissions";
export const COMICCOIN_FAUCET_COMIC_SUBMISSION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/comic-submission/{id}";
export const COMICCOIN_FAUCET_COMIC_SUBMISSION_CUSTOMER_SWAP_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/comic-submissions/operation/set-customer";
export const COMICCOIN_FAUCET_COMIC_SUBMISSION_CREATE_COMMENT_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/comic-submissions/operation/create-comment";
export const COMICCOIN_FAUCET_COMIC_SUBMISSION_FILE_ATTACHMENTS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/comic-submission/{id}/file-attachments";

export const COMICCOIN_FAUCET_COMIC_SUBMISSIONS_COUNT_TOTAL_CREATED_TODAY_BY_USER_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/comic-submissions/count-total-created-today-by-user?user_id={user_id}";
export const COMICCOIN_FAUCET_COMIC_SUBMISSIONS_COUNT_BY_FILTER_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/comic-submissions/count";
export const COMICCOIN_FAUCET_COMIC_SUBMISSIONS_COUNT_COINS_REWARD_BY_FILTER_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/comic-submissions/count-coins-reward";
export const COMICCOIN_FAUCET_COMIC_SUBMISSIONS_JUDGE_OPERATION_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/comic-submissions/judge-operation";
export const COMICCOIN_FAUCET_COMIC_SUBMISSIONS_TOTAL_COINS_AWARDED_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/comic-submissions/total-coins-awarded";

export const COMICCOIN_FAUCET_USERS_COUNT_JOINED_THIS_WEEK_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/users/count-joined-this-week";

export const COMICCOIN_FAUCET_BALANCE_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/faucet/balance";

export const COMICCOIN_FAUCET_ATTACHMENTS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/attachments";
export const COMICCOIN_FAUCET_ATTACHMENT_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/attachment/{id}";

export const COMICCOIN_FAUCET_OFFERS_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/offers";
export const COMICCOIN_FAUCET_OFFER_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/offer/{id}";
export const COMICCOIN_FAUCET_OFFER_BY_SERVICE_TYPE_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/offer/service-type/{serviceType}";
export const COMICCOIN_FAUCET_OFFER_SELECT_OPTIONS_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/offers/select-options";
export const COMICCOIN_FAUCET_CREDITS_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/credits";
export const COMICCOIN_FAUCET_CREDIT_API_ENDPOINT = HTTP_API_SERVER + "/api/v1/credit/{id}";
export const COMICCOIN_FAUCET_CREDIT_COMPLETE_STRIPE_CHECKOUT_SESSION_API_ENDPOINT =
  HTTP_API_SERVER +
  "/api/v1/stripe/complete-checkout-session?session_id={sessionID}";
export const COMICCOIN_FAUCET_CREDIT_CREATE_STRIPE_CHECKOUT_SESSION_FOR_COMIC_SUBMISSION_API_ENDPOINT =
  HTTP_API_SERVER +
  "/api/v1/stripe/create-checkout-session-for-comic-submission/{id}";
export const COMICCOIN_FAUCET_CREDIT_CANCEL_SUBSCRIPTION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/stripe/cancel-subscription";
export const COMICCOIN_FAUCET_CREDIT_PAYMENT_PROCESSOR_STRIPE_INVOICES_API_ENDPOINT =
  HTTP_API_SERVER +
  "/api/v1/stripe/invoices?user_id={userID}&cursor={cursor}&page_size={pageSize}";
export const COMICCOIN_FAUCET_USER_PURCHASES_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/user-purchases";
export const COMICCOIN_FAUCET_2FA_GENERATE_OTP_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/otp/generate";
export const COMICCOIN_FAUCET_2FA_GENERATE_OTP_AND_QR_CODE_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/otp/generate-qr-code";
export const COMICCOIN_FAUCET_2FA_VERIFY_OTP_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/otp/verify";
export const COMICCOIN_FAUCET_2FA_VALIDATE_OTP_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/otp/validate";
export const COMICCOIN_FAUCET_2FA_DISABLED_OTP_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/otp/disable";
export const COMICCOIN_FAUCET_2FA_RECOVERY_OTP_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/otp/recovery";
export const COMICCOIN_FAUCET_USER_ARCHIVE_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/users/operation/archive";
export const COMICCOIN_FAUCET_USER_UPGRADE_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/users/operation/upgrade";
export const COMICCOIN_FAUCET_USER_DOWNGRADE_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/users/operation/downgrade";
export const COMICCOIN_FAUCET_USER_AVATAR_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/users/operation/avatar";
export const COMICCOIN_FAUCET_USER_CHANGE_2FA_OPERATION_API_URL =
  HTTP_API_SERVER + "/api/v1/users/operations/change-2fa";
export const COMICCOIN_FAUCET_USER_CHANGE_PASSWORD_OPERATION_API_ENDPOINT =
  HTTP_API_SERVER + "/api/v1/users/operations/change-password";
