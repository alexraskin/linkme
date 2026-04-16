# linkme

A personal link tree built with [Astro](https://astro.build) and deployed on [Cloudflare Workers](https://workers.cloudflare.com).

## Setup

```sh
npm install
```

## Configuration

Edit `src/config.ts` to set your profile and links:

```ts
export const profile = {
  name: "Your Name",
  handle: "@yourhandle",
  bio: "Your bio",
  avatar: null, // set to "/avatar.jpg" for an image in public/
};

export const links = [
  { text: "GitHub", url: "https://github.com/you" },
  { text: "Twitter", url: "https://x.com/you" },
];
```

To use a profile photo, place an image in the `public/` directory and set `avatar` to its path.

## Development

```sh
npm run dev
```

## Deploy

```sh
npm run deploy
```

## License

unlicense
