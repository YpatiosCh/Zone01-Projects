function filterEntries(obj, callback) {
  return Object.fromEntries(
    Object.entries(obj).filter(entry => callback(entry))
  );
}

function mapEntries(obj, callback) {
  return Object.fromEntries(
    Object.entries(obj).map(entry => callback(entry))
  );
}

function reduceEntries(obj, callback, initial) {
  return Object.entries(obj).reduce(callback, initial);
}


function totalCalories(cart) {
  const total = reduceEntries(
    cart,
    (acc, [item, grams]) => {
      const per100g = nutritionDB[item]?.calories || 0;
      return acc + (per100g * grams) / 100;
    },
    0
  );
  return Math.round(total * 10) / 10; // round to 1 decimal place
}


function lowCarbs(cart) {
  // Keep only items with total carbs < 50g
  return filterEntries(
    cart,
    ([item, grams]) => {
      const carbsPer100g = nutritionDB[item]?.carbs || 0;
      return (carbsPer100g * grams) / 100 < 50;
    }
  );
}

function cartTotal(cart) {
  return mapEntries(
    cart,
    ([item, grams]) => {
      const nutrition = nutritionDB[item] || {};
      const scaledNutrition = {};
      for (const [key, value] of Object.entries(nutrition)) {
        // Multiply and round to 3 decimals
        scaledNutrition[key] = Math.round((value * grams) * 1000 / 100) / 1000;
      }
      return [item, scaledNutrition];
    }
  );
}