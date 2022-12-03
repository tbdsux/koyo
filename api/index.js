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
	const { website } = req.body;
	if (!website) {
		res.status(400).json({ error: true, message: 'Website url not defined in body.', code: 400 });
		return;
	}

	try {
		const browser = await puppeteer.launch({
			args: chromium.args,
			defaultViewport: chromium.defaultViewport,
			executablePath: await chromium.executablePath,
			headless: chromium.headless,
			ignoreHTTPSErrors: true
		});
		const page = await browser.newPage();
		await page.goto(website, { waitUntil: 'domcontentloaded' });

		const data = await page.screenshot({ encoding: 'base64' });
		const base64Data = data.replace(/^data:image\/png;base64,/, '');
		const img = Buffer.from(base64Data, 'base64');

		await browser.close();

		res.writeHead(200, {
			'Content-Type': 'image/png',
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
