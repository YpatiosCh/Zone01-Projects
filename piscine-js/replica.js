function replica(target, ...sources) {
  for (const source of sources) {
    for (const key in source) {
      if (!source.hasOwnProperty(key)) continue;

      const sourceValue = source[key];
      const targetValue = target[key];

      // If both values are objects, recurse
      if (
        sourceValue &&
        typeof sourceValue === 'object' &&
        !Array.isArray(sourceValue) &&
        targetValue &&
        typeof targetValue === 'object' &&
        !Array.isArray(targetValue)
      ) {
        replica(targetValue, sourceValue);
      }
      // If array or primitive, replace
      else if (Array.isArray(sourceValue)) {
        // Deep copy the array
        target[key] = sourceValue.map((item) =>
          typeof item === 'object' && item !== null ? replica({}, item) : item
        );
      } else {
        // Primitive or function or null
        target[key] = sourceValue;
      }
    }
  }
  return target;
}
