import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faEllipsis,
  faTasks,
  faTachometer,
  faPlus,
  faArrowLeft,
  faCheckCircle,
  faGauge,
  faPencil,
  faCubes,
  faEye,
  faArrowRight,
  faTrashCan,
  faArrowUpRightFromSquare,
  faFile,
  faDownload,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";
import {
  PIN_OBJECT_STATES,
  PAGE_SIZE_OPTIONS,
} from "../../../../../Constants/FieldOptions";

import { getCollectionDetailAPI } from "../../../../../API/Collection";
import {
  getNFTMetadataListAPI,
  deleteNFTMetadataAPI,
} from "../../../../../API/NFTMetadata";
import FormErrorBox from "../../../../Reusable/FormErrorBox";
import FormInputField from "../../../../Reusable/FormInputField";
import FormTextareaField from "../../../../Reusable/FormTextareaField";
import FormRadioField from "../../../../Reusable/FormRadioField";
import FormMultiSelectField from "../../../../Reusable/FormMultiSelectField";
import FormSelectField from "../../../../Reusable/FormSelectField";
import FormCheckboxField from "../../../../Reusable/FormCheckboxField";
import PageLoadingContent from "../../../../Reusable/PageLoadingContent";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../../AppState";
import AdminCollectionDetailForNFTMetadataListDesktop from "./DetailForNFTMetadataListDektop";
import AdminCollectionDetailForNFTMetadataListMobile from "./DetailForNFTMetadataListMobile";
import AlertBanner from "../../../../Reusable/EveryPage/AlertBanner";

function AdminCollectionDetailForNFTMetadataList() {
  ////
  //// URL Parameters.
  ////

  const { id } = useParams();

  ////
  //// Global state.
  ////

  const [topAlertMessage, setTopAlertMessage] =
    useRecoilState(topAlertMessageState);
  const [topAlertStatus, setTopAlertStatus] =
    useRecoilState(topAlertStatusState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [collection, setCollection] = useState({});
  const [tabIndex, setTabIndex] = useState(1);
  const [nftAssets, setNFTMetadatas] = useState("");
  const [selectedNFTMetadataForDeletion, setSelectedNFTMetadataForDeletion] =
    useState("");
  const [pageSize, setPageSize] = useState(10); // Pagination
  const [previousCursors, setPreviousCursors] = useState([]); // Pagination
  const [nextCursor, setNextCursor] = useState(""); // Pagination
  const [currentCursor, setCurrentCursor] = useState(""); // Pagination

  ////
  //// Event handling.
  ////

  const fetchNFTMetadataList = (cur, collectionID, limit) => {
    setFetching(true);
    setErrors({});

    let params = new Map();
    params.set("collection_id", id);
    params.set("page_size", limit);
    if (cur !== "") {
      params.set("cursor", cur);
    }

    getNFTMetadataListAPI(
      params,
      onNFTMetadataListSuccess,
      onNFTMetadataListError,
      onNFTMetadataListDone,
      onUnauthorized,
    );
  };

  const onNextClicked = (e) => {
    console.log("onNextClicked");
    let arr = [...previousCursors];
    arr.push(currentCursor);
    setPreviousCursors(arr);
    setCurrentCursor(nextCursor);
  };

  const onPreviousClicked = (e) => {
    console.log("onPreviousClicked");
    let arr = [...previousCursors];
    const previousCursor = arr.pop();
    setPreviousCursors(arr);
    setCurrentCursor(previousCursor);
  };

  const onSelectNFTMetadataForDeletion = (e, nftAsset) => {
    console.log("onSelectNFTMetadataForDeletion", nftAsset);
    setSelectedNFTMetadataForDeletion(nftAsset);
  };

  const onDeselectNFTMetadataForDeletion = (e) => {
    console.log("onDeselectNFTMetadataForDeletion");
    setSelectedNFTMetadataForDeletion("");
  };

  const onDeleteConfirmButtonClick = (e) => {
    console.log("onDeleteConfirmButtonClick"); // For debugging purposes only.

    deleteNFTMetadataAPI(
      selectedNFTMetadataForDeletion.requestid,
      onNFTMetadataDeleteSuccess,
      onNFTMetadataDeleteError,
      onNFTMetadataDeleteDone,
      onUnauthorized,
    );
    setSelectedNFTMetadataForDeletion("");
  };

  ////
  //// API.
  ////

  // Collection details.

  function onCollectionDetailSuccess(response) {
    console.log("onCollectionDetailSuccess: Starting...");
    setCollection(response);
  }

  function onCollectionDetailError(apiErr) {
    console.log("onCollectionDetailError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onCollectionDetailDone() {
    console.log("onCollectionDetailDone: Starting...");
    setFetching(false);
  }

  // NFTMetadata list.

  function onNFTMetadataListSuccess(response) {
    console.log("onNFTMetadataListSuccess: Starting...");
    if (response.results !== null) {
      setNFTMetadatas(response);
      if (response.hasNextPage) {
        setNextCursor(response.nextCursor); // For pagination purposes.
      }
    }
  }

  function onNFTMetadataListError(apiErr) {
    console.log("onNFTMetadataListError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onNFTMetadataListDone() {
    console.log("onNFTMetadataListDone: Starting...");
    setFetching(false);
  }

  // NFTMetadata delete.

  function onNFTMetadataDeleteSuccess(response) {
    console.log("onNFTMetadataDeleteSuccess: Starting..."); // For debugging purposes only.

    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("NFTMetadata deleted");
    setTimeout(() => {
      console.log(
        "onDeleteConfirmButtonClick: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Fetch again an updated list.
    fetchNFTMetadataList(currentCursor, id, pageSize);
  }

  function onNFTMetadataDeleteError(apiErr) {
    console.log("onNFTMetadataDeleteError: Starting..."); // For debugging purposes only.
    setErrors(apiErr);

    // Update notification.
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed deleting");
    setTimeout(() => {
      console.log(
        "onNFTMetadataDeleteError: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onNFTMetadataDeleteDone() {
    console.log("onNFTMetadataDeleteDone: Starting...");
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

    if (mounted) {
      window.scrollTo(0, 0); // Start the page at the top of the page.

      setFetching(true);
      getCollectionDetailAPI(
        id,
        onCollectionDetailSuccess,
        onCollectionDetailError,
        onCollectionDetailDone,
        onUnauthorized,
      );
      fetchNFTMetadataList(currentCursor, id, pageSize);
    }

    return () => {
      mounted = false;
    };
  }, [currentCursor, id, pageSize]);

  ////
  //// Component rendering.
  ////

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <>
      <div class="container">
        <section class="section">
          {/* Desktop Breadcrumbs */}
          <nav class="breadcrumb is-hidden-touch" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to="/admin/dashboard" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faGauge} />
                  &nbsp;Admin Dashboard
                </Link>
              </li>
              <li class="">
                <Link to="/admin/collections" aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faCubes} />
                  &nbsp;Collections
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail (NFT Metadata)
                </Link>
              </li>
            </ul>
          </nav>

          {/* Mobile Breadcrumbs */}
          <nav class="breadcrumb is-hidden-desktop" aria-label="breadcrumbs">
            <ul>
              <li class="">
                <Link to={`/admin/collections`} aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                  &nbsp;Back to Collections
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modals */}
          <div
            class={`modal ${selectedNFTMetadataForDeletion ? "is-active" : ""}`}
          >
            <div class="modal-background"></div>
            <div class="modal-card">
              <header class="modal-card-head">
                <p class="modal-card-title">Are you sure?</p>
                <button
                  class="delete"
                  aria-label="close"
                  onClick={onDeselectNFTMetadataForDeletion}
                ></button>
              </header>
              <section class="modal-card-body">
              You are about to <b>delete</b> this pin; the data will be permanently deleted and no
              longer appear on your dashboard. This action cannot be undone. Are you sure
              you would like to continue?
              </section>
              <footer class="modal-card-foot">
                <button
                  class="button is-success"
                  onClick={onDeleteConfirmButtonClick}
                >
                  Confirm
                </button>
                <button
                  class="button"
                  onClick={onDeselectNFTMetadataForDeletion}
                >
                  Cancel
                </button>
              </footer>
            </div>
          </div>

          {/* Page banner */}
          {collection && collection.status === 100 && (
            <AlertBanner message="Archived" status="info" />
          )}

          {/* Page */}
          <nav class="box">
            <div class="columns">
              <div class="column">
                <p class="title is-4">
                  <FontAwesomeIcon className="fas" icon={faCubes} />
                  &nbsp;Collection
                </p>
              </div>
              {collection && collection.status === 1 && (
                <div class="column has-text-right">
                  <Link
                    to={`/admin/collection/${id}/nft-metadata/add`}
                    class="button is-small is-success is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faPlus} />
                    &nbsp;Add NFT Metadata
                  </Link>
                </div>
              )}
            </div>
            <FormErrorBox errors={errors} />

            {/* <p class="pb-4">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                {collection && (
                  <div class="container">
                    <div class="tabs is-medium is-size-7-mobile">
                      <ul>
                        <li>
                          <Link to={`/admin/collection/${collection.id}`}>Detail</Link>
                        </li>
                        <li class="is-active">
                          <Link to={`/admin/collection/${collection.id}/nft-metadata`}>
                            <b>NFT Metadata</b>
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/collection/${collection.id}/nft-assets`}>
                            NFT Assets
                          </Link>
                        </li>
                        <li>
                          <Link to={`/admin/collection/${collection.id}/more`}>
                            More&nbsp;&nbsp;
                            <FontAwesomeIcon
                              className="mdi"
                              icon={faEllipsis}
                            />
                          </Link>
                        </li>
                      </ul>
                    </div>

                    {!isFetching &&
                    nftAssets &&
                    nftAssets.results &&
                    (nftAssets.results.length > 0 ||
                      previousCursors.length > 0) ? (
                      <div class="container">
                        {/*
                            ##################################################################
                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A DESKTOP SCREEN.
                            ##################################################################
                        */}
                        <div class="is-hidden-touch">
                          <AdminCollectionDetailForNFTMetadataListDesktop
                            collectionID={id}
                            listData={nftAssets}
                            setPageSize={setPageSize}
                            pageSize={pageSize}
                            previousCursors={previousCursors}
                            onPreviousClicked={onPreviousClicked}
                            onNextClicked={onNextClicked}
                            onSelectNFTMetadataForDeletion={
                              onSelectNFTMetadataForDeletion
                            }
                          />
                        </div>

                        {/*
                            ###########################################################################
                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A TABLET OR MOBILE SCREEN.
                            ###########################################################################
                        */}
                        <div class="is-fullwidth is-hidden-desktop">
                          <AdminCollectionDetailForNFTMetadataListMobile
                            collectionID={id}
                            listData={nftAssets}
                            setPageSize={setPageSize}
                            pageSize={pageSize}
                            previousCursors={previousCursors}
                            onPreviousClicked={onPreviousClicked}
                            onNextClicked={onNextClicked}
                            onSelectNFTMetadataForDeletion={
                              onSelectNFTMetadataForDeletion
                            }
                          />
                        </div>
                      </div>
                    ) : (
                      <div class="container">
                        <article class="message is-dark">
                          <div class="message-body">
                            No NFT metadata.{" "}
                            <b>
                              <Link to={`/admin/collection/${id}/nft-metadata/add`}>
                                Click here&nbsp;
                                <FontAwesomeIcon
                                  className="mdi"
                                  icon={faArrowRight}
                                />
                              </Link>
                            </b>{" "}
                            to get started creating a new nft metadata.
                          </div>
                        </article>
                      </div>
                    )}

                    <div class="columns pt-5">
                      <div class="column is-half">
                        <Link class="button is-fullwidth-mobile" to={`/collections`}>
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back to Collections
                        </Link>
                      </div>
                      <div class="column is-half has-text-right">
                        {collection && collection.status === 1 && <Link
                          to={`/admin/collection/${id}/nft-metadata/add`}
                          class="button is-primary is-fullwidth-mobile"
                        >
                          <FontAwesomeIcon className="fas" icon={faPlus} />
                          &nbsp;Add NFT Metadata
                        </Link>}
                      </div>
                    </div>
                  </div>
                )}
              </>
            )}
          </nav>

          {/* Bottom Page Logout Link  */}
          {/*
          <div className="has-text-right has-text-grey">
            <Link to={`/admin/collection/${id}/nft-metadata/add-via-ws`} className="has-text-grey">
              Add NFTMetadata via Web-Service API&nbsp;
              <FontAwesomeIcon className="mdi" icon={faArrowRight} />
            </Link>
          </div>
          */}
        </section>
      </div>
    </>
  );
}

export default AdminCollectionDetailForNFTMetadataList;
