import { useState, useEffect, useRef } from "react";
import { Navigate } from "react-router-dom";
import PageLoadingContent from "../Reusable/PageLoadingContent";
import { ListAllPendingSignedTransactions } from "../../../wailsjs/go/main/App";

function SendCoinProcessingView() {
    ////
    //// Component states.
    ////

    const [forceURL, setForceURL] = useState("");
    const intervalRef = useRef(null); // Mutable reference to store interval ID.

    ////
    //// Event handling.
    ////

    const backgroundPollingTick = () => {
        ListAllPendingSignedTransactions().then((listResp) => {
            if (listResp.length > 0) {
                console.log("SendCoinProcessingView: tick", new Date().getTime(), listResp);
            } else {
                console.log("SendCoinProcessingView: tick: done");
                if (intervalRef.current) {
                    clearInterval(intervalRef.current); // Clear the interval.
                    intervalRef.current = null; // Reset the ref.
                }
                setForceURL("/send-success"); // Redirect to success page.
            }
        });
    };

    ////
    //// Misc.
    ////

    useEffect(() => {
        // Start polling when the component mounts.
        intervalRef.current = setInterval(backgroundPollingTick, 2000);

        return () => {
            // Clear the interval when the component unmounts.
            if (intervalRef.current) {
                clearInterval(intervalRef.current);
                intervalRef.current = null;
            }
        };
    }, []); // Empty dependency array ensures this effect runs only once on mount.

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

export default SendCoinProcessingView;
