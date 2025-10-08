function filterShortStateName(arr) {
    return arr.filter(word => word.length < 7);
}

function filterStartVowel(arr) {
    return arr.filter(word => startsWithVowel(word));
}

function startsWithVowel(str) {
  return /^[aeiou]/i.test(str);
}

function filter5Vowels(arr) {
    return arr.filter(word => has5vowels(word));
}

function has5vowels(str) {
    const vowels = str.match(/[aeiou]/gi); // match all vowels, case-insensitive
    return vowels !== null && vowels.length >= 5;
}

function filter1DistinctVowel(arr) {
    return arr.filter(word => hasUniqueVowel(word));
}

function hasUniqueVowel(str) {
    const vowels = str.toLowerCase().match(/[aeiou]/gi);
    if (!vowels) return false;

    const vowelSet = new Set(vowels);
    return vowelSet.size === 1;
}

function multiFilter(arr) {
  return arr.filter(obj => {
    const capitalOK = obj.capital && obj.capital.length >= 8;
    const nameOK = obj.name && !/^[aeiou]/i.test(obj.name);
    const tagOK = obj.tag && /[aeiou]/i.test(obj.tag);
    const regionOK = obj.region !== 'South';

    return capitalOK && nameOK && tagOK && regionOK;
  });
}


