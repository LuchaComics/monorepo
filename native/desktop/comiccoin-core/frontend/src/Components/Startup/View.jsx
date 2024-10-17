import {useState, useEffect} from 'react';
import { Link } from "react-router-dom";


function StartupView() {


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

            <div id="result" className="result">Startup page ... <Link to="/pick-storage-location-on-startup">Go to page</Link>
            </div>

        </div>
    )
}

export default StartupView
