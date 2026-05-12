import { ROUTES } from "../constants.js";
import { fetchProfileData, logoutUser } from "../services/auth.js";
import { navigate } from "../router.js";
import { getAuthState, setAuthState } from "../state.js";
import {
  ensureStatusMessage,
  getAppRoot,
  removeElementById,
  updateStatus,
} from "../ui/dom.js";
import { setLayoutMode } from "../ui/layout.js";
import {
  deriveUserProfile,
  extractTotalXP,
  formatISODate,
  formatXPAmount,
  deriveCollaboratorCounts,
  deriveProgressSeries,
  deriveProjectList,
} from "../utils/profile.js";
import { createCollaboratorChart } from "../charts/collaborators.js";
import { createProgressChart } from "../charts/progress.js";

export function renderProfileView() {
  const appRoot = getAppRoot();

  if (!appRoot) {
    return;
  }

  removeElementById("login-form");
  removeElementById("login-panel");

  setLayoutMode("profile");
  updateStatus("");

  const statusMessage = ensureStatusMessage();
  if (statusMessage) {
    statusMessage.classList.remove("status-message--login");
  }

  ensureLogoutButton();

  const container = ensureProfileContainer(appRoot);
  const state = getAuthState();

  if (state.profileData) {
    renderProfileSections(state.profileData, container);
  } else {
    container.innerHTML = "";
    const skeleton = document.createElement("div");
    skeleton.className = "profile-loading";
    skeleton.textContent = "Loading profile data…";
    container.appendChild(skeleton);
    updateStatus("Loading profile…");
  }

  void loadProfileData(container);
}

function ensureProfileContainer(root) {
  let container = document.getElementById("profile-data");
  if (!container) {
    container = document.createElement("div");
    container.id = "profile-data";
    container.className = "profile-data";
    root.appendChild(container);
  }
  const statusMessage = document.getElementById("status-message");
  if (statusMessage && statusMessage.parentElement === root) {
    const nextSibling = statusMessage.nextSibling;
    if (nextSibling !== container) {
      root.insertBefore(container, nextSibling);
    }
  } else if (!container.parentElement) {
    root.appendChild(container);
  }

  return container;
}

function ensureLogoutButton() {
  let button = document.getElementById("logout-button");
  if (!button) {
    button = document.createElement("button");
    button.id = "logout-button";
    button.type = "button";
  }

  button.className = "logout-button logout-button--floating";
  button.disabled = false;
  button.textContent = "Log Out";

  button.removeEventListener("click", handleLogout);
  button.addEventListener("click", handleLogout, { once: false });

  if (button.parentElement !== document.body) {
    document.body.appendChild(button);
  }

  return button;
}

async function handleLogout(event) {
  const button = event.currentTarget;
  updateStatus("");

  button.disabled = true;
  button.textContent = "Logging out...";

  const result = await logoutUser();

  if (!result.ok) {
    button.disabled = false;
    button.textContent = "Log Out";
    updateStatus(result.message, "error");
    return;
  }

  navigate(ROUTES.LOGIN, {
    replace: true,
    renderOptions: {
      statusMessage: result.message,
      statusType: "success",
    },
  });
}

async function loadProfileData(container) {
  const result = await fetchProfileData();

  if (result?.status === 401 || result?.status === 403) {
    updateStatus("Session expired. Please sign in again.", "error");
    container.textContent = "Session expired.";
    navigate(ROUTES.LOGIN, {
      replace: true,
      renderOptions: {
        statusMessage: "Session expired. Please sign in again.",
        statusType: "error",
      },
    });
    return;
  }

  if (!result?.ok) {
    updateStatus(result?.message || "Unable to load profile.", "error");
    container.textContent =
      result?.message || "Unable to load profile. Please try again.";
    return;
  }

  setAuthState("authenticated", result.data);
  updateStatus("");
  renderProfileSections(result.data, container);
}

function renderProfileSections(data, container) {
  container.innerHTML = "";

  const profile = deriveUserProfile(data?.user);
  const totalXP = extractTotalXP(data?.xp);
  const collaboratorCounts = deriveCollaboratorCounts(data?.collaborators);
  const progressSeries = deriveProgressSeries(data?.progress);
  const projects = deriveProjectList(data?.projects);

  const page = document.createElement("div");
  page.className = "profile-page";
  container.appendChild(page);

  const header = document.createElement("header");
  header.className = "profile-header";
  page.appendChild(header);

  const summary = createProfileSummary(profile, totalXP, projects);
  if (summary) {
    header.appendChild(summary);
  } else {
    const fallback = document.createElement("div");
    fallback.className = "profile-summary profile-summary--empty";
    fallback.textContent = "No profile information available.";
    header.appendChild(fallback);
  }

  const charts = document.createElement("section");
  charts.className = "profile-charts";

  const progressChart = createProgressChart(progressSeries);
  if (progressChart) {
    charts.appendChild(progressChart);
  }

  const collabChart = createCollaboratorChart(collaboratorCounts);
  if (collabChart) {
    charts.appendChild(collabChart);
  }

  if (charts.childElementCount > 0) {
    page.appendChild(charts);
  }

}

function createProfileSummary(profile, totalXP, projects) {
  if (!profile && typeof totalXP !== "number" && (!projects || projects.length === 0)) {
    return null;
  }

  const summary = document.createElement("section");
  summary.className = "profile-summary";

  const left = document.createElement("div");
  left.className = "profile-summary__left";
  summary.appendChild(left);

  const userInfo = createUserInfoSection(profile);
  if (userInfo) {
    left.appendChild(userInfo);
  }

  if (typeof totalXP === "number") {
    left.appendChild(createXpSection(totalXP));
  }

  const projectsSection = createProjectsSection(projects);
  if (projectsSection) {
    summary.appendChild(projectsSection);
  }

  return summary;
}

function createUserInfoSection(profile) {
  if (!profile) {
    return null;
  }

  const container = document.createElement("div");
  container.className = "profile-summary__user";

  const name = document.createElement("h2");
  name.className = "profile-summary__name";
  name.textContent = profile.fullName || profile.login || "Your profile";
  container.appendChild(name);

  if (profile.login) {
    const loginEl = document.createElement("p");
    loginEl.className = "profile-summary__login";
    loginEl.textContent = `@${profile.login}`;
    container.appendChild(loginEl);
  }

  const metaList = document.createElement("div");
  metaList.className = "profile-summary__meta";
  container.appendChild(metaList);

  const metaEntries = [
    { label: "Email", value: profile.email },
    { label: "Phone", value: profile.phone },
    { label: "Campus", value: profile.campus },
    { label: "Joined", value: formatISODate(profile.createdAt) },
  ].filter((entry) => entry.value);

  metaEntries.forEach((entry) => {
    metaList.appendChild(createMetaRow(entry.label, entry.value));
  });

  return container;
}

function createXpSection(totalXP) {
  const wrapper = document.createElement("div");
  wrapper.className = "profile-summary__xp";

  const label = document.createElement("span");
  label.className = "profile-summary__xp-label";
  label.textContent = "Total XP";

  const value = document.createElement("span");
  value.className = "profile-summary__xp-value";
  value.textContent = formatXPAmount(totalXP);

  wrapper.append(label, value);
  return wrapper;
}

function createProjectsSection(projects) {
  if (!projects || projects.length === 0) {
    return null;
  }

  const container = document.createElement("div");
  container.className = "profile-summary__right";

  const title = document.createElement("h3");
  title.className = "profile-summary__projects-title";
  title.textContent = "Projects Completed";
  container.appendChild(title);

  const list = document.createElement("ul");
  list.className = "profile-summary__projects-list";

  projects.forEach((project) => {
    const item = document.createElement("li");
    item.textContent = project;
    list.appendChild(item);
  });

  container.appendChild(list);
  return container;
}

function createMetaRow(labelText, valueText) {
  const wrapper = document.createElement("div");
  wrapper.className = "profile-summary__meta-row";

  const label = document.createElement("span");
  label.className = "profile-summary__meta-label";
  label.textContent = labelText;

  const value = document.createElement("span");
  value.className = "profile-summary__meta-value";
  value.textContent = valueText;

  wrapper.append(label, value);
  return wrapper;
}
