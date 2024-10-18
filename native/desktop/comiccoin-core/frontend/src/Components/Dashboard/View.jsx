import {useState, useEffect} from 'react';
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faGauge,
  faArrowRight,
  faUsers,
  faBarcode,
  faCubes
} from "@fortawesome/free-solid-svg-icons";


import logo from '../../assets/images/CPS-logo-2023-square.webp';


function DashboardView() {

    const totalCoins = 1000;
  const recentTransactions = [
    {
      id: 1,
      date: '2023-02-20',
      type: 'Received',
      amount: 100,
      sender: '0x1234567890',
      receiver: '0x9876543210',
    },
    {
      id: 2,
      date: '2023-02-19',
      type: 'Sent',
      amount: 200,
      sender: '0x9876543210',
      receiver: '0x1234567890',
    },
    {
      id: 3,
      date: '2023-02-18',
      type: 'Received',
      amount: 50,
      sender: '0x1234567890',
      receiver: '0x9876543210',
    },
  ];

  const nonFungibleTokens = [
    {
      id: 1,
      name: 'Token 1',
      description: 'This is a rare token',
      image: logo,
    },
    {
      id: 2,
      name: 'Token 2',
      description: 'This is a common token',
      image:logo,
    },
    {
      id: 3,
      name: 'Token 3',
      description: 'This is a unique token',
      image: logo,
    },
  ];


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
        <>
          <div class="container">
            <section class="section">

              <nav class="breadcrumb" aria-label="breadcrumbs">
                <ul>
                  <li class="is-active">
                    <Link to="/dashboard" aria-current="page">
                      <FontAwesomeIcon className="fas" icon={faGauge} />
                      &nbsp;Overview
                    </Link>
                  </li>
                </ul>
              </nav>

              <nav class="box">
                <div class="columns">
                  <div class="column">
                    <h1 class="title is-4">
                      <FontAwesomeIcon className="fas" icon={faGauge} />
                      &nbsp;Overview
                    </h1>
                  </div>
                </div>

                <div className="columns">
          <div className="column is-8">
            <div className="box">
              <h1 className="title">Summary</h1>
              <div className="columns">
                <div className="column">
                  <h2 className="subtitle">Total Coins</h2>
                  <p className="is-size-1">{totalCoins}</p>
                </div>
              </div>
              <h2 className="subtitle">Recent Transactions</h2>
              <table className="table is-fullwidth">
                <thead>
                  <tr>
                    <th>Date</th>
                    <th>Type</th>
                    <th>Amount</th>
                    <th>Sender</th>
                    <th>Receiver</th>
                  </tr>
                </thead>
                <tbody>
                  {recentTransactions.map((transaction) => (
                    <tr key={transaction.id}>
                      <td>{transaction.date}</td>
                      <td>{transaction.type}</td>
                      <td>{transaction.amount}</td>
                      <td>{transaction.sender}</td>
                      <td>{transaction.receiver}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
      <div className="column is-4">
        <div className="box">
          <h1 className="title">Non-Fungible Tokens</h1>
          <ul>
            {nonFungibleTokens.map((token) => (
              <li key={token.id}>
                <div className="media">
                  <div className="media-left">
                    <img src={token.image} alt={token.name} class="image is-64x64" />
                  </div>
                  <div className="media-content">
                    <h2>{token.name}</h2>
                    <p>{token.description}</p>
                  </div>
                </div>
              </li>
            ))}
          </ul>
        </div>
      </div>

    </div>





              </nav>
            </section>
          </div>
        </>
    )
}

export default DashboardView
