<div align="center">
    <img src="./assets/banner.png">
    <h1>koyo</h1>
    <p>Web screenshot service api</p>
</div>

## New Features

- Whitehole integration from https://alpha.deta.space/discovery/@mikhailsdv/black_hole-3kf

## Known Issues

- Full page screenshot doesn't work on some websites.

## Development

This project is a monorepo using `pnpm`.

- Clone the repository

  ```sh
  git clone https://github.com/TheBoringDude/koyo.git
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

**&copy; 2022 | TheBoringDude**
