// words splits string by spaces and returns an array of those words
function words(x) {
    return x.split(' ');
}

// sentence joins the array of strings into a sentence
function sentence(x) {
    return x.join(' ');
}

// yell returns the string in uppercase
function yell(x) {
    return x.toUpperCase();
}

// whisper returns the string in lowercase, surrounded by asterisks (*)
function whisper(x) {
    return `*${x.toLowerCase()}*`;
}

// capitalize capitalizes the first letter of the string
function capitalize(x) {
    if (x.length === 0) return '';
    return x[0].toUpperCase() + x.slice(1).toLowerCase();
}