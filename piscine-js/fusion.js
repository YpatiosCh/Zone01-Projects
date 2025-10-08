function fusion(obj1, obj2) {
  const result = {};

  const keys = new Set([...Object.keys(obj1), ...Object.keys(obj2)]);

  for (const key of keys) {
    const val1 = obj1[key];
    const val2 = obj2[key];

    // Only in one object
    if (!(key in obj1)) {
      result[key] = val2;
    } else if (!(key in obj2)) {
      result[key] = val1;
    }
    // Both values exist
    else if (Array.isArray(val1) && Array.isArray(val2)) {
      result[key] = val1.concat(val2);
    } else if (typeof val1 === "string" && typeof val2 === "string") {
      result[key] = val1 + " " + val2;
    } else if (typeof val1 === "number" && typeof val2 === "number") {
      result[key] = val1 + val2;
    } else if (isObject(val1) && isObject(val2)) {
      result[key] = fusion(val1, val2); // recursive merge
    } else {
      result[key] = val2; // Type mismatch or override
    }
  }

  return result;
}

function isObject(value) {
  return typeof value === 'object' && value !== null && !Array.isArray(value);
}
