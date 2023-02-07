<div align="center">
    <img src="./assets/banner.png">
    <h1>koyo</h1>
    <p>Web screenshot service api</p>
</div>

## [CLI](./cli/koyo)

## API

You will need to create an `API Key` from your Koyo space app before communicating through the screenshot api outside of space app.

### Structure

- `/api/screenshot`

  Main screenshot api. Options can be viewed and modifed through your space app instance.

- `/api/drive/files`

  List saved screenshot files from the default app drive. Filename is sorted with `{iso-date-string}-{website-hostname}.{image-type}` (Example: `2023-02-07T09:49:11.890Z-www.google.com.png`)

- `/api/drive/files/{filename}`

  Endpoint for saved screenshot `filename`. It will return the image buffer / data if it exists otherwise, will return an `Error 404` json object.

## Changelog

### v0.3.0

- Added [`API Keys`](https://deta.space/docs/en/basics/micros#api-keys) for security purposes.
- Added **Save To Drive** screenshot api option. Can also view saved screenshots.
- Add more features for [CLI](./cli/koyo/README.md)

### v0.2.0

- Bug fixes and some minor adjustments and additions
- Download image button

### v0.1.1

- Whitehole integration from https://alpha.deta.space/discovery/@mikhailsdv/black_hole-3kf

## Known Issues

- Full page screenshot doesn't work on some websites.

## Development

This project is a monorepo using `pnpm`.

- Clone the repository

  ```sh
  git clone https://github.com/tbdsux/koyo.git
  ```

- Install the dependencies

  ```sh
  pnpm install
  ```

- Start development. This will run the both the api and frontend projects together simultaneously.

  ```sh
  pnpm dev
  ```

### Project structure

- `api` - Expressjs, main screenshot api.

- `website` - Sveltekit website frontend.

##

**&copy; 2022 | tbdsux**
