function RNA(dna) {
  let result = '';
  for (let i = 0; i < dna.length; i++) {
    const base = dna[i];
    if (base === 'G') result += 'C';
    else if (base === 'C') result += 'G';
    else if (base === 'T') result += 'A';
    else if (base === 'A') result += 'U';
  }
  return result;
}

function DNA(rna) {
  let result = '';
  for (let i = 0; i < rna.length; i++) {
    const base = rna[i];
    if (base === 'C') result += 'G';
    else if (base === 'G') result += 'C';
    else if (base === 'A') result += 'T';
    else if (base === 'U') result += 'A';
  }
  return result;
}
