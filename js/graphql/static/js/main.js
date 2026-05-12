import { ROUTES } from "./constants.js";
import { configureRouter, initRouter } from "./router.js";
import { bootstrapAuthState } from "./services/auth.js";
import { renderLoginView } from "./views/login.js";
import { renderProfileView } from "./views/profile.js";

async function start() {
  configureRouter({
    loginRenderer: renderLoginView,
    profileRenderer: renderProfileView,
  });

  const initialPath =
    typeof window !== "undefined" ? window.location.pathname : ROUTES.LOGIN;

  const bootstrapResult = await bootstrapAuthState();

  let renderOptions;
  if (
    bootstrapResult &&
    bootstrapResult.status === "anonymous" &&
    bootstrapResult.message
  ) {
    renderOptions = {
      statusMessage: bootstrapResult.message,
      statusType: "error",
    };
  }

  initRouter({
    initialPath,
    renderOptions,
  });
}

function run() {
  start().catch((err) => {
    // eslint-disable-next-line no-console
    console.error("Failed to initialize application", err);
  });
}

if (document.readyState === "loading") {
  document.addEventListener("DOMContentLoaded", run);
} else {
  run();
}
