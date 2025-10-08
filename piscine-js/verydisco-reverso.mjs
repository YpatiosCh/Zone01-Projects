import { readFile } from "fs/promises";


try {

    const fileName = process.argv[2];
    if (!fileName) {
        throw new Error("No file name provided");
    }

    const content = await readFile(fileName, "utf8");
    if (!content) {
        throw new Error("File is empty or not found");
    }

    const words = content.split(" ");
    for (let i = 0; i < words.length; i++) {
        const half = Math.floor(words[i].length / 2);
        const firstHalf = words[i].slice(0, half);
        const secondHalf = words[i].slice(half);
        words[i] = secondHalf + firstHalf;
    }
    const output = words.join(" ");
    console.log(output);
} catch (error) {
    console.error("Error reading file:", error.message);
}

