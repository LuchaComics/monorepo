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
import AdminTenantUpdate from "./Components/Admin/Tenant/Update";

// Collections
import AdminCollectionList from "./Components/Admin/Collection/List/View";
import AdminCollectionAdd from "./Components/Admin/Collection/Add/View";
import AdminCollectionDetail from "./Components/Admin/Collection/Detail/View";
import AdminCollectionUpdate from "./Components/Admin/Collection/Update/View";
import AdminCollectionDetailForNFTMetadataList from "./Components/Admin/Collection/Detail/NFTMetadata/ListView";
import AdminCollectionNFTMetadataAdd from "./Components/Admin/Collection/Detail/NFTMetadata/Add";
// import AdminCollectionNFTMetadataAddViaWebService from "./Components/Admin/Collection/Detail/NFTMetadata/AddViaWS";
// import AdminCollectionNFTMetadataDetail from "./Components/Admin/Collection/Detail/NFTMetadata/Detail";
// import AdminCollectionNFTMetadataUpdate from "./Components/Admin/Collection/Detail/NFTMetadata/Update";

import AdminCollectionDetailForNFTAssetList from "./Components/Admin/Collection/Detail/NFTAsset/ListView";
import AdminCollectionNFTAssetAdd from "./Components/Admin/Collection/Detail/NFTAsset/Add";
import AdminCollectionNFTAssetAddViaWebService from "./Components/Admin/Collection/Detail/NFTAsset/AddViaWS";
import AdminCollectionNFTAssetDetail from "./Components/Admin/Collection/Detail/NFTAsset/Detail";
import AdminCollectionNFTAssetUpdate from "./Components/Admin/Collection/Detail/NFTAsset/Update";
import AdminCollectionDetailMore from "./Components/Admin/Collection/Detail/More/View";

// Users
import AdminUserList from "./Components/Admin/User/List/View";
import AdminUserAdd from "./Components/Admin/User/Add/View";
import AdminUserDetail from "./Components/Admin/User/Detail/View";
import AdminUserDetailForCommentList from "./Components/Admin/User/Detail/Comment/View";
import AdminUserDetailMore from "./Components/Admin/User/Detail/More/View";
import AdminUserArchiveOperation from "./Components/Admin/User/Detail/More/Archive/View";
import AdminUserUnarchiveOperation from "./Components/Admin/User/Detail/More/Unarchive/View";
import AdminUserDeleteOperation from "./Components/Admin/User/Detail/More/Delete/View";
import AdminUserMoreOperationChangePassword from "./Components/Admin/User/Detail/More/ChangePassword/View";
import AdminUserMoreOperation2FAToggle from "./Components/Admin/User/Detail/More/2FA/View";
import AdminUserUpdate from "./Components/Admin/User/Update/View";

// IPFS
import AdminIPFSDashboard from "./Components/Admin/IPFS/Dashboard";

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
                  {/* Tenants */}
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

                  {/* Collection */}
                  <Route
                    exact
                    path="/admin/collections"
                    element={<AdminCollectionList />}
                  />
                  <Route
                    exact
                    path="/admin/collections/add"
                    element={<AdminCollectionAdd />}
                  />
                  <Route
                    exact
                    path="/admin/collection/:id"
                    element={<AdminCollectionDetail />}
                  />
                  <Route
                    exact
                    path="/admin/collection/:id/edit"
                    element={<AdminCollectionUpdate />}
                  />
                  <Route
                    exact
                    path="/admin/collection/:id/more"
                    element={<AdminCollectionDetailMore />}
                  />

                  {/* Collection NFT Metadata */}
                  <Route
                    exact
                    path="/admin/collection/:id/nft-metadata"
                    element={<AdminCollectionDetailForNFTMetadataList />}
                  />
                  <Route
                    exact
                    path="/admin/collection/:id/nft-metadata/add"
                    element={<AdminCollectionNFTMetadataAdd />}
                  />
                  {/*
                  <Route
                    exact
                    path="/admin/collection/:id/nfts/add-via-ws"
                    element={<AdminCollectionNFTMetadataAddViaWebService />}
                  />
                  <Route
                    exact
                    path="/admin/collection/:id/nft/:rid"
                    element={<AdminCollectionNFTMetadataDetail />}
                  />
                  <Route
                    exact
                    path="/admin/collection/:id/nft/:rid/edit"
                    element={<AdminCollectionNFTMetadataUpdate />}
                  />
                  */}

                  {/* Collection NFT Asset */}
                  <Route
                    exact
                    path="/admin/collection/:id/nft-assets"
                    element={<AdminCollectionDetailForNFTAssetList />}
                  />
                  <Route
                    exact
                    path="/admin/collection/:id/nft-assets/add"
                    element={<AdminCollectionNFTAssetAdd />}
                  />
                  {/*
                  <Route
                    exact
                    path="/admin/collection/:id/nfts/add-via-ws"
                    element={<AdminCollectionNFTMetadataAddViaWebService />}
                  />
                  <Route
                    exact
                    path="/admin/collection/:id/nft/:rid"
                    element={<AdminCollectionNFTMetadataDetail />}
                  />
                  <Route
                    exact
                    path="/admin/collection/:id/nft/:rid/edit"
                    element={<AdminCollectionNFTMetadataUpdate />}
                  />
                  */}

                  {/* Collection Pins */}
                  <Route
                    exact
                    path="/admin/collection/:id/pins"
                    element={<AdminCollectionDetailForNFTAssetList />}
                  />
                  <Route
                    exact
                    path="/admin/collection/:id/pins/add"
                    element={<AdminCollectionNFTAssetAdd />}
                  /><Route
                    exact
                    path="/admin/collection/:id/pins/add-via-ws"
                    element={<AdminCollectionNFTAssetAddViaWebService />}
                  />
                  <Route
                    exact
                    path="/admin/collection/:id/pin/:rid"
                    element={<AdminCollectionNFTAssetDetail />}
                  />
                  <Route
                    exact
                    path="/admin/collection/:id/pin/:rid/edit"
                    element={<AdminCollectionNFTAssetUpdate />}
                  />
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

                  {/* IPFS */}
                  <Route
                    exact
                    path="/admin/ipfs"
                    element={<AdminIPFSDashboard />}
                  />

                  {/* Dashboard */}
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
