const map = (arr, func) => {
    const result = [];

    for (let i = 0; i < arr.length; i++) {
        result.push(func(arr[i], i, arr));
    }

    return result;
}

const flatMap = (arr, func) => {
    const result = [];
    for (let i = 0; i < arr.length; i++) {
        const mapped = func(arr[i], i, arr);
        // If the result is an array, spread it into the result
        if (Array.isArray(mapped)) {
            result.push(...mapped);
        } else {
            result.push(mapped);
        }
    }

    return result;
}