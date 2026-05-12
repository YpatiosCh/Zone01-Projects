import { ROUTES } from "../constants.js";
import { loginUser } from "../services/auth.js";
import { navigate } from "../router.js";
import {
  getAppRoot,
  removeElementById,
  updateStatus,
  ensureStatusMessage,
} from "../ui/dom.js";
import { setLayoutMode } from "../ui/layout.js";

function buildLoginForm() {
  const form = document.createElement("form");
  form.id = "login-form";
  form.className = "login-form";
  form.method = "post";
  form.action = "/api/auth/login";
  form.noValidate = true;

  form.innerHTML = "";

  const subtitle = document.createElement("p");
  subtitle.className = "login-form__subtitle";
  subtitle.textContent = "Sign in to view your Zone01 profile.";
  form.appendChild(subtitle);

  const identifierGroup = document.createElement("div");
  identifierGroup.className = "login-form__group";
  identifierGroup.innerHTML = `
    <label for="identifier">Identifier</label>
    <input type="text" id="identifier" name="identifier" autocomplete="username" required />
  `;
  form.appendChild(identifierGroup);

  const passwordGroup = document.createElement("div");
  passwordGroup.className = "login-form__group";
  passwordGroup.innerHTML = `
    <label for="password">Password</label>
    <input type="password" id="password" name="password" autocomplete="current-password" required />
  `;
  form.appendChild(passwordGroup);

  const actions = document.createElement("div");
  actions.className = "login-form__actions";

  const submit = document.createElement("button");
  submit.type = "submit";
  submit.className = "login-form__submit";
  submit.textContent = "Sign In";
  actions.appendChild(submit);

  form.appendChild(actions);

  const error = document.createElement("p");
  error.id = "login-error";
  error.className = "form-error";
  error.setAttribute("role", "alert");
  form.appendChild(error);

  form.addEventListener("submit", handleLoginSubmit, { once: false });

  return form;
}

export function renderLoginView(options = {}) {
  const { statusMessage: statusText = "", statusType } = options;
  const appRoot = getAppRoot();

  if (!appRoot) {
    return;
  }

  removeElementById("login-form");
  removeElementById("login-panel");
  removeElementById("logout-button");
  removeElementById("profile-data");

  setLayoutMode("login");

  const form = buildLoginForm();

  const panel = document.createElement("div");
  panel.id = "login-panel";
  panel.className = "login-panel";
  panel.appendChild(form);

  const statusEl = ensureStatusMessage();
  if (statusEl) {
    statusEl.classList.add("status-message--login");
    panel.appendChild(statusEl);
  }

  appRoot.appendChild(panel);

  updateStatus(statusText, statusType);
}

async function handleLoginSubmit(event) {
  event.preventDefault();

  const form = event.currentTarget;
  const identifierField = form.querySelector("#identifier");
  const passwordField = form.querySelector("#password");
  const submitButton = form.querySelector('button[type="submit"]');
  const errorMessage = form.querySelector("#login-error");

  if (!identifierField || !passwordField || !submitButton) {
    return;
  }

  if (errorMessage) {
    errorMessage.textContent = "";
  }

  const identifier = identifierField.value.trim();
  const password = passwordField.value;

  if (!identifier || !password) {
    if (errorMessage) {
      errorMessage.textContent = "Identifier and password are required.";
    }
    return;
  }

  try {
    submitButton.disabled = true;
    submitButton.textContent = "Signing in...";

    const result = await loginUser(identifier, password);
    if (!result.ok) {
      if (errorMessage) {
        errorMessage.textContent = result.message;
      }
      return;
    }

    navigate(ROUTES.PROFILE, { replace: true });
  } finally {
    submitButton.disabled = false;
    submitButton.textContent = "Sign In";
  }
}
