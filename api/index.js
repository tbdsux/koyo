const express = require('express');
const chromium = require('@sparticuz/chromium');
const { chromium: playwright } = require('playwright-core');
const puppeteer = require('puppeteer-core');
const cors = require('cors');
const FormData = require('form-data');
const fetch = require('cross-fetch');
const { Drive } = require('deta');

const app = express();
const port = process.env.PORT ?? 8080;

const imageTypes = ['png', 'jpeg', 'webp'];
const drive = Drive('__default__');

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
	const { height, width, imageType, fullPage, driver, whiteholeUrl, saveToDrive, saveNoOutput } =
		req.query;

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
				executablePath: await chromium.executablePath(),
				headless: 'new',
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
				executablePath: await chromium.executablePath(),
				headless: true
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
					res.status(500).json({ error: true, message: error, code: 500 });
					return;
				}
			} catch (e) {
				res.status(500).json({ error: true, message: String(e), code: 500 });
				return;
			}
		}

		// save to drive if true
		if (saveToDrive === 'true') {
			drive.put(`${new Date().toISOString()}-${new URL(website).hostname}.${imageType}`, {
				data: img,
				contentType: `image/${imageType}`
			});

			// if noOutput is set, do not return the screenshot image
			if (saveNoOutput === 'true') {
				res.status(200).json({ error: false, message: 'Successfully saved screenshot to Drive.' });
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

app.get('/drive/files', async (req, res) => {
	let result = await drive.list();
	let allFiles = result.names;
	let last = result.paging?.last;

	while (last) {
		// provide last from previous call
		result = await drive.list({ last: result.paging.last });

		allFiles = allFiles.concat(result.names);

		// update last
		last = result.paging.last;
	}

	res.status(200).json({ error: false, data: allFiles, code: 200 });
});

app.get(`/drive/files/:filename`, async (req, res) => {
	const { filename } = req.params;

	const blob = await drive.get(filename);

	if (blob == null) {
		res.status(404).json({ error: true, message: 'Screenshot not found.', code: 404 });
		return;
	}

	const buffer = await blob.arrayBuffer();

	res.type(blob.type);
	res.send(Buffer.from(buffer));
});

app.listen(port, () => {
	console.log(`App listening on http://localhost:${port}`);
});
