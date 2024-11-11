import {useState, useEffect} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {Greet, GetPageID} from "../wailsjs/go/main/App";

function Welcome() {
    const [resultText, setResultText] = useState("Please enter your name WelcomeWelcomeWelcomeWelcome below 👇");
    const [name, setName] = useState('');
    const [pageID, setPageID] = useState("PageID");
    const updateName = (e) => setName(e.target.value);
    const updateResultText = (result) => setResultText(result);
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
            <img src={logo} id="logo" alt="logo"/>
            <div id="result" className="result">{resultText}</div>
            <div id="input" className="input-box">
                <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text"/>
                <button className="btn" onClick={greet}>Greet!!!</button>
            </div>
        </div>
    )
}

export default Welcome
