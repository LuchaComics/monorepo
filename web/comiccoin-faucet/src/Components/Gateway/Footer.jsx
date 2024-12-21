import React from "react";
import { Link } from "react-router-dom";
import { Coins, BookOpen, Camera, Gift, Github, ArrowRight } from 'lucide-react';

function Footer() {
  return (
        <footer className="bg-gradient-to-r from-purple-700 to-indigo-800 text-white py-8">
          <div className="max-w-7xl mx-auto px-4">
            <div className="grid md:grid-cols-3 gap-8 mb-8">
              {/* Project Links */}
              <div className="text-center md:text-left">
                <h3 className="font-bold mb-4">Project Resources</h3>
                <ul className="space-y-2">
                  <li>
                    <a href="https://github.com/LuchaComics/monorepo" className="hover:text-purple-200 flex items-center justify-center md:justify-start gap-2">
                      <Github className="h-4 w-4" />
                      <span>GitHub Repository</span>
                    </a>
                  </li>
                  <li>
                    <a href="https://comiccoin.ca" className="hover:text-purple-200">
                      Project Website
                    </a>
                  </li>
                  {/*
                  <li>
                    <a href="https://groups.google.com/g/comic-coin-community" className="hover:text-purple-200">
                      Join Discussion Group
                    </a>
                  </li>
                  */}
                </ul>
              </div>
              {/* Account Links */}
              <div className="text-center md:text-left">
                <h3 className="font-bold mb-4">Account</h3>
                <ul className="space-y-2">
                  <li>
                    <Link to="/register" className="hover:text-purple-200">
                      Register
                    </Link>
                  </li>
                  <li>
                    <Link to="/login" className="hover:text-purple-200">
                      Login
                    </Link>
                  </li>
                </ul>
              </div>
              {/* Legal Links */}
              <div className="text-center md:text-left">
                <h3 className="font-bold mb-4">Legal</h3>
                <ul className="space-y-2">
                  <li>
                    <a href="/terms" className="hover:text-purple-200">
                      Terms of Service
                    </a>
                  </li>
                  <li>
                    <a href="/privacy" className="hover:text-purple-200">
                      Privacy Policy
                    </a>
                  </li>
                </ul>
              </div>
            </div>
            <div className="text-center pt-4 border-t border-purple-500">
              <p>© 2024 ComicCoin Faucet. All rights reserved.</p>
            </div>
          </div>
        </footer>
  );
}

export default Footer;
