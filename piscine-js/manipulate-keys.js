function filterKeys(obj, callback) {
  return Object.fromEntries(
    Object.entries(obj).filter(([key]) => callback(key))
  );
}

function mapKeys(obj, callback) {
  return Object.fromEntries(
    Object.entries(obj).map(([key, value]) => [callback(key), value])
  );
}

function reduceKeys(obj, callback, initial) {
  const keys = Object.keys(obj);
  return initial !== undefined
    ? keys.reduce(callback, initial)
    : keys.reduce(callback);
}
