import {useState, useEffect} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet, GetPageID} from "../wailsjs/go/main/App";
import { Link } from "react-router-dom";

function App() {
    const [resultText, setResultText] = useState("Please enter your name below 👇");
    const [name, setName] = useState('');
    const updateName = (e) => setName(e.target.value);
    const updateResultText = (result) => setResultText(result);


    const [pageID, setPageID] = useState("PageID");
    const updatePageID = (result) => setPageID(result);

    function greet() {
        Greet(name).then(updateResultText);
        GetPageID().then(updatePageID);
    }

    useEffect(() => {
      let mounted = true;

      if (mounted) {
            window.scrollTo(0, 0); // Start the page at the top of the page.

            GetPageID().then(updatePageID);
      }


      return () => {
        mounted = false;
      };
    }, []);

    console.log("---> ", pageID);


    return (
        <div id="App">

            <div id="result" className="result">You are ready <Link to="/welcome">tesT</Link>
            </div>

        </div>
    )
}

export default App
