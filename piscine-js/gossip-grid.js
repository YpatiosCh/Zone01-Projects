import { gossips } from './gossip-grid.data.js';

export const grid = () => {
    // Create ranges
    ranges();
    
    // Create form card
    const form = document.createElement('form');
    form.className = 'gossip';
    
    const textarea = document.createElement('textarea');
    const button = document.createElement('button');
    button.textContent = 'Share gossip!';
    button.type = 'submit';
    
    button.addEventListener('click', (e) => {
        e.preventDefault();
        const newGossip = textarea.value.trim();
        if (newGossip) {
            // Add to beginning of gossips array
            gossips.unshift(newGossip);
            
            // Remove all existing gossip cards except the form
            document.querySelectorAll('.gossip').forEach((card, i) => {
                if (i > 0) card.remove();
            });
            
            // Clear textarea
            textarea.value = '';
            
            // Re-render all gossips
            renderGossips();
        }
    });
    
    form.appendChild(textarea);
    form.appendChild(button);
    document.body.appendChild(form);
    
    // Render initial gossips
    renderGossips();
};

function renderGossips() {
    gossips.forEach(gossip => {
        const div = document.createElement('div');
        div.className = 'gossip';
        div.textContent = gossip;
        document.body.appendChild(div);
    });
}

function ranges() {
    const rangesDiv = document.createElement('div');
    rangesDiv.className = 'ranges';
    
    // Width range
    const widthRange = document.createElement('input');
    widthRange.type = 'range';
    widthRange.id = 'width';
    widthRange.min = '200';
    widthRange.max = '800';
    widthRange.value = '250';
    widthRange.addEventListener('input', (e) => {
        document.querySelectorAll('.gossip').forEach(card => {
            card.style.width = e.target.value + 'px';
        });
    });
    
    // Font size range
    const fontSizeRange = document.createElement('input');
    fontSizeRange.type = 'range';
    fontSizeRange.id = 'fontSize';
    fontSizeRange.min = '20';
    fontSizeRange.max = '40';
    fontSizeRange.value = '20';
    fontSizeRange.addEventListener('input', (e) => {
        document.querySelectorAll('.gossip').forEach(card => {
            card.style.fontSize = e.target.value + 'px';
        });
    });
    
    // Background range
    const backgroundRange = document.createElement('input');
    backgroundRange.type = 'range';
    backgroundRange.id = 'background';
    backgroundRange.min = '20';
    backgroundRange.max = '75';
    backgroundRange.value = '50';
    backgroundRange.addEventListener('input', (e) => {
        document.querySelectorAll('.gossip').forEach(card => {
            card.style.background = `hsl(280, 50%, ${e.target.value}%)`;
        });
    });
    
    rangesDiv.appendChild(widthRange);
    rangesDiv.appendChild(fontSizeRange);
    rangesDiv.appendChild(backgroundRange);
    document.body.appendChild(rangesDiv);
}