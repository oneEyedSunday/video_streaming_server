const fs = require('fs');
const path = require('path');
const http = require('http');

const { promisify } = require('util');


const PORT = parseInt(process.env.PORT, 10) || 4000

/**
 * @type {import('http').RequestListener}
 */
const listener = (req, res) => {
    res.setHeader('Content-Type', 'text/html');
    const servePt = path.join(__dirname, 'public/index.html');
    fs.readFile(servePt, (err, data) => {
        if (err) {
            res.writeHead(404, 'Not Found');
            res.write('404: File Not Found!');
            return res.end();
        }

        res.statusCode = 200;

        res.write(data);
        return res.end();
    });
};

/**
 * @type {import('http').ServerOptions}
 */
const httpOpts = {};

const app = http.createServer(httpOpts, listener);

app.listen(PORT, () => {
    console.info(`Server up at http://localhost:${PORT}`);
});
