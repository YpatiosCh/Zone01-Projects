function matchCron(cronString, date) {
  const [minC, hourC, dayC, monthC, weekC] = cronString.split(' ');

  // Extract components from the date
  const minutes = date.getMinutes();
  const hours = date.getHours();
  const dayOfMonth = date.getDate();
  const month = date.getMonth() + 1; // JS: 0-11, Cron: 1-12
  const dayOfWeek = ((date.getDay() + 6) % 7) + 1; // JS: 0 (Sun) to 6 (Sat) → Cron: 1 (Mon) to 7 (Sun)

  // Helper function to check match
  const match = (cronVal, actual) =>
    cronVal === '*' || Number(cronVal) === actual;

  // Compare each cron field with the corresponding date value
  return (
    match(minC, minutes) &&
    match(hourC, hours) &&
    match(dayC, dayOfMonth) &&
    match(monthC, month) &&
    match(weekC, dayOfWeek)
  );
}
