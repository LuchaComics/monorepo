export const FINDING_OPTIONS = [
  { value: "pr", label: "Poor" },
  { value: "fr", label: "Fair" },
  { value: "gd", label: "Good" },
  { value: "vg", label: "Very good" },
  { value: "fn", label: "Fine" },
  { value: "vf", label: "Very Fine" },
  { value: "nm", label: "Near Mint" },
];

export const FINDING_WITH_EMPTY_OPTIONS = [
  { value: "", label: "Please select" }, // EMPTY OPTION
  ...FINDING_OPTIONS,
];

export const OVERALL_NUMBER_GRADE_OPTIONS = [
  { value: 0.5, label: "0.5" },
  { value: 1.0, label: "1.0" },
  { value: 1.5, label: "1.5" },
  { value: 2.0, label: "2.0" },
  { value: 2.5, label: "2.5" },
  { value: 3.0, label: "3.0" },
  { value: 3.5, label: "3.5" },
  { value: 4.0, label: "4.0" },
  { value: 4.5, label: "4.5" },
  { value: 5.0, label: "5.0" },
  { value: 5.5, label: "5.5" },
  { value: 6.0, label: "6.0" },
  { value: 6.5, label: "6.5" },
  { value: 7.0, label: "7.0" },
  { value: 7.5, label: "7.5" },
  { value: 8.0, label: "8.0" },
  { value: 8.5, label: "8.5" },
  { value: 9.0, label: "9.0" },
  { value: 9.2, label: "9.2" },
  { value: 9.4, label: "9.4" },
  { value: 9.6, label: "9.6" },
  { value: 9.8, label: "9.8" },
  { value: 10, label: "10" },
];

export const OVERALL_NUMBER_GRADE_WITH_EMPTY_OPTIONS = [
  { value: 0, label: "Please select" }, // EMPTY OPTION
  ...OVERALL_NUMBER_GRADE_OPTIONS,
];

export const PUBLISHER_NAME_OPTIONS = [
  { value: 2, label: "DC Comics" },
  { value: 3, label: "Marvel Comics" },
  { value: 4, label: "Image" },
  { value: 5, label: "Studiocomix Press" },
  { value: 6, label: "Lucha" },
  { value: 7, label: "Boom! Studios" },
  { value: 8, label: "Dark Horse Comics" },
  { value: 9, label: "IDW Publishing" },
  { value: 10, label: "Archie Comics" },
  { value: 11, label: "MUCAHI" },
  { value: 12, label: "Novedades Editores SA de CV" },
  { value: 13, label: "Marvel Comics Mexico" },
  { value: 14, label: "Panini" },
  { value: 15, label: "Aardvark-Vanaheim" },
  { value: 16, label: "Scattered Comics" },
  { value: 17, label: "Xochitl" },
  { value: 18, label: "CCBA" },
  { value: 19, label: "Dell" },
  { value: 20, label: "Goldkey" },
  { value: 21, label: "Charlton Comics" },
  { value: 22, label: "Harvey Comics" },
  { value: 23, label: "Star Comics" },
  { value: 24, label: "Now Comics" },
  { value: 25, label: "Mirage Studios" },
  { value: 26, label: "Eclipse Comics" },
  { value: 27, label: "Malibu Comics" },
  { value: 1, label: "Other (Please specify)" },
];

export const PUBLISHER_NAME_WITH_EMPTY_OPTIONS = [
  { value: 0, label: "Please select" }, // EMPTY OPTION
  ...PUBLISHER_NAME_OPTIONS,
];

export const CPS_PERCENTAGE_GRADE_WITH_EMPTY_OPTIONS = [
  { value: 0, label: "Please select" }, // EMPTY OPTION
  { value: 5, label: "5%" },
  { value: 10, label: "10%" },
  { value: 15, label: "15%" },
  { value: 20, label: "20%" },
  { value: 25, label: "25%" },
  { value: 30, label: "30%" },
  { value: 35, label: "35%" },
  { value: 40, label: "40%" },
  { value: 45, label: "45%" },
  { value: 50, label: "50%" },
  { value: 55, label: "55%" },
  { value: 60, label: "60%" },
  { value: 65, label: "65%" },
  { value: 70, label: "70%" },
  { value: 75, label: "75%" },
  { value: 80, label: "80%" },
  { value: 85, label: "85%" },
  { value: 90, label: "90%" },
  { value: 92, label: "92%" },
  { value: 94, label: "94%" },
  { value: 96, label: "96%" },
  { value: 98, label: "98%" },
  { value: 100, label: "100%" },
];

export const CPS_PERCENTAGE_GRADE_OPTIONS = [
  { value: 5, label: "5%" },
  { value: 10, label: "10%" },
  { value: 15, label: "15%" },
  { value: 20, label: "20%" },
  { value: 25, label: "25%" },
  { value: 30, label: "30%" },
  { value: 35, label: "35%" },
  { value: 40, label: "40%" },
  { value: 45, label: "45%" },
  { value: 50, label: "50%" },
  { value: 55, label: "55%" },
  { value: 60, label: "60%" },
  { value: 65, label: "65%" },
  { value: 70, label: "70%" },
  { value: 75, label: "75%" },
  { value: 80, label: "80%" },
  { value: 85, label: "85%" },
  { value: 90, label: "90%" },
  { value: 92, label: "92%" },
  { value: 94, label: "94%" },
  { value: 96, label: "96%" },
  { value: 98, label: "98%" },
  { value: 100, label: "100%" },
];

export const HOW_DID_YOU_HEAR_ABOUT_US_OPTIONS = [
  { value: 2, label: "My local comic book shop" },
  { value: 3, label: "CPS website" },
  { value: 4, label: "Comic Con booth" },
  { value: 5, label: "Friend" },
  { value: 6, label: "Social media" },
  { value: 7, label: "Blog post article" },
  { value: 1, label: "Other (Please specify)" },
];

export const HOW_DID_YOU_HEAR_ABOUT_US_WITH_EMPTY_OPTIONS = [
  { value: 0, label: "Please select" }, // EMPTY OPTION
  ...HOW_DID_YOU_HEAR_ABOUT_US_OPTIONS,
];

export const HOW_LONG_HAS_YOUR_STORE_BEEN_OPERATING_FOR_OPTIONS = [
  { value: 1, label: "1 year" },
  { value: 2, label: "2-5 years" },
  { value: 3, label: "5-9 years" },
  { value: 4, label: "10+ years" },
];

export const HOW_LONG_HAS_YOUR_STORE_BEEN_OPERATING_FOR_WITH_EMPTY_OPTIONS = [
  { value: 0, label: "Please select" }, // EMPTY OPTION
  ...HOW_LONG_HAS_YOUR_STORE_BEEN_OPERATING_FOR_OPTIONS,
];

export const ISSUE_COVER_YEAR_OPTIONS = [
  { value: 0, label: "Please select" }, // EMPTY OPTION
  { value: 1, label: "No Cover Date Year" },
  { value: 2025, label: "2025" },
  { value: 2024, label: "2024" },
  { value: 2023, label: "2023" },
  { value: 2022, label: "2022" },
  { value: 2021, label: "2021" },
  { value: 2020, label: "2020" },
  { value: 2019, label: "2019" },
  { value: 2018, label: "2018" },
  { value: 2017, label: "2017" },
  { value: 2016, label: "2016" },
  { value: 2015, label: "2015" },
  { value: 2014, label: "2014" },
  { value: 2013, label: "2013" },
  { value: 2012, label: "2012" },
  { value: 2011, label: "2011" },
  { value: 2010, label: "2010" },
  { value: 2009, label: "2009" },
  { value: 2008, label: "2008" },
  { value: 2007, label: "2007" },
  { value: 2006, label: "2006" },
  { value: 2005, label: "2005" },
  { value: 2004, label: "2004" },
  { value: 2003, label: "2003" },
  { value: 2002, label: "2002" },
  { value: 2001, label: "2001" },
  { value: 2000, label: "2000" },
  { value: 1999, label: "1999" },
  { value: 1998, label: "1998" },
  { value: 1997, label: "1997" },
  { value: 1996, label: "1996" },
  { value: 1995, label: "1995" },
  { value: 1994, label: "1994" },
  { value: 1993, label: "1993" },
  { value: 1992, label: "1992" },
  { value: 1991, label: "1991" },
  { value: 1990, label: "1990" },
  { value: 1989, label: "1989" },
  { value: 1988, label: "1988" },
  { value: 1987, label: "1987" },
  { value: 1986, label: "1986" },
  { value: 1985, label: "1985" },
  { value: 1984, label: "1984" },
  { value: 1983, label: "1983" },
  { value: 1982, label: "1982" },
  { value: 1981, label: "1981" },
  { value: 1980, label: "1980" },
  { value: 1979, label: "1979" },
  { value: 1978, label: "1978" },
  { value: 1977, label: "1977" },
  { value: 1976, label: "1976" },
  { value: 1975, label: "1975" },
  { value: 1974, label: "1974" },
  { value: 1973, label: "1973" },
  { value: 1972, label: "1972" },
  { value: 1971, label: "1971" },
  { value: 1970, label: "1970" },
  { value: 1969, label: "1969" },
  { value: 1968, label: "1968" },
  { value: 1967, label: "1967" },
  { value: 1966, label: "1966" },
  { value: 1965, label: "1965" },
  { value: 1964, label: "1964" },
  { value: 1963, label: "1963" },
  { value: 1962, label: "1962" },
  { value: 1961, label: "1961" },
  { value: 1960, label: "1960" },
  { value: 1959, label: "1959" },
  { value: 1958, label: "1958" },
  { value: 1957, label: "1957" },
  { value: 1956, label: "1956" },
  { value: 1955, label: "1955" },
  { value: 1954, label: "1954" },
  { value: 1953, label: "1953" },
  { value: 1952, label: "1952" },
  { value: 1951, label: "1951" },
  { value: 1950, label: "1950" },
  { value: 1949, label: "1949" },
  { value: 1948, label: "1948" },
  { value: 1947, label: "1947" },
  { value: 1946, label: "1946" },
  { value: 1945, label: "1945" },
  { value: 1944, label: "1944" },
  { value: 1943, label: "1943" },
  { value: 1942, label: "1942" },
  { value: 1941, label: "1941" },
  { value: 1940, label: "1940" },
  { value: 1939, label: "1939" },
  { value: 1938, label: "1938" },
  { value: 1937, label: "1937" },
  { value: 1936, label: "1936" },
  { value: 1935, label: "1935" },
  { value: 1934, label: "1934" },
  { value: 1933, label: "1933" },
  { value: 1932, label: "1932" },
  { value: 1931, label: "1931" },
  { value: 1939, label: "1930" },
  { value: 1929, label: "1929" },
  { value: 1928, label: "1928" },
  { value: 1927, label: "1927" },
  { value: 1926, label: "1926" },
  { value: 1925, label: "1925" },
  { value: 1924, label: "1924" },
  { value: 1923, label: "1923" },
  { value: 1922, label: "1922" },
  { value: 1921, label: "1921" },
  { value: 1920, label: "1920" },
  { value: 1919, label: "1919" },
  { value: 1918, label: "1918" },
  { value: 1917, label: "1917" },
  { value: 1916, label: "1916" },
  { value: 1915, label: "1915" },
  { value: 1914, label: "1914" },
  { value: 1913, label: "1913" },
  { value: 1912, label: "1912" },
  { value: 1911, label: "1911" },
  { value: 1910, label: "1910" },
  { value: 1909, label: "1909" },
  { value: 1908, label: "1908" },
  { value: 1907, label: "1907" },
  { value: 1906, label: "1906" },
  { value: 1905, label: "1905" },
  { value: 1904, label: "1904" },
  { value: 1903, label: "1903" },
  { value: 1902, label: "1902" },
  { value: 1901, label: "1901" },
  { value: 1900, label: "1900" },
  { value: 2, label: "1899 or before" },
];

export const ISSUE_COVER_YEAR_WITH_EMPTY_OPTIONS = [
  { value: 0, label: "Please select" }, // EMPTY OPTION
  ...ISSUE_COVER_YEAR_OPTIONS,
];

export const ISSUE_COVER_MONTH_OPTIONS = [
  { value: 1, label: "January" },
  { value: 2, label: "February" },
  { value: 3, label: "March" },
  { value: 4, label: "April" },
  { value: 5, label: "May" },
  { value: 6, label: "June" },
  { value: 7, label: "July" },
  { value: 8, label: "August" },
  { value: 9, label: "September" },
  { value: 10, label: "October" },
  { value: 11, label: "November" },
  { value: 12, label: "December" },
  { value: 13, label: "No Cover Date Month" },
];

export const ISSUE_COVER_MONTH_WITH_EMPTY_OPTIONS = [
  { value: 0, label: "Please select" }, // EMPTY OPTION
  ...ISSUE_COVER_MONTH_OPTIONS,
];

export const USER_STATE_OPTIONS = [
  { value: 1, label: "Pending" },
  { value: 2, label: "Active" },
  { value: 3, label: "Error" },
  { value: 100, label: "Archived" },
];

export const USER_STATE_WITH_EMPTY_OPTIONS = [
  { value: "", label: "Please select" }, // EMPTY OPTION
  ...USER_STATE_OPTIONS,
];

export const USER_ROLES = {
  1: "Root",
  2: "Tenant Owner/Manager",
  3: "Customer",
};

export const USER_ROLE_OPTIONS = [
  { value: 1, label: "Root" },
  { value: 2, label: "Staff" },
  { value: 3, label: "Customer" },
];

export const USER_ROLE_WITH_EMPTY_OPTIONS = [
  { value: "", label: "Please select" }, // EMPTY OPTION
  ...USER_ROLE_OPTIONS,
];

export const USER_STATUS_OPTIONS = [
  { value: 1, label: "Active" },
  { value: 2, label: "Archived" },
];

export const USER_STATUS_WITH_EMPTY_OPTIONS = [
  { value: "", label: "Please select" }, // EMPTY OPTION
  ...USER_STATUS_OPTIONS,
];

export const PAGE_SIZE_OPTIONS = [
  { value: 2, label: "2 Rows" },
  { value: 5, label: "5 Rows" },
  { value: 10, label: "10 Rows" },
  { value: 25, label: "25 Rows" },
  { value: 50, label: "50 Rows" },
  { value: 100, label: "100 Rows" },
  { value: 250, label: "250 Rows" },
];

export const PIN_OBJECT_STATES = {
  "queued": "Queued",
  "pinning": "Pinning",
  "pinned": "Pinned",
  "failed": "Failed",
};

export const SPECIAL_DETAILS_OPTIONS = [
  { value: 3, label: "Direct Edition" },
  { value: 4, label: "Newsstand Edition" },
  // { value: 5, label: 'Variant Cover' },
  { value: 2, label: "Regular Edition" },
  { value: 6, label: "Canadian Price Variant" },
  { value: 7, label: "Facsimile" },
  // { value: 8, label: 'Reprint' },
  { value: 1, label: "Other" },
];

export const SPECIAL_DETAILS_WITH_EMPTY_OPTIONS = [
  { value: "", label: "Please select" }, // EMPTY OPTION
  ...SPECIAL_DETAILS_OPTIONS,
];

export const SERVICE_TYPE_OPTIONS = [
  { value: 1, label: "Pre-Screening Service" },
  { value: 2, label: "CPS Pedigree Service" },
  { value: 3, label: "CPS Capsule" },
  { value: 4, label: "CPS Capsule Indie Mint Gem" },
  { value: 5, label: "CPS Capsule Signature Collection" },
  { value: 6, label: "CPS Capsule U-Grade" },
  { value: 7, label: "CPS Capsule U-Grade Signature Collection" },
];

export const SERVICE_TYPE_WITH_EMPTY_OPTIONS = [
  { value: "", label: "Please select" }, // EMPTY OPTION
  ...SERVICE_TYPE_OPTIONS,
];

export const RETAILER_AVAILABLE_SERVICE_TYPE_OPTIONS = [
  { value: 1, label: "Pre-Screening Service" },
  { value: 2, label: "CPS Pedigree Service" },
  { value: 6, label: "CPS Capsule U-Grade" },
];

export const RETAILER_AVAILABLE_SERVICE_TYPE_WITH_EMPTY_OPTIONS = [
  { value: "", label: "Please select" }, // EMPTY OPTION
  ...RETAILER_AVAILABLE_SERVICE_TYPE_OPTIONS,
];

export const SIGNATURE_BASED_SERVICE_TYPE_OPTIONS = [
  { value: 4, label: "CPS Capsule Indie Mint Gem" },
  { value: 5, label: "CPS Capsule Signature Collection" },
];

export const SIGNATURE_BASED_SERVICE_TYPE_WITH_EMPTY_OPTIONS = [
  { value: "", label: "Please select" }, // EMPTY OPTION
  ...SIGNATURE_BASED_SERVICE_TYPE_OPTIONS,
];

export const SUBMISSION_STATES = {
  1: "Waiting to Receive",
  2: "Received",
  3: "Pending",
  4: "In Process",
  5: "Complete",
  6: "Shipped/Ready for pickup",
  7: "Completed by Retail Partner",
  10: "Error",
  11: "Archived",
};

export const SUBMISSION_STATUS_OPTIONS = [
  { value: 1, label: "Waiting to Receive" },
  { value: 2, label: "Received" },
  { value: 3, label: "Pending" },
  { value: 4, label: "In Process" },
  { value: 5, label: "Complete" },
  { value: 6, label: "Shipped/Ready for pickup" },
  { value: 7, label: "Completed by Retail Partner" },
  { value: 10, label: "Error" },
];

export const SUBMISSION_STATUS_WITH_EMPTY_OPTIONS = [
  { value: "", label: "Please select" }, // EMPTY OPTION
  ...SUBMISSION_STATUS_OPTIONS,
];

////
//// This section is for list filters.
////

export const USER_ROLE_LIST_OPTIONS = [
  { value: 0, label: "All" },
  { value: 1, label: "Root" },
  { value: 2, label: "Staff" },
  { value: 3, label: "Customer" },
];

export const ORGANIZATION_STATUS_LIST_OPTIONS = [
  { value: 0, label: "All" },
  { value: 1, label: "Waiting to Receive" },
  { value: 2, label: "Received" },
  { value: 3, label: "Pending" },
  { value: 4, label: "In Process" },
  { value: 5, label: "Complete" },
  { value: 6, label: "Shipped/Ready for pickup" },
  { value: 10, label: "Error" },
  { value: 11, label: "Archived" },
];

export const SUBMISSION_STATUS_LIST_OPTIONS = [
  { value: 0, label: "All" },
  { value: 1, label: "Waiting to Receive" },
  { value: 2, label: "Received" },
  { value: 3, label: "Pending" },
  { value: 4, label: "In Process" },
  { value: 5, label: "Complete" },
  { value: 6, label: "Shipped/Ready for pickup" },
  { value: 7, label: "Completed by Retail Partner" },
  { value: 10, label: "Error" },
];

export const USER_STATUS_LIST_OPTIONS = [
  { value: 0, label: "All" },
  { value: 1, label: "Active" },
  { value: 100, label: "Archived" },
];

export const SUBMISSION_KEY_ISSUE_OPTIONS = [
  { value: 2, label: "1st appearance of" },
  { value: 3, label: "1st cameo of" },
  { value: 4, label: "1st Team Up of" },
  { value: 5, label: "1st appearance in solo title" },
  { value: 6, label: "Death of" },
  { value: 7, label: "Wedding of" },
  { value: 8, label: "Birth of" },
  { value: 9, label: "Iconic cover by" },
  { value: 10, label: "Origin of" },
  { value: 1, label: "Other" },
];

export const SUBMISSION_KEY_ISSUE_WITH_EMPTY_OPTIONS = [
  { value: "", label: "Please select" }, // EMPTY OPTION
  ...SUBMISSION_KEY_ISSUE_OPTIONS,
];

export const SUBMISSION_PRINTING_OPTIONS = [
  { value: 1, label: "1st Print " },
  { value: 2, label: "2nd Print" },
  { value: 3, label: "3rd print" },
  { value: 4, label: "4th print" },
  { value: 5, label: "5th prin" },
  { value: 6, label: "6th print" },
  { value: 7, label: "7th print" },
  { value: 8, label: "8th print" },
  { value: 9, label: "9th print" },
  { value: 10, label: "10th print" },
  { value: 100, label: "Ashcan" },
  { value: 101, label: "Unknown" },
];

export const SUBMISSION_PRINTING_WITH_EMPTY_OPTIONS = [
  { value: "", label: "Please select" }, // EMPTY OPTION
  ...SUBMISSION_PRINTING_OPTIONS,
];

export const ESTIMATED_SUBMISSION_PER_MONTH_OPTIONS = [
  { value: 1, label: "1-4" },
  { value: 2, label: "5-10" },
  { value: 3, label: "11-20" },
  { value: 4, label: "21-49" },
  { value: 5, label: "50+" },
];

export const ESTIMATED_SUBMISSION_PER_MONTH_WITH_EMPTY_OPTIONS = [
  { value: "", label: "Please select" }, // EMPTY OPTION
  ...ESTIMATED_SUBMISSION_PER_MONTH_OPTIONS,
];

export const USER_SPECIAL_COLLECTION_OPTIONS = [
  // Membership Rank No
  { value: 1, label: "0-0001" },
  { value: 2, label: "1-0001" },
  { value: 3, label: "2-0001" },
  { value: 4, label: "3-0001" },
  { value: 5, label: "4-0001" },
];

export const USER_SPECIAL_COLLECTION_WITH_EMPTY_OPTIONS = [
  { value: "", label: "Please select" }, // EMPTY OPTION
  ...USER_SPECIAL_COLLECTION_OPTIONS,
];

export const ORGANIZATION_LEVEL_OPTIONS = [
  { value: 1, label: "1" },
  { value: 2, label: "2" },
  { value: 3, label: "3" },
];

export const ORGANIZATION_LEVEL_WITH_EMPTY_OPTIONS = [
  { value: "", label: "Please select" }, // EMPTY OPTION
  ...ORGANIZATION_LEVEL_OPTIONS,
];

const OFFER_PAY_FREQUENCY_OPTIONS = [
  { value: 1, label: "One-Time" },
  { value: 2, label: "Day" },
  { value: 3, label: "Week" },
  { value: 4, label: "Month" },
  { value: 5, label: "Year" },
];

export const OFFER_PAY_FREQUENCY_WITH_EMPTY_OPTIONS = [
  { value: "", label: "Please select" }, // EMPTY OPTION
  ...OFFER_PAY_FREQUENCY_OPTIONS,
];

export const BUSINESS_FUNCTION_OPTIONS = [
  { value: 1, label: "Grant Comic Book Submission Services" },
];

export const BUSINESS_FUNCTION_WITH_EMPTY_OPTIONS = [
  { value: 0, label: "Please select" }, // EMPTY OPTION
  ...BUSINESS_FUNCTION_OPTIONS,
];

export const OFFER_STATUS_OPTIONS = [
  // Membership Rank No
  { value: 1, label: "Pending" },
  { value: 2, label: "Active" },
  { value: 3, label: "Archive" },
];

export const OFFER_STATUS_WITH_EMPTY_OPTIONS = [
  { value: "", label: "Please select" }, // EMPTY OPTION
  ...OFFER_STATUS_OPTIONS,
];

export const CREDIT_BUSINESS_FUNCTION_OPTIONS = [
  { value: 1, label: "Grant Free Submission" },
];

export const CREDIT_BUSINESS_FUNCTION_WITH_EMPTY_OPTIONS = [
  { value: "", label: "Please select" }, // EMPTY OPTION
  ...CREDIT_BUSINESS_FUNCTION_OPTIONS,
];

export const CREDIT_STATUS_STATES = {
  1: "Available",
  2: "Consumed",
};

export const CREDIT_BUSINESS_FUNCTION_STATES = {
  1: "Grant Free Submission",
};

export const CREDIT_STATUS_OPTIONS = [
  // Membership Rank No
  { value: 1, label: "Available" },
  { value: 2, label: "Consumed" },
];

export const CREDIT_STATUS_WITH_EMPTY_OPTIONS = [
  { value: "", label: "Please select" }, // EMPTY OPTION
  ...CREDIT_STATUS_OPTIONS,
];

export const NUMBER_OF_CREDITS_OPTIONS = [
  { value: 1, label: "1 Credit" },
  { value: 2, label: "2 Credits" },
  { value: 3, label: "3 Credits" },
  { value: 4, label: "4 Credits" },
  { value: 5, label: "5 Credits" },
  { value: 6, label: "6 Credits" },
  { value: 7, label: "7 Credits" },
  { value: 8, label: "8 Credits" },
  { value: 9, label: "9 Credits" },
  { value: 19, label: "10 Credits" },
];

export const NUMBER_OF_CREDITS_WITH_EMPTY_OPTIONS = [
  { value: "", label: "Please select" }, // EMPTY OPTION
  ...NUMBER_OF_CREDITS_OPTIONS,
];

export const PAYMENT_PROCESSOR_OPTIONS = [
  // Membership Rank No
  { value: 1, label: "Stripe, Inc." },
];

export const PAYMENT_PROCESSOR_WITH_EMPTY_OPTIONS = [
  { value: "", label: "Please select" }, // EMPTY OPTION
  ...PAYMENT_PROCESSOR_OPTIONS,
];

export const NFTS_STATUSES = {
  1: "Active",
  0: "Archived",
};
