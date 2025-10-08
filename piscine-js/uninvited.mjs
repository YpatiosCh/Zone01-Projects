import http from "http";
import fs from "fs/promises";
import path from "path";

const PORT = 5000;
const GUESTS_DIR = './guests';

const server = http.createServer(async (req, res) => {
    // handle only POST requests
    if (req.method !== 'POST') {
        res.writeHead(405, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Method Not Allowed' }));
        return;
    }

    // get the name from url
    const guestName = req.url.slice(1);
    if (!guestName) {
        res.writeHead(400, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'guest name is required' }));
        return;
    }

    // read the request body
    let body = '';
    req.on('data', chunk => {
        body += chunk.toString();
    });

    // handle the end of the request (on.end event is fired when the request body is fully received)
    req.on('end', async () => {
        try {
            // construct the file path
            const filePath = path.join(GUESTS_DIR, `${guestName}.json`);

            // write the guest data to a file
            await fs.writeFile(filePath, body, 'utf8');

            res.writeHead(201, { 'Content-Type': 'application/json' });
            res.end(body);

        } catch (error) {
            res.writeHead(500, { 'Content-Type': 'application/json' });
            res.end(JSON.stringify({ error: 'server failed' }));
        }
    });
});

// start the server
server.listen(PORT, () => {
    console.log(`Server is running on http://localhost:${PORT}`);
});
