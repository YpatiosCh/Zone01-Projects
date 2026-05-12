import { ROUTES } from "./constants.js";
import { getAuthState } from "./state.js";

let renderLogin;
let renderProfile;
let currentRoute = null;

export function configureRouter({ loginRenderer, profileRenderer }) {
  renderLogin = loginRenderer;
  renderProfile = profileRenderer;
}

export function navigate(path, options = {}) {
  handleRoute(path, { ...options, updateHistory: true });
}

export function initRouter({ initialPath, renderOptions } = {}) {
  if (typeof window !== "undefined") {
    window.addEventListener("popstate", () => {
      handleRoute(window.location.pathname, { updateHistory: false });
    });
  }

  const targetPath = initialPath ?? ROUTES.LOGIN;
  handleRoute(targetPath, { replace: true, renderOptions });
}

function normalizePath(path) {
  if (!path) {
    return ROUTES.HOME;
  }

  if (typeof window !== "undefined") {
    try {
      const url = new URL(path, window.location.origin);
      return url.pathname || ROUTES.HOME;
    } catch (err) {
      // fall through
    }
  }

  return path.startsWith("/") ? path : `/${path}`;
}

function applyHistory(path, replace) {
  if (typeof window === "undefined" || !window.history) {
    return;
  }

  if (window.location.pathname === path && !replace) {
    return;
  }

  const method = replace ? "replaceState" : "pushState";
  window.history[method](window.history.state, "", path);
}

function handleRoute(
  path,
  { replace = false, updateHistory = true, renderOptions } = {},
) {
  if (!renderLogin || !renderProfile) {
    throw new Error("Router renderers not configured.");
  }

  let target = normalizePath(path);
  let shouldReplace = replace;
  const isAuthenticated = getAuthState().status === "authenticated";

  if (target === ROUTES.HOME) {
    target = isAuthenticated ? ROUTES.PROFILE : ROUTES.LOGIN;
    shouldReplace = true;
  }

  if (target === ROUTES.PROFILE && !isAuthenticated) {
    target = ROUTES.LOGIN;
    shouldReplace = true;
  }

  if (target === ROUTES.LOGIN && isAuthenticated) {
    target = ROUTES.PROFILE;
    shouldReplace = true;
  }

  if (updateHistory) {
    applyHistory(target, shouldReplace);
  }

  const requiresRender =
    currentRoute !== target ||
    (target === ROUTES.LOGIN && renderOptions !== undefined);

  if (!requiresRender) {
    return;
  }

  currentRoute = target;

  if (target === ROUTES.PROFILE) {
    renderProfile();
    return;
  }

  renderLogin(renderOptions);
}
