function isFriday(date) {
    return new Date(date).getDay() === 5;
}

function isWeekend(date) {
    let day = new Date(date).getDay();

    return day === 0 || day === 6;
}

function isLeapYear(date) {
    let year = new Date(date).getFullYear();

    return (year % 4 === 0 && year % 100 !== 0) || year % 400 === 0;
}

function isLastDayOfMonth(date) {
  let givenDate = new Date(date);

  let nextDay = new Date(givenDate);
  nextDay.setDate(givenDate.getDate() + 1);

  return nextDay.getMonth() !== givenDate.getMonth();
}