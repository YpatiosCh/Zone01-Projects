const body = document.body;
const MODES = ["mode-login", "mode-profile"];

export function setLayoutMode(mode) {
  if (!body) {
    return;
  }
  body.classList.remove(...MODES);
  const normalized = typeof mode === "string" ? `mode-${mode}` : null;
  if (normalized && MODES.includes(normalized)) {
    body.classList.add(normalized);
  }
}