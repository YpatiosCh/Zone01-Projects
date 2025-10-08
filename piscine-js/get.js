function get(src, path) {
  const keys = path.split('.'); // split the path into parts
  let current = src;

  for (let key of keys) {
    if (current === null || current === undefined) {
      return undefined; // if the path breaks at any point
    }
    current = current[key];
  }

  return current;
}
