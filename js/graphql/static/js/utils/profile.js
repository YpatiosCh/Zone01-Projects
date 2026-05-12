const NAME_PARTS = ["firstName", "lastName"];

function coalesceString(value) {
  return typeof value === "string" ? value.trim() : "";
}

export function deriveUserProfile(userResponse) {
  if (!userResponse || typeof userResponse !== "object") {
    return null;
  }

  const users = userResponse?.data?.user;
  const record = Array.isArray(users) && users.length > 0 ? users[0] : null;
  if (!record) {
    return null;
  }

  const attrs = record.attrs ?? {};
  const fullName = NAME_PARTS.map((key) => coalesceString(attrs[key]))
    .filter(Boolean)
    .join(" ")
    .trim();

  return {
    id: record.id ?? null,
    login: coalesceString(record.login),
    campus: coalesceString(record.campus),
    email: coalesceString(attrs.email),
    phone: coalesceString(attrs.phoneNumber),
    fullName: fullName || null,
    createdAt: coalesceString(record.createdAt),
  };
}

export function extractTotalXP(xpResponse) {
  if (!xpResponse || typeof xpResponse !== "object") {
    return null;
  }
  const amount =
    xpResponse?.data?.transaction_aggregate?.aggregate?.sum?.amount;

  return typeof amount === "number" ? amount : null;
}

export function formatXPAmount(amount) {
  if (typeof amount !== "number" || !Number.isFinite(amount) || amount < 0) {
    return "—";
  }

  const KB = 1000;
  const MB = KB * 1000;

  if (amount >= MB) {
    return `${(amount / MB).toFixed(2)} MB`;
  }

  if (amount >= KB) {
    return `${(amount / KB).toFixed(2)} KB`;
  }

  return `${amount.toFixed(0)} B`;
}

export function deriveCollaboratorCounts(collaboratorsResponse) {
  if (!collaboratorsResponse || typeof collaboratorsResponse !== "object") {
    return [];
  }

  const records = collaboratorsResponse?.data?.group_user;
  if (!Array.isArray(records) || records.length === 0) {
    return [];
  }

  const counts = new Map();

  records.forEach((entry) => {
    const login = entry?.user?.login;
    if (!login) {
      return;
    }
    const normalized = String(login).trim();
    if (!normalized) {
      return;
    }
    counts.set(normalized, (counts.get(normalized) || 0) + 1);
  });

  return Array.from(counts.entries())
    .map(([login, count]) => ({ login, count }))
    .sort((a, b) => {
      if (b.count !== a.count) {
        return b.count - a.count;
      }
      return a.login.localeCompare(b.login);
    });
}

export function deriveProgressSeries(progressResponse) {
  if (!progressResponse || typeof progressResponse !== "object") {
    return [];
  }

  const records = progressResponse?.data?.transaction;
  if (!Array.isArray(records) || records.length === 0) {
    return [];
  }

  const byDate = new Map();

  records.forEach((entry) => {
    const amount = typeof entry?.amount === "number" ? entry.amount : null;
    const createdAt = entry?.createdAt;
    if (amount === null || amount < 0 || !createdAt) {
      return;
    }

    const timestamp = Date.parse(createdAt);
    if (Number.isNaN(timestamp)) {
      return;
    }

    const dateKey = new Date(timestamp);
    const normalizedKey = new Date(
      Date.UTC(dateKey.getUTCFullYear(), dateKey.getUTCMonth(), dateKey.getUTCDate()),
    ).toISOString();

    byDate.set(normalizedKey, (byDate.get(normalizedKey) || 0) + amount);
  });

  return Array.from(byDate.entries())
    .map(([isoDate, amount]) => ({
      date: new Date(isoDate),
      value: amount,
    }))
    .sort((a, b) => a.date - b.date);
}

function formatProjectName(raw) {
  if (typeof raw !== "string" || !raw.trim()) {
    return null;
  }
  const segments = raw.trim().split("/");
  const name = segments[segments.length - 1];
  if (!name) {
    return null;
  }
  return name.replace(/[-_]+/g, " ").trim();
}

export function deriveProjectList(projectsResponse) {
  if (!projectsResponse || typeof projectsResponse !== "object") {
    return [];
  }

  const records = projectsResponse?.data?.transaction;
  if (!Array.isArray(records) || records.length === 0) {
    return [];
  }

  const unique = new Set();
  records.forEach((entry) => {
    const formatted = formatProjectName(entry?.path);
    if (formatted) {
      unique.add(formatted);
    }
  });

  return Array.from(unique).sort((a, b) => a.localeCompare(b));
}

export function formatISODate(isoString) {
  if (!isoString) {
    return null;
  }

  const date = new Date(isoString);
  if (Number.isNaN(date.getTime())) {
    return null;
  }

  return date.toLocaleDateString(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}
