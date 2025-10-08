let lastCircle = null;
let box = null;

// Add event listeners to the document body
document.body.addEventListener('click', (e) => {
    createCircle(e);
});

document.body.addEventListener('mousemove', (e) => {
    moveCircle(e);
});

export const createCircle = (event) => {
    if (!event) return;
    
    // create a new circle div
    const circle = document.createElement('div');
    circle.className = 'circle';
    circle.style.background = 'white';
    circle.style.position = 'absolute';
    circle.style.zIndex = '1000';

    // position the circle at the mouse click position
    // subtract half the circle size to center it on the cursor
    circle.style.left = `${event.clientX - 25}px`;
    circle.style.top = `${event.clientY - 25}px`;

    // add the circle to the body
    document.body.appendChild(circle);

    // update the last circle reference
    lastCircle = circle;
};

export const moveCircle = (event) => {
    if (!event || !lastCircle) return;
    
    const newLeft = event.clientX - 25;
    const newTop = event.clientY - 25;
    
    // check if the circle is trapped
    if (isCircleTrapped(lastCircle)) {
        // If trapped, only allow movement within the box
        const constrainedPosition = constrainToBox(newLeft, newTop);
        lastCircle.style.left = `${constrainedPosition.x}px`;
        lastCircle.style.top = `${constrainedPosition.y}px`;
    } else {
        // If not trapped, move freely
        lastCircle.style.left = `${newLeft}px`;
        lastCircle.style.top = `${newTop}px`;
        
        // check if the circle is now inside the box
        if (isCircleInsideBox(lastCircle)) {
            // turn purple and trap it
            lastCircle.style.background = `var(--purple)`;
            lastCircle.dataset.trapped = 'true';
        }
    }
};

export const setBox = () => {
    // create the box element 
    box = document.createElement('div');
    box.className = 'box';

    // center the box on the page
    box.style.position = 'absolute';
    box.style.left = '50%';
    box.style.top = '50%';
    box.style.transform = 'translate(-50%, -50%)';

    // add the box to the body
    document.body.appendChild(box);
};

// Helper function to constrain circle position to stay within the box
const constrainToBox = (x, y) => {
    if (!box) return { x, y };
    
    const boxRect = box.getBoundingClientRect();
    const circleRadius = 25;
    
    // Calculate box boundaries (account for 1px border)
    const leftEdge = boxRect.left + 1;
    const rightEdge = boxRect.right - 1;
    const topEdge = boxRect.top + 1;
    const bottomEdge = boxRect.bottom - 1;
    
    // Constrain x position
    const minX = leftEdge;
    const maxX = rightEdge - 50; // subtract full circle width
    const constrainedX = Math.max(minX, Math.min(maxX, x));
    
    // Constrain y position
    const minY = topEdge;
    const maxY = bottomEdge - 50; // subtract full circle height
    const constrainedY = Math.max(minY, Math.min(maxY, y));
    
    return { x: constrainedX, y: constrainedY };
};

// Helper function to check if a circle is entirely inside the box
const isCircleInsideBox = (circle) => {
    if (!box) return false;
    
    const circleRect = circle.getBoundingClientRect();
    const boxRect = box.getBoundingClientRect();
    
    // Circle is 50px diameter, so radius is 25px
    const circleRadius = 25;
    
    // Check if the entire circle (including its radius) is inside the box
    // Account for the 1px border by adding 1px margin
    const leftEdge = boxRect.left + 1;
    const rightEdge = boxRect.right - 1;
    const topEdge = boxRect.top + 1;
    const bottomEdge = boxRect.bottom - 1;
    
    const circleCenterX = circleRect.left + circleRadius;
    const circleCenterY = circleRect.top + circleRadius;
    
    return (
        circleCenterX - circleRadius >= leftEdge &&
        circleCenterX + circleRadius <= rightEdge &&
        circleCenterY - circleRadius >= topEdge &&
        circleCenterY + circleRadius <= bottomEdge
    );
};

// Helper function to check if a circle is trapped (already purple)
const isCircleTrapped = (circle) => {
    return circle.dataset.trapped === 'true';
};