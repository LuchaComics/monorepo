import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
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
  SUBMISSION_STATES,
} from "../../../Constants/FieldOptions";

function RetailerCustomerDetailForComicSubmissionListDesktop(props) {
  const {
    listData,
    setPageSize,
    pageSize,
    previousCursors,
    onPreviousClicked,
    onNextClicked,
    onSelectComicSubmissionForDeletion,
  } = props;
  return (
    <div class="b-table">
      <div class="table-wrapper has-mobile-cards">
        <table class="is-fullwidth is-striped is-hoverable is-fullwidth table">
          <thead>
            <tr>
              <th>Title</th>
              <th>Vol</th>
              <th>No</th>
              <th>Status</th>
              <th>Created</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {listData &&
              listData.results &&
              listData.results.map(function (datum, i) {
                return (
                  <tr>
                    <td data-label="Title">{datum.seriesTitle}</td>
                    <td data-label="Vol">{datum.issueVol}</td>
                    <td data-label="No">{datum.issueNo}</td>
                    <td data-label="State">
                      {SUBMISSION_STATES[datum.status]}
                    </td>
                    <td data-label="Created">{datum.createdAt}</td>
                    <td class="is-actions-cell">
                      <div class="buttons is-right">
                        <Link
                          to={`/submissions/comic/${datum.id}`}
                          target="_blank"
                          rel="noreferrer"
                          class="button is-small is-primary"
                          type="button"
                        >
                          View&nbsp;
                          <FontAwesomeIcon
                            className="fas"
                            icon={faArrowUpRightFromSquare}
                          />
                        </Link>
                        <Link
                          to={`/submissions/comic/${datum.id}/edit`}
                          target="_blank"
                          rel="noreferrer"
                          class="button is-small is-warning"
                          type="button"
                        >
                          Edit&nbsp;
                          <FontAwesomeIcon
                            className="fas"
                            icon={faArrowUpRightFromSquare}
                          />
                        </Link>
                        <button
                          onClick={(e, ses) =>
                            onSelectComicSubmissionForDeletion(e, datum)
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

export default RetailerCustomerDetailForComicSubmissionListDesktop;
