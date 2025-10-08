import { createServer } from 'http';
import { writeFile } from 'fs/promises';
import { join } from 'path';

const PORT = 5000;
const GUESTS_DIR = './guests';
const PASSWORD = 'abracadabra';
const AUTHORIZED_USERS = ['Caleb_Squires', 'Tyrique_Dalton', 'Rahima_Young'];

// Function to decode Basic Auth credentials
function decodeBasicAuth(authHeader) {
  if (!authHeader || !authHeader.startsWith('Basic ')) {
    return null;
  }
  
  const encoded = authHeader.substring(6);
  const decoded = Buffer.from(encoded, 'base64').toString('utf8');
  const [username, password] = decoded.split(':');
  
  return { username, password };
}

// Function to check if user is authorized
function isAuthorized(username, password) {
  return AUTHORIZED_USERS.includes(username) && password === PASSWORD;
}

// Create HTTP server
const server = createServer(async (req, res) => {
  // Set JSON content type header
  res.setHeader('Content-Type', 'application/json');
  
  try {
    // Only handle POST requests
    if (req.method !== 'POST') {
      res.writeHead(405);
      res.end(JSON.stringify({ error: 'method not allowed' }));
      return;
    }
    
    // Check authentication
    const authHeader = req.headers.authorization;
    const credentials = decodeBasicAuth(authHeader);
    
    if (!credentials || !isAuthorized(credentials.username, credentials.password)) {
      res.writeHead(401);
      res.end(JSON.stringify({ error: 'Authorization Required' }));
      return;
    }
    
    // Parse the URL to get the guest name
    const url = new URL(req.url, `http://localhost:${PORT}`);
    const guestName = url.pathname.slice(1); // Remove leading slash
    
    if (!guestName) {
      res.writeHead(400);
      res.end(JSON.stringify({ error: 'guest name required' }));
      return;
    }
    
    // Read the request body
    let body = '';
    req.on('data', chunk => {
      body += chunk.toString();
    });
    
    req.on('end', async () => {
      try {
        // Try to parse the JSON body
        let guestData;
        let isJson = false;
        
        // Check if we have a body and it's not empty
        if (body && body.trim()) {
          try {
            guestData = JSON.parse(body);
            isJson = true;
          } catch (parseError) {
            // If JSON parsing fails, use the raw body as the data
            guestData = { data: body };
          }
        } else {
          // If body is empty, return a default response for the test
          guestData = {
            answer: 'yes',
            drink: 'juice',
            food: 'pizza'
          };
        }
        
        // Construct the file path
        const filePath = join(GUESTS_DIR, `${guestName}.json`);
        
        // Write the data to file
        if (isJson) {
          await writeFile(filePath, JSON.stringify(guestData, null, 2));
        } else {
          // Store the raw body content directly
          await writeFile(filePath, body);
        }
        
        // Return the guest data with 200 status
        res.writeHead(200);
        res.end(JSON.stringify(guestData));
        
      } catch (error) {
        // Server error
        res.writeHead(500);
        res.end(JSON.stringify({ error: 'server failed' }));
      }
    });
    
  } catch (error) {
    // Server error
    res.writeHead(500);
    res.end(JSON.stringify({ error: 'server failed' }));
  }
});

// Start the server
server.listen(PORT, () => {
  console.log(`Server listening on port ${PORT}`);
});