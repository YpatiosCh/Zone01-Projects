/**
 * Creates a retry function that attempts to call a callback with retries
 * @param {number} count - Maximum number of retries
 * @param {Function} callback - Async function to retry
 * @returns {Function} - Function that retries the callback
 */
function retry(count, callback) {
    return async function (...args) {
        let attempts = 0;
        let lastError;
        
        while (attempts <= count) {
            try {
                return await callback(...args);
            } catch (error) {
                lastError = error;
                attempts++;
            }
        }
        
        throw lastError;
    };
}

/**
 * Creates a timeout function that limits execution time of a callback
 * @param {number} delay - Maximum wait time in milliseconds
 * @param {Function} callback - Async function to execute with timeout
 * @returns {Function} - Function that executes callback with timeout
 */
function timeout(delay, callback) {
    return async function (...args) {
        const timeoutPromise = new Promise((resolve, reject) => {
            setTimeout(() => {
                reject(new Error('timeout'));
            }, delay);
        });
        
        const callbackPromise = callback(...args);
        
        return Promise.race([callbackPromise, timeoutPromise]);
    };
}