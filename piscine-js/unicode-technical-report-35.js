function format(date, formatString) {
    // Helper function to pad numbers with leading zeros
    const pad = (num, length) => String(num).padStart(length, '0');
    
    // Month names for formatting
    const monthNames = [
        'January', 'February', 'March', 'April', 'May', 'June',
        'July', 'August', 'September', 'October', 'November', 'December'
    ];
    
    const monthNamesShort = [
        'Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun',
        'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'
    ];
    
    // Weekday names for formatting
    const weekdayNames = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];
    const weekdayNamesShort = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat'];
    
    // Extract date components
    const year = date.getFullYear();
    const month = date.getMonth() + 1; // getMonth() returns 0-11
    const day = date.getDate();
    const weekday = date.getDay(); // 0 = Sunday, 1 = Monday, etc.
    const hours24 = date.getHours();
    const hours12 = hours24 === 0 ? 12 : hours24 > 12 ? hours24 - 12 : hours24;
    const minutes = date.getMinutes();
    const seconds = date.getSeconds();
    const ampm = hours24 < 12 ? 'AM' : 'PM';
    
    // Create a map of format tokens to their replacements
    const formatMap = {
        // Year
        'yyyy': pad(Math.abs(year), 4),
        'y': String(Math.abs(year)),
        
        // Era
        'GGGG': year > 0 ? 'Anno Domini' : 'Before Christ',
        'G': year > 0 ? 'AD' : 'BC',
        
        // Month
        'MMMM': monthNames[month - 1],
        'MMM': monthNamesShort[month - 1],
        'MM': pad(month, 2),
        'M': String(month),
        
        // Day
        'dd': pad(day, 2),
        'd': String(day),
        
        // Weekday
        'EEEE': weekdayNames[weekday],
        'E': weekdayNamesShort[weekday],
        
        // Hour (12-hour)
        'hh': pad(hours12, 2),
        'h': String(hours12),
        
        // Hour (24-hour)
        'HH': pad(hours24, 2),
        'H': String(hours24),
        
        // Minute
        'mm': pad(minutes, 2),
        'm': String(minutes),
        
        // Second
        'ss': pad(seconds, 2),
        's': String(seconds),
        
        // AM/PM
        'a': ampm
    };
    
    // Replace format tokens in the string
    // We need to be careful not to replace tokens that are part of larger strings
    // Process the string character by character to identify actual tokens
    let result = '';
    let i = 0;
    
    while (i < formatString.length) {
        let tokenFound = false;
        
        // Check for tokens starting from longest to shortest
        const sortedTokens = Object.keys(formatMap).sort((a, b) => b.length - a.length);
        
        for (const token of sortedTokens) {
            if (formatString.substring(i, i + token.length) === token) {
                result += formatMap[token];
                i += token.length;
                tokenFound = true;
                break;
            }
        }
        
        if (!tokenFound) {
            result += formatString[i];
            i++;
        }
    }
    
    return result;
}