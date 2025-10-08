import { colors } from './fifty-shades-of-cold.data.js';

export const generateClasses = () => {
    // Create a style element
    const style = document.createElement('style');
    
    // Generate CSS classes for each color
    let cssRules = '';
    colors.forEach(color => {
        cssRules += `.${color} { background: ${color}; }\n`;
    });
    
    // Debug: log the number of colors and first few rules
    console.log(`Generated ${colors.length} color classes`);
    console.log('First few CSS rules:', cssRules.slice(0, 200));
    
    // Set the CSS content
    style.textContent = cssRules;
    
    // Add the style element to the head
    document.head.appendChild(style);
    
    // Debug: verify the style was added
    console.log('Style element added to head');
};

export const generateColdShades = () => {
    const body = document.querySelector('body');
    
    // Define cold color keywords
    const coldKeywords = ['aqua', 'blue', 'turquoise', 'green', 'cyan', 'navy', 'purple'];
    
    // Filter colors that contain any of the cold keywords
    const coldColors = colors.filter(color => {
        return coldKeywords.some(keyword => color.includes(keyword));
    });
    
    // Create a div for each cold color
    coldColors.forEach(color => {
        const div = document.createElement('div');
        div.className = color;
        div.textContent = color;
        body.appendChild(div);
    });
};

export const choseShade = (clickedShade) => {
    // Get all divs
    const allDivs = document.querySelectorAll('div');
    
    // Replace the class of each div with the clicked shade
    allDivs.forEach(div => {
        // Find the current color class and replace it with the clicked shade
        colors.forEach(color => {
            if (div.classList.contains(color)) {
                div.classList.replace(color, clickedShade);
            }
        });
    });
};