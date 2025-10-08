export const compose = () => {
    // add keydown event listener to the document
    document.addEventListener('keydown', (event) => {
        // get the key
        const key = event.key;

        // check if its a lowercase letter
        if (key >= 'a' && key <= 'z') {
            // create a new note div
            const noteDiv = document.createElement('div');
            noteDiv.className = 'note';
            noteDiv.textContent = key;

            // generate unique backround color based on the key
            const backroundColor = generateColor(key);
            noteDiv.style.backgroundColor = backroundColor;

            // add the note to the body
            document.body.appendChild(noteDiv);
        } else if (key === 'Backspace') { // handle backspace
            // get all notes 
            const allNotes = document.querySelectorAll('.note');
            // get the last note
            const lastNote = allNotes[allNotes.length -1];
            // remove the last note
            lastNote.remove();
        } else if (key === 'Escape') { // handle escape
            // get all notes
            const allNotes = document.querySelectorAll('.note');
            // remove all notes
            allNotes.forEach(note => note.remove());
        }
    });
};

const generateColor = (key) => {
    // convert key to a number 
    const keyNumber = key.charCodeAt(0) - 96; // 'a' is 97, so 'a' becomes 1

    // Generate HSL color with unique hue based on the key
    const hue = (keyNumber * 360 / 26) % 360; // Distribute across color wheel
    const saturation = 70; // Fixed saturation for vibrant colors
    const lightness = 50; // Fixed lightness
    
    return `hsl(${hue}, ${saturation}%, ${lightness}%)`;
}