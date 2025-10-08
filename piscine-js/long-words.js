const longWords = (arr) => {
    return arr.every(word => checkStrAndLength5(word));
}

const checkStrAndLength5 = (s) => {
    if (typeof s === "string" && s.length >= 5) {
        return true;
    }
    return false;
}

const oneLongWord = (arr) => {
    return arr.some(word => checkStrAndLength10(word));
}

const checkStrAndLength10 = (s) => {
    if (typeof s === "string" && s.length >= 10) {
        return true;
    }
    return false;
}

const noLongWords = (arr) => {
    return arr.every(word => checkStrAndLength7(word));
}

const checkStrAndLength7 = (s) => {
    if (typeof s === "string" && s.length >= 7) {
        return false;
    }
    return true;
}