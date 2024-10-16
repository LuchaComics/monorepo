import {useState, useEffect} from 'react';
import { HashRouter, Routes, Route } from "react-router-dom";

import logo from './assets/images/logo-universal.png';
import './App.css';
import {GetPageID} from "../wailsjs/go/main/App";

import Welcome from './Welcome';
import App from './App';
import NotFoundError from "./NotFoundError";


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
                <Route path="/" element={<App />} exact />
                <Route path="/welcome" element={<Welcome />} exact />
                <Route path="*" element={<NotFoundError />} />
            </Routes>
        </HashRouter>
    )
}

export default AppRoute
