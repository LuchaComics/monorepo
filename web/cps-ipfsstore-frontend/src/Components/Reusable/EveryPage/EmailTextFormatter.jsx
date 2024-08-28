import React from "react";
import { Link } from "react-router-dom";


function EmailTextFormatter(props) {
    const { value } = props;
    if (value.includes("customer_") || value.includes("@cpsdata.ca")) {
        return <>-</>
    }
    return (
        <Link to={`mailto:${value}`}>{value}</Link>
    );
}

export default EmailTextFormatter;
