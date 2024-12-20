import React, { useState } from 'react';
import {
  Coins, Home, Image, History, Wallet,
  Settings, HelpCircle, LogOut, Plus,
  MinusCircle, Menu, X, BookOpen,
  MessageCircle, Shield, AlertCircle
} from 'lucide-react';

import Topbar from "../../../Components/Navigation/Topbar";


const HelpPage = () => {
  const [isNavOpen, setIsNavOpen] = useState(false);
  const [openSection, setOpenSection] = useState('getting-started');

  const navigation = [
    { name: 'Dashboard', icon: Home, current: false },
    { name: 'Submit Comic', icon: Image, current: false },
    { name: 'My Submissions', icon: History, current: false },
    { name: 'My Wallet', icon: Wallet, current: false },
    { name: 'Help', icon: HelpCircle, current: true },
    { name: 'Settings', icon: Settings, current: false },
  ];

  const helpSections = [
    {
      id: 'getting-started',
      title: 'Getting Started',
      icon: BookOpen,
      content: [
        {
          question: 'What is ComicCoin Faucet?',
          answer: 'ComicCoin Faucet is a platform where comic book collectors can submit their comic book information to receive free ComicCoins. These coins can be used within our ecosystem for various comic-related activities and trades.'
        },
        {
          question: 'How do I start earning ComicCoins?',
          answer: 'Getting started is easy! Simply register your account, verify your email, and start submitting your comic book information. Each verified submission earns you ComicCoins.'
        },
        {
          question: 'What information do I need to submit?',
          answer: 'For each submission, you\'ll need to provide: The comic book title, issue number, cover photos (front and back), publication date, and condition details. The more accurate your information, the faster your submission will be processed.'
        }
      ]
    },
    {
      id: 'submission-guidelines',
      title: 'Submission Guidelines',
      icon: Image,
      content: [
        {
          question: 'What makes a good submission?',
          answer: 'A good submission includes clear, well-lit photos of your comic book, accurate details about its condition, and complete publication information. Make sure your photos show all corners and edges clearly.'
        },
        {
          question: 'How long does verification take?',
          answer: 'Most submissions are verified within 24 hours. Complex cases or high-value comics may take up to 48 hours for thorough verification.'
        },
        {
          question: 'Why was my submission rejected?',
          answer: 'Submissions may be rejected if the photos are unclear, information is incomplete, or if it\'s a duplicate submission. You\'ll receive specific feedback for each rejection.'
        }
      ]
    },
    {
      id: 'rewards',
      title: 'Rewards & ComicCoins',
      icon: Coins,
      content: [
        {
          question: 'How many ComicCoins can I earn?',
          answer: 'The number of ComicCoins awarded depends on the rarity and condition of your comic. Common issues typically earn 25-50 coins, while rare or valuable comics can earn significantly more.'
        },
        {
          question: 'When do I receive my ComicCoins?',
          answer: 'ComicCoins are credited to your wallet immediately after your submission is approved. You\'ll receive a notification when the coins are added to your balance.'
        },
        {
          question: 'What can I do with ComicCoins?',
          answer: 'ComicCoins can be used to access exclusive content, participate in community events, get priority verification, and trade with other collectors in our marketplace (coming soon).'
        }
      ]
    },
    {
      id: 'support',
      title: 'Support & Community',
      icon: MessageCircle,
      content: [
        {
          question: 'How do I get help?',
          answer: 'Our support team is available 24/7 through the contact form. For faster responses, check our FAQ or join our community Discord server where experienced members can help.'
        },
        {
          question: 'Is my information secure?',
          answer: 'Yes! We use industry-standard encryption to protect your data. Your personal information and comic details are never shared without your permission.'
        }
      ]
    }
  ];

  return (
    <div className="min-h-screen bg-purple-50">
      {/* Navigation */}
      <Topbar currentPage="Help" />

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Help Header */}
        <div className="flex items-center mb-8">
          <HelpCircle className="h-8 w-8 text-purple-600 mr-3" />
          <h1 className="text-3xl font-bold text-purple-800" style={{fontFamily: 'Comic Sans MS, cursive'}}>
            Help Center
          </h1>
        </div>

        {/* Search Bar */}
        <div className="bg-white rounded-xl shadow-lg p-6 mb-8 border-2 border-purple-200">
          <h2 className="text-xl text-purple-800 mb-4">How can we help you today?</h2>
          <input
            type="text"
            placeholder="Search for help topics..."
            className="w-full p-3 border-2 border-purple-200 rounded-lg focus:outline-none focus:border-purple-400"
          />
        </div>

        {/* Help Sections - Now in a single column */}
        <div className="space-y-6">
          {helpSections.map((section) => (
            <div key={section.id} className="bg-white rounded-xl shadow-lg border-2 border-purple-200 overflow-hidden">
              <button
                onClick={() => setOpenSection(openSection === section.id ? '' : section.id)}
                className="w-full p-6 flex items-center justify-between bg-purple-50 hover:bg-purple-100 transition-colors"
              >
                <div className="flex items-center space-x-3">
                  <section.icon className="h-6 w-6 text-purple-600" />
                  <h3 className="text-lg font-semibold text-purple-800">{section.title}</h3>
                </div>
                {openSection === section.id ? (
                  <MinusCircle className="h-5 w-5 text-purple-600" />
                ) : (
                  <Plus className="h-5 w-5 text-purple-600" />
                )}
              </button>

              {openSection === section.id && (
                <div className="p-6 space-y-6">
                  {section.content.map((item, index) => (
                    <div key={index} className="space-y-2">
                      <h4 className="font-medium text-purple-800">{item.question}</h4>
                      <p className="text-gray-600 text-sm">{item.answer}</p>
                    </div>
                  ))}
                </div>
              )}
            </div>
          ))}
        </div>

        {/* Need More Help Section */}
        <div className="mt-8 bg-white rounded-xl shadow-lg p-6 border-2 border-purple-200">
          <div className="flex items-center mb-4">
            <AlertCircle className="h-6 w-6 text-purple-600 mr-2" />
            <h2 className="text-xl font-semibold text-purple-800">Still Need Help?</h2>
          </div>
          <p className="text-gray-600 mb-4">
            Can't find what you're looking for? Our support team is here to help you 24/7.
          </p>
          <button className="px-4 py-2 bg-purple-600 text-white rounded-lg hover:bg-purple-700 transition-colors">
            Contact Support
          </button>
        </div>
      </main>
    </div>
  );
};

export default HelpPage;
