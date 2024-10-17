import {useState, useEffect} from 'react';
import { Link } from "react-router-dom";


function DashboardView() {


    useEffect(() => {
      let mounted = true;

      if (mounted) {
            window.scrollTo(0, 0); // Start the page at the top of the page.
      }


      return () => {
        mounted = false;
      };
    }, []);



    return (
        <div id="App">

            <div id="result" className="result">App dashboard... <Link to="/send">Go to send</Link>
            </div>

        </div>
    )
}

export default DashboardView
