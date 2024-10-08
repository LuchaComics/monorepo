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
  faCircleCheck,
  faArchive
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";
import { DateTime } from "luxon";

import FormErrorBox from "../../../Reusable/FormErrorBox";
import { PAGE_SIZE_OPTIONS, USER_ROLES } from "../../../../Constants/FieldOptions";

/*
Display for both tablet and mobile.
*/
function AdminUserListMobile(props) {
  const {
    listData,
    setPageSize,
    pageSize,
    previousCursors,
    onPreviousClicked,
    onNextClicked,
  } = props;
  return (
    <>
      {listData &&
        listData.results &&
        listData.results.map(function (datum, i) {
          return (
            <div class="mb-5" key={`${datum.id}_mobile`}>
              <hr />
              <strong>Name:</strong>&nbsp;{datum.name}
              <br />
              <br />
              <strong>Email:</strong>&nbsp;
              <Link to={`mailto:${datum.email}`}>{datum.email}</Link>
              <br />
              <br />
              <strong>Phone:</strong>&nbsp;
              {datum.phone ? (
                <Link to={`tel:${datum.phone}`}>{datum.phone}</Link>
              ) : (
                <>-</>
              )}
              <br />
              <br />
              <strong>Store:</strong>&nbsp;
              {datum.storeId !== "000000000000000000000000" && (
                <Link
                  to={`/admin/store/${datum.storeId}`}
                  target="_blank"
                  rel="noreferrer"
                  class="is-small"
                >
                  {datum.storeName}&nbsp;
                  <FontAwesomeIcon
                    className="fas"
                    icon={faArrowUpRightFromSquare}
                  />
                </Link>
              )}
              <br />
              <br />
              <strong>Role:</strong>&nbsp;{USER_ROLES[datum.role]}
              <br />
              <br />
              <strong>Joined:</strong>&nbsp;{datum.createdAt}
              <br />
              <br />
              <strong>Status:</strong>&nbsp;{datum.status === 1 ? <><FontAwesomeIcon className="mdi" icon={faCircleCheck} />&nbsp;Active</> : <><FontAwesomeIcon className="mdi" icon={faArchive} />&nbsp;Archived</>}
              <br />
              <br />
              {/* Tablet only */}
              <div class="is-hidden-mobile pt-2">
                <div className="buttons is-right">
                  <Link
                    to={`/admin/user/${datum.id}`}
                    class="button is-small is-primary"
                    type="button"
                  >
                    <FontAwesomeIcon className="mdi" icon={faEye} />
                    &nbsp;View
                  </Link>
                </div>
              </div>
              {/* Mobile only */}
              <div class="is-hidden-tablet pt-2">
                <div class="columns is-mobile">
                  <div class="column">
                    <Link
                      to={`/admin/user/${datum.id}`}
                      class="button is-small is-primary is-fullwidth"
                      type="button"
                    >
                      <FontAwesomeIcon className="mdi" icon={faEye} />
                      &nbsp;View
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

export default AdminUserListMobile;
