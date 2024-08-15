import { React, useState } from "react";
import "bulma/css/bulma.min.css";
import "./index.css";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { RecoilRoot } from "recoil";

//--------------//
// Admin Portal //
//--------------//

// Offer
import AdminOfferUpdate from "./Components/Admin/Offer/Update";
import AdminOfferDetail from "./Components/Admin/Offer/Detail";
import AdminOfferList from "./Components/Admin/Offer/List";
import AdminOfferAdd from "./Components/Admin/Offer/Add";

// Store
import AdminStoreList from "./Components/Admin/Store/List";
import AdminStoreAdd from "./Components/Admin/Store/Add";
import AdminStoreDetail from "./Components/Admin/Store/Detail";
import AdminStoreDetailForPurchaseList from "./Components/Admin/Store/UserPurchases/List";
import AdminStoreDetailForComicSubmission from "./Components/Admin/Store/DetailForComicSubmission";
import AdminStoreDetailForUserList from "./Components/Admin/Store/DetailForUserList";
import AdminStoreDetailForCommentList from "./Components/Admin/Store/DetailForCommentList";
import AdminStoreDetailForAttachmentList from "./Components/Admin/Store/DetailForAttachmentList";
import AdminStoreAttachmentAdd from "./Components/Admin/Store/Attachment/Add";
import AdminStoreAttachmentDetail from "./Components/Admin/Store/Attachment/Detail";
import AdminStoreAttachmentUpdate from "./Components/Admin/Store/Attachment/Update";
import AdminStoreUpdate from "./Components/Admin/Store/Update";

// Dashboard
import AdminDashboard from "./Components/Admin/Dashboard";

// Registry
import AdminRegistrySearch from "./Components/Admin/Registry/Search";
import AdminRegistryResult from "./Components/Admin/Registry/Result";

// Comic Submission
import AdminSubmissionPickTypeForAdd from "./Components/Admin/Submission/PickTypeForAdd";
import AdminComicSubmissionList from "./Components/Admin/Submission/Comic/List";
import AdminComicSubmissionAddStep1WithSearch from "./Components/Admin/Submission/Comic/AddStep1WithSearch";
import AdminComicSubmissionAddStep1WithResult from "./Components/Admin/Submission/Comic/AddStep1WithResult";
import AdminComicSubmissionAddStep1WithStarredCustomer from "./Components/Admin/Submission/Comic/AddStep1WithStarred";
import AdminComicSubmissionAddStep2 from "./Components/Admin/Submission/Comic/AddStep2";
import AdminComicSubmissionAddStep3 from "./Components/Admin/Submission/Comic/AddStep3";
import AdminComicSubmissionDetail from "./Components/Admin/Submission/Comic/Detail";
import AdminComicSubmissionDetailForCommentList from "./Components/Admin/Submission/Comic/DetailForCommentList";
import AdminComicSubmissionDetailForCustomer from "./Components/Admin/Submission/Comic/DetailForCustomer";
import AdminComicSubmissionDetailForPDFFile from "./Components/Admin/Submission/Comic/DetailForPDFFile";
import AdminSubmissionAttachmentAdd from "./Components/Admin/Submission/Comic/Attachment/Add";
import AdminSubmissionAttachmentDetail from "./Components/Admin/Submission/Comic/Attachment/Detail";
import AdminSubmissionAttachmentUpdate from "./Components/Admin/Submission/Comic/Attachment/Update";
import AdminComicSubmissionDetailForAttachmentList from "./Components/Admin/Submission/Comic/DetailForAttachmentList";
import AdminComicSubmissionUpdateForComicSubmission from "./Components/Admin/Submission/Comic/UpdateSubmission";
import AdminComicSubmissionUpdatePickCustomerWithResult from "./Components/Admin/Submission/Comic/UpdatePickCustomerWithResult";
import AdminComicSubmissionUpdatePickCustomerWithSearch from "./Components/Admin/Submission/Comic/UpdatePickCustomerWithSearch";
import AdminSubmissionLaunchpad from "./Components/Admin/Submission/Launchpad";

// Users
import AdminUserList from "./Components/Admin/User/List/View";
import AdminUserAdd from "./Components/Admin/User/Add/View";
import AdminUserDetail from "./Components/Admin/User/Detail/View";
import AdminUserDetailForComicSubmissionList from "./Components/Admin/User/Detail/ComicSubmission/ListView";
import AdminUserDetailForCommentList from "./Components/Admin/User/Detail/Comment/View";
import AdminUserDetailForAttachmentList from "./Components/Admin/User/Detail/Attachment/ListView";
import AdminUserAttachmentAdd from "./Components/Admin/User/Detail/Attachment/Add";
import AdminUserAttachmentDetail from "./Components/Admin/User/Detail/Attachment/Detail";
import AdminUserAttachmentUpdate from "./Components/Admin/User/Detail/Attachment/Update";
import AdminUserCreditList from "./Components/Admin/User/Detail/Credit/List";
import AdminUserCreditAdd from "./Components/Admin/User/Detail/Credit/Add";
import AdminUserCreditDetail from "./Components/Admin/User/Detail/Credit/Detail";
import AdminUserCreditUpdate from "./Components/Admin/User/Detail/Credit/Update";
import AdminUserDetailMore from "./Components/Admin/User/Detail/More/View";
import AdminUserArchiveOperation from "./Components/Admin/User/Detail/More/Archive/View";
import AdminUserUnarchiveOperation from "./Components/Admin/User/Detail/More/Unarchive/View";
import AdminUserDeleteOperation from "./Components/Admin/User/Detail/More/Delete/View";
import AdminUserMoreOperationChangePassword from "./Components/Admin/User/Detail/More/ChangePassword/View";
import AdminUserMoreOperation2FAToggle from "./Components/Admin/User/Detail/More/2FA/View";
import AdminUserUpdate from "./Components/Admin/User/Update/View";

//-----------------//
// Retailer Portal //
//-----------------//

// Dashboard
import RetailerDashboard from "./Components/Retailer/Dashboard";

// Registry
import RetailerRegistrySearch from "./Components/Retailer/Registry/Search";
import RetailerRegistryResult from "./Components/Retailer/Registry/Result";

// Comic Submission
import RetailerSubmissionPickTypeForAdd from "./Components/Retailer/Submission/PickForAdd";
import RetailerComicSubmissionList from "./Components/Retailer/Submission/Comic/List";
import RetailerComicSubmissionAddStep1WithSearch from "./Components/Retailer/Submission/Comic/AddStep01WithSearch";
import RetailerComicSubmissionAddStep1WithResult from "./Components/Retailer/Submission/Comic/AddStep01WithResult";
import RetailerComicSubmissionAddStep1WithStarredCustomer from "./Components/Retailer/Submission/Comic/AddStep01WithStarred";
import RetailerComicSubmissionAddStep2 from "./Components/Retailer/Submission/Comic/AddStep02";
import RetailerComicSubmissionAddStep3 from "./Components/Retailer/Submission/Comic/AddStep03";
import RetailerComicSubmissionAddStep4 from "./Components/Retailer/Submission/Comic/AddStep04";
import RetailerComicSubmissionAddStep5 from "./Components/Retailer/Submission/Comic/AddStep05";
import RetailerComicSubmissionAddStep6 from "./Components/Retailer/Submission/Comic/AddStep06";
import RetailerComicSubmissionAddStep7 from "./Components/Retailer/Submission/Comic/AddStep07";
import RetailerComicSubmissionAddStep8 from "./Components/Retailer/Submission/Comic/AddStep08";
import RetailerComicSubmissionAddStep9 from "./Components/Retailer/Submission/Comic/AddStep09";
import RetailerComicSubmissionAddStep10 from "./Components/Retailer/Submission/Comic/AddStep10";
import RetailerComicSubmissionAddStep11CheckoutLaunchpad from "./Components/Retailer/Submission/Comic/AddStep11CheckoutLaunchpad";
import RetailerComicSubmissionAddStep12Confirmation from "./Components/Retailer/Submission/Comic/AddStep12Confirmation";
// import RetailerComicSubmissionAddStepXXX from "./Components/Retailer/Submission/Comic/AddStepX"; // DEPRECATED: PLEASE REMOVE SOON...
// import RetailerComicSubmissionAddStepYYY from "./Components/Retailer/Submission/Comic/AddStepY"; // DEPRECATED: PLEASE REMOVE SOON...
import RetailerComicSubmissionDetail from "./Components/Retailer/Submission/Comic/Detail";
import RetailerSubmissionLaunchpad from "./Components/Retailer/Submission/Launchpad";
import RetailerSubmissionAttachmentAdd from "./Components/Retailer/Submission/Comic/Attachment/Add";
import RetailerSubmissionAttachmentDetail from "./Components/Retailer/Submission/Comic/Attachment/Detail";
import RetailerSubmissionAttachmentUpdate from "./Components/Retailer/Submission/Comic/Attachment/Update";
import RetailerComicSubmissionDetailForAttachmentList from "./Components/Retailer/Submission/Comic/DetailForAttachmentList";
import RetailerComicSubmissionDetailForCommentList from "./Components/Retailer/Submission/Comic/DetailForCommentList";
import RetailerComicSubmissionDetailForCustomer from "./Components/Retailer/Submission/Comic/DetailForCustomer";
import RetailerComicSubmissionDetailForPDFFile from "./Components/Retailer/Submission/Comic/DetailForPDFFile";
import RetailerComicSubmissionUpdateForComicSubmission from "./Components/Retailer/Submission/Comic/UpdateSubmission";
import RetailerComicSubmissionUpdatePickCustomerWithResult from "./Components/Retailer/Submission/Comic/UpdatePickCustomerWithResult";
import RetailerComicSubmissionUpdatePickCustomerWithSearch from "./Components/Retailer/Submission/Comic/UpdatePickCustomerWithSearch";

// Customer
import RetailerCustomerList from "./Components/Retailer/Customer/List";
import RetailerCustomerAdd from "./Components/Retailer/Customer/Add";
import RetailerCustomerDetail from "./Components/Retailer/Customer/Detail";
import RetailerCustomerDetailForComicSubmissionList from "./Components/Retailer/Customer/DetailForComicSubmissionList";
import RetailerCustomerDetailForCommentList from "./Components/Retailer/Customer/DetailForCommentList";
import RetailerCustomerUpdate from "./Components/Retailer/Customer/Update";
import RetailerCustomerDetailForAttachmentList from "./Components/Retailer/Customer/DetailForAttachmentList";
import RetailerCustomerAttachmentAdd from "./Components/Retailer/Customer/Attachment/Add";
import RetailerCustomerAttachmentDetail from "./Components/Retailer/Customer/Attachment/Detail";
import RetailerCustomerAttachmentUpdate from "./Components/Retailer/Customer/Attachment/Update";

//---------------//
// Common System //
//---------------//

// Account
import AccountDetail from "./Components/Profile/Detail/View";
import AccountUpdate from "./Components/Profile/Update/View";
import AccountMoreLaunchpad from "./Components/Profile/More/LaunchpadView";
import AccountTwoFactorAuthenticationDetail from "./Components/Profile/2FA/View";
import AccountEnableTwoFactorAuthenticationStep1 from "./Components/Profile/2FA/EnableStep1View";
import AccountEnableTwoFactorAuthenticationStep2 from "./Components/Profile/2FA/EnableStep2View";
import AccountEnableTwoFactorAuthenticationStep3 from "./Components/Profile/2FA/EnableStep3View";
import AccountTwoFactorAuthenticationBackupCodeGenerate from "./Components/Profile/2FA/BackupCodeGenerateView";
import AccountMoreOperationChangePassword from "./Components/Profile/More/Operation/ChangePassword/View";

// Gateway
import LogoutRedirector from "./Components/Gateway/LogoutRedirector";
import Login from "./Components/Gateway/Login";
import RegisterLaunchpad from "./Components/Gateway/Register/Launchpad";
import RegisterAsStoreOwner from "./Components/Gateway/Register/StoreOwner";
import RegisterAsCustomer from "./Components/Gateway/Register/Customer";
import RegisterSuccessful from "./Components/Gateway/RegisterSuccessful";
import EmailVerification from "./Components/Gateway/EmailVerification";
import ForgotPassword from "./Components/Gateway/ForgotPassword";
import PasswordReset from "./Components/Gateway/PasswordReset";

// Gateway (2FA specific)
import TwoFactorAuthenticationWizardStep1 from "./Components/Gateway/2FA/Step1";
import TwoFactorAuthenticationWizardStep2 from "./Components/Gateway/2FA/Step2";
import TwoFactorAuthenticationWizardStep3 from "./Components/Gateway/2FA/Step3";
import TwoFactorAuthenticationBackupCodeGenerate from "./Components/Gateway/2FA/BackupCodeGenerate";
import TwoFactorAuthenticationBackupCodeRecovery from "./Components/Gateway/2FA/BackupCodeRecovery";
import TwoFactorAuthenticationValidateOnLogin from "./Components/Gateway/2FA/ValidateOnLogin";

// Store
import StoreDetail from "./Components/Store/Detail";
import StoreUpdate from "./Components/Store/Update";
import StorePurchaseList from "./Components/Store/PurchaseList";
import StoreCreditList from "./Components/Store/CreditList";

// Navigation
import TopAlertBanner from "./Components/Misc/TopAlertBanner";
import Sidebar from "./Components/Menu/Sidebar";
import Topbar from "./Components/Menu/Top";

// Redirectors.
import AnonymousCurrentUserRedirector from "./Components/Misc/AnonymousCurrentUserRedirector";
import TwoFactorAuthenticationRedirector from "./Components/Misc/TwoFactorAuthenticationRedirector";
import CPSRNRedirector from "./Components/Gateway/CPSRNRedirector";

// Public Registery
import PublicRegistrySearch from "./Components/Gateway/RegistrySearch";
import PublicRegistryResult from "./Components/Gateway/RegistryResult";

// Public Generic
import Index from "./Components/Gateway/Index";
import Terms from "./Components/Gateway/Terms";
import Privacy from "./Components/Gateway/Privacy";
import NotImplementedError from "./Components/Misc/NotImplementedError";
import NotFoundError from "./Components/Misc/NotFoundError";
import DashboardHelp from "./Components/Misc/DashboardHelp";

//-----------------//
// Customer Portal //
//-----------------//

// Dashboard
import CustomerDashboard from "./Components/Customer/Dashboard";

// Registry
import CustomerRegistrySearch from "./Components/Customer/Registry/Search";
import CustomerRegistryResult from "./Components/Customer/Registry/Result";

// Comic Submission
import CustomerSubmissionPickTypeForAdd from "./Components/Customer/Submission/PickForAdd";
import CustomerComicSubmissionList from "./Components/Customer/Submission/Comic/List";
import CustomerComicSubmissionAddStep1 from "./Components/Customer/Submission/Comic/AddStep1";
import CustomerComicSubmissionAddStep2 from "./Components/Customer/Submission/Comic/AddStep2";
import CustomerComicSubmissionAddStep3 from "./Components/Customer/Submission/Comic/AddStep3";
import CustomerComicSubmissionDetail from "./Components/Customer/Submission/Comic/Detail";
import CustomerSubmissionLaunchpad from "./Components/Customer/Submission/Launchpad";
import CustomerSubmissionAttachmentAdd from "./Components/Customer/Submission/Comic/Attachment/Add";
import CustomerSubmissionAttachmentDetail from "./Components/Customer/Submission/Comic/Attachment/Detail";
import CustomerSubmissionAttachmentUpdate from "./Components/Customer/Submission/Comic/Attachment/Update";
import CustomerComicSubmissionDetailForAttachmentList from "./Components/Customer/Submission/Comic/DetailForAttachmentList";
import CustomerComicSubmissionDetailForCommentList from "./Components/Customer/Submission/Comic/DetailForCommentList";
import CustomerComicSubmissionDetailForPDFFile from "./Components/Customer/Submission/Comic/DetailForPDFFile";
import CustomerComicSubmissionUpdateForComicSubmission from "./Components/Customer/Submission/Comic/UpdateSubmission";

//-----------------//
//    App Routes   //
//-----------------//

function AppRoute() {
  return (
    <div class="is-widescreen is-size-5-widescreen is-size-6-tablet is-size-7-mobile">
      {/*
            NOTES FOR ABOVE
            USE THE FOLLOWING TEXT SIZES BASED ON DEVICE TYPE
            - is-size-5-widescreen
            - is-size-6-tablet
            - is-size-7-mobile
        */}
      <RecoilRoot>
        <Router>
          <AnonymousCurrentUserRedirector />
          <TwoFactorAuthenticationRedirector />
          <TopAlertBanner />
          <Topbar />
          <div class="columns">
            <Sidebar />
            <div class="column">
              <section class="main-content columns is-fullheight">
                <Routes>
                  <Route
                    exact
                    path="/admin/offer/:id/update"
                    element={<AdminOfferUpdate />}
                  />
                  <Route
                    exact
                    path="/admin/offer/:id"
                    element={<AdminOfferDetail />}
                  />
                  <Route
                    exact
                    path="/admin/offers/add"
                    element={<AdminOfferAdd />}
                  />
                  <Route
                    exact
                    path="/admin/offers"
                    element={<AdminOfferList />}
                  />
                  <Route
                    exact
                    path="/admin/stores"
                    element={<AdminStoreList />}
                  />
                  <Route
                    exact
                    path="/admin/stores/add"
                    element={<AdminStoreAdd />}
                  />
                  <Route
                    exact
                    path="/admin/store/:id"
                    element={<AdminStoreDetail />}
                  />
                  <Route
                    exact
                    path="/admin/store/:id/purchases"
                    element={<AdminStoreDetailForPurchaseList />}
                  />
                  <Route
                    exact
                    path="/admin/store/:id/comics"
                    element={<AdminStoreDetailForComicSubmission />}
                  />
                  <Route
                    exact
                    path="/admin/store/:id/users"
                    element={<AdminStoreDetailForUserList />}
                  />
                  <Route
                    exact
                    path="/admin/store/:id/edit"
                    element={<AdminStoreUpdate />}
                  />
                  <Route
                    exact
                    path="/admin/store/:id/comments"
                    element={<AdminStoreDetailForCommentList />}
                  />
                  <Route
                    exact
                    path="/admin/store/:id/attachments"
                    element={<AdminStoreDetailForAttachmentList />}
                  />
                  <Route
                    exact
                    path="/admin/store/:id/attachment/:aid"
                    element={<AdminStoreAttachmentDetail />}
                  />
                  <Route
                    exact
                    path="/admin/store/:id/attachments/add"
                    element={<AdminStoreAttachmentAdd />}
                  />
                  <Route
                    exact
                    path="/admin/store/:id/attachment/:aid/edit"
                    element={<AdminStoreAttachmentUpdate />}
                  />
                  <Route
                    exact
                    path="/admin/registry"
                    element={<AdminRegistrySearch />}
                  />
                  <Route
                    exact
                    path="/admin/registry/:cpsn"
                    element={<AdminRegistryResult />}
                  />
                  <Route
                    exact
                    path="/admin/submissions/pick-type-for-add"
                    element={<AdminSubmissionPickTypeForAdd />}
                  />
                  <Route
                    exact
                    path="/admin/submissions"
                    element={<AdminSubmissionLaunchpad />}
                  />
                  <Route
                    exact
                    path="/admin/submissions/comics"
                    element={<AdminComicSubmissionList />}
                  />
                  <Route
                    exact
                    path="/admin/submissions/comics/add/search"
                    element={<AdminComicSubmissionAddStep1WithSearch />}
                  />
                  <Route
                    exact
                    path="/admin/submissions/comics/add/results"
                    element={<AdminComicSubmissionAddStep1WithResult />}
                  />
                  <Route
                    exact
                    path="/admin/submissions/comics/add/starred"
                    element={
                      <AdminComicSubmissionAddStep1WithStarredCustomer />
                    }
                  />
                  <Route
                    exact
                    path="/admin/submissions/comics/add"
                    element={<AdminComicSubmissionAddStep2 />}
                  />
                  <Route
                    exact
                    path="/admin/submissions/comics/add/:id/confirmation"
                    element={<AdminComicSubmissionAddStep3 />}
                  />
                  <Route
                    exact
                    path="/admin/submissions/comic/:id"
                    element={<AdminComicSubmissionDetail />}
                  />
                  <Route
                    exact
                    path="/admin/submissions/comic/:id/edit"
                    element={<AdminComicSubmissionUpdateForComicSubmission />}
                  />
                  <Route
                    exact
                    path="/admin/submissions/comic/:id/cust/search"
                    element={
                      <AdminComicSubmissionUpdatePickCustomerWithSearch />
                    }
                  />
                  <Route
                    exact
                    path="/admin/submissions/comic/:id/cust/results"
                    element={
                      <AdminComicSubmissionUpdatePickCustomerWithResult />
                    }
                  />
                  <Route
                    exact
                    path="/admin/submissions/comic/:id/comments"
                    element={<AdminComicSubmissionDetailForCommentList />}
                  />
                  <Route
                    exact
                    path="/admin/submissions/comic/:id/cust"
                    element={<AdminComicSubmissionDetailForCustomer />}
                  />
                  <Route
                    exact
                    path="/admin/submissions/comic/:id/file"
                    element={<AdminComicSubmissionDetailForPDFFile />}
                  />
                  <Route
                    exact
                    path="/admin/submissions/comic/:id/attachments"
                    element={<AdminComicSubmissionDetailForAttachmentList />}
                  />
                  <Route
                    exact
                    path="/admin/submissions/comic/:id/attachment/:aid/edit"
                    element={<AdminSubmissionAttachmentUpdate />}
                  />
                  <Route
                    exact
                    path="/admin/submissions/comic/:id/attachment/:aid"
                    element={<AdminSubmissionAttachmentDetail />}
                  />
                  <Route
                    exact
                    path="/admin/submissions/comic/:id/attachments/add"
                    element={<AdminSubmissionAttachmentAdd />}
                  />
                  <Route
                    exact
                    path="/admin/submissions/cards"
                    element={<NotImplementedError />}
                  />
                  <Route
                    exact
                    path="/admin/users"
                    element={<AdminUserList />}
                  />
                  <Route
                    exact
                    path="/admin/users/add"
                    element={<AdminUserAdd />}
                  />
                  <Route
                    exact
                    path="/admin/user/:id"
                    element={<AdminUserDetail />}
                  />
                  <Route
                    exact
                    path="/admin/user/:id/comics"
                    element={<AdminUserDetailForComicSubmissionList />}
                  />
                  <Route
                    exact
                    path="/admin/user/:id/edit"
                    element={<AdminUserUpdate />}
                  />
                  <Route
                    exact
                    path="/admin/user/:id/comments"
                    element={<AdminUserDetailForCommentList />}
                  />
                  <Route
                    exact
                    path="/admin/user/:id/attachments"
                    element={<AdminUserDetailForAttachmentList />}
                  />
                  <Route
                    exact
                    path="/admin/user/:id/attachment/:aid"
                    element={<AdminUserAttachmentDetail />}
                  />
                  <Route
                    exact
                    path="/admin/user/:id/attachments/add"
                    element={<AdminUserAttachmentAdd />}
                  />
                  <Route
                    exact
                    path="/admin/user/:id/attachment/:aid/edit"
                    element={<AdminUserAttachmentUpdate />}
                  />
                  <Route
                    exact
                    path="/admin/user/:id/credits"
                    element={<AdminUserCreditList />}
                  />
                  <Route
                    exact
                    path="/admin/user/:id/credits/add"
                    element={<AdminUserCreditAdd />}
                  />
                  <Route
                    exact
                    path="/admin/user/:id/credit/:cid"
                    element={<AdminUserCreditDetail />}
                  />
                  <Route
                    exact
                    path="/admin/user/:id/credit/:cid/edit"
                    element={<AdminUserCreditUpdate />}
                  />
                  <Route
                    exact
                    path="/admin/user/:id/more"
                    element={<AdminUserDetailMore />}
                  />
                  <Route
                    exact
                    path="/admin/user/:id/more/archive"
                    element={<AdminUserArchiveOperation />}
                  />
                  <Route
                    exact
                    path="/admin/user/:id/more/unarchive"
                    element={<AdminUserUnarchiveOperation />}
                  />
                  <Route
                    exact
                    path="/admin/user/:id/more/permadelete"
                    element={<AdminUserDeleteOperation />}
                  />
                  <Route
                    exact
                    path="/admin/user/:id/more/change-password"
                    element={<AdminUserMoreOperationChangePassword />}
                  />
                  <Route
                    exact
                    path="/admin/user/:id/more/change-2fa"
                    element={<AdminUserMoreOperation2FAToggle />}
                  />
                  <Route
                    exact
                    path="/admin/dashboard"
                    element={<AdminDashboard />}
                  />
                  <Route
                    exact
                    path="/dashboard"
                    element={<RetailerDashboard />}
                  />
                  <Route
                    exact
                    path="/registry"
                    element={<RetailerRegistrySearch />}
                  />
                  <Route
                    exact
                    path="/registry/:cpsn"
                    element={<RetailerRegistryResult />}
                  />
                  <Route
                    exact
                    path="/submissions"
                    element={<RetailerSubmissionLaunchpad />}
                  />
                  <Route
                    exact
                    path="/submissions/pick-type-for-add"
                    element={<RetailerSubmissionPickTypeForAdd />}
                  />
                  <Route
                    exact
                    path="/submissions/comics"
                    element={<RetailerComicSubmissionList />}
                  />
                  <Route
                    exact
                    path="/submissions/comics/add/step-1/search"
                    element={<RetailerComicSubmissionAddStep1WithSearch />}
                  />
                  <Route
                    exact
                    path="/submissions/comics/add/step-1/results"
                    element={<RetailerComicSubmissionAddStep1WithResult />}
                  />
                  <Route
                    exact
                    path="/submissions/comics/add/step-1/starred"
                    element={
                      <RetailerComicSubmissionAddStep1WithStarredCustomer />
                    }
                  />
                  <Route
                    exact
                    path="/submissions/comics/add/step-2"
                    element={<RetailerComicSubmissionAddStep2 />}
                  />
                  <Route
                    exact
                    path="/submissions/comics/add/step-3"
                    element={<RetailerComicSubmissionAddStep3 />}
                  />
                  <Route
                    exact
                    path="/submissions/comics/add/step-4"
                    element={<RetailerComicSubmissionAddStep4 />}
                  />
                  <Route
                    exact
                    path="/submissions/comics/add/step-5"
                    element={<RetailerComicSubmissionAddStep5 />}
                  />
                  <Route
                    exact
                    path="/submissions/comics/add/step-6"
                    element={<RetailerComicSubmissionAddStep6 />}
                  />
                  <Route
                    exact
                    path="/submissions/comics/add/step-7"
                    element={<RetailerComicSubmissionAddStep7 />}
                  />
                  <Route
                    exact
                    path="/submissions/comics/add/step-8"
                    element={<RetailerComicSubmissionAddStep8 />}
                  />
                  <Route
                    exact
                    path="/submissions/comics/add/step-9"
                    element={<RetailerComicSubmissionAddStep9 />}
                  />
                  <Route
                    exact
                    path="/submissions/comics/add/step-10"
                    element={<RetailerComicSubmissionAddStep10 />}
                  />
                  <Route
                    exact
                    path="/submissions/comics/add/checkout"
                    element={<RetailerComicSubmissionAddStep11CheckoutLaunchpad />}
                  />
                  <Route
                    exact
                    path="/submissions/comics/add/confirmation"
                    element={<RetailerComicSubmissionAddStep12Confirmation />}
                  />
                  {/*
                  ----- DEPRECATED: PLEASE REMOVE SOON... ----
                  <Route
                    exact
                    path="/submissions/comics/add/:id"
                    element={<RetailerComicSubmissionAddStepXXX />}
                  />
                  <Route
                    exact
                    path="/submissions/comics/add/:id/confirmation"
                    element={<RetailerComicSubmissionAddStepYYY />}
                  />
                  */}
                  <Route
                    exact
                    path="/submissions/comic/:id"
                    element={<RetailerComicSubmissionDetail />}
                  />
                  <Route
                    exact
                    path="/submissions/comic/:id/edit"
                    element={
                      <RetailerComicSubmissionUpdateForComicSubmission />
                    }
                  />
                  <Route
                    exact
                    path="/submissions/comic/:id/cust/search"
                    element={
                      <RetailerComicSubmissionUpdatePickCustomerWithSearch />
                    }
                  />
                  <Route
                    exact
                    path="/submissions/comic/:id/cust/results"
                    element={
                      <RetailerComicSubmissionUpdatePickCustomerWithResult />
                    }
                  />
                  <Route
                    exact
                    path="/submissions/comic/:id/attachment/:aid/edit"
                    element={<RetailerSubmissionAttachmentUpdate />}
                  />
                  <Route
                    exact
                    path="/submissions/comic/:id/attachment/:aid"
                    element={<RetailerSubmissionAttachmentDetail />}
                  />
                  <Route
                    exact
                    path="/submissions/comic/:id/attachments/add"
                    element={<RetailerSubmissionAttachmentAdd />}
                  />
                  <Route
                    exact
                    path="/submissions/comic/:id/attachments"
                    element={<RetailerComicSubmissionDetailForAttachmentList />}
                  />
                  <Route
                    exact
                    path="/submissions/comic/:id/comments"
                    element={<RetailerComicSubmissionDetailForCommentList />}
                  />
                  <Route
                    exact
                    path="/submissions/comic/:id/cust"
                    element={<RetailerComicSubmissionDetailForCustomer />}
                  />
                  <Route
                    exact
                    path="/submissions/comic/:id/file"
                    element={<RetailerComicSubmissionDetailForPDFFile />}
                  />
                  <Route
                    exact
                    path="/submissions/cards/add/search"
                    element={<NotImplementedError />}
                  />
                  <Route
                    exact
                    path="/submissions/cards"
                    element={<NotImplementedError />}
                  />
                  <Route
                    exact
                    path="/customers"
                    element={<RetailerCustomerList />}
                  />
                  <Route
                    exact
                    path="/customers/add"
                    element={<RetailerCustomerAdd />}
                  />
                  <Route
                    exact
                    path="/customer/:id"
                    element={<RetailerCustomerDetail />}
                  />
                  <Route
                    exact
                    path="/customer/:id/comics"
                    element={<RetailerCustomerDetailForComicSubmissionList />}
                  />
                  <Route
                    exact
                    path="/customer/:id/edit"
                    element={<RetailerCustomerUpdate />}
                  />
                  <Route
                    exact
                    path="/customer/:id/comments"
                    element={<RetailerCustomerDetailForCommentList />}
                  />
                  <Route
                    exact
                    path="/customer/:id/attachments"
                    element={<RetailerCustomerDetailForAttachmentList />}
                  />
                  <Route
                    exact
                    path="/customer/:id/attachments/add"
                    element={<RetailerCustomerAttachmentAdd />}
                  />
                  <Route
                    exact
                    path="/customer/:id/attachment/:aid"
                    element={<RetailerCustomerAttachmentDetail />}
                  />
                  <Route
                    exact
                    path="/customer/:id/attachment/:aid/edit"
                    element={<RetailerCustomerAttachmentUpdate />}
                  />
                  <Route exact path="/account" element={<AccountDetail />} />
                  <Route
                    exact
                    path="/account/update"
                    element={<AccountUpdate />}
                  />
                  <Route
                    exact
                    path="/account/more"
                    element={<AccountMoreLaunchpad />}
                  />
                  <Route
                    exact
                    path="/account/more/change-password"
                    element={<AccountMoreOperationChangePassword />}
                  />
                  <Route
                    exact
                    path="/account/2fa"
                    element={<AccountTwoFactorAuthenticationDetail />}
                  />
                  <Route
                    exact
                    path="/account/2fa/setup/step-1"
                    element={<AccountEnableTwoFactorAuthenticationStep1 />}
                  />
                  <Route
                    exact
                    path="/account/2fa/setup/step-2"
                    element={<AccountEnableTwoFactorAuthenticationStep2 />}
                  />
                  <Route
                    exact
                    path="/account/2fa/setup/step-3"
                    element={<AccountEnableTwoFactorAuthenticationStep3 />}
                  />
                  <Route
                    exact
                    path="/account/2fa/backup-code"
                    element={
                      <AccountTwoFactorAuthenticationBackupCodeGenerate />
                    }
                  />
                  <Route exact path="/store/" element={<StoreDetail />} />
                  <Route exact path="/store/update" element={<StoreUpdate />} />
                  <Route
                    exact
                    path="/store/:id/purchases"
                    element={<StorePurchaseList />}
                  />
                  <Route
                    exact
                    path="/store/:id/credits"
                    element={<StoreCreditList />}
                  />
                  <Route
                    exact
                    path="/register"
                    element={<RegisterLaunchpad />}
                  />
                  {/*
                  <Route
                    exact
                    path="/register"
                    element={<Register />}
                  />
                  <Route
                    exact
                    path="/register/v2"
                    element={<RegisterLaunchpad />}
                  />
                  */}
                  <Route
                    exact
                    path="/register/store"
                    element={<RegisterAsStoreOwner />}
                  />
                  <Route
                    exact
                    path="/register/user"
                    element={<RegisterAsCustomer />}
                  />
                  <Route
                    exact
                    path="/register-successful"
                    element={<RegisterSuccessful />}
                  />
                  <Route exact path="/login" element={<Login />} />
                  <Route
                    exact
                    path="/login/2fa/step-1"
                    element={<TwoFactorAuthenticationWizardStep1 />}
                  />
                  <Route
                    exact
                    path="/login/2fa/step-2"
                    element={<TwoFactorAuthenticationWizardStep2 />}
                  />
                  <Route
                    exact
                    path="/login/2fa/step-3"
                    element={<TwoFactorAuthenticationWizardStep3 />}
                  />
                  <Route
                    exact
                    path="/login/2fa/backup-code"
                    element={<TwoFactorAuthenticationBackupCodeGenerate />}
                  />
                  <Route
                    exact
                    path="/login/2fa/backup-code-recovery"
                    element={<TwoFactorAuthenticationBackupCodeRecovery />}
                  />
                  <Route
                    exact
                    path="/login/2fa"
                    element={<TwoFactorAuthenticationValidateOnLogin />}
                  />
                  <Route exact path="/logout" element={<LogoutRedirector />} />
                  <Route exact path="/verify" element={<EmailVerification />} />
                  <Route exact path="/terms" element={<Terms />} />
                  <Route exact path="/privacy" element={<Privacy />} />
                  <Route
                    exact
                    path="/forgot-password"
                    element={<ForgotPassword />}
                  />
                  <Route
                    exact
                    path="/password-reset"
                    element={<PasswordReset />}
                  />
                  <Route
                    exact
                    path="/cpsrn-result"
                    element={<PublicRegistryResult />}
                  />
                  <Route
                    exact
                    path="/cpsrn-registry"
                    element={<PublicRegistrySearch />}
                  />
                  <Route
                    exact
                    path="/c/dashboard"
                    element={<CustomerDashboard />}
                  />
                  <Route
                    exact
                    path="/c/registry"
                    element={<CustomerRegistrySearch />}
                  />
                  <Route
                    exact
                    path="/c/registry/:cpsn"
                    element={<CustomerRegistryResult />}
                  />
                  <Route
                    exact
                    path="/c/submissions"
                    element={<CustomerSubmissionLaunchpad />}
                  />
                  <Route
                    exact
                    path="/c/submissions/pick-type-for-add"
                    element={<CustomerSubmissionPickTypeForAdd />}
                  />
                  <Route
                    exact
                    path="/c/submissions/comics"
                    element={<CustomerComicSubmissionList />}
                  />
                  <Route
                    exact
                    path="/c/submissions/comics/add"
                    element={<CustomerComicSubmissionAddStep1 />}
                  />
                  <Route
                    exact
                    path="/c/submissions/comics/add/:id"
                    element={<CustomerComicSubmissionAddStep2 />}
                  />
                  <Route
                    exact
                    path="/c/submissions/comics/add/:id/confirmation"
                    element={<CustomerComicSubmissionAddStep3 />}
                  />
                  <Route
                    exact
                    path="/c/submissions/comic/:id"
                    element={<CustomerComicSubmissionDetail />}
                  />
                  <Route
                    exact
                    path="/c/submissions/comic/:id/edit"
                    element={
                      <CustomerComicSubmissionUpdateForComicSubmission />
                    }
                  />
                  <Route
                    exact
                    path="/c/submissions/comic/:id/attachment/:aid/edit"
                    element={<CustomerSubmissionAttachmentUpdate />}
                  />
                  <Route
                    exact
                    path="/c/submissions/comic/:id/attachment/:aid"
                    element={<CustomerSubmissionAttachmentDetail />}
                  />
                  <Route
                    exact
                    path="/c/submissions/comic/:id/attachments/add"
                    element={<CustomerSubmissionAttachmentAdd />}
                  />
                  <Route
                    exact
                    path="/c/submissions/comic/:id/attachments"
                    element={<CustomerComicSubmissionDetailForAttachmentList />}
                  />
                  <Route
                    exact
                    path="/c/submissions/comic/:id/comments"
                    element={<CustomerComicSubmissionDetailForCommentList />}
                  />
                  <Route
                    exact
                    path="/c/submissions/comic/:id/file"
                    element={<CustomerComicSubmissionDetailForPDFFile />}
                  />
                  <Route
                    exact
                    path="/help"
                    element={<DashboardHelp />}
                  />
                  <Route
                    exact
                    path="/cpsrn"
                    element={<CPSRNRedirector />}
                  />
                  <Route exact path="/" element={<Index />} />
                  <Route path="*" element={<NotFoundError />} />
                </Routes>
              </section>
              <div>
                {/* DEVELOPERS NOTE: Mobile tab-bar menu can go here */}
              </div>
              <footer class="footer is-hidden">
                <div class="container">
                  <div class="content has-text-centered">
                    <p>Hello</p>
                  </div>
                </div>
              </footer>
            </div>
          </div>
        </Router>
      </RecoilRoot>
    </div>
  );
}

export default AppRoute;
