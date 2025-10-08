/**
 * Interpolation function that generates points from start to end position
 * @param {Object} config - Configuration object
 * @param {number} config.step - Number of interpolation steps
 * @param {number} config.start - Starting position
 * @param {number} config.end - Ending position (not included)
 * @param {Function} config.callback - Function to call with each point [x, y]
 * @param {number} config.duration - Total duration for all steps
 */
function interpolation({
    step = 0,
    start = 0,
    end = 0,
    callback = () => {},
    duration = 0,
} = {}) {
    const delta = (end - start) / step;
    let current = start;
    let i = 0;
    const timer = setInterval(() => {
        if (i < step) {
            callback([current, (duration / step) * (i + 1)]);
            current += delta;
            i++;
        } else {
            clearInterval(timer);
        }
    }, duration / step);
}