import { useState, useEffect } from "react";
import { HashRouter, Routes, Route } from "react-router-dom";
import { RecoilRoot } from "recoil";

// CSS App Styling Override and extra.
import "./App.css";

// MENU
import Topbar from "./Components/Menu/Top";
import BottomTabBar from "./Components/Menu/BottomBar";

// CORE VIEWS
import InitializeView from "./Components/OnAppStart/InitializeView";
import NotFoundErrorView from "./Components/Other/NotFoundErrorView";
import PickDataDirectoryView from "./Components/OnAppStart/PickDataDirectoryView";
import StartupView from "./Components/OnAppStart/StartupView";
import CreateYourFirstWalletView from "./Components/OnAppStart/CreateYourFirstWalletView";
import DashboardView from "./Components/Dashboard/View";
import SendCoinView from "./Components/Send/SendCoinView";
import SendCoinProcessingView from "./Components/Send/SendCoinProcessingView";
import SendCoinSuccessView from "./Components/Send/SendCoinSuccessView";
import ReceiveView from "./Components/Receive/View";
import MoreView from "./Components/More/View";
import ListWalletsView from "./Components/More/Wallets/ListView";
import CreateWalletView from "./Components/More/Wallets/CreateView";
import ListTransactionsView from "./Components/More/Transactions/ListView";
import TransactionDetailView from "./Components/More/Transactions/DetailView";
import ListTokensView from "./Components/More/Tokens/ListView";
import TokenDetailView from "./Components/More/Tokens/DetailView";
import TokenTransferView from "./Components/More/Tokens/TransferView";
import TokenTransferSuccessView from "./Components/More/Tokens/TransferSuccessView";
import TokenBurnView from "./Components/More/Tokens/BurnView";
import SettingsView from "./Components/More/Settings/View";

function AppRoute() {
  return (
    <div className="min-h-screen bg-slate-50 flex flex-col">
      <RecoilRoot>
        <HashRouter basename={"/"}>
          {/* <TopAlertBanner /> */}

          {/* Top Navigation */}
          <Topbar />

          <Routes>
            <Route path="/" element={<InitializeView />} exact />
            <Route
              path="/pick-data-directory"
              element={<PickDataDirectoryView />}
              exact
            />
            <Route path="/startup" element={<StartupView />} exact />
            <Route
              path="/create-your-first-wallet"
              element={<CreateYourFirstWalletView />}
              exact
            />
            <Route path="/dashboard" element={<DashboardView />} exact />
            <Route path="/send" element={<SendCoinView />} exact />
            <Route
              path="/send-processing"
              element={<SendCoinProcessingView />}
              exact
            />
            <Route
              path="/send-success"
              element={<SendCoinSuccessView />}
              exact
            />
            <Route path="/receive" element={<ReceiveView />} exact />
            <Route path="/more" element={<MoreView />} exact />
            <Route path="/more/wallets" element={<ListWalletsView />} exact />
            <Route
              path="/more/wallets/add"
              element={<CreateWalletView />}
              exact
            />
            <Route
              path="/more/transactions"
              element={<ListTransactionsView />}
              exact
            />
            <Route
              path="/more/transaction/:timestamp"
              element={<TransactionDetailView />}
              exact
            />
            <Route path="/more/tokens" element={<ListTokensView />} exact />
            <Route
              path="/more/token/:tokenID"
              element={<TokenDetailView />}
              exact
            />
            <Route
              path="/more/token/:tokenID/burn"
              element={<TokenBurnView />}
              exact
            />
            <Route
              path="/more/token/:tokenID/transfer"
              element={<TokenTransferView />}
              exact
            />
            <Route
              path="/more/token/:tokenID/transfer-success"
              element={<TokenTransferSuccessView />}
              exact
            />
            <Route path="/more/settings" element={<SettingsView />} exact />
            <Route path="*" element={<NotFoundErrorView />} />
          </Routes>

          {/* Bottom Navigation */}
          <BottomTabBar />
        </HashRouter>
      </RecoilRoot>
    </div>
  );
}

export default AppRoute;
