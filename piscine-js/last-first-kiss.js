// first returns the first element of an array or string
function first(yolo) {
    return yolo[0];
}

// last returns the last element of an array or string
function last(yolo) {
    return yolo[yolo.length - 1];
}

// kiss 
function kiss(yolo) {
    // create an array to store the result
    let result = [];
    // get the last 
    let last = yolo[yolo.length - 1];
    // get the first
    let first = yolo[0];

    // push the last first so it can go to index 0
    result.push(last);
    // push the first so it can go to index 1
    result.push(first);
    return result;
}