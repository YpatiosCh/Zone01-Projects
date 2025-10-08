function sunnySunday(date) {
    // Reference date: January 1, 0001 is a Monday (day 0 in our 6-day system)
    const referenceDate = new Date(1, 0, 1); // January 1, 0001
    
    // Calculate the difference in days between the input date and reference date
    const timeDiff = date.getTime() - referenceDate.getTime();
    const daysDiff = Math.floor(timeDiff / (1000 * 60 * 60 * 24));
    
    // In our 6-day week system: Monday=0, Tuesday=1, Wednesday=2, Thursday=3, Friday=4, Saturday=5
    const dayIndex = ((daysDiff % 6) + 6) % 6; // Handle negative numbers properly
    
    // Map the day index to weekday names
    const weekdays = ["Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"];
    
    return weekdays[dayIndex];
}