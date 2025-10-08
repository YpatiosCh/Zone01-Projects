function adder(arr, initialValue = 0) {
    return arr.reduce((acc, current) => {
        return acc + current
    },initialValue);
}

function sumOrMul(numbers, initialValue = 0) {

    return numbers.reduce((accumulator, currentValue) => {
        if (currentValue % 2 === 0) {
            // Even number - multiply
            return accumulator * currentValue;
        } else {
        // Odd number - add
            return accumulator + currentValue;
        }
    }, initialValue);
}

function funcExec(functions, initialValue = 0) {
  return functions.reduce((accumulator, currentFunction) => {
    return currentFunction(accumulator);
  }, initialValue);
}
