import {useState, useEffect} from 'react';
import { Link } from "react-router-dom";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faTasks,
  faGauge,
  faArrowRight,
  faUsers,
  faBarcode,
  faCubes,
  faInbox
} from "@fortawesome/free-solid-svg-icons";
import { useRecoilState } from "recoil";

import logo from '../../assets/images/CPS-logo-2023-square.webp';
import {QRCodeSVG} from 'qrcode.react';
import { currentOpenWalletAtAddressState } from "../../AppState";


function ReceiveView() {
    ////
    //// Global State
    ////

    const [currentOpenWalletAtAddress] = useRecoilState(currentOpenWalletAtAddressState);

    ////
    //// Component states.
    ////

    ////
    //// Event handling.
    ////

    ////
    //// Misc.
    ////

    useEffect(() => {
      let mounted = true;

      if (mounted) {
            window.scrollTo(0, 0); // Start the page at the top of the page.
      }

      return () => {
        mounted = false;
      };
    }, []);

    ////
    //// Component rendering.
    ////

    return (
        <>
          <div class="container">
            <section class="section">
              <nav class="box">
                <div class="columns">
                  <div class="column">
                    <h1 class="title is-4">
                      <FontAwesomeIcon className="fas" icon={faInbox} />
                      &nbsp;Receive ComicCoins
                    </h1>
                  </div>
                </div>
                <p class="has-text-grey">
                    Senders can scan this QRcode and quickly get a copy of your address.
                </p>
                <p class="has-text-grey">
                  Your current wallet address is:
                  <b> {currentOpenWalletAtAddress}</b>.
                </p>
                  <p>&nbsp;</p>
                <div className="columns is-centered">
                  <div class="column is-half">
                <figure class="image">
                    {/* https://www.npmjs.com/package/qrcode.react */}
                    <QRCodeSVG value={currentOpenWalletAtAddress} size={375} />
                </figure>
                </div>
                </div>

              </nav>
            </section>
          </div>
        </>
    )
}

export default ReceiveView