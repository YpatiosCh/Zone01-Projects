function multiply(a, b) {
    let result = 0;

    // flag to see if multiplier is negative
    const negative = b < 0;

    // create absolute number of multiplier
    const absB = Math.abs(b);

    // multiply by iteratively adding to result
    for (let i = 0; i < absB; i++) {
        result += a;
    }

    // and now perfection
    return negative ? -result : result;

    // if negative is false return result, if true return -result. So far, thats the most dope shit of js
}

function divide(a, b) {
    let result = 0;

    // if any is 0 then return 0 
    if (a === 0 || b === 0) return result;

    // check if any of those is negative
    const negative = (a < 0) !== (b < 0);

    // get absolute of both
    let absA = Math.abs(a);
    let absB = Math.abs(b);


    // divide 
    while (absA >= absB) {
        absA -= absB;
        result++
    }

    // return with sign
    return negative ? -result : result;
}

function modulo(a, b) {
  if (a === 0 || b === 0) return 0;

  const absA = Math.abs(a);
  const absB = Math.abs(b);
  let remainder = absA;

  while (remainder >= absB) {
    remainder -= absB;
  }

  return a < 0 ? -remainder : remainder;
}
