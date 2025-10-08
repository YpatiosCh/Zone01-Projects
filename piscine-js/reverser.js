function reverse(s) {
    let result = [];

    let n = s.length-1;

    for (let i = n; i >= 0; i--) {
        result.push(s[i]);
    }

    return typeof s === 'string' ? result.join('') : result;
}