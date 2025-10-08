export const generateLetters = () => {
    // get the body to later append the letters
    const body = document.querySelector('body');
    // create uppercase alphabet 
    const alphabet = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ';

    // calculate font size increment from 11px to 130px over 120 letters
    const minFontSize = 11;
    const maxFontSize = 130;
    const fontSizeIncrement = (maxFontSize - minFontSize) / (120 - 1);

    // create the 120 letters
    for (let i = 0; i < 120; i++) {
        // create the div element
        const div = document.createElement('div');

        // pick random letter from alphabet
        const randomLetter = alphabet[Math.floor(Math.random() * alphabet.length)];
        // add it to the created div
        div.textContent = randomLetter;

        // calculate font size for the current div
        const fontSize = minFontSize + (i * fontSizeIncrement);
        div.style.fontSize = `${fontSize}px`;

        // set font weight based on thirds
        let fontWeight;
        if (i < 40) {
            fontWeight = 300;
        } else if (i < 80) {
            fontWeight = 400;
        } else {
            fontWeight = 600;
        }
        div.style.fontWeight = fontWeight;

        // append to body
        body.append(div);
    }
};