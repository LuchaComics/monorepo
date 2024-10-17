import {useState, useEffect} from 'react';
import { HashRouter, Routes, Route } from "react-router-dom";

import logo from './assets/images/logo-universal.png';
import './App.css';
import {GetPageID} from "../wailsjs/go/main/App";

import Welcome from './Welcome';
import App from './App';
import StartupView from "./Components/Startup/View";
import PickStorageLocationOnStartupView from "./Components/Startup/PickStorageLocationOnStartup/View";
import DashboardView from "./Components/Dashboard/View";
import SendView from "./Components/Send/View";
import ReceiveView from "./Components/Receive/View";
import TransactionsView from "./Components/Transactions/View";
import NotFoundError from "./Components/Other/NotFoundError";


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
        <HashRouter basename={"/"}>
            <Routes>
                <Route path="/" element={<StartupView />} exact />
                <Route path="/pick-storage-location-on-startup" element={<PickStorageLocationOnStartupView />} exact />
                <Route path="/dashboard" element={<DashboardView />} exact />
                <Route path="/send" element={<SendView />} exact />
                <Route path="/receive" element={<ReceiveView />} exact />
                <Route path="/transactions" element={<TransactionsView />} exact />
                <Route path="*" element={<NotFoundError />} />
            </Routes>
        </HashRouter>
    )
}

export default AppRoute
