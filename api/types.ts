export interface ScreenshotApiOptions {
	height?: string;
	width?: string;
	imageType?: 'png' | 'jpeg' | 'webp';
	fullPage?: string;
	driver?: 'playwright' | 'puppeteer';
}

export interface ScreenshotBodyOptions {
	website: string;
}
