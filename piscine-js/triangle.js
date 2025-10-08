function triangle(char, n) {
    let triangle = '';

    for (let i = 1; i <= n; i++) {
        if (i === n) {
            triangle += char.repeat(i);
        } else {
            triangle += char.repeat(i) + '\n';
        }
    }

    return triangle;
}