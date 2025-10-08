//  returns an array of strings from the city key
function citiesOnly(objArr) {
    return objArr.map(obj => {
        return obj.city;
    }); 
}

//  accepts an array of strings, and returns a new array of strings. 
// The returned array will be the same as the argument, except the first letter of every word must be capitalized.
function upperCasingStates(arr) {
    
    // .map :
    // 1. Loops through each item in the array.
    // 2. Applies a function to each item.
    // 3. Returns a new array with the results of that function.
    return arr.map(state => {
        return state.split(' ')
        .map(word => word.charAt(0).toUpperCase() + word.slice(1))
        .join(' ');
    });
}

// convert an array of fahrenheit to an array of celcius
function fahrenheitToCelsius(fTemps) {
    // return the array .map creates
    
    return fTemps.map(f => {
        // Remove the "°F" and convert to number
        const fahrenheit = parseInt(f);
        // Convert to Celsius and round down
        const celsius = Math.floor((fahrenheit - 32) * 5 / 9);
        // Return as string with "°C"
        return `${celsius}°C`;
    });
}

function trimTemp(objArr) {
  return objArr.map(obj => ({
    ...obj,
    temperature: obj.temperature.replaceAll(" ", "")
  }));
}


function tempForecasts(objArr) {

    return objArr.map(obj => {
        // convert to celsius
        const tempF = parseInt(obj.temperature);
        const celsius = Math.floor((tempF - 32)* 5 / 9);

        const state = upperCasingStates([obj.state])[0];

        return `${celsius}°Celsius in ${obj.city}, ${state}`;
    });
}