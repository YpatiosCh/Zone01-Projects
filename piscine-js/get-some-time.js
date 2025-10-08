function firstDayWeek(weekNumber, year) {
  // Parse the year string to a number and get Jan 1st of that year
  const jan1 = new Date(`${year}-01-01`);

  // Get the day of the week (0 = Sunday, 1 = Monday, ..., 6 = Saturday)
  let dayOfWeek = jan1.getDay();

  // Adjust so that Monday is considered 0 (so we can shift backwards to Monday)
  // JavaScript: Sunday is 0, so we map it to 6, and shift others accordingly
  const shift = (dayOfWeek + 6) % 7;

  // Get the Monday of Week 1 (could be in previous year)
  const week1Monday = new Date(jan1);
  week1Monday.setDate(jan1.getDate() - shift);

  // Now calculate the Monday of the desired week
  const desiredMonday = new Date(week1Monday);
  desiredMonday.setDate(week1Monday.getDate() + (weekNumber - 1) * 7);

  // If the resulting Monday is before Jan 1st of that year, return 01-01-[year]
  if (desiredMonday < jan1) {
    return `01-01-${year}`;
  }

  // Format date as dd-mm-yyyy
  const dd = String(desiredMonday.getDate()).padStart(2, '0');
  const mm = String(desiredMonday.getMonth() + 1).padStart(2, '0'); // months are 0-indexed
  const yyyy = String(desiredMonday.getFullYear()).padStart(4, '0');

  return `${dd}-${mm}-${yyyy}`;
}
