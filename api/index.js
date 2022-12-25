const express = require('express');
const chromium = require('@sparticuz/chromium');
const { chromium: playwright } = require('playwright-core');
const puppeteer = require('puppeteer');
const cors = require('cors');
const FormData = require('form-data');
const fetch = require('cross-fetch');

const app = express();
const port = 8080;

const imageTypes = ['png', 'jpeg', 'webp'];

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
	const { height, width, imageType, fullPage, driver, whiteholeUrl } = req.query;

	const _height = !isNaN(Number(height)) ? Number(height) : 1280;
	const _width = !isNaN(Number(width)) ? Number(width) : 800;
	const _fullPage = fullPage === 'true';
	const _usePuppeteer = driver === 'puppeteer';
	const _whiteholeUrl = whiteholeUrl ?? '';
	let _imageType = imageTypes.includes(imageType) ? imageType : 'png';

	/** @type {import('puppeteer').Viewport} */
	const viewport = {
		width: _width,
		height: _height
	};

	/** @type {import('puppeteer').Browser | import('playwright-core').Browser | null} */
	let browser = null;

	try {
		/** @type {Buffer} */
		let img;

		if (_usePuppeteer) {
			browser = await puppeteer.launch({
				args: chromium.args,
				executablePath: await chromium.executablePath,
				headless: chromium.headless,
				ignoreHTTPSErrors: true
			});
			const page = await browser.newPage();
			await page.setViewport(viewport);

			await page.goto(website, { waitUntil: 'load' });

			const data = await page.screenshot({
				encoding: 'base64',
				fullPage: _fullPage,
				type: _imageType
			});
			img = Buffer.from(data, 'base64');
		} else {
			_imageType = 'png'; // playwright doesn't support webp

			browser = await playwright.launch({
				args: chromium.args,
				executablePath: await chromium.executablePath,
				headless: chromium.headless
			});
			const context = await browser.newContext();
			const page = await context.newPage();
			await page.setViewportSize(viewport);

			await page.goto(website, { waitUntil: 'load' });
			img = await page.screenshot({ fullPage: _fullPage, type: _imageType });
		}

		// close browser
		await browser.close();
		browser = null;

		// send to whitehole if exists
		if (_whiteholeUrl != '') {
			const form = new FormData();
			form.append('photo', img, `${website}.${imageType}`);

			try {
				const r = await fetch(_whiteholeUrl, {
					method: 'POST',
					body: form
				});

				const data = await r.json();
				const { error, status } = data;
				if (!status) {
					res.status(500).json({ error: false, message: error, code: 500 });
					return;
				}
			} catch (e) {
				res.status(500).json({ error: true, message: String(e), code: 500 });
				return;
			}
		}

		res.writeHead(200, {
			'Content-Type': `image/${_imageType}`,
			'Content-Length': img.length
		});
		res.end(img);
	} catch (e) {
		try {
			if (browser != null) {
				// make sure to close the browser
				await browser.close();
			}
		} catch (e) {
			console.error(e);
		}

		res.status(500).json({ error: true, message: String(e), code: 500 });
	}
});

app.listen(port, () => {
	console.log(`App listening on http://localhost:${port}`);
});
