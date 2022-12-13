import { dev } from '$app/environment';
import { env } from '$env/dynamic/private';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	const baseUrl = dev ? 'http://localhost:8080' : `https://${env.DETA_SPACE_APP_HOSTNAME}/api`;

	return {
		baseUrl
	};
};
