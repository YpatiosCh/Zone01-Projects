import { readdir } from "fs/promises";
import { readFile } from "fs/promises";

try {
    const dirPath = process.argv[2];
    if (!dirPath) {
        dirPath = "./"
    }

    const files = await readdir(dirPath);
    if (files.length === 0) {
        throw new Error("Directory is empty or not found");
    }

    console.log(files.length);
} catch (error) {
    console.error("Error reading directory:", error.message);
}