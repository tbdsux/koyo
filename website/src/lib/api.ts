import { dev } from '$app/environment';

export const apiUrl = dev ? 'http://localhost:8080' : '/api';

export interface APIResponse {
	error: boolean;
	code: number;
	message?: string;
}

export interface APILoadFilesDataProps extends APIResponse {
	data: string[];
}
