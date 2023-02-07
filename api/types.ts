export interface ScreenshotApiOptions {
	height?: string;
	width?: string;
	imageType?: 'png' | 'jpeg' | 'webp';
	fullPage?: string;
	driver?: 'playwright' | 'puppeteer';
	whiteholeUrl?: string;
	saveToDrive?: string;
	saveNoOutput?: string;
}

export interface ScreenshotBodyOptions {
	website: string;
}
