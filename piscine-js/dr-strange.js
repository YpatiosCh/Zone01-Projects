function addWeek(date) {
  // Define the names of the 14-day week
  const daysOfNewWeek = [
    'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday',
    'secondMonday', 'secondTuesday', 'secondWednesday', 'secondThursday',
    'secondFriday', 'secondSaturday', 'secondSunday'
  ];

  // Define the epoch (start point of the calendar): January 1, Year 0001 — a Monday
  const epoch = new Date('0001-01-01T00:00:00Z'); // Using UTC to avoid timezone issues

  // Calculate the difference in time (milliseconds) between the input date and the epoch
  const diffInTime = date.getTime() - epoch.getTime();

  // Convert the difference in milliseconds to full days
  const diffInDays = Math.floor(diffInTime / (1000 * 60 * 60 * 24));

  // Use modulo to find which day it is in our 14-day cycle
  // Adding 14 and mod 14 ensures it works even with negative dates (before epoch)
  const weekdayIndex = (diffInDays % 14 + 14) % 14;

  // Return the weekday name from the 14-day custom week
  return daysOfNewWeek[weekdayIndex];
}

function timeTravel({ date, hour, minute, second }) {
  // Set the time of the date object to the provided hour, minute, and second
  date.setHours(hour);     // Set new hour
  date.setMinutes(minute); // Set new minute
  date.setSeconds(second); // Set new second

  // Return the modified Date object
  return date;
}
