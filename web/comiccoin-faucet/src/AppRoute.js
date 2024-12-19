import { React, useState } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { RecoilRoot } from "recoil";

// Gateway
import Register from "./Components/Gateway/Register";
import RegisterSuccessful from "./Components/Gateway/RegisterSuccessful";
import Login from "./Components/Gateway/Login";
import LogoutRedirector from "./Components/Gateway/LogoutRedirector";
// import RegisterLaunchpad from "./Components/Gateway/Register/Launchpad";
// import RegisterAsStoreOwner from "./Components/Gateway/Register/StoreOwner";
// import RegisterAsCustomer from "./Components/Gateway/Register/Customer";
import EmailVerification from "./Components/Gateway/EmailVerification";
// import ForgotPassword from "./Components/Gateway/ForgotPassword";
// import PasswordReset from "./Components/Gateway/PasswordReset";

// // Gateway (2FA specific)
// import TwoFactorAuthenticationWizardStep1 from "./Components/Gateway/2FA/Step1";
// import TwoFactorAuthenticationWizardStep2 from "./Components/Gateway/2FA/Step2";
// import TwoFactorAuthenticationWizardStep3 from "./Components/Gateway/2FA/Step3";
// import TwoFactorAuthenticationBackupCodeGenerate from "./Components/Gateway/2FA/BackupCodeGenerate";
// import TwoFactorAuthenticationBackupCodeRecovery from "./Components/Gateway/2FA/BackupCodeRecovery";
// import TwoFactorAuthenticationValidateOnLogin from "./Components/Gateway/2FA/ValidateOnLogin";

// // Redirectors.
import AnonymousCurrentUserRedirector from "./Components/Misc/AnonymousCurrentUserRedirector";
import TwoFactorAuthenticationRedirector from "./Components/Misc/TwoFactorAuthenticationRedirector";
import FaucetAddWalletRedirector from "./Components/Gateway/FaucetAddWalletRedirector";

// Public Generic
import Index from "./Components/Gateway/Index";
import Terms from "./Components/Gateway/Terms";
import Privacy from "./Components/Gateway/Privacy";
import NotImplementedError from "./Components/Misc/NotImplementedError";
import NotFoundError from "./Components/Misc/NotFoundError";
// import DashboardHelp from "./Components/Misc/DashboardHelp";

// On first-time-use wizard.
import UserAddWalletToFaucet from "./Components/Gateway/FaucetAddWallet";
import UserAddWalletToFaucetSuccess from "./Components/Gateway/FaucetAddWalletSuccess";

// User facing app
import DashboardPage from "./Components/User/Dashboard";
import SubmitComicPage from "./Components/User/Submission/Submit";
import SubmitComicSuccessPage from "./Components/User/Submission/SubmitSuccess";
import ApplyForVerificationPage from "./Components/User/Verification/View";
import SubmissionsPage from "./Components/User/Submission/List";
import MyWalletPage from "./Components/User/MyWallet/View";
import HelpPage from "./Components/User/Help/View";
import SettingsPage from "./Components/User/Settings/SettingsView";
import InfoPage from "./Components/User/Settings/InfoView";
import PassPage from "./Components/User/Settings/PassView";

// Admin facing app
import AdminDashboardPage from "./Components/Admin/Dashboard";
import AdminSettingsPage from "./Components/Admin/Settings/SettingsView";
import AdminInfoPage from "./Components/Admin/Settings/InfoView";
import AdminChangePasswordPage from "./Components/Admin/Settings/PassView";

// //-----------------//
// //    App Routes   //
// //-----------------//

function AppRoute() {
  return (
    <div>
      <RecoilRoot>
        <Router>
          <AnonymousCurrentUserRedirector />
          <TwoFactorAuthenticationRedirector />
          <FaucetAddWalletRedirector />
          <Routes>
            <Route exact path="/register" element={<Register />} />
            <Route
              exact
              path="/register-successful"
              element={<RegisterSuccessful />}
            />
            <Route exact path="/login" element={<Login />} />
            <Route exact path="/logout" element={<LogoutRedirector />} />
            <Route exact path="/verify" element={<EmailVerification />} />
            <Route exact path="/terms" element={<Terms />} />
            <Route exact path="/privacy" element={<Privacy />} />
            <Route
              exact
              path="/add-my-wallet-to-faucet"
              element={<UserAddWalletToFaucet />}
            />
            <Route
              exact
              path="/added-my-wallet-to-faucet-successfully"
              element={<UserAddWalletToFaucetSuccess />}
            />
            <Route exact path="/dashboard" element={<DashboardPage />} />
            <Route exact path="/submit" element={<SubmitComicPage />} />
            <Route
              exact
              path="/submit/success"
              element={<SubmitComicSuccessPage />}
            />
            <Route
              exact
              path="/apply-for-verification"
              element={<ApplyForVerificationPage />}
            />
            <Route exact path="/submissions" element={<SubmissionsPage />} />
            <Route exact path="/my-wallet" element={<MyWalletPage />} />
            <Route exact path="/help" element={<HelpPage />} />
            <Route exact path="/settings" element={<SettingsPage />} />
            <Route exact path="/settings/info" element={<InfoPage />} />
            <Route exact path="/settings/pass" element={<PassPage />} />
            <Route exact path="/admin/dashboard" element={<AdminDashboardPage />} />
            <Route exact path="/admin/settings" element={<AdminSettingsPage />} />
            <Route exact path="/admin/settings/info" element={<AdminInfoPage />} />
            <Route exact path="/admin/settings/pass" element={<AdminChangePasswordPage />} />
            <Route exact path="/" element={<Index />} />
            <Route path="*" element={<NotFoundError />} />
          </Routes>
        </Router>
      </RecoilRoot>
    </div>
  );
}

export default AppRoute;
