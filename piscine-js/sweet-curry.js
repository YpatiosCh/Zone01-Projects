// mult2: multiplies two numbers
//  Takes one number, returns a function that takes the second number and multiplies them.
const mult2 = (a) => {
    return (b) => {
        return a * b;
    };
};

// add3: adds three numbers
// Takes one number, returns a function that takes the second number, 
// which returns a function that takes the third number and adds all three.
const add3 = (a) => {
    return (b) => {
        return (c) => {
            return a + b + c;
        };
    };
};

// sub4: subtracts four numbers in order (a - b - c - d)
// Takes one number, returns a function that takes the second number, which returns a function that takes the third number, 
// which returns a function that takes the fourth number and subtracts them in order (a - b - c - d).
const sub4 = (a) => {
    return (b) => {
        return (c) => {
            return (d) => {
                return a - b - c - d;
            };
        };
    };
};