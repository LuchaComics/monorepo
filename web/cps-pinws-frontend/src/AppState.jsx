import {
  RecoilRoot,
  atom,
  selector,
  useRecoilState,
  useRecoilValue,
  AtomEffect,
} from "recoil";
import { recoilPersist } from "recoil-persist";

// Control whether the hamburer menu icon was clicked or not. This state is
// needed by 'TopNavigation' an 'SideNavigation' components.
export const onHamburgerClickedState = atom({
  key: "onHamburgerClicked", // unique ID (with respect to other atoms/selectors)
  default: false, // default value (aka initial value)
});

// Control what message to display at the top as a banner in the app.
export const topAlertMessageState = atom({
  key: "topBannerAlertMessage",
  default: "",
});

// Control what type of message to display at the top as a banner in the app.
export const topAlertStatusState = atom({
  key: "topBannerAlertStatus",
  default: "success",
});

// https://github.com/polemius/recoil-persist
const { persistAtom } = recoilPersist();

export const currentUserState = atom({
  key: "currentUser",
  default: null,
  effects_UNSTABLE: [persistAtom],
});

export const currentOTPResponseState = atom({
  key: "currentOTPResponse",
  default: null,
  effects_UNSTABLE: [persistAtom],
});

// --- Customers List --- //

// Control whether to show filters for the list.
export const customerFilterShowState = atom({
  key: "customerFilterShowState",
  default: false,
  effects_UNSTABLE: [persistAtom],
});

export const customerFilterTemporarySearchTextState = atom({
  key: "customerFilterTemporarySearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const customerFilterActualSearchTextState = atom({
  key: "customerFilterActualSearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const customerFilterStatusState = atom({
  key: "customerFilterStatusState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});
export const customerFilterJoinedAfterState = atom({
  key: "customerFilterJoinedAfterState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

// --- Submission List --- //

// Control whether to show filters for the list.
export const submissionFilterShowState = atom({
  key: "submissionFilterShowState",
  default: false,
  effects_UNSTABLE: [persistAtom],
});

export const submissionFilterTemporarySearchTextState = atom({
  key: "submissionFilterTemporarySearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const submissionFilterActualSearchTextState = atom({
  key: "submissionFilterActualSearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const submissionFilterStatusState = atom({
  key: "submissionFilterStatusState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});
export const submissionFilterJoinedAfterState = atom({
  key: "submissionFilterJoinedAfterState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});
export const submissionFilterTenantIDState = atom({
  key: "submissionFilterTenantIDState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});
export const submissionFilterTenantNameState = atom({
  key: "submissionFilterTenantNameState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

// --- User List --- //

// Control whether to show filters for the list.
export const userFilterShowState = atom({
  key: "userFilterShowState",
  default: false,
  effects_UNSTABLE: [persistAtom],
});

export const userFilterTemporarySearchTextState = atom({
  key: "userFilterTemporarySearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const userFilterActualSearchTextState = atom({
  key: "userFilterActualSearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const userFilterRoleState = atom({
  key: "userFilterRoleState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const userFilterStatusState = atom({
  key: "userFilterStatusState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const userFilterJoinedAfterState = atom({
  key: "userFilterJoinedAfterState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

// --- Tenant List --- //

// Control whether to show filters for the list.
export const tenantFilterShowState = atom({
  key: "tenantFilterShowState",
  default: false,
  effects_UNSTABLE: [persistAtom],
});

export const tenantFilterTemporarySearchTextState = atom({
  key: "tenantFilterTemporarySearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const tenantFilterActualSearchTextState = atom({
  key: "tenantFilterActualSearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const tenantFilterRoleState = atom({
  key: "tenantFilterRoleState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const tenantFilterStatusState = atom({
  key: "tenantFilterStatusState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const tenantFilterJoinedAfterState = atom({
  key: "tenantFilterJoinedAfterState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

// --- Offers --- //

// Control whether to show filters for the list.
export const offersFilterShowState = atom({
  key: "offersFilterShowState",
  default: false,
  effects_UNSTABLE: [persistAtom],
});

export const offersFilterTemporarySearchTextState = atom({
  key: "offersFilterTemporarySearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const offersFilterActualSearchTextState = atom({
  key: "offersFilterActualSearchTextState",
  default: "",
  effects_UNSTABLE: [persistAtom],
});

export const offersFilterStatusState = atom({
  key: "offersFilterStatusState",
  default: 2,
  effects_UNSTABLE: [persistAtom],
});

// --- Submissions --- //

export const ADD_COMIC_SUBMISSION_STATE_DEFAULT = {
    seriesTitle: "",
    issueVol: "",
    issueNo: "",
    issueCoverYear: 0,
    issueCoverMonth: 0,
    publisherName: 0,
    publisherNameOther: "",
    isKeyIssue: "",
    keyIssue: 0,
    keyIssueOther: "",
    keyIssueDetail: "",
    isInternationalEdition: false,
    isVariantCover: false,
    variantCoverDetail: "",
    printing: 1,
    primaryLabelDetails: 2,
    primaryLabelDetailsOther: "",
    specialNotes: "",
    gradingNotes: "",
    signatures: [],
    creasesFinding: "",
    tearsFinding: "",
    missingPartsFinding: "",
    stainsFinding: "",
    distortionFinding: "",
    paperQualityFinding: "",
    spineFinding: "",
    coverFinding: "",
    gradingScale: 0,
    overallLetterGrade: "",
    isOverallLetterGradeNearMintPlus: false,
    overallNumberGrade: "",
    cpsPercentageGrade: "",
    showsSignsOfTamperingOrRestoration: 2,
    status: 1, //1=Waiting to Receive, 7=Completed by Retail Partner
    serviceType: 0,
    tenantID: "",
    collectibleType: 1, // 1=Comic, 2=Card
    customerID: "",
}

// Control whether to show filters for the list.
export const addComicSubmissionState = atom({
  key: "addComicSubmission",
  default: ADD_COMIC_SUBMISSION_STATE_DEFAULT,
  effects_UNSTABLE: [persistAtom],
});
