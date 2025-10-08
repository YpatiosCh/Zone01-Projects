
function currify(func) {
    return function curried(...args) {
        if (args.length >= func.length) {
            // If we have enough arguments, call the original function
            return func(...args);
        } else {
            // If we don't have enough arguments, return a new function
            return function(...nextArgs) {
                return curried(...args, ...nextArgs);
            };
        }
    };
}

/*
The currify function works by:

Checking the argument count: It compares the number of arguments received (args.length) with the number of parameters the original function expects (func.length)
Two possible outcomes:

Enough arguments: If we have enough arguments, it calls the original function immediately with all the arguments
Not enough arguments: If we need more arguments, it returns a new function that will collect more arguments and recursively call curried with the combined arguments


Flexible usage: The curried function can be called:

One argument at a time: mult2Curried(2)(2)
Multiple arguments at once: mult2Curried(2, 2)
Partial application: add3Curried(1, 2)(3)

*/