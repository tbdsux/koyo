const express = require('express');
const chromium = require('@sparticuz/chromium');
const puppeteer = require('puppeteer');
const cors = require('cors');

const app = express();
const port = 8080;

app.use(cors());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

app.get('/', (req, res) => {
	res.send('Hello World!');
});

app.post('/screenshot', async (req, res) => {
	/** @type {import('./types').ScreenshotBodyOptions} */
	const { website } = req.body;
	if (!website) {
		res.status(400).json({ error: true, message: 'Website url not defined in body.', code: 400 });
		return;
	}

	/** @type {import('./types').ScreenshotApiOptions} */
	const { height, width, imageType, fullPage } = req.query;

	const _height = !isNaN(Number(height)) ? Number(height) : 1280;
	const _width = !isNaN(Number(width)) ? Number(width) : 800;
	const _fullPage = fullPage === 'true';

	/** @type {import('puppeteer').Viewport} */
	const viewport = {
		width: _width,
		height: _height
	};

	try {
		const browser = await puppeteer.launch({
			args: chromium.args,
			executablePath: await chromium.executablePath,
			headless: chromium.headless,
			ignoreHTTPSErrors: true
		});
		const page = await browser.newPage();
		await page.goto(website, { waitUntil: 'networkidle2' });

		if (_fullPage) {
			viewport.height = 16834;
		}
		await page.setViewport(viewport);

		const data = await page.screenshot({ encoding: 'base64', _fullPage, type: imageType });
		const img = Buffer.from(data, 'base64');

		await browser.close();

		res.writeHead(200, {
			'Content-Type': `image/${imageType}`,
			'Content-Length': img.length
		});
		res.end(img);
	} catch (e) {
		res.status(500).json({ error: true, message: String(e), code: 500 });
	}
});

app.listen(port, () => {
	console.log(`Example app listening on http://localhost:${port}`);
});
