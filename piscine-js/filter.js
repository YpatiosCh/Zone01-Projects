function filter(array, callback) {
  const result = [];
  for (let i = 0; i < array.length; i++) {
    if (callback(array[i], i, array)) {
      result.push(array[i]);
    }
  }
  return result;
}

function reject(array, callback) {
  const result = [];
  for (let i = 0; i < array.length; i++) {
    if (!callback(array[i], i, array)) {
      result.push(array[i]);
    }
  }
  return result;
}

function partition(array, callback) {
  const truthy = [];
  const falsy = [];
  
  for (let i = 0; i < array.length; i++) {
    if (callback(array[i], i, array)) {
      truthy.push(array[i]);
    } else {
      falsy.push(array[i]);
    }
  }
  
  return [truthy, falsy];
}