function fold(array, func, accumulator) {
    let result = accumulator;
    for (let i = 0; i < array.length; i++) {
        result = func(result, array[i]);
    }
    return result;
}

function foldRight(array, func, accumulator) {
    let result = accumulator;
    for (let i = array.length - 1; i >= 0; i--) {
        result = func(result, array[i]);
    }
    return result;
}

function reduce(array, func) {
    let result = array[0];
    for (let i = 1; i < array.length; i++) {
        result = func(result, array[i]);
    }
    return result;
}

function reduceRight(array, func) {  
    let result = array[array.length - 1];
    for (let i = array.length - 2; i >= 0; i--) {
        result = func(result, array[i]);
    }
    return result;
}