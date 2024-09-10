import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faDownload,
  faCalendarMinus,
  faCalendarPlus,
  faDumbbell,
  faCalendar,
  faGauge,
  faSearch,
  faEye,
  faPencil,
  faTrashCan,
  faPlus,
  faArrowRight,
  faTable,
  faArrowUpRightFromSquare,
  faFilter,
  faRefresh,
  faCalendarCheck,
  faCollections,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { DateTime } from "luxon";

import FormErrorBox from "../../../../Reusable/FormErrorBox";
import {
  PAGE_SIZE_OPTIONS,
  NFT_METADATA_STATUSES,
} from "../../../../../Constants/FieldOptions";

/*
Display for both tablet and mobile.
*/
function AdminCollectionDetailForNFTMetadataListMobile(props) {
  const {
    collectionID,
    listData,
    setPageSize,
    pageSize,
    previousCursors,
    onPreviousClicked,
    onNextClicked,
    onSelectNFTMetadataForDeletion,
  } = props;
  return (
    <>
      {listData &&
        listData.results &&
        listData.results.map(function (datum, i) {
          return (
            <div class="mb-5">
              {i !== 0 && <hr />}
              <strong>Name:</strong>&nbsp;{(datum.tokenId !== undefined && datum.tokenId !== null && datum.tokenId !== "") ? datum.tokenId : "-"}
              <br />
              <br />
              <strong>Name:</strong>&nbsp;{datum.name ? datum.name : "-"}
              <br />
              <br />
              <strong>Status:</strong>&nbsp;{NFT_METADATA_STATUSES[datum.status]}
              <br />
              <br />
              <strong>Created At:</strong>&nbsp;{datum.createdAt}
              <br />
              <br />
              <strong>Modified At:</strong>&nbsp;{datum.modifiedAt}
              <br />
              <br />
              {/* Tablet only */}
              <div class="is-hidden-mobile pt-2">
                <div className="buttons is-right">
                  <Link
                    to={`/admin/collection/${collectionID}/nft-metadatum/${datum.id}`}
                    class="button is-small is-primary"
                    type="button"
                  >
                    View
                  </Link>
                  <Link
                    to={`/admin/collection/${collectionID}/nft-metadatum/${datum.id}/edit`}
                    class="button is-small is-warning"
                    type="button"
                  >
                    Edit
                  </Link>
                  <button
                    onClick={(e, ses) =>
                      onSelectNFTMetadataForDeletion(e, datum)
                    }
                    class="button is-small is-danger"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faTrashCan} />
                    &nbsp;Delete
                  </button>
                </div>
              </div>
              {/* Mobile only */}
              <div class="is-hidden-tablet pt-2">
                <div class="columns is-mobile">
                  <div class="column">
                    <Link
                      to={`/admin/collection/${collectionID}/nft-metadatum/${datum.id}`}
                      class="button is-small is-primary is-fullwidth"
                      type="button"
                    >
                      View
                    </Link>
                  </div>
                  <div class="column">
                    <Link
                      to={`/admin/collection/${collectionID}/nft-metadatum/${datum.id}/edit`}
                      class="button is-small is-warning is-fullwidth"
                      type="button"
                    >
                      Edit
                    </Link>
                  </div>
                  <div class="column">
                    <button
                      onClick={(e, ses) =>
                        onSelectNFTMetadataForDeletion(e, datum)
                      }
                      class="button is-small is-danger is-fullwidth"
                      type="button"
                    >
                      <FontAwesomeIcon className="mdi" icon={faTrashCan} />
                      &nbsp;Delete
                    </button>
                  </div>
                </div>
              </div>
            </div>
          );
        })}

      <div class="columns is-mobile pt-4">
        <div class="column is-half">
          <span class="select">
            <select
              class={`input has-text-grey-light`}
              name="pageSize"
              onChange={(e) => setPageSize(parseInt(e.target.value))}
            >
              {PAGE_SIZE_OPTIONS.map(function (option, i) {
                return (
                  <option
                    selected={pageSize === option.value}
                    value={option.value}
                  >
                    {option.label}
                  </option>
                );
              })}
            </select>
          </span>
        </div>
        <div class="column is-half has-text-right">
          {previousCursors.length > 0 && (
            <button class="button" onClick={onPreviousClicked}>
              Previous
            </button>
          )}
          {listData.hasNextPage && (
            <>
              <button class="button" onClick={onNextClicked}>
                Next
              </button>
            </>
          )}
        </div>
      </div>
    </>
  );
}

export default AdminCollectionDetailForNFTMetadataListMobile;
