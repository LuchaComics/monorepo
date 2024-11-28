import React, { useState, useEffect } from "react";
import { Navigate } from "react-router-dom";
import { useRecoilState } from 'recoil';
import Scroll from 'react-scroll';

import FormErrorBox from "./FormErrorBox";
import { postNFTAssetCreateAPI, deleteNFTAssetAPI } from "../../API/NFTAsset";
import { topAlertMessageState, topAlertStatusState } from "../../AppState";


function FormNFTAssetField({
    label,
    name,
    placeholder,
    filename="",
    setFilename=null,
    nftAssetID="",
    setNFTAssetID=null,
    errorText,
    validationText,
    helpText,
    maxWidth,
    disabled=false
}) {
    ////
    //// Global state.
    ////

    const [topAlertMessage, setTopAlertMessage] = useRecoilState(topAlertMessageState);
    const [topAlertStatus, setTopAlertStatus] = useRecoilState(topAlertStatusState);

    ////
    //// Component states.
    ////

    const [errors, setErrors] = useState({});
    const [isFetching, setFetching] = useState(false);
    const [forceURL, setForceURL] = useState("");

    ////
    //// Event handling.
    ////

    const onHandleFileChange = (event) => {
        console.log("onHandleFileChange: Starting...")
        setFetching(true);
        setErrors({});

        const selectedFile = event.target.files[0];

        const formData = new FormData();
        formData.append('file', selectedFile);

        // Extract filename.
        const filename = selectedFile.name;
        console.log('Filename:', filename);

        const mimeType = selectedFile.type || "application/octet-stream"; // Fallback to octet-stream if type is not detected
        console.log('MIME Type:', mimeType);


        // Convert the FormData object to binary data and pass it directly
         // Note: You may need to adjust the handling depending on the type of file and how it needs to be processed
         const reader = new FileReader();
         reader.onload = () => {
           const fileBinaryData = reader.result;

           postNFTAssetCreateAPI(
               filename,
               mimeType, // Pass the detected MIME type
               fileBinaryData, // Pass binary data instead of FormData
               onCreateSuccess,
               onCreateError,
               onCreateDone,
               onUnauthorized
           );

           console.log("onSubmitClick: Finished.");
           setFetching(false); // Reset fetching state after API call
         };
         reader.onerror = () => {
           console.error('Error reading file');
           setErrors({ file: 'Error reading file' });
           setFetching(false);
         };

         reader.readAsArrayBuffer(selectedFile); // Read file as binary data
        console.log("onSubmitClick: Finished.");
    };

    const onDeleteClick = () => {
        console.log("onDeleteClick: Starting...")
        setFetching(true);
        setErrors({});

        deleteNFTAssetAPI(
            nftAssetID,
            onDeleteSuccess,
            onDeleteError,
            onDeleteDone,
            onUnauthorized
        );
        console.log("onDeleteClick: Finished")
    }

    ////
    //// API.
    ////

    // --- Create --- //

    function onCreateSuccess(response){
        // For debugging purposes only.
        console.log("onCreateSuccess: Starting...");
        console.log("onCreateSuccess: ", response);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("File uploaded");
        setTopAlertStatus("success");
        setTimeout(() => {
            console.log("onAddSuccess: Delayed for 2 seconds.");
            console.log("onAddSuccess: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        setNFTAssetID(response.id);
        setFilename(response.filename);
    }

    function onCreateError(apiErr) {
        console.log("onCreateError: Starting...");
        setErrors(apiErr);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Failed submitting");
        setTopAlertStatus("danger");
        setTimeout(() => {
            console.log("onAddError: Delayed for 2 seconds.");
            console.log("onAddError: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onCreateDone() {
        console.log("onCreateDone: Starting...");
        setFetching(false);
    }

    // --- Delete --- //

    function onDeleteSuccess(response){
        // For debugging purposes only.
        console.log("onDeleteSuccess: Starting...");
        console.log(response);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("File deleted");
        setTopAlertStatus("success");
        setTimeout(() => {
            console.log("onAddSuccess: Delayed for 2 seconds.");
            console.log("onAddSuccess: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        setNFTAssetID("");
        setFilename("");
    }

    function onDeleteError(apiErr) {
        console.log("onDeleteError: Starting...");
        setErrors(apiErr);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("Failed submitting");
        setTopAlertStatus("danger");
        setTimeout(() => {
            console.log("onAddError: Delayed for 2 seconds.");
            console.log("onAddError: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        // The following code will cause the screen to scroll to the top of
        // the page. Please see ``react-scroll`` for more information:
        // https://github.com/fisshy/react-scroll
        var scroll = Scroll.animateScroll;
        scroll.scrollToTop();
    }

    function onDeleteDone() {
        console.log("onDeleteDone: Starting...");
        setFetching(false);
    }

    // --- All --- //

    const onUnauthorized = () => {
      setForceURL("/login?unauthorized=true"); // If token expired or collection is not logged in, redirect back to login.
    };

    ////
    //// Misc.
    ////

    useEffect(() => {
        let mounted = true;
    //
    //     if (mounted) {
    //         window.scrollTo(0, 0);  // Start the page at the top of the page.
    //     }
    //
        return () => { mounted = false; }
    }, []);

    ////
    //// Component rendering.
    ////


    if (forceURL !== "") {
        return <Navigate to={forceURL} />;
    }

    let classNameText = "input";
    if (errorText) {
        classNameText = "input is-danger";
    }
    return (
        <div class="field pb-4">
            <FormErrorBox errors={errors} />
            <label class="label">{label}</label>
            {isFetching
                ?
                <>
                <b>Uploading...</b>
                </>
                :
                <div class="control">
                    {nftAssetID !== undefined && nftAssetID !== null && nftAssetID !== ""
                        ?
                        <>{filename}&nbsp;<button className="is-fullwidth-mobile button is-small is-danger" type="button" onClick={onDeleteClick}>Delete</button></>
                        :
                        <input class={classNameText}
                                name={name}
                                type={"file"}
                         placeholder={placeholder}
                               value={nftAssetID}
                            onChange={onHandleFileChange}
                               style={{maxWidth:maxWidth}}
                            disabled={disabled}
                        autoComplete="off" />
                    }
                </div>
            }
            {errorText &&
                <p class="help is-danger">{errorText}</p>
            }
            {helpText &&
                <p class="help">{helpText}</p>
            }
        </div>
    );
}

export default FormNFTAssetField;
