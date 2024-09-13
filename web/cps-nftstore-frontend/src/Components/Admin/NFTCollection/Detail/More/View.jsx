import React, { useState, useEffect } from "react";
import { Link, Navigate } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faChevronRight,
  faMobile,
  faKey,
  faBuildingCollection,
  faImage,
  faPaperclip,
  faAddressCard,
  faSquarePhone,
  faTasks,
  faTachometer,
  faPlus,
  faArrowLeft,
  faCheckCircle,
  faCubes,
  faGauge,
  faPencil,
  faEye,
  faIdCard,
  faAddressBook,
  faContactCard,
  faChartPie,
  faBuilding,
  faEllipsis,
  faArchive,
  faBoxOpen,
  faTrashCan,
  faHomeCollection,
  faTable,
  faArrowRight,
  faHammer
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { useParams } from "react-router-dom";

import { getCollectionDetailAPI } from "../../../../../API/NFTCollection";
import FormErrorBox from "../../../../Reusable/FormErrorBox";
import DataDisplayRowText from "../../../../Reusable/DataDisplayRowText";
import DataDisplayRowSelect from "../../../../Reusable/DataDisplayRowSelect";
import AlertBanner from "../../../../Reusable/EveryPage/AlertBanner";
import BubbleLink from "../../../../Reusable/EveryPage/BubbleLink";
import PageLoadingContent from "../../../../Reusable/PageLoadingContent";
import {
  topAlertMessageState,
  topAlertStatusState,
} from "../../../../../AppState";
import { COMMERCIAL_CUSTOMER_TYPE_OF_ID } from "../../../../../Constants/App";
import {
  addCustomerState,
  ADD_CUSTOMER_STATE_DEFAULT,
  currentUserState,
} from "../../../../../AppState";
import {
  Collection_PHONE_TYPE_OF_OPTIONS_WITH_EMPTY_OPTIONS,
  Collection_TYPE_OF_FILTER_OPTIONS,
  Collection_ORGANIZATION_TYPE_OPTIONS,
} from "../../../../../Constants/FieldOptions";
import {
  EXECUTIVE_ROLE_ID,
  MANAGEMENT_ROLE_ID,
} from "../../../../../Constants/App";
import AdminNFTCollectionDetailMoreMobile from "./MobileView";
import AdminNFTCollectionDetailMoreDesktop from "./DesktopView";

function AdminNFTCollectionDetailMore() {
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
  const [currentCollection] = useRecoilState(currentUserState);

  ////
  //// Component states.
  ////

  const [errors, setErrors] = useState({});
  const [isFetching, setFetching] = useState(false);
  const [forceURL, setForceURL] = useState("");
  const [collection, setCollection] = useState({});

  ////
  //// Event handling.
  ////

  //

  ////
  //// API.
  ////

  function onSuccess(response) {
    console.log("onSuccess: Starting...");
    setCollection(response);
  }

  function onError(apiErr) {
    console.log("onError: Starting...");
    setErrors(apiErr);

    // The following code will cause the screen to scroll to the top of
    // the page. Please see ``react-scroll`` for more information:
    // https://github.com/fisshy/react-scroll
    var scroll = Scroll.animateScroll;
    scroll.scrollToTop();
  }

  function onDone() {
    console.log("onDone: Starting...");
    setFetching(false);
  }

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
      getCollectionDetailAPI(id, onSuccess, onError, onDone, onUnauthorized);
    }

    return () => {
      mounted = false;
    };
  }, [id]);

  ////
  //// Component rendering.
  ////

  if (forceURL !== "") {
    return <Navigate to={forceURL} />;
  }

  return (
    <>
      <div className="container">
        <section className="section">
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
                  &nbsp;Detail (More)
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

          {/* Page banner */}
          {collection && collection.status === 100 && (
            <AlertBanner message="Archived" status="info" />
          )}

          {/* Page */}
          <nav className="box">
            {/* Title + Options */}
            {collection && (
              <div className="columns">
                <div className="column">
                  <p className="title is-4">
                    <FontAwesomeIcon className="fas" icon={faCubes} />
                    &nbsp;NFT Collection
                  </p>
                </div>
                <div className="column has-text-right"></div>
              </div>
            )}

            {/* <p className="pb-4">Please fill out all the required fields before submitting this form.</p> */}

            {isFetching ? (
              <PageLoadingContent displayMessage={"Loading..."} />
            ) : (
              <>
                <FormErrorBox errors={errors} />

                {collection && (
                  <div className="container">
                    {/* Tab Navigation */}
                    <div className="tabs is-medium is-size-7-mobile">
                      <ul>
                        <li class="">
                          <Link to={`/admin/collection/${id}`}>Detail</Link>
                        </li>
                        <li>
                          <Link to={`/admin/collection/${id}/nfts`}>
                            NFTs
                          </Link>
                        </li>
                        <li className="is-active">
                          <Link>
                            <strong>
                              More&nbsp;&nbsp;
                              <FontAwesomeIcon
                                className="mdi"
                                icon={faEllipsis}
                              />
                            </strong>
                          </Link>
                        </li>
                      </ul>
                    </div>

                    {collection.smartContractStatus === 1 ? (
                        <>
                            <AdminNFTCollectionDetailMoreDesktop
                              collection={collection}
                            />
                            <AdminNFTCollectionDetailMoreMobile
                              collection={collection}
                              currentCollection={currentCollection}
                            />
                        </>
                    ) : (
                        <section class="hero is-medium has-background-white-ter">
                          <div class="hero-body">
                            <p class="title">
                              <FontAwesomeIcon className="fas" icon={faHammer} />
                              &nbsp;Ready for minting
                            </p>
                            <p class="subtitle">
                              Your NFT collection is now running on the ethereum blockchain and you are ready to start minting NFTs!{" "}
                              <b>
                                <Link to={`/admin/collection/${id}/nfts/add`}>
                                  Click here&nbsp;
                                  <FontAwesomeIcon
                                    className="mdi"
                                    icon={faArrowRight}
                                  />
                                </Link>
                              </b>{" "}
                              to get started creating your first NFT collection.
                            </p>
                          </div>
                        </section>
                    )}

                    {/* Bottom Navigation */}
                    <div className="columns pt-5">
                      <div className="column is-half">
                        <Link
                          className="button is-fullwidth-mobile"
                          to={`/admin/collections`}
                        >
                          <FontAwesomeIcon className="fas" icon={faArrowLeft} />
                          &nbsp;Back to NFT Collections
                        </Link>
                      </div>
                      <div className="column is-half has-text-right"></div>
                    </div>
                  </div>
                )}
              </>
            )}
          </nav>
        </section>
      </div>
    </>
  );
}

export default AdminNFTCollectionDetailMore;
