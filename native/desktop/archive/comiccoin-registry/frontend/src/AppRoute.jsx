import {useState, useEffect} from 'react';
import { HashRouter, Routes, Route } from "react-router-dom";
import 'bulma/css/bulma.min.css';
import { RecoilRoot } from "recoil";

// CSS App Styling Override and extra.
import './App.css';

// MENU
import Topbar from "./Components/Menu/Top";
import BottomTabBar  from "./Components/Menu/BottomBar";

// CORE VIEWS
import InitializeView from "./Components/Other/InitializeView";
import NotFoundErrorView from "./Components/Other/NotFoundErrorView";
import PickDataDirectoryView from "./Components/Other/PickDataDirectoryView";
import NFTSToreSetupView from "./Components/Other/NFTSToreSetupView";
import StartupView from "./Components/Other/StartupView";
import DashboardView from "./Components/Dashboard/View";
import ListTokensView from "./Components/Tokens/ListView";
import CreateTokenView from "./Components/Tokens/CreateView";
import TokenDetailView from "./Components/Tokens/DetailView";


function AppRoute() {
    // const [pageID, setPageID] = useState("PageID");
    // const updatePageID = (result) => setPageID(result);

    // useEffect(() => {
    //   let mounted = true;
    //
    //   if (mounted) {
    //         window.scrollTo(0, 0); // Start the page at the top of the page.
    //
    //
    //   }
    //
    //       GetPageID().then(updatePageID);
    //
    //
    //   return () => {
    //     mounted = false;
    //   };
    // }, []);
    //
    // console.log("---> ", pageID);


    return (
        <div class="is-widescreen is-size-5-widescreen is-size-6-tablet is-size-7-mobile theme-light">
            {/*
                NOTES FOR ABOVE
                USE THE FOLLOWING TEXT SIZES BASED ON DEVICE TYPE
                - is-size-5-widescreen
                - is-size-6-tablet
                - is-size-7-mobile

                FURTHERMORE, WE ARE FORCING LIGHT MODE.
            */}
            <RecoilRoot>
                <HashRouter basename={"/"}>
                    {/*
                        <AnonymousCurrentUserRedirector />
                        <TwoFactorAuthenticationRedirector />
                        <TopAlertBanner />
                    */}
                    <Topbar />
                    <section class="main-content pt-6 is-fullheight">
                        <Routes>
                            <Route path="/" element={<InitializeView />} exact />
                            <Route path="/pick-data-directory" element={<PickDataDirectoryView />} exact />
                            <Route path="/config-nftstore" element={<NFTSToreSetupView />} exact />
                            <Route path="/startup" element={<StartupView />} exact />
                            <Route path="/dashboard" element={<DashboardView />} exact />
                            <Route path="/tokens" element={<ListTokensView />} exact />
                            <Route path="/tokens/new" element={<CreateTokenView />} exact />
                            <Route path="/token/:id" element={<TokenDetailView />} exact />
                            {/*
                            <Route path="/wallets" element={<ListWalletsView />} exact />
                            <Route path="/wallet/add" element={<CreateWalletView />} exact />
                            <Route path="/send" element={<SendCoinSubmissionView />} exact />
                            <Route path="/send-success" element={<SendCoinSuccessView />} exact />
                            <Route path="/receive" element={<ReceiveView />} exact />
                            <Route path="/more" element={<MoreView />} exact />
                            <Route path="/more/transactions" element={<ListTransactionsView />} exact />
                            <Route path="/transactions" element={<ListTransactionsView />} exact />
                            <Route path="/tokens" element={<ListTokensView />} exact />
                            <Route path="/settings" element={<SettingsView />} exact />
                            */}
                            <Route path="*" element={<NotFoundErrorView />} />
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
                    <BottomTabBar />
                </HashRouter>
            </RecoilRoot>
        </div>
    )
}

export default AppRoute
