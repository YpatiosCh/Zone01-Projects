function join(arr, joinner) {
    if (joinner === null) {
        joinner = ",";
    }
    let result = arr[0].toString();
    for (let i = 1; i < arr.length; i++) {
        result += joinner + arr[i];
    }
    return result;
}

function split(str, sep) {
    if (sep === null) {
        sep = ",";
    }
    let result = [];
    if (sep === "") {
        for (let i = 0; i < str.length; i++) {
            result.push(str[i]);
        }
        return result;
    }
    let end = str.indexOf(sep);
    while (end > -1) {
        end = str.indexOf(sep);
        if (end === -1) {
            break;
        }
        result.push(str.slice(0, end));
        str = str.slice(end + sep.length);
    }
    result.push(str);
    return result;
}