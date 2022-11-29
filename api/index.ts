import express, { Express, Request, Response } from 'express';
import { BaseAPIResponse } from '@koyo/types';

const app: Express = express();
const port = 8080;

app.get('/', (req: Request, res: Response<BaseAPIResponse>) => {
	res.json({ error: false, message: 'Hello World' });
});

app.listen(port, () => {
	console.log(`⚡️[server]: Server is running at http://localhost:${port}`);
});
