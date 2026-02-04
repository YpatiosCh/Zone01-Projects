function updateGlowEffect() {
    // Get the user input from the textarea and trim any extra spaces
    const inputText = document.getElementById('text').value.trim();

    // If the input is empty, reset the highlighted text
    if (!inputText) {
        const asciiDiv = document.querySelector('.ascii-background');
        asciiDiv.innerHTML = asciiDiv.textContent;
        return;
    }

    // Get the ascii-background div
    const asciiDiv = document.querySelector('.ascii-background');

    // Get the content of the ascii-background div
    const originalText = asciiDiv.textContent;

    // Find the start and end of the "MUST use these:" section
    const startPattern = 'Yes, yes, yes, MUST use these:';
    const endPattern = 'No, no, no, MUST NOT use this one either:';

    const startIndex = originalText.indexOf(startPattern);
    const endIndex = originalText.indexOf(endPattern);

    if (startIndex === -1 || endIndex === -1) {
        console.error('MUST use these section not found.');
        return;
    }

    // Extract the part to highlight
    const partToHighlight = originalText.substring(startIndex + startPattern.length, endIndex).trim();

    // Escape the input for regex
    const escapedInputText = inputText.replace(/[.*+?^=!:${}()|\[\]\/\\]/g, "\\$&");

    // Create a regex pattern
    const regex = new RegExp(`[${escapedInputText}]`, 'gi');

    // Highlight matching characters
    const highlightedPart = partToHighlight.replace(regex, match => {
        return `<span class="glow-text">${match}</span>`;
    });

    // Update the content with the highlighted part
    const updatedText = originalText.substring(0, startIndex + startPattern.length) +
                        highlightedPart +
                        originalText.substring(endIndex);

    // Use innerHTML to apply the update
    asciiDiv.innerHTML = updatedText;
}

// Attach the event listener to the textarea
document.getElementById('text').addEventListener('input', updateGlowEffect);