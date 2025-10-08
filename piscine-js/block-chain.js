
/**
 * Creates a block in the blockchain.
 * @param {any} data - The data to store in the block (any valid JSON).
 * @param {object} prev - The previous block in the chain.
 *                         Defaults to the genesis block if not provided.
 * @returns {object} The newly created block with index, hash, data, prev, and chain method.
 */
function blockChain(data, prev = { index: 0, hash: '0' }) {
  // Calculate the current block's index by adding 1 to the previous block's index.
  const index = prev.index + 1;

  // Create the hash input string by concatenating:
  // - The current block index (number)
  // - The previous block's hash (string)
  // - The stringified JSON data of this block
  const hashInput = index + prev.hash + JSON.stringify(data);

  // Calculate the hash for the current block using the provided hashCode function.
  const hash = hashCode(hashInput);

  // Return the new block object with all required properties:
  return {
    index,      // The block’s position in the chain
    hash,       // The unique hash of this block
    data,       // The stored data in the block
    prev,       // Reference to the previous block in the chain

    // The chain method allows creating the next block by calling:
    // currentBlock.chain(nextData)
    // It calls blockChain recursively, passing the nextData as new data,
    // and 'this' (the current block) as the previous block.
    chain(nextData) {
      return blockChain(nextData, this);
    }
  };
}


/*

// hashCode takes a string and returns a short hash string.
// It works by processing each character in the string and combining their char codes
// into a single number using bitwise operations, then converts that number to base 36.

const hashCode = str =>
  (
    // Convert the string into an array of characters.
    [...str]
    // Use reduce to process each character (c) and accumulate a hash number (h).
    .reduce((h, c) => 
      // Shift h 5 bits left (multiply by 32), subtract h, add the char code of c.
      // Then bitwise AND with h to keep the number in a 32-bit range.
      (h = (h << 5) - h + c.charCodeAt(0)) & h, 0)
  )
  // Convert to an unsigned 32-bit integer.
  >>> 0
  // Convert the number to a base-36 string (numbers + letters).
  .toString(36);


*/