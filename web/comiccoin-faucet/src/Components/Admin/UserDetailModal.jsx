import React from 'react';
import { X } from 'lucide-react';

const UserDetailModal = ({ user, onClose, onAccept, onReject }) => {
  if (!user) return null;

  const renderField = (label, value) => {
    if (value === undefined || value === null || value === '') return null;
    return (
      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700">{label}</label>
        <div className="mt-1 text-sm text-gray-900">{value.toString()}</div>
      </div>
    );
  };

  const renderBooleanField = (label, value) => {
    if (value === undefined || value === null) return null;
    return renderField(label, value ? 'Yes' : 'No');
  };

  const renderNumericField = (label, value) => {
    if (value === undefined || value === null) return null;
    return renderField(label, value.toString());
  };

  const getHowLongLabel = (value) => {
    const options = {
      1: 'Less than 1 year',
      2: '1-2 years',
      3: '3-5 years',
      4: '5-10 years',
      5: 'More than 10 years'
    };
    return options[value] || value.toString();
  };

  const getYesNoLabel = (value) => {
    const options = {
      1: 'Yes',
      2: 'No',
      3: 'Not Sure'
    };
    return options[value] || value.toString();
  };

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div className="relative top-20 mx-auto p-8 border w-full max-w-4xl shadow-lg rounded-2xl bg-white">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold text-purple-800">User Profile Review</h2>
          <button onClick={onClose} className="text-gray-500 hover:text-gray-700">
            <X className="w-6 h-6" />
          </button>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 max-h-[60vh] overflow-y-auto p-4">
          <div>
            <h3 className="text-lg font-semibold text-purple-700 mb-4">Personal Information</h3>
            {renderField('Email', user.email)}
            {renderField('First Name', user.firstName)}
            {renderField('Last Name', user.lastName)}
            {renderField('Phone', user.phone)}
            {renderField('Country', user.country)}
            {renderField('Region', user.region)}
            {renderField('City', user.city)}
            {renderField('Postal Code', user.postalCode)}
            {renderField('Address Line 1', user.addressLine1)}
            {renderField('Address Line 2', user.addressLine2)}

            <h3 className="text-lg font-semibold text-purple-700 mt-6 mb-4">Comic Collecting Experience</h3>
            {renderField('How Long Collecting', getHowLongLabel(user.howLongCollectingComicBooksForGrading))}
            {renderField('Previously Submitted for Grading', getYesNoLabel(user.hasPreviouslySubmittedComicBookForGrading))}
            {renderField('Owned Graded Comics', getYesNoLabel(user.hasOwnedGradedComicBooks))}
            {renderField('Has Regular Comic Shop', getYesNoLabel(user.hasRegularComicBookShop))}
          </div>

          <div>
            <h3 className="text-lg font-semibold text-purple-700 mb-4">Shipping Information</h3>
            {renderBooleanField('Has Shipping Address', user.hasShippingAddress)}
            {user.hasShippingAddress && (
              <>
                {renderField('Shipping Name', user.shippingName)}
                {renderField('Shipping Phone', user.shippingPhone)}
                {renderField('Shipping Country', user.shippingCountry)}
                {renderField('Shipping Region', user.shippingRegion)}
                {renderField('Shipping City', user.shippingCity)}
                {renderField('Shipping Postal Code', user.shippingPostalCode)}
                {renderField('Shipping Address Line 1', user.shippingAddressLine1)}
                {renderField('Shipping Address Line 2', user.shippingAddressLine2)}
              </>
            )}

            <h3 className="text-lg font-semibold text-purple-700 mt-6 mb-4">Additional Information</h3>
            {renderNumericField('How Did You Hear About Us', user.howDidYouHearAboutUs)}
            {user.howDidYouHearAboutUs === 'other' && renderField('Other Source', user.howDidYouHearAboutUsOther)}
            {renderBooleanField('Agreed to Terms of Service', user.agreeTermsOfService)}
            {renderBooleanField('Agreed to Promotions', user.agreePromotions)}
            {renderField('Purchased from Auction Site', getYesNoLabel(user.hasPreviouslyPurchasedFromAuctionSite))}
            {renderField('Purchased from Facebook Marketplace', getYesNoLabel(user.hasPreviouslyPurchasedFromFacebookMarketplace))}
            {renderField('Attends Comic Cons', getYesNoLabel(user.hasRegularlyAttendedComicConsOrCollectibleShows))}
          </div>
        </div>

        <div className="flex justify-end space-x-4 mt-6 pt-4 border-t">
          <button
            onClick={onClose}
            className="px-4 py-2 text-gray-600 hover:text-gray-800 font-medium rounded-lg border border-gray-300 hover:bg-gray-50"
          >
            Close
          </button>
          <button
            onClick={onReject}
            className="px-4 py-2 bg-red-500 text-white font-medium rounded-lg hover:bg-red-600"
          >
            Reject
          </button>
          <button
            onClick={onAccept}
            className="px-4 py-2 bg-green-500 text-white font-medium rounded-lg hover:bg-green-600"
          >
            Accept
          </button>
        </div>
      </div>
    </div>
  );
};

export default UserDetailModal;
