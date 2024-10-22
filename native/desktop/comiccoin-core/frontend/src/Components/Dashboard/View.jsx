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



              <nav class="box">
                <div class="columns">
                  <div class="column">
                    <h1 class="title is-4">
                      <FontAwesomeIcon className="fas" icon={faGauge} />
                      &nbsp;Overview
                    </h1>
                  </div>
                </div>

                <nav class="level">
              <div class="level-item has-text-centered">
                <div>
                  <p class="heading">Coins</p>
                  <p class="title">3,456</p>
                </div>
              </div>
              <div class="level-item has-text-centered">
                <div>
                  <p class="heading">Tokens Count</p>
                  <p class="title">123</p>
                </div>
              </div>

            </nav>

            <h1 className="subtitle is-4 pt-5">Recent Transactions</h1>
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
                <div className="has-text-right">
                <Link to={`/transactions`}>See More&nbsp;<FontAwesomeIcon className="fas" icon={faArrowRight} /></Link>
                </div>
              </nav>
            </section>
          </div>
        </>
    )
}

export default DashboardView
