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

import FormErrorBox from "../../../Reusable/FormErrorBox";
import {
  PAGE_SIZE_OPTIONS,
  SUBMISSION_STATES,
} from "../../../../Constants/FieldOptions";

function RetailerComicSubmissionListDesktop(props) {
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
        <table class="is-fullwidth is-striped is-hoverable is-fullwidth table">
          <thead>
            <tr>
              <th>CPSR #</th>
              <th>Title</th>
              <th>Vol</th>
              <th>No</th>
              <th>Status</th>
              <th>Customer</th>
              <th>Created</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {listData &&
              listData.results &&
              listData.results.map(function (submission, i) {
                return (
                  <tr>
                    <td data-label="CPSR #">{submission.cpsrn}</td>
                    <td data-label="Title">{submission.seriesTitle}</td>
                    <td data-label="Vol">{submission.issueVol}</td>
                    <td data-label="No">{submission.issueNo}</td>
                    <td data-label="State">
                      {SUBMISSION_STATES[submission.status]}
                    </td>
                    <td data-label="Customer">
                      {submission.customerId !== undefined &&
                      submission.customerId !== null &&
                      submission.customerId !== "" ? (
                        <>
                          {submission.customerId !==
                            "000000000000000000000000" && (
                            <Link
                              to={`/customer/${submission.customerId}`}
                              target="_blank"
                              rel="noreferrer"
                              class="is-small"
                            >
                              {submission.customerFirstName}&nbsp;
                              {submission.customerLastName}&nbsp;
                              <FontAwesomeIcon
                                className="fas"
                                icon={faArrowUpRightFromSquare}
                              />
                            </Link>
                          )}
                        </>
                      ) : (
                        <>-</>
                      )}
                    </td>
                    <td data-label="Created">{submission.createdAt}</td>
                    <td class="is-actions-cell">
                      <div class="buttons is-right">
                        <Link
                          to={`/submissions/comic/${submission.id}`}
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

export default RetailerComicSubmissionListDesktop;
