function flow(functions) {
    return function(...args) {
        let result = functions[0](...args);
        
        for (let i = 1; i < functions.length; i++) {
            result = functions[i](result);
        }
        
        return result;
    };
}

/* The flow function creates a new function that:

Takes an array of functions as input
Returns a new function that accepts any number of arguments
Applies the functions in sequence:

The first function gets all the original arguments
Each subsequent function gets the result of the previous function as its single argument

*/