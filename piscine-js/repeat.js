function repeat(str, num) {
    let result = '';
    if (num < 0) {
        return 'RangeError';
    } else if (num === 0) {
        return result;
    } else {
        for (let i= 0; i < num; i++) {
            result += str;
        }
    }

    return result;
}