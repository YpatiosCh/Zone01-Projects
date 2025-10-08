const pick = (obj, keys) => {
    // create the result object
    const result = {};

    // create/get the array from keys
    const keyArray = Array.isArray(keys) ? keys : [keys];

    // loop through keyArray to see if obj has any of its keys
    for (const key of keyArray) {
        if (obj.hasOwnProperty(key)) {
            result[key] = obj[key];
        }
    }

    return result;
}

function omit(obj, keys) {
  const result = {};
  const keyArray = Array.isArray(keys) ? keys : [keys];

  for (const key in obj) {
    if (Object.prototype.hasOwnProperty.call(obj, key) && !keyArray.includes(key)) {
      result[key] = obj[key];
    }
  }

  return result;
}