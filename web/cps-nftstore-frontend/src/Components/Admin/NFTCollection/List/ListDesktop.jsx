import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
faCircleCheck,
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
  faArchive
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { DateTime } from "luxon";

import FormErrorBox from "../../../Reusable/FormErrorBox";
import { PAGE_SIZE_OPTIONS, USER_ROLES } from "../../../../Constants/FieldOptions";

function AdminNFTCollectionListDesktop(props) {
  const {
    listData,
    setPageSize,
    pageSize,
    previousCursors,
    onPreviousClicked,
    onNextClicked,
  } = props;
  return (
    <div class="b-table">
      <div class="table-wrapper has-mobile-cards">
        <table class="is-fullwidth is-striped is-hoverable is-fullwidth is-size-7-desktop is-size-6-widescreen table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Tokens</th>
              <th>Created</th>
              <th>Status</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {listData &&
              listData.results &&
              listData.results.map(function (collection, i) {
                return (
                  <tr key={`${collection.id}_desktop`}>
                    <td data-label="Name">{collection.name}</td>
                    <td data-label="Tokens Count">{collection.tokensCount}</td>
                    <td data-label="Created">{collection.createdAt}</td>
                    <td data-label="Status" className="has-text-centered">
                        {collection.status === 1 ? <FontAwesomeIcon className="mdi" icon={faCircleCheck} /> : <FontAwesomeIcon className="mdi" icon={faArchive} />}
                    </td>
                    <td class="is-actions-cell">
                      <div class="buttons is-right">
                        <Link
                          to={`/admin/collection/${collection.id}`}
                          class="button is-small is-primary"
                          type="button"
                        >
                          <FontAwesomeIcon className="mdi" icon={faEye} />
                          &nbsp;View
                        </Link>
                      </div>
                    </td>
                  </tr>
                );
              })}
          </tbody>
        </table>

        <div class="columns">
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
      </div>
    </div>
  );
}

export default AdminNFTCollectionListDesktop;
