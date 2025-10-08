const input = process.argv[2];

const words = input.split(' ');

for (let i = 0; i < words.length; i++) {
    const half = Math.ceil(words[i].length / 2);
    const firstHalf = words[i].slice(0, half);
    const secondHalf = words[i].slice(half);
    words[i] = secondHalf + firstHalf;
}

const output = words.join(' ');
console.log(output);