import { readdir, readFile, writeFile } from "fs/promises";
import { join } from "path";

const outputFile = "vip.txt";

try {
    const dirPath = process.argv[2] || ".";

    // Get all files in directory
    const files = await readdir(dirPath);

    // Filter .json files and read their contents
    const guests = [];

    for (const file of files) {
        if (!file.endsWith(".json")) continue;

        const [first, last] = file.replace(".json", "").split("_");
        const filePath = join(dirPath, file);

        try {
            const content = await readFile(filePath, "utf8");
            const data = JSON.parse(content);

            if (data.answer === "yes") {
                guests.push({ first, last });
            }
        } catch (err) {
            // Skip malformed JSON or unreadable files
            continue;
        }
    }

    // Sort guests alphabetically by last name
    guests.sort((a, b) => a.last.localeCompare(b.last));

    // Format the output
    const lines = guests.map((g, i) => `${i + 1}. ${g.last} ${g.first}`);

    // Always write vip.txt, even if empty
    await writeFile(outputFile, lines.join("\n"));

} catch (error) {
    console.error("Error:", error.message);
}
