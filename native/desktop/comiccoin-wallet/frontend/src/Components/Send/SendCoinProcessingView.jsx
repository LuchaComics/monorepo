import {useState, useEffect} from 'react';
import { Link, Navigate } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faGauge,
  faArrowRight,
  faUsers,
  faBarcode,
  faCubes,
  faPaperPlane,
  faTimesCircle,
  faCheckCircle,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import logo from '../../assets/images/CPS-logo-2023-square.webp';
import FormErrorBox from "../Reusable/FormErrorBox";
import FormInputField from "../Reusable/FormInputField";
import FormRadioField from "../Reusable/FormRadioField";
import FormTextareaField from "../Reusable/FormTextareaField";
import { ListAllPendingSignedTransactions } from "../../../wailsjs/go/main/App";
import { currentOpenWalletAtAddressState } from "../../AppState";
import PageLoadingContent from "../Reusable/PageLoadingContent";


function SendCoinProcessingView() {
    ////
    //// Global State
    ////

    ////
    //// Component states.
    ////

    // GUI States.
    const [errors, setErrors] = useState({});
    const [forceURL, setForceURL] = useState("");
    const [isLoading, setIsLoading] = useState(false);
    const [intervalId, setIntervalId] = useState(null);

    ////
    //// Event handling.
    ////

    // Function will make a call to check to see if our node is running
    // and if our backend says the node is running then we will redirect
    // to another page.
    const backgroundPollingTick = (e) => {
        ListAllPendingSignedTransactions().then( (listResp)=>{
            if (listResp.length > 0) {
                console.log("SendCoinProcessingView: tick", count, new Date().getTime(), listResp);
            } else {
                console.log("SendCoinProcessingView: tick: done");
                clearInterval(intervalId);
                setIntervalId(null);
                setForceURL("/send-success");
            }
        })
    }

    ////
    //// API.
    ////

    const onSubmitClick = (e) => {
        e.preventDefault();


    }

    ////
    //// Misc.
    ////

    useEffect(() => {
      let mounted = true;

      if (mounted) {
          window.scrollTo(0, 0); // Start the page at the top of the page.

          const interval = setInterval(() => backgroundPollingTick(), 2000);
          setIntervalId(interval);
      }

      return () => {
          mounted = false;
      };
    }, []);

    ////
    //// Component rendering.
    ////

    ////
    //// Component rendering.
    ////

    if (forceURL !== "") {
      return <Navigate to={forceURL} />;
    }

    return (
        <>
            <PageLoadingContent displayMessage="Processing..." />
        </>
    );
}

export default SendCoinProcessingView
