import { dev } from '$app/environment';
import { env } from '$env/dynamic/private';
import type { APILoadFilesDataProps } from '$lib/api';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
	const baseUrl = dev ? 'http://localhost:4200/api' : `https://${env.DETA_SPACE_APP_HOSTNAME}/api`;

	console.log(baseUrl);

	const res = await fetch(baseUrl + '/drive/files', {
		headers: {
			'X-API-KEY': env.DETA_API_KEY || ''
		}
	});
	const driveData: APILoadFilesDataProps = await res.json();

	return {
		baseUrl,
		driveData
	};
};
