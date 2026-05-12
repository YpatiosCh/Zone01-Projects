const authState = {
  status: "unknown",
  profileData: null,
};

export function getAuthState() {
  return authState;
}

export function setAuthState(status, profileData = null) {
  authState.status = status;
  authState.profileData = profileData;
}

export function clearProfileData() {
  authState.profileData = null;
}
