import { React, useState } from "react";
import "bulma/css/bulma.min.css";
import "./index.css";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { RecoilRoot } from "recoil";

//--------------//
// Admin Portal //
//--------------//

// Dashboard
import AdminDashboard from "./Components/Admin/Dashboard";

// Tenant
import AdminTenantList from "./Components/Admin/Tenant/List";
import AdminTenantAdd from "./Components/Admin/Tenant/Add";
import AdminTenantDetail from "./Components/Admin/Tenant/Detail";
import AdminTenantDetailForUserList from "./Components/Admin/Tenant/DetailForUserList";
import AdminTenantDetailForCommentList from "./Components/Admin/Tenant/DetailForCommentList";
import AdminTenantDetailForPinObjectList from "./Components/Admin/Tenant/DetailForPinObjectList";
import AdminTenantPinObjectAdd from "./Components/Admin/Tenant/PinObject/Add";
import AdminTenantPinObjectDetail from "./Components/Admin/Tenant/PinObject/Detail";
import AdminTenantPinObjectUpdate from "./Components/Admin/Tenant/PinObject/Update";
import AdminTenantUpdate from "./Components/Admin/Tenant/Update";

// Projects
import AdminProjectList from "./Components/Admin/Project/List/View";
import AdminProjectAdd from "./Components/Admin/Project/Add/View";
import AdminProjectDetail from "./Components/Admin/Project/Detail/View";
import AdminProjectUpdate from "./Components/Admin/Project/Update/View";
import AdminProjectDetailForPinObjectList from "./Components/Admin/Project/Detail/PinObject/ListView";
import AdminProjectPinObjectAdd from "./Components/Admin/Project/Detail/PinObject/Add";
// import AdminTenantPinObjectDetail from "./Components/Admin/Tenant/PinObject/Detail";
// import AdminTenantPinObjectUpdate from "./Components/Admin/Tenant/PinObject/Update";

// Users
import AdminUserList from "./Components/Admin/User/List/View";
import AdminUserAdd from "./Components/Admin/User/Add/View";
import AdminUserDetail from "./Components/Admin/User/Detail/View";
import AdminUserDetailForCommentList from "./Components/Admin/User/Detail/Comment/View";
import AdminUserDetailForPinObjectList from "./Components/Admin/User/Detail/PinObject/ListView";
import AdminUserPinObjectAdd from "./Components/Admin/User/Detail/PinObject/Add";
import AdminUserPinObjectDetail from "./Components/Admin/User/Detail/PinObject/Detail";
import AdminUserPinObjectUpdate from "./Components/Admin/User/Detail/PinObject/Update";
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
                    path="/admin/tenants"
                    element={<AdminTenantList />}
                  />
                  <Route
                    exact
                    path="/admin/tenants/add"
                    element={<AdminTenantAdd />}
                  />
                  <Route
                    exact
                    path="/admin/tenant/:id"
                    element={<AdminTenantDetail />}
                  />
                  <Route
                    exact
                    path="/admin/tenant/:id/users"
                    element={<AdminTenantDetailForUserList />}
                  />
                  <Route
                    exact
                    path="/admin/tenant/:id/edit"
                    element={<AdminTenantUpdate />}
                  />
                  <Route
                    exact
                    path="/admin/tenant/:id/comments"
                    element={<AdminTenantDetailForCommentList />}
                  />
                  <Route
                    exact
                    path="/admin/tenant/:id/pins"
                    element={<AdminTenantDetailForPinObjectList />}
                  />
                  <Route
                    exact
                    path="/admin/tenant/:id/pin/:aid"
                    element={<AdminTenantPinObjectDetail />}
                  />
                  <Route
                    exact
                    path="/admin/tenant/:id/pin/add"
                    element={<AdminTenantPinObjectAdd />}
                  />
                  <Route
                    exact
                    path="/admin/tenant/:id/pin/:aid/edit"
                    element={<AdminTenantPinObjectUpdate />}
                  />
                  {/* Project */}

                  <Route
                    exact
                    path="/admin/projects"
                    element={<AdminProjectList />}
                  />
                  <Route
                    exact
                    path="/admin/projects/add"
                    element={<AdminProjectAdd />}
                  />
                  <Route
                    exact
                    path="/admin/project/:id"
                    element={<AdminProjectDetail />}
                  />
                  <Route
                    exact
                    path="/admin/project/:id/edit"
                    element={<AdminProjectUpdate />}
                  />
                  <Route
                    exact
                    path="/admin/project/:id/pins"
                    element={<AdminProjectDetailForPinObjectList />}
                  />
                  <Route
                    exact
                    path="/admin/project/:id/pins/add"
                    element={<AdminProjectPinObjectAdd />}
                  />
                  {/*
                  <Route
                    exact
                    path="/admin/tenant/:id/pin/:aid"
                    element={<AdminTenantPinObjectDetail />}
                  />
                  <Route
                    exact
                    path="/admin/tenant/:id/pin/:aid/edit"
                    element={<AdminTenantPinObjectUpdate />}
                  />
                  */}

                  {/* User */}
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
                    path="/admin/user/:id/pinobjects"
                    element={<AdminUserDetailForPinObjectList />}
                  />
                  <Route
                    exact
                    path="/admin/user/:id/pinobject/:aid"
                    element={<AdminUserPinObjectDetail />}
                  />
                  <Route
                    exact
                    path="/admin/user/:id/pinobjects/add"
                    element={<AdminUserPinObjectAdd />}
                  />
                  <Route
                    exact
                    path="/admin/user/:id/pinobject/:aid/edit"
                    element={<AdminUserPinObjectUpdate />}
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
