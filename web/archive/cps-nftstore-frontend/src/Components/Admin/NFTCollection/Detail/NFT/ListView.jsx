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

import { getCollectionDetailAPI } from "../../../../../API/NFTCollection";
import {
  getNFTListAPI,
  deleteNFTAPI,
} from "../../../../../API/NFT";
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
import AdminNFTCollectionDetailForNFTListDesktop from "./ListDektopView";
import AdminNFTCollectionDetailForNFTListMobile from "./ListMobileView";
import AlertBanner from "../../../../Reusable/EveryPage/AlertBanner";

function AdminNFTCollectionDetailForNFTList() {
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
  const [nftAssets, setNFTs] = useState("");
  const [selectedNFTForDeletion, setSelectedNFTForDeletion] =
    useState("");
  const [pageSize, setPageSize] = useState(10); // Pagination
  const [previousCursors, setPreviousCursors] = useState([]); // Pagination
  const [nextCursor, setNextCursor] = useState(""); // Pagination
  const [currentCursor, setCurrentCursor] = useState(""); // Pagination

  ////
  //// Event handling.
  ////

  const fetchNFTList = (cur, collectionID, limit) => {
    setFetching(true);
    setErrors({});

    let params = new Map();
    params.set("collection_id", id);
    params.set("page_size", limit);
    if (cur !== "") {
      params.set("cursor", cur);
    }

    getNFTListAPI(
      params,
      onNFTListSuccess,
      onNFTListError,
      onNFTListDone,
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

  const onSelectNFTForDeletion = (e, nftAsset) => {
    console.log("onSelectNFTForDeletion", nftAsset);
    setSelectedNFTForDeletion(nftAsset);
  };

  const onDeselectNFTForDeletion = (e) => {
    console.log("onDeselectNFTForDeletion");
    setSelectedNFTForDeletion("");
  };

  const onArchiveConfirmButtonClick = (e) => {
    console.log("onArchiveConfirmButtonClick"); // For debugging purposes only.

    deleteNFTAPI(
      selectedNFTForDeletion.requestid,
      onNFTArchiveSuccess,
      onNFTArchiveError,
      onNFTArchiveDone,
      onUnauthorized,
    );
    setSelectedNFTForDeletion("");
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

  // NFT list.

  function onNFTListSuccess(response) {
    console.log("onNFTListSuccess: Starting...");
    if (response.results !== null) {
      setNFTs(response);
      if (response.hasNextPage) {
        setNextCursor(response.nextCursor); // For pagination purposes.
      }
    }
  }

  function onNFTListError(apiErr) {
    console.log("onNFTListError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onNFTListDone() {
    console.log("onNFTListDone: Starting...");
    setFetching(false);
  }

  // NFT archive.

  function onNFTArchiveSuccess(response) {
    console.log("onNFTArchiveSuccess: Starting..."); // For debugging purposes only.

    // Update notification.
    setTopAlertStatus("success");
    setTopAlertMessage("NFT archived");
    setTimeout(() => {
      console.log(
        "onArchiveConfirmButtonClick: topAlertMessage, topAlertStatus:",
        topAlertMessage,
        topAlertStatus,
      );
      setTopAlertMessage("");
    }, 2000);

    // Fetch again an updated list.
    fetchNFTList(currentCursor, id, pageSize);
  }

  function onNFTArchiveError(apiErr) {
    console.log("onNFTArchiveError: Starting..."); // For debugging purposes only.
    setErrors(apiErr);

    // Update notification.
    setTopAlertStatus("danger");
    setTopAlertMessage("Failed deleting");
    setTimeout(() => {
      console.log(
        "onNFTArchiveError: topAlertMessage, topAlertStatus:",
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

  function onNFTArchiveDone() {
    console.log("onNFTArchiveDone: Starting...");
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
      fetchNFTList(currentCursor, id, pageSize);
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
                  &nbsp;NFT Collections
                </Link>
              </li>
              <li class="is-active">
                <Link aria-current="page">
                  <FontAwesomeIcon className="fas" icon={faEye} />
                  &nbsp;Detail (NFTs)
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
                  &nbsp;Back to NFT Collections
                </Link>
              </li>
            </ul>
          </nav>

          {/* Modals */}
          <div
            class={`modal ${selectedNFTForDeletion ? "is-active" : ""}`}
          >
            <div class="modal-background"></div>
            <div class="modal-card">
              <header class="modal-card-head">
                <p class="modal-card-title">Are you sure?</p>
                <button
                  class="archive"
                  aria-label="close"
                  onClick={onDeselectNFTForDeletion}
                ></button>
              </header>
              <section class="modal-card-body">
                  You are about to <b>archive</b> this MFT metadata; it will no
                  longer appear on your dashboard. This action cannot be undone. Are you sure
                  you would like to continue?
              </section>
              <footer class="modal-card-foot">
                <button
                  class="button is-success"
                  onClick={onArchiveConfirmButtonClick}
                >
                  Confirm
                </button>
                <button
                  class="button"
                  onClick={onDeselectNFTForDeletion}
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
                  &nbsp;NFT Collection
                </p>
              </div>
              {collection && collection.status === 1 && (
                <div class="column has-text-right">
                  <Link
                    to={`/admin/collection/${id}/nfts/add/step-1`}
                    class="button is-small is-success is-fullwidth-mobile"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faPlus} />
                    &nbsp;New NFT
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
                          <Link to={`/admin/collection/${collection.id}/nfts`}>
                            <b>NFTs</b>
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
                          <AdminNFTCollectionDetailForNFTListDesktop
                            collectionID={id}
                            listData={nftAssets}
                            setPageSize={setPageSize}
                            pageSize={pageSize}
                            previousCursors={previousCursors}
                            onPreviousClicked={onPreviousClicked}
                            onNextClicked={onNextClicked}
                            onSelectNFTForDeletion={
                              onSelectNFTForDeletion
                            }
                          />
                        </div>

                        {/*
                            ###########################################################################
                            EVERYTHING INSIDE HERE WILL ONLY BE DISPLAYED ON A TABLET OR MOBILE SCREEN.
                            ###########################################################################
                        */}
                        <div class="is-fullwidth is-hidden-desktop">
                          <AdminNFTCollectionDetailForNFTListMobile
                            collectionID={id}
                            listData={nftAssets}
                            setPageSize={setPageSize}
                            pageSize={pageSize}
                            previousCursors={previousCursors}
                            onPreviousClicked={onPreviousClicked}
                            onNextClicked={onNextClicked}
                            onSelectNFTForDeletion={
                              onSelectNFTForDeletion
                            }
                          />
                        </div>
                      </div>
                    ) : (
                      <div class="container">
                        {collection.smartContractStatus === 2 ?
                            <>
                                <article class="message is-dark">
                                  <div class="message-body">
                                    No NFTs.{" "}
                                    <b>
                                      <Link to={`/admin/collection/${id}/nfts/add/step-1`}>
                                        Click here&nbsp;
                                        <FontAwesomeIcon
                                          className="mdi"
                                          icon={faArrowRight}
                                        />
                                      </Link>
                                    </b>{" "}
                                    to get started creating your first NFT.
                                  </div>
                                </article>
                            </>
                            :
                            <>
                                <article class="message is-warning">
                                  <div class="message-body">
                                    NFT collection not deployed to blockchain.{" "}
                                    <b>
                                      <Link to={`/admin/collection/${id}/more/deploy`}>
                                        Click here&nbsp;
                                        <FontAwesomeIcon
                                          className="mdi"
                                          icon={faArrowRight}
                                        />
                                      </Link>
                                    </b>{" "}
                                    to get started deploying.
                                  </div>
                                </article>
                            </>
                        }
                      </div>
                    )}

                    <div class="columns pt-5">
                      <div class="column is-half">
                        <Link class="button is-fullwidth-mobile" to={`/collections`}>
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back to NFT Collections
                        </Link>
                      </div>
                      <div class="column is-half has-text-right">
                        {collection && collection.status === 1 && <Link
                          to={`/admin/collection/${id}/nfts/add/step-1`}
                          class="button is-primary is-fullwidth-mobile"
                        >
                          <FontAwesomeIcon className="fas" icon={faPlus} />
                          &nbsp;New NFT
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
            <Link to={`/admin/collection/${id}/nfts/add/step-1-via-ws`} className="has-text-grey">
              NFT via Web-Service API&nbsp;
              <FontAwesomeIcon className="mdi" icon={faArrowRight} />
            </Link>
          </div>
          */}
        </section>
      </div>
    </>
  );
}

export default AdminNFTCollectionDetailForNFTList;
