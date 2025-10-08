function slice(value, start = 0, end = value.length) {
  const result = [];

  // Handle negative indexes
  const len = value.length;
  if (start < 0) start = Math.max(len + start, 0);
  if (end < 0) end = len + end;

  // Ensure end doesn't exceed length
  end = Math.min(end, len);

  for (let i = start; i < end; i++) {
    result.push(value[i]);
  }

  // Return same type as input
  return typeof value === 'string' ? result.join('') : result;
}
