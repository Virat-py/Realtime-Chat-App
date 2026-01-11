// src/routes/dashboard/+page.js
import { redirect } from '@sveltejs/kit';

export const load = () => {
	throw redirect(302, '/login');
};
