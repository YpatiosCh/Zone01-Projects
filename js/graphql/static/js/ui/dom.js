const appRoot = document.getElementById("app");
let statusMessage = document.getElementById("status-message");

export function getAppRoot() {
  return appRoot;
}

export function ensureStatusMessage() {
  if (statusMessage) {
    return statusMessage;
  }

  statusMessage = document.createElement("p");
  statusMessage.id = "status-message";
  statusMessage.className = "status-message";
  statusMessage.setAttribute("aria-live", "polite");

  return statusMessage;
}

export function updateStatus(message = "", type) {
  const el = ensureStatusMessage();
  if (!el) {
    return;
  }

  el.textContent = message;
  el.classList.remove("status-error", "status-success");

  if (type === "error") {
    el.classList.add("status-error");
  } else if (type === "success") {
    el.classList.add("status-success");
  }
}

export function removeElementById(id) {
  const el = document.getElementById(id);
  if (el) {
    el.remove();
  }
}
