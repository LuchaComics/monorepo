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
  ATTACHMENT_STATES,
} from "../../../Constants/FieldOptions";

/*
Display for both tablet and mobile.
*/
function AdminStoreDetailForAttachmentListMobile(props) {
  const {
    submissionID,
    listData,
    setPageSize,
    pageSize,
    previousCursors,
    onPreviousClicked,
    onNextClicked,
    onSelectAttachmentForDeletion,
  } = props;
  return (
    <>
      {listData &&
        listData.results &&
        listData.results.map(function (datum, i) {
          return (
            <div class="mb-5">
              {i !== 0 && <hr />}
              <strong>Name:</strong>&nbsp;{datum.name}
              <br />
              <br />
              <strong>Status:</strong>&nbsp;{ATTACHMENT_STATES[datum.status]}
              <br />
              <br />
              <strong>Created:</strong>&nbsp;{datum.createdAt}
              <br />
              <br />
              {/* Tablet only */}
              <div class="is-hidden-mobile pt-2">
                <div class="buttons is-right">
                  <Link
                    to={`/submissions/comic/${submissionID}/attachment/${datum.id}`}
                    class="button is-small is-primary"
                    type="button"
                  >
                    View
                  </Link>
                  <Link
                    to={`/submissions/comic/${submissionID}/attachment/${datum.id}/edit`}
                    class="button is-small is-warning"
                    type="button"
                  >
                    Edit
                  </Link>
                  <button
                    onClick={(e, ses) =>
                      onSelectAttachmentForDeletion(e, datum)
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
                      to={`/submissions/comic/${submissionID}/attachment/${datum.id}`}
                      class="button is-small is-primary is-fullwidth"
                      type="button"
                    >
                      View
                    </Link>
                  </div>
                  <div class="column">
                    <Link
                      to={`/submissions/comic/${submissionID}/attachment/${datum.id}/edit`}
                      class="button is-small is-warning is-fullwidth"
                      type="button"
                    >
                      Edit
                    </Link>
                  </div>
                  <div class="column">
                    <button
                      onClick={(e, ses) =>
                        onSelectAttachmentForDeletion(e, datum)
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

export default AdminStoreDetailForAttachmentListMobile;
