// Array to Set
function arrToSet(arr) {
  return new Set(arr);
}

// Array to String
function arrToStr(arr) {
  return arr.join('');
}

// Set to Array
function setToArr(set) {
  return [...set];
}

// Set to String
function setToStr(set) {
  return [...set].join('');
}

// String to Array
function strToArr(str) {
  return str.split('');
}

// String to Set
function strToSet(str) {
  return new Set(str);
}

// Map to Object
function mapToObj(map) {
  const obj = {};
  for (const [key, value] of map.entries()) {
    obj[key] = value;
  }
  return obj;
}

// Object to Array (returns array of values)
function objToArr(obj) {
  return Object.values(obj);
}

// Object to Map
function objToMap(obj) {
  return new Map(Object.entries(obj));
}

// Array to Object (index as keys)
function arrToObj(arr) {
  const obj = {};
  arr.forEach((val, i) => {
    obj[i] = val;
  });
  return obj;
}

// String to Object (index as keys)
function strToObj(str) {
  const obj = {};
  for (let i = 0; i < str.length; i++) {
    obj[i] = str[i];
  }
  return obj;
}

function superTypeOf(value) {
  if (value === null) return 'null';
  if (Array.isArray(value)) return 'Array';
  if (value instanceof Set) return 'Set';
  if (value instanceof Map) return 'Map';
  if (value instanceof Function) return 'Function';
  if (typeof value === 'object') return 'Object';
  return typeof value === 'string' ? 'String'
       : typeof value === 'number' ? 'Number'
       : typeof value === 'boolean' ? 'Boolean'
       : typeof value === 'undefined' ? 'undefined'
       : typeof value;
}

