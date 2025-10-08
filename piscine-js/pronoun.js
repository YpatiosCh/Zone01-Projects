function pronoun(str) {
  const PRONOUNS = ['i', 'you', 'he', 'she', 'it', 'they', 'we'];
  const result = {};
  const words = str.toLowerCase().match(/\b\w+\b/g) || [];

  for (let i = 0; i < words.length; i++) {
    const current = words[i];

    if (PRONOUNS.includes(current)) {
      if (!result[current]) {
        result[current] = { word: [], count: 0 };
      }
      result[current].count++;

      // Find the next non-pronoun word
      let j = i + 1;
      while (j < words.length && PRONOUNS.includes(words[j])) {
        j++;
      }

      // Only add the next word if it is *immediately* after current (no pronouns between)
      if (j === i + 1 && j < words.length) {
        // next word immediately follows current pronoun and is not a pronoun
        result[current].word.push(words[j]);
      }
    }
  }

  return result;
}
