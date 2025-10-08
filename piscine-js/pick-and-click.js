export const pick = () => {
    // create the hsl display element 
    const hslDiv = document.createElement('div');
    hslDiv.className = 'hsl text';
    hslDiv.style.position = 'fixed';
    hslDiv.style.top = '50%';
    hslDiv.style.left = '50%';
    hslDiv.style.transform = 'translate(-50%, -50%)';
    document.body.appendChild(hslDiv);

    // create the hue display element 
    const hueDiv = document.createElement('div');
    hueDiv.className = 'hue text';
    document.body.appendChild(hueDiv);

    // create the luminosity display element
    const luminosityDiv = document.createElement('div');
    luminosityDiv.className = 'luminosity text';
    document.body.appendChild(luminosityDiv);

    // create svg for crosshairs
    const svg = document.createElementNS('http://www.w3.org/2000/svg', 'svg');
    document.body.appendChild(svg);

    // create X axis line (vertical)
    const axisX = document.createElementNS('http://www.w3.org/2000/svg', 'line');
    axisX.id = 'axisX';
    axisX.setAttribute('y1', '0');
    axisX.setAttribute('y2', '100vh');
    svg.appendChild(axisX);

    // Create Y axis line (horizontal)
    const axisY = document.createElementNS('http://www.w3.org/2000/svg', 'line');
    axisY.id = 'axisY';
    axisY.setAttribute('x1', '0');
    axisY.setAttribute('x2', '100vw');
    svg.appendChild(axisY);

    // Current HSL values
    let currentHSL = '';

    // Mouse move event handler
    const handleMouseMove = (event) => {
        const x = event.clientX;
        const y = event.clientY;
        
        // Calculate hue based on X position (0-360)
        const hue = Math.round((x / window.innerWidth) * 360);
        
        // Calculate luminosity based on Y position (0-100)
        // Invert Y so top is bright (100%) and bottom is dark (0%)
        const luminosity = Math.round((1 - y / window.innerHeight) * 100);
        
        // Fixed saturation at 50%
        const saturation = 50;
        
        // Create HSL string
        currentHSL = `hsl(${hue}, ${saturation}%, ${luminosity}%)`;
        
        // Update background color
        document.body.style.background = currentHSL;
        
        // Update display elements
        hslDiv.textContent = currentHSL;
        hueDiv.textContent = hue;
        luminosityDiv.textContent = luminosity;
        
        // Update crosshair lines
        axisX.setAttribute('x1', x);
        axisX.setAttribute('x2', x);
        axisY.setAttribute('y1', y);
        axisY.setAttribute('y2', y);
    };

    // Click event handler to copy to clipboard
    const handleClick = async () => {
        try {
            await navigator.clipboard.writeText(currentHSL);
            console.log('HSL value copied to clipboard:', currentHSL);
        } catch (err) {
            console.error('Failed to copy to clipboard:', err);
        }
    };

    // Add event listeners
    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('click', handleClick);

    // Initialize with center position
    const centerX = window.innerWidth / 2;
    const centerY = window.innerHeight / 2;
    handleMouseMove({ clientX: centerX, clientY: centerY });
}