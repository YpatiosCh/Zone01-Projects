function dogYears(planet, ageInSeconds) {
  const orbitalPeriods = {
    earth: 1.0,
    mercury: 0.2408467,
    venus: 0.61519726,
    mars: 1.8808158,
    jupiter: 11.862615,
    saturn: 29.447498,
    uranus: 84.016846,
    neptune: 164.79132,
  };

  const SECONDS_IN_EARTH_YEAR = 31557600; // 365.25 days

  const humanYears = ageInSeconds / (SECONDS_IN_EARTH_YEAR * orbitalPeriods[planet.toLowerCase()]);
  const dogYears = humanYears * 7;

  return Number(dogYears.toFixed(2));
}
