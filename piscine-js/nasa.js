function nasa(n) {
    let result = '';
    let sep = ' ';

    for (let i = 1; i <= n; i++) {
        if (i % 3 === 0 && i % 5 === 0) {
            result += 'NASA';
        } else if (i % 3 === 0) {
            result += 'NA';
        } else if (i % 5 === 0) {
            result += 'SA';
        } else {
            result += i.toString();
        }

        if (i !== n) {
            result += sep;
        }
    }

    return result;
}