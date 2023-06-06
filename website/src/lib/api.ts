export interface APIResponse {
	error: boolean;
	code: number;
	message?: string;
}

export interface APILoadFilesDataProps extends APIResponse {
	data: string[];
}
