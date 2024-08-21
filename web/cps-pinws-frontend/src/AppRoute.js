import { React, useState } from "react";
import "bulma/css/bulma.min.css";
import "./index.css";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { RecoilRoot } from "recoil";

//--------------//
// Admin Portal //
//--------------//

// Store
import AdminStoreList from "./Components/Admin/Store/List";
import AdminStoreAdd from "./Components/Admin/Store/Add";
import AdminStoreDetail from "./Components/Admin/Store/Detail";
import AdminStoreDetailForUserList from "./Components/Admin/Store/DetailForUserList";
import AdminStoreDetailForCommentList from "./Components/Admin/Store/DetailForCommentList";
import AdminStoreDetailForAttachmentList from "./Components/Admin/Store/DetailForAttachmentList";
import AdminStoreAttachmentAdd from "./Components/Admin/Store/Attachment/Add";
import AdminStoreAttachmentDetail from "./Components/Admin/Store/Attachment/Detail";
import AdminStoreAttachmentUpdate from "./Components/Admin/Store/Attachment/Update";
import AdminStoreUpdate from "./Components/Admin/Store/Update";

// Dashboard
import AdminDashboard from "./Components/Admin/Dashboard";

// Users
import AdminUserList from "./Components/Admin/User/List/View";
import AdminUserAdd from "./Components/Admin/User/Add/View";
import AdminUserDetail from "./Components/Admin/User/Detail/View";
import AdminUserDetailForCommentList from "./Components/Admin/User/Detail/Comment/View";
import AdminUserDetailForAttachmentList from "./Components/Admin/User/Detail/Attachment/ListView";
import AdminUserAttachmentAdd from "./Components/Admin/User/Detail/Attachment/Add";
import AdminUserAttachmentDetail from "./Components/Admin/User/Detail/Attachment/Detail";
import AdminUserAttachmentUpdate from "./Components/Admin/User/Detail/Attachment/Update";
import AdminUserDetailMore from "./Components/Admin/User/Detail/More/View";
import AdminUserArchiveOperation from "./Components/Admin/User/Detail/More/Archive/View";
import AdminUserUnarchiveOperation from "./Components/Admin/User/Detail/More/Unarchive/View";
import AdminUserDeleteOperation from "./Components/Admin/User/Detail/More/Delete/View";
import AdminUserMoreOperationChangePassword from "./Components/Admin/User/Detail/More/ChangePassword/View";
import AdminUserMoreOperation2FAToggle from "./Components/Admin/User/Detail/More/2FA/View";
import AdminUserUpdate from "./Components/Admin/User/Update/View";


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
import ForgotPassword from "./Components/Gateway/ForgotPassword";
import PasswordReset from "./Components/Gateway/PasswordReset";

// Gateway (2FA specific)
import TwoFactorAuthenticationWizardStep1 from "./Components/Gateway/2FA/Step1";
import TwoFactorAuthenticationWizardStep2 from "./Components/Gateway/2FA/Step2";
import TwoFactorAuthenticationWizardStep3 from "./Components/Gateway/2FA/Step3";
import TwoFactorAuthenticationBackupCodeGenerate from "./Components/Gateway/2FA/BackupCodeGenerate";
import TwoFactorAuthenticationBackupCodeRecovery from "./Components/Gateway/2FA/BackupCodeRecovery";
import TwoFactorAuthenticationValidateOnLogin from "./Components/Gateway/2FA/ValidateOnLogin";

// Navigation
import TopAlertBanner from "./Components/Misc/TopAlertBanner";
import Sidebar from "./Components/Menu/Sidebar";
import Topbar from "./Components/Menu/Top";

// Redirectors.
import AnonymousCurrentUserRedirector from "./Components/Misc/AnonymousCurrentUserRedirector";
import TwoFactorAuthenticationRedirector from "./Components/Misc/TwoFactorAuthenticationRedirector";

// Public Generic
import Index from "./Components/Gateway/Index";
import Terms from "./Components/Gateway/Terms";
import Privacy from "./Components/Gateway/Privacy";
import NotImplementedError from "./Components/Misc/NotImplementedError";
import NotFoundError from "./Components/Misc/NotFoundError";
import DashboardHelp from "./Components/Misc/DashboardHelp";

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
                    path="/help"
                    element={<DashboardHelp />}
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
