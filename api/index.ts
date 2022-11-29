import express, { Express, Request, Response } from "express";
import { BaseAPIResponse } from "@koyo/types";

const app: Express = express();
const port = process.env.PORT ?? 8080;

app.get("/", (req: Request, res: Response<BaseAPIResponse>) => {
  res.send({ error: false, message: "Hello World" });
});

app.listen(port, () => {
  console.log(`⚡️[server]: Server is running at https://localhost:${port}`);
});
