import React from "react";

function PageLoadingContent(props) {
  const { displayMessage } = props;
  return (
    <div>
      <div className="loader-wrapper is-centered">
        <div
          className="loader is-loading is-centered"
          style={{ height: "80px", width: "80px", borderColor: "#333" }}
        ></div>
      </div>
      <div className="has-text-centered" style={{ fontSize: "24px", fontWeight: "bold", color: "#333" }}>
        {displayMessage}
      </div>
    </div>
  );
}
export default PageLoadingContent;
