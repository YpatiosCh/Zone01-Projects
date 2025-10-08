import { readdir } from "fs/promises";

try {
    const dirPath = process.argv[2] || "./";
    const files = await readdir(dirPath);

    if (files.length === 0) {
        throw new Error("Directory is empty or not found");
    }

    // Only keep .json files and map them to objects with extracted names
    const guests = files
        .filter(file => file.endsWith('.json'))
        .map(file => {
            const [first, last] = file.replace('.json', '').split('_');
            return { file, first, last };
        });

    // Sort by last name
    guests.sort((a, b) => a.last.localeCompare(b.last));

    // Print in "Last First" format
    guests.forEach((guest, index) => {
        console.log(`${index + 1}. ${guest.last} ${guest.first}`);
    });

} catch (error) {
    console.error("Error reading directory:", error.message);
}
