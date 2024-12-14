import React, { useState } from "react";
import {
  Coins,
  Home,
  Image,
  History,
  Wallet,
  Settings,
  HelpCircle,
  LogOut,
  Clock,
  CheckCircle,
  XCircle,
  Menu,
  X,
  ChevronLeft,
  ChevronRight,
} from "lucide-react";

import Topbar from "../../../Components/Navigation/Topbar";

const SubmitComicSuccessPage = () => {
  const [isNavOpen, setIsNavOpen] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedSubmission, setSelectedSubmission] = useState(null);
  const itemsPerPage = 12;

  return (
    <div className="min-h-screen bg-purple-50">
      <Topbar currentPage="Submit Comic" />
      TODO: IMPLEMENT
    </div>
  );
};

export default SubmitComicSuccessPage;
