export interface ScreenshotApiOptions {
	height?: string;
	width?: string;
	imageType?: 'png' | 'jpeg' | 'webp';
	fullPage?: string;
	driver?: 'playwright' | 'puppeteer';
	whiteholeUrl?: string;
}

export interface ScreenshotBodyOptions {
	website: string;
}
