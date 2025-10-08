import fs from "fs/promises";

try {

    // Get command line arguments
    const [inputFile, mode, outputFile] = process.argv.slice(2);

    // Validate input
    if (!inputFile || !mode || !['encode', 'decode'].includes(mode)) {
        throw new Error("Usage: node tell-me-vip.mjs <inputFile> <encode|decode> <outputFile>");
    }

    // Process content based on mode
    if (mode === 'encode') {

        // Read the input file
        const content = await fs.readFile(inputFile);
        if (!content) {
            throw new Error("Input file is empty or not found");
        }
        // encode to base64
        const encodedContent = content.toString('base64');
        await fs.writeFile(outputFile || "cypher.txt", encodedContent, 'utf8');
    } else if (mode === 'decode') {
        // Read the input file as UTF-8
        let content = await fs.readFile(inputFile, 'utf8');
        if (!content) {
            throw new Error("Input file is empty or not found");
        }
        // remove any trailing newline/whitespaces characters
        content = content.replace(/\s+/g, '');

        // convert base64 string to Buffer
        const buffer = Buffer.from(content, 'base64');
        await fs.writeFile(outputFile || "clear.txt", buffer);
    }

} catch (error) {
    console.error("Error:", error.message);
}


