export interface ScreenshotApiOptions {
	height: string;
	width: string;
	imageType: 'png' | 'jpeg' | 'webp';
	fullPage: string;
}

export interface ScreenshotBodyOptions {
	website: string;
}
