# linkme

A simple, fast link tree style page built with Go. All assets are embedded at build time into a single binary.

## Usage

Edit `config.yaml` to customize your links, then build and run:

```bash
make build
make run
```

## Docker

```bash
make docker-build
make docker-run
```

## Configuration

All configuration is done in `config.yaml`. See the file for available options:

- Profile: name, subtitle, description, avatar
- Links: main link buttons with icons and colors
- Sections: grouped links with headings
- Socials: footer social icons
- Meta: SEO and favicon settings

## Assets

Place your images in the `assets/` directory:

- `avatar.png` - profile picture
- `favicon.ico` - favicon
- `icons/` - custom link icons

## License

unlicense
