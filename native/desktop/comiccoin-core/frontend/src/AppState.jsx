import {
  RecoilRoot,
  atom,
  selector,
  useRecoilState,
  useRecoilValue,
  // AtomEffect,
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

export const currentOpenWalletAtAddressState = atom({
  key: "currentOpenWalletAtAddress",
  default: "",
  effects_UNSTABLE: [persistAtom],
});
