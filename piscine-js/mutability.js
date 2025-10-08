// clone1 and clone2 must keep original values
const clone1 = { ...person };           // shallow copy
const clone2 = Object.assign({}, person); // another shallow copy

// samePerson must reflect changes to `person`
const samePerson = person;              // reference, not copy

// Modify person
person.age += 1;
person.country = 'FR';