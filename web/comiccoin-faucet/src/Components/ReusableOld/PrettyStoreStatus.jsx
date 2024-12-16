import React from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faMinusCircle,
  faArchive,
  faBug,
  faBullhorn,
  faCheckCircle,
} from "@fortawesome/free-solid-svg-icons";

function PrettyStoreStatus({ status }) {
  switch (status) {
    case 1:
      return (
        <>
          <FontAwesomeIcon className="mdi" icon={faBullhorn} />
          &nbsp;Pending
        </>
      );
      break;
    case 2:
      return (
        <>
          <FontAwesomeIcon className="mdi" icon={faCheckCircle} />
          &nbsp;Active
        </>
      );
      break;
    case 3:
      return (
        <>
          <FontAwesomeIcon className="mdi" icon={faMinusCircle} />
          &nbsp;Rejected
        </>
      );
      break;
    case 4:
      return (
        <>
          <FontAwesomeIcon className="mdi" icon={faBug} />
          &nbsp;Error
        </>
      );
      break;
    case 5:
      return (
        <>
          <FontAwesomeIcon className="mdi" icon={faArchive} />
          &nbsp;Archived
        </>
      );
      break;
    default:
  }
}

export default PrettyStoreStatus;
