const vowels = /[aeiou]/gi

function vowelDots(str) {
  return str.replace(vowels, (v) => v + '.')
}

/*

What is (v) => v + '.'?
This is an arrow function in JavaScript. It's a shorter way of writing a function.


v is the parameter of the function (in this case, each vowel matched by the regular expression).

v + '.' is the return value — it adds a dot after the vowel.

This function is passed to .replace() so that every vowel in the string gets replaced with itself + '.'.

JavaScript automatically calls your function for every match found by the regex — and it passes the matched string (the part it found) as the first argument.

*/
