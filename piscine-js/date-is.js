function isValid(date) {
    if (new Date(date).toString() === "Invalid Date") {
        return false;
    }

    if (!(date instanceof Date) && typeof date !== "number") {
        return false;
    }

    return true;
}

function isAfter(date1, date2) {
    if (date1 > date2) {
        return true;
    }
    return false;
}

function isBefore(date1, date2) {
    if (date1 < date2) {
        return true;
    }
    return false;
}

function isFuture(date) {
    if (!isValid(date)) {
        return false;
    }

    if (new Date(date).getTime() > Date.now()) {
        return true;
    }

    return false;
}

function isPast(date) {
    if (!isValid(date)) {
        return false;
    }

    if (new Date(date).getTime() < Date.now()) {
        return true;
    }

    return false;
}