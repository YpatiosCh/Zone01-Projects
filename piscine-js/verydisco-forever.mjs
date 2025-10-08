import { writeFile } from 'fs/promises';

try {
    const input = process.argv[2];
    if (!input) {
        throw new Error('No input provided');
    }
    const words = input.split(' ');
    for (let i = 0; i < words.length; i++) {
        const half = Math.ceil(words[i].length / 2);
        const firstHalf = words[i].slice(0, half);
        const secondHalf = words[i].slice(half);
        words[i] = secondHalf + firstHalf;
    }
    const output = words.join(' ');

    await writeFile('verydisco-forever.txt', output, 'utf8');
} catch (error) {
    console.error('Error writing to file:', error);
}
