// get gets the key from the object and returns its value
function get(key) {
  return sourceObject[key];
}

// set sets the key in the object to the value and returns the value
function set(key, value) {
  sourceObject[key] = value;
  return value;
}

// const protects the reference to the object, 
// so we can still modify its properties but we cannot reassign the object