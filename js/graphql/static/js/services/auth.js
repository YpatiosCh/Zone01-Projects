import { setAuthState } from "../state.js";

async function requestProfileData() {
  try {
    const response = await fetch("/api/profile", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify({}),
    });

    const contentType = response.headers.get("content-type") || "";
    const isJson = contentType.includes("application/json");
    const payload = isJson ? await response.json().catch(() => null) : null;

    if (!response.ok) {
      const message =
        (payload && (payload.message || payload.error)) ||
        `Request failed with status ${response.status}.`;
      return {
        ok: false,
        status: response.status,
        message,
      };
    }

    if (!payload?.success) {
      const message =
        payload?.message ||
        payload?.error ||
        "Unexpected response while loading profile.";
      return {
        ok: false,
        status: response.status,
        message,
      };
    }

    return {
      ok: true,
      status: response.status,
      data: payload.data,
    };
  } catch (err) {
    return {
      ok: false,
      status: 0,
      message: "Network error while loading profile. Please try again.",
    };
  }
}

export async function bootstrapAuthState() {
  const result = await requestProfileData();

  if (result?.ok) {
    setAuthState("authenticated", result.data);
    return { status: "authenticated" };
  }

  if (result?.status === 401 || result?.status === 403) {
    setAuthState("anonymous");
    return { status: "anonymous" };
  }

  setAuthState("anonymous");
  return {
    status: "anonymous",
    message: result?.message,
  };
}

export async function fetchProfileData() {
  const result = await requestProfileData();

  if (result?.ok) {
    setAuthState("authenticated", result.data);
  } else if (result?.status === 401 || result?.status === 403) {
    setAuthState("anonymous");
  }

  return result;
}

export async function loginUser(identifier, password) {
  try {
    const response = await fetch("/api/auth/login", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify({ identifier, password }),
    });

    const payload = await response.json().catch(() => null);

    if (!response.ok || !payload?.success) {
      const message =
        (payload && (payload.message || payload.error)) ||
        "Unable to sign in. Please try again.";
      return { ok: false, message };
    }

    setAuthState("authenticated");
    return { ok: true };
  } catch (err) {
    return {
      ok: false,
      message: "Something went wrong. Please check your network and try again.",
    };
  }
}

export async function logoutUser() {
  try {
    const response = await fetch("/api/auth/logout", {
      method: "POST",
      credentials: "include",
    });

    const contentType = response.headers.get("content-type") || "";
    const payload = contentType.includes("application/json")
      ? await response.json().catch(() => null)
      : null;

    if (!response.ok) {
      const message =
        (payload && (payload.message || payload.error)) ||
        "Unable to log out. Please try again.";
      return { ok: false, message };
    }

    const message =
      (payload && (payload.message || payload.status)) ||
      "You have been logged out.";
    setAuthState("anonymous");
    return { ok: true, message };
  } catch (err) {
    return {
      ok: false,
      message: "Network error while logging out. Please try again.",
    };
  }
}
