function indexOf(arr, val, index) {
    // lets see first if the index to start from was provided 
    if (index === undefined) {
        index = 0;
    }

    // start from index to check if value matches
    for (let i = index; i < arr.length; i++) {
        if (arr[i] === val) {
            return i;
        }
    }

    // if no match return -1
    return -1;
}

function lastIndexOf(arr, val, index) {
    // check if index to start from was provided
    if (index === undefined) {
        index = arr.length-1;
    }

    // iterate backwards to check for matches
    for (let i = index; i >= 0; i--) {
        if (arr[i] === val) {
            return i;
        }
    }

    // if no matches found return -1
    return -1
}


function includes(arr, val) {

    for (let i = 0; i < arr.length; i++) {
        if (arr[i] === val) {
            return true;
        }
    }

    return false;
}