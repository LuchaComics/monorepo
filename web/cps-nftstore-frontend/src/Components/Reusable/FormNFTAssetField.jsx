import React, { useState, useEffect } from "react";
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

    ////
    //// Event handling.
    ////

    const onHandleFileChange = (event) => {
        console.log("onHandleFileChange: Starting...")
        setFetching(true);
        setErrors({});

        const formData = new FormData();
        formData.append('file', event.target.files[0]);
        // formData.append('ownership_id', "");
        formData.append('ownership_type', "1");

        postNFTAssetCreateAPI(
            formData,
            onCreateSuccess,
            onCreateError,
            onCreateDone
        );
        console.log("onSubmitClick: Finished.")
    };

    const onDeleteClick = () => {
        console.log("onDeleteClick: Starting...")
        setFetching(true);
        setErrors({});

        deleteNFTAssetAPI(
            nftAssetID,
            onDeleteSuccess,
            onDeleteError,
            onDeleteDone
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
        console.log(response);

        // Add a temporary banner message in the app and then clear itself after 2 seconds.
        setTopAlertMessage("File uploaded");
        setTopAlertStatus("success");
        setTimeout(() => {
            console.log("onAddSuccess: Delayed for 2 seconds.");
            console.log("onAddSuccess: topAlertMessage, topAlertStatus:", topAlertMessage, topAlertStatus);
            setTopAlertMessage("");
        }, 2000);

        setNFTAssetID(response.id);
        setFilename(response.meta.filename);
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
