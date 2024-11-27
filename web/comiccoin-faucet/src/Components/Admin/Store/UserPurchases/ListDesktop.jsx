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
import { PAGE_SIZE_OPTIONS } from "../../../../Constants/FieldOptions";

function AdminStorePurchaseListDesktop(props) {
  const {
    listData,
    setPageSize,
    pageSize,
    previousCursors,
    onPreviousClicked,
    onNextClicked,
  } = props;
  return (
    <div className="b-table">
      <div className="table-wrapper has-mobile-cards">
        <table className="is-fullwidth is-striped is-hoverable is-fullwidth table">
          <thead>
            <tr>
              <th>Description</th>
              <th>Subtotal</th>
              <th>Tax</th>
              <th>Total</th>
              <th>Submission</th>
              <th>Purchase Date</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {listData &&
              listData.results &&
              listData.results.map(function (datum, i) {
                return (
                  <tr key={`desktop_${datum.id}`}>
                    <td data-label="Description">{datum.offerName}</td>
                    <td data-label="Subtotal">${datum.amountSubtotal}</td>
                    <td data-label="Tax">${datum.amountTax}</td>
                    <td data-label="Total">${datum.amountTotal}</td>
                    <td data-label="Submission">
                      <Link
                        target="_blank"
                        rel="noreferrer"
                        to={`/admin/submissions/comic/${datum.comicSubmissionId}`}
                      >
                        {datum.comicSubmissionSeriesTitle}&nbsp;
                        {datum.comicSubmissionSeriesIssueVol}&nbsp;
                        {datum.comicSubmissionSeriesIssueNo}&nbsp;
                        <FontAwesomeIcon
                          className="fas"
                          icon={faArrowUpRightFromSquare}
                        />
                      </Link>
                    </td>
                    <td data-label="Purchase Date">
                      {DateTime.fromISO(
                        datum.paymentProcessorPurchasedAt,
                      ).toLocaleString(DateTime.DATETIME_MED)}
                    </td>
                    <td className="is-actions-cell">
                      <div className="buttons is-right">
                        <Link
                          target="_blank"
                          rel="noreferrer"
                          to={datum.paymentProcessorReceiptUrl}
                          className="button is-small is-dark"
                          type="button"
                        >
                          View Receipt&nbsp;
                          <FontAwesomeIcon
                            className="fas"
                            icon={faArrowUpRightFromSquare}
                          />
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

export default AdminStorePurchaseListDesktop;
