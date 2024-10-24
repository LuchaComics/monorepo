import React from "react";

function PageLoadingContent(props) {
  const { displayMessage } = props;
  return (
    <div className="columns is-centered is-vcentered is-fullheight" style={{ height: "100vh" }}>
      <div className="column is-narrow" style={{ marginBottom: "15vh" }}>
        <div className="loader is-loading is-centered" style={{ height: "80px", width: "80px", borderColor: "#333" }}></div>
        <div className="has-text-centered" style={{ fontSize: "24px", fontWeight: "bold", color: "#333", marginTop: "20px" }}>
          {displayMessage}
        </div>
      </div>
    </div>
  );
}
export default PageLoadingContent;
