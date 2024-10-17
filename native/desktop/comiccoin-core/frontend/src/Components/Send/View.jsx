import {useState, useEffect} from 'react';
import { Link } from "react-router-dom";


function SendView() {


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

            <div id="result" className="result">Send Page ... <Link to="/receive">Go to receive</Link>
            </div>

        </div>
    )
}

export default SendView
