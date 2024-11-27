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
  faUsers,
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { DateTime } from "luxon";

import FormErrorBox from "../../Reusable/FormErrorBox";
import {
  PAGE_SIZE_OPTIONS,
  ATTACHMENT_STATES,
} from "../../../Constants/FieldOptions";

function RetailerCustomerDetailForAttachmentDesktop(props) {
  const {
    userID,
    listData,
    setPageSize,
    pageSize,
    previousCursors,
    onPreviousClicked,
    onNextClicked,
    onSelectAttachmentForDeletion,
  } = props;
  return (
    <div class="b-table">
      <div class="table-wrapper has-mobile-cards">
        <table class="is-fullwidth is-striped is-hoverable is-fullwidth table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Status</th>
              <th>Created</th>
              <th>File</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {listData &&
              listData.results &&
              listData.results.map(function (attachment, i) {
                return (
                  <tr>
                    <td data-label="Title">{attachment.name}</td>
                    <td data-label="State">
                      {ATTACHMENT_STATES[attachment.status]}
                    </td>
                    <td data-label="Created">{attachment.createdAt}</td>
                    <td data-label="File">
                      <a
                        href={attachment.objectUrl}
                        target="_blank"
                        rel="noreferrer"
                        class=""
                      >
                        <FontAwesomeIcon className="mdi" icon={faDownload} />
                        &nbsp;Download File
                      </a>
                    </td>
                    <td class="is-actions-cell">
                      <div class="buttons is-right">
                        <Link
                          to={`/customer/${userID}/attachment/${attachment.id}`}
                          class="button is-small is-primary"
                          type="button"
                        >
                          View
                        </Link>
                        <Link
                          to={`/customer/${userID}/attachment/${attachment.id}/edit`}
                          class="button is-small is-warning"
                          type="button"
                        >
                          Edit
                        </Link>
                        <button
                          onClick={(e, ses) =>
                            onSelectAttachmentForDeletion(e, attachment)
                          }
                          class="button is-small is-danger"
                          type="button"
                        >
                          <FontAwesomeIcon className="mdi" icon={faTrashCan} />
                          &nbsp;Delete
                        </button>
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

export default RetailerCustomerDetailForAttachmentDesktop;
