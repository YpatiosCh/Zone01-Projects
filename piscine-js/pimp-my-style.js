import { styles } from './pimp-my-style.data.js'

let currentIndex = 0;
let isRemoving = false;

export const pimp = () => {
    // get the button 
    const button = document.querySelector('button');

    if (!isRemoving) {
        // adding phase 
        if (currentIndex < styles.length) {
            // add current style class
            button.classList.add(styles[currentIndex]);
            currentIndex++;

            // if we have added all the classes, switch to removing phase
            if (currentIndex === styles.length) {
                isRemoving = true;
                button.classList.add('unpimp');
            }
        }
    } else {
        // removing phase
        if (currentIndex > 0) {
            currentIndex--;
            // remove the class at current index
            button.classList.remove(styles[currentIndex]);


            // if we have removed all classes, switch back to adding mode
            if (currentIndex === 0) {
                isRemoving = false;
                button.classList.remove('unpimp');
            }
        }
    }
};