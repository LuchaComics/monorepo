import React, { useState, useEffect } from "react";
import { Link } from "react-router-dom";
import Scroll from "react-scroll";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faChevronRight,
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

import FormErrorBox from "../Reusable/FormErrorBox";
import {
  PAGE_SIZE_OPTIONS,
  CREDIT_STATUS_STATES,
  CREDIT_BUSINESS_FUNCTION_STATES,
} from "../../Constants/FieldOptions";

/*
Display for both tablet and mobile.
*/
function StoreCreditListMobile(props) {
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
    <>
      {listData &&
        listData.results &&
        listData.results.map(function (datum, i) {
          return (
            <div class="mb-5">
              {i !== 0 && <hr />}
              <strong>Business Function:</strong>&nbsp;
              {CREDIT_BUSINESS_FUNCTION_STATES[datum.businessFunction]}
              <br />
              <br />
              <strong>Offer:</strong>&nbsp;{datum.offerName}
              <br />
              <br />
              <strong>Status:</strong>&nbsp;{CREDIT_STATUS_STATES[datum.status]}
              <br />
              <br />
              <strong>Created:</strong>&nbsp;{datum.createdAt}
              <br />
              <br />
              {/* Tablet only */}
              <div class="is-hidden-mobile pt-2">
                <div className="buttons is-right">
                  <Link
                    to={`/admin/user/${datum.userId}/credit/${datum.id}`}
                    class="button is-small is-primary"
                    type="button"
                  >
                    View&nbsp;
                    <FontAwesomeIcon className="fas" icon={faChevronRight} />
                  </Link>
                </div>
              </div>
              {/* Mobile only */}
              <div class="is-hidden-tablet pt-2">
                <div class="columns is-mobile">
                  <div class="column">
                    <Link
                      to={`/admin/user/${datum.userId}/credit/${datum.id}`}
                      class="button is-small is-primary is-fullwidth"
                      type="button"
                    >
                      View&nbsp;
                      <FontAwesomeIcon className="fas" icon={faChevronRight} />
                    </Link>
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

export default StoreCreditListMobile;
