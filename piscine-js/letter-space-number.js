function letterSpaceNumber(str) {
    let arr = str.match(/[a-z] [0-9](?![a-z0-9])/gi);
    return arr !== null ? arr : [];
}



/*
What it matches:

[a-zA-Z]:  match a single letter

' ':       space

\d:        match one digit

(?!\d):    ensures it's not followed by another digit

*/