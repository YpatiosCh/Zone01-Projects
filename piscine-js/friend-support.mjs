import http from "http";
import fs from "fs/promises";
import path from "path";

const PORT = 5000;
const GUESTS_DIR = './guests';

const server = http.createServer(async (req, res) => { 
    // handle only GET requests
    if (req.method !== 'GET') {
        res.writeHead(405, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'Method Not Allowed'}));
        return;
    }

    // get the guest name 
    const guestName = req.url.slice(1);
    if (!guestName) {
        res.writeHead(404, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify({ error: 'guest not found' }));
        return;
    }

    // construct the file path
    const filePath = path.join(GUESTS_DIR, `${guestName}.json`);

    try {
        const data = await fs.readFile(filePath, 'utf8');
        const guestData = JSON.parse(data);
        res.writeHead(200, { 'Content-Type': 'application/json' });
        res.end(JSON.stringify(guestData));

    } catch (error) {

        if (error.code === 'ENOENT') {
            res.writeHead(404, { 'Content-Type': 'application/json' });
            res.end(JSON.stringify({ error: 'guest not found' }));
            return;
        } else {
            res.writeHead(500, { 'Content-Type': 'application/json' });
            res.end(JSON.stringify({ error: 'server failed' }));
            return;
        }        
    }
});

// start the server
server.listen(PORT, () => {
    console.log(`Server is running on http://localhost:${PORT}`);
});