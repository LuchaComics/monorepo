import {useState, useEffect} from 'react';
import { Link } from "react-router-dom";


function TransactionsView() {


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

            <div id="result" className="result">Transactions Page ... <Link to="/">Go to Startup</Link>
            </div>

        </div>
    )
}

export default TransactionsView
