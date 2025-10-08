function deepCopy(input) {
  // Primitives and functions
  if (input === null || typeof input !== 'object') {
    return input;
  }

  // Handle RegExp
  if (input instanceof RegExp) {
    return new RegExp(input);
  }

  // Handle Date
  if (input instanceof Date) {
    return new Date(input);
  }

  // Handle Function — return same reference
  if (typeof input === 'function') {
    return input;
  }

  // Handle Array
  if (Array.isArray(input)) {
    return input.map(deepCopy);
  }

  // Handle plain Object
  const copy = {};
  for (const key in input) {
    if (input.hasOwnProperty(key)) {
      copy[key] = deepCopy(input[key]);
    }
  }

  return copy;
}
